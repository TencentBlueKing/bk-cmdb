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

	"github.com/emicklei/go-restful"
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

	// get stream to config by db id from gse, if not found, register it, else update it
	streamTo, err := s.generateGseConfigStreamTo(header, user, defErr, rid)
	if err != nil {
		_ = resp.WriteError(http.StatusOK, err)
		return
	}

	streamToID, streamTos, err := s.gseConfigQueryStreamTo(header, user, defErr, rid)
	if err != nil {
		_ = resp.WriteError(http.StatusOK, err)
		return
	}

	if len(streamTos) == 0 {
		if streamToID, err = s.gseConfigAddStreamTo(streamTo, header, user, defErr, rid); err != nil {
			_ = resp.WriteError(http.StatusOK, err)
			return
		}
	} else if len(streamTos) != 1 {
		blog.ErrorJSON("get multiple stream to(%s), rid: %s", streamTos, rid)
		result := &metadata.RespError{Msg: defErr.CCErrorf(common.CCErrCommParamsInvalid, "stream to id")}
		_ = resp.WriteError(http.StatusOK, result)
		return
	} else {
		if !reflect.DeepEqual(streamTos[0].StreamTo, *streamTo) {
			if err := s.gseConfigUpdateStreamTo(streamTo, streamToID, header, user, defErr, rid); err != nil {
				_ = resp.WriteError(http.StatusOK, err)
				return
			}
		}
	}

	// get already registered channel by host snap data id from gse, if not found, register it, else update it
	channel, err := s.generateGseConfigChannel(streamToID, header, user, defErr, rid)
	if err != nil {
		_ = resp.WriteError(http.StatusOK, err)
		return
	}

	channels, err := s.gseConfigQueryRoute(header, user, defErr, rid)
	if err != nil {
		blog.Errorf("query gse channel failed, ** skip this error for not exist case **, err: %v, rid: %s", err, rid)
		// TODO clarify this error when gse returns a specified error code
		// _ = resp.WriteError(http.StatusOK, err)
		// return
	}

	if len(channels) == 0 {
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

// generateGseConfigStreamTo generate host snap stream to config by snap redis config
func (s *Service) generateGseConfigStreamTo(header http.Header, user string, defErr errors.DefaultCCErrorIf, rid string) (
	*metadata.GseConfigStreamTo, error) {

	if snapStreamTo != nil {
		return snapStreamTo, nil
	}

	redisStreamToAddresses := make([]metadata.GseConfigStorageAddress, 0)
	snapRedisAddresses := strings.Split(s.Config.SnapRedis.Address, ",")
	for _, addr := range snapRedisAddresses {
		ipPort := strings.Split(addr, ":")
		if len(ipPort) != 2 {
			blog.Errorf("host snap redis address is invalid, addr: %s, rid: %s", addr, rid)
			return nil, &metadata.RespError{Msg: defErr.CCErrorf(common.CCErrCommParamsInvalid, "snap redis")}
		}

		port, err := strconv.ParseInt(ipPort[1], 10, 64)
		if err != nil {
			blog.Errorf("parse snap redis address port failed, err: %v, port: %s, rid: %s", err, ipPort[1], rid)
			return nil, &metadata.RespError{Msg: defErr.CCErrorf(common.CCErrCommParamsInvalid, "snap redis")}
		}

		redisStreamToAddresses = append(redisStreamToAddresses, metadata.GseConfigStorageAddress{
			IP:   ipPort[0],
			Port: port,
		})
	}

	snapStreamTo = &metadata.GseConfigStreamTo{
		Name:       snapStreamToName,
		ReportMode: metadata.GseConfigReportModeRedis,
		Redis: &metadata.GseConfigStreamToRedis{
			StorageAddresses: redisStreamToAddresses,
			Password:         s.Config.SnapRedis.Password,
			MasterName:       s.Config.SnapRedis.MasterName,
		},
	}
	return snapStreamTo, nil
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
		return streamToID.HostSnap, nil, &metadata.RespError{Msg: defErr.CCErrorf(common.CCErrCommMigrateFailed, err.Error())}
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

// generateGseConfigStreamTo generate host snap stream to config by snap redis config
func (s *Service) generateGseConfigChannel(streamToID int64, header http.Header, user string,
	defErr errors.DefaultCCErrorIf, rid string) (*metadata.GseConfigChannel, error) {

	if snapChannel != nil {
		return snapChannel, nil
	}

	cfgCond := map[string]interface{}{
		"_id": common.ConfigAdminID,
	}
	cfg := make(map[string]string)
	err := s.db.Table(common.BKTableNameSystem).Find(cfgCond).Fields(common.ConfigAdminValueField).One(s.ctx, &cfg)
	if nil != err {
		blog.Errorf("get config admin failed, err: %v", err)
		return nil, err
	}

	configAdmin := new(metadata.ConfigAdmin)
	if err := json.Unmarshal([]byte(cfg[common.ConfigAdminValueField]), configAdmin); err != nil {
		blog.Errorf("unmarshal config admin failed, err: %v, config: %s", err, cfg[common.ConfigAdminValueField])
		return nil, err
	}

	bizCond := map[string]interface{}{common.BKAppNameField: configAdmin.Backend.SnapshotBizName}
	biz := new(metadata.BizBasicInfo)
	if err := s.db.Table(common.BKTableNameBaseApp).Find(bizCond).One(s.ctx, biz); err != nil {
		blog.Errorf("get snap biz failed, err: %v, biz name: %s", err, configAdmin.Backend.SnapshotBizName)
		return nil, err
	}

	snapChannel = &metadata.GseConfigChannel{
		Metadata: metadata.GseConfigAddRouteMetadata{
			PlatName:  metadata.GseConfigPlatBkmonitor,
			ChannelID: s.Config.SnapDataID,
		},
		Route: []metadata.GseConfigRoute{{
			Name: snapRouteName,
			StreamTo: metadata.GseConfigRouteStreamTo{
				StreamToID: streamToID,
				Redis: &metadata.GseConfigRouteRedis{
					ChannelName: fmt.Sprintf("snapshot%d", biz.BizID),
					// compatible for the older version of gse that uses DataSet+BizID as channel name
					DataSet: "snapshot",
					BizID:   biz.BizID,
				},
			},
		}},
	}
	return snapChannel, nil
}

// gseConfigQueryRoute get channel by host snap data id from gse
func (s *Service) gseConfigQueryRoute(header http.Header, user string, defErr errors.DefaultCCErrorIf, rid string) (
	[]metadata.GseConfigChannel, error) {

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

	channels, err := esb.EsbClient().GseSrv().ConfigQueryRoute(s.ctx, header, params)
	if err != nil {
		blog.ErrorJSON("query route from gse failed, err: %s, params: %s, rid: %s", err, params, rid)
		return nil, &metadata.RespError{Msg: defErr.CCErrorf(common.CCErrCommMigrateFailed, err.Error())}
	}

	return channels, nil
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
