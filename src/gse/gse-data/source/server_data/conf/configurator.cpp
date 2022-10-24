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

#include "configurator.h"
#include <vector>

#include <netinet/in.h>
#include <sys/socket.h>

#include "bbx/gse_errno.h"
#include "conftor/zkconftor.h"
#include "log/log.h"
#include "tools/filesystem.h"
#include "tools/hostinfo.h"
#include "tools/macros.h"
#include "tools/net.h"
#include "tools/rapidjson_macro.h"
#include "tools/strings.h"

#include "api/channelid_def.h"
#include "api/channelid_struct.h"
#include "api/error_code.h"

#include "balance_config.h"
#include "channel_id_config.h"
#include "dataserver.h"
#include "ops/op_healthz.h"
#include "tools/finally.hpp"
#include "tools/strings.h"

#include "utils.h"

namespace gse {
namespace data {

// set json tag
const MetaType *ServiceNodeMeta::properties[] = {new Property<ServiceNodeMeta, std::string>(&ServiceNodeMeta::m_zoneId, "zone_id"),
                                                 new Property<ServiceNodeMeta, std::string>(&ServiceNodeMeta::m_cityId, "city_id"),
                                                 new Property<ServiceNodeMeta, std::string>(&ServiceNodeMeta::m_clusterName, "cluster_name"),
                                                 new Property<ServiceNodeMeta, std::string>(&ServiceNodeMeta::m_clusterId, "cluster_id"),
                                                 new Property<ServiceNodeMeta, std::string>(&ServiceNodeMeta::m_serviceName, "service_name"),
                                                 new Property<ServiceNodeMeta, std::string>(&ServiceNodeMeta::m_serviceId, "service_id"),
                                                 new Property<ServiceNodeMeta, std::string>(&ServiceNodeMeta::m_advertiseIp, "node_ip")};

// set json tag
const MetaType *ServiceNode::properties[] = {new Property<ServiceNode, std::string>(&ServiceNode::m_serviceIP, "service_ip"),
                                             new Property<ServiceNode, uint32_t>(&ServiceNode::m_port, "port"),
                                             new Property<ServiceNode, std::string>(&ServiceNode::m_protocol, "protocal"),
                                             new Property<ServiceNode, bool>(&ServiceNode::m_ssl, "ssl")};

static uint32_t kDefaultTglogChannelId = 0;

Configurator::Configurator()
{
    m_conftor = NULL;
    m_channelIdZKConftor = NULL;
    m_dataConf = NULL;
    m_ptrEventThread = NULL;

    m_cpuUsage = 0;
    m_hostPerformance = 0;
    m_hostLoadBance = 0;
    m_loadweight = 0;
    m_dloadweight = 0;

    m_platIdManager = NULL;

    m_discoverZkAcl = false;
    m_channelIdZkAcl = false;

    m_ptrConfigFunc[CONFITEMFLAG_DATACONF] = Configurator::updateDataConf;
    m_ptrConfigFunc[CONFITEMFLAG_BALANCE_CONFIG] = Configurator::updateBalanceConfig;

    m_ptrConfigFunc[CONFITEMFLAG_CHANNELID_LIST] = Configurator::handleChannelIdChildListChangeEvt;
    m_ptrConfigFunc[CONFITEMFLAG_CHANNELID_CONFIG] = Configurator::updateChannelIDConfig;

    m_ptrConfigFunc[CONFITEMFLAG_TGLOG_CHANNEL_ID_VALUE] = Configurator::handleTglogChannelIdChangeEvtFromZK;
    //--------------plat id---------------------------
    m_ptrConfigFunc[CONFITEMFLAG_PLAT_ID_CONFIG_LIST] = Configurator::handlePlatIdChildListChangeEvt;
    m_ptrConfigFunc[CONFITEMFLAG_PLAT_ID_CONFIG_VALUE] = Configurator::updatePlatIdConfig;
    // OPS config
    m_ptrConfigFunc[CONFITEMFLAG_OPS_SERVICE_CONFIG_LIST] = Configurator::HandleOpsConfigChildListChangeEvt;
    m_ptrConfigFunc[CONFITEMFLAG_OPS_SERVICE_CONFIG_VALUE] = Configurator::UpdateOpsConfigValue;
}

Configurator::~Configurator()
{
    if (m_conftor != NULL)
    {
        m_conftor->Stop();
        delete m_conftor;
        m_conftor = NULL;
    }

    if (m_dataConf != NULL)
    {
        delete m_dataConf;
        m_dataConf = NULL;
    }

    if (m_channelIdZKConftor != NULL)
    {
        m_channelIdZKConftor->Stop();
        delete m_channelIdZKConftor;
        m_channelIdZKConftor = NULL;
    }

    if (m_platIdManager != NULL)
    {
        m_platIdManager->Stop();
        delete m_platIdManager;
        m_platIdManager = NULL;
    }

    for (std::size_t idx = 0; idx < m_callbacks.size(); ++idx)
    {
        delete m_callbacks.at(idx);
    }
    m_callbacks.clear();
}

ChannelIdManager *Configurator::GetPlatIdManager()
{
    return m_platIdManager;
}

void Configurator::reportBalance(evutil_socket_t fd, short what, void *args)
{
    Configurator *ptr_this = reinterpret_cast<Configurator *>(args);
    ptr_this->updateDataServerLoadBalanceNode();
}

void Configurator::ReportSystemLoad(evutil_socket_t fd, short what, void *args)
{
    Configurator *ptrThis = reinterpret_cast<Configurator *>(args);
    ptrThis->ReportLoadInfo();
}

// this v2 service load info
void Configurator::ReportLoadInfo()
{
    std::string serviceNodeJson = std::move(BuildServiceNodeJson());
    std::string nodePath = std::string(ZK_SERIVCE_NODE_PATH) + "/" + DataServer::GetNodeId();

    int ret = m_conftor->SetConfItem(nodePath, serviceNodeJson);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("failed to set service node load info, path:%s", ZK_SERIVCE_NODE_PATH);
        CreateV2ServiceNode();
        return;
    }

    LOG_DEBUG("successfully set service node(%s) config info(%s)", nodePath.c_str(), serviceNodeJson.c_str());
    return;
}

uint32_t Configurator::getDefaultTglogChannelId()
{
    return kDefaultTglogChannelId;
}

int Configurator::Init(std::shared_ptr<DataProcessConfig> configPtr)
{
    LOG_DEBUG("init from remote zk config center");

    if (configPtr->m_clusterInfoConfig.m_advertiseIp != "")
    {
        m_selfIp = configPtr->m_clusterInfoConfig.m_advertiseIp;
    }
    else
    {
        m_selfIp = gse::tools::net::GetMachineIp();
    }

    m_clusterName = configPtr->m_clusterInfoConfig.m_clusterName;
    m_instanceId = configPtr->m_clusterInfoConfig.m_instanceId;
    m_zoneID = configPtr->m_clusterInfoConfig.m_zoneId;
    m_cityID = configPtr->m_clusterInfoConfig.m_cityId;

    m_watchPath = configPtr->m_tglogConfig.m_watchPath;
    m_localDataFlowConfig = configPtr->m_configFilePath;

    m_dataConf = new DataConf();

    int ret = initFromLocalConfig(m_localDataFlowConfig);
    if (configPtr->m_zooKeeperConfig.m_serviceDiscoverZkHost.empty())
    {
        LOG_DEBUG("Configurator init finished without zk");
        return ret;
    }

    m_systemResourceMonitor = std::make_shared<SystemResourceMonitor>();
    m_systemResourceMonitor->SetEthName(configPtr->m_balanceConfig.m_netDevName);
    m_systemResourceMonitor->SetNetDevMaxSpeed(configPtr->m_balanceConfig.m_netDevMaxSpeed);

    m_systemResourceMonitor->Start();

    m_systemConnectionMonitor = std::make_shared<SystemConnectionMonitor>();
    m_systemConnectionMonitor->SetMaxConnectionCount(configPtr->m_balanceConfig.m_maxAgentCount);

    ret = StartConftor();
    if (GSE_SUCCESS != ret)
    {
        return ret;
    }

    m_platIdManager = new ChannelIdManager(true);
    m_platIdManager->Start();

    // watch configure item
    watchConfigsFromZK();

    //启动成功后注册服务发现节点
    ret = createBaseConfItem();
    if (GSE_SUCCESS != ret)
    {
        StopConftor();
        return ret;
    }

    ret = CreateV2ServiceNode();
    if (GSE_SUCCESS != ret)
    {
        return ret;
    }

    return GSE_SUCCESS;
}

int Configurator::initFromLocalConfig(const std::string &dataFlowConf)
{
    std::string configContent;
    int ret_value = gse::tools::filesystem::ReadFromFile(dataFlowConf, configContent);
    if (GSE_SUCCESS != ret_value)
    {
        LOG_ERROR("it is failed to read the dataflow config by the absolute path ( %s ) ", SAFE_CSTR(dataFlowConf.c_str()));
        return ret_value;
    }

    return updateDataFlowConf(configContent);
}

// node value: {"cityid":"", "regid":""}
void Configurator::updateLocationFromZK()
{
    std::string path = "/gse/config/ip2city/" + m_selfIp;
    std::string pos;

    int ret = m_conftor->GetConfItem(path, pos, NULL, NULL, CONFITEMFLAG_UNSET);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("fail to get cityId and regId from zk node[%s]", SAFE_CSTR(path.c_str()));
        return;
    }

    Json::Reader reader(Json::Features::strictMode());
    Json::Value value;
    if (!reader.parse(pos, value, false))
    {
        LOG_ERROR("fail to parse json[%s]", SAFE_CSTR(pos.c_str()));
        return;
    }

    m_cityID = value.get("cityid", "default").asString();
    m_zoneID = value.get("regid", "default").asString();

    return;
}

DataFlowConf *Configurator::GetDataFlowConf()
{
    DataFlowConf *pDataFlowConf = NULL;

    pDataFlowConf = m_dataConf->GetDataFlowConf();

    return pDataFlowConf;
}

