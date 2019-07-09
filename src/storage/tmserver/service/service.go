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

package service

import (
	"context"
	"net/http"

	restful "github.com/emicklei/go-restful"

	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/storage/mongodb"
	"configcenter/src/storage/rpc"
	"configcenter/src/storage/tmserver/app/options"
	"configcenter/src/storage/tmserver/core"
	"configcenter/src/storage/tmserver/core/session"
	"configcenter/src/storage/types"
)

// Service service methods
type Service interface {
	WebService() *restful.WebService
	SetConfig(engin *backbone.Engine, db mongodb.Client, txnCfg options.TransactionConfig) error
}

// New create a new service instance
func New(ip string, port uint) Service {

	return &coreService{
		listenIP:   ip,
		listenPort: port,
	}
}

type coreService struct {
	engine     *backbone.Engine
	rpc        *rpc.Server
	dbProxy    mongodb.Client
	core       core.Core
	listenIP   string
	listenPort uint
}

func (s *coreService) SetConfig(engin *backbone.Engine, db mongodb.Client, txnCfg options.TransactionConfig) error {

	// set config
	s.engine = engin
	s.dbProxy = db
	s.rpc = rpc.NewServer()

	// init all handlers
	s.rpc.Handle(types.CommandRDBOperation, s.DBOperation)
	s.rpc.HandleStream(types.CommandWatchTransactionOperation, s.WatchTransaction)

	// create a new core instance
	sess, err := session.New(
		core.ContextParams{
			Context:  context.Background(),
			ListenIP: s.listenIP,
		}, txnCfg, db, s.listenIP)
	if err != nil {
		return err
	}
	go func() {
		if err := sess.Run(); err != nil {
			blog.Errorf("tmserver stoped with error: %v", err)
		} else {
			blog.Errorf("tmserver stoped")
		}
	}()
	s.core = core.New(sess, db)
	return nil
}

func (s *coreService) WebService() *restful.WebService {

	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)
	restful.SetLogger(&blog.GlogWriter{})
	restful.TraceLogger(&blog.GlogWriter{})

	ws := &restful.WebService{}
	ws.Path("/txn/v3").Filter(s.engine.Metric().RestfulMiddleWare)

	ws.Route(ws.Method(http.MethodConnect).Path("rpc").To(func(req *restful.Request, resp *restful.Response) {
		if sub, ok := resp.ResponseWriter.(*restful.Response); ok {
			s.rpc.ServeHTTP(sub.ResponseWriter, req.Request)
			return
		}
		s.rpc.ServeHTTP(resp.ResponseWriter, req.Request)
	}))

	return ws
}
