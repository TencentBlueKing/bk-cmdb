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
#include "thread/thread.h"
#include "receiver.h"
namespace gse { 
namespace dataserver {



class TcpReceiver : public Receiver
{
public:
    TcpReceiver();
    virtual ~TcpReceiver();

public:
    int Start();
    int Stop();
    void Join();

private:
    // implement virtual function
    void OnConnectionEvent(const std::shared_ptr<gse::net::TcpConnection> &conn);
    void MsgHandler(const std::shared_ptr<gse::net::TcpConnection> &conn, const std::shared_ptr<gse::net::Message> &msg);
    uint32_t getMsgHeadLen();
    void setProtoMsgHead();
protected:
    bool loadServerCert(const std::string& certRootPath);

private:
    //TcpServerWorker* m_tcpwork;
    std::shared_ptr<gse::net::TcpServer> m_tcpServer;
    std::shared_ptr<std::thread> m_tcpServerThread;

};

}
}
#endif // !_GSE_DATA_TCPRECEIVER_H_
