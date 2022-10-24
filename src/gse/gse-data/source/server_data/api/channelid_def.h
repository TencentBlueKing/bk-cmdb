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

#ifndef _GSE_DATA_CHANNELID_DEF_H_
#define _GSE_DATA_CHANNELID_DEF_H_

#include <string>

#define MONITOR_RESERVER_CHANNEL_ID_START_1 1000
#define MONITOR_RESERVER_CHANNEL_ID_END_1 1020

#define MONITOR_RESERVER_CHANNEL_ID_START_2 1100000
#define MONITOR_RESERVER_CHANNEL_ID_END_2 1199999

#define GSE_RESERVER_CHANNEL_ID_START_1 1200000
#define GSE_RESERVER_CHANNEL_ID_END_1 1210000

#ifndef ZK_CHANNEL_ID_OFFSET
#define ZK_CHANNEL_ID_OFFSET 0x80000
#endif

#ifndef ZK_STREAM_TO_ID_OFFSET
#define ZK_STREAM_TO_ID_OFFSET 0x400
#endif

// config tglog/659 or tdm
#ifndef ZK_CHANNEL_ID_CONFIG_BASE_PATH
#define ZK_CHANNEL_ID_CONFIG_BASE_PATH (std::string("/gse/config/server/dataserver/channelid"))
#endif

#ifndef ZK_CHANNEL_ID_CONFIG_PATH
#define ZK_CHANNEL_ID_CONFIG_PATH(id) (std::string("/gse/config/server/dataserver/channelid/") + std::string(id))
#endif

#ifndef ZK_CHANNEL_ID_CONFIG_FILTER_PATH
#define ZK_CHANNEL_ID_CONFIG_FILTER_PATH(channelID) (ZK_CHANNEL_ID_CONFIG_PATH(channelID) + "/filter")
#endif

#ifndef ZK_CHANNEL_ID_CONFIG_METADATA_PATH
#define ZK_CHANNEL_ID_CONFIG_METADATA_PATH(channelID) (ZK_CHANNEL_ID_CONFIG_PATH(channelID) + "/metadata")
#endif

#ifndef ZK_CHANNEL_ID_CONFIG_WITH_FILTERNAME_PATH
#define ZK_CHANNEL_ID_CONFIG_WITH_FILTERNAME_PATH(channelID, filterName) (ZK_CHANNEL_ID_CONFIG_FILTER_PATH(channelID) + "/" + std::string(filterName))
#endif

#ifndef ZK_CHANNEL_ID_CONFIG_CHANNEL_PATH
#define ZK_CHANNEL_ID_CONFIG_CHANNEL_PATH(channelID) (ZK_CHANNEL_ID_CONFIG_PATH(channelID) + "/channel")
#endif

#ifndef ZK_CHANNEL_ID_CONFIG_WITH_CHANNELNAME_PATH
#define ZK_CHANNEL_ID_CONFIG_WITH_CHANNELNAME_PATH(channelID, channelName) (ZK_CHANNEL_ID_CONFIG_CHANNEL_PATH(channelID) + "/" + std::string(channelName))
#endif

#ifndef ZK_PLAT_ID_CONFIG_BASE_PATH
#define ZK_PLAT_ID_CONFIG_BASE_PATH (std::string("/gse/config/server/dataserver/platids"))
#endif

#ifndef ZK_PLAT_ID_CONFIG_PATH
#define ZK_PLAT_ID_CONFIG_PATH(platId) (ZK_PLAT_ID_CONFIG_BASE_PATH + "/" + std::string(platId))
#endif

#ifndef ZK_PLAT_ID_CONFIG_FILTER_PATH
#define ZK_PLAT_ID_CONFIG_FILTER_PATH(platId) (ZK_PLAT_ID_CONFIG_PATH(platId) + "/filter")
#endif

#ifndef ZK_PLAT_ID_CONFIG_METADATA_PATH
#define ZK_PLAT_ID_CONFIG_METADATA_PATH(platId) (ZK_PLAT_ID_CONFIG_PATH(platId) + "/metadata")
#endif

