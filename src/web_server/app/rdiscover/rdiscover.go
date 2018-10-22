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
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"configcenter/src/common/RegisterDiscover"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
)

// RegDiscover register and discover
type RegDiscover struct {
	ip      string
	port    uint
	isSSL   bool
	rd      *RegisterDiscover.RegDiscover
	rootCtx context.Context
	cancel  context.CancelFunc
	apiSevs []*types.APIServerServInfo
	apiLock sync.RWMutex
}

// NewRegDiscover create a RegDiscover object
func NewRegDiscover(zkserv string, ip string, port uint, isSSL bool) *RegDiscover {
	return &RegDiscover{
		ip:      ip,
		port:    port,
		isSSL:   isSSL,
		rd:      RegisterDiscover.NewRegDiscoverEx(zkserv, 10*time.Second),
		apiSevs: []*types.APIServerServInfo{},
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

	// register web server
	if err := r.registerWebServer(); err != nil {
		blog.Errorf("fail to register web(%s), err:%s", r.ip, err.Error())
		return err
	}

	// here: discover other servers
	/// api-server
	apiServPath := types.CC_SERV_BASEPATH + "/" + types.CC_MODULE_APISERVER
	apiServEvent, err := r.rd.DiscoverService(apiServPath)
	if err != nil {
		blog.Errorf("fail to register discover for api server. err:%s", err.Error())
		return err
	}

	for {
		select {
		case apiServEnv := <-apiServEvent:
			r.discoverAPIServer(apiServEnv.Server)
		case <-r.rootCtx.Done():
			blog.Warn("register and discover serv done")
			return nil
		}
	}
}

func (r *RegDiscover) GetServer(servType string) (string, error) {
	switch servType {
	case types.CC_MODULE_APISERVER:
		return r.GetAPIServ()
	}

	err := fmt.Errorf("there is no server discover for type(%s)", servType)
	blog.Errorf("%s", err.Error())

	return "", err
}

// GetAPIServ get api-server
func (r *RegDiscover) GetAPIServ() (string, error) {
	r.apiLock.RLock()
	defer r.apiLock.RUnlock()

	lServ := len(r.apiSevs)
	if lServ <= 0 {
		err := fmt.Errorf("there is no api servers")
		blog.Errorf("%s", err.Error())
		return "", err
	}

	// rand
	rand.Seed(int64(time.Now().Nanosecond()))
	servInfo := r.apiSevs[rand.Intn(lServ)]

	host := servInfo.Scheme + "://" + servInfo.IP + ":" + strconv.Itoa(int(servInfo.Port))

	return host, nil
}

// Stop the register and discover
func (r *RegDiscover) Stop() error {
	r.cancel()

	r.rd.Stop()

	return nil
}

func (r *RegDiscover) registerWebServer() error {
	webServInfo := new(types.WebServerInfo)

	webServInfo.IP = r.ip
	webServInfo.Port = r.port
	webServInfo.Scheme = "http"
	if r.isSSL {
		webServInfo.Scheme = "https"
	}

	webServInfo.Version = version.GetVersion()
	webServInfo.Pid = os.Getpid()

	data, err := json.Marshal(webServInfo)
	if err != nil {
		blog.Error("fail to marshal web server info to json. err:%s", err.Error())
		return err
	}

	path := types.CC_SERV_BASEPATH + "/" + types.CC_MODULE_WEBSERVER + "/" + r.ip

	return r.rd.RegisterAndWatchService(path, data)
}

func (r *RegDiscover) discoverAPIServer(servInfos []string) error {
	blog.Infof("discover apiserver(%v)", servInfos)

	apiServs := []*types.APIServerServInfo{}
	for _, serv := range servInfos {
		apiServ := new(types.APIServerServInfo)
		if err := json.Unmarshal([]byte(serv), apiServ); err != nil {
			blog.Warnf("fail to do json unmarshal(%s), err:%s", serv, err.Error())
			continue
		}

		apiServs = append(apiServs, apiServ)
	}

	r.apiLock.Lock()
	defer r.apiLock.Unlock()
	r.apiSevs = apiServs

	return nil
}
