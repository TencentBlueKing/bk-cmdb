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
	"context"
	"fmt"
	"os"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/tools/cmdb_ctl/app/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewTopoCheckCommand())
}

type topoCheckConf struct {
	bizID int64
}

func NewTopoCheckCommand() *cobra.Command {
	conf := new(topoCheckConf)

	cmd := &cobra.Command{
		Use:   "topo",
		Short: "check business topo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTopoCheck(conf)
		},
	}

	conf.addFlags(cmd)

	return cmd
}

func (c *topoCheckConf) addFlags(cmd *cobra.Command) {
	cmd.Flags().Int64Var(&c.bizID, "bizId", 2, "blueking business id. default value is 2")
}

type topoCheckService struct {
	service         *config.Service
	bizID           int64
	objectParentMap map[string]string
	modelIDs        []string
	instanceMap     map[string]*topoInstance
}

type topoInstance struct {
	ObjectID         string
	InstanceID       int64
	ParentInstanceID int64
	Default          int64
}

func newTopoCheckService(mongoURI string, mongoRsName string, bizID int64) (*topoCheckService, error) {
	service, err := config.NewMongoService(mongoURI, mongoRsName)
	if err != nil {
		return nil, err
	}
	return &topoCheckService{
		service:         service,
		bizID:           bizID,
		objectParentMap: make(map[string]string),
		instanceMap:     make(map[string]*topoInstance),
		modelIDs:        make([]string, 0),
	}, nil
}

func runTopoCheck(c *topoCheckConf) error {
	srv, err := newTopoCheckService(config.Conf.MongoURI, config.Conf.MongoRsName, c.bizID)
	if err != nil {
		return err
	}
	return srv.checkTopo()
}

func (s *topoCheckService) checkTopo() error {
	fmt.Println("=====================\nstart check")
	if err := s.searchMainlineModel(); err != nil {
		return err
	}
	if err := s.searchMainlineInstance(); err != nil {
		return err
	}
	s.checkMainlineInstanceTopo()
	fmt.Println("end check")
	return nil
}

func (s *topoCheckService) searchMainlineModel() error {
	fmt.Println("start searching mainline model")
	associations := make([]metadata.Association, 0)
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.AssociationKindIDField, Val: common.AssociationKindMainline})
	if err := s.service.DbProxy.Table(common.BKTableNameObjAsst).Find(cond.ToMapStr()).All(context.Background(), &associations); err != nil {
		return fmt.Errorf("query topo model mainline association from db failed, %+v", err)
	}
	for _, association := range associations {
		s.modelIDs = append(s.modelIDs, association.ObjectID)
		s.modelIDs = append(s.modelIDs, association.AsstObjID)
		if _, exist := s.objectParentMap[association.ObjectID]; !exist {
			s.objectParentMap[association.ObjectID] = association.AsstObjID
		}
	}
	s.modelIDs = util.StrArrayUnique(s.modelIDs)
	return nil
}

