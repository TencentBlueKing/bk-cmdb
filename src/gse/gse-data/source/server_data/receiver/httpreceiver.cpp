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
#include "net/http/http_server.h"
#include "bbx/gse_errno.h"
#include "tools/macros.h"
namespace gse { 
namespace dataserver {

HttpReceiver::HttpReceiver()
{
    m_httpserver = NULL;
}

HttpReceiver::~HttpReceiver()
{
}

uint32_t HttpReceiver::getChannelIDFromURI(const char *uri)
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

void HttpReceiver::makeResponse(int errorCode, const std::string &message)
{
    std::string result("{\"code\":");
    result.append(gse::tools::strings::ToString(errorCode));
    result.append(",\"message\":\"");
    result.append(message);
    result.append("\"}");
}

void HttpReceiver::OnHttpMessageHandler(gse::net::http::HttpMessagePtr message, std::string &response)
{
    std::string content;
    message->GetMessage(content);
    OPMetric::HttpMsgInc();
    if (content.empty())
    {
        LOG_ERROR("it is to read nothing from the http request");
        makeResponse(GSE_ERROR, std::string("it is to read nothing from the http request"));
        return;
    }

    gse::net::http::HttpURIPtr urisptr = message->GetHttpURIPtr();
    std::string uri = urisptr->GetURI();
    uint32_t channel_id = getChannelIDFromURI(uri.c_str());

    // create data cell, it' will be deleted in exporter layer
    DataCell *pDataCell = new DataCell();

    // copy data
    pDataCell->SetSourceIp(gse::tools::strings::StringToUint32(message->GetClientIP()));
    //pDataCell->SetSourcePort(message->getClientPort());
    pDataCell->SetServerIp(gse::tools::net::StringToIp(m_servIp, true));
    pDataCell->SetServerPort(m_servPort);
    pDataCell->SetChannelProtocol("HTTP");
    pDataCell->SetChannelID(channel_id);

    if (GSE_SUCCESS != pDataCell->CopyData(content.c_str(), content.size()))
    {
        LOG_ERROR("it is failed to copy the body data (%s) int a datacell", SAFE_CSTR(content.c_str()));
        makeResponse(GSE_ERROR, std::string("it is failed to copy the body data into a datacell"));
        delete pDataCell;
        return;
    }

    if (NULL == m_fnRecvData)
    {
        LOG_ERROR("it is failed to transport data into next module, the callback is null");
        makeResponse(GSE_ERROR, std::string("it is failed to transport data into next module, the callback is null"));
        delete pDataCell;
        return;
    }

    m_fnRecvData(pDataCell, m_pCaller);
    makeResponse(GSE_SUCCESS, std::string("success"));

}

int HttpReceiver::Start()
{
    if (NULL == m_recevierConf)
    {
        LOG_ERROR("the configure of http receiver is empty, please check configure");
        return GSE_SYSTEMERROR;
    }

    m_servIp = m_recevierConf->m_bind;
    m_servPort = m_recevierConf->m_port;
    if ("" == m_servIp)
    {
        m_servIp = gse::tools::net::GetMachineIp();
    }

    auto msg_handler = std::bind(&HttpReceiver::OnHttpMessageHandler, this, std::placeholders::_1, std::placeholders::_2);

    m_httpserver = new gse::net::http::HttpServer(m_servIp, m_servPort, m_recevierConf->m_workThreadNum);
    gse::net::http::HttpURIPtr dataid_uri(new gse::net::http::HttpURI(SEND_DATA_BY_DATAID));
    gse::net::http::HttpHandlerPtr dataid_handler(new gse::net::http::HttpHandler(gse::net::http::HttpMethodPost,
                                                        dataid_uri, msg_handler));

    std::string errormsg;
    m_httpserver->RegistHandler(dataid_handler, errormsg);
    if (errormsg != "")
    {
        LOG_ERROR("register dataid handler failed, errormsg:%s", errormsg.c_str());
        return GSE_ERROR;
    }

    auto httpserver_thread = std::bind(&gse::net::http::HttpServer::Start, m_httpserver);
    m_listenThread = std::thread(httpserver_thread);

    LOG_INFO("http receiver[%s] has start on (ip:[%s], port:[%d]) with worker thread[%d]", SAFE_CSTR(m_recevierConf->m_name.c_str()), SAFE_CSTR(m_servIp.c_str()), m_servPort, m_recevierConf->m_workThreadNum);

    return GSE_SUCCESS;
}

int HttpReceiver::Stop()
{
    return GSE_SUCCESS;
}

void HttpReceiver::Join()
{
    m_listenThread.join();
}
}
}
