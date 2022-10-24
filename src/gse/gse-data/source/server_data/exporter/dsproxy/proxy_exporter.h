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

#ifndef _GSE_PROXY_EXPORTER_H_
#define _GSE_PROXY_EXPORTER_H_

#include <string>
#include <vector>

#include <evhtp.h>

#include "datacell.h"
#include "exporter/exporter.h"
#include "net/udp/udp_client.h"

#include "net/tcp_client.h"
#include "protocol.h"
#include "safe/lock.h"
#include "safe/ssl.h"

#include "datastruct/safe_map.h"

namespace gse {
namespace data {

#ifndef _GSE_NET_DATAITEM_H_
#define _GSE_NET_DATAITEM_H_

typedef struct NetSendDataItem_
{
    char *m_buf;
    int m_bufLen;

    NetSendDataItem_(char *buf, int bufLen)
    {
        m_buf = new char[bufLen];
        memcpy(m_buf, buf, bufLen);
        m_bufLen = bufLen;
    }

    ~NetSendDataItem_()
    {
        if (m_buf != NULL)
        {
            delete[] m_buf;
            m_buf = NULL;
        }
    }
} NetSendDataItem;

#endif

/* error code defines*/
#define GSE_PROXY_OK (0)
#define GSE_PROXY_ERROR (-1)

class ProxyWriter
{
public:
    ProxyWriter();
    virtual ~ProxyWriter();

public:
    virtual int Start() = 0;
    virtual int Stop() = 0;
    virtual int Write(DataCell *pDataCell) = 0;
    virtual void SetConf(DSProxyExporterConf *conf) = 0;
};

class ProxyTcpWriter : public ProxyWriter
{
public:
    typedef std::vector<gse::net::TcpClient *> TcpClientPoolVector;

public:
    ProxyTcpWriter();
    virtual ~ProxyTcpWriter();

public:
    virtual int Start() override;
    virtual int Stop() override;
    virtual int Write(DataCell *pDataCell) override;
    virtual void SetConf(DSProxyExporterConf *conf) override;

private:
    bool startTCPClient();
    bool WriteToByTCP(const char *ptrData, uint32_t dataLength, DataCell *pDataCell);
    void connectCallback(net::TcpClient *tcpclient, const net::TcpConnectionPtr connection);

    void ReportOps(DataCell *pDataCell, gse::net::TcpClient *ptrTcpClient);

private:
    bool m_needSendByProtocolV1;
    int m_nextTcpClient;

    std::vector<Address> m_addresses;
    std::string m_certPath;
    std::string m_passwdFilePath;
    std::string m_keyFilePath;
    std::string m_caFilePath;
    std::string m_proxyVersion;
    int m_connectionNumEachAddress;
    bool m_heartbeat;
    bool m_fillChannelid;

    gse::safe::RWLock m_tcpClientLock;
    std::map<std::string, gse::net::TcpClient *> m_tcpClients;
    std::vector<std::string> m_extentions;
};

class ProxyUdpWriter : public ProxyWriter
{
public:
    typedef std::vector<gse::net::udp::UDPClient *> UdpClientPoolVector;

public:
    ProxyUdpWriter();
    virtual ~ProxyUdpWriter();

public:
    virtual int Start() override;
    virtual int Stop() override;
    virtual int Write(DataCell *pDataCell) override;
    virtual void SetConf(DSProxyExporterConf *conf) override;

    bool writeToByUDP(const char *ptrData, uint32_t dataLength, DataCell *pDataCell);
    bool startUDPClient();

private:
    bool m_needSendByProtocolV1;
    int m_nextUdpClient;
    int m_connectionNumEachAddress;

    std::vector<Address> m_addresses;
    UdpClientPoolVector m_udpClientVector;
    bool m_fillChannelid;

    std::vector<std::string> m_extentions;
};

class HTTPClientWorker : public gse::tools::thread::Thread
{
public:
    HTTPClientWorker();
    ~HTTPClientWorker();

