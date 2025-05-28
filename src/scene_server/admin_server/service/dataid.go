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

package service

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"configcenter/pkg/tenant"
	"configcenter/pkg/tenant/logics"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/driver/mongodb"

	"github.com/emicklei/go-restful/v3"
)

const (
	snapStreamToName   = "cc_hostsnap_streamto"
	snapRouteName      = "cc_hostsnap_route"
	gseStreamToIDDBKey = "gse_stream_to_id"
	gseDataIDDBKey     = "gse_data_id"
)

// migrateDataID register, update or delete cc related data id in gse
func (s *Service) migrateDataID(req *restful.Request, resp *restful.Response) {
	kit := rest.NewKitFromHeader(req.Request.Header, s.CCErr)

	// validate tenant id
	if !logics.ValidatePlatformTenantMode(kit.TenantID, s.Config.EnableMultiTenantMode) {
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{
			ErrCode: common.CCErrAPICheckTenantInvalid,
			Msg:     fmt.Errorf("tenant %s is not system tenant, cannot migrate data id", kit.TenantID),
		})
		return
	}

	if !s.Config.MigrateDataID {
		blog.Errorf("migrate data id config is not set, rid: %s", kit.Rid)
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{
			Msg: kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "hostsnap.migrateDataID"),
		})
		return
	}

	// migrate data id for all tenants
	err := tenant.ExecForAllTenants(func(tenantID string) error {
		kit = kit.NewKit().WithTenant(tenantID)
		return s.upsertTenantDataID(kit)
	})
	if err != nil {
		_ = resp.WriteError(http.StatusOK, err)
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(nil))
}

func (s *Service) upsertTenantDataID(kit *rest.Kit) error {
	// upsert stream to config to gse
	if s.Config.SnapStreamToID == 0 {
		streamToID, err := s.upsertGseConfigStreamTo(kit)
		if err != nil {
			return err
		}
		s.Config.SnapStreamToID = streamToID
	}

	// generate gse config channel for current tenant
	channel, err := s.generateGseConfigChannel(kit, s.Config.SnapStreamToID)
	if err != nil {
		return err
	}

	// tenant channel id not exist, add to gse
	if channel.Metadata.ChannelID == 0 {
		return s.gseConfigAddRoute(kit, channel)
	}

	// get gse channel by data id, upsert config channel to gse
	channels, exists, err := s.gseConfigQueryRoute(kit, channel.Metadata.ChannelID)
	if err != nil {
		// gse returns a common error for not exist case before ee1.7.18/ce3.6.18, so we ignore it for compatibility
		blog.Errorf("query gse channel failed, **skip this error for not exist case**, err: %v, rid: %s", err, kit.Rid)
	}

	if !exists {
		return s.gseConfigAddRoute(kit, channel)
	}

	if len(channels) != 1 {
		blog.ErrorJSON("get multiple channel by data id %s, channels: %s, rid: %s", channel.Metadata.ChannelID,
			channels, kit.Rid)
		return &metadata.RespError{Msg: kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "hostsnap.dataid")}
	}

	if !reflect.DeepEqual(channels[0], *channel) {
		return s.gseConfigUpdateRoute(kit, channel)
	}

	return nil
}

// generateGseConfigStreamTo generate host snap stream to config by snap redis config
func (s *Service) generateGseConfigStreamTo(kit *rest.Kit) (*metadata.GseConfigStreamTo, error) {
	snapStreamTo := &metadata.GseConfigStreamTo{
		Name: snapStreamToName,
	}

	switch metadata.GseConfigReportMode(s.Config.SnapReportMode) {
	case "", metadata.GseConfigReportModeRedis:
		return s.generateRedisStreamTo(kit, snapStreamTo)
	case metadata.GseConfigReportModeKafka:
		return s.generateKafkaStreamTo(kit, snapStreamTo)
	}

	return nil, fmt.Errorf("can not support this SnapReportMode type: %s", s.Config.SnapReportMode)
}

