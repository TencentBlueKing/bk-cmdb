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

#ifndef _GSE_DATA_OPS_COLLECTION_H_
#define _GSE_DATA_OPS_COLLECTION_H_

#include <rapidjson/document.h>     // rapidjson's DOM-style API
#include <rapidjson/prettywriter.h> // for stringify JSON
#include <rapidjson/stringbuffer.h>

#include "concurrentqueue/concurrentqueue.h"

#include <event2/event.h>
#include <json/json.h>
#include <list>
#include <unordered_map>
#include <vector>

#include "conf/dataconf.h"
#include "datastruct/safe_map.h"
#include "safe/lock.h"
#include "tools/thread.h"

#include "datacell.h"
#include "eventthread/event_thread.h"
#include "report_client.h"
//#include "dynamic_protocol/monitorTag.h"
//#include "opscollection/ops_client.h"

//#include "ops/ops_client.h"
#include "conf/channel_id_config.h"
#include "conf/confItem.h"
#include "conf/dataconf.h"
#include "ops/op_metric.h"
#include "ops/ops_client.h"

namespace gse {
namespace data {

#define OPS_ERROR_DATAID 1
#define OPS_ERROR_CHANNELID 2
#define OPS_ERROR_DECODE 3
#define OPS_ERROR_PULSAR_ERROR 6
#define OPS_ERROR_TRANSFER_TCP_ERROR 7
#define OPS_ERROR_TRANSFER_UDP_ERROR 8

#define OPS_ERROR_KAFKA_ERROR 100
#define OPS_ERROR_KAFKA_NO_TOPIC_ERROR 101

#define OPS_ERROR_REDIS_ERROR 200

enum BKDataMonitorType
{
    BKDATA_MONITOR_INPUT = 1,
    BKDATA_MONITOR_OUTPT = 2,
    BKDATA_MONITOR_DROP = 3
};

/**
 * key value 结构定义
 */
typedef std::map<std::string, std::string> StringKeyValue;

/**
 * 普通的tag 结构定义
 */
typedef struct _CommonTag
{
    _CommonTag()
        : m_isValid(false), m_tag(""){};
    bool m_isValid;
    /**
     * StringKeyValue 结构
     */
    StringKeyValue m_desc;
    /**
     * desc 里面所有value 的拼接，由| 分割
     */
    std::string m_tag;
} CommonTag;

/**
 * info node struct definition
 */
typedef struct _NodeInfo
{
    _NodeInfo()
        : m_component(""), m_module(""){};
    /**
     * storm|ds|ts|agent|cmdline|...
     */
    std::string m_component;
    /**
     * module collector|realtime|...
     */
    std::string m_module;
    /**
     * custom tag
     */
    StringKeyValue m_customTag;
    /**
     * logical tag
     */
    CommonTag m_logicalTag;
    /**
     * physical tag ,indentify a host
     */
    CommonTag m_physicalTag;
} NodeInfo;

/**
 * stream struct definition
 */
typedef struct _Stream
{
    _Stream()
        : m_component(""), m_module(""){};
    /**
     * storm|ds|ts|agent|cmdline|...
     */
    std::string m_component;
    /**
     * module collector|realtime|...
     */
    std::string m_module;
    /**
     * logical tag definition
     */
    CommonTag m_logicalTag;
} Stream;

/**
 * location node
 */
typedef struct _NodeLocation
{
    /**
     * downstream, next node in the flow
     */
    std::list<Stream> m_downStream;
    /**
     * upstream, up node in the flow
     */
    std::list<Stream> m_upStream;
} NodeLocation;

/* metrics struct definition */

typedef struct _Input
{
    _Input()
        : m_isValid(false), m_totalCnt(0), m_incrementCnt(0), m_incrementBytes(0){};
    bool m_isValid;
    uint64_t m_totalCnt;     // total_cnt 共处理多少条输入数据，累计值，每次打点输出不清0
    uint64_t m_incrementCnt; // total_cnt_increment 最近1分钟共处理多少条输入数据
    uint64_t m_incrementBytes;
    std::map<string, uint64_t> m_tags; // <tag, count>, tags 按数据上的tag信息分别统计每种tag输入了多少条
} Input;

typedef Input Output;

typedef struct _DropObj
{
    _DropObj()
        : m_count(0), m_reason(""), m_errcode(0){};
    // 按抛弃的原因code分类，分别统计条数， 并把具体的原因写入reason字段
    uint64_t m_count;
    uint32_t m_errcode;
    std::string m_reason;
} DropObj;
typedef std::map<std::string, DropObj> DropImp;

typedef struct _DataLoss
{
    Input m_input;
    Output m_output;
    DropImp m_drop;
} DataLoss;

typedef struct _OutputEx
{
    _OutputEx()
        : m_isValid(false), m_totalCnt(0), m_incrementCnt(0), m_incrementBytes(0){};
    bool m_isValid;
    uint64_t m_totalCnt;     // total_cnt 共处理多少条输入数据，累计值，每次打点输出不清0
    uint64_t m_incrementCnt; // total_cnt_increment 最近1分钟共处理多少条输入数据
    uint64_t m_incrementBytes;
    std::map<string, uint64_t> m_tags; // <tag, count>, tags 按数据上的tag信息分别统计每种tag输入了多少条
    std::string m_type;
    std::string m_address;
} OutputEx;

typedef struct _MultiDataLoss
{
    Input *m_input;
    std::map<std::string, OutputEx *> m_output;
    std::map<std::string, DropObj *> m_drop;
    _MultiDataLoss()
    {
        m_input = NULL;
    }

    ~_MultiDataLoss()
    {
        if (m_input != NULL)
        {
            delete m_input;
            m_input = NULL;
        }

        for (std::map<std::string, OutputEx *>::iterator it = m_output.begin(); it != m_output.end(); it++)
        {
            if (it->second != NULL)
            {
                delete it->second;
            }
        }

        for (std::map<std::string, DropObj *>::iterator it = m_drop.begin(); it != m_drop.end(); it++)
        {
            if (it->second != NULL)
            {
                delete it->second;
            }
        }
    }
} MultiDataLoss;

typedef struct _DelayTime
{
    time_t m_local; //数据接收时间
    time_t m_recv;  //协议头时间
    int64_t m_delay;
    _DelayTime()
    {
        m_delay = 0;
        m_recv = 0;
        m_local = 0;
    }
} DelayTime;

typedef struct _DataDelay
{
    // windows time is 60, default 60s
    // waiting time not use now, default 0
    _DataDelay()
        : m_windowTime(60), m_waitingTime(0){};
    int m_windowTime;
    int m_waitingTime;
    DelayTime m_minDelay;
    DelayTime m_maxDelay; // min_delay and max_delay are the same now
} DataDelay;

typedef struct _DataMonitor
{
    DataLoss m_dataLoss;
    DataDelay m_dataDelay;
} DataMonitor;

class OPSMonitorCustomTags
{
public:
    std::unordered_map<std::string, std::string> m_customTags;
};

class Count
{
public:
    Count()
        : m_count(0), m_bytes(0) {}
    uint64_t m_count;
    uint64_t m_bytes;
};
class OpsMsgCount
{
public:
    Count m_input;
    Count m_output;
    Count m_drop;
};

class OpsDataStat
{
public:
    OpsDataStat() {}
    ~OpsDataStat()
    {
    }

