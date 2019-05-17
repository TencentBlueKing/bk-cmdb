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

package config

import (
	"context"
	"sync"

	"configcenter/src/common/backbone/service_mange/zk"
	"configcenter/src/common/confregdiscover"
	"configcenter/src/common/types"
	"configcenter/src/framework/core/log"
)

// ConfCenter discover configure changed. get, update configures
type ConfCenter struct {
	confRegDiscv *confregdiscover.ConfRegDiscover
	rootCtx      context.Context
	cancel       context.CancelFunc
	ctx          []byte
	ctxLock      sync.RWMutex
}

// NewConfCenter create a ConfCenter object
func NewConfCenter(client *zk.ZkClient) *ConfCenter {
	ctx, cancel := client.WithCancel()

	return &ConfCenter{
		ctx:          nil,
		rootCtx:      ctx,
		cancel:       cancel,
		confRegDiscv: confregdiscover.NewConfRegDiscover(client),
	}
}

// Ping to ping server
func (cc *ConfCenter) Ping() error {
	return cc.confRegDiscv.Ping()
}

// Start the configure center module service
func (cc *ConfCenter) Start() error {

	// here: discover itselft configure
	confPath := types.CC_SERVCONF_BASEPATH + "/" + types.CC_MODULE_OBJECTCONTROLLER
	confEvent, err := cc.confRegDiscv.DiscoverConfig(confPath)
	if err != nil {
		log.Errorf("fail to discover configure for objectcontroller service. err:%s", err.Error())
		return err
	}

	for {
		select {
		case confEvn := <-confEvent:
			cc.dealConfChangeEvent(confEvn.Data)
		case <-cc.rootCtx.Done():
			log.Warning("configure discover service done")
			return nil
		}
	}
}

// GetConfigureCtx fetch the configure
func (cc *ConfCenter) GetConfigureCxt() []byte {
	cc.ctxLock.RLock()
	defer cc.ctxLock.RUnlock()

	return cc.ctx
}

func (cc *ConfCenter) dealConfChangeEvent(data []byte) error {
	log.Infof("%s configure has changed", types.CC_MODULE_OBJECTCONTROLLER)

	cc.ctxLock.Lock()
	defer cc.ctxLock.Unlock()

	cc.ctx = data

	return nil
}