    int SendMsg(DataCell *ptr_datacell);
    void SetConf(const DSProxyExporterConf &conf);
    int Init();
    int ThreadFun();
    void stop();

    static void MsgEventHandler(evutil_socket_t fd, short which, void *v);

private:
    bool InitEvent();
    bool SSLInit();
    bool PostMsg(const char *msg, const int len);
    bool ProcessQueneMsg(evutil_socket_t fd, short which);
    void CleanMsgQuene();
    void FreeEvent();
    void NotifyRecvMsgEvent();
    void StopBaseLoop();

    evhtp_connection_t *GetConnection(const std::string &host);
    void CloseConnection();
    void PushConnection(const std::string &host, evhtp_connection_t *conn);
    bool LoadClientCert();
    bool SelectNextHost(std::string &ip, uint16_t &port);

    void RemoveInvalidConnection(evhtp_connection_t *connection);

    static int PrintHeader(evhtp_header_t *header, void *arg);
    static void HandleReponse(evhtp_request_t *request, void *arg);

    static void PosHttpRequestConnectErrorCb(evhtp_connection_t *connection, evhtp_error_flags errtype, void *arg);
    static void PosHttpRequestErrorCb(evhtp_request_t *req, evhtp_error_flags errtype, void *arg);

    static evhtp_res HTTPConnectionFreeCallBack(evhtp_connection_t *c, void *arg);
    static evhtp_res OnResponseDataCallBack(evhtp_request_t *request, void *arg);

private:
    bool m_https;
    struct event_base *m_evbase;
    struct event *m_msgHandleEvent;
    evhtp_ssl_ctx_t *m_sslCtx;
    int m_msgNotifyPipe[2];
    bool m_exit;
    moodycamel::ConcurrentQueue<NetSendDataItem *> *m_msgQuene;
    // config
    std::map<std::string, evhtp_connection_t *> m_httpConns;

    std::string m_certPath;
    std::string m_certKey;
    std::string m_passwdFilePath;
    std::vector<Address> m_addresses;
    std::string m_thirdPartyCertPasswd;
    std::string m_thirdPartyCertFile;
    std::string m_thirdPartyKeyFile;
    std::string m_strAddress;
    bool m_isThirdPartyCert;
    bool m_isKeepAlive;
    std::string m_requestURI;
};

class ProxyHTTPWriter : public ProxyWriter
{
public:
    ProxyHTTPWriter();
    virtual ~ProxyHTTPWriter();

public:
    virtual int Start() override;
    virtual int Stop() override;
    virtual int Write(DataCell *pDataCell) override;
    virtual void SetConf(DSProxyExporterConf *conf) override;

private:
    DSProxyExporterConf m_conf;
    std::vector<std::shared_ptr<HTTPClientWorker>> m_clientWorkers;
    int m_maxWorkerNum;
};

class ProxyExporter : public Exporter
{
public:
public:
    ProxyExporter();
    virtual ~ProxyExporter();

public:
    int Start();
    int Stop();
    int Write(DataCell *pDataCell);

public:
    void SetExporterName(const std::string &name);

private:
    // bool startWithChannelID(ChannelIDStorage* ptrChannelIDConfig);
    bool startWithDataFlow(ExporterConf *ptrExporterConf);

    bool startWithChannelID(ChannelIdExporterConfig *ptrChannelIDConfig);
    bool startProxyExporter(DSProxyExporterConf *ptr_proxy_exporter_conf);

    ProxyWriter *CreateWriterByProtocol(const std::string &protocol);
    bool FilterByPlatId(DataCell *pDataCell);

private:
    std::string m_protocol;
    ProxyWriter *m_writer;
    std::map<int, int> m_cfgSupportPlatIds;

    std::string m_exporterName;

private:
    // TODO: 需要优化，当前满足tglog 实施推送数据增加静态扩展信息配置使用
    std::vector<std::string> m_extensions;
};

} // namespace data
} // namespace gse
#endif