func (s *topoCheckService) searchMainlineInstance() error {
	fmt.Println("start searching mainline instance")
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.BKAppIDField, Val: s.bizID})
	// search business instances
	num, err := s.service.DbProxy.Table(common.BKTableNameBaseApp).Find(cond.ToMapStr()).Count(context.Background())
	if err != nil {
		return fmt.Errorf("get business instances by business id: %d failed, err: %+v", s.bizID, err)
	}
	if num != 1 {
		_, _ = fmt.Fprintf(os.Stderr, "business id: %d has too many(num = %d) business instances\n", s.bizID, num)
	} else {
		s.instanceMap[fmt.Sprintf("%s:%d", common.BKInnerObjIDApp, s.bizID)] = &topoInstance{
			ObjectID:         common.BKInnerObjIDApp,
			InstanceID:       s.bizID,
			ParentInstanceID: 0,
		}
	}
	// search set instances
	setInstances := make([]map[string]interface{}, 0)
	err = s.service.DbProxy.Table(common.BKTableNameBaseSet).Find(cond.ToMapStr()).All(context.Background(), &setInstances)
	if err != nil {
		return fmt.Errorf("get set instances by business id: %d failed, err: %+v", s.bizID, err)
	}
	for _, set := range setInstances {
		setID, err := util.GetInt64ByInterface(set[common.BKSetIDField])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "parse setID: %+v to int64 failed, err: %+v, set instance: %+v\n", set[common.BKSetIDField], err, set)
			continue
		}
		parentInstanceID, err := util.GetInt64ByInterface(set[common.BKInstParentStr])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "parse parentInstanceID:%+v to int64 failed, err: %+v, set instance: %+v\n", set[common.BKInstParentStr], err, set)
			continue
		}
		defaultFieldValue, err := util.GetInt64ByInterface(set[common.BKDefaultField])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "parse default field failed, default: %+v, err: %+v, set instance: %+v\n", set[common.BKDefaultField], err, set)
			continue
		}
		s.instanceMap[fmt.Sprintf("%s:%d", common.BKInnerObjIDSet, setID)] = &topoInstance{
			ObjectID:         common.BKInnerObjIDSet,
			InstanceID:       setID,
			ParentInstanceID: parentInstanceID,
			Default:          defaultFieldValue,
		}
	}
	// search module instances
	moduleInstances := make([]map[string]interface{}, 0)
	err = s.service.DbProxy.Table(common.BKTableNameBaseModule).Find(cond.ToMapStr()).All(context.Background(), &moduleInstances)
	if err != nil {
		return fmt.Errorf("get module instances by business id: %d failed, err: %+v", s.bizID, err)
	}
	for _, module := range moduleInstances {
		moduleID, err := util.GetInt64ByInterface(module[common.BKModuleIDField])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "parse moduleID: %+v to int64 failed, err: %+v, module instance: %+v\n", module[common.BKModuleIDField], err, module)
			continue
		}
		parentInstanceID, err := util.GetInt64ByInterface(module[common.BKInstParentStr])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "parse parentInstanceID:%+v to int64 failed, err: %+v, module instance: %+v\n", module[common.BKInstParentStr], err, module)
			continue
		}
		defaultFieldValue, err := util.GetInt64ByInterface(module[common.BKDefaultField])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "parse default field failed, default: %+v, err: %+v, module instance: %+v\n", module[common.BKDefaultField], err, module)
			continue
		}
		s.instanceMap[fmt.Sprintf("%s:%d", common.BKInnerObjIDModule, moduleID)] = &topoInstance{
			ObjectID:         common.BKInnerObjIDModule,
			InstanceID:       moduleID,
			ParentInstanceID: parentInstanceID,
			Default:          defaultFieldValue,
		}
	}
	// search mainline instances
	mainlineInstances := make([]map[string]interface{}, 0)
	cond = mongo.NewCondition()
	cond.Element(&mongo.In{Key: common.BKObjIDField, Val: s.modelIDs})
	_, metaCond := cond.Embed(metadata.BKMetadata)
	_, labelCond := metaCond.Embed(metadata.BKLabel)
	labelCond.Element(&mongo.Eq{Key: common.BKAppIDField, Val: strconv.FormatInt(s.bizID, 10)})
	err = s.service.DbProxy.Table(common.BKTableNameBaseInst).Find(cond.ToMapStr()).All(context.Background(), &mainlineInstances)
	if err != nil {
		return fmt.Errorf("get mainline instances by business id: %d failed, err: %+v", s.bizID, err)
	}
	for _, instance := range mainlineInstances {
		instanceID, err := util.GetInt64ByInterface(instance[common.BKInstIDField])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "parse instanceID: %+v to int64 failed, err: %+v, mainline instance: %+v\n", instance[common.BKInstIDField], err, instance)
			continue
		}
		parentInstanceID, err := util.GetInt64ByInterface(instance[common.BKInstParentStr])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "parse parentInstanceID:%+v to int64 failed, err: %+v, mainline instance: %+v\n", instance[common.BKInstParentStr], err, instance)
			continue
		}
		defaultFieldValue, err := util.GetInt64ByInterface(instance[common.BKDefaultField])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "parse default field failed, default: %+v, err: %+v, mainline instance: %+v\n", instance[common.BKDefaultField], err, instance)
			continue
		}
		objectID := util.GetStrByInterface(instance[common.BKObjIDField])
		s.instanceMap[fmt.Sprintf("%s:%d", objectID, instanceID)] = &topoInstance{
			ObjectID:         objectID,
			InstanceID:       instanceID,
			ParentInstanceID: parentInstanceID,
			Default:          defaultFieldValue,
		}
	}
	return nil
}

