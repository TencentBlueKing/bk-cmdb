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

#include "net/udp/udp_server.h"
#include "tools/strings.h"
#include "tools/net.h"
#include "bbx/gse_errno.h"
#include "datacell.h"
#include "tools/macros.h"
#include "opscollection/ops.h"
#include "conf/configurator.h"

#include "dataserver.h"

namespace gse { 
namespace dataserver {


UdpReceiver::UdpReceiver()
{
    m_serverIp = 0;
    m_serverPort = 0;
}

UdpReceiver::~UdpReceiver()
{
    //

}

int UdpReceiver::Start()
{
    if (NULL == m_recevierConf)
    {
        LOG_ERROR("the configure of tcp receiver is empty, please check configure");
        return GSE_SYSTEMERROR;
    }
    OPMetric::UdpMsgInc();
    std::string bindip = m_recevierConf->m_bind;
    int port = m_recevierConf->m_port;
    if ("" == bindip)
    {
        bindip = gse::tools::net::GetMachineIp();
    }

    m_serverIp = gse::tools::strings::StringToInt32(bindip);
    m_serverPort = port;

    auto msghandler = std::bind(&UdpReceiver::MsgHandler, this, std::placeholders::_1, std::placeholders::_2);

    m_udpServer = std::make_shared<gse::net::udp::UdpServer>();
    m_udpServer->Init(bindip, port);
    m_udpServer->SetBufferSize(m_recevierConf->m_recvBufSize);
    m_udpServer->SetMessageHandler(msghandler);
    m_udpServer->EnableWorkThread(m_recevierConf->m_workThreadNum);

    auto udpserver_thread = std::bind(&gse::net::udp::UdpServer::Start, m_udpServer);
    m_udpServerThread = std::make_shared<std::thread>(udpserver_thread);

    LOG_INFO("udp receiver[%s] has start on (ip:[%s], port:[%d]) with worker thread[%d]", SAFE_CSTR(m_recevierConf->m_name.c_str()), SAFE_CSTR(bindip.c_str()), m_serverPort, m_recevierConf->m_workThreadNum);

    return GSE_SUCCESS;
}

int UdpReceiver::Stop()
{
    LOG_INFO("udp receiver[%s] has start on (ip:[%s], port:[%d]) will stop", SAFE_CSTR(m_recevierConf->m_name.c_str()), SAFE_CSTR(m_servIp.c_str()), m_servPort);
    return GSE_SUCCESS;
}

void UdpReceiver::Join()
{
    m_udpServerThread->join();
}


void UdpReceiver::MsgHandler(gse::eventloop::EventLoop* loop, gse::net::udp::UdpMessagePtr & msg)
{
    LOG_DEBUG("recv msg from %s:%d, msg:%s", msg->Remote().ToIPString().c_str(), msg->Remote().ToPort(), msg->Buffer());
    uint32_t channeldid = Configurator::getDefaultTglogChannelId();

    DataCell* pDataCell = new DataCell();
    pDataCell->SetSourceIp(msg->Remote().ToIP());
    pDataCell->SetSourcePort(msg->Remote().ToPort());

    pDataCell->SetServerIp(m_serverIp);
    pDataCell->SetServerPort(m_serverPort);
    pDataCell->SetChannelProtocol("UDP");
    pDataCell->SetChannelID(channeldid);
    char * buf = msg->Buffer();
    int len = msg->GetBufferLen();
    if (GSE_SUCCESS != pDataCell->CopyData(buf, len))
    {
        LOG_WARN("fail to copy udp data into data cell(buf:0x%x, buflen:%d, sourceip:%s,sourceport:%d)",
            buf, len, msg->Remote().ToIPString().c_str(), msg->Remote().ToPort());
        DataServer::Instance().GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_LOST_STATE));
        delete pDataCell;
        pDataCell = NULL;
        return;
    }

    // report ops
    //m_ptrOPSReport->ReportReceiverOPS(pDataCell);

    // transmit data cell to channel
    if (NULL == m_fnRecvData)
    {
        LOG_WARN("the pointer of deal recv data function is NULL in udp receiver[%s] on (ip:[%s], port:[%d])", SAFE_CSTR(m_recevierConf->m_name.c_str()), SAFE_CSTR(m_servIp.c_str()), m_servPort);
        DataServer::Instance().GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_LOST_STATE));
        delete pDataCell;
        pDataCell = NULL;
        return;
    }

    m_fnRecvData(pDataCell, m_pCaller);


}
}
}
