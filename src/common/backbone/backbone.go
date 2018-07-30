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

package backbone

import (
	"context"
	"fmt"
	"sync"

	"configcenter/src/apimachinery"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
)

func NewBackbone(ctx context.Context, zkAddr string, procName string, confPath string, procHandler cc.ProcHandlerFunc, c *Config) (*Engine, error) {
	disc, err := NewServcieDiscovery(zkAddr)
	if err != nil {
		return nil, fmt.Errorf("new service discover failed, err:%v", err)
	}

	engine, err := New(c, disc)
	if err != nil {
		return nil, fmt.Errorf("new engine failed, err: %v", err)
	}

	handler := &cc.CCHandler{
		OnProcessUpdate:  procHandler,
		OnLanguageUpdate: engine.onLanguageUpdate,
		OnErrorUpdate:    engine.onErrorUpdate,
	}

	err = cc.NewConfigCenter(ctx, zkAddr, procName, confPath, handler)
	if err != nil {
		return nil, fmt.Errorf("new config center failed, err: %v", err)
	}

	if err := ListenServer(c.Server); err != nil {
		return nil, err
	}

	return engine, nil
}

func New(c *Config, disc ServiceDiscoverInterface) (*Engine, error) {
	if err := disc.Register(c.RegisterPath, c.RegisterInfo); err != nil {
		return nil, err
	}

	return &Engine{
		CoreAPI:  c.CoreAPI,
		SvcDisc:  disc,
		Language: language.NewFromCtx(language.EmptyLanguageSetting),
		CCErr:    errors.NewFromCtx(errors.EmptyErrorsSetting),
	}, nil
}

type Engine struct {
	sync.Mutex
	CoreAPI  apimachinery.ClientSetInterface
	SvcDisc  ServiceDiscoverInterface
	Language language.CCLanguageIf
	CCErr    errors.CCErrorIf
}

func (e *Engine) onLanguageUpdate(previous, current map[string]language.LanguageMap) {
	e.Lock()
	defer e.Unlock()
	if e.Language == nil {
		e.Language = language.NewFromCtx(current)
		blog.Infof("load language config success.")
		return
	}
	e.Language.Load(current)
	blog.V(3).Infof("load new language config success.")
}

func (e *Engine) onErrorUpdate(previous, current map[string]errors.ErrorCode) {
	e.Lock()
	defer e.Unlock()
	if e.CCErr == nil {
		e.CCErr = errors.NewFromCtx(current)
		blog.Infof("load error code config success.")
		return
	}
	e.CCErr.Load(current)
	blog.V(3).Infof("load new error config success.")
}

func (e *Engine) Ping() error {
	return e.SvcDisc.Ping()
}
