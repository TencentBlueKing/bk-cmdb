package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/emicklei/go-restful"

	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/scene_server/datacollection/app/options"
	"configcenter/src/scene_server/datacollection/datacollection"
	svc "configcenter/src/scene_server/datacollection/service"
	"configcenter/src/storage/mgoclient"
	"configcenter/src/storage/redisclient"
)

func Run(ctx context.Context, op *options.ServerOption) error {
	svrInfo, err := newServerInfo(op)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	c := &util.APIMachineryConfig{
		ZkAddr:    op.ServConf.RegDiscover,
		QPS:       1000,
		Burst:     2000,
		TLSConfig: nil,
	}

	machinery, err := apimachinery.NewApiMachinery(c)
	if err != nil {
		return fmt.Errorf("new api machinery failed, err: %v", err)
	}

	service := new(svc.Service)
	server := backbone.Server{
		ListenAddr: svrInfo.IP,
		ListenPort: svrInfo.Port,
		Handler:    restful.NewContainer().Add(service.WebService()),
		TLS:        backbone.TLSConfig{},
	}

	regPath := fmt.Sprintf("%s/%s/%s", types.CC_SERV_BASEPATH, types.CC_MODULE_DATACOLLECTION, svrInfo.IP)
	bonC := &backbone.Config{
		RegisterPath: regPath,
		RegisterInfo: *svrInfo,
		CoreAPI:      machinery,
		Server:       server,
	}

	process := new(DCServer)
	engine, err := backbone.NewBackbone(ctx, op.ServConf.RegDiscover,
		types.CC_MODULE_DATACOLLECTION,
		op.ServConf.ExConfig,
		process.onHostConfigUpdate,
		bonC)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	service.Engine = engine
	process.Core = engine
	process.Service = service
	for {
		if process.Config == nil {
			time.Sleep(time.Second * 2)
			blog.V(3).Info("config not found, retry 2s later")
			continue
		}

		db, err := mgoclient.NewFromConfig(process.Config.MongoDB)
		if err != nil {
			return fmt.Errorf("connect mongo server failed %s", err.Error())
		}
		err = db.Open()
		if err != nil {
			return fmt.Errorf("connect mongo server failed %s", err.Error())
		}
		process.Service.SetDB(db)

		cache, err := redisclient.NewFromConfig(process.Config.CCRedis)
		if err != nil {
			return fmt.Errorf("connect redis server failed %s", err.Error())
		}
		process.Service.SetCache(cache)
		break
	}

	datacollection.NewDataCollection(process.Config, process.Core)

	select {}
	return nil
}

type DCServer struct {
	Core    *backbone.Engine
	Config  *options.Config
	Service *svc.Service
}

func (h *DCServer) onHostConfigUpdate(previous, current cc.ProcessConfig) {
	if len(current.ConfigMap) > 0 {
		h.Config = new(options.Config)
		dbprefix := "mongodb"
		h.Config.MongoDB.Address = current.ConfigMap[dbprefix+".host"]
		h.Config.MongoDB.User = current.ConfigMap[dbprefix+".usr"]
		h.Config.MongoDB.Password = current.ConfigMap[dbprefix+".pwd"]
		h.Config.MongoDB.Database = current.ConfigMap[dbprefix+".database"]
		h.Config.MongoDB.Port = current.ConfigMap[dbprefix+".port"]
		h.Config.MongoDB.MaxOpenConns = current.ConfigMap[dbprefix+".maxOpenConns"]
		h.Config.MongoDB.MaxIdleConns = current.ConfigMap[dbprefix+".maxIDleConns"]

		ccredisPrefix := "redis"
		h.Config.CCRedis.Address = current.ConfigMap[ccredisPrefix+".host"]
		h.Config.CCRedis.Password = current.ConfigMap[ccredisPrefix+".pwd"]
		h.Config.CCRedis.Database = current.ConfigMap[ccredisPrefix+".database"]
		h.Config.CCRedis.Port = current.ConfigMap[ccredisPrefix+".port"]

		snapPrefix := "snap-redis"
		h.Config.SnapRedis.Address = current.ConfigMap[snapPrefix+".host"]
		h.Config.SnapRedis.Password = current.ConfigMap[snapPrefix+".pwd"]
		h.Config.SnapRedis.Database = current.ConfigMap[snapPrefix+".database"]
		h.Config.SnapRedis.Port = current.ConfigMap[snapPrefix+".port"]

		discoverPrefix := "discover-redis"
		h.Config.DiscoverRedis.Address = current.ConfigMap[discoverPrefix+".host"]
		h.Config.DiscoverRedis.Password = current.ConfigMap[discoverPrefix+".pwd"]
		h.Config.DiscoverRedis.Database = current.ConfigMap[discoverPrefix+".database"]
		h.Config.DiscoverRedis.Port = current.ConfigMap[discoverPrefix+".port"]
	}
}

func newServerInfo(op *options.ServerOption) (*types.ServerInfo, error) {
	ip, err := op.ServConf.GetAddress()
	if err != nil {
		return nil, err
	}

	port, err := op.ServConf.GetPort()
	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	info := &types.ServerInfo{
		IP:       ip,
		Port:     port,
		HostName: hostname,
		Scheme:   "http",
		Version:  version.GetVersion(),
		Pid:      os.Getpid(),
	}
	return info, nil
}