func (s *Service) generateRedisStreamTo(kit *rest.Kit, snapStreamTo *metadata.GseConfigStreamTo) (
	*metadata.GseConfigStreamTo, error) {

	snapStreamTo.ReportMode = metadata.GseConfigReportModeRedis
	redisStreamToAddresses := make([]metadata.GseConfigStorageAddress, 0)
	snapRedisAddresses := strings.Split(s.Config.SnapRedis.Address, ",")
	for _, addr := range snapRedisAddresses {
		ipPort := strings.Split(addr, ":")
		if len(ipPort) != 2 || ipPort[0] == "" {
			blog.Errorf("host snap redis address is invalid, addr: %s, rid: %s", addr, kit.Rid)
			return nil, &metadata.RespError{Msg: kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "snap redis")}
		}

		port, err := strconv.ParseInt(ipPort[1], 10, 64)
		if err != nil {
			blog.Errorf("parse snap redis address port failed, port: %s, err: %v, rid: %s", ipPort[1], err, kit.Rid)
			return nil, &metadata.RespError{Msg: kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "snap redis")}
		}

		redisStreamToAddresses = append(redisStreamToAddresses, metadata.GseConfigStorageAddress{
			IP:   ipPort[0],
			Port: port,
		})
	}
	snapStreamTo.Redis = &metadata.GseConfigStreamToRedis{
		StorageAddresses: redisStreamToAddresses,
		Password:         s.Config.SnapRedis.Password,
		MasterName:       s.Config.SnapRedis.MasterName,
		SentinelPasswd:   s.Config.SnapRedis.SentinelPassword,
	}

	// The special logic here is to be compatible with the changes of the gse, it is necessary to explicitly specify
	// whether the mode is sentinel or single.
	if s.Config.SnapRedis.MasterName == "" {
		snapStreamTo.Redis.Mode = common.RedisSingleMode
	} else {
		snapStreamTo.Redis.Mode = common.RedisSentinelMode
	}
	return snapStreamTo, nil
}

func (s *Service) generateKafkaStreamTo(kit *rest.Kit, snapStreamTo *metadata.GseConfigStreamTo) (
	*metadata.GseConfigStreamTo, error) {

	if err := s.Config.SnapKafka.Check(); err != nil {
		blog.Errorf("kafka config is error, err: %v, rid: %s,", err, kit.Rid)
		return nil, err
	}

	snapStreamTo.ReportMode = metadata.GseConfigReportModeKafka
	kafkaStreamToAddresses := make([]metadata.GseConfigStorageAddress, 0)
	for _, addr := range s.Config.SnapKafka.Brokers {
		ipPort := strings.Split(addr, ":")
		if len(ipPort) != 2 || ipPort[0] == "" {
			blog.Errorf("host snap kafka address is invalid, addr: %s, rid: %s", addr, kit.Rid)
			return nil, &metadata.RespError{Msg: kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "snap kafka")}
		}

		port, err := strconv.ParseInt(ipPort[1], 10, 64)
		if err != nil {
			blog.Errorf("parse snap kafka address port failed, port: %s, err: %v, rid: %s", ipPort[1], err, kit.Rid)
			return nil, &metadata.RespError{Msg: kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "snap kafka")}
		}

		kafkaStreamToAddresses = append(kafkaStreamToAddresses, metadata.GseConfigStorageAddress{
			IP:   ipPort[0],
			Port: port,
		})
	}
	snapStreamTo.Kafka = &metadata.GseConfigStreamToKafka{
		StorageAddresses: kafkaStreamToAddresses,
	}
	if s.Config.SnapKafka.User != "" && s.Config.SnapKafka.Password != "" {
		snapStreamTo.Kafka.SaslUsername = s.Config.SnapKafka.User
		snapStreamTo.Kafka.SaslPassword = s.Config.SnapKafka.Password
		snapStreamTo.Kafka.SaslMechanisms = "SCRAM-SHA-512"
		snapStreamTo.Kafka.SecurityProtocol = "SASL_PLAINTEXT"
	}
	return snapStreamTo, nil
}

