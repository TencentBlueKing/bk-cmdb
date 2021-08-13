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

#include "conftor/zkconftor.h"
#include "log/log.h"
#include "bbx/gse_errno.h"
#include "tools/macros.h"
#include "tools/hostinfo.h"
#include "tools/strings.h"
#include "tools/filesystem.h"
#include "tools/net.h"

#include "channel_id_config.h"

#include "utils.h"

namespace gse { 
namespace dataserver {

#define DATA_BASE_PATH "/gse/config/server/dataserver"

static uint32_t kDefaultTglogChannelId = 0;

Configurator::Configurator()
{
    m_conftor = NULL;
    m_bkdataZKConftor = NULL;
    m_channelIdZKConftor = NULL;
    m_dataConf = NULL;
    m_ptrEventThread = NULL;

    m_cpuUsage = 0;
    m_hostPerformance = 0 ;
    m_hostLoadBance = 0;
    m_loadweight =0 ;
    m_dloadweight = 0;

    m_channelIdManager = NULL;
    m_channelIdExporterManager = NULL;
    m_ptrConfigFunc[CONFITEMFLAG_DATACONF] = Configurator::updateDataConf;
    m_ptrConfigFunc[CONFITEMFLAG_BALANCE_CONFIG] = Configurator::updateBalanceConfig;
//start-----------------------版本不需要支持data id时需要删除-----------------------
    m_ptrConfigFunc[CONFITEMFLAG_STORAGE] = Configurator::updateStorage;
    m_ptrConfigFunc[CONFITEMFLAG_STORAGE_CONFIG] = Configurator::updateStorageConfig;
    m_ptrConfigFunc[CONFITEMFLAG_DATAID] = Configurator::updateDataID;
    m_ptrConfigFunc[CONFITEMFLAG_DATAID_CONFIG] = Configurator::updateDataIDConfig;

    m_ptrConfigFunc[CONFITEMFLAG_DATAID_FROM_BKDATA] = Configurator::updateDataIDFromBKData;
    m_ptrConfigFunc[CONFITEMFLAG_DATAID_CONFIG_FROM_BKDATA] = Configurator::updateDataIDConfigFromBKData;
    m_ptrConfigFunc[CONFITEMFLAG_STORAGE_FROM_BKDATA] = Configurator::updateStorageFromBKData;
    m_ptrConfigFunc[CONFITEMFLAG_STORAGE_CONFIG_FROM_BKDATA] = Configurator::updateStorageConfigFromBKData;
//end-----------------------版本不需要支持data id时需要删除-----------------------

    m_ptrConfigFunc[CONFITEMFLAG_CHANNELID_LIST] = Configurator::handleChannelIdChildListChangeEvt;
    m_ptrConfigFunc[CONFITEMFLAG_CHANNELID_CONFIG] = Configurator::updateChannelIDConfig;

    m_ptrConfigFunc[CONFITEMFLAG_STREAMTO_CONFIG_LIST] = Configurator::handleExporterChildListChangeEvtFromZK;
    m_ptrConfigFunc[CONFITEMFLAG_STREAMTO_CONFIG_VALUE] = Configurator::handleExporterValueChangeEvtFromZK;

    m_ptrConfigFunc[CONFITEMFLAG_TGLOG_CHANNEL_ID_VALUE] = Configurator::handleTglogChannelIdChangeEvtFromZK;
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

    if (NULL != m_bkdataZKConftor)
    {
        m_bkdataZKConftor->Stop();
        delete m_bkdataZKConftor;
        m_bkdataZKConftor = NULL;
    }

    if (m_channelIdZKConftor != NULL)
    {
        m_channelIdZKConftor->Stop();
        delete m_channelIdZKConftor;
        m_channelIdZKConftor = NULL;
    }

    if (m_channelIdManager != NULL)
    {
        m_channelIdManager->stop();
        delete m_channelIdManager;
        m_channelIdManager = NULL;
    }

    if (m_channelIdExporterManager != NULL)
    {
        m_channelIdExporterManager->stop();
        delete m_channelIdExporterManager;
        m_channelIdExporterManager = NULL;
    }

    for (std::size_t idx = 0; idx < m_callbacks.size(); ++idx)
    {
        delete m_callbacks.at(idx);
    }
    m_callbacks.clear();
}

DataID *Configurator::GetStorageByDataID(uint32_t dataID)
{

    // Attention:兼容V1版本DS处理逻辑而存在
    DataID* ptr_dataid = NULL;
    m_dataIDSV1.Find(dataID, ptr_dataid);
    return ptr_dataid;
}

DataStorage *Configurator::GetStorageByIndex(int storageIndex)
{
    StorageConfigVector *ptr_storage_vector = NULL;
    m_storagesV1.Find(storageIndex, ptr_storage_vector);
    if (NULL == ptr_storage_vector)
    {
        LOG_DEBUG("can not found the storage by the index[%d]", storageIndex);
        return NULL;
    }

    DataStorage *ptr_storage = new DataStorage();
    ptr_storage->m_isDataID = true;
    ptr_storage->m_storage.m_ptrDataIDConfig = ptr_storage_vector;
    return ptr_storage;
}

DataStorage *Configurator::GetStorageByChannelID(uint32_t channelID)
{
    ChannelIDConfig *ptr_channelid_config = NULL;
    m_channelIDS.Find(channelID, ptr_channelid_config);
    if (NULL == ptr_channelid_config)
    {
        LOG_DEBUG("can not found the channelid config by the channelid (%d)", channelID);
        return NULL;
    }

    DataStorage *ptr_storage = new DataStorage();
    ptr_storage->m_isDataID = false;
    ptr_storage->m_storage.m_ptrChannelIDStorage = ptr_channelid_config->ToChannelIDStorage();
    return ptr_storage;
}

void Configurator::GetAllV1StorageID(std::vector<std::string> &storage_id_list)
{
    std::map<StorageIndex, StorageConfigVector *> tmpmap;
    m_storagesV1.Copy(tmpmap);
    for (auto it = tmpmap.begin(); it != tmpmap.end(); it++)
    {
        storage_id_list.push_back(gse::tools::strings::ToString(it->first));
    }

    return ;
}


ChannelIdManager* Configurator::GetChannelIdManager()
{
    return m_channelIdManager;
}

void Configurator::GetAllChannelID(std::vector<std::string> &channelid_list)
{
    std::map<ChannelIDType, ChannelIDConfig *> tmpmap;
    m_channelIDS.Copy(tmpmap);
    for (auto it = tmpmap.begin(); it != tmpmap.end(); it++)
    {
        channelid_list.push_back(gse::tools::strings::ToString(it->first));
    }
    return ;
}
void Configurator::GetAllDataID(std::vector<std::string> &dataid_list)
{
    std::map<DataIDType, DataID *> tmpmap;
    m_dataIDSV1.Copy(tmpmap);
    for (auto it = tmpmap.begin(); it != tmpmap.end(); it++)
    {
        dataid_list.push_back(gse::tools::strings::ToString(it->first));
    }
    return ;
}

DataStorage *Configurator::GetAllStorages()
{
    DataStorage *ptr_storage = new DataStorage();

    // dataid storage config
    std::vector<StorageConfigVector *> storage_config_vector;
    m_storagesV1.AddVector(storage_config_vector);
    std::size_t max_count = storage_config_vector.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        StorageConfigVector *ptr_storage_config = storage_config_vector.at(idx);
        DataStorage *ptr_channelid_storage = new DataStorage();
        ptr_channelid_storage->m_isDataID = true;
        ptr_channelid_storage->m_storage.m_ptrDataIDConfig = ptr_storage_config;
        ptr_storage->SetNext(ptr_channelid_storage);

        LOG_DEBUG("get dataid storage config(curent storage:%p, next:%p)", ptr_storage, ptr_channelid_storage);
    }

