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

#ifndef _GSE_DATA_CHANNELID_REQUEST_API_H_
#define _GSE_DATA_CHANNELID_REQUEST_API_H_

#include <vector>

#include <string>

#include "channelid_api.h"
#include "net/http/http_handler.hpp"

namespace gse {
namespace data {

class ChannelIdRequest
{

public:
    Json::Value ToJSONValue();
    bool ToDataId(string &jsonStr);
    inline std::string ToJSON()
    {
        Json::FastWriter writer;
        return writer.write(ToJSONValue());
    }

    bool Valid(ApiError &error);
    bool ToZK(ApiError &error);
    bool GenerateChannelId(ApiError &error);

    void MakeResponse(std::string &response, ApiError &error);
    bool ParseRequest(const Json::Value &requestJson, ApiError &error);

public:
    ChannelIdRequest(const std::string &requestId)
        : m_requestId(requestId)
    {
        m_ptrMetadata = std::make_shared<Metadata>();
        m_ptrOperation = std::make_shared<Operation>();
    }

    ~ChannelIdRequest()
    {
    }

public:
    std::shared_ptr<Metadata> m_ptrMetadata;
    std::shared_ptr<Operation> m_ptrOperation;

    std::vector<std::shared_ptr<Channel>> m_channels;
    std::vector<std::shared_ptr<StreamFilter>> m_streamFilters;

    std::string m_originData;

    std::string m_channelId; // response channelid
    std::string m_requestId;
};

const std::string ChannelIdRequsetHandlerURI = "/dataroute/v1/add_route";
const std::string ChannelIdRequsetHandlerURIV2 = "/api/v2/data/config/add_route";

class ChannelIdRequestAPIHandler : public gse::net::http::HTTPHandler
{

public:
    ChannelIdRequestAPIHandler(const std::string &version = "v2");
    virtual ~ChannelIdRequestAPIHandler();

public:
private:
    int OnPost(gse::net::http::HTTPMessagePtr message, std::string &response) noexcept;

    bool HandleRequest(gse::net::http::HTTPMessagePtr message, ApiError &error, std::string &response);

private:
    std::string m_requestId;

    const std::string m_metricsApiName = {"add_route"};
    const std::string m_metricsReponseMillseconds = {"response_milliseconds"};
    const std::string m_metricsApiCount = {"api_requests_total"};
};

} // namespace data
} // namespace gse

#endif // _GSE_DATA_CHANNELID_REQUEST_API_H_