    std::unordered_map<std::string, std::string> m_customTags;
    MultiDataLoss m_dataLoss;
    DataDelay m_delay;
    DataDelay m_outputDelay;
};

class OPSStats
{
public:
    OPSStats()
        : m_count(1){};
    ~OPSStats(){};

public:
    int m_count;
    std::string m_state;
    std::string m_protocol;
    std::string m_szChannelID;
    Json::Value m_jsonResultValue;
};

enum OpsMsgType
{
    EN_METRIC_TYPE = 0,
    EN_ACCOUNT_TYPE = 1,
};

class OpsStatWorker
{
public:
    OpsStatWorker();
    ~OpsStatWorker();
    inline void SetGseConf(std::shared_ptr<DataProcessConfig> configPtr)
    {
        m_configPtr = configPtr;
    }

    void SetOpsServerAddress(std::vector<OPSServerAddress> ops_servers);
    void SetOpsParam(std::string &report_ip, int channel_id);

private:
    DISALLOW_COPY_AND_ASSIGN(OpsStatWorker);

public:
    static void DealData(evutil_socket_t fd, short what, void *args);
    void reportToMonitor(evutil_socket_t fd);
    static void perMinReportMonitor(evutil_socket_t fd, short what, void *args);
    static void OpsQueueSizeReportEvent(evutil_socket_t fd, short what, void *args);

