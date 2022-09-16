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

#include <stdlib.h>

#include <atomic>

#include <evhtp/evhtp.h>

#include "bbx/gse_errno.h"
#include "dataserver.h"
#include "net/inet_addr.h"
#include "net/status.h"
#include "net/tcp_client.h"
#include "net/tcp_connection.h"
#include "protocol.h"
#include "safe/secure_config.h"
#include "tools/atomic.h"
#include "tools/error.h"
#include "tools/finally.hpp"
#include "tools/macros.h"
#include "tools/pipe.h"
#include "tools/system.h"

namespace gse {
namespace data {

ProxyWriter::ProxyWriter()
{
}

ProxyWriter::~ProxyWriter()
{
}
ProxyExporter::ProxyExporter()
    : m_writer(NULL)
{
}

ProxyExporter::~ProxyExporter()
{
    if (m_writer != NULL)
    {
        m_writer->Stop();
        delete m_writer;
        m_writer = NULL;
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

    case ExporterConfigWrapper::ChannelIDConfType:
        m_upConfLock.RLock();
        if (!startWithChannelID(m_ptrConfWrapper->m_conf.m_ptrChannelIdExporterConfig))
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
    if (m_writer != NULL)
    {
        m_writer->Stop();
        delete m_writer;
        m_writer = NULL;
    }

    return GSE_SUCCESS;
}

int ProxyExporter::Write(DataCell *pDataCell)
{
    if (NULL == pDataCell)
    {
        LOG_WARN("the parameter is valid. the pointer of data cell is [%x]", pDataCell);
        return GSE_ERROR;
    }

    if (!FilterByPlatId(pDataCell))
    {
        LOG_INFO("datacell len[%d] not send to exporter[%s], plat id not config in the exporter", pDataCell->GetBufferLen(), m_exporterName.c_str());
        return GSE_SUCCESS;
    }

    return m_writer->Write(pDataCell);
}

bool ProxyExporter::startWithChannelID(ChannelIdExporterConfig *ptrChannelIDConfig)
{
    auto ptr_proxy_exporter = ptrChannelIDConfig->m_streamToCluster.m_dsproxy;
    DSProxyExporterConf proxy_exporter_conf;
    proxy_exporter_conf.m_certPath = ptr_proxy_exporter->m_certPath;
    proxy_exporter_conf.m_passwdFilePath = ptr_proxy_exporter->m_passwdFilePath;
    proxy_exporter_conf.m_proxyVersion = ptr_proxy_exporter->m_proxyVersion;
    proxy_exporter_conf.m_proxyProtocol = ptr_proxy_exporter->m_proxyProtocol;
    proxy_exporter_conf.m_connectionNumEachAddress = ptr_proxy_exporter->m_connectionNumEachAddress;
    proxy_exporter_conf.m_heartbeat = ptr_proxy_exporter->m_heartbeat;
    proxy_exporter_conf.m_noblock = ptr_proxy_exporter->m_noblock;
    proxy_exporter_conf.m_fillChannelid = ptr_proxy_exporter->m_fillChannelid;

    proxy_exporter_conf.m_httpURI = ptr_proxy_exporter->m_httpURI;
    proxy_exporter_conf.m_thirdPartyCertFile = ptr_proxy_exporter->m_thirdPartyCertFile;
    proxy_exporter_conf.m_thirdPartyCertPasswd = ptr_proxy_exporter->m_thirdPartyCertPasswd;
    proxy_exporter_conf.m_thirdPartyKeyFile = ptr_proxy_exporter->m_thirdPartyKeyFile;
    proxy_exporter_conf.m_isThirdPartyCert = ptr_proxy_exporter->m_isThirdPartyCert;
    proxy_exporter_conf.m_addresses = ptr_proxy_exporter->m_addresses;

    SetExporterName(ptrChannelIDConfig->m_streamToCluster.m_name);

    LOG_DEBUG("start channelid proxy exporter, certpath:%s, passwdfilepath:%s, proxyversion:%s, connection num:%d "
              "proxyprotocol:%s, heartbeat:%d, fillchannelid:%d, httpuri:%s, thirdpartycertpasswd:%s, thirdpartykeyfile:%s, thirdpartykeyfile:%s"
              "isthirdpartycert:%d",
              proxy_exporter_conf.m_certPath.c_str(), proxy_exporter_conf.m_passwdFilePath.c_str(), proxy_exporter_conf.m_proxyVersion.c_str(),
              proxy_exporter_conf.m_connectionNumEachAddress, proxy_exporter_conf.m_proxyProtocol.c_str(),
              proxy_exporter_conf.m_heartbeat, proxy_exporter_conf.m_fillChannelid, proxy_exporter_conf.m_httpURI.c_str(),
              proxy_exporter_conf.m_thirdPartyCertFile.c_str(), proxy_exporter_conf.m_thirdPartyCertPasswd.c_str(),
              proxy_exporter_conf.m_thirdPartyKeyFile.c_str(), proxy_exporter_conf.m_isThirdPartyCert);
    return startProxyExporter(&proxy_exporter_conf);
}

ProxyWriter *ProxyExporter::CreateWriterByProtocol(const std::string &protocol)
{
    ProxyWriter *writer = NULL;
    bool tcp = gse::tools::strings::CompareStringCase(protocol, PROXY_PROTOCOL_TCP);
    bool udp = gse::tools::strings::CompareStringCase(protocol, PROXY_PROTOCOL_UDP);
    bool http = gse::tools::strings::CompareStringCase(protocol, PROXY_PROTOCOL_HTTP);

    if (tcp)
    {
        writer = new ProxyTcpWriter();
    }
    else if (udp)
    {
        writer = new ProxyUdpWriter();
    }
    else if (http)
    {
        writer = new ProxyHTTPWriter();
    }
    else
    {
        LOG_ERROR(" failed to create proxy writer, unkown protocol:%s", protocol.c_str());
    }
    return writer;
}

bool ProxyExporter::startProxyExporter(DSProxyExporterConf *bkdataProxyconfigPtr)
{
    if (m_writer != NULL)
    {
        delete m_writer;
        m_writer = NULL;
    }

    m_writer = CreateWriterByProtocol(bkdataProxyconfigPtr->m_proxyProtocol);
    if (m_writer == NULL)
    {
        LOG_ERROR("failed to start proxy export by protocol(%s)", bkdataProxyconfigPtr->m_proxyProtocol.c_str());
        return false;
    }

    m_writer->SetConf(bkdataProxyconfigPtr);
    return m_writer->Start();
}

bool ProxyExporter::FilterByPlatId(DataCell *pDataCell)
{
    if (m_cfgSupportPlatIds.empty())
    {
        return true;
    }

    int platid = GetPlatNum(pDataCell->GetChannelID());
    std::map<int, int>::iterator it = m_cfgSupportPlatIds.find(platid);
    if (it != m_cfgSupportPlatIds.end())
    {
        LOG_DEBUG("find platid in dataflow config: %d", platid);
        return true;
    }

    return false;
}

void ProxyExporter::SetExporterName(const string &name)
{
    m_exporterName = name;
}

bool ProxyExporter::startWithDataFlow(ExporterConf *ptrExporterConf)
{
    DSProxyExporterConf *ptr_proxy_exporter_conf = ptrExporterConf->m_dsProxyConf;
    for (auto it : ptr_proxy_exporter_conf->m_platids)
    {
        m_cfgSupportPlatIds.insert(std::make_pair(it, it));
    }

    SetExporterName(ptrExporterConf->m_name);
    return startProxyExporter(ptr_proxy_exporter_conf);
}

void PackV2Msg(DataCell *pDataCell, const std::vector<std::string> &extensions, bool fill_channelid, GSEDataPackage &package)
{
    GSEDataPackageExtension extension;
    static std::atomic<unsigned long> cnt_id(0);
    std::string id = gse::tools::strings::ToString(cnt_id++); // TODO: 为tglog 实时数据添加的数据id头
    std::string timestamp = gse::tools::strings::ToString(pDataCell->GetCreationTimestamp());
    extension.ResetBufferPosition();
    if (fill_channelid)
    {
        extension.SetChannelID(pDataCell->GetChannelID());
    }
    else
    {
        extension.SetChannelID(0);
    }

    extension.SetTimestamp(pDataCell->GetArrivedTimestamp());
    extension.AddExtension(id.c_str(), id.size());
    extension.AddExtension(timestamp.c_str(), timestamp.size());
    std::size_t max_count = extensions.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        std::string value = extensions.at(idx);
        if (value.empty())
        {
            continue;
        }

        extension.AddExtension(value.c_str(), value.length());
    }

    package.ResetBufferPosition();
    package.SetType(DS_EXT_TYPE);
    package.SetVersion(DS_EXT_VERSION);
    package.SetMagicNumber(DS_MAGIC_NUM);
    package.SetBody(pDataCell->GetDataBuf(), pDataCell->GetDataBufLen(), &extension);
}

int PackV1Msg(DataCell *pDataCell, DynamicalProtocol &dynamic_protocol, bool fill_channelid, char *&buffer)
{
    DynamicalProtocolHead head;
    // msghead
    head.m_msgHead.msgType = REPORT_DYNAMICAL_PROTOCOL_TYPE;
    LOG_DEBUG("pack gse v1 msg, fillchannelid:%d", fill_channelid);
    if (fill_channelid)
    {
        LOG_DEBUG("Fill channelid , channelid:%d", pDataCell->GetChannelID());
        head.m_msgHead.lMsg.dataid = pDataCell->GetChannelID();
    }
    else
    {
        head.m_msgHead.lMsg.dataid = 0;
        LOG_DEBUG("Fill channelid , use default channelid:%d", 0);
    }

    head.m_msgHead.lMsg.utctime = time(NULL);
    head.m_msgHead.lMsg.resv[0] = 0;
    head.m_msgHead.lMsg.resv[1] = 0;

    // add nodeid
    //  set Dynamical protocol
    int ret = dynamic_protocol.setDynamicalProtocolHead(head, pDataCell->GetDataBuf(), pDataCell->GetDataBufLen());
    if (ret != 0)
    {
        LOG_ERROR("set dynamical protocol head failed");
        return 0;
    }
    int len = 0;
    dynamic_protocol.binaryData(buffer, len);
    return len;
}

ProxyTcpWriter::ProxyTcpWriter()
{
    m_nextTcpClient = 0;
    m_needSendByProtocolV1 = false;
    m_fillChannelid = true;
    m_heartbeat = true;
    m_connectionNumEachAddress = 1;
}

ProxyTcpWriter::~ProxyTcpWriter()
{
    Stop();
}

void ProxyTcpWriter::connectCallback(gse::net::TcpClient *tcpclient, const gse::net::TcpConnectionPtr connection)
{
    LOG_DEBUG("connection status callback, connection:%s, status:%d", connection->PeerAddress().ToIPPort().c_str(), connection->status());
    if (connection->IsDisconnected())
    {
        m_tcpClientLock.WLock();
        auto it = m_tcpClients.find(connection->GetNetID());
        if (it != m_tcpClients.end())
        {
            m_tcpClients.erase(connection->GetNetID());
        }
        m_tcpClientLock.UnLock();

        LOG_WARN("connection disconnected, peer host [%s:%d], netid:%s",
                 connection->PeerAddress().ToIPString().c_str(), connection->PeerAddress().ToPort(),
                 connection->GetNetID().c_str());
        gse::tools::system::Sleep(1);
        int ret = tcpclient->Connect();
        if (ret != 0)
        {
            LOG_ERROR("failed to reconnect remote server, ret:%d", ret);
        }
    }

    if (connection->IsConnected())
    {
        LOG_DEBUG("connect peer host [%s:%d] success, netid:%s",
                  connection->PeerAddress().ToIPString().c_str(),
                  connection->PeerAddress().ToPort(),
                  connection->GetNetID().c_str());

        m_tcpClientLock.WLock();

        m_tcpClients.insert(std::make_pair(connection->GetNetID(), tcpclient));
        m_tcpClientLock.UnLock();
    }
}

bool ProxyTcpWriter::startTCPClient()
{
    bool useSSL = false;
#ifdef _GSE_SSL_
    useSSL = true;
#define _GSE_DATA_EDTION_ _GSE_SSL_

#else
    if (!m_caFilePath.empty() &&
        !m_certPath.empty() &&
        !m_keyFilePath.empty() &&
        !m_passwdFilePath.empty())
    {
        useSSL = true;
#define _GSE_DATA_EDTION_ GSE_EDITION_EE
    }

#endif

    std::size_t max_count = m_addresses.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        Address *address = &m_addresses.at(idx);
        for (int client_idx = 0; client_idx < m_connectionNumEachAddress; ++client_idx)
        {
            gse::net::TcpClient *clientPtr = new gse::net::TcpClient(address->m_ip, address->m_port);
            if (useSSL)
            {
                if (!clientPtr->EnableSSL(m_caFilePath.c_str(), m_certPath.c_str(), m_keyFilePath.c_str(), m_passwdFilePath.c_str(), _GSE_DATA_EDTION_))
                {
                    LOG_ERROR("Failed to start ssl tcp client, ca:%s, certpath:%s, keypath:%s, passwd:%s",
                              m_caFilePath.c_str(), m_certPath.c_str(), m_keyFilePath.c_str(), m_passwdFilePath.c_str());
                    return false;
                }
            }

            clientPtr->SetConnectionCallback(std::bind(&ProxyTcpWriter::connectCallback, this, clientPtr, std::placeholders::_1));
            int ret = clientPtr->Connect();
        }
    }

