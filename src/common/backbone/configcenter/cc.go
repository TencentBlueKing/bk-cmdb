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

package configcenter

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	"configcenter/src/common/blog"
	crd "configcenter/src/common/confregdiscover"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/types"
)

var confC *CC

func NewConfigCenter(ctx context.Context, disc crd.ConfRegDiscvIf, confPath string, handler *CCHandler) error {
	return New(ctx, confPath, disc, handler)
}

func New(ctx context.Context, confPath string, disc crd.ConfRegDiscvIf, handler *CCHandler) error {
	confC = &CC{
		ctx:           ctx,
		disc:          disc,
		handler:       handler,
		previousProc:  new(ProcessConfig),
		previousLang:  make(map[string]language.LanguageMap),
		previousError: make(map[string]errors.ErrorCode),
	}

	// parse config only from file
	if len(confPath) != 0 {
		return LoadConfigFromLocalFile(confPath, handler)
	}

	if err := confC.run(); err != nil {
		return err
	}

	confC.sync()

	return nil
}

type ProcHandlerFunc func(previous, current ProcessConfig)

type CCHandler struct {
	OnProcessUpdate  ProcHandlerFunc
	OnExtraUpdate    ProcHandlerFunc
	OnLanguageUpdate func(previous, current map[string]language.LanguageMap)
	OnErrorUpdate    func(previous, current map[string]errors.ErrorCode)
	OnMongodbUpdate  func(previous, current ProcessConfig)
	OnRedisUpdate    func(previous, current ProcessConfig)
}

type CC struct {
	sync.Mutex
	// used to stop the config center gracefully.
	ctx             context.Context
	disc            crd.ConfRegDiscvIf
	handler         *CCHandler
	procName        string
	previousProc    *ProcessConfig
	previousExtra   *ProcessConfig
	previousMongodb *ProcessConfig
	previousRedis   *ProcessConfig
	previousLang    map[string]language.LanguageMap
	previousError   map[string]errors.ErrorCode
}

func (c *CC) run() error {
	commonConfPath := fmt.Sprintf("%s/%s", types.CC_SERVCONF_BASEPATH, types.CCConfigureCommon)
	commonConfEvent, err := c.disc.Discover(commonConfPath)
	if err != nil {
		return err
	}

	extraConfPath := fmt.Sprintf("%s/%s", types.CC_SERVCONF_BASEPATH, types.CCConfigureExtra)
	extraConfEvent, err := c.disc.Discover(extraConfPath)
	if err != nil {
		return err
	}

	mongodbConfPath := fmt.Sprintf("%s/%s", types.CC_SERVCONF_BASEPATH, types.CCConfigureMongo)
	mongodbConfEvent, err := c.disc.Discover(mongodbConfPath)
	if err != nil {
		return err
	}

	redisConfPath := fmt.Sprintf("%s/%s", types.CC_SERVCONF_BASEPATH, types.CCConfigureRedis)
	redisConfEvent, err := c.disc.Discover(redisConfPath)
	if err != nil {
		return err
	}

	langEvent, err := c.disc.Discover(types.CC_SERVLANG_BASEPATH)
	if err != nil {
		return err
	}

	errEvent, err := c.disc.Discover(types.CC_SERVERROR_BASEPATH)
	if err != nil {
		return err
	}

	go func() {
		select {
		case pEvent := <-commonConfEvent:
			c.onProcChange(pEvent)
		case pEvent := <-extraConfEvent:
			c.onExtraChange(pEvent)
		case pEvent := <-mongodbConfEvent:
			c.onMongodbChange(pEvent)
		case pEvent := <-redisConfEvent:
			c.onRedisChange(pEvent)
		case eEvent := <-errEvent:
			c.onErrorChange(eEvent)
		case langEvent := <-langEvent:
			c.onLanguageChange(langEvent)
		case <-c.ctx.Done():
			blog.Warnf("config center event watch stopped because of context done.")
			return
		}
	}()
	return nil
}

func (c *CC) onProcChange(cur *crd.DiscoverEvent) {
	if cur.Err != nil {
		blog.Errorf("config center received event that %s config has changed, but got err: %v", types.CCConfigureCommon, cur.Err)
		return
	}

	now := ParseConfigWithData(cur.Data)
	c.Lock()
	defer c.Unlock()
	prev := c.previousProc
	c.previousProc = now
	if err := SetCommonFromByte(now.ConfigData); err != nil {
		blog.Errorf("add updated configuration error: %v", err)
		return
	}
	if c.handler != nil {
		if c.handler.OnProcessUpdate != nil {
			go c.handler.OnProcessUpdate(*prev, *now)
		}
	}
}

