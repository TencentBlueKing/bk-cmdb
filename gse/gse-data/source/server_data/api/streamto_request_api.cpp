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

#include "streamto_request_api.h"

#include <json/json.h>
#include <json/reader.h>

#include "api_metrics.h"
#include "bbx/gse_errno.h"
#include "channelid_zk_api.h"
#include "common/logger.hpp"
#include "error_code.h"
#include "log/log.h"
#include "net/http/http_common.hpp"
#include "tools/finally.hpp"
#include "tools/maths.h"

namespace gse {
namespace data {

StreamToClusterRequest::StreamToClusterRequest(const string requestId)
    : m_requestId(requestId)
{
    m_ptrMetadata = std::make_shared<StreamToMetadata>();
    m_ptrOperation = std::make_shared<Operation>();
    m_ptrStreamToCluster = std::make_shared<StreamToCluster>();
}

bool StreamToClusterRequest::ParseRequest(const Json::Value &requestJson, ApiError &error)
{
    std::string errorMsg;
    if (!requestJson.isMember("metadata"))
    {
        error = MetadataNotSet;
        return false;
    }

    if (!requestJson.isMember("stream_to"))
    {
        error = StreamToNotSet;
        return false;
    }

    if (!m_ptrMetadata->ParseMetadata(requestJson["metadata"], error))
    {
        return false;
    }

    if (requestJson.isMember("operation"))
    {
        if (!m_ptrOperation->ParseOperation(requestJson["operation"], error))
        {
            return false;
        }
    }

    if (!m_ptrStreamToCluster->Parse(requestJson["stream_to"], error))
    {
        return false;
    }

    return true;
}

bool StreamToClusterRequest::GenerateStreamToClusterId(ApiError &error)
{
    int64_t streamToId = ChannelIdZkApiInst()->GenerateStreamToClusterId(m_ptrMetadata->m_requestPlatName, error);
    if (streamToId < 0)
    {
        return false;
    }

    m_streamToClusterId = gse::tools::strings::ToString(streamToId);
    return true;
}
bool StreamToClusterRequest::Valid(std::string &errorMsg)
{
    return true;
}

bool StreamToClusterRequest::ToDataIdStorageConfigJson(std::string &jsString)
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

bool StreamToClusterRequest::ToZK(ApiError &error)
{
    //指定id已存在,update
    std::string errorMsg;
    int streamToId = gse::tools::strings::StringToInt64(m_streamToClusterId);
    if (!ChannelIdZkApiInst()->ExistStreamToClusterId(m_streamToClusterId))
    {
        if (!ChannelIdZkApiInst()->CreateStreamToClusterNode(streamToId, m_ptrStreamToCluster, error))
        {
            return false;
        }

        if (!ChannelIdZkApiInst()->CreateStreamToClusterMetadataNode(streamToId, m_ptrMetadata, error))
        {
            return false;
        }
        // set stream_to_id
        m_ptrMetadata->m_streamToId = gse::tools::strings::StringToUint32(m_streamToClusterId);
        ChannelIdZkApiInst()->CreateStreamToIdQueryIndex(m_ptrMetadata, m_ptrStreamToCluster->m_reportMode, error);
    }
    else
    {
        if (!ChannelIdZkApiInst()->UpdateStreamToClusterConfig(streamToId, m_ptrStreamToCluster, error))
        {
            return false;
        }

        if (!ChannelIdZkApiInst()->UpdateStreamToClusterMeta(streamToId, m_ptrMetadata, error))
        {
            return false;
        }
    }

    std::string strDataIdStorage;
    if (ToDataIdStorageConfigJson(strDataIdStorage))
    {
        if (!CreateClusterIdToZk(m_streamToClusterId, strDataIdStorage, error))
        {
            BLOG_ERROR(m_requestId, "failed create clusterid config, error:%s", GET_ERROR_MESSAGE(error).c_str());
            return false;
        }
    }

    return true;
}

void StreamToClusterRequest::MakeResponse(std::string &response, ApiError &error)
{
    Json::Value data;
    data["stream_to_id"] = Json::Value(gse::tools::strings::StringToUint32(m_streamToClusterId));
    data["name"] = Json::Value(m_ptrStreamToCluster->m_name);

    Json::Value responseJson;
    ChannelIdCommonAPI::MakeResponse(responseJson, error, data);
    response.assign(responseJson.toStyledString());
}

ChannelStreamToRequestAPIHandler::ChannelStreamToRequestAPIHandler(const std::string &version)
{
}
ChannelStreamToRequestAPIHandler::~ChannelStreamToRequestAPIHandler()
{
}

bool ChannelStreamToRequestAPIHandler::HandleRequest(gse::net::http::HTTPMessagePtr message, ApiError &error, std::string &response)
{
    Json::Value requestJson;
    Json::Reader reader(Json::Features::strictMode());

    std::string errorMsg;
    std::shared_ptr<StreamToClusterRequest> ptrRequest = std::make_shared<StreamToClusterRequest>(m_requestId);

    auto _ = gse::tools::defer::finally([&]() {
        ptrRequest->MakeResponse(response, error);
        BLOG_INFO(m_requestId, "send response uri:%s, reponse(%s)",
                  message->GetURI().c_str(), response.c_str());
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

    if (!ptrRequest->Valid(errorMsg))
    {
        BLOG_ERROR(m_requestId, "failed to check request param, error:%s, request(%s)", GET_ERROR_MESSAGE(error).c_str(), message->GetBody().c_str());
        return false;
    }

    if (!ptrRequest->GenerateStreamToClusterId(error))
    {
        BLOG_ERROR(m_requestId, "failed to gennerate channelid, error:%s, request(%s)", GET_ERROR_MESSAGE(error).c_str(), message->GetBody().c_str());
        return false;
    }

    if (!ptrRequest->ToZK(error))
    {
        BLOG_ERROR(m_requestId, "failed to write config, error:%s, request(%s)", GET_ERROR_MESSAGE(error).c_str(), message->GetBody().c_str());
        return false;
    }

    return true;
}

int ChannelStreamToRequestAPIHandler::OnPost(gse::net::http::HTTPMessagePtr message, std::string &response) noexcept
{
    auto headerReqID = message->GetHeaderParameter(BK_APIGW_REQ_ID);
    m_requestId = headerReqID ? headerReqID->GetValue() : BK_API_UNKOWN_REO_ID;
    BLOG_INFO(m_requestId, "request uri:%s, message(%s)", message->GetURI().c_str(), message->GetBody().c_str());

    int errCode = GSE_SUCCESS;

    std::string errorMsg;
    uint64_t beginTimestamp = gse::tools::time::GetUTCMillsec();
    ApiError error;

    auto _ = gse::tools::defer::finally([&]() {
        int errcode = GET_ERROR_CODE(error);
        uint64_t endTimestamp = gse::tools::time::GetUTCMillsec();
        uint64_t costTime = (endTimestamp >= beginTimestamp) ? (endTimestamp >= beginTimestamp) : 0;
        ChannelIdApiMetricsInst()->CounterMetrics(m_metricsReponseMillseconds, m_metricsApiName, costTime, errcode);
        ChannelIdApiMetricsInst()->CounterMetrics(m_metricsApiCount, m_metricsApiName, 1, errcode);
    });

    if (!HandleRequest(message, error, response))
    {
        BLOG_ERROR(m_requestId, "failed to handle request, error:%s", GET_ERROR_MESSAGE(error).c_str());
        return HTTP_STATUS_CODE_400;
    }

    BLOG_INFO(m_requestId, "successfully handle request, uri:%s, message(%s)s", message->GetURI().c_str(), message->GetBody().c_str());
    return HTTP_STATUS_CODE_200;
}

} // namespace data
} // namespace gse
