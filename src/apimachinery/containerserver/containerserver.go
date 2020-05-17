/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under,
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package containerserver

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/metadata"
)

// ContainerServerClientInterface client interface for container server
type ContainerServerClientInterface interface {
	CreatePod(ctx context.Context, h http.Header, bizID int64, data interface{}) (*metadata.CreatedOneOptionResult, error)
	CreateManyPod(ctx context.Context, h http.Header, bizID int64, data interface{}) (*metadata.CreatedManyOptionResult, error)
	UpdatePod(ctx context.Context, h http.Header, bizID int64, data interface{}) (*metadata.UpdatedOptionResult, error)
	DeletePod(ctx context.Context, h http.Header, bizID int64, data interface{}) (*metadata.DeletedOptionResult, error)
	ListPods(ctx context.Context, h http.Header, bizID int64, data interface{}) (*metadata.ListPodsResult, error)
}

// NewContainerServerClientInterface create container server client interface
func NewContainerServerClientInterface(c *util.Capability, version string) ContainerServerClientInterface {
	base := fmt.Sprintf("/container/%s", version)
	return &containerServer{
		client: rest.NewRESTClient(c, base),
	}
}

type containerServer struct {
	client rest.ClientInterface
}
