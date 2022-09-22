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

#include "channelid_request_api.h"

#include <json/json.h>

#include "api/api_metrics.h"
#include "bbx/gse_errno.h"
#include "channelid_zk_api.h"
#include "common/logger.hpp"
#include "log/log.h"
#include "tools/finally.hpp"
#include "tools/maths.h"
#include "tools/time.h"

namespace gse {
namespace data {

bool ChannelIdRequest::Valid(ApiError &error)
{
    if (!m_ptrMetadata->Valid(error))
    {
        return false;
    }

    if (!m_ptrOperation->Valid(error))
    {
        return false;
    }

    std::set<std::string> names;
    for (auto channel : m_channels)
    {
        if (!channel->Valid(error))
        {
            return false;
        }
        names.insert(channel->m_name);
    }

    if (names.size() < m_channels.size())
    {
        BLOG_ERROR(m_requestId, "channel name has repeat");
        error = ChannelNameRepeated;
        return false;
    }

    for (auto streamFilter : m_streamFilters)
    {
        if (!streamFilter->Valid(error))
        {
            return false;
        }
    }

    return true;
}

bool ChannelIdRequest::GenerateChannelId(ApiError &error)
{
    uint32_t channelId = m_ptrMetadata->m_channelID;

    if (channelId > 0)
    {
        std::string strChannelId = gse::tools::strings::ToString(m_ptrMetadata->m_channelID);
        if (GetRealChannelID(channelId) >= MAX_CHANNEL_ID_INDEX)
        {
            error = ChannelIdOutofRange;
            BLOG_ERROR(m_requestId, "the specified channelid (%u) out of range[1, 2097151]", channelId);
            return false;
        }

        BLOG_ERROR(m_requestId, "specified channelid (%u) for the plat(%s)", channelId, m_ptrMetadata->m_requestPlatName.c_str());
        if (ChannelIdZkApiInst()->ExistChannelId(strChannelId, error))
        {
            error = ChannelIdExist;
            BLOG_ERROR(m_requestId, "the specified channelid (%u) already exists , please update", channelId);
            return false;
        }

        if (!ChannelIdCommApiInst()->IsBuiltinChannelId(channelId))
        {
            error = ChannelIdNotBuildIn;
            BLOG_ERROR(m_requestId, "the specified channelid (%u) is not a built-in channelid", channelId);
            return false;
        }

        //        if (isBeyondLimit(channelId, m_ptrMetadata->m_requestPlatName))
        //        {
        //            char str_errmsg[200] = {0};
        //            snprintf(str_errmsg, sizeof(str_errmsg), "Channelid:%d out of range, bkmonitor{[%d,%d], [%d, %d]}, gse [%d, %d]", ptr_request->m_ptrMetadata->m_channelID,
        //                     MONITOR_RESERVER_CHANNEL_ID_START_1, MONITOR_RESERVER_CHANNEL_ID_END_1,
        //                     MONITOR_RESERVER_CHANNEL_ID_START_2, MONITOR_RESERVER_CHANNEL_ID_END_2,
        //                     GSE_RESERVER_CHANNEL_ID_START_1, GSE_RESERVER_CHANNEL_ID_END_1);

        //            return GSE_ERROR;
        //        }

        m_channelId = gse::tools::strings::ToString(m_ptrMetadata->m_channelID);
    }
    else
    {
        channelId = ChannelIdZkApiInst()->GenerateChannelId(m_ptrMetadata->m_requestPlatName, error);
        BLOG_INFO(m_requestId, "generate a new channel id (%u) for the plat(%s)", channelId, m_ptrMetadata->m_requestPlatName.c_str());
    }

    if (channelId < 0)
    {
        error = ChannelIdExist;
        BLOG_ERROR(m_requestId, "failed to alloc channelid, error:%s", GET_ERROR_MESSAGE(error).c_str());
        return false;
    }
    else
    {
        m_channelId = gse::tools::strings::ToString(channelId);
        m_ptrMetadata->m_channelID = channelId;
    }

    return true;
}

bool ChannelIdRequest::ParseRequest(const Json::Value &requestJson, ApiError &error)
{
    std::string errorMsg;
    // parse metadata
    if (!requestJson.isMember("metadata"))
    {
        error = MetadataNotSet;
        return false;
    }

    if (!m_ptrMetadata->ParseMetadata(requestJson["metadata"], error))
    {
        BLOG_ERROR(m_requestId, "failed to parse metadata, error:%s", GET_ERROR_MESSAGE(error).c_str());
        return false;
    }

    if (requestJson.isMember("operation"))
    {
        if (!m_ptrOperation->ParseOperation(requestJson["operation"], error))
        {
            return false;
        }
    }

    // parse route array
    if (requestJson.isMember("route"))
    {
        if (!requestJson["route"].isArray())
        {
            error = ChannelRouteNotARRAY;
            BLOG_ERROR(m_requestId, "failed to parse channel, the route is not an array, input json:%s", requestJson.toStyledString().c_str());
            return false;
        }

        for (int i = 0; i < requestJson["route"].size(); i++)
        {
            auto channel = std::make_shared<Channel>();
            if (!channel->Parse(requestJson["route"][i], error))
            {
                error = ChannelRouteNotARRAY;
                return false;
            }
            m_channels.push_back(channel);
        }
    }

    if (requestJson.isMember("stream_filters"))
    {
        if (!requestJson["stream_filters"].isArray())
        {
            error = FilterNotARRAY;
            BLOG_ERROR(m_requestId, "failed to parse stream_filters, stream_filters is not an array, json:%s", requestJson.toStyledString().c_str());
            return false;
        }

        for (int i = 0; i < requestJson["stream_filters"].size(); i++)
        {
            auto streamFilter = std::make_shared<StreamFilter>();
            if (!streamFilter->Parse(requestJson["stream_filters"][i], error))
            {
                BLOG_ERROR(m_requestId, "failed to parse stream_filter, json:%s", requestJson["stream_filters"][i].toStyledString().c_str());
                return false;
            }

            m_streamFilters.push_back(streamFilter);
        }
    }

    return true;
}

bool ChannelIdRequest::ToZK(ApiError &error)
{
    if (!ChannelIdZkApiInst()->CreateChannelIdPath(m_channelId, error))
    {
        BLOG_ERROR(m_requestId, "failed to create channelid path, channelid:%s, error:%s", m_channelId.c_str(), GET_ERROR_MESSAGE(error).c_str());
        return false;
    }

    if (!ChannelIdZkApiInst()->CreateChannels(m_channelId, m_channels, error))
    {
        BLOG_ERROR(m_requestId, "failed to create channel, channelid:%s, error:%s", m_channelId.c_str(), GET_ERROR_MESSAGE(error).c_str());
        return false;
    }

    if (!ChannelIdZkApiInst()->CreateMetaData(m_channelId, m_ptrMetadata, error))
    {
        BLOG_ERROR(m_requestId, "failed to create metadata, channelid:%s, error:%s", m_channelId.c_str(), GET_ERROR_MESSAGE(error).c_str());
        return false;
    }

    if (!ChannelIdZkApiInst()->CreateFilters(m_channelId, m_streamFilters, error))
    {
        BLOG_ERROR(m_requestId, "failed to create filter, channelid:%s, error:%s", m_channelId.c_str(), GET_ERROR_MESSAGE(error).c_str());
        return false;
    }

    // update value to notify
    if (!ChannelIdZkApiInst()->UpdateChannelIdTime(m_channelId, error))
    {
        BLOG_ERROR(m_requestId, "failed to update channelid value, channelid:%s, error:%s", m_channelId.c_str(), GET_ERROR_MESSAGE(error).c_str());
        return false;
    }

    for (auto channel : m_channels)
    {
        std::string strStreamToId = gse::tools::strings::ToString(channel->m_streamTo.m_streamToId);
        ChannelIdZkApiInst()->CreateChannelIdQueryIndex(m_ptrMetadata, strStreamToId, error);
    }

    if (ChannelIdCommApiInst()->IsTglogPlat(m_ptrMetadata->m_requestPlatName))
    {
        // TGLOG的channelid存在，不需要再创建了，直接返回已存在的dataid
        uint32_t iChannelid = 0;
        if (ChannelIdZkApiInst()->GetTglogChannelId(m_ptrMetadata, iChannelid, error))
        {
            m_channelId = gse::tools::strings::ToString(iChannelid);
            return true;
        }

        if (!ChannelIdZkApiInst()->SaveTlogChannelId(m_ptrMetadata, error))
        {
            return false;
        }

        return true;
    }
    else
    {
        std::string dataIdJsonStr;
        if (ToDataId(dataIdJsonStr))
        {
            if (!SaveToDataIdZk(m_channelId, dataIdJsonStr, error))
            {
                BLOG_ERROR(m_requestId, "failed to create dataid, error:%s", GET_ERROR_MESSAGE(error).c_str());
                return false;
            }
        }
    }

    return true;
}

bool ChannelIdRequest::ToDataId(std::string &jsonStr)
{
    Json::Value dataIdJson;
    Json::Value dataIds;

    for (auto channel : m_channels)
    {
        std::string reportMode = channel->m_streamTo.m_reportMode;
        if (reportMode == ChannelIDOperationRequestMethodKafka)
        {
            KafkaTopic &kafkaTopic = channel->m_streamTo.m_kafkaTopic;
            dataIdJson["data_set"] = Json::Value(kafkaTopic.m_dataSet);
            dataIdJson["partition"] = Json::Value(kafkaTopic.m_partition);
            dataIdJson["biz_id"] = Json::Value(kafkaTopic.m_bizId);
            dataIdJson["cluster_index"] = Json::Value(channel->m_streamToId);
            dataIdJson["server_id"] = Json::Value(-1);
            dataIdJson["msg_system"] = Json::Value(KAFKA_COMMON);
            dataIds.append(dataIdJson);
        }
        else if (reportMode == ChannelIDOperationRequestMethodRedis)
        {
            RedisChannel &redisChannel = channel->m_streamTo.m_redisChannel;
            dataIdJson["data_set"] = Json::Value(redisChannel.m_dataSet);
            dataIdJson["partition"] = Json::Value(0);
            dataIdJson["biz_id"] = Json::Value(redisChannel.m_bizId);
            dataIdJson["cluster_index"] = Json::Value(channel->m_streamToId);
            dataIdJson["server_id"] = Json::Value(-1);
            dataIdJson["msg_system"] = Json::Value(REDIS_SENTINEL_PUB);
            dataIds.append(dataIdJson);
        }
    }

    jsonStr.assign(dataIds.toStyledString());
    if (dataIds.size() > 0)
    {
        BLOG_INFO(m_requestId, "dataid(%s) config value(%s)", m_channelId.c_str(), jsonStr.c_str());
        return true;
    }
    return false;
}

Json::Value ChannelIdRequest::ToJSONValue()
{
    Json::Value data;

    data["metadata"] = m_ptrMetadata->ToJSONValue();
    data["operation"] = m_ptrOperation->ToJSONValue();
    data["stream_filters"].resize(0);
    data["route"].resize(0);

    for (auto channel : m_channels)
    {
        data["route"].append(channel->ToJSONValue());
    }

    for (auto filter : m_streamFilters)
    {
        data["stream_filters"].append(filter->ToJSONValue());
    }

    return data;
}

void ChannelIdRequest::MakeResponse(std::string &response, ApiError &error)
{
    Json::Value data;
    data["channel_id"] = Json::Value(gse::tools::strings::StringToInt32(m_channelId));
    uint64_t proofingId = ((uint64_t)(gse::tools::strings::StringToUint64(m_channelId) << 32) | gse::tools::time::GetUTCSecond());
    // data["proofing_id"] = Json::Value(proofingId);

    Json::Value responseJson;
    ChannelIdCommonAPI::MakeResponse(responseJson, error, data);
    response.assign(responseJson.toStyledString());
}

ChannelIdRequestAPIHandler::ChannelIdRequestAPIHandler(const std::string &version)
{
}
ChannelIdRequestAPIHandler::~ChannelIdRequestAPIHandler()
{
}

bool ChannelIdRequestAPIHandler::HandleRequest(gse::net::http::HTTPMessagePtr message, ApiError &error, std::string &response)
{
    Json::Value requestJson;
    Json::Reader reader(Json::Features::strictMode());
    std::shared_ptr<ChannelIdRequest> ptrRequest = std::make_shared<ChannelIdRequest>(m_requestId);

    auto _ = gse::tools::defer::finally([&]() {
        ptrRequest->MakeResponse(response, error);
        BLOG_INFO(m_requestId, "send response uri:%s, reponse(%s)", message->GetURI().c_str(), response.c_str());
    });

    if (!reader.parse(message->GetBody().c_str(), requestJson, false))
    {
        BLOG_ERROR(m_requestId, "failed to parse request json,json invalid, request(%s)", message->GetBody().c_str());
        error = InputParamJsonInvalid;
        return false;
    }

    if (!ptrRequest->ParseRequest(requestJson, error))
    {
        BLOG_ERROR(m_requestId, "failed to check param, error:%s, request(%s)", GET_ERROR_MESSAGE(error).c_str(), message->GetBody().c_str());
        return false;
    }

    if (!ptrRequest->Valid(error))
    {
        BLOG_ERROR(m_requestId, "failed to check param, error:%s, request(%s)", GET_ERROR_MESSAGE(error).c_str(), message->GetBody().c_str());
        return false;
    }

    if (!ptrRequest->GenerateChannelId(error))
    {
        BLOG_ERROR(m_requestId, "failed to gennerate channelid, error:%s", GET_ERROR_MESSAGE(error).c_str());
        return false;
    }

    std::string errorMsg;
    if (!ptrRequest->ToZK(error))
    {
        BLOG_ERROR(m_requestId, "failed to write to zk, error:%s, request(%s)", GET_ERROR_MESSAGE(error).c_str(), message->GetBody().c_str());
        return false;
    }

    return true;
}

int ChannelIdRequestAPIHandler::OnPost(gse::net::http::HTTPMessagePtr message, std::string &response) noexcept
{
    auto headerReqID = message->GetHeaderParameter(BK_APIGW_REQ_ID);
    m_requestId = headerReqID ? headerReqID->GetValue() : BK_API_UNKOWN_REO_ID;
    BLOG_INFO(m_requestId, "request uri:%s, message(%s)", message->GetURI().c_str(), message->GetBody().c_str());

    uint64_t beginTimestamp = gse::tools::time::GetUTCMillsec();
    ApiError error;
    auto _ = gse::tools::defer::finally([&]() {
        uint64_t endTimestamp = gse::tools::time::GetUTCMillsec();
        uint64_t costTime = (endTimestamp >= beginTimestamp) ? (endTimestamp >= beginTimestamp) : 0;
        ChannelIdApiMetricsInst()->CounterMetrics(m_metricsReponseMillseconds, m_metricsApiName, costTime, GET_ERROR_CODE(error));
        ChannelIdApiMetricsInst()->CounterMetrics(m_metricsApiCount, m_metricsApiName, 1, GET_ERROR_CODE(error));
    });

    if (!HandleRequest(message, error, response))
    {
        BLOG_ERROR(m_requestId, "failed to handle channeld request, error:%s", GET_ERROR_MESSAGE(error).c_str());
        return HTTP_STATUS_CODE_400;
    }

    return HTTP_STATUS_CODE_200;
}

} // namespace data
} // namespace gse
