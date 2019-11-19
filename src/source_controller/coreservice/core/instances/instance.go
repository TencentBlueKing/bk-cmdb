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
	redis "gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/eventclient"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

var _ core.InstanceOperation = (*instanceManager)(nil)

type instanceManager struct {
	dbProxy   dal.RDB
	dependent OperationDependences
	validator validator
	Cache     *redis.Client
	EventCli  eventclient.Client
}

// New create a new instance manager instance
func New(dbProxy dal.RDB, dependent OperationDependences, cache *redis.Client) core.InstanceOperation {
	return &instanceManager{
		dbProxy:   dbProxy,
		dependent: dependent,
		EventCli:  eventclient.NewClientViaRedis(cache, dbProxy),
	}
}

func (m *instanceManager) instCnt(ctx core.ContextParams, objID string, cond mapstr.MapStr) (cnt uint64, exists bool, err error) {
	tableName := common.GetInstTableName(objID)
	cnt, err = m.dbProxy.Table(tableName).Find(cond).Count(ctx)
	exists = 0 != cnt
	return cnt, exists, err
}

func (m *instanceManager) CreateModelInstance(ctx core.ContextParams, objID string, inputParam metadata.CreateModelInstance) (*metadata.CreateOneDataResult, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	err := m.validCreateInstanceData(ctx, objID, inputParam.Data)
	if nil != err {
		blog.Errorf("CreateModelInstance failed, valid error: %+v, rid: %s", err, rid)
		return nil, err
	}
	id, err := m.save(ctx, objID, inputParam.Data)
	if err != nil {
		blog.ErrorJSON("CreateModelInstance create objID(%s) instance error. err:%s, data:%s, rid:%s", objID, err.Error(), inputParam.Data, ctx.ReqID)
		return nil, err
	}

	instIDFieldName := common.GetInstIDField(objID)
	// 处理事件数据的
	eh := m.NewEventClient(objID)
	err = eh.SetCurDataAndPush(ctx, objID, metadata.EventActionCreate, mapstr.MapStr{instIDFieldName: id})
	if err != nil {
		blog.ErrorJSON("CreateModelInstance  event push instance current data error. err:%s, objID:%s inst id:%s, rid:%s", err, objID, id, ctx.ReqID)
		return &metadata.CreateOneDataResult{Created: metadata.CreatedDataResult{ID: id}}, err
	}
	return &metadata.CreateOneDataResult{Created: metadata.CreatedDataResult{ID: id}}, err
}

func (m *instanceManager) CreateManyModelInstance(ctx core.ContextParams, objID string, inputParam metadata.CreateManyModelInstance) (*metadata.CreateManyDataResult, error) {
	var newIDs []uint64
	dataResult := &metadata.CreateManyDataResult{}
	for itemIdx, item := range inputParam.Datas {
		item.Set(common.BKOwnerIDField, ctx.SupplierAccount)
		err := m.validCreateInstanceData(ctx, objID, item)
		if nil != err {
			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.(errors.CCErrorCoder).GetCode()),
				Data:        item,
				OriginIndex: int64(itemIdx),
			})
			continue
		}
		item.Set(common.BKOwnerIDField, ctx.SupplierAccount)
		id, err := m.save(ctx, objID, item)
		if nil != err {
			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.(errors.CCErrorCoder).GetCode()),
				Data:        item,
				OriginIndex: int64(itemIdx),
			})
			continue
		}

		dataResult.Created = append(dataResult.Created, metadata.CreatedDataResult{
			ID: id,
		})
		newIDs = append(newIDs, id)

	}
	instIDFieldName := common.GetInstIDField(objID)
	// 处理事件数据的
	eh := m.NewEventClient(objID)
	err := eh.SetCurDataAndPush(ctx, objID, metadata.EventActionCreate, condition.CreateCondition().Field(instIDFieldName).In(newIDs).ToMapStr())
	if err != nil {
		blog.ErrorJSON("CreateManyModelInstance  event push instance current data error. err:%s, objID:%s inst id:%s, rid:%s", err, objID, newIDs, ctx.ReqID)
		return dataResult, err
	}
	return dataResult, nil
}

func (m *instanceManager) UpdateModelInstance(ctx core.ContextParams, objID string, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error) {
	instIDFieldName := common.GetInstIDField(objID)
	inputParam.Condition.Set(common.BKOwnerIDField, ctx.SupplierAccount)
	origins, _, err := m.getInsts(ctx, objID, inputParam.Condition)
	if nil != err {
		blog.Errorf("UpdateModelInstance failed, get inst failed, err: %v, rid:%s", err, ctx.ReqID)
		return nil, err
	}

	if len(origins) == 0 {
		blog.Errorf("UpdateModelInstance failed, no instance found. model: %s, condition:%+v, rid:%s", objID, inputParam.Condition, ctx.ReqID)
		return nil, ctx.Error.Error(common.CCErrCommNotFound)
	}

	// 处理事件数据的
	eh := m.NewEventClient(objID)

	var instMedataData metadata.Metadata
	instMedataData.Label = make(metadata.Label)
	for key, val := range inputParam.Condition {
		if metadata.BKMetadata == key {
			bizID := metadata.GetBusinessIDFromMeta(val)
			if "" != bizID {
				instMedataData.Label.Set(metadata.LabelBusinessID, bizID)
			}
			continue
		}
	}
	if inputParam.Condition.Exists(metadata.BKMetadata) {
		inputParam.Condition.Set(metadata.BKMetadata, instMedataData)
	}

	for _, origin := range origins {
		instIDI := origin[instIDFieldName]
		instID, _ := util.GetInt64ByInterface(instIDI)
		err := m.validUpdateInstanceData(ctx, objID, inputParam.Data, instMedataData, uint64(instID))
		if nil != err {
			blog.Errorf("update model instance validate error :%v ,rid:%s", err, ctx.ReqID)
			return nil, err
		}
		// 设置实例变更前数据
		eh.SetPreData(instID, origin)
	}

	if nil != err {
		blog.Errorf("update model instance validate error :%v ,rid:%s", err, ctx.ReqID)
		return &metadata.UpdatedCount{}, err
	}
	cnt, err := m.update(ctx, objID, inputParam.Data, inputParam.Condition)
	if err != nil {
		blog.ErrorJSON("UpdateModelInstance update objID(%s) inst error. err:%s, condition:%s, rid:%s", objID, inputParam.Condition, ctx.ReqID)
		return nil, err
	}
	err = eh.SetCurDataAndPush(ctx, objID, metadata.EventActionUpdate, inputParam.Condition)
	if err != nil {
		blog.ErrorJSON("UpdateModelInstance  event push instance current data error. err:%s, condition:%s, rid:%s", err, inputParam.Condition, ctx.ReqID)
		return nil, err
	}

	return &metadata.UpdatedCount{Count: cnt}, nil
}

