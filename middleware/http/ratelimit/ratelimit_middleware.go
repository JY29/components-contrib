/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ratelimit

import (
	"fmt"
	"net/http"
	"strconv"

	tollbooth "github.com/didip/tollbooth/v7"
	libstring "github.com/didip/tollbooth/v7/libstring"

	"github.com/JY29/components-contrib/middleware"
	"github.com/dapr/kit/logger"
)

// Metadata is the ratelimit middleware config.
type rateLimitMiddlewareMetadata struct {
	MaxRequestsPerSecond float64 `json:"maxRequestsPerSecond"`
}

const (
	maxRequestsPerSecondKey = "maxRequestsPerSecond"

	// Defaults.
	defaultMaxRequestsPerSecond = 100
)

// NewRateLimitMiddleware returns a new ratelimit middleware.
func NewRateLimitMiddleware(_ logger.Logger) middleware.Middleware {
	return &Middleware{}
}

// Middleware is an ratelimit middleware.
type Middleware struct{}

// GetHandler returns the HTTP handler provided by the middleware.
func (m *Middleware) GetHandler(metadata middleware.Metadata) (func(next http.Handler) http.Handler, error) {
	meta, err := m.getNativeMetadata(metadata)
	if err != nil {
		return nil, err
	}

	limiter := tollbooth.NewLimiter(meta.MaxRequestsPerSecond, nil)

	return func(next http.Handler) http.Handler {
		// Adapted from toolbooth.LimitHandler
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// The tollbooth library requires a remote IP. If this isn't present in the request's headers, then we need to set a value for X-Forwarded-For or the rate limiter won't work
			remoteIP := libstring.RemoteIP(limiter.GetIPLookups(), limiter.GetForwardedForIndexFromBehind(), r)
			remoteIP = libstring.CanonicalizeIP(remoteIP)
			if remoteIP == "" {
				// Forcefully set a remote IP
				r.Header.Set("X-Forwarded-For", "0.0.0.0")
			}

			httpError := tollbooth.LimitByRequest(limiter, w, r)
			if httpError != nil {
				limiter.ExecOnLimitReached(w, r)
				if limiter.GetOverrideDefaultResponseWriter() {
					return
				}
				w.Header().Add("Content-Type", limiter.GetMessageContentType())
				w.WriteHeader(httpError.StatusCode)
				w.Write([]byte(httpError.Message))
				return
			}

			// There's no rate-limit error, serve the next handler.
			next.ServeHTTP(w, r)
		})
	}, nil
}

func (m *Middleware) getNativeMetadata(metadata middleware.Metadata) (*rateLimitMiddlewareMetadata, error) {
	var middlewareMetadata rateLimitMiddlewareMetadata

	middlewareMetadata.MaxRequestsPerSecond = defaultMaxRequestsPerSecond
	if val, ok := metadata.Properties[maxRequestsPerSecondKey]; ok {
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing ratelimit middleware property %s: %w", maxRequestsPerSecondKey, err)
		}
		if f <= 0 {
			return nil, fmt.Errorf("ratelimit middleware property %s must be a positive value", maxRequestsPerSecondKey)
		}
		middlewareMetadata.MaxRequestsPerSecond = f
	}

	return &middlewareMetadata, nil
}