void Configurator::SetChannelIdZkClient(std::shared_ptr<gse::discover::zkapi::ZkApi> zkClient, bool zkAcl)
{
    m_channelIdZkClient = zkClient;
    m_channelIdZkAcl = zkAcl;
}

void Configurator::SetDiscoverZkClient(std::shared_ptr<gse::discover::zkapi::ZkApi> zkClient, bool zkAcl)
{
    m_discoverZkClient = zkClient;
    m_discoverZkAcl = zkAcl;
}

int Configurator::StartConftor()
{
    StopConftor();

    if (NULL == m_conftor)
    {
        m_conftor = new ZkConftor(m_discoverZkClient.get(), m_discoverZkAcl);
    }

    int ret = m_conftor->Start();
    if (GSE_SUCCESS != ret)
    {
        LOG_ERROR("failed to start discover zk configtor");

        OpHealthZ::AddInitHealthInfo("zookeeper", "failed to start zookeeper", -1);
        return ret;
    }

    if (m_discoverZkClient != m_channelIdZkClient)
    {
        m_channelIdZKConftor = new ZkConftor(m_channelIdZkClient.get(), m_channelIdZkAcl);
        ret = m_channelIdZKConftor->Start();
        if (GSE_SUCCESS != ret)
        {
            LOG_ERROR("failed to connect the channelid zk host");
            OpHealthZ::AddInitHealthInfo("channelid_zookeeper", "failed to start channelid zookeeper", -1);

            return ret;
        }
    }
    else
    {
        m_channelIdZKConftor = m_conftor;
    }

    m_ptrEventThread = new EventThread();
    int registerResponseVal = m_ptrEventThread->RegisterTimerPersistEvent(Configurator::reportBalance, this, 60); // 1 minutes
    if (GSE_SUCCESS > registerResponseVal)
    {
        LOG_ERROR("failed to register timer, error code %d", registerResponseVal);
        return registerResponseVal;
    }

    registerResponseVal = m_ptrEventThread->RegisterTimerPersistEvent(Configurator::ReportSystemLoad, this, 60);
    if (GSE_SUCCESS > registerResponseVal)
    {
        LOG_ERROR("failed to register timer, error code %d", registerResponseVal);
        return registerResponseVal;
    }

    int startEventThreadVal = m_ptrEventThread->Start();
    if (GSE_SUCCESS > startEventThreadVal)
    {
        LOG_ERROR("failed to start event thread, error code %d", startEventThreadVal);
        return startEventThreadVal;
    }

    return GSE_SUCCESS;
}

void Configurator::StopConftor()
{
    if (m_conftor != NULL)
    {
        m_conftor->Stop();
        delete m_conftor;
        m_conftor = NULL;
    }

    if (NULL != m_ptrEventThread)
    {
        m_ptrEventThread->Stop();
        m_ptrEventThread->Join();
        delete m_ptrEventThread;
        m_ptrEventThread = NULL;
    }

    m_systemResourceMonitor->Stop();
    m_systemResourceMonitor->Join();

    if (m_platIdManager != NULL)
    {
        m_platIdManager->Stop();
        delete m_platIdManager;
        m_platIdManager = NULL;
    }
}

int Configurator::createBaseConfItem()
{
    if (NULL == m_conftor)
    {
        LOG_WARN("the pointer of conftor is NULL, please initialize conftor first");
        return GSE_ERROR;
    }

    // TODO:兼容1.0 版本agent
    DeleteEphemeralZkNodes();

    string value = "{\"fd\":1,\"rb\":0}";
    string strrst;

    int ret = GSE_SUCCESS;
    int zkRet = GSE_SUCCESS;

    std::string nodePath(DATA_BASE_PATH);
    std::string empty_value;

    ret = m_conftor->SetConfItem(nodePath, empty_value);
    if (ret != GSE_SUCCESS)
    {
        ret = m_conftor->CreateConfItemWithParents(nodePath, empty_value);
        if (ret != GSE_SUCCESS)
        {
            return ret;
        }
    }

    string allPath = nodePath + "/all";
    ret = m_conftor->SetConfItem(allPath, value);
    if (ret != GSE_SUCCESS)
    {
        ret = m_conftor->CreateConfItemWithParents(allPath, value);
        if (ret != GSE_SUCCESS)
        {
            return ret;
        }
    }

    string ipPath = allPath + "/" + m_selfIp;
    ret = m_conftor->SetConfItem(ipPath, m_selfIp);
    if (ret != GSE_SUCCESS)
    {
        ret = m_conftor->CreateConfItemWithParents(ipPath, value, true);
        if (ret != GSE_SUCCESS)
        {
            return ret;
        }
    }

    std::string service_node(nodePath);

    if (!m_zoneID.empty())
    {
        service_node += "/" + m_zoneID;
    }

    if (!m_cityID.empty())
    {
        service_node += "/" + m_cityID;
    }

    if (!m_selfIp.empty())
    {
        service_node += "/" + m_selfIp;
    }

    zkRet = m_conftor->CreateConfItemWithParents(service_node, value, true);
    if (zkRet != GSE_SUCCESS)
    {
        LOG_ERROR("failed to create zk node[%s], ret[%d]", service_node.c_str(), zkRet);
        return GSE_ERROR;
    }

    if (!m_clusterName.empty())
    {
        std::string clusterNodePath = nodePath + "/" + m_clusterName + "/" + m_selfIp;
        if (m_conftor->SetConfItem(clusterNodePath, value) != GSE_SUCCESS)
        {
            if (m_conftor->CreateConfItemWithParents(clusterNodePath, value, true) != GSE_SUCCESS)
            {
                LOG_ERROR("failed to create zk node[%s] value %s", clusterNodePath.c_str(), value.c_str());
                return GSE_ERROR;
            }
        }
    }

    updateDataServerLoadBalanceNode();

    return GSE_SUCCESS;

    return ret;
}

int Configurator::watchDataFlow()
{

    if (!m_localDataFlowConfig.empty())
    {
        return GSE_SUCCESS;
    }

    // dataflow
    std::string rootpath(DATA_BASE_PATH);
    std::string path = rootpath + "/all/" + m_selfIp;
    if (m_clusterName != "")
    {
        path = rootpath + "/cluster/" + m_clusterName + "/" + m_selfIp;
    }

    if (m_instanceId != "")
    {
        path = path + "/" + m_instanceId;
    }
    path = path + "/dataflow";

    std::string value;
    int ret = m_conftor->GetConfItem(path, value, Configurator::watchConfCallBack, this, CONFITEMFLAG_DATACONF);
    if (ret == GSE_SUCCESS)
    {
        updateDataFlowConf(value);
    }
    return ret;
}

int Configurator::updateChannelIDConfigFromZK(uint32_t channelID)
{

    std::string value;
    /*
    int ret = m_conftor->GetConfItem(nodePath, value, Configurator::watchConfCallBack, this, CONFITEMFLAG_CHANNELID_CONFIG);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("it is failed to get the node (%s) value from the zk, maybe it's a invalid json", SAFE_CSTR(nodePath.c_str()));
        return ret;
    }
    */
    return updateChannelID(channelID, value);
}

void Configurator::handleTglogChannelIdChangeEvtFromZK(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    LOG_INFO("tglog channelid value change event, zk path:%s", confItem.m_Key.c_str());
    if (confItem.m_valueType == CONFITEMVALUE_TYPE_DELETE)
    {
        if (kDefaultTglogChannelId != 0)
        {
            pConftor->DeleteChannelID(kDefaultTglogChannelId);
        }

        LOG_DEBUG("delete tglog channel id config(%d)", kDefaultTglogChannelId);
    }

    if (GSE_SUCCESS != pConftor->WatchTglogChannelId())
    {
        LOG_ERROR("failed to update tglog channel id from zk, %s", confItem.m_Key.c_str());
        return;
    }
    return;
}

int Configurator::WatchTglogChannelId()
{
    std::string tglog_channel_id_notify_path = ZK_CHANNEL_ID_TGLOG_NOTIFY_PATH(m_watchPath);
    std::string tglog_channel_id_json;

    int ret = m_channelIdZKConftor->GetConfItem(tglog_channel_id_notify_path, tglog_channel_id_json, Configurator::watchConfCallBack, this, CONFITEMFLAG_TGLOG_CHANNEL_ID_VALUE);
    if (ret != GSE_SUCCESS)
    {
        LOG_WARN("Tglog Channelid is not configured, zk path:%s", tglog_channel_id_notify_path.c_str());
        ret = m_channelIdZKConftor->ExistConfItem(tglog_channel_id_notify_path, Configurator::watchConfCallBack, this, CONFITEMFLAG_TGLOG_CHANNEL_ID_VALUE);
        return ret;
    }

    /*
    {
        "channelid":1
    }
    */

    Json::Value property_cfg_json;
    Json::Reader channel_cfg_json_reader;
    if (!channel_cfg_json_reader.parse(tglog_channel_id_json, property_cfg_json))
    {
        LOG_ERROR("the channel id %s 's config is invalid", tglog_channel_id_notify_path.c_str());
        return false;
    }
    uint32_t channel_id = 0;

    if (property_cfg_json.isMember("channelid"))
    {
        channel_id = property_cfg_json.get("channelid", -1).asInt();
    }
    else
    {
        LOG_ERROR("Tglog Channelid's config is not valid json, zk path:%s", tglog_channel_id_notify_path.c_str());
        return GSE_ERROR;
    }

    ret = updateChannelIDConfigFromZK(channel_id);
    kDefaultTglogChannelId = channel_id;
    return ret;
}