#ifndef ZK_PLAT_ID_CONFIG_WITH_FILTERNAME_PATH
#define ZK_PLAT_ID_CONFIG_WITH_FILTERNAME_PATH(platId, filterName) (ZK_PLAT_ID_CONFIG_PATH(platId) + "/" + std::string(filterName))
#endif

#ifndef ZK_PLAT_ID_CONFIG_EXPORTER_PATH
#define ZK_PLAT_ID_CONFIG_EXPORTER_PATH(platId) (ZK_PLAT_ID_CONFIG_PATH(platId) + "/exporter")
#endif

#ifndef ZK_PLAT_ID_CONFIG_WITH_EXPORTERNAME_PATH
#define ZK_PLAT_ID_CONFIG_WITH_EXPORTERNAME_PATH(platId, exporterName) (ZK_PLAT_ID_CONFIG_PATH(platId) + "/" + std::string(exporterName))
#endif

#ifndef ZK_PLAT_ID_CONFIG_CHANNEL_PATH
#define ZK_PLAT_ID_CONFIG_CHANNEL_PATH(platId) (ZK_PLAT_ID_CONFIG_PATH(platId) + "/channel")
#endif

#ifndef ZK_PLAT_ID_CONFIG_WITH_CHANNELNAME_PATH
#define ZK_PLAT_ID_CONFIG_WITH_CHANNELNAME_PATH(platId, channelName) (ZK_PLAT_ID_CONFIG_PATH(platId) + "/" + std::string(channelName))
#endif

// TGLOG channelid watch path

#ifndef ZK_CHANNEL_ID_TGLOG_BASE_PATH
#define ZK_CHANNEL_ID_TGLOG_BASE_PATH std::string("/gse/config/server/dataserver/tglog/")
#endif

#ifndef ZK_CHANNEL_ID_TGLOG_NOTIFY_PATH
#define ZK_CHANNEL_ID_TGLOG_NOTIFY_PATH(watch_path) (std::string("/gse/config/server/dataserver/tglog/") + watch_path)
#endif

//--------------------------------------------------
// streamId
#ifndef ZK_STREAM_ID_CONFIG_BASE_PATH
#define ZK_STREAM_ID_CONFIG_BASE_PATH std::string("/gse/config/server/dataserver/streamto")
#endif

#ifndef ZK_STREAM_ID_CONFIG_PATH
#define ZK_STREAM_ID_CONFIG_PATH(streamToId) ZK_STREAM_ID_CONFIG_BASE_PATH + "/" + std::string(streamToId)
#endif

#ifndef ZK_STREAM_TO_ID_CONFIG_METADATA_PATH
#define ZK_STREAM_TO_ID_CONFIG_METADATA_PATH(streamToId) (ZK_STREAM_ID_CONFIG_PATH(streamToId) + "/metadata")
#endif

// ---------------分割线-------------------------
#ifndef ZK_CHANNEL_ID_WATCH_BASE_PATH
#define ZK_CHANNEL_ID_WATCH_BASE_PATH "/gse/config/server/dataserver/channelids"
#endif

#ifndef ZK_CHANNEL_ID_WATCH_PATH
#define ZK_CHANNEL_ID_WATCH_PATH(channelID) (std::string("/gse/config/server/dataserver/channelids/") + std::string(channelID))
#endif

#ifndef ZK_PLAT_ID_WATCH_PATH
#define ZK_PLAT_ID_WATCH_PATH(platid) (std::string("/gse/config/server/dataserver/platids/") + std::string(platid))
#endif

const std::string kSupportPlats = "[tglog|tdm|bkmonitor|tgdp|iegdata|cmdb|gse]";

//-------------------------------------------------------------------
// query index

#ifndef ZK_CHANNEL_ID_CONFIG_WITH_LABLE_PATH
#define ZK_CHANNEL_ID_CONFIG_WITH_LABLE_PATH() (std::string("/gse/config/server/dataserver/channelid/index"))
#endif

#ifndef ZK_CHANNEL_ID_CONFIG_WITH_LABLE_BIZID_PATH
#define ZK_CHANNEL_ID_CONFIG_WITH_LABLE_BIZID_PATH(bizid) (ZK_CHANNEL_ID_CONFIG_WITH_LABLE_PATH() + "/bk_biz_id/" + std::string(bizid))
#endif

#ifndef ZK_CHANNEL_ID_CONFIG_WITH_LABLE_ODM_PATH
#define ZK_CHANNEL_ID_CONFIG_WITH_LABLE_ODM_PATH(odm) (ZK_CHANNEL_ID_CONFIG_WITH_LABLE_PATH() + "/odm/" + std::string(odm))
#endif

#ifndef ZK_CHANNEL_ID_CONFIG_WITH_TYPE_PATH
#define ZK_CHANNEL_ID_CONFIG_WITH_TYPE_PATH(type) (ZK_CHANNEL_ID_CONFIG_WITH_LABLE_PATH() + "/type/" + std::string(type))
#endif

#ifndef ZK_CHANNEL_ID_CONFIG_WITH_LABLE_PLATNAME_PATH
#define ZK_CHANNEL_ID_CONFIG_WITH_LABLE_PLATNAME_PATH(platname) (ZK_CHANNEL_ID_CONFIG_WITH_LABLE_PATH() + "/plat_name/" + std::string(platname))
#endif

#ifndef ZK_CHANNEL_ID_CONFIG_WITH_LABLE_STREAMTO_ID_PATH
#define ZK_CHANNEL_ID_CONFIG_WITH_LABLE_STREAMTO_ID_PATH(streamToId, channelId) (ZK_CHANNEL_ID_CONFIG_WITH_LABLE_PATH() + "/stream_to_id/" + std::string(streamToId) + "/" + std::string(channelId))
#endif

#ifndef ZK_CHANNEL_ID_ORIGIN_PATH
//#define ZK_CHANNEL_ID_ORIGIN_PATH(platName) (std::string("/gse/config/server/dataserver/etc/channelid/origin/") + std::string(platName))
/// orgin channelid
#define ZK_CHANNEL_ID_ORIGIN_PATH() std::string("/gse/config/server/dataserver/etc/channelid/origin")
#endif

#ifndef ZK_STREAM_TO_ID_ORIGIN_PATH
#define ZK_STREAM_TO_ID_ORIGIN_PATH() std::string("/gse/config/server/dataserver/etc/streamto/origin")
#endif

//--------------------------------------------------
#ifndef ZK_STREAM_ID_CONFIG_PATH
#define ZK_STREAM_ID_CONFIG_PATH(streamToId) (std::string("/gse/config/server/dataserver/streamto/") + std::string(streamToId))
#endif

#ifndef ZK_STREAM_TO_ID_CONFIG_METADATA_PATH
#define ZK_STREAM_TO_ID_CONFIG_METADATA_PATH(streamToId) (ZK_STREAM_ID_CONFIG_PATH(streamToId) + "/metadata")
#endif

#ifndef ZK_STREAM_TO_ID_CONFIG_STREAM_PATH
#define ZK_STREAM_TO_ID_CONFIG_STREAM_PATH(streamToId) (ZK_STREAM_ID_CONFIG_PATH(streamToId))
#endif

//-------------------------------------------------------------------
// query index

#ifndef ZK_STREAM_TO_ID_CONFIG_WITH_LABLE_PATH
#define ZK_STREAM_TO_ID_CONFIG_WITH_LABLE_PATH() (std::string("/gse/config/server/dataserver/streamto/index"))
#endif

#ifndef ZK_STREAM_TO_ID_CONFIG_WITH_LABLE_BIZID_PATH
#define ZK_STREAM_TO_ID_CONFIG_WITH_LABLE_BIZID_PATH(bizid) (ZK_STREAM_TO_ID_CONFIG_WITH_LABLE_PATH() + "/bk_biz_id/" + std::string(bizid))
#endif

#ifndef ZK_STREAM_TO_ID_CONFIG_WITH_LABLE_ODM_PATH
#define ZK_STREAM_TO_ID_CONFIG_WITH_LABLE_ODM_PATH(odm) (ZK_STREAM_TO_ID_CONFIG_WITH_LABLE_PATH() + "/odm/" + std::string(odm))
#endif

