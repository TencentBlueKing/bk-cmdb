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
#include "loadbalance_schedule.h"

#include "ops/ops.h"

#include "api/channelid_delete_api.h"
#include "api/channelid_query_api.h"
#include "api/channelid_request_api.h"
#include "api/channelid_update_api.h"

#include "api/streamto_delete_api.h"
#include "api/streamto_query_api.h"
#include "api/streamto_request_api.h"
#include "api/streamto_update_api.h"

#include "discover/zkapi/zk_api.h"
#include "log/log.h"
#include "ops/op_healthz.h"
#include "time_center.h"
#include "tools/maths.h"
#include "tools/system.h"

#include "utils.h"

namespace gse {
namespace data {
static std::string gNodeId;

DataServer::DataServer()
    : m_upgradeFlag(false),
      m_dataFlow(nullptr),
      m_conftor(nullptr),
      m_opsReport(nullptr),
      m_opsHttpServer(nullptr),
      m_opsHealthZ(nullptr),
      m_scheduler(nullptr),
      m_controlMsgProcessor(nullptr),
      m_dataReportMsgProcessor(nullptr),
      m_opsMsgProcessor(nullptr),
      m_channelIdApiSerivce(nullptr),
      m_channelIdZkClient(nullptr),
      m_discoverZkClient(nullptr),
      m_channelIdApiZkClient(nullptr),
      m_channelIdApiMetrics(nullptr),
      m_channelIdCommonApi(nullptr),
      m_exit(false)

{
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
    return DataServer::Instance().UpgradeFlag();
}

bool DataServer::UpgradeFlag()
{
    return m_upgradeFlag;
}
OpsCollection *DataServer::OpsReportClient()
{
    return m_opsReport;
}

OpsCollection *DataServer::GetOpsReportClient()
{
    return DataServer::Instance().OpsReportClient();
}

void DataServer::StartOpsReportClient()
{
    // start ops
    //-------------------------------------------
    std::string nodeIp = DataServer::GetConfigPtr()->GetAdvertiseIp();
    if (nodeIp == "")
    {
        nodeIp = gse::tools::net::GetMachineIp();
    }
    // TODO:

    OpsConf opsConf;
    m_opsReport = new OpsCollection();
    m_opsReport->SetGseConf(m_configPtr);

    m_opsReport->Open(opsConf.m_maxQueneSize, opsConf.m_opsThreadCount, m_configPtr->GetOpsFlag(),
                      nodeIp, opsConf.m_opsChannelID, opsConf.m_opsAddress);
    m_opsReport->SetAgentOpsChannelid(opsConf.m_agentOpsChannelID);

    //-------------------------------------------
}
const std::string DataServer::GetNodeId()
{
    return gNodeId;
}

const std::string DataServer::GetVersion()
{
    std::string version;
#ifdef _VERSION_
    version.assign(std::string(_VERSION_));
#else
    version.assign("1.0");
#endif

#ifndef _GSE_SSL_
    version.append("i");
#endif
    return version;
}

void DataServer::GracefullyQuit()
{
    LOG_INFO("gracefully quit data servcie");
    gse::tools::system::Sleep(120);
    _exit(128);
    return;
}

void DataServer::InitMsgProcessor()
{
    m_controlMsgProcessor = std::make_shared<SignalControlMsgProcessor>();
    m_controlMsgProcessor->SetScheduler(m_scheduler);

    m_dataReportMsgProcessor = std::make_shared<DataReportMsgProcessor>();
    m_opsMsgProcessor = std::make_shared<OpsReportMsgProcessor>();
}

bool DataServer::StartZkSerivce(std::shared_ptr<DataProcessConfig> configPtr)
{
    auto zookeeperConfigPtr = configPtr->GetZookeeperConfig();
    if (!zookeeperConfigPtr->m_serviceDiscoverZkHost.empty())
    {
        m_discoverZkClient = CreateZkClient(zookeeperConfigPtr->m_serviceDiscoverZkHost,
                                            zookeeperConfigPtr->m_serviceDiscoverZkAuth,
                                            zookeeperConfigPtr->m_timeout);
        if (m_discoverZkClient == nullptr)
        {
            LOG_ERROR("failed to init discover zk client, zkhost:%s", zookeeperConfigPtr->m_serviceDiscoverZkHost.c_str());
            return false;
        }

        OpHealthZ::AddInitHealthInfo("zookeeper", "ok", 0);
    }

    if (!zookeeperConfigPtr->m_channelIdConfigZkHost.empty())
    {
        m_channelIdZkClient = CreateZkClient(zookeeperConfigPtr->m_channelIdConfigZkHost,
                                             zookeeperConfigPtr->m_channelIdConfigZkAuth,
                                             zookeeperConfigPtr->m_timeout);

        if (m_channelIdZkClient == nullptr)
        {
            LOG_ERROR("failed to init channelid zk client, zkhost:%s", zookeeperConfigPtr->m_channelIdConfigZkHost.c_str());
            return false;
        }

        OpHealthZ::AddInitHealthInfo("channelid_zookeeper", "ok", 0);
    }
    else
    {
        m_channelIdZkClient = m_discoverZkClient;
    }

    auto dataIdZkConfig = configPtr->GetDataIdZkConfig();
    for (auto dataIdZk : *dataIdZkConfig)
    {
        auto dataIdZkClient = std::make_shared<DataIdZkClient>();
        dataIdZkClient->m_ptrEventZkClient = CreateZkClient(dataIdZk.m_zkhost, dataIdZk.m_zkAuth, kDefaultZkTimeoutMs);
        if (dataIdZkClient->m_ptrEventZkClient == nullptr)
        {
            LOG_ERROR("failed to init dataid zk client, zkhost:%s", dataIdZk.m_zkhost.c_str());
            return false;
        }

        dataIdZkClient->m_dataIdZkPath = dataIdZk.m_dataIdPath;
        dataIdZkClient->m_storageZkPath = dataIdZk.m_storagePath;
        dataIdZkClient->m_zkAuth = dataIdZk.m_zkAuth;
        m_dataIdZkClients.push_back(dataIdZkClient);
        dataIdZkClient->Init();
        LOG_DEBUG("init dataid zkhost:%s, dataid path:%s, storage path:%s",
                  dataIdZk.m_zkhost.c_str(),
                  dataIdZk.m_dataIdPath.c_str(),
                  dataIdZk.m_storagePath.c_str());
    }

    return true;
}

void DataServer::StopZkService()
{
    m_discoverZkClient->ApiClose();
    m_discoverZkClient->ApiJoin();
    LOG_DEBUG("stop discover zookeeper service");

    m_channelIdZkClient->ApiClose();
    m_channelIdZkClient->ApiJoin();
    LOG_DEBUG("stop channelid zookeeper service");

    return;
}

void DataServer::StopConftor()
{
    m_conftor->StopConftor();
}

std::shared_ptr<gse::discover::zkapi::ZkApi> DataServer::CreateZkClient(const std::string &zkHost, const std::string &zkAuth, int timeoutMs)
{
    auto client = std::make_shared<gse::discover::zkapi::ZkApi>();

    int flags = 0;
    int ret = client->ApiSetup();
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("failed tto setup zk client");
        return nullptr;
    }