    // channelid storage config
//    std::vector<ChannelIDConfig *> channelid_config_vector;
//    m_channelIDS.addVector(channelid_config_vector);
//    max_count = channelid_config_vector.size();
//    for (std::size_t idx = 0; idx < max_count; ++idx)
//    {
//        ChannelIDConfig *ptr_channelid_config = channelid_config_vector.at(idx);
//        DataStorage *ptr_channelid_storage = new DataStorage();
//        ptr_channelid_storage->m_isDataID = false;
//        ptr_channelid_storage->m_storage.m_ptrChannelIDStorage = ptr_channelid_config->ToChannelIDStorage();
//        ptr_storage->SetNext(ptr_channelid_storage);
//    }

    return ptr_storage;
}

void Configurator::WatchUpdateEvent(WatchEventFunc callback, void *args)
{
    EventCallbackParams *ptr_callback = new EventCallbackParams();
    ptr_callback->m_ptrCallbackArgs = args;
    ptr_callback->m_eventCallbackFunc = callback;
    LOG_DEBUG("push update watch event:%p", callback);
    m_callbacks.push_back(ptr_callback);
}

void Configurator::cleanDataID(evutil_socket_t fd, short what, void *args)
{
    Configurator *ptr_this = reinterpret_cast<Configurator *>(args);

    // clear dataid
    int max_count = ptr_this->m_toDeleteDataIDS.Size();
    LOG_DEBUG("clean the invalid dataid by timer, the invalid dataid size %d", max_count);

    while (--max_count >= 0)
    {
        DataID *ptr_dataid = NULL;
        ptr_this->m_toDeleteDataIDS.Pop(ptr_dataid);
        if (NULL == ptr_dataid)
        {
            break;
        }

        if (!ptr_dataid->IsNeedDelete())
        {
            ptr_this->m_toDeleteDataIDS.Push(ptr_dataid);
            continue;
        }
        LOG_DEBUG("stop the invalid dataid (%u) by timer", ptr_dataid->m_dataId);
        delete ptr_dataid;
    }

    // clear channelid
    max_count = ptr_this->m_toDeleteChannelIDS.Size();
    LOG_DEBUG("clean the invalid channel id by timer, the invalid dataid size %d", max_count);

    while (--max_count >= 0)
    {
        ChannelIDConfig *ptr_channelid_conf = NULL;
        ptr_this->m_toDeleteChannelIDS.Pop(ptr_channelid_conf);
        if (NULL == ptr_channelid_conf)
        {
            break;
        }

        if (!ptr_channelid_conf->IsNeedDelete())
        {
            ptr_this->m_toDeleteChannelIDS.Push(ptr_channelid_conf);
            continue;
        }
        LOG_DEBUG("stop the invalid channelid config (%s) by timer", SAFE_CSTR(ptr_channelid_conf->m_originData.c_str()));
        delete ptr_channelid_conf;
    }
}

void Configurator::reportBalance(evutil_socket_t fd, short what, void *args)
{
    Configurator *ptr_this = reinterpret_cast<Configurator *>(args);
    ptr_this->updateDataServerLoadBalanceNode();
}

uint32_t Configurator::getDefaultTglogChannelId()
{
    return kDefaultTglogChannelId;
}

int Configurator::Init(const std::string &dataFlowConf, const std::string &confHost, const std::string &password,
                       const std::string &bkdataZK,
                       const std::string &channelid_zkhost, const std::string &channelid_zkauth,
                       const std::string &selfIp, const std::string &regionID, const std::string &cityID,
                       const std::string &clusterName, const std::string &instanceId, const std::string &watchpath)
{
    LOG_DEBUG("init from remote zk config center");
    m_confHost = confHost;
    m_selfIp = selfIp;
    m_clusterName = clusterName;
    m_instanceId = instanceId;
    m_regionID = regionID;
    m_cityID = cityID;
    m_bkdataZK = bkdataZK;
    m_confHostPassword = password;
    m_localDataFlowConfig = dataFlowConf;
    m_watchPath = watchpath;

    m_channelidZkHost = channelid_zkhost;
    m_channelidZkAuth = channelid_zkauth;

    m_dataConf = new DataConf();

    if (!m_localDataFlowConfig.empty())
    {
        int ret = initFromLocalConfig(m_localDataFlowConfig);
        if (m_confHost.empty())
        {
            LOG_DEBUG("Configurator init finished without zk");
            return ret;
        }
    }


    // connect configure host
    int ret = startConftor();
    if (GSE_SUCCESS != ret)
    {
        return ret;
    }

    m_channelIdExporterManager = new ChannelIdExporterManager();

    m_channelIdExporterManager->Start();

    m_channelIdManager = new ChannelIdManager(m_channelIdExporterManager);
    m_channelIdManager->Start();
    // watch configure item
    watchConfigsFromZK();

    //启动成功后注册服务发现节点
    ret = createBaseConfItem();
    if (GSE_SUCCESS != ret)
    {
        stopConftor();
        return ret;
    }


    return GSE_SUCCESS;
}

int Configurator::initFromLocalConfig(const std::string &dataFlowConf)
{
    std::string config_content;
    int ret_value = gse::tools::filesystem::ReadFromFile(dataFlowConf, config_content);
    if (GSE_SUCCESS != ret_value)
    {
        LOG_ERROR("it is failed to read the dataflow config by the absolute path ( %s ) ", SAFE_CSTR(dataFlowConf.c_str()));
        return ret_value;
    }

    return updateDataFlowConf(config_content);
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
    m_regionID = value.get("regid", "default").asString();

    return;
}

DataFlowConf *Configurator::GetDataFlowConf()
{
    DataFlowConf *pDataFlowConf = NULL;

    pDataFlowConf = m_dataConf->GetDataFlowConf();

    return pDataFlowConf;
}

int Configurator::startConftor()
{
    stopConftor();

    if (NULL == m_conftor)
    {
        ZkConftorParam conftorParam;

        conftorParam.m_ZkHost = m_confHost;
        conftorParam.m_HostIP = m_selfIp;
        conftorParam.m_BasePath = DATA_BASE_PATH;
        conftorParam.m_password = m_confHostPassword;

        m_conftor = new ZkConftor(conftorParam);
    }

    int ret = m_conftor->Start();
    if (GSE_SUCCESS != ret)
    {
        LOG_ERROR("fail to connect configure host[%s], ret=[%d]", SAFE_CSTR(m_confHost.c_str()), ret);
        return ret;
    }

    {
        if (!m_bkdataZK.empty())
        {
            ZkConftorParam conftorParam;
            conftorParam.m_ZkHost = m_bkdataZK;
            m_bkdataZKConftor = new ZkConftor(conftorParam);
            int ret = m_bkdataZKConftor->Start();
            if (GSE_SUCCESS != ret)
            {
                LOG_ERROR("faile to connect the bkdata configure hosts[%s]", SAFE_CSTR(m_bkdataZK.c_str()));
                return ret;
            }
        }
    }

    if (!m_channelidZkHost.empty())
    {
        ZkConftorParam conftorParam;
        conftorParam.m_ZkHost = m_channelidZkHost;
        conftorParam.m_ZkAuth = m_channelidZkAuth;
        conftorParam.m_password = m_channelidZkAuth;
        m_channelIdZKConftor = new ZkConftor(conftorParam);
        int ret = m_channelIdZKConftor->Start();
        if (GSE_SUCCESS != ret)
        {
            LOG_ERROR("faile to connect the channelid zk hosts[%s]", SAFE_CSTR(conftorParam.m_ZkHost.c_str()));
            return ret;
        }
    }
    else
    {
        m_channelIdZKConftor = m_conftor;
    }

    m_ptrEventThread = new rgse::GseEventThread();
    int registerResponseVal = m_ptrEventThread->registerTimerPersistEvent(Configurator::reportBalance, this, 60); // 1 minutes
    if (GSE_SUCCESS > registerResponseVal)
    {
        LOG_ERROR("failed to register timer, error code %d", registerResponseVal);
        return registerResponseVal;
    }

    registerResponseVal = m_ptrEventThread->registerTimerPersistEvent(Configurator::cleanDataID, this, 60); // 1 minutes
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

    LOG_INFO("success to connect configure host[%s]", SAFE_CSTR(m_confHost.c_str()));

    return GSE_SUCCESS;
}

void Configurator::stopConftor()
{
    if (m_conftor != NULL)
    {
        m_conftor->Stop();
        delete m_conftor;
        m_conftor = NULL;
    }

    if (NULL != m_ptrEventThread)
    {
        m_ptrEventThread->stop();
        delete m_ptrEventThread;
        m_ptrEventThread = NULL;
    }

    if (m_channelIdExporterManager != NULL)
    {
        m_channelIdExporterManager->stop();
        delete m_channelIdExporterManager;
        m_channelIdExporterManager = NULL;
    }

    if (m_channelIdManager != NULL)
    {
        m_channelIdManager->stop();
        delete m_channelIdManager;
        m_channelIdManager = NULL;
    }

    if (m_bkdataZKConftor != NULL)
    {
        m_bkdataZKConftor->Stop();
    }
}

int Configurator::createBaseConfItem()
{
    if (NULL == m_conftor)
    {
        LOG_WARN("the pointer of conftor is NULL, please initialize conftor first");
        return GSE_ERROR;
    }

    std::string value = "{\"fd\":1,\"rb\":0}";

    int ret = GSE_SUCCESS;

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
        ret = m_conftor->CreateConfItemWithParents(ipPath, value);
        if (ret != GSE_SUCCESS)
        {
            return ret;
        }
    }

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
int Configurator::updateStorageConfig(StorageIndex storageIndex, const std::string &context)
{
    // parse node to config struct
    StorageConfigVector storageConfigs;
    if (parseStorageNode(storageIndex, context, storageConfigs) != GSE_SUCCESS)
    {
        LOG_ERROR("failed to parse the storage config %s, it may be a invalid json.", SAFE_CSTR(context.c_str()));
        return GSE_JSON_INVALID;
    }

    StorageConfigVector *storage_vector = new StorageConfigVector();
    storage_vector->assign(storageConfigs.begin(), storageConfigs.end());
    StorageConfigVector *ptr_tmp_storagevector = NULL;
    m_storagesV1.Find(storageIndex, ptr_tmp_storagevector);
    LOG_DEBUG("update the storage index %d, config(%s)", storageIndex, context.c_str());
    m_storagesV1.Push(storageIndex, storage_vector);
    if (NULL != ptr_tmp_storagevector)
    {
        delete ptr_tmp_storagevector;
    }

    // notition
    if (0 != m_callbacks.size())
    {
        // ATTETION:此处仅更新 storage index ，dataid 设置为 0 表示不需要更新，后期将dataid 和channelid 合并后此处逻辑可以不要
        std::size_t max_count = m_callbacks.size();
        for (std::size_t idx = 0; idx < max_count; ++idx)
        {
            EventCallbackParams *ptr_callback = m_callbacks.at(idx);
            ptr_callback->m_eventCallbackFunc(ptr_callback->m_ptrCallbackArgs, storageIndex, 0);
        }
    }
    return GSE_SUCCESS;
}
int Configurator::updateStorageConfigFromZK(StorageIndex storageIndex, const std::string &nodePath)
{
    // get node from zk
    std::string nodeValue;
    if (m_conftor->GetConfItem(nodePath, nodeValue, Configurator::watchConfCallBack, this, CONFITEMFLAG_STORAGE_CONFIG) != GSE_SUCCESS)
    {
        LOG_ERROR("CAN NOT GET USER CLUSTER INDEX, KEY:%s VALUE:%s", nodePath.c_str(), SAFE_CSTR(nodeValue.c_str()));
        return GSE_ERROR;
    }

    return updateStorageConfig(storageIndex, nodeValue);
}

int Configurator::updateStorageConfigFromBKDataZK(StorageIndex storageIndex, const std::string &nodePath)
{
    if (NULL == m_bkdataZKConftor)
    {
        return GSE_SUCCESS;
    }

    // get node from zk
    std::string nodeValue;
    if (m_bkdataZKConftor->GetConfItem(nodePath, nodeValue, Configurator::watchConfCallBack, this, CONFITEMFLAG_STORAGE_CONFIG_FROM_BKDATA) != GSE_SUCCESS)
    {
        LOG_ERROR("CAN NOT GET USER CLUSTER INDEX, KEY:%s VALUE:%s", nodePath.c_str(), SAFE_CSTR(nodeValue.c_str()));
        return GSE_ERROR;
    }

    return updateStorageConfig(storageIndex, nodeValue);
}

int Configurator::watchStorageFromBKData()
{
    if (NULL == m_bkdataZKConftor)
    {
        return GSE_SUCCESS;
    }

    std::vector<string> userNodeList;
    std::string strStorageKey("/config/leaf/kafka");
    int ret = m_bkdataZKConftor->GetChildConfItem(strStorageKey, userNodeList, Configurator::watchConfCallBack, this, CONFITEMFLAG_STORAGE_FROM_BKDATA);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("can not update from zk:%s", strStorageKey.c_str());
        return ret;
    }

