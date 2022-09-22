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

#include "ops_client.h"

#include <stdlib.h>
#include <unistd.h>

#include "bbx/gse_errno.h"
#include "log/log.h"
#include "net/udp/udp_client.h"
#include "tools/atomic.h"
#include "tools/error.h"

#include "protocol.h"

namespace gse {
namespace data {

#define LOCAL_HOST "127.0.0.1"

OPSClient::OPSClient()
{
    m_usedUDPClientIndex = 0;
}

OPSClient::~OPSClient()
{
    CloseUDP();
}

int OPSClient::OpenUDP(const std::vector<OPSServerAddress> &addresses)
{
    std::size_t address_cnt = addresses.size();
    for (std::size_t address_idx = 0; address_idx < address_cnt; ++address_idx)
    {
        OPSServerAddress address_info = addresses.at(address_idx);
        std::string host_ip = address_info.GetHostIP();
        uint16_t host_port = address_info.GetHostPort();

        gse::net::udp::UDPClient *ptr_udpclient = new gse::net::udp::UDPClient();
        bool result_ok = ptr_udpclient->Connect(host_port, host_ip);
        if (!result_ok)
        {
            delete ptr_udpclient;
            LOG_ERROR("it is failed to connet to the remote ops server(%s:%d)", host_ip.c_str(), host_port);
            continue;
        }
        m_udpClients.push_back(ptr_udpclient);
    }

    return GSE_SUCCESS;
}

int OPSClient::CloseUDP()
{
    std::size_t udp_client_cnt = m_udpClients.size();
    for (std::size_t idx = 0; idx < udp_client_cnt; ++idx)
    {
        gse::net::udp::UDPClient *ptr_udpclient = m_udpClients.at(idx);
        delete ptr_udpclient;
    }
    m_udpClients.clear();
    return GSE_SUCCESS;
}

int OPSClient::SendToByUDP(uint32_t channelID, const char *ptrData, int dataLen)
{
    if (m_udpClients.empty())
    {
        LOG_INFO("not found andy valid udp connection to send ops data");
        return GSE_ERROR;
    }
    GSEDataPackage gse_data_package;
    gse_data_package.SetChannelId(channelID);
    gse_data_package.Pack(ptrData, dataLen);

    std::size_t to_used_idx = gse::tools::atomic::AtomAdd(&m_usedUDPClientIndex, 1) % m_udpClients.size();

    gse::net::udp::UDPClient *ptr_udpclient = m_udpClients.at(to_used_idx);

    return ptr_udpclient->Send(gse_data_package.GetDataPointer(), gse_data_package.GetDataLength());
}

} // namespace data
} // namespace gse