type dbStreamToID struct {
	HostSnap int64 `bson:"host_snap"`
}

// gseConfigQueryStreamTo get host snap stream to id from db, then get stream to config from gse
func (s *Service) gseConfigQueryStreamTo(kit *rest.Kit) (int64, []metadata.GseConfigAddStreamToParams, error) {
	cond := map[string]interface{}{"_id": gseStreamToIDDBKey}

	streamToID := new(dbStreamToID)
	if err := s.db.Shard(kit.SysShardOpts()).Table(common.BKTableNameSystem).Find(cond).One(s.ctx,
		&streamToID); err != nil {
		if mongodb.IsNotFoundError(err) {
			return 0, make([]metadata.GseConfigAddStreamToParams, 0), nil
		}
		blog.Errorf("get stream to id from db failed, err: %v, rid: %s", err, kit.Rid)
		return 0, nil, &metadata.RespError{Msg: kit.CCError.Error(common.CCErrCommDBSelectFailed)}
	}

	commonOperation := metadata.GseConfigOperation{
		OperatorName: kit.User,
	}
	params := &metadata.GseConfigQueryStreamToParams{
		Condition: metadata.GseConfigQueryStreamToCondition{
			GseConfigStreamToCondition: metadata.GseConfigStreamToCondition{
				StreamToID: streamToID.HostSnap,
				PlatName:   metadata.GseConfigPlatBkmonitor,
			},
		},
		Operation: commonOperation,
	}

	streamTos, err := s.GseClient.ConfigQueryStreamTo(s.ctx, kit.Header, params)
	if err != nil {
		blog.ErrorJSON("query stream to from gse failed, err: %s, params: %s, rid: %s", err, params, kit.Rid)
		return streamToID.HostSnap, nil, &metadata.RespError{Msg: kit.CCError.CCErrorf(common.CCErrCommMigrateFailed,
			err.Error())}
	}

	return streamToID.HostSnap, streamTos, nil
}

// gseConfigAddStreamTo add host snap redis stream to config to gse, then add or update stream to id in db
func (s *Service) gseConfigAddStreamTo(kit *rest.Kit, streamTo *metadata.GseConfigStreamTo) (
	int64, error) {

	params := &metadata.GseConfigAddStreamToParams{
		Metadata: metadata.GseConfigAddStreamToMetadata{
			PlatName: metadata.GseConfigPlatBkmonitor,
		},
		Operation: metadata.GseConfigOperation{
			OperatorName: kit.User,
		},
		StreamTo: *streamTo,
	}

	addStreamResult, err := s.GseClient.ConfigAddStreamTo(s.ctx, kit.Header, params)
	if err != nil {
		blog.ErrorJSON("add stream to gse failed, err: %s, params: %s, rid: %s", err, params, kit.Rid)
		return 0, &metadata.RespError{Msg: kit.CCError.CCErrorf(common.CCErrCommMigrateFailed, err.Error())}
	}

	cond := map[string]interface{}{"_id": gseStreamToIDDBKey}

	streamToID := dbStreamToID{
		HostSnap: addStreamResult.StreamToID,
	}

	if err := s.db.Shard(kit.SysShardOpts()).Table(common.BKTableNameSystem).Upsert(s.ctx, cond,
		&streamToID); err != nil {
		blog.Errorf("upsert stream to id %d to db failed, err: %v, rid: %s", addStreamResult.StreamToID, err, kit.Rid)
		return addStreamResult.StreamToID, &metadata.RespError{Msg: kit.CCError.Error(common.CCErrCommDBSelectFailed)}
	}

	return addStreamResult.StreamToID, nil
}