func (c *CC) onExtraChange(cur *crd.DiscoverEvent) {
	if cur.Err != nil {
		blog.Errorf("config center received event that %s config has changed, but got err: %v", types.CCConfigureExtra, cur.Err)
		return
	}

	now := ParseConfigWithData(cur.Data)
	c.Lock()
	defer c.Unlock()
	prev := c.previousExtra
	if prev == nil {
		prev = &ProcessConfig{}
	}
	c.previousExtra = now
	if err := SetExtraFromByte(now.ConfigData); err != nil {
		blog.Errorf("add updated extra configuration error: %v", err)
		return
	}
	if c.handler != nil {
		if c.handler.OnExtraUpdate != nil {
			go c.handler.OnExtraUpdate(*prev, *now)
		}
	}
}

func (c *CC) onMongodbChange(cur *crd.DiscoverEvent) {
	if cur.Err != nil {
		blog.Errorf("config center received event that %s config has changed, but got err: %v", types.CCConfigureCommon, cur.Err)
		return
	}
	now := ParseConfigWithData(cur.Data)
	c.Lock()
	defer c.Unlock()
	prev := c.previousMongodb
	if prev == nil {
		prev = &ProcessConfig{}
	}
	c.previousMongodb = now
	if c.handler != nil {
		if c.handler.OnMongodbUpdate != nil {
			go c.handler.OnMongodbUpdate(*prev, *now)
		}
	}
}

func (c *CC) onRedisChange(cur *crd.DiscoverEvent) {
	if cur.Err != nil {
		blog.Errorf("config center received event that %s config has changed, but got err: %v", types.CCConfigureCommon, cur.Err)
		return
	}
	now := ParseConfigWithData(cur.Data)
	c.Lock()
	defer c.Unlock()
	prev := c.previousRedis
	if prev == nil {
		prev = &ProcessConfig{}
	}
	c.previousRedis = now
	if c.handler != nil {
		if c.handler.OnRedisUpdate != nil {
			go c.handler.OnRedisUpdate(*prev, *now)
		}
	}
}

func (c *CC) onErrorChange(cur *crd.DiscoverEvent) {

	if cur.Err != nil {
		blog.Errorf("config center received event that *ERROR CODE* config has changed, but got err: %v", cur.Err)
		return
	}

	now := make(map[string]errors.ErrorCode)
	if err := json.Unmarshal(cur.Data, &now); err != nil {
		blog.Errorf("config center received event that *ERROR CODE* config has changed, but unmarshal err: %v", err)
		return
	}

	c.Lock()
	defer c.Unlock()
	prev := c.previousError
	c.previousError = now

	if c.handler != nil {
		go c.handler.OnErrorUpdate(prev, deepCopyError(now))
	}
}

func (c *CC) onLanguageChange(cur *crd.DiscoverEvent) {
	if cur.Err != nil {
		blog.Errorf("config center received event that *LANGUAGE* config has changed, but got err: %v", cur.Err)
		return
	}

	now := make(map[string]language.LanguageMap)
	if err := json.Unmarshal(cur.Data, &now); err != nil {
		blog.Errorf("config center received event that *LANGUAGE* config has changed, but unmarshal err: %v", err)
		return
	}

	c.Lock()
	defer c.Unlock()
	prev := c.previousLang
	c.previousLang = now

	if c.handler != nil {
		go c.handler.OnLanguageUpdate(prev, deepCopyLanguage(now))
	}
}

func (c *CC) sync() {
	blog.Infof("start sync config from config center.")
	c.syncProc()
	c.syncExtra()
	c.syncMongodb()
	c.syncRedis()
	c.syncLang()
	c.syncErr()
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			default:

			}
			// sync the data from zk, and compare if it has been changed.
			// then call their handler.
			c.syncProc()
			c.syncExtra()
			c.syncMongodb()
			c.syncRedis()
			c.syncLang()
			c.syncErr()
			time.Sleep(15 * time.Second)
		}
	}()
}

func (c *CC) syncProc() {
	blog.V(5).Infof("start sync proc config from config center.")
	procPath := fmt.Sprintf("%s/%s", types.CC_SERVCONF_BASEPATH, types.CCConfigureCommon)
	data, err := c.disc.Read(procPath)
	if err != nil {
		blog.Errorf("sync process config failed, node: %s, err: %v", procPath, err)
		return
	}

	conf := ParseConfigWithData([]byte(data))

	c.Lock()
	if reflect.DeepEqual(conf, c.previousProc) {
		blog.V(5).Infof("sync process config, but nothing is changed.")
		c.Unlock()
		return
	}

	event := &crd.DiscoverEvent{
		Err:  nil,
		Data: []byte(data),
	}

	c.Unlock()
	c.onProcChange(event)

}

