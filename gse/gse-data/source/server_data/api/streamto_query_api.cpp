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

#include "streamto_query_api.h"

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

StreamTolIdQuery::StreamTolIdQuery(const std::string &requestId)
    : m_requestId(requestId)
{
    m_ptrOperation = std::make_shared<Operation>();
    m_ptrCondition = std::make_shared<StreamIdCondition>();
}

bool StreamTolIdQuery::ReadFromZK(ApiError &error)
{
    // query by channel id
    if (m_ptrCondition->m_streamToId != 0)
    {
        std::string strStreamToId = gse::tools::strings::ToString(m_ptrCondition->m_streamToId);
        bool bExist = ChannelIdZkApiInst()->ExistStreamToClusterId(strStreamToId);
        if (!bExist)
        {
            error = StreamToIdNotExist;
            return false;
        }

        auto ptrStreamToConfig = ChannelIdZkApiInst()->ReadStreamToIdConfig(strStreamToId, error);
        if (ptrStreamToConfig == nullptr)
        {
            return false;
        }

        Json::Value data = ptrStreamToConfig->m_ptrStreamToCluster->ToJSONValue();
        data["stream_to_id"] = Json::Value(m_ptrCondition->m_streamToId);
        m_result.append(data);
        BLOG_INFO(m_requestId, "read stream_to_id[%d] config from zk, result(%s)", m_ptrCondition->m_streamToId, m_result.toStyledString().c_str());
        return true;
    }

    error = StreamToIdNotExist;
    return false;
}

bool StreamTolIdQuery::ParseRequest(const Json::Value &requestJson, ApiError &error)
{
    if (!requestJson.isMember("condition"))
    {
        error = CondtionNotSet;
        return false;
    }

    if (!m_ptrCondition->Parse(requestJson["condition"], error))
    {
        return false;
    }

    if (requestJson.isMember("operation"))
    {
        m_ptrOperation->ParseOperation(requestJson["operation"], error);
    }
    return true;
}

void StreamTolIdQuery::MakeResponse(std::string &response, ApiError &error)
{
    Json::Value data;
    data.resize(0);
    Json::Value responseJson;
    ChannelIdCommonAPI::MakeResponse(responseJson, error, m_result);
    response.assign(responseJson.toStyledString());
}

StreamToIdQueryAPIHandler::StreamToIdQueryAPIHandler(const std::string &version)
{
}

StreamToIdQueryAPIHandler::~StreamToIdQueryAPIHandler()
{
}

bool StreamToIdQueryAPIHandler::HandleRequest(gse::net::http::HTTPMessagePtr message, ApiError &error, std::string &response)
{
    Json::Value requestJson;
    Json::Reader reader(Json::Features::strictMode());
    std::shared_ptr<StreamTolIdQuery> ptrRequest = std::make_shared<StreamTolIdQuery>(m_requestId);

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
        BLOG_ERROR(m_requestId, "failed to parse request json, error:%s, request(%s)", GET_ERROR_MESSAGE(error).c_str(), message->GetBody().c_str());
        return false;
    }

    if (!ptrRequest->ReadFromZK(error))
    {
        BLOG_ERROR(m_requestId, "failed to read config, error:%s, request(%s)", GET_ERROR_MESSAGE(error).c_str(), message->GetBody().c_str());
        return false;
    }

    return true;
}

int StreamToIdQueryAPIHandler::OnPost(gse::net::http::HTTPMessagePtr message, std::string &response) noexcept
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
