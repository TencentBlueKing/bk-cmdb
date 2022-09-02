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

#ifndef _PULSAR_EXPORTER_H_
#define _PULSAR_EXPORTER_H_

#include <vector>

#include "exporter/exporter.h"
#include "datacell.h"
#include "pulsar_producer.h"
#include "filter/datafilter.h"
#include "conf/confItem.h"

namespace gse {
namespace data {

class PulsarExporter : public Exporter
{
public:
    PulsarExporter();
    virtual ~PulsarExporter();
public:
    int Start();
    int Stop();
    int Write(DataCell *pDataCell);

private:
    void clear();
    int createPulsarProducers(); 
    
    bool startWithChannelID(ChannelIdExporterConfig *ptrChannelIDConfig);
    bool startWithDataFlow(ExporterConf* ptrExporterConf);
    static void pulsarPoll(int fd, short what, void* v);
    void toPulsarTopics(const std::vector<std::string> &topicnames, std::vector<std::string> &newtopicnames);

private:
    std::vector<PulsarProducer*> m_pulsarPorducers;
    uint32_t m_nextProducerId;
    int32_t m_producerNum;
    std::string m_topicName;
    std::string m_serivce_url;
    std::string m_tlsTrustCertsFilePath;
    std::string m_tlsKeyFilePath;
    std::string m_tlsCertFilePath;
    std::string m_token;
};
}
}
#endif
