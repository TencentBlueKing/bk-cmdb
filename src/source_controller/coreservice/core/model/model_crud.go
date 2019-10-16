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

package model

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
)

func (m *modelManager) count(ctx core.ContextParams, cond universalsql.Condition) (uint64, error) {

	cnt, err := m.dbProxy.Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).Count(ctx)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute database count operation by the condition (%#v), error info is %s", ctx.ReqID, cond.ToMapStr(), err.Error())
		return 0, ctx.Error.Errorf(common.CCErrObjectDBOpErrno, err.Error())
	}

	return cnt, err
}

func (m *modelManager) save(ctx core.ContextParams, model *metadata.Object) (id uint64, err error) {

	id, err = m.dbProxy.NextSequence(ctx, common.BKTableNameObjDes)
	if err != nil {
		blog.Errorf("request(%s): it is failed to make sequence id on the table (%s), error info is %s", ctx.ReqID, common.BKTableNameObjDes, err.Error())
		return id, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}
	model.ID = int64(id)
	model.OwnerID = ctx.SupplierAccount

	if nil == model.LastTime {
		model.LastTime = &metadata.Time{}
		model.LastTime.Time = time.Now()
	}
	if nil == model.CreateTime {
		model.CreateTime = &metadata.Time{}
		model.CreateTime.Time = time.Now()
	}

	err = m.dbProxy.Table(common.BKTableNameObjDes).Insert(ctx, model)
	return id, err
}

func (m *modelManager) update(ctx core.ContextParams, data mapstr.MapStr, cond universalsql.Condition) (cnt uint64, err error) {

	cnt, err = m.count(ctx, cond)
	if nil != err {
		return 0, err
	}

	if 0 == cnt {
		return 0, nil
	}

	data.Set(metadata.ModelFieldLastTime, time.Now())
	models := make([]metadata.Object, 0)
	err = m.dbProxy.Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).All(ctx, &models)
	if nil != err {
		blog.Errorf("find models failed, filter: %+v, err: %s, rid: %s", cond.ToMapStr(), err.Error(), ctx.ReqID)
		return 0, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	if objName, exist := data[common.BKObjNameField]; exist == true && len(util.GetStrByInterface(objName)) > 0 {
		for _, model := range models {
			modelName := data[common.BKObjNameField]

			// 检查模型名称重复
			modelNameUniqueFilter := map[string]interface{}{
				common.BKObjNameField: modelName,
				common.BKFieldID: map[string]interface{}{
					common.BKDBNE: model.ID,
				},
			}
			bizFilter := metadata.PublicAndBizCondition(model.Metadata)
			for key, value := range bizFilter {
				modelNameUniqueFilter[key] = value
			}
			sameNameCount, err := m.dbProxy.Table(common.BKTableNameObjDes).Find(modelNameUniqueFilter).Count(ctx)
			if err != nil {
				blog.Errorf("check whether same name model exists failed, name: %s, filter: %+v, err: %s, rid: %s", modelName, modelNameUniqueFilter, err.Error(), ctx.ReqID)
				return 0, err
			}
			if sameNameCount > 0 {
				blog.Warnf("update model failed, field `%s` duplicated, rid: %s", modelName, ctx.ReqID)
				return 0, ctx.Error.Errorf(common.CCErrCommDuplicateItem, modelName)
			}

			// 一次更新多个模型的时候，唯一校验需要特别小心
			filter := map[string]interface{}{common.BKFieldID: model.ID}
			err = m.dbProxy.Table(common.BKTableNameObjDes).Update(ctx, filter, data)
			if nil != err {
				blog.Errorf("request(%s): it is failed to execute database update operation on the table (%s), error info is %s", ctx.ReqID, common.BKTableNameObjDes, err.Error())
				return 0, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
			}
		}
		return cnt, nil
	}

	err = m.dbProxy.Table(common.BKTableNameObjDes).Update(ctx, cond.ToMapStr(), data)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute database update operation on the table (%s), error info is %s", ctx.ReqID, common.BKTableNameObjDes, err.Error())
		return 0, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return cnt, err
}

func (m *modelManager) search(ctx core.ContextParams, cond universalsql.Condition) ([]metadata.Object, error) {

	dataResult := make([]metadata.Object, 0)
	if err := m.dbProxy.Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).All(ctx, &dataResult); nil != err {
		blog.Errorf("request(%s): it is failed to find all models by the condition (%#v), error info is %s", ctx.ReqID, cond.ToMapStr(), err.Error())
		return dataResult, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return dataResult, nil
}

func (m *modelManager) searchReturnMapStr(ctx core.ContextParams, cond universalsql.Condition) ([]mapstr.MapStr, error) {

	dataResult := make([]mapstr.MapStr, 0)
	if err := m.dbProxy.Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).All(ctx, &dataResult); nil != err {
		blog.Errorf("request(%s): it is failed to find all models by the condition (%#v), error info is %s", ctx.ReqID, cond.ToMapStr(), err.Error())
		return dataResult, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return dataResult, nil
}

