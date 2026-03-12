/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/backbone/service_mange/zk"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/tools/cmdb_ctl/app/config"
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/rs/xid"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
)

const defaultPageSize = 200

type objTopoLevelOperation = string

const (
	createOperation objTopoLevelOperation = "create"
	updateOperation objTopoLevelOperation = "update"
	deleteOperation objTopoLevelOperation = "delete"
)

type topoChangeConf struct {
	ParentObjID  string `json:"parent_obj_id"`
	ChildObjID   string `json:"child_obj_id"`
	CurrentObjID string `json:"current_obj_id"`
	NewObjID     string `json:"new_obj_id"` // when update,current_obj_id -> new_obj_id
	Operation    string `json:"operation"`  //add update delete
}

func (c *topoChangeConf) addFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.ParentObjID, "parent-obj-id", "biz", "parent-obj-id")
	cmd.Flags().StringVar(&c.ChildObjID, "child-obj-id", "set", "child-obj-id")
	cmd.Flags().StringVar(&c.CurrentObjID, "current-obj-id", "", "current-obj-id")
	cmd.Flags().StringVar(&c.NewObjID, "new-obj-id", "", "new-obj-id")
	cmd.Flags().StringVar(&c.Operation, "operation", "", "operation：create，update，delete")
}

// NewTopoLevelCheckCommand TODO
func NewTopoLevelCheckCommand() *cobra.Command {
	conf := new(topoChangeConf)
	cmd := &cobra.Command{
		Use:   "level",
		Short: "check business topo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTopoLevelCheck(conf)
		},
	}

	conf.addFlags(cmd)

	return cmd
}

func runTopoLevelCheck(c *topoChangeConf) error {

	switch c.Operation {
	case createOperation:
		if len(c.CurrentObjID) == 0 || len(c.ParentObjID) == 0 || len(c.ChildObjID) == 0 {
			return fmt.Errorf("invalid create parameter")
		}
	case deleteOperation:
		if len(c.CurrentObjID) == 0 {
			return fmt.Errorf("invalid delete parameter")
		}
	case updateOperation:
		if len(c.CurrentObjID) == 0 || len(c.NewObjID) == 0 {
			return fmt.Errorf("invalid update parameter")
		}
		if c.CurrentObjID == c.NewObjID {
			return fmt.Errorf("invalid update parameter")
		}
	default:
		return fmt.Errorf("invalid operation")
	}

	srv, err := newTopoChangeService(config.Conf.MongoURI, config.Conf.MongoRsName)
	if err != nil {
		return err
	}
	srv.conf = c
	return srv.changeTopo(c)
}

type topoChangeService struct {
	service *config.Service
	// 从业务信息中取出来业务所在的租户
	supplierAccount string
	conf            *topoChangeConf
}

func newTopoChangeService(mongoURI string, mongoRsName string) (*topoChangeService, error) {
	service, err := config.NewMongoService(mongoURI, mongoRsName)
	if err != nil {
		return nil, err
	}
	return &topoChangeService{
		service:         service,
		supplierAccount: "0",
	}, nil
}

func (s *topoChangeService) newKit() *rest.Kit {
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("X-Bkcmdb-User", "admin")
	header.Add("X-Bkcmdb-Supplier-Account", "0")
	header.Add("BK_USER", "admin")
	header.Add("HTTP_BLUEKING_SUPPLIER_ID", "0")
	header.Add(common.BKHTTPLanguage, "zh")
	header.Add(common.BKHTTPCCRequestID, xid.New().String())
	kit := rest.NewKitFromHeader(header, errors.NewFromCtx(nil))
	return kit
}