    ret = client->ZkInit(zkHost, NULL, timeoutMs, -1, NULL, this, flags, zkAuth);
    if (ret != GSE_SUCCESS)
    {
        return nullptr;
    }

    return client;
}

bool DataServer::StartChannelIdApiService()
{
    auto channelIdApiService = DataServer::GetConfigPtr()->GetChannelIdApiServiceConfig();

    if (channelIdApiService && channelIdApiService->m_valid)
    {
        // init zk service
        auto zookeeperConfig = DataServer::GetConfigPtr()->GetZookeeperConfig();
        bool channelIdZkAcl = zookeeperConfig->m_channelIdConfigZkAuth.empty() ? false : true;
        m_channelIdApiZkClient = std::make_shared<ChannelIdZkAPI>(m_channelIdZkClient, channelIdZkAcl);
        if (!m_channelIdApiZkClient->Init())
        {
            LOG_ERROR("failed to initialize channelid api zk client");
            return false;
        }
        LOG_INFO("start channelid zkapi service success");

        // init api servcie
        m_channelIdCommonApi = std::make_shared<ChannelIdCommonAPI>();

        m_channelIdApiSerivce = std::make_shared<HTTPAPIService>(channelIdApiService->m_bindIp,
                                                                 channelIdApiService->m_bindPort,
                                                                 channelIdApiService->m_threadNum);

        auto channelIdRequestHandler = std::make_shared<ChannelIdRequestAPIHandler>("v1");
        auto channelIdUpdateHandler = std::make_shared<ChannelIdUpdateAPIHandler>("v1");
        auto channelIdDeleteHandler = std::make_shared<ChannelIdDeleteAPIHandler>("v1");
        auto channelIdQueryHandler = std::make_shared<ChannelIdQueryAPIHandler>("v1");

        m_channelIdApiSerivce->RegisterHttpURI(ChannelIdRequsetHandlerURI, channelIdRequestHandler);
        m_channelIdApiSerivce->RegisterHttpURI(ChannelIdUpdateHandlerURI, channelIdUpdateHandler);
        m_channelIdApiSerivce->RegisterHttpURI(ChannelIdDeleteHandlerURI, channelIdDeleteHandler);
        m_channelIdApiSerivce->RegisterHttpURI(ChannelIdQueryHandlerURI, channelIdQueryHandler);

        auto channelIdRequestHandlerV2 = std::make_shared<ChannelIdRequestAPIHandler>();
        auto channelIdUpdateHandlerV2 = std::make_shared<ChannelIdUpdateAPIHandler>();
        auto channelIdDeleteHandlerV2 = std::make_shared<ChannelIdDeleteAPIHandler>();
        auto channelIdQueryHandlerV2 = std::make_shared<ChannelIdQueryAPIHandler>();

        m_channelIdApiSerivce->RegisterHttpURI(ChannelIdRequsetHandlerURIV2, channelIdRequestHandlerV2);
        m_channelIdApiSerivce->RegisterHttpURI(ChannelIdUpdateHandlerURIV2, channelIdUpdateHandlerV2);
        m_channelIdApiSerivce->RegisterHttpURI(ChannelIdDeleteHandlerURIV2, channelIdDeleteHandlerV2);
        m_channelIdApiSerivce->RegisterHttpURI(ChannelIdQueryHandlerURIV2, channelIdQueryHandlerV2);

        auto streamToIdRequestHandler = std::make_shared<ChannelStreamToRequestAPIHandler>("v1");
        auto streamToIdUpdateHandler = std::make_shared<ChannelStreamToUpdateAPIHandler>("v1");
        auto streamToIdDeleteHandler = std::make_shared<StreamToDeleteAPIHandler>("v1");
        auto streamToIdQueryHandler = std::make_shared<StreamToIdQueryAPIHandler>("v1");

        m_channelIdApiSerivce->RegisterHttpURI(StreamToIdRequsetHandlerURI, streamToIdRequestHandler);
        m_channelIdApiSerivce->RegisterHttpURI(StreamToIdUpdateHandlerURI, streamToIdUpdateHandler);
        m_channelIdApiSerivce->RegisterHttpURI(StreamToDeleteHandlerURI, streamToIdDeleteHandler);
        m_channelIdApiSerivce->RegisterHttpURI(StreamToIdQueryHandlerURI, streamToIdQueryHandler);

        auto streamToIdRequestHandlerV2 = std::make_shared<ChannelStreamToRequestAPIHandler>();
        auto streamToIdUpdateHandlerV2 = std::make_shared<ChannelStreamToUpdateAPIHandler>();
        auto streamToIdDeleteHandlerV2 = std::make_shared<StreamToDeleteAPIHandler>();
        auto streamToIdQueryHandlerV2 = std::make_shared<StreamToIdQueryAPIHandler>();

        m_channelIdApiSerivce->RegisterHttpURI(StreamToIdRequsetHandlerURIV2, streamToIdRequestHandlerV2);
        m_channelIdApiSerivce->RegisterHttpURI(StreamToIdUpdateHandlerURIV2, streamToIdUpdateHandlerV2);
        m_channelIdApiSerivce->RegisterHttpURI(StreamToDeleteHandlerURIV2, streamToIdDeleteHandlerV2);
        m_channelIdApiSerivce->RegisterHttpURI(StreamToIdQueryHandlerURIV2, streamToIdQueryHandlerV2);

        m_channelIdApiSerivce->Start();

        LOG_INFO("start channelid http api service success");

        return true;
    }
    else
    {
        return true;
    }
}

int DataServer::StartDataflow()
{
    DataFlowConf *pDataFlowConf = NULL;
    pDataFlowConf = m_conftor->GetDataFlowConf();
    if (NULL == pDataFlowConf)
    {
        LOG_ERROR("it is failed to get data flow config from zk");
        delete m_conftor;
        m_conftor = NULL;
        return GSE_ERROR;
    }

    m_dataFlow = new DataFlow();
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

    return GSE_SUCCESS;
}

void DataServer::StopDataflow()
{
    if (m_dataFlow != NULL)
    {
        m_dataFlow->Stop();
    }
}
int DataServer::Start(std::shared_ptr<DataProcessConfig> configPtr, bool upgrade)
{
    // Sets the rate at which we release unused memory to the system.
    // Zero means we never release memory back to the system.  Increase
    // this flag to return memory faster; decrease it to return memory
    // slower.  Reasonable rates are in the range [0,10].  (Currently
    // only implemented in tcmalloc).
    MallocExtension::instance()->SetMemoryReleaseRate(2.0);
    LOG_INFO("Malloc release Rate:%f", MallocExtension::instance()->GetMemoryReleaseRate());

    TimeCenter::Instance()->Start();

    gNodeId = std::move(gse::tools::maths::GetRandomUuid());
    LOG_DEBUG("Data server node id:%s", gNodeId.c_str());
    m_configPtr = configPtr;
    m_upgradeFlag = upgrade;
    m_advertiseIP = configPtr->GetAdvertiseIp();

    if (!StartZkSerivce(configPtr))
    {
        LOG_ERROR("failed to initialize zookeeper client");
        return GSE_ERROR;
    }

    StartStreamExporterManager();
    StartChannelIdManager();

    LOG_DEBUG("data server start, upgrade flag:%d", upgrade);

    // first start metrics http,wrap ops service func
    m_opsHealthZ = std::make_shared<OpHealthZ>();
    m_stackService = std::make_shared<Stack>();

    m_opsHttpServer = std::make_shared<OPSHttpServer>(configPtr->GetPrometheusBindIp(), configPtr->GetPrometheusListenerPort());
    m_opsHttpServer->Start();
    m_opsHttpServer->RegisterMsgHandler(m_opsHealthZ->GetURI(), m_opsHealthZ);
    m_opsHttpServer->RegisterMsgHandler(m_stackService->GetURI(), m_stackService);

    OPMetric::OPMetricInst()->MetricEventInit(m_opsHttpServer->GetHttpServer());

    bool discoverZkAcl = configPtr->GetZookeeperConfig()->m_serviceDiscoverZkAuth.empty() ? false : true;
    bool channelIdZkAcl = configPtr->GetZookeeperConfig()->m_channelIdConfigZkAuth.empty() ? false : true;

    // TODO: need reconstruct
    m_conftor = new Configurator();
    m_conftor->SetDiscoverZkClient(m_discoverZkClient, discoverZkAcl);
    m_conftor->SetChannelIdZkClient(m_channelIdZkClient, channelIdZkAcl);

    if (GSE_SUCCESS != m_conftor->Init(configPtr))
    {
        LOG_ERROR("failed to initialize the configurator");

        return GSE_ERROR;
    }

    StartScheduler();
    InitMsgProcessor();

    if (GSE_SUCCESS != StartDataflow())
    {
        LOG_ERROR("failed to initialize dataflow");
        return GSE_ERROR;
    }

    StartOpsReportClient();

    // channelMetrics
    m_channelIdApiMetrics = std::make_shared<ChannelIdApiMetrics>(configPtr);
    m_channelIdApiMetrics->Start(m_opsHttpServer->GetHttpServer());
    if (!StartChannelIdApiService())
    {
        LOG_ERROR("failed to start channelid http serivice");
        return GSE_ERROR;
    }

    LOG_INFO("data server start");

    GuardLoop();
    return GSE_SUCCESS;
}

int DataServer::Stop()
{
    StopScheduler();
    StopChannelIDAPIService();
    StopOpsServiceAPISerivce();
    StopZkService();
    // TOOD: not impls stop
    //     TimeCenter::Instance()->Stop();
    //     TimeCenter::Instance()->Join();

    StopDataflow();
    StopChannelIdManager();
    LOG_DEBUG("channelid manager stopped");
    StopStreamExporterManager();
    StopConftor();

    LOG_INFO("data server stoped");
    m_exit = true;
    return GSE_SUCCESS;
}

void DataServer::GuardLoop()
{
    LOG_DEBUG("data server main guard loop is running now");

    while (!m_exit)
    {
        // NOTE: add some others async work here, and blocking here as
        // the main service guard loop.
        gse::tools::system::Sleep(1);
    }
    LOG_WARN("data server main guard loop is stopped");
}

void DataServer::Join()
{
    if (m_dataFlow != NULL)
    {
        m_dataFlow->Join();
    }

    if (m_opsHttpServer != nullptr)
    {
        m_opsHttpServer->Join();
        LOG_DEBUG("ops http service joined");
    }
}

bool DataServer::StartScheduler()
{
    auto zkConfig = DataServer::GetConfigPtr()->GetZookeeperConfig();
    if (zkConfig->m_serviceDiscoverZkHost.empty())
    {
        LOG_INFO("no discover zk, don't start scheduler module");
        return false;
    }

    m_scheduler = std::make_shared<LoadBalanceSchedule>(zkConfig->m_serviceDiscoverZkHost, zkConfig->m_serviceDiscoverZkAuth);
    m_scheduler->Start();

    LOG_INFO("start scheduler module success");

    return true;
}

void DataServer::StopScheduler()
{
    m_scheduler->Stop();
}

bool DataServer::StartStreamExporterManager()
{
    m_streamExporterManager = std::make_shared<ChannelIdStreamExporterManager>();
    m_streamExporterManager->Start();

    m_streamExporterLoader = std::make_shared<StreamExporterIDLoader>(m_channelIdZkClient);
    m_streamExporterLoader->SetStreamExporterManager(m_streamExporterManager);
    m_streamExporterLoader->LoadStreamExporterConfig();

    return true;
}

void DataServer::StopStreamExporterManager()
{
    m_streamExporterManager->Stop();
    m_streamExporterManager->Join();
    return;
}

bool DataServer::StartChannelIdManager()
{
    m_channelIdManager = std::make_shared<ChannelIdManager>();
    m_channelIdManager->Start();
    return true;
}

void DataServer::StopChannelIdManager()
{
    m_channelIdManager->Stop();
    m_channelIdManager->Join();
}

void DataServer::StopChannelIDAPIService()
{
    m_channelIdApiSerivce->Stop();
}

void DataServer::StopOpsServiceAPISerivce()
{
    m_opsHttpServer->Stop();
}

std::shared_ptr<LoadBalanceSchedule> DataServer::Scheduler()
{
    return m_scheduler;
}

std::shared_ptr<LoadBalanceSchedule> DataServer::GetScheduler()
{
    return DataServer::Instance().Scheduler();
}

std::shared_ptr<ChannelIdZkAPI> DataServer::ChannelIdZkApiInst()
{
    return m_channelIdApiZkClient;
}

std::shared_ptr<ChannelIdZkAPI> DataServer::GetChannelIdZkApiInst()
{
    return DataServer::Instance().ChannelIdZkApiInst();
}

std::shared_ptr<gse::discover::zkapi::ZkApi> DataServer::DiscoverZkApiClient()
{
    return m_discoverZkClient;
}

std::shared_ptr<gse::discover::zkapi::ZkApi> DataServer::GetDiscoverZkApiClientInst()
{
    return DataServer::Instance().DiscoverZkApiClient();
}

std::shared_ptr<ChannelIdApiMetrics> DataServer::ChannelIdApiMetricsInst()
{
    return m_channelIdApiMetrics;
}

std::shared_ptr<ChannelIdApiMetrics> DataServer::GetChannelIdApiMetricsInst()
{
    return DataServer::Instance().ChannelIdApiMetricsInst();
}

std::vector<std::shared_ptr<DataIdZkClient>> DataServer::DataIdZkApiInsts()
{
    return m_dataIdZkClients;
}

std::vector<std::shared_ptr<DataIdZkClient>> DataServer::GetDataIdZkApiInst()
{
    return DataServer::Instance().DataIdZkApiInsts();
}

std::shared_ptr<ChannelIdCommonAPI> DataServer::ChannelIdCommonApiInst()
{
    return m_channelIdCommonApi;
}

std::shared_ptr<ChannelIdCommonAPI> DataServer::GetChannelIdCommonApiInst()
{
    return DataServer::Instance().ChannelIdCommonApiInst();
}

std::shared_ptr<ChannelIdStreamExporterManager> DataServer::StreamExporterInst()
{
    return m_streamExporterManager;
}

std::shared_ptr<ChannelIdStreamExporterManager> DataServer::GetStreamExporterManagerInst()
{
    return DataServer::Instance().StreamExporterInst();
}

std::shared_ptr<ChannelIdManager> DataServer::ChannelIDMangerInst()
{
    return m_channelIdManager;
}
std::shared_ptr<ChannelIdManager> DataServer::GetChannelIdManagerInst()
{
    return DataServer::Instance().ChannelIDMangerInst();
}

std::string DataServer::GetAdvertiseIPInst()
{
    return m_advertiseIP;
}

std::string DataServer::GetAdvertiseIP()
{
    return DataServer::Instance().GetAdvertiseIPInst();
}

std::shared_ptr<DataProcessConfig> DataServer::GetConfigPtrInst()
{
    return m_configPtr;
}

std::shared_ptr<DataProcessConfig> DataServer::GetConfigPtr()
{
    return DataServer::Instance().GetConfigPtrInst();
}

} // namespace data
} // namespace gse
