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

#include <vector>
#include <unordered_map>
#include <atomic>
#include <gperftools/malloc_extension.h>

#include "tools/process.h"
#include "bbx/gse_errno.h"
#include "log/log.h"
#include "common/time.h"
#include "conf/bkdata_config_v1.h"
#include "time_center.h"


namespace gse {
namespace dataserver {


//metrics
//---------------------------------------------------------------------
const std::string data_cpu= "data_cpu";
const std::string data_mem= "data_mem";
const std::string data_disk_io_read_bytes = "data_disk_io_read_bytes";
const std::string data_disk_io_write_bytes = "data_disk_io_write_bytes";
const std::string data_net_flow_tx_bytes = "data_net_flow_tx_bytes";
const std::string data_net_flow_rx_bytes = "data_net_flow_rx_bytes";
const std::string data_net_flow_rx_pakcages = "data_net_flow_rx_pakcages";
const std::string data_net_flow_tx_pakcages = "data_net_flow_tx_pakcages";

const std::string data_net_flow_total_pakcages = "data_net_flow_total_pakcages";
const std::string data_net_flow_total_bytes = "data_net_flow_total_bytes";

const std::string data_fd_count = "data_fd_count";
const std::string data_tcp_count = "data_tcp_count";
const std::string data_thread= "data_threads_count";
const std::string total_msg_count = "data_total_msg_count";
const std::string name_total_udp_msg_count = "data_total_udp_msg_count";
const std::string name_total_tcp_msg_count = "data_total_tcp_msg_count";
const std::string name_total_http_msg_count = "data_total_http_msg_count";
const std::string name_total_quene_msg_count = "data_quene_msg_count";
const std::string name_total_kafka_msg_count = "data_total_kafka_msg_count";
const std::string name_total_redis_msg_count = "data_total_redis_msg_count";
const std::string name_total_log_msg_count = "data_total_log_msg_count";
const std::string name_total_proxy_msg_count = "data_total_proxy_msg_count";
const std::string name_total_pulsar_msg_count = "data_total_pulsar_msg_count";
const std::string name_total_udp_proxy_msg_count = "data_total_udp_proxy_msg_count";

const std::string name_tcp_connection_count = "data_tcp_connection_count";

const std::string name_tcp_transfer_queune_msg_count = "data_tcp_proxy_queune_msg_count";
const std::string name_udp_receiver_queune_msg_count = "data_udp_receiver_queune_msg_count";
const std::string name_worker_queune_msg_count = "data_worker_queune_msg_count";

//---------------channelid stat----------------

static std::unordered_map<std::string, unsigned long > channel_id_msg_stat;
static std::unordered_map<std::string, unsigned long > channel_id_error_msg_stat;

//end channelid stat------------
const std::string msg_max_delay = "msg_msg_delay";
const std::string channel_id_count = "channel_id_msg_count";
const std::string total_write_kafka_msg_count = "total_write_kafka_msg_count";
const std::string kafka_connection_count = "kafak_connection_count";

//--------------------------------------------------------
static std::atomic<unsigned long> total_udp_msg_count(0);
static std::atomic<unsigned long> total_udp_error_msg_count(0);
static std::atomic<unsigned long> total_tcp_msg_count(0);
static std::atomic<unsigned long> total_http_msg_count(0);
//-------------------------------------------------------

static std::atomic<unsigned long> total_kafka_msg_count(0);
static std::atomic<unsigned long> total_redis_msg_count(0);
static std::atomic<unsigned long> total_log_msg_count(0);
static std::atomic<unsigned long> total_proxy_msg_count(0);
static std::atomic<unsigned long> total_udp_proxy_msg_count(0);
static std::atomic<unsigned long> total_pulsar_msg_count(0);

static long udp_msg_count;
static long tcp_msg_count;
static long http_msg_count;
static long kafka_msg_count;
static long redis_msg_count;
static long log_msg_count;
static long proxy_msg_count;
static long pulsar_msg_count;


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

long OPMetric::GetKafakMsgCount()
{
    return kafka_msg_count;
}

long OPMetric::GetRedisMsgCount()
{
    return redis_msg_count;
}

long OPMetric::GetLogMsgCount()
{
    return log_msg_count;
}

long OPMetric::GetProxyMsgCount()
{
    return proxy_msg_count;
}

long OPMetric::GetPulsarMsgCount()
{
    return pulsar_msg_count;
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

void OPMetric::ChannelIdErrorMsgStat(std::string &channelid, uint32_t error_code)
{
    std::string key = channelid + "_" + gse::tools::strings::ToString(error_code);

    if (channel_id_error_msg_stat.find(key) == channel_id_error_msg_stat.end())
    {
        channel_id_error_msg_stat.insert(std::make_pair(key, 1));
    }
    else
    {
        channel_id_error_msg_stat[key]++;
    }
}

void OPMetric::ConnectionCountInc()
{
    //tcp_connection_count
}

void OPMetric::UdpInvalidMsgInc()
{
    total_udp_error_msg_count++;
}

void OPMetric::UdpMsgInc()
{
    total_udp_msg_count++;
}

void OPMetric::TcpMsgInc()
{
    total_tcp_msg_count++;
    tcp_msg_count++;
}

void OPMetric::HttpMsgInc()
{
    total_http_msg_count++;
    http_msg_count++;
}

void OPMetric::ProxyMsgInc()
{
    total_proxy_msg_count++;
    proxy_msg_count++;
}

void OPMetric::ProxyUdpMsgInc()
{
    total_udp_proxy_msg_count++;
}

void OPMetric::LogMsgInc()
{
    total_log_msg_count++;
    log_msg_count++;
}

void OPMetric::KafkaMsgInc()
{
    total_kafka_msg_count++;
}

void OPMetric::PulsarMsgInc()
{
    total_pulsar_msg_count++;
    pulsar_msg_count++;
}

void OPMetric::RedisMsgInc()
{
    total_redis_msg_count++;
    redis_msg_count++;
}

void OPMetric::ReportUdpMsgMetric(uint64_t value)
{
    int thread_id = gse::tools::thread::GetThreadID();
    std::string metric_name_1 = "name_total_udp_msg_count";
    ReportQueneCount(metric_name_1, thread_id, value);
}

void OPMetric::cleanMsgStat()
{
    udp_msg_count += total_udp_msg_count.load();
    tcp_msg_count += total_tcp_msg_count.load();
    http_msg_count += total_http_msg_count.load();
    kafka_msg_count += total_kafka_msg_count.load();
    redis_msg_count += total_redis_msg_count.load();
    log_msg_count += total_log_msg_count.load();
    proxy_msg_count += total_proxy_msg_count.load();
    pulsar_msg_count += total_proxy_msg_count.load();


    total_udp_msg_count = 0;
    total_tcp_msg_count = 0;
    total_udp_error_msg_count = 0;
    total_kafka_msg_count = 0;
    total_redis_msg_count = 0;
    total_log_msg_count = 0;
    total_proxy_msg_count = 0;
    total_pulsar_msg_count = 0;
    total_udp_proxy_msg_count = 0;

}

OPMetric::OPMetric()
{
    m_ptrMetricQueue = new moodycamel::ConcurrentQueue<OpsMetricMsg *>();
    m_ptrMetricEventThread = NULL;
    m_opsMsgHandleEventId = -1;
    m_cpuUsage = 0.0;
    m_exporter = NULL;

//    m_metricRegistry = std::make_shared<prometheus::Registry>();

//    m_exporter = new gse::bbx::PrometheusExporter(m_gseCfg.m_promethusBindIp, m_gseCfg.m_promethusPort, 4);

//    m_exporter->RegisterMetricCollectable(m_metricRegistry);
}


OPMetric& OPMetric::Instance()
{
    static OPMetric single;
    return single;
}

OPMetric::~OPMetric()
{
    if (NULL != m_ptrMetricEventThread)
    {
        m_ptrMetricEventThread->stop();
        delete m_ptrMetricEventThread;
        m_ptrMetricEventThread = NULL;
    }

    if (m_ptrMetricQueue != NULL)
    {
        delete m_ptrMetricQueue;
    }
}

void OPMetric::ReportQueneCount(std::string &metric_name, int thread_id, uint64_t queue_size)
{
    OpsMetricMsg *ptr_metricmsg = new OpsMetricMsg();
    ptr_metricmsg->m_metricName = metric_name;
    std::string id = metric_name + gse::tools::strings::ToString(thread_id);
    ptr_metricmsg->m_metricItem = std::make_pair(id, queue_size);
    EnQueneMetricMsg(ptr_metricmsg);
    return;
}

bool OPMetric::EnQueneMetricMsg(OpsMetricMsg *ptr)
{
    int max_retry_times = 0;
    while (!m_ptrMetricQueue->enqueue(ptr)) {
        max_retry_times++;
        if (max_retry_times > 5)
        {
            delete ptr;
            return false;
        }
    }
    return true;
}

void OPMetric::HandleMetricOpsMsg(OpsMetricMsg *ptr_ops)
{
    std::map<std::string, std::map<std::string, uint64_t>>::iterator it;
    it = m_metrics.find(ptr_ops->m_metricName);
    if (it != m_metrics.end())
    {
        std::map<std::string, uint64_t> perthread_item  = it->second;

        perthread_item[ptr_ops->m_metricItem.first] = ptr_ops->m_metricItem.second;
    }
    else
    {
        std::map<std::string, uint64_t> perthread_item;
        perthread_item[ptr_ops->m_metricItem.first] = ptr_ops->m_metricItem.second;
        m_metrics[ptr_ops->m_metricName] = perthread_item;
    }
}

void OPMetric::HandleMetricMsgTimerEvent(evutil_socket_t fd)
{
    OpsMetricMsg *ptr_metric_msg = NULL;
    int count = 0;

    while (m_ptrMetricQueue->try_dequeue(ptr_metric_msg))
    {
        if (ptr_metric_msg != NULL)
        {
            HandleMetricOpsMsg(ptr_metric_msg);
            delete ptr_metric_msg;
            ptr_metric_msg = NULL;
        }
    }
}


void OPMetric::BaseMetric()
{

    // counter familiy

//    std::string metric_name("base_metric");
//    Label label;
//    label.AddKeyValue("ip", m_reportIP);
//    label.AddKeyValue("cluster_name", m_gseCfg.getDataSvrClusterName());
//    label.AddKeyValue("instance_id", m_gseCfg.getDataSvrInstanceId());

//    pid_t pid = getpid();
//    double pmem = ProResource::get_pmem(pid);
//    int threadNum = tgse::getThreadNum(pid);

//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_cpu, label,
//        std::make_pair(m_cpuUsage, 0) , Gauge);

//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_mem, label,
//        std::make_pair(pmem, 0) , Gauge);

//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_thread, label,
//        std::make_pair(threadNum, 0) , Gauge);

//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_disk_io_read_bytes, label,
//        std::make_pair(m_diskIoMetric.read_bytes, 0) , Gauge);

//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_disk_io_write_bytes, label,
//        std::make_pair(m_diskIoMetric.write_bytes, 0) , Gauge);


//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_net_flow_tx_bytes, label,
//        std::make_pair(m_netMetric.net_wbytes, 0) , Gauge);


//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_net_flow_rx_bytes, label,
//        std::make_pair(m_netMetric.net_rbytes, 0) , Gauge);


//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_net_flow_rx_pakcages, label,
//        std::make_pair(m_netMetric.net_rpackets, 0) , Gauge);


//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_net_flow_tx_pakcages, label,
//        std::make_pair(m_netMetric.net_wpackets, 0) , Gauge);

//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_net_flow_total_pakcages, label,
//        std::make_pair(m_netMetric.net_packets, 0) , Gauge);


//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_net_flow_total_bytes, label,
//        std::make_pair(m_netMetric.net_bytes, 0) , Gauge);
}

void OPMetric::SumMetricStatItem()
{

//    Label label;
//    label.AddKeyValue("ip", m_reportIP);
//    label.AddKeyValue("cluster_name", m_gseCfg.getDataSvrClusterName());
//    label.AddKeyValue("instance_id", m_gseCfg.getDataSvrInstanceId());

//    std::map<std::string, std::map<std::string, uint64_t>>::iterator it;
//    for (it = m_metrics.begin(); it != m_metrics.end(); it++)
//    {
//        std::map<std::string, uint64_t>::iterator it_item;
//        uint64_t total = 0;
//        for (it_item = it->second.begin(); it_item != it->second.end(); it_item++)
//        {
//            total += it_item->second;
//        }

//        PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, it->first, label,
//            std::make_pair(total, 0) , Gauge);
//    }

//    m_metrics.clear();
}

void OPMetric::MsgMetric()
{
//    Label label;
//    label.AddKeyValue("ip", m_reportIP);
//    label.AddKeyValue("cluster_name", m_gseCfg.getDataSvrClusterName());
//    label.AddKeyValue("instance_id", m_gseCfg.getDataSvrInstanceId());

//    unsigned long udp_count = total_udp_msg_count.load();
//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, name_total_udp_msg_count, label,
//        std::make_pair(udp_count, 0) , Gauge);

//    unsigned long tcp_count = total_tcp_msg_count.load();
//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, name_total_tcp_msg_count, label,
//        std::make_pair(tcp_count, 0) , Gauge);

//    unsigned long http_count = total_http_msg_count.load();
//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, name_total_http_msg_count, label,
//        std::make_pair(http_count, 0) , Gauge);

//    unsigned long kafka_count = total_kafka_msg_count.load();
//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, name_total_kafka_msg_count, label,
//        std::make_pair(kafka_count, 0) , Gauge);

//    unsigned long redis_count = total_redis_msg_count.load();
//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, name_total_redis_msg_count, label,
//        std::make_pair(redis_count, 0) , Gauge);

//    unsigned long log_count = total_log_msg_count.load();
//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, name_total_log_msg_count, label,
//        std::make_pair(log_count, 0) , Gauge);

//    unsigned long proxy_count = total_proxy_msg_count.load();
//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, name_total_proxy_msg_count, label,
//        std::make_pair(proxy_count, 0) , Gauge);

//    unsigned long proxy_udp_count = total_udp_proxy_msg_count.load();
//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, name_total_udp_proxy_msg_count, label,
//        std::make_pair(proxy_udp_count, 0) , Gauge);

//    unsigned long pulsar_count = total_pulsar_msg_count.load();
//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, name_total_pulsar_msg_count, label,
//        std::make_pair(pulsar_count, 0) , Gauge);

//    LOG_DEBUG("Message metric, udp:%d, tcp:%d, http:%d, kafka:%d, redis:%d, log:%d, proxy:%d, pulsar:%d", udp_count, tcp_count, http_count, kafka_count, redis_count, log_count, proxy_count, pulsar_count);
//    cleanMsgStat();
}


void OPMetric::channelidStatMetric()
{
//    Label label;
//    label.AddKeyValue("ip", m_reportIP);
//    label.AddKeyValue("cluster_name", m_gseCfg.getDataSvrClusterName());
//    label.AddKeyValue("instance_id", m_gseCfg.getDataSvrInstanceId());

//    for (std::unordered_map<std::string,  unsigned long>::iterator it_channelid = channel_id_msg_stat.begin(); it_channelid != channel_id_msg_stat.end(); it_channelid++)
//    {
//        PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, it_channelid->first, label, std::make_pair(it_channelid->second, 0), Gauge);
//    }
//    channel_id_msg_stat.clear();

//    for (std::unordered_map<std::string,  unsigned long>::iterator it_channelid = channel_id_error_msg_stat.begin(); it_channelid != channel_id_error_msg_stat.end(); it_channelid++)
//    {
//        PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, it_channelid->first, label, std::make_pair(it_channelid->second, 0), Gauge);
//    }
//    channel_id_error_msg_stat.clear();
}

void OPMetric::SystemStats()
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

    m_cpuUsage = 100.0f*(procputime - last_procputime)/(totalcputime - last_totalcputime);
    last_totalcputime = totalcputime;
    last_procputime = procputime;

    GetNetTrafficByPid(m_netMetric, pid);
    GetDiskIOByPid(m_diskIoMetric, pid);
}

void OPMetric::ReportMetricStats()
{
    BaseMetric();
    MsgMetric();
    channelidStatMetric();
    SumMetricStatItem();
}

void OPMetric::DealMetricMsgRequest(evutil_socket_t fd, short what, void *args)
{
    if (NULL == args)
    {
        return;
    }

    OPMetric *ptr_this = reinterpret_cast<OPMetric *>(args);
    ptr_this->HandleMetricMsgTimerEvent(fd);
}


void OPMetric::ReportStats(evutil_socket_t fd, short what, void *args)
{
	if (NULL == args)
	{
		return;
	}

    OPMetric *ptr_this = reinterpret_cast<OPMetric *>(args);
    ptr_this->ReportMetricStats();
}

void OPMetric::BaseMetricStatTimer(evutil_socket_t fd, short what, void *args)
{
    if (NULL == args)
    {
        return;
    }

    OPMetric *ptr_this = reinterpret_cast<OPMetric *>(args);
    ptr_this->SystemStats();
}


void OPMetric::ReleaseTCFreeMemory(evutil_socket_t fd, short what, void *args)
{
    LOG_DEBUG("Begin release idle tc memory");
	MallocExtension::instance()->ReleaseFreeMemory();
}

void OPMetric::MetricEventInit()
{
    m_reportIP = m_gseCfg.getDataSvrIp();
    m_ptrMetricEventThread = new rgse::GseEventThread();
    int registerResponseVal = m_ptrMetricEventThread->registerTimerPersistEvent(OPMetric::DealMetricMsgRequest, this, 1);
    if (registerResponseVal < 0)
    {
        LOG_ERROR("failed to register event function callback , error code %d", m_opsMsgHandleEventId);
        return;
    }

    registerResponseVal = m_ptrMetricEventThread->registerTimerPersistEvent(OPMetric::ReleaseTCFreeMemory, this, 120);
    if (GSE_SUCCESS > registerResponseVal)
    {
        LOG_ERROR("failed to register timer, error code %d", registerResponseVal);
        return ;
    }

    registerResponseVal = m_ptrMetricEventThread->registerTimerPersistEvent(OPMetric::ReportStats, this, 10);
    if (GSE_SUCCESS > registerResponseVal)
    {
        LOG_ERROR("failed to register timer, error code %d", registerResponseVal);
        return;
    }

    registerResponseVal = m_ptrMetricEventThread->registerTimerPersistEvent(OPMetric::BaseMetricStatTimer, this, 1);
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

int OPMetric::GetNetTrafficByPid(struct NetMetricStat &data, pid_t pid)
{
    FILE * fp;
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
        // eh1
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

int OPMetric::GetDiskIOByPid(struct DisKIoMetricStat &data, pid_t pid)
{
    FILE * fp;
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


//BaseMetric::BaseMetric(prometheus::Registry *registry)
//{

//    prometheus::Label label;
//    label.Add("ip", m_cfg.m_serverIp);
//    label.Add("cluster_name", m_cfg.m_clusterName);
//    label.Add("instance_id", m_cfg.m_instanceId);

//    std::map<std::string, std::string> labels;
//    labels.insert(std::pair<std::string, std::string>("ip", m_cfg.m_serverIp));
//    labels.insert(std::pair<std::string, std::string>("cluster_name", m_cfg.m_clusterName));
//    labels.insert(std::pair<std::string, std::string>("instance_id", m_cfg.m_instanceId));

//    m_cpu = prometheus::BuildCounter().Name(data_cpu).Add(labels).Register(*registry);
//    m_mem = prometheus::BuildCounter().Name(data_mem).Add(labels).Register(*registry);
//    m_threadCount = prometheus::BuildCounter().Name(data_thread).Add(labels).Register(*registry);

//    m_diskIoReadBytes = prometheus::BuildCounter().Name(data_disk_io_read_bytes).Add(labels).Register(*registry);
//    m_diskIoWriteBytes = prometheus::BuildCounter().Name(data_disk_io_write_bytes).Add(labels).Register(*registry);

//    m_NetFlowRxBytes = prometheus::BuildCounter().Name(data_net_flow_rx_bytes).Add(labels).Register(*registry);
//    m_NetFlowTxBytes = prometheus::BuildCounter().Name(data_net_flow_tx_bytes).Add(labels).Register(*registry);

//    m_NetFlowRxPakcages = prometheus::BuildCounter().Name(data_net_flow_rx_pakcages).Add(labels).Register(*registry);
//    m_NetFlowTxPakcages = prometheus::BuildCounter().Name(data_net_flow_tx_pakcages).Add(labels).Register(*registry);

//    m_NetFlowTotalPakcages = prometheus::BuildCounter().Name(data_net_flow_total_pakcages).Add(labels).Register(*registry);
//    m_NetFlowTotalBytes = prometheus::BuildCounter().Name(data_net_flow_total_bytes).Add(labels).Register(*registry);
//}

void BaseMetric::AddCounter()
{

//    pid_t pid = getpid();
//    //double pmem = ProResource::get_pmem(pid);
//    double pmem  =  tools::process::GetProcessMemUsage(pid);
//    int threadNum = tgse::getThreadNum(pid);
//    m_mem.Increment(pmem);
//    m_cpu.Increment();
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
//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_cpu, label,
//        std::make_pair(m_cpuUsage, 0) , Gauge);

//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_mem, label,
//        std::make_pair(pmem, 0) , Gauge);

//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_thread, label,
//        std::make_pair(threadNum, 0) , Gauge);

//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_disk_io_read_bytes, label,
//        std::make_pair(m_diskIoMetric.read_bytes, 0) , Gauge);

//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_disk_io_write_bytes, label,
//        std::make_pair(m_diskIoMetric.write_bytes, 0) , Gauge);


//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_net_flow_tx_bytes, label,
//        std::make_pair(m_netMetric.net_wbytes, 0) , Gauge);


//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_net_flow_rx_bytes, label,
//        std::make_pair(m_netMetric.net_rbytes, 0) , Gauge);


//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_net_flow_rx_pakcages, label,
//        std::make_pair(m_netMetric.net_rpackets, 0) , Gauge);


//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_net_flow_tx_pakcages, label,
//        std::make_pair(m_netMetric.net_wpackets, 0) , Gauge);

//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_net_flow_total_pakcages, label,
//        std::make_pair(m_netMetric.net_packets, 0) , Gauge);


//    PrometheusHttpThread::instance().updateMetrics(MetricsUpdate, data_net_flow_total_bytes, label,
//        std::make_pair(m_netMetric.net_bytes, 0) , Gauge);
//
}
}
}