int Configurator::watchChannelID()
{
    // TGLOG watch, TGLOG only support one channel id
    if (m_watchPath != "")
    {
        return WatchTglogChannelId();
    }

    std::vector<std::string> channelIdList;
    std::string channelIdRootPath = ZK_CHANNEL_ID_CONFIG_BASE_PATH;
    LOG_DEBUG("channel id root path:%s", channelIdRootPath.c_str());

    int ret = m_channelIdZKConftor->GetChildConfItem(channelIdRootPath, channelIdList, Configurator::watchConfCallBack, this, CONFITEMFLAG_CHANNELID_LIST);
    if (ret != GSE_SUCCESS)
    {
        LOG_WARN("failed to get channel id root path from zk:%s", channelIdRootPath.c_str());
        ret = m_channelIdZKConftor->ExistConfItem(channelIdRootPath, Configurator::watchConfCallBack, this, CONFITEMFLAG_CHANNELID_LIST);
        return ret;
    }

    LOG_INFO("start to read the channel id config from zk");
    for (int channel_id_idx = 0; channel_id_idx < channelIdList.size(); ++channel_id_idx)
    {
        std::string strChannelID(channelIdList[channel_id_idx]);
        if (strChannelID.compare("index") == 0)
        {
            continue;
        }

        uint32_t channelId = gse::tools::strings::StringToUint32(strChannelID);
        // if exist not update
        if (FindChannelID(channelId))
        {
            LOG_DEBUG("channelid(%d) exist, don't need update", channelId);
            continue;
        }

        ret = updateChannelIDConfigFromZK(channelId);
        if (GSE_SUCCESS != ret)
        {
            LOG_ERROR("failed to update channelid config from zk, path:%s", ZK_CHANNEL_ID_WATCH_PATH(strChannelID).c_str());
            continue;
        }
    }

    LOG_INFO("finish reading the channel id config from zk node");
    channelIdList.clear();
    // register  watch path
    m_channelIdZKConftor->GetChildConfItem(channelIdRootPath, channelIdList, Configurator::watchConfCallBack, this, CONFITEMFLAG_CHANNELID_LIST);
    return GSE_SUCCESS;
}

int Configurator::watchPlatID()
{
    std::vector<std::string> plat_id_list;
    std::string root_path = ZK_PLAT_ID_CONFIG_BASE_PATH;
    LOG_DEBUG("plat id root path:%s", root_path.c_str());

    int ret = m_channelIdZKConftor->GetChildConfItem(root_path, plat_id_list, Configurator::watchConfCallBack, this, CONFITEMFLAG_PLAT_ID_CONFIG_LIST);
    if (ret != GSE_SUCCESS)
    {
        LOG_WARN("failed to get channel id root path from zk:%s", root_path.c_str());
        ret = m_channelIdZKConftor->ExistConfItem(root_path, Configurator::watchConfCallBack, this, CONFITEMFLAG_PLAT_ID_CONFIG_LIST);
        return ret;
    }

    LOG_INFO("Start to read the plat id config from zk node");
    for (int plat_id_idx = 0; plat_id_idx < plat_id_list.size(); ++plat_id_idx)
    {
        uint32_t plat_id = gse::tools::strings::StringToUint32(plat_id_list[plat_id_idx]);
        ret = updatePlatID(plat_id);
        if (GSE_SUCCESS != ret)
        {
            LOG_ERROR("can not update plat config from zk, %s/%s", root_path.c_str(), plat_id_list[plat_id_idx].c_str());
        }
    }

    LOG_INFO("Finished reading the plat id config from zk node");
    plat_id_list.clear();
    m_channelIdZKConftor->GetChildConfItem(root_path, plat_id_list, Configurator::watchConfCallBack, this, CONFITEMFLAG_PLAT_ID_CONFIG_LIST);
    return GSE_SUCCESS;
}

int Configurator::watchOpsServiceConfig()
{
    std::vector<std::string> opsServiceNodeConfigs;
    std::string opsConfigRootPath = DATA_OPS_SERVICE_CONF_PATH;
    LOG_DEBUG("ops config root path:%s", opsConfigRootPath.c_str());

    int ret = m_conftor->GetChildConfItem(opsConfigRootPath, opsServiceNodeConfigs, NULL, this, CONFITEMFLAG_OPS_SERVICE_CONFIG_LIST);
    if (ret != GSE_SUCCESS)
    {
        LOG_WARN("failed to get ops service config from zk:%s", opsConfigRootPath.c_str());
        ret = m_conftor->ExistConfItem(opsConfigRootPath, Configurator::watchConfCallBack, this, CONFITEMFLAG_OPS_SERVICE_CONFIG_LIST);
        return ret;
    }
    std::string opsServiceNodePath;
    for (std::vector<std::string>::iterator it = opsServiceNodeConfigs.begin(); it != opsServiceNodeConfigs.end(); ++it)
    {
        opsServiceNodePath.clear();
        opsServiceNodePath.append(opsConfigRootPath).append("/").append(*it);

        std::string opsConfigJsonValue;
        ret = GetOpsConfig(opsServiceNodePath, opsConfigJsonValue);
        if (ret != GSE_SUCCESS)
        {
            continue;
        }

        OpsServiceConfig opsServiceConfig;
        ret = opsServiceConfig.ParseJsonConfig(opsConfigJsonValue);
        if (ret != GSE_SUCCESS)
        {
            return ret;
        }

        m_opsConfig.Push(opsServiceConfig.m_serivceId, opsServiceConfig);
    }

    m_conftor->GetChildConfItem(opsConfigRootPath, opsServiceNodeConfigs, Configurator::watchConfCallBack, this, CONFITEMFLAG_OPS_SERVICE_CONFIG_LIST);
    return GSE_SUCCESS;
}

int Configurator::watchBalanceConfig()
{
    // 获取负载配置参数
    std::string balancfgPath;
    std::string cfgNode;
    do
    {
        balancfgPath = "/gse/config/etc/dataserver/" + m_selfIp + "/balancecfg";
        int ret = m_conftor->GetConfItem(balancfgPath, cfgNode, Configurator::watchConfCallBack, this, CONFITEMFLAG_BALANCE_CONFIG);
        if (ret == GSE_SUCCESS)
        {
            break;
        }

        LOG_INFO("it is failed to get balance config from the zk node (%s)", balancfgPath.c_str());
        balancfgPath = "/gse/config/etc/dataserver/all/balancecfg";
        ret = m_conftor->GetConfItem(balancfgPath, cfgNode, Configurator::watchConfCallBack, this, CONFITEMFLAG_BALANCE_CONFIG);
        if (ret == GSE_SUCCESS)
        {
            break;
        }

        LOG_ERROR("it is failed to get balance config from the node (%s), please the node in zk", balancfgPath.c_str());
    } while (false);

    return updateDataServerBalanceConfig(balancfgPath, cfgNode);
}

int Configurator::updateDataServerBalanceConfig(const std::string &path, const std::string &context)
{
    // 解析负载配置参数
    LOG_DEBUG("balance config: %s", (context.c_str()));
    Json::Reader reader(Json::Features::strictMode());
    if (!reader.parse(context, m_cfgLoadBalanceValue, false))
    {
        LOG_ERROR("parse to json failed: %s", (context.c_str()));
        return GSE_JSON_INVALID;
    }

    m_ethName = m_cfgLoadBalanceValue.get("netdev", "eth1").asString();
    if (m_systemResourceMonitor)
    {
        m_systemResourceMonitor->SetEthName(m_ethName);
    }

    std::string value;
    int ret = m_conftor->GetConfItem(path, value, Configurator::watchConfCallBack, this, CONFITEMFLAG_BALANCE_CONFIG);
    if (ret != GSE_SUCCESS)
    {
        LOG_WARN("failed to get node:%s value, ret:%d", path.c_str(), ret);
    }
    return GSE_SUCCESS;
}

void Configurator::DeleteEphemeralZkNodes()
{
    std::string dsNode = "/gse/config/server/dataserver/" + m_zoneID + "/" + m_cityID + "/" + m_selfIp;
    m_conftor->DeleteConfItem(dsNode);

    if (!m_clusterName.empty())
    {
        string dsNodeNew = "/gse/config/server/dataserver/cluster/" + m_clusterName + "/" + m_selfIp;
        m_conftor->DeleteConfItem(dsNodeNew);
    }
}

