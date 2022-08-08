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

#ifndef _GSE_DATA_STREAMTO_QUERY_API_H_
#define _GSE_DATA_STREAMTO_QUERY_API_H_

#include "api/channelid_api.h"
#include "api/channelid_struct.h"
#include "api/channelid_def.h"
#include "error_code.h"
#include "net/http/http_handler.hpp"

namespace gse {
namespace data {

class StreamTolIdQuery
{
public:
    Json::Value ToJSONValue();
    inline std::string ToJSON()
    {
        Json::FastWriter writer;
        return writer.write(ToJSONValue());
    }

    bool Valid(ApiError &error);
    bool ReadFromZK(ApiError &error);
    void MakeResponse(std::string &response, ApiError &error);
    bool ParseRequest(const Json::Value &requestJson, ApiError &error);
    std::string m_requestId;

public:
    StreamTolIdQuery(const std::string &requestId);

    ~StreamTolIdQuery()
    {
    }

public:
    std::shared_ptr<Operation> m_ptrOperation;
    std::shared_ptr<StreamIdCondition> m_ptrCondition;
    Json::Value m_result;
};

const std::string StreamToIdQueryHandlerURI = "/dataroute/v1/query_streamto";
const std::string StreamToIdQueryHandlerURIV2 = "/api/v2/data/config/query_streamto";

class StreamToIdQueryAPIHandler : public gse::net::http::HTTPHandler
{

public:
    StreamToIdQueryAPIHandler(const std::string &version = "v2");
    virtual ~StreamToIdQueryAPIHandler();

public:
private:
    int OnPost(gse::net::http::HTTPMessagePtr message, std::string &response) noexcept;

    bool HandleRequest(gse::net::http::HTTPMessagePtr message, ApiError &error, std::string &response);

private:
    std::string m_requestId;

    const std::string m_metricsApiName = {"query_streamto"};
    const std::string m_metricsReponseMillseconds = {"response_milliseconds"};
    const std::string m_metricsApiCount = {"api_requests_total"};
};

} // namespace data
} // namespace gse

#endif // STREAMTO_QUERY_API_H
