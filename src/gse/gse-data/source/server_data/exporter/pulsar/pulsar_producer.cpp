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

#include "pulsar_producer.h"

#include "bbx/gse_errno.h"
#include "log/log.h"
#include "ops/op_metric.h"
#include "tools/finally.hpp"
#include "utils.h"

namespace gse {
namespace data {

void PulsarProducer::pulsar_logger(pulsar_logger_level_t level, const char *file, int line, const char *message,
                                   void *ctx)
{
    LOG_DEBUG("pulsar_logger, LOGLEVEL:%d,file:%s,line:%d,message:%s", level, file, line, message);
}
PulsarProducer::PulsarProducer()
{
    m_client = NULL;
    m_conf = NULL;
    m_ssl = false;
    m_auth = NULL;
}

PulsarProducer::~PulsarProducer()
{
    for (auto it : m_producers)
    {
        delete it.second;
    }

    m_producers.clear();
    LOG_DEBUG("clean pulsar producer");
}

void PulsarProducer::getDefaultClientConfig()
{
    LOG_DEBUG("producer default config:");
    LOG_DEBUG("io threads:%d", pulsar_client_configuration_get_io_threads(m_conf));
    LOG_DEBUG("message listener threads:%d", pulsar_client_configuration_get_message_listener_threads(m_conf));
    LOG_DEBUG("operation timeout seconds:%d", pulsar_client_configuration_get_operation_timeout_seconds(m_conf));
    LOG_DEBUG("tls trust certs file path:%s", pulsar_client_configuration_get_tls_trust_certs_file_path(m_conf));
    // LOG_DEBUG("tls trust certs file path:%d", pulsar_client_configuration_get_tls_trust_certs_file_path(m_conf));
}

pulsar_authentication_t *PulsarProducer::pulsarAuthenticationCreate(const std::string &trust_certificate, const std::string &keycertificate, const std::string &certtifile)
{
    // pulsar_authentication_create_t *auth = (pulsar_authentication_create)
    return NULL;
}

int PulsarProducer::createProducer(const std::string &service_url, const std::string &certificatePath, const std::string &privateKeyPath, const std::string &token)
{

    m_conf = pulsar_client_configuration_create();

    if (certificatePath != "")
    {
        m_auth = pulsar_authentication_tls_create(certificatePath.c_str(), privateKeyPath.c_str());
        pulsar_client_configuration_set_auth(m_conf, m_auth);
    }

    if (!token.empty())
    {
        m_auth = pulsar_authentication_token_create(token.c_str());
        pulsar_client_configuration_set_auth(m_conf, m_auth);
        LOG_DEBUG("create pulsar client token:%s", token.c_str());
    }

    pulsar_client_configuration_set_logger(m_conf, PulsarProducer::pulsar_logger, (void *)this);
    m_client = pulsar_client_create(service_url.c_str(), m_conf);
    LOG_DEBUG("create pulsar client(%s) success", service_url.c_str());
    getDefaultClientConfig();
    return GSE_SUCCESS;
}

void PulsarProducer::handleSendCallback(pulsar_result result, pulsar_message_id_t *msgId, void *ctx)
{
    if (msgId == nullptr)
    {
        LOG_ERROR("failed to process pulsar msg callback, msgid invalid, send msg result:%s", pulsar_result_str(result));
        OPMetric::AddSendMsgCounter(kOutputPulsar, kOutputPulsar, result, 1);
        return;
    }

    pulsar_message_t *message = (pulsar_message_t *)ctx;
    if (message == nullptr)
    {
        LOG_ERROR("failed to send msg result:%s, message invalid", pulsar_result_str(result));
        OPMetric::AddSendMsgCounter(kOutputPulsar, kOutputPulsar, result, 1);
        return;
    }

    uint64_t costTime = 0;

    uint64_t beginTimestamp = pulsar_message_get_event_timestamp(message);
    uint64_t endTimestamp = gse::tools::time::GetUTCMillsec();
    costTime = (endTimestamp >= beginTimestamp) ? (endTimestamp - beginTimestamp) : 0;
    OPMetric::AddMessageQueneRequestMillsecondsMetrics("pulasr", costTime);

    if (result != pulsar_result_Ok)
    {
        LOG_ERROR("send msg to pulsar topic[%s] failed, error(%d:%s)",
                  pulsar_message_get_topic_name(message), result, pulsar_result_str(result));
        OPMetric::AddSendMsgCounter(kOutputPulsar, kOutputPulsar, result, 1);
    }
    else
    {
        LOG_DEBUG("send msg to pulsar topic[%s] success", pulsar_message_get_topic_name(message));
        OPMetric::AddSendMsgBytesCounter(kOutputPulsar, kOutputPulsar, pulsar_message_get_length(message));
        OPMetric::AddSendMsgCounter(kOutputPulsar, kOutputPulsar, result, 1);
    }

    if (message != nullptr)
    {
        pulsar_message_free(message);
    }

    pulsar_message_id_free(msgId);
}

pulsar_producer_t *PulsarProducer::findProducer(const std::string &topic)
{
    auto it = m_producers.find(topic);
    if (it == m_producers.end())
    {
        return NULL;
    }
    return it->second->m_producer;
}

pulsar_producer_t *PulsarProducer::newProducer(const string &topic)
{
    pulsar_producer_t *producer = NULL;
    producer = findProducer(topic);
    if (producer == NULL)
    {
        pulsar_producer_configuration_t *producerConf = pulsar_producer_configuration_create();
        pulsar_producer_configuration_set_batching_enabled(producerConf, 1);
        pulsar_producer_configuration_set_max_pending_messages(producerConf, 2000000);
        pulsar_producer_configuration_set_max_pending_messages_across_partitions(producerConf, 4000000);

        LOG_DEBUG("pulsar producer set max pending messages:%d", 2000000);
        LOG_DEBUG("Create producer for topic(%s), m_client:%p, m_conf:%p", topic.c_str(), m_client, m_conf);
        try
        {
            pulsar_result err = pulsar_client_create_producer(m_client, topic.c_str(), producerConf, &producer);
            if (err != pulsar_result_Ok)
            {
                LOG_ERROR("failed to create producer:%s", pulsar_result_str(err));
                pulsar_producer_configuration_free(producerConf);
                return NULL;
            }

            LOG_DEBUG("Create producer for topic(%s) success", topic.c_str());
            PulsarProducerObject *ptrPuslarProducer = new PulsarProducerObject();
            ptrPuslarProducer->m_producer = producer;
            ptrPuslarProducer->m_producer_conf = producerConf;
            m_producers.insert(std::pair<std::string, PulsarProducerObject *>(topic, ptrPuslarProducer));
        }
        catch (const std::exception &e)
        {
            LOG_ERROR("create producer exception:%s", e.what());
            return NULL;
        }
    }
    else
    {
        LOG_DEBUG("producer has created, name:%s", pulsar_producer_get_producer_name(producer));
    }

    return producer;
}

int PulsarProducer::excuteProduce(const std::string &topic, int partition, const std::string &key, const std::string &data)
{
    pulsar_producer_t *producer = newProducer(topic);
    if (producer == NULL)
    {
        LOG_DEBUG("Failed to create producer for topic(%s)", topic.c_str());
        return GSE_ERROR;
    }

    uint64_t beginTimestamp = gse::tools::time::GetUTCMillsec();

    pulsar_message_t *message = pulsar_message_create();
    pulsar_message_set_content(message, (const void *)data.c_str(), data.length());
    pulsar_message_set_property(message, "message_key", key.c_str());
    pulsar_message_set_event_timestamp(message, beginTimestamp);
    pulsar_producer_send_async(producer, message, PulsarProducer::handleSendCallback, (void *)message);
    return GSE_SUCCESS;
}

void PulsarProducer::closeProducer()
{
    pulsar_client_close(m_client);
    pulsar_client_free(m_client);
    pulsar_client_configuration_free(m_conf);
    LOG_DEBUG("free pulsar client");
}

} // namespace data
} // namespace gse
