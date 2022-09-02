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

#ifndef _GSE_DATA_SERVER_H_
#define _GSE_DATA_SERVER_H_

#include "conf/configurator.h"
#include "conf/dataconf.h"
#include "dataflow.h"
#include "ops/op_healthz.h"
#include "ops/ops.h"

#include "api/api_metrics.h"
#include "api/channelid_zk_api.h"
#include "api/data_api_serivce.h"
#include "conf/stream_exporter_loader.h"
#include "discover/zkapi/zk_api.h"
#include "loadbalance_schedule.h"
#include "ops/stack.h"
#include "receiver/tcpreceiver.h"

namespace gse {
namespace data {

class DataServer
{
public:
    virtual ~DataServer();
    DataServer();

public:
    static DataServer& Instance();

public:
    int Start(std::shared_ptr<DataProcessConfig> configPtr, bool upgrade);
    int Stop();
    void Join();

public:
    static OpsCollection* GetOpsReportClient();
    static const std::string GetNodeId();
    static const std::string GetVersion();

public:
    // api
    static std::shared_ptr<LoadBalanceSchedule> GetScheduler();
    static std::shared_ptr<ChannelIdZkAPI> GetChannelIdZkApiInst();
    static std::shared_ptr<gse::discover::zkapi::ZkApi> GetDiscoverZkApiClientInst();

    static std::shared_ptr<ChannelIdApiMetrics> GetChannelIdApiMetricsInst();
    static std::vector<std::shared_ptr<DataIdZkClient>> GetDataIdZkApiInst();
    static std::shared_ptr<ChannelIdCommonAPI> GetChannelIdCommonApiInst();

    static std::shared_ptr<ChannelIdStreamExporterManager> GetStreamExporterManagerInst();
    static std::shared_ptr<ChannelIdManager> GetChannelIdManagerInst();
    static bool GetUpgradeFlag();
    static void GracefullyQuit();

    static std::string GetAdvertiseIP();
    static std::shared_ptr<DataProcessConfig> GetConfigPtr();

private:
    bool StartZkSerivce(std::shared_ptr<DataProcessConfig> configPtr);
    // load channelid and stream
    bool StartChannelIdManager();
    bool StartStreamExporterManager();

    // ops service
    void StartOpsReportClient();
    // start Schededuler
    bool StartScheduler();
    bool StartChannelIdApiService();
    void InitMsgProcessor();

    // start dataflow will init receiver and exporter
    int StartDataflow();

    // uninit
    void StopZkService();

    void StopConftor();

    void StopScheduler();

    void StopChannelIDAPIService();

    void StopOpsServiceAPISerivce();

    void StopDataflow();

    void StopChannelIdManager();
    void StopStreamExporterManager();
    void GuardLoop();

private:
    DataServer(const DataServer& src);
    DataServer& operator=(const DataServer& src);

    OpsCollection* OpsReportClient();

    std::shared_ptr<LoadBalanceSchedule> Scheduler();
    std::shared_ptr<ChannelIdZkAPI> ChannelIdZkApiInst();
    std::shared_ptr<ChannelIdApiMetrics> ChannelIdApiMetricsInst();
    std::shared_ptr<ChannelIdCommonAPI> ChannelIdCommonApiInst();
    std::shared_ptr<ChannelIdStreamExporterManager> StreamExporterInst();
    std::shared_ptr<ChannelIdManager> ChannelIDMangerInst();

    std::vector<std::shared_ptr<DataIdZkClient>> DataIdZkApiInsts();
    std::shared_ptr<gse::discover::zkapi::ZkApi> DiscoverZkApiClient();
    std::string GetAdvertiseIPInst();
    std::shared_ptr<DataProcessConfig> GetConfigPtrInst();

    bool UpgradeFlag();

    std::shared_ptr<gse::discover::zkapi::ZkApi> CreateZkClient(const std::string& zkHost, const std::string& zkAuth, int timeoutMs);

private:
    std::string m_advertiseIP;
    std::shared_ptr<DataProcessConfig> m_configPtr;

private:
    DataFlow* m_dataFlow;
    Configurator* m_conftor;
    bool m_upgradeFlag;
    bool m_exit;

    OpsCollection* m_opsReport;
    std::shared_ptr<OPSHttpServer> m_opsHttpServer;
    std::shared_ptr<Stack> m_stackService;
    std::shared_ptr<OpHealthZ> m_opsHealthZ;
    std::shared_ptr<LoadBalanceSchedule> m_scheduler;
    std::shared_ptr<SignalControlMsgProcessor> m_controlMsgProcessor;
    std::shared_ptr<DataReportMsgProcessor> m_dataReportMsgProcessor;
    std::shared_ptr<OpsReportMsgProcessor> m_opsMsgProcessor;
    std::shared_ptr<HTTPAPIService> m_channelIdApiSerivce;

    std::shared_ptr<gse::discover::zkapi::ZkApi> m_channelIdZkClient;
    std::shared_ptr<gse::discover::zkapi::ZkApi> m_discoverZkClient;
    std::shared_ptr<ChannelIdZkAPI> m_channelIdApiZkClient;
    std::shared_ptr<ChannelIdApiMetrics> m_channelIdApiMetrics;
    std::shared_ptr<ChannelIdCommonAPI> m_channelIdCommonApi;
    std::vector<std::shared_ptr<DataIdZkClient>> m_dataIdZkClients;

    std::shared_ptr<ChannelIdStreamExporterManager> m_streamExporterManager;
    std::shared_ptr<ChannelIdManager> m_channelIdManager;

    std::shared_ptr<StreamExporterIDLoader> m_streamExporterLoader;
};

} // namespace data
} // namespace gse
#endif
