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

#ifndef _GSE_DATA_REDISEXPORTER_H_
#define _GSE_DATA_REDISEXPORTER_H_

#include "exporter/exporter.h"
#include "datacell.h"
#include "redis_sentinel_publisher.h"
#include "redis_pub_producer.h"
namespace gse { 
namespace data {

class RedisExporter : public Exporter
{
public:
    RedisExporter();
    virtual ~RedisExporter();

public:
    int Start();
    int Stop();
    int Write(DataCell *pDataCell);

private:
    bool startWithChannelID(ChannelIdExporterConfig *ptrChannelIDConfig);
    bool startWithDataFlow(ExporterConf* ptrExporterConf);

private:
    RedisSentinelPublisher*   m_ptrSentinelPubliser;
    RedisPublishProducer*     m_ptrPubliser;
};

}
}
#endif
