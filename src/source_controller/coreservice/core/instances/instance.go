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
	"strconv"
	"strings"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/mongodb/instancemapping"
	"configcenter/src/thirdparty/hooks"
)

var _ core.InstanceOperation = (*instanceManager)(nil)

type instanceManager struct {
	dependent OperationDependences
	language  language.CCLanguageIf
	clientSet apimachinery.ClientSetInterface
}

// New create a new instance manager instance
func New(dependent OperationDependences, language language.CCLanguageIf,
	clientSet apimachinery.ClientSetInterface) core.InstanceOperation {
	return &instanceManager{
		dependent: dependent,
		language:  language,
		clientSet: clientSet,
	}
}

func (m *instanceManager) instCnt(kit *rest.Kit, objID string, cond mapstr.MapStr) (cnt uint64, exists bool,
	err error) {
	tableName := common.GetInstTableName(objID, kit.SupplierAccount)
	cnt, err = mongodb.Client().Table(tableName).Find(cond).Count(kit.Ctx)
	exists = 0 != cnt
	return cnt, exists, err
}

// CreateModelInstance TODO
func (m *instanceManager) CreateModelInstance(kit *rest.Kit, objID string,
	inputParam metadata.CreateModelInstance) (*metadata.CreateOneDataResult, error) {
	rid := util.ExtractRequestIDFromContext(kit.Ctx)

	inputParam.Data.Set(common.BKOwnerIDField, kit.SupplierAccount)
	bizID, err := m.getBizIDFromInstance(kit, objID, inputParam.Data, common.ValidCreate, 0)
	if err != nil {
		blog.Errorf("CreateModelInstance failed, getBizIDFromInstance err:%v, objID:%s, data:%#v, rid:%s", err, objID,
			inputParam.Data, kit.Rid)
		return nil, err
	}
	validator, err := m.newValidator(kit, objID, bizID)
	if err != nil {
		blog.Errorf("CreateModelInstance failed, newValidator err:%v, objID:%s, data:%#v, rid:%s", err, objID,
			inputParam.Data, kit.Rid)
		return nil, err
	}

	tableData := mapstr.New()
	for key, val := range inputParam.Data {
		if validator.properties[key].PropertyType == common.FieldTypeInnerTable {
			tableData[key] = val
		}
	}

	err = m.validCreateInstanceData(kit, objID, inputParam.Data, validator)
	if nil != err {
		blog.Errorf("CreateModelInstance failed, validCreateInstanceData error:%v, objID:%s, data:%#v, rid:%s", err,
			objID, inputParam.Data, rid)
		return nil, err
	}

	id, err := m.save(kit, objID, inputParam.Data)
	if err != nil {
		blog.ErrorJSON("CreateModelInstance failed, save error:%v, objID:%s, data:%s, rid:%s",
			err, objID, inputParam.Data, kit.Rid)
		return nil, err
	}

	// attach the instance with the quoted instances
	if err = m.dependent.AttachQuotedInst(kit, objID, id, tableData); err != nil {
		return nil, err
	}

	return &metadata.CreateOneDataResult{Created: metadata.CreatedDataResult{ID: id}}, err
}

