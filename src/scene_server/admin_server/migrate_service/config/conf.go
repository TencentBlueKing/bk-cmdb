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
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/language"
	"encoding/json"

	"configcenter/src/common/blog"
	"configcenter/src/common/confregdiscover"
	"configcenter/src/common/errors"
	"configcenter/src/common/types"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

// ConfCenter discover configure changed. get, update configures
type ConfCenter struct {
	confRegDiscv *confregdiscover.ConfRegDiscover
	rootCtx      context.Context
	cancel       context.CancelFunc
	ctx          []byte
	errorcode    map[string]errors.ErrorCode
	ctxLock      sync.RWMutex
}

// NewConfCenter create a ConfCenter object
func NewConfCenter(serv string) *ConfCenter {
	return &ConfCenter{
		ctx:          nil,
		confRegDiscv: confregdiscover.NewConfRegDiscover(serv),
	}
}

// Start the configure center module service
func (cc *ConfCenter) Start(confDir, errres string) error {
	// create root context
	cc.rootCtx, cc.cancel = context.WithCancel(context.Background())

	// start configure register and discover service
	if err := cc.confRegDiscv.Start(); err != nil {
		blog.Errorf("fail to start config register and discover service. err:%s", err.Error())
		return err
	}

	// save configures
	cc.WriteConfs2Center(confDir)

	// here: no need to discover itselft configure
	// confPath := types.CC_SERVCONF_BASEPATH + "/" + types.CC_MODULE_MIGRATE
	// confEvent, err := cc.confRegDiscv.DiscoverConfig(confPath)

	if err := cc.WriteErrorRes2Center(errres); err != nil {
		blog.Errorf("fail to write languate packages to center, err:%s", err.Error())
	} else {
		blog.Infof("writed languate packages to center %v", types.CC_SERVERROR_BASEPATH)
	}

	errorResEvent, err := cc.confRegDiscv.DiscoverConfig(types.CC_SERVERROR_BASEPATH)
	if err != nil {
		blog.Errorf("fail to discover configure for migrate service. err:%s", err.Error())
		return err
	}

	go func() {
		for {
			select {
			case confEvn := <-errorResEvent:
				cc.dealErrorResEvent(confEvn.Data)
			case <-cc.rootCtx.Done():
				blog.Warn("configure discover service done")
			}
		}
	}()
	return nil
}

// Stop the configure center
func (cc *ConfCenter) Stop() error {
	cc.cancel()

	cc.confRegDiscv.Stop()

	return nil
}

// GetConfigureCtx fetch the configure
func (cc *ConfCenter) GetConfigureCxt() []byte {
	cc.ctxLock.RLock()
	defer cc.ctxLock.RUnlock()

	return cc.ctx
}

func (cc *ConfCenter) GetErrorCxt() map[string]errors.ErrorCode {
	cc.ctxLock.RLock()
	defer cc.ctxLock.RUnlock()

	return cc.errorcode
}

func (cc *ConfCenter) dealConfChangeEvent(data []byte) error {
	blog.Info("%s configure has changed", types.CC_MODULE_MIGRATE)

	cc.ctxLock.Lock()
	defer cc.ctxLock.Unlock()

	cc.ctx = data

	return nil
}

func (cc *ConfCenter) dealErrorResEvent(data []byte) error {
	blog.Info("language has changed")

	cc.ctxLock.Lock()
	defer cc.ctxLock.Unlock()

	errorcode := map[string]errors.ErrorCode{}
	err := json.Unmarshal(data, &errorcode)
	if err != nil {
		return err
	}

	cc.errorcode = errorcode
	a := api.GetAPIResource()
	if a.Error != nil {
		a.Error.Load(errorcode)
	}

	blog.InfoJSON("loaded language package: %s", errorcode)

	return nil
}

func (cc *ConfCenter) WriteErrorRes2Center(errorres string) error {
	info, err := os.Stat(errorres)
	if os.ErrNotExist == err {
		return fmt.Errorf("directory %s not exists", errorres)
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not directory", errorres)
	}

	errcode, err := errors.LoadErrorResourceFromDir(errorres)
	if err != nil {
		return fmt.Errorf("load error resource error: %s", err)
	}

	data, err := json.Marshal(errcode)
	key := types.CC_SERVERROR_BASEPATH
	return cc.confRegDiscv.Write(key, data)
}

func (cc *ConfCenter) WriteLanguageRes2Center(languageres string) error {
	info, err := os.Stat(languageres)
	if os.ErrNotExist == err {
		return fmt.Errorf("directory %s not exists", languageres)
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not directory", languageres)
	}

	languagepack, err := language.LoadLanguageResourceFromDir(languageres)
	if err != nil {
		return fmt.Errorf("load error resource error: %s", err)
	}

	data, err := json.Marshal(languagepack)
	key := types.CC_SERVLANG_BASEPATH
	return cc.confRegDiscv.Write(key, data)
}

//WriteConfs2Center save configurs into center.
// parameter[confRootPath] define the configurs root path, the specification name of the configure \
// file is [modulename].conf \
func (cc *ConfCenter) WriteConfs2Center(confRootPath string) error {
	modules := make([]string, 0)

	modules = append(modules, types.CC_MODULE_APISERVER)
	modules = append(modules, types.CC_MODULE_AUDITCONTROLLER)
	modules = append(modules, types.CC_MODULE_DATACOLLECTION)
	modules = append(modules, types.CC_MODULE_HOST)
	modules = append(modules, types.CC_MODULE_HOSTCONTROLLER)
	// modules = append(modules, types.CC_MODULE_MIGRATE)
	modules = append(modules, types.CC_MODULE_OBJECTCONTROLLER)
	modules = append(modules, types.CC_MODULE_PROC)
	modules = append(modules, types.CC_MODULE_PROCCONTROLLER)
	modules = append(modules, types.CC_MODULE_TOPO)
	modules = append(modules, types.CC_MODULE_WEBSERVER)
	modules = append(modules, types.CC_MODULE_EVENTSERVER)

	for _, moduleName := range modules {
		filePath := confRootPath + "/" + moduleName + ".conf"
		key := types.CC_SERVCONF_BASEPATH + "/" + moduleName
		if err := cc.writeConfigure(filePath, key); err != nil {
			blog.Warnf("fail to write configure of module(%s) into center", moduleName)
			continue
		}
	}

	return nil
}

func (cc *ConfCenter) writeConfigure(confFilePath, key string) error {
	confFile, err := os.Open(confFilePath)
	if err != nil {
		blog.Errorf("fail to open file(%s), err(%s)", confFilePath, err.Error())
		return err
	}
	defer confFile.Close()

	data, err := ioutil.ReadAll(confFile)
	if err != nil {
		blog.Errorf("fail to read all data from config file(%s), err:%s", confFilePath, err.Error())
		return err
	}

	blog.Debug("write configure(%s), key(%s), data(%s)", confFilePath, key, data)
	if err := cc.confRegDiscv.Write(key, data); err != nil {
		blog.Errorf("fail to write configure(%s) data into center. err:%s", key, err.Error())
		return err
	}

	return nil
}
