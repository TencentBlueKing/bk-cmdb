/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metrics"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/common/watch"
	"configcenter/src/scene_server/admin_server/upgrader"
	// import upgrader
	_ "configcenter/src/scene_server/admin_server/upgrader/imports"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/redis"
	daltypes "configcenter/src/storage/dal/types"
	streamtypes "configcenter/src/storage/stream/types"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	rootCmd.AddCommand(NewMigrateCommand())
}

// NewMigrateCommand new tool command for migration
func NewMigrateCommand() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "migrate cmdb data",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.PersistentFlags().StringVar(&configPath, "config-path", "", "migrate tool config file path")
	cmd.PersistentFlags().Var(auth.EnableAuthFlag, "enable-auth",
		"The auth center enable status, true for enabled, false for disabled")

	cmd.AddCommand(&cobra.Command{
		Use:   "db",
		Short: "migrate cmdb data in database",
		RunE: func(cmd *cobra.Command, args []string) error {
			srv, err := newMigrateService(configPath)
			if err != nil {
				return err
			}
			return srv.migrateDB()
		},
	})

	specifyVersionReq := new(MigrateSpecifyVersionRequest)
	specifyVersionCmd := &cobra.Command{
		Use:   "specify-version",
		Short: "migrate cmdb data of specify version",
		RunE: func(cmd *cobra.Command, args []string) error {
			srv, err := newMigrateService(configPath)
			if err != nil {
				return err
			}
			return srv.migrateSpecifyVersion(specifyVersionReq)
		},
	}
	specifyVersionReq.addFlags(specifyVersionCmd)

	cmd.AddCommand(specifyVersionCmd)

	return cmd
}

type migrateService struct {
	db      dal.DB
	watchDB dal.DB
	cache   redis.Client
	iam     *iam.IAM
}

func newMigrateService(configPath string) (*migrateService, error) {
	if !strings.HasSuffix(configPath, ".yaml") {
		return nil, fmt.Errorf("config path %s is invalid, should be a yaml file", configPath)
	}
	configPath = strings.TrimSuffix(configPath, ".yaml")

	if err := cc.SetMigrateFromFile(configPath); err != nil {
		return nil, fmt.Errorf("parse migration config from file[%s] failed, err: %v", configPath, err)
	}

	configDir, err := cc.String("confs.dir")
	if err != nil {
		return nil, fmt.Errorf("get migration config directory from file[%s] failed, err: %v", configPath, err)
	}

	// load mongodb, redis and common config from configure directory
	mongodbPath := configDir + "/" + types.CCConfigureMongo
	if err = cc.SetMongodbFromFile(mongodbPath); err != nil {
		return nil, fmt.Errorf("parse mongodb config from file[%s] failed, err: %v", mongodbPath, err)
	}

	redisPath := configDir + "/" + types.CCConfigureRedis
	if err = cc.SetRedisFromFile(redisPath); err != nil {
		return nil, fmt.Errorf("parse redis config from file[%s] failed, err: %v", redisPath, err)
	}

	commonPath := configDir + "/" + types.CCConfigureCommon
	if err = cc.SetCommonFromFile(commonPath); err != nil {
		return nil, fmt.Errorf("parse common config from file[%s] failed, err: %v", commonPath, err)
	}

	svc := new(migrateService)

	// new mongodb client
	dbConf, err := cc.Mongo("mongodb")
	if err != nil {
		return nil, fmt.Errorf("get mongodb config failed, err: %v", err)
	}

	svc.db, err = local.NewMgo(dbConf.GetMongoConf(), time.Minute)
	if err != nil {
		return nil, fmt.Errorf("new mongodb client failed, err: %v", err)
	}

	// new watch mongodb client
	watchDBConf, err := cc.Mongo("watch")
	if err != nil {
		return nil, fmt.Errorf("get watch mongodb config failed, err: %v", err)
	}

	svc.watchDB, err = local.NewMgo(watchDBConf.GetMongoConf(), time.Minute)
	if err != nil {
		return nil, fmt.Errorf("new watch mongodb client failed, err: %v", err)
	}

	// new redis client
	redisConf, err := cc.Redis("redis")
	if err != nil {
		return nil, fmt.Errorf("get redis config failed, err: %v", err)
	}

	svc.cache, err = redis.NewFromConfig(redisConf)
	if err != nil {
		return nil, fmt.Errorf("new redis client failed, err: %v", err)
	}

	// new iam client
	if auth.EnableAuthorize() {
		iamConf, err := iam.ParseConfigFromKV("iam", nil)
		if err != nil && auth.EnableAuthorize() {
			return nil, fmt.Errorf("parse iam config failed, err: %v", err)
		}

		metricService := metrics.NewService(metrics.Config{ProcessName: "migrate_tool"})
		svc.iam, err = iam.NewIAM(iamConf, metricService.Registry())
		if err != nil {
			return nil, fmt.Errorf("new iam client failed, err: %v", err)
		}
	}

	return svc, nil
}

