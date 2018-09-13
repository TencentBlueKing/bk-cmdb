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

func NewConfigCenter(ctx context.Context, zkAddr string, procName string, confPath string, handler *CCHandler) error {
	disc := crd.NewZkRegDiscover(zkAddr, 10*time.Second)
	return New(ctx, procName, confPath, disc, handler)
}

func New(ctx context.Context, procName string, confPath string, disc crd.ConfRegDiscvIf, handler *CCHandler) error {
	confC = &CC{
		ctx:           ctx,
		disc:          disc,
		handler:       handler,
		procName:      procName,
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
	OnLanguageUpdate func(previous, current map[string]language.LanguageMap)
	OnErrorUpdate    func(previous, current map[string]errors.ErrorCode)
}

type CC struct {
	sync.Mutex
	// used to stop the config center gracefully.
	ctx           context.Context
	disc          crd.ConfRegDiscvIf
	handler       *CCHandler
	procName      string
	previousProc  *ProcessConfig
	previousLang  map[string]language.LanguageMap
	previousError map[string]errors.ErrorCode
}

func (c *CC) run() error {
	if err := c.disc.Start(); err != nil {
		return fmt.Errorf("start discover config center failed, err: %v", err)
	}

	procPath := fmt.Sprintf("%s/%s", types.CC_SERVCONF_BASEPATH, c.procName)
	procEvent, err := c.disc.Discover(procPath)
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
		case pEvent := <-procEvent:
			c.onProcChange(pEvent)
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
	blog.V(5).Infof("config center received event that *%s* config has changed. event: %s", c.procName, string(cur.Data))

	if cur.Err != nil {
		blog.Errorf("config center received event that %s config has changed, but got err: %v", c.procName, cur.Err)
		return
	}

	now, err := ParseConfigWithData(cur.Data)
	if err != nil {
		blog.Errorf("config center received event that *%s* config has changed, but parse failed, err: %v", c.procName, err)
		return
	}

	c.Lock()
	defer c.Unlock()
	prev := c.previousProc
	c.previousProc = now
	if c.handler != nil {
		go c.handler.OnProcessUpdate(*prev, *now)
	}
	blog.V(5).Infof("config center received event that *%s* config has changed. prev: %v, cur: %v", c.procName, *prev, *now)
}

func (c *CC) onErrorChange(cur *crd.DiscoverEvent) {
	blog.V(5).Infof("config center received event that *ERROR CODE* config has changed. event: %s", string(cur.Data))

	if cur.Err != nil {
		blog.Errorf("config center received event that *ERROR CODE* config has changed, but got err: %v", cur.Err)
		return
	}

	now := make(map[string]errors.ErrorCode)
	if err := json.Unmarshal(cur.Data, &now); err != nil {
		blog.Errorf("config center received %s event that *ERROR CODE* config has changed, but unmarshal err: %v", c.procName, err)
		return
	}

	c.Lock()
	defer c.Unlock()
	prev := c.previousError
	c.previousError = now

	if c.handler != nil {
		go c.handler.OnErrorUpdate(prev, deepCopyError(now))
	}
	blog.V(5).Infof("config center received event that *ERROR CODE* config has changed. prev: %v, cur: %v", prev, now)
}

func (c *CC) onLanguageChange(cur *crd.DiscoverEvent) {
	blog.V(5).Infof("config center received event that *LANGUAGE* config has changed. event: %s", string(cur.Data))

	if cur.Err != nil {
		blog.Errorf("config center received event that *LANGUAGE* config has changed, but got err: %v", cur.Err)
		return
	}

	now := make(map[string]language.LanguageMap)
	if err := json.Unmarshal(cur.Data, &now); err != nil {
		blog.Errorf("config (%s) center received event that *LANGUAGE* config has changed, but unmarshal err: %v", c.procName, err)
		return
	}

	c.Lock()
	defer c.Unlock()
	prev := c.previousLang
	c.previousLang = now

	if c.handler != nil {
		go c.handler.OnLanguageUpdate(prev, deepCopyLanguage(now))
	}
	blog.V(5).Infof("config center received event that *LANGUAGE* config has changed. prev: %v, cur: %v", prev, now)
}

func (c *CC) sync() {
	blog.Infof("start sync config from config center.")
	c.syncProc()
	c.syncLang()
	c.syncErr()
	ticker := time.NewTicker(15 * time.Second)
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			case <-ticker.C:
			}

			// sync the data from zk, and compare if it has been changed.
			// then call their handler.
			c.syncProc()
			c.syncLang()
			c.syncErr()
		}
	}()
}

func (c *CC) syncProc() {
	blog.V(5).Infof("start sync proc config from config center.")
	procPath := fmt.Sprintf("%s/%s", types.CC_SERVCONF_BASEPATH, c.procName)
	data, err := c.disc.Read(procPath)
	if err != nil {
		blog.Errorf("sync process config failed, err: %v", err)
		return
	}

	conf, err := ParseConfigWithData([]byte(data))
	if err != nil {
		blog.Errorf("config center sync process[%s] config, but parse failed, err: %v", c.procName, err)
		return
	}

	c.Lock()
	if reflect.DeepEqual(conf, c.previousProc) {
		blog.V(4).Infof("sync process config, but nothing is changed.")
		c.Unlock()
		return
	}
	blog.V(4).Infof("sync process[%s] config, before change is: %+#v", c.procName, *(c.previousProc))
	blog.V(4).Infof("sync process[%s] config, after change is: %+#v", c.procName, *conf)

	event := &crd.DiscoverEvent{
		Err:  nil,
		Data: []byte(data),
	}

	c.Unlock()
	c.onProcChange(event)

}

func (c *CC) syncLang() {
	blog.V(5).Infof("start sync lang config from config center.")
	data, err := c.disc.Read(types.CC_SERVLANG_BASEPATH)
	if err != nil {
		blog.Errorf("sync process config failed, err: %v", err)
		return
	}

	lang := make(map[string]language.LanguageMap)
	if err := json.Unmarshal([]byte(data), &lang); err != nil {
		blog.Errorf("sync %s *LANGUAGE* config, but unmarshal failed, err: %v", c.procName, err)
		return
	}

	c.Lock()

	if reflect.DeepEqual(lang, c.previousLang) {
		blog.V(5).Infof("sync language config, but nothing is changed.")
		c.Unlock()
		return
	}

	blog.V(5).Infof("sync language config, before change is: %v", c.previousLang)
	blog.V(5).Infof("sync language config, after change is: %v", lang)

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
		blog.Errorf("sync process config failed, err: %v", err)
		return
	}

	errCode := make(map[string]errors.ErrorCode)
	if err := json.Unmarshal([]byte(data), &errCode); err != nil {
		blog.Errorf("sync %s error code config, but unmarshal failed, err: %v", c.procName, err)
		return
	}

	c.Lock()
	if reflect.DeepEqual(errCode, c.previousError) {
		blog.V(5).Infof("sync error code config, but nothing is changed.")
		c.Unlock()
		return
	}

	blog.V(5).Infof("sync language config, before change is: %v", c.previousError)
	blog.V(5).Infof("sync language config, after change is: %v", errCode)

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
