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

#include "proxy_exporter.h"
#include "bbx/gse_errno.h"
#include "tools/macros.h"
#include "tools/atomic.h"
#include "net/status.h"
#include "net/inet_addr.h"
#include "net/tcp_connection.h"
#include "net/tcp_client.h"

#include "protocol.h"
#include "dataserver.h"

#include <atomic>
long recv_pack_count = 0;
std::atomic<long> send_tcp_pack_count(0);

namespace gse { 
namespace dataserver {

ProxyExporter::ProxyExporter()
{
    m_nextTcpClient = 0;
    m_nextUdpClient = 0;
    m_ptrSSL = NULL;
    m_needSendByProtocolV1 = false;
    m_needSendByProtocolTCP = true;
    m_fillChannelid = true;
}

ProxyExporter::~ProxyExporter()
{
    if (NULL != m_ptrSSL)
    {
        delete m_ptrSSL;
    }
}

int ProxyExporter::Start()
{

    switch (m_ptrConfWrapper->m_exporterConfTypeEnum)
    {
    case ExporterConfigWrapper::DataFlowConfType:
        m_upConfLock.RLock();
        if (!startWithDataFlow(m_ptrConfWrapper->m_conf.m_exporterConf))
        {
            m_upConfLock.UnLock();
            return GSE_ERROR;
        }
        m_upConfLock.UnLock();
        return GSE_SUCCESS;
    case ExporterConfigWrapper::DataIDConfType:
        m_upConfLock.RLock();
        if (!startWithDataID(m_ptrConfWrapper->m_conf.m_ptrDataIDConfig))
        {
            m_upConfLock.UnLock();
            return GSE_ERROR;
        }
        m_upConfLock.UnLock();
        return GSE_SUCCESS;
    default:
        break;
    }

    LOG_INFO("it is failed to start exporter[%s] with the config type [%d]", SAFE_CSTR(m_name.c_str()), m_ptrConfWrapper->m_exporterConfTypeEnum);

    return GSE_ERROR;
}

int ProxyExporter::Stop()
{
    // close tcp connection
    std::size_t max_count =0;

    max_count = m_tcpClientVector.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        gse::net::TcpClient *ptr_tcp_client = m_tcpClientVector.at(idx);
        ptr_tcp_client->Disconnect();
        delete ptr_tcp_client;
    }
    m_tcpClientVector.clear();

    // close udp connection
    max_count = m_udpClientVector.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        gse::net::udp::UDPClient *ptr_udp_client = m_udpClientVector.at(idx);
        delete ptr_udp_client;
    }
    m_udpClientVector.clear();
    return GSE_SUCCESS;
}

bool ProxyExporter::writeToByTCP(const char *ptrData, uint32_t dataLength, DataCell *pDataCell)
{
    std::size_t tcpclient_max_count = m_tcpClientVector.size();
    if (0 == tcpclient_max_count)
    {
        LOG_WARN("there is not any valid connection used to send data");
        return GSE_ERROR;
    }

    int next_idx = ::abs(gse::tools::atomic::AtomAddAfter(&m_nextTcpClient) % tcpclient_max_count);
    gse::net::TcpClient *ptr_tcp_client = m_tcpClientVector.at(next_idx);

    if (ptr_tcp_client->Send(ptrData, dataLength) < 0)
    {
        std::string  remote_server_ip = ptr_tcp_client->GetServerIp();
        int port = ptr_tcp_client->GetServerPort();
        LOG_ERROR("it is failed to send data to the remote server, will to reconnect the remote server(%s:%d) ", SAFE_CSTR(remote_server_ip.c_str()),port);
        return false;
    }

    LOG_DEBUG("send mesage to [%s:%d], mesage len[%d]", ptr_tcp_client->GetServerIp().c_str(), ptr_tcp_client->GetServerPort(), dataLength);
    std::string exportortag = "proxy|tcp|" + m_name + "|" + ptr_tcp_client->GetServerIp() + "|" + gse::tools::strings::ToString(ptr_tcp_client->GetServerPort());
    pDataCell->SetExportorName(exportortag);

    return true;
}