func (m *modelManager) delete(ctx core.ContextParams, cond universalsql.Condition) (uint64, error) {

	cnt, err := m.count(ctx, cond)
	if nil != err {
		return 0, err
	}

	if 0 == cnt {
		return 0, nil
	}

	if err = m.dbProxy.Table(common.BKTableNameObjDes).Delete(ctx, cond.ToMapStr()); nil != err {
		blog.Errorf("request(%s): it is failed to execute a deletion operation on the table (%s), error info is %s", ctx.ReqID, common.BKTableNameObjDes, err.Error())
		return 0, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return cnt, nil
}

// cascadeDelete 删除模型的字段，分组，唯一校验。模型等。
func (m *modelManager) cascadeDelete(ctx core.ContextParams, cond universalsql.Condition) (uint64, error) {

	modelItems, err := m.search(ctx, cond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute a cascade model deletion operation by the condition (%#v), error info is %s", ctx.ReqID, cond.ToMapStr(), err.Error())
		return 0, err
	}

	// 按照bk_obj_id删除的时候。业务下私有模型bk_obj_id相同。将会出现bug
	targetObjIDS := make([]string, 0)
	for _, modelItem := range modelItems {
		targetObjIDS = append(targetObjIDS, modelItem.ObjectID)
	}
	if len(targetObjIDS) == 0 {
		return 0, nil
	}

	if err := m.canCascadeDelete(ctx, targetObjIDS); err != nil {
		return 0, err
	}

	delCond := mongo.NewCondition()
	delCond.Element(mongo.Field(common.BKObjIDField).In(targetObjIDS))
	delCondMap := util.SetQueryOwner(delCond.ToMapStr(), ctx.SupplierAccount)

	// delete model property group
	if err := m.dbProxy.Table(common.BKTableNamePropertyGroup).Delete(ctx, delCondMap); err != nil {
		blog.ErrorJSON("delete model attribute group error. err:%s, cond:%s, rid:%s", err.Error(), delCondMap, ctx.ReqID)
		return 0, ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}

	// delete model property attribute
	if err := m.dbProxy.Table(common.BKTableNameObjAttDes).Delete(ctx, delCondMap); err != nil {
		blog.ErrorJSON("delete model attribute error. err:%s, cond:%s, rid:%s", err.Error(), delCondMap, ctx.ReqID)
		return 0, ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}

	// delete model unique
	if err := m.dbProxy.Table(common.BKTableNameObjUnique).Delete(ctx, delCondMap); err != nil {
		blog.ErrorJSON("delete model unique error. err:%s, cond:%s, rid:%s", err.Error(), delCondMap, ctx.ReqID)
		return 0, ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}

	// delete model
	if err := m.dbProxy.Table(common.BKTableNameObjDes).Delete(ctx, delCondMap); err != nil {
		blog.ErrorJSON("delete model unique error. err:%s, cond:%s, rid:%s", err.Error(), delCondMap, ctx.ReqID)
		return 0, ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}

	return uint64(len(targetObjIDS)), nil
}

// canCascadeDelete 判断是否可以删除
// 1. 检查是否内置模型
// 2. 是否包含实例
// 3. 是否有关联关系
func (m *modelManager) canCascadeDelete(ctx core.ContextParams, targetObjIDS []string) (err error) {
	// notice inner model not can delete
	for _, objID := range targetObjIDS {
		if util.IsInnerObject(objID) {
			return ctx.Error.Errorf(common.CCErrCoreServiceNotAllowDeleteErr, m.modelAttribute.getLangObjID(ctx, objID))
		}
	}

	// has instance
	instanceFilter := map[string]interface{}{
		common.BKObjIDField: map[string]interface{}{
			common.BKDBIN: targetObjIDS,
		},
		common.BkSupplierAccount: ctx.SupplierAccount,
	}
	cnt, err := m.dbProxy.Table(common.BKTableNameBaseInst).Find(instanceFilter).Count(ctx)
	if err != nil {
		blog.ErrorJSON("canCascadeDelete failed, count model instance failed, error. cond:%s, err:%s, rid:%s", instanceFilter, err.Error(), ctx.ReqID)
		return ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}
	if cnt > 0 {
		return ctx.Error.Error(common.CCErrCoreServiceModelHasInstanceErr)
	}

	// has model association, 不检查关联关系的是否有实例化。
	asstCond := mongo.NewCondition()
	asstCond.Or(
		mongo.Field(common.BKObjIDField).In(targetObjIDS),
		mongo.Field(common.BKAsstObjIDField).In(targetObjIDS),
	)
	asstCondMap := util.SetQueryOwner(asstCond.ToMapStr(), ctx.SupplierAccount)
	cnt, err = m.dbProxy.Table(common.BKTableNameObjAsst).Find(asstCondMap).Count(ctx)
	if err != nil {
		blog.ErrorJSON("canCascadeDelete failed, count model association failed, cond:%s, err:%s, rid:%s", asstCondMap, err.Error(), ctx.ReqID)
		return ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}
	if cnt > 0 {
		return ctx.Error.Error(common.CCErrCoreServiceModelHasAssociationErr)
	}

	return nil
}
