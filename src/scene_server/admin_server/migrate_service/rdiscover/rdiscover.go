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

package rdiscover

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"configcenter/src/common/RegisterDiscover"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/version"

	"context"
)

// RegDiscover register and discover
type RegDiscover struct {
	ip        string
	port      uint
	isSSL     bool
	rd        *RegisterDiscover.RegDiscover
	rootCtx   context.Context
	cancel    context.CancelFunc
	topoServs []*types.TopoServInfo
	topoLock  sync.RWMutex
	procServs []*types.ProcServInfo
	procLock  sync.RWMutex
}

// NewRegDiscover create a RegDiscover object
func NewRegDiscover(zkserv string, ip string, port uint, isSSL bool) *RegDiscover {
	return &RegDiscover{
		ip:        ip,
		port:      port,
		isSSL:     isSSL,
		rd:        RegisterDiscover.NewRegDiscoverEx(zkserv, 10*time.Second),
		topoServs: []*types.TopoServInfo{},
		procServs: []*types.ProcServInfo{},
	}
}

// Ping to ping server
func (cc *RegDiscover) Ping() error {
	return cc.rd.Ping()
}

// Start the register and discover
func (r *RegDiscover) Start() error {

	//create root context
	r.rootCtx, r.cancel = context.WithCancel(context.Background())
	//start regdiscover
	if err := r.rd.Start(); err != nil {
		blog.Error("fail to start register and discover serv. err:%s", err.Error())
		return err
	}

	// register migrate server
	if err := r.registerMigrate(); err != nil {
		blog.Error("fail to register migrate(%s), err:%s", r.ip, err.Error())
		return err
	}

	for {
		select {
		case <-r.rootCtx.Done():
			blog.Warn("register and discover serv done")
			return nil
		}
	}
}

// Stop the register and discover
func (r *RegDiscover) Stop() error {
	r.cancel()

	r.rd.Stop()

	return nil
}

func (r *RegDiscover) GetServer(servType string) (string, error) {
	switch servType {
	// case types.CC_MODULE_PROC:
	// 	return r.GetProcServ()
	}

	err := fmt.Errorf("there is no server discover for type(%s)", servType)
	blog.Errorf("%s", err.Error())
	return "", err
}

func (r *RegDiscover) registerMigrate() error {
	migrateServInfo := new(types.MigrateServInfo)

	migrateServInfo.IP = r.ip
	migrateServInfo.Port = r.port
	migrateServInfo.Scheme = "http"
	if r.isSSL {
		migrateServInfo.Scheme = "https"
	}

	migrateServInfo.Version = version.GetVersion()
	migrateServInfo.Pid = os.Getpid()

	data, err := json.Marshal(migrateServInfo)
	if err != nil {
		blog.Error("fail to marshal Migrate server info to json. err:%s", err.Error())
		return err
	}

	path := types.CC_SERV_BASEPATH + "/" + types.CC_MODULE_MIGRATE + "/" + r.ip

	return r.rd.RegisterAndWatchService(path, data)
}