func (s *migrateService) migrateDB() error {
	ctx := context.Background()
	if err := s.createWatchDBChainCollections(ctx); err != nil {
		return err
	}

	updateCfg := &upgrader.Config{
		OwnerID: common.BKDefaultOwnerID,
		User:    common.CCSystemOperatorUserName,
	}

	preVersion, finishedVersions, err := upgrader.Upgrade(ctx, s.db, s.cache, s.iam, updateCfg)
	if err != nil {
		return fmt.Errorf("db upgrade failed, err: %v", err)
	}

	currentVersion := preVersion
	if len(finishedVersions) > 0 {
		currentVersion = finishedVersions[len(finishedVersions)-1]
	}

	result := MigrationResult{
		Data:             "migrate success",
		PreVersion:       preVersion,
		CurrentVersion:   currentVersion,
		FinishedVersions: finishedVersions,
	}

	res, err := json.Marshal(result)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", string(res))
	return nil
}

// MigrationResult is the migration result
type MigrationResult struct {
	Data             interface{} `json:"data"`
	PreVersion       string      `json:"pre_version"`
	CurrentVersion   string      `json:"current_version"`
	FinishedVersions []string    `json:"finished_migrations"`
}

// dbChainTTLTime the ttl time seconds of the db event chain, used to set the ttl index of mongodb
const dbChainTTLTime = 5 * 24 * 60 * 60

func (s *migrateService) createWatchDBChainCollections(ctx context.Context) error {
	// create watch token table to store the last watch token info for every collection
	exists, err := s.watchDB.HasTable(ctx, common.BKTableNameWatchToken)
	if err != nil {
		return fmt.Errorf("check if table %s exists failed, err: %v", common.BKTableNameWatchToken, err)
	}

	if !exists {
		err = s.watchDB.CreateTable(ctx, common.BKTableNameWatchToken)
		if err != nil && !s.watchDB.IsDuplicatedError(err) {
			return fmt.Errorf("create table %s failed, err: %v", common.BKTableNameWatchToken, err)
		}
	}

	// create watch chain node table and init the last token info as empty for all collections
	cursorTypes := watch.ListCursorTypes()
	for _, cursorType := range cursorTypes {
		key, err := event.GetResourceKeyWithCursorType(cursorType)
		if err != nil {
			return fmt.Errorf("get resource key with cursor type %s failed, err: %v", cursorType, err)
		}

		exists, err := s.watchDB.HasTable(ctx, key.ChainCollection())
		if err != nil {
			return fmt.Errorf("check if table %s exists failed, err: %v", key.ChainCollection(), err)
		}

		if !exists {
			err = s.watchDB.CreateTable(ctx, key.ChainCollection())
			if err != nil && !s.watchDB.IsDuplicatedError(err) {
				return fmt.Errorf("create table %s failed, err: %v", key.ChainCollection(), err)
			}
		}

		if err = s.createWatchIndexes(ctx, cursorType, key); err != nil {
			return err
		}

		if err = s.createWatchToken(ctx, key); err != nil {
			return err
		}
	}
	return nil
}

func (s *migrateService) createWatchIndexes(ctx context.Context, cursorType watch.CursorType, key event.Key) error {
	indexes := []daltypes.Index{
		{Name: "index_id", Keys: bson.D{{common.BKFieldID, -1}}, Background: true, Unique: true},
		{Name: "index_cursor", Keys: bson.D{{common.BKCursorField, -1}}, Background: true, Unique: true},
		{Name: "index_cluster_time", Keys: bson.D{{common.BKClusterTimeField, -1}}, Background: true,
			ExpireAfterSeconds: dbChainTTLTime},
	}

	if cursorType == watch.ObjectBase || cursorType == watch.MainlineInstance || cursorType == watch.InstAsst {
		subResourceIndex := daltypes.Index{
			Name: "index_sub_resource", Keys: bson.D{{common.BKSubResourceField, 1}}, Background: true,
		}
		indexes = append(indexes, subResourceIndex)
	}

	existIndexArr, err := s.watchDB.Table(key.ChainCollection()).Indexes(ctx)
	if err != nil {
		return fmt.Errorf("get exist indexes for table %s failed, err: %v", key.ChainCollection(), err)
	}

	existIdxMap := make(map[string]bool)
	for _, index := range existIndexArr {
		existIdxMap[index.Name] = true
	}

	for _, index := range indexes {
		if _, exist := existIdxMap[index.Name]; exist {
			continue
		}

		err = s.watchDB.Table(key.ChainCollection()).CreateIndex(ctx, index)
		if err != nil && !s.watchDB.IsDuplicatedError(err) {
			return fmt.Errorf("create indexes for table %s failed, err: %v", key.ChainCollection(), err)
		}
	}
	return nil
}

