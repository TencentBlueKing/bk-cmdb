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

#include "bbx/gse_errno.h"
#include "db/redisapi/async_factory.h"
#include "db/redisapi/factory.h"
#include "db/redisapi/options.hpp"
#include "log/log.h"
#include "tools/finally.hpp"
#include "tools/macros.h"

namespace gse {
namespace data {

RedisSentinelPublisher::RedisSentinelPublisher(const std::string &host, int port, const std::string &passwd)
{
    m_host = host;
    m_port = port;
    m_masterPasswd = passwd;
    m_redisSentinal = nullptr;
}

RedisSentinelPublisher::~RedisSentinelPublisher()
{
    Stop();
}

int RedisSentinelPublisher::Init()
{
    gse::redis::RedisSentinelOptions options(m_masterName, "default", m_masterPasswd, m_sentienlPasswd);
    options.SetNodes({{m_host, m_port}});

    m_redisSentinal = gse::redis::AsyncRedisFactory::CreateSentinelMaster(options);
    if (m_redisSentinal == nullptr)
    {
        LOG_ERROR("failed to create redis sentinal client, api return nullptr");
        return GSE_ERROR;
    }

    if (!m_redisSentinal->IsOK())
    {
        LOG_ERROR("failed to create redis sentinal client, error:%s", m_redisSentinal->Error());
        return GSE_ERROR;
    }

    LOG_DEBUG("create sentinal redis client(%s:%d), mastername:%s", m_host.c_str(), m_port, m_masterName.c_str());

    return GSE_SUCCESS;
}

void RedisSentinelPublisher::Stop()
{
    m_redisSentinal.reset();
}

void RedisSentinelPublisher::SetMasterName(const std::string &masterName)
{
    m_masterName = masterName;
}

void RedisSentinelPublisher::SetSentinelPasswd(const std::string &passwd)
{
    m_sentienlPasswd = passwd;
}

gse::redis::RedisErrorCode RedisSentinelPublisher::Produce(const string &key, const string &value)
{
    m_clientLock.RLock();
    auto _ = gse::tools::defer::finally([this]() {
        m_clientLock.UnLock();
    });

    LOG_DEBUG("publish %s %s to redis", key.c_str(), value.c_str());
    return m_redisSentinal->Publish(key, value);
}

void RedisSentinelPublisher::GetHost(std::string &host, int &port)
{
    host = m_host;
    port = m_port;
}

} // namespace data
} // namespace gse
