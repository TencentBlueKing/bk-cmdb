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
 
package metric

import (
	"configcenter/src/common/metric"
	"configcenter/src/framework/core/httpserver"
	"configcenter/src/framework/core/option"
	"github.com/emicklei/go-restful"
)

type Manager struct {
	ms []metric.Action
}

var _ Metric = &Manager{}

func NewManager(opt *option.Options) Metric {
	conf := metric.Config{
		ModuleName:    opt.AppName,
		ServerAddress: opt.Addrport,
	}
	ms := metric.NewMetricController(conf, healthMetric)
	manager := &Manager{
		ms: ms,
	}

	return manager
}

// Actions returns metricActions
func (m *Manager) Actions() []httpserver.Action {
	var httpactions []httpserver.Action
	for _, a := range m.ms {
		httpactions = append(httpactions, httpserver.Action{Method: a.Method, Path: a.Path, Handler: func(req *restful.Request, resp *restful.Response) {
			a.HandlerFunc(resp.ResponseWriter, req.Request)
		}})
	}
	return httpactions
}

// HealthMetric check netservice is health
func healthMetric() metric.HealthMeta {
	meta := metric.HealthMeta{IsHealthy: true}
	return meta
}
