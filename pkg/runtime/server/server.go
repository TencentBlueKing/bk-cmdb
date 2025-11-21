/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

// Package server defines run cmdb server related logics.
package server

import (
	"context"
	"fmt"
	"log/slog"

	clientv3 "go.etcd.io/etcd/client/v3"

	etcdcli "github.com/TencentBlueKing/bk-cmdb/pkg/client/etcd"
	cc "github.com/TencentBlueKing/bk-cmdb/pkg/config-center"
	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
	etcdconf "github.com/TencentBlueKing/bk-cmdb/pkg/config-center/etcd"
	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/options"
	cerr "github.com/TencentBlueKing/bk-cmdb/pkg/errors"
	"github.com/TencentBlueKing/bk-cmdb/pkg/i18n"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	"github.com/TencentBlueKing/bk-cmdb/pkg/metrics"
	sd "github.com/TencentBlueKing/bk-cmdb/pkg/service-discovery"
	"github.com/TencentBlueKing/bk-cmdb/pkg/service-discovery/etcd"
	"github.com/TencentBlueKing/bk-cmdb/pkg/trace"
)

// Options defines the common options to create new server.
type Options struct {
	*options.Options
	ConfOpts *ConfigOptions
}

// LogValue implements slog.LogValuer interface.
func (o *Options) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("log options", o.Log),
		slog.Any("server info", o.Server),
		slog.Any("etcd endpoints", o.Etcd.Endpoints),
		slog.Any("config options", o.ConfOpts),
	)
}

// Validate validates the common server options.
func (o *Options) Validate() error {
	if o == nil {
		return fmt.Errorf("new common server options is nil")
	}

	if err := o.Options.Validate(); err != nil {
		return fmt.Errorf("common server options validate failed: %w", err)
	}

	if err := o.ConfOpts.Validate(o.Server.Name); err != nil {
		return fmt.Errorf("config options validate failed: %w", err)
	}
	return nil
}

// CommonInfo defines a common server info.
type CommonInfo struct {
	// Server is the common server info.
	Server *config.ServerInfo
	// ConfigInfo is the common config info.
	*ConfigInfo
	// Metrics is the common metrics instance.
	Metrics *metrics.Service
}