    for (int userIndex = 0; userIndex < userNodeList.size(); ++userIndex)
    {
        string userCluster(userNodeList[userIndex]);
        string strUserKey(strStorageKey + "/" + userCluster);
        LOG_INFO("read the storage config from zk node :%s", strUserKey.c_str());
        ret = updateStorageConfigFromBKDataZK(split_storage_index(userCluster.c_str()), strUserKey);
        if (GSE_SUCCESS != ret)
        {
            LOG_ERROR("can not update from zk, %s", strUserKey.c_str());
        }
    }

    return GSE_SUCCESS;
}
int Configurator::watchStorage()
{
    watchStorageFromBKData();

    std::vector<string> userNodeList;
    std::string strStorageKey("/gse/config/etc/dataserver/storage/all");
    int ret = m_conftor->GetChildConfItem(strStorageKey, userNodeList, Configurator::watchConfCallBack, this, CONFITEMFLAG_STORAGE);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("can not update from zk:%s", strStorageKey.c_str());
        return ret;
    }

    for (int userIndex = 0; userIndex < userNodeList.size(); ++userIndex)
    {
        string userCluster(userNodeList[userIndex]);
        string strUserKey(strStorageKey + "/" + userCluster);
        LOG_INFO("read the storage config from zk node :%s", strUserKey.c_str());
        ret = updateStorageConfigFromZK(split_storage_index(userCluster.c_str()), strUserKey);
        if (GSE_SUCCESS != ret)
        {
            LOG_ERROR("can not update from zk, %s", strUserKey.c_str());
        }
    }

    return GSE_SUCCESS;
}


int Configurator::watchExporter()
{
    watchStorageFromBKData();

    std::vector<string> userNodeList;
    std::string strStorageKey("/gse/config/etc/dataserver/storage/all");
    int ret = m_conftor->GetChildConfItem(strStorageKey, userNodeList, Configurator::watchConfCallBack, this, CONFITEMFLAG_STORAGE);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("can not update from zk:%s", strStorageKey.c_str());
        return ret;
    }

    for (int userIndex = 0; userIndex < userNodeList.size(); ++userIndex)
    {
        string userCluster(userNodeList[userIndex]);
        string strUserKey(strStorageKey + "/" + userCluster);
        LOG_INFO("read the storage config from zk node :%s", strUserKey.c_str());
        ret = updateStorageConfigFromZK(split_storage_index(userCluster.c_str()), strUserKey);
        if (GSE_SUCCESS != ret)
        {
            LOG_ERROR("can not update from zk, %s", strUserKey.c_str());
        }
    }

    return GSE_SUCCESS;
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

int Configurator::deleteStreamToIdConfigFromZK(std::string &stream_to_id)
{
    std::string property_value;
    ChannelIdExporterConfig *ptr_stream_to_id_config = new ChannelIdExporterConfig();

    std::string stream_to_id_path = ZK_STREAM_ID_CONFIG_PATH(stream_to_id);
    uint32_t u_stream_to_id = gse::tools::strings::StringToUint32(stream_to_id);
    uint32_t *ptr_stream_to_id = new uint32_t;
    *ptr_stream_to_id = u_stream_to_id;
    ZkEvent *event = new ZkEvent();
    event->m_eventType = ZK_EVENT_DELETE;
    event->m_msg = (void*)ptr_stream_to_id;
    if (m_channelIdExporterManager->UpdateExporterConfig(event) != GSE_SUCCESS)
    {
        delete ptr_stream_to_id_config;
        delete event;
        return GSE_ERROR;
    }

    return GSE_SUCCESS;
}

