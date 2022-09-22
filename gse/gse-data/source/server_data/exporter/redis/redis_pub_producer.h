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

#ifndef _DATA_REDIS_PUB_PRODUCER_H_
#define _DATA_REDIS_PUB_PRODUCER_H_

#include <string>

#include "db/redisapi/factory.h"
#include "safe/lock.h"

namespace gse {
namespace data {

using namespace std;

class RedisPublishProducer
{
public:
    RedisPublishProducer(const std::string &host, int port, const std::string &passwd);
    ~RedisPublishProducer();

public:
    // init will connect redis
    int init();
    // publish to redis
    redis::RedisErrorCode produce(const std::string &key, const std::string &value);

    void GetHost(std::string &host, int &port);

private:
    std::string m_host;
    int m_port;
    string m_passwd;
    gse::redis::RedisStandalonePtr m_redisStandalone;
    gse::safe::RWLock m_clientLock;
};

} // namespace data
} // namespace gse
#endif
