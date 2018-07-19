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

package operation

import (
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/metric"
	gtypes "configcenter/src/common/types"
	"configcenter/src/scene_server/topo_server/core/types"
)

type HealthOperationInterface interface {
	Health(params types.ContextParams) (*metric.HealthResponse, error)
}

func NewHealthOperation(client apimachinery.ClientSetInterface) HealthOperationInterface {
	return &health{clientSet: client}
}

type health struct {
	clientSet apimachinery.ClientSetInterface
}

func (h *health) Health(params types.ContextParams) (*metric.HealthResponse, error) {
	meta := metric.HealthMeta{IsHealthy: true}

	// zk health status
	zkItem := metric.HealthItem{IsHealthy: true, Name: gtypes.CCFunctionalityServicediscover}
	if err := params.Engin.Ping(); err != nil {
		zkItem.IsHealthy = false
		zkItem.Message = err.Error()
	}
	meta.Items = append(meta.Items, zkItem)

	// object controller
	objCtr := metric.HealthItem{IsHealthy: true, Name: gtypes.CC_MODULE_OBJECTCONTROLLER}
	if _, err := params.Engin.CoreAPI.Healthz().HealthCheck(gtypes.CC_MODULE_OBJECTCONTROLLER); err != nil {
		objCtr.IsHealthy = false
		objCtr.Message = err.Error()
	}
	meta.Items = append(meta.Items, objCtr)

	// audit controller
	auditCtrl := metric.HealthItem{IsHealthy: true, Name: gtypes.CC_MODULE_AUDITCONTROLLER}
	if _, err := params.Engin.CoreAPI.Healthz().HealthCheck(gtypes.CC_MODULE_AUDITCONTROLLER); err != nil {
		auditCtrl.IsHealthy = false
		auditCtrl.Message = err.Error()
	}
	meta.Items = append(meta.Items, auditCtrl)

	for _, item := range meta.Items {
		if item.IsHealthy == false {
			meta.IsHealthy = false
			meta.Message = "topo server is unhealthy"
			break
		}
	}

	info := metric.HealthInfo{
		Module:     gtypes.CC_MODULE_TOPO,
		HealthMeta: meta,
		AtTime:     gtypes.Now(),
	}

	return &metric.HealthResponse{
		Code:    common.CCSuccess,
		Data:    info,
		OK:      meta.IsHealthy,
		Result:  meta.IsHealthy,
		Message: meta.Message,
	}, nil

}