func (c *CC) syncExtra() {
	blog.V(5).Infof("start sync proc config from config center.")
	procPath := fmt.Sprintf("%s/%s", types.CC_SERVCONF_BASEPATH, types.CCConfigureExtra)
	data, err := c.disc.Read(procPath)
	if err != nil {
		blog.Errorf("sync *extra* config failed, node: %s, err: %v", procPath, err)
		return
	}

	conf := ParseConfigWithData([]byte(data))

	c.Lock()
	if reflect.DeepEqual(conf, c.previousExtra) {
		blog.V(5).Infof("sync *extra* config, but nothing is changed.")
		c.Unlock()
		return
	}

	event := &crd.DiscoverEvent{
		Err:  nil,
		Data: []byte(data),
	}

	c.Unlock()

	c.onExtraChange(event)

}

func (c *CC) syncMongodb() {
	blog.V(5).Infof("start sync mongo config from config center.")
	mongoPath := fmt.Sprintf("%s/%s", types.CC_SERVCONF_BASEPATH, types.CCConfigureMongo)
	data, err := c.disc.Read(mongoPath)
	if err != nil {
		blog.Errorf("sync *mongo* config failed, node: %s, err: %v", mongoPath, err)
		return
	}

	conf := ParseConfigWithData([]byte(data))
	c.Lock()
	if reflect.DeepEqual(conf, c.previousMongodb) {
		blog.V(5).Infof("sync *mongo* config, but nothing is changed.")
		c.Unlock()
		return
	}
	event := &crd.DiscoverEvent{
		Err:  nil,
		Data: []byte(data),
	}

	c.Unlock()

	c.onMongodbChange(event)

}

func (c *CC) syncRedis() {
	blog.V(5).Infof("start sync redis config from config center.")
	redisPath := fmt.Sprintf("%s/%s", types.CC_SERVCONF_BASEPATH, types.CCConfigureRedis)
	data, err := c.disc.Read(redisPath)
	if err != nil {
		blog.Errorf("sync *redis* config failed, node: %s, err: %v", redisPath, err)
		return
	}

	conf := ParseConfigWithData([]byte(data))

	c.Lock()
	if reflect.DeepEqual(conf, c.previousRedis) {
		blog.V(5).Infof("sync *redis* config, but nothing is changed.")
		c.Unlock()
		return
	}

	event := &crd.DiscoverEvent{
		Err:  nil,
		Data: []byte(data),
	}

	c.Unlock()

	c.onRedisChange(event)

}

func (c *CC) syncLang() {
	blog.V(5).Infof("start sync lang config from config center.")
	data, err := c.disc.Read(types.CC_SERVLANG_BASEPATH)
	if err != nil {
		blog.Errorf("sync process config failed, node: %s, err: %v", types.CC_SERVLANG_BASEPATH, err)
		return
	}

	lang := make(map[string]language.LanguageMap)
	if err := json.Unmarshal([]byte(data), &lang); err != nil {
		blog.Errorf("sync *LANGUAGE* config, but unmarshal failed, err: %v", err)
		return
	}

	c.Lock()

	if reflect.DeepEqual(lang, c.previousLang) {
		blog.V(5).Infof("sync language config, but nothing is changed.")
		c.Unlock()
		return
	}

	event := &crd.DiscoverEvent{
		Err:  nil,
		Data: []byte(data),
	}
	c.Unlock()
	c.onLanguageChange(event)
}

func (c *CC) syncErr() {
	blog.V(5).Infof("start sync error config from config center.")
	data, err := c.disc.Read(types.CC_SERVERROR_BASEPATH)
	if err != nil {
		blog.Errorf("sync process config failed, node: %s, err: %v", types.CC_SERVERROR_BASEPATH, err)
		return
	}

	errCode := make(map[string]errors.ErrorCode)
	if err := json.Unmarshal([]byte(data), &errCode); err != nil {
		blog.Errorf("sync error code config, but unmarshal failed, err: %v", err)
		return
	}

	c.Lock()
	if reflect.DeepEqual(errCode, c.previousError) {
		blog.V(5).Infof("sync error code config, but nothing is changed.")
		c.Unlock()
		return
	}

	event := &crd.DiscoverEvent{
		Err:  nil,
		Data: []byte(data),
	}
	c.Unlock()
	c.onErrorChange(event)
}

func deepCopyError(source map[string]errors.ErrorCode) map[string]errors.ErrorCode {
	copy := make(map[string]errors.ErrorCode)
	if source == nil {
		return copy
	}

	for k, v := range source {
		copy[k] = v
	}

	return copy
}

func deepCopyLanguage(source map[string]language.LanguageMap) map[string]language.LanguageMap {
	copy := make(map[string]language.LanguageMap)
	if source == nil {
		return copy
	}

	for k, v := range source {
		copy[k] = v
	}

	return copy
}
