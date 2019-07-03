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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/selector"
	"configcenter/src/source_controller/coreservice/core"
)

func (p *processOperation) AddLabel(ctx core.ContextParams, option selector.LabelOperateOption) errors.CCErrorCoder {
	if field, err := option.Labels.Validate(); err != nil {
		blog.Infof("addLabel failed, validate failed, field:%s, err: %+v, rid: %s", field, err, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, "label."+field)
	}

	for _, instanceID := range option.InstanceIDs {
		instance, err := p.GetServiceInstance(ctx, instanceID)
		if err != nil {
			blog.Errorf("AddLabel failed, get service instance failed, instanceID: %+v, err: %+v, rid: %s", instanceID, err, ctx.ReqID)
			return err
		}
		instance.Labels.AddLabel(option.Labels)
		if _, err := p.UpdateServiceInstance(ctx, instanceID, *instance); err != nil {
			blog.Errorf("AddLabel failed, update service instance failed, instanceID: %+v, err: %+v, rid: %s", instanceID, err, ctx.ReqID)
			return err
		}
	}
	return nil
}

func (p *processOperation) RemoveLabel(ctx core.ContextParams, option selector.LabelOperateOption) errors.CCErrorCoder {
	for _, instanceID := range option.InstanceIDs {
		instance, err := p.GetServiceInstance(ctx, instanceID)
		if err != nil {
			blog.Errorf("RemoveLabel failed, get service instance failed, instanceID: %+v, err: %+v, rid: %s", instanceID, err, ctx.ReqID)
			return err
		}
		instance.Labels.RemoveLabel(option.Labels)
		if _, err := p.UpdateServiceInstance(ctx, instanceID, *instance); err != nil {
			blog.Errorf("RemoveLabel failed, update service instance failed, instanceID: %+v, err: %+v, rid: %s", instanceID, err, ctx.ReqID)
			return err
		}
	}
	return nil
}
