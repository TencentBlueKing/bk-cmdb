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

#ifndef _GSE_DATA_KAFKA_EXPORTER_H_
#define _GSE_DATA_KAFKA_EXPORTER_H_

#include <vector>
#include "exporter/exporter.h"
#include "datacell.h"
#include "kafka_producer.h"
#include "dataStruct/safe_map.h"
#include "filter/datafilter.h"
#include "conf/confItem.h"
#include "eventthread/gseEventThread.h"
namespace gse { 
namespace dataserver {

class KafkaExporter : public Exporter
{
public:
    KafkaExporter();
    virtual ~KafkaExporter();
public:
    static void KafkaPoll(int fd, short what, void* v);
public:
    int Start();
    int Stop();
    int Write(DataCell *pDataCell);

private:
    void clear();
    int createKafkaProducers();

    bool startWithChannelID(ChannelIdExporterConfig *ptrChannelIDConfig);
    bool startWithDataFlow(ExporterConf* ptrExporterConf);
    bool startWithDataID(StorageConfigType* ptrStorageConfig);

private:
    std::vector<KafkaProducer*> m_vKafkaProducer;
    uint32_t m_nextProducerId;
    std::string m_kafkaBrokers;
    int32_t m_kafkaMaxQueue;
    int32_t m_kafkaMaxMessageBytes;
    int32_t m_producerNum;
    rgse::GseEventThread m_eventManager;
    std::string m_defaultTopicName;
    std::string m_token;

    std::string m_selfIp;
    KafkaConfig m_kafkaConfig;
};

}
}
#endif
