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

package configures

import (
	"encoding/json"
	"time"

	"context"
	"fmt"
	"io/ioutil"
	"os"

	"configcenter/src/common/blog"
	"configcenter/src/common/confregdiscover"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/types"
)

// ConfCenter discover configure changed. get, update configures
type ConfCenter struct {
	confRegDiscv confregdiscover.ConfRegDiscvIf
	ctx          context.Context
}

// NewConfCenter create a ConfCenter object
func NewConfCenter(ctx context.Context, serv string) *ConfCenter {
	return &ConfCenter{
		ctx:          ctx,
		confRegDiscv: confregdiscover.NewZkRegDiscover(serv, time.Second*60),
	}
}

// Ping to ping server
func (cc *ConfCenter) Ping() error {
	return cc.confRegDiscv.Ping()
}

// Start the configure center module service
func (cc *ConfCenter) Start(confDir, errres, languageres string) error {
	// start configure register and discover service
	if err := cc.confRegDiscv.Start(); err != nil {
		blog.Errorf("fail to start config register and discover service. err:%s", err.Error())
		return err
	}

	// save configures
	if err := cc.writeConfs2Center(confDir); err != nil {
		blog.Errorf("fail to write configures to center, err:%s", err.Error())
		return err
	} else {
		blog.Infof("writed all configures resource to center %v", types.CC_SERVCONF_BASEPATH)
	}

	if err := cc.writeErrorRes2Center(errres); err != nil {
		blog.Errorf("fail to write error resource to center, err:%s", err.Error())
		return err
	} else {
		blog.Infof("writed error resource to center %v", types.CC_SERVERROR_BASEPATH)
	}

	if err := cc.writeLanguageRes2Center(languageres); err != nil {
		blog.Errorf("fail to write languate packages to center, err:%s", err.Error())
		return err
	} else {
		blog.Infof("writed languate packages to center %v", types.CC_SERVLANG_BASEPATH)
	}

	// TODO discover config file change
	go func() {
		select {
		case <-cc.ctx.Done():
			cc.Stop()
		}
	}()
	return nil
}

// Stop the configure center
func (cc *ConfCenter) Stop() error {
	cc.confRegDiscv.Stop()
	return nil
}

func (cc *ConfCenter) writeErrorRes2Center(errorres string) error {
	info, err := os.Stat(errorres)
	if os.ErrNotExist == err {
		return fmt.Errorf("directory %s not exists", errorres)
	}
	if err != nil {
		return fmt.Errorf("Stat directory %s faile, %s", errorres, err.Error())
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

func (cc *ConfCenter) writeLanguageRes2Center(languageres string) error {
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
		return fmt.Errorf("load language resource error: %s", err)
	}

	data, err := json.Marshal(languagepack)
	key := types.CC_SERVLANG_BASEPATH
	return cc.confRegDiscv.Write(key, data)
}

//WriteConfs2Center save configurs into center.
// parameter[confRootPath] define the configurs root path, the specification name of the configure \
// file is [modulename].conf \
func (cc *ConfCenter) writeConfs2Center(confRootPath string) error {
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
		} else {
			blog.Infof("writed configure to center %s", key)
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

	blog.V(3).Infof("write configure(%s), key(%s), data(%s)", confFilePath, key, data)
	if err := cc.confRegDiscv.Write(key, data); err != nil {
		blog.Errorf("fail to write configure(%s) data into center. err:%s", key, err.Error())
		return err
	}

	return nil
}
