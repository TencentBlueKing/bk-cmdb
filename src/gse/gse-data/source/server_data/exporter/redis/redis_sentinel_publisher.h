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

#ifndef _DATA_REDIS_SENTINEL_PUB_H_
#define _DATA_REDIS_SENTINEL_PUB_H_

#include <string>
#include "safe/lock.h"
#include "eventthread/gseEventThread.h"
#include "db/redisapi/sync_redis.h"
#include "redis_pub_producer.h"
namespace gse { 
namespace dataserver {

class RedisSentinelPublisher
{
public:
    RedisSentinelPublisher(const std::string &host, int port, const std::string &passwd);
    ~RedisSentinelPublisher();

public:
    static void perMinHandler(int fd, short what, void* param);

public:
    int init();
    // publish to redis
    int produce(const std::string &key, const std::string &value);

    void setMasterName(const std::string &mastername);

    void getRedisMasterHostAndPort(string &host, int &port);
private:
    int checkMaster();
    int createRedisClient();
    int getRedisMaster(std::string &host, int &port);
    void stop();

private:
    std::string m_host;
    int m_port;
    std::string m_masterHost;
    int m_masterPort;
    std::string m_masterPasswd;
    std::string m_masterName;
    RedisPublishProducer *m_client;
    gse::safe::RWLock m_clientLock;
    rgse::GseEventThread m_eventManager;
};

}
}
#endif