bool ProxyExporter::writeToByUDP(const char *ptrData, uint32_t dataLength, DataCell *pDataCell)
{
    std::size_t udpclient_max_count = m_udpClientVector.size();
    if (0 == udpclient_max_count)
    {
        LOG_WARN("there is not any valid connection used to send data");
        return GSE_ERROR;
    }

    int next_idx = abs(int(gse::tools::atomic::AtomAddAfter(&m_nextUdpClient) % udpclient_max_count));
    gse::net::udp::UDPClient *ptr_udp_client = m_udpClientVector.at(next_idx);

    if (ptr_udp_client->Send(ptrData, dataLength) < 0)
    {

        LOG_ERROR("it is failed to send data to the remote server, the data length is %d", dataLength);
        return false;
    }
    std::string exportortag = "proxy|udp|" + m_name + "|" + ptr_udp_client->GetServerIp() + "|" + gse::tools::strings::ToString(ptr_udp_client->GetPort());
    pDataCell->SetExportorName(exportortag);
    return true;
}

int ProxyExporter::Write(DataCell *pDataCell)
{
    if(NULL == pDataCell)
    {
        LOG_WARN("the parameter is valid. the pointer of data cell is [%x]", pDataCell);
        return GSE_ERROR;
    }

    std::string extension = "proxy";
    pDataCell->PushExtension(extension);
    static uint64_t cnt_id = 0;
    // construct the extension
//    GseDataPkgCodec gse_data_codec;
//    GseDataPkgCodecV1 gse_data_codec_v1;

    GSEDataPackage           m_gseDataPackage;
    GSEDataPackageV1        gse_data_packagev1;
    int channelid = 0;

    char* ptr_msg = nullptr;
    int msg_len = 0;

    if (!m_needSendByProtocolV1)
    {
        std::string id = gse::tools::strings::ToString(cnt_id++);

        if (m_fillChannelid)
        {
            channelid= pDataCell->GetChannelID();
        }

        m_gseDataPackage.SetId(id);
        m_gseDataPackage.SetChannelId(channelid);

        m_gseDataPackage.SetExterntions(m_extensions);

        m_gseDataPackage.SetCreateTimestamp(pDataCell->GetCreationTimestamp());
        m_gseDataPackage.SetArrivedTimestamp(pDataCell->GetArrivedTimestamp());
        m_gseDataPackage.Pack(pDataCell->GetDataBuf(), pDataCell->GetDataBufLen());
        msg_len = m_gseDataPackage.GetDataLength();
        ptr_msg = m_gseDataPackage.GetDataPointer();
    }
    else
    {
        if (m_fillChannelid)
        {
            channelid= pDataCell->GetChannelID();
        }
        gse_data_packagev1.SetChannelId(channelid);

        gse_data_packagev1.Pack(pDataCell->GetDataBuf(), pDataCell->GetDataBufLen());
        msg_len = gse_data_packagev1.GetDataLength();
        ptr_msg = gse_data_packagev1.GetDataPointer();

        LOG_DEBUG("send data to next ds by proxy for the channel id (%d) the data length is (%d), the package data is %d",
                  pDataCell->GetChannelID(), pDataCell->GetDataBufLen(), msg_len);
    }

    if (m_needSendByProtocolTCP)
    {
        if (!writeToByTCP(ptr_msg, msg_len, pDataCell))
        {
            LOG_ERROR("it is failed to send any data for the channel id (%d) to next ds by tcp", pDataCell->GetChannelID());
            return GSE_ERROR;
        }
    }
    else
    {

        if (!writeToByUDP(ptr_msg, msg_len, pDataCell))
        {
            LOG_ERROR("it is failed to send any data for the channel id (%d) to next ds by udp ", pDataCell->GetChannelID());
            return GSE_ERROR;
        }
    }

    DataServer::Instance().GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_OUTPUT_STATE));
    return GSE_SUCCESS;
}

void ProxyExporter::connectCallback(const gse::net::TcpConnectionPtr &connection)
{
    if (connection->IsDisconnected())
    {
        LOG_WARN("connection disconnected, peer host [%s:%d]", connection->PeerAddress().ToIPString().c_str(), connection->PeerAddress().ToPort());
    }

    if (connection->IsConnected())
    {
        LOG_DEBUG("connect peer host [%s:%d] success", connection->PeerAddress().ToIPString().c_str(), connection->PeerAddress().ToPort());
    }
}