    return true;
}
int ProxyTcpWriter::Start()
{
    return startTCPClient();
}

int ProxyTcpWriter::Stop()
{

    for (auto &it : m_tcpClients)
    {
        gse::net::TcpClient *ptrTcpClient = it.second;
        ptrTcpClient->Disconnect();
        delete ptrTcpClient;
        ptrTcpClient = NULL;
    }
    m_tcpClients.clear();
}

void ProxyTcpWriter::ReportOps(DataCell *pDataCell, gse::net::TcpClient *ptrTcpClient)
{
    std::string address = ptrTcpClient->GetServerIp() + "|" + gse::tools::strings::ToString(ptrTcpClient->GetServerPort());
    pDataCell->SetOutputType("tcp");
    pDataCell->SetOutputAddress(address);
    pDataCell->SetOutputTimestamp(::time(NULL));
    std::string exportortag = "proxy|tcp|" + address;
    pDataCell->SetExportorName(exportortag);
}

bool ProxyTcpWriter::WriteToByTCP(const char *ptrData, uint32_t dataLength, DataCell *pDataCell)
{
    m_tcpClientLock.RLock();
    std::size_t maxCount = m_tcpClients.size();
    m_tcpClientLock.UnLock();
    if (0 == maxCount)
    {
        LOG_WARN("there is not any valid connection used to send data");
        return false;
    }

    int retry_times = 0;
    int nextIdx = ::abs(static_cast<int>(gse::tools::atomic::AtomAddAfter(&m_nextTcpClient) % maxCount));
    m_tcpClientLock.RLock();
    auto iter = m_tcpClients.begin();
    std::advance(iter, nextIdx);

    if (m_tcpClients.end() == iter)
    {
        ++iter;
    }

    gse::net::TcpClient *ptrTcpClient = iter->second;
    m_tcpClientLock.UnLock();
    bool bSuccess = false;
    do
    {
        int ret = 0;
        ret = ptrTcpClient->Send(ptrData, dataLength);
        if (ret >= 0)
        {
            bSuccess = true;
            break;
        }
        else
        {
            std::string remote_server_ip = ptrTcpClient->GetServerIp();
            int port = ptrTcpClient->GetServerPort();
            LOG_ERROR("it is failed to send data to the remote server, will to reconnect the remote server(%s:%d), ret:%d",
                      remote_server_ip.c_str(), port, ret);
            retry_times++;
        }
    } while (retry_times < maxCount);

    if (!bSuccess)
    {
        LOG_ERROR("failed to send mesage, mesage len:%d, retry times:%d", dataLength, retry_times);

        return false;
    }

    LOG_DEBUG("send mesage to [%s:%d], mesage len[%d]", ptrTcpClient->GetServerIp().c_str(), ptrTcpClient->GetServerPort(), dataLength);

    ReportOps(pDataCell, ptrTcpClient);
    return true;
}

int ProxyTcpWriter::Write(DataCell *pDataCell)
{
    std::string type;
    int errcode = GSE_PROXY_OK;
    auto _ = gse::tools::defer::finally([&]() {
        OPMetric::AddSendMsgCounter(type, kOutputUDP, errcode, 1);
    });

    if (!m_needSendByProtocolV1)
    {
        type = "v2";
        GSEDataPackage gse_datapackage;
        PackV2Msg(pDataCell, m_extentions, m_fillChannelid, gse_datapackage);
        LOG_DEBUG("send data to next ds by proxy for the channel id (%d) the data length is (%d), the package data length is %d",
                  pDataCell->GetChannelID(), pDataCell->GetDataBufLen(), gse_datapackage.GetDataLength());

        if (!WriteToByTCP(gse_datapackage.GetDataPointer(), gse_datapackage.GetDataLength(), pDataCell))
        {
            LOG_ERROR("it is failed to send any data for the channel id (%d) to next ds by tcp", pDataCell->GetChannelID());
            pDataCell->SetErrorMsg("failed to send by tcp", OPS_ERROR_TRANSFER_TCP_ERROR);
            DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_LOST_STATE));
            errcode = GSE_PROXY_ERROR;
            return GSE_ERROR;
        }

