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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/thirdparty/hooks"
)

var _ core.InstanceOperation = (*instanceManager)(nil)

type instanceManager struct {
	dependent OperationDependences
	language  language.CCLanguageIf
	clientSet apimachinery.ClientSetInterface
}

// New create a new instance manager instance
func New(dependent OperationDependences, language language.CCLanguageIf, clientSet apimachinery.ClientSetInterface) core.InstanceOperation {
	return &instanceManager{
		dependent: dependent,
		language:  language,
		clientSet: clientSet,
	}
}

func (m *instanceManager) instCnt(kit *rest.Kit, objID string, cond mapstr.MapStr) (cnt uint64, exists bool, err error) {
	tableName := common.GetInstTableName(objID)
	cnt, err = mongodb.Client().Table(tableName).Find(cond).Count(kit.Ctx)
	exists = 0 != cnt
	return cnt, exists, err
}

func (m *instanceManager) CreateModelInstance(kit *rest.Kit, objID string, inputParam metadata.CreateModelInstance) (*metadata.CreateOneDataResult, error) {
	rid := util.ExtractRequestIDFromContext(kit.Ctx)

	inputParam.Data.Set(common.BKOwnerIDField, kit.SupplierAccount)
	bizID, err := m.getBizIDFromInstance(kit, objID, inputParam.Data, common.ValidCreate, 0)
	if err != nil {
		blog.Errorf("CreateModelInstance failed, getBizIDFromInstance err:%v, objID:%s, data:%#v, rid:%s", err, objID, inputParam.Data, kit.Rid)
		return nil, err
	}
	validator, err := m.newValidator(kit, objID, bizID)
	if err != nil {
		blog.Errorf("CreateModelInstance failed, newValidator err:%v, objID:%s, data:%#v, rid:%s", err, objID, inputParam.Data, kit.Rid)
		return nil, err
	}

	err = m.validCreateInstanceData(kit, objID, inputParam.Data, validator)
	if nil != err {
		blog.Errorf("CreateModelInstance failed, validCreateInstanceData error:%v, objID:%s, data:%#v, rid:%s", err, objID, inputParam.Data, rid)
		return nil, err
	}

	id, err := m.save(kit, objID, inputParam.Data)
	if err != nil {
		blog.ErrorJSON("CreateModelInstance failed, save error:%v, objID:%s, data:%#v, rid:%s", err, objID, inputParam.Data, kit.Rid)
		return nil, err
	}

	return &metadata.CreateOneDataResult{Created: metadata.CreatedDataResult{ID: id}}, err
}

func (m *instanceManager) CreateManyModelInstance(kit *rest.Kit, objID string, inputParam metadata.CreateManyModelInstance) (*metadata.CreateManyDataResult, error) {
	dataResult := &metadata.CreateManyDataResult{}
	allValidators := make(map[int64]*validator)
	for _, item := range inputParam.Datas {
		if item == nil {
			blog.ErrorJSON("the model instance data can't be empty, input data:%s rid:%s", inputParam.Datas, kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrCommInstDataNil, "modelInstance")
		}
		item.Set(common.BKOwnerIDField, kit.SupplierAccount)
		bizID, err := m.getBizIDFromInstance(kit, objID, item, common.ValidCreate, 0)
		if err != nil {
			blog.Errorf("CreateManyModelInstance failed, getBizIDFromInstance err:%v, objID:%s, data:%#v, rid:%s", err, objID, item, kit.Rid)
			return nil, err
		}
		if allValidators[bizID] == nil {
			validator, err := m.newValidator(kit, objID, bizID)
			if err != nil {
				blog.Errorf("CreateManyModelInstance failed, newValidator err:%v, objID:%s, bizID:%d, rid:%s", err, objID, bizID, kit.Rid)
				return nil, err
			}
			allValidators[bizID] = validator
		}

		err = m.validCreateInstanceData(kit, objID, item, allValidators[bizID])
		if nil != err {
			blog.Errorf("CreateManyModelInstance failed, validCreateInstanceData err:%v, objID:%s, item:%#v, rid:%s", err, objID, item, kit.Rid)
			return nil, err
		}

		id, err := m.save(kit, objID, item)
		if nil != err {
			blog.Errorf("CreateManyModelInstance failed, save err:%v, objID:%s, item:%#v, rid:%s", err, objID, item, kit.Rid)
			return nil, err
		}

		dataResult.Created = append(dataResult.Created, metadata.CreatedDataResult{
			ID: id,
		})
	}

	return dataResult, nil
}

