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

	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	mgo "configcenter/src/storage/mongodb/driver"
	"configcenter/src/storage/tmserver/app/options"
	"configcenter/src/storage/tmserver/service"

	"github.com/emicklei/go-restful"
)

// Run main goroute
func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := newServerInfo(op)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	coreService := service.New(svrInfo.IP, svrInfo.Port)
	tmServer := &Server{}
	input := &backbone.BackboneParameter{
		ConfigUpdate: tmServer.onConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		Regdiscv:     op.ServConf.RegDiscover,
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	tmServer.engin = engine
	tmServer.coreService = coreService

	// connect to the mongodb
	errCh := make(chan error, 1)

	for {
		tmServer.configLock.Lock()
		if tmServer.config == nil {
			tmServer.configLock.Unlock()
			blog.Infof("config is empty, retry 2s later")
			time.Sleep(time.Second * 2)
			continue
		}
		tmServer.configLock.Unlock()

		// connect db
		blog.Infof("connecting to mongo %v", tmServer.config.MongoDB.BuildURI())
		db := mgo.NewClient(tmServer.config.MongoDB.BuildURI())
		err = db.Open()
		if nil != err {
			return fmt.Errorf("connect mongo server failed %s", err.Error())
		}

		blog.Infof("connected to mongo %v", tmServer.config.MongoDB.BuildURI())

		// set logics service
		coreService.SetConfig(engine, db, tmServer.config.Transaction)
		break
	}
	tmServer.engin = engine
	handler := restful.NewContainer().Add(coreService.WebService())
	err = backbone.StartServer(ctx, cancel, engine, handler, true)
	if err != nil {
		return err
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
