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

#ifndef CONF_COMMON_H
#define CONF_COMMON_H
#include <string>
namespace gse {
namespace data {

class KafkaConfig
{
public:
    std::string m_securityProtocol;
    std::string m_saslMechanisms;
    std::string m_saslUserName;
    std::string m_saslPasswd;
    int m_messageMaxBytes;
    int m_queueBufferingMaxMessages;
    std::string m_requestRequiredAcks;
    int m_queueBufferingMaxMs;
    int m_msgTimeoutMs;
    std::string m_partitioner;
    std::string m_clientid;
    KafkaConfig()
    {
        m_messageMaxBytes = 10000000; // 10MB
        m_requestRequiredAcks = "1";
        m_queueBufferingMaxMs = 0;
        m_queueBufferingMaxMessages = 0;
        m_msgTimeoutMs = 0;
        m_partitioner = "murmur2_random";
    }
};

} // namespace data
} // namespace gse
#endif // CONF_COMMON_H
