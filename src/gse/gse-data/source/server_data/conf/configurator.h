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

#ifndef _GSE_DATA_CONFIGURATOR_H_
#define _GSE_DATA_CONFIGURATOR_H_

#include <string.h>
#include <string>

#include "dataStruct/safe_map.h"
#include "dataStruct/safe_queue.h"
#include "tools/net.h"
#include "tools/strings.h"
#include "conftor/conftor.h"
#include "dataconf.h"
#include "tools/thread.h"
#include "eventthread/gseEventThread.h"

#include "conf/channel_id_config.h"
#include "conf/bkdata_config_v1.h"

namespace gse { 
namespace dataserver {

enum ConfItemFlag
{
    CONFITEMFLAG_UNSET = -1,
    CONFITEMFLAG_DATACONF = 0,
    CONFITEMFLAG_CHANNELID_LIST = 1,
    CONFITEMFLAG_CHANNELID_CONFIG = 2,
    CONFITEMFLAG_CHANNELID_FILTER_CONFIG = 3,
    CONFITEMFLAG_CHANNELID_EXPORTER_CONFIG = 4,
    CONFITEMFLAG_CHANNELID_CHANNEL_CONFIG = 5,
    CONFITEMFLAG_STORAGE = 6,
    CONFITEMFLAG_STORAGE_CONFIG = 7,
    CONFITEMFLAG_DATAID = 8,
    CONFITEMFLAG_DATAID_CONFIG = 9,
    CONFITEMFLAG_BALANCE_CONFIG = 10,

    CONFITEMFLAG_SCHEDULE_CONFIG = 11,
    CONFITEMFLAG_DATAID_FROM_BKDATA = 12,
    CONFITEMFLAG_DATAID_CONFIG_FROM_BKDATA = 13,
    CONFITEMFLAG_STORAGE_FROM_BKDATA = 14,
    CONFITEMFLAG_STORAGE_CONFIG_FROM_BKDATA = 15,

    CONFITEMFLAG_STREAMTO_CONFIG_LIST = 16,
    CONFITEMFLAG_STREAMTO_CONFIG_VALUE = 17,

    CONFITEMFLAG_TGLOG_CHANNEL_ID_VALUE = 18,

    CONFITEMFLAG_MAX_NUM
};

inline int split_storage_index(const char *target)
{

    if (NULL == target)
    {
        return 0;
    }
    std::string data(target);

    char *ptr = strtok(const_cast<char *>(data.c_str()), "_");
    while (NULL != ptr)
    {
        int val = atoi(ptr);
        if (val > 0)
        {
            return val;
        }
        ptr = strtok(NULL, "_");
    }
    return 0;
}

inline uint32_t split_channel_id(const char *target)
{
    if (NULL == target)
    {
        return 0;
    }
    std::string data(target);

    std::size_t pos = data.find_last_of("/");
    if (pos != data.npos)
    {
        return gse::tools::strings::StringToUint32(data.substr(pos + 1));
    }
    return 0;
}

typedef struct _EventCallbackParams{
    IDToStorage::WatchEventFunc  m_eventCallbackFunc;
    void*           m_ptrCallbackArgs;
}EventCallbackParams;



class Configurator : public IDToStorage
{
public:
    typedef void (*PTRConfigWatchFunc)(WatchConfItem &confItem);


public:
    Configurator();
    virtual ~Configurator();

public:
    // overloaded
    virtual DataID *GetStorageByDataID(uint32_t dataID);
    virtual DataStorage *GetAllStorages();
    virtual DataStorage *GetStorageByIndex(int storageIndex);
    virtual DataStorage *GetStorageByChannelID(uint32_t channelID);
    virtual void WatchUpdateEvent(WatchEventFunc callback, void *args);

public:
    // init function

    // Init init self from remote config center
    int Init(const std::string &dataFlowConf, const std::string &confHost, const std::string &password, const std::string &bkdataZK, const string &channelid_zkhost, const string &channelid_zkauth, const std::string &selfIp, const std::string &regionID, const std::string &cityID, const std::string &clusterName, const std::string &instanceId, const std::string &watchpath);

public:
    // static functions
    static void reportBalance(evutil_socket_t fd, short what, void *args);
    static void cleanDataID(evutil_socket_t fd, short what, void *args);
    static uint32_t getDefaultTglogChannelId();

public:
    // common functions
    void GetAllChannelID(std::vector<std::string> &channelid_list);
    void GetAllDataID(std::vector<std::string> &dataid_list);
    void GetAllV1StorageID(std::vector<std::string> &storage_id_list);