int Configurator::updateStreamToIdConfigFromZK(std::string &stream_to_id)
{
    std::string property_value;
    ChannelIdExporterConfig *ptr_stream_to_id_config = new ChannelIdExporterConfig();

    std::string stream_to_id_path = ZK_STREAM_ID_CONFIG_PATH(stream_to_id);
    uint32_t u_stream_to_id = gse::tools::strings::StringToUint32(stream_to_id);

    m_channelIdZKConftor->GetConfItem(stream_to_id_path, property_value, Configurator::watchConfCallBack, this, CONFITEMFLAG_STREAMTO_CONFIG_VALUE);

    LOG_DEBUG("get exporter:%s ,value:%s", stream_to_id_path.c_str(), property_value.c_str());
    Json::Value property_cfg_json;
    Json::Reader channel_cfg_json_reader(Json::Features::strictMode());
    std::string error_msg;
    if (!channel_cfg_json_reader.parse(property_value, property_cfg_json))
    {
        error_msg = "the channel id (" + stream_to_id_path + ")'s config is invalid";
        LOG_DEBUG("the channel id (%s)'s config is invalid", error_msg.c_str());
        delete ptr_stream_to_id_config;
        return GSE_ERROR;
    }

    if (!m_channelIDConfigFactory.ParseExporterConfig(property_cfg_json, ptr_stream_to_id_config, error_msg))
    {
        delete ptr_stream_to_id_config;
        return GSE_ERROR;
    }
    ptr_stream_to_id_config->m_streamToId = u_stream_to_id;
    ZkEvent *event = new ZkEvent();
    event->m_eventType = ZK_EVENT_CHANGE;
    event->m_msg = (void*)ptr_stream_to_id_config;

    if (m_channelIdExporterManager->UpdateExporterConfig(event) != GSE_SUCCESS)
    {
        delete event;
        delete ptr_stream_to_id_config;
        return GSE_ERROR;
    }

    return GSE_SUCCESS;
}

void Configurator::handleTglogChannelIdChangeEvtFromZK(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    LOG_INFO("tglog channelid value change event, zk path:%s", confItem.m_Key.c_str());
    if (confItem.m_valueType == ZK_EVENT_DELETE)
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
    //TGLOG watch, TGLOG only support one channel id
    if (m_watchPath != "")
    {
        return WatchTglogChannelId();
    }

    std::vector<std::string> channel_id_list;
    std::string root_path = ZK_CHANNEL_ID_CONFIG_BASE_PATH;
    LOG_DEBUG("channel id root path:%s", root_path.c_str());

    int ret = m_channelIdZKConftor->GetChildConfItem(root_path, channel_id_list, Configurator::watchConfCallBack, this, CONFITEMFLAG_CHANNELID_LIST);
    if (ret != GSE_SUCCESS)
    {
        LOG_WARN("failed to get channel id root path from zk:%s", root_path.c_str());
        ret = m_channelIdZKConftor->ExistConfItem(root_path, Configurator::watchConfCallBack, this, CONFITEMFLAG_CHANNELID_LIST);
        return ret;
    }

    LOG_INFO("Start to read the channel id config from zk node");
    for (int channel_id_idx = 0; channel_id_idx < channel_id_list.size(); ++channel_id_idx)
    {
        std::string channelID(channel_id_list[channel_id_idx]);
        std::string str_channelid_path = root_path + "/" + channelID;
        uint32_t channel_id = gse::tools::strings::StringToUint32(channelID);
        ret = updateChannelIDConfigFromZK(channel_id);
        if (GSE_SUCCESS != ret)
        {
            LOG_ERROR("can not update channelid config from zk, %s", str_channelid_path.c_str());
        }
    }

    LOG_INFO("Finish reading the channel id config from zk node");
    channel_id_list.clear();
    //register  watch path
    m_channelIdZKConftor->GetChildConfItem(root_path, channel_id_list, Configurator::watchConfCallBack, this, CONFITEMFLAG_CHANNELID_LIST);
    return GSE_SUCCESS;
}


int Configurator::watchStreamToID()
{
    std::vector<std::string> streamto_id_list;
    std::string root_path = ZK_STREAM_ID_CONFIG_BASE_PATH;
    LOG_DEBUG("streamto id root path:%s", root_path.c_str());

    int ret = m_channelIdZKConftor->GetChildConfItem(root_path, streamto_id_list, NULL, this, CONFITEMFLAG_STREAMTO_CONFIG_LIST);
    if (ret != GSE_SUCCESS)
    {
        LOG_WARN("failed to get streamto id root path from zk:%s", root_path.c_str());
        ret = m_channelIdZKConftor->ExistConfItem(root_path, Configurator::watchConfCallBack, this, CONFITEMFLAG_STREAMTO_CONFIG_LIST);
        return ret;
    }

    for (int streamto_id_idx = 0; streamto_id_idx < streamto_id_list.size(); ++streamto_id_idx)
    {
        std::string stream_to_id(streamto_id_list[streamto_id_idx]);
        std::string stream_to_id_zk_path  = root_path + "/" + stream_to_id;
        LOG_INFO("read the streamto id config from zk node :%s", stream_to_id_zk_path.c_str());
        ret = updateStreamToIdConfigFromZK(stream_to_id);
        if (GSE_SUCCESS != ret)
        {
            LOG_ERROR("failed to update stream to id config, zk path:%s", stream_to_id_zk_path.c_str());
            continue;
        }
    }

    //register path watch
    m_channelIdZKConftor->GetChildConfItem(root_path, streamto_id_list, Configurator::watchConfCallBack, this, CONFITEMFLAG_STREAMTO_CONFIG_LIST);
    return GSE_SUCCESS;
}

int Configurator::watchBalanceConfig()
{
    // 获取负载配置参数
    std::string cfgNode;
    do
    {

        std::string balancfg_key("/gse/config/etc/dataserver/" + m_selfIp + "/balancecfg");
        int ret = m_conftor->GetConfItem(balancfg_key, cfgNode, Configurator::watchConfCallBack, this, CONFITEMFLAG_BALANCE_CONFIG);
        if (ret == GSE_SUCCESS)
        {
            break;
        }

        LOG_WARN("it is failed to get balance config from the node (%s), please the node in zk", SAFE_CSTR(balancfg_key.c_str()));

        balancfg_key = "/gse/config/etc/dataserver/all/balancecfg";
        ret = m_conftor->GetConfItem(balancfg_key, cfgNode, Configurator::watchConfCallBack, this, CONFITEMFLAG_BALANCE_CONFIG);
        if (ret == GSE_SUCCESS)
        {
            break;
        }

        LOG_ERROR("it is failed to get balance config from the node (%s), please the node in zk", SAFE_CSTR(balancfg_key.c_str()));
        return ret;
    } while (false);

    return updateDataServerBalanceConfig(cfgNode);
}

