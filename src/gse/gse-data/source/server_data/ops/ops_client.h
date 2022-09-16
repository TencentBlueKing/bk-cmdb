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

#ifndef _GSE_DATA_OPS_OPSCLIENT_H_
#define _GSE_DATA_OPS_OPSCLIENT_H_

#include <arpa/inet.h>
#include <string>
#include <sys/socket.h>
#include <vector>

#include "net/udp/udp_client.h"
#include "protocol.h"

namespace gse {
namespace data {

class OPSServerAddress
{
public:
    OPSServerAddress(const std::string &ip, uint16_t port)
        : m_hostIP(ip), m_hostPort(port){};

public:
    inline std::string GetHostIP() { return m_hostIP; };
    inline uint16_t GetHostPort() { return m_hostPort; };

private:
    std::string m_hostIP;
    uint16_t m_hostPort;
};

class OPSClient
{

public:
    OPSClient();
    ~OPSClient();

public:
    int OpenUDP(const std::vector<OPSServerAddress> &addresses);
    int SendToByUDP(uint32_t channelID, const char *ptrData, int dataLen);
    int CloseUDP();

private:
    std::vector<net::udp::UDPClient *> m_udpClients;
    uint64_t m_usedUDPClientIndex;
};

} // namespace data
} // namespace gse
#endif /* SOURCE_GSE_DATA_SERVER_DATA_DSCLIENT_OPSCLIENT_H_ */
