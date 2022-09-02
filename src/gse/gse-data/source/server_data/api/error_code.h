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

#ifndef _GSE_DATA_ERROR_CODE_H_
#define _GSE_DATA_ERROR_CODE_H_
#include "bbx/bbx.h"
#include <string>

namespace gse {
namespace data {

#define GSE_DATA_ERROR_CONFIG_SERVER_START GSE_DATA_ERROR_BASE
#define GSE_DATA_ERROR_CONFIG_SERVER_CHANNEL_ID_NOT_EXIST (GSE_DATA_ERROR_CONFIG_SERVER_START + 1)
#define GSE_DATA_ERROR_CONFIG_SERVER_CHANNEL_ID_EXIST (GSE_DATA_ERROR_CONFIG_SERVER_START + 2)
#define GSE_DATA_ERROR_STREAM_TO_ID_NOT_EXIST (GSE_DATA_ERROR_CONFIG_SERVER_START + 3)
#define GSE_DATA_ERROR_STREAM_TO_ID_EXIST (GSE_DATA_ERROR_CONFIG_SERVER_START + 4)

#define GSE_DATA_ERROR_REQUEST_PARAM_INVALID (GSE_DATA_ERROR_CONFIG_SERVER_START + 100)

#define GSE_DATA_ERROR_PARAM_ERROR (GSE_DATA_ERROR_CONFIG_SERVER_START + 100)
#define GSE_DATA_ERROR_INTERNAL_LOGIC_ERROR (GSE_DATA_ERROR_CONFIG_SERVER_START + 200)

#define GSE_DATA_ERROR_INTERNAL_ZK_ERROR (GSE_DATA_ERROR_CONFIG_SERVER_START + 300)

#define GSE_DATA_ERROR_RUNTIME_ERROR (GSE_DATA_ERROR_CONFIG_SERVER_START + 400)

#define GET_ERROR_MESSAGE(Errors) std::get<1>(Errors)
#define GET_ERROR_CODE(Errors) std::get<0>(Errors)

using ApiError = std::tuple<int, std::string>;

const std::tuple<int, std::string> ChannelIdExist(GSE_DATA_ERROR_CONFIG_SERVER_CHANNEL_ID_EXIST, std::string("channelid exist"));
const std::tuple<int, std::string> ChannelIdNotExist(GSE_DATA_ERROR_CONFIG_SERVER_CHANNEL_ID_NOT_EXIST, std::string("channelid not exist"));

const std::tuple<int, std::string> StreamToIdNotExist(GSE_DATA_ERROR_STREAM_TO_ID_NOT_EXIST, std::string("stream_to_id not exist"));

//-------------

const std::tuple<int, std::string> ChannelIdNotSet(GSE_DATA_ERROR_REQUEST_PARAM_INVALID, std::string("channel_id not seted"));
const std::tuple<int, std::string> TopNameInvalid(GSE_DATA_ERROR_PARAM_ERROR, std::string("topic name invalid"));

const std::tuple<int, std::string> ChannelNameInvalid(GSE_DATA_ERROR_PARAM_ERROR, std::string("redis channel name invalid"));
const std::tuple<int, std::string> FilterNameInvalid(GSE_DATA_ERROR_PARAM_ERROR, std::string("filter name invalid"));

const std::tuple<int, std::string> MetadataNotSet(GSE_DATA_ERROR_PARAM_ERROR, std::string("metadata is not seted"));
const std::tuple<int, std::string> MetadataJsonInvalid(GSE_DATA_ERROR_PARAM_ERROR, std::string("metadata's json invalid"));
const std::tuple<int, std::string> PlatnameNotSet(GSE_DATA_ERROR_PARAM_ERROR, std::string("plat_name is not seted"));
const std::tuple<int, std::string> PlatnameNotSupport(GSE_DATA_ERROR_PARAM_ERROR, std::string("unsupported platform name"));

const std::tuple<int, std::string> ChannelRouteNotSet(GSE_DATA_ERROR_PARAM_ERROR, std::string("route is not seted"));
const std::tuple<int, std::string> ChannelRouteNameNotSet(GSE_DATA_ERROR_PARAM_ERROR, std::string("route name is not seted"));

const std::tuple<int, std::string> AddressNotExist(GSE_DATA_ERROR_PARAM_ERROR, std::string("address is not set"));
const std::tuple<int, std::string> AddressNotARRAY(GSE_DATA_ERROR_PARAM_ERROR, std::string("address is not array"));

//参数错误
const std::tuple<int, std::string> ChannelIdOutofRange(GSE_DATA_ERROR_PARAM_ERROR, std::string("channel_id out of range[1, 2097151]"));
const std::tuple<int, std::string> ChannelIdNotBuildIn(GSE_DATA_ERROR_PARAM_ERROR, std::string("channelid not a built-in"));
const std::tuple<int, std::string> ChannelRouteNotARRAY(GSE_DATA_ERROR_PARAM_ERROR, std::string("route is not array"));
const std::tuple<int, std::string> FilterNotARRAY(GSE_DATA_ERROR_PARAM_ERROR, std::string("stream filter is not array"));
const std::tuple<int, std::string> FilterNameNotString(GSE_DATA_ERROR_PARAM_ERROR, std::string("stream filter name is not string"));
const std::tuple<int, std::string> ChannelNameNotString(GSE_DATA_ERROR_PARAM_ERROR, std::string("channel route name is not string"));
const std::tuple<int, std::string> ChannelNameRepeated(GSE_DATA_ERROR_PARAM_ERROR, std::string("channel name is repeated"));

const std::tuple<int, std::string> StreamFilterNameNotSet(GSE_DATA_ERROR_PARAM_ERROR, std::string("streamfilter name is not seted"));
const std::tuple<int, std::string> StreamFilterIndexInvalid(GSE_DATA_ERROR_PARAM_ERROR, std::string("streamfilter field_index invalid"));

const std::tuple<int, std::string> CondtionNotSet(GSE_DATA_ERROR_PARAM_ERROR, std::string("condition is not seted"));

const std::tuple<int, std::string> OperationNotSet(GSE_DATA_ERROR_PARAM_ERROR, std::string("operation is not seted"));
const std::tuple<int, std::string> OperationParamInvalid(GSE_DATA_ERROR_PARAM_ERROR, std::string("operation param type error"));

const std::tuple<int, std::string> SpecificationNotSet(GSE_DATA_ERROR_PARAM_ERROR, std::string("specification is not seted"));
const std::tuple<int, std::string> SpecificationNotSetStreamTo(GSE_DATA_ERROR_PARAM_ERROR, std::string("specification not set stream_to"));

const std::tuple<int, std::string> KafkaNotSet(GSE_DATA_ERROR_PARAM_ERROR, std::string("kafka not set"));
const std::tuple<int, std::string> PulsarNotSet(GSE_DATA_ERROR_PARAM_ERROR, std::string("pulsar not set"));
const std::tuple<int, std::string> RedisNotSet(GSE_DATA_ERROR_PARAM_ERROR, std::string("redis not set"));
const std::tuple<int, std::string> ProxyNotSet(GSE_DATA_ERROR_PARAM_ERROR, std::string("proxy not set"));

const std::tuple<int, std::string> UnsupportReportMode(GSE_DATA_ERROR_PARAM_ERROR, std::string("unsupport report_mode"));
const std::tuple<int, std::string> AddressInvalid(GSE_DATA_ERROR_PARAM_ERROR, std::string("address invalid"));
const std::tuple<int, std::string> PortInvalid(GSE_DATA_ERROR_PARAM_ERROR, std::string("port invalid"));
const std::tuple<int, std::string> UnsupportRedisClusterMode(GSE_DATA_ERROR_PARAM_ERROR, std::string("unsupport redis cluster mode"));

const std::tuple<int, std::string> StreamToNotSet(GSE_DATA_ERROR_PARAM_ERROR, std::string("stream_to not set"));
const std::tuple<int, std::string> StreamToIdNotSet(GSE_DATA_ERROR_PARAM_ERROR, std::string("stream_to_id not set"));
const std::tuple<int, std::string> UnsupportDeleteMethod(GSE_DATA_ERROR_PARAM_ERROR, std::string("unsupport delete method"));
const std::tuple<int, std::string> InputParamJsonInvalid(GSE_DATA_ERROR_PARAM_ERROR, std::string("input param json invalid"));

//内部逻辑错误

const std::tuple<int, std::string> NoPermissionUpdate(GSE_DATA_ERROR_INTERNAL_LOGIC_ERROR, std::string("No permission to update channel id"));
const std::tuple<int, std::string> NoPermissionDelete(GSE_DATA_ERROR_INTERNAL_LOGIC_ERROR, std::string("No permission to delete channel id"));

const std::tuple<int, std::string> ReadFromConfigJsonInvalid(GSE_DATA_ERROR_INTERNAL_LOGIC_ERROR, std::string("read from config, json invalid"));

//内部zk操作失败错误
const std::tuple<int, std::string> ZooKeeperUnit(GSE_DATA_ERROR_INTERNAL_ZK_ERROR, std::string("the zookeeper server is not connected"));
// const std::tuple<int, std::string> ZooKeeperSetError(GSE_DATA_ERROR_INTERNAL_ZK_ERROR, std::string("zookeeper set node failed"));
// const std::tuple<int, std::string> ZooKeeperCreateError(GSE_DATA_ERROR_INTERNAL_ZK_ERROR, std::string("zookeeper create node failed"));
const std::tuple<int, std::string> ZooKeeperGetError(GSE_DATA_ERROR_INTERNAL_ZK_ERROR, std::string("zookeeper get node failed"));
const std::tuple<int, std::string> ZooKeeperGetChildError(GSE_DATA_ERROR_INTERNAL_ZK_ERROR, std::string("zookeeper get child node failed"));
const std::tuple<int, std::string> ZooKeeperDeleteError(GSE_DATA_ERROR_INTERNAL_ZK_ERROR, std::string("zookeeper delete node failed"));
const std::tuple<int, std::string> ZooKeeperExistError(GSE_DATA_ERROR_INTERNAL_ZK_ERROR, std::string("zookeeper exist node failed"));

const std::tuple<int, std::string> ZooKeeperCreateError(GSE_DATA_ERROR_RUNTIME_ERROR, std::string("zookeeper create node failed"));
const std::tuple<int, std::string> ZooKeeperSetError(GSE_DATA_ERROR_RUNTIME_ERROR, std::string("zookeeper set node failed"));

class APIErrorCode
{
public:
};

} // namespace data
} // namespace gse

#endif // ERROR_CODE_H
