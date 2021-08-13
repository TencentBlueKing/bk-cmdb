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

#include "tcpreceiver.h"
#include <memory>
#include <functional>

#include "conf/configurator.h"
#include "tools/macros.h"
#include "dataserver.h"
namespace gse { 
namespace dataserver {


TcpReceiver::TcpReceiver()
{
}

TcpReceiver::~TcpReceiver()
{
    //
}

void TcpReceiver::MsgHandler(const std::shared_ptr<gse::net::TcpConnection> &conn, const std::shared_ptr<gse::net::Message> &msg)
{
    DataCell *pDataCell = new DataCell();

    gse::net::InetAddr peeraddr = conn->PeerAddress();
    gse::net::InetAddr localaddr = conn->LocalAddress();
    //pDataCell->SetSourceIp((uint32_t)conn.GetSourceIp());
    pDataCell->SetSourceIp(peeraddr.ToIP());
    pDataCell->SetSourcePort(peeraddr.ToPort());
    pDataCell->SetServerIp(localaddr.ToIP());
    pDataCell->SetServerPort(localaddr.ToPort());
    pDataCell->SetChannelProtocol("TCP");
    const char* msgbuff = msg->GetBody();
    int msglen = msg->GetBodyLen();
	OPMetric::TcpMsgInc();
    if (GSE_SUCCESS != pDataCell->CopyData(msgbuff, msglen))
    {
        LOG_WARN("fail to copy tcp data into data cell(buf:[0x%x], buflen:[%d], sourceip:[%d],sourceport:[%d])  on (ip: [%s], port: [%d])",
                    msg->GetBody(), msg->GetMsgLen(), peeraddr.ToIPString().c_str(),
                    peeraddr.ToPort(), localaddr.ToIPString().c_str(), localaddr.ToPort());
        //m_ptrOPSReport->PutOpsData(pDataCell->ToOPS(kLostState));
        delete pDataCell;
        pDataCell = NULL;
        return;
    }

    // report ops
    //m_ptrOPSReport->ReportReceiverOPS(pDataCell);

    // transmit data cell to channel
    if (NULL == m_fnRecvData)
    {
        LOG_WARN("the pointer of deal recv data function is NULL in tcp receiver[%s] on (ip: [%s], port: [%d])", SAFE_CSTR(m_recevierConf->m_name.c_str()), SAFE_CSTR(m_servIp.c_str()), m_servPort);
        //m_ptrOPSReport->PutOpsData(pDataCell->ToOPS(kLostState));
        delete pDataCell;
        pDataCell = NULL;
        return;
    }

    LOG_DEBUG("receiver msg(buf:0x%x, len:%d) from client (ip: [%s], port: [%d])", msgbuff, msglen, peeraddr.ToIPString().c_str(), peeraddr.ToPort());
    m_fnRecvData(pDataCell, m_pCaller);
}

void TcpReceiver::OnConnectionEvent(const std::shared_ptr<gse::net::TcpConnection> &conn)
{
    if (conn->IsConnected())
	{
		LOG_DEBUG("new connection create, client:%s:%d", conn->PeerAddress().ToIPString().c_str(), conn->PeerAddress().ToPort());
	}
	else
	{
		LOG_WARN("disconnected from:%s:%d", conn->PeerAddress().ToIPString().c_str(), conn->PeerAddress().ToPort());
	}
}


void TcpReceiver::setProtoMsgHead()
{
    if (NULL == m_recevierConf)
    {
        LOG_ERROR("the configure of tcp receiver is empty, please check configure");
        return;
    }

    switch (m_recevierConf->m_protoStack)
    {
    case PS_TYPE_GSEDATA:
        m_tcpServer->SetProtoType(gse::net::EnmProtoType::DATA_PROTO);
        LOG_DEBUG("the protocol stack is %d, set proto type is DATA_PROTO", m_recevierConf->m_protoStack);
        break;
    case PS_TYPE_GSEDATA_V1:
        m_tcpServer->SetProtoType(gse::net::EnmProtoType::DATA_PROTO);
        LOG_DEBUG("the protocol stack is %d, set proto type is DATA_PROTO", m_recevierConf->m_protoStack);
        break;
    case PS_TYPE_GSEDATA_V1_FOR_TGLOG_PROXY:
        //len = sizeof(GseDataTglogMsgHead);
        break;
    case PS_TYPE_GSEDATA_V2:
        m_tcpServer->SetProtoType(gse::net::EnmProtoType::GENERIC_PROTO);
        LOG_DEBUG("the protocol stack is %d, set proto type is GENERIC_PROTO", m_recevierConf->m_protoStack);
        break;
    default:
        LOG_WARN("unknow the data protocol type[%d] in tcp receiver[%s] on (ip: [%s], port: [%d])", m_recevierConf->m_protocol, SAFE_CSTR(m_recevierConf->m_name.c_str()), SAFE_CSTR(m_servIp.c_str()), m_servPort);
        break;
    }
    LOG_DEBUG("the protocol stack is %d", m_recevierConf->m_protoStack);
    return;
}

int TcpReceiver::Start()
{
    if (NULL == m_recevierConf)
    {
        LOG_ERROR("the configure of tcp receiver is empty, please check configure");
        return GSE_SYSTEMERROR;
    }

    m_servIp = m_recevierConf->m_bind;
    m_servPort = m_recevierConf->m_port;
    if ("" == m_servIp)
    {
        m_servIp = gse::tools::net::GetMachineIp();
    }

    m_tcpServer = std::make_shared<gse::net::TcpServer>(m_servIp, m_servPort, m_recevierConf->m_workThreadNum);
    auto msg_handler = std::bind(&TcpReceiver::MsgHandler, this, std::placeholders::_1, std::placeholders::_2);
    auto connection_handler = std::bind(&TcpReceiver::OnConnectionEvent, this, std::placeholders::_1);

    m_tcpServer->SetMessageCallback(msg_handler);
    m_tcpServer->SetConnectionCallback(connection_handler);
    m_tcpServer->SetMaxMsgDiffTime(-1);
    //m_tcpServer->SetProtoType(gse::net::EnmProtoType::DATA_PROTO);
    setProtoMsgHead();
    bool flag = DataServer::Instance().GetUpgradeFlag();

    m_tcpServer->SetUpgradeFlag(flag);

    LOG_DEBUG("tcp server start with upgrade flag:%d", flag);


    if (!m_recevierConf->m_certPath.empty())
    {
        std::string cafile = m_recevierConf->m_certPath + "/gseca.crt";
        std::string certfile = m_recevierConf->m_certPath + "/gse_server.crt";
        std::string keyfile = m_recevierConf->m_certPath + "/gse_server.key";
        std::string passwdfile = m_recevierConf->m_certPath + "/cert_encrypt.key";

        //m_tcpserver->EnableSSL(cafile.c_str(), certfile.c_str(), keyfile.c_str(), passwdfile.c_str());
        m_tcpServer->EnableSSL(cafile.c_str(), certfile.c_str(), keyfile.c_str(), passwdfile.c_str());
        LOG_DEBUG("Enable ssl, ca:%s, cert:%s, keyfile:%s, passwdfile:%s", cafile.c_str(), certfile.c_str(), keyfile.c_str(), passwdfile.c_str());
        return GSE_SYSTEMERROR;
    }

    auto tcpserver_thread = std::bind(&gse::net::TcpServer::Start, m_tcpServer);
    m_tcpServerThread = std::make_shared<std::thread>(tcpserver_thread);

    LOG_INFO("tcp receiver[%s] has start on (ip: [%s], port: [%d]) with worker thread[%d]", SAFE_CSTR(m_recevierConf->m_name.c_str()), SAFE_CSTR(m_servIp.c_str()), m_servPort, m_recevierConf->m_workThreadNum);

    return GSE_SUCCESS;
}

int TcpReceiver::Stop()
{
    //close();
    LOG_INFO("tcp receiver[%s] on (ip: [%s], port: [%d]) will stop", SAFE_CSTR(m_recevierConf->m_name.c_str()), SAFE_CSTR(m_servIp.c_str()), m_servPort);
    return GSE_SUCCESS;
}

void TcpReceiver::Join()
{
    //
    m_tcpServerThread->join();

}


}
}