int Configurator::updateDataServerBalanceConfig(const std::string &context)
{
    // 解析负载配置参数
    LOG_INFO("balance config: %s", (context.c_str()));
    Json::Reader reader(Json::Features::strictMode());
    if (!reader.parse(context, m_cfgLoadBalanceValue, false))
    {
        LOG_ERROR("parse to json failed: %s", (context.c_str()));
        return GSE_JSON_INVALID;
    }

    return GSE_SUCCESS;
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
        LOG_WARN("CAN NOT GET NET DEV INFO");
        return GSE_ERROR;
    }

    LOG_INFO("NET DEV NAME:%s RX_BYTES:%llu", currNetDevItems.at(0).devname, currNetDevItems.at(0).rx_bytes);
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
    //int netwide = LocalStatus::sysapi_get_net_wide(m_cfgLoadBalanceValue.get("netdev", "eth1").asCString());gse::tools::GetNetwide

    int netwide = getNetWide(ethname.c_str());
    //int diskmaxio = LocalStatus::sysapi_get_disk_maxio();

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

    m_hostPerformance = markcpu  + markmem + marknet;

    // 计算当前网络流量，单位 Mb
    float diskUsage = 0;
    int mem_usage = total - free;
    float netUsagePer = (((netReadBytes * 8) / 1024) / 1024) / dtimestamp; // Mb

    // 为各个，cpu 按照使用率计算占用性能指数，并以 cpup 为单位进行打分
    float cpuUsageMark = (m_cfgLoadBalanceValue.get("cpur", 0).asFloat() * m_cpuUsage * cpunum * cpufreq) / m_cfgLoadBalanceValue.get("cpup", 1.0).asFloat();
    float memUsageMark = m_cfgLoadBalanceValue.get("memr", 0).asFloat() * mem_usage / m_cfgLoadBalanceValue.get("memp", 1.0).asFloat();
    float netUsageMark = m_cfgLoadBalanceValue.get("netr", 0).asFloat() * netUsagePer / m_cfgLoadBalanceValue.get("netp", 1.0).asFloat();

    m_hostLoadBance = cpuUsageMark + memUsageMark  + netUsageMark;

    LOG_INFO("origin usage cpu:%f mem:%d disk:%f net:%f", m_cpuUsage, mem_usage, netUsagePer);

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

    // set all node
    string dsAllNode = "/gse/config/server/dataserver/all/" + m_selfIp;
    if (m_conftor->SetConfItem(dsAllNode, value) != GSE_SUCCESS)
    {
        if (m_conftor->CreateConfItemWithParents(dsAllNode, value) != GSE_SUCCESS)
        {
            LOG_ERROR("set node %s  value %s failed", SAFE_CSTR(dsAllNode.c_str()), SAFE_CSTR(value.c_str()));
        }
    }

    //set origin node
    std::string dsNode = "/gse/config/server/dataserver/" + m_regionID + "/" + m_cityID + "/" + m_selfIp;
    if (m_conftor->SetConfItem(dsNode, value) != GSE_SUCCESS)
    {
        if (m_conftor->CreateConfItemWithParents(dsNode, value) != GSE_SUCCESS)
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
            if (m_conftor->CreateConfItemWithParents(dsNodeNew, value) != GSE_SUCCESS)
            {
                LOG_ERROR("set node %s value %s failed", SAFE_CSTR(dsNodeNew.c_str()), SAFE_CSTR(value.c_str()));
            }
        }
    }

    return GSE_SUCCESS;
}

int Configurator::updateDataIDConfigFromZK(const std::string &nodePath, const std::string &dataID)
{
    LOG_INFO("will to get the dataid config by the node:%s", SAFE_CSTR(nodePath.c_str()));
    std::string nodeValue;
    int ret = m_conftor->GetConfItem(nodePath, nodeValue, Configurator::watchConfCallBack, this, CONFITEMFLAG_DATAID_CONFIG);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("it is failed to get dataid config from the node (%s), please the node in zk", SAFE_CSTR(nodePath.c_str()));
        return ret;
    }
    uint32_t data_id = gse::tools::strings::StringToUint32(dataID);
    LOG_INFO("the dataid node (%s->%u) config is %s ", SAFE_CSTR(nodePath.c_str()), data_id, SAFE_CSTR(nodeValue.c_str()));
    return updateDataID(data_id, nodeValue);
}

int Configurator::updateDataIDConfigFromBKDataZK(const std::string &nodePath, const std::string &dataID)
{
    if (NULL == m_bkdataZKConftor)
    {
        return GSE_SUCCESS;
    }
    LOG_INFO("will to get the bkdata dataid config by the node:%s", SAFE_CSTR(nodePath.c_str()));
    std::string nodeValue;
    int ret = m_bkdataZKConftor->GetConfItem(nodePath, nodeValue, Configurator::watchConfCallBack, this, CONFITEMFLAG_DATAID_CONFIG_FROM_BKDATA);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("it is failed to get dataid config from the node (%s), please the node in zk", SAFE_CSTR(nodePath.c_str()));
        return ret;
    }
    uint32_t data_id = gse::tools::strings::StringToUint32(dataID);
    LOG_INFO("the bkdata dataid node (%s->%u) config is %s ", SAFE_CSTR(nodePath.c_str()), data_id, SAFE_CSTR(nodeValue.c_str()));
    return updateDataID(data_id, nodeValue);
}

int Configurator::watchDataIDFromBKData()
{
    if (NULL == m_bkdataZKConftor)
    {
        return GSE_SUCCESS;
    }

    std::vector<string> userNodeList;

    string strDataIdKey("/config/leaf/data");
    int ret = m_bkdataZKConftor->GetChildConfItem(strDataIdKey, userNodeList, Configurator::watchConfCallBack, this, CONFITEMFLAG_DATAID_FROM_BKDATA);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("it is failed to get dataid config from zk, please to check the node (/config/leaf/data) in the zk");
        return ret;
    }

    for (int userIndex = 0; userIndex < userNodeList.size(); ++userIndex)
    {
        std::string userDataId(userNodeList[userIndex]);
        if (!gse::tools::strings::IsNumber(userDataId))
        {
            LOG_FATAL("DATA ID IS NOT INT:%s", SAFE_CSTR(userDataId.c_str()));
            continue;
        }

        std::string strUserKey(strDataIdKey + "/" + userDataId);
        ret = updateDataIDConfigFromBKDataZK(strUserKey, userDataId);
        if (GSE_SUCCESS != ret)
        {
            LOG_ERROR("can not update from zk, %s", strUserKey.c_str());
        }
    }

    return GSE_SUCCESS;
}
int Configurator::watchDataID()
{
    watchDataIDFromBKData();

    std::vector<string> userNodeList;

    string strDataIdKey("/gse/config/etc/dataserver/data");
    int ret = m_conftor->GetChildConfItem(strDataIdKey, userNodeList, Configurator::watchConfCallBack, this, CONFITEMFLAG_DATAID);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("it is failed to get dataid config from zk, please to check the node (/gse/config/etc/dataserver/data) in the zk");
        return ret;
    }

    for (int userIndex = 0; userIndex < userNodeList.size(); ++userIndex)
    {
        std::string userDataId(userNodeList[userIndex]);
        if (!gse::tools::strings::IsNumber(userDataId))
        {
            LOG_FATAL("DATA ID IS NOT INT:%s", SAFE_CSTR(userDataId.c_str()));
            continue;
        }

        std::string strUserKey(strDataIdKey + "/" + userDataId);
        ret = updateDataIDConfigFromZK(strUserKey, userDataId);
        if (GSE_SUCCESS != ret)
        {
            LOG_ERROR("can not update from zk, %s", strUserKey.c_str());
        }
    }

    return GSE_SUCCESS;
}

