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
	"fmt"
	"gopkg.in/redis.v5"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"gopkg.in/mgo.v2/bson"
)

func (s *coreService) TransferHostToInnerModule(ctx *rest.Contexts) {
	inputData := &metadata.TransferHostToInnerModule{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	exceptionArr, err := s.core.HostOperation().TransferToInnerModule(ctx.Kit, inputData)
	if err != nil {
		blog.ErrorJSON("TransferHostToDefaultModule  error. err:%s, exception:%s, rid:%s", err.Error(), exceptionArr, ctx.Kit.Rid)
		ctx.RespEntityWithError(exceptionArr, err)
		return
	}
	ctx.RespEntity(nil)
}

func (s *coreService) TransferHostToNormalModule(ctx *rest.Contexts) {
	inputData := &metadata.HostsModuleRelation{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	exceptionArr, err := s.core.HostOperation().TransferToNormalModule(ctx.Kit, inputData)
	if err != nil {
		blog.ErrorJSON("TransferHostModule  error. err:%s, exception:%s, rid:%s", err.Error(), exceptionArr, ctx.Kit.Rid)
		ctx.RespEntityWithError(exceptionArr, err)
		return
	}
	ctx.RespEntity(nil)
}

func (s *coreService) TransferHostToAnotherBusiness(ctx *rest.Contexts) {
	inputData := &metadata.TransferHostsCrossBusinessRequest{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	exceptionArr, err := s.core.HostOperation().TransferToAnotherBusiness(ctx.Kit, inputData)
	if err != nil {
		blog.ErrorJSON("TransferHostCrossBusiness  error. err:%s, input:%s, exception:%s, rid:%s", err.Error(), inputData, exceptionArr, ctx.Kit.Rid)
		ctx.RespEntityWithError(exceptionArr, err)
		return
	}
	ctx.RespEntity(nil)
}

func (s *coreService) RemoveFromModule(ctx *rest.Contexts) {
	inputData := &metadata.RemoveHostsFromModuleOption{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	exceptionArr, err := s.core.HostOperation().RemoveFromModule(ctx.Kit, inputData)
	if err != nil {
		blog.ErrorJSON("RemoveFromModule error. err:%s, input:%s, exception:%s, rid:%s", err.Error(), inputData, exceptionArr, ctx.Kit.Rid)
		ctx.RespEntityWithError(exceptionArr, err)
		return
	}
	ctx.RespEntity(nil)
}

func (s *coreService) GetHostModuleRelation(ctx *rest.Contexts) {
	inputData := &metadata.HostModuleRelationRequest{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	relationArr, err := s.core.HostOperation().GetHostModuleRelation(ctx.Kit, inputData)
	if err != nil {
		blog.ErrorJSON("GetHostModuleRelation  error. err:%s, rid:%s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(relationArr)
}

func (s *coreService) DeleteHostFromSystem(ctx *rest.Contexts) {
	inputData := &metadata.DeleteHostRequest{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	exceptionArr, err := s.core.HostOperation().DeleteFromSystem(ctx.Kit, inputData)
	if err != nil {
		blog.ErrorJSON("DeleteHost  error. err:%s, exception:%s, rid:%s", err.Error(), exceptionArr, ctx.Kit.Rid)
		ctx.RespEntityWithError(exceptionArr, err)
		return
	}
	ctx.RespEntity(nil)
}

func (s *coreService) HostIdentifier(ctx *rest.Contexts) {
	inputData := &metadata.SearchHostIdentifierParam{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	hostIdentifierArr, err := s.core.HostOperation().Identifier(ctx.Kit, inputData)

	if err != nil {
		blog.InfoJSON("Identifier host identifier handle error. err:%s, input:%s, rid:%s", err.Error(), inputData, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(metadata.SearchHostIdentifierData{Info: hostIdentifierArr, Count: len(hostIdentifierArr)})
}

// TransferHostModuleDep is a TransferHostModule dependence
func (s *coreService) TransferHostModuleDep(kit *rest.Kit, input *metadata.HostsModuleRelation) ([]metadata.ExceptionResult, error) {
	exceptionArr, err := s.core.HostOperation().TransferToNormalModule(kit, input)
	if err != nil {
		blog.ErrorJSON("TransferHostModule  error. err:%s, exception:%s, rid:%s", err.Error(), exceptionArr, kit.Rid)
		return exceptionArr, err
	}
	return nil, nil
}

func (s *coreService) GetHostByID(ctx *rest.Contexts) {
	hostID, err := strconv.Atoi(ctx.Request.PathParameter("bk_host_id"))
	if err != nil {
		blog.Errorf("GetHostByID failed, get host by id, but got invalid host id, hostID: %s, err: %+v, rid: %s", hostID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKHostIDField))
		return
	}

	result := make(map[string]interface{}, 0)
	condition := common.KvMap{common.BKHostIDField: hostID}
	condition = util.SetModOwner(condition, ctx.Kit.SupplierAccount)
	err = s.db.Table(common.BKTableNameBaseHost).Find(condition).One(ctx.Kit.Ctx, &result)
	// TODO: return error for not found and deal error with all callers
	if err != nil && !s.db.IsNotFoundError(err) {
		blog.Errorf("GetHostByID failed, get host by id[%d] failed, err: %+v, rid: %s", hostID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
	}

	ctx.RespEntity(result)
}

func (s *coreService) GetHosts(ctx *rest.Contexts) {
	var dat metadata.ObjQueryInput
	if err := ctx.DecodeInto(&dat); err != nil {
		ctx.RespAutoError(err)
		return
	}

	condition := util.ConvParamsTime(dat.Condition)
	var cond map[string]interface{}
	switch value := condition.(type) {
	case map[string]interface{}:
		cond = value
	case mapstr.MapStr:
		cond = value
	case common.KvMap:
		cond = value
	default:
		out, err := bson.Marshal(condition)
		if err != nil {
			blog.Errorf("SetModOwner failed condition %#v, error %s", condition, err.Error())
		}
		err = bson.Unmarshal(out, &cond)
		if err != nil {
			blog.Errorf("SetModOwner failed condition %#v, error %s", condition, err.Error())
		}
	}
	condition = util.SetModOwner(cond, ctx.Kit.SupplierAccount)
	fieldArr := util.SplitStrField(dat.Fields, ",")
	result, err := s.getObjectByCondition(ctx, common.BKInnerObjIDHost, fieldArr, condition, dat.Sort, dat.Start, dat.Limit)
	if err != nil {
		blog.Errorf("get object failed type:%s,input:%v error:%v, rid: %s", common.BKInnerObjIDHost, dat, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostSelectInst))
		return
	}

	count, err := s.db.Table(common.BKTableNameBaseHost).Find(condition).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("get object failed type:%s ,input: %v error: %v, rid: %s", common.BKInnerObjIDHost, dat, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostSelectInst))
		return
	}

	ctx.RespEntity(metadata.HostInfo{
		Count: int(count),
		Info:  result,
	})
}

func (s *coreService) getObjectByCondition(ctx *rest.Contexts, objType string, fields []string, condition interface{}, sort string, skip, limit int) ([]mapstr.MapStr, error) {
	results := make([]mapstr.MapStr, 0)
	tName := common.GetInstTableName(objType)

	dbInst := s.db.Table(tName).Find(condition).Sort(sort).Start(uint64(skip)).Limit(uint64(limit))
	if 0 < len(fields) {
		dbInst.Fields(fields...)
	}
	if err := dbInst.All(ctx.Kit.Ctx, &results); err != nil {
		blog.Errorf("failed to query the inst , error info %s, rid: %s", err.Error(), ctx.Kit.Rid)
		return nil, err
	}

	// translate language for default name
	language := util.GetLanguage(ctx.Kit.Header)
	lang := s.language.CreateDefaultCCLanguageIf(language)
	if m, ok := defaultNameLanguagePkg[objType]; nil != lang && ok {
		for index, info := range results {
			l := m[fmt.Sprint(info["default"])]
			if len(l) >= 3 {
				results[index][l[1]] = util.FirstNotEmptyString(lang.Language(l[0]), fmt.Sprint(info[l[1]]), fmt.Sprint(info[l[2]]))
			}
		}
	}

	return results, nil
}

func (s *coreService) GetHostSnap(ctx *rest.Contexts) {
	hostID := ctx.Request.PathParameter(common.BKHostIDField)
	key := common.RedisSnapKeyPrefix + hostID
	result, err := s.cache.Get(key).Result()
	if nil != err && err != redis.Nil {
		blog.Errorf("get host snapshot failed, hostID: %v, err: %v, rid: %s", hostID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostGetSnapshot))
	}

	ctx.RespEntity(metadata.HostSnap{
		Data: result,
	})
}
