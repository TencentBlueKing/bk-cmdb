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

#include <list>
#include <memory>
#include <mutex> // std::mutex
#include <vector>

#include <rapidjson/document.h>     // rapidjson's DOM-style API
#include <rapidjson/prettywriter.h> // for stringify JSON
#include <rapidjson/stringbuffer.h>

#include <prometheus/counter.h>
#include <prometheus/gauge.h>
#include <prometheus/registry.h>

#include "concurrentqueue/concurrentqueue.h"

#include <json/json.h>
#include <unordered_map>

#include "bbx/prometheus/prometheus_exporter.h"
#include "bbx/prometheus/prometheus_handler.h"
#include "conf/dataconf.h"
#include "datacell.h"
#include "eventthread/event_thread.h"
#include "net/http/http_server.hpp"
#include "report_client.h"
#include "tools/finally.hpp"

namespace gse {
namespace data {

enum EnMessageCodecErrorCode
{
    MESSAGE_CODEC_SUCCESS = 0,
    MESSAGE_CODEC_BODYLEN_ERROR = 1,
    MESSAGE_CODEC_MAGIC_NUMBER_ERROR,
    MESSAGE_CODEC_CHANNELID_NOT_SET_ERROR,
    MESSAGE_CODEC_DYNAMIC_PROTOCOL_PARSE_ERROR,
    MESSAGE_CODEC_TDM_PROTOCOL_SOURCE_ID_ERROR,
    MESSAGE_CODEC_TDM_PROTOCOL_EVENT_NAME_ERROR,
    MESSAGE_CODEC_TDM_PROTOCOL_BUSINESS_ID_ERROR,

};

static const std::string kRecvMessageDirection = "recv";
static const std::string kDispathMessageDirection = "dispath";
static const std::string kSendMessageDirection = "send";

static const std::string kDefaultMsgType = "0";

static const std::string kOutputKafka = "kafka";
static const std::string kOutputPulsar = "pulsar";
static const std::string kOutputRedis = "redis";
static const std::string kOutputTCP = "tcp";
static const std::string kOutputHTTP = "http";
static const std::string kOutputUDP = "udp";

static const std::string kCodecTDM = "tdm_codec";
static const std::string kCodecV1 = "codec_v1";
static const std::string kCodecTGLOG = "codec_tglog";
static const std::string kCodecV2 = "codec_v2";
static const std::string kGenneralCodec = "genneral_codec";

#define TIME_TO_SECOND(timeMS) (timeMS / 1000)
#define SECOND_TO_MS(timeMS) (timeMS * 1000)

#define TIME_DEALY_SECOND(beginTime, delay)                                     \
    int ____LOCAL_TIME____ = TIME_TO_SECOND(gse::tools::time::GetUTCMillsec()); \
    delay = beginTime > ____LOCAL_TIME____ ? beginTime - ____LOCAL_TIME____ : 0;

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

class Metrics
{
public:
    virtual void InitBaseMetric() = 0;
    virtual void ReportMetrics() = 0;
    virtual void CollectMetrics() = 0;
    virtual void SetConfig(std::shared_ptr<DataProcessConfig> configPtr) = 0;
    virtual std::shared_ptr<prometheus::Registry> GetRegistry() = 0;
};

class ServiceModeReportMetrics : public Metrics
{
public:
    ServiceModeReportMetrics();
    ~ServiceModeReportMetrics();
    void InitBaseMetric();
    void ReportMetrics();
    void CollectMetrics();
    void HandleReportMetricRequest(OpsMetricMsg *ptr_ops);
    void SetConfig(std::shared_ptr<DataProcessConfig> configPtr);
    std::shared_ptr<prometheus::Registry> GetRegistry();

private:
    std::map<std::string, std::map<std::string, uint64_t>> m_metrics;
    std::shared_ptr<prometheus::Registry> m_registry;
    std::map<std::string, std::string> m_labels;
    std::shared_ptr<DataProcessConfig> m_configPtr;
};

class MessageMetric : public Metrics
{
public:
    MessageMetric();
    virtual ~MessageMetric();
    void InitBaseMetric();
    void ReportMetrics();
    void CollectMetrics();
    void SetConfig(std::shared_ptr<DataProcessConfig> configPtr);
    void Clean();
    std::shared_ptr<prometheus::Registry> GetRegistry();

private:
    void MessageBytesMetric();

private:
    std::shared_ptr<DataProcessConfig> m_configPtr;
    std::map<std::string, std::string> m_labels;
    std::shared_ptr<prometheus::Registry> m_registry;
};

class BaseMetric : public Metrics
{
public:
    BaseMetric();
    ~BaseMetric();
    void InitBaseMetric();
    void ReportMetrics();
    void CollectMetrics();
    void SystemStats();

    std::shared_ptr<prometheus::Registry> GetRegistry();