// CreateManyModelInstance create model instances
func (m *instanceManager) CreateManyModelInstance(kit *rest.Kit,
	objID string, inputParam metadata.CreateManyModelInstance) (
	*metadata.CreateManyDataResult, error) {

	dataResult := new(metadata.CreateManyDataResult)
	if len(inputParam.Datas) == 0 {
		return dataResult, nil
	}

	instValidators, err := m.getValidatorsFromInstances(kit, objID, inputParam.Datas, common.ValidCreate)
	if err != nil {
		blog.Errorf("get inst(%#v) validators failed, err: %v, obj: %s, rid:%s", err, objID, inputParam.Datas, kit.Rid)
		return nil, err
	}

	for index, item := range inputParam.Datas {
		if item == nil {
			blog.ErrorJSON("the model instance data can't be empty, input data: %s rid: %s", inputParam.Datas, kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrCommInstDataNil, "modelInstance")
		}
		item.Set(common.BKOwnerIDField, kit.SupplierAccount)

		validator := instValidators[index]
		if validator == nil {
			blog.Errorf("get validator failed, objID: %s, inst: %#v, rid: %s", err, objID, item, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommNotFound)
		}

		err = m.validCreateInstanceData(kit, objID, item, validator)
		if err != nil {
			blog.Errorf("valid create instance data(%#v) failed, err: %v, obj: %s, rid: %s", err, item, objID, kit.Rid)
			// 由于此err返回的类型可能是mongo返回的error，也可能是经过转化之后的CCError，当返回值是mongo返回的error的场景下没有
			// GetCode方法。
			var errCode int64
			if errInfo, ok := err.(errors.CCErrorCoder); ok {
				errCode = int64(errInfo.GetCode())
			} else {
				errCode = common.CCErrorUnknownOrUnrecognizedError
			}

			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        errCode,
				Data:        item,
				OriginIndex: int64(index),
			})
			continue
		}

		id, err := m.save(kit, objID, item)
		if err != nil {
			blog.Errorf("create instance failed, err: %v, objID: %s, item: %#v, rid: %s", err, objID, item, kit.Rid)
			if id != 0 {
				dataResult.Repeated = append(dataResult.Repeated, metadata.RepeatedDataResult{
					Data:        mapstr.MapStr{"err_msg": err.Error()},
					OriginIndex: int64(index),
				})
			} else {
				var errCode int64
				if errInfo, ok := err.(errors.CCErrorCoder); ok {
					errCode = int64(errInfo.GetCode())
				} else {
					errCode = common.CCErrorUnknownOrUnrecognizedError
				}

				dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
					Message:     err.Error(),
					Code:        errCode,
					Data:        item,
					OriginIndex: int64(index),
				})
			}
			continue
		}

		dataResult.Created = append(dataResult.Created, metadata.CreatedDataResult{
			ID:          id,
			OriginIndex: int64(index),
		})
	}

	return dataResult, nil
}

// BatchCreateModelInstance batch create model instance, if one of instances fails to create, an error is returned.
func (m *instanceManager) BatchCreateModelInstance(kit *rest.Kit, objID string,
	inputParam *metadata.BatchCreateModelInstOption) (*metadata.BatchCreateInstRespData, error) {

	result := new(metadata.BatchCreateInstRespData)
	if len(inputParam.Data) == 0 {
		return result, nil
	}

	instValidators, err := m.getValidatorsFromInstances(kit, objID, inputParam.Data, common.ValidCreate)
	if err != nil {
		blog.Errorf("get inst(%#v) validators failed, err: %v, obj: %s, data: %v, rid: %s", err, objID, inputParam.Data,
			kit.Rid)
		return nil, err
	}

	for idx := range inputParam.Data {
		if inputParam.Data[idx] == nil {
			blog.ErrorJSON("the model instance data can't be empty, input data: %s, rid: %s", inputParam.Data, kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrCommInstDataNil, "modelInstance")
		}
		inputParam.Data[idx].Set(common.BKOwnerIDField, kit.SupplierAccount)

		validator := instValidators[idx]
		if validator == nil {
			blog.Errorf("get validator failed, objID: %s, inst: %#v, rid: %s", objID, inputParam.Data[idx], kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommNotFound)
		}
		err = m.validCreateInstanceData(kit, objID, inputParam.Data[idx], validator)
		if err != nil {
			blog.Errorf("valid create instance data(%#v) failed, err: %v, obj: %s, rid: %s", err, inputParam.Data[idx],
				objID, kit.Rid)
			return nil, err
		}
	}

	ids, err := m.batchSave(kit, objID, inputParam.Data)
	if err != nil {
		return nil, err
	}
	result.IDs = make([]int64, len(ids))
	for idx := range ids {
		result.IDs[idx] = int64(ids[idx])
	}

	return result, nil
}