func (s *topoChangeService) getExistAssoc(currentObjID string) ([]metadata.Association, error) {
	// 不允许删除内置层级
	builtinObjIDs := []string{common.BKInnerObjIDApp, common.BKInnerObjIDSet,
		common.BKInnerObjIDModule, common.BKInnerObjIDHost}
	if slices.Contains(builtinObjIDs, currentObjID) {
		return nil, fmt.Errorf("bk_mainline: builtin object[%s] cannot be deleted", currentObjID)
	}

	var result []metadata.Association
	// topo正确
	filter := mapstr.MapStr{
		common.AssociationKindIDField: common.AssociationKindMainline,
	}
	err := s.service.DbProxy.Table(common.BKTableNameObjAsst).Find(filter).
		All(context.Background(), &result)
	if err != nil {
		return nil, fmt.Errorf("bk_mainline find err:%w", err)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("bk_mainline find no assoc")
	}
	//保证child->parent存在且唯一 common.AssociationKindIDField: bk_mainline,

	return result, nil
}

func (s *topoChangeService) changeTopo(c *topoChangeConf) error {
	clientSet, err := GetClientSetInterface()
	if err != nil {
		return err
	}

	switch c.Operation {
	case createOperation:
		err = s.runCreateLevelMainline(clientSet, c.CurrentObjID, c.ParentObjID, c.ChildObjID)

	case deleteOperation:
		err = s.runDeleteLevelMainline(clientSet, c.CurrentObjID)

	case updateOperation:
		err = s.runUpdateLevelMainline(clientSet, c.CurrentObjID, c.NewObjID)
	}
	if err != nil {
		log.Printf("[%v Topo] Failed !!,%v", c.Operation, err)
		return err
	}
	log.Printf("[%v Topo] Success !!,%+v", c.Operation, s.conf)
	return nil
}

// create
func (s *topoChangeService) runCreateLevelMainline(clientSet apimachinery.ClientSetInterface,
	currentObjID, parentObjID, childObjID string) error {

	result, err := s.getExistAssoc(currentObjID)
	if err != nil {
		return err
	}

	objectParentMap := make(map[string]string)
	for _, association := range result {
		if _, ok := objectParentMap[association.AsstObjID]; ok {
			return fmt.Errorf("check create bk_mainline already duplicate: %s", association.AsstObjID)
		}

		objectParentMap[association.AsstObjID] = association.ObjectID

		if association.AsstObjID == parentObjID && association.ObjectID != childObjID {
			return fmt.Errorf("bk_mainline:%s -> %s has no this level", parentObjID, childObjID)
		}
		if association.AsstObjID == currentObjID || association.ObjectID == currentObjID {
			return fmt.Errorf("bk_mainline:%s has already exist", currentObjID)
		}
	}

	kit := s.newKit()

	return s.createTopoLevel(kit, clientSet, currentObjID, parentObjID, childObjID)
}

func (s *topoChangeService) createTopoLevel(kit *rest.Kit, clientSet apimachinery.ClientSetInterface,
	currentObjID, parentObjID, childObjID string) error {

	return clientSet.CoreService().Txn().AutoRunTxn(kit.Ctx, kit.Header, func() error {
		err := s.createTopoChangeAsst(kit, clientSet, currentObjID, parentObjID, childObjID)
		if err != nil {
			return err
		}

		//查询所有parent
		parentCount, err := s.service.DbProxy.Table(common.GetInstTableName(parentObjID, "0")).Find(bson.M{}).Count(kit.Ctx)
		if err != nil {
			return fmt.Errorf("get parent model count failed, err: %v", err)
		}

		parentInstTable := common.GetInstTableName(parentObjID, "0")
		instIDField := metadata.GetInstIDFieldByObjID(parentObjID)

		for i := 0; i < 1+int(parentCount)/defaultPageSize; i++ {
			var parentList []mapstr.MapStr
			err := s.service.DbProxy.Table(parentInstTable).Find(bson.M{}).Fields(
				common.BKAppIDField, common.BKObjIDField, instIDField,
			).Sort(instIDField).Start(uint64(i*defaultPageSize)).Limit(defaultPageSize).All(kit.Ctx, &parentList)

			if err != nil {
				return fmt.Errorf("find parent inst err:%w", err)
			}
			// 插入模型实例，parentA -> childA  => parentA -> currentA(复制parentA) -> childA
			for _, parent := range parentList {
				parentID, err := metadata.GetInstID(parentObjID, parent)
				if err != nil {
					return fmt.Errorf("find parent id err:%w", err)
				}
				bizID, err := metadata.GetBizID(parent)
				if err != nil {
					return fmt.Errorf("find biz id err:%w", err)
				}

				currentObjInst := mapstr.MapStr{
					common.BKAppIDField:    bizID,
					common.BKObjIDField:    currentObjID,
					common.BKInstNameField: currentObjID + strconv.Itoa(int(parentID)),
					common.BKOwnerIDField:  "0",
					common.BKInstIDField:   parentID,
					common.BKParentIDField: parentID,
					common.LastTimeField:   metadata.Time{Time: time.Now()},
					common.CreateTimeField: metadata.Time{Time: time.Now()},
				}
				// 插入自定义模型
				err = s.service.DbProxy.Table(common.GetInstTableName(currentObjID, "0")).Insert(kit.Ctx, currentObjInst)
				if err != nil {
					return fmt.Errorf("insert current obj err:%w", err)
				}
			}
		}

		return nil
	})
}

