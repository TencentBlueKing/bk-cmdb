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

#ifndef _GSE_DATA_LOGEXPORTER_H_
#define _GSE_DATA_LOGEXPORTER_H_

#include <string>
#include <event.h>
#include "tools/thread.h"
#include "datastruct/safe_map.h"
#include "datastruct/safe_queue.h"
#include "safe/lock.h"
#include "exporter/exporter.h"
#include "datacell.h"
#include "log/log.h"
namespace gse { 
namespace data {

class LogExporter : public Exporter
{
public:
    typedef uint64_t                                LogSink;
    typedef std::string                        SourceID;
    typedef gse::datastruct::SafeMap<const std::string, LogSink>      SourceIDToLogSinkMap;
public:
    LogExporter();
    virtual ~LogExporter();

public:
    int Start();
    int Stop();
    int Write(DataCell *pDataCell);
private:
    bool startWithChannelID(ChannelIdExporterConfig *ptrChannelIDConfig);
    bool startWithDataFlow(ExporterConf* ptrExporterConf);

    LogSink createLogSink(const std::string& source_id);

private:
    int m_logSink;

private:
    bool                   m_needAllSourceIDS;
    LogSink                m_defaultLogSink;
    SourceIDToLogSinkMap   m_sourceIDToSink;
private:

    uint32_t     m_queueLimitMaxCount;
};

}
}
#endif