// UpdateModelInstance update model instances
func (m *instanceManager) UpdateModelInstance(kit *rest.Kit, objID string, inputParam metadata.UpdateOption) (
	*metadata.UpdatedCount, error) {

	inputParam.Condition = util.SetModOwner(inputParam.Condition, kit.SupplierAccount)
	origins, _, err := m.getInsts(kit, objID, inputParam.Condition)
	if err != nil {
		blog.Errorf("get inst failed, err: %v, objID: %s, data: %#v, rid: %s", err, objID, inputParam.Data, kit.Rid)
		return nil, err
	}

	if len(origins) == 0 {
		blog.Errorf("Update model instance failed, no instance found, condition: %#v, objID: %s, data: %#v, rid: %s",
			inputParam.Condition, objID, inputParam.Data, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommNotFound)
	}

	instValidators, err := m.getValidatorsFromInstances(kit, objID, origins, common.ValidUpdate)
	if err != nil {
		blog.Errorf("get inst validators failed, err: %v, objID: %s, data: %#v, rid:%s", err, objID, origins, kit.Rid)
		return nil, err
	}

	isMainline, err := m.isMainlineObject(kit, objID)
	if err != nil {
		return nil, err
	}

	for index, origin := range origins {
		validator := instValidators[index]
		if validator == nil {
			blog.Errorf("get validator failed, objID: %s, inst: %#v, rid: %s", err, objID, origin, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommNotFound)
		}

		// it is not allowed to update multiple records if the updateData has a unique field
		if index == 0 && len(origins) > 1 {
			if err := validator.validUpdateUniqFieldInMulti(kit, inputParam.Data, m); err != nil {
				blog.Errorf("update unique field in multiple records, err: %v, updateData: %#v, rid:%s",
					err, inputParam.Data, kit.Rid)
				return nil, err
			}
		}

		if err := hooks.UpdateProcessBindInfoHook(kit, objID, origin, inputParam.Data); err != nil {
			return nil, err
		}

		if err = m.validUpdateInstanceData(kit, objID, inputParam.Data, origin, validator, inputParam.CanEditAll,
			isMainline); err != nil {
			blog.Errorf("update instance validation failed, err: %v, objID: %s, update data: %#v, inst: %#v, rid: %s",
				err, objID, inputParam.Data, origin, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, err.Error())
		}
	}

	err = m.update(kit, objID, inputParam.Data, inputParam.Condition)
	if err != nil {
		blog.Errorf("update objID(%s) inst failed, err: %v, condition: %#v, data: %#v rid: %s", objID, err,
			inputParam.Condition, inputParam.Data, kit.Rid)
		return nil, err
	}

	if objID == common.BKInnerObjIDHost {
		if err := m.updateHostProcessBindIP(kit, inputParam.Data, origins); err != nil {
			return nil, err
		}
	}

	return &metadata.UpdatedCount{Count: uint64(len(origins))}, nil
}

// updateHostProcessBindIP if hosts' ips are updated, update processes which binds the changed ip
func (m *instanceManager) updateHostProcessBindIP(kit *rest.Kit, data mapstr.MapStr, origins []mapstr.MapStr) error {
	updatedHostFirstIPMap := make(map[string]string)
	for _, hostField := range metadata.ProcBindIPHostFieldMap {
		ip, ipExist := data[hostField]
		if ipExist {
			updatedHostFirstIPMap[hostField] = getFirstIP(ip)
		}
	}

	// get all hosts whose first ip changes
	hostIDs := make([]int64, 0)

	for _, origin := range origins {
		for field, updatedFirstIP := range updatedHostFirstIPMap {
			if getFirstIP(origin[field]) == updatedFirstIP {
				continue
			}

			hostID, err := util.GetInt64ByInterface(origin[common.BKHostIDField])
			if err != nil {
				blog.Errorf("host ID invalid, err: %v, host: %+v, rid: %s", err, origin, kit.Rid)
				return err
			}
			hostIDs = append(hostIDs, hostID)
			break
		}
	}

	if len(hostIDs) == 0 {
		return nil
	}

	// get process templates and their relation to process ids
	procTempMap, processTemplates, err := m.getUpdateIPHostProcTempInfo(kit, hostIDs)
	if err != nil {
		return err
	}

	if len(procTempMap) == 0 || len(processTemplates) == 0 {
		return nil
	}

	// update process whose process template has the bind ip that is changed
	for _, processTemplate := range processTemplates {
		updateProc := make(map[string]interface{})

		for index, value := range processTemplate.Property.BindInfo.Value {
			if value.Std == nil {
				continue
			}

			ip := value.Std.IP
			if !metadata.IsAsDefaultValue(ip.AsDefaultValue) {
				continue
			}

			if ip.Value != nil {
				updatedIP, exist := updatedHostFirstIPMap[metadata.ProcBindIPHostFieldMap[*ip.Value]]
				if exist {
					updateProc[common.BKProcBindInfo+"."+strconv.Itoa(index)+"."+common.BKIP] = updatedIP
				}
			}
		}

		if len(updateProc) != 0 {
			if err := m.updateProcessBindIP(kit, updateProc, procTempMap[processTemplate.ID]); err != nil {
				blog.Errorf("update process bind ip failed, err: %v, rid: %s", err, kit.Rid)
				return err
			}
		}
	}

	return nil
}

