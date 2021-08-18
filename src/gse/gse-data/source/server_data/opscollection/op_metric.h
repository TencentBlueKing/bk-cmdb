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

#ifndef _GSE_DATA_OPS_METRIC_H_
#define _GSE_DATA_OPS_METRIC_H_

#include <rapidjson/document.h>		// rapidjson's DOM-style API
#include <rapidjson/prettywriter.h>	// for stringify JSON
#include <rapidjson/stringbuffer.h>



#include <vector>
#include <list>
#include <unordered_map>
#include <json/json.h>
#include <concurrentqueue.h>

#include "eventthread/gseEventThread.h"
#include "conf/dataconf.h"

#include "datacell.h"
#include "report_client.h"
//#include "dynamic_protocol/monitorTag.h"
//#include "ops/ops_client.h"

#include "bbx/prometheus/prometheus_exporter.h"
#include "bbx/prometheus/prometheus_handler.h"

//#include "gseConfig.h"


namespace gse {
namespace dataserver {


class OpsMetricMsg
{
public:
    std::string m_metricName;
    std::pair<std::string, uint64_t> m_metricItem;
};


struct NetMetricStat
{
    uint64_t net_bytes;
    uint64_t net_packets;
    uint64_t net_rbytes;
    uint64_t net_wbytes;
    uint64_t net_rpackets;
    uint64_t net_wpackets;
    NetMetricStat()
    {
        net_bytes = 0;
        net_packets = 0;
        net_rbytes = 0;
        net_wbytes = 0;
        net_rpackets = 0;
        net_wpackets = 0;
    }
};

struct DisKIoMetricStat
{
    uint64_t read_bytes;
    uint64_t write_bytes;
    DisKIoMetricStat()
    {
        read_bytes = 0;
        write_bytes = 0;
    }
};


struct IoMetricStat
{
    uint64_t net_bytes;
    uint64_t net_packets;
    uint64_t net_rbytes;
    uint64_t net_wbytes;
    uint64_t net_rpackets;
    uint64_t net_wpackets;
    IoMetricStat()
    {
        net_bytes = 0;
        net_packets = 0;
        net_rbytes = 0;
        net_wbytes = 0;
        net_rpackets = 0;
        net_wpackets = 0;
    }
};

class BaseMetric
{
public:
//    promethus::Counter m_cpu;
//    promethus::Counter m_mem;
//    promethus::Counter m_threadCount;
//    promethus::Counter m_diskIoReadBytes;
//    promethus::Counter m_diskIoReadBytes;
//    promethus::Counter m_diskIoWriteBytes;
//    promethus::Counter m_NetFlowRxBytes;
//    promethus::Counter m_NetFlowTxBytes;
//    promethus::Counter m_NetFlowRxPakcages;
//    promethus::Counter m_NetFlowTxPakcages;
//    promethus::Counter m_NetFlowTotalPakcages;
//    promethus::Counter m_NetFlowTotalBytes;
//    BaseMetric(prometheus::Registry* registry);
    void AddCounter();
    DataProcessConfig m_cfg;
};


class OPMetric
{

public:
    ~OPMetric();
    inline void SetGseConf(DataProcessConfig &cfg)
    {
        m_gseCfg = cfg;
    }

    static OPMetric& Instance();

public:
    void MetricEventInit();
public:
    static void DealData(evutil_socket_t fd, short what, void *args);
    static void DealMetricMsgRequest(evutil_socket_t fd, short what, void *args);
    static void ReportStats(evutil_socket_t fd, short what, void *args);
    static void BaseMetricStatTimer(evutil_socket_t fd, short what, void *args);
    //--------------------------------------------
    void ChannelIdMsgStat(std::string &channelid);
    void ChannelIdErrorMsgStat(std::string &channelid, uint32_t error_code);
    static long GetChannelidMsgCount(std::string &channelid);
    static long GetChannelidErrorMsgCount(std::string &channelid);
    //-------------------------------------------
    //--------------------------------
    static void UdpMsgInc();
    static void UdpInvalidMsgInc();
    static void TcpMsgInc();
    static void HttpMsgInc();
    static void ProxyMsgInc();
    static void ProxyUdpMsgInc();
    static void LogMsgInc();
    static void KafkaMsgInc();
    static void PulsarMsgInc();
    static void RedisMsgInc();
    static void ReleaseTCFreeMemory(evutil_socket_t fd, short what, void *args);

    static void ConnectionCountInc();

    static long GetUdpMsgCount();
    static long GetTcpMsgCount();
    static long GetHttpMsgCount();
    static long GetKafakMsgCount();
    static long GetRedisMsgCount();
    static long GetLogMsgCount();
    static long GetProxyMsgCount();
    static long GetPulsarMsgCount();

    void ReportQueneCount(std::string &metric_name, int thread_id, uint64_t queue_size);
    void ReportUdpMsgMetric(uint64_t value);
    void HandleMetricData();

private:
    OPMetric();
    void HandleMetricMsgTimerEvent(evutil_socket_t fd);
    void HandleMetricOpsMsg(OpsMetricMsg *ptr_ops);
    void ReportMetricStats();
    void SystemStats();
    bool EnQueneMetricMsg(OpsMetricMsg *ptr);

    //promethus--------------------------

    void BaseMetric();
    void MsgMetric();
    void cleanMsgStat();
    void channelidStatMetric();

    void SumMetricStatItem();
    int GetDiskIOByPid(DisKIoMetricStat &data, pid_t pid);
    int GetNetTrafficByPid(struct NetMetricStat &data, pid_t pid);

private:
    rgse::GseEventThread *m_ptrMetricEventThread;
    int m_opsMsgHandleEventId;
    DataProcessConfig m_gseCfg;
    NetMetricStat m_netMetric;
    DisKIoMetricStat m_diskIoMetric;
    double m_cpuUsage;
    moodycamel::ConcurrentQueue<OpsMetricMsg *>* m_ptrMetricQueue;
    //name <perthread_stat_name, value>
    std::string m_reportIP;
    std::map<std::string, std::map<std::string, uint64_t>> m_metrics;

    //prometheus::Registry m_metricRegistry;
    gse::bbx::PrometheusExporter *m_exporter;
};

}
}
#endif
