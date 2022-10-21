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

#ifndef _GSE_DATA_TCPRECEIVER_H_
#define _GSE_DATA_TCPRECEIVER_H_

#include "net/tcp_server.h"
#include "receiver.h"
#include "tools/thread.h"

#include "loadbalance_schedule.h"
#include "net/message.h"
#include "net/message_factory.hpp"
#include "net/msg_processor.h"
#include "net/tcp_connection.h"

namespace gse {
namespace data {

class TcpReceiver : public Receiver
{
public:
    TcpReceiver();
    virtual ~TcpReceiver();

public:
    int Start();
    int Stop();
    void Join();

    void AddMsgProcessor(int msgType, std::shared_ptr<gse::net::MsgProcessor> processor);
    void SetSheduler(std::shared_ptr<LoadBalanceSchedule> schedule);

private:
    // implement virtual function
    void OnConnectionEvent(const std::shared_ptr<net::TcpConnection> conn);
    void MsgHandler(const std::shared_ptr<net::TcpConnection> conn, const std::shared_ptr<net::Message> msg);
    uint32_t getMsgHeadLen();
    void InitProto();

    int StartMigrationSerivce();
    int MigrationConnections();
    void RestoreConnectionCallback(evutil_socket_t fd,
                                   char *evInputBuffer, size_t evInputBufferSize,
                                   char *evOutputBuffer, size_t evOuputBufferSize,
                                   const gse::net::InetAddr &local, const gse::net::InetAddr &remote);

    void MigrationTcpListennerFd();
    int GetListenerFd();
    void StartTcpServer();

    void RegisterMsgProcessor();

protected:
    bool loadServerCert(const std::string &certRootPath);

private:
    std::shared_ptr<gse::net::TcpServer> m_tcpServer;
    std::unique_ptr<std::thread> m_tcpServerThread;
    std::unique_ptr<gse::net::MigrationClient> m_migrationClient;
    std::unique_ptr<gse::net::MigrationServer> m_migrationServer;

    std::shared_ptr<LoadBalanceSchedule> m_scheduler;
    std::map<int, std::shared_ptr<gse::net::MsgProcessor>> m_processors;
    gse::net::EnmProtoType m_protoType;
    int m_listennerFd;
};

class SignalControlMsgProcessor : public gse::net::MsgProcessor
{
public:
    SignalControlMsgProcessor();
    virtual ~SignalControlMsgProcessor();

    void SetScheduler(std::shared_ptr<LoadBalanceSchedule> scheduler);
    void SetProtoType(gse::net::EnmProtoType protoType);

public:
    void ProcessMsg(const gse::net::TcpConnectionPtr conn, const gse::net::MessagePtr msg) override;
    void SendAccessRespose(const std::shared_ptr<gse::net::TcpConnection> conn, const char *resp, int len);

private:
    void MarshalKeepaliveReponse(std::string &value);

private:
    std::shared_ptr<LoadBalanceSchedule> m_scheduler;
    gse::net::EnmProtoType m_protoType;
};

class DataReportMsgProcessor : public gse::net::MsgProcessor
{
public:
    DataReportMsgProcessor();
    virtual ~DataReportMsgProcessor();

public:
    void ProcessMsg(const gse::net::TcpConnectionPtr conn, const gse::net::MessagePtr msg) override;

    void SetRecvDataCallback(RecvDataCallBack fnRecvData, void *pCaller);

private:
    RecvDataCallBack m_fnRecvData;
    void *m_pCaller;

private:
};

class OpsReportMsgProcessor : public gse::net::MsgProcessor
{
public:
    OpsReportMsgProcessor();
    ~OpsReportMsgProcessor();

public:
    void ProcessMsg(const gse::net::TcpConnectionPtr conn, const gse::net::MessagePtr msg);

    void SetRecvDataCallback(RecvDataCallBack fnRecvData, void *pCaller);

private:
    RecvDataCallBack m_fnRecvData;
    void *m_pCaller;

private:
};

} // namespace data
} // namespace gse
#endif // !_GSE_DATA_TCPRECEIVER_H_