int Configurator::watchConfigsFromZK()
{

    if (m_regionID.empty() || m_cityID.empty())
    {
        updateLocationFromZK();
    }

    int ret = GSE_SUCCESS;
    ret = watchDataFlow();
    if (ret != GSE_SUCCESS)
    {
        return ret;
    }

    ret = watchStreamToID();
    if ( ret != GSE_SUCCESS)
    {
        return ret;
    }

    ret = watchChannelID();
    if (ret != GSE_SUCCESS)
    {
        return ret;
    }

//    if (ret = watchStorage() && ret != GSE_SUCCESS)
//    {
//        return ret;
//    }

//    if (ret = watchDataID() && ret != GSE_SUCCESS)
//    {
//        return ret;
//    }
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


void Configurator::ParseMetaConfig(std::string &path, int rc, std::string &property_value, ChannelIDConfig* ptr_channelid_config)
{
    Json::Value property_cfg_json;
    Json::Reader channel_cfg_json_reader;
    std::string error_msg;
    ptr_channelid_config->AddResponseCount();
    if (!channel_cfg_json_reader.parse(property_value, property_cfg_json))
    {
        LOG_ERROR("the channel id (%s)'s meta config(%s) is invalid", path.c_str(), property_value.c_str());
        ptr_channelid_config->SetError();
        DoUpdateChannelIdConfigRequest(ptr_channelid_config);
        return;
    }

    if (!m_channelIDConfigFactory.ParseMetadata(property_cfg_json, ptr_channelid_config->m_ptrMetadata, error_msg))
    {
        ptr_channelid_config->SetError();
        LOG_ERROR("failed to parse channel id (%s)'s meta config(%s), error(%s)", path.c_str(), property_value.c_str(), error_msg.c_str());
        DoUpdateChannelIdConfigRequest(ptr_channelid_config);
        return;
    }

    DoUpdateChannelIdConfigRequest(ptr_channelid_config);
    LOG_DEBUG("parse meta config success, path(%s) value(%s)", path.c_str(), property_value.c_str());
}
//typedef void (*FnZkGetValueCallBack)(std::string &path, int rc, const char *value, int32_t value_len, const struct Stat *stat, const void *data);
void Configurator::GetMetaValueCallBack(std::string &path, int rc, const char *value, int32_t value_len, const struct Stat *stat, const void* ptr_callbackobj)
{
    if (ptr_callbackobj == NULL)
    {
        LOG_ERROR("get mata value input param invalid , callback obj is null", ptr_callbackobj);
        return;
    }

    if (value == NULL || value_len <=0)
    {
        LOG_ERROR("get mata value invalid, value le:%d", value_len);
        return;
    }

    std::string str_value(value, value_len);
    ChannelIdZkCallBackObj *ptr_channelId_zk_callbakc_obj = (ChannelIdZkCallBackObj *)ptr_callbackobj;
    Configurator *callobj = ptr_channelId_zk_callbakc_obj->m_ptrConfigurator;
    callobj->ParseMetaConfig(path, rc, str_value, ptr_channelId_zk_callbakc_obj->m_ptrChannelidConfig);
    delete ptr_channelId_zk_callbakc_obj;
}


bool Configurator::readMetadata(const std::string &metadataPath, ChannelIDConfig *ptr_channelid_config, std::string &errorMsg)
{
    ChannelIdZkCallBackObj *ptr_callbackobj = new ChannelIdZkCallBackObj();
    ptr_callbackobj->m_ptrChannelidConfig = ptr_channelid_config;
    ptr_callbackobj->m_ptrConfigurator = this;
   int ret = m_channelIdZKConftor->GetConfItemAsync(metadataPath, NULL, NULL, CONFITEMFLAG_UNSET, Configurator::GetMetaValueCallBack, ptr_callbackobj);
    if (GSE_SUCCESS != ret)
    {
        LOG_ERROR("it is failed to read the config for the zk node path (%s)", SAFE_CSTR(metadataPath.c_str()));
        ptr_channelid_config->SetError();
        ptr_channelid_config->AddResponseCount();
        return false;
    }
    return true;
}

void Configurator::UpdateFilterValue(std::string &path, int rc, std::string &property_value, void* ptr_callbackobj)
{
    ChannelIDConfig *ptr_channelid_config = (ChannelIDConfig *)ptr_callbackobj;
    std::string channel_id_str = ptr_channelid_config->m_strChannelId;

    ptr_channelid_config->AddResponseCount();
    if (rc != 0)
    {
        LOG_ERROR("get filter[%s] return error:%d", path.c_str(), rc);
        ptr_channelid_config->SetError();
        DoUpdateChannelIdConfigRequest(ptr_channelid_config);
        return;
    }

    LOG_DEBUG("get filters:%s ,value:%s", path.c_str(), property_value.c_str());
    Json::Value property_cfg_json;
    Json::Reader channel_cfg_json_reader;
    if (!channel_cfg_json_reader.parse(property_value, property_cfg_json))
    {
        LOG_ERROR("get filter[%s] is invalid json", path.c_str(), property_value.c_str());
        ptr_channelid_config->SetError();
        DoUpdateChannelIdConfigRequest(ptr_channelid_config);
        return;
    }

    Json::Value tmp;
    std::string errorMsg;
    tmp.append(property_cfg_json);
    if (!m_channelIDConfigFactory.ParseStreamFilter(tmp, ptr_channelid_config->m_streamFilter, errorMsg))
    {
        ptr_channelid_config->SetError();
        DoUpdateChannelIdConfigRequest(ptr_channelid_config);
        LOG_ERROR("failed to parse stream filter, path[%s], value[%s] error[%s]", path.c_str(), property_value.c_str(), errorMsg.c_str());
        return;
    }

    DoUpdateChannelIdConfigRequest(ptr_channelid_config);
    return;
}
//typedef void (*FnZkGetValueCallBack)(std::string &path, int rc, const char *value, int32_t value_len, const struct Stat *stat, const void *data);
void Configurator::FilterValueResultCallBack(std::string &path, int rc, const char *value, int32_t value_len, const struct Stat *stat, const void *ptr_callbackobj)
{
    if (ptr_callbackobj == NULL)
    {
        return;
    }

    if (value == NULL || value_len <=0)
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

void Configurator::GetFiltersValue(std::string &path, int rc, std::vector<std::string> &values, ChannelIDConfig * ptr_channelid_config)
{
    if (ptr_channelid_config == NULL)
    {
        return;
    }

    //TOP 1 level response
    ptr_channelid_config->AddResponseCount();
    if (values.size() == 0)
    {
        //channels must > 0
        LOG_INFO("No filter rules for channelid(%s) configured", path.c_str());
        DoUpdateChannelIdConfigRequest(ptr_channelid_config);
        return;
    }

    LOG_DEBUG("filters[%s] has child node:%d", path.c_str(), values.size());
    std::string channel_id_str = ptr_channelid_config->m_strChannelId;
    std::size_t max_cnt = values.size();
    ptr_channelid_config->AddNeedAckCount(max_cnt);
    std::string filter_path = ZK_CHANNEL_ID_CONFIG_FILTER_PATH(ptr_channelid_config->m_strChannelId);

    for (std::size_t idx = 0; idx < max_cnt; ++idx)
    {
        std::string property_value;
        std::string zk_node_path(filter_path + "/" + values.at(idx));

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

void Configurator::GetFiltersResultCallBack(std::string &path, int rc, std::vector<std::string> &values, const void* ptr_callbackobj)
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

bool Configurator::readExporters(const std::string &exporterPath, std::vector<StreamConfig *> &exporters, std::string &errorMsg)
{
    std::vector<std::string> property_vector;
    int ret = m_conftor->GetChildConfItem(exporterPath, property_vector, NULL, NULL, CONFITEMFLAG_UNSET);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("it is failed to read the channelid(%s)'s property", SAFE_CSTR(exporterPath.c_str()));
        return false;
    }

    std::size_t max_cnt = property_vector.size();
    for (std::size_t idx = 0; idx < max_cnt; ++idx)
    {
        std::string property_value;
        std::string zk_node_path(exporterPath + "/" + property_vector.at(idx));
        int ret = m_conftor->GetConfItem(zk_node_path, property_value, NULL, NULL, CONFITEMFLAG_UNSET);
        if (GSE_SUCCESS != ret)
        {
            LOG_ERROR("it is failed to read the config for the zk node path (%s)", SAFE_CSTR(zk_node_path.c_str()));
            continue;
        }
        LOG_DEBUG("get exporters:%s ,value:%s", zk_node_path.c_str(), property_value.c_str());
        Json::Value property_cfg_json;
        Json::Reader channel_cfg_json_reader;
        if (!channel_cfg_json_reader.parse(property_value, property_cfg_json))
        {
            errorMsg = "the channel id (" + exporterPath + ")'s config is invalid";
            continue;
        }
        Json::Value tmp;
        tmp.append(property_cfg_json);
        if (!m_channelIDConfigFactory.ParseStreamConfig(tmp, exporters, errorMsg))
        {
            return false;
        }
    }

    return true;
}


//typedef void (*FnZkGetChildCallBack)(std::string &path, int rc, std::string<std::string> &values, void* data);
void Configurator::GetChannelConfigValueCallBack(std::string &path, int rc, std::string &property_value, void* ptr_callbackobj)
{
    ChannelIDConfig *ptr_channelid_config = (ChannelIDConfig *)ptr_callbackobj;
    std::string channel_id_str = ptr_channelid_config->m_strChannelId;

    ptr_channelid_config->AddResponseCount();
    LOG_DEBUG("get channel[%s] value(%s)", path.c_str(), property_value.c_str());
    if (rc != 0)
    {
        LOG_ERROR("get value return error:%d", rc);
        ptr_channelid_config->SetError();
        DoUpdateChannelIdConfigRequest(ptr_channelid_config);
        return;
    }

    Json::Value property_cfg_json;
    Json::Reader channel_cfg_json_reader;
    if (!channel_cfg_json_reader.parse(property_value, property_cfg_json))
    {
        ptr_channelid_config->SetError();
        LOG_ERROR("the channel id (%s)'s config is invalid", path.c_str());
        ptr_channelid_config->SetError();
        DoUpdateChannelIdConfigRequest(ptr_channelid_config);
        return;
    }

    Json::Value tmp;
    tmp.append(property_cfg_json);
    std::string errorMsg;
    if (!m_channelIDConfigFactory.ParseStreamTo(tmp, ptr_channelid_config->m_streamTo, errorMsg))
    {
        LOG_ERROR("the channel id (%s)'s config is invalid, errmsg(%s)", path.c_str(), errorMsg.c_str());
        ptr_channelid_config->SetError();
        DoUpdateChannelIdConfigRequest(ptr_channelid_config);
        return;
    }

    DoUpdateChannelIdConfigRequest(ptr_channelid_config);
    return;
}
void Configurator::ChannelValueResultCallBack(std::string &path, int rc, const char *value, int32_t value_len, const struct Stat *stat, const void* ptr_callbackobj)
{
    if (ptr_callbackobj == NULL)
    {
        LOG_ERROR("get channel value result failed ,callback obj is null");
        return;
    }

    if (value == NULL || value_len <=0)
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

void Configurator::GetChannelsConfigValue(std::string &path, int rc, std::vector<std::string> &values, ChannelIDConfig * ptr_channelid_config)
{
    if (ptr_channelid_config == NULL)
    {
        LOG_ERROR("channelid object invalid, path:%s", path.c_str());
        return;
    }

    ptr_channelid_config->AddResponseCount();
    if (values.size() == 0)
    {
        //channels must > 0
        LOG_ERROR("No channel rules for channelid(%s) configured", path.c_str());
        ptr_channelid_config->SetError();
        DoUpdateChannelIdConfigRequest(ptr_channelid_config);
        return  ;
    }

    LOG_DEBUG("channel[%s] has child node:%d", path.c_str(), values.size());
    std::string channel_id_str = ptr_channelid_config->m_strChannelId;
    std::size_t max_cnt = values.size();
    ptr_channelid_config->AddNeedAckCount(max_cnt);
    std::string channel_path = ZK_CHANNEL_ID_CONFIG_CHANNEL_PATH(ptr_channelid_config->m_strChannelId);

    for (std::size_t idx = 0; idx < max_cnt; ++idx)
    {
        std::string property_value;
        std::string zk_node_path(channel_path + "/" + values.at(idx));

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

bool Configurator::readChannels(const std::string &channels_path, ChannelIDConfig* ptr_channelid_config, std::string &errorMsg)
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
    event->m_msg = (void*)ptr_channel_id;
    return m_channelIdManager->Update(event);
}

void Configurator::DelayFreeChannelIdConfig(ChannelIDConfig *ptr_channelid_config)
{
    m_channelIdManager->FreeChannelIdPtr(ptr_channelid_config);
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
            LOG_DEBUG("failed to parse channelid(%s) config, config invalid, ptr config:%p", ptr_channelid_config->m_strChannelId.c_str(), ptr_channelid_config);
            DelayFreeChannelIdConfig(ptr_channelid_config);
            return;
        }
    }
}

bool Configurator::CanSendUpdateEventMsg(ChannelIDConfig *ptr_channelid_config)
{

    LOG_DEBUG("check channelid[%s] need request count[%d], recv response count[%d], finish flag[%d], success flag[%d]"
              ,ptr_channelid_config->m_strChannelId.c_str()
              ,ptr_channelid_config->m_zkReqResponseCount.m_requestCount, ptr_channelid_config->m_zkReqResponseCount.m_responseCount
              ,ptr_channelid_config->IsComplete(), ptr_channelid_config->IsSuccess());

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


    //return (all_done && ptr_channelid_config->m_success);
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
    event->m_msg = (void*)ptr_channelid_config;
    m_channelIdManager->Update(event);
    return;
}
int Configurator::updateChannelID(uint32_t channelID, const std::string &context)
{
    std::string errMsg;
    std::string value;

    std::string channel_id_str = gse::tools::strings::ToString(channelID);
    std::string channel_id_path = ZK_CHANNEL_ID_CONFIG_PATH(channel_id_str);
    ChannelIDConfig *ptr_channelid_config = new ChannelIDConfig(channel_id_str);
    std::vector<std::string> list;
    int max_child_node_count = 3;//filter, channel, meta
    ptr_channelid_config->AddNeedAckCount(max_child_node_count);
    //read channel
    if (!readChannels(ZK_CHANNEL_ID_CONFIG_CHANNEL_PATH(channel_id_str), ptr_channelid_config, errMsg))
    {
        LOG_ERROR("it is failed to read channel config for channelid(%s), error info is %s", SAFE_CSTR(channel_id_str.c_str()), SAFE_CSTR(errMsg.c_str()));
        m_channelIdZKConftor->GetConfItemAsync(channel_id_path, Configurator::watchConfCallBack, this, CONFITEMFLAG_CHANNELID_CONFIG, NULL, NULL);
        delete ptr_channelid_config;
        return GSE_ERROR;
    }

    //read meta
    if (!readMetadata(ZK_CHANNEL_ID_CONFIG_METADATA_PATH(channel_id_str), ptr_channelid_config, errMsg))
    {
        LOG_ERROR("it is failed to read metadata config for channelid(%s), error info is %s", SAFE_CSTR(channel_id_str.c_str()), SAFE_CSTR(errMsg.c_str()));
        m_channelIdZKConftor->GetConfItemAsync(channel_id_path, Configurator::watchConfCallBack, this, CONFITEMFLAG_CHANNELID_CONFIG, NULL, NULL);
        return GSE_ERROR;
    }

    //read filter
    if (!readFilters(ZK_CHANNEL_ID_CONFIG_FILTER_PATH(channel_id_str), ptr_channelid_config, errMsg))
    {
        LOG_ERROR("it is failed to read filter config for channelid(%s), error info is %s", SAFE_CSTR(channel_id_str.c_str()), SAFE_CSTR(errMsg.c_str()));
        m_channelIdZKConftor->GetConfItem(channel_id_path, value, Configurator::watchConfCallBack, this, CONFITEMFLAG_CHANNELID_CONFIG);
        return GSE_ERROR;
    }

    m_channelIdZKConftor->GetConfItemAsync(channel_id_path, Configurator::watchConfCallBack, this, CONFITEMFLAG_CHANNELID_CONFIG, NULL, NULL);
    return GSE_SUCCESS;
}

int Configurator::updateDataID(uint32_t dataID, const std::string &context)
{
    DataID *ptrDataid = parseToDataID(context);
    if (NULL == ptrDataid)
    {
        LOG_ERROR("it is failed to parse the dataid  (%u) config, please to check the node (%s) value in the zk", dataID, SAFE_CSTR(context.c_str()));
        return GSE_JSON_INVALID;
    }

    ptrDataid->m_dataId = dataID;
    DataID *ptr_data_id = NULL;
    m_dataIDSV1.Find(dataID, ptr_data_id);
    m_dataIDSV1.Push(ptrDataid->m_dataId, ptrDataid);
    if (NULL != ptr_data_id)
    {
        ptr_data_id->SetNeedDelete();
        m_toDeleteDataIDS.Push(ptr_data_id);
    }

    // NOTIFICATION:
    if (0 != m_callbacks.size())
    {
        std::size_t max_count = m_callbacks.size();
        for (std::size_t idx = 0; idx < max_count; ++idx)
        {
            EventCallbackParams *ptr_callback = m_callbacks.at(idx);
            ptr_callback->m_eventCallbackFunc(ptr_callback->m_ptrCallbackArgs, -1, dataID);
        }
    }
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
}


void Configurator::handleChannelIdChildListChangeEvt(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    // check add a new dataid
    std::size_t max_count = confItem.m_Values.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        std::string channel_id = confItem.m_Values.at(idx);
        LOG_INFO("channelid channged. key[%s], value[%s], child node count[%d]", SAFE_CSTR(confItem.m_Key.c_str()), SAFE_CSTR(channel_id.c_str()), max_count);
        ChannelIDType i_channel_id = gse::tools::strings::StringToUint32(channel_id);
        if (GSE_SUCCESS != pConftor->updateChannelIDConfigFromZK(i_channel_id))
        {
            LOG_ERROR("it is failed to update the channel id(%s)'s config from zookeeper", SAFE_CSTR(channel_id.c_str()));
        }
    }

    pConftor->watchChannelIdConfig(confItem);

}
void Configurator::updateChannelIDConfig(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    std::string value = confItem.m_Values[0];
    LOG_INFO("watch channelid configure channged. key[%s], value[%s]", SAFE_CSTR(confItem.m_Key.c_str()), SAFE_CSTR(value.c_str()));
    if (confItem.m_valueType == CONFITEMVALUE_TYPE_DELETE)
    {
        pConftor->DeleteChannelID(split_channel_id(confItem.m_Key.c_str()));
    }
    else
    {
        pConftor->updateChannelID(split_channel_id(confItem.m_Key.c_str()), value);
    }
}
void Configurator::updateStorage(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    std::size_t max_count = confItem.m_Values.size();
    //LOG_INFO("watch storage config channged. key[%s], value[%s]", SAFE_CSTR(confItem.m_Key.c_str()), SAFE_CSTR(value.c_str()));
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        std::string cluster_id = confItem.m_Values.at(idx);
        std::string cluster_key(confItem.m_Key + "/" + cluster_id);
        // add new storage
        LOG_INFO("read the storage config from zk node :%s", cluster_key.c_str());
        if (GSE_SUCCESS != pConftor->updateStorageConfigFromZK(split_storage_index(cluster_id.c_str()), cluster_key))
        {
            LOG_ERROR("can not update from zk, %s", cluster_key.c_str());
            continue;
        }
    }
}
void Configurator::updateStorageConfig(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    std::string value = confItem.m_Values[0];
    std::size_t pos = confItem.m_Key.find_last_of("/");
    int storage_index = 0;
    if (pos != confItem.m_Key.npos)
    {
        storage_index = split_storage_index(confItem.m_Key.substr(pos).c_str());
    }

    LOG_INFO("read the storage config from zk node :%s", confItem.m_Key.c_str());
    if (GSE_SUCCESS != pConftor->updateStorageConfig(storage_index, value))
    {
        LOG_ERROR("can not update from zk, %s", confItem.m_Key.c_str());
    }
}
void Configurator::updateDataID(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    // check add a new dataid
    std::size_t max_count = confItem.m_Values.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        std::string data_id = confItem.m_Values.at(idx);
        LOG_INFO("dataid channged. key[%s], value[%s]", SAFE_CSTR(confItem.m_Key.c_str()), SAFE_CSTR(data_id.c_str()));
        DataIDType i_data_id = gse::tools::strings::StringToUint32(data_id);
        if (GSE_SUCCESS != pConftor->updateDataIDConfigFromZK(confItem.m_Key + "/" + data_id, data_id))
        {
            LOG_ERROR("it is failed to update the dataid(%s)'s config from zookeeper", SAFE_CSTR(data_id.c_str()));
        }
    }

}
void Configurator::updateDataIDConfig(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    std::string value = confItem.m_Values[0];
    LOG_INFO("watch dataid config channged. key[%s], value[%s]", SAFE_CSTR(confItem.m_Key.c_str()), SAFE_CSTR(value.c_str()));
    pConftor->updateDataID(split_channel_id(confItem.m_Key.c_str()), value);
}
void Configurator::updateBalanceConfig(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    std::string value = confItem.m_Values[0];
    LOG_INFO("watch balance config channged. key[%s], value[%s]", SAFE_CSTR(confItem.m_Key.c_str()), SAFE_CSTR(value.c_str()));
    pConftor->updateDataServerBalanceConfig(value);
}

