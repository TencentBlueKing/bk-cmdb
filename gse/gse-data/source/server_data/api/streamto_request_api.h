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

#ifndef _GSE_DATA_STREAMTOID_REQUEST_API_H_
#define _GSE_DATA_STREAMTOID_REQUEST_API_H_

#include "api/channelid_api.h"
#include "api/channelid_struct.h"
#include "api/channelid_def.h"
#include "error_code.h"
#include "net/http/http_handler.hpp"

namespace gse {
namespace data {

class StreamToClusterRequest
{
public:
    StreamToClusterRequest(const std::string requestId);

public:
    Json::Value ToJSONValue()
    {
    }

    inline std::string ToJSON()
    {
        Json::FastWriter writer;
        return writer.write(ToJSONValue());
    }

    bool Valid(std::string &errorMsg);
    bool ToZK(ApiError &error);
    bool GenerateStreamToClusterId(ApiError &error);

    bool ParseRequest(const Json::Value &requestJson, ApiError &error);

    void MakeResponse(std::string &response, ApiError &error);
    bool ToDataIdStorageConfigJson(std::string &jsString);

public:
    std::shared_ptr<StreamToMetadata> m_ptrMetadata;
    std::shared_ptr<Operation> m_ptrOperation;

    std::shared_ptr<StreamToCluster> m_ptrStreamToCluster;
    std::string m_originData;

    std::string m_streamToClusterId; // response channelid

    std::string m_requestId;
};

const std::string StreamToIdRequsetHandlerURI = "/dataroute/v1/add_streamto";
const std::string StreamToIdRequsetHandlerURIV2 = "/api/v2/data/config/add_streamto";

class ChannelStreamToRequestAPIHandler : public gse::net::http::HTTPHandler
{

public:
    ChannelStreamToRequestAPIHandler(const std::string &version = "v2");
    virtual ~ChannelStreamToRequestAPIHandler();

public:
private:
    int OnPost(gse::net::http::HTTPMessagePtr message, std::string &response) noexcept;
    bool HandleRequest(gse::net::http::HTTPMessagePtr message, ApiError &error, string &response);

private:
    std::string m_requestId;
    const std::string m_metricsApiName = {"add_streamto"};
    const std::string m_metricsReponseMillseconds = {"response_milliseconds"};
    const std::string m_metricsApiCount = {"api_requests_total"};
};

} // namespace data
} // namespace gse

#endif // STREAMTOID_REQUEST_API_H
