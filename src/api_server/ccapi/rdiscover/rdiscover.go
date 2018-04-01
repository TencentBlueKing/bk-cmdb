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
	ip         string
	port       uint
	isSSL      bool
	rd         *RegisterDiscover.RegDiscover
	rootCtx    context.Context
	cancel     context.CancelFunc
	hostServs  []*types.HostServerInfo
	hostLock   sync.RWMutex
	topoServs  []*types.TopoServInfo
	eventLock  sync.RWMutex
	eventServs []*types.EventServInfo
	topoLock   sync.RWMutex
	procServs  []*types.ProcServInfo
	procLock   sync.RWMutex
}

// NewRegDiscover create a RegDiscover object
func NewRegDiscover(zkserv string, ip string, port uint, isSSL bool) *RegDiscover {
	return &RegDiscover{
		ip:        ip,
		port:      port,
		isSSL:     isSSL,
		rd:        RegisterDiscover.NewRegDiscoverEx(zkserv, 10*time.Second),
		hostServs: []*types.HostServerInfo{},
		topoServs: []*types.TopoServInfo{},
		procServs: []*types.ProcServInfo{},
	}
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
	/// host server
	hostPath := types.CC_SERV_BASEPATH + "/" + types.CC_MODULE_HOST
	hostEvent, err := r.rd.DiscoverService(hostPath)
	if err != nil {
		blog.Errorf("fail to register discover for host_server. err:%s", err.Error())
		return err
	}

	/// topo server
	topoPath := types.CC_SERV_BASEPATH + "/" + types.CC_MODULE_TOPO
	topoEvent, err := r.rd.DiscoverService(topoPath)
	if err != nil {
		blog.Errorf("fail to register discover for topo_server. err:%s", err.Error())
		return err
	}

	/// proc_server
	procPath := types.CC_SERV_BASEPATH + "/" + types.CC_MODULE_PROC
	procEvent, err := r.rd.DiscoverService(procPath)
	if err != nil {
		blog.Errorf("fail to register discover for proc_server. err:%s", err.Error())
		return err
	}

	/// event_server
	eventPath := types.CC_SERV_BASEPATH + "/" + types.CC_MODULE_EVENTSERVER
	eventEvent, err := r.rd.DiscoverService(eventPath)
	if err != nil {
		blog.Errorf("fail to register discover for proc_server. err:%s", err.Error())
		return err
	}

	for {
		select {
		case hostEnv := <-hostEvent:
			r.discoverHostServ(hostEnv.Server)
		case topoEnv := <-topoEvent:
			r.discoverTopoServ(topoEnv.Server)
		case procEnv := <-procEvent:
			r.discoverProcServ(procEnv.Server)
		case eventEnv := <-eventEvent:
			r.discoverEventServ(eventEnv.Server)
		case <-r.rootCtx.Done():
			blog.Warn("register and discover serv done")
			return nil
		}
	}
}

func (r *RegDiscover) GetServer(servType string) (string, error) {
	switch servType {
	case types.CC_MODULE_PROC:
		return r.GetProcServ()
	case types.CC_MODULE_HOST:
		return r.GetHostServ()
	case types.CC_MODULE_TOPO:
		return r.GetTopoServ()
	case types.CC_MODULE_EVENTSERVER:
		return r.GetEventServ()
	}

	err := fmt.Errorf("there is no server discover for type(%s)", servType)
	blog.Errorf("%s", err.Error())

	return "", err
}

// GetHostServ fetch host server
func (r *RegDiscover) GetHostServ() (string, error) {
	r.hostLock.RLock()
	defer r.hostLock.RUnlock()

	lServ := len(r.hostServs)
	if lServ <= 0 {
		err := fmt.Errorf("there is no host servers")
		blog.Errorf("%s", err.Error())
		return "", err
	}

	//rand
	rand.Seed(int64(time.Now().Nanosecond()))
	servInfo := r.hostServs[rand.Intn(lServ)]

	host := servInfo.Scheme + "://" + servInfo.IP + ":" + strconv.Itoa(int(servInfo.Port))

	return host, nil
}

