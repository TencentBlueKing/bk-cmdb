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

#ifndef _GSE_DATA_CHANNELID_ZK_API_H_
#define _GSE_DATA_CHANNELID_ZK_API_H_

#include <string>

#include "channelid_def.h"
#include "discover/zkapi/zk_api.h"

#include "channelid_api.h"

namespace gse {
namespace data {

#include <memory>

class ChannelIdZkAPI;
class DataIdZkAPI;

static std::shared_ptr<ChannelIdZkAPI> kChannelIdZkInstance;
static std::shared_ptr<DataIdZkAPI> kDataIdZkInstance;

std::shared_ptr<ChannelIdZkAPI> ChannelIdZkApiInst();

class ChannelIdZkAPI
{
public:
    bool Init();
    ChannelIdZkAPI(std::shared_ptr<gse::discover::zkapi::ZkApi> zkClient, bool acl);

public:
    // API

    bool ReadPlatNumber(const std::string platName, std::string& platid, std::string& errorMsg);

    std::string GenerateOriginChanelId(const std::string& platName, std::string& errorMsg);
    std::string GetSequeceNumber(std::string& buff);

    int GetPlatId(const std::string& platName, std::string& errorMsg);
    int GenerateChannelId(const std::string& platName, ApiError& error);
    uint32_t GetOriginChannelID(uint32_t channelid);
    bool ExistChannelId(const std::string& channelId, ApiError& error);

    // channelid api
    bool CreateChannelIdPath(const std::string& strChannelId, ApiError& error);

    bool CreateMetaData(const std::string& strChannelId, const std::shared_ptr<Metadata> metaData, ApiError& error, bool create = true);
    bool CreateFilters(const std::string& strChannelId, const std::vector<std::shared_ptr<StreamFilter>>& filters, ApiError& error, bool create = true);
    bool CreateChannels(const std::string& strChannelId, const std::vector<std::shared_ptr<Channel>> m_channels, ApiError& error, bool create = true);

    bool UpdateChannelIdTime(const std::string& strChannelId, ApiError& error);

    bool CreateChannelIdNotify(const std::string& strChannelId, const string& value, std::string& errorMsg, bool create = true);
    bool CreateChannelIdQueryIndex(std::shared_ptr<Metadata> ptrMetadata, const std::string& streamToId, ApiError& error);

    bool GetTglogChannelId(const std::shared_ptr<Metadata> ptrMetadata, uint32_t& realChannelId, ApiError& error);
    bool SaveTlogChannelId(const std::shared_ptr<Metadata> ptrMetadata, ApiError& error);
    bool DeleteTglogChannelId(int bizid, std::string& odm);

    bool DeleteOriginChannelID(const std::string& originChannelID, ApiError& error);
    bool DeleteChannelId(int64_t channelId, ApiError& error);
    bool DeleteBySpecification(int64_t channelId, std::shared_ptr<DeleteSpecification> deleteSpecification, ApiError& error);

    bool DeleteChannelIdBizIdIndex(const std::string& bizId, const std::string& channelId);
    bool DeleteChannelIdPlatNameIndex(const std::string& platName, const std::string& channelId);
    bool DeleteChannelIdOdmNameIndex(const std::string& label, const std::string& channelId);
    bool DeleteChannelIdStreamToIndex(const std::string& streamToId, const std::string& channelId);
    // read from zk
    bool ReadFilters(const std::string& channelId, std::vector<std::shared_ptr<StreamFilter>>& filters, ApiError& error);
    bool ReadMetadata(const std::string& channelId, std::shared_ptr<Metadata> metadata, ApiError& error);
    bool ReadChannels(const std::string& channelId, std::vector<std::shared_ptr<Channel>>& channels, ApiError& error);

    std::shared_ptr<ChannelIdConfig> ReadChannelIdConfig(const std::string& channelId, const std::string& platName, ApiError& error);
    // streamTo api

    bool ExistStreamToClusterId(const std::string& streamToId);

    int64_t GenerateStreamToClusterId(const std::string& platName, ApiError& error);
    int64_t GenerateOriginStreamToClusterId(const std::string& platName, std::string& errorMsg);

    bool CreateStreamToIdQueryIndex(std::shared_ptr<StreamToMetadata> ptr_metadata, const std::string& type, ApiError& error);

    // streamTo Api
    bool CreateStreamToClusterNode(int64_t streamToClusterId, std::shared_ptr<StreamToCluster> ptrStreamToCluster, ApiError& error);
    bool CreateStreamToClusterMetadataNode(int64_t streamToClusterId, std::shared_ptr<StreamToMetadata> ptrMeta, ApiError& error);