        OPMetric::AddSendMsgBytesCounter(type, kOutputTCP, gse_datapackage.GetDataLength());
        DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_OUTPUT_STATE));
        return GSE_SUCCESS;
    }

    type = "v1";
    DynamicalProtocol dynamic_protocol;
    int32_t len = 0;
    char *ptrData = NULL;
    len = PackV1Msg(pDataCell, dynamic_protocol, m_fillChannelid, ptrData);
    LOG_DEBUG("send data to next ds by proxy for the channel id (%d) the data length is (%d), the package data is %d", pDataCell->GetChannelID(), pDataCell->GetDataBufLen(), len);
    if (!WriteToByTCP(ptrData, len, pDataCell))
    {
        LOG_ERROR("it is failed to send any data for the channel id (%d) to next ds by tcp ", pDataCell->GetChannelID());
        pDataCell->SetErrorMsg("failed to send by tcp", OPS_ERROR_TRANSFER_TCP_ERROR);
        DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_LOST_STATE));
        errcode = GSE_PROXY_ERROR;
        return GSE_ERROR;
    }

    OPMetric::AddSendMsgBytesCounter(type, kOutputTCP, len);
    DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_OUTPUT_STATE));
    return GSE_SUCCESS;
}

void ProxyTcpWriter::SetConf(DSProxyExporterConf *conf)
{
    m_fillChannelid = conf->m_fillChannelid;
    m_needSendByProtocolV1 = gse::tools::strings::CompareStringCase(conf->m_proxyVersion, PROXY_VERSION_1);
    m_connectionNumEachAddress = conf->m_connectionNumEachAddress;
    m_addresses = conf->m_addresses;
    m_certPath = conf->m_certPath;
    m_caFilePath = conf->m_caFilePath;
    m_keyFilePath = conf->m_keyfilePath;
    m_passwdFilePath = conf->m_passwdFilePath;
    m_heartbeat = conf->m_heartbeat;
    m_extentions = conf->m_extentions;

    LOG_DEBUG("set tcp proxy config, heartbeat:%d, fillchannelid:%d, protocol:%s, connection num:%d, host address count:%d, cert path:%s",
              m_heartbeat, conf->m_fillChannelid, conf->m_proxyVersion.c_str(), m_connectionNumEachAddress, m_addresses.size(), m_certPath.c_str());
    return;
}

