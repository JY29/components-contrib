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

package blobstorage

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/JY29/components-contrib/state"
	"github.com/dapr/kit/logger"
)

func TestInit(t *testing.T) {
	m := state.Metadata{}
	s := NewAzureBlobStorageStore(logger.NewLogger("logger")).(*StateStore)
	t.Run("Init with valid metadata", func(t *testing.T) {
		m.Properties = map[string]string{
			"accountName":   "acc",
			"accountKey":    "e+Dnvl8EOxYxV94nurVaRQ==",
			"containerName": "dapr",
		}
		err := s.Init(m)
		assert.Nil(t, err)
		assert.Equal(t, "https://acc.blob.core.windows.net/dapr", s.containerClient.URL())
	})

	t.Run("Init with missing metadata", func(t *testing.T) {
		m.Properties = map[string]string{
			"invalidValue": "a",
		}
		err := s.Init(m)
		assert.NotNil(t, err)
		assert.Equal(t, err, fmt.Errorf("missing or empty accountName field from metadata"))
	})

	t.Run("Init with invalid account name", func(t *testing.T) {
		m.Properties = map[string]string{
			"accountName":   "invalid-account",
			"accountKey":    "e+Dnvl8EOxYxV94nurVaRQ==",
			"containerName": "dapr",
		}
		s.Init(m)
		err := s.Ping()
		assert.NotNil(t, err)
	})
}

func TestFileName(t *testing.T) {
	t.Run("Valid composite key", func(t *testing.T) {
		key := getFileName("app_id||key")
		assert.Equal(t, "key", key)
	})

	t.Run("No delimiter present", func(t *testing.T) {
		key := getFileName("key")
		assert.Equal(t, "key", key)
	})
}