int Configurator::updateDataServerLoadBalanceNode()
{
    std::string ethname = m_cfgLoadBalanceValue.get("netdev", "eth1").asString();

    static std::vector<gse::os::hostinfo::stats_net_dev> lastNetDevItems;
    static int lastTimeStamp = 0;
    std::vector<gse::os::hostinfo::stats_net_dev> currNetDevItems;

    gse::tools::hostinfo::ReadNetDev(currNetDevItems, ethname.c_str());

    uint64_t usedCpu;
    uint64_t totalCpu;

    static uint64_t lastUsedCpu;
    static uint64_t lastTotalCpu;

    uint32_t pid = gse::tools::process::GetProcessID();
    usedCpu = gse::tools::process::GetProcessCpuTime(pid);
    totalCpu = gse::tools::process::GetTotalCpuTime();

    if (currNetDevItems.size() <= 0)
    {
        LOG_WARN("can't get net device info, ethname:%s", ethname.c_str());
        return GSE_ERROR;
    }

    LOG_INFO("netdevice name:%s rx_bytes:%llu", currNetDevItems.at(0).devname, currNetDevItems.at(0).rx_bytes);
    if (lastNetDevItems.size() <= 0)
    {
        lastNetDevItems = currNetDevItems;
        lastTimeStamp = ::time(NULL);
        lastUsedCpu = usedCpu;
        lastTotalCpu = totalCpu;
        return GSE_SUCCESS;
    }

    float netReadBytes = currNetDevItems.at(0).rx_bytes - lastNetDevItems.at(0).rx_bytes;
    int nowtime = ::time(NULL);
    int dtimestamp = nowtime - lastTimeStamp;
    if (dtimestamp <= 0)
    {
        LOG_WARN("diff time invalid, now:%d, last time:%d", nowtime, lastTimeStamp);
        return GSE_ERROR;
    }

    lastTimeStamp = nowtime;
    lastNetDevItems = currNetDevItems;

    uint64_t dtotal = totalCpu - lastTotalCpu;
    if (0 != dtotal)
    {
        m_cpuUsage = (usedCpu - lastUsedCpu) * 1.0 / dtotal;
    }
    else
    {
        return GSE_SUCCESS;
    }

    lastUsedCpu = usedCpu;
    lastTotalCpu = totalCpu;

    // C(Si) = K1 *N* C(Ci) + K2 * C(Di) + K4 * C(Ni), i=1,2,3,..,n

    int cpunum = gse::tools::hostinfo::GetCpuNum();
    double cpufreq = gse::tools::hostinfo::GetCpuFreq();
    int netwide = GetNetDevSpeed(ethname.c_str());
    int total = 0, free = 0, vtotal = 0, vfree = 0;
    int buffer = 0, cache = 0;
    gse::tools::hostinfo::GetMemoryInfo(&total, &free, &vtotal, &vfree, &buffer, &cache);

    LOG_INFO("cpunum:%d cpufreq:%f netwide:%d  mem_total:%d mem_free:%d mem_vtotal:%d mem_vfree:%d",
             cpunum, cpufreq, netwide, total, free, vtotal, vfree);
    //  netp  cpup memp diskp  分别用于代表 一分 所对应的 带宽，cpu 频率，disk io 速率。

    int allCpuFreq = cpunum * cpufreq;
    float markcpu = m_cfgLoadBalanceValue.get("cpuk", 0).asFloat() * allCpuFreq / m_cfgLoadBalanceValue.get("cpup", 1.0).asFloat();
    float markmem = m_cfgLoadBalanceValue.get("memk", 0).asFloat() * total / m_cfgLoadBalanceValue.get("memp", 1.0).asFloat();
    float marknet = m_cfgLoadBalanceValue.get("netk", 0).asFloat() * netwide / m_cfgLoadBalanceValue.get("netp", 1.0).asFloat();

    m_hostPerformance = markcpu + markmem + marknet;

    // 计算当前网络流量，单位 Mb
    float diskUsage = 0;
    int memUsage = total - free;
    float netUsagePer = (((netReadBytes * 8) / 1024) / 1024) / dtimestamp; // Mb

    // 为各个，cpu 按照使用率计算占用性能指数，并以 cpup 为单位进行打分
    float cpuUsageMark = (m_cfgLoadBalanceValue.get("cpur", 0).asFloat() * m_cpuUsage * cpunum * cpufreq) / m_cfgLoadBalanceValue.get("cpup", 1.0).asFloat();
    float memUsageMark = m_cfgLoadBalanceValue.get("memr", 0).asFloat() * memUsage / m_cfgLoadBalanceValue.get("memp", 1.0).asFloat();
    float netUsageMark = m_cfgLoadBalanceValue.get("netr", 0).asFloat() * netUsagePer / m_cfgLoadBalanceValue.get("netp", 1.0).asFloat();

    m_hostLoadBance = cpuUsageMark + memUsageMark + netUsageMark;

    LOG_INFO("origin usage cpu:%f mem:%d disk:%f net:%f", m_cpuUsage, memUsage, netUsagePer);

    Json::Value data;
    data["hostperformance"] = m_hostPerformance;
    data["loadbalance"] = m_hostLoadBance;
    if (0 != m_hostPerformance)
    {
        m_loadweight = m_hostLoadBance / m_hostPerformance;
    }
    data["loadweight"] = m_loadweight;
    m_dloadweight = (m_cfgLoadBalanceValue.get("weightmax", 0.8).asFloat() - m_loadweight) * m_hostPerformance;
    data["dloadweight"] = m_dloadweight;

    LOG_INFO("performance:%f mark_cpu:%f mark_mem:%f  mark_net:%f loadbalance:%f mark_usage_cpu:%f mark_usage_mem:%f mark_usage_net:%f load_weight:%f dload_weight:%f",
             m_hostPerformance, markcpu, markmem, marknet, m_hostLoadBance, cpuUsageMark, memUsageMark, netUsageMark, m_loadweight, m_dloadweight);

    Json::FastWriter fwriter;
    string value = fwriter.write(data);

    // set origin node
    std::string dsNode = "/gse/config/server/dataserver/" + m_zoneID + "/" + m_cityID + "/" + m_selfIp;
    if (m_conftor->SetConfItem(dsNode, value) != GSE_SUCCESS)
    {
        if (m_conftor->CreateConfItemWithParents(dsNode, value, true) != GSE_SUCCESS)
        {
            LOG_ERROR("set node %s value %s failed", SAFE_CSTR(dsNode.c_str()), SAFE_CSTR(value.c_str()));
        }
    }

    // set new node
    if (!m_clusterName.empty())
    {
        string dsNodeNew = "/gse/config/server/dataserver/cluster/" + m_clusterName + "/" + m_selfIp;
        if (m_conftor->SetConfItem(dsNodeNew, value) != GSE_SUCCESS)
        {
            if (m_conftor->CreateConfItemWithParents(dsNodeNew, value, true) != GSE_SUCCESS)
            {
                LOG_ERROR("set node %s value %s failed", SAFE_CSTR(dsNodeNew.c_str()), SAFE_CSTR(value.c_str()));
            }
        }
    }

    return GSE_SUCCESS;
}

int Configurator::watchConfigsFromZK()
{

    if (m_zoneID.empty() || m_cityID.empty())
    {
        updateLocationFromZK();
    }

    int ret = GSE_SUCCESS;
    ret = watchDataFlow();
    if (ret != GSE_SUCCESS)
    {
        return ret;
    }
    ret = watchChannelID();
    if (ret != GSE_SUCCESS)
    {
        return ret;
    }

    ret = watchPlatID();
    if (ret != GSE_SUCCESS)
    {
        return ret;
    }

    ret = CreateOpsServiceConfigNode();
    if (GSE_SUCCESS != ret)
    {
        return ret;
    }

    ret = watchOpsServiceConfig();
    if (ret != GSE_SUCCESS)
    {
        return ret;
    }

    ret = watchBalanceConfig();
    if (ret != GSE_SUCCESS)
    {
        return ret;
    }

    return ret;
}

int Configurator::updateDataFlowConf(const std::string &context)
{
    int ret = GSE_SUCCESS;

    if (NULL == m_dataConf)
    {
        LOG_WARN("can't parse dataflow configure, because the object dataconf is null");
        return GSE_ERROR;
    }

    ret = m_dataConf->DecodeConf(context);
    if (GSE_SUCCESS != ret)
    {
        return ret;
    }

    return GSE_SUCCESS;
}

void Configurator::ParseMetaConfig(std::string &path, int rc, std::string &propertyValue, ChannelIDConfig *ptr_channelid_config)
{
    ApiError error;
    ptr_channelid_config->AddResponseCount();

    Json::Value jsonValue;
    Json::Reader reader;
    if (!reader.parse(propertyValue, jsonValue))
    {
        LOG_ERROR("meta is invalid json, node path:%s", path.c_str(), propertyValue.c_str(), path.c_str());
        ptr_channelid_config->SetError();
        DoUpdateChannelIdConfigRequest(ptr_channelid_config);
        return;
    }

    if (!ptr_channelid_config->m_ptrMetadata->ParseMetadata(jsonValue, error))
    {
        LOG_ERROR("failed to parse channel id (%s)'s meta config(%s), error:%s", path.c_str(), propertyValue.c_str(), GET_ERROR_MESSAGE(error).c_str());
        ptr_channelid_config->SetError();
        DoUpdateChannelIdConfigRequest(ptr_channelid_config);
        return;
    }

    DoUpdateChannelIdConfigRequest(ptr_channelid_config);
    LOG_DEBUG("parse meta config success, path(%s) value(%s)", path.c_str(), propertyValue.c_str());
}
// typedef void (*FnZkGetValueCallBack)(std::string &path, int rc, const char *value, int32_t value_len, const struct Stat *stat, const void *data);
void Configurator::GetMetaValueCallBack(std::string &path, int rc, const char *value, int32_t valueLen, const struct Stat *stat, const void *ptrCallbackobj)
{
    if (ptrCallbackobj == NULL)
    {
        LOG_ERROR("get mata value input param invalid , callback obj is null", ptrCallbackobj);
        return;
    }

    std::string strValue;
    if (value == NULL || valueLen <= 0)
    {
        LOG_WARN("get mata value invalid, valueLen:%d", valueLen);
    }
    else
    {
        strValue.assign(value, valueLen);
    }

    ChannelIdZkCallBackObj *ptrChannelIdZkCallbakcObj = (ChannelIdZkCallBackObj *)ptrCallbackobj;
    Configurator *callobj = ptrChannelIdZkCallbakcObj->m_ptrConfigurator;
    callobj->ParseMetaConfig(path, rc, strValue, ptrChannelIdZkCallbakcObj->m_ptrChannelidConfig);
    delete ptrChannelIdZkCallbakcObj;
}

bool Configurator::readMetadata(const std::string &metadataPath, ChannelIDConfig *ptrChannelidConfig, std::string &errorMsg)
{
    ChannelIdZkCallBackObj *ptrCallbackobj = new ChannelIdZkCallBackObj();
    ptrCallbackobj->m_ptrChannelidConfig = ptrChannelidConfig;
    ptrCallbackobj->m_ptrConfigurator = this;
    int ret = m_channelIdZKConftor->GetConfItemAsync(metadataPath, NULL, NULL, CONFITEMFLAG_UNSET, Configurator::GetMetaValueCallBack, ptrCallbackobj);
    if (GSE_SUCCESS != ret)
    {
        LOG_ERROR("it is failed to read the config for the zk node path (%s)", metadataPath.c_str());
        ptrChannelidConfig->SetError();
        ptrChannelidConfig->AddResponseCount();
        return false;
    }
    return true;
}

void Configurator::UpdateFilterValue(std::string &path, int rc, std::string &propertyValue, void *ptr_callbackobj)
{
    std::string errorMsg;
    ApiError error;
    ChannelIDConfig *ptr_channelid_config = (ChannelIDConfig *)ptr_callbackobj;

    ptr_channelid_config->AddResponseCount();
    StreamFilter *ptrStreamFilter = new StreamFilter();

    Json::Value propertyCfgJson;
    Json::Reader channel_cfg_json_reader;
    if (!channel_cfg_json_reader.parse(propertyValue, propertyCfgJson))
    {
        LOG_ERROR("filter is invalid json, path:%s, value:%s", path.c_str(), propertyValue.c_str());
        ptr_channelid_config->SetError();
        DoUpdateChannelIdConfigRequest(ptr_channelid_config);
        return;
    }

    if (!ptrStreamFilter->Parse(propertyCfgJson, error))
    {
        LOG_ERROR("failed to parse filter[%s], path:%s, error:%s", propertyValue.c_str(), path.c_str(), GET_ERROR_MESSAGE(error).c_str());
        ptr_channelid_config->SetError();
        DoUpdateChannelIdConfigRequest(ptr_channelid_config);
        return;
    }

    ptr_channelid_config->m_streamFilter.push_back(ptrStreamFilter);

    DoUpdateChannelIdConfigRequest(ptr_channelid_config);
    return;
}