ProxyUdpWriter::ProxyUdpWriter()
    : m_needSendByProtocolV1(false),
      m_nextUdpClient(0),
      m_connectionNumEachAddress(1),
      m_fillChannelid(true)
{
}

ProxyUdpWriter::~ProxyUdpWriter()
{
}

bool ProxyUdpWriter::startUDPClient()
{
    std::size_t max_count = m_addresses.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        Address address = m_addresses.at(idx);
        for (int client_idx = 0; client_idx < m_connectionNumEachAddress; ++client_idx)
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

int ProxyUdpWriter::Start()
{
    return startUDPClient();
}

int ProxyUdpWriter::Stop()
{
    int max_count = m_udpClientVector.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        gse::net::udp::UDPClient *ptr_udp_client = m_udpClientVector.at(idx);
        delete ptr_udp_client;
    }
    m_udpClientVector.clear();
}

bool ProxyUdpWriter::writeToByUDP(const char *ptrData, uint32_t dataLength, DataCell *pDataCell)
{
    std::size_t udpclient_max_count = m_udpClientVector.size();
    if (0 == udpclient_max_count)
    {
        LOG_WARN("there is not any valid connection used to send data");
        return false;
    }

    int next_idx = abs(static_cast<int>(gse::tools::atomic::AtomAddAfter(&m_nextUdpClient) % udpclient_max_count));
    gse::net::udp::UDPClient *ptr_udp_client = m_udpClientVector.at(next_idx);

    if (ptr_udp_client->Send(ptrData, dataLength) < 0)
    {
        LOG_ERROR("it is failed to send data to the remote server, the data length is %d", dataLength);
        return false;
    }

    std::string address = ptr_udp_client->GetServerIp() + "|" + gse::tools::strings::ToString(ptr_udp_client->GetPort());

    std::string exportortag = "proxy|udp|" + address;
    pDataCell->SetExportorName(exportortag);

    pDataCell->SetOutputType("udp");
    pDataCell->SetOutputAddress(address);
    pDataCell->SetOutputTimestamp(::time(NULL));
    return true;
}

