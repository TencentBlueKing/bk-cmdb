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

#include "channelid_delete_api.h"

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
ChannelIdDelete::ChannelIdDelete(const std::string &requestId)
    : m_requestId(requestId)
{
    m_ptrOperation = std::make_shared<Operation>();
    m_ptrCondition = std::make_shared<Condition>();
    m_ptrDeleteSpecification = std::make_shared<DeleteSpecification>();
}

void ChannelIdDelete::MakeResponse(std::string &response, ApiError &error)
{
    Json::Value data;
    data.resize(0);
    Json::Value responseJson;
    ChannelIdCommonAPI::MakeResponse(responseJson, error, data);
    response.assign(responseJson.toStyledString());
}

bool ChannelIdDelete::DeleteFromZK(std::string &result, ApiError &error)
{
    std::string strChannelId = gse::tools::strings::ToString(m_ptrCondition->m_channelID);
    auto ptrChannelIdConfg = ChannelIdZkApiInst()->ReadChannelIdConfig(strChannelId, m_ptrCondition->m_requestPlatName, error);
    if (ptrChannelIdConfg == nullptr)
    {
        BLOG_ERROR(m_requestId, "failed to read channelid(%s) config, error:%s", strChannelId.c_str(), GET_ERROR_MESSAGE(error).c_str());
        return false;
    }

    if (m_ptrCondition->m_requestPlatName != ptrChannelIdConfg->m_ptrMetadata->m_requestPlatName)
    {
        BLOG_ERROR(m_requestId, "channel id(%s) does not belong to the platform(%s), real plat name(%s)",
                   strChannelId.c_str(), m_ptrCondition->m_requestPlatName.c_str(),
                   ptrChannelIdConfg->m_ptrMetadata->m_requestPlatName.c_str());
        error = NoPermissionDelete;
        return false;
    }

    std::string method = m_ptrOperation->m_method;

    if (method == ChannelIDOperationDeleteMethodAll)
    {
        if (!ChannelIdZkApiInst()->DeleteChannelId(m_ptrCondition->m_channelID, error))
        {
            BLOG_ERROR(m_requestId, "failed to delete channelid, error:%s", GET_ERROR_MESSAGE(error).c_str());
        }

        std::string platName = m_ptrCondition->m_requestPlatName;
        // delete channelid index
        if (!platName.empty())
        {
            ChannelIdZkApiInst()->DeleteChannelIdPlatNameIndex(platName, strChannelId);
        }

        if (m_ptrCondition->m_label.m_valid && m_ptrCondition->m_label.m_bkBizId > 0)
        {
            ChannelIdZkApiInst()->DeleteChannelIdBizIdIndex(gse::tools::strings::ToString(m_ptrCondition->m_label.m_bkBizId), strChannelId);
        }

        if (m_ptrCondition->m_label.m_valid && !m_ptrCondition->m_label.m_odm.empty())
        {
            ChannelIdZkApiInst()->DeleteChannelIdOdmNameIndex(m_ptrCondition->m_label.m_odm, strChannelId);
        }

        for (auto channel : ptrChannelIdConfg->m_channels)
        {
            ChannelIdZkApiInst()->DeleteChannelIdStreamToIndex(gse::tools::strings::ToString(channel->m_streamTo.m_streamToId), strChannelId);
        }

        if (ChannelIdCommonAPI::IsTglogPlat(platName))
        {
            ChannelIdZkApiInst()->DeleteTglogChannelId(m_ptrCondition->m_label.m_bkBizId, m_ptrCondition->m_label.m_odm);
        }
        else
        {
            ApiError tempError;
            DeleteDataIdZk(gse::tools::strings::ToString(m_ptrCondition->m_channelID), tempError);
        }
    }
    else if (method == ChannelIDOperationDeleteMethodSpecification)
    {
        ApiError tempError;
        if (!ChannelIdZkApiInst()->DeleteBySpecification(m_ptrCondition->m_channelID, m_ptrDeleteSpecification, tempError))
        {
            BLOG_ERROR(m_requestId, "failed to delete specification, tempError:%s", GET_ERROR_MESSAGE(tempError).c_str());
            return false;
        }
    }
    else
    {
        error = UnsupportDeleteMethod;
        return false;
    }

    return true;
}

bool ChannelIdDelete::ParseRequest(const Json::Value &requestJson, ApiError &error)
{
    if (!requestJson.isMember("condition"))
    {
        error = CondtionNotSet;
        BLOG_ERROR(m_requestId, "failed to parse channeld delete request, condition not set, json(%s)", requestJson.toStyledString().c_str());

        return false;
    }

    if (!requestJson.isMember("operation"))
    {
        error = OperationNotSet;
        return false;
    }

    if (requestJson.isMember("operation"))
    {
        m_ptrOperation->ParseOperation(requestJson["operation"], error);
    }

    if (!m_ptrCondition->Parse(requestJson["condition"], error))
    {
        BLOG_ERROR(m_requestId, "failed to parse channeld delete request, condition error, json(%s)", requestJson.toStyledString().c_str());
        return false;
    }

    if (m_ptrOperation->m_method == "specification")
    {
        if (!requestJson.isMember("specification"))
        {
            error = SpecificationNotSet;
            BLOG_ERROR(m_requestId, "failed to parse channeld delete request, specification not set, json(%s)", requestJson.toStyledString().c_str());
            return false;
        }

        if (!m_ptrDeleteSpecification->Parse(requestJson["specification"], error))
        {

            BLOG_ERROR(m_requestId, "failed to parse channeld delete request, specification error, json(%s), error:%s",
                       requestJson.toStyledString().c_str(), GET_ERROR_MESSAGE(error).c_str());
            return false;
        }
    }

    return true;
}

ChannelIdDeleteAPIHandler::ChannelIdDeleteAPIHandler(const std::string &version)
{
}

ChannelIdDeleteAPIHandler::~ChannelIdDeleteAPIHandler()
{
}

bool ChannelIdDeleteAPIHandler::HandleRequest(gse::net::http::HTTPMessagePtr message, ApiError &error, std::string &response)
{
    Json::Value requestJson;
    Json::Reader reader(Json::Features::strictMode());
    std::shared_ptr<ChannelIdDelete> ptrRequest = std::make_shared<ChannelIdDelete>(m_requestId);

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

    std::string result;
    if (!ptrRequest->DeleteFromZK(result, error))
    {
        BLOG_ERROR(m_requestId, "failed to delete config, error:%s, request(%s)", GET_ERROR_MESSAGE(error).c_str(), message->GetBody().c_str());
        return false;
    }

    return true;
}

int ChannelIdDeleteAPIHandler::OnPost(gse::net::http::HTTPMessagePtr message, std::string &response) noexcept
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
