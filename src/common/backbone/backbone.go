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
	"net/http"
	"sync"
	"time"

	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/backbone/service_mange/zk"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/types"
)

// BackboneParameter Used to constrain different services to ensure
// consistency of service startup capabilities
type BackboneParameter struct {
	// ConfigUpdate handle process config change
	ConfigUpdate cc.ProcHandlerFunc

	// servvice component addr
	Regdiscv string
	// config path
	ConfigPath string
	// http server paramster
	SrvInfo *types.ServerInfo
}

func newManageSrvClient(ctx context.Context, manageSrvAddr string) (*zk.ZkClient, error) {
	client := zk.NewZkClient(manageSrvAddr, 5*time.Second)
	if err := client.Start(); err != nil {
		return nil, fmt.Errorf("connect regdiscv [%s] failed: %v", manageSrvAddr, err)
	}
	if err := client.Ping(); err != nil {
		return nil, fmt.Errorf("connect regdiscv [%s] failed: %v", manageSrvAddr, err)
	}

	return client, nil
}

func newConfig(ctx context.Context, srvInfo *types.ServerInfo, discovery discovery.DiscoveryInterface, apiMachinerConfig *util.APIMachineryConfig) (*Config, error) {

	machinery, err := apimachinery.NewApiMachinery(apiMachinerConfig, discovery)
	if err != nil {
		return nil, fmt.Errorf("new api machinery failed, err: %v", err)
	}
	regPath := fmt.Sprintf("%s/%s/%s", types.CC_SERV_BASEPATH, common.GetIdentification(), srvInfo.IP)

	bonC := &Config{
		RegisterPath: regPath,
		RegisterInfo: *srvInfo,
		CoreAPI:      machinery,
	}

	return bonC, nil
}

func parameterValid(input *BackboneParameter) error {
	if input.Regdiscv == "" {
		return fmt.Errorf("regdiscv can not be emtpy")
	}
	if input.SrvInfo.IP == "" {
		return fmt.Errorf("addrport ip can not be emtpy")
	}
	if input.SrvInfo.Port <= 0 || input.SrvInfo.Port > 65535 {
		return fmt.Errorf("addrport port must be 1-65535")
	}

	if input.ConfigUpdate == nil {
		return fmt.Errorf("service config change funcation can not be emtpy")
	}

	return nil
}

func NewBackbone(ctx context.Context, input *BackboneParameter) (*Engine, error) {
	if err := parameterValid(input); err != nil {
		return nil, err
	}
	common.SetServerInfo(input.SrvInfo)
	client, err := newManageSrvClient(ctx, input.Regdiscv)
	if err != nil {
		return nil, fmt.Errorf("connect regdiscv [%s] failed: %v", input.Regdiscv, err)
	}
	discoveryInterface, err := discovery.NewDiscoveryInterface(client)
	if err != nil {
		return nil, fmt.Errorf("connect regdiscv [%s] failed: %v", input.Regdiscv, err)
	}
	disc, err := NewServcieDiscovery(client)
	if err != nil {
		return nil, fmt.Errorf("new service discover failed, err:%v", err)
	}

	apiMachineryConfig := &util.APIMachineryConfig{
		QPS:       1000,
		Burst:     2000,
		TLSConfig: nil,
	}
	c, err := newConfig(ctx, input.SrvInfo, discoveryInterface, apiMachineryConfig)
	if err != nil {
		return nil, err
	}
	engine, err := New(c, disc)
	if err != nil {
		return nil, fmt.Errorf("new engine failed, err: %v", err)
	}
	engine.client = client
	engine.apiMachinerConfig = apiMachineryConfig
	engine.discovery = discoveryInterface
	engine.ServiceManageInterface = discoveryInterface
	engine.srvInfo = input.SrvInfo

	handler := &cc.CCHandler{
		OnProcessUpdate:  input.ConfigUpdate,
		OnLanguageUpdate: engine.onLanguageUpdate,
		OnErrorUpdate:    engine.onErrorUpdate,
	}

	err = cc.NewConfigCenter(ctx, client, common.GetIdentification(), input.ConfigPath, handler)
	if err != nil {
		return nil, fmt.Errorf("new config center failed, err: %v", err)
	}

	return engine, nil
}

func StartServer(ctx context.Context, e *Engine, HTTPHandler http.Handler) error {
	e.server = Server{
		ListenAddr: e.srvInfo.IP,
		ListenPort: e.srvInfo.Port,
		Handler:    HTTPHandler,
		TLS:        TLSConfig{},
	}

	if err := ListenServer(e.server); err != nil {
		return err
	}
	return nil
}

func New(c *Config, disc ServiceDiscoverInterface) (*Engine, error) {
	if err := disc.Register(c.RegisterPath, c.RegisterInfo); err != nil {
		return nil, err
	}

	return &Engine{
		ServerInfo: c.RegisterInfo,
		CoreAPI:    c.CoreAPI,
		SvcDisc:    disc,
		Language:   language.NewFromCtx(language.EmptyLanguageSetting),
		CCErr:      errors.NewFromCtx(errors.EmptyErrorsSetting),
		CCCtx:      newCCContext(),
	}, nil
}

type Engine struct {
	client *zk.ZkClient
	sync.Mutex
	ServerInfo             types.ServerInfo
	CoreAPI                apimachinery.ClientSetInterface
	SvcDisc                ServiceDiscoverInterface
	Language               language.CCLanguageIf
	CCErr                  errors.CCErrorIf
	CCCtx                  CCContextInterface
	ServiceManageInterface discovery.ServiceManageInterface
	apiMachinerConfig      *util.APIMachineryConfig
	discovery              discovery.DiscoveryInterface
	server                 Server
	srvInfo                *types.ServerInfo
}

func (e *Engine) Discovery() discovery.DiscoveryInterface {
	return e.discovery
}

func (e *Engine) ApiMachineryConfig() *util.APIMachineryConfig {
	return e.apiMachinerConfig
}

func (e *Engine) ServiceManageClient() *zk.ZkClient {
	return e.client
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
