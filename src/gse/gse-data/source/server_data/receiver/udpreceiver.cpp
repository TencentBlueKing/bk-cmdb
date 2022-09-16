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

#include "udpreceiver.h"

#include "bbx/gse_errno.h"
#include "conf/configurator.h"
#include "datacell.h"
#include "net/udp/udp_server.h"
#include "ops/ops.h"
#include "tools/macros.h"
#include "tools/net.h"
#include "tools/strings.h"

#include "dataserver.h"

namespace gse {
namespace data {

static const std::string moduleName = "udpserver";
UdpReceiver::UdpReceiver()
{
    m_serverPort = 0;
    m_listennerFd = -1;
}

UdpReceiver::~UdpReceiver()
{
}

int UdpReceiver::Start()
{
    if (NULL == m_recevierConf)
    {
        LOG_ERROR("the configure of tcp receiver is empty, please check configure");
        return GSE_SYSTEMERROR;
    }

    std::string bindip = m_recevierConf->m_bind;
    m_serverPort = m_recevierConf->m_port;

    auto msghandler = std::bind(&UdpReceiver::MsgHandler, this, std::placeholders::_1, std::placeholders::_2);
    m_udpServer = std::make_shared<gse::net::udp::UdpServer>();

    bool flag = DataServer::GetUpgradeFlag();
    int listennerFd = -1;
    if (flag)
    {
        listennerFd = MigrationUdpListennerFd();
    }

    std::string host = m_recevierConf->m_bind + ":" + gse::tools::strings::ToString(m_serverPort);
    m_udpServer->SetBufferSize(m_recevierConf->m_recvBufSize);
    bool ret = false;
    if (listennerFd > 0)
    {
        std::vector<int> fds = {listennerFd};
        ret = m_udpServer->Init(fds);
    }
    else
    {
        ret = m_udpServer->Init(bindip, m_serverPort);
    }

    if (!ret)
    {
        std::string errmsg = "errmsg:" + std::string(strerror(errno));
        OpHealthZ::AddInitHealthInfo(moduleName + "_" + m_recevierConf->m_name, "failed to start udp server on " + host + ", " + errmsg, -1);
        LOG_ERROR("Failed to init udp server, errmsg:%s", errmsg.c_str());
        return GSE_ERROR;
    }

    m_udpServer->SetMessageHandler(msghandler);
    if (m_recevierConf->m_workThreadNum > 1)
    {
        m_udpServer->EnableWorkThread(m_recevierConf->m_workThreadNum);
    }

    auto udpserver_thread = std::bind(&gse::net::udp::UdpServer::Start, m_udpServer);
    m_udpServerThread = std::make_shared<std::thread>(udpserver_thread);

    LOG_INFO("udp receiver[%s] has start on (ip:[%s], port:[%d]) with worker thread[%d]", SAFE_CSTR(m_recevierConf->m_name.c_str()), SAFE_CSTR(bindip.c_str()), m_serverPort, m_recevierConf->m_workThreadNum);
    OpHealthZ::AddInitHealthInfo(moduleName + "_" + m_recevierConf->m_name, "ok", 0);
    if (flag)
    {
        StartMigrationSerivce();
    }

    return GSE_SUCCESS;
}

int UdpReceiver::Stop()
{
    LOG_INFO("udp receiver[%s] has start on (ip:[%s], port:[%d]) will stop", SAFE_CSTR(m_recevierConf->m_name.c_str()), SAFE_CSTR(m_servIp.c_str()), m_servPort);
    m_udpServer->Stop(true);
    return GSE_SUCCESS;
}

void UdpReceiver::Join()
{
    m_udpServerThread->join();
}

void UdpReceiver::MsgHandler(gse::eventloop::EventLoop* loop, gse::net::udp::UdpMessagePtr msg)
{
    OPMetric::UdpMsgInc(1);

    LOG_DEBUG("recv msg from %s:%d, msg:%s", msg->Remote().ToIPString().c_str(), msg->Remote().ToPort(), msg->Buffer());
    uint32_t channeldid = Configurator::getDefaultTglogChannelId();

    DataCell* pDataCell = new DataCell();
    pDataCell->SetSourceIp(msg->Remote().ToIPString());
    pDataCell->SetSourcePort(msg->Remote().ToPort());
    pDataCell->SetServerIP(DataServer::GetAdvertiseIP());
    pDataCell->SetServerPort(m_serverPort);
    pDataCell->SetChannelProtocol("UDP");
    pDataCell->SetChannelID(channeldid);
    char* buf = msg->Buffer();
    int len = msg->GetBufferLen();

    if (GSE_SUCCESS != pDataCell->CopyData(buf, len))
    {
        LOG_WARN("fail to copy udp data into data cell(buf:0x%x, buflen:%d, sourceip:%s,sourceport:%d)",
                 buf, len, msg->Remote().ToIPString().c_str(), msg->Remote().ToPort());
        DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_LOST_STATE));
        delete pDataCell;
        pDataCell = NULL;
        return;
    }

    if (NULL == m_fnRecvData)
    {
        LOG_WARN("the pointer of deal recv data function is NULL in udp receiver[%s] on (ip:[%s], port:[%d])", SAFE_CSTR(m_recevierConf->m_name.c_str()), SAFE_CSTR(m_servIp.c_str()), m_servPort);
        DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_LOST_STATE));
        delete pDataCell;
        pDataCell = NULL;
        return;
    }

    m_fnRecvData(pDataCell, m_pCaller);
}

int UdpReceiver::MigrationUdpListennerFd()
{
    m_migrationClient = std::unique_ptr<gse::net::MigrationClient>(new gse::net::MigrationClient(m_serverPort, 20));

    if (m_migrationClient->ConnectDomainSocket() != GSE_SUCCESS)
    {
        LOG_WARN("failed to connect domain socket");
    }
    else
    {
        m_listennerFd = m_migrationClient->MigrateListennerFd();
        LOG_DEBUG("migrate udp server listenner socket fd:%d", m_listennerFd);
    }

    return m_listennerFd;
}

int UdpReceiver::StartMigrationSerivce()
{
    int domainsocketListennerFd = m_migrationClient->MigrateDomainSocketListenner();
    auto pFuncGetListernnerFd = std::bind(&gse::net::udp::UdpServer::GetListennerFd, m_udpServer);

    auto pFuncStopListenner = std::bind(&gse::net::udp::UdpServer::StopListenner, m_udpServer);

    m_migrationServer = std::unique_ptr<gse::net::MigrationServer>(new gse::net::MigrationServer(m_serverPort));
    m_migrationServer->SetGetListennerFdCallback(pFuncGetListernnerFd);
    m_migrationServer->SetStopListennerCallback(pFuncStopListenner);

    m_migrationServer->SetFinishedCallback(&DataServer::GracefullyQuit);
    return m_migrationServer->StartMigrationService(domainsocketListennerFd);
}

} // namespace data
} // namespace gse
