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

package kubernetes

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/JY29/components-contrib/bindings"
	kubeclient "github.com/JY29/components-contrib/internal/authentication/kubernetes"
	"github.com/dapr/kit/logger"
)

type kubernetesInput struct {
	kubeClient   kubernetes.Interface
	namespace    string
	resyncPeriod time.Duration
	logger       logger.Logger
}

type EventResponse struct {
	Event  string   `json:"event"`
	OldVal v1.Event `json:"oldVal"`
	NewVal v1.Event `json:"newVal"`
}

// NewKubernetes returns a new Kubernetes event input binding.
func NewKubernetes(logger logger.Logger) bindings.InputBinding {
	return &kubernetesInput{logger: logger}
}

func (k *kubernetesInput) Init(metadata bindings.Metadata) error {
	client, err := kubeclient.GetKubeClient()
	if err != nil {
		return err
	}
	k.kubeClient = client

	return k.parseMetadata(metadata)
}

func (k *kubernetesInput) parseMetadata(metadata bindings.Metadata) error {
	if val, ok := metadata.Properties["namespace"]; ok && val != "" {
		k.namespace = val
	} else {
		return errors.New("namespace is missing in metadata")
	}
	if val, ok := metadata.Properties["resyncPeriodInSec"]; ok && val != "" {
		intval, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			k.logger.Warnf("invalid resyncPeriodInSec %s; %v; defaulting to 10s", val, err)
			k.resyncPeriod = time.Second * 10
		} else {
			k.resyncPeriod = time.Second * time.Duration(intval)
		}
	}

	return nil
}

func (k *kubernetesInput) Read(ctx context.Context, handler bindings.Handler) error {
	watchlist := cache.NewListWatchFromClient(
		k.kubeClient.CoreV1().RESTClient(),
		"events",
		k.namespace,
		fields.Everything(),
	)
	resultChan := make(chan EventResponse)
	_, controller := cache.NewInformer(
		watchlist,
		&v1.Event{},
		k.resyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if obj != nil {
					resultChan <- EventResponse{
						Event:  "add",
						NewVal: *(obj.(*v1.Event)),
						OldVal: v1.Event{},
					}
				} else {
					k.logger.Warnf("Nil Object in Add handle %v", obj)
				}
			},
			DeleteFunc: func(obj interface{}) {
				if obj != nil {
					resultChan <- EventResponse{
						Event:  "delete",
						OldVal: *(obj.(*v1.Event)),
						NewVal: v1.Event{},
					}
				} else {
					k.logger.Warnf("Nil Object in Delete handle %v", obj)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				if oldObj != nil && newObj != nil {
					resultChan <- EventResponse{
						Event:  "update",
						OldVal: *(oldObj.(*v1.Event)),
						NewVal: *(newObj.(*v1.Event)),
					}
				} else {
					k.logger.Warnf("Nil Objects in Update handle %v %v", oldObj, newObj)
				}
			},
		},
	)

	// Start the controller in backgound
	stopCh := make(chan struct{})
	go controller.Run(stopCh)

	// Watch for new messages and for context cancellation
	go func() {
		var (
			obj  EventResponse
			data []byte
			err  error
		)
		for {
			select {
			case obj = <-resultChan:
				data, err = json.Marshal(obj)
				if err != nil {
					k.logger.Errorf("Error marshalling event %w", err)
				} else {
					handler(ctx, &bindings.ReadResponse{
						Data: data,
					})
				}
			case <-ctx.Done():
				close(stopCh)
				return
			}
		}
	}()

	return nil
}