#ifndef ZK_STREAM_TO_ID_CONFIG_WITH_LABLE_PLATNAME_PATH
#define ZK_STREAM_TO_ID_CONFIG_WITH_LABLE_PLATNAME_PATH(platname) (ZK_STREAM_TO_ID_CONFIG_WITH_LABLE_PATH() + "/plat_name/" + std::string(platname))
#endif

#ifndef ZK_STREAM_TO_ID_CONFIG_WITH_TYPE_PATH
#define ZK_STREAM_TO_ID_CONFIG_WITH_TYPE_PATH(type) (ZK_STREAM_TO_ID_CONFIG_WITH_LABLE_PATH() + "/type/" + std::string(type))
#endif

#ifndef ZK_CHANNEL_ID_PLAT_PATH
#define ZK_CHANNEL_ID_PLAT_PATH "/gse/config/server/configserver/etc/channelid/plats"
#endif

#ifndef ZK_CHANNEL_ID_PLAT_PATH_WITH_NAME
#define ZK_CHANNEL_ID_PLAT_PATH_WITH_NAME (std::string("/gse/config/server/configserver/etc/channelid/plats"))
#endif

#ifndef ZK_CHANNEL_ID_TGLOG_NOTIFY_BASE_PATH
#define ZK_CHANNEL_ID_TGLOG_NOTIFY_BASE_PATH (std::string("/gse/config/server/dataserver/tglog"))
#endif

//-----------------------------------
#ifndef MakeChannelID
#define MakeChannelID(platNum, originChannelID) ((platNum << 20) | originChannelID) // 高12位作为 平台ID
#endif

//-----------------------------------
#ifndef MakeUnifyChannelID
#define MakeUnifyChannelID(platID, originChannelID) MakeChannelID(platID, (ZK_CHANNEL_ID_OFFSET + originChannelID))
#endif

#ifndef MakeUnifyStreamToID
#define MakeUnifyStreamToID(originStreamToId) (ZK_STREAM_TO_ID_OFFSET + originStreamToId)
#endif

#ifndef GetOriginStreamToID
#define GetOriginStreamToID(streamToId) ((streamToId >= ZK_STREAM_TO_ID_OFFSET) ? (streamToId - ZK_STREAM_TO_ID_OFFSET) : streamToId)
#endif

// high 12bits plat id
#ifndef GetPlatNum
#define GetPlatNum(channelID) ((channelID >> 20) & 0x3FF)
#endif

// low 20bits real id
#ifndef GetRealChannelID
#define GetRealChannelID(channelID) (channelID & 0xFFFFF)
#endif

#ifndef SUPPORT_PLAT_NAME_TGDP
#define SUPPORT_PLAT_NAME_TGDP "tgdp"
#endif

#ifndef SUPPORT_PLAT_NAME_TGLOG
#define SUPPORT_PLAT_NAME_TGLOG "tglog"
#endif

#ifndef SUPPORT_PLAT_NAME_TDM
#define SUPPORT_PLAT_NAME_TDM "tdm"
#endif

#ifndef SUPPORT_PLAT_NAME_BKMONITOR
#define SUPPORT_PLAT_NAME_BKMONITOR "bkmonitor"
#endif

#ifndef SUPPORT_PLAT_NAME_GSE
#define SUPPORT_PLAT_NAME_GSE "gse"
#endif

#ifndef SUPPORT_PLAT_NAME_DATAMORE
#define SUPPORT_PLAT_NAME_DATAMORE "datemore"
#endif

#ifndef SUPPORT_PLAT_NAME_IEG_DATA
#define SUPPORT_PLAT_NAME_IEG_DATA "iegdata"
#endif

#ifndef SUPPORT_PLAT_NAME_CMDB
#define SUPPORT_PLAT_NAME_CMDB "cmdb"
#endif

// TDM 不支持动态分配channelid，所有tdm业务都是 1
#ifndef PLAT_TDM_ID
#define PLAT_TDM_ID 0
#endif

#ifndef PLAT_TGDP_ID
#define PLAT_TGDP_ID 0
#endif

#ifndef PLAT_TGLOG_ID
#define PLAT_TGLOG_ID 3
#endif

#ifndef PLAT_IEG_DATA_ID
#define PLAT_IEG_DATA_ID 2
#endif

#ifndef PLAT_BKMONITOR_ID
#define PLAT_BKMONITOR_ID 1
#endif

