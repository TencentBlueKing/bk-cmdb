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

package monitor

import (
	"configcenter/src/common/metric"
)

func me() {

	// MetricServer
	conf := metric.Config{
		ModuleName:    "name",
		ServerAddress: "",
	}
	metricActions := metric.NewMetricController(conf, HealthMetric)
	as := []*httpserver.Action{}
	for _, metricAction := range metricActions {
		as = append(as, &httpserver.Action{Verb: common.HTTPSelectGet, Path: metricAction.Path, Handler: func(req *restful.Request, resp *restful.Response) {
			metricAction.HandlerFunc(resp.ResponseWriter, req.Request)
		}})
	}

	ccAPI.httpServ.RegisterWebServer("/", nil, as)
}

// HealthMetric check netservice is health
func HealthMetric() metric.HealthMeta {
	a := api.GetAPIResource()
	meta := metric.HealthMeta{IsHealthy: true}

	for _, item := range meta.Items {
		if item.IsHealthy == false {
			meta.IsHealthy = false
			meta.Message = "hostserver is not healthy"
			break
		}
	}

	return meta
}
