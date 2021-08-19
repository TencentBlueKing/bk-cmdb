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

#ifndef _GSE_DATA_UDP_RECEIVER_H_
#define _GSE_DATA_UDP_RECEIVER_H_

#include <unistd.h>
#include "net/udp/udp_server.h"
#include "eventloop/event_loop.h"
#include "receiver/receiver.h"
namespace gse { 
namespace dataserver {


class UdpReceiver : public Receiver
{
public:
    UdpReceiver();
    virtual ~UdpReceiver();

    int Start();
    int Stop();
    void Join();
protected:
private:
    void MsgHandler(gse::eventloop::EventLoop* loop, gse::net::udp::UdpMessagePtr & msg);

private:
    std::shared_ptr<gse::net::udp::UdpServer> m_udpServer;
    std::shared_ptr<std::thread> m_udpServerThread;
    //UdpServerWorker* m_udpserver;
    int m_serverIp;
    int m_serverPort;

};

}
}
#endif // !_GSE_DATA_UDP_RECEIVER_H_