    ChannelIdManager *GetChannelIdManager();
    DataFlowConf *GetDataFlowConf();
public:
    void GetChannelsConfigValue(std::string &path, int rc, std::vector<std::string> &values, ChannelIDConfig *ptr_channelid_config);

    void GetChannelConfigValueCallBack(std::string &path, int rc, std::string &property_value, void* ptr_callbackobj);
    void ParseMetaConfig(std::string &path, int rc, std::string &property_value, ChannelIDConfig* ptr_channelid_config);
    void GetFiltersValue(std::string &path, int rc, std::vector<std::string> &values, ChannelIDConfig * ptr_channelid_config);
    void UpdateFilterValue(std::string &path, int rc, std::string &property_value, void* ptr_callbackobj);
private:
    // InitFromLocalConfig init self from local config
    // Maybe this function will be replaced by another way
    int initFromLocalConfig(const std::string &dataFlowConf);

    static void watchConfCallBack(WatchConfItem &confItem, void *lpWatcher);

    static void channelNodeChangeCallBack(WatchConfItem &confItem, void *lpWatcher);
   
    static void updateDataConf(WatchConfItem &confItem, void *lpWatcher);
    static void handleChannelIdChildListChangeEvt(WatchConfItem &confItem, void *lpWatcher);
    static void updateChannelIDConfig(WatchConfItem &confItem, void *lpWatcher);
    static void updateStorage(WatchConfItem &confItem, void *lpWatcher);
    static void updateStorageConfig(WatchConfItem &confItem, void *lpWatcher);
    static void updateDataID(WatchConfItem &confItem, void *lpWatcher);
    static void updateDataIDConfig(WatchConfItem &confItem, void *lpWatcher);
    static void updateBalanceConfig(WatchConfItem &confItem, void *lpWatcher);
    static void updateDataIDFromBKData(WatchConfItem &confItem, void *lpWatcher);
    static void updateDataIDConfigFromBKData(WatchConfItem &confItem, void *lpWatcher);
    static void updateStorageFromBKData(WatchConfItem &confItem, void *lpWatcher);
    static void updateStorageConfigFromBKData(WatchConfItem &confItem, void *lpWatcher);
    static void handleExporterChildListChangeEvtFromZK(WatchConfItem &confItem, void *lpWatcher);
    static void handleExporterValueChangeEvtFromZK(WatchConfItem &confItem, void *lpWatcher);
    static void handleTglogChannelIdChangeEvtFromZK(WatchConfItem &confItem, void *lpWatcher);

    static void GetChannelsListResultCallBack(std::string &path, int rc, std::vector<std::string> &values, const void* ptr_callbackobj);
    static void GetFiltersResultCallBack(std::string &path, int rc, std::vector<std::string> &values, const void* ptr_callbackobj);

    static void ChannelValueResultCallBack(std::string &path, int rc, const char *value, int32_t value_len, const struct Stat *stat, const void* ptr_callbackobj);
    static void GetMetaValueCallBack(std::string &path, int rc, const char *value, int32_t value_len, const struct Stat *stat, const void *ptr_callbackobj);
    static void FilterValueResultCallBack(std::string &path, int rc, const char *value, int32_t value_len, const struct Stat *stat, const void* ptr_callbackobj);

private:
    // update config from zk

    int watchConfigsFromZK();
    int watchDataFlow();
    int watchStorage();
    int watchExporter();
    // watchStorageFromBKData 仅用于兼容存量bkdata 里关于dataid 的配置，未来需要被废弃
    int watchStorageFromBKData();
    int watchChannelID();
    int watchStreamToID();
    int WatchTglogChannelId();
    int watchDataID();
    // watchDataIDFromBKData 仅用于兼容存量bkdata 里关于dataid 的配置，未来需要被废弃
    int watchDataIDFromBKData();

    int watchBalanceConfig();
    void watchChannelIdConfig(WatchConfItem &confItem);
    int updateStorageConfigFromZK(StorageIndex storageIndex, const std::string &nodePath);
    int updateStorageConfigFromBKDataZK(StorageIndex storageIndex, const std::string &nodePath);
    int updateDataIDConfigFromZK(const std::string &nodePath, const std::string &dataID);
    int updateDataIDConfigFromBKDataZK(const std::string &nodePath, const std::string &dataID);
    int updateChannelIDConfigFromZK(uint32_t channelID);