bool ProxyExporter::startTCPClient(DSProxyExporterConf *ptrProxyExporterConf)
{

    //m_tcpClient.reset(new gse::net::TcpClient(ptrProxyExporterConf->m_addresses, m_serverPort));
    std::size_t max_count = ptrProxyExporterConf->m_addresses.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        Address address = ptrProxyExporterConf->m_addresses.at(idx);
        for (int client_idx = 0; client_idx < ptrProxyExporterConf->m_connectionNumEachAddress; ++client_idx)
        {
            gse::net::TcpClient *ptr_tcp_client = new gse::net::TcpClient(address.m_ip, address.m_port);
            ptr_tcp_client->SetProtoType(gse::net::EnmProtoType::DATA_PROTO);
            if (ptrProxyExporterConf->m_certPath != "")
            {
                std::string cafile = ptrProxyExporterConf->m_certPath + "/" + "gseca.crt";
                std::string certfile = ptrProxyExporterConf->m_certPath + "/" + "gse_server.crt";
                std::string keyfile = ptrProxyExporterConf->m_certPath + "/" + "gse_server.key";
                std::string passwdfile = ptrProxyExporterConf->m_certPath + "/" + "cert_encrypt.key";
                ptr_tcp_client->EnableSSL(cafile.c_str(), certfile.c_str(), keyfile.c_str(), passwdfile.c_str());
            }

            int ret = ptr_tcp_client->Connect();
            ptr_tcp_client->SetConnectionCallback(std::bind(&ProxyExporter::connectCallback, this, std::placeholders::_1));
            m_tcpClientVector.push_back(ptr_tcp_client);
        }
    }

    return true;
}

bool ProxyExporter::startUDPClient(DSProxyExporterConf *ptrProxyExporterConf)
{
    std::size_t max_count = ptrProxyExporterConf->m_addresses.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        Address address = ptrProxyExporterConf->m_addresses.at(idx);
        for (int client_idx = 0; client_idx < ptrProxyExporterConf->m_connectionNumEachAddress; ++client_idx)
        {

            gse::net::udp::UDPClient *ptr_udp_client = new gse::net::udp::UDPClient();
            if (!ptr_udp_client->Connect(address.m_port, address.m_ip))
            {
                LOG_ERROR("it is failed to connect the remote address(%s:%d) by udp", SAFE_CSTR(address.m_ip.c_str()), address.m_port);
                delete ptr_udp_client;
                break;
            }

            m_udpClientVector.push_back(ptr_udp_client);

        }
    }

    return true;
}

bool ProxyExporter::startWithChannelID(ChannelIDStorage *ptrChannelIDConfig)
{

    return false;
}

bool ProxyExporter::startWithDataFlow(ExporterConf *ptrExporterConf)
{
    DSProxyExporterConf *ptr_proxy_exporter_conf = ptrExporterConf->m_dsProxyConf;
    m_needSendByProtocolV1 = gse::tools::strings::CompareStringCase(ptr_proxy_exporter_conf->m_proxyVersion, PROXY_VERSION_1);
    m_needSendByProtocolTCP = gse::tools::strings::CompareStringCase(ptr_proxy_exporter_conf->m_proxyProtocol, PROXY_PROTOCOL_TCP);
    m_fillChannelid = ptr_proxy_exporter_conf->m_fillChannelid;
    m_extensions = ptrExporterConf->m_extensions;
    if (m_needSendByProtocolTCP)
    {
        LOG_INFO("it will to start the tcp client by protocol version (%s) by the protocol (%s)", SAFE_CSTR(ptr_proxy_exporter_conf->m_proxyVersion.c_str()), SAFE_CSTR(ptr_proxy_exporter_conf->m_proxyProtocol.c_str()));
        return startTCPClient(ptr_proxy_exporter_conf);
    }
    else
    {
        // now only TCP and UDP
        LOG_INFO("it will to start the udp client by protocol version (%s) by the protocol (%s)", SAFE_CSTR(ptr_proxy_exporter_conf->m_proxyVersion.c_str()), SAFE_CSTR(ptr_proxy_exporter_conf->m_proxyProtocol.c_str()));
        return startUDPClient(ptr_proxy_exporter_conf);
    }

    return false;
}

bool ProxyExporter::startWithDataID(StorageConfigType *ptrStorageConfig)
{
    return false;
}
}
}
