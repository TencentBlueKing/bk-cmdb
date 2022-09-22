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

#include "balance_config.h"

#include "ops/op_metric.h"
#include "tools/hostinfo.h"
#include "tools/process.h"
#include "tools/system.h"
#include "utils.h"

namespace gse {
namespace data {

SystemResourceMonitor::SystemResourceMonitor() noexcept
    : m_cpuUsage(0.0), m_memUsage(0.0), m_netUsage(0.0),
      m_maxSpeed(kDefaultNetDevMaxSpeed),
      m_ethName(kDefaultDevName), m_exit(false)
{
}

SystemResourceMonitor::SystemResourceMonitor(const std::string &ethName) noexcept
    : m_cpuUsage(0.0), m_memUsage(0.0), m_netUsage(0.0),
      m_maxSpeed(kDefaultNetDevMaxSpeed),
      m_ethName(ethName), m_exit(false)
{
}

void SystemResourceMonitor::SetEthName(const std::string &ethName)
{
    m_ethName = ethName;
}

void SystemResourceMonitor::SetNetDevMaxSpeed(int maxSpeed)
{
    m_maxSpeed = maxSpeed;
}

SystemResourceMonitor::~SystemResourceMonitor()
{
}

void SystemResourceMonitor::GetSystemResourceUsage()
{
    while (!m_exit)
    {
        m_cpuUsage = CpuUsage();
        m_memUsage = MemUsage();
        m_netUsage = NetUsage();
        gse::tools::system::MillSleep(1000);
    }
}

double SystemResourceMonitor::MemUsage()
{
    return gse::tools::process::GetCurrentProcessMemUsage();
}

uint64_t SystemResourceMonitor::GetMaxNetBandFromConfig()
{
    return m_maxSpeed;
}

double SystemResourceMonitor::CpuUsage()
{
    uint64_t usedCpu;
    uint64_t totalCpu;
    double cpuUsage = 0.0;

    static uint64_t lastUsedCpu;
    static uint64_t lastTotalCpu;

    uint32_t pid = gse::tools::process::GetProcessID();
    usedCpu = gse::tools::process::GetProcessCpuTime(pid);
    totalCpu = gse::tools::process::GetTotalCpuTime();

    uint64_t dtotal = totalCpu - lastTotalCpu;
    if (0 != dtotal)
    {
        cpuUsage = ((usedCpu - lastUsedCpu) * 100.00) / dtotal;
    }

    lastUsedCpu = usedCpu;
    lastTotalCpu = totalCpu;
    // LOG_DEBUG("cpu usage:%f, last cpu usage:%d, current cpu usage:%d, total cpu usage:%d", cpuUsage, lastUsedCpu, usedCpu, dtotal);
    return cpuUsage;
}

double SystemResourceMonitor::NetUsage()
{
    static std::vector<gse::os::hostinfo::stats_net_dev> lastNetDevItems;
    static int lastTimeStamp = 0;

    std::vector<gse::tools::hostinfo::gse_stats_net_dev> currNetDevItems;
    gse::tools::hostinfo::ReadNetDev(currNetDevItems, m_ethName.c_str());

    if (currNetDevItems.size() <= 0)
    {
        LOG_WARN("failed to get eth device, ethname:%s", m_ethName.c_str());
        return 0.0;
    }

    if (lastNetDevItems.size() <= 0)
    {
        lastNetDevItems = currNetDevItems;
        LOG_DEBUG("system startup, calc net usage in next cycle");
        return 0.0;
    }

    float netReadBytes = currNetDevItems.at(0).rx_bytes - lastNetDevItems.at(0).rx_bytes;
    int nowtime = ::time(NULL);
    int dtimestamp = nowtime - lastTimeStamp;
    if (dtimestamp <= 0)
    {
        LOG_WARN("diff time invalid, now:%d, last time:%d", nowtime, lastTimeStamp);
        return 0.0;
    }

    lastTimeStamp = nowtime;
    lastNetDevItems = std::move(currNetDevItems);
    const int bits = 8;
    float netUsagePer = (netReadBytes * bits / 1024 / 1024) / dtimestamp;
    return (netUsagePer * 100 / m_maxSpeed);
}

int SystemResourceMonitor::Start()
{
    int speed = GetNetDevSpeed(m_ethName.c_str());
    if (speed > 0)
    {
        m_maxSpeed = speed;
    }

    LOG_INFO("get netusage ethname:%s, max speed:%d", m_ethName.c_str(), m_maxSpeed);

    auto getSystemResource = std::bind(&SystemResourceMonitor::GetSystemResourceUsage, this);
    m_systemResourceMonitorThread = std::thread(getSystemResource);
}

void SystemResourceMonitor::Stop()
{
    m_exit = true;
}

void SystemResourceMonitor::Join()
{
    if (m_systemResourceMonitorThread.joinable())
    {
        m_systemResourceMonitorThread.join();
    }
}

double SystemResourceMonitor::GetCpuUsage() const
{
    return m_cpuUsage;
}

double SystemResourceMonitor::GetMemUsage() const
{
    return m_memUsage;
}

double SystemResourceMonitor::GetNetUsage() const
{
    return m_netUsage;
}

SystemConnectionMonitor::SystemConnectionMonitor()
    : m_maxConnectionCount(kDefaultMaxAgentCount)
{
}

SystemConnectionMonitor::~SystemConnectionMonitor()
{
}

void SystemConnectionMonitor::SetMaxConnectionCount(uint64_t maxConnectionCount)
{
    m_maxConnectionCount = maxConnectionCount;
}

uint64_t SystemConnectionMonitor::GetConnectionCount()
{
    return OPMetric::GetConnectionCount();
}

uint64_t SystemConnectionMonitor::GetMaxConnectionCount() const
{
    return m_maxConnectionCount;
}

} // namespace data

} // namespace gse