    int Start();
    bool PushOpsQueue(DataCellOPS *ptr);

private:
    void convertToPrometheus(evutil_socket_t fd);
    void OpsQueueSizeReport();
    void handleMonitorReport(evutil_socket_t fd);
    void HandleMetricMsg(evutil_socket_t fd);
    void updateStats(DataCellOPS *ptrDataCellOPS);
    void ReportMetricStats();
    void calcDataCount(DataCellOPS *ptrDataCellOPS);
    void calcDataAccumulationCount(DataCellOPS *ptrDataCellOPS);
    uint32_t toMinutes(time_t timestamp);
    void splitMonitorKey(uint64_t key, uint32_t &channelid, uint32_t &timestamp);
    void copyIncreasementCount(DataLoss *count, DataLoss *inc_count);
    uint64_t MakeMonitorReportKey(DataCellOPS *ptrDataCellOPS);
    void MakeKeyTag(DataCellOPS *ptr_ops, std::string &tag);
    void StatRecvPkg(DataCellOPS *ptr_ops);
    OpsMsgCount *GetTotalMsgStat(DataCellOPS *ptr_ops);
    void RecvTotalMsgStat(DataCellOPS *ptr_ops);
    void calcDelay(DataCellOPS *ptr_ops, OpsDataStat *dataloss);
    void calcOutputDelay(DataCellOPS *ptr_ops, OpsDataStat *dataloss);
    void CountMsg(DataCellOPS *ptr_ops, OpsDataStat *dataloss);
    void CalcRecvCount(DataCellOPS *ptr_ops, MultiDataLoss *count);
    void CalcSendAccount(DataCellOPS *ptr_ops, MultiDataLoss *count);
    void CalcDropCount(DataCellOPS *ptr_ops, MultiDataLoss *count);
    void ReportAllChannelIdStat();
    void ConvertToReportJsonString(std::string &key, OpsDataStat *ops_data_stat, std::string &report_json);
    void makeLogicalTag(std::string &channleid, rapidjson::Document &logical_tag);
    void makeCustomTagJson(std::unordered_map<std::string, std::string> &custom_tags, rapidjson::Document &custom_tag);
    void makeDataMonitorJson(OpsDataStat *ptr_data_loss, rapidjson::Document &data_monitor_json);
    void makeNodeInfoJson(rapidjson::Document &node_info);
    void makePhysicalTagJson(rapidjson::Document &physical_tag_json);
    void MakeCustomTag(DataCellOPS *ptr_ops, OpsDataStat *dataloss, std::string &key);

private:
    int m_threadCount;

    EventThread *m_ptrEventThread;
    std::vector<OPSServerAddress> m_opsServerAddress;
    int m_opsMsgHandleEventId;
    ReportClient *m_ptrReportClient;
    OPSClient m_opsClient;
    bool m_hasOpsServer;
    int m_channelID;
    int m_agentOpsChannelID;
    std::string m_customTagName;
    std::string m_reportIP;
    std::map<std::string, OPSStats *> m_opsStats;

    std::shared_ptr<DataProcessConfig> m_configPtr;

    NodeInfo m_nodeInfo;
    typedef map<uint64_t, OPSMonitorCustomTags> ReportMonitorCustomTagMap;
    typedef map<uint64_t, DataLoss> ReportMonitorMap;
    typedef map<uint64_t, DataLoss> ReportMonitorSumMap;

    ReportMonitorSumMap m_channelIdMap;
    ReportMonitorMap m_reportMonitorMap;
    std::unordered_map<std::string, OpsDataStat *> m_allMsgStat;

    std::unordered_map<std::string, OpsMsgCount *> m_allMsgTotalStat;
    rapidjson::Document m_nodeInfoJson;
    std::unordered_map<uint64_t, OPSMonitorCustomTags *> m_reportMonitorCustomTagMap;
    bool m_enableOps;
    moodycamel::ConcurrentQueue<DataCellOPS *> *m_ptrOpsQueue;
};

class OpsCollection
{

public:
    OpsCollection()
    {
        m_threadCount = 8;
        m_ptrReportClient = NULL;
        m_hasOpsServer = false;
        m_channelID = -1;
        m_agentOpsChannelID = -1;
        m_enableOps = false;
    };
    ~OpsCollection(){};

    inline void SetGseConf(std::shared_ptr<DataProcessConfig> &configPtr)
    {
        m_configPtr = configPtr;
    }

public:
    int Open(int maxQueueCnt, int thread_count, bool enableops, std::string &reportIP, int channelID, std::vector<OpsAddress> &server_addresses);
    void Close();
    int PutOpsData(DataCellOPS *ptrData);
    int TransferOpsData(DataCell *ptrData);
    void SetAgentOpsChannelid(int channelid);

private:
    DISALLOW_COPY_AND_ASSIGN(OpsCollection);

private:
    bool PushOpsQueue(DataCellOPS *ptr);

private:
    int m_threadCount;
    std::shared_ptr<DataProcessConfig> m_configPtr;
    ReportClient *m_ptrReportClient;
    OPSClient m_opsClient;

    bool m_hasOpsServer;
    int m_channelID;
    int m_agentOpsChannelID;
    std::string m_customTagName;
    std::string m_reportIP;

    std::vector<OpsStatWorker *> m_opsWorkers;
    bool m_enableOps;
};
} // namespace data
} // namespace gse
#endif
