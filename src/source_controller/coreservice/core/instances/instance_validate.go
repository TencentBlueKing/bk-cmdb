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
	stderr "errors"
	"fmt"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/valid"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/thirdparty/hooks"
)

var updateIgnoreKeys = []string{
	common.BKOwnerIDField,
	common.BKDefaultField,
	common.BKInstParentStr,
	common.BKAppIDField,
	common.BKDataStatusField,
	common.BKInstIDField,
}

var createIgnoreKeys = []string{
	common.BKOwnerIDField,
	common.BKDefaultField,
	common.BKInstParentStr,
	common.BKAppIDField,
	common.BKInstIDField,
	common.BKDataStatusField,
	common.CreateTimeField,
	common.LastTimeField,
	common.BKProcIDField,
}

func (m *instanceManager) getBizIDFromInstance(kit *rest.Kit, objID string, instanceData mapstr.MapStr, validTye string,
	instID int64) (int64, error) {
	bizID, err := m.fetchBizIDFromInstance(kit, objID, instanceData, validTye, instID)
	if err != nil {
		blog.Errorf("fetchBizIDFromInstance failed, err: %v, rid: %s", err, kit.Rid)
		return 0, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "bk_biz_id")
	}
	if err := m.validBizID(kit, bizID); err != nil {
		blog.Errorf("validBizID failed, err %v, rid: %s", err, kit.Rid)
		return 0, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "bk_biz_id")
	}
	return bizID, nil
}

func (m *instanceManager) fetchBizIDFromInstance(kit *rest.Kit, objID string, instanceData mapstr.MapStr,
	validTye string, instID int64) (int64, error) {
	switch objID {
	case common.BKInnerObjIDHost:
		if validTye == common.ValidUpdate {
			bizID, err := m.getHostRelatedBizID(kit, instID)
			if err != nil {
				blog.Errorf("getHostRelatedBizID failed, hostID: %d, err: %s, rid: %s", instID, err, kit.Rid)
				return 0, kit.CCError.CCErrorf(common.CCErrCommGetBusinessIDByHostIDFailed)
			}
			return bizID, nil
		}
		fallthrough
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
	case common.BKInnerObjIDBizSet, common.BKInnerObjIDProject, common.BKInnerObjIDPlat:
		return 0, nil
	default:
		biz, exist := instanceData[common.BKAppIDField]
		if exist == false {
			return 0, nil
		}
		bizID, err := util.GetInt64ByInterface(biz)
		if err != nil {
			return 0, err
		}
		return bizID, nil
	}
}

