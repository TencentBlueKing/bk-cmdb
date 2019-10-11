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
	"configcenter/src/common/metrics"
	"configcenter/src/common/types"
)

// connect svcManager retry connect time
const maxRetry = 200

// BackboneParameter Used to constrain different services to ensure
// consistency of service startup capabilities
type BackboneParameter struct {
	// ConfigUpdate handle process config change
	ConfigUpdate cc.ProcHandlerFunc

	// service component addr
	Regdiscv string
	// config path
	ConfigPath string
	// http server parameter
	SrvInfo *types.ServerInfo
}

func newSvcManagerClient(ctx context.Context, svcManagerAddr string) (*zk.ZkClient, error) {
	var err error
	for retry := 0; retry < maxRetry; retry++ {
		client := zk.NewZkClient(svcManagerAddr, 5*time.Second)
		if err = client.Start(); err != nil {
			blog.Errorf("connect regdiscv [%s] failed: %v", svcManagerAddr, err)
			time.Sleep(time.Second * 2)
			continue
		}

		if err = client.Ping(); err != nil {
			blog.Errorf("connect regdiscv [%s] failed: %v", svcManagerAddr, err)
			time.Sleep(time.Second * 2)
			continue
		}

		return client, nil
	}

	return nil, err
}

func newConfig(ctx context.Context, srvInfo *types.ServerInfo, discovery discovery.DiscoveryInterface, apiMachineryConfig *util.APIMachineryConfig) (*Config, error) {

	machinery, err := apimachinery.NewApiMachinery(apiMachineryConfig, discovery)
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

func validateParameter(input *BackboneParameter) error {
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
	if err := validateParameter(input); err != nil {
		return nil, err
	}

	metricService := metrics.NewService(metrics.Config{ProcessName: common.GetIdentification(), ProcessInstance: input.SrvInfo.Instance()})

	common.SetServerInfo(input.SrvInfo)
	client, err := newSvcManagerClient(ctx, input.Regdiscv)
	if err != nil {
		return nil, fmt.Errorf("connect regdiscv [%s] failed: %v", input.Regdiscv, err)
	}
	serviceDiscovery, err := discovery.NewServiceDiscovery(client)
	if err != nil {
		return nil, fmt.Errorf("connect regdiscv [%s] failed: %v", input.Regdiscv, err)
	}
	disc, err := NewServiceRegister(client)
	if err != nil {
		return nil, fmt.Errorf("new service discover failed, err:%v", err)
	}

	apiMachineryConfig := &util.APIMachineryConfig{
		QPS:       1000,
		Burst:     2000,
		TLSConfig: nil,
	}
	c, err := newConfig(ctx, input.SrvInfo, serviceDiscovery, apiMachineryConfig)
	if err != nil {
		return nil, err
	}
	engine, err := New(c, disc)
	if err != nil {
		return nil, fmt.Errorf("new engine failed, err: %v", err)
	}
	engine.client = client
	engine.apiMachineryConfig = apiMachineryConfig
	engine.discovery = serviceDiscovery
	engine.ServiceManageInterface = serviceDiscovery
	engine.srvInfo = input.SrvInfo
	engine.metric = metricService

	handler := &cc.CCHandler{
		OnProcessUpdate:  input.ConfigUpdate,
		OnLanguageUpdate: engine.onLanguageUpdate,
		OnErrorUpdate:    engine.onErrorUpdate,
	}

	err = cc.NewConfigCenter(ctx, client, common.GetIdentification(), input.ConfigPath, handler)
	if err != nil {
		return nil, fmt.Errorf("new config center failed, err: %v", err)
	}

	err = handleNotice(ctx, client.Client(), input.SrvInfo.Instance())
	if err != nil {
		return nil, fmt.Errorf("handle notice failed, err: %v", err)
	}

	return engine, nil
}

func StartServer(ctx context.Context, cancel context.CancelFunc, e *Engine, HTTPHandler http.Handler, pprofEnabled bool) error {
	e.server = Server{
		ListenAddr:   e.srvInfo.IP,
		ListenPort:   e.srvInfo.Port,
		Handler:      e.Metric().HTTPMiddleware(HTTPHandler),
		TLS:          TLSConfig{},
		PProfEnabled: pprofEnabled,
	}

	return ListenAndServe(e.server, e.SvcDisc, cancel)
}

func New(c *Config, disc ServiceRegisterInterface) (*Engine, error) {
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
	CoreAPI            apimachinery.ClientSetInterface
	apiMachineryConfig *util.APIMachineryConfig

	client                 *zk.ZkClient
	ServiceManageInterface discovery.ServiceManageInterface
	SvcDisc                ServiceRegisterInterface
	discovery              discovery.DiscoveryInterface
	metric                 *metrics.Service

	sync.Mutex

	ServerInfo types.ServerInfo
	server     Server
	srvInfo    *types.ServerInfo

	Language language.CCLanguageIf
	CCErr    errors.CCErrorIf
	CCCtx    CCContextInterface
}

func (e *Engine) Discovery() discovery.DiscoveryInterface {
	return e.discovery
}

func (e *Engine) ApiMachineryConfig() *util.APIMachineryConfig {
	return e.apiMachineryConfig
}

func (e *Engine) ServiceManageClient() *zk.ZkClient {
	return e.client
}

func (e *Engine) Metric() *metrics.Service {
	return e.metric
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
