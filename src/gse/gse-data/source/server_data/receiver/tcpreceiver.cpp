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
#include <functional>
#include <memory>

#include "conf/configurator.h"
#include "dataserver.h"
#include "inner/protocol_head.h"
#include "loadbalance_schedule.h"
#include "net/message_factory.hpp"
#include "net/migration/migration_client.h"
#include "net/migration/migration_server.h"
#include "ops/op_healthz.h"
#include "tools/macros.h"
#include "tools/rapidjson_macro.h"

namespace gse {
namespace data {

const static std::string moduleName = "tcpserver";

TcpReceiver::TcpReceiver()
    : m_migrationServer(nullptr), m_migrationClient(nullptr), m_listennerFd(-1)
{
}

TcpReceiver::~TcpReceiver()
{
}

void TcpReceiver::OnConnectionEvent(const std::shared_ptr<gse::net::TcpConnection> conn)
{
    if (conn->IsConnected())
    {
        LOG_DEBUG("new connection create, client:%s:%d", conn->PeerAddress().ToIPString().c_str(), conn->PeerAddress().ToPort());
        OPMetric::ConnectionCountInc(conn->GetID());
    }
    else if (conn->IsDisconnected())
    {
        LOG_INFO("tcp connection(%s:%d) disconnected", conn->PeerAddress().ToIPString().c_str(), conn->PeerAddress().ToPort());
        OPMetric::ConnectionCountDec(conn->GetID());
    }
}

void TcpReceiver::InitProto()
{
    if (NULL == m_recevierConf)
    {
        LOG_ERROR("the configure of tcp receiver is empty, please check configure");
        return;
    }

    switch (m_recevierConf->m_protoStack)
    {
    case PS_TYPE_GSEDATA:
        m_protoType = gse::net::EnmProtoType::DATA_DYNAMIC_PROTO;
        m_tcpServer->SetProtoType(m_protoType);
        LOG_DEBUG("the protocol stack is %d, set proto type is DATA_PROTO", m_recevierConf->m_protoStack);
        break;
    case PS_TYPE_GSEDATA_V1:
        m_tcpServer->SetProtoType(gse::net::EnmProtoType::DATA_PROTO);
        m_protoType = gse::net::EnmProtoType::DATA_PROTO;
        LOG_DEBUG("the protocol stack is %d, set proto type is DATA_PROTO", m_recevierConf->m_protoStack);
        break;
    case PS_TYPE_GSEDATA_V1_FOR_TGLOG_PROXY:
        m_tcpServer->SetProtoType(gse::net::EnmProtoType::DATA_TGLOG_PROTO);
        m_protoType = gse::net::EnmProtoType::DATA_TGLOG_PROTO;
        LOG_DEBUG("the protocol stack is %d, set proto type is DATA_TGLOG_PROTO", m_recevierConf->m_protoStack);
        break;
    case PS_TYPE_GSEDATA_GENERAL:
        m_tcpServer->SetProtoType(gse::net::EnmProtoType::DATA_GENERAL_PROTO);
        m_protoType = gse::net::EnmProtoType::DATA_GENERAL_PROTO;
        LOG_DEBUG("the protocol stack is %d, set proto type is DATA_GENERAL_PROTO", m_recevierConf->m_protoStack);
        break;
    default:
        LOG_WARN("unknow the data protocol type[%d] in tcp receiver[%s] on (ip: [%s], port: [%d])", m_recevierConf->m_protocol, SAFE_CSTR(m_recevierConf->m_name.c_str()), SAFE_CSTR(m_servIp.c_str()), m_servPort);
        break;
    }
    return;
}

void TcpReceiver::MigrationTcpListennerFd()
{
    m_migrationClient = std::unique_ptr<gse::net::MigrationClient>(new gse::net::MigrationClient(m_servPort, 20));

    if (m_migrationClient->ConnectDomainSocket() != GSE_SUCCESS)
    {
        LOG_WARN("failed to connect domain socket");
    }
    else
    {
        m_listennerFd = m_migrationClient->MigrateListennerFd();
    }
}

int TcpReceiver::StartMigrationSerivce()
{
    int domainsocket_listenner_fd = m_migrationClient->MigrateDomainSocketListenner();
    auto pFuncGetListernnerFd = std::bind(&gse::net::TcpServer::GetListenerFd, m_tcpServer);
    auto pFuncStopListenner = std::bind(&gse::net::TcpServer::StopListenner, m_tcpServer);

    m_migrationServer = std::unique_ptr<gse::net::MigrationServer>(new gse::net::MigrationServer(m_servPort));
    m_migrationServer->SetGetListennerFdCallback(pFuncGetListernnerFd);
    m_migrationServer->SetStopListennerCallback(pFuncStopListenner);
    auto pFuncGetTcpConnection = std::bind(&gse::net::TcpServer::MigrationConnections, m_tcpServer, std::placeholders::_1);
    m_migrationServer->SetGetTcpConnectionCallback(pFuncGetTcpConnection);
    auto pFuncGetTcpConnectioCount = std::bind(&gse::net::TcpServer::GetConnectionCount, m_tcpServer);
    m_migrationServer->SetGetTcpConnectionCountCallback(pFuncGetTcpConnectioCount);
    m_migrationServer->SetFinishedCallback(&DataServer::GracefullyQuit);
    return m_migrationServer->StartMigrationService(domainsocket_listenner_fd);
}

int TcpReceiver::GetListenerFd()
{
    return m_listennerFd;
}

void TcpReceiver::RestoreConnectionCallback(evutil_socket_t fd,
                                            char *evInputBuffer, size_t evInputBufferSize,
                                            char *evOutputBuffer, size_t evOuputBufferSize,
                                            const gse::net::InetAddr &local, const gse::net::InetAddr &remote)
{
    m_tcpServer->RestoreConnection(fd, evInputBuffer, evInputBufferSize, evOutputBuffer, evOuputBufferSize, local, remote, m_tcpServer->GetMaxMessageSize());
}

int TcpReceiver::MigrationConnections()
{

    auto restoreConnectionCb = std::bind(&TcpReceiver::RestoreConnectionCallback, this,
                                         std::placeholders::_1, std::placeholders::_2,
                                         std::placeholders::_3, std::placeholders::_4,
                                         std::placeholders::_5, std::placeholders::_6,
                                         std::placeholders::_7);

    m_migrationClient->SetRestoreMigrationConnectonCallbck(restoreConnectionCb);

    int ret = m_migrationClient->MigrateConnections();

    if (ret != 0)
    {
        LOG_ERROR("failed to migration connections");
    }

    LOG_DEBUG("Migration connection successful");
    return ret;
}

void TcpReceiver::StartTcpServer()
{
    m_tcpServer->Start();
    std::string host = " " + m_recevierConf->m_bind + ":" + gse::tools::strings::ToString(m_recevierConf->m_port) + " ";
    std::string errmsg = "errmsg:" + std::string(strerror(errno));
    OpHealthZ::AddInitHealthInfo(moduleName + "_" + m_recevierConf->m_name, "failed to start tcp server on " + host + ", " + errmsg, -1);
}

void TcpReceiver::SetSheduler(std::shared_ptr<LoadBalanceSchedule> schedule)
{
    m_scheduler = schedule;
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
    m_tcpServer = std::shared_ptr<gse::net::TcpServer>(new gse::net::TcpServer(m_servIp, m_servPort, m_recevierConf->m_workThreadNum, gse::net::CONNID_TYPE_IP_PORT));
    InitProto();

    //信令消息
    auto controlMsgProcessor = std::make_shared<SignalControlMsgProcessor>();
    controlMsgProcessor->SetScheduler(m_scheduler);
    controlMsgProcessor->SetProtoType(m_protoType);

    auto dataReportMsgProcessor = std::make_shared<DataReportMsgProcessor>();
    dataReportMsgProcessor->SetRecvDataCallback(m_fnRecvData, m_pCaller);

    auto opsReportMsgProcessor = std::make_shared<OpsReportMsgProcessor>();
    opsReportMsgProcessor->SetRecvDataCallback(m_fnRecvData, m_pCaller);

    m_tcpServer->RegistMsgProcessor(DATA_ACCESS_REQUEST, controlMsgProcessor);
    m_tcpServer->RegistMsgProcessor(DATA_KEEPALIVE_REQUEST, controlMsgProcessor);

    m_tcpServer->RegistDefaultMsgProcessor(dataReportMsgProcessor);
    auto connection_handler = std::bind(&TcpReceiver::OnConnectionEvent, this, std::placeholders::_1);
    m_tcpServer->SetConnectionCallback(connection_handler);
    m_tcpServer->SetMaxMsgDiffTime(-1);

    bool flag = DataServer::GetUpgradeFlag();

    LOG_DEBUG("tcp server start with upgrade flag:%d", flag);
    bool useSSL = false;
    bool bret = false;
#ifdef _GSE_SSL_

    bret = m_tcpServer->EnableSSL(m_recevierConf->m_caPath.c_str(),
                                  m_recevierConf->m_certPath.c_str(),
                                  m_recevierConf->m_keyPath.c_str(),
                                  m_recevierConf->m_passwdPath.c_str(),
                                  _GSE_SSL_);

#else
    if (!m_recevierConf->m_certPath.empty() &&
        !m_recevierConf->m_caPath.empty() &&
        !m_recevierConf->m_keyPath.empty())
    {
        useSSL = true;

        bret = m_tcpServer->EnableSSL(m_recevierConf->m_caPath.c_str(),
                                      m_recevierConf->m_certPath.c_str(),
                                      m_recevierConf->m_keyPath.c_str(),
                                      m_recevierConf->m_passwdPath.c_str(),
                                      GSE_EDITION_EE);
    }

#endif
    if (useSSL)
    {
        if (!bret)
        {
            OpHealthZ::AddInitHealthInfo(moduleName + "_" + m_recevierConf->m_name, "failed to load cert", -1);
            LOG_DEBUG("Failed to load ssl cert, ca:%s, cert:%s, keyfile:%s, passwdfile:%s", m_recevierConf->m_caPath.c_str(),
                      m_recevierConf->m_certPath.c_str(),
                      m_recevierConf->m_keyPath.c_str(),
                      m_recevierConf->m_passwdPath.c_str());
            return GSE_SYSTEMERROR;
        }

        LOG_DEBUG("Enable ssl, ca:%s, cert:%s, keyfile:%s, passwdfile:%s", m_recevierConf->m_caPath.c_str(),
                  m_recevierConf->m_certPath.c_str(),
                  m_recevierConf->m_keyPath.c_str(),
                  m_recevierConf->m_passwdPath.c_str());
    }

    if (flag)
    {
        MigrationTcpListennerFd();
    }

    m_tcpServer->SetListennerFd(GetListenerFd());
    auto tcpserverThread = std::bind(&TcpReceiver::StartTcpServer, this);
    m_tcpServerThread = std::unique_ptr<std::thread>(new std::thread(tcpserverThread));

    LOG_INFO("tcp receiver[%s] has start on (ip: [%s], port: [%d]) with worker thread[%d]", SAFE_CSTR(m_recevierConf->m_name.c_str()), SAFE_CSTR(m_servIp.c_str()), m_servPort, m_recevierConf->m_workThreadNum);
    m_tcpServer->SetMaxMessageSize(m_recevierConf->m_maxMessageLen > 0 ? m_recevierConf->m_maxMessageLen : TCP_DEFAULT_MESSAGE_LEN);
    if (flag)
    {
        MigrationConnections();
        int ret = StartMigrationSerivce();
        if (ret < 0)
        {
            LOG_ERROR("failed to start migration service, ret:%d", ret);
        }
    }

    gse::net::InetAddr localAddr;
    m_tcpServer->GetLocalAddr(localAddr);

    OpHealthZ::AddInitHealthInfo(moduleName + "_" + m_recevierConf->m_name, "ok", 0);

    return GSE_SUCCESS;
}

int TcpReceiver::Stop()
{
    LOG_INFO("tcp receiver[%s] on (ip: [%s], port: [%d]) is stopping", SAFE_CSTR(m_recevierConf->m_name.c_str()), SAFE_CSTR(m_servIp.c_str()), m_servPort);
    m_tcpServer->Stop();
    Join();

    LOG_INFO("tcp receiver[%s] on (ip: [%s], port: [%d]) stop finished", SAFE_CSTR(m_recevierConf->m_name.c_str()), SAFE_CSTR(m_servIp.c_str()), m_servPort);

    return GSE_SUCCESS;
}

void TcpReceiver::Join()
{
    m_tcpServerThread->join();
    LOG_DEBUG("tcp receiver(%s) joined,", m_recevierConf->m_name.c_str());
}

void TcpReceiver::AddMsgProcessor(int msgType, std::shared_ptr<gse::net::MsgProcessor> processor)
{
    auto it = m_processors.find(msgType);
    if (it == m_processors.end())
    {
        m_processors[msgType] = processor;
    }
    else
    {
        LOG_WARN("repeat register msgprocessor, type:%d", msgType);
    }
}

void TcpReceiver::RegisterMsgProcessor()
{
    for (auto it : m_processors)
    {
        m_tcpServer->RegistMsgProcessor(it.first, it.second);
    }
}

SignalControlMsgProcessor::SignalControlMsgProcessor()
{
}

SignalControlMsgProcessor::~SignalControlMsgProcessor()
{
}

void SignalControlMsgProcessor::SetScheduler(std::shared_ptr<LoadBalanceSchedule> scheduler)
{
    m_scheduler = scheduler;
}

void SignalControlMsgProcessor::SetProtoType(net::EnmProtoType protoType)
{
    m_protoType = protoType;
}

void SignalControlMsgProcessor::ProcessMsg(const gse::net::TcpConnectionPtr conn, const gse::net::MessagePtr msg)
{
    int type = msg->GetMsgNum();
    LOG_DEBUG("recv control message,type:%ld", type);
    switch (type)
    {
    case DATA_ACCESS_REQUEST:
    {

        auto callback = [=](const char *resp, int len) {
            this->SendAccessRespose(conn, resp, len);
        };
        OPMetric::AccessRequestInc();
        bool ret = m_scheduler->EnqueuAccessRequest(callback, msg);
        if (!ret)
        {
            OPMetric::AccessResponseFailedInc();
        }
        break;
    }
    case DATA_KEEPALIVE_REQUEST:
    {
        // send reponse
        std::string keepaliveResponse;
        MarshalKeepaliveReponse(keepaliveResponse);

        auto responseMessage = gse::net::MessageFactory::CreateNew(m_protoType);
        responseMessage->SetMessageType(DATA_KEEPALIVE_RESPONSE);
        responseMessage->BodyFromString(keepaliveResponse);
        LOG_DEBUG("send keepalive response(%s), proto:%d", keepaliveResponse.c_str(), m_protoType);
        conn->Send(responseMessage);
        break;
    }
    default:
        break;
    }

    return;
}

void SignalControlMsgProcessor::MarshalKeepaliveReponse(std::string &value)
{
    rapidjson::StringBuffer strBuf;
    rapidjson::Writer<rapidjson::StringBuffer> writer(strBuf);

    // root start.
    writer.StartObject();

    RAPIDJSON_SET_STRING(writer, "service_id", DataServer::GetNodeId().c_str());
    RAPIDJSON_SET_STRING(writer, "version", DataServer::GetVersion().c_str());
    RAPIDJSON_SET_INT64(writer, "hearttime", TimeCenter::Instance()->GetDateTime());

    // root end.
    writer.EndObject();

    // marshal to string.
    value = strBuf.GetString();

    return;
}

void SignalControlMsgProcessor::SendAccessRespose(const std::shared_ptr<gse::net::TcpConnection> conn, const char *resp, int len)
{
    OPMetric::AccessResponseInc();
    auto responseMessage = gse::net::MessageFactory::CreateNew(m_protoType);
    responseMessage->SetMessageType(DATA_ACCESS_RESPONSE);
    responseMessage->BodyFromBuff(resp, len);
    LOG_DEBUG("send access response(%s)", resp);
    conn->Send(responseMessage);
}

DataReportMsgProcessor::DataReportMsgProcessor()
    : m_fnRecvData(NULL), m_pCaller(NULL)
{
}

DataReportMsgProcessor::~DataReportMsgProcessor()
{
}

void DataReportMsgProcessor::ProcessMsg(const gse::net::TcpConnectionPtr conn, const gse::net::MessagePtr msg)
{
    DataCell *pDataCell = new DataCell();

    gse::net::InetAddr peeraddr = conn->PeerAddress();
    gse::net::InetAddr localaddr = conn->LocalAddress();
    pDataCell->SetSourceIp(peeraddr.ToIPString());
    pDataCell->SetSourcePort(peeraddr.ToPort());
    pDataCell->SetServerIP(localaddr.ToIPString());
    pDataCell->SetServerPort(localaddr.ToPort());
    pDataCell->SetChannelProtocol("TCP");
    //框架中没有GetMsg,GetBody获取的是完整消息，包含head和body
    const char *msgbuff = msg->GetBody();
    int msglen = msg->GetBodyLen();
    if (GSE_SUCCESS != pDataCell->CopyData(msgbuff, msglen))
    {
        LOG_WARN("fail to copy tcp data into data cell(buf:[0x%x], buflen:[%d], sourceip:[%d],sourceport:[%d])  on (ip: [%s], port: [%d])",
                 msg->GetBody(), msg->GetMsgLen(), peeraddr.ToIPString().c_str(),
                 peeraddr.ToPort(), localaddr.ToIPString().c_str(), localaddr.ToPort());
        delete pDataCell;
        pDataCell = NULL;
        return;
    }

    if (NULL == m_fnRecvData)
    {
        LOG_WARN("the pointer of deal recv data function is NULL");
        delete pDataCell;
        pDataCell = NULL;
        return;
    }

    LOG_DEBUG("receiver msg(buf:0x%x, len:%d) from client (ip: [%s], port: [%d])", msgbuff, msglen, peeraddr.ToIPString().c_str(), peeraddr.ToPort());
    m_fnRecvData(pDataCell, m_pCaller);
}

void DataReportMsgProcessor::SetRecvDataCallback(RecvDataCallBack fnRecvData, void *pCaller)
{
    m_fnRecvData = fnRecvData;
    m_pCaller = pCaller;
}

OpsReportMsgProcessor::OpsReportMsgProcessor()
    : m_fnRecvData(NULL), m_pCaller(NULL)
{
}

OpsReportMsgProcessor::~OpsReportMsgProcessor()
{
}

void OpsReportMsgProcessor::ProcessMsg(const gse::net::TcpConnectionPtr conn, const gse::net::MessagePtr msg)
{
    DataCell *pDataCell = new DataCell();

    gse::net::InetAddr peeraddr = conn->PeerAddress();
    gse::net::InetAddr localaddr = conn->LocalAddress();
    pDataCell->SetSourceIp(peeraddr.ToIPString());
    pDataCell->SetSourcePort(peeraddr.ToPort());
    pDataCell->SetServerIP(localaddr.ToIPString());
    pDataCell->SetServerPort(localaddr.ToPort());
    pDataCell->SetChannelProtocol("TCP");
    pDataCell->SetOpsMsg(true);
    pDataCell->SetOpsServiceId(ntohl(msg->GetMsgNum()));

    const char *msgbuff = msg->GetBody();
    int msglen = msg->GetBodyLen();
    if (GSE_SUCCESS != pDataCell->CopyData(msgbuff, msglen))
    {
        LOG_WARN("fail to copy tcp data into data cell(buf:[0x%x], buflen:[%d], sourceip:[%d],sourceport:[%d])  on (ip: [%s], port: [%d])",
                 msg->GetBody(), msg->GetMsgLen(), peeraddr.ToIPString().c_str(),
                 peeraddr.ToPort(), localaddr.ToIPString().c_str(), localaddr.ToPort());
        delete pDataCell;
        pDataCell = NULL;
        return;
    }

    if (NULL == m_fnRecvData)
    {
        LOG_WARN("the pointer of deal recv data function is NULL");
        delete pDataCell;
        pDataCell = NULL;
        return;
    }

    LOG_DEBUG("receiver msg(buf:0x%x, len:%d) from client (ip: [%s], port: [%d])", msgbuff, msglen, peeraddr.ToIPString().c_str(), peeraddr.ToPort());
    m_fnRecvData(pDataCell, m_pCaller);
}

void OpsReportMsgProcessor::SetRecvDataCallback(RecvDataCallBack fnRecvData, void *pCaller)
{
    m_fnRecvData = fnRecvData;
    m_pCaller = pCaller;
}

} // namespace data
} // namespace gse
