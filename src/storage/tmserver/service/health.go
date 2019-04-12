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

package service

import (
	restful "github.com/emicklei/go-restful"
)

// Healthz health check methods
func (s *coreService) Healthz(req *restful.Request, resp *restful.Response) {
	/*
		meta := metric.HealthMeta{IsHealthy: true}

			// zk health status
			zkItem := metric.HealthItem{IsHealthy: true, Name: types.CCFunctionalityServicediscover}
			if err := s.engine.Ping(); err != nil {
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
	*/
}