int ProxyUdpWriter::Write(DataCell *pDataCell)
{
    int errcode = GSE_PROXY_OK;

    std::string type;
    auto _ = gse::tools::defer::finally([&]() {
        OPMetric::AddSendMsgCounter(type, kOutputUDP, errcode, 1);
    });

    if (!m_needSendByProtocolV1)
    {
        type = "v2";
        GSEDataPackage gse_datapackage;
        PackV2Msg(pDataCell, m_extentions, m_fillChannelid, gse_datapackage);

        LOG_DEBUG("send data to next ds by proxy for the channel id (%d) the data length is (%d), the package data length is %d",
                  pDataCell->GetChannelID(), pDataCell->GetDataBufLen(), gse_datapackage.GetDataLength());

        if (!writeToByUDP(gse_datapackage.GetDataPointer(), gse_datapackage.GetDataLength(), pDataCell))
        {
            LOG_ERROR("it is failed to send any data for the channel id (%d) to next ds by udp ", pDataCell->GetChannelID());
            pDataCell->SetErrorMsg("failed to send by udp", OPS_ERROR_TRANSFER_TCP_ERROR);
            DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_LOST_STATE));
            errcode = GSE_PROXY_ERROR;
            return GSE_ERROR;
        }

        OPMetric::AddSendMsgBytesCounter(type, kOutputUDP, gse_datapackage.GetDataLength());
        DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_OUTPUT_STATE));
        return GSE_SUCCESS;
    }

    type = "v1";
    DynamicalProtocol dynamical_protocol;

    int32_t len = 0;
    char *ptrData = NULL;
    len = PackV1Msg(pDataCell, dynamical_protocol, m_fillChannelid, ptrData);

    LOG_DEBUG("send data to next ds by proxy for the channel id (%d) the data length is (%d), the package data is %d", pDataCell->GetChannelID(), pDataCell->GetDataBufLen(), len);
    if (!writeToByUDP(ptrData, len, pDataCell))
    {
        LOG_ERROR("it is failed to send any data for the channel id (%d) to next ds by udp ", pDataCell->GetChannelID());
        pDataCell->SetErrorMsg("failed to send by udp", OPS_ERROR_TRANSFER_UDP_ERROR);
        DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_LOST_STATE));
        errcode = GSE_PROXY_ERROR;
        return GSE_ERROR;
    }

    OPMetric::AddSendMsgBytesCounter(type, kOutputUDP, len);
    DataServer::GetOpsReportClient()->PutOpsData(pDataCell->ToOPS(EN_OUTPUT_STATE));
    return GSE_SUCCESS;
}

void ProxyUdpWriter::SetConf(DSProxyExporterConf *conf)
{
    m_fillChannelid = conf->m_fillChannelid;
    m_needSendByProtocolV1 = gse::tools::strings::CompareStringCase(conf->m_proxyVersion, PROXY_VERSION_1);
    m_connectionNumEachAddress = conf->m_connectionNumEachAddress;
    m_addresses = conf->m_addresses;
    m_extentions = conf->m_extentions;
}

HTTPClientWorker::HTTPClientWorker()
    : m_evbase(NULL), m_msgHandleEvent(NULL), m_sslCtx(NULL), m_exit(false), m_isKeepAlive(false),
      m_https(false), m_isThirdPartyCert(false)
{
    m_msgNotifyPipe[0] = -1;
    m_msgNotifyPipe[1] = -1;
    m_msgQuene = new moodycamel::ConcurrentQueue<NetSendDataItem *>();
}

HTTPClientWorker::~HTTPClientWorker()
{
    stop();
}

bool HTTPClientWorker::LoadClientCert()
{
    std::string cert_file_path = m_certPath + "/gse_agent.crt";
    std::string cert_key_file_path = m_certPath + "/gse_agent.key";
    std::string certPasswd;

    if (m_isThirdPartyCert)
    {
        cert_file_path = m_thirdPartyCertFile;
        cert_key_file_path = m_thirdPartyKeyFile;
        certPasswd = m_thirdPartyCertPasswd;
    }
    else
    {
        std::string keyfilename = "gse_agent.key";
        gse::safe::SecureConfig cfg;
#ifdef _GSE_SSL_
        int ret = cfg.GetCertKeyFilePassword(m_passwdFilePath, keyfilename, _GSE_SSL_, certPasswd);
#else
        int ret = cfg.GetCertKeyFilePassword(m_passwdFilePath, keyfilename, GSE_EDITION_EE, certPasswd);
#endif
        if (ret != GSE_SUCCESS)
        {
            LOG_ERROR("failed to get cert(%s) passwd", keyfilename.c_str());
            return false;
        }
    }

#if OPENSSL_VERSION_NUMBER < 0x10100000L
    m_sslCtx = SSL_CTX_new(SSLv23_client_method());
#else
    m_sslCtx = SSL_CTX_new(TLS_client_method());
#endif
    if (NULL == m_sslCtx)
    {
        int err = ERR_get_error();
        LOG_ERROR("failed to new client sslctx, error:%s", ERR_reason_error_string(err));
        return false;
    }

    if (!certPasswd.empty())
    {
        SSL_CTX_set_default_passwd_cb_userdata(m_sslCtx, (void *)certPasswd.c_str());
    }

    if (!cert_key_file_path.empty())
    {
        /* client private key file defined, so use it */
        int res = SSL_CTX_use_PrivateKey_file(m_sslCtx, cert_key_file_path.c_str(), SSL_FILETYPE_PEM);

        if (res == 0)
        {
            int err = ERR_get_error();
            LOG_ERROR("failed to call SSL_CTX_use_PrivateKey_file,error: %s", ERR_reason_error_string(err));
            return false;
        }
    }

    if (!cert_file_path.empty())
    {
        /* client cert key file defined, use it */
        int res = SSL_CTX_use_certificate_file(m_sslCtx, cert_file_path.c_str(), SSL_FILETYPE_PEM);

        if (res == 0)
        {
            int err = ERR_get_error();
            LOG_ERROR("failed to call SSL_CTX_use_certificate_file,error: %s", ERR_reason_error_string(err));
            return false;
        }
    }

    return true;
}