func (m *instanceManager) SearchModelInstance(ctx core.ContextParams, objID string, inputParam metadata.QueryCondition) (*metadata.QueryResult, error) {
	condition, err := mongo.NewConditionFromMapStr(inputParam.Condition)
	if nil != err {
		blog.Errorf("SearchModelInstance failed, parse condition failed, inputParam: %+v, err: %+v, rid: %s", inputParam, err, ctx.ReqID)
		return &metadata.QueryResult{}, err
	}
	ownerIDArr := []string{ctx.SupplierAccount, common.BKDefaultOwnerID}
	condition.Element(&mongo.In{Key: common.BKOwnerIDField, Val: ownerIDArr})
	inputParam.Condition = condition.ToMapStr()

	blog.V(9).Infof("search instance with parameter: %+v, rid: %s", inputParam, ctx.ReqID)
	instItems, err := m.searchInstance(ctx, objID, inputParam)
	if nil != err {
		blog.Errorf("search instance error [%v], rid: %s", err, ctx.ReqID)
		return &metadata.QueryResult{}, err
	}

	dataResult := &metadata.QueryResult{}
	dataResult.Count, err = m.countInstance(ctx, objID, inputParam.Condition)
	if nil != err {
		blog.Errorf("count instance error [%v], rid: %s", err, ctx.ReqID)
		return &metadata.QueryResult{}, err
	}
	dataResult.Info = instItems

	return dataResult, nil
}

func (m *instanceManager) DeleteModelInstance(ctx core.ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {
	tableName := common.GetInstTableName(objID)
	instIDFieldName := common.GetInstIDField(objID)
	inputParam.Condition.Set(common.BKOwnerIDField, ctx.SupplierAccount)
	origins, _, err := m.getInsts(ctx, objID, inputParam.Condition)
	if nil != err {
		return &metadata.DeletedCount{}, err
	}

	// 处理事件数据的
	eh := m.NewEventClient(objID)

	for _, origin := range origins {
		instID, err := util.GetInt64ByInterface(origin[instIDFieldName])
		if nil != err {
			return nil, err
		}
		exists, err := m.dependent.IsInstAsstExist(ctx, objID, uint64(instID))
		if nil != err {
			return nil, err
		}
		if exists {
			return &metadata.DeletedCount{}, ctx.Error.Error(common.CCErrorInstHasAsst)
		}
		eh.SetPreData(instID, origin)
	}
	err = m.dbProxy.Table(tableName).Delete(ctx, inputParam.Condition)
	if nil != err {
		blog.ErrorJSON("DeleteModelInstance delete objID(%s) instance error. err:%s, coniditon:%s, rid:%s", objID, err.Error(), inputParam.Condition, ctx.ReqID)
		return &metadata.DeletedCount{}, err
	}
	err = eh.Push(ctx, objID, metadata.EventActionDelete)
	if err != nil {
		blog.ErrorJSON("DeleteModelInstance push delete objType(%s) instance to event server error. data:%s, rid:%s", objID, origins, ctx.ReqID)
		return &metadata.DeletedCount{Count: uint64(len(origins))}, ctx.Error.CCErrorf(common.CCErrCoreServiceEventPushEventFailed)

	}
	return &metadata.DeletedCount{Count: uint64(len(origins))}, nil
}

func (m *instanceManager) CascadeDeleteModelInstance(ctx core.ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {
	tableName := common.GetInstTableName(objID)
	instIDFieldName := common.GetInstIDField(objID)
	origins, _, err := m.getInsts(ctx, objID, inputParam.Condition)
	blog.V(5).Infof("cascade delete model instance get inst error:%v, rid: %s", origins, ctx.ReqID)
	if nil != err {
		blog.Errorf("cascade delete model instance get inst error:%v, rid: %s", err, ctx.ReqID)
		return &metadata.DeletedCount{}, err
	}

	for _, origin := range origins {
		instID, err := util.GetInt64ByInterface(origin[instIDFieldName])
		if nil != err {
			return &metadata.DeletedCount{}, err
		}
		err = m.dependent.DeleteInstAsst(ctx, objID, uint64(instID))
		if nil != err {
			return &metadata.DeletedCount{}, err
		}
	}
	inputParam.Condition.Set(common.BKOwnerIDField, ctx.SupplierAccount)
	err = m.dbProxy.Table(tableName).Delete(ctx, inputParam.Condition)
	if nil != err {
		return &metadata.DeletedCount{}, err
	}
	return &metadata.DeletedCount{Count: uint64(len(origins))}, nil
}
