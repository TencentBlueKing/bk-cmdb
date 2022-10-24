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

#include "service.h"

#include <event2/thread.h>
#include <iostream>
#include <unistd.h>

#include "version.hpp"

#include "log/log.h"
#include "safe/ssl.h"
#include "tools/strings.h"
#include "tools/time.h"

#include "conf/dataconf.h"
#include "dataserver.h"

namespace gse {
namespace data {

DataService::DataService()
{
}

DataService::~DataService()
{
}

void DataService::ServiceStart(int id)
{
    evthread_use_pthreads();

    auto configPtr = std::make_shared<DataProcessConfig>();
    if (!configPtr->LoadConfig(m_configFile))
    {
        fprintf(stderr, "load config failed, config file:%s", m_configFile.c_str());
        return;
    }

    if (!gse::log::Log::Instance().InitLog("data", configPtr->m_loggerConfig.m_level,
                                           configPtr->m_loggerConfig.m_path,
                                           configPtr->m_loggerConfig.m_fileSize,
                                           configPtr->m_loggerConfig.m_fileCount,
                                           DFT_LOG_PATTERN,
                                           configPtr->m_loggerConfig.m_rotate,
                                           configPtr->m_loggerConfig.m_flushIntervalMs))
    {
        fprintf(stderr, "init log failed, level:%s, logpath:%s, log file size:%d, log file num:%d, rotate:%d, flush intreval:%d",
                configPtr->m_loggerConfig.m_level.c_str(),
                configPtr->m_loggerConfig.m_path.c_str(),
                configPtr->m_loggerConfig.m_fileSize,
                configPtr->m_loggerConfig.m_fileCount,
                configPtr->m_loggerConfig.m_rotate,
                configPtr->m_loggerConfig.m_flushIntervalMs);
        return;
    }

    SetProcPidPath(configPtr->m_pidFilePath);
    SetProcessName("data");
    bool upgrade = GetUpgradeFlag();
    LOG_DEBUG("service start, upgrade flag:%d", upgrade);

    DataServer::Instance().Start(configPtr, upgrade);
    LOG_INFO("service exit");
    return;
}

void DataService::ServiceStop()
{
    DataServer::Instance().Stop();
}

void DataService::ServiceHandle()
{
    return;
}

const std::string DataService::Version()
{
    std::string gitHash;
    std::string buildTime;
    std::string version;
    version = DataServer::GetVersion();
#ifdef _GITHASH_
    gitHash.assign(std::string(_GITHASH_));
#else
    gitHash.assign("");
#endif

#ifdef GSE_MAIN_VERSION_COMMIT_HASH
    gitHash.append("-");
    gitHash.append(GSE_MAIN_VERSION_COMMIT_HASH);
#endif

#ifdef _BUILDTIME_
    buildTime.assign(std::string(_BUILDTIME_));
#else
    buildTime.assign(gse::tools::time::LocalTimeToString);
#endif

    std::string verboseVersion;
    verboseVersion.append("Version: ");
    verboseVersion.append(version);
    verboseVersion.append("\n");
    verboseVersion.append("GitHash: ");
    verboseVersion.append(gitHash);
    verboseVersion.append("\n");
    verboseVersion.append("BuildTime: ");
    verboseVersion.append(buildTime);
    verboseVersion.append("\n");
    return verboseVersion;
}

void DataService::Reconfigure()
{
    if (m_logLevel.compare("debug") == 0)
    {
        std::cout << "debug --> warn" << std::endl;
        m_logLevel = "warn";
    }
    else if (m_logLevel.compare("warn") == 0)
    {
        std::cout << "warn --> error" << std::endl;
        m_logLevel = "error";
    }
    else if (m_logLevel.compare("error") == 0)
    {
        std::cout << "error --> debug" << std::endl;
        m_logLevel = "debug";
    }

    gse::log::Log::Instance().ResetLogLevel(m_logLevel);
}

void DataService::AddPrivateOptions(po::options_description &config)
{
    config.add_options()("config,f", po::value<std::string>()->implicit_value("./etc/gse_data.json"), "local config file");
}

void DataService::HandlePrivateFlags(po::variables_map &vm)
{
    m_configFile = "./etc/gse_data.json";
    if (vm.count("config"))
    {
        m_configFile = vm["config"].as<string>();
    }
}

} // namespace data
} // namespace gse
