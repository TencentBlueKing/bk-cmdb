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

package instances

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

var updateIgnoreKeys = []string{
	common.BKOwnerIDField,
	common.BKDefaultField,
	common.BKInstParentStr,
	common.BKAppIDField,
	common.BKDataStatusField,
	common.BKSupplierIDField,
	common.BKInstIDField,
}

var createIgnoreKeys = []string{
	common.BKOwnerIDField,
	common.BKDefaultField,
	common.BKInstParentStr,
	common.BKAppIDField,
	common.BKSupplierIDField,
	common.BKInstIDField,
	common.BKDataStatusField,
	common.CreateTimeField,
	common.LastTimeField,
	common.BKProcIDField,
}

func FetchBizIDFromInstance(objID string, instanceData mapstr.MapStr) (int64, error) {
	switch objID {
	case common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule, common.BKInnerObjIDProc:
		biz, exist := instanceData[common.BKAppIDField]
		if exist == false {
			return 0, nil
		}
		bizID, err := util.GetInt64ByInterface(biz)
		if err != nil {
			return 0, err
		}
		return bizID, nil
	case common.BKInnerObjIDPlat:
		return 0, nil
	default:
		if _, exist := instanceData[common.MetadataField]; exist == false {
			return 0, nil
		}
		return metadata.ParseBizIDFromData(instanceData)
	}
}

func (m *instanceManager) validBizID(ctx core.ContextParams, bizID int64) error {
	if bizID == 0 {
		return nil
	}
	cond := map[string]interface{}{
		common.BKAppIDField: bizID,
	}
	cnt, err := m.countInstance(ctx, common.BKInnerObjIDApp, cond)
	if err != nil {
		blog.Errorf("search instance biz error %v, rid: %s", err, ctx.ReqID)
		return err
	}
	if cnt != 1 {
		blog.Errorf("biz %d invalid, rid: %s", bizID, ctx.ReqID)
		return ctx.Error.Errorf(common.CCErrCommParamsIsInvalid, "bk_biz_id")
	}
	return nil
}

func (m *instanceManager) validCreateInstanceData(ctx core.ContextParams, objID string, instanceData mapstr.MapStr) error {
	bizID, err := FetchBizIDFromInstance(objID, instanceData)
	if err != nil {
		blog.Errorf("validCreateInstanceData failed, FetchBizIDFromInstance failed, err: %+v, rid: %s", err, ctx.ReqID)
		return ctx.Error.Errorf(common.CCErrCommParamsIsInvalid, "bk_biz_id")
	}

	err = m.validBizID(ctx, bizID)
	if err != nil {
		blog.Errorf("valid biz id error %v, rid: %s", err, ctx.ReqID)
		return err
	}
	valid, err := NewValidator(ctx, m.dependent, objID, bizID)
	if nil != err {
		blog.Errorf("init validator failed %s, rid: %s", err.Error(), ctx.ReqID)
		return err
	}
	for _, key := range valid.requirefields {
		if _, ok := instanceData[key]; !ok {
			blog.Errorf("field [%s] in required for model [%s], input data: %+v, rid: %s", key, objID, instanceData, ctx.ReqID)
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)
		}
	}
	FillLostedFieldValue(ctx.Context, instanceData, valid.propertyslice)
	var instMedataData metadata.Metadata
	instMedataData.Label = make(metadata.Label)
	for key, val := range instanceData {
		if key == common.BKObjIDField {
			// common instance always has no property bk_obj_id, but this field need save to db
			blog.V(9).Infof("skip verify filed: %s, rid: %s", key, ctx.ReqID)
			continue
		}
		if metadata.BKMetadata == key {
			bizID := metadata.GetBusinessIDFromMeta(val)
			if "" != bizID {
				instMedataData.Label.Set(metadata.LabelBusinessID, bizID)
			}
			continue
		}
		if util.InStrArr(createIgnoreKeys, key) {
			// ignore the key field
			continue
		}
		property, ok := valid.propertys[key]
		if !ok {
			delete(instanceData, key)
			continue
			// blog.Errorf("field [%s] is not a valid property for model [%s], rid: %s", key, objID, ctx.ReqID)
			// return valid.errif.CCErrorf(common.CCErrCommParamsIsInvalid, key)
		}
		rawErr := property.Validate(ctx, val, key)
		if rawErr.ErrCode > 0 {
			ccErr := rawErr.ToCCError(valid.errif)
			blog.Errorf("validCreateInstanceData failed, key: %s, value: %s, err: %s, rid: %s", key, val, ccErr.Error(), ctx.ReqID)
			return ccErr
		}
	}
	if instanceData.Exists(metadata.BKMetadata) {
		instanceData.Set(metadata.BKMetadata, instMedataData)
	}

	// module instance's name must coincide with template
	if objID == common.BKInnerObjIDModule {
		if err := m.validateModuleCreate(ctx, instanceData, valid); err != nil {
			if blog.V(9) {
				blog.InfoJSON("validateModuleCreate failed, module: %s, err: %s, rid: %s", instanceData, err, ctx.ReqID)
			}
			return err
		}
	}
	return valid.validCreateUnique(ctx, instanceData, instMedataData, m)
}

