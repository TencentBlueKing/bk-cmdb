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

#include "gseprocess.h"

#include <event2/thread.h>
#include <iostream>
#include <unistd.h>


#include "log/log.h"
#include "tools/strings.h"
#include "tools/time.h"
#include "safe/ssl.h"

#include "conf/dataconf.h"
#include "dataserver.h"



//#include "config/client.h"
#include "tools/json_property.hpp"
using namespace gse::config;
using namespace gse::tools::json;
using namespace gse::dataserver;

GseProcess::GseProcess()
{
}

GseProcess::~GseProcess()
{
    //
}
/*
    std::string m_configZk;
    std::string m_eventZk;
    std::string m_zkauth;
    std::string m_clusterName;
    std::string m_instanceId;
    std::string m_dataflowfile;
    std::string m_regionID;
    std::string m_cityID;
    std::string m_watchpath;
    std::string m_certPath;
    std::string m_passwdfile;
    std::string m_serverIp;

    std::string m_runtimePath;
    std::string m_logPath;
    std::string m_logLevel;
    unsigned int m_filesize;
    unsigned int m_filenum;


    m_dataflowfile = cfgValue.get("dataflow", "").asString();
    m_eventZk = cfgValue.get("eventzkhost", "").asString();
    m_configZk = cfgValue.get("zkhost", "").asString();
    m_zkauth = cfgValue.get("zkauth", "").asString();
    m_serverIp = cfgValue.get("datasvrip", "0.0.0.0").asString();

    m_regionID = cfgValue.get("dftregid", "").asString();
    m_cityID = cfgValue.get("dftcityid", "").asString();
    m_passwdfile = cfgValue.get("password_keyfile", "").asString();
    m_certPath = cfgValue.get("cert", "").asString();
    m_watchpath = cfgValue.get("watchpath", "").asString();
    m_logPath = cfgValue.get("log", "").asString();
    m_filesize = cfgValue.get("logfilesize", 10).asInt();
    m_filenum = cfgValue.get("logfilenum", 100).asInt();
    m_logLevel = cfgValue.get("level", "debug").asString();

*/
/*
const MetaType *DataProcessConfig::properties[] = {
    new Property<DataProcessConfig, std::string>(&DataProcessConfig::m_dataflowfile, "dataflow"),
    new Property<DataProcessConfig, std::string>(&DataProcessConfig::m_eventZk, "eventzkhost"),
    new Property<DataProcessConfig, std::string>(&DataProcessConfig::m_configZk, "zkhost"),
    new Property<DataProcessConfig, std::string>(&DataProcessConfig::m_zkauth, "zkauth"),
    new Property<DataProcessConfig, std::string>(&DataProcessConfig::m_serverIp, "datasvrip"),
    new Property<DataProcessConfig, std::string>(&DataProcessConfig::m_regionID, "dftregid"),
    new Property<DataProcessConfig, std::string>(&DataProcessConfig::m_cityID, "dftcityid"),
    new Property<DataProcessConfig, std::string>(&DataProcessConfig::m_passwdfile, "password_keyfile"),
    new Property<DataProcessConfig, std::string>(&DataProcessConfig::m_certPath, "cert"),
    new Property<DataProcessConfig, std::string>(&DataProcessConfig::m_watchpath, "watchpath"),
    new Property<DataProcessConfig, std::string>(&DataProcessConfig::m_logPath, "log"),
    new Property<DataProcessConfig, unsigned int>(&DataProcessConfig::m_filesize, "logfilesize"),
    new Property<DataProcessConfig, unsigned int>(&DataProcessConfig::m_filenum, "logfilenum"),
    new Property<DataProcessConfig, std::string>(&DataProcessConfig::m_logLevel, "level")};

*/

void GseProcess::ServiceStart(int id)
{
    evthread_use_pthreads();
    std::string errorMsg;
    //gse::config::Client client;

    //DataProcessConfig processcfg = client.LoadFromFile<DataProcessConfig>(m_configfile, errorMsg);

    //std::shared_ptr<DataProcessConfig> processcfg = std::make_shared<DataProcessConfig>(config);
    DataProcessConfig processcfg(m_configfile);
    std::string name = "data";
    if (!gse::log::Log::Instance().InitLog(name, processcfg.m_logLevel,
                                           processcfg.m_logPath,
                                           processcfg.m_filesize,
                                           processcfg.m_filenum))
    {
        fprintf(stderr, "Init log failed, level:%s, logpath:%s, log file size:%d, log file num:%d",
                processcfg.m_logLevel.c_str(), processcfg.m_logPath.c_str(), processcfg.m_filesize, processcfg.m_filenum);
        return;
    }

    SetProcPidPath(processcfg.m_pidFilePath);
    SetProcessName(name);
    bool upgrade = GetUpgradeFlag();
    LOG_DEBUG("DataServer serivce start, upgrade flag:%d", upgrade);
    m_serverptr = std::make_shared<DataServer>();
    m_serverptr->Start(processcfg, upgrade);
    m_serverptr->Join();

    LOG_INFO("gse data server exit");
    return;
}

void GseProcess::ServiceStop()
{
    DataServer::Instance().Stop();
}

void GseProcess::ServiceHandle()
{
    return;
}

const std::string GseProcess::Version()
{
    return "TODO:get version from compile macro";
}

void GseProcess::Reconfigure()
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

void GseProcess::AddPrivateOptions(po::options_description &config)
{
    config.add_options()("file,f", po::value<string>(), "set config path...");
}

void GseProcess::HandlePrivateFlags(po::variables_map &vm)
{
    if (vm.count("file"))
    {
        m_configfile = vm["file"].as<string>();
    }
}

int32_t main(int argc, char **argv)
{
    GseProcess *process = new GseProcess();
    process->Start(argc, argv);
    return GSE_SUCCESS;
}
