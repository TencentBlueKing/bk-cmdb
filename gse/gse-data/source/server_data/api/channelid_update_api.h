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

#ifndef _GSE_DATA_CHANNELID_UPDATE_API_H_
#define _GSE_DATA_CHANNELID_UPDATE_API_H_

#include <memory>
#include <string>

#include <json/json.h>

#include "channelid_api.h"
#include "channelid_def.h"
#include "discover/zkapi/zk_api.h"
#include "net/http/http_handler.hpp"

namespace gse {
namespace data {

class ChannelIdUpdate
{
public:
    ChannelIdUpdate(const std::string &requestId);
    bool Parse(const Json::Value &requestJson, ApiError &error);
    bool Valid(ApiError &error);
    bool ToZK(ApiError &error);
    void MakeResponse(std::string &response, ApiError &error);
    bool ToDataId(std::string &jsonStr);

private:
    std::shared_ptr<Condition> m_ptrCondition;
    std::shared_ptr<Operation> m_ptrOperation;
    std::shared_ptr<UpdateSpecification> m_ptrSpecification;
    std::string m_originData;
    std::string m_requestId;
};

const std::string ChannelIdUpdateHandlerURI = "/dataroute/v1/update_route";
const std::string ChannelIdUpdateHandlerURIV2 = "/api/v2/data/config/update_route";

class ChannelIdUpdateAPIHandler : public gse::net::http::HTTPHandler
{
public:
    ChannelIdUpdateAPIHandler(const std::string &version = "v2");
    virtual ~ChannelIdUpdateAPIHandler();

public:
private:
    void MakeResponse(int errorCode, const std::string &message, std::string &response) noexcept;
    int OnPost(gse::net::http::HTTPMessagePtr message, std::string &response) noexcept;
    bool HandleRequest(gse::net::http::HTTPMessagePtr message, ApiError &error, std::string &response);

private:
    std::string m_requestId;

    const std::string m_metricsApiName = {"update_route"};
    const std::string m_metricsReponseMillseconds = {"response_milliseconds"};
    const std::string m_metricsApiCount = {"api_requests_total"};
};

} // namespace data
} // namespace gse
#endif
