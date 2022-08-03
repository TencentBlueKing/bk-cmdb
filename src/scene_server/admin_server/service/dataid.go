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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/resource/esb"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful/v3"
)

const (
	snapStreamToName   = "cc_hostsnap_streamto"
	snapRouteName      = "cc_hostsnap_route"
	gseStreamToIDDBKey = "gse_stream_to_id"
)

var (
	snapStreamTo *metadata.GseConfigStreamTo
	snapChannel  *metadata.GseConfigChannel
)

type snapshotVersion string

const (
	oldVersion snapshotVersion = "old_version"
	newVersion snapshotVersion = "new_version"
)

// migrateDataID register, update or delete cc related data id in gse
func (s *Service) migrateDataID(req *restful.Request, resp *restful.Response) {
	header := req.Request.Header
	rid := util.GetHTTPCCRequestID(header)
	user := util.GetUser(header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	if s.Config.SnapDataID == 0 {
		blog.Errorf("host snap data id not set in configuration, rid: %s", rid)
		result := &metadata.RespError{Msg: defErr.CCErrorf(common.CCErrCommParamsNeedSet, "hostsnap.dataid")}
		_ = resp.WriteError(http.StatusOK, result)
		return
	}

	// upsert stream to config to gse
	streamToID, err := s.UpsertGseConfigStreamTo(header, user, defErr, rid, newVersion)
	if err != nil {
		_ = resp.WriteError(http.StatusOK, err)
		return
	}

	// upsert config channel to gse
	channel, err := s.generateGseConfigChannel(streamToID, s.Config.SnapDataID, rid, newVersion)
	if err != nil {
		_ = resp.WriteError(http.StatusOK, err)
		return
	}

	channels, exists, err := s.gseConfigQueryRoute(header, user, defErr, rid)
	if err != nil {
		// gse returns a common error for not exist case before ee1.7.18/ce3.6.18, so we ignore it for compatibility
		blog.Errorf("query gse channel failed, ** skip this error for not exist case **, err: %v, rid: %s", err, rid)
	}

	if !exists {
		if err := s.gseConfigAddRoute(channel, header, user, defErr, rid); err != nil {
			_ = resp.WriteError(http.StatusOK, err)
			return
		}
	} else if len(channels) != 1 {
		blog.ErrorJSON("get multiple channel by data id %s, channels: %s, rid: %s", s.Config.SnapDataID, channels, rid)
		result := &metadata.RespError{Msg: defErr.CCErrorf(common.CCErrCommParamsInvalid, "hostsnap.dataid")}
		_ = resp.WriteError(http.StatusOK, result)
		return
	} else {
		if !reflect.DeepEqual(channels[0], *channel) {
			if err := s.gseConfigUpdateRoute(channel, header, user, defErr, rid); err != nil {
				_ = resp.WriteError(http.StatusOK, err)
				return
			}
		}
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(nil))
}

// migrateOldDataID migrate old version data id, register it when it is not exist
func (s *Service) migrateOldDataID(req *restful.Request, resp *restful.Response) {
	header := req.Request.Header
	rid := util.GetHTTPCCRequestID(header)
	user := util.GetUser(header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	if err := s.migrateOldVersionDataID(header, user, defErr, rid); err != nil {
		_ = resp.WriteError(http.StatusOK, err)
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(nil))
}

// migrateOldVersionDataID old version data id is registered using script with gse version < 3.1, but new version of gse
// only allows registering data id by http interface, and the script can not be used, so we need to compensate by
// registering it in cc if it is not already registered before by script in the former version
func (s *Service) migrateOldVersionDataID(header http.Header, user string, defErr errors.DefaultCCErrorIf,
	rid string) error {

	const oldDataID = 1001

	// get already registered channel by host snap data id from gse, if not found, register it with its stream to
	commonOperation := metadata.GseConfigOperation{
		OperatorName: user,
	}
	queryParams := &metadata.GseConfigQueryRouteParams{
		Condition: metadata.GseConfigRouteCondition{
			ChannelID: oldDataID,
		},
		Operation: commonOperation,
	}

	channels, isDataIDExists, err := esb.EsbClient().GseSrv().ConfigQueryRoute(s.ctx, header, queryParams)
	if err != nil {
		// gse returns a common error for not exist case before ee1.7.18/ce3.6.18, so we ignore it for compatibility
		blog.Errorf("query gse channel failed, ** skip this error for not exist case **, err: %v, rid: %s", err, rid)
	}

	// if old data id has channels, we need to check if they are registered by cc or by other system like bk-monitor
	var existsPlatName metadata.GseConfigPlatName
	if isDataIDExists {
		bizID, err := s.getSnapBizID(rid)
		if err != nil {
			return err
		}

		// check if channel name is snapshot+snap biz id to confirm if it is registered by cc, skip in this situation
		for _, channel := range channels {
			existsPlatName = channel.Metadata.PlatName
			for _, route := range channel.Route {
				if route.StreamTo.Redis == nil {
					continue
				}
				if route.StreamTo.Redis.ChannelName == fmt.Sprintf("snapshot%d", bizID) ||
					(route.StreamTo.Redis.BizID == bizID && route.StreamTo.Redis.DataSet == "snapshot") {
					blog.Infof("old gse data id is already exist, skip registering it, rid: %s", rid)
					return nil
				}
			}
		}
	}

	// old stream to and channel is the same with the new one except for the data id, generate in the same way
	streamToID, err := s.UpsertGseConfigStreamTo(header, user, defErr, rid, oldVersion)
	if err != nil {
		return err
	}

	oldChannel, err := s.generateGseConfigChannel(streamToID, oldDataID, rid, oldVersion)
	if err != nil {
		blog.Errorf("generate gse channel failed, err: %v, stream to id: %d, rid: %s", err, streamToID, rid)
		return err
	}

	// update the exist data id's corresponding channel, add the route to it
	if isDataIDExists {
		params := &metadata.GseConfigUpdateRouteParams{
			Condition: metadata.GseConfigRouteCondition{
				ChannelID: oldDataID,
				PlatName:  existsPlatName,
			},
			Operation: metadata.GseConfigOperation{
				OperatorName: user,
			},
			Specification: metadata.GseConfigUpdateRouteSpecification{
				Route:         oldChannel.Route,
				StreamFilters: oldChannel.StreamFilters,
			},
		}

		err := esb.EsbClient().GseSrv().ConfigUpdateRoute(s.ctx, header, params)
		if err != nil {
			blog.Errorf("update old data id route to gse failed, err: %v, params: %#v, rid: %s", err, params, rid)
			return &metadata.RespError{Msg: defErr.CCErrorf(common.CCErrCommMigrateFailed, err.Error())}
		}
	} else {
		if err := s.gseConfigAddRoute(oldChannel, header, user, defErr, rid); err != nil {
			blog.Errorf("add route to gse failed, err: %v, channel: %v, rid: %s", err, oldChannel, rid)
			return err
		}
	}

	return nil
}

// generateGseConfigStreamTo generate host snap stream to config by snap redis config
func (s *Service) generateGseConfigStreamTo(header http.Header, user string, defErr errors.DefaultCCErrorIf,
	rid string, version snapshotVersion) (*metadata.GseConfigStreamTo, error) {

	if snapStreamTo != nil {
		return snapStreamTo, nil
	}
	snapStreamTo = &metadata.GseConfigStreamTo{
		Name: snapStreamToName,
	}

	if version == oldVersion || s.Config.SnapReportMode == "" ||
		metadata.GseConfigReportMode(s.Config.SnapReportMode) == metadata.GseConfigReportModeRedis {

		snapStreamTo.ReportMode = metadata.GseConfigReportModeRedis
		redisStreamToAddresses := make([]metadata.GseConfigStorageAddress, 0)
		snapRedisAddresses := strings.Split(s.Config.SnapRedis.Address, ",")
		for _, addr := range snapRedisAddresses {
			ipPort := strings.Split(addr, ":")
			if len(ipPort) != 2 || ipPort[0] == "" {
				blog.Errorf("host snap redis address is invalid, addr: %s, rid: %s", addr, rid)
				return nil, &metadata.RespError{Msg: defErr.CCErrorf(common.CCErrCommParamsInvalid, "snap redis")}
			}

			port, err := strconv.ParseInt(ipPort[1], 10, 64)
			if err != nil {
				blog.Errorf("parse snap redis address port failed, port: %s, err: %v, rid: %s", ipPort[1], err, rid)
				return nil, &metadata.RespError{Msg: defErr.CCErrorf(common.CCErrCommParamsInvalid, "snap redis")}
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

	if metadata.GseConfigReportMode(s.Config.SnapReportMode) == metadata.GseConfigReportModeKafka {
		if err := s.Config.SnapKafka.Check(); err != nil {
			blog.Errorf("kafka config is error, err: %v, rid: %s,", err, rid)
			return nil, err
		}

		snapStreamTo.ReportMode = metadata.GseConfigReportModeKafka
		kafkaStreamToAddresses := make([]metadata.GseConfigStorageAddress, 0)
		for _, addr := range s.Config.SnapKafka.Brokers {
			ipPort := strings.Split(addr, ":")
			if len(ipPort) != 2 || ipPort[0] == "" {
				blog.Errorf("host snap kafka address is invalid, addr: %s, rid: %s", addr, rid)
				return nil, &metadata.RespError{Msg: defErr.CCErrorf(common.CCErrCommParamsInvalid, "snap kafka")}
			}

			port, err := strconv.ParseInt(ipPort[1], 10, 64)
			if err != nil {
				blog.Errorf("parse snap kafka address port failed, port: %s, err: %v, rid: %s", ipPort[1], err, rid)
				return nil, &metadata.RespError{Msg: defErr.CCErrorf(common.CCErrCommParamsInvalid, "snap kafka")}
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

	return nil, fmt.Errorf("can not support this SnapReportMode type: %s", s.Config.SnapReportMode)
}

type dbStreamToID struct {
	HostSnap int64 `bson:"host_snap"`
}

// gseConfigQueryStreamTo get host snap stream to id from db, then get stream to config from gse
func (s *Service) gseConfigQueryStreamTo(header http.Header, user string, defErr errors.DefaultCCErrorIf, rid string) (
	int64, []metadata.GseConfigAddStreamToParams, error) {

	cond := map[string]interface{}{"_id": gseStreamToIDDBKey}
	streamToID := new(dbStreamToID)
	if err := s.db.Table(common.BKTableNameSystem).Find(cond).One(s.ctx, &streamToID); err != nil {
		if s.db.IsNotFoundError(err) {
			return 0, make([]metadata.GseConfigAddStreamToParams, 0), nil
		}
		blog.Errorf("get stream to id from db failed, err: %v, rid: %s", err, rid)
		return 0, nil, &metadata.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)}
	}

	commonOperation := metadata.GseConfigOperation{
		OperatorName: user,
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

	streamTos, err := esb.EsbClient().GseSrv().ConfigQueryStreamTo(s.ctx, header, params)
	if err != nil {
		blog.ErrorJSON("query stream to from gse failed, err: %s, params: %s, rid: %s", err, params, rid)
		return streamToID.HostSnap, nil, &metadata.RespError{Msg: defErr.CCErrorf(common.CCErrCommMigrateFailed,
			err.Error())}
	}

	return streamToID.HostSnap, streamTos, nil
}

// gseConfigAddStreamTo add host snap redis stream to config to gse, then add or update stream to id in db
func (s *Service) gseConfigAddStreamTo(streamTo *metadata.GseConfigStreamTo, header http.Header, user string,
	defErr errors.DefaultCCErrorIf, rid string) (int64, error) {

	params := &metadata.GseConfigAddStreamToParams{
		Metadata: metadata.GseConfigAddStreamToMetadata{
			PlatName: metadata.GseConfigPlatBkmonitor,
		},
		Operation: metadata.GseConfigOperation{
			OperatorName: user,
		},
		StreamTo: *streamTo,
	}

	addStreamResult, err := esb.EsbClient().GseSrv().ConfigAddStreamTo(s.ctx, header, params)
	if err != nil {
		blog.ErrorJSON("add stream to gse failed, err: %s, params: %s, rid: %s", err, params, rid)
		return 0, &metadata.RespError{Msg: defErr.CCErrorf(common.CCErrCommMigrateFailed, err.Error())}
	}

	cond := map[string]interface{}{"_id": gseStreamToIDDBKey}
	streamToID := dbStreamToID{
		HostSnap: addStreamResult.StreamToID,
	}

	if err := s.db.Table(common.BKTableNameSystem).Upsert(s.ctx, cond, &streamToID); err != nil {
		blog.Errorf("upsert stream to id %d to db failed, err: %v, rid: %s", addStreamResult.StreamToID, err, rid)
		return addStreamResult.StreamToID, &metadata.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)}
	}

	return addStreamResult.StreamToID, nil
}

// gseConfigUpdateStreamTo update host snap redis stream to config to gse
func (s *Service) gseConfigUpdateStreamTo(streamTo *metadata.GseConfigStreamTo, streamToID int64, header http.Header,
	user string, defErr errors.DefaultCCErrorIf, rid string) error {

	params := &metadata.GseConfigUpdateStreamToParams{
		Condition: metadata.GseConfigStreamToCondition{
			StreamToID: streamToID,
			PlatName:   metadata.GseConfigPlatBkmonitor,
		},
		Operation: metadata.GseConfigOperation{
			OperatorName: user,
		},
		Specification: metadata.GseConfigUpdateStreamToSpecification{
			StreamTo: *streamTo,
		},
	}

	err := esb.EsbClient().GseSrv().ConfigUpdateStreamTo(s.ctx, header, params)
	if err != nil {
		blog.ErrorJSON("update stream to gse failed, err: %s, params: %s, rid: %s", err, params, rid)
		return &metadata.RespError{Msg: defErr.CCErrorf(common.CCErrCommMigrateFailed, err.Error())}
	}
	return nil
}

// getSnapBizID get the biz id that host snap uses
func (s *Service) getSnapBizID(rid string) (int64, error) {
	cfgCond := map[string]interface{}{
		"_id": common.ConfigAdminID,
	}
	cfg := make(map[string]string)
	err := s.db.Table(common.BKTableNameSystem).Find(cfgCond).Fields(common.ConfigAdminValueField).One(s.ctx, &cfg)
	if nil != err {
		blog.Errorf("get config admin failed, err: %v, rid: %s", err, rid)
		return 0, err
	}

	configAdmin := new(metadata.PlatformSettingConfig)
	if err := json.Unmarshal([]byte(cfg[common.ConfigAdminValueField]), configAdmin); err != nil {
		blog.Errorf("unmarshal config admin(%s) failed, err: %v, rid: %s", cfg[common.ConfigAdminValueField], err, rid)
		return 0, err
	}

	bizCond := map[string]interface{}{common.BKAppNameField: configAdmin.Backend.SnapshotBizName}
	biz := new(metadata.BizBasicInfo)
	if err := s.db.Table(common.BKTableNameBaseApp).Find(bizCond).One(s.ctx, biz); err != nil {
		blog.Errorf("get snap biz by name(%s) failed, err: %v, rid: %s", configAdmin.Backend.SnapshotBizName, err, rid)
		return 0, err
	}

	return biz.BizID, nil
}

// generateGseConfigChannel generate host snap stream to config by snap redis config
func (s *Service) generateGseConfigChannel(streamToID, dataID int64, rid string,
	version snapshotVersion) (*metadata.GseConfigChannel, error) {

	if snapChannel != nil {
		snapChannel.Metadata.ChannelID = dataID
		return snapChannel, nil
	}

	bizID, err := s.getSnapBizID(rid)
	if err != nil {
		return nil, err
	}

	snapChannel = &metadata.GseConfigChannel{
		Metadata: metadata.GseConfigAddRouteMetadata{
			PlatName:  metadata.GseConfigPlatBkmonitor,
			ChannelID: dataID,
		},
		Route: []metadata.GseConfigRoute{{
			Name: snapRouteName,
			StreamTo: metadata.GseConfigRouteStreamTo{
				StreamToID: streamToID,
			},
		}},
	}

	if version == oldVersion || s.Config.SnapReportMode == "" ||
		metadata.GseConfigReportMode(s.Config.SnapReportMode) == metadata.GseConfigReportModeRedis {

		snapChannel.Route[0].StreamTo.Redis = &metadata.GseConfigRouteRedis{
			ChannelName: fmt.Sprintf("snapshot%d", bizID),
			// compatible for the older version of gse that uses DataSet+BizID as channel name
			DataSet: "snapshot",
			BizID:   bizID,
		}
		return snapChannel, nil
	}

	if metadata.GseConfigReportMode(s.Config.SnapReportMode) == metadata.GseConfigReportModeKafka {
		snapChannel.Route[0].StreamTo.Kafka = &metadata.GseConfigRouteKafka{
			TopicName: fmt.Sprintf("snapshot%d", bizID),
			// compatible for the older version of gse that uses DataSet+BizID as channel name
			DataSet:   "snapshot",
			BizID:     bizID,
			Partition: s.Config.SnapKafka.Partition,
		}
		return snapChannel, nil
	}

	return nil, fmt.Errorf("can not support this SnapReportMode type: %s", s.Config.SnapReportMode)
}

// gseConfigQueryRoute get channel by host snap data id from gse
func (s *Service) gseConfigQueryRoute(header http.Header, user string, defErr errors.DefaultCCErrorIf, rid string) (
	[]metadata.GseConfigChannel, bool, error) {

	commonOperation := metadata.GseConfigOperation{
		OperatorName: user,
	}
	params := &metadata.GseConfigQueryRouteParams{
		Condition: metadata.GseConfigRouteCondition{
			ChannelID: s.Config.SnapDataID,
			PlatName:  metadata.GseConfigPlatBkmonitor,
		},
		Operation: commonOperation,
	}

	channels, exists, err := esb.EsbClient().GseSrv().ConfigQueryRoute(s.ctx, header, params)
	if err != nil {
		blog.ErrorJSON("query route from gse failed, err: %s, params: %s, rid: %s", err, params, rid)
		return nil, false, &metadata.RespError{Msg: defErr.CCErrorf(common.CCErrCommMigrateFailed, err.Error())}
	}

	return channels, exists, nil
}

// gseConfigAddRoute add host snap channel to gse
func (s *Service) gseConfigAddRoute(channel *metadata.GseConfigChannel, header http.Header, user string,
	defErr errors.DefaultCCErrorIf, rid string) error {

	params := &metadata.GseConfigAddRouteParams{
		Operation: metadata.GseConfigOperation{
			OperatorName: user,
		},
		GseConfigChannel: *channel,
	}

	_, err := esb.EsbClient().GseSrv().ConfigAddRoute(s.ctx, header, params)
	if err != nil {
		blog.ErrorJSON("add channel to gse failed, err: %s, params: %s, rid: %s", err, params, rid)
		return &metadata.RespError{Msg: defErr.CCErrorf(common.CCErrCommMigrateFailed, err.Error())}
	}

	return nil
}

// gseConfigUpdateRoute update host snap redis channel to gse
func (s *Service) gseConfigUpdateRoute(channel *metadata.GseConfigChannel, header http.Header, user string,
	defErr errors.DefaultCCErrorIf, rid string) error {

	params := &metadata.GseConfigUpdateRouteParams{
		Condition: metadata.GseConfigRouteCondition{
			ChannelID: s.Config.SnapDataID,
			PlatName:  metadata.GseConfigPlatBkmonitor,
		},
		Operation: metadata.GseConfigOperation{
			OperatorName: user,
		},
		Specification: metadata.GseConfigUpdateRouteSpecification{
			Route:         channel.Route,
			StreamFilters: channel.StreamFilters,
		},
	}

	err := esb.EsbClient().GseSrv().ConfigUpdateRoute(s.ctx, header, params)
	if err != nil {
		blog.ErrorJSON("update channel to gse failed, err: %s, params: %s, rid: %s", err, params, rid)
		return &metadata.RespError{Msg: defErr.CCErrorf(common.CCErrCommMigrateFailed, err.Error())}
	}
	return nil
}

// UpsertGseConfigStreamTo get stream to config by db id from gse, if not found, register it, else update it
func (s *Service) UpsertGseConfigStreamTo(header http.Header, user string, defErr errors.DefaultCCErrorIf,
	rid string, version snapshotVersion) (int64, error) {

	streamTo, err := s.generateGseConfigStreamTo(header, user, defErr, rid, version)
	if err != nil {
		return 0, err
	}

	streamToID, streamTos, err := s.gseConfigQueryStreamTo(header, user, defErr, rid)
	if err != nil {
		return 0, err
	}

	if len(streamTos) == 0 {
		if streamToID, err = s.gseConfigAddStreamTo(streamTo, header, user, defErr, rid); err != nil {
			return 0, err
		}
	} else if len(streamTos) != 1 {
		blog.ErrorJSON("get multiple stream to(%s), rid: %s", streamTos, rid)
		return 0, defErr.CCErrorf(common.CCErrCommParamsInvalid, "stream to id")
	} else {
		if !reflect.DeepEqual(streamTos[0].StreamTo, *streamTo) {
			if err := s.gseConfigUpdateStreamTo(streamTo, streamToID, header, user, defErr, rid); err != nil {
				return 0, err
			}
		}
	}
	return streamToID, nil
}
