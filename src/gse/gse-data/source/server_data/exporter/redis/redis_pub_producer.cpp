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

#include "bbx/gse_errno.h"
#include "db/redisapi/factory.h"
#include "db/redisapi/options.hpp"
#include "log/log.h"
#include "tools/finally.hpp"

namespace gse {
namespace data {

RedisPublishProducer::RedisPublishProducer(const std::string &host, int port, const std::string &passwd)
{
    m_host = host;
    m_port = port;
    m_passwd = passwd;
}

RedisPublishProducer::~RedisPublishProducer()
{
    m_redisStandalone.reset();
}

int RedisPublishProducer::init()
{
    gse::redis::RedisOptions options(m_host, m_port, "default", m_passwd);

    m_redisStandalone = gse::redis::RedisFactory::CreateRedisStandalone(options);
    if (m_redisStandalone == nullptr)
    {
        LOG_ERROR("failed to create redis standalone client, api return nullptr");
        return GSE_ERROR;
    }

    if (!m_redisStandalone->IsOK())
    {
        LOG_ERROR("failed to create redis sentinal client, error:%s", m_redisStandalone->Error());
        return GSE_ERROR;
    }

    LOG_DEBUG("create sentinal redis client(%s:%d) success", m_host.c_str(), m_port);
    return GSE_SUCCESS;
}

gse::redis::RedisErrorCode RedisPublishProducer::produce(const std::string &key, const std::string &value)
{
    m_clientLock.RLock();
    auto _ = gse::tools::defer::finally([this]() {
        m_clientLock.UnLock();
    });

    LOG_DEBUG("publish %s %s to redis", key.c_str(), value.c_str());
    return m_redisStandalone->Publish(key, value);
}

void RedisPublishProducer::GetHost(std::string &host, int &port)
{
    host = m_host;
    port = m_port;
}
} // namespace data
} // namespace gse