func (m *instanceManager) UpdateModelInstance(kit *rest.Kit, objID string, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error) {
	instIDFieldName := common.GetInstIDField(objID)
	inputParam.Condition = util.SetModOwner(inputParam.Condition, kit.SupplierAccount)
	origins, _, err := m.getInsts(kit, objID, inputParam.Condition)
	if nil != err {
		blog.Errorf("UpdateModelInstance failed, get inst failed, err:%v, objID:%s, data:%#v, rid:%s", err, objID, inputParam.Data, kit.Rid)
		return nil, err
	}

	if len(origins) == 0 {
		blog.Errorf("UpdateModelInstance failed, no instance found, condition:%#v, objID:%s, data:%#v, rid:%s",
			inputParam.Condition, objID, inputParam.Data, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommNotFound)
	}

	allValidators := make(map[int64]*validator)
	for idx, origin := range origins {
		instIDI := origin[instIDFieldName]
		instID, _ := util.GetInt64ByInterface(instIDI)
		bizID, err := m.getBizIDFromInstance(kit, objID, origin, common.ValidUpdate, instID)
		if err != nil {
			blog.Errorf("UpdateModelInstance failed, getBizIDFromInstance err:%v, objID:%s, data:%#v, rid:%s", err, objID, origin, kit.Rid)
			return nil, err
		}
		if allValidators[bizID] == nil {
			validator, err := m.newValidator(kit, objID, bizID)
			if err != nil {
				blog.Errorf("UpdateModelInstance failed, newValidator err:%v, objID:%s, bizID:%d, rid:%s", err, objID, bizID, kit.Rid)
				return nil, err
			}
			allValidators[bizID] = validator
		}

		// it is not allowed to update multiple records if the updateData has a unique field
		if idx == 0 && len(origins) > 1 {
			valid := allValidators[bizID]
			if err := valid.validUpdateUniqFieldInMulti(kit, inputParam.Data, m); err != nil {
				blog.Errorf("UpdateModelInstance failed, validUpdateUniqFieldInMulti error:%v, updateData: %#v, rid:%s",
					err, inputParam.Data, kit.Rid)
				return nil, err
			}
		}

		err = m.validUpdateInstanceData(kit, objID, inputParam.Data, origin, allValidators[bizID], instID, inputParam.CanEditAll)
		if nil != err {
			blog.Errorf("update model instance validate error:%v, objID:%s, updateData: %#v, instData:%#v, rid:%s",
				err, objID, inputParam.Data, origin, kit.Rid)
			return nil, err
		}
	}

	err = m.update(kit, objID, inputParam.Data, inputParam.Condition)
	if err != nil {
		blog.ErrorJSON("UpdateModelInstance update objID(%s) inst error. err:%s, data:%#v, condition:%s, rid:%s",
			objID, err, inputParam.Condition, inputParam.Data, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommDBUpdateFailed)
	}

	if objID == common.BKInnerObjIDHost {
		if err := m.updateHostProcessBindIP(kit, inputParam.Data, origins); err != nil {
			return nil, err
		}
	}

	return &metadata.UpdatedCount{Count: uint64(len(origins))}, nil
}