bool HTTPClientWorker::SSLInit()
{
    if (m_certPath.empty() && !m_isThirdPartyCert)
    {
        m_https = false;
        return false;
    }

    m_https = true;
    LOG_DEBUG("http client init, ssl:%d", m_isThirdPartyCert);
    return LoadClientCert();
}

int HTTPClientWorker::PrintHeader(evhtp_header_t *header, void *arg)
{
    LOG_DEBUG("%s: %s", header->key, header->val);
    return 0;
}

void HTTPClientWorker::HandleReponse(evhtp_request_t *request, void *arg)
{
    //    auto _ = gse::tools::defer::finally([&request]() {
    //        size_t bytes = evbuffer_get_length(request->buffer_in);
    //        LOG_DEBUG("receive reply bytes: %zu", bytes);
    //        evhtp_unset_all_hooks(&request->hooks);
    //        evbuffer_drain(request->buffer_out, -1);
    //        evbuffer_drain(request->buffer_in, -1);
    //        if (request->uri)
    //        {
    //            free(request->uri);
    //            request->uri = NULL;
    //        }
    //        evhtp_request_free(request);
    //        LOG_DEBUG("free request:%p", request);
    //    });

    if (request->status != EVHTP_RES_200)
    {
        LOG_ERROR("recv http response error, code: %d", request->status);

        // evhtp_headers_for_each(request->headers_in, HTTPClientWorker::PrintHeader, NULL);

        return;
    }

    LOG_DEBUG("recv http response success");
}

bool HTTPClientWorker::SelectNextHost(std::string &ip, uint16_t &port)
{
    static std::atomic<unsigned long> index(0);
    LOG_DEBUG("available http server count:%d", m_addresses.size());
    if (m_addresses.size() == 0)
    {
        return false;
    }
    if (m_addresses.size() > 1)
    {
        index = (index++) % m_addresses.size();
        ip = m_addresses[index].m_ip;
        port = m_addresses[index].m_port;
        return true;
    }
    else
    {
        ip = m_addresses[0].m_ip;
        port = m_addresses[0].m_port;
        return true;
    }
    return false;
}
// typedef void (* evhtp_hook_err_cb)(evhtp_request_t * req, evhtp_error_flags errtype, void * arg);
void HTTPClientWorker::PosHttpRequestErrorCb(evhtp_request_t *req, evhtp_error_flags errtype, void *arg)
{
    LOG_ERROR("failed to send http request, error:%d", errtype);
}
// typedef evhtp_res (* evhtp_hook_conn_err_cb)(evhtp_connection_t * connection, evhtp_error_flags errtype, void * arg);
void HTTPClientWorker::PosHttpRequestConnectErrorCb(evhtp_connection_t *connection, evhtp_error_flags errtype, void *arg)
{
}

evhtp_connection_t *HTTPClientWorker::GetConnection(const std::string &host)
{
    auto it = m_httpConns.find(host);
    if (it != m_httpConns.end())
    {
        return it->second;
    }

    return NULL;
}

void HTTPClientWorker::CloseConnection()
{
    for (auto &it : m_httpConns)
    {
        evhtp_connection_free(it.second);
    }

    m_httpConns.clear();
}

void HTTPClientWorker::PushConnection(const std::string &host, evhtp_connection_t *conn)
{
    auto it = m_httpConns.find(host);
    if (it == m_httpConns.end())
    {
        m_httpConns.insert(std::pair<std::string, evhtp_connection_t *>(host, conn));
    }
}

void HTTPClientWorker::RemoveInvalidConnection(evhtp_connection_t *connection)
{
    for (auto it = m_httpConns.begin(); it != m_httpConns.end();)
    {
        if (it->second == connection)
        {
            m_httpConns.erase(it);
            LOG_WARN("connection:%p freed, remove from connection pool", connection);
            break;
        }
        else
        {
            it++;
        }
    }
}

evhtp_res HTTPClientWorker::HTTPConnectionFreeCallBack(evhtp_connection_t *connection, void *arg)
{
    HTTPClientWorker *self = (HTTPClientWorker *)arg;
    self->RemoveInvalidConnection(connection);
}

evhtp_res HTTPClientWorker::OnResponseDataCallBack(
    evhtp_request_t *request, void *arg)
{

    size_t n_bytes = evbuffer_get_length(request->buffer_in);
    LOG_DEBUG("Bytes arrived: %zu", n_bytes);
    evbuffer_drain(request->buffer_out, -1);
    return EVHTP_RES_OK;
}

