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

#include "dataserver.h"
#include <gperftools/malloc_extension.h>

#include "conf/confItem.h"

#include "opscollection/ops.h"
//#include "dbgshell/dbgshell.h"
//#include "dbgshell/dbgshell.h"
#include "log/log.h"
#include "time_center.h"
#include "opscollection/ops.h"

#include "utils.h"

namespace gse { 
namespace dataserver {
static bool kUpgradeFlag = false;

DataServer::DataServer()
{
    m_dataFlow = NULL;
    m_conftor = NULL;
    m_opsReport = NULL;
}

DataServer::~DataServer()
{
    if (m_dataFlow != NULL)
    {
        delete m_dataFlow;
        m_dataFlow = NULL;
    }

    if (m_conftor != NULL)
    {
        delete m_conftor;
        m_conftor = NULL;
    }
}

DataServer &DataServer::Instance()
{
    static DataServer inst;
    return inst;
}
bool DataServer::GetUpgradeFlag()
{
    return kUpgradeFlag;
}

OpsCollection* DataServer::GetOpsReportClient()
{
    return m_opsReport;
}

void DataServer::StartOpsReportClient(DataFlowConf* ptr_dataflow, DataProcessConfig &gse_cfg)
{
    //start ops
    //-------------------------------------------
    std::string self_ip = gse_cfg.m_serverIp;
    if (self_ip == "")
    {
        self_ip = gse::tools::net::GetMachineIp();
    }

    OpsConf ops_conf;
    ptr_dataflow->GetOpsConf(ops_conf);
    m_opsReport = new OpsCollection();
    m_opsReport->SetGseConf(gse_cfg);

    m_opsReport->Open(ops_conf.m_maxQueneSize, ops_conf.m_opsThreadCount, gse_cfg.getOpsFlag(),
                      self_ip,  ops_conf.m_opsChannelID, ops_conf.m_opsAddress);
    m_opsReport->SetAgentOpsChannelid(ops_conf.m_agentOpsChannelID);
    //-------------------------------------------
}

int DataServer::Start(DataProcessConfig &conf, bool upgrade)
{
    // Sets the rate at which we release unused memory to the system.
    // Zero means we never release memory back to the system.  Increase
    // this flag to return memory faster; decrease it to return memory
    // slower.  Reasonable rates are in the range [0,10].  (Currently
    // only implemented in tcmalloc).
    MallocExtension::instance()->SetMemoryReleaseRate(2.0);
    LOG_ERROR("Malloc release Rate:%f", MallocExtension::instance()->GetMemoryReleaseRate());

    // start dataflow
    if (m_dataFlow != NULL)
    {
        m_dataFlow->Stop();
        delete m_dataFlow;
        m_dataFlow = NULL;
    }

//    if(PrometheusHttpThread::serviceSetup(conf.get_prometheus_http_svr_ip(), conf.get_prometheus_datasvr_port()) != GSE_SUCCESS)
//    {
//        LOG_ERROR("failed to start prometheus server");
//    }

    std::string confHost = conf.m_configZk;
    std::string serverIp = conf.m_serverIp;
    std::string clusterName = conf.m_clusterName;
    std::string instanceId = conf.m_instanceId;
    std::string dataflow_config = conf.m_dataflowfile;
    std::string zk_auth = conf.m_zkauth;
    std::string regionID = conf.m_regionID;
    std::string cityID = conf.m_cityID;

    std::string watchpath = conf.m_watchpath;
    std::string bkdata_zk_host = conf.m_eventZk;

    std::string channelid_zk_host = conf.m_channelidZkHost;
    std::string channelid_zk_auth = conf.m_channelidZkAuth;


    kUpgradeFlag = upgrade;
    LOG_DEBUG("data server start, upgrade flag:%d", upgrade);
    m_conftor = new Configurator();
 
    if (GSE_SUCCESS != m_conftor->Init(dataflow_config, confHost, zk_auth, bkdata_zk_host,
                                       channelid_zk_host, channelid_zk_auth, serverIp, regionID, cityID, clusterName, instanceId, watchpath))
    {
        LOG_ERROR("fail to initialize the configurator");
        return GSE_ERROR;
    }

    // get dataflow conf
    DataFlowConf *pDataFlowConf = NULL;
    pDataFlowConf = m_conftor->GetDataFlowConf();
    if (NULL == pDataFlowConf)
    {
        LOG_ERROR("it is failed to get data flow config from zk");
        delete m_conftor;
        m_conftor = NULL;
        return GSE_ERROR;
    }

    StartOpsReportClient(pDataFlowConf, conf);

    m_dataFlow = new DataFlow();
    m_dataFlow->SetGseConf(conf);
    m_dataFlow->UpdateConf(pDataFlowConf);
    m_dataFlow->UpdateIDToStorage(m_conftor);

    if (GSE_SUCCESS != m_dataFlow->Start())
    {
        LOG_ERROR("fail to start data flow");
        delete m_conftor;
        m_conftor = NULL;
        delete m_dataFlow;
        m_dataFlow = NULL;
        return GSE_ERROR;
    }

//    DbgShellHandler::ServiceSetup(conf);

//    OpsDbgShell::Instance().SetConftor(m_conftor);
//    OpsDbgShell::Instance().Init();

    OPMetric::Instance().SetGseConf(conf);
    OPMetric::Instance().MetricEventInit();
    LOG_INFO("data server start");

    return GSE_SUCCESS;
}

int DataServer::Stop()
{
    if (m_dataFlow != NULL)
    {
        m_dataFlow->Stop();
    }
    TimeCenter::Instance()->Stop();
    LOG_INFO("data server stoped");

    return GSE_SUCCESS;
}

void DataServer::Join()
{
    if (m_dataFlow != NULL)
    {
        m_dataFlow->Join();
    }
}


}
}
