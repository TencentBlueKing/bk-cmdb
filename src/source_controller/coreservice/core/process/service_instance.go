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

package process

import (
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (p *processOperation)CreateServiceInstance(ctx core.ContextParams, template metadata.ServiceInstance) (*metadata.ServiceInstance, error) {
	return nil, nil
}

func (p *processOperation)GetServiceInstance(ctx core.ContextParams, templateID int64) (*metadata.ServiceInstance, error) {
	return nil, nil
}

func (p *processOperation)UpdateServiceInstance(ctx core.ContextParams, templateID int64, template metadata.ServiceInstance) (*metadata.ServiceInstance, error) {
	return nil, nil
}

func (p *processOperation)ListServiceInstance(ctx core.ContextParams, bizID int64, serviceTemplateID int64, hostID int64) (*metadata.MultipleServiceInstance, error) {
	return nil, nil
}

func (p *processOperation) DeleteServiceInstance(ctx core.ContextParams, processTemplateID int64) error {
	return nil
}