// GetTopoServ fetch topo server
func (r *RegDiscover) GetTopoServ() (string, error) {
	r.topoLock.RLock()
	defer r.topoLock.RUnlock()

	lServ := len(r.topoServs)
	if lServ <= 0 {
		err := fmt.Errorf("there is no topo servers")
		blog.Errorf("%s", err.Error())
		return "", err
	}

	//rand
	rand.Seed(int64(time.Now().Nanosecond()))
	servInfo := r.topoServs[rand.Intn(lServ)]

	host := servInfo.Scheme + "://" + servInfo.IP + ":" + strconv.Itoa(int(servInfo.Port))

	return host, nil
}

// GetProcServ fetch topo server
func (r *RegDiscover) GetProcServ() (string, error) {
	r.procLock.RLock()
	defer r.procLock.RUnlock()
	lServ := len(r.procServs)
	if lServ <= 0 {
		err := fmt.Errorf("there is no proc servers")
		blog.Errorf("%s", err.Error())
		return "", err
	}

	//rand
	rand.Seed(int64(time.Now().Nanosecond()))
	servInfo := r.procServs[rand.Intn(lServ)]

	host := servInfo.Scheme + "://" + servInfo.IP + ":" + strconv.Itoa(int(servInfo.Port))

	return host, nil
}

// GetTopoServ fetch topo server
func (r *RegDiscover) GetEventServ() (string, error) {
	r.eventLock.RLock()
	defer r.eventLock.RUnlock()

	lServ := len(r.eventServs)
	if lServ <= 0 {
		err := fmt.Errorf("there is no event servers")
		blog.Errorf("%s", err.Error())
		return "", err
	}

	//rand
	rand.Seed(int64(time.Now().Nanosecond()))
	servInfo := r.eventServs[rand.Intn(lServ)]

	host := servInfo.Scheme + "://" + servInfo.IP + ":" + strconv.Itoa(int(servInfo.Port))

	return host, nil
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

func (r *RegDiscover) discoverHostServ(servInfos []string) error {
	blog.Infof("discover host_server(%v)", servInfos)

	hosts := []*types.HostServerInfo{}
	for _, serv := range servInfos {
		host := new(types.HostServerInfo)
		if err := json.Unmarshal([]byte(serv), host); err != nil {
			blog.Warnf("fail to do json unmarshal(%s), err:%s", serv, err.Error())
			continue
		}

		hosts = append(hosts, host)
	}

	r.hostLock.Lock()
	defer r.hostLock.Unlock()
	r.hostServs = hosts

	return nil
}

func (r *RegDiscover) discoverTopoServ(servInfos []string) error {
	blog.Infof("discover topo_server(%v)", servInfos)

	topoServs := []*types.TopoServInfo{}
	for _, serv := range servInfos {
		topo := new(types.TopoServInfo)
		if err := json.Unmarshal([]byte(serv), topo); err != nil {
			blog.Warnf("fail to do json unmarshal(%s), err:%s", serv, err.Error())
			continue
		}

		topoServs = append(topoServs, topo)
	}

	r.topoLock.Lock()
	defer r.topoLock.Unlock()
	r.topoServs = topoServs

	return nil
}

func (r *RegDiscover) discoverProcServ(servInfos []string) error {
	blog.Infof("discover proc_server(%v)", servInfos)

	procServs := []*types.ProcServInfo{}
	for _, serv := range servInfos {
		proc := new(types.ProcServInfo)
		if err := json.Unmarshal([]byte(serv), proc); err != nil {
			blog.Warnf("fail to do json unmarshal(%s), err:%s", serv, err.Error())
			continue
		}

		procServs = append(procServs, proc)
	}

	r.procLock.Lock()
	defer r.procLock.Unlock()
	r.procServs = procServs

	return nil
}

func (r *RegDiscover) discoverEventServ(servInfos []string) error {
	blog.Infof("discover event_server(%v)", servInfos)

	eventServs := []*types.EventServInfo{}
	for _, serv := range servInfos {
		event := new(types.EventServInfo)
		if err := json.Unmarshal([]byte(serv), event); err != nil {
			blog.Warnf("fail to do json unmarshal(%s), err:%s", serv, err.Error())
			continue
		}

		eventServs = append(eventServs, event)
	}

	r.eventLock.Lock()
	defer r.eventLock.Unlock()
	r.eventServs = eventServs

	return nil
}
