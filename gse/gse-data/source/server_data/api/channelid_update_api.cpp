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
#include "channelid_update_api.h"

#include <json/json.h>

#include "api/api_metrics.h"
#include "channelid_zk_api.h"
#include "common/logger.hpp"
#include "error_code.h"
#include "log/log.h"
#include "tools/finally.hpp"
#include "tools/maths.h"

namespace gse {
namespace data {

ChannelIdUpdate::ChannelIdUpdate(const std::string &requestId)
    : m_requestId(requestId)
{
    m_ptrCondition = std::make_shared<Condition>();
    m_ptrOperation = std::make_shared<Operation>();
    m_ptrSpecification = std::make_shared<UpdateSpecification>();
}

bool ChannelIdUpdate::Parse(const Json::Value &requestJson, ApiError &error)
{
    if (!requestJson.isMember("condition"))
    {
        error = CondtionNotSet;
        return false;
    }

    if (!requestJson.isMember("specification"))
    {
        error = SpecificationNotSet;
        return false;
    }

    // optional
    if (!requestJson.isMember("operation"))
    {
        error = OperationNotSet;
        return false;
    }

    if (!m_ptrCondition->Parse(requestJson["condition"], error))
    {
        BLOG_ERROR(m_requestId, "failed to parse conditon, error:%s", GET_ERROR_MESSAGE(error).c_str());
        return false;
    }

    if (requestJson.isMember("operation"))
    {
        m_ptrOperation->ParseOperation(requestJson["operation"], error);
    }

    if (!m_ptrSpecification->Parse(requestJson["specification"], error))
    {
        BLOG_ERROR(m_requestId, "failed to parse specification, error:%s", GET_ERROR_MESSAGE(error).c_str());
        return false;
    }

    return true;
}

bool ChannelIdUpdate::Valid(ApiError &error)
{

    if (!m_ptrOperation->Valid(error))
    {
        BLOG_ERROR(m_requestId, "failed to check operation, error:%s", GET_ERROR_MESSAGE(error).c_str());
        return false;
    }

    for (auto streamFilter : m_ptrSpecification->m_streamFilters)
    {
        if (!streamFilter->Valid(error))
        {
            BLOG_ERROR(m_requestId, "failed to check stream filter, error:%s", GET_ERROR_MESSAGE(error).c_str());
            return false;
        }
    }

    for (auto channel : m_ptrSpecification->m_channels)
    {
        if (!channel->Valid(error))
        {
            BLOG_ERROR(m_requestId, "failed to check stream filter, error:%s", GET_ERROR_MESSAGE(error).c_str());
            return false;
        }
    }
    return true;
}

bool ChannelIdUpdate::ToZK(ApiError &error)
{
    int64_t channelId = m_ptrCondition->m_channelID;
    std::string strChannelId = gse::tools::strings::ToString(channelId);

    auto ptrChannelIdConfg = ChannelIdZkApiInst()->ReadChannelIdConfig(strChannelId, m_ptrCondition->m_requestPlatName, error);
    if (ptrChannelIdConfg == nullptr)
    {
        BLOG_ERROR(m_requestId, "failed to read channelid(%d) config, error:%s", channelId, GET_ERROR_MESSAGE(error).c_str());
        return false;
    }

    if (m_ptrCondition->m_requestPlatName != ptrChannelIdConfg->m_ptrMetadata->m_requestPlatName)
    {
        BLOG_ERROR(m_requestId, "channel id(%d) does not belong to the platform(%s), real plat name(%s)",
                   channelId, m_ptrCondition->m_requestPlatName.c_str(),
                   ptrChannelIdConfg->m_ptrMetadata->m_requestPlatName.c_str());
        error = NoPermissionUpdate;
        return false;
    }

    for (auto &channel : m_ptrSpecification->m_channels)
    {
        bool bFound = false;
        for (auto &channelLocal : ptrChannelIdConfg->m_channels)
        {
            BLOG_INFO(m_requestId, "channel name:%s, specification name:%s",
                      channel->m_name.c_str(), channelLocal->m_name.c_str());
            if (channelLocal->m_name == channel->m_name)
            {
                channelLocal->CopyFrom(channel.get());
                bFound = true;
                break;
            }
        }

        if (!bFound)
        {
            BLOG_INFO(m_requestId, "not find stream config rule name:%s, add new config", channel->m_name.c_str());
            ptrChannelIdConfg->m_channels.push_back(channel->Clone());
        }
    }

    for (auto &filter : m_ptrSpecification->m_streamFilters)
    {
        bool bFound = false;
        for (auto &filterLocal : ptrChannelIdConfg->m_streamFilters)
        {
            BLOG_INFO(m_requestId, "filter name:%s,type:%s, specification filter name:%s", filter->m_name.c_str(), filterLocal->m_name.c_str());
            if (filter->m_name == filterLocal->m_name)
            {
                filterLocal->CopyFrom(filter.get());
                bFound = true;
                break;
            }
        }

        if (!bFound)
        {
            BLOG_INFO(m_requestId, "not find filter name:%s, add new config", filter->m_name.c_str());
            ptrChannelIdConfg->m_streamFilters.push_back(filter->Clone());
        }
    }

    if (!ChannelIdZkApiInst()->CreateMetaData(strChannelId, ptrChannelIdConfg->m_ptrMetadata, error, false))
    {
        return false;
    }

    if (!ChannelIdZkApiInst()->CreateFilters(strChannelId, ptrChannelIdConfg->m_streamFilters, error, false))
    {
        return false;
    }

    if (!ChannelIdZkApiInst()->CreateChannels(strChannelId, ptrChannelIdConfg->m_channels, error, false))
    {
        return false;
    }
    if (!ChannelIdZkApiInst()->UpdateChannelIdTime(strChannelId, error))
    {
        return false;
    }

    std::string strDataIdConfigJson;
    if (ToDataId(strDataIdConfigJson))
    {
        ApiError tempError;
        UpdateDataIdZk(strChannelId, strDataIdConfigJson, tempError);
    }

    return true;
}

bool ChannelIdUpdate::ToDataId(std::string &jsonStr)
{
    Json::Value dataIdJson;
    Json::Value dataIds;

    for (auto channel : m_ptrSpecification->m_channels)
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
        BLOG_INFO(m_requestId, "dataid(%d) config value(%s)", m_ptrCondition->m_channelID, jsonStr.c_str());
        return true;
    }
    return false;
}

void ChannelIdUpdate::MakeResponse(std::string &response, ApiError &error)
{
    Json::Value data;
    data.resize(0);
    Json::Value responseJson;
    ChannelIdCommonAPI::MakeResponse(responseJson, error, data);
    response.assign(responseJson.toStyledString());
}

ChannelIdUpdateAPIHandler::ChannelIdUpdateAPIHandler(const std::string &version)
{
}

ChannelIdUpdateAPIHandler::~ChannelIdUpdateAPIHandler()
{
}

int ChannelIdUpdateAPIHandler::OnPost(gse::net::http::HTTPMessagePtr message, std::string &response) noexcept
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
        return HTTP_STATUS_CODE_400;
    }

    return HTTP_STATUS_CODE_200;
}