func (s *topoChangeService) createTopoChangeAsst(kit *rest.Kit, clientSet apimachinery.ClientSetInterface,
	currentObjID, parentObjID, childObjID string) error {

	currentObjIDCount, err := s.service.DbProxy.Table(common.BKTableNameObjDes).Find(bson.M{
		common.BKObjIDField: currentObjID,
	}).Count(kit.Ctx)
	if err != nil {
		return fmt.Errorf("find current object[%s] failed: %v", currentObjID, err)
	}
	//已存在该模型
	if currentObjIDCount > 0 {
		//检查是否已经存在数据
		count, err := s.service.DbProxy.Table(common.GetInstTableName(currentObjID, "0")).Find(bson.M{}).Count(kit.Ctx)
		if err != nil {
			return fmt.Errorf("find current object[%s] failed: %v", currentObjID, err)
		}
		if count > 0 {
			return fmt.Errorf("%s has already with data", common.GetInstTableName(currentObjID, "0"))
		}
	} else {
		currentObj := metadata.Object{
			ID:         0,
			ObjCls:     "bk_uncategorized",
			ObjIcon:    "icon-cc-default",
			ObjectID:   currentObjID,
			ObjectName: currentObjID,
			OwnerID:    "0",
		}
		resp, err := clientSet.TopoServer().Object().CreateObject(kit.Ctx, kit.Header, currentObj)
		if err != nil {
			return fmt.Errorf("create model[%s] failed: %v", currentObjID, err)
		}
		if !resp.Result {
			return fmt.Errorf("create model[%s] failed: %v", currentObjID, resp)
		}
		err = s.service.DbProxy.DropTable(kit.Ctx, common.GetInstTableName(currentObjID, "0"))
		if err != nil {
			return fmt.Errorf("drop object %s table failed, err: %v", currentObjID, err)
		}
	}
	// child -> parent => child -> current
	_, err = clientSet.CoreService().Association().DeleteModelAssociation(kit.Ctx, kit.Header, &metadata.DeleteOption{
		Condition: mapstr.MapStr{
			common.AssociationKindIDField: common.AssociationKindMainline,
			common.BKAsstObjIDField:       parentObjID,
			common.BKObjIDField:           childObjID,
		},
	})
	if err != nil {
		return err
	}
	falseOption := false
	association := metadata.Association{
		OwnerID:         "0",
		AssociationName: fmt.Sprintf("%s_%s_%s", childObjID, common.AssociationKindMainline, currentObjID),
		ObjectID:        childObjID,
		AsstObjID:       currentObjID,
		AsstKindID:      common.AssociationKindMainline,
		Mapping:         "1:1",
		OnDelete:        "none",
		IsPre:           &falseOption,
	}
	_, err = clientSet.CoreService().Association().CreateMainlineModelAssociation(kit.Ctx,
		kit.Header, &metadata.CreateModelAssociation{
			Spec: association})
	if err != nil {
		return err
	}
	association.AssociationName = fmt.Sprintf("%s_%s_%s", currentObjID, common.AssociationKindMainline, parentObjID)
	association.ObjectID = currentObjID
	association.AsstObjID = parentObjID
	association.AsstKindID = common.AssociationKindMainline
	_, err = clientSet.CoreService().Association().CreateMainlineModelAssociation(kit.Ctx,
		kit.Header, &metadata.CreateModelAssociation{Spec: association})
	if err != nil {
		return err
	}
	return nil
}

