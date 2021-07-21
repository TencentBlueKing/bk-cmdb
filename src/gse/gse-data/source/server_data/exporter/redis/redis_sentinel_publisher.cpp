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

#include "redis_sentinel_publisher.h"
#include <vector>
#include "log/log.h"
#include "bbx/gse_errno.h"
#include "tools/macros.h"
namespace gse { 
namespace dataserver {


RedisSentinelPublisher::RedisSentinelPublisher(const std::string &host, int port, const std::string &passwd)
{
    m_client = NULL;
    m_host = host;
    m_port = port;
    m_masterPasswd = passwd;
    m_masterPort = -1;
}

RedisSentinelPublisher::~RedisSentinelPublisher()
{
    stop();

    m_clientLock.WLock();
    if (m_client != NULL)
    {
        delete m_client;
        m_client = NULL;
    }
    m_clientLock.UnLock();
}

int RedisSentinelPublisher::init()
{
    // get redis master when start, and check every 1 min
    checkMaster();

    int iRet = m_eventManager.registerTimerPersistEvent(RedisSentinelPublisher::perMinHandler, this, 60);
    if (iRet < 0)
    {
        LOG_ERROR("fail ot start one timer to check redis master. iRet(%d)", iRet);
        return GSE_ERROR;
    }

    return m_eventManager.Start();
}

void RedisSentinelPublisher::stop()
{
    m_eventManager.clear();
    m_eventManager.stop();
    m_eventManager.Join();
}

void RedisSentinelPublisher::perMinHandler(int fd, short what, void* param)
{
    RedisSentinelPublisher *ptrThis = reinterpret_cast<RedisSentinelPublisher *>(param);
    ptrThis->checkMaster();
}

int RedisSentinelPublisher::checkMaster()
{
    string host;
    int port;

    int ret = getRedisMaster(host, port);
    if (ret != GSE_SUCCESS)
    {
        LOG_WARN("fail to get redis master info. iRet(%d)", ret);
        return ret;
    }

    if (m_masterHost == host && m_masterPort == port)
    {
        return GSE_SUCCESS;
    }

    // master changed
    m_masterHost = host;
    m_masterPort = port;
    LOG_INFO("the redis master has changed, the new redis master is %s:%d", SAFE_CSTR(m_masterHost.c_str()), m_masterPort);

    ret = createRedisClient();
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("fail to create a new redis client which connected to reids master(%s:%d), iRet:%d", SAFE_CSTR(m_masterHost.c_str()), m_masterPort, ret);
        return ret;
    }

    return GSE_SUCCESS;
}

int RedisSentinelPublisher::createRedisClient()
{
    m_clientLock.WLock();
    // release old client
    if (m_client != NULL)
    {
        delete m_client;
        m_client = NULL;
    }

    // passwd is configed on zk
    m_client = new RedisPublishProducer(m_masterHost, m_masterPort, m_masterPasswd);
    int ret = m_client->init();
    if (ret != GSE_SUCCESS)
    {
        delete m_client;
        m_client = NULL;
        m_clientLock.UnLock();
        return ret;
    }
    m_clientLock.UnLock();

    return GSE_SUCCESS;
}

void RedisSentinelPublisher::setMasterName(const std::string &mastername)
{
    m_masterName = mastername;
}

void RedisSentinelPublisher::getRedisMasterHostAndPort(string &host, int &port)
{
    host = m_masterHost;
    port = m_masterPort;
    return;
}
int RedisSentinelPublisher::getRedisMaster(string &host, int &port)
{
    // connect and 'SENTINEL get-master-addr-by-name mymaster' get master ip port
    gse::redis::sync::Redis client(m_host, m_port);
    std::string errmsg;
    bool bret = client.Connect(errmsg);
    if (!bret)
    {
        LOG_ERROR("connect sentinel redis failed, errmsg:%s", errmsg.c_str());
        return GSE_ERROR;
    }

    const string QueryCmd = "SENTINEL get-master-addr-by-name " + m_masterName;
    vector<string> result;
    std::string errormsg;
    bret = client.ExecRedisCommand(QueryCmd, result, errormsg);
    if (!bret)
    {
        LOG_ERROR("SENTINEL get-master-addr-by-name %s faile, errormsg:%s", m_masterName.c_str(), errormsg.c_str());
        return GSE_ERROR;
    }

    client.Close();

    // result is ip and port
    // > SENTINEL get-master-addr-by-name mymaster
    // 1) "ip"
    // 2) "port"
    //
    if (result.size() != 2)
    {
        LOG_ERROR("SENTINEL get-master-addr-by-name %s result error", m_masterName.c_str());
        return GSE_ERROR;
    }
    host = result[0];
    port = gse::tools::strings::StringToInt16(result[1]);

    return GSE_SUCCESS;
}

int RedisSentinelPublisher::produce(const string &key, const string &value)
{
    m_clientLock.RLock();
    if (m_client == NULL)
    {
        m_clientLock.UnLock();
        return GSE_ERROR;
    }
    int ret = m_client->produce(key, value);
    m_clientLock.UnLock();
    return ret;
}

}
}