func (s *migrateService) createWatchToken(ctx context.Context, key event.Key) error {
	filter := map[string]interface{}{
		"_id": key.Collection(),
	}

	count, err := s.watchDB.Table(common.BKTableNameWatchToken).Find(filter).Count(ctx)
	if err != nil {
		return fmt.Errorf("check if last watch token exists failed, err: %v, filter: %+v", err, filter)
	}

	if count > 0 {
		return nil
	}

	if key.Collection() == event.HostIdentityKey.Collection() {
		// host identity's watch token is different with other identity.
		// only set coll is ok, the other fields is useless
		data := mapstr.MapStr{
			"_id":                              key.Collection(),
			common.BKTableNameBaseHost:         watch.LastChainNodeData{Coll: common.BKTableNameBaseHost},
			common.BKTableNameModuleHostConfig: watch.LastChainNodeData{Coll: common.BKTableNameModuleHostConfig},
			common.BKTableNameBaseProcess:      watch.LastChainNodeData{Coll: common.BKTableNameBaseProcess},
		}
		if err = s.watchDB.Table(common.BKTableNameWatchToken).Insert(ctx, data); err != nil {
			return fmt.Errorf("init last watch token failed, err: %v, data: %+v", err, data)
		}
		return nil
	}

	if key.Collection() == event.BizSetRelationKey.Collection() {
		// biz set relation's watch token is generated in the same way with the host identity's watch token
		data := mapstr.MapStr{
			"_id":                        key.Collection(),
			common.BKTableNameBaseApp:    watch.LastChainNodeData{Coll: common.BKTableNameBaseApp},
			common.BKTableNameBaseBizSet: watch.LastChainNodeData{Coll: common.BKTableNameBaseBizSet},
			common.BKFieldID:             0,
			common.BKTokenField:          "",
		}
		if err = s.watchDB.Table(common.BKTableNameWatchToken).Insert(ctx, data); err != nil {
			return fmt.Errorf("init last biz set relation watch token failed, err: %v, data: %+v", err, data)
		}
		return nil
	}

	data := watch.LastChainNodeData{
		Coll:  key.Collection(),
		Token: "",
		StartAtTime: streamtypes.TimeStamp{
			Sec:  uint32(time.Now().Unix()),
			Nano: 0,
		},
	}
	if err = s.watchDB.Table(common.BKTableNameWatchToken).Insert(ctx, data); err != nil {
		return fmt.Errorf("init last watch token failed, err: %v, data: %+v", err, data)
	}
	return nil
}

func (s *migrateService) migrateSpecifyVersion(input *MigrateSpecifyVersionRequest) error {
	updateCfg := &upgrader.Config{
		OwnerID: common.BKDefaultOwnerID,
		User:    common.CCSystemOperatorUserName,
	}

	if input.CommitID != version.CCGitHash {
		return fmt.Errorf("commit id %s is not the same with current version %s", input.CommitID, version.CCGitHash)
	}

	err := upgrader.UpgradeSpecifyVersion(context.Background(), s.db, s.cache, s.iam, updateCfg, input.Version)
	if err != nil {
		return fmt.Errorf("db upgrade specify failed, err: %v", err)
	}

	fmt.Printf("migrate success, version: %s\n", input.Version)
	return nil
}

// MigrateSpecifyVersionRequest migrate specify version request
type MigrateSpecifyVersionRequest struct {
	CommitID string `json:"commit_id"`
	Version  string `json:"version"`
}

func (m *MigrateSpecifyVersionRequest) addFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&m.CommitID, "commit-id", "", "the commit id of this tool")
	cmd.Flags().StringVar(&m.Version, "version", "", "version to migrate")
}
