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

#include "kafka_producer.h"

#include "bbx/gse_errno.h"
#include "log/log.h"
#include "tools/error.h"
#include "tools/filesystem.h"
#include "tools/macros.h"
#include "utils.h"
#include <fstream>
#include <iostream>
#include <stdio.h>

#include "ops/op_metric.h"
#include "tools/finally.hpp"

namespace gse {
namespace data {

std::string KafkaProducer::m_runtimeDataDirector("rt");

static void logger(const rd_kafka_t *rk, int level,
                   const char *fac, const char *buf)
{
    static int i = 0;
    i++;
    if (i >= 10)
    {
        LOG_ERROR("KAFKA-%i-%s: %s %s", level, fac, rd_kafka_name(rk), buf);
        i = 0;
    }
}

// save stat data on local and remote
// remote send every 1min
static int stats_cb(rd_kafka_t *rk, char *json, size_t json_len, void *opaque)
{
    // LOG_DEBUG("kafka stats callback, rdkafks(%s)", SAFE_CSTR(json));
    if (NULL != json)
    {
        // GLOG_DATA("timestamp:%s data: %s", tgse::localTimeToString().c_str(), json);

        /* save into file */
        string pName(rd_kafka_name(rk));
        string rtpath = "./rt/";
        string file = rtpath + pName + ".json";
        ofstream ofs;
        ofs.open(file.c_str());
        LOG_DEBUG("kafka runtime data director is %s", file.c_str());

        if (!gse::tools::filesystem::IsDirExist(rtpath.c_str()))
        {
            gse::tools::filesystem::DirTreeMake(rtpath.c_str());
        }
        if (!ofs.is_open())
        {
            LOG_ERROR("[kafka prodducer] open file[%s] fail,errno[%d] errorstr %s", file.c_str(), errno, gse::tools::error::ErrnoToStr(errno).c_str());
            return 0;
        }

        ofs.write(json, json_len);
        ofs.close();
    }

    return 0;
}

void KafkaProducer::MsgDeliverCb(rd_kafka_t *rk,
                                 const rd_kafka_message_t *rkmessage, void *opaque)
{
    OPMetric::AddSendMsgCounter(kOutputKafka, kOutputKafka, rkmessage->err, 1);
    uint64_t costTime = 0;

    if (rkmessage->_private != NULL)
    {
        uint64_t *beginTimestamp = (uint64_t *)rkmessage->_private;
        uint64_t endTimestamp = gse::tools::time::GetUTCMillsec();
        costTime = (endTimestamp >= *beginTimestamp) ? (endTimestamp - *beginTimestamp) : 0;

        OPMetric::AddMessageQueneRequestMillsecondsMetrics(kOutputKafka, costTime);
        free(rkmessage->_private);
    }

    if (rkmessage->err)
    {
        LOG_ERROR("delivery message to topic:%s failed, error:%s",
                  rd_kafka_topic_name(rkmessage->rkt),
                  rd_kafka_err2str(rkmessage->err));
    }
    else
    {
        LOG_DEBUG("deliver message success, elapsed time:%d(ms)), (topic %s, %d bytes, partition %d, offset %d)",
                  costTime, rd_kafka_topic_name(rkmessage->rkt), rkmessage->len,
                  rkmessage->partition, rkmessage->offset);
        OPMetric::AddSendMsgBytesCounter(kOutputKafka, kOutputKafka, rkmessage->len);
    }
    /* The rkmessage is destroyed automatically by librdkafka */
}

KafkaProducer::KafkaProducer()
{
    this->m_rdKafa = NULL;
    this->conf = NULL;
    m_maxMessageBytes = 0;
    m_lastlogTime.tv_sec = 0;
    m_lastlogTime.tv_usec = 0;
}

KafkaProducer::~KafkaProducer()
{
}

int KafkaProducer::CreateProducer(const string &broker)
{
    m_brokers = broker;
    char errstr[512] = {0};

    conf = rd_kafka_conf_new();
    if (NULL == conf)
    {
        LOG_ERROR("fail to create kafka conf object");
        return GSE_SYSTEMERROR;
    }

    rd_kafka_conf_set_log_cb(conf, logger);
    rd_kafka_conf_set_stats_cb(conf, stats_cb);
    rd_kafka_conf_set_dr_msg_cb(conf, &KafkaProducer::MsgDeliverCb);

    if (m_kafkaConfig.m_queueBufferingMaxMessages > 0)
    {
        if (RD_KAFKA_CONF_OK != rd_kafka_conf_set(conf, "queue.buffering.max.messages", gse::tools::strings::ToString(m_kafkaConfig.m_queueBufferingMaxMessages).c_str(), errstr, sizeof(errstr)))
        {
            LOG_ERROR("set kakfa queue.buffering.max.messages(%d) failed, error:%s", m_kafkaConfig.m_queueBufferingMaxMessages, errstr);
            return GSE_SYSTEMERROR;
        }

        LOG_DEBUG("set kakfa queue.buffering.max.messages(%d) success", m_kafkaConfig.m_queueBufferingMaxMessages);
    }

    if (RD_KAFKA_CONF_OK != rd_kafka_conf_set(conf, "broker.version.fallback", "0.8.2", errstr, sizeof(errstr)))
    {
        LOG_ERROR("set kakfa broker.version.fallback(%s) failed, error:%s", "0.8.2", errstr);
        return GSE_SYSTEMERROR;
    }

    if (m_kafkaConfig.m_messageMaxBytes > 0)
    {
        if (rd_kafka_conf_set(conf, "message.max.bytes", gse::tools::strings::ToString(m_kafkaConfig.m_messageMaxBytes).c_str(), errstr, sizeof(errstr)) != RD_KAFKA_CONF_OK)
        {
            LOG_ERROR("set kakfa message.max.bytes(%d) failed, error:%s", m_kafkaConfig.m_messageMaxBytes, errstr);
            return GSE_SYSTEMERROR;
        }
        LOG_DEBUG("set kakfa message.max.bytes(%d) success", m_kafkaConfig.m_messageMaxBytes);
    }

    if (m_kafkaConfig.m_requestRequiredAcks != "")
    {
        if (rd_kafka_conf_set(conf, "request.required.acks", m_kafkaConfig.m_requestRequiredAcks.c_str(), errstr, sizeof(errstr)) != RD_KAFKA_CONF_OK)
        {
            LOG_ERROR("set request.required.acks(%s) failed, error:%s", m_kafkaConfig.m_requestRequiredAcks.c_str(), errstr);
            return GSE_SYSTEMERROR;
        }
        LOG_DEBUG("set request.required.acks(%s) success", m_kafkaConfig.m_requestRequiredAcks.c_str());
    }

    if (m_kafkaConfig.m_queueBufferingMaxMs > 0)
    {
        if (rd_kafka_conf_set(conf, "queue.buffering.max.ms", gse::tools::strings::ToString(m_kafkaConfig.m_queueBufferingMaxMs).c_str(), errstr, sizeof(errstr)) != RD_KAFKA_CONF_OK)
        {
            LOG_ERROR("set kafaka queue.buffering.max.ms:%d failed, error:%s", m_kafkaConfig.m_queueBufferingMaxMs, errstr);
            return GSE_SYSTEMERROR;
        }
        LOG_DEBUG("set queue.buffering.max.ms(%d) success", m_kafkaConfig.m_queueBufferingMaxMs);
    }

    if (m_kafkaConfig.m_msgTimeoutMs > 0)
    {
        if (rd_kafka_conf_set(conf, "message.timeout.ms", gse::tools::strings::ToString(m_kafkaConfig.m_msgTimeoutMs).c_str(), errstr, sizeof(errstr)) != RD_KAFKA_CONF_OK)
        {
            LOG_ERROR("set kafaka message.timeout.ms:%d failed, error:%s", m_kafkaConfig.m_msgTimeoutMs, errstr);
            return GSE_SYSTEMERROR;
        }
        LOG_DEBUG("set message.timeout.ms(%d) success", m_kafkaConfig.m_msgTimeoutMs);
    }

    if (m_kafkaConfig.m_partitioner != "")
    {
        if (rd_kafka_conf_set(conf, "partitioner", m_kafkaConfig.m_partitioner.c_str(), errstr, sizeof(errstr)) != RD_KAFKA_CONF_OK)
        {
            LOG_ERROR("set kafaka partitioner:%s failed, error:%s", m_kafkaConfig.m_partitioner.c_str(), errstr);
            return GSE_SYSTEMERROR;
        }
        LOG_DEBUG("set partitioner(%s) success", m_kafkaConfig.m_partitioner.c_str());
    }

    if (m_kafkaConfig.m_clientid != "")
    {
        if (rd_kafka_conf_set(conf, "client.id", m_kafkaConfig.m_clientid.c_str(), errstr, sizeof(errstr)) != RD_KAFKA_CONF_OK)
        {
            LOG_ERROR("set kafaka client.id:%s failed, error:%s", m_kafkaConfig.m_partitioner.c_str(), errstr);
            return GSE_SYSTEMERROR;
        }
        LOG_DEBUG("set client.id(%s) success", m_kafkaConfig.m_clientid.c_str());
    }

    // default plaintext
    if (m_kafkaConfig.m_saslUserName != "")
    {
        if (rd_kafka_conf_set(conf, "security.protocol", m_kafkaConfig.m_securityProtocol.c_str(), errstr, sizeof(errstr) != RD_KAFKA_CONF_OK))
        {
            LOG_ERROR("set kakfa security.protocol(%s) failed, error:%s", m_kafkaConfig.m_securityProtocol.c_str(), errstr);
            return GSE_SYSTEMERROR;
        }
        LOG_DEBUG("set kakfa security.protocol(%s) success", m_kafkaConfig.m_securityProtocol.c_str());

        if (rd_kafka_conf_set(conf, "sasl.mechanisms", m_kafkaConfig.m_saslMechanisms.c_str(), errstr, sizeof(errstr) != RD_KAFKA_CONF_OK))
        {
            LOG_ERROR("set kakfa security.protocol(%s) failed, error:%s", m_kafkaConfig.m_saslMechanisms.c_str(), errstr);
            return GSE_SYSTEMERROR;
        }
        LOG_DEBUG("set kakfa sasl.mechanisms(%s) success", m_kafkaConfig.m_saslMechanisms.c_str());

        if (rd_kafka_conf_set(conf, "sasl.username", m_kafkaConfig.m_saslUserName.c_str(), errstr, sizeof(errstr) != RD_KAFKA_CONF_OK))
        {
            LOG_ERROR("set kakfa security.protocol(%s) failed, error:%s", m_kafkaConfig.m_saslUserName.c_str(), errstr);
            return GSE_SYSTEMERROR;
        }
        LOG_DEBUG("set kakfa sasl.username(%s) success", m_kafkaConfig.m_saslUserName.c_str());

        if (rd_kafka_conf_set(conf, "sasl.password", m_kafkaConfig.m_saslPasswd.c_str(), errstr, sizeof(errstr) != RD_KAFKA_CONF_OK))
        {
            LOG_ERROR("set kakfa security.protocol(%s) failed, error:%s", m_kafkaConfig.m_saslPasswd.c_str(), errstr);
            return GSE_SYSTEMERROR;
        }
        LOG_DEBUG("set kakfa sasl.password(%s) success", m_kafkaConfig.m_saslPasswd.c_str());
    }

    if (rd_kafka_conf_set(conf, "statistics.interval.ms", "600000", errstr, sizeof(errstr)) != RD_KAFKA_CONF_OK)
    {
        LOG_ERROR("set kafka statistics.interval.ms failed, error:%s", errstr);
        return GSE_SYSTEMERROR;
    }

    m_rdKafa = rd_kafka_new(RD_KAFKA_PRODUCER, conf, errstr, sizeof(errstr));
    if (m_rdKafa == NULL)
    {
        LOG_ERROR("kafka producer create producer fail, error:%s, broker addr :%s", errstr, m_brokers.c_str());
        return GSE_SYSTEMERROR;
    }

    if (rd_kafka_brokers_add(m_rdKafa, m_brokers.c_str()) == 0)
    {
        LOG_ERROR("kafka producer create producer fail for no valid broker found in %s", m_brokers.c_str());
        return GSE_SYSTEMERROR;
    }

    return GSE_SUCCESS;
}

void KafkaProducer::SetMaxMessageBytes(int bytes)
{
    m_maxMessageBytes = bytes;
}

void KafkaProducer::SetKafkaConfig(KafkaConfig &kafa_conf)
{
    m_kafkaConfig = kafa_conf;
}

int KafkaProducer::ExcuteProduce(const string &topic, int partition, const std::string &key, const string &value, const std::string &auxiliary)
{
    int res;
    rd_kafka_topic_t *rkt = NULL;
    int errcode = 0;

    uint64_t *beginTimestamp = (uint64_t *)malloc(sizeof(uint64_t));
    *beginTimestamp = gse::tools::time::GetUTCMillsec();
    auto _ = gse::tools::defer::finally([&]() {
        // uint64_t endTimestamp = gse::tools::time::GetUTCMillsec();
        // uint64_t costTime = (endTimestamp >= beginTimestamp) ? (endTimestamp >= beginTimestamp) : 0;
        if (errcode != 0)
        {
            OPMetric::AddSendMsgCounter(kOutputKafka, kOutputKafka, errcode, 1);
        }

        // OPMetric::AddMessageQueneRequestMillsecondsMetrics("kafka", costTime);
    });

    rkt = rd_kafka_topic_new(m_rdKafa, topic.c_str(), NULL);
    if (NULL == rkt)
    {
        errcode = errno;
        {
            static int i = 0;
            i++;
            if (i >= 30)
            {
                LOG_ERROR("kafka create topic err,topic[%s],partion[%d] errno %d, errinfo %s",
                          topic.c_str(), partition, errcode, gse::tools::error::ErrnoToStr(errcode).c_str());
                i = 0;
            }
        }

        return GSE_ERROR;
    }

    LOG_DEBUG("kafka key:%s value:%s", key.c_str(), value.c_str());
    res = rd_kafka_produce(rkt, RD_KAFKA_PARTITION_UA, RD_KAFKA_MSG_F_COPY, const_cast<char *>(value.c_str()), value.size(), key.c_str(), key.size(), beginTimestamp);

    if (res == -1)
    {
        // errcode will be one of reasons
        // ENOBUFS, EMSGSIZE, ESRCH, ENOENT
        struct timeval nowtv;
        errcode = errno;

        static int limit = 0;
        if (limit < 200)
        {
            LOG_ERROR("kafka produce msg failed, topic[%s], key [%s], partion[%d] errno [%d], errinfo [%s], auxiliary[%s], valuelen[%d]",
                      topic.c_str(), key.c_str(), partition, errcode, rd_kafka_err2str((rd_kafka_resp_err_t)res),
                      auxiliary.c_str(), value.size());
            gse::tools::time::GetTimeOfDay(&m_lastlogTime, NULL);
        }

        limit++;
        gse::tools::time::GetTimeOfDay(&nowtv, NULL);
        int diff = nowtv.tv_sec - m_lastlogTime.tv_sec;
        if (diff > 10)
        {
            limit = 0;
        }
        rd_kafka_topic_destroy(rkt);
        return errcode;
    }

    LOG_INFO("kafka produce success, topic:%s partition:%d queue cnt:%d ", topic.c_str(), partition, rd_kafka_outq_len(m_rdKafa));

    return GSE_SUCCESS;
}

int KafkaProducer::GetProducerQueueCount()
{
    int ret = -1;

    if (m_rdKafa)
    {
        ret = rd_kafka_outq_len(m_rdKafa);
    }

    return ret;
}

void KafkaProducer::CloseProducer()
{
    rd_kafka_destroy(m_rdKafa);
}

void KafkaProducer::KafkaPoll()
{
    rd_kafka_poll(m_rdKafa, 0);
}

} // namespace data
} // namespace gse