// gseConfigUpdateStreamTo update host snap redis stream to config to gse
func (s *Service) gseConfigUpdateStreamTo(kit *rest.Kit, streamTo *metadata.GseConfigStreamTo, streamToID int64) error {
	params := &metadata.GseConfigUpdateStreamToParams{
		Condition: metadata.GseConfigStreamToCondition{
			StreamToID: streamToID,
			PlatName:   metadata.GseConfigPlatBkmonitor,
		},
		Operation: metadata.GseConfigOperation{
			OperatorName: kit.User,
		},
		Specification: metadata.GseConfigUpdateStreamToSpecification{
			StreamTo: *streamTo,
		},
	}

	err := s.GseClient.ConfigUpdateStreamTo(s.ctx, kit.Header, params)
	if err != nil {
		blog.ErrorJSON("update stream to gse failed, err: %s, params: %s, rid: %s", err, params, kit.Rid)
		return &metadata.RespError{Msg: kit.CCError.CCErrorf(common.CCErrCommMigrateFailed, err.Error())}
	}
	return nil
}

// generateGseConfigChannel generate host snap gse config channel by snap redis/kafka config
func (s *Service) generateGseConfigChannel(kit *rest.Kit, streamToID int64) (*metadata.GseConfigChannel, error) {
	// get tenant data id from db
	cond := map[string]interface{}{"_id": gseDataIDDBKey}
	dataIDInfo := new(metadata.DataIDInfo)
	err := s.db.Shard(kit.SysShardOpts()).Table(common.BKTableNameSystem).Find(cond).
		Fields("host_snap."+kit.TenantID).One(s.ctx, &dataIDInfo)
	if err != nil && !mongodb.IsNotFoundError(err) {
		blog.Errorf("get stream to id from db failed, err: %v, rid: %s", err, kit.Rid)
		return nil, &metadata.RespError{Msg: kit.CCError.Error(common.CCErrCommDBSelectFailed)}
	}

	dataID := int64(0)
	if len(dataIDInfo.HostSnap) > 0 {
		dataID = dataIDInfo.HostSnap[kit.TenantID]
	}

	// generate gse config channel
	snapChannel := &metadata.GseConfigChannel{
		Metadata: metadata.GseConfigAddRouteMetadata{
			PlatName:  metadata.GseConfigPlatBkmonitor,
			ChannelID: dataID,
		},
		Route: []metadata.GseConfigRoute{{
			Name: snapRouteName + ":" + kit.TenantID,
			StreamTo: metadata.GseConfigRouteStreamTo{
				StreamToID: streamToID,
			},
		}},
	}

	switch metadata.GseConfigReportMode(s.Config.SnapReportMode) {
	case "", metadata.GseConfigReportModeRedis:
		snapChannel.Route[0].StreamTo.Redis = &metadata.GseConfigRouteRedis{
			ChannelName: common.SnapshotChannelName,
		}
		return snapChannel, nil
	case metadata.GseConfigReportModeKafka:
		snapChannel.Route[0].StreamTo.Kafka = &metadata.GseConfigRouteKafka{
			TopicName: common.SnapshotChannelName,
			Partition: s.Config.SnapKafka.Partition,
		}
		return snapChannel, nil
	}

	return nil, fmt.Errorf("can not support this SnapReportMode type: %s", s.Config.SnapReportMode)
}

// gseConfigQueryRoute get channel by host snap data id from gse
func (s *Service) gseConfigQueryRoute(kit *rest.Kit, dataID int64) ([]metadata.GseConfigChannel, bool, error) {
	commonOperation := metadata.GseConfigOperation{
		OperatorName: kit.User,
	}
	params := &metadata.GseConfigQueryRouteParams{
		Condition: metadata.GseConfigRouteCondition{
			ChannelID: dataID,
			PlatName:  metadata.GseConfigPlatBkmonitor,
		},
		Operation: commonOperation,
	}

	channels, exists, err := s.GseClient.ConfigQueryRoute(s.ctx, kit.Header, params)
	if err != nil {
		blog.ErrorJSON("query route from gse failed, err: %s, params: %s, rid: %s", err, params, kit.Rid)
		return nil, false, &metadata.RespError{Msg: kit.CCError.CCErrorf(common.CCErrCommMigrateFailed, err.Error())}
	}

	return channels, exists, nil
}

