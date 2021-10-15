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
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/metrics"
	"configcenter/src/common/registerdiscover"
	"configcenter/src/common/types"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/thirdparty/monitor"

	"github.com/rs/xid"
)

// BackboneParameter Used to constrain different services to ensure
// consistency of service startup capabilities
type BackboneParameter struct {
	// ConfigUpdate handle process config change
	ConfigUpdate cc.ProcHandlerFunc
	ExtraUpdate  cc.ProcHandlerFunc

	// service component addr
	Regdiscv string
	// config path
	ConfigPath string
	// http server parameter
	SrvInfo *types.ServerInfo
}

// NewBackbone create a backbone object
func NewBackbone(ctx context.Context, input *BackboneParameter) (*Engine, error) {
	// validate backbone config
	if err := validateParameter(input); err != nil {
		return nil, err
	}

	// init server info
	common.SetServerInfo(input.SrvInfo)

	// new metrics service
	metricService := metrics.NewService(
		metrics.Config{
			ProcessName:     common.GetIdentification(),
			ProcessInstance: input.SrvInfo.Instance(),
		})

	// new register and discover base
	rd, err := newRegDiscv(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("new register and discover (%s) failed, err: %v", input.Regdiscv, err)
	}

	// new service discovery for api machinery
	serviceDiscovery, err := discovery.NewServiceDiscovery(rd)
	if err != nil {
		return nil, fmt.Errorf("new service discovery failed, err: %v", err)
	}

	// new api machinery
	apiMachineryConfig := &util.APIMachineryConfig{
		QPS:       1000,
		Burst:     2000,
		TLSConfig: nil,
	}
	apiMachinery, err := apimachinery.NewApiMachinery(apiMachineryConfig, serviceDiscovery)
	if err != nil {
		return nil, fmt.Errorf("new api machinery failed, err: %v", err)
	}

	// new service register for backbone engine
	serviceRegister, err := NewServiceRegister(rd)
	if err != nil {
		return nil, fmt.Errorf("new service discover failed, err: %v", err)
	}

	// new backbone engine
	engine, err := newEngine(input.SrvInfo, apiMachinery, serviceRegister)
	if err != nil {
		return nil, fmt.Errorf("new engine failed, err: %v", err)
	}
	engine.regdiscv = rd
	engine.apiMachineryConfig = apiMachineryConfig
	engine.discovery = serviceDiscovery
	engine.ServiceManageInterface = serviceDiscovery
	engine.srvInfo = input.SrvInfo
	engine.metric = metricService

	// new config center
	handler := &cc.CCHandler{
		// 扩展这个函数， 新加传递错误
		OnProcessUpdate:  input.ConfigUpdate,
		OnExtraUpdate:    input.ExtraUpdate,
		OnLanguageUpdate: engine.onLanguageUpdate,
		OnErrorUpdate:    engine.onErrorUpdate,
		OnMongodbUpdate:  engine.onMongodbUpdate,
		OnRedisUpdate:    engine.onRedisUpdate,
	}
	err = cc.NewConfigCenter(ctx, rd, input.ConfigPath, handler)
	if err != nil {
		return nil, fmt.Errorf("new config center failed, err: %v", err)
	}

	// start discover event handler
	err = handleNotice(ctx, rd, input.SrvInfo.Instance())
	if err != nil {
		return nil, fmt.Errorf("handle notice failed, err: %v", err)
	}

	// init monitor
	if err := monitor.InitMonitor(); err != nil {
		return nil, fmt.Errorf("init monitor failed, err: %v", err)
	}

	return engine, nil
}

// StartServer start http server and register to register and discover
func StartServer(ctx context.Context, cancel context.CancelFunc, e *Engine, httpHandler http.Handler,
	pprofEnabled bool) error {
	e.server = Server{
		ListenAddr:   e.srvInfo.IP,
		ListenPort:   e.srvInfo.Port,
		Handler:      e.Metric().HTTPMiddleware(httpHandler),
		TLS:          TLSConfig{},
		PProfEnabled: pprofEnabled,
	}

	if err := ListenAndServe(e.server, e.register, cancel); err != nil {
		return err
	}

	// wait for a while to see if ListenAndServe in goroutine is successful
	// to avoid registering an invalid server address
	time.Sleep(time.Second)

	return e.register.Register(e.RegisterPath, *e.srvInfo)
}