bool ChannelIdUpdateAPIHandler::HandleRequest(gse::net::http::HTTPMessagePtr message, ApiError &error, std::string &response)
{
    Json::Value requestJson;
    Json::Reader reader(Json::Features::strictMode());
    std::string errorMsg;
    int errorCode = -1;
    std::shared_ptr<ChannelIdUpdate> ptrUpdateRequest = std::make_shared<ChannelIdUpdate>(m_requestId);

    auto _ = gse::tools::defer::finally([&]() {
        ptrUpdateRequest->MakeResponse(response, error);
        BLOG_INFO(m_requestId, "send response uri:%s, reponse(%s)", message->GetURI().c_str());
    });

    if (!reader.parse(message->GetBody().c_str(), requestJson, false))
    {
        BLOG_ERROR(m_requestId, "failed to parse request json,json invalid [%s]", message->GetBody().c_str());
        error = InputParamJsonInvalid;
        return false;
    }

    if (!ptrUpdateRequest->Parse(requestJson, error))
    {
        BLOG_ERROR(m_requestId, "failed to parse request, error:%s, request(%s)", GET_ERROR_MESSAGE(error).c_str(), message->GetBody().c_str());
        return false;
    }

    if (!ptrUpdateRequest->Valid(error))
    {
        BLOG_ERROR(m_requestId, "failed to check request param, error:%s, request(%s)", errorMsg.c_str(), message->GetBody().c_str());
        return false;
    }

    if (!ptrUpdateRequest->ToZK(error))
    {
        BLOG_ERROR(m_requestId, "failed to update config, error:%s, request(%s)", GET_ERROR_MESSAGE(error).c_str(), message->GetBody().c_str());
        return false;
    }

    return true;
}

} // namespace data
} // namespace gse