void Configurator::updateDataIDFromBKData(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    // check add a new dataid
    std::size_t max_count = confItem.m_Values.size();
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        std::string data_id = confItem.m_Values.at(idx);
        LOG_INFO("dataid channged. key[%s], value[%s]", SAFE_CSTR(confItem.m_Key.c_str()), SAFE_CSTR(data_id.c_str()));
        DataIDType i_data_id = gse::tools::strings::StringToUint32(data_id);
        if (GSE_SUCCESS != pConftor->updateDataIDConfigFromBKDataZK(confItem.m_Key + "/" + data_id, data_id))
        {
            LOG_ERROR("it is failed to update the dataid(%s)'s config from zookeeper", SAFE_CSTR(data_id.c_str()));
        }
    }

}

void Configurator::updateDataIDConfigFromBKData(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    updateDataIDConfig(confItem, lpWatcher);
}

void Configurator::updateStorageFromBKData(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    std::size_t max_count = confItem.m_Values.size();
    //LOG_INFO("watch storage config channged. key[%s], value[%s]", SAFE_CSTR(confItem.m_Key.c_str()), SAFE_CSTR(value.c_str()));
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        std::string cluster_id = confItem.m_Values.at(idx);
        std::string cluster_key(confItem.m_Key + "/" + cluster_id);
        // add new storage
        LOG_INFO("read the storage config from zk node :%s", cluster_key.c_str());
        if (GSE_SUCCESS != pConftor->updateStorageConfigFromBKDataZK(split_storage_index(cluster_id.c_str()), cluster_key))
        {
            LOG_ERROR("can not update from zk, %s", cluster_key.c_str());
            continue;
        }
    }
}

