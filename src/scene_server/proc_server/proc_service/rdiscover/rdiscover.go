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
	ip              string
	port            uint
	isSSL           bool
	rd              *RegisterDiscover.RegDiscover
	rootCtx         context.Context
	cancel          context.CancelFunc
	hostCtrlServs   []*types.HostControllerServInfo
	hostCtrlLock    sync.RWMutex
	objectCtrlServs []*types.ObjectControllerServInfo
	objCtrlLock     sync.RWMutex
	auditCtrlServs  []*types.AuditControllerServInfo
	auditCtrlLock   sync.RWMutex
	procCtrlServs   []*types.ProcControllerServInfo
	procCtrlLock    sync.RWMutex
}

// NewRegDiscover create a RegDiscover object
func NewRegDiscover(zkserv string, ip string, port uint, isSSL bool) *RegDiscover {
	return &RegDiscover{
		ip:              ip,
		port:            port,
		isSSL:           isSSL,
		rd:              RegisterDiscover.NewRegDiscoverEx(zkserv, 10*time.Second),
		hostCtrlServs:   []*types.HostControllerServInfo{},
		objectCtrlServs: []*types.ObjectControllerServInfo{},
		auditCtrlServs:  []*types.AuditControllerServInfo{},
		procCtrlServs:   []*types.ProcControllerServInfo{},
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

	// register process server
	if err := r.registerProc(); err != nil {
		blog.Error("fail to register proc(%s), err:%s", r.ip, err.Error())
		return err
	}

	// here: discover other services
	/// host-ctrl server
	hostCtrlPath := types.CC_SERV_BASEPATH + "/" + types.CC_MODULE_HOSTCONTROLLER
	hostCtrlEvent, err := r.rd.DiscoverService(hostCtrlPath)
	if err != nil {
		blog.Errorf("fail to register discover for hostctrl server. err:%s", err.Error())
		return err
	}

	/// object-ctrl server
	objCtrlPath := types.CC_SERV_BASEPATH + "/" + types.CC_MODULE_OBJECTCONTROLLER
	objCtrlEvent, err := r.rd.DiscoverService(objCtrlPath)
	if err != nil {
		blog.Errorf("fail to register discover for objectctrl server. err:%s", err.Error())
		return err
	}

	/// AuditController server
	auditCtrlPath := types.CC_SERV_BASEPATH + "/" + types.CC_MODULE_AUDITCONTROLLER
	auditCtrlEvent, err := r.rd.DiscoverService(auditCtrlPath)
	if err != nil {
		blog.Errorf("fail to register discover to auditctrl server. err:%s", err.Error())
		return err
	}

	/// ProcController server
	procCtrlPath := types.CC_SERV_BASEPATH + "/" + types.CC_MODULE_PROCCONTROLLER
	procCtrlEvent, err := r.rd.DiscoverService(procCtrlPath)
	if err != nil {
		blog.Errorf("fail to register discover to procctrl server. err:%s", err.Error())
		return err
	}

	for {
		select {
		case hostCtrlEnv := <-hostCtrlEvent:
			r.discoverHostCtrlServ(hostCtrlEnv.Server)
		case objCtrlEnv := <-objCtrlEvent:
			r.discoverObjectCtrlServ(objCtrlEnv.Server)
		case auditCtrlEnv := <-auditCtrlEvent:
			r.discoverAuditCtrlServ(auditCtrlEnv.Server)
		case procCtrlEnv := <-procCtrlEvent:
			r.discoverProcCtrlServ(procCtrlEnv.Server)
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

//GetServer fetch server
func (r *RegDiscover) GetServer(servType string) (string, error) {

	switch servType {
	case types.CC_MODULE_AUDITCONTROLLER:
		return r.GetAuditCtrlServ()
	case types.CC_MODULE_HOSTCONTROLLER:
		return r.GetHostCtrlServ()
	case types.CC_MODULE_OBJECTCONTROLLER:
		return r.GetObjectCtrlServ()
	case types.CC_MODULE_PROCCONTROLLER:
		return r.GetProcCtrlServ()
	}

	err := fmt.Errorf("there is no server discover for type(%s)", servType)
	blog.Errorf("%s", err.Error())
	return "", err
}

//GetHostCtrlServ get hostcontroller server info
func (r *RegDiscover) GetHostCtrlServ() (string, error) {

	r.hostCtrlLock.RLock()
	defer r.hostCtrlLock.RUnlock()

	lServ := len(r.hostCtrlServs)
	if lServ <= 0 {
		err := fmt.Errorf("there is no host-ctrl servers")
		blog.Errorf("%s", err.Error())
		return "", err
	}

	//rand
	rand.Seed(int64(time.Now().Nanosecond()))
	servInfo := r.hostCtrlServs[rand.Intn(lServ)]

	host := servInfo.Scheme + "://" + servInfo.IP + ":" + strconv.Itoa(int(servInfo.Port))

	return host, nil
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

//GetAuditCtrlServ fetch auditcontroller server info
func (r *RegDiscover) GetAuditCtrlServ() (string, error) {

	r.auditCtrlLock.RLock()
	defer r.auditCtrlLock.RUnlock()

	lServ := len(r.auditCtrlServs)
	if lServ <= 0 {
		err := fmt.Errorf("there is no audit-ctrl servers")
		blog.Errorf("%s", err.Error())
		return "", err
	}

	//rand
	rand.Seed(int64(time.Now().Nanosecond()))
	servInfo := r.auditCtrlServs[rand.Intn(lServ)]

	host := servInfo.Scheme + "://" + servInfo.IP + ":" + strconv.Itoa(int(servInfo.Port))

	return host, nil
}

//GetProcCtrlServ fetch proccontroller server info
func (r *RegDiscover) GetProcCtrlServ() (string, error) {

	r.procCtrlLock.RLock()
	defer r.procCtrlLock.RUnlock()

	lServ := len(r.procCtrlServs)
	if lServ <= 0 {
		err := fmt.Errorf("there is no proc-ctrl servers")
		blog.Errorf("%s", err.Error())
		return "", err
	}

	//rand
	rand.Seed(int64(time.Now().Nanosecond()))
	servInfo := r.procCtrlServs[rand.Intn(lServ)]

	host := servInfo.Scheme + "://" + servInfo.IP + ":" + strconv.Itoa(int(servInfo.Port))

	return host, nil
}

func (r *RegDiscover) registerProc() error {
	procServInfo := new(types.ProcServInfo)

	procServInfo.IP = r.ip
	procServInfo.Port = r.port
	procServInfo.Scheme = "http"
	if r.isSSL {
		procServInfo.Scheme = "https"
	}

	procServInfo.Version = version.GetVersion()
	procServInfo.Pid = os.Getpid()

	data, err := json.Marshal(procServInfo)
	if err != nil {
		blog.Error("fail to marshal Proc server info to json. err:%s", err.Error())
		return err
	}

	path := types.CC_SERV_BASEPATH + "/" + types.CC_MODULE_PROC + "/" + r.ip

	return r.rd.RegisterAndWatchService(path, data)
}

func (r *RegDiscover) discoverHostCtrlServ(servInfos []string) error {
	blog.Infof("discover HostCtrl(%v)", servInfos)

	hostCtrlServs := []*types.HostControllerServInfo{}
	for _, serv := range servInfos {
		hostCtrl := new(types.HostControllerServInfo)
		if err := json.Unmarshal([]byte(serv), hostCtrl); err != nil {
			blog.Warnf("fail to do json unmarshal(%s), err:%s", serv, err.Error())
			continue
		}

		hostCtrlServs = append(hostCtrlServs, hostCtrl)
	}

	r.hostCtrlLock.Lock()
	defer r.hostCtrlLock.Unlock()
	r.hostCtrlServs = hostCtrlServs

	return nil
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

func (r *RegDiscover) discoverAuditCtrlServ(servInfos []string) error {
	blog.Infof("discover AuditCtrl(%v)", servInfos)

	auditCtrlServs := []*types.AuditControllerServInfo{}
	for _, serv := range servInfos {
		auditCtrl := new(types.AuditControllerServInfo)
		if err := json.Unmarshal([]byte(serv), auditCtrl); err != nil {
			blog.Warnf("fail to do json unmarshal(%s), err:%s", serv, err.Error())
			continue
		}

		auditCtrlServs = append(auditCtrlServs, auditCtrl)
	}

	r.auditCtrlLock.Lock()
	defer r.auditCtrlLock.Unlock()
	r.auditCtrlServs = auditCtrlServs

	return nil
}

func (r *RegDiscover) discoverProcCtrlServ(servInfos []string) error {
	blog.Infof("discover ProcCtrl(%v)", servInfos)

	procCtrlServs := []*types.ProcControllerServInfo{}
	for _, serv := range servInfos {
		procCtrl := new(types.ProcControllerServInfo)
		if err := json.Unmarshal([]byte(serv), procCtrl); err != nil {
			blog.Warnf("fail to do json unmarshal(%s), err:%s", serv, err.Error())
			continue
		}

		procCtrlServs = append(procCtrlServs, procCtrl)
	}

	r.procCtrlLock.Lock()
	defer r.procCtrlLock.Unlock()
	r.procCtrlServs = procCtrlServs

	return nil
}