void Configurator::FilterValueResultCallBack(std::string &path, int rc, const char *value, int32_t value_len, const struct Stat *stat, const void *ptr_callbackobj)
{
    if (ptr_callbackobj == NULL)
    {
        return;
    }

    if (value == NULL || value_len <= 0)
    {
        LOG_DEBUG("get filter value invalid");
        return;
    }

    std::string str_value(value, value_len);
    ChannelIdZkCallBackObj *ptr_channelId_zk_callbakc_obj = (ChannelIdZkCallBackObj *)ptr_callbackobj;
    Configurator *callobj = ptr_channelId_zk_callbakc_obj->m_ptrConfigurator;
    callobj->UpdateFilterValue(path, rc, str_value, ptr_channelId_zk_callbakc_obj->m_ptrChannelidConfig);
    delete ptr_channelId_zk_callbakc_obj;
}

void Configurator::GetFiltersValue(std::string &path, int rc, std::vector<std::string> &values, ChannelIDConfig *ptr_channelid_config)
{
    if (ptr_channelid_config == NULL)
    {
        return;
    }

    // TOP 1 level response
    ptr_channelid_config->AddResponseCount();

    if (rc != ZOK)
    {
        LOG_INFO("get filters failed, path:%s, zk return code:%d", path.c_str(), rc);
        DoUpdateChannelIdConfigRequest(ptr_channelid_config);
        return;
    }

    if (values.size() == 0)
    {
        LOG_INFO("no filter rules for channelid's filter(%s) configured", path.c_str());
        DoUpdateChannelIdConfigRequest(ptr_channelid_config);
        return;
    }

    LOG_DEBUG("filters[%s] has child node:%d", path.c_str(), values.size());
    std::string channel_id_str = ptr_channelid_config->m_strChannelId;
    std::size_t max_cnt = values.size();
    ptr_channelid_config->AddNeedAckCount(max_cnt);

    for (std::size_t idx = 0; idx < max_cnt; ++idx)
    {
        std::string property_value;
        std::string zk_node_path(path + "/" + values.at(idx));

        ChannelIdZkCallBackObj *ptr_callbackobj = new ChannelIdZkCallBackObj();
        ptr_callbackobj->m_ptrChannelidConfig = ptr_channelid_config;
        ptr_callbackobj->m_ptrConfigurator = this;
        int ret = m_channelIdZKConftor->GetConfItemAsync(zk_node_path, NULL, NULL, CONFITEMFLAG_UNSET, Configurator::FilterValueResultCallBack, ptr_callbackobj);
        if (GSE_SUCCESS != ret)
        {
            LOG_ERROR("it is failed to read the config for the zk node path (%s)", SAFE_CSTR(zk_node_path.c_str()));
            ptr_channelid_config->SetError();
            ptr_channelid_config->AddResponseCount();
            DoUpdateChannelIdConfigRequest(ptr_channelid_config);
            delete ptr_callbackobj;
            return;
        }
    }

    return;
}

void Configurator::GetFiltersResultCallBack(std::string &path, int rc, std::vector<std::string> &values, const void *ptr_callbackobj)
{
    if (ptr_callbackobj == NULL)
    {
        LOG_ERROR("get filter list, callback obj invalid, get channelid path:%s, result:%d", path.c_str(), rc);
        return;
    }

    ChannelIdZkCallBackObj *ptr_channelId_zk_callbakc_obj = (ChannelIdZkCallBackObj *)ptr_callbackobj;
    Configurator *callobj = ptr_channelId_zk_callbakc_obj->m_ptrConfigurator;
    callobj->GetFiltersValue(path, rc, values, ptr_channelId_zk_callbakc_obj->m_ptrChannelidConfig);
    delete ptr_channelId_zk_callbakc_obj;
    return;
}

bool Configurator::readFilters(const std::string &filterPath, ChannelIDConfig *ptr_channelid_config, std::string &errorMsg)
{
    ChannelIdZkCallBackObj *ptr_callbackobj = new ChannelIdZkCallBackObj();
    ptr_callbackobj->m_ptrChannelidConfig = ptr_channelid_config;
    ptr_callbackobj->m_ptrConfigurator = this;

    LOG_DEBUG("start read filter config[%s]", filterPath.c_str());
    int ret = m_channelIdZKConftor->GetChildConfItemAsync(filterPath, NULL, NULL, CONFITEMFLAG_UNSET, Configurator::GetFiltersResultCallBack, ptr_callbackobj);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("it is failed to read the filter(%s)'s property, zk return(%d)", SAFE_CSTR(filterPath.c_str()), ret);
        ptr_channelid_config->AddResponseCount();
        delete ptr_callbackobj;
        return false;
    }
    return true;
}

// typedef void (*FnZkGetChildCallBack)(std::string &path, int rc, std::string<std::string> &values, void* data);
void Configurator::GetChannelConfigValueCallBack(std::string &path, int rc, std::string &propertyValue, void *ptr_callbackobj)
{
    ChannelIDConfig *ptr_channelid_config = (ChannelIDConfig *)ptr_callbackobj;
    std::string channel_id_str = ptr_channelid_config->m_strChannelId;

    ptr_channelid_config->AddResponseCount();
    LOG_DEBUG("get channel[%s] value(%s)", path.c_str(), propertyValue.c_str());
    if (rc != 0)
    {
        LOG_ERROR("get value return error:%d", rc);
        ptr_channelid_config->SetError();
        DoUpdateChannelIdConfigRequest(ptr_channelid_config);
        return;
    }

    Json::Value property_cfg_json;
    Json::Reader channel_cfg_json_reader;
    if (!channel_cfg_json_reader.parse(propertyValue, property_cfg_json))
    {
        LOG_ERROR("channel is invalid json", path.c_str(), propertyValue.c_str());
        ptr_channelid_config->SetError();
        DoUpdateChannelIdConfigRequest(ptr_channelid_config);
        return;
    }

    Channel *ptrChannel = new Channel();
    ApiError error;
    if (!ptrChannel->Parse(property_cfg_json, error))
    {
        LOG_ERROR("failed to parse channel json[%s], path:%s, error:%s", propertyValue.c_str(), path.c_str(), GET_ERROR_MESSAGE(error).c_str());
        delete ptrChannel;
        ptr_channelid_config->SetError();
        DoUpdateChannelIdConfigRequest(ptr_channelid_config);
        return;
    }

    ptr_channelid_config->m_channels.push_back(ptrChannel);
    DoUpdateChannelIdConfigRequest(ptr_channelid_config);
    return;
}
void Configurator::ChannelValueResultCallBack(std::string &path, int rc, const char *value, int32_t value_len, const struct Stat *stat, const void *ptr_callbackobj)
{
    if (ptr_callbackobj == NULL)
    {
        LOG_ERROR("get channel value result failed ,callback obj is null");
        return;
    }

    if (value == NULL || value_len <= 0)
    {
        LOG_ERROR("get channel value invalid, value len:%", value_len);
        return;
    }

    std::string str_value(value, value_len);
    ChannelIdZkCallBackObj *ptr_channelId_zk_callbakc_obj = (ChannelIdZkCallBackObj *)ptr_callbackobj;
    Configurator *callobj = ptr_channelId_zk_callbakc_obj->m_ptrConfigurator;
    callobj->GetChannelConfigValueCallBack(path, rc, str_value, ptr_channelId_zk_callbakc_obj->m_ptrChannelidConfig);
    delete ptr_channelId_zk_callbakc_obj;
}

void Configurator::GetChannelsConfigValue(std::string &path, int rc, std::vector<std::string> &values, ChannelIDConfig *ptr_channelid_config)
{
    if (ptr_channelid_config == NULL)
    {
        LOG_ERROR("channelid object invalid, path:%s", path.c_str());
        return;
    }

    ptr_channelid_config->AddResponseCount();
    LOG_DEBUG("channel[%s] has child node:%d", path.c_str(), values.size());
    std::string channel_id_str = ptr_channelid_config->m_strChannelId;
    std::size_t max_cnt = values.size();
    ptr_channelid_config->AddNeedAckCount(max_cnt);

    if (rc != ZOK)
    {
        LOG_WARN("get channel failed, path:%s, zk return code:%d", path.c_str(), rc);
        ptr_channelid_config->SetError();
        DoUpdateChannelIdConfigRequest(ptr_channelid_config);
        return;
    }

    if (values.size() == 0)
    {
        // channels must > 0
        LOG_WARN("No channel rules for channelid(%s) configured", path.c_str());
        ptr_channelid_config->SetError();
        DoUpdateChannelIdConfigRequest(ptr_channelid_config);
        return;
    }

    for (std::size_t idx = 0; idx < max_cnt; ++idx)
    {
        std::string property_value;
        std::string zk_node_path(path + "/" + values.at(idx));

        ChannelIdZkCallBackObj *ptr_callbackobj = new ChannelIdZkCallBackObj();
        ptr_callbackobj->m_ptrChannelidConfig = ptr_channelid_config;
        ptr_callbackobj->m_ptrConfigurator = this;
        int ret = m_channelIdZKConftor->GetConfItemAsync(zk_node_path, NULL, NULL, CONFITEMFLAG_UNSET, Configurator::ChannelValueResultCallBack, ptr_callbackobj);
        if (GSE_SUCCESS != ret)
        {
            LOG_ERROR("it is failed to read the config for the zk node path (%s)", SAFE_CSTR(zk_node_path.c_str()));
            ptr_channelid_config->SetError();
            ptr_channelid_config->AddResponseCount();
            DoUpdateChannelIdConfigRequest(ptr_channelid_config);
            return;
        }
    }

    return;
}

