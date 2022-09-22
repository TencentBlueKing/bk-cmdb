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

#include "httpreceiver.h"

#include <memory>

#include "bbx/gse_errno.h"
#include "dataserver.h"
#include "net/http/http_server.hpp"
#include "tools/macros.h"

namespace gse {
namespace data {

HttpReceiverHandler::HttpReceiverHandler()
    : m_callback(nullptr), m_caller(nullptr)
{
}

HttpReceiverHandler::HttpReceiverHandler(RecvDataCallBack fnRecvData, void *pCaller)
    : m_callback(fnRecvData), m_caller(pCaller)
{
}

HttpReceiverHandler::~HttpReceiverHandler()
{
}

int HttpReceiverHandler::OnPost(gse::net::http::HTTPMessagePtr message, std::string &response) noexcept
{
    std::string content;
    content = message->GetBody();
    OPMetric::HttpMsgInc(content.size());
    if (content.empty())
    {
        LOG_ERROR("it is to read nothing from the http request");
        makeResponse(GSE_ERROR, std::string("it is to read nothing from the http request"), response);
        return EVHTP_RES_SERVERR;
    }

    std::string uri = message->GetURI();
    uint32_t channel_id = getChannelIDFromURI(uri.c_str());

    DataCell *pDataCell = new DataCell();
    std::string clientIP = message->GetClientIP();
    pDataCell->SetSourceIp(clientIP);
    pDataCell->SetServerIP(DataServer::GetAdvertiseIP());
    pDataCell->SetServerPort(m_servPort);
    pDataCell->SetChannelProtocol("HTTP");
    pDataCell->SetChannelID(channel_id);

    if (GSE_SUCCESS != pDataCell->CopyData(content.c_str(), content.size()))
    {
        LOG_ERROR("it is failed to copy the body data (%s) int a datacell", SAFE_CSTR(content.c_str()));
        makeResponse(GSE_ERROR, std::string("it is failed to copy the body data into a datacell"), response);
        delete pDataCell;
        return EVHTP_RES_SERVERR;
    }

    if (NULL == m_callback)
    {
        LOG_ERROR("it is failed to transport data into next module, the callback is null");
        makeResponse(GSE_ERROR, std::string("it is failed to transport data into next module, the callback is null"), response);
        delete pDataCell;
        return EVHTP_RES_SERVERR;
    }

    m_callback(pDataCell, m_caller);
    makeResponse(GSE_SUCCESS, std::string("success"), response);
    return HTTP_STATUS_CODE_200;
}

void HttpReceiverHandler::SetServerIp(const std ::string &ip)
{
    m_servIp = ip;
}
void HttpReceiverHandler::SetServerPort(uint16_t port)
{
    m_servPort = port;
}

void HttpReceiverHandler::makeResponse(int errorCode, const std::string &message, std::string &response) noexcept
{
    response.assign("{\"code\":");
    response.append(gse::tools::strings::ToString(errorCode));
    response.append(",\"message\":\"");
    response.append(message);
    response.append("\"}");
}

uint32_t HttpReceiverHandler::getChannelIDFromURI(const char *uri) noexcept
{
    std::string uri_str(uri);
    std::size_t pos = uri_str.find_last_of("/");
    LOG_DEBUG("split the uri (%s), the channel id pos is (%u)", uri, pos);
    if ((pos + 1) != std::string::npos)
    {
        return gse::tools::strings::StringToUint32(uri_str.substr(pos + 1));
    }
    return 0;
}

HttpReceiver::HttpReceiver()
    : m_listennerFd(-1), m_httpserver(NULL)
{
}

HttpReceiver::~HttpReceiver()
{
}

int HttpReceiver::StartMigrationSerivce()
{
    m_migrationClient = std::unique_ptr<gse::net::MigrationClient>(new gse::net::MigrationClient(m_servPort, 20));
    int domainListenFd = -1;
    if (m_migrationClient->ConnectDomainSocket() != GSE_SUCCESS)
    {
        LOG_WARN("failed to connect domain socket");
    }
    else
    {
        m_listennerFd = m_migrationClient->MigrateListennerFd();
        domainListenFd = m_migrationClient->MigrateDomainSocketListenner();
    }

    auto pFuncGetListernnerFd = std::bind(&gse::net::http::HTTPServer::GetEvhtpListennerFd, m_httpserver);
    auto pFuncStopListenner = std::bind(&gse::net::http::HTTPServer::StoppingListenner, m_httpserver);

    m_migrationServer = std::unique_ptr<gse::net::MigrationServer>(new gse::net::MigrationServer(m_servPort));
    m_migrationServer->SetGetListennerFdCallback(pFuncGetListernnerFd);
    m_migrationServer->SetStopListennerCallback(pFuncStopListenner);
    m_migrationServer->SetFinishedCallback(&DataServer::GracefullyQuit);

    return m_migrationServer->StartMigrationService(domainListenFd);
}

int HttpReceiver::Start()
{
    m_servIp = m_recevierConf->m_bind;
    m_servPort = m_recevierConf->m_port;

    m_httpserver = new gse::net::http::HTTPServer(m_servIp, m_servPort, m_recevierConf->m_workThreadNum);

    auto httpHandler = std::make_shared<HttpReceiverHandler>(m_fnRecvData, m_pCaller);
    httpHandler->SetServerIp(m_servIp);
    httpHandler->SetServerPort(m_servPort);
    m_httpserver->RegisterHandler(SEND_DATA_BY_DATAID, httpHandler);
    m_httpserver->SetBacklogSize(m_recevierConf->m_backlogSize);

    if (DataServer::GetUpgradeFlag())
    {
        StartMigrationSerivce();
        m_httpserver->SetListenerFd(m_listennerFd);
    }

    auto httpserver_thread = std::bind(&gse::net::http::HTTPServer::Start, m_httpserver);
    m_listenThread = std::thread(httpserver_thread);

    LOG_INFO("http receiver[%s] has start on (ip:[%s], port:[%d]) with worker thread[%d]", SAFE_CSTR(m_recevierConf->m_name.c_str()), SAFE_CSTR(m_servIp.c_str()), m_servPort, m_recevierConf->m_workThreadNum);

    return GSE_SUCCESS;
}

int HttpReceiver::Stop()
{
    m_httpserver->Stop();
    return GSE_SUCCESS;
}

void HttpReceiver::Join()
{
    m_listenThread.join();
    LOG_DEBUG("http receiver(%s) joined,", m_recevierConf->m_name.c_str());
}
} // namespace data
} // namespace gse
