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

package hazelcast

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mdata "github.com/JY29/components-contrib/metadata"
	"github.com/JY29/components-contrib/pubsub"
)

func TestValidateMetadata(t *testing.T) {
	t.Run("return error when required servers is empty", func(t *testing.T) {
		fakeMetaData := pubsub.Metadata{Base: mdata.Base{
			Properties: map[string]string{
				hazelcastServers: "",
			},
		}}

		m, err := parseHazelcastMetadata(fakeMetaData)

		// assert
		assert.Error(t, err)
		assert.Empty(t, m)
	})
}