void Configurator::GetChannelsListResultCallBack(std::string &path, int rc, std::vector<std::string> &values, const void *ptr_callbackobj)
{
    if (ptr_callbackobj == NULL)
    {
        LOG_ERROR("get channelid list, callback obj invalid, get channelid path:%s, result:%d", path.c_str(), rc);
        return;
    }

    ChannelIdZkCallBackObj *ptr_channelId_zk_callbakc_obj = (ChannelIdZkCallBackObj *)ptr_callbackobj;
    Configurator *callobj = ptr_channelId_zk_callbakc_obj->m_ptrConfigurator;
    ChannelIDConfig *ptr_channelid_config = ptr_channelId_zk_callbakc_obj->m_ptrChannelidConfig;
    if (callobj == NULL)
    {
        LOG_ERROR("callback caller obj invalid, path:%s, result:%d", path.c_str(), rc);
        return;
    }
    callobj->GetChannelsConfigValue(path, rc, values, ptr_channelid_config);
    delete ptr_channelId_zk_callbakc_obj;
    return;
}

bool Configurator::readChannels(const std::string &channels_path, ChannelIDConfig *ptr_channelid_config, std::string &errorMsg)
{
    ChannelIdZkCallBackObj *ptr_callbackobj = new ChannelIdZkCallBackObj();
    ptr_callbackobj->m_ptrChannelidConfig = ptr_channelid_config;
    ptr_callbackobj->m_ptrConfigurator = this;
    LOG_DEBUG("start read channel config[%s]", channels_path.c_str());
    int ret = m_channelIdZKConftor->GetChildConfItemAsync(channels_path, NULL, NULL, CONFITEMFLAG_UNSET, Configurator::GetChannelsListResultCallBack, ptr_callbackobj);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("it is failed to read the channelid(%s)'s property, zk return(%d)", SAFE_CSTR(channels_path.c_str()), ret);
        ptr_channelid_config->AddResponseCount();
        ptr_channelid_config->SetError();
        delete ptr_callbackobj;
        return false;
    }
    return true;
}

int Configurator::DeleteChannelID(uint32_t channelID)
{
    ZkEvent *event = new ZkEvent();
    event->m_eventType = ZK_EVENT_DELETE;
    uint32_t *ptr_channel_id = new uint32_t;
    *ptr_channel_id = channelID;
    event->m_msg = (void *)ptr_channel_id;
    return DataServer::GetChannelIdManagerInst()->Update(event);
}

bool Configurator::FindChannelID(uint32_t channelId)
{
    return DataServer::GetChannelIdManagerInst()->Find(channelId);
}

int Configurator::deletePlatID(uint32_t platid)
{
    ZkEvent *event = new ZkEvent();
    event->m_eventType = ZK_EVENT_DELETE;
    uint32_t *ptr_channel_id = new uint32_t;
    *ptr_channel_id = platid;
    event->m_msg = (void *)ptr_channel_id;
    return m_platIdManager->Update(event);
}

void Configurator::DelayFreeChannelIdConfig(ChannelIDConfig *ptr_channelid_config)
{
    if (!ptr_channelid_config->m_ptrMetadata->m_isPlatId)
    {
        DataServer::GetChannelIdManagerInst()->FreeChannelIdPtr(ptr_channelid_config);
    }
    else
    {
        m_platIdManager->FreeChannelIdPtr(ptr_channelid_config);
    }
}

void Configurator::DoUpdateChannelIdConfigRequest(ChannelIDConfig *ptr_channelid_config)
{
    if (CanSendUpdateEventMsg(ptr_channelid_config))
    {
        LOG_DEBUG("end read channel config(%s) config event msg, ptr config:%p", ptr_channelid_config->m_strChannelId.c_str(), ptr_channelid_config);
        SendUpdateChannelIdConfigEventMsg(ptr_channelid_config);
    }
    else
    {
        if (!ptr_channelid_config->IsSuccess())
        {
            LOG_DEBUG("failed to parse channelid(%s) config", ptr_channelid_config->m_strChannelId.c_str());
            DelayFreeChannelIdConfig(ptr_channelid_config);
            return;
        }
    }
}

bool Configurator::CanSendUpdateEventMsg(ChannelIDConfig *ptr_channelid_config)
{

    LOG_DEBUG("check channelid[%s] need request count[%d], recv response count[%d], finish flag[%d], success flag[%d]", ptr_channelid_config->m_strChannelId.c_str(), ptr_channelid_config->m_zkReqResponseCount.m_requestCount, ptr_channelid_config->m_zkReqResponseCount.m_responseCount, ptr_channelid_config->IsComplete(), ptr_channelid_config->IsSuccess());

    if (ptr_channelid_config->m_zkReqResponseCount.Finished())
    {
        if (ptr_channelid_config->IsSuccess())
        {
            return true;
        }
        else
        {
            return false;
        }
    }
    else
    {
        if (!ptr_channelid_config->IsSuccess())
        {
            return false;
        }
        else
        {
            return false;
        }
    }

    // return (all_done && ptr_channelid_config->m_success);
}

bool Configurator::IsFinishAndFailed(ChannelIDConfig *ptr_channelid_config)
{
    if (ptr_channelid_config->IsComplete())
    {
        return true;
    }
    else
    {
        return false;
    }
}
void Configurator::SendUpdateChannelIdConfigEventMsg(ChannelIDConfig *ptr_channelid_config)
{
    ZkEvent *event = new ZkEvent();
    event->m_eventType = ZK_EVENT_CHANGE;
    event->m_msg = (void *)ptr_channelid_config;
    if (ptr_channelid_config->m_ptrMetadata->m_isPlatId)
    {
        m_platIdManager->Update(event);
    }
    else
    {
        DataServer::GetChannelIdManagerInst()->Update(event);
    }

    return;
}

int Configurator::updateChannelID(uint32_t channelID, const std::string &context)
{
    std::string errMsg;
    std::string value;

    std::string channelIdStr = gse::tools::strings::ToString(channelID);
    std::string channelIdPath = ZK_CHANNEL_ID_CONFIG_PATH(channelIdStr);
    ChannelIDConfig *ptr_channelid_config = new ChannelIDConfig(channelIdStr);
    std::vector<std::string> list;
    int max_child_node_count = 3; // filter, channel, meta
    ptr_channelid_config->AddNeedAckCount(max_child_node_count);
    // read channel
    if (!readChannels(ZK_CHANNEL_ID_CONFIG_CHANNEL_PATH(channelIdStr), ptr_channelid_config, errMsg))
    {
        LOG_ERROR("it is failed to read channel config for channelid(%s), error info is %s", channelIdStr.c_str(), errMsg.c_str());
        m_channelIdZKConftor->GetConfItemAsync(channelIdPath,
                                               Configurator::watchConfCallBack,
                                               this, CONFITEMFLAG_CHANNELID_CONFIG, NULL, NULL);
        delete ptr_channelid_config;
        return GSE_ERROR;
    }

    // read meta
    if (!readMetadata(ZK_CHANNEL_ID_CONFIG_METADATA_PATH(channelIdStr), ptr_channelid_config, errMsg))
    {
        LOG_ERROR("it is failed to read metadata config for channelid(%s), error info is %s", channelIdPath.c_str(), errMsg.c_str());
        m_channelIdZKConftor->GetConfItemAsync(channelIdPath,
                                               Configurator::watchConfCallBack,
                                               this, CONFITEMFLAG_CHANNELID_CONFIG, NULL, NULL);
        return GSE_ERROR;
    }

    // read filter
    if (!readFilters(ZK_CHANNEL_ID_CONFIG_FILTER_PATH(channelIdStr), ptr_channelid_config, errMsg))
    {
        LOG_ERROR("it is failed to read filter config for channelid(%s), error info is %s", channelIdPath.c_str(), errMsg.c_str());
        m_channelIdZKConftor->GetConfItemAsync(channelIdPath,
                                               Configurator::watchConfCallBack,
                                               this, CONFITEMFLAG_CHANNELID_CONFIG, NULL, NULL);
        return GSE_ERROR;
    }

    m_channelIdZKConftor->GetConfItemAsync(channelIdPath,
                                           Configurator::watchConfCallBack,
                                           this, CONFITEMFLAG_CHANNELID_CONFIG, NULL, NULL);
    return GSE_SUCCESS;
}

int Configurator::updatePlatID(uint32_t plat_id)
{
    std::string errMsg;
    std::string value;

    std::string plat_id_str = gse::tools::strings::ToString(plat_id);
    std::string plat_id_path = ZK_PLAT_ID_CONFIG_PATH(plat_id_str);

    // platid 与channelid 共用相同的数据结构
    ChannelIDConfig *ptr_platIdConfig = new ChannelIDConfig(plat_id_str);
    std::vector<std::string> list;
    int max_child_node_count = 3; // filter, channel, meta
    ptr_platIdConfig->AddNeedAckCount(max_child_node_count);
    if (!readChannels(ZK_PLAT_ID_CONFIG_CHANNEL_PATH(plat_id_str), ptr_platIdConfig, errMsg))
    {
        LOG_ERROR("it is failed to read channel config for channelid(%s), error info is %s", SAFE_CSTR(plat_id_str.c_str()), SAFE_CSTR(errMsg.c_str()));
        m_channelIdZKConftor->GetConfItemAsync(plat_id_path, Configurator::watchConfCallBack, this, CONFITEMFLAG_PLAT_ID_CONFIG_VALUE, NULL, NULL);
        delete ptr_platIdConfig;
        return GSE_ERROR;
    }

    // read meta
    if (!readMetadata(ZK_PLAT_ID_CONFIG_METADATA_PATH(plat_id_str), ptr_platIdConfig, errMsg))
    {
        LOG_ERROR("it is failed to read metadata config for channelid(%s), error info is %s", SAFE_CSTR(plat_id_str.c_str()), SAFE_CSTR(errMsg.c_str()));
        m_channelIdZKConftor->GetConfItemAsync(plat_id_path, Configurator::watchConfCallBack, this, CONFITEMFLAG_PLAT_ID_CONFIG_VALUE, NULL, NULL);
        return GSE_ERROR;
    }

    // read filter
    if (!readFilters(ZK_PLAT_ID_CONFIG_FILTER_PATH(plat_id_str), ptr_platIdConfig, errMsg))
    {
        LOG_ERROR("it is failed to read filter config for channelid(%s), error info is %s", SAFE_CSTR(plat_id_str.c_str()), SAFE_CSTR(errMsg.c_str()));
        m_channelIdZKConftor->GetConfItem(plat_id_path, value, Configurator::watchConfCallBack, this, CONFITEMFLAG_PLAT_ID_CONFIG_VALUE);
        return GSE_ERROR;
    }

    m_channelIdZKConftor->GetConfItemAsync(plat_id_path, Configurator::watchConfCallBack, this, CONFITEMFLAG_PLAT_ID_CONFIG_VALUE, NULL, NULL);
    return GSE_SUCCESS;
}