// gseConfigAddRoute add host snap channel to gse
func (s *Service) gseConfigAddRoute(kit *rest.Kit, channel *metadata.GseConfigChannel) error {
	params := &metadata.GseConfigAddRouteParams{
		Operation: metadata.GseConfigOperation{
			OperatorName: kit.User,
		},
		GseConfigChannel: *channel,
	}

	res, err := s.GseClient.ConfigAddRoute(s.ctx, kit.Header, params)
	if err != nil {
		blog.ErrorJSON("add channel to gse failed, err: %s, params: %s, rid: %s", err, params, kit.Rid)
		return &metadata.RespError{Msg: kit.CCError.CCErrorf(common.CCErrCommMigrateFailed, err.Error())}
	}

	// save tenant data id
	cond := map[string]interface{}{"_id": gseDataIDDBKey}
	dataIDInfo := map[string]any{
		"host_snap." + kit.TenantID: res.ChannelID,
	}
	err = s.db.Shard(kit.SysShardOpts()).Table(common.BKTableNameSystem).Upsert(s.ctx, cond, dataIDInfo)
	if err != nil {
		blog.Errorf("upsert data id info: %+v to db failed, err: %v, rid: %s", dataIDInfo, err, kit.Rid)
		return &metadata.RespError{Msg: kit.CCError.Error(common.CCErrCommDBSelectFailed)}
	}

	return nil
}

// gseConfigUpdateRoute update host snap redis channel to gse
func (s *Service) gseConfigUpdateRoute(kit *rest.Kit, channel *metadata.GseConfigChannel) error {
	params := &metadata.GseConfigUpdateRouteParams{
		Condition: metadata.GseConfigRouteCondition{
			ChannelID: channel.Metadata.ChannelID,
			PlatName:  metadata.GseConfigPlatBkmonitor,
		},
		Operation: metadata.GseConfigOperation{
			OperatorName: kit.User,
		},
		Specification: metadata.GseConfigUpdateRouteSpecification{
			Route:         channel.Route,
			StreamFilters: channel.StreamFilters,
		},
	}

	err := s.GseClient.ConfigUpdateRoute(s.ctx, kit.Header, params)
	if err != nil {
		blog.ErrorJSON("update channel to gse failed, err: %s, params: %s, rid: %s", err, params, kit.Rid)
		return &metadata.RespError{Msg: kit.CCError.CCErrorf(common.CCErrCommMigrateFailed, err.Error())}
	}
	return nil
}

// upsertGseConfigStreamTo get stream to config by db id from gse, if not found, register it, else update it
func (s *Service) upsertGseConfigStreamTo(kit *rest.Kit) (int64, error) {
	streamTo, err := s.generateGseConfigStreamTo(kit)
	if err != nil {
		return 0, err
	}

	streamToID, streamTos, err := s.gseConfigQueryStreamTo(kit)
	if err != nil {
		return 0, err
	}

	if len(streamTos) == 0 {
		return s.gseConfigAddStreamTo(kit, streamTo)
	}

	if len(streamTos) != 1 {
		blog.ErrorJSON("get multiple stream to(%s), rid: %s", streamTos, kit.Rid)
		return 0, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "stream to id")
	}

	if !reflect.DeepEqual(streamTos[0].StreamTo, *streamTo) {
		if err = s.gseConfigUpdateStreamTo(kit, streamTo, streamToID); err != nil {
			return 0, err
		}
	}
	return streamToID, nil
}
