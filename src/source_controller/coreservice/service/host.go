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
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *coreService) TransferHostToInnerModule(ctx *rest.Contexts) {
	inputData := &metadata.TransferHostToInnerModule{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	err := s.core.HostOperation().TransferToInnerModule(ctx.Kit, inputData)
	if err != nil {
		blog.Errorf("transfer host to default module failed. err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
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
	err := s.core.HostOperation().TransferToNormalModule(ctx.Kit, inputData)
	if err != nil {
		blog.Errorf("transfer host to normal module failed. err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
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

	err := s.core.HostOperation().TransferToAnotherBusiness(ctx.Kit, inputData)
	if err != nil {
		blog.ErrorJSON("transfer host across business failed. err: %s, input: %s, rid: %s", err.Error(), inputData, ctx.Kit.Rid)
		ctx.RespAutoError(err)
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
	err := s.core.HostOperation().RemoveFromModule(ctx.Kit, inputData)
	if err != nil {
		blog.ErrorJSON("remove host from module failed. err: %s, input: %s, rid: %s", err.Error(), inputData, ctx.Kit.Rid)
		ctx.RespAutoError(err)
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

	err := s.core.HostOperation().DeleteFromSystem(ctx.Kit, inputData)
	if err != nil {
		blog.ErrorJSON("delete host error. err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
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
		blog.InfoJSON("Identifier host identifier handle error. err: %s, input: %s, rid:%s", err.Error(), inputData,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(metadata.SearchHostIdentifierData{Info: hostIdentifierArr, Count: len(hostIdentifierArr)})
}

// TransferHostModuleDep is a TransferHostModule dependence
func (s *coreService) TransferHostModuleDep(kit *rest.Kit, input *metadata.HostsModuleRelation) error {
	err := s.core.HostOperation().TransferToNormalModule(kit, input)
	if err != nil {
		blog.Errorf("transfer host to normal module failed. err: %s, rid: %s", err.Error(), kit.Rid)
		return err
	}
	return nil
}

func (s *coreService) GetHostByID(ctx *rest.Contexts) {
	hostID, err := strconv.Atoi(ctx.Request.PathParameter("bk_host_id"))
	if err != nil {
		blog.Errorf("GetHostByID failed, get host by id, but got invalid host id, hostID: %s, err: %+v, rid: %s", hostID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKHostIDField))
		return
	}

	result := make(metadata.HostMapStr, 0)
	condition := common.KvMap{common.BKHostIDField: hostID}
	condition = util.SetModOwner(condition, ctx.Kit.SupplierAccount)
	err = mongodb.Client().Table(common.BKTableNameBaseHost).Find(condition).One(ctx.Kit.Ctx, &result)
	// TODO: return error for not found and deal error with all callers
	if err != nil && !mongodb.Client().IsNotFoundError(err) {
		blog.Errorf("GetHostByID failed, get host by id[%d] failed, err: %+v, rid: %s", hostID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
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

	result := make([]metadata.HostMapStr, 0)
	dbInst := mongodb.Client().Table(common.BKTableNameBaseHost).Find(condition).Sort(dat.Sort).Start(uint64(dat.Start)).Limit(uint64(dat.Limit))
	if 0 < len(fieldArr) {
		dbInst.Fields(fieldArr...)
	}
	if err := dbInst.All(ctx.Kit.Ctx, &result); err != nil {
		blog.ErrorJSON("failed to query the host , cond: %s err: %s, rid: %s", cond, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	var finalCount uint64

	if !dat.DisableCounter {
		count, err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(condition).Count(ctx.Kit.Ctx)
		if err != nil {
			blog.Errorf("get object failed type:%s ,input: %v error: %v, rid: %s", common.BKInnerObjIDHost, dat, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostSelectInst))
			return
		}
		finalCount = count
	}

	info := make([]mapstr.MapStr, len(result))
	for index, host := range result {
		info[index] = mapstr.MapStr(host)
	}
	ctx.RespEntity(metadata.HostInfo{
		Count: int(finalCount),
		Info:  info,
	})
}

func (s *coreService) GetHostSnap(ctx *rest.Contexts) {
	hostID := ctx.Request.PathParameter(common.BKHostIDField)
	key := common.RedisSnapKeyPrefix + hostID
	result, err := redis.Client().Get(ctx.Kit.Ctx, key).Result()
	if nil != err && !redis.IsNilErr(err) {
		blog.Errorf("get host snapshot failed, hostID: %v, err: %v, rid: %s", hostID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostGetSnapshot))
		return
	}

	ctx.RespEntity(metadata.HostSnap{
		Data: result,
	})
}

func (s *coreService) GetHostSnapBatch(ctx *rest.Contexts) {
	input := metadata.HostSnapBatchInput{}
	if err := ctx.DecodeInto(&input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if len(input.HostIDs) == 0 {
		ctx.RespEntity(map[int64]string{})
		return
	}

	keys := []string{}
	for _, id := range input.HostIDs {
		keys = append(keys, common.RedisSnapKeyPrefix+strconv.FormatInt(id, 10))
	}

	res, err := redis.Client().MGet(ctx.Kit.Ctx, keys...).Result()
	if err != nil {
		if redis.IsNilErr(err) {
			ctx.RespEntity(map[int64]string{})
			return
		}
		blog.Errorf("get host snapshot failed, keys: %#v, err: %v, rid: %s", keys, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostGetSnapshot))
		return
	}

	ret := make(map[int64]string)
	for i, hostID := range input.HostIDs {
		if res[i] == nil {
			ret[hostID] = ""
			continue
		}
		value, ok := res[i].(string)
		if !ok {
			blog.Errorf("GetHostSnapBatch failed, hostID: %d, value in redis is not type string, but tyep: %T, value:%#v, rid: %s", hostID, res[i], res[i], ctx.Kit.Rid)
			ret[hostID] = ""
			continue
		}
		ret[hostID] = value
	}

	ctx.RespEntity(ret)
}

// GetDistinctHostIDsByTopoRelation get all  host ids by topology relation condition
func (s *coreService) GetDistinctHostIDsByTopoRelation(ctx *rest.Contexts) {
	inputData := &metadata.DistinctHostIDByTopoRelationRequest{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	hostIDArr, err := s.core.HostOperation().GetDistinctHostIDsByTopoRelation(ctx.Kit, inputData)
	if err != nil {
		blog.ErrorJSON("GetDistinctHostIDsByTopoRelation  error. err:%s, rid:%s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(metadata.DistinctID{IDArr: hostIDArr})
}

func (s *coreService) TransferHostResourceDirectory(ctx *rest.Contexts) {
	input := &metadata.TransferHostResourceDirectory{}
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	err := s.core.HostOperation().TransferResourceDirectory(ctx.Kit, input)
	if err != nil {
		blog.ErrorJSON("TransferHostResourceDirectory  error. err:%s, input:%s, rid:%s", err.Error(), input, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}
