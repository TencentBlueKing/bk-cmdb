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

#ifndef __GSE_OPS_DEAL_H__
#define __GSE_OPS_DEAL_H__

#include <string>
#include <json/json.h>
#include "utilTools.h"
#include "gseEventThread.h"
#include "monitorTag.h"
namespace gse { 
namespace dataserver {

/**
 * @ brief 运营数据对象
 *
 */
typedef struct _OpsDataObject
{
//    const char* m_name;
//    int m_cloudid;
//    int m_bizid;
//    uint32_t m_type;
    int m_dataId;
    string m_ctx;
//    string m_clientIp;
} OpsDataObject;
typedef std::string TagKey;
typedef std::map<TagKey, int> MonitorTag;
struct OpsCount
{
    size_t m_count;
    size_t m_size;
    std::string m_dsip;
};

class OpsDeal
{
public:
    typedef map<int, OpsCount> StatusSizeMap;
    typedef map<string, StatusSizeMap> IpStatusMap;
    typedef map<int, IpStatusMap> DataIdIpMap;
    typedef map<string, uint32_t> IpBizidMap;

public:
    OpsDeal();
    ~OpsDeal();

public:
    static void opsDealEventHandler(int fd, short what, void* args);
    static void reportDataMsgHandler(int fd, short what, void* args);

public:

    int setup();
    int syncWait();
    void pushData(OpsDataObject* ptrDataObject);
    int stop();

    void updateOpDataId(int opType, int dataId);
    int findDataIdByOpType(int opType);
    void deleteOpDataId(int dataId);

private:
    int init();
    void reportDataFlow();
    void dataFlowCount(Json::Value &dataValue);
    void sendTo(int opType, const string &msg);

    // monitor tag methods
    bool monitorTagFlowCount(Json::Value &dataValue);
    void calcRecvCount(DataLoss *count, const string &tag);
    void calcSendCount(DataLoss *count, const string &tag);
    void calcDropCount(DataLoss *count, const string &reason);
    void reportMonitorTag();
    void updateTimeTag();

private:
    rgse::GseEventThread m_gseEventThread;
    int m_notifyFD[2];
    IpBizidMap m_ipBizidMap;
    DataIdIpMap m_dataFlows;
    map<int,int> m_opType2DataIdMap;
    gse::safe::RWLock m_opType2DataIdMapLock;

    string m_timeTag;
    string m_localIp; // local ip
    NodeInfo m_nodeInfo;
    typedef map<int, DataLoss> DataIdMap;
    DataIdMap m_dataMonitorMap;
};

# endif

}
}
