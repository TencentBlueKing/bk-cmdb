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

#ifndef _DATASERVER_OPS_COMPUTE_OPSSERVER_H_
#define _DATASERVER_OPS_COMPUTE_OPSSERVER_H_

#include "gseUdpServer.h"
#include "opsDeal.h"
#include "gseMsg.h"
#include "opszkclient.h"
namespace gse { 
namespace dataserver {
/**
 * @brief 运营服务器
 */
class OpsServer
{
public:
    static OpsServer &Instance();
    ~OpsServer();

public:
    /**
     * @brief 消息回调
     * @param ptrData 数据
     * @param dataLength 消息长度
     * @param hostip 主机ip
     * @param hostport 主机port
     *
     */
    static void OpsDealCallback(char * ptrData, uint32_t dataLength, uint32_t hostip, uint16_t hostport);

    /**
     * @brief 处理消息数据
     * @param dataid 数据id
     * @param data 数据
     */
    static void dealDataOps(int dataId, const string& data);

public:
    int start();
    int stop();

    OpsDeal m_opsDeal;
private:
    OpsServer();

    GseUdpServer m_udpServer;
    OpsZKClient m_zk;
};

}
}
#endif

