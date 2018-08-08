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
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/metric"
	"configcenter/src/common/types"
	"configcenter/src/txnframe/mongobyc"
	"configcenter/src/txnframe/rpc"
	"configcenter/src/txnframe/server/manager"
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
	"os"
)

type Service struct {
	*backbone.Engine
	txnService *TXRPC
	man        *manager.TxnManager
	db         mongobyc.Client
}

func (s *Service) SetEngine(engine *backbone.Engine) {
	s.Engine = engine
	s.txnService.SetEngine(engine)
}

func (s *Service) SetDB(db mongobyc.Client) {
	s.db = db
}

func (s *Service) SetMan(man *manager.TxnManager) {
	s.man = man
	s.txnService.SetMan(man)
}

func (s *Service) WebService() *restful.WebService {
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)
	restful.SetLogger(log.New(os.Stdout, "restful", log.LstdFlags))
	restful.TraceLogger(log.New(os.Stdout, "restful", log.LstdFlags))
	ws := new(restful.WebService)
	ws.Path("/txn/v3")

	rpcsrv := rpc.NewServer()
	txnService := NewTXRPC(rpcsrv)
	s.txnService = txnService

	ws.Route(ws.Method(http.MethodConnect).Path("rpc").To(func(req *restful.Request, resp *restful.Response) {
		blog.Infof("requeting rpc")
		rpcsrv.ServeHTTP(resp.ResponseWriter, req.Request)
	}))

	return ws
}

func (s *Service) Healthz(req *restful.Request, resp *restful.Response) {
	meta := metric.HealthMeta{IsHealthy: true}

	// zk health status
	zkItem := metric.HealthItem{IsHealthy: true, Name: types.CCFunctionalityServicediscover}
	if err := s.Engine.Ping(); err != nil {
		zkItem.IsHealthy = false
		zkItem.Message = err.Error()
	}
	meta.Items = append(meta.Items, zkItem)

	// // mongodb
	// meta.Items = append(meta.Items, metric.NewHealthItem(types.CCFunctionalityServicediscover, s.db.Ping()))

	// // redis
	// meta.Items = append(meta.Items, metric.NewHealthItem(types.CCFunctionalityServicediscover, s.cache.Ping().Err()))

	for _, item := range meta.Items {
		if item.IsHealthy == false {
			meta.IsHealthy = false
			meta.Message = "txn server is unhealthy"
			break
		}
	}

	info := metric.HealthInfo{
		Module:     types.CC_MODULE_TXC,
		HealthMeta: meta,
		AtTime:     metadata.Now(),
	}

	answer := metric.HealthResponse{
		Code:    common.CCSuccess,
		Data:    info,
		OK:      meta.IsHealthy,
		Result:  meta.IsHealthy,
		Message: meta.Message,
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteEntity(answer)
}