// delete
func (s *topoChangeService) runDeleteLevelMainline(clientSet apimachinery.ClientSetInterface,
	currentObjID string) error {

	result, err := s.getExistAssoc(currentObjID)
	if err != nil {
		return err
	}

	var currentChildObjID string
	var currentGrandParentObjID string
	objectParentMap := make(map[string]string)
	for _, association := range result {

		if _, ok := objectParentMap[association.AsstObjID]; ok {
			return fmt.Errorf("check delete bk_mainline already duplicate: %s", association.AsstObjID)
		}
		objectParentMap[association.AsstObjID] = association.ObjectID

		// current 作为 parent（bk_asst_obj_id == currentObjID）=> child 指向 current
		if association.AsstObjID == currentObjID {
			currentChildObjID = association.ObjectID
		}
		// current 作为 child（bk_obj_id == currentObjID）=> current 指向 grandParent
		if association.ObjectID == currentObjID {
			currentGrandParentObjID = association.AsstObjID
		}
	}
	if len(currentGrandParentObjID) == 0 {
		return fmt.Errorf("bk_mainline: cannot find grandParent of %s", currentObjID)
	}
	if len(currentChildObjID) == 0 {
		return fmt.Errorf("bk_mainline: cannot find child of %s, maybe it is a leaf node", currentObjID)
	}

	kit := s.newKit()

	return s.deleteTopoLevel(kit, clientSet, currentObjID, currentGrandParentObjID, currentChildObjID)
}

// deleteTopoLevel 事务执行删除拓扑层级
// grandparent -> current -> child  =>  grandparent -> child
func (s *topoChangeService) deleteTopoLevel(kit *rest.Kit, clientSet apimachinery.ClientSetInterface,
	currentObjID, grandParentObjID, childObjID string) error {

	return clientSet.CoreService().Txn().AutoRunTxn(kit.Ctx, kit.Header, func() error {

		// ----------------------------------------
		// Step 1: 冲突检测
		// 将 current 下的子实例 bk_parent_id 改为 grandParent 的实例 id，
		// 若 grandParent 实例下已有同名 child 实例，则冲突
		// ----------------------------------------
		if err := s.checkDeleteConflict(kit, currentObjID, grandParentObjID, childObjID); err != nil {
			return err
		}

		// ----------------------------------------
		// Step 2: 修改关联关系
		// 删除: child -> current, current -> grandParent 两条 mainline
		// 新建: child -> grandParent
		// ----------------------------------------
		if err := s.deleteTopoChangeAsst(kit, clientSet, currentObjID, grandParentObjID, childObjID); err != nil {
			return err
		}

		// ----------------------------------------
		// Step 3: 迁移子实例的 bk_parent_id 指向 grandParent 实例
		// ----------------------------------------
		if err := s.migrateChildInstToGrandParent(kit, currentObjID, grandParentObjID, childObjID); err != nil {
			return err
		}

		// ----------------------------------------
		// Step 4: 清空并删除 current 层实例表
		// ----------------------------------------
		//err := s.service.DbProxy.DropTable(kit.Ctx, common.GetInstTableName(currentObjID, "0"))
		//if err != nil {
		//	return fmt.Errorf("drop object %s inst table failed, err: %v", currentObjID, err)
		//}

		return nil
	})
}

