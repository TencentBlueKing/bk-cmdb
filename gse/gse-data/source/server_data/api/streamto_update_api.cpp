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

#include "streamto_update_api.h"

#include <json/json.h>

#include "api_metrics.h"
#include "channelid_zk_api.h"
#include "common/logger.hpp"
#include "error_code.h"
#include "tools/finally.hpp"
#include "tools/maths.h"

namespace gse {
namespace data {

StreamToUpdate::StreamToUpdate(const std::string &requestId)
    : m_requestId(requestId)
{
    m_ptrOperation = std::make_shared<Operation>();
    m_ptrCondition = std::make_shared<StreamToCondition>();
    m_ptrStreamToCluster = std::make_shared<StreamToCluster>();
}

bool StreamToUpdate::ParseRequest(const Json::Value &requestJson, ApiError &error)
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
    if (requestJson.isMember("operation"))
    {
        m_ptrOperation->ParseOperation(requestJson["operation"], error);
    }

    if (!m_ptrCondition->Parse(requestJson["condition"], error))
    {
        BLOG_ERROR(m_requestId, "failed to parse conditon, error:%s", GET_ERROR_MESSAGE(error).c_str());
        return false;
    }

    if (!requestJson["specification"].isMember("stream_to"))
    {
        error = SpecificationNotSetStreamTo;
        BLOG_ERROR(m_requestId, "specification not set stream_to");
        return false;
    }
    if (!m_ptrStreamToCluster->Parse(requestJson["specification"]["stream_to"], error))
    {
        BLOG_ERROR(m_requestId, "failed to parse stream_to, error:%s", GET_ERROR_MESSAGE(error).c_str());
        return false;
    }

    return true;
}

void StreamToUpdate::MakeResponse(std::string &response, ApiError &error)
{
    Json::Value data;
    data.resize(0);
    Json::Value responseJson;
    ChannelIdCommonAPI::MakeResponse(responseJson, error, data);
    response.assign(responseJson.toStyledString());
}

bool StreamToUpdate::Valid(ApiError &error)
{
    if (!m_ptrCondition->Valid())
    {
        return false;
    }

    if (!m_ptrStreamToCluster->Valid(error))
    {
        return false;
    }

    return true;
}

bool StreamToUpdate::ToDataIdStorageConfigJson(std::string &jsString)
{
    Json::Value data;
    Json::Value jsValue;
    std::string reportMode = m_ptrStreamToCluster->m_reportMode;
    if (reportMode == ChannelIDOperationRequestMethodKafka)
    {
        for (auto addr : m_ptrStreamToCluster->m_kafkaCluster->m_addresses)
        {
            jsValue["host"] = Json::Value(addr.m_ip);
            jsValue["port"] = Json::Value(addr.m_port);
            jsValue["cluster_index"] = Json::Value(gse::tools::strings::StringToUint32(m_streamToClusterId));
            jsValue["type"] = Json::Value(KAFKA_COMMON);
            data.append(jsValue);
        }
    }
    else if (reportMode == ChannelIDOperationRequestMethodRedis)
    {
        for (auto addr : m_ptrStreamToCluster->m_redisCluster->m_addresses)
        {
            jsValue["host"] = Json::Value(addr.m_ip);
            jsValue["port"] = Json::Value(addr.m_port);
            jsValue["cluster_index"] = Json::Value(gse::tools::strings::StringToUint32(m_streamToClusterId));
            if (m_ptrStreamToCluster->m_redisCluster->m_mode == kRedisSentinel)
            {
                jsValue["type"] = Json::Value(REDIS_SENTINEL_PUB);
            }
            else
            {
                jsValue["type"] = Json::Value(REDIS_PUB);
            }

            data.append(jsValue);
        }
    }
    if (data.size() > 0)
    {
        Json::FastWriter json_fwriter;
        jsString.assign(json_fwriter.write(data));
        return true;
    }

    return false;
}

bool StreamToUpdate::ToZK(ApiError &error)
{
    if (!ChannelIdZkApiInst()->UpdateStreamToClusterConfig(m_ptrCondition->m_streamToId, m_ptrStreamToCluster, error))
    {
        return false;
    }

    BLOG_INFO(m_requestId, "successfully update stream to config(%s)", m_ptrStreamToCluster->ToJSON().c_str());

    std::string strDataIdStorage;
    if (ToDataIdStorageConfigJson(strDataIdStorage))
    {
        ApiError tempError;
        if (!UpdateClusterIdToZk(gse::tools::strings::ToString(m_ptrCondition->m_streamToId), strDataIdStorage, tempError))
        {
            BLOG_ERROR(m_requestId, "failed to update clusterid config, error:%s", GET_ERROR_MESSAGE(error).c_str());
            return false;
        }
    }
    return true;
}

ChannelStreamToUpdateAPIHandler::ChannelStreamToUpdateAPIHandler(const std::string &version)
{
}
ChannelStreamToUpdateAPIHandler::~ChannelStreamToUpdateAPIHandler()
{
}

bool ChannelStreamToUpdateAPIHandler::HandleRequest(gse::net::http::HTTPMessagePtr message, ApiError &error, std::string &response)
{
    Json::Value requestJson;
    Json::Reader reader(Json::Features::strictMode());

    std::string errorMsg;
    std::shared_ptr<StreamToUpdate> ptrRequest = std::make_shared<StreamToUpdate>(m_requestId);

    auto _ = gse::tools::defer::finally([&]() {
        ptrRequest->MakeResponse(response, error);
        BLOG_INFO(m_requestId, "send response uri:%s, reponse(%s)", message->GetURI().c_str(), response.c_str());
    });

    if (!reader.parse(message->GetBody().c_str(), requestJson, false))
    {
        error = InputParamJsonInvalid;
        BLOG_ERROR(m_requestId, "failed to parse request json,json invalid, request(%s)", message->GetBody().c_str());
        return false;
    }

    if (!ptrRequest->ParseRequest(requestJson, error))
    {
        BLOG_ERROR(m_requestId, "failed to parse request json, error:%s, request(%s)", GET_ERROR_MESSAGE(error).c_str(), message->GetBody().c_str());
        return false;
    }

    if (!ptrRequest->Valid(error))
    {
        BLOG_ERROR(m_requestId, "input param invalid, error:%s, request(%s)", GET_ERROR_MESSAGE(error).c_str(), message->GetBody().c_str());
        return false;
    }

    if (!ptrRequest->ToZK(error))
    {
        BLOG_ERROR(m_requestId, "failed to write to zk, error:%s, request(%s)", GET_ERROR_MESSAGE(error).c_str(), message->GetBody().c_str());
        return false;
    }

    return true;
}

int ChannelStreamToUpdateAPIHandler::OnPost(gse::net::http::HTTPMessagePtr message, std::string &response) noexcept
{
    auto headerReqID = message->GetHeaderParameter(BK_APIGW_REQ_ID);
    m_requestId = headerReqID ? headerReqID->GetValue() : BK_API_UNKOWN_REO_ID;
    BLOG_INFO(m_requestId, "request uri:%s, message(%s)", message->GetURI().c_str(), message->GetBody().c_str());

    int errCode = GSE_SUCCESS;

    std::string errorMsg;
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
        BLOG_ERROR(m_requestId, "failed to handle request, error:%s", GET_ERROR_MESSAGE(error).c_str());
        return HTTP_STATUS_CODE_400;
    }

    BLOG_INFO(m_requestId, "successfully handle request, uri:%s, message(%s)", message->GetURI().c_str(), message->GetBody().c_str());
    return HTTP_STATUS_CODE_200;
}

} // namespace data
} // namespace gse