func getFirstIP(ip interface{}) string {
	switch t := ip.(type) {
	case string:
		index := strings.Index(t, ",")
		if index == -1 {
			return t
		}

		return t[:index]
	case []string:
		if len(t) == 0 {
			return ""
		}

		return t[0]
	case []interface{}:
		if len(t) == 0 {
			return ""
		}

		return util.GetStrByInterface(t[0])
	}
	return util.GetStrByInterface(ip)
}

// getUpdateIPHostProcTempInfo get update ip host related process and template info to update process bind ip
func (m *instanceManager) getUpdateIPHostProcTempInfo(kit *rest.Kit, hostIDs []int64) (map[int64][]int64,
	[]metadata.ProcessTemplate, error) {

	// get hosts related process and template relations
	procRelations := make([]metadata.ProcessInstanceRelation, 0)
	procRelationFilter := mapstr.MapStr{common.BKHostIDField: mapstr.MapStr{common.BKDBIN: hostIDs}}

	err := mongodb.Client().Table(common.BKTableNameProcessInstanceRelation).Find(procRelationFilter).Fields(
		common.BKHostIDField, common.BKProcessIDField, common.BKProcessTemplateIDField).All(kit.Ctx, &procRelations)
	if err != nil {
		blog.Errorf("get process relation failed, err: %v, hostIDs: %+v, rid: %s", err, hostIDs, kit.Rid)
		return nil, nil, err
	}

	if len(procRelations) == 0 {
		return make(map[int64][]int64), make([]metadata.ProcessTemplate, 0), nil
	}

	procTemplateIDs := make([]int64, len(procRelations))
	procTempMap := make(map[int64][]int64)
	for index, relation := range procRelations {
		procTemplateIDs[index] = relation.ProcessTemplateID
		procTempMap[relation.ProcessTemplateID] = append(procTempMap[relation.ProcessTemplateID], relation.ProcessID)
	}

	// get all processes whose template has corresponding bind info
	processTemplates := make([]metadata.ProcessTemplate, 0)
	processTemplateFilter := map[string]interface{}{
		common.BKFieldID:                      map[string]interface{}{common.BKDBIN: procTemplateIDs},
		"property.bind_info.as_default_value": true,
	}

	err = mongodb.Client().Table(common.BKTableNameProcessTemplate).Find(processTemplateFilter).Fields(
		common.BKFieldID, "property.bind_info").All(kit.Ctx, &processTemplates)
	if err != nil {
		blog.Errorf("get process template failed, ids: %+v, err: %v, rid: %s", err, procTemplateIDs, kit.Rid)
		return nil, nil, err
	}

	return procTempMap, processTemplates, nil
}

// updateProcessBindIP update processes using changed ip
func (m *instanceManager) updateProcessBindIP(kit *rest.Kit, data map[string]interface{}, processIDs []int64) error {
	processFilter := map[string]interface{}{common.BKProcessIDField: map[string]interface{}{common.BKDBIN: processIDs}}

	if err := mongodb.Client().Table(common.BKTableNameBaseProcess).Update(kit.Ctx, processFilter, data); err != nil {
		blog.Errorf("update process failed, err: %v, processIDs: %+v, data: %+v, rid: %s", err, processIDs, data,
			kit.Rid)
		return err
	}

	return nil
}

// SearchModelInstance search model instance
func (m *instanceManager) SearchModelInstance(kit *rest.Kit, objID string, inputParam metadata.QueryCondition) (
	*metadata.QueryResult, error) {

	blog.V(9).Infof("search instance with parameter: %+v, rid: %s", inputParam, kit.Rid)

	tableName := common.GetInstTableName(objID, kit.SupplierAccount)
	if common.IsObjectInstShardingTable(tableName) {
		if inputParam.Condition == nil {
			inputParam.Condition = mapstr.MapStr{}
		}
		objIDCond, ok := inputParam.Condition[common.BKObjIDField]
		if ok && objIDCond != objID {
			blog.V(9).Infof("searchInstance condition's bk_obj_id: %s not match objID: %s, rid: %s", objIDCond, objID,
				kit.Rid)
			return nil, nil
		}
		inputParam.Condition[common.BKObjIDField] = objID
	}
	inputParam.Condition = util.SetQueryOwner(inputParam.Condition, kit.SupplierAccount)

	if inputParam.TimeCondition != nil {
		var err error
		inputParam.Condition, err = inputParam.TimeCondition.MergeTimeCondition(inputParam.Condition)
		if err != nil {
			blog.ErrorJSON("merge time condition failed, error: %s, input: %s, rid: %s", err, inputParam, kit.Rid)
			return nil, err
		}
	}

	// parse vip fields for processes
	fields, vipFields := hooks.ParseVIPFieldsForProcessHook(inputParam.Fields, tableName)

	return m.searchModelInstance(kit, objID, inputParam, fields, vipFields)
}