// checkDeleteConflict 检测删除层级后，子节点挂到祖父节点时是否有同名冲突
// 对应 JS 中 Step 3 全量冲突检测逻辑
func (s *topoChangeService) checkDeleteConflict(kit *rest.Kit,
	currentObjID, grandParentObjID, childObjID string) error {

	currentInstTable := common.GetInstTableName(currentObjID, "0")
	childInstTable := common.GetInstTableName(childObjID, "0")

	instIDField := metadata.GetInstIDFieldByObjID(currentObjID)
	childInstNameField := metadata.GetInstNameFieldName(childObjID)

	currentCount, err := s.service.DbProxy.Table(currentInstTable).Find(bson.M{}).Count(kit.Ctx)
	if err != nil {
		return fmt.Errorf("get current obj[%s] inst count failed, err: %v", currentObjID, err)
	}

	for i := 0; i < 1+int(currentCount)/defaultPageSize; i++ {
		var currentList []mapstr.MapStr
		err := s.service.DbProxy.Table(currentInstTable).Find(bson.M{}).Fields(
			common.BKAppIDField, instIDField, common.BKParentIDField,
		).Sort(instIDField).Start(uint64(i*defaultPageSize)).Limit(defaultPageSize).All(kit.Ctx, &currentList)
		if err != nil {
			return fmt.Errorf("find current obj[%s] inst failed, err: %w", currentObjID, err)
		}

		for _, currentInst := range currentList {
			currentInstID, err := metadata.GetInstID(currentObjID, currentInst)
			if err != nil {
				return fmt.Errorf("get inst id of obj[%s] failed, err: %w", currentObjID, err)
			}
			bizID, err := metadata.GetBizID(currentInst)
			if err != nil {
				return fmt.Errorf("get biz id of obj[%s] inst failed, err: %w", currentObjID, err)
			}
			grandParentInstID := currentInst[common.BKParentIDField]

			// 查出该 current 实例下所有 child 实例
			var childList []mapstr.MapStr
			err = s.service.DbProxy.Table(childInstTable).Find(bson.M{
				common.BKParentIDField: currentInstID,
				common.BKAppIDField:    bizID,
			}).Fields(metadata.GetInstIDFieldByObjID(childObjID), childInstNameField).
				All(kit.Ctx, &childList)
			if err != nil {
				return fmt.Errorf("find child obj[%s] inst failed, err: %w", childObjID, err)
			}

			// 逐一检测：child 同名实例是否已存在于 grandParent 下
			for _, childInst := range childList {
				childInstID, err := metadata.GetInstID(childObjID, childInst)
				if err != nil {
					return fmt.Errorf("get inst id of child obj[%s] failed, err: %w", childObjID, err)
				}
				childInstName := childInst[childInstNameField]

				var conflictInst []mapstr.MapStr
				err = s.service.DbProxy.Table(childInstTable).Find(bson.M{
					common.BKParentIDField: grandParentInstID,
					common.BKAppIDField:    bizID,

					childInstNameField: childInstName,
					metadata.GetInstIDFieldByObjID(childObjID): bson.M{
						"$ne": childInstID,
					},
				}).Fields(metadata.GetInstIDFieldByObjID(childObjID), childInstNameField).
					All(kit.Ctx, &conflictInst)
				if err != nil {
					return fmt.Errorf("conflict check for child obj[%s] inst[%v] failed, err: %w",
						childObjID, childInstID, err)
				}
				if len(conflictInst) > 0 {
					return fmt.Errorf(
						"conflict: child obj[%s] inst[%v] name[%v] already exists under grandParent inst[%v], "+
							"please rename or merge before deleting topo level[%s]",
						childObjID, childInstID, childInstName, grandParentInstID, currentObjID,
					)
				}
			}
		}
	}
	return nil
}

