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
	"configcenter/src/auth/authcenter"
	"configcenter/src/common"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/backbone/service_mange/zk"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/metrics"
	"configcenter/src/common/types"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/redis"

	"github.com/rs/xid"
)

// connect svcManager retry connect time
const maxRetry = 200

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
		OnExtraUpdate:    input.ExtraUpdate,
		OnLanguageUpdate: engine.onLanguageUpdate,
		OnErrorUpdate:    engine.onErrorUpdate,
	}

	err = cc.NewConfigCenter(ctx, client, input.ConfigPath, handler)
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

var (
	redisConf = make(map[string]redis.Config, 0)
)

func (e *Engine) WithRedis(prefixes ...string) (redis.Config, error) {
	// use default prefix if no prefix is specified, or use the first prefix
	var prefix string
	if len(prefixes) == 0 {
		prefix = "redis"
	} else {
		prefix = prefixes[0]
	}
	// only initialize once, not allow hot update
	if conf, exist := redisConf[prefix]; exist {
		return conf, nil
	}
	data, err := e.client.Client().Get(fmt.Sprintf("%s/%s", types.CC_SERVCONF_BASEPATH, types.CCConfigureRedis))
	if err != nil {
		blog.Errorf("get redis config failed, err: %s", err.Error())
		return redis.Config{}, err
	}
	conf, err := cc.ParseConfigWithData([]byte(data))
	if err != nil {
		blog.Errorf("parse redis config failed, err: %s, data: %s", err.Error(), data)
		return redis.Config{}, err
	}
	redisConf[prefix] = redis.ParseConfigFromKV(prefix, conf.ConfigMap)
	return redisConf[prefix], nil
}

var (
	mongoConf = make(map[string]mongo.Config, 0)
)

func (e *Engine) WithMongo(prefixes ...string) (mongo.Config, error) {
	var prefix string
	if len(prefixes) == 0 {
		prefix = "mongodb"
	} else {
		prefix = prefixes[0]
	}
	if conf, exist := mongoConf[prefix]; exist {
		return conf, nil
	}
	data, err := e.client.Client().Get(fmt.Sprintf("%s/%s", types.CC_SERVCONF_BASEPATH, types.CCConfigureMongo))
	if err != nil {
		blog.Errorf("get mongo config failed, err: %s", err.Error())
		return mongo.Config{}, err
	}
	conf, err := cc.ParseConfigWithData([]byte(data))
	if err != nil {
		blog.Errorf("parse mongo config failed, err: %s, data: %s", err.Error(), data)
		return mongo.Config{}, err
	}
	mongoConf[prefix] = mongo.ParseConfigFromKV(prefix, conf.ConfigMap)
	return mongoConf[prefix], nil
}

var (
	authConf = make(map[string]authcenter.AuthConfig, 0)
)

func (e *Engine) WithAuth(prefixes ...string) (authcenter.AuthConfig, error) {
	var prefix string
	if len(prefixes) == 0 {
		prefix = "auth"
	} else {
		prefix = prefixes[0]
	}
	if conf, exist := authConf[prefix]; exist {
		return conf, nil
	}
	data, err := e.client.Client().Get(fmt.Sprintf("%s/%s", types.CC_SERVCONF_BASEPATH, types.CCConfigureCommon))
	if err != nil {
		blog.Errorf("get common config failed, err: %s", err.Error())
		return authcenter.AuthConfig{}, err
	}
	conf, err := cc.ParseConfigWithData([]byte(data))
	if err != nil {
		blog.Errorf("parse common config failed, err: %s, data: %s", err.Error(), data)
		return authcenter.AuthConfig{}, err
	}
	authConf[prefix], err = authcenter.ParseConfigFromKV(prefix, conf.ConfigMap)
	if err != nil {
		blog.Errorf("parse auth center config failed: %s", err.Error())
		return authcenter.AuthConfig{}, err
	}
	return authConf[prefix], nil
}
