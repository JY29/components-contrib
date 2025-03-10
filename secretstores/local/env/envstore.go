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

package env

import (
	"context"
	"os"
	"reflect"
	"strings"

	"github.com/JY29/components-contrib/metadata"
	"github.com/JY29/components-contrib/secretstores"
	"github.com/dapr/kit/logger"
)

var _ secretstores.SecretStore = (*envSecretStore)(nil)

type envSecretStore struct {
	logger logger.Logger
}

// NewEnvSecretStore returns a new env var secret store.
func NewEnvSecretStore(logger logger.Logger) secretstores.SecretStore {
	return &envSecretStore{
		logger: logger,
	}
}

// Init creates a Local secret store.
func (s *envSecretStore) Init(metadata secretstores.Metadata) error {
	return nil
}

// GetSecret retrieves a secret from env var using provided key.
func (s *envSecretStore) GetSecret(ctx context.Context, req secretstores.GetSecretRequest) (secretstores.GetSecretResponse, error) {
	return secretstores.GetSecretResponse{
		Data: map[string]string{
			req.Name: os.Getenv(req.Name),
		},
	}, nil
}

// BulkGetSecret retrieves all secrets in the store and returns a map of decrypted string/string values.
func (s *envSecretStore) BulkGetSecret(ctx context.Context, req secretstores.BulkGetSecretRequest) (secretstores.BulkGetSecretResponse, error) {
	r := map[string]map[string]string{}

	for _, element := range os.Environ() {
		envVariable := strings.SplitN(element, "=", 2)
		r[envVariable[0]] = map[string]string{envVariable[0]: envVariable[1]}
	}

	return secretstores.BulkGetSecretResponse{
		Data: r,
	}, nil
}

// Features returns the features available in this secret store.
func (s *envSecretStore) Features() []secretstores.Feature {
	return []secretstores.Feature{} // No Feature supported.
}

func (s *envSecretStore) GetComponentMetadata() map[string]string {
	type unusedMetadataStruct struct{}
	metadataStruct := unusedMetadataStruct{}
	metadataInfo := map[string]string{}
	metadata.GetMetadataInfoFromStructType(reflect.TypeOf(metadataStruct), &metadataInfo)
	return metadataInfo
}