// getHostRelatedBizID 根据主机ID获取所属业务ID
func (m *instanceManager) getHostRelatedBizID(kit *rest.Kit, hostID int64) (bizID int64, ccErr errors.CCErrorCoder) {
	rid := kit.Rid
	filter := map[string]interface{}{
		common.BKHostIDField: hostID,
	}

	relation := new(metadata.ModuleHost)
	if err := mongodb.Client().Table(common.BKTableNameModuleHostConfig).Find(filter).Fields(common.BKAppIDField).
		One(kit.Ctx, relation); err != nil {
		blog.Errorf("getHostRelatedBizID failed, db get failed, hostID: %d, err: %s, rid: %s", hostID, err.Error(), rid)
		return 0, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	return relation.AppID, nil
}

func (m *instanceManager) validBizID(kit *rest.Kit, bizID int64) error {
	if bizID == 0 {
		return nil
	}
	cond := map[string]interface{}{
		common.BKAppIDField: bizID,
	}
	cnt, err := m.countInstance(kit, common.BKInnerObjIDApp, cond)
	if err != nil {
		blog.Errorf("search instance biz error %v, rid: %s", err, kit.Rid)
		return err
	}
	if cnt != 1 {
		blog.Errorf("biz %d invalid, rid: %s", bizID, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "bk_biz_id")
	}
	return nil
}

func (m *instanceManager) newValidator(kit *rest.Kit, objID string, bizID int64) (*validator, error) {
	validator, err := NewValidator(kit, m.dependent, objID, bizID, m.language)
	if nil != err {
		blog.Errorf("newValidator failed , NewValidator err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	return validator, nil
}

func (m *instanceManager) validCreateInstanceData(kit *rest.Kit, objID string, instanceData mapstr.MapStr,
	valid *validator) error {
	for _, key := range valid.requireFields {
		if _, ok := instanceData[key]; !ok {
			blog.Errorf("field [%s] in required for model [%s], input data: %+v, rid: %s", key, objID, instanceData,
				kit.Rid)
			return valid.errIf.Errorf(common.CCErrCommParamsNeedSet, key)
		}
	}
	if err := FillLostFieldValue(kit.Ctx, instanceData, valid.propertySlice); err != nil {
		return err
	}

	if err := m.validCloudID(kit, objID, instanceData); err != nil {
		return err
	}

	isMainline, err := m.isMainlineObject(kit, objID)
	if err != nil {
		return err
	}

	if err := m.validMainlineInstanceData(kit, objID, instanceData, isMainline); err != nil {
		return err
	}

	err = hooks.ValidateBizBsTopoHook(kit, objID, instanceData, nil, common.ValidCreate, m.clientSet)
	if err != nil {
		blog.Errorf("validate biz bk_bs_topo attribute failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	if err := hooks.ValidateHostBsInfoHook(kit, objID, instanceData); err != nil {
		blog.Errorf("validate host attribute hook failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	err = m.validateCreateInstValue(kit, instanceData, valid)
	if err != nil {
		return err
	}

	skip, err := hooks.IsSkipValidateHook(kit, objID, instanceData)
	if err != nil {
		blog.Errorf("check is skip validate %s hook failed, err: %v, rid: %s", objID, err, kit.Rid)
		return err
	}

	if skip {
		return nil
	}

	if err := m.changeStringToTime(instanceData, valid.propertySlice); err != nil {
		blog.Errorf("there is an error in converting the time type string to the time type, err: %v, rid: %s", err,
			kit.Rid)
		return err
	}

	switch objID {
	case common.BKInnerObjIDModule:
		// module instance's name must coincide with template
		if err := m.validateModuleCreate(kit, instanceData, valid); err != nil {
			blog.Errorf("validate create module failed, module: %v, err: %v, rid: %s", instanceData, err, kit.Rid)
			return err
		}

	case common.BKInnerObjIDHost:
		if err := m.validateHostCreate(kit, instanceData, valid); err != nil {
			blog.Errorf("validate create host failed, host: %v, err: %v, rid: %s", instanceData, err, kit.Rid)
			return err
		}
	}
	return nil
}

func (m *instanceManager) validateCreateInstValue(kit *rest.Kit, instanceData mapstr.MapStr, valid *validator) error {
	for key, val := range instanceData {
		if key == common.BKObjIDField {
			// common instance always has no property bk_obj_id, but this field need save to db
			blog.V(9).Infof("skip verify filed: %s, rid: %s", key, kit.Rid)
			continue
		}
		if util.InStrArr(createIgnoreKeys, key) {
			// ignore the key field
			continue
		}
		property, ok := valid.properties[key]
		if !ok {
			delete(instanceData, key)
			continue
		}
		if value, ok := val.(string); ok {
			val = strings.TrimSpace(value)
			instanceData[key] = val
		}

		rawErr := property.Validate(kit.Ctx, val, key)
		if rawErr.ErrCode != 0 {
			blog.Errorf("validCreateInstanceData failed, key: %s, value: %s, err: %s, rid: %s", key, val,
				kit.CCError.Error(rawErr.ErrCode), kit.Rid)
			return rawErr.ToCCError(kit.CCError)
		}
		// 在Validate里面没有对枚举引用的值进行校验，只是对其数据类型做了基本的校验，
		// 因为对引用值是否存在校验需要查询数据库，所以放在Validate里面校验不太合适
		if property.PropertyType == common.FieldTypeEnumQuote {
			if err := m.validInstIDs(kit, property, val); err != nil {
				return err
			}
		}

		// remove inner table value
		if property.PropertyType == common.FieldTypeInnerTable {
			delete(instanceData, property.PropertyID)
		}
	}
	return nil
}

func (m *instanceManager) validateModuleCreate(kit *rest.Kit, instanceData mapstr.MapStr, valid *validator) error {
	svcTplIDIf, exist := instanceData[common.BKServiceTemplateIDField]
	if exist == false {
		return valid.errIf.Errorf(common.CCErrCommParamsNeedSet, common.BKServiceTemplateIDField)
	}
	svcTplID, err := util.GetInt64ByInterface(svcTplIDIf)
	if err != nil {
		return valid.errIf.Errorf(common.CCErrCommParamsNeedInt, common.BKServiceTemplateIDField)
	}
	if svcTplID == common.ServiceTemplateIDNotSet {
		return nil
	}
	svcCategoryIDIf, exist := instanceData[common.BKServiceCategoryIDField]
	if exist == false {
		return valid.errIf.Errorf(common.CCErrCommParamsNeedSet, common.BKServiceCategoryIDField)
	}
	svcCategoryID, err := util.GetInt64ByInterface(svcCategoryIDIf)
	if err != nil {
		return valid.errIf.Errorf(common.CCErrCommParamsNeedInt, common.BKServiceCategoryIDField)
	}
	bizIDIf, exist := instanceData[common.BKAppIDField]
	if exist == false {
		return valid.errIf.Errorf(common.CCErrCommParamsNeedSet, common.BKAppIDField)
	}
	bizID, err := util.GetInt64ByInterface(bizIDIf)
	if err != nil {
		return valid.errIf.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)
	}
	tpl := metadata.ServiceTemplate{}
	filter := map[string]interface{}{
		common.BKFieldID:                svcTplID,
		common.BKServiceCategoryIDField: svcCategoryID,
		common.BKAppIDField:             bizID,
	}
	if err := mongodb.Client().Table(common.BKTableNameServiceTemplate).Find(filter).One(kit.Ctx, &tpl); err != nil {
		return valid.errIf.Errorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
	}
	instanceData[common.BKModuleNameField] = tpl.Name

	return nil
}

func (m *instanceManager) validUpdateInstanceData(kit *rest.Kit, objID string, updateData, instanceData mapstr.MapStr,
	valid *validator, canEditAll, isMainline bool) error {

	if err := m.validCloudID(kit, objID, updateData); err != nil {
		return err
	}

	if err := hooks.ValidUpdateCloudIDHook(kit, objID, instanceData, updateData); err != nil {
		return err
	}

	if err := m.validMainlineInstanceData(kit, objID, updateData, isMainline); err != nil {
		return err
	}

	err := hooks.ValidateBizBsTopoHook(kit, objID, instanceData, updateData, common.ValidUpdate, m.clientSet)
	if err != nil {
		blog.Errorf("validate biz bk_bs_topo attribute failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	if err := hooks.ValidateHostBsInfoHook(kit, objID, updateData); err != nil {
		blog.Errorf("validate host bk_bs_info attribute failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	isInnerModel := common.IsInnerModel(objID) || isMainline

	err = m.validOneUpdateInstKeyVal(kit, valid, updateData, isInnerModel, canEditAll)
	if err != nil {
		return err
	}

	if err := m.changeStringToTime(updateData, valid.propertySlice); err != nil {
		blog.Errorf("there is an error in converting the time type string to the time type, err: %v, rid: %s", err,
			kit.Rid)
		return err
	}

	skip, err := hooks.IsSkipValidateHook(kit, objID, instanceData)
	if err != nil {
		blog.Errorf("check is skip validate %s hook failed, err: %v, rid: %s", objID, err, kit.Rid)
		return err
	}

	if skip {
		return nil
	}

	return nil
}

func (m *instanceManager) validOneUpdateInstKeyVal(kit *rest.Kit, valid *validator, updateData mapstr.MapStr,
	isInnerModel, canEditAll bool) error {

	for key, val := range updateData {
		if isInnerModel && util.InStrArr(updateIgnoreKeys, key) {
			// ignore the key field
			continue
		}

		property, ok := valid.properties[key]
		if !ok || (!property.IsEditable && !canEditAll) {
			delete(updateData, key)
			continue
		}

		// right now inner table should be updated as quoted instance, cannot update in source instance
		if property.PropertyType == common.FieldTypeInnerTable {
			delete(updateData, key)
			continue
		}

		if value, ok := val.(string); ok {
			val = strings.TrimSpace(value)
			updateData[key] = val
		}
		rawErr := property.Validate(kit.Ctx, val, key)
		if rawErr.ErrCode != 0 {
			blog.ErrorJSON("validUpdateInstanceData failed, err: %s, val: %s, key:%s, rid: %s",
				rawErr.ToCCError(kit.CCError), val, key, kit.Rid)
			return rawErr.ToCCError(kit.CCError)
		}
		// 在Validate里面没有对枚举引用的值进行校验，只是对其数据类型做了基本的校验，
		// 因为对引用值是否存在校验需要查询数据库，所以放在Validate里面校验不太合适
		if property.PropertyType == common.FieldTypeEnumQuote {
			if err := m.validInstIDs(kit, property, val); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *instanceManager) isMainlineObject(kit *rest.Kit, objID string) (bool, error) {
	// judge whether it is an inner mainline model
	if common.IsInnerMainlineModel(objID) {
		return true, nil
	}

	// if not inner mainline model, then judge whether it is a self defined mainline layer
	mainlineCond := map[string]interface{}{
		common.AssociationKindIDField: common.AssociationKindMainline,
		common.BKAsstObjIDField:       objID,
	}
	cnt, err := mongodb.Client().Table(common.BKTableNameObjAsst).Find(mainlineCond).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count mainline association failed, err: %v, cond: %#v, rid: %s", err, mainlineCond, kit.Rid)
		return false, err
	}

	if cnt > 0 {
		return true, nil
	}
	return false, nil
}

func (m *instanceManager) validMainlineInstanceData(kit *rest.Kit, objID string, instanceData mapstr.MapStr,
	isMainline bool) error {

	if !isMainline {
		return nil
	}

	// validate instance name
	nameField := metadata.GetInstNameFieldName(objID)
	if nameVal, exist := instanceData[nameField]; exist {
		name, ok := nameVal.(string)
		if !ok {
			return kit.CCError.CCErrorf(common.CCErrCommParamsNeedString, nameField)
		}

		name, err := valid.ValidTopoNameField(name, nameField, kit.CCError)
		if err != nil {
			return err
		}
		instanceData[nameField] = name
	}

	// validate bk_parent_id
	if instanceData.Exists(common.BKParentIDField) {
		parentID, err := util.GetInt64ByInterface(instanceData[common.BKParentIDField])
		if err != nil || parentID <= 0 {
			blog.Errorf("invalid bk_parent_id value: %#v, err: %v, rid: %s", instanceData[common.BKParentIDField], err,
				kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKParentIDField)
		}
	}

	return nil
}

// validCloudID valid the bk_cloud_id
func (m *instanceManager) validCloudID(kit *rest.Kit, objID string, instanceData mapstr.MapStr) error {
	if objID == common.BKInnerObjIDHost {
		if instanceData.Exists(common.BKCloudIDField) {
			if cloudID, err := instanceData.Int64(common.BKCloudIDField); err != nil || cloudID < 0 {
				blog.Errorf("invalid bk_cloud_id value: %#v, err: %v, rid: %s", instanceData[common.BKCloudIDField],
					err, kit.Rid)
				return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudIDField)
			}
		}
	}

	return nil
}

func (m *instanceManager) changeStringToTime(valData mapstr.MapStr, properties []metadata.Attribute) error {
	for _, field := range properties {
		if field.PropertyType != common.FieldTypeTime {
			continue
		}
		val, ok := valData[field.PropertyID]
		if ok == false || val == nil {
			continue
		}

		_, ok = val.(time.Time)
		if ok {
			continue
		}

		valStr, ok := val.(string)
		if ok == false {
			return stderr.New("it is not a string of time type")
		}
		if timeType, isTime := util.IsTime(valStr); isTime {
			valData[field.PropertyID] = util.Str2Time(valStr, timeType)
			continue
		}
		return stderr.New("can not convert value from string type to time type")
	}
	return nil
}

// getValidatorsFromInstances get validators from instances, returns the mapping of instance index to its validator
func (m *instanceManager) getValidatorsFromInstances(kit *rest.Kit, objID string, instanceData []mapstr.MapStr,
	validTye string) ([]*validator, error) {

	instLen := len(instanceData)
	if instLen == 0 {
		return make([]*validator, 0), nil
	}

	bizIDs := make([]int64, instLen)
	needSearchBizInstIDs := make([]int64, 0)
	needSearchBizInstIDIndexMap := make(map[int64]int)
	for index, instance := range instanceData {
		switch objID {
		case common.BKInnerObjIDPlat, common.BKInnerObjIDBizSet:
			bizIDs[index] = 0
		case common.BKInnerObjIDHost:
			if validTye == common.ValidUpdate {
				hostID, err := util.GetInt64ByInterface(instance[common.BKHostIDField])
				if err != nil {
					blog.Errorf("parse host id failed, err: %v, data: %#v, rid: %s", err, instance, kit.Rid)
					return nil, err
				}
				needSearchBizInstIDs = append(needSearchBizInstIDs, hostID)
				needSearchBizInstIDIndexMap[hostID] = index
			}
			fallthrough
		default:
			biz, exist := instance[common.BKAppIDField]
			if !exist {
				bizIDs[index] = 0
				break
			}
			bizID, err := util.GetInt64ByInterface(biz)
			if err != nil {
				blog.Errorf("parse biz id failed, err: %v, obj: %s, data: %#v, rid: %s", err, objID, instance, kit.Rid)
				return nil, err
			}
			bizIDs[index] = bizID
		}
	}

	// get biz id for hosts that need to be updated from db
	if len(needSearchBizInstIDs) > 0 && objID == common.BKInnerObjIDHost {
		filter := map[string]interface{}{
			common.BKHostIDField: map[string]interface{}{common.BKDBIN: needSearchBizInstIDs},
		}

		relations := make([]metadata.ModuleHost, 0)
		if err := mongodb.Client().Table(common.BKTableNameModuleHostConfig).Find(filter).Fields(
			common.BKAppIDField, common.BKHostIDField).All(kit.Ctx, &relations); err != nil {
			blog.Errorf("get hosts(%v) related bizID failed, err: %v, rid: %s", needSearchBizInstIDs, err, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
		}

		for _, relation := range relations {
			bizIDs[needSearchBizInstIDIndexMap[relation.HostID]] = relation.AppID
		}
	}

	if err := m.validBizIDs(kit, bizIDs); err != nil {
		blog.Errorf("valid biz ids(%+v) failed, err %v, rid: %s", bizIDs, err, kit.Rid)
		return nil, err
	}

	bizValidatorMap, err := NewValidators(kit, m.dependent, objID, bizIDs, m.language)
	if err != nil {
		blog.Errorf("new validators failed, err: %v, objID: %s, bizIDs: %v, rid: %s", err, objID, bizIDs, kit.Rid)
		return nil, err
	}

	instValidators := make([]*validator, instLen)
	for index, bizID := range bizIDs {
		instValidators[index] = bizValidatorMap[bizID]
	}
	return instValidators, nil
}

// validBizIDs validate if all the biz id's corresponding biz exists
func (m *instanceManager) validBizIDs(kit *rest.Kit, bizIDs []int64) error {
	uniqueBizIDs := make([]int64, 0)
	bizIDMap := make(map[int64]struct{})
	for _, bizID := range bizIDs {
		if bizID == 0 {
			continue
		}
		if _, exists := bizIDMap[bizID]; exists {
			continue
		}
		bizIDMap[bizID] = struct{}{}
		uniqueBizIDs = append(uniqueBizIDs, bizID)
	}

	if len(bizIDs) == 0 {
		return nil
	}

	cond := map[string]interface{}{
		common.BKAppIDField: map[string]interface{}{common.BKDBIN: uniqueBizIDs},
	}
	cnt, err := m.countInstance(kit, common.BKInnerObjIDApp, cond)
	if err != nil {
		blog.Errorf("search biz failed, err: %v, cond: %#v, rid: %s", err, cond, kit.Rid)
		return err
	}

	if int(cnt) != len(uniqueBizIDs) {
		blog.Errorf("instance biz ids(%+v) contains invalid biz, rid: %s", uniqueBizIDs, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, common.BKAppIDField)
	}
	return nil
}

func (m *instanceManager) validateHostCreate(kit *rest.Kit, instanceData mapstr.MapStr, valid *validator) error {
	// at least one of bk_host_innerip and bk_host_innerip_v6 attribute needs to be passed, because can not validate in
	// db, validate it here
	innerIPv4, ipv4Exist := instanceData[common.BKHostInnerIPField]
	innerIPv6, ipv6Exist := instanceData[common.BKHostInnerIPv6Field]
	if (!ipv4Exist || innerIPv4 == "") && (!ipv6Exist || innerIPv6 == "") {
		return valid.errIf.Errorf(common.CCErrCommAtLeastSetOneVal, common.BKHostInnerIPField,
			common.BKHostInnerIPv6Field)
	}
	return nil
}

// valid enum quote inst id is exist
func (m *instanceManager) validInstIDs(kit *rest.Kit, property metadata.Attribute, val interface{}) error {
	if property.Option == nil {
		return fmt.Errorf("option params is invalid")
	}

	valIDs, ok := val.([]interface{})
	if !ok {
		blog.Errorf("convert val to interface slice failed, val type: %T, rid: %s", val, kit.Rid)
		return fmt.Errorf("convert val to interface slice failed, val: %v", val)
	}

	if len(valIDs) == 0 {
		if property.IsRequired {
			blog.Errorf("enum quote inst id is null, rid: %s", kit.Rid)
			return fmt.Errorf("enum quote inst id is null, please set the correct value")
		}
		return nil
	}

	if property.IsMultiple == nil {
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKIsMultipleField)
	}
	if !(*property.IsMultiple) && len(valIDs) != 1 {
		blog.Errorf("enum quote is single choice, but inst id is multiple, rid: %s", kit.Rid)
		return kit.CCError.CCError(common.CCErrCommParamsNeedSingleChoice)
	}

	valEnumIDMap := make(map[int64]struct{}, 0)
	for _, valID := range valIDs {
		valEnumID, err := util.GetInt64ByInterface(valID)
		if err != nil {
			blog.Errorf("get valEnumID failed, valID type is %T, err: %v, rid: %s", valID, err, kit.Rid)
			return err
		}

		if valEnumID == 0 {
			return fmt.Errorf("enum quote instID is %d, it is illegal", valEnumID)
		}
		valEnumIDMap[valEnumID] = struct{}{}
	}

	if len(valEnumIDMap) == 0 {
		return fmt.Errorf("enum quote instID is null, valEnumIDMap: %v", valEnumIDMap)
	}

	arrOption, err := metadata.ParseEnumQuoteOption(kit.Ctx, property.Option)
	if len(arrOption) == 0 {
		return fmt.Errorf("parse enum quote option data, but is null")
	}
	var quoteObjID string
	for _, o := range arrOption {
		if quoteObjID == "" {
			quoteObjID = o.ObjID
		} else if quoteObjID != o.ObjID {
			return fmt.Errorf("enum quote objID not unique, objID: %s", quoteObjID)
		}
	}

	valEnumIDs := make([]int64, 0)
	for valEnumID := range valEnumIDMap {
		valEnumIDs = append(valEnumIDs, valEnumID)
	}
	tableName := common.GetInstTableName(quoteObjID, kit.SupplierAccount)
	cond := map[string]interface{}{
		common.GetInstIDField(quoteObjID): map[string]interface{}{
			common.BKDBIN: valEnumIDs,
		},
	}
	cnt, err := mongodb.Client().Table(tableName).Find(cond).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count inst failed, err: %v, cond: %#v, rid: %s", err, cond, kit.Rid)
		return err
	}
	if len(valEnumIDs) != int(cnt) {
		return fmt.Errorf("inst not exist, rid: %s", kit.Rid)
	}

	return nil
}
