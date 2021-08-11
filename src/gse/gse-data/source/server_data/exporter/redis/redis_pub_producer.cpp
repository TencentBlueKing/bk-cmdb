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

#include "redis_pub_producer.h"

#include "db/redisapi/async_redis.h"
#include "bbx/gse_errno.h"
namespace gse { 
namespace dataserver {

RedisPublishProducer::RedisPublishProducer(const std::string &host, int port, const std::string &passwd)
{
    m_host = host;
    m_port = port;
    m_passwd = passwd;
    //m_client = new RedisWriter(true);
    m_client = new gse::redis::RedisMsgWriter(m_host, port);
    if (!m_passwd.empty())
    {
        m_client->SetPassword(m_passwd);
    }
}

RedisPublishProducer::~RedisPublishProducer()
{
    if (m_client != NULL)
    {
        //m_client->close();
        m_client->Stop();
        m_client->Join();
        delete m_client;
        m_client = NULL;
    }
}

int RedisPublishProducer::init()
{
    m_client->Start();
    return GSE_SUCCESS;
}

int RedisPublishProducer::produce(const std::string &key, const std::string &value)
{
    int ret = GSE_SUCCESS;
    m_client->SendMsg("publish", key.c_str(), value.c_str(), value.size());
    return ret;
}

void RedisPublishProducer::GetHost(std::string &host, int &port)
{
    host = m_host;
    port = m_port;
}
}
}