void Configurator::updateDataConf(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    std::string value = confItem.m_Values[0];
    LOG_INFO("watch dataflow configure channged. key[%s], value[%s]", SAFE_CSTR(confItem.m_Key.c_str()), SAFE_CSTR(value.c_str()));
    pConftor->updateDataFlowConf(value);
}

void Configurator::watchChannelIdConfig(WatchConfItem &confItem)
{
    //注册变化通知回调
    std::vector<std::string> list;
    int ret = m_channelIdZKConftor->GetChildConfItem(confItem.m_Key, list, Configurator::watchConfCallBack, this, CONFITEMFLAG_CHANNELID_LIST);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("failed to get %s child config, ret:%d", confItem.m_Key.c_str(), ret);
        ret = m_channelIdZKConftor->ExistConfItem(confItem.m_Key, Configurator::watchConfCallBack, this, CONFITEMFLAG_CHANNELID_LIST);

        if (ret != GSE_SUCCESS)
        {
            LOG_ERROR("failed to check exist path:%s, ret:%d", confItem.m_Key.c_str(), ret);
        }
    }

    return;
}

void Configurator::watchPlatIdConfig(WatchConfItem &confItem)
{
    std::vector<std::string> list;
    int ret = m_channelIdZKConftor->GetChildConfItem(confItem.m_Key, list, Configurator::watchConfCallBack, this, CONFITEMFLAG_PLAT_ID_CONFIG_LIST);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("failed to get %s child config, ret:%d", confItem.m_Key.c_str(), ret);
        ret = m_channelIdZKConftor->ExistConfItem(confItem.m_Key, Configurator::watchConfCallBack, this, CONFITEMFLAG_PLAT_ID_CONFIG_LIST);

        if (ret != GSE_SUCCESS)
        {
            LOG_ERROR("failed to check exist path:%s, ret:%d", confItem.m_Key.c_str(), ret);
        }
    }
}

void Configurator::handleChannelIdChildListChangeEvt(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    if (confItem.m_valueType == CONFITEMVALUE_TYPE_CREATE)
    {
        LOG_DEBUG("channelid node:%s create, get child nodes again", confItem.m_Key.c_str());
        pConftor->watchChannelIdConfig(confItem);
        return;
    }

    std::size_t max_count = confItem.m_Values.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        std::string channelId = confItem.m_Values.at(idx);
        LOG_INFO("channelid channged. key[%s], value[%s], child node count[%d]", SAFE_CSTR(confItem.m_Key.c_str()), SAFE_CSTR(channelId.c_str()), max_count);
        // check
        if (channelId.compare("index") == 0)
        {
            continue;
        }

        if (!gse::tools::strings::IsNumber(channelId))
        {
            continue;
        }

        ChannelIDType iChannelId = gse::tools::strings::StringToUint32(channelId);

        if (pConftor->FindChannelID(iChannelId))
        {
            LOG_DEBUG("channelid(%d) is exist, don't need update", iChannelId);
            continue;
        }

        if (GSE_SUCCESS != pConftor->updateChannelIDConfigFromZK(iChannelId))
        {
            LOG_ERROR("it is failed to update the channel id(%s)'s config from zookeeper", SAFE_CSTR(channelId.c_str()));
        }
    }

    pConftor->watchChannelIdConfig(confItem);
}
void Configurator::handlePlatIdChildListChangeEvt(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    if (confItem.m_valueType == CONFITEMVALUE_TYPE_CREATE)
    {
        LOG_DEBUG("platid node:%s create, get child nodes again");
        pConftor->watchPlatIdConfig(confItem);
        return;
    }

    std::size_t max_count = confItem.m_Values.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        std::string strPlatId = confItem.m_Values.at(idx);
        LOG_INFO("plat id channged. key[%s], value[%s], child node count[%d]", SAFE_CSTR(confItem.m_Key.c_str()), SAFE_CSTR(strPlatId.c_str()), max_count);

        if (!gse::tools::strings::IsNumber(strPlatId))
        {
            continue;
        }

        if (strPlatId.compare("index") == 0)
        {
            continue;
        }

        uint32_t plat_id = gse::tools::strings::StringToInt32(strPlatId);
        if (GSE_SUCCESS != pConftor->updatePlatID(plat_id))
        {
            LOG_ERROR("it is failed to update the channel id(%s)'s config from zookeeper", SAFE_CSTR(strPlatId.c_str()));
        }
    }

    pConftor->watchPlatIdConfig(confItem);
}
void Configurator::updateChannelIDConfig(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    std::string value = confItem.m_Values[0];
    LOG_INFO("watch channelid configure channged. key[%s], type[%d], value[%s]",
             confItem.m_Key.c_str(), confItem.m_valueType, value.c_str());

    if (confItem.m_valueType == CONFITEMVALUE_TYPE_DELETE)
    {
        pConftor->DeleteChannelID(split_channel_id(confItem.m_Key.c_str()));
    }
    else
    {
        pConftor->updateChannelID(split_channel_id(confItem.m_Key.c_str()), value);
    }
}

void Configurator::updatePlatIdConfig(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    std::string value = confItem.m_Values[0];
    LOG_INFO("watch platID configure channged. key[%s], value[%s]", SAFE_CSTR(confItem.m_Key.c_str()), SAFE_CSTR(value.c_str()));
    if (confItem.m_valueType == CONFITEMVALUE_TYPE_DELETE)
    {
        pConftor->deletePlatID(split_channel_id(confItem.m_Key.c_str()));
    }
    else
    {
        pConftor->updatePlatID(split_channel_id(confItem.m_Key.c_str()));
    }
}

void Configurator::UpdateOpsServiceConfig(const std::string &cfgJson)
{
    OpsServiceConfig opsServiceCfg;
    int ret = opsServiceCfg.ParseJsonConfig(cfgJson);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("failed to parse ops config json, input json:%s", cfgJson.c_str());
        return;
    }

    m_opsConfig.Push(opsServiceCfg.m_serivceId, opsServiceCfg);
}

int Configurator::GetChannelIdByOpsServiceId(int serviceId)
{
    OpsServiceConfig opsServiceConfig;
    bool bFind = m_opsConfig.Find(serviceId, opsServiceConfig);
    if (bFind)
    {
        LOG_DEBUG("get ops channelid:%d by serviceid:%d", opsServiceConfig.m_channelId, serviceId);
        return opsServiceConfig.m_channelId;
    }

    LOG_WARN("service_id:%d is not configured", serviceId);

    return 0;
}

void Configurator::WatchOpsConfigValueChangeEvt(const std::string &path)
{
    //(const std::string &key, std::string &value, FnWatchConf pFnWatchConf, void *lpWatcher, int confItemFlag)
    std::string value;
    m_conftor->GetConfItem(path, value, Configurator::watchConfCallBack, this, CONFITEMFLAG_OPS_SERVICE_CONFIG_VALUE);
}

void Configurator::WatchOpsConfigChildNodeChangeEvt(const std::string &path)
{
    std::vector<std::string> values;
    //(const std::string &key, std::vector<std::string> &values, FnWatchConf pFnWatchConf, void *lpWatcher, int confItemFlag)
    m_conftor->GetChildConfItem(path, values, Configurator::watchConfCallBack, this, CONFITEMFLAG_OPS_SERVICE_CONFIG_LIST);
}

void Configurator::UpdateOpsConfigValue(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    std::string value = confItem.m_Values[0];
    LOG_INFO("watch ops configure channged. key[%s], value[%s]", SAFE_CSTR(confItem.m_Key.c_str()), SAFE_CSTR(value.c_str()));

    pConftor->UpdateOpsServiceConfig(value);
    pConftor->WatchOpsConfigValueChangeEvt(confItem.m_Key);
}

void Configurator::HandleOpsConfigChildListChangeEvt(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    if (confItem.m_valueType == CONFITEMVALUE_TYPE_CREATE)
    {
        LOG_DEBUG("node:%s create, get child nodes again");
        pConftor->WatchOpsConfigChildNodeChangeEvt(confItem.m_Key);
        return;
    }

    std::size_t max_count = confItem.m_Values.size();

    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        std::string opsConfNode = confItem.m_Values.at(idx);
        std::string opsConfNodePath = std::string(DATA_OPS_SERVICE_CONF_PATH) + "/" + opsConfNode;
        LOG_DEBUG("ops config node path:%s", opsConfNodePath.c_str());
        std::string configValue;
        int ret = pConftor->GetOpsConfig(opsConfNodePath, configValue);
        if (ret != GSE_SUCCESS)
        {
            LOG_ERROR("failed to get ops config, zk path:%s, ret:%d", opsConfNodePath.c_str(), ret);
            continue;
        }

        pConftor->UpdateOpsServiceConfig(configValue);
    }

    pConftor->WatchOpsConfigChildNodeChangeEvt(confItem.m_Key);
}

int Configurator::GetOpsConfig(const std::string &zkNodePath, std::string &configValue)
{
    return m_conftor->GetConfItem(zkNodePath, configValue, Configurator::watchConfCallBack, this, CONFITEMFLAG_OPS_SERVICE_CONFIG_VALUE);
}