// NewServerInfoWithSvcDisc creates a new server info with service discovery.
func NewServerInfoWithSvcDisc(ctx context.Context, opts *Options, services []config.ServiceName) (
	*CommonInfo, sd.ServiceDiscovery, error) {

	svrInfo, err := newServerInfo(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	// create etcd service discovery instance
	registryOpt := &etcd.RegistryOption{
		Service: opts.Server,
	}
	discoveryOpt := &etcd.DiscoveryOption{
		Cluster:  opts.Server.Cluster,
		Services: services,
	}
	serviceDiscovery, err := etcd.NewServiceDiscovery(ctx, svrInfo.etcdCli, registryOpt, discoveryOpt)
	if err != nil {
		log.Error(ctx, "new service discovery failed", log.E(err), "reg", registryOpt, "dis", discoveryOpt)
		return nil, nil, err
	}

	return svrInfo.info, serviceDiscovery, nil
}

// NewServerInfoWithRegistry creates a new server info with registry.
func NewServerInfoWithRegistry(ctx context.Context, opts *Options) (*CommonInfo, sd.Registry, error) {
	svrInfo, err := newServerInfo(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	// create etcd service registry instance
	registryOpt := &etcd.RegistryOption{
		Service: opts.Server,
	}
	registry, err := etcd.NewRegistry(ctx, svrInfo.etcdCli, registryOpt)
	if err != nil {
		log.Error(ctx, "new service registry failed", log.E(err), "opt", registryOpt)
		return nil, nil, err
	}

	return svrInfo.info, registry, nil
}

// newServerInfo creates a new common server info.
func newServerInfo(ctx context.Context, opts *Options) (*serverInfo, error) {
	// 设置服务名称
	config.SetServiceName(opts.Server.Name)

	if err := opts.Validate(); err != nil {
		log.Error(ctx, "validate options failed", log.E(err), "opt", opts)
		return nil, err
	}

	// 日志初始化
	handler := log.NewContextualHandler(opts.Log)
	log.SetDefault(handler)

	// create etcd client
	etcdCli, err := etcdcli.New(opts.Etcd)
	if err != nil {
		log.Error(ctx, "new etcd client failed", log.E(err), "conf", opts.Etcd)
		return nil, err
	}

	// create config info
	confInfo, err := newConfigInfo(ctx, opts.ConfOpts, etcdCli)
	if err != nil {
		return nil, err
	}

	// trace初始化
	if err = trace.SetupTrace(ctx, confInfo.traceConf); err != nil {
		log.Error(ctx, "setup trace failed", log.E(err), "opt", confInfo.traceConf)
		return nil, err
	}

	// init i18n manager and error manager
	m, err := i18n.NewI18nManager(ctx, &i18n.Options{})
	if err != nil {
		log.Error(ctx, "init i18n manager failed", log.E(err))
		return nil, err
	}
	i18n.SetDefaultManager(m)

	errorManager := cerr.NewErrorManager()
	cerr.SetDefaultErrorManager(errorManager)

	// create metrics instance
	metricsConf := &metrics.Config{
		ProcessName: opts.Server.Name,
		Host:        opts.Server.IP,
	}
	metricsSvc := metrics.NewService(metricsConf)

	return &serverInfo{
		etcdCli: etcdCli,
		info: &CommonInfo{
			Server:     opts.Server,
			ConfigInfo: confInfo,
			Metrics:    metricsSvc,
		},
	}, nil
}

type serverInfo struct {
	etcdCli *clientv3.Client
	info    *CommonInfo
}

// ConfigOptions is the common config options.
type ConfigOptions struct {
	// NeededConfigs are the needed config types.
	NeededConfigs []cc.ConfigType
	// Directory is the config directory, right now only admin server uses the config directory to read the config.
	Directory string
}

// Validate validates the common server options.
func (o *ConfigOptions) Validate(service config.ServiceName) error {
	if o == nil {
		return fmt.Errorf("config options is nil")
	}

	if len(o.NeededConfigs) == 0 {
		return fmt.Errorf("needed configs is empty")
	}

	if service == config.AdminServer && len(o.Directory) == 0 {
		return fmt.Errorf("admin server config directory is empty")
	}
	return nil
}

// ConfigInfo is the common config info.
type ConfigInfo struct {
	// CC is the config center.
	CC cc.ConfigCenter
	// TLSConf is the tls config.
	TLSConf *config.TLSConfig
	// traceConf is the trace config.
	traceConf *trace.Option
}

// newConfigInfo creates a new config info.
func newConfigInfo(ctx context.Context, opts *ConfigOptions, etcdCli *clientv3.Client) (*ConfigInfo, error) {
	var configCenter cc.ConfigCenter

	if config.GetServiceName() == config.AdminServer {
		// create config registry instance and registers config
		confReg, err := etcdconf.NewRegistry(etcdCli)
		if err != nil {
			log.Error(ctx, "new config registry failed", log.E(err))
			return nil, err
		}

		writer := cc.NewRegistryWriter(confReg, opts.Directory, opts.NeededConfigs)
		if err = writer.RunConfigWrite(ctx); err != nil {
			log.Error(ctx, "run config write failed", log.E(err))
			return nil, err
		}
		configCenter = writer
	} else {
		// create config discovery instance and reads config
		confDisc, err := etcdconf.NewDiscovery(etcdCli)
		if err != nil {
			log.Error(ctx, "new config discovery failed", log.E(err))
			return nil, err
		}

		reader := cc.NewReader(confDisc, opts.NeededConfigs)
		if err = reader.RunConfigRead(ctx); err != nil {
			log.Error(ctx, "run config read failed", log.E(err))
			return nil, err
		}
		configCenter = reader
	}

	// get tls config
	tlsConfig, _, err := cc.GetPtr[config.TLSConfig](configCenter, cc.CommonConfType, "tls")
	if err != nil {
		log.Error(ctx, "get tls config failed", log.E(err))
		return nil, err
	}

	// get trace config
	traceConfig, _, err := cc.GetPtr[trace.Option](configCenter, cc.CommonConfType, "trace")
	if err != nil {
		log.Error(ctx, "get trace config failed", log.E(err))
		return nil, err
	}

	return &ConfigInfo{
		CC:        configCenter,
		TLSConf:   tlsConfig,
		traceConf: traceConfig,
	}, nil
}
