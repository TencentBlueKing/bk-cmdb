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

#ifndef _GSE_DATA_STREAMTO_UPDATE_API_H_
#define _GSE_DATA_STREAMTO_UPDATE_API_H_

#include "api/channelid_api.h"
#include "api/channelid_struct.h"
#include "api/channelid_def.h"
#include "error_code.h"
#include "net/http/http_handler.hpp"

namespace gse {
namespace data {

class StreamToUpdate
{
public:
    StreamToUpdate(const std::string &requestId);

public:
    Json::Value ToJSONValue();

    inline std::string ToJSON()
    {
        Json::FastWriter writer;
        return writer.write(ToJSONValue());
    }

    bool Valid(ApiError &error);
    bool ToZK(ApiError &error);

    bool ParseRequest(const Json::Value &requestJson, ApiError &error);

    void MakeResponse(std::string &response, ApiError &error);
    bool ToDataIdStorageConfigJson(std::string &jsString);

public:
    std::shared_ptr<StreamToCondition> m_ptrCondition;
    std::shared_ptr<Operation> m_ptrOperation;

    std::shared_ptr<StreamToCluster> m_ptrStreamToCluster;
    std::string m_originData;

    std::string m_streamToClusterId; // response channelid
    std::string m_requestId;

    // std::vector<std::shared_ptr<StreamTo>> m_streamTo;
};

const std::string StreamToIdUpdateHandlerURI = "/dataroute/v1/update_streamto";
const std::string StreamToIdUpdateHandlerURIV2 = "/api/v2/data/config/update_streamto";

class ChannelStreamToUpdateAPIHandler : public gse::net::http::HTTPHandler
{

public:
    ChannelStreamToUpdateAPIHandler(const std::string &version = "v2");
    virtual ~ChannelStreamToUpdateAPIHandler();

public:
private:
    int OnPost(gse::net::http::HTTPMessagePtr message, std::string &response) noexcept;
    bool HandleRequest(gse::net::http::HTTPMessagePtr message, ApiError &error, std::string &response);

private:
    std::string m_requestId;
    const std::string m_metricsApiName = {"update_streamto"};
    const std::string m_metricsReponseMillseconds = {"response_milliseconds"};
    const std::string m_metricsApiCount = {"api_requests_total"};
};

} // namespace data
} // namespace gse
#endif
