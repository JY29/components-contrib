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

package postmark

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/JY29/components-contrib/bindings"
	"github.com/dapr/kit/logger"
)

func TestParseMetadata(t *testing.T) {
	logger := logger.NewLogger("test")

	t.Run("Has correct metadata", func(t *testing.T) {
		m := bindings.Metadata{}
		m.Properties = map[string]string{"serverToken": "abc", "accountToken": "123", "emailFrom": "test1@example.net", "emailTo": "test2@example.net", "subject": "hello"}
		r := Postmark{logger: logger}
		pMeta, err := r.parseMetadata(m)
		assert.Nil(t, err)
		assert.Equal(t, "abc", pMeta.ServerToken)
		assert.Equal(t, "123", pMeta.AccountToken)
		assert.Equal(t, "test1@example.net", pMeta.EmailFrom)
		assert.Equal(t, "test2@example.net", pMeta.EmailTo)
		assert.Equal(t, "hello", pMeta.Subject)
	})
}
