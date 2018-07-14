/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package hostcontroller

import (
	"fmt"

	"configcenter/src/apimachinery/hostcontroller/favorite"
	"configcenter/src/apimachinery/hostcontroller/history"
	"configcenter/src/apimachinery/hostcontroller/host"
	"configcenter/src/apimachinery/hostcontroller/module"
	"configcenter/src/apimachinery/hostcontroller/user"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
)

type HostCtrlClientInterface interface {
	Favorite() favorite.FavoriteInterface
	History() history.HistoryInterface
	Host() host.HostInterface
	Module() module.ModuleInterface
	User() user.UserInterface
}

func NewHostCtrlClientInterface(c *util.Capability, version string) HostCtrlClientInterface {
	base := fmt.Sprintf("/host/%s", version)
	return &hostctl{
		client: rest.NewRESTClient(c, base),
	}
}

type hostctl struct {
	client rest.ClientInterface
}

func (h *hostctl) Favorite() favorite.FavoriteInterface {
	return favorite.NewFavoriteInterface(h.client)
}

func (h *hostctl) History() history.HistoryInterface {
	return history.NewHistoryInterface(h.client)
}

func (h *hostctl) Host() host.HostInterface {
	return host.NewHostInterface(h.client)
}

func (h *hostctl) Module() module.ModuleInterface {
	return module.NewModuleInterface(h.client)
}

func (h *hostctl) User() user.UserInterface {
	return user.NewUserInterface(h.client)
}