func (s *topoCheckService) checkMainlineInstanceTopo() {
	fmt.Println("start checking mainline instance topo")
	for _, instance := range s.instanceMap {
		if instance.ParentInstanceID == 0 {
			continue
		}
		var parentKey string
		if instance.ObjectID == common.BKInnerObjIDSet && instance.Default == 1 {
			// `空闲机池` 是一种特殊的set，它用来包含空闲机和故障机两个模块，它的父节点直接是业务（不论是否有自定义层级）
			parentKey = fmt.Sprintf("%s:%d", common.BKInnerObjIDApp, instance.ParentInstanceID)
		} else {
			parentObjectID := s.objectParentMap[instance.ObjectID]
			parentKey = fmt.Sprintf("%s:%d", parentObjectID, instance.ParentInstanceID)
		}
		// check whether parent instance exist, if not, try to get it at best.
		_, exist := s.instanceMap[parentKey]
		if exist {
			continue
		}
		mongoCondition := mongo.NewCondition()
		parentObjectID := s.objectParentMap[instance.ObjectID]
		parentIDField := common.GetInstIDField(parentObjectID)
		mongoCondition.Element(&mongo.Eq{Key: parentIDField, Val: instance.ParentInstanceID})
		missedInstances := make([]map[string]interface{}, 0)
		parentTable := common.GetInstTableName(parentObjectID)
		err := s.service.DbProxy.Table(parentTable).Find(mongoCondition.ToMapStr()).All(context.Background(), &missedInstances)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "find missing parent intance for object %s and instance: %d failed, err: %+v, parentObjectID: %s, ParentInstanceID: %d\n",
				instance.ObjectID, instance.InstanceID, err, parentObjectID, instance.ParentInstanceID)
			continue
		}
		if len(missedInstances) == 0 {
			_, _ = fmt.Fprintf(os.Stderr, "instance: %d of model: %s found no parent instance by parentObjectID %s and parentInstanceID: %d\n",
				instance.InstanceID, instance.ObjectID, parentObjectID, instance.ParentInstanceID)
			continue
		}
		if len(missedInstances) > 1 {
			_, _ = fmt.Fprintf(os.Stderr, "instance: %d of model: %s found too many(%d) parent instances: %+v by parentObjectID %s and parentInstanceID: %d\n",
				instance.InstanceID, instance.ObjectID, len(missedInstances), missedInstances, parentObjectID, instance.ParentInstanceID)
			continue
		}
		missedInstance := missedInstances[0]
		instanceID, err := util.GetInt64ByInterface(missedInstance[common.BKInstIDField])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "parse instanceID: %+v to int64 failed, err: %+v, %s instance: %+v\n", missedInstance[common.BKInstIDField], err, parentObjectID, missedInstance)
			continue
		}
		parentInstanceID, err := util.GetInt64ByInterface(missedInstance[common.BKInstParentStr])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "parse parentInstanceID:%+v to int64 failed, err: %+v, %s instance: %+v\n", missedInstance[common.BKInstParentStr], err, parentObjectID, missedInstance)
			continue
		}
		defaultFieldValue, err := util.GetInt64ByInterface(missedInstance[common.BKDefaultField])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "parse default field failed, default: %+v, err: %+v, %s instance: %+v\n", missedInstance[common.BKDefaultField], err, parentObjectID, missedInstance)
			continue
		}
		objectID := util.GetStrByInterface(missedInstance[common.BKObjIDField])
		s.instanceMap[fmt.Sprintf("%s:%d", objectID, instanceID)] = &topoInstance{
			ObjectID:         objectID,
			InstanceID:       instanceID,
			ParentInstanceID: parentInstanceID,
			Default:          defaultFieldValue,
		}
	}
}