func newRegDiscv(ctx context.Context, input *BackboneParameter) (*registerdiscover.RegDiscv, error) {
	regdiscvConf := &registerdiscover.Config{
		Host: input.Regdiscv,
		// TODO: get User, Passwd, TLS from flag
		TLS: nil,
	}
	return registerdiscover.NewRegDiscv(regdiscvConf)
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
	if input.ConfigUpdate == nil && input.ExtraUpdate == nil {
		return fmt.Errorf("service config change funcation can not be emtpy")
	}
	// to prevent other components which doesn't set it from failing
	if input.SrvInfo.RegisterIP == "" {
		input.SrvInfo.RegisterIP = input.SrvInfo.IP
	}
	if input.SrvInfo.UUID == "" {
		input.SrvInfo.UUID = xid.New().String()
	}
	return nil
}

func newEngine(srvInfo *types.ServerInfo, coreApi apimachinery.ClientSetInterface, reg ServiceRegisterInterface) (
	*Engine, error) {

	regPath := fmt.Sprintf("%s/%s/%s", types.CCDiscoverBaseEndpoint, common.GetIdentification(), srvInfo.IP)

	return &Engine{
		RegisterPath: regPath,
		CoreAPI:      coreApi,
		register:     reg,
		Language:     language.NewFromCtx(language.EmptyLanguageSetting),
		CCErr:        errors.NewFromCtx(errors.EmptyErrorsSetting),
		CCCtx:        newCCContext(),
	}, nil
}

type Engine struct {
	CoreAPI                apimachinery.ClientSetInterface
	apiMachineryConfig     *util.APIMachineryConfig
	regdiscv               *registerdiscover.RegDiscv
	register               ServiceRegisterInterface
	discovery              discovery.DiscoveryInterface
	ServiceManageInterface discovery.ServiceManageInterface
	metric                 *metrics.Service

	sync.Mutex

	RegisterPath string
	server       Server
	srvInfo      *types.ServerInfo

	Language language.CCLanguageIf
	CCErr    errors.CCErrorIf
	CCCtx    CCContextInterface
}

// Discovery return service discovery interface
func (e *Engine) Discovery() discovery.DiscoveryInterface {
	return e.discovery
}

// ApiMachineryConfig return api machinery config
func (e *Engine) ApiMachineryConfig() *util.APIMachineryConfig {
	return e.apiMachineryConfig
}

// RegDiscv return register and discover
func (e *Engine) RegDiscv() *registerdiscover.RegDiscv {
	return e.regdiscv
}

// Metric return metrics service
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

func (e *Engine) onMongodbUpdate(previous, current cc.ProcessConfig) {
	e.Lock()
	defer e.Unlock()
	if err := cc.SetMongodbFromByte(current.ConfigData); err != nil {
		blog.Errorf("parse mongo config failed, err: %s, data: %s", err.Error(), string(current.ConfigData))
	}
}

func (e *Engine) onRedisUpdate(previous, current cc.ProcessConfig) {
	e.Lock()
	defer e.Unlock()
	if err := cc.SetRedisFromByte(current.ConfigData); err != nil {
		blog.Errorf("parse redis config failed, err: %s, data: %s", err.Error(), string(current.ConfigData))
	}
}

// Ping verify register and discover accessibility
func (e *Engine) Ping() error {
	return e.register.Ping()
}

// WithRedis return redis config
func (e *Engine) WithRedis(prefixes ...string) (redis.Config, error) {
	// use default prefix if no prefix is specified, or use the first prefix
	var prefix string
	if len(prefixes) == 0 {
		prefix = "redis"
	} else {
		prefix = prefixes[0]
	}

	return cc.Redis(prefix)
}

// WithMongo return mongo config
func (e *Engine) WithMongo(prefixes ...string) (mongo.Config, error) {
	var prefix string
	if len(prefixes) == 0 {
		prefix = "mongodb"
	} else {
		prefix = prefixes[0]
	}

	return cc.Mongo(prefix)
}
