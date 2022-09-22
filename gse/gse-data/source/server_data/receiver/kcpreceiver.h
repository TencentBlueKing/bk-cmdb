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

#ifndef _GSE_DATA_KCP_RECEIVER_H_
#define _GSE_DATA_KCP_RECEIVER_H_

#include "net/kcp/gseKcpServer.h"
#include "receiver.h"
namespace gse { 
namespace data {

class KcpReceiver : public Receiver, public GseKcpServer
{
public:
    KcpReceiver();
    virtual ~KcpReceiver();

    int Start();
    int Stop();
    void Join();

protected:
private:
    static void receiveHandler(char *buf, uint32_t len, uint32_t srcIp, uint16_t srcPort, void* pCaller, uint32_t serverIP, uint16_t serverPort);
    void localReceiveHanler(char *buf, uint32_t len, uint32_t srcIp, uint16_t srcPort, uint32_t serverIP, uint16_t serverPort);
};

}
}
#endif // !_GSE_DATA_KCP_RECEIVER_H_