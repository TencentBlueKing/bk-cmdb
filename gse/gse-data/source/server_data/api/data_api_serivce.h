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

#ifndef _GSE_DATA_API_SERVICE_H_
#define _GSE_DATA_API_SERVICE_H_
#include <memory>

#include "net/http/http_handler.hpp"
#include "net/http/http_server.hpp"

#include "net/migration/migration_client.h"
#include "net/migration/migration_server.h"

namespace gse {
namespace data {

using HttpURI = std::string;

class HTTPAPIService
{

public:
    HTTPAPIService(const std::string &host, int port, int threadNum);
    virtual ~HTTPAPIService();
    void Start();

    void Stop();
    void RegisterHttpURI(const std::string &uri, std::shared_ptr<gse::net::http::HTTPHandler> handler);

    int StartMigrationSerivce();

private:
    // need set
    std::string m_host;
    int m_port;
    int m_threadNum;

    // inner
    std::shared_ptr<gse::net::http::HTTPServer> m_httpServer;

    std::map<std::string, std::shared_ptr<gse::net::http::HTTPHandler>> m_uriHandlers;
    std::shared_ptr<std::thread> m_serviceThread;

    std::unique_ptr<gse::net::MigrationClient> m_migrationClient;
    std::unique_ptr<gse::net::MigrationServer> m_migrationServer;
    int m_listennerFd;
};

} // namespace data
} // namespace gse

#endif
