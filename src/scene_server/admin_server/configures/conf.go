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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"configcenter/src/common/backbone/service_mange/zk"
	"configcenter/src/common/blog"
	"configcenter/src/common/confregdiscover"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/types"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ConfCenter discover configure changed. get, update configures
type ConfCenter struct {
	confRegDiscv confregdiscover.ConfRegDiscvIf
	ctx          context.Context
}

// NewConfCenter create a ConfCenter object
func NewConfCenter(ctx context.Context, client *zk.ZkClient) *ConfCenter {
	return &ConfCenter{
		ctx:          ctx,
		confRegDiscv: confregdiscover.NewZkRegDiscover(client),
	}
}

// Ping to ping server
func (cc *ConfCenter) Ping() error {
	return cc.confRegDiscv.Ping()
}

// Start the configure center module service
func (cc *ConfCenter) Start(confDir, errRes, languageRes string) error {

	// save configures
	if err := cc.writeConfs2Center(confDir); err != nil {
		blog.Errorf("fail to write configures to center, err:%s", err.Error())
		return err
	} else {
		blog.Infof("write all configures resource to center %v success", types.CC_SERVCONF_BASEPATH)
	}

	if err := cc.writeErrorRes2Center(errRes); err != nil {
		blog.Errorf("fail to write error resource to center, err:%s", err.Error())
		return err
	} else {
		blog.Infof("write error resource to center %v success", types.CC_SERVERROR_BASEPATH)
	}

	if err := cc.writeLanguageRes2Center(languageRes); err != nil {
		blog.Errorf("fail to write language packages to center, err:%s", err.Error())
		return err
	} else {
		blog.Infof("write language packages to center %v success", types.CC_SERVLANG_BASEPATH)
	}

	// TODO discover config file change
	go func() {
		select {
		case <-cc.ctx.Done():
		}
	}()
	return nil
}

func (cc *ConfCenter) writeErrorRes2Center(errorres string) error {
	info, err := os.Stat(errorres)
	if os.ErrNotExist == err {
		return fmt.Errorf("directory %s not exists", errorres)
	}
	if err != nil {
		return fmt.Errorf("stat directory %s faile, %s", errorres, err.Error())
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not directory", errorres)
	}

	errcode, err := errors.LoadErrorResourceFromDir(errorres)
	if err != nil {
		return fmt.Errorf("load error resource error: %s", err)
	}

	data, err := json.Marshal(errcode)
	if err != nil {
		return fmt.Errorf("unmarshal resource failed, err: %s", err)
	}
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
	if err != nil {
		return err
	}
	key := types.CC_SERVLANG_BASEPATH
	return cc.confRegDiscv.Write(key, data)
}

// WriteConfs2Center save configures into center.
// parameter[confRootPath] define the configure's root path, the configure files are
// redis.conf, mongodb.conf，common.conf，extra.conf
func (cc *ConfCenter) writeConfs2Center(confRootPath string) error {
	configs := []string{
		types.CCConfigureRedis,
		types.CCConfigureMongo,
		types.CCConfigureCommon,
		types.CCConfigureExtra,
	}

	confFileSuffix := ".yaml"

	for _, configName := range configs {
		filePath := filepath.Join(confRootPath, configName+confFileSuffix)
		key := types.CC_SERVCONF_BASEPATH + "/" + configName
		if err := cc.writeConfigure(filePath, key); err != nil {
			blog.Warnf("fail to write configure of %s into center", configName)
			continue
		} else {
			blog.Infof("write configure to center %s success", key)
			cc.listenFileChange(key,filePath)
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

	// check the configuration in the file
	if err := cc.checkFile(confFilePath); err != nil {
		blog.Errorf("There is a problem in configuration file %s, err:%s", confFilePath, err)
		os.Exit(1)
	}

	blog.V(3).Infof("write configure(%s), key(%s), data(%s)", confFilePath, key, data)
	if err := cc.confRegDiscv.Write(key, data); err != nil {
		blog.Errorf("fail to write configure(%s) data into center. err:%s", key, err.Error())
		return err
	}

	return nil
}

var redisViper *viper.Viper
var mongodbViper *viper.Viper
var commonViper *viper.Viper
var extraViper *viper.Viper

//此方法给adminserver实现热更新,监听每个文件，当文件发生更改时，将改后的数据重新写到注册中心
func (cc *ConfCenter) listenFileChange(configcenterPath string,filePath string) {
	v := viper.New()
	base := path.Base(filePath)
	split := strings.Split(base, ".")
	fileName := split[0]
	v.SetConfigName(fileName)
	v.AddConfigPath(path.Dir(filePath))
	err := v.ReadInConfig()
	if err != nil {
		blog.Warnf("fail to read configure from %s ", base)
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		if err := cc.writeConfigure(filePath, configcenterPath); err != nil {
			blog.Warnf("fail to write configure of %s into center", base)
		} else {
			blog.Infof("write configure to center %s success", configcenterPath)
		}
	})
	if fileName == types.CCConfigureRedis {
		redisViper = v
	} else if fileName == types.CCConfigureMongo {
		mongodbViper = v
	} else if fileName == types.CCConfigureCommon {
		commonViper = v
	} else if fileName == types.CCConfigureExtra {
		extraViper = v
	}
}