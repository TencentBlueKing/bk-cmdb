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

#include "streamto_delete_api.h"

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

StreamToIdDelete::StreamToIdDelete(const std::string &requestId)
    : m_requestId(requestId)
{
    m_ptrCondition = std::make_shared<StreamToCondition>();
    m_ptrOperation = std::make_shared<Operation>();
}

bool StreamToIdDelete::ParseRequest(const Json::Value &requestJson, ApiError &error)
{
    if (!requestJson.isMember("condition"))
    {
        error = CondtionNotSet;
        BLOG_ERROR(m_requestId, "failed to parse streamto delete request, condition not set, json(%s)", requestJson.toStyledString().c_str());

        return false;
    }

    if (requestJson.isMember("operation"))
    {
        m_ptrOperation->ParseOperation(requestJson["operation"], error);
    }

    if (!m_ptrCondition->Parse(requestJson["condition"], error))
    {
        BLOG_ERROR(m_requestId, "failed to parse stream delete request, condition error, json(%s)", requestJson.toStyledString().c_str());
        return false;
    }

    return true;
}

void StreamToIdDelete::MakeResponse(std::string &response, ApiError &error)
{
    Json::Value data;
    data.resize(0);
    Json::Value responseJson;
    ChannelIdCommonAPI::MakeResponse(responseJson, error, data);
    response.assign(responseJson.toStyledString());
}

bool StreamToIdDelete::DeleteFromZK(ApiError &error)
{
    std::string strStreamToId = gse::tools::strings::ToString(m_ptrCondition->m_streamToId);
    std::string platName = m_ptrCondition->m_requestPlatName;

    if (m_ptrCondition->m_streamToId != 0)
    {
        if (!ChannelIdZkApiInst()->ExistStreamToClusterId(strStreamToId))
        {
            error = StreamToIdNotExist;
            return false;
        }

        auto streamToIdConfig = ChannelIdZkApiInst()->ReadStreamToIdConfig(strStreamToId, error);
        if (streamToIdConfig == nullptr)
        {
            BLOG_ERROR(m_requestId, "failed to read stream(%s) config from zk, error:%s", strStreamToId.c_str(), GET_ERROR_MESSAGE(error).c_str());
            return false;
        }

        if (streamToIdConfig->m_ptrMetadata->m_requestPlatName != m_ptrCondition->m_requestPlatName)
        {
            error = NoPermissionDelete;
            BLOG_ERROR(m_requestId, "failed to delete stream_to_id(%s), the id don't belong to the plat(%s)",
                       strStreamToId.c_str(),
                       streamToIdConfig->m_ptrMetadata->m_requestPlatName.c_str());
            return false;
        }

        if (!ChannelIdZkApiInst()->DeleteStreamToID(strStreamToId, error))
        {
            BLOG_ERROR(m_requestId, "failed to delete stream_to_id(%s)", strStreamToId.c_str(), GET_ERROR_MESSAGE(error).c_str());
            return false;
        }

        ChannelIdZkApiInst()->DeleteStreamToBkBizIdIndex(gse::tools::strings::ToString(streamToIdConfig->m_ptrMetadata->m_label.m_bkBizId), strStreamToId, error);
        ChannelIdZkApiInst()->DeleteStreamToPlatNameIndex(streamToIdConfig->m_ptrMetadata->m_requestPlatName, strStreamToId, error);
        ChannelIdZkApiInst()->DeleteClusterTypeIndex(streamToIdConfig->m_ptrStreamToCluster->m_reportMode, strStreamToId, error);

        ChannelIdZkApiInst()->DeleteOriginStreamToID(strStreamToId, error);
        // ignore delete dataid error
        ApiError tmpError;
        DeleteClusterIdFromZk(strStreamToId, tmpError);
        return true;

        //        if (canWriteDataId(platname))
        //        {
        //            deleteStorageConfig(stream_to_id, errorMsg);
        //        }
    }
    else
    {
        error = StreamToIdNotSet;
        return false;
    }

    return true;
}

StreamToDeleteAPIHandler::StreamToDeleteAPIHandler(const std::string &version)
{
}
StreamToDeleteAPIHandler::~StreamToDeleteAPIHandler()
{
}

bool StreamToDeleteAPIHandler::HandleRequest(gse::net::http::HTTPMessagePtr message, ApiError &error, std::string &response)
{
    Json::Value requestJson;
    Json::Reader reader(Json::Features::strictMode());

    std::string errorMsg;
    std::shared_ptr<StreamToIdDelete> ptrRequest = std::make_shared<StreamToIdDelete>(m_requestId);

    auto _ = gse::tools::defer::finally([&]() {
        ptrRequest->MakeResponse(response, error);
        BLOG_INFO(m_requestId, "send response uri:%s, reponse(%s)", message->GetURI().c_str(), response.c_str());
    });

    if (!reader.parse(message->GetBody().c_str(), requestJson, false))
    {
        BLOG_ERROR(m_requestId, "failed to parse request json,json invalid [%s]", message->GetBody().c_str());
        error = InputParamJsonInvalid;
        return false;
    }

    if (!ptrRequest->ParseRequest(requestJson, error))
    {
        BLOG_ERROR(m_requestId, "failed to parse request json, error:%s, request(%s)", GET_ERROR_MESSAGE(error).c_str(), message->GetBody().c_str());
        return false;
    }

    //    if (ptrRequest->Valid(error))
    //    {
    //        return false;
    //    }
    if (!ptrRequest->DeleteFromZK(error))
    {
        BLOG_ERROR(m_requestId, "failed to delete stream_to from zk, stream_to_id(%d), error:%s, request(%s)",
                   ptrRequest->m_ptrCondition->m_streamToId, GET_ERROR_MESSAGE(error).c_str(), message->GetBody().c_str());
        return false;
    }

    return true;
}

int StreamToDeleteAPIHandler::OnPost(gse::net::http::HTTPMessagePtr message, std::string &response) noexcept
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
