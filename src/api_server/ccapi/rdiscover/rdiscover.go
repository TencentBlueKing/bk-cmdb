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
	"math/rand"
	"os"
	"strconv"
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
	ip      string
	port    uint
	isSSL   bool
	rd      *RegisterDiscover.RegDiscover
	rootCtx context.Context
	cancel  context.CancelFunc

	serverLock sync.RWMutex
	servers    map[string][]types.ServerInfo
}

// NewRegDiscover create a RegDiscover object
func NewRegDiscover(zkserv string, ip string, port uint, isSSL bool) *RegDiscover {
	return &RegDiscover{
		ip:      ip,
		port:    port,
		isSSL:   isSSL,
		rd:      RegisterDiscover.NewRegDiscoverEx(zkserv, 10*time.Second),
		servers: map[string][]types.ServerInfo{},
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

	// register apiserver server
	if err := r.registerAPIServer(); err != nil {
		blog.Error("fail to register apiserver(%s), err:%s", r.ip, err.Error())
		return err
	}

	// here: discover other services
	for module := range types.AllModule {
		if module == types.CC_MODULE_APISERVER {
			continue
		}
		modulePath := types.CC_SERV_BASEPATH + "/" + types.CC_MODULE_HOST
		rdEvent, err := r.rd.DiscoverService(modulePath)
		if err != nil {
			blog.Errorf("fail to register discover for host_server. err:%s", err.Error())
			return err
		}
		go r.discover(module, rdEvent)
	}

	<-r.rootCtx.Done()
	blog.Warn("register and discover serv done")
	return nil
}

// GetServer returns the discovered server
func (r *RegDiscover) GetServer(servType string) (string, error) {
	r.serverLock.RLock()
	defer r.serverLock.RUnlock()

	servers, ok := r.servers[servType]
	if !ok {
		err := fmt.Errorf("there is no server discover for type(%s)", servType)
		blog.Errorf("%s", err.Error())
		return "", err
	}

	lServ := len(servers)
	if lServ <= 0 {
		err := fmt.Errorf(`there is no "%s" servers`, servType)
		blog.Errorf("%s", err.Error())
		return "", err
	}

	//rand
	rand.Seed(int64(time.Now().Nanosecond()))
	servInfo := servers[rand.Intn(lServ)]

	host := servInfo.Scheme + "://" + servInfo.IP + ":" + strconv.Itoa(int(servInfo.Port))

	return host, nil

}

func (r *RegDiscover) discover(module string, rdEvent <-chan *RegisterDiscover.DiscoverEvent) {
	for {
		servInfos := <-rdEvent
		blog.Infof("discover host_server(%#v)", servInfos)

		hosts := []types.ServerInfo{}
		for _, serv := range servInfos.Server {
			host := types.ServerInfo{}
			if err := json.Unmarshal([]byte(serv), &host); err != nil {
				blog.Warnf("fail to do json unmarshal(%s), err:%s", serv, err.Error())
				continue
			}

			hosts = append(hosts, host)
		}

		r.serverLock.Lock()
		r.servers[module] = hosts
		r.serverLock.Unlock()
	}
}

// Stop the register and discover
func (r *RegDiscover) Stop() error {
	r.cancel()

	r.rd.Stop()

	return nil
}

func (r *RegDiscover) registerAPIServer() error {
	apiServInfo := new(types.APIServerServInfo)

	apiServInfo.IP = r.ip
	apiServInfo.Port = r.port
	apiServInfo.Scheme = "http"
	if r.isSSL {
		apiServInfo.Scheme = "https"
	}

	apiServInfo.Version = version.GetVersion()
	apiServInfo.Pid = os.Getpid()

	data, err := json.Marshal(apiServInfo)
	if err != nil {
		blog.Errorf("fail to marshal APIServer server info to json. err:%s", err.Error())
		return err
	}

	path := types.CC_SERV_BASEPATH + "/" + types.CC_MODULE_APISERVER + "/" + r.ip

	return r.rd.RegisterAndWatchService(path, data)
}
