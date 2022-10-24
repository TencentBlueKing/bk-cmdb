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

#ifndef _PULSAR_PRODUCER_H_
#define _PULSAR_PRODUCER_H_

#include <ctype.h>
#include <signal.h>
#include <sys/queue.h>

#include <list>
#include <memory>
#include <string>
#include <unordered_map>

#include <pulsar/c/authentication.h>
#include <pulsar/c/client.h>

#include "eventthread/event_thread.h"
#include "log/log.h"

namespace gse {
namespace data {

using namespace std;

class PulsarProducerObject
{
public:
    PulsarProducerObject()
    {
        m_producer = NULL;
        m_producer_conf = NULL;
    };

    ~PulsarProducerObject()
    {
        pulsar_producer_close(m_producer);
        pulsar_producer_free(m_producer);
        m_producer = NULL;

        pulsar_producer_configuration_free(m_producer_conf);
        m_producer_conf = NULL;

        LOG_DEBUG("free pulsar producer");
    }

    PulsarProducerObject(pulsar_producer_t *producer, pulsar_producer_configuration_t *conf)
        : m_producer(producer), m_producer_conf(conf) {}

    pulsar_producer_t *m_producer;
    pulsar_producer_configuration_t *m_producer_conf;
};

class PulsarProducer
{

public:
    PulsarProducer();
    ~PulsarProducer();

public:
    int createProducer(const string &service_url, const string &certificatePath, const string &privateKeyPath, const std::string &token);

    int getProducerState();
    int excuteProduce(const string &topic, int partition, const std::string &key, const string &value);
    void closeProducer();
    int getProducerQueueCount();
    pulsar_authentication_t *pulsarAuthenticationCreate(const std::string &trust_certificate, const std::string &keycertificate, const std::string &certtifile);

private:
    void getDefaultClientConfig();
    pulsar_producer_t *newProducer(const string &topic);

    pulsar_producer_t *findProducer(const std::string &topic);
    static void handleSendCallback(pulsar_result result, pulsar_message_id_t *msgId, void *ctx);
    static void pulsar_logger(pulsar_logger_level_t level, const char *file, int line, const char *message,
                              void *ctx);

private:
    string m_url;
    bool m_ssl;
    pulsar_client_configuration_t *m_conf;
    pulsar_client_t *m_client;
    pulsar_authentication_t *m_auth;
    gse::safe::RWLock m_rwlock;
    std::unordered_map<std::string, PulsarProducerObject *> m_producers;
};

} // namespace data
} // namespace gse
#endif