#ifndef PLAT_GSE_ID
#define PLAT_GSE_ID 4
#endif

#ifndef PLAT_DATAMORE_ID
#define PLAT_DATAMORE_ID 5
#endif

#ifndef PLAT_CMDB_ID
#define PLAT_CMDB_ID 6
#endif

#ifndef ChannelIDOperationRequestMethodFile
#define ChannelIDOperationRequestMethodFile "file"
#endif

#ifndef ChannelIDOperationRequestMethodKafka
#define ChannelIDOperationRequestMethodKafka "kafka"
#endif

#ifndef ChannelIDOperationRequestMethodPulsar
#define ChannelIDOperationRequestMethodPulsar "pulsar"
#endif

#ifndef ChannelIDOperationRequestMethodRedis
#define ChannelIDOperationRequestMethodRedis "redis"
#endif

#ifndef ChannelIDOperationRequestMethodDsProxy
#define ChannelIDOperationRequestMethodDsProxy "proxy"
#endif

#ifndef ChannelIDOperationUpdateMethodAppend
#define ChannelIDOperationUpdateMethodAppend "append"
#endif

#ifndef ChannelIDOperationUpdateMethodSet
#define ChannelIDOperationUpdateMethodSet "set"
#endif

#ifndef ChannelIDOperationDeleteMethodAll
#define ChannelIDOperationDeleteMethodAll "all"
#endif

#ifndef ChannelIDOperationDeleteMethodSpecification
#define ChannelIDOperationDeleteMethodSpecification "specification"
#endif

#ifndef SUPPORT_PLAT_NAME_TGDP
#define SUPPORT_PLAT_NAME_TGDP "tgdp"
#endif

#ifndef SUPPORT_PLAT_NAME_TGLOG
#define SUPPORT_PLAT_NAME_TGLOG "tglog"
#endif

#ifndef SUPPORT_PLAT_NAME_TDM
#define SUPPORT_PLAT_NAME_TDM "tdm"
#endif

#ifndef SUPPORT_PLAT_NAME_BKMONITOR
#define SUPPORT_PLAT_NAME_BKMONITOR "bkmonitor"
#endif

#ifndef SUPPORT_PLAT_NAME_GSE
#define SUPPORT_PLAT_NAME_GSE "gse"
#endif

#ifndef SUPPORT_PLAT_NAME_DATAMORE
#define SUPPORT_PLAT_NAME_DATAMORE "datemore"
#endif

#ifndef SUPPORT_PLAT_NAME_IEG_DATA
#define SUPPORT_PLAT_NAME_IEG_DATA "iegdata"
#endif

#ifndef SUPPORT_PLAT_NAME_CMDB
#define SUPPORT_PLAT_NAME_CMDB "cmdb"
#endif

// TDM 不支持动态分配channelid，所有tdm业务都是 1
#ifndef PLAT_TDM_ID
#define PLAT_TDM_ID 0
#endif

#ifndef PLAT_TGDP_ID
#define PLAT_TGDP_ID 0
#endif

#ifndef PLAT_TGLOG_ID
#define PLAT_TGLOG_ID 3
#endif

#ifndef PLAT_IEG_DATA_ID
#define PLAT_IEG_DATA_ID 2
#endif

#ifndef PLAT_BKMONITOR_ID
#define PLAT_BKMONITOR_ID 1
#endif

#ifndef PLAT_GSE_ID
#define PLAT_GSE_ID 4
#endif

#ifndef PLAT_DATAMORE_ID
#define PLAT_DATAMORE_ID 5
#endif

#ifndef PLAT_CMDB_ID
#define PLAT_CMDB_ID 6
#endif

#define BK_ERROR_CODE(errorcode) (-1 == (errorcode) ? 1000101 : (0 < (errorcode) && 1000000 > (errorcode) ? (1000000 + (errorcode)) : (errorcode)))

#define MAX_CHANNEL_ID_INDEX 0xFFFFF

#define BK_API_REO_ID "X-Bkapi-Request-Id"
#define BK_API_UNKOWN_REO_ID "Unkonw-Request-Id"

#endif // CHANNELID_DEF_H