    int updateStreamToIdConfigFromZK(std::string &stream_to_id);
    int deleteStreamToIdConfigFromZK(std::string &stream_to_id);
    int updateStorageConfig();
    int updateDataFlowConf(const std::string &context);
    int updateStorageConfig(StorageIndex storageIndex, const std::string &context);
    int updateChannelID(uint32_t channelId, const std::string &context);
    int DeleteChannelID(uint32_t channelId);

    // updateDataID compatible with existing configurations about dataid (now it is channel id)
    int updateDataID(uint32_t dataID, const std::string &context);
    
    int updateDataServerBalanceConfig(const std::string &context);
    int updateDataServerLoadBalanceNode();

    void updateLocationFromZK();

    void watchStreamToChildNodes();

    //private ,support async zk api get channelid
    bool CanSendUpdateEventMsg(ChannelIDConfig *ptr_channelid_config);
    void SendUpdateChannelIdConfigEventMsg(ChannelIDConfig *ptr_channelid_config);
    bool IsFinishAndFailed(ChannelIDConfig *ptr_channelid_config);
    void DoUpdateChannelIdConfigRequest(ChannelIDConfig *ptr_channelid_config);
    void DelayFreeChannelIdConfig(ChannelIDConfig *ptr_channelid_config);
private:
    // read channelid config
    bool readMetadata(const std::string& metadataPath, ChannelIDConfig *ptr_channelid_config, std::string& errorMsg);
    bool readFilters(const std::string& fileterPath, ChannelIDConfig *ptr_channelid_config, std::string& errorMsg);
    bool readExporters(const std::string& exporterPath, std::vector<StreamConfig*>& exporters, std::string& errorMsg);
    bool readChannels(const std::string& channelPath, ChannelIDConfig *ptr_channelid_config, std::string& errorMsg );

private:
    // set config into zk

    int createBaseConfItem();

private:
    int startConftor();
    void stopConftor();

private:
    Conftor *m_conftor;
    DataConf *m_dataConf;
    std::string m_confHost;
    std::string m_confHostPassword;
    std::string m_selfIp;
    std::string m_clusterName;
    std::string m_instanceId;
    std::string m_localDataFlowConfig;


    std::string m_channelidZkHost;
    std::string m_channelidZkAuth;
    std::string m_watchPath;

private:
    // 兼容 V1 版本DS 的配置
    Json::Value m_cfgLoadBalanceValue;
    StorageConfigMap m_storagesV1;
    DataIDConfigMap m_dataIDSV1;
    gse::datastruct::SafeQueue<DataID*>  m_toDeleteDataIDS;
    gse::datastruct::SafeQueue<ChannelIDConfig*> m_toDeleteChannelIDS;
    gse::datastruct::SafeQueue<ChannelIdExporterConfig*> m_toDeleteStreamToIDS;

    gse::datastruct::SafeMap<ChannelIDType, ChannelIDConfig *> m_channelIDS;
    gse::datastruct::SafeMap<StreamToIDType, ChannelIdExporterConfig *> m_streamToIdConfig;

    std::string m_regionID;
    std::string m_cityID;
    rgse::GseEventThread *m_ptrEventThread;

private:
    // V1 版本DS负载均衡的参数配置
    float m_cpuUsage;
    float m_hostPerformance;
    float m_hostLoadBance;
    float m_loadweight;
    float m_dloadweight;

    BaseCfg m_baseCfg;
    std::string m_bkdataZK;
    Conftor *m_bkdataZKConftor;

    Conftor *m_channelIdZKConftor;
    ChannelIdManager *m_channelIdManager;
    ChannelIdExporterManager *m_channelIdExporterManager;
private:
    // factory
    ChannelIDConfigFactory m_channelIDConfigFactory;
    void (*m_ptrConfigFunc[CONFITEMFLAG_MAX_NUM])(WatchConfItem &confItem, void *lpWatcher);
    std::vector<EventCallbackParams*> m_callbacks;
};



struct ChannelIdZkCallBackObj
{
    Configurator *m_ptrConfigurator;
    ChannelIDConfig *m_ptrChannelidConfig;

    ChannelIdZkCallBackObj()
    {
        m_ptrConfigurator = NULL;
        m_ptrChannelidConfig = NULL;
    }
};

}
}
#endif //_GSE_DATA_CONFIGURATOR_H_