func (m *instanceManager) searchModelInstance(kit *rest.Kit, objID string, inputParam metadata.QueryCondition,
	fields []string, vipFields []string) (*metadata.QueryResult, error) {

	tableName := common.GetInstTableName(objID, kit.SupplierAccount)
	instItems := make([]mapstr.MapStr, 0)
	query := mongodb.Client().Table(tableName).Find(inputParam.Condition).Start(uint64(inputParam.Page.Start)).
		Limit(uint64(inputParam.Page.Limit)).Sort(inputParam.Page.Sort)

	instItems, instErr := FindInst(kit, fields, query, objID)
	if instErr != nil {
		blog.Errorf("search instance failed, err: %v, rid: %s", instErr, kit.Rid)
		return nil, instErr
	}

	var finalCount uint64
	if !inputParam.DisableCounter {
		count, err := m.countInstance(kit, objID, inputParam.Condition)
		if err != nil {
			blog.Errorf("search model instances count err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}
		finalCount = count
	}

	// set vip info for processes
	instItems, instErr = hooks.SetVIPInfoForProcessHook(kit, instItems, vipFields, tableName, mongodb.Client())
	if instErr != nil {
		return nil, instErr
	}

	dataResult := &metadata.QueryResult{
		Count: finalCount,
		Info:  instItems,
	}

	return dataResult, nil
}

// CountModelInstances counts target model instances num.
func (m *instanceManager) CountModelInstances(kit *rest.Kit,
	objID string, input *metadata.Condition) (*metadata.CommonCountResult, error) {

	if len(objID) == 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKObjIDField)
	}

	if input.TimeCondition != nil {
		var err error
		input.Condition, err = input.TimeCondition.MergeTimeCondition(input.Condition)
		if err != nil {
			blog.Errorf("merge time condition failed, error: %v, input: %s, rid: %s", err, input, kit.Rid)
			return nil, err
		}
	}

	count, err := m.countInstance(kit, objID, input.Condition)
	if err != nil {
		blog.Errorf("count model instances failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, err
	}
	result := &metadata.CommonCountResult{Count: count}

	return result, nil
}

// DeleteModelInstance TODO
func (m *instanceManager) DeleteModelInstance(kit *rest.Kit, objID string,
	inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {
	instIDs := make([]int64, 0)
	tableName := common.GetInstTableName(objID, kit.SupplierAccount)
	instIDFieldName := common.GetInstIDField(objID)

	inputParam.Condition.Set(common.BKOwnerIDField, kit.SupplierAccount)
	inputParam.Condition = util.SetModOwner(inputParam.Condition, kit.SupplierAccount)

	origins, _, err := m.getInsts(kit, objID, inputParam.Condition)
	if nil != err {
		return &metadata.DeletedCount{}, err
	}

	for _, origin := range origins {
		instID, err := util.GetInt64ByInterface(origin[instIDFieldName])
		if nil != err {
			return nil, err
		}
		instIDs = append(instIDs, instID)

		exists, err := m.dependent.IsInstAsstExist(kit, objID, uint64(instID))
		if nil != err {
			return nil, err
		}
		if exists {
			return &metadata.DeletedCount{}, kit.CCError.Error(common.CCErrorInstHasAsst)
		}
	}

	// delete object instance data.
	err = mongodb.Client().Table(tableName).Delete(kit.Ctx, inputParam.Condition)
	if nil != err {
		blog.ErrorJSON("DeleteModelInstance delete objID(%s) instance error. err:%s, coniditon:%s, rid:%s", objID,
			err.Error(), inputParam.Condition, kit.Rid)
		return &metadata.DeletedCount{}, err
	}

	// delete object instance mapping.
	if metadata.IsCommon(objID) {
		if err := instancemapping.Delete(kit.Ctx, instIDs); err != nil {
			blog.Errorf("delete object %s instance mapping failed, err: %s, instance: %v, rid: %s",
				objID, err.Error(), instIDs, kit.Rid)
			return nil, err
		}
	}

	// delete these instances' quoted instances
	if err = m.dependent.DeleteQuotedInst(kit, objID, instIDs); err != nil {
		return nil, err
	}

	return &metadata.DeletedCount{Count: uint64(len(origins))}, nil
}

// CascadeDeleteModelInstance TODO
func (m *instanceManager) CascadeDeleteModelInstance(kit *rest.Kit, objID string,
	inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {
	instIDs := make([]int64, 0)
	tableName := common.GetInstTableName(objID, kit.SupplierAccount)
	instIDFieldName := common.GetInstIDField(objID)

	origins, _, err := m.getInsts(kit, objID, inputParam.Condition)
	if nil != err {
		blog.Errorf("cascade delete model instance get inst error:%v, rid: %s", err, kit.Rid)
		return &metadata.DeletedCount{}, err
	}

	for _, origin := range origins {
		instID, err := util.GetInt64ByInterface(origin[instIDFieldName])
		if nil != err {
			return &metadata.DeletedCount{}, err
		}
		instIDs = append(instIDs, instID)

		err = m.dependent.DeleteInstAsst(kit, objID, uint64(instID))
		if nil != err {
			return &metadata.DeletedCount{}, err
		}
	}

	// delete object instance data.
	inputParam.Condition = util.SetModOwner(inputParam.Condition, kit.SupplierAccount)
	err = mongodb.Client().Table(tableName).Delete(kit.Ctx, inputParam.Condition)
	if nil != err {
		return &metadata.DeletedCount{}, err
	}

	// delete object instance mapping.
	if metadata.IsCommon(objID) {
		if err := instancemapping.Delete(kit.Ctx, instIDs); err != nil {
			blog.Errorf("delete object %s instance mapping failed, err: %s, instance: %v, rid: %s",
				objID, err.Error(), instIDs, kit.Rid)
			return nil, err
		}
	}

	// delete these instances' quoted instances
	if err = m.dependent.DeleteQuotedInst(kit, objID, instIDs); err != nil {
		return nil, err
	}

	return &metadata.DeletedCount{Count: uint64(len(origins))}, nil
}

// FindInst find instance
func FindInst(kit *rest.Kit, fields []string, query types.Find, objID string) ([]mapstr.MapStr, error) {
	existCreateAt := true
	existCreateTime := true
	existUpdateAt := true
	existLastTime := true

	if len(fields) != 0 {
		// 旧数据用create_time和last_time分别记录了实例的创建和更新时间，如果bk_created_at和bk_updated_at字段没值，需要把旧值赋过来
		fieldMap := make(map[string]struct{})
		for _, field := range fields {
			fieldMap[field] = struct{}{}
		}

		_, existCreateAt = fieldMap[common.BKCreatedAt]
		_, existCreateTime = fieldMap[common.CreateTimeField]
		if existCreateAt && !existCreateTime {
			fields = append(fields, common.CreateTimeField)
		}

		_, existUpdateAt = fieldMap[common.BKUpdatedAt]
		_, existLastTime = fieldMap[common.LastTimeField]
		if existUpdateAt && !existLastTime {
			fields = append(fields, common.LastTimeField)
		}
	}

	query.Fields(fields...)

	insts := make([]mapstr.MapStr, 0)
	var err error
	if objID == common.BKInnerObjIDHost {
		hosts := make([]metadata.HostMapStr, 0)
		err = query.All(kit.Ctx, &hosts)
		for _, host := range hosts {
			insts = append(insts, mapstr.MapStr(host))
		}
	} else {
		err = query.All(kit.Ctx, &insts)
	}

	if err != nil {
		blog.Errorf("search instance failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	for idx, inst := range insts {
		if existCreateAt && inst[common.BKCreatedAt] == nil {
			inst[common.BKCreatedAt] = inst[common.CreateTimeField]
		}

		if existUpdateAt && inst[common.BKUpdatedAt] == nil {
			inst[common.BKUpdatedAt] = inst[common.LastTimeField]
		}

		if !existCreateTime {
			inst.Remove(common.CreateTimeField)
		}

		if !existLastTime {
			inst.Remove(common.LastTimeField)
		}

		insts[idx] = inst
	}

	return insts, nil
}
