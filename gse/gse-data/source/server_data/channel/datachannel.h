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

#ifndef _GSE_DATA_DATACHANNEL_H_
#define _GSE_DATA_DATACHANNEL_H_

#include "datastruct/safe_map.h"
#include "datastruct/safe_queue.h"
#include "safe/lock.h"
#include <map>
#include <string>
#include <vector>

#include "channelworker.h"
#include "codec/codec.h"
#include "conf/confItem.h"
#include "conf/configurator.h"
#include "conf/dataconf.h"
#include "datacell.h"
#include "eventthread/event_thread.h"
#include "exporter/exporter.h"
#include "receiver/receiver.h"
namespace gse {
namespace data {
using namespace std;

class DataChannel
{
public:
    typedef int StorageIndex;
    typedef std::string StorageName;
    typedef uint32_t ChannelIDType;
    typedef uint32_t DataIDType;
    typedef gse::datastruct::SafeMap<StorageIndex, Exporter *> StorageIndexToExporterMap;
    typedef gse::datastruct::SafeMap<StorageName, Exporter *> StorageNameToExporterMap;
    typedef gse::datastruct::SafeMap<DataIDType, Exporter *> DataIDToExporterMap;
    typedef gse::datastruct::SafeMap<ChannelIDType, Exporter *> ChannelIDToExporterMap;

public:
    DataChannel();
    virtual ~DataChannel();

public:
    inline void SetGseConf(std::shared_ptr<DataProcessConfig> cfg)
    {
        m_configPtr = cfg;
    }

    void SetConftor(Configurator *conftor)
    {
        m_conftor = conftor;
    }

public:
    int Start();
    int Stop();
    void Join();
    void UpdateConf(const ChannelConf &channelConf);
    int DecodeMsg(DataCell *pDataCell);
    int ExportData(DataCell *pDataCell);

    void SetOps(OpsCollection *ptr_ops_report);

    void SetDataCellChannelidByServiceId(DataCell *pDataCell);

private:
    static void cleanExporter(evutil_socket_t fd, short what, void *args);
    static void handlRecvDataCell(DataCell *pDataCell, void *pCaller);
    void localHanlerRecvDataCell(DataCell *pDataCell);

private:
    int init();
    int StartChannelWorker();
    int StopChannelWorker();
    int startReceiver();
    int StopReceiver();

    int directExport(DataCell *pDataCell);
    int dispatchChannelWorker(DataCell *pDataCell);

private:
private:
    int CreateDataflowExporter();
    void ClearDataflowExporter();

private:
    void dataflowExporterWrite(DataCell *pDataCell);

private:
    std::shared_ptr<DataProcessConfig> m_configPtr;

private:
    std::string m_channelName;
    std::string m_receiverName;
    gse::safe::RWLock m_upConfLock;
    ChannelConf *m_channelConf;
    bool m_isKeepTiming;
    uint32_t m_nextWorkerId;
    uint32_t m_numWorker;
    std::vector<ChannelWorker *> m_vChannelWorkers;
    Receiver *m_receiver;
    Codec *m_codec;
    OpsCollection *m_opsReport;
    gse::datastruct::SafeQueue<Exporter *> m_toDeletedExporter;

    Configurator *m_conftor;

private:
    std::vector<Exporter *> m_vDataflowExporter;

    // 以下两个变量用于兼容V1.0 的dataid 相关的配置信息
    StorageIndexToExporterMap m_storageIndexToExporter;
    ChannelIDToExporterMap m_channelIDToExporter;

    EventThread *m_ptrEventThread;

private:
    IDToStorage *m_ptrIDToStorage;
};

} // namespace data
} // namespace gse

#endif
