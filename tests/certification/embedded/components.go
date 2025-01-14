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

package embedded

import (
	"github.com/dapr/dapr/pkg/runtime"
	"github.com/dapr/kit/logger"

	// Name resolutions.
	nrConsul "github.com/JY29/components-contrib/nameresolution/consul"
	nrKubernetes "github.com/JY29/components-contrib/nameresolution/kubernetes"
	nrMdns "github.com/JY29/components-contrib/nameresolution/mdns"

	nrLoader "github.com/dapr/dapr/pkg/components/nameresolution"
)

func CommonComponents(log logger.Logger) []runtime.Option {
	registry := nrLoader.NewRegistry()
	registry.Logger = log
	registry.RegisterComponent(nrMdns.NewResolver, "mdns")
	registry.RegisterComponent(nrKubernetes.NewResolver, "kubernetes")
	registry.RegisterComponent(nrConsul.NewResolver, "consul")
	return []runtime.Option{
		runtime.WithNameResolutions(registry),
	}
}
