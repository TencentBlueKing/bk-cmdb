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

#include "op_metric.h"

#include <atomic>
#include <memory>
#include <unordered_map>
#include <unordered_set>
#include <vector>

#include <gperftools/malloc_extension.h>

#include "bbx/gse_errno.h"
#include "common/time.h"
#include "dataserver.h"
#include "log/log.h"
#include "time_center.h"
#include "tools/hostinfo.h"
#include "tools/process.h"

namespace gse {
namespace data {

// metrics
//---------------------------------------------------------------------
const std::string data_cpu = "process_cpu_usage_ratio";
const std::string data_mem = "process_mem_usage";
const std::string data_mem_bytes = "process_mem_usage_bytes_total";
const std::string name_process_run_seconds_total = "process_run_seconds_total";

const std::string data_fd_count = "process_open_fds";
const std::string data_thread = "process_thread_total";
const std::string data_mem_usage_bytes_total = "process_mem_usage_bytes_total";
const std::string name_tcp_connection_count = "process_connection_count"; // counter

const std::string data_disk_io_read_bytes = "process_disk_io_read_bytess_count";
const std::string data_disk_io_write_bytes = "process_disk_io_write_bytes_count";
const std::string data_net_flow_tx_bytes = "process_net_flow_tx_bytes_count";
const std::string data_net_flow_rx_bytes = "process_net_flow_rx_bytes_count";
const std::string data_net_flow_rx_pakcages = "process_net_flow_rx_pakcages_count";
const std::string data_net_flow_tx_pakcages = "process_net_flow_tx_pakcages_count";
const std::string data_net_flow_total_pakcages = "process_net_flow_total_pakcages_count";
const std::string data_net_flow_total_bytes = "process_net_flow_total_bytes_count";

const std::string name_total_udp_msg_count = "recv_udp_msg_count";
const std::string name_total_http_msg_count = "recv_http_msg_count";
const std::string name_total_quene_msg_count = "data_quene_msg_count";

const std::string name_total_proxy_msg_count = "send_tcp_proxy_msg_count";
const std::string name_total_udp_proxy_msg_count = "send_udp_proxy_msg_count";
const std::string name_total_http_proxy_msg_count = "send_http_proxy_msg_count";
const std::string name_total_http_proxy_failed_msg_count = "send_http_proxy_failed_msg_count";

const std::string name_total_log_msg_count = "write_log_count";

const std::string name_total_request_msg_count = "request_msg_count";

//-----------bytes---------------------
const std::string name_total_udp_recv_msg_bytets = "data_total_udp_recv_msg_bytets";
const std::string name_total_http_recv_msg_bytets = "data_total_http_recv_msg_bytets";

const std::string name_channel_id_count = "channel_id_count";
const std::string name_channel_id_memusage = "channel_id_memusage";
const std::string name_stream_to_id_count = "stream_to_id_count";

const std::string name_recv_msg_metrics = "recv_msg_counter";
const std::string name_recv_msg_bytes_metrics = "recv_msg_bytes_counter";
const std::string name_send_msg_count = "send_msg_counter";
const std::string name_send_msg_bytes_metrics = "send_msg_bytes_counter";

const std::string name_msg_delay_metrics = "msg_delay_time_ms_counter";

const std::string name_message_quene_request_metrics = "message_queue_request_counter";
const std::string name_message_quene_request_millseconds_metrics = "message_queue_request_millseconds";

//-----------------------------------------------------------------

static std::unordered_map<std::string, unsigned long> channel_id_msg_stat;
static std::unordered_map<std::string, unsigned long> channel_id_error_msg_stat;

static std::mutex kMetricMutex;
static std::unordered_set<std::string> kConnectionList;

// end channelid stat------------
const std::string kafka_connection_count = "kafak_connection_count";

//--------------------------------------------------------
static std::atomic<unsigned long> totalUdpMsgCount(0);
static std::atomic<unsigned long> totalUdpErrorMsgCount(0);
static std::atomic<unsigned long> totalHttpMsgCount(0);
static std::atomic<unsigned long> totalUdpMsgBytes(0);
static std::atomic<unsigned long> totalHttpMsgBytes(0);
//-------------------------------------------------------

static std::atomic<unsigned long> totalLogMsgCount(0);
//------------------------------------------------------

// request response msg
static std::atomic<unsigned long> totalAccessRequestMsgCount(0);
static std::atomic<unsigned long> totalAccessResponseMsgCount(0);
static std::atomic<unsigned long> totalAccessResponseFailedMsgCount(0);
static std::atomic<unsigned long> totalConnectionCount(0);

static long udp_msg_count;
static long tcp_msg_count;
static long http_msg_count;
static long log_msg_count;

//全局counter
static prometheus::Family<prometheus::Counter> *kTcpMsgMsgCounterFamily;
static prometheus::Family<prometheus::Counter> *kRecvMsgBytesCounterFamily;
static prometheus::Family<prometheus::Counter> *kMsgDelayCounterFamily;

static prometheus::Family<prometheus::Counter> *kSendMsgBytesCounterFamily;
static prometheus::Family<prometheus::Counter> *kSendMsgCounterFamily;

static prometheus::Family<prometheus::Counter> *kMessageQueueRequestCounterFamily;
static prometheus::Family<prometheus::Counter> *kMessageQueueRequestMillsecondsFamily;
static prometheus::Family<prometheus::Gauge> *kChannelIdCounterFamily;
static prometheus::Family<prometheus::Gauge> *kChannelIdMemoryUsageFamily;

void OPMetric::AddRecvMsgCounter(const std::string &msgType, const std::string &codec, int error, int count)
{
    std::map<std::string, std::string> codeLabel = {{"message_type", msgType},
                                                    {"codec", codec},
                                                    {"error", gse::tools::strings::ToString(error)}};

    if (kTcpMsgMsgCounterFamily != nullptr)
    {
        auto &counter = kTcpMsgMsgCounterFamily->Add(codeLabel);
        counter.Increment(count);
    }
}

void OPMetric::AddRecvMsgBytesCounter(const std::string &msgType, const std::string &codec, int bytes)
{
    std::map<std::string, std::string> codeLabel = {{"message_type", msgType},
                                                    {"codec", codec}};

    if (kRecvMsgBytesCounterFamily != nullptr)
    {
        auto &counter = kRecvMsgBytesCounterFamily->Add(codeLabel);
        counter.Increment(bytes);
    }
}

void OPMetric::AddSendMsgBytesCounter(const std::string &msgType, const std::string &outputType, int bytes)
{
    std::map<std::string, std::string> codeLabel = {{"message_type", msgType},
                                                    {"output_type", outputType}};

    if (kSendMsgBytesCounterFamily != nullptr)
    {
        auto &counter = kSendMsgBytesCounterFamily->Add(codeLabel);
        counter.Increment(bytes);
    }
}

void OPMetric::AddSendMsgCounter(const std::string &msgType, const std::string &outputType, int error, int count)
{
    std::map<std::string, std::string> codeLabel = {{"message_type", msgType},
                                                    {"output_type", outputType},
                                                    {"error", gse::tools::strings::ToString(error)}};

    if (kSendMsgCounterFamily != nullptr)
    {
        auto &counter = kSendMsgCounterFamily->Add(codeLabel);
        counter.Increment(count);
    }
}

void OPMetric::AddMsgDelayMsCounter(const std::string &msgType, const std::string &direction, const std::string &codec, int delay)
{
    std::map<std::string, std::string> codeLabel = {{"message_type", msgType},
                                                    {"codec", codec},
                                                    {"direction", direction}};

    if (kMsgDelayCounterFamily != nullptr)
    {
        auto &counter = kMsgDelayCounterFamily->Add(codeLabel);
        if (delay > 0)
        {
            counter.Increment(delay);
        }
    }
}

void OPMetric::AddMessageQueneRequestCounter(const std::string &queneType, const string &direction, int errorCode, int count)
{
    std::map<std::string, std::string> codeLabel = {{"quene_type", queneType},
                                                    {"error_code", gse::tools::strings::ToString(errorCode)},
                                                    {"direction", direction}};

    if (kMessageQueueRequestCounterFamily != nullptr)
    {
        auto &counter = kMessageQueueRequestCounterFamily->Add(codeLabel);
        counter.Increment(count);
    }
}

void OPMetric::AddMessageQueneRequestMillsecondsMetrics(const std::string &queneType, int count)
{
    std::map<std::string, std::string> codeLabel = {{"output_type", queneType}};

    if (kMessageQueueRequestMillsecondsFamily != nullptr)
    {
        auto &counter = kMessageQueueRequestMillsecondsFamily->Add(codeLabel);
        counter.Increment(count);
    }
}

void OPMetric::AddChannelIdMetrics(const std::string &type, int count)
{
    std::map<std::string, std::string> codeLabel = {{"type", type}};

    if (kChannelIdCounterFamily != nullptr)
    {
        auto &gauge = kChannelIdCounterFamily->Add(codeLabel);
        if (count >= 0)
        {
            gauge.Increment(count);
        }
        else
        {
            gauge.Decrement(-1 * count);
        }
    }
}

void OPMetric::AddChannelIdUsageMemoryMetrics(const std::string &type, int count)
{
    std::map<std::string, std::string> codeLabel = {{"type", type}};

    if (kChannelIdMemoryUsageFamily != nullptr)
    {
        auto &gauge = kChannelIdMemoryUsageFamily->Add(codeLabel);
        if (count >= 0)
        {
            gauge.Increment(count);
        }
        else
        {
            gauge.Decrement(-1 * count);
        }
    }
}

long OPMetric::GetUdpMsgCount()
{
    return udp_msg_count;
}
long OPMetric::GetTcpMsgCount()
{
    return tcp_msg_count;
}
long OPMetric::GetHttpMsgCount()
{
    return http_msg_count;
}

long OPMetric::GetLogMsgCount()
{
    return log_msg_count;
}

void OPMetric::ChannelIdMsgStat(std::string &channelid)
{
    std::string name = "channelid_" + channelid;
    if (channel_id_msg_stat.find(name) == channel_id_msg_stat.end())
    {
        channel_id_msg_stat.insert(std::make_pair(name, 1));
    }
    else
    {
        channel_id_msg_stat[name]++;
    }
}

long OPMetric::GetChannelidMsgCount(std::string &channelid)
{
    std::string name = "channelid_" + channelid;
    if (channel_id_msg_stat.find(name) == channel_id_msg_stat.end())
    {
        return 0;
    }

    return channel_id_msg_stat[name];
}

void OPMetric::ConnectionCountInc(const std::string &connectionId)
{
    std::lock_guard<std::mutex> guard(kMetricMutex);
    kConnectionList.emplace(connectionId);
}

void OPMetric::ConnectionCountDec(const std::string &connectionId)
{
    std::lock_guard<std::mutex> guard(kMetricMutex);
    kConnectionList.erase(connectionId);
}

size_t OPMetric::GetConnectionCount()
{
    std::lock_guard<std::mutex> guard(kMetricMutex);
    return kConnectionList.size();
}

void OPMetric::UdpInvalidMsgInc()
{
    totalUdpErrorMsgCount++;
}

void OPMetric::UdpMsgInc(int bytes)
{
    totalUdpMsgCount++;
    totalUdpMsgBytes += bytes;
}

void OPMetric::HttpMsgInc(int bytes)
{
    totalHttpMsgCount++;
    totalHttpMsgBytes += bytes;
}

void OPMetric::LogMsgInc()
{
    totalLogMsgCount++;
    log_msg_count++;
}

void OPMetric::AccessRequestInc()
{
    totalAccessRequestMsgCount++;
}

void OPMetric::AccessResponseInc()
{
    totalAccessResponseMsgCount++;
}

void OPMetric::AccessResponseFailedInc()
{
    totalAccessResponseFailedMsgCount++;
}

void OPMetric::ReportUdpMsgMetric(uint64_t value)
{
    int thread_id = gse::tools::thread::GetThreadID();
    std::string metric_name_1 = "name_total_udp_msg_count";
    ReportQueneCount(metric_name_1, thread_id, value);
}

void OPMetric::cleanMsgStat()
{
    udp_msg_count += totalUdpMsgCount.load();
    http_msg_count += totalHttpMsgCount.load();
    log_msg_count += totalLogMsgCount.load();

    totalUdpMsgCount = 0;
    totalUdpErrorMsgCount = 0;
    totalHttpMsgCount = 0;
    totalLogMsgCount = 0;

    totalHttpMsgBytes = 0;
    totalUdpMsgBytes = 0;
}

OPMetric *OPMetric::OPMetricInst()
{
    return &OPMetric::Instance();
}

OPMetric::OPMetric()
{
    m_ptrMetricQueue = std::make_shared<moodycamel::ConcurrentQueue<OpsMetricMsg *>>();
}

OPMetric &OPMetric::Instance()
{
    static OPMetric single;
    return single;
}

OPMetric::~OPMetric()
{
    if (nullptr != m_ptrMetricEventThread)
    {
        m_ptrMetricEventThread->Stop();
    }

    if (nullptr != m_ptrMetricQueue)
    {
        OpsMetricMsg *ptrMetricMsg = nullptr;
        while (m_ptrMetricQueue->try_dequeue(ptrMetricMsg))
        {
            if (ptrMetricMsg != nullptr)
            {
                delete ptrMetricMsg;
                ptrMetricMsg = nullptr;
            }
        }
    }
}

void OPMetric::ReportQueneCount(std::string &metric_name, int thread_id, uint64_t queue_size)
{
    OpsMetricMsg *ptr_metricmsg = new OpsMetricMsg();
    ptr_metricmsg->m_metricName = metric_name;
    std::string id = metric_name + gse::tools::strings::ToString(thread_id);
    ptr_metricmsg->m_metricItem = std::make_pair(thread_id, queue_size);
    EnQueneMetricMsg(ptr_metricmsg);
    return;
}

bool OPMetric::EnQueneMetricMsg(OpsMetricMsg *ptr)
{
    int max_retry_times = 0;
    while (!m_ptrMetricQueue->enqueue(ptr))
    {
        max_retry_times++;
        if (max_retry_times > 5)
        {
            delete ptr;
            return false;
        }
    }
    return true;
}

void OPMetric::HandleMetricMsgTimerEvent(evutil_socket_t fd)
{
    OpsMetricMsg *ptr_metric_msg = NULL;
    int count = 0;

    while (m_ptrMetricQueue->try_dequeue(ptr_metric_msg))
    {
        if (ptr_metric_msg != NULL)
        {
            m_serviceMetrics->HandleReportMetricRequest(ptr_metric_msg);
            delete ptr_metric_msg;
            ptr_metric_msg = NULL;
        }
    }
}

void OPMetric::ReportMetricStats()
{
    for (auto it : m_metrics)
    {
        it->ReportMetrics();
    }
}

void OPMetric::DealMetricMsgRequest(evutil_socket_t fd, short what, void *args)
{
    if (NULL == args)
    {
        return;
    }

    OPMetric *ptrThis = reinterpret_cast<OPMetric *>(args);
    ptrThis->HandleMetricMsgTimerEvent(fd);
}

void OPMetric::ReportStats(evutil_socket_t fd, short what, void *args)
{
    if (NULL == args)
    {
        return;
    }

    OPMetric *ptrThis = reinterpret_cast<OPMetric *>(args);
    ptrThis->ReportMetricStats();
}

void OPMetric::CollectMetrics()
{
    for (auto &it : m_metrics)
    {
        it->CollectMetrics();
    }
}

void OPMetric::BaseMetricStatTimer(evutil_socket_t fd, short what, void *args)
{
    if (NULL == args)
    {
        return;
    }

    OPMetric *ptrThis = reinterpret_cast<OPMetric *>(args);
    ptrThis->CollectMetrics();
}

void OPMetric::ReleaseTCFreeMemory(evutil_socket_t fd, short what, void *args)
{
    LOG_DEBUG("Begin release idle tc memory");
    MallocExtension::instance()->ReleaseFreeMemory();
}

void OPMetric::StartMetrics()
{
    auto baseMetric = std::make_shared<BaseMetric>();
    auto messageMetrics = std::make_shared<MessageMetric>();
    m_serviceMetrics = std::make_shared<ServiceModeReportMetrics>();
    m_metrics.push_back(baseMetric);
    m_metrics.push_back(messageMetrics);
    m_metrics.push_back(m_serviceMetrics);

    for (auto &it : m_metrics)
    {
        it->SetConfig(DataServer::GetConfigPtr());
        it->InitBaseMetric();
    }
}

std::shared_ptr<gse::bbx::PrometheusExporter> OPMetric::GetPromethusExporter()
{
    return m_promethusExporter;
}

void OPMetric::MetricEventInit(std::shared_ptr<gse::net::http::HTTPServer> htpServer)
{
    StartMetricExporter(htpServer);

    m_ptrMetricEventThread = std::make_shared<EventThread>();
    int registerResponseVal = m_ptrMetricEventThread->RegisterTimerPersistEvent(OPMetric::DealMetricMsgRequest, this, 1);
    if (GSE_SUCCESS > registerResponseVal)
    {
        LOG_ERROR("failed to register event function callback , error code %d", registerResponseVal);
        return;
    }

    registerResponseVal = m_ptrMetricEventThread->RegisterTimerPersistEvent(OPMetric::ReleaseTCFreeMemory, this, 120);
    if (GSE_SUCCESS > registerResponseVal)
    {
        LOG_ERROR("failed to register timer, error code %d", registerResponseVal);
        return;
    }

    registerResponseVal = m_ptrMetricEventThread->RegisterTimerPersistEvent(OPMetric::ReportStats, this, 1);
    if (GSE_SUCCESS > registerResponseVal)
    {
        LOG_ERROR("failed to register timer, error code %d", registerResponseVal);
        return;
    }

    registerResponseVal = m_ptrMetricEventThread->RegisterTimerPersistEvent(OPMetric::BaseMetricStatTimer, this, 1);
    if (GSE_SUCCESS > registerResponseVal)
    {
        LOG_ERROR("failed to register timer, error code %d", registerResponseVal);
        return;
    }

    int startEventThreadVal = m_ptrMetricEventThread->Start();
    if (GSE_SUCCESS > startEventThreadVal)
    {
        LOG_ERROR("failed to start event thread, error code %d", startEventThreadVal);
        return;
    }
}

int OPMetric::StartMetricExporter(std::shared_ptr<gse::net::http::HTTPServer> htpServer)
{
    StartMetrics();
    m_promethusExporter = std::make_shared<gse::bbx::PrometheusExporter>(htpServer);
    for (auto r : m_metrics)
    {
        m_promethusExporter->RegisterMetricCollectable(r->GetRegistry());
    }

    if (GSE_SUCCESS != m_promethusExporter->Start())
    {
        LOG_ERROR("failed to start prometheus exporter(%s:%d)\n", m_gseCfg.GetPrometheusBindIp().c_str(),
                  m_gseCfg.GetPrometheusListenerPort());
        return GSE_ERROR;
    }

    return GSE_SUCCESS;
}

int BaseMetric::GetNetTrafficByPid(struct NetMetricStat &data, pid_t pid)
{
    FILE *fp;
    char netfile[1024] = {0};
    sprintf(netfile, "/proc/%d/net/dev", pid);
    if ((fp = fopen(netfile, "r")) == NULL)
    {
        return -1;
    }
    char line[1024];
    char devname[512];
    size_t pos = 0;
    while ((fgets(line, 1024, fp) != NULL))
    {
        pos = strcspn(line, ":");
        if (pos == 0 || pos > 16)
            continue;

        ::memcpy(devname, line, pos);
        if (strncmp("  eth1", devname, 6) == 0)
            break;
    }
    fclose(fp);
    static unsigned long long last_packets = 0, last_bytes = 0;
    static unsigned long long last_rbytes = 0, last_wbytes = 0, last_rpackets = 0, last_wpackets = 0;
    unsigned long long r_bytes, t_bytes, r_packets, t_packets;
    unsigned int tmp;
    sscanf(line + pos + 1, "%llu %llu %u %u %u %u %u %u %llu %llu %u %u %u %u %u %u",
           &r_bytes, &r_packets, &tmp, &tmp, &tmp, &tmp, &tmp, &tmp,
           &t_bytes, &t_packets, &tmp, &tmp, &tmp, &tmp, &tmp, &tmp);
    data.net_bytes = r_bytes + t_bytes - last_bytes;
    data.net_packets = r_packets + t_packets - last_packets;
    last_bytes = r_bytes + t_bytes;
    last_packets = r_packets + t_packets;
    data.net_rbytes = r_bytes - last_rbytes;
    data.net_wbytes = t_bytes - last_wbytes;
    data.net_rpackets = r_packets - last_rpackets;
    data.net_wpackets = t_packets - last_wpackets;
    last_rbytes = r_bytes;
    last_wbytes = t_bytes;
    last_rpackets = r_packets;
    last_wpackets = t_packets;
    return 0;
}

int BaseMetric::GetDiskIOByPid(struct DisKIoMetricStat &data, pid_t pid)
{
    FILE *fp;
    char netfile[1024] = {0};
    sprintf(netfile, "/proc/%d/io", pid);
    if ((fp = fopen(netfile, "r")) == NULL)
    {
        return -1;
    }
    char line[1024];
    char name[512];
    size_t pos = 0;
    static unsigned long long last_r = 0, last_w = 0;
    unsigned long long r_bytes = 0, w_bytes = 0;
    while ((fgets(line, 1024, fp) != NULL))
    {
        pos = strcspn(line, ":");
        if (pos == 0 || pos > 16)
        {
            continue;
        }

        ::memcpy(name, line, pos);
        if (strncmp("read_bytes", name, 4) == 0)
        {
            r_bytes = atoi(line + pos + 1);
            continue;
        }
        else if (strncmp("write_bytes", name, 4) == 0)
        {
            w_bytes = atoi(line + pos + 1);
            // read_bytes before write_bytes
            break;
        }
    }
    fclose(fp);
    data.read_bytes = r_bytes - last_r;
    data.write_bytes = w_bytes - last_w;
    last_r = r_bytes;
    last_w = w_bytes;
    return 0;
}

BaseMetric::BaseMetric()
{
    m_registry = std::make_shared<prometheus::Registry>();
    m_startupTimeSecond = gse::tools::time::GetUTCSecond();
}

BaseMetric::~BaseMetric()
{
}

std::shared_ptr<prometheus::Registry> BaseMetric::GetRegistry()
{
    return m_registry;
}

void BaseMetric::InitBaseMetric()
{
    m_labels.insert(std::pair<std::string, std::string>("module", "data"));
    m_labels.insert(std::pair<std::string, std::string>("node", m_configPtr->GetAdvertiseIp()));
    m_labels.insert(std::pair<std::string, std::string>("cluster_name", m_configPtr->GetDataSvrClusterName()));
    m_labels.insert(std::pair<std::string, std::string>("instance_id", m_configPtr->GetDataSvrInstanceId()));
}

void BaseMetric::CollectMetrics()
{
    SystemStats();
}

void BaseMetric::SetConfig(std::shared_ptr<DataProcessConfig> configPtr)
{
    m_configPtr = configPtr;
}

void BaseMetric::SystemStats()
{
    static bool is_init = false;
    static bool last;

    pid_t pid = getpid();
    static uint64_t last_totalcputime;
    static uint64_t last_procputime;
    if (!is_init)
    {
        last_totalcputime = gse::tools::process::GetTotalCpuTime();
        last_procputime = gse::tools::process::GetProcessCpuTime(pid);
        is_init = true;
    }

    uint64_t totalcputime = gse::tools::process::GetTotalCpuTime();
    uint64_t procputime = gse::tools::process::GetProcessCpuTime(pid);

    m_cpuUsage = 100.0f * (procputime - last_procputime) / (totalcputime - last_totalcputime);
    last_totalcputime = totalcputime;
    last_procputime = procputime;

    GetNetTrafficByPid(m_netMetric, pid);
    GetDiskIOByPid(m_diskIoMetric, pid);
}

int BaseMetric::GetProcessThreadCount(pid_t pid)
{
    int threadCount = 0;
    char fileName[64] = {0};
    char lineBuff[1024] = {0};

    FILE *fd = NULL;

    sprintf(fileName, "/proc/%d/status", pid);

    fd = fopen(fileName, "r");
    if (fd == NULL)
    {
        return threadCount;
    }
    char *line = NULL;
    do
    {
        line = fgets(lineBuff, sizeof(lineBuff), fd);
        if (NULL == line)
        {
            break;
        }

        if (strncmp(lineBuff, "Threads", strlen("Threads")) == 0)
        {
            sscanf(lineBuff, "Threads: %ld", &threadCount);
            break;
        }
    } while (line != NULL);

    fclose(fd);
    return threadCount;
}

void BaseMetric::ReportMetrics()
{
    pid_t pid = tools::process::GetProcessID();
    double pmem = tools::process::GetProcessMemUsage(pid);
    int threadNum = GetProcessThreadCount(pid);

    auto &cpuGaugeFamily = prometheus::BuildGauge().Name(data_cpu).Help("cpu usage").Register(*m_registry);
    auto &cpuGauge = cpuGaugeFamily.Add(m_labels);
    cpuGauge.Set(m_cpuUsage);

    auto &memGaugeFamily = prometheus::BuildGauge().Name(data_mem).Help("mem usage").Register(*m_registry);
    auto &memGauge = memGaugeFamily.Add(m_labels);
    memGauge.Set(pmem);

    auto memSize = tools::process::GetProcessMemSize(pid);
    auto &memSizeGaugeFamily = prometheus::BuildGauge().Name(data_mem_bytes).Help("process used mem size").Register(*m_registry);
    auto &memSizeGauge = memSizeGaugeFamily.Add(m_labels);
    memSizeGauge.Set(memSize);

    auto processRunTimes = gse::tools::time::GetUTCSecond() - m_startupTimeSecond;

    auto &processRunTimesGaugeFamily = prometheus::BuildGauge().Name(name_process_run_seconds_total).Help("rocess run seconds total").Register(*m_registry);
    auto &processRunTimesGauge = processRunTimesGaugeFamily.Add(m_labels);
    processRunTimesGauge.Set(processRunTimes);

    auto &threadCountGaugeFamily = prometheus::BuildGauge().Name(data_thread).Help("process thread count").Register(*m_registry);
    auto &threadCountGauge = threadCountGaugeFamily.Add(m_labels);
    threadCountGauge.Set(threadNum);

    auto &diskIoReadBytesCounterFamily = prometheus::BuildCounter().Name(data_disk_io_read_bytes).Help("process read disk bytes").Register(*m_registry);
    auto &diskIoReadBytesCounter = diskIoReadBytesCounterFamily.Add(m_labels);
    diskIoReadBytesCounter.Increment(m_diskIoMetric.read_bytes);

    auto &diskIoWriteBytesCounterFamily = prometheus::BuildCounter().Name(data_disk_io_write_bytes).Help("process write disk bytes").Register(*m_registry);
    auto &diskIoWriteBytesCounter = diskIoWriteBytesCounterFamily.Add(m_labels);
    diskIoWriteBytesCounter.Increment(m_diskIoMetric.write_bytes);

    //    m_m_netFlowTxBytesGudgeFamily = &prometheus::BuildGauge().Name(data_net_flow_tx_bytes).Help("").Register(*registry);
    //    m_netFlowTxBytesGauge = &m_m_netFlowTxBytesGudgeFamily.Add(m_labels);

    //    m_netFlowRxBytesGugeFamily = &prometheus::BuildGauge().Name(data_net_flow_rx_bytes).Help("").Register(*registry);
    //    m_netFlowRxBytesGauge = &m_netFlowRxBytesGugeFamily.Add(m_labels);

    //    m_netFlowRxPakcagesGugeFamily = &prometheus::BuildGauge().Name(data_net_flow_rx_pakcages).Help("").Register(*registry);
    //    m_netFlowRxPakcagesGauge = &m_netFlowRxBytesGugeFamily.Add(m_labels);

    //    m_netFlowTxPakcagesGugeFamily = &prometheus::BuildGauge().Name(data_net_flow_tx_pakcages).Help("").Register(*registry);
    //    m_netFlowTxPakcagesGauge = &m_netFlowTxPakcagesGugeFamily.Add(m_labels);

    //    m_diskIoReadBytesGauge->Set(m_diskIoMetric.read_bytes);
    //    m_diskIoWriteBytesGauge->Set(m_diskIoMetric.write_bytes);
    //    m_netFlowRxBytesGauge->Set(m_netMetric.net_rbytes);
    //    m_netFlowRxBytesGauge->Set(m_netMetric.net_rpackets);
    //    m_netFlowRxBytesGauge->Set(m_netMetric.net_wbytes);
    //    m_netFlowRxBytesGauge->Set(m_netMetric.net_wpackets);
}

MessageMetric::MessageMetric()
{
    m_registry = std::make_shared<prometheus::Registry>();
    m_labels.clear();
}

MessageMetric::~MessageMetric()
{
}

std::shared_ptr<prometheus::Registry> MessageMetric::GetRegistry()
{
    return m_registry;
}

void MessageMetric::InitBaseMetric()
{
    m_labels.insert(std::pair<std::string, std::string>("module", "data"));
    m_labels.insert(std::pair<std::string, std::string>("cluster_name", m_configPtr->GetDataSvrClusterName()));
    m_labels.insert(std::pair<std::string, std::string>("instance_id", m_configPtr->GetDataSvrInstanceId()));
    m_labels.insert(std::pair<std::string, std::string>("node", m_configPtr->GetAdvertiseIp()));

    kTcpMsgMsgCounterFamily = &prometheus::BuildCounter().Name(name_recv_msg_metrics).Help("recv msg metrics").Labels(m_labels).Register(*m_registry);

    kRecvMsgBytesCounterFamily = &prometheus::BuildCounter().Name(name_recv_msg_bytes_metrics).Help("recv msg bytes metrics").Labels(m_labels).Register(*m_registry);
    kSendMsgBytesCounterFamily = &prometheus::BuildCounter().Name(name_send_msg_bytes_metrics).Help("send msg bytes metrics").Labels(m_labels).Register(*m_registry);
    kSendMsgCounterFamily = &prometheus::BuildCounter().Name(name_send_msg_count).Help("send msg counter metrics").Labels(m_labels).Register(*m_registry);

    kMsgDelayCounterFamily = &prometheus::BuildCounter().Name(name_msg_delay_metrics).Help("message delay ms metrics").Labels(m_labels).Register(*m_registry);

    kMessageQueueRequestCounterFamily = &prometheus::BuildCounter().Name(name_message_quene_request_metrics).Help("message quene request couter metrics").Labels(m_labels).Register(*m_registry);

    kMessageQueueRequestMillsecondsFamily = &prometheus::BuildCounter().Name(name_message_quene_request_millseconds_metrics).Help("message quene request millseconds metrics").Labels(m_labels).Register(*m_registry);

    kChannelIdCounterFamily = &prometheus::BuildGauge().Name(name_channel_id_count).Help("channelid count").Labels(m_labels).Register(*m_registry);

    kChannelIdMemoryUsageFamily = &prometheus::BuildGauge().Name(name_channel_id_memusage).Help("channelid mem usage").Labels(m_labels).Register(*m_registry);
}

void MessageMetric::MessageBytesMetric()
{
    unsigned long udpMsgBytes = totalUdpMsgBytes.load();
    auto &udpMsgBytesCounterFamily = prometheus::BuildCounter().Name(name_total_udp_recv_msg_bytets).Help("recv udp msg bytes").Register(*m_registry);
    auto &recvUdpMsgBytesCounter = udpMsgBytesCounterFamily.Add(m_labels);
    recvUdpMsgBytesCounter.Increment(udpMsgBytes);

    unsigned long httpMsgBytes = totalHttpMsgBytes.load();
    auto &httpMsgBytesCounterFamily = prometheus::BuildCounter().Name(name_total_http_recv_msg_bytets).Help("recv http msg bytes").Register(*m_registry);
    auto &recvHttpMsgBytesCounter = httpMsgBytesCounterFamily.Add(m_labels);
    recvHttpMsgBytesCounter.Increment(httpMsgBytes);
}
void MessageMetric::ReportMetrics()
{
    unsigned long udpMsgCount = totalUdpMsgCount.load();
    unsigned long httpMsgCount = totalHttpMsgCount.load();

    unsigned long logLineCount = totalLogMsgCount.load();

    unsigned long accessRequestMsgCount = totalAccessRequestMsgCount.load();
    unsigned long accessResponseMsgCount = totalAccessResponseMsgCount.load();
    unsigned long accessResponseFailedMsgCount = totalAccessResponseFailedMsgCount.load();

    auto apiLabel = m_labels;
    apiLabel.insert({"apiname", "access_request"});
    auto &msgRequest = prometheus::BuildCounter().Name(name_total_request_msg_count).Help("access request msg count").Register(*m_registry);
    auto &requestMsgCounter = msgRequest.Add(apiLabel);
    requestMsgCounter.Increment(accessRequestMsgCount);

    apiLabel = m_labels;
    apiLabel.insert({"apiname", "access_response"});
    apiLabel.insert({"code", "success"});
    auto &msgResponse = prometheus::BuildCounter().Name(name_total_request_msg_count).Help("access request msg count").Register(*m_registry);
    auto &responseMsgCounter = msgResponse.Add(apiLabel);
    responseMsgCounter.Increment(accessResponseMsgCount);

    apiLabel = m_labels;
    apiLabel.insert({"apiname", "access_response"});
    apiLabel.insert({"code", "failed"});
    auto &msgResponseFailed = prometheus::BuildCounter().Name(name_total_request_msg_count).Help("access request msg count").Register(*m_registry);
    auto &responseFailedMsgCounter = msgResponseFailed.Add(apiLabel);
    responseFailedMsgCounter.Increment(accessResponseFailedMsgCount);

    auto &recvUdpMsgCounterFamily = prometheus::BuildCounter().Name(name_total_udp_msg_count).Help("recv udp msg count").Register(*m_registry);
    auto &recvUdpMsgCounter = recvUdpMsgCounterFamily.Add(m_labels);
    recvUdpMsgCounter.Increment(udpMsgCount);

    auto &recvHttpMsgCounterFamily = prometheus::BuildCounter().Name(name_total_http_msg_count).Help("recv http msg count").Register(*m_registry);
    auto &recvHttpMsgCounter = recvHttpMsgCounterFamily.Add(m_labels);
    recvHttpMsgCounter.Increment(httpMsgCount);

    auto &writeLogCounterFamily = prometheus::BuildCounter().Name(name_total_log_msg_count).Help("write log line count").Register(*m_registry);
    auto &writeLogMsgCounter = writeLogCounterFamily.Add(m_labels);
    writeLogMsgCounter.Increment(logLineCount);

    auto &connectionGaugeFamily = prometheus::BuildGauge().Name(name_tcp_connection_count).Help("tcp connection count").Register(*m_registry);
    auto &connectionGauge = connectionGaugeFamily.Add(m_labels);

    connectionGauge.Set(kConnectionList.size());

    MessageBytesMetric();
    Clean();
}

void MessageMetric::CollectMetrics()
{
}

void MessageMetric::SetConfig(std::shared_ptr<DataProcessConfig> configPtr)
{
    m_configPtr = configPtr;
}

void MessageMetric::Clean()
{
    udp_msg_count += totalUdpMsgCount.load();
    http_msg_count += totalHttpMsgCount.load();
    log_msg_count += totalLogMsgCount.load();

    totalUdpMsgCount = 0;
    totalUdpErrorMsgCount = 0;
    totalHttpMsgCount = 0;
    totalLogMsgCount = 0;
    totalAccessRequestMsgCount = 0;
    totalAccessResponseMsgCount = 0;
    totalAccessResponseFailedMsgCount = 0;

    totalHttpMsgBytes = 0;
    totalUdpMsgBytes = 0;
}

ServiceModeReportMetrics::ServiceModeReportMetrics()
{
    m_registry = std::make_shared<prometheus::Registry>();
}

ServiceModeReportMetrics::~ServiceModeReportMetrics()
{
}

void ServiceModeReportMetrics::InitBaseMetric()
{
    m_labels = {{"module", "data"},
                {"cluster_name", m_configPtr->GetDataSvrClusterName()},
                {"instance_id", m_configPtr->GetDataSvrInstanceId()},
                {"node", m_configPtr->GetAdvertiseIp()}};
}

std::shared_ptr<prometheus::Registry> ServiceModeReportMetrics::GetRegistry()
{
    return m_registry;
}
void ServiceModeReportMetrics::HandleReportMetricRequest(OpsMetricMsg *ptr_ops)
{
    std::map<std::string, std::map<std::string, uint64_t>>::iterator it;
    it = m_metrics.find(ptr_ops->m_metricName);
    if (it != m_metrics.end())
    {
        std::map<std::string, uint64_t> perthread_item = it->second;

        perthread_item[ptr_ops->m_metricItem.first] = ptr_ops->m_metricItem.second;
    }
    else
    {
        std::map<std::string, uint64_t> perthread_item;
        perthread_item[ptr_ops->m_metricItem.first] = ptr_ops->m_metricItem.second;
        m_metrics[ptr_ops->m_metricName] = perthread_item;
    }
}

void ServiceModeReportMetrics::CollectMetrics()
{
}

void ServiceModeReportMetrics::SetConfig(std::shared_ptr<DataProcessConfig> configPtr)
{
    m_configPtr = configPtr;
}

void ServiceModeReportMetrics::ReportMetrics()
{
    std::map<std::string, std::map<std::string, uint64_t>>::iterator it;
    for (it = m_metrics.begin(); it != m_metrics.end(); it++)
    {
        std::map<std::string, uint64_t>::iterator itMetric;
        uint64_t total = 0;
        for (itMetric = it->second.begin(); itMetric != it->second.end(); itMetric++)
        {
            total += itMetric->second;
        }

        auto &counterFamily = prometheus::BuildCounter().Name(it->first).Help("").Register(*m_registry);
        auto &counter = counterFamily.Add(m_labels);
        counter.Increment(total);
    }

    m_metrics.clear();
}

OPSHttpServer::OPSHttpServer()
    : m_listennerFd(-1), m_port(0)
{
}

OPSHttpServer::OPSHttpServer(const std::string &ip, uint16_t port)
    : m_serverIp(ip), m_port(port), m_listennerFd(-1)
{
}

OPSHttpServer::~OPSHttpServer()
{
}

int OPSHttpServer::Start()
{
    std::string bindIp = m_serverIp;
    if ("" == m_serverIp)
    {
        bindIp = gse::tools::net::GetMachineIp();
    }

    m_httpServer = std::make_shared<gse::net::http::HTTPServer>(bindIp, m_port, 2);

    if (DataServer::GetUpgradeFlag())
    {
        StartMigrationSerivce();
    }

    auto htpServerThread = std::bind(&gse::net::http::HTTPServer::Start, m_httpServer);
    m_listenThread = std::thread(htpServerThread);
    return true;
}

void OPSHttpServer::Stop()
{
    m_httpServer->Stop();
    Join();
}

int OPSHttpServer::Join()
{
    if (m_listenThread.joinable())
    {
        m_listenThread.join();
    }
    return GSE_SUCCESS;
}

int OPSHttpServer::StartMigrationSerivce()
{
    m_migrationClient = std::unique_ptr<gse::net::MigrationClient>(new gse::net::MigrationClient(m_port, 20));
    int domainsocket_listenner_fd = -1;
    if (m_migrationClient->ConnectDomainSocket() != GSE_SUCCESS)
    {
        LOG_WARN("failed to connect domain socket");
    }
    else
    {
        m_listennerFd = m_migrationClient->MigrateListennerFd();
        domainsocket_listenner_fd = m_migrationClient->MigrateDomainSocketListenner();
    }

    auto pFuncGetListernnerFd = std::bind(&gse::net::http::HTTPServer::GetEvhtpListennerFd, m_httpServer);
    auto pFuncStopListenner = std::bind(&gse::net::http::HTTPServer::StoppingListenner, m_httpServer);

    m_migrationServer = std::unique_ptr<gse::net::MigrationServer>(new gse::net::MigrationServer(m_port));
    m_migrationServer->SetGetListennerFdCallback(pFuncGetListernnerFd);
    // m_migrationServer->SetStopListennerCallback(pFuncStopListenner);
    m_httpServer->SetListenerFd(m_listennerFd);
    // std::string m_ipcFilePath
    LOG_DEBUG("Migrate ops http service, listener socket fd:%d", m_listennerFd);
    return m_migrationServer->StartMigrationService(domainsocket_listenner_fd);
}

std::shared_ptr<gse::net::http::HTTPServer> OPSHttpServer::GetHttpServer()
{
    return m_httpServer;
}

bool OPSHttpServer::RegisterMsgHandler(const std::string &uri, std::shared_ptr<net::http::HTTPHandler> handler)
{
    return m_httpServer->RegisterHandler(uri, handler);
}
} // namespace data
} // namespace gse
