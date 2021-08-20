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

#ifndef _GSE_DATA_PROXY_EXPORTER_H_
#define _GSE_DATA_PROXY_EXPORTER_H_

#include <string>
#include <vector>
#include "exporter/exporter.h"
#include "datacell.h"
//#include "net/tcp_client.h"
#include "net/udp/udp_client.h"
//#include "codec/gsedata_pkg/gsedata_package.h"
//#include "codec/gsedata_pkg/gsedata_extension.h"

#include "protocol.h"
#include "safe/lock.h"
#include "safe/ssl.h"
#include "net/tcp_client.h"
namespace gse { 
namespace dataserver {

class ProxyTcpClient
{
     unique_ptr<gse::net::TcpClient> m_tcpClient;
};

class ProxyExporter : public Exporter
{
public:
    typedef std::vector<gse::net::TcpClient*> TcpClientPoolVector;
    typedef std::vector<gse::net::udp::UDPClient*> UdpClientPoolVector;

public:
    ProxyExporter();
    virtual ~ProxyExporter();

public:
    int Start();
    int Stop();
    int Write(DataCell *pDataCell);

private:
    bool writeToByTCP(const char* ptrData, uint32_t dataLength, DataCell *pDataCell);
    bool writeToByUDP(const char* ptrData, uint32_t dataLength, DataCell *pDataCell);
    bool startTCPClient(DSProxyExporterConf* ptrProxyExporterConf);
    bool startUDPClient(DSProxyExporterConf* ptrProxyExporterConf);
    void connectCallback(const gse::net::TcpConnectionPtr &connection);

private:
    bool startWithChannelID(ChannelIDStorage* ptrChannelIDConfig);
    bool startWithDataFlow(ExporterConf* ptrExporterConf);
    bool startWithDataID(StorageConfigType* ptrStorageConfig);

private:
    bool                 m_needSendByProtocolV1;
    bool                 m_needSendByProtocolTCP;
    bool                 m_fillChannelid;

    int                  m_nextTcpClient;
    int                  m_nextUdpClient;
    TcpClientPoolVector  m_tcpClientVector;
    UdpClientPoolVector  m_udpClientVector;
    //GseSSL* m_ptrSSL;
    gse::safe::SSLContext *m_ptrSSL;
    gse::safe::RWLock           m_rwlock;
    unique_ptr<gse::net::TcpClient> m_tcpClient;

private:
    std::vector<std::string> m_extensions;

};

}
}
#endif
