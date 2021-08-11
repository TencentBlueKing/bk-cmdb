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

#include <ctype.h>
#include <signal.h>
#include <string>
#include <list>
#include <sys/queue.h>
#include "eventthread/gseEventThread.h"
#include "conf/conf_common.h"

namespace gse {
namespace dataserver {
#ifdef __cplusplus
extern "C" {
#endif
#include <rdkafka.h>
#ifdef __cplusplus
}
#endif


using namespace std;

/**
 * @brief kafka producer
 *        生产者
 */
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
    int createProducer(const string &broker);
    int getProducerState();
    int excuteProduce(const string &topic, int partition, const std::string& key, const string& value, const std::string& auxiliary);
    void closeProducer();
    int getProducerQueueCount();

    void SetKafkaConfig(KafkaConfig &kafa_conf);
    void SetMaxMessageBytes(int bytes);
public:
   void KafkaPoll();
public:
    /**
     * @brief 地址集
     */
    string s_brokers;
    /**
     * @brief kafka 对象
     */
    rd_kafka_t *rk;

    rd_kafka_conf_t *conf;
    //rd_kafka_topic_conf_t *topic_conf;
public:
    static std::string m_runtimeDataDirector;

private:
    struct timeval m_lastlogTime;
    int32_t m_maxMessageBytes;
    KafkaConfig m_kafkaConfig;

};

}
}
#endif