void Configurator::updateBalanceConfig(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    std::string value = confItem.m_Values[0];
    LOG_INFO("watch balance config channged. key[%s], value[%s]", SAFE_CSTR(confItem.m_Key.c_str()), SAFE_CSTR(value.c_str()));
    pConftor->updateDataServerBalanceConfig(confItem.m_Key, value);
}

void Configurator::watchConfCallBack(WatchConfItem &confItem, void *lpWatcher)
{
    if (confItem.m_Values.size() <= 0)
    {
        LOG_ERROR("the config[%s] has changed, but value is empty", SAFE_CSTR(confItem.m_Key.c_str()));
        return;
    }

    if (confItem.m_confItemFlag <= CONFITEMFLAG_UNSET || confItem.m_confItemFlag >= CONFITEMFLAG_MAX_NUM)
    {
        LOG_ERROR("the config[%s] has changed, but don't know the configure item flag[%d]", SAFE_CSTR(confItem.m_Key.c_str()), confItem.m_confItemFlag);
        return;
    }

    LOG_DEBUG("recv watch the config key (%s) child change event, the type (%d), the value is (%s)", SAFE_CSTR(confItem.m_Key.c_str()), confItem.m_confItemFlag, SAFE_CSTR(confItem.m_Values[0].c_str()));

    Configurator *pConftor = (Configurator *)lpWatcher;
    pConftor->m_ptrConfigFunc[confItem.m_confItemFlag](confItem, lpWatcher);
}

void Configurator::channelNodeChangeCallBack(WatchConfItem &confItem, void *lpWatcher)
{
    if (confItem.m_Values.size() <= 0)
    {
        LOG_ERROR("the config[%s] has changed, but value is empty", SAFE_CSTR(confItem.m_Key.c_str()));
        return;
    }

    if (confItem.m_confItemFlag <= CONFITEMFLAG_UNSET || confItem.m_confItemFlag >= CONFITEMFLAG_MAX_NUM)
    {
        LOG_ERROR("the config[%s] has changed, but don't know the configure item flag[%d]", SAFE_CSTR(confItem.m_Key.c_str()), confItem.m_confItemFlag);
        return;
    }

    LOG_DEBUG("the config key (%s), the type (%d), the value is (%s)", SAFE_CSTR(confItem.m_Key.c_str()), confItem.m_confItemFlag, SAFE_CSTR(confItem.m_Values[0].c_str()));
    if (confItem.m_valueType == CONFITEMVALUE_TYPE_CREATE)
    {
    }
    else if (confItem.m_valueType == CONFITEMVALUE_TYPE_DELETE)
    {
    }

    Configurator *pConftor = (Configurator *)lpWatcher;
    pConftor->m_ptrConfigFunc[confItem.m_confItemFlag](confItem, lpWatcher);
}

template <typename T>
void Configurator::ToJsonObj(rapidjson::Writer<rapidjson::StringBuffer> &writer, T object)
{
    std::size_t count = sizeof(T::properties) / sizeof(T::properties[0]);

    writer.StartObject();
    for (int i = 0; i < count; ++i)
    {
        const std::type_info &targetType = T::properties[i]->TypeID();
        if (targetType == typeid(std::string))
        {
            const Property<T, std::string> *targetProperty = reinterpret_cast<const Property<T, std::string> *>(T::properties[i]);
            writer.Key(targetProperty->name);
            writer.String((object.*(targetProperty->member)).data(), rapidjson::SizeType((object.*(targetProperty->member)).size()));
            continue;
        }

        if (targetType == typeid(int))
        {
            const Property<T, int> *targetProperty = reinterpret_cast<const Property<T, int> *>(T::properties[i]);
            writer.Key(targetProperty->name);
            writer.Int(object.*(targetProperty->member));
            continue;
        }

        if (targetType == typeid(uint32_t))
        {
            const Property<T, uint32_t> *targetProperty = reinterpret_cast<const Property<T, uint32_t> *>(T::properties[i]);
            writer.Key(targetProperty->name);
            writer.Uint(object.*(targetProperty->member));
            continue;
        }

        if (targetType == typeid(double))
        {
            const Property<T, double> *targetProperty = reinterpret_cast<const Property<T, double> *>(T::properties[i]);
            writer.Key(targetProperty->name);
            writer.Double(object.*(targetProperty->member));
            continue;
        }

        if (targetType == typeid(bool))
        {
            const Property<T, bool> *targetProperty = reinterpret_cast<const Property<T, bool> *>(T::properties[i]);
            writer.Key(targetProperty->name);
            writer.Bool(object.*(targetProperty->member));
            continue;
        }
    }

    writer.EndObject();
}

bool Configurator::GetServiceConfig(const std::string &serviceName, ServiceNode &serviceConfig)
{
    DataFlowConf *ptrDataFlowConf = GetDataFlowConf();
    auto _ = gse::tools::defer::finally([ptrDataFlowConf]() {
        delete ptrDataFlowConf;
    });

    if (ptrDataFlowConf == NULL)
    {
        LOG_ERROR("failed to create service  node,dataflow not config");
        return false;
    }

    for (auto it : ptrDataFlowConf->m_channelsConf)
    {
        ChannelConf *channeldConfig = it.second;
        if (channeldConfig != NULL)
        {
            LOG_DEBUG("channel config, channelname:%s, receivername:%s, receiverConf:%p ", channeldConfig->m_name.c_str(), channeldConfig->m_receiverName.c_str(), channeldConfig->m_receiverConf);
            ReceiverConf *receiverConf = m_dataConf->findReceiverConf(channeldConfig->m_receiverName);
            if (receiverConf == NULL)
            {
                continue;
            }

            if (receiverConf->m_name == serviceName)
            {
                serviceConfig.m_serviceIP = DataServer::GetConfigPtr()->GetAdvertiseIp();
                serviceConfig.m_port = receiverConf->m_port;

                if (receiverConf->m_caPath != "" && receiverConf->m_keyPath != "")
                {
                    serviceConfig.m_ssl = true;
                }

                serviceConfig.m_proto = receiverConf->Protostack();
                serviceConfig.m_protocol = receiverConf->ProtocolIDToName();
                LOG_DEBUG("find service config:%s", serviceName.c_str());
                return true;
            }
        }
    }

    return false;
}

std::string Configurator::BuildServiceNodeJson()
{
    rapidjson::Document doc(rapidjson::kObjectType);

    ServiceNode agentService;
    bool hasAgentService = GetServiceConfig("agent_service", agentService);
    rapidjson::Document agentServiceJson;
    if (hasAgentService)
    {
        agentServiceJson = std::move(agentService.ToJsonValue());
        doc.AddMember("agent_service", agentServiceJson, doc.GetAllocator());

        LOG_DEBUG("find agent service config:%s", agentService.ToJsonStr().c_str());
    }

    ServiceNode opsSerivce;
    bool hasOps = GetServiceConfig("ops_service", opsSerivce);
    rapidjson::Document opsServiceJson;
    if (hasOps)
    {
        opsServiceJson = std::move(opsSerivce.ToJsonValue());
        doc.AddMember("ops_service", opsServiceJson, doc.GetAllocator());
        LOG_DEBUG("find ops service config:%s", opsSerivce.ToJsonStr().c_str());
    }

    ServiceNodeMeta meta;
    meta.m_serviceId = DataServer::GetNodeId();
    meta.m_cityId = m_cityID;
    meta.m_zoneId = m_zoneID;
    meta.m_clusterId = m_clusterName;
    meta.m_clusterName = m_clusterName;
    meta.m_advertiseIp = m_selfIp;
    meta.m_serviceName = "data";

    auto metaJson = std::move(meta.ToJsonValue());
    doc.AddMember("metadata", metaJson, doc.GetAllocator());

    NodeLoadBalance loadBalance;
    if (m_systemResourceMonitor != nullptr)
    {
        loadBalance.m_cpuUsage = m_systemResourceMonitor->GetCpuUsage();
        loadBalance.m_memUsage = m_systemResourceMonitor->GetMemUsage();
        loadBalance.m_netUsage = m_systemResourceMonitor->GetNetUsage();
    }
    else
    {
        LOG_ERROR("system resource monitor module not init");
    }

    if (m_systemConnectionMonitor != nullptr)
    {
        loadBalance.m_maxConnectionCount = m_systemConnectionMonitor->GetMaxConnectionCount();
        loadBalance.m_connectionCount = m_systemConnectionMonitor->GetConnectionCount();
    }
    else
    {
        LOG_ERROR("system connection monitor module not init");
    }

    auto loadbalance = std::move(loadBalance.ToJsonValue());
    doc.AddMember("load", loadbalance, doc.GetAllocator());

    rapidjson::StringBuffer strBuffer;
    rapidjson::Writer<rapidjson::StringBuffer> writer(strBuffer);

    doc.Accept(writer);
    LOG_DEBUG("create  service v2 node json:%s", strBuffer.GetString());
    return std::string(strBuffer.GetString());
}

int Configurator::CreateV2ServiceNode()
{
    // create base node
    std::string nodePath = std::string(ZK_SERIVCE_NODE_PATH);
    std::string value = gse::tools::time::GetUTCTimeString();
    int ret = m_conftor->CreateConfItemWithParents(nodePath, value, false);
    if (ret != GSE_SUCCESS)
    {
        LOG_INFO("failed to create base path:%s", ZK_SERIVCE_NODE_PATH);
    }

    nodePath = std::string(ZK_SERIVCE_NODE_PATH) + "/" + DataServer::GetNodeId();

    std::string nodeJson = std::move(BuildServiceNodeJson());
    ret = m_conftor->CreateConfItemWithParents(nodePath, nodeJson, true);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("failed to create node path:%s", ZK_SERIVCE_NODE_PATH);
        return GSE_ERROR;
    }

    return GSE_SUCCESS;
}

int Configurator::CreateOpsServiceConfigNode()
{
    std::string nodePath = std::string(DATA_OPS_SERVICE_CONF_PATH);
    std::string value;
    int ret = m_conftor->CreateConfItemWithParents(nodePath, value, false);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("failed to create zk path:%s, ret:%d", nodePath.c_str(), ret);
        return GSE_ERROR;
    }

    return GSE_SUCCESS;
}

} // namespace data
} // namespace gse