// deleteTopoChangeAsst 修改 mainline 关联关系
// 对应 JS 中 Step 4 的 cc_ObjAsst 操作
func (s *topoChangeService) deleteTopoChangeAsst(kit *rest.Kit, clientSet apimachinery.ClientSetInterface,
	currentObjID, grandParentObjID, childObjID string) error {

	// 删除 child -> current 的 mainline 关联
	_, err := clientSet.CoreService().Association().DeleteModelAssociation(kit.Ctx, kit.Header, &metadata.DeleteOption{
		Condition: mapstr.MapStr{
			common.AssociationKindIDField: common.AssociationKindMainline,
			common.BKAsstObjIDField:       currentObjID,
			common.BKObjIDField:           childObjID,
		},
	})
	if err != nil {
		return fmt.Errorf("delete mainline assoc [%s->%s] failed, err: %v",
			childObjID, currentObjID, err)
	}

	// 删除 current -> grandParent 的 mainline 关联
	_, err = clientSet.CoreService().Association().DeleteModelAssociation(kit.Ctx, kit.Header, &metadata.DeleteOption{
		Condition: mapstr.MapStr{
			common.AssociationKindIDField: common.AssociationKindMainline,
			common.BKAsstObjIDField:       grandParentObjID,
			common.BKObjIDField:           currentObjID,
		},
	})
	if err != nil {
		return fmt.Errorf("delete mainline assoc [%s->%s] failed, err: %v", currentObjID, grandParentObjID, err)
	}
	// 新建 child -> grandParent 的 mainline 关联（补回直连）
	falseOption := false
	_, err = clientSet.CoreService().Association().CreateMainlineModelAssociation(kit.Ctx, kit.Header,
		&metadata.CreateModelAssociation{
			Spec: metadata.Association{
				ID:              0,
				OwnerID:         "0",
				AssociationName: fmt.Sprintf("%s_%s_%s", childObjID, common.AssociationKindMainline, grandParentObjID),
				ObjectID:        childObjID,
				AsstObjID:       grandParentObjID,
				AsstKindID:      common.AssociationKindMainline,
				Mapping:         "1:1",
				OnDelete:        "none",
				IsPre:           &falseOption,
			},
		})
	if err != nil {
		return fmt.Errorf("create mainline assoc [%s->%s] failed, err: %v", childObjID, grandParentObjID, err)
	}

	return nil
}

// migrateChildInstToGrandParent 将 current 下所有 child 实例的 bk_parent_id 改指向 grandParent 实例
func (s *topoChangeService) migrateChildInstToGrandParent(kit *rest.Kit,
	currentObjID, grandParentObjID, childObjID string) error {

	currentInstTable := common.GetInstTableName(currentObjID, "0")
	childInstTable := common.GetInstTableName(childObjID, "0")

	instIDField := metadata.GetInstIDFieldByObjID(currentObjID)

	currentCount, err := s.service.DbProxy.Table(currentInstTable).Find(bson.M{}).Count(kit.Ctx)
	if err != nil {
		return fmt.Errorf("get current obj[%s] inst count failed, err: %v", currentObjID, err)
	}

	for i := 0; i < 1+int(currentCount)/defaultPageSize; i++ {
		var currentList []mapstr.MapStr
		err := s.service.DbProxy.Table(currentInstTable).Find(bson.M{}).Fields(
			common.BKAppIDField, instIDField, common.BKParentIDField,
		).Sort(instIDField).Start(uint64(i*defaultPageSize)).Limit(defaultPageSize).All(kit.Ctx, &currentList)

		if err != nil {
			return fmt.Errorf("find current obj[%s] inst failed, err: %w", currentObjID, err)
		}

		for _, currentInst := range currentList {
			currentInstID, err := metadata.GetInstID(currentObjID, currentInst)
			if err != nil {
				return fmt.Errorf("get inst id of obj[%s] failed, err: %w", currentObjID, err)
			}
			bizID, err := metadata.GetBizID(currentInst)
			if err != nil {
				return fmt.Errorf("get biz id of obj[%s] inst failed, err: %w", currentObjID, err)
			}
			grandParentInstID := currentInst[common.BKParentIDField]

			// 将该 current 实例下所有 child 实例的父节点改为 grandParent 实例
			err = s.service.DbProxy.Table(childInstTable).Update(kit.Ctx,
				bson.M{
					common.BKParentIDField: currentInstID,
					common.BKAppIDField:    bizID,
				},
				bson.M{
					common.BKParentIDField: grandParentInstID,
				},
			)
			if err != nil {
				return fmt.Errorf(
					"migrate child obj[%s] inst under current inst[%v] to grandParent inst[%v] failed, err: %w",
					childObjID, currentInstID, grandParentInstID, err,
				)
			}
		}
	}
	return nil
}

