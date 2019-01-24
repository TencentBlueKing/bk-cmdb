/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package app

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	mgo "configcenter/src/storage/mongodb/driver"
	"configcenter/src/storage/tmserver/app/options"
	"configcenter/src/storage/tmserver/service"

	restful "github.com/emicklei/go-restful"
)

// Run main goroute
func Run(ctx context.Context, op *options.ServerOption) error {

	svrInfo, err := newServerInfo(op)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	discover, err := discovery.NewDiscoveryInterface(op.ServConf.RegDiscover)
	if err != nil {
		return fmt.Errorf("connect zookeeper [%s] failed: %v", op.ServConf.RegDiscover, err)
	}

	c := &util.APIMachineryConfig{
		QPS:       1000,
		Burst:     2000,
		TLSConfig: nil,
	}

	machinery, err := apimachinery.NewApiMachinery(c, discover)
	if err != nil {
		return fmt.Errorf("new api machinery failed, err: %v", err)
	}

	coreService := service.New(svrInfo.IP, svrInfo.Port)
	server := backbone.Server{
		ListenAddr: svrInfo.IP,
		ListenPort: svrInfo.Port,
		Handler:    restful.NewContainer().Add(coreService.WebService()),
		TLS:        backbone.TLSConfig{},
	}

	regPath := fmt.Sprintf("%s/%s/%s", types.CC_SERV_BASEPATH, types.CC_MODULE_TXC, svrInfo.IP)
	bonC := &backbone.Config{
		RegisterPath: regPath,
		RegisterInfo: *svrInfo,
		CoreAPI:      machinery,
		Server:       server,
	}

	tmServer := &Server{}
	engine, err := backbone.NewBackbone(ctx, op.ServConf.RegDiscover,
		types.CC_MODULE_TXC,
		op.ServConf.ExConfig,
		tmServer.onConfigUpdate,
		discover,
		bonC)

	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	tmServer.engin = engine
	tmServer.coreService = coreService

	// connect to the mongodb
	errCh := make(chan error)

	for {
		if tmServer.config == nil {
			time.Sleep(time.Second * 2)
			continue
		}

		// connect db
		db := mgo.NewClient(tmServer.config.MongoDB.BuildURI())
		err = db.Open()
		if nil != err {
			errCh <- err
			return fmt.Errorf("connect mongo server failed %s", err.Error())
		}

		// set core service
		coreService.SetConfig(engine, db, tmServer.config.Transaction)

		break
	}

	// waiting to exit
	select {
	case <-ctx.Done():
	case err = <-errCh:
		blog.V(3).Infof("distribution routine stoped %s", err.Error())
	}
	blog.V(3).Infof("process stoped")

	return nil
}
