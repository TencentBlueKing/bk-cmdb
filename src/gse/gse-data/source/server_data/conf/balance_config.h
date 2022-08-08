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

#ifndef _GSE_DATA_CONFIG_BALANCE_H_
#define _GSE_DATA_CONFIG_BALANCE_H_

#include <memory>
#include <thread>

#include <rapidjson/document.h>     // rapidjson's DOM-style API
#include <rapidjson/prettywriter.h> // for stringify JSON
#include <rapidjson/stringbuffer.h>

#include "dataconf.h"

namespace gse {
namespace data {

class SystemResourceMonitor final
{
public:
    SystemResourceMonitor() noexcept;
    SystemResourceMonitor(const std::string &ethName) noexcept;
    virtual ~SystemResourceMonitor();
    int Start();
    void Stop();
    void Join();

    double GetCpuUsage() const;
    double GetMemUsage() const;
    double GetNetUsage() const;

    void SetEthName(const std::string &ethName);
    void SetNetDevMaxSpeed(int maxSpeed);

private:
    void GetSystemResourceUsage();

    void ReportLoad();
    double CpuUsage();
    double MemUsage();
    double NetUsage();

    uint64_t GetMaxNetBandFromConfig();

private:
    double m_cpuUsage;
    double m_memUsage;
    double m_netUsage;

    uint64_t m_maxSpeed;
    std::string m_ethName;
    bool m_exit;

    std::thread m_systemResourceMonitorThread;
};

class SystemConnectionMonitor final
{
public:
    SystemConnectionMonitor();
    virtual ~SystemConnectionMonitor();

    uint64_t GetMaxConnectionCount() const;
    uint64_t GetConnectionCount();
    void SetMaxConnectionCount(uint64_t maxConnectionCount);

private:
    uint64_t m_maxConnectionCount;
};

} // namespace data
} // namespace gse

#endif
