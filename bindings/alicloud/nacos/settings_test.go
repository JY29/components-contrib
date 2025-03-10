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

package nacos_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JY29/components-contrib/bindings/alicloud/nacos"
)

func TestParseMetadata(t *testing.T) { //nolint:paralleltest
	props := map[string]string{
		"endpoint":        "a",
		"region":          "b",
		"namespace":       "c",
		"accessKey":       "d",
		"secretKey":       "e",
		"updateThreadNum": "3",
	}

	var settings nacos.Settings
	err := settings.Decode(props)
	require.NoError(t, err)
	assert.Equal(t, "a", settings.Endpoint)
	assert.Equal(t, "b", settings.RegionID)
	assert.Equal(t, "c", settings.NamespaceID)
	assert.Equal(t, "d", settings.AccessKey)
	assert.Equal(t, "e", settings.SecretKey)
	assert.Equal(t, 3, settings.UpdateThreadNum)
}