func (m *instanceManager) validateModuleCreate(ctx core.ContextParams, instanceData mapstr.MapStr, valid *validator) error {
	svcTplIDIf, exist := instanceData[common.BKServiceTemplateIDField]
	if exist == false {
		return valid.errif.Errorf(common.CCErrCommParamsNeedSet, common.BKServiceTemplateIDField)
	}
	svcTplID, err := util.GetInt64ByInterface(svcTplIDIf)
	if err != nil {
		return valid.errif.Errorf(common.CCErrCommParamsNeedInt, common.BKServiceTemplateIDField)
	}
	if svcTplID == common.ServiceTemplateIDNotSet {
		return nil
	}
	svcCategoryIDIf, exist := instanceData[common.BKServiceCategoryIDField]
	if exist == false {
		return valid.errif.Errorf(common.CCErrCommParamsNeedSet, common.BKServiceCategoryIDField)
	}
	svcCategoryID, err := util.GetInt64ByInterface(svcCategoryIDIf)
	if err != nil {
		return valid.errif.Errorf(common.CCErrCommParamsNeedInt, common.BKServiceCategoryIDField)
	}
	bizIDIf, exist := instanceData[common.BKAppIDField]
	if exist == false {
		return valid.errif.Errorf(common.CCErrCommParamsNeedSet, common.BKAppIDField)
	}
	bizID, err := util.GetInt64ByInterface(bizIDIf)
	if err != nil {
		return valid.errif.Errorf(common.CCErrCommParamsNeedInt, common.MetadataLabelBiz)
	}
	tpl := metadata.ServiceTemplate{}
	filter := map[string]interface{}{
		common.BKFieldID:                svcTplID,
		common.BKServiceCategoryIDField: svcCategoryID,
		common.BKAppIDField:             bizID,
	}
	if err := m.dbProxy.Table(common.BKTableNameServiceTemplate).Find(filter).One(ctx.Context, &tpl); err != nil {
		return valid.errif.Errorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
	}
	instanceData[common.BKModuleNameField] = tpl.Name

	return nil
}

// getHostRelatedBizID 根据主机ID获取所属业务ID
func getHostRelatedBizID(ctx core.ContextParams, dbProxy dal.DB, hostID int64) (bizID int64, ccErr errors.CCErrorCoder) {
	rid := ctx.ReqID
	filter := map[string]interface{}{
		common.BKHostIDField: hostID,
	}
	hostConfig := make([]metadata.ModuleHost, 0)
	if err := dbProxy.Table(common.BKTableNameModuleHostConfig).Find(filter).All(ctx.Context, &hostConfig); err != nil {
		blog.Errorf("getHostRelatedBizID failed, db get failed, hostID: %d, err: %s, rid: %s", hostID, err.Error(), rid)
		return 0, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}
	if len(hostConfig) == 0 {
		blog.Errorf("host module config empty, hostID: %d, rid: %s", hostID, rid)
		return 0, ctx.Error.CCErrorf(common.CCErrHostModuleConfigFailed, hostID)
	}
	bizID = hostConfig[0].AppID
	for _, item := range hostConfig {
		if item.AppID != bizID {
			blog.Errorf("getHostRelatedBizID failed, get multiple bizID, hostID: %d, hostConfig: %+v, rid: %s", hostID, hostConfig, rid)
			return 0, ctx.Error.CCErrorf(common.CCErrCommGetMultipleObject)
		}
	}
	return bizID, nil
}

func (m *instanceManager) validUpdateInstanceData(ctx core.ContextParams, objID string, instanceData mapstr.MapStr, instMetaData metadata.Metadata, instID uint64) error {
	updateData, err := m.getInstDataByID(ctx, objID, instID, m)
	if err != nil {
		blog.ErrorJSON("validUpdateInstanceData failed, getInstDataByID failed, err: %s, objID: %s, instID: %s, rid: %s", err, instID, objID, ctx.ReqID)
		return err
	}
	var bizID int64
	if objID != common.BKInnerObjIDHost {
		bizID, err = FetchBizIDFromInstance(objID, updateData)
		if err != nil {
			blog.ErrorJSON("validUpdateInstanceData failed, FetchBizIDFromInstance failed, err: %s, data: %s, rid: %s", err, updateData, ctx.ReqID)
			return ctx.Error.Errorf(common.CCErrCommParamsIsInvalid, "bk_biz_id")
		}
	} else {
		bizID, err = getHostRelatedBizID(ctx, m.dbProxy, int64(instID))
		if err != nil {
			blog.ErrorJSON("validUpdateInstanceData failed, getHostRelatedBizID failed, hostID: %d, err: %s, rid: %s", instID, err, ctx.ReqID)
			return ctx.Error.CCErrorf(common.CCErrCommGetBusinessIDByHostIDFailed)
		}
	}

	valid, err := NewValidator(ctx, m.dependent, objID, bizID)
	if nil != err {
		blog.Errorf("init validator failed %s, rid: %s", err.Error(), ctx.ReqID)
		return err
	}

	for key, val := range instanceData {

		if util.InStrArr(updateIgnoreKeys, key) {
			// ignore the key field
			continue
		}

		property, ok := valid.propertys[key]
		if !ok {
			delete(instanceData, key)
			continue
		}
		rawErr := property.Validate(ctx, val, key)
		if rawErr.ErrCode > 0 {
			return rawErr.ToCCError(valid.errif)
		}
	}

	for key, val := range instanceData {
		updateData[key] = val
	}
	bizID, err = FetchBizIDFromInstance(objID, updateData)
	if err != nil {
		blog.ErrorJSON("validUpdateInstanceData failed, FetchBizIDFromInstance failed, err: %s, data: %s, rid: %s", err, updateData, ctx.ReqID)
		return ctx.Error.Errorf(common.CCErrCommParamsIsInvalid, "bk_biz_id")
	}
	if bizID != 0 {
		err = m.validBizID(ctx, bizID)
		if err != nil {
			blog.Errorf("valid biz id error %v, rid: %s", err, ctx.ReqID)
			return err
		}
	}

	return valid.validUpdateUnique(ctx, updateData, instMetaData, instID, m)
}
