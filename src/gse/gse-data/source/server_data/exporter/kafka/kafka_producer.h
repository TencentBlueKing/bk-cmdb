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

#ifndef _GSE_KAFKA_PRODUCER_H_
#define _GSE_KAFKA_PRODUCER_H_

#include "conf/conf_common.h"
#include "eventthread/event_thread.h"
#include <ctype.h>
#include <list>
#include <signal.h>
#include <string>
#include <sys/queue.h>

namespace gse {
namespace data {
#ifdef __cplusplus
extern "C" {
#endif
#include <librdkafka/rdkafka.h>
#ifdef __cplusplus
}
#endif

using namespace std;

class KafkaProducer
{

public:
    enum KFKA_P_STATE
    {
        KAFKA_P_DOWN = 1,
        KAFKA_P_CONNECT,
        KAFKA_P_UP
    };

public:
    KafkaProducer();
    ~KafkaProducer();

public:
    int CreateProducer(const string &broker);
    int ExcuteProduce(const string &topic, int partition, const std::string &key, const string &value, const std::string &auxiliary);
    void CloseProducer();
    int GetProducerQueueCount();

    void SetKafkaConfig(KafkaConfig &kafa_conf);
    void SetMaxMessageBytes(int bytes);

public:
    void KafkaPoll();

private:
    static void MsgDeliverCb(rd_kafka_t *rk,
                             const rd_kafka_message_t *rkmessage, void *opaque);

public:
    static std::string m_runtimeDataDirector;

private:
    string m_brokers;
    rd_kafka_t *m_rdKafa;
    rd_kafka_conf_t *conf;

    struct timeval m_lastlogTime;
    int32_t m_maxMessageBytes;
    KafkaConfig m_kafkaConfig;
};

} // namespace data
} // namespace gse
#endif