bool HTTPClientWorker::PostMsg(const char *msg, const int len)
{
    evhtp_connection_t *conn = NULL;
    evhtp_request_t *request = NULL;
    std::string host;
    uint16_t port = 0;
    int errcode = 0;

    auto _ = gse::tools::defer::finally([&]() {
        OPMetric::AddSendMsgCounter(kOutputHTTP, kOutputHTTP, errcode, 1);
    });

    if (!SelectNextHost(host, port))
    {
        LOG_ERROR("failed to post http messeage, no available host");
        errcode = GSE_PROXY_ERROR;
        return false;
    }

    LOG_DEBUG("will post http msg to %s:%d, URI:%s", host.c_str(), port, m_requestURI.c_str());
    std::string strProto;
    if (conn == NULL)
    {
        if (m_https)
        {
            conn = evhtp_connection_ssl_new(m_evbase,
                                            host.c_str(), port, m_sslCtx);
            strProto = "https";
        }
        else
        {
            conn = evhtp_connection_new(m_evbase, host.c_str(), port);
            strProto = "http";
        }

        if (conn == NULL)
        {
            LOG_ERROR("failed to new evhtp connection, https:%d, host:%s:%d", m_https, host.c_str(), port);
            errcode = GSE_PROXY_ERROR;
            return false;
        }
        // PushConnection(httpHost, conn);

        //        evhtp_connection_set_hook(conn,
        //                                  evhtp_hook_on_connection_fini,
        //                                  (evhtp_hook)HTTPClientWorker::HTTPConnectionFreeCallBack, (void *)this);
    }

    request = evhtp_request_new(HTTPClientWorker::HandleReponse, m_evbase);
    if (request == NULL)
    {
        LOG_ERROR("failed to new evhtp request");
        errcode = GSE_PROXY_ERROR;
        return false;
    }

    evhtp_headers_add_header(request->headers_out,
                             evhtp_header_new("Host", "localhost", 0, 0));
    evhtp_headers_add_header(request->headers_out,
                             evhtp_header_new("User-Agent", "gse-data", 0, 0));

    evhtp_headers_add_header(request->headers_out,
                             evhtp_header_new("Connection", "close", 0, 0));

    //    evhtp_request_set_hook(request, evhtp_hook_on_chunks_complete,
    //                           (evhtp_hook)HTTPClientWorker::OnResponseDataCallBack, (void *)this);
    // evhtp_request_set_keepalive(request, 1);
    std::string str_len = gse::tools::strings::ToString(len);
    evhtp_headers_add_header(request->headers_out,
                             evhtp_header_new("Content-Length", str_len.c_str(), 0, 0));

    evhtp_make_request(conn, request, htp_method_POST, m_requestURI.c_str());
    int ret = evbuffer_add(bufferevent_get_output(conn->bev), (char *)msg, len);
    if (ret != 0)
    {
        LOG_ERROR("failed to http request to remote host:%s:%d, request uri:%s", host.c_str(), port, m_requestURI.c_str());
        errcode = GSE_PROXY_ERROR;
        return false;
    }

    OPMetric::AddSendMsgBytesCounter(kOutputHTTP, kOutputHTTP, str_len.size());

    //    evhtp_safe_free(request, evhtp_request_free);
    //    evhtp_safe_free(conn, evhtp_connection_free);
    return true;
}

bool HTTPClientWorker::ProcessQueneMsg(evutil_socket_t fd, short which)
{
    char buff[256];
    int count = 0;
    int max_read_times = 0;
    while ((count = read(fd, buff, sizeof(buff))) == sizeof(buff))
    {
        max_read_times++;
        if (max_read_times > 3)
        {
            break;
        }
    }

    NetSendDataItem *item;
    while (m_msgQuene->try_dequeue(item))
    {
        if (!PostMsg(item->m_buf, item->m_bufLen))
        {
            LOG_ERROR("failed to post msg, message length:%d", item->m_bufLen);
        }
        delete item;
    }
}

void HTTPClientWorker::MsgEventHandler(evutil_socket_t fd, short which, void *v)
{
    HTTPClientWorker *self = static_cast<HTTPClientWorker *>(v);
    self->ProcessQueneMsg(fd, which);
}

void HTTPClientWorker::StopBaseLoop()
{
    if (m_evbase != NULL)
    {
        event_base_loopexit(m_evbase, NULL);
    }
}
void HTTPClientWorker::FreeEvent()
{
    if (m_msgHandleEvent != NULL)
    {
        event_free(m_msgHandleEvent);
        m_msgHandleEvent = NULL;
    }
    if (m_msgNotifyPipe[0] > 0)
    {
        gse::tools::pipe::ClosePipe(m_msgNotifyPipe[0]);
        m_msgNotifyPipe[0] = -1;
    }

    if (m_msgNotifyPipe[1] > 0)
    {
        gse::tools::pipe::ClosePipe(m_msgNotifyPipe[1]);
        m_msgNotifyPipe[1] = -1;
    }

    if (m_sslCtx != NULL)
    {
        SSL_CTX_free(m_sslCtx);
        m_sslCtx = NULL;
    }
}

void HTTPClientWorker::CleanMsgQuene()
{
    NetSendDataItem *msg = NULL;
    if (m_msgQuene == NULL)
    {
        return;
    }

    while (m_msgQuene->try_dequeue(msg))
    {
        if (msg != NULL)
        {
            delete msg;
        }
    }

    delete m_msgQuene;
    m_msgQuene = NULL;
}

bool HTTPClientWorker::InitEvent()
{
    m_evbase = event_base_new();
    if (m_evbase == NULL)
    {
        int tmperrno = gse_socket_errno;
        LOG_ERROR("failed to new event base, errno=[%d], errstr=[%s]", tmperrno, SAFE_CSTR(gse::tools::error::ErrnoToStr(tmperrno).c_str()));
        return false;
    }

    if (gse::tools::pipe::CreatePipe(m_msgNotifyPipe) < 0)
    {
        int tmperrno = gse_socket_errno;

        LOG_ERROR("failed to create msg notify pipe, errno=[%d], errstr=[%s]", tmperrno, SAFE_CSTR(gse::tools::error::ErrnoToStr(tmperrno).c_str()));
        return false;
    }

    gse::tools::net::SocketSetnonblocking(m_msgNotifyPipe[0]);
    gse::tools::net::SocketSetnonblocking(m_msgNotifyPipe[1]);

    m_msgHandleEvent = event_new(m_evbase, m_msgNotifyPipe[0], EV_PERSIST | EV_READ, HTTPClientWorker::MsgEventHandler, this);
    if (NULL == m_msgHandleEvent)
    {
        int temperrno = gse_errno;
        LOG_ERROR("failed to event_new msg handle event, [%d:%s]", temperrno, gse::tools::error::ErrnoToStr(temperrno).c_str());
        return false;
    }

    event_add(m_msgHandleEvent, NULL);
    LOG_DEBUG("Http worker init event success");
    return true;
}
int HTTPClientWorker::Init()
{
    return GSE_SUCCESS;
}

