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

#ifndef _GSE_DATA_HTTP_RECEIVER_H_
#define _GSE_DATA_HTTP_RECEIVER_H_

#include "net/http/http_common.hpp"
#include "net/http/http_message.hpp"
#include "net/migration/migration_client.h"
#include "net/migration/migration_server.h"
#include "receiver.h"

namespace gse {
namespace data {

#ifndef SEND_DATA_BY_DATAID
#define SEND_DATA_BY_DATAID "/gse/v1/senddata/"
#endif

class HttpReceiverHandler : public gse::net::http::HTTPHandler
{

public:
    HttpReceiverHandler();
    HttpReceiverHandler(RecvDataCallBack fnRecvData, void *pCaller);
    virtual ~HttpReceiverHandler();
    void SetServerIp(const std ::string &ip);
    void SetServerPort(uint16_t port);

private:
    void makeResponse(int errorCode, const std::string &message, std::string &response) noexcept;
    int OnPost(gse::net::http::HTTPMessagePtr message, std::string &response) noexcept;
    uint32_t getChannelIDFromURI(const char *uri) noexcept;

private:
    RecvDataCallBack m_callback;
    void *m_caller;

    std::string m_servIp;
    uint16_t m_servPort;
};

class HttpReceiver : public Receiver,
                     public gse::net::http::HTTPHandler
{
public:
    HttpReceiver();
    virtual ~HttpReceiver();

public:
    int Start() override;
    int Stop() override;
    void Join() override;

private:
    int StartMigrationSerivce();

private:
    gse::net::http::HTTPServer *m_httpserver;
    std::thread m_listenThread;

    std::unique_ptr<gse::net::MigrationClient> m_migrationClient;
    std::unique_ptr<gse::net::MigrationServer> m_migrationServer;
    int m_listennerFd;
};

} // namespace data
} // namespace gse
#endif