    bool UpdateStreamToClusterConfig(int64_t streamToClusterId, std::shared_ptr<StreamToCluster> ptrStreamToCluster, ApiError& error);
    bool UpdateStreamToClusterMeta(int64_t streamToClusterId, std::shared_ptr<StreamToMetadata> ptrMeta, ApiError& error);

    bool ReadStreamToIdZkdata(const std::string& streamToId, std::shared_ptr<StreamToCluster> ptrStreamToConfig, ApiError& error);
    bool ReadStreamToMetadata(const std::string& streamToId, std::shared_ptr<StreamToMetadata> metadata, ApiError& error);
    std::shared_ptr<StreamToIdConfig> ReadStreamToIdConfig(const std::string& streamToId, ApiError& error);

    bool DeleteStreamToID(const std::string& streamToId, ApiError& error);
    bool DeleteStreamToPlatNameIndex(const std::string& platName, const std::string& streamToId, ApiError& error);
    bool DeleteStreamToBkBizIdIndex(const std::string& bkBizid, const std::string& streamToId, ApiError& error);
    bool DeleteClusterTypeIndex(const std::string& clusterType, const std::string& strStreamToId, ApiError& error);

    bool DeleteOriginStreamToID(const std::string& strStreamToId, ApiError& error);
    std::shared_ptr<gse::discover::zkapi::ZkApi> GetZkApiClient();

private:
    bool CreateNode(const std::string& nodePath, std::string& value, ApiError& error);
    bool DeleteRecursive(const std::string& path, int32_t version /*= -1*/, ApiError& error);
    bool CreateBaseNode();

private:
    std::shared_ptr<gse::discover::zkapi::ZkApi> m_ptrEventZkClient;
    std::string m_zkAuth;

    // std::vector<std::shared_ptr<DataIdZkClient>>* m_dataidZkClient;
    bool m_acl;

    std::map<std::string, uint32_t> m_supportPlats = {
        {SUPPORT_PLAT_NAME_TGDP, PLAT_BKMONITOR_ID},
        {SUPPORT_PLAT_NAME_TGLOG, PLAT_TGLOG_ID},
        {SUPPORT_PLAT_NAME_TDM, PLAT_TDM_ID},
        {SUPPORT_PLAT_NAME_BKMONITOR, PLAT_BKMONITOR_ID},
        {SUPPORT_PLAT_NAME_GSE, PLAT_GSE_ID},
        {SUPPORT_PLAT_NAME_IEG_DATA, PLAT_IEG_DATA_ID},
        {SUPPORT_PLAT_NAME_DATAMORE, PLAT_DATAMORE_ID},
        {SUPPORT_PLAT_NAME_CMDB, PLAT_CMDB_ID},
    };
};

class DataIdZkAPI
{

public:
    bool CreateDataIdConfig();

private:
    std::shared_ptr<gse::discover::zkapi::ZkApi> m_ptrEventZkClient;
    std::string m_zkAuth;
};

bool SaveToDataIdZk(const std::string& dataId, const std::string& value, ApiError& error);
bool DeleteDataIdZk(const std::string& dataId, ApiError& error);
bool UpdateDataIdZk(const std::string& dataId, const std::string& value, ApiError& error);

bool CreateClusterIdToZk(const std::string& clusterId, const std::string& value, ApiError& error);
bool DeleteClusterIdFromZk(const std::string& clusterId, ApiError& error);
bool UpdateClusterIdToZk(const std::string& dataId, const std::string& value, ApiError& error);

class DataIdZkClient
{
public:
    bool Init();
    bool DeleteDataId(const std::string& dataId, ApiError& error);
    bool SaveDataId(const std::string& dataId, const string& value, ApiError& error);
    bool CreateNode(const std::string& nodePath, std::string& value, ApiError& error);
    bool UpdateDataId(const std::string& dataId, const std::string& value, ApiError& error);

    bool CreateClusterIdConfig(const std::string& clusterIndex, const std::string& value, ApiError& error);
    bool UpdateClusterIdConfig(const std::string& clusterIndex, const std::string& value, ApiError& error);
    bool DeleteClusterIdConfig(const std::string& clusterIndex, ApiError& error);

public:
    std::shared_ptr<gse::discover::zkapi::ZkApi> m_ptrEventZkClient;
    std::string m_zkAuth;
    std::string m_dataIdZkPath;
    std::string m_storageZkPath;

    bool m_acl;
};

std::vector<std::shared_ptr<DataIdZkClient>> DataIdZkApiInst();

} // namespace data
} // namespace gse

#endif // CHANNELID_ZK_API_H