void Configurator::watchStreamToChildNodes()
{
    std::vector<std::string> streamto_id_list;
    std::string root_path = ZK_STREAM_ID_CONFIG_BASE_PATH;

    m_channelIdZKConftor->GetChildConfItem(root_path, streamto_id_list, Configurator::watchConfCallBack, this, CONFITEMFLAG_STREAMTO_CONFIG_LIST);
}

void Configurator::handleExporterChildListChangeEvtFromZK(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    std::size_t max_count = confItem.m_Values.size();
    LOG_DEBUG("watch exporter config channged. key[%s], count[%d]", SAFE_CSTR(confItem.m_Key.c_str()), max_count);
    for (std::size_t idx = 0; idx < max_count; ++idx)
    {
        std::string index = confItem.m_Values.at(idx);
        std::string expoter_key(confItem.m_Key + "/" + index);
        LOG_INFO("read the exporter config from zk node :%s", expoter_key.c_str());
        if (GSE_SUCCESS != pConftor->updateStreamToIdConfigFromZK(index))
        {
            LOG_ERROR("failed to update exporter config from zk, %s", expoter_key.c_str());
            continue;
        }
    }
    pConftor->watchStreamToChildNodes();
}


void Configurator::handleExporterValueChangeEvtFromZK(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    LOG_DEBUG("watch exporter config channged, type[%d] key[%s], value[%s]", confItem.m_valueType, SAFE_CSTR(confItem.m_Key.c_str()), SAFE_CSTR(confItem.m_Values[0].c_str()));
    int stream_to_id = split_channel_id(confItem.m_Key.c_str());
    std::string str_stream_to_id = gse::tools::strings::ToString(stream_to_id);
    if (confItem.m_valueType == CONFITEMVALUE_TYPE_DELETE)
    {
        LOG_DEBUG("delete stream to id[%s]", str_stream_to_id.c_str());
        pConftor->deleteStreamToIdConfigFromZK(str_stream_to_id);
        return;
    }


    LOG_INFO("read the exporter config from zk node :%s, index:%s", confItem.m_Key.c_str(), str_stream_to_id.c_str());
    if (GSE_SUCCESS != pConftor->updateStreamToIdConfigFromZK(str_stream_to_id))
    {
        LOG_ERROR("failed to update exporter config from zk, %s", confItem.m_Key.c_str());
        return;
    }

    return;
}

void Configurator::updateStorageConfigFromBKData(WatchConfItem &confItem, void *lpWatcher)
{
    Configurator *pConftor = (Configurator *)lpWatcher;
    updateStorageConfig(confItem, lpWatcher);
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

}
}