func (s *topoChangeService) updateTopoLevel(kit *rest.Kit, clientSet apimachinery.ClientSetInterface,
	currentObjID, grandParentObjID, childObjID, newObjID string) error {

	return clientSet.CoreService().Txn().AutoRunTxn(kit.Ctx, kit.Header, func() error {
		err := s.runDeleteLevelMainline(clientSet, currentObjID)
		if err != nil {
			return err
		}
		err = s.runCreateLevelMainline(clientSet, newObjID, grandParentObjID, childObjID)
		if err != nil {
			return err
		}
		return nil
	})
}

// update
func (s *topoChangeService) runUpdateLevelMainline(clientSet apimachinery.ClientSetInterface,
	currentObjID string, newObjID string) error {

	result, err := s.getExistAssoc(currentObjID)
	if err != nil {
		return err
	}

	var currentChildObjID string
	var currentGrandParentObjID string

	objectParentMap := make(map[string]string)
	for _, association := range result {
		if _, ok := objectParentMap[association.AsstObjID]; ok {
			return fmt.Errorf("check update bk_mainline already duplicate: %s", association.AsstObjID)
		}
		objectParentMap[association.AsstObjID] = association.ObjectID

		if association.AsstObjID == currentObjID {
			currentChildObjID = association.ObjectID
		}
		if association.ObjectID == currentObjID {
			currentGrandParentObjID = association.AsstObjID
		}
	}
	if len(currentGrandParentObjID) == 0 {
		return fmt.Errorf("bk_mainline: cannot find grandParent of %s", currentObjID)
	}
	if len(currentChildObjID) == 0 {
		return fmt.Errorf("bk_mainline: cannot find child of %s, maybe it is a leaf node", currentObjID)
	}

	kit := s.newKit()
	return s.updateTopoLevel(kit, clientSet, currentObjID, currentGrandParentObjID, currentChildObjID, newObjID)
}

// GetClientSetInterface 构建apimachinery.ClientSetInterface提供CMDB内置服务
func GetClientSetInterface() (apimachinery.ClientSetInterface, error) {
	client := zk.NewZkClient(config.Conf.ZkAddr, 1*time.Minute)
	if err := client.Start(); err != nil {
		return nil, fmt.Errorf("connect regdiscv [%s] failed: %v", config.Conf.ZkAddr, err)
	}
	if err := client.Ping(); err != nil {
		return nil, fmt.Errorf("connect regdiscv [%s] failed: %v", config.Conf.ZkAddr, err)
	}
	serviceDiscovery, err := discovery.NewServiceDiscovery(client)
	if err != nil {
		return nil, fmt.Errorf("connect regdiscv [%s] failed: %v", config.Conf.ZkAddr, err)
	}
	apiMachineryConfig := &util.APIMachineryConfig{
		QPS:       1000,
		Burst:     2000,
		TLSConfig: nil,
	}
	clientSet, err := apimachinery.NewApiMachinery(apiMachineryConfig, serviceDiscovery)
	if err != nil {
		return nil, fmt.Errorf("new api machinery failed, err: %v", err)
	}
	time.Sleep(time.Second)
	return clientSet, nil
}
