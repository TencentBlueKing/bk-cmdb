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

#include <vector>
#include <string>
#include <map>
#include "dataStruct/safe_map.h"
#include "dataStruct/safe_queue.h"
#include "safe/lock.h"

#include "conf/confItem.h"
#include "conf/dataconf.h"
#include "conf/configurator.h"
#include "channelworker.h"
#include "datacell.h"
#include "receiver/receiver.h"
#include "exporter/exporter.h"
#include "codec/codec.h"
#include "eventthread/gseEventThread.h"
namespace gse { 
namespace dataserver {
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
    inline void SetGseConf(DataProcessConfig &cfg)
    {
        m_gseCfg = cfg;
    }

    void SetConftor(Configurator* conftor)
    {
        m_conftor = conftor;
    }

public:
    int Start();
    int Stop();
    void Join();
    void UpdateConf(const ChannelConf &channelConf);
    void UpdateIDToStorage(IDToStorage *ptrIDToStorage);
    int DecodeMsg(DataCell *pDataCell);
    int ExportData(DataCell *pDataCell);

    void SetOps(OpsCollection* ptr_ops_report);
    void stopExporter(uint32_t channelid);
private:
    static void cleanExporter(evutil_socket_t fd, short what, void *args);
    static void handleConfigUpdateEvent(void *args, int storageIndex, uint32_t channelID);
    static void handlRecvDataCell(DataCell *pDataCell, void *pCaller);
    void localHanlerRecvDataCell(DataCell *pDataCell);
private:
    int init();
    int startChannelWorker();
    int stopChannelWorker();
    int startReceiver();
    int stopReceiver();

    int directExport(DataCell *pDataCell);
    int dispatchChannelWorker(DataCell *pDataCell);

private:
    // common exporter manager methods
    void createExporter();
    Exporter *createDataIDExporter(StorageConfigVector *ptrStorage);
    void updateDataIDExporter(StorageConfigVector *ptrStorage);
    Exporter *createChannelIDExporter(ChannelIDStorage *ptrStorage);
    void updateChannelIDExporter(ChannelIDStorage *ptrStorage);

private:
    int createDataflowExporter();
    void clearDataflowExporter();

private:
    void dataflowExporterWrite(DataCell *pDataCell);
    

private:
    DataProcessConfig m_gseCfg;
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


    Configurator* m_conftor;

private:
    std::vector<Exporter *> m_vDataflowExporter;

    // 以下两个变量用于兼容V1.0 的dataid 相关的配置信息
    StorageIndexToExporterMap m_storageIndexToExporter;
    ChannelIDToExporterMap m_channelIDToExporter;

    rgse::GseEventThread *m_ptrEventThread;

private:
    IDToStorage *m_ptrIDToStorage;
};

}
}

#endif
