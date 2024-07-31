/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package iam

import (
	"context"
	"net/http"
	"time"

	iamcli "configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	headerutil "configcenter/src/common/http/header/util"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/logics"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/storage/driver/mongodb"
)

const (
	// 同步周期最小值
	syncIAMPeriodMinutesMin = 1
	// 同步周期默认值
	syncIAMPeriodMinutesDefault = 5
)

// syncor used to sync iam
type syncor struct {
	// 同步周期
	SyncIAMPeriodMinutes int
	// db mongodb实例连接，用于判断是否数据库初始化已完成，防止和模型实例权限迁移的upgrader冲突
	db dal.RDB
}

// NewSyncor TODO
func NewSyncor() *syncor {
	return &syncor{}
}

// SetSyncIAMPeriod set the sync period
func (s *syncor) SetSyncIAMPeriod(periodMinutes int) {
	s.SyncIAMPeriodMinutes = periodMinutes
	if s.SyncIAMPeriodMinutes < syncIAMPeriodMinutesMin {
		s.SyncIAMPeriodMinutes = syncIAMPeriodMinutesDefault
	}
	blog.Infof("sync iam period is %d minutes", s.SyncIAMPeriodMinutes)
}

// SetDB set db
func (s *syncor) SetDB(db dal.RDB) {
	s.db = db
}

// newHeader 创建IAM同步需要的header
func newHeader() http.Header {
	header := headerutil.BuildHeader(common.BKIAMSyncUser, common.BKSuperOwnerID)
	httpheader.SetLanguage(header, "cn")
	return header
}

// newKit 创建新的Kit
func newKit() *rest.Kit {
	header := newHeader()
	ctx := util.NewContextFromHTTPHeader(header)
	rid := httpheader.GetRid(header)
	user := httpheader.GetUser(header)
	supplierAccount := httpheader.GetSupplierAccount(header)
	defaultCCError := util.GetDefaultCCError(header)

	return &rest.Kit{
		Rid:             rid,
		Header:          header,
		Ctx:             ctx,
		CCError:         defaultCCError,
		User:            user,
		SupplierAccount: supplierAccount,
	}
}

// SyncIAM sync the system instances resource between CMDB and IAM
func (s *syncor) SyncIAM(iamCli *iamcli.IAM, redisCli redis.Client, lgc *logics.Logics) {
	if !auth.EnableAuthorize() {
		return
	}
	time.Sleep(time.Minute)

	// 等待数据库初始化完成，防止和模型实例权限迁移的upgrader冲突
	rid := util.GenerateRID()
	for dbReady := false; !dbReady; {
		var err error
		dbReady, err = upgrader.DBReady(context.Background(), s.db)
		if err != nil {
			blog.Errorf("sync iam, check whether db initialization is complete failed, err: %v, rid: %s", err, rid)
			time.Sleep(5 * time.Second)
			continue
		}
		if !dbReady {
			blog.Warnf("sync iam, but db initialization is not complete, rid: %s", rid)
			time.Sleep(5 * time.Second)
			continue
		}
	}

	for {
		// new kit with a different rid, header
		kit := newKit()

		// only master can run it
		if !lgc.ServiceManageInterface.IsMaster() {
			blog.V(4).Infof("it is not master, skip sync iam, rid: %s", kit.Rid)
			time.Sleep(time.Minute)
			continue
		}

		blog.Infof("start sync iam, rid: %s", kit.Rid)

		objects, err := GetCustomObjects(kit.Ctx, s.db)
		if err != nil {
			blog.Errorf("sync iam failed, get custom objects err: %s ,rid: %s", err, kit.Rid)
			time.Sleep(time.Duration(s.SyncIAMPeriodMinutes) * time.Minute)
			continue
		}

		if err := iamCli.SyncIAMSysInstances(kit.Ctx, redisCli, objects); err != nil {
			blog.Errorf("sync iam failed, sync iam system instances err: %s ,rid: %s", err, kit.Rid)
			time.Sleep(time.Duration(s.SyncIAMPeriodMinutes) * time.Minute)
			continue
		}

		blog.Infof("finish sync iam successfully, rid:%s", kit.Rid)
		time.Sleep(time.Duration(s.SyncIAMPeriodMinutes) * time.Minute)
	}
}

// GetCustomObjects get all custom objects(without inner and mainline objects that authorize separately)
func GetCustomObjects(ctx context.Context, db dal.DB) ([]metadata.Object, error) {
	// get mainline objects
	associations := make([]metadata.Association, 0)
	filter := mapstr.MapStr{
		common.AssociationKindIDField: common.AssociationKindMainline,
	}

	err := db.Table(common.BKTableNameObjAsst).Find(filter).Fields(common.BKObjIDField).All(ctx, &associations)
	if err != nil {
		blog.Errorf("get mainline associations failed, err: %v", err)
		return nil, err
	}

	// get all excluded objectIDs
	excludeObjIDs := []string{
		common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule,
		common.BKInnerObjIDHost, common.BKInnerObjIDProc, common.BKInnerObjIDPlat,
	}
	for _, association := range associations {
		if !metadata.IsCommon(association.ObjectID) {
			excludeObjIDs = append(excludeObjIDs, association.ObjectID)
		}
	}

	// get all custom objects
	objects := make([]metadata.Object, 0)
	condition := map[string]interface{}{
		common.BKIsPre: false,
		common.BKObjIDField: map[string]interface{}{
			common.BKDBNIN: excludeObjIDs,
		},
	}
	if err := db.Table(common.BKTableNameObjDes).Find(condition).All(ctx, &objects); err != nil {
		blog.Errorf("get all custom objects failed, err: %v", err)
		return nil, err
	}

	// 表格字段类型的object不注册到权限中心，这里需要将他们过滤出来
	cnt, err := mongodb.Client().Table(common.BKTableNameModelQuoteRelation).Find(nil).Count(ctx)
	if err != nil {
		blog.Errorf("count cc_ModelQuoteRelation failed, err: %v", err)
		return nil, err
	}
	if cnt == 0 {
		return objects, nil
	}

	relationObjMap := make(map[string]struct{})
	for i := uint64(0); i < cnt; i += common.BKMaxLimitSize {
		relations := make([]metadata.ModelQuoteRelation, 0)
		err = mongodb.Client().Table(common.BKTableNameModelQuoteRelation).Find(nil).Start(i).
			Limit(common.BKMaxLimitSize).Fields(common.BKDestModelField).All(ctx, &relations)
		if err != nil {
			blog.Errorf("list model quote relations failed, err: %v", err)
			return nil, err
		}

		for _, relation := range relations {
			relationObjMap[relation.DestModel] = struct{}{}
		}
	}

	result := make([]metadata.Object, 0)
	for _, object := range objects {
		if _, ok := relationObjMap[object.ObjectID]; ok {
			continue
		}

		result = append(result, object)
	}

	return result, nil
}
