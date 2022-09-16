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

#include "data_api_serivce.h"

#include <chrono> // std::chrono::milliseconds
#include <future> // std::async, std::future

#include "dataserver.h"
#include "log/log.h"
#include "ops/op_healthz.h"

namespace gse {
namespace data {

HTTPAPIService::HTTPAPIService(const std::string &host, int port, int threadNum)
    : m_host(host), m_port(port), m_threadNum(threadNum),
      m_listennerFd(NULL_SOCK)

{
}

HTTPAPIService::~HTTPAPIService()
{
}

int HTTPAPIService::StartMigrationSerivce()
{
    m_migrationClient = std::unique_ptr<gse::net::MigrationClient>(new gse::net::MigrationClient(m_port, 20));
    int domainsocket_listenner_fd = -1;
    if (m_migrationClient->ConnectDomainSocket() != GSE_SUCCESS)
    {
        LOG_WARN("failed to connect channelid api domain socket");
    }
    else
    {
        m_listennerFd = m_migrationClient->MigrateListennerFd();
        domainsocket_listenner_fd = m_migrationClient->MigrateDomainSocketListenner();
    }

    auto pFuncGetListernnerFd = std::bind(&gse::net::http::HTTPServer::GetEvhtpListennerFd, m_httpServer);
    auto pFuncStopListenner = std::bind(&gse::net::http::HTTPServer::StoppingListenner, m_httpServer);

    m_migrationServer = std::unique_ptr<gse::net::MigrationServer>(new gse::net::MigrationServer(m_port));
    m_migrationServer->SetGetListennerFdCallback(pFuncGetListernnerFd);
    // m_migrationServer->SetStopListennerCallback(pFuncStopListenner);
    //  m_migrationServer->SetFinishedCallback(&DataServer::GracefullyQuit);

    return m_migrationServer->StartMigrationService(domainsocket_listenner_fd);
}

void HTTPAPIService::RegisterHttpURI(const std::string &uri, std::shared_ptr<gse::net::http::HTTPHandler> handler)
{
    auto it = m_uriHandlers.find(uri);
    if (it == m_uriHandlers.end())
    {
        m_uriHandlers.insert(std::make_pair(uri, handler));
    }
    else
    {
    }
}

void HTTPAPIService::Start()
{
    m_httpServer = std::make_shared<gse::net::http::HTTPServer>(m_host, m_port, m_threadNum);
#ifdef _GSE_SSL_
    // m_httpServer->EnableSSL("../cert/gseca.crt", "../cert/gse_agent.crt", "../cert/gse_agent.key", "../cert/cert_encrypt.key", _GSE_SSL_);
#endif

    for (auto handler : m_uriHandlers)
    {
        bool ret = m_httpServer->RegisterHandler(handler.first, handler.second);
        if (!ret)
        {
            LOG_ERROR("failed to register http handler, uri:%s", handler.first.c_str());
        }
    }

    if (DataServer::GetUpgradeFlag())
    {
        StartMigrationSerivce();
        m_httpServer->SetListenerFd(m_listennerFd);
    }

    m_serviceThread = std::make_shared<std::thread>(&gse::net::http::HTTPServer::Start, m_httpServer);

    //    int fd = m_httpServer->GetEvhtpListennerFd();
    //    if (fd == -1)
    //    {
    //        LOG_ERROR("failed to start config api service");
    //        std::string hostAddr = m_host + ":" + gse::tools::strings::ToString(m_port);
    //        OpHealthZ::AddInitHealthInfo("", "start config api server on" + hostAddr + "error", -1);
    //        return;
    //    }

    std::string hostAddr = m_host + ":" + gse::tools::strings::ToString(m_port);
    OpHealthZ::AddInitHealthInfo("channelid_api", "ok", 0);

    LOG_DEBUG("start config api service succesfully");
}

void HTTPAPIService::Stop()
{
    m_httpServer->Stop();
    if (m_serviceThread->joinable())
    {
        m_serviceThread->join();
    }
}

} // namespace data
} // namespace gse