int HTTPClientWorker::ThreadFun()
{
    int rc = -1;
    LOG_DEBUG("Http client worker thread run...");
    while (rc < 0 && !m_exit)
    {
        SSLInit();
        InitEvent();
        rc = event_base_loop(m_evbase, 0);
        LOG_ERROR("Http client worker event base loop abnomal break. ret=[%d]", rc);
        FreeEvent();
        event_base_free(m_evbase);
        m_evbase = NULL;
    }
}

void HTTPClientWorker::NotifyRecvMsgEvent()
{
    int retry = 0;
    char notifyC = '1';

    while (retry < 3)
    {
        if (gse::tools::pipe::WritePipe(m_msgNotifyPipe[1], &notifyC, sizeof(notifyC)) < 0)
        {
            retry++;
            int tmp_errno = gse_errno;
            LOG_WARN("fail to send data notify. retry[%d], errno[%d], errstr[%s]", retry, tmp_errno, gse::tools::error::ErrnoToStr(tmp_errno).c_str());
        }
        else
        {
            break;
        }
    }
    return;
}

void HTTPClientWorker::SetConf(const DSProxyExporterConf &conf)
{
    m_requestURI = conf.m_httpURI;
    m_certPath = conf.m_certPath;
    // m_passwdFilePath = conf;

    m_addresses = conf.m_addresses;
    m_strAddress.clear();
    for (int i = 0; i < m_addresses.size(); i++)
    {
        m_strAddress.append(m_addresses[i].m_ip);
        m_strAddress.append(":");
        m_strAddress.append(gse::tools::strings::ToString(m_addresses[i].m_port));
        if (i < m_addresses.size())
        {
            m_strAddress.append(",");
        }
    }
    m_thirdPartyCertPasswd = conf.m_thirdPartyCertPasswd;
    m_thirdPartyCertFile = conf.m_thirdPartyCertFile;
    m_thirdPartyKeyFile = conf.m_thirdPartyKeyFile;
    m_isThirdPartyCert = conf.m_isThirdPartyCert;
    m_passwdFilePath = conf.m_passwdFilePath;
}

int HTTPClientWorker::SendMsg(DataCell *ptr_datacell)
{
    int max_retry_times = 0;
    if (m_exit)
    {
        LOG_ERROR("http client worker exit");
        return GSE_ERROR;
    }

    if (m_msgQuene->size_approx() >= MAX_QUEUE_SIZE)
    {
        LOG_ERROR("send http msg failed, the queue exceeds the max limit(%d), queue msg count:%d", MAX_QUEUE_SIZE, m_msgQuene->size_approx());

        return GSE_ERROR;
    }

    NetSendDataItem *msg = new NetSendDataItem((char *)ptr_datacell->GetDataBuf(), ptr_datacell->GetDataBufLen());

    while (!m_msgQuene->enqueue(msg))
    {
        max_retry_times++;
        if (max_retry_times > 5)
        {
            delete msg;
            return GSE_ERROR;
        }
    }

    std::string exportortag = "proxy|http|" + m_strAddress;
    ptr_datacell->SetExportorName(exportortag);
    std::string str_proto = "http";
    if (m_https)
    {
        str_proto = "https";
    }
    ptr_datacell->SetOutputType(str_proto);
    ptr_datacell->SetOutputAddress(m_strAddress);
    ptr_datacell->SetOutputTimestamp(::time(NULL));
    DataServer::GetOpsReportClient()->PutOpsData(ptr_datacell->ToOPS(EN_OUTPUT_STATE));
    NotifyRecvMsgEvent();
    return GSE_SUCCESS;
}

void HTTPClientWorker::stop()
{
    LOG_DEBUG("Stop http client worker");
    m_exit = true;
    CleanMsgQuene();
    CloseConnection();
    FreeEvent();
    StopBaseLoop();
}

ProxyHTTPWriter::ProxyHTTPWriter()
    : m_maxWorkerNum(32)
{
}

ProxyHTTPWriter::~ProxyHTTPWriter()
{
}

int ProxyHTTPWriter::Start()
{
    for (int i = 0; i < m_maxWorkerNum; ++i)
    {
        std::shared_ptr<HTTPClientWorker> worker = std::make_shared<HTTPClientWorker>();
        worker->SetConf(m_conf);
        worker->Start();
        m_clientWorkers.push_back(worker);
        LOG_DEBUG("Start http writer:%d", i);
    }
}

int ProxyHTTPWriter::Stop()
{
    for (int i = 0; i < m_maxWorkerNum; ++i)
    {
        m_clientWorkers[i]->stop();
    }

    m_clientWorkers.clear();
    m_clientWorkers.shrink_to_fit();
    return GSE_SUCCESS;
}

int ProxyHTTPWriter::Write(DataCell *pDataCell)
{
    static std::atomic<int> idx(0);
    std::atomic_fetch_add(&idx, 1);
    int i = idx % m_clientWorkers.size();
    int ret = m_clientWorkers[i]->SendMsg(pDataCell);
    return ret;
}

void ProxyHTTPWriter::SetConf(DSProxyExporterConf *conf)
{
    m_conf = *conf;
}

} // namespace data
} // namespace gse
