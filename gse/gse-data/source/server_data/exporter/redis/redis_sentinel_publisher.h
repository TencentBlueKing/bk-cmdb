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

#ifndef _DATA_REDIS_SENTINEL_PUBLISHER_H_
#define _DATA_REDIS_SENTINEL_PUBLISHER_H_

#include <string>

#include "db/redisapi/async_factory.h"
#include "db/redisapi/factory.h"
#include "eventthread/event_thread.h"
#include "safe/lock.h"

#include "redis_pub_producer.h"
namespace gse {
namespace data {

class RedisSentinelPublisher
{
public:
    RedisSentinelPublisher(const std::string &host, int port, const std::string &passwd);
    ~RedisSentinelPublisher();

public:
    int Init();
    redis::RedisErrorCode Produce(const std::string &key, const std::string &value);
    void GetHost(std::string &host, int &port);
    void SetMasterName(const std::string &masterName);
    void SetSentinelPasswd(const std::string &passwd);

private:
    void Stop();

private:
    std::string m_host;
    int m_port;
    std::string m_sentienlPasswd;
    std::string m_masterPasswd;
    std::string m_masterName;
    gse::safe::RWLock m_clientLock;
    gse::redis::AsyncRedisSentinelPtr m_redisSentinal;
    gse::redis::RedisSentinelPtr m_syncRedisSentinal;
};

} // namespace data
} // namespace gse
#endif