// updateHostProcessBindIP if hosts' ips are updated, update processes which binds the changed ip
func (m *instanceManager) updateHostProcessBindIP(kit *rest.Kit, updateData mapstr.MapStr, origins []mapstr.MapStr) error {
	innerIP, innerIPExist := updateData[common.BKHostInnerIPField]
	outerIP, outerIPExist := updateData[common.BKHostOuterIPField]

	firstInnerIP := getFirstIP(innerIP)
	firstOuterIP := getFirstIP(outerIP)

	// get all hosts whose first ip changes
	innerIPUpdatedHostMap := make(map[int64]bool)
	outerIPUpdatedHostMap := make(map[int64]bool)
	hostIDs := make([]int64, 0)
	var err error

	for _, origin := range origins {
		var hostID int64

		if innerIPExist && getFirstIP(origin[common.BKHostInnerIPField]) != firstInnerIP {
			hostID, err = util.GetInt64ByInterface(origin[common.BKHostIDField])
			if err != nil {
				blog.Errorf("host ID invalid, err: %v, host: %+v, rid: %s", err, origin, kit.Rid)
				return err
			}
			innerIPUpdatedHostMap[hostID] = true
		}

		if outerIPExist && getFirstIP(origin[common.BKHostOuterIPField]) != firstOuterIP {
			if hostID == 0 {
				hostID, err = util.GetInt64ByInterface(origin[common.BKHostIDField])
				if err != nil {
					blog.Errorf("host ID invalid, err: %v, host: %+v, rid: %s", err, origin, kit.Rid)
					return err
				}
			}
			outerIPUpdatedHostMap[hostID] = true
		}

		if hostID != 0 {
			hostIDs = append(hostIDs, hostID)
		}
	}

	if len(hostIDs) == 0 {
		return nil
	}

	// get hosts related process and template relations
	processRelations := make([]metadata.ProcessInstanceRelation, 0)
	processRelationFilter := map[string]interface{}{common.BKHostIDField: map[string]interface{}{common.BKDBIN: hostIDs}}

	err = mongodb.Client().Table(common.BKTableNameProcessInstanceRelation).Find(processRelationFilter).Fields(
		common.BKHostIDField, common.BKProcessIDField, common.BKProcessTemplateIDField).All(kit.Ctx, &processRelations)
	if err != nil {
		blog.Errorf("get process relation failed, err: %v, hostIDs: %+v, rid: %s", err, hostIDs, kit.Rid)
		return err
	}

	if len(processRelations) == 0 {
		return nil
	}

	processTemplateIDs := make([]int64, len(processRelations))
	processTemplateMap := make(map[int64][]int64)
	for index, relation := range processRelations {
		processTemplateIDs[index] = relation.ProcessTemplateID
		processTemplateMap[relation.ProcessTemplateID] = append(processTemplateMap[relation.ProcessTemplateID], relation.ProcessID)
	}

	// get all processes whose templates has corresponding bind ip
	processTemplates := make([]metadata.ProcessTemplate, 0)
	processTemplateFilter := map[string]interface{}{
		common.BKFieldID:                      map[string]interface{}{common.BKDBIN: processTemplateIDs},
		"property.bind_info.as_default_value": true,
	}

	err = mongodb.Client().Table(common.BKTableNameProcessTemplate).Find(processTemplateFilter).Fields(
		common.BKFieldID, "property.bind_info").All(kit.Ctx, &processTemplates)
	if err != nil {
		blog.Errorf("get process template failed, err: %v, processTemplateIDs: %+v, rid: %s", err, processTemplateIDs, kit.Rid)
		return err
	}

	for _, processTemplate := range processTemplates {
		data := make(map[string]interface{})

		for index, value := range processTemplate.Property.BindInfo.Value {
			if value.Std == nil {
				continue
			}

			ip := value.Std.IP
			if !metadata.IsAsDefaultValue(ip.AsDefaultValue) {
				continue
			}

			if ip.Value != nil {
				if innerIPExist && *ip.Value == metadata.BindInnerIP {
					data[common.BKProcBindInfo+"."+strconv.Itoa(index)+"."+common.BKIP] = firstInnerIP
				}

				if outerIPExist && *ip.Value == metadata.BindOuterIP {
					data[common.BKProcBindInfo+"."+strconv.Itoa(index)+"."+common.BKIP] = firstOuterIP
				}
			}
		}

		if len(data) != 0 {
			if err := m.updateProcessBindIP(kit, data, processTemplateMap[processTemplate.ID]); err != nil {
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

// updateHostProcessBindIP update processes using changed ip
func (m *instanceManager) updateProcessBindIP(kit *rest.Kit, data map[string]interface{}, processIDs []int64) error {
	processFilter := map[string]interface{}{common.BKProcessIDField: map[string]interface{}{common.BKDBIN: processIDs}}

	if err := mongodb.Client().Table(common.BKTableNameBaseProcess).Update(kit.Ctx, processFilter, data); err != nil {
		blog.Errorf("update process failed, err: %v, processIDs: %+v, data: %+v, rid: %s", err, processIDs, data, kit.Rid)
		return err
	}

	return nil
}

func (m *instanceManager) SearchModelInstance(kit *rest.Kit, objID string, inputParam metadata.QueryCondition) (*metadata.QueryResult, error) {
	blog.V(9).Infof("search instance with parameter: %+v, rid: %s", inputParam, kit.Rid)

	tableName := common.GetInstTableName(objID)
	if tableName == common.BKTableNameBaseInst {
		if inputParam.Condition == nil {
			inputParam.Condition = mapstr.MapStr{}
		}
		objIDCond, ok := inputParam.Condition[common.BKObjIDField]
		if ok && objIDCond != objID {
			blog.V(9).Infof("searchInstance condition's bk_obj_id: %s not match objID: %s, rid: %s", objIDCond, objID, kit.Rid)
			return nil, nil
		}
		inputParam.Condition[common.BKObjIDField] = objID
	}
	inputParam.Condition = util.SetQueryOwner(inputParam.Condition, kit.SupplierAccount)

	// parse vip fields for processes
	fields, vipFields := hooks.ParseVIPFieldsForProcessHook(inputParam.Fields, tableName)

	instItems := make([]mapstr.MapStr, 0)
	query := mongodb.Client().Table(tableName).Find(inputParam.Condition).Start(uint64(inputParam.Page.Start)).
		Limit(uint64(inputParam.Page.Limit)).
		Sort(inputParam.Page.Sort).
		Fields(fields...)
	var instErr error
	if objID == common.BKInnerObjIDHost {
		hosts := make([]metadata.HostMapStr, 0)
		instErr = query.All(kit.Ctx, &hosts)
		for _, host := range hosts {
			instItems = append(instItems, mapstr.MapStr(host))
		}
	} else {
		instErr = query.All(kit.Ctx, &instItems)
	}
	if instErr != nil {
		blog.Errorf("search instance error [%v], rid: %s", instErr, kit.Rid)
		return nil, instErr
	}

	var finalCount uint64

	if !inputParam.DisableCounter {
		count, countErr := mongodb.Client().Table(tableName).Find(inputParam.Condition).Count(kit.Ctx)
		if countErr != nil {
			blog.Errorf("count instance error [%v], rid: %s", countErr, kit.Rid)
			return nil, countErr
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

func (m *instanceManager) DeleteModelInstance(kit *rest.Kit, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {
	tableName := common.GetInstTableName(objID)
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
		exists, err := m.dependent.IsInstAsstExist(kit, objID, uint64(instID))
		if nil != err {
			return nil, err
		}
		if exists {
			return &metadata.DeletedCount{}, kit.CCError.Error(common.CCErrorInstHasAsst)
		}
	}

	err = mongodb.Client().Table(tableName).Delete(kit.Ctx, inputParam.Condition)
	if nil != err {
		blog.ErrorJSON("DeleteModelInstance delete objID(%s) instance error. err:%s, coniditon:%s, rid:%s", objID, err.Error(), inputParam.Condition, kit.Rid)
		return &metadata.DeletedCount{}, err
	}

	return &metadata.DeletedCount{Count: uint64(len(origins))}, nil
}

func (m *instanceManager) CascadeDeleteModelInstance(kit *rest.Kit, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {
	tableName := common.GetInstTableName(objID)
	instIDFieldName := common.GetInstIDField(objID)
	origins, _, err := m.getInsts(kit, objID, inputParam.Condition)
	blog.V(5).Infof("cascade delete model instance get inst error:%v, rid: %s", origins, kit.Rid)
	if nil != err {
		blog.Errorf("cascade delete model instance get inst error:%v, rid: %s", err, kit.Rid)
		return &metadata.DeletedCount{}, err
	}

	for _, origin := range origins {
		instID, err := util.GetInt64ByInterface(origin[instIDFieldName])
		if nil != err {
			return &metadata.DeletedCount{}, err
		}
		err = m.dependent.DeleteInstAsst(kit, objID, uint64(instID))
		if nil != err {
			return &metadata.DeletedCount{}, err
		}
	}
	inputParam.Condition = util.SetModOwner(inputParam.Condition, kit.SupplierAccount)
	err = mongodb.Client().Table(tableName).Delete(kit.Ctx, inputParam.Condition)
	if nil != err {
		return &metadata.DeletedCount{}, err
	}
	return &metadata.DeletedCount{Count: uint64(len(origins))}, nil
}
