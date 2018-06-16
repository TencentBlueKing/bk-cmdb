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
	"os"
	"time"

	"configcenter/src/common/RegisterDiscover"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/version"

	"context"
	"fmt"
	"sync"
	"math/rand"
	"strconv"
)

// RegDiscover register and discover
type RegDiscover struct {
	ip      string
	port    uint
	isSSL   bool
	rd      *RegisterDiscover.RegDiscover
	rootCtx context.Context
	cancel  context.CancelFunc

	objectCtrlServs []*types.ObjectControllerServInfo
	objCtrlLock     sync.RWMutex
}

// NewRegDiscover create a RegDiscover object
func NewRegDiscover(zkserv string, ip string, port uint, isSSL bool) *RegDiscover {
	return &RegDiscover{
		ip:              ip,
		port:            port,
		isSSL:           isSSL,
		rd:              RegisterDiscover.NewRegDiscoverEx(zkserv, 10*time.Second),
		objectCtrlServs: []*types.ObjectControllerServInfo{},
	}
}

// Ping the server
func (r *RegDiscover) Ping() error {
	return r.rd.Ping()
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

	// register data collection
	if err := r.registerDataCollection(); err != nil {
		blog.Error("fail to register datacollection(%s), err:%s", r.ip, err.Error())
		return err
	}

	// here: discover other services
	// object-ctrl server
	objCtrlPath := types.CC_SERV_BASEPATH + "/" + types.CC_MODULE_OBJECTCONTROLLER
	objCtrlEvent, err := r.rd.DiscoverService(objCtrlPath)
	if err != nil {
		blog.Errorf("fail to register discover for objectctrl server. err:%s", err.Error())
		return err
	}

	for {
		select {
		case objCtrlEnv := <-objCtrlEvent:
			r.discoverObjectCtrlServ(objCtrlEnv.Server)
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

func (r *RegDiscover) registerDataCollection() error {
	dataCollection := new(types.DataCollectionControllerServInfo)

	dataCollection.IP = r.ip
	dataCollection.Port = r.port
	dataCollection.Scheme = "http"
	if r.isSSL {
		dataCollection.Scheme = "https"
	}

	dataCollection.Version = version.GetVersion()
	dataCollection.Pid = os.Getpid()

	data, err := json.Marshal(dataCollection)
	if err != nil {
		blog.Error("fail to marshal DataCollection server info to json. err:%s", err.Error())
		return err
	}

	path := types.CC_SERV_BASEPATH + "/" + types.CC_MODULE_DATACOLLECTION + "/" + r.ip

	return r.rd.RegisterAndWatchService(path, data)
}

// GetServer fetch server info
func (r *RegDiscover) GetServer(servType string) (string, error) {
	switch servType {
	case types.CC_MODULE_OBJECTCONTROLLER:
		return r.GetObjectCtrlServ()
	}

	err := fmt.Errorf("there is no server discover for type(%s)", servType)
	blog.Errorf("%s", err.Error())
	return "", err
}

//GetObjectCtrlServ fetch objectcontroller server info
func (r *RegDiscover) GetObjectCtrlServ() (string, error) {

	r.objCtrlLock.RLock()
	defer r.objCtrlLock.RUnlock()

	lServ := len(r.objectCtrlServs)
	if lServ <= 0 {
		err := fmt.Errorf("there is no object-ctrl servers")
		blog.Errorf("%s", err.Error())
		return "", err
	}

	//rand
	rand.Seed(int64(time.Now().Nanosecond()))
	servInfo := r.objectCtrlServs[rand.Intn(lServ)]

	host := servInfo.Scheme + "://" + servInfo.IP + ":" + strconv.Itoa(int(servInfo.Port))

	return host, nil
}

func (r *RegDiscover) discoverObjectCtrlServ(servInfos []string) error {
	blog.Infof("discover ObjectCtrl(%v)", servInfos)

	objCtrlServs := []*types.ObjectControllerServInfo{}
	for _, serv := range servInfos {
		objCtrl := new(types.ObjectControllerServInfo)
		if err := json.Unmarshal([]byte(serv), objCtrl); err != nil {
			blog.Warnf("fail to do json unmarshal(%s), err:%s", serv, err.Error())
			continue
		}

		objCtrlServs = append(objCtrlServs, objCtrl)
	}

	r.objCtrlLock.Lock()
	defer r.objCtrlLock.Unlock()
	r.objectCtrlServs = objCtrlServs

	return nil
}