    void SetConfig(std::shared_ptr<DataProcessConfig> configPtr);
    void Clean();

private:
    int GetNetTrafficByPid(struct NetMetricStat &data, pid_t pid);
    int GetDiskIOByPid(struct DisKIoMetricStat &data, pid_t pid);
    int GetProcessThreadCount(pid_t pid);

private:
    std::shared_ptr<prometheus::Registry> m_registry;
    std::map<std::string, std::string> m_labels;

    time_t m_startupTimeSecond;
    NetMetricStat m_netMetric;
    DisKIoMetricStat m_diskIoMetric;
    double m_cpuUsage;
    std::shared_ptr<DataProcessConfig> m_configPtr;
};

class OPMetric
{

public:
    virtual ~OPMetric();
    inline void SetGseConf(DataProcessConfig &cfg)
    {
        m_gseCfg = cfg;
    }

    static OPMetric &Instance();

    static OPMetric *OPMetricInst();
    std::shared_ptr<gse::bbx::PrometheusExporter> GetPromethusExporter();

public:
    void MetricEventInit(std::shared_ptr<net::http::HTTPServer> htpServer);

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
    static void UdpMsgInc(int bytes);
    static void UdpInvalidMsgInc();
    static void HttpMsgInc(int bytes);
    static void LogMsgInc();

    static void AccessRequestInc();
    static void AccessResponseInc();
    static void AccessResponseFailedInc();

    static void ReleaseTCFreeMemory(evutil_socket_t fd, short what, void *args);

    static void ConnectionCountInc(const std::string &connectionId);
    static void ConnectionCountDec(const std::string &connectionId);
    static size_t GetConnectionCount();

    static long GetUdpMsgCount();
    static long GetTcpMsgCount();
    static long GetHttpMsgCount();
    static long GetLogMsgCount();

    static void AddRecvMsgCounter(const std::string &msgType, const std::string &codec, int error, int count);
    static void AddRecvMsgBytesCounter(const std::string &msgType, const std::string &codec, int bytes);
    static void AddMsgDelayMsCounter(const std::string &msgType, const std::string &direction, const std::string &codec, int delay);

    static void AddSendMsgBytesCounter(const std::string &msgType, const std::string &outputType, int bytes);
    static void AddSendMsgCounter(const std::string &msgType, const std::string &outputType, int error, int count);

    static void AddMessageQueneRequestCounter(const std::string &queneType, const std::string &direction, int errorCode, int count);

    static void AddMessageQueneRequestMillsecondsMetrics(const std::string &queneType, int count);
    static void AddChannelIdMetrics(const std::string &type, int count);
    static void AddChannelIdUsageMemoryMetrics(const std::string &type, int count);

    void ReportQueneCount(std::string &metric_name, int thread_id, uint64_t queue_size);
    void ReportUdpMsgMetric(uint64_t value);
    void HandleMetricData();
    void CollectMetrics();

private:
    OPMetric();
    void HandleMetricMsgTimerEvent(evutil_socket_t fd);
    void ReportMetricStats();
    bool EnQueneMetricMsg(OpsMetricMsg *ptr);
    void HandleMetricOpsMsg(OpsMetricMsg *ptr_ops);
    // prometheus--------------------------
    int StartMetricExporter(std::shared_ptr<net::http::HTTPServer> htpServer);
    void cleanMsgStat();
    void channelidStatMetric();

    void StartMetrics();

private:
    std::shared_ptr<EventThread> m_ptrMetricEventThread;
    DataProcessConfig m_gseCfg;

    std::shared_ptr<moodycamel::ConcurrentQueue<OpsMetricMsg *>> m_ptrMetricQueue;
    std::string m_reportIP;
    std::list<std::shared_ptr<Metrics>> m_metrics;
    std::shared_ptr<ServiceModeReportMetrics> m_serviceMetrics;
    std::shared_ptr<std::thread> m_metricExporterThread;

    std::shared_ptr<gse::bbx::PrometheusExporter> m_promethusExporter;
};

class OPSHttpServer
{
public:
    OPSHttpServer();
    OPSHttpServer(const std::string &ip, uint16_t port);
    ~OPSHttpServer();
    bool RegisterMsgHandler(const string &URI, std::shared_ptr<gse::net::http::HTTPHandler> handler);
    int Start();
    void Stop();
    int Join();

    std::shared_ptr<gse::net::http::HTTPServer> GetHttpServer();

private:
    int StartMigrationSerivce();

private:
    std::shared_ptr<gse::net::http::HTTPServer> m_httpServer;
    std::unique_ptr<gse::net::MigrationClient> m_migrationClient;
    std::unique_ptr<gse::net::MigrationServer> m_migrationServer;
    std::thread m_listenThread;

    int m_listennerFd;
    std::string m_serverIp;
    uint16_t m_port;
};

} // namespace data
} // namespace gse
#endif
