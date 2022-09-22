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
#include "channelid_zk_api.h"

#include <stdio.h>
#include <string.h>
#include <time.h>

#include <iomanip>  // std::setfill, std::setw
#include <iostream> // std::cout, std::endl
#include <string>

#include "channelid_def.h"
#include "dataserver.h"
#include "error_code.h"
#include "log/log.h"
#include "tools/strings.h"

namespace gse {
namespace data {

std::shared_ptr<ChannelIdZkAPI> ChannelIdZkApiInst()
{
    return DataServer::GetChannelIdZkApiInst();
}

std::vector<std::shared_ptr<DataIdZkClient>> DataIdZkApiInst()
{
    return DataServer::GetDataIdZkApiInst();
}
ChannelIdZkAPI::ChannelIdZkAPI(std::shared_ptr<gse::discover::zkapi::ZkApi> zkClient, bool acl)
    : m_acl(acl)
{
    m_ptrEventZkClient = zkClient;
}

bool ChannelIdZkAPI::Init()
{
    // create zk base path
    bool ret = CreateBaseNode();
    return ret;
}

bool ChannelIdZkAPI::CreateBaseNode()
{
    std::string buff;
    ApiError error;
    std::string zkChannelIDOriginPath(ZK_CHANNEL_ID_ORIGIN_PATH());
    if (!CreateNode(zkChannelIDOriginPath, buff, error))
    {
        LOG_ERROR("failed to create the node path (%s), error info is %s",
                  zkChannelIDOriginPath.c_str(), GET_ERROR_MESSAGE(error).c_str());
        return false;
    }

    std::string zkStreamToIDOriginPath(ZK_STREAM_TO_ID_ORIGIN_PATH());
    if (!CreateNode(zkStreamToIDOriginPath, buff, error))
    {
        LOG_ERROR("it is failed to create the node path (%s), error info is %s",
                  zkStreamToIDOriginPath.c_str(), GET_ERROR_MESSAGE(error).c_str());
        return false;
    }

    std::string zkStreamToIDPath(ZK_STREAM_ID_CONFIG_BASE_PATH);
    if (!CreateNode(zkStreamToIDPath, buff, error))
    {
        LOG_ERROR("it is failed to create the node path (%s), error info is %s",
                  zkStreamToIDPath.c_str(), GET_ERROR_MESSAGE(error).c_str());
        return false;
    }

    std::string zkPlatIDPath(ZK_PLAT_ID_CONFIG_BASE_PATH);
    if (!CreateNode(zkPlatIDPath, buff, error))
    {
        LOG_ERROR("it is failed to create the node path (%s), error info is %s",
                  zkPlatIDPath.c_str(), GET_ERROR_MESSAGE(error).c_str());
        return false;
    }

    std::string zkChannelIDPath(ZK_CHANNEL_ID_CONFIG_BASE_PATH);
    if (!CreateNode(zkChannelIDPath, buff, error))
    {
        LOG_ERROR("it is failed to create the node path (%s), error info is %s",
                  zkChannelIDPath.c_str(), GET_ERROR_MESSAGE(error).c_str());
        return false;
    }
    return true;
}

bool ChannelIdZkAPI::ReadPlatNumber(const std::string platName, std::string& platid, std::string& errorMsg)
{
    if (NULL == m_ptrEventZkClient)
    {
        errorMsg = "the zookeeper is not connected";
        LOG_ERROR("it is failed to read the value of the plat name (%s), because of the zookeeper server is not connected", platName.c_str());
        return false;
    }

    std::string platnameNodePath = std::string(ZK_CHANNEL_ID_PLAT_PATH_WITH_NAME) + "/" + platName;

    int ret = GSE_SUCCESS;
    ret = m_ptrEventZkClient->ZkGet(platnameNodePath, platid, NULL, NULL, NULL);
    if (ret != GSE_SUCCESS)
    {
        errorMsg = " the plat '" + platName + "' is not registered";
        LOG_ERROR("failed to get platid, zknode:%s", platnameNodePath.c_str());
        return false;
    }
    if (gse::tools::strings::IsNumber(platid))
    {
        return true;
    }
    errorMsg.assign("platid is not a number");
    LOG_ERROR("platid is not a number:%s", platid.c_str());
    return false;
}

bool ChannelIdZkAPI::CreateNode(const std::string& nodePath, std::string& value, ApiError& error)
{
    if (NULL == m_ptrEventZkClient)
    {
        error = ZooKeeperUnit;
        LOG_ERROR("it is failed to create the node(%s), because of the zookeeper server is not connected", nodePath.c_str());
        return false;
    }

    std::string buff;
    std::vector<std::string> nodes;
    ChannelIdCommonAPI::SplitZkPath(nodePath, nodes);
    std::string path;
    for (std::size_t size = nodes.size(), idx = 0; idx != size; ++idx)
    {
        if (0 != idx)
        {
            path.append("/");
        }
        path.append(nodes.at(idx));
        if (!m_ptrEventZkClient->ZkExists(path, NULL, NULL, NULL))
        {
            int ret = 0;
            if (idx == (size - 1))
            {
                LOG_DEBUG("create leaf node:%s, value:%s", path.c_str(), value.c_str());
                ret = m_ptrEventZkClient->ZkCreateNormal(path, value, buff, m_acl);
            }
            else
            {
                ret = m_ptrEventZkClient->ZkCreateNormal(path, " ", buff, m_acl);
            }

            if (ret != GSE_SUCCESS)
            {
                error = ZooKeeperCreateError;
                LOG_ERROR("failed to create the zk node[%s], ret:%d ", path.c_str(), ret);
                return false;
            }
        }
    }

    return true;
}

std::string ChannelIdZkAPI::GetSequeceNumber(std::string& buff)
{
    std::size_t len = strlen(buff.c_str());
    buff.resize(len);
    std::size_t index = buff.rfind("/");
    std::string subPath = buff.substr(index + 1);

    return subPath;
}

std::string ChannelIdZkAPI::GenerateOriginChanelId(const std::string& platName, std::string& errorMsg)
{
    std::string buff;
    std::string channelidValue;

    std::string zkChannelIDOriginPath(ZK_CHANNEL_ID_ORIGIN_PATH());

    std::string sequenceNode = zkChannelIDOriginPath + std::string("/");
    int ret = m_ptrEventZkClient->ZkCreateSequence(sequenceNode, platName, buff);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("failed to create sequence node, node:%s, zkerror:%s", sequenceNode.c_str(), m_ptrEventZkClient->ZkError(ret));
        return "";
    }

    channelidValue = GetSequeceNumber(buff);
    LOG_DEBUG("generate an origin channelid (%s)", channelidValue.c_str());
    return channelidValue;
}

int ChannelIdZkAPI::GetPlatId(const std::string& platName, std::string& errorMsg)
{
    auto it = m_supportPlats.find(platName);
    if (it != m_supportPlats.end())
    {
        return it->second;
    }
    return -1;
}

bool ChannelIdZkAPI::ExistChannelId(const std::string& channelId, ApiError& error)
{
    if (NULL == m_ptrEventZkClient)
    {
        LOG_ERROR("it is failed to check channelid exist(%s), because of the zookeeper server is not connected", channelId.c_str());
        error = ZooKeeperUnit;
        return false;
    }

    std::string channelIdPath = ZK_CHANNEL_ID_CONFIG_PATH(channelId);
    bool ret = m_ptrEventZkClient->ZkExists(channelIdPath, NULL, NULL, NULL);
    if (!ret)
    {
        return false;
    }

    return true;
}

bool ChannelIdZkAPI::ExistStreamToClusterId(const std::string& streamToId)
{
    if (NULL == m_ptrEventZkClient)
    {
        LOG_ERROR("it is failed to check streamtoid exist(%s), because of the zookeeper server is not connected", streamToId.c_str());
        return false;
    }

    std::string streamToIdPath = ZK_STREAM_ID_CONFIG_PATH(streamToId);
    bool ret = m_ptrEventZkClient->ZkExists(streamToIdPath, NULL, NULL, NULL);
    if (!ret)
    {
        return false;
    }

    return true;
}

int ChannelIdZkAPI::GenerateChannelId(const std::string& platName, ApiError& error)
{
    int channelId = 0;
    int retries = 0;
    // TODO:
    const int maxReties = 100;
    std::string errorMsg;

    do
    {
        std::string originPlatChannelId = GenerateOriginChanelId(platName, errorMsg);
        if (!errorMsg.empty())
        {
            return -1;
        }

        int platId = 0;
        if ((platId = GetPlatId(platName, errorMsg)) < 0)
        {
            LOG_ERROR("invalid plat name:%s", platName.c_str());
            errorMsg.assign("Unsupported platform name");
            return -1;
        }

        if (gse::tools::strings::StringToInt64(originPlatChannelId) >= MAX_CHANNEL_ID_INDEX)
        {
            LOG_ERROR("alloc channelid failed, byond max %d", MAX_CHANNEL_ID_INDEX);
            char errmsg[256] = {0};
            snprintf(errmsg, sizeof(errmsg), "failed to alloc channelid, beyond the range:[1 - %d]", MAX_CHANNEL_ID_INDEX);
            errorMsg.assign(errmsg);
            return -1;
        }

        channelId = MakeUnifyChannelID(platId, gse::tools::strings::StringToInt32(originPlatChannelId));
        std::string strChannelId = gse::tools::strings::ToString(channelId);
        if (ExistChannelId(strChannelId, error))
        {
            LOG_ERROR("channelid(%d) has alloced, realloc ", channelId);
            retries++;
        }
        else
        {
            return channelId;
        }

        if (retries >= maxReties)
        {
            LOG_ERROR("alloc channelid failed,the maximum retry times is exceeded(%d)", maxReties);
            errorMsg.assign("failed to alloc channelid");
            return -1;
        }

    } while (true);

    return -1;
}

uint32_t ChannelIdZkAPI::GetOriginChannelID(uint32_t channelid)
{
    uint32_t realId = GetRealChannelID(channelid);
    uint32_t originId = ((realId >= ZK_CHANNEL_ID_OFFSET) ? (realId - ZK_CHANNEL_ID_OFFSET) : realId);
    return originId;
}

int64_t ChannelIdZkAPI::GenerateOriginStreamToClusterId(const std::string& platName, std::string& errorMsg)
{
    std::string buff;
    std::string streamToClusterId;

    std::string zkStreamToIDOriginPath(ZK_STREAM_TO_ID_ORIGIN_PATH());

    int ret = m_ptrEventZkClient->ZkCreateSequence(zkStreamToIDOriginPath + std::string("/"), platName, buff);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("it is failed to create the node");
        return -1;
    }

    streamToClusterId = GetSequeceNumber(buff);
    LOG_DEBUG("generate a new origin stream to id (%s)", streamToClusterId.c_str());
    return gse::tools::strings::StringToInt64(streamToClusterId);
}

int64_t ChannelIdZkAPI::GenerateStreamToClusterId(const std::string& platName, ApiError& error)
{
    std::string errorMsg;
    int64_t originStreamToId = GenerateOriginStreamToClusterId(platName, errorMsg);
    if (!errorMsg.empty())
    {
        return -1;
    }

    if (originStreamToId < 0)
    {
        return -1;
    }

    return MakeUnifyStreamToID(originStreamToId);
}

bool ChannelIdZkAPI::CreateChannelIdPath(const std::string& strChannelId, ApiError& error)
{
    std::string buff;
    std::string nowTime = gse::tools::time::GetUTCTimeString();
    std::string strChannelIdPath = ZK_CHANNEL_ID_CONFIG_PATH(strChannelId);
    int ret = m_ptrEventZkClient->ZkCreateNormal(strChannelIdPath, nowTime, buff, m_acl);
    if (ret != GSE_SUCCESS && ret != GSE_ZK_NODE_EXIST)
    {
        error = ZooKeeperCreateError;
        LOG_ERROR("it is failed to create channel path:%s", strChannelIdPath.c_str());
        return false;
    }

    return true;
}

bool ChannelIdZkAPI::CreateMetaData(const std::string& strChannelId, const std::shared_ptr<Metadata> metaData, ApiError& error, bool create)
{
    std::string buff;
    if (create)
    {
        int ret = m_ptrEventZkClient->ZkCreateNormal(ZK_CHANNEL_ID_CONFIG_METADATA_PATH(strChannelId), metaData->ToJSON(), buff, m_acl);
        if (ret != GSE_SUCCESS)
        {
            return false;
        }
    }
    else
    {
        int ret = m_ptrEventZkClient->ZkSet(ZK_CHANNEL_ID_CONFIG_METADATA_PATH(strChannelId), metaData->ToJSON(), -1, NULL);

        if (ret != GSE_SUCCESS)
        {
            error = ZooKeeperSetError;
            LOG_ERROR("it is failed to set the channel id(%s)'s metadata", strChannelId.c_str());
            return false;
        }
    }

    return true;
}

bool ChannelIdZkAPI::CreateFilters(const std::string& strChannelId, const std::vector<std::shared_ptr<StreamFilter>>& filters, ApiError& error, bool create)
{
    std::string buff;

    int filterCount = filters.size();
    if (create)
    {

        int ret = m_ptrEventZkClient->ZkCreateNormal(ZK_CHANNEL_ID_CONFIG_FILTER_PATH(strChannelId), gse::tools::strings::ToString(filterCount), buff, m_acl);
        if (ret != GSE_SUCCESS && ret != GSE_ZK_NODE_EXIST)
        {
            error = ZooKeeperCreateError;
            LOG_ERROR("it is failed to create the node path(%s)",
                      ZK_CHANNEL_ID_CONFIG_METADATA_PATH(strChannelId).c_str());
            return false;
        }
    }

    for (auto filter : filters)
    {
        std::string filterPath = ZK_CHANNEL_ID_CONFIG_FILTER_PATH(strChannelId) + "/" + filter->m_name;

        if (create)
        {
            int ret = m_ptrEventZkClient->ZkCreateNormal(filterPath, filter->ToJSON(), buff, m_acl);
            if (ret != GSE_SUCCESS)
            {
                error = ZooKeeperCreateError;
                LOG_ERROR("failed to create filter node, path:%s, %s", filterPath.c_str(), GET_ERROR_MESSAGE(error).c_str());
                return false;
            }

            LOG_DEBUG("successfully create filter, path:%s, value:%s", filterPath.c_str(), filter->ToJSON().c_str());
        }
        else
        {
            int ret = m_ptrEventZkClient->ZkSet(filterPath, filter->ToJSON(), -1, NULL);
            if (ret != GSE_SUCCESS)
            {
                int ret = m_ptrEventZkClient->ZkCreateNormal(filterPath, filter->ToJSON(), buff, m_acl);
                if (ret != GSE_SUCCESS)
                {
                    error = ZooKeeperSetError;
                    LOG_ERROR("failed to create filter node, path:%s, %s", filterPath.c_str(), GET_ERROR_MESSAGE(error).c_str());
                    return false;
                }
            }
        }
    }

    return true;
}

bool ChannelIdZkAPI::UpdateChannelIdTime(const std::string& strChannelId, ApiError& error)
{
    std::string channelPath = ZK_CHANNEL_ID_CONFIG_PATH(strChannelId);
    int ret = m_ptrEventZkClient->ZkSet(channelPath, gse::tools::time::GetUTCTimeString(), -1, NULL);
    if (ret != GSE_SUCCESS)
    {
        error = ZooKeeperSetError;
        LOG_ERROR("failed to set channelid(%s), path:%s", strChannelId.c_str(), channelPath.c_str());
        return false;
    }

    return true;
}
bool ChannelIdZkAPI::CreateChannels(const std::string& strChannelId, const std::vector<std::shared_ptr<Channel>> m_channels, ApiError& error, bool create)
{
    std::string buff;

    int channelCount = m_channels.size();

    if (create)
    {
        int ret = m_ptrEventZkClient->ZkCreateNormal(ZK_CHANNEL_ID_CONFIG_CHANNEL_PATH(strChannelId), gse::tools::strings::ToString(channelCount), buff, m_acl);
        if (ret != GSE_SUCCESS && ret != GSE_ZK_NODE_EXIST)
        {
            error = ZooKeeperCreateError;
            LOG_ERROR("it is failed to create the channel id(%s)'s channel, path:%s",
                      strChannelId.c_str(), ZK_CHANNEL_ID_CONFIG_CHANNEL_PATH(strChannelId).c_str());
            return false;
        }
    }

    for (auto channel : m_channels)
    {
        std::string channelPath = ZK_CHANNEL_ID_CONFIG_WITH_CHANNELNAME_PATH(strChannelId, channel->m_name);

        if (create)
        {
            int ret = m_ptrEventZkClient->ZkCreateNormal(channelPath, channel->ToJSON(), buff, m_acl);
            if (ret != GSE_SUCCESS)
            {
                error = ZooKeeperCreateError;
                LOG_ERROR("failed to create channel node, path:%s", channelPath.c_str());
                return false;
            }
            LOG_DEBUG("successfully create filter, path:%s, value:%s", channelPath.c_str(), channel->ToJSON().c_str());
        }
        else
        {
            int ret = m_ptrEventZkClient->ZkSet(channelPath, channel->ToJSON(), -1, NULL);
            if (ret != GSE_SUCCESS)
            {

                int ret = m_ptrEventZkClient->ZkCreateNormal(channelPath, channel->ToJSON(), buff, m_acl);
                if (ret != GSE_SUCCESS)
                {
                    error = ZooKeeperCreateError;
                    LOG_ERROR("failed to create channel node, path:%s", channelPath.c_str());
                    return false;
                }
            }
        }
    }

    return true;
}

bool ChannelIdZkAPI::CreateChannelIdQueryIndex(std::shared_ptr<Metadata> ptrMetadata, const std::string& streamToId, ApiError& error)
{
    std::string strChannelid = gse::tools::strings::ToString(ptrMetadata->m_channelID);
    int ret = GSE_SUCCESS;

    if (!ptrMetadata->m_requestPlatName.empty())
    {
        std::string platIndexPath = ZK_CHANNEL_ID_CONFIG_WITH_LABLE_PLATNAME_PATH(ptrMetadata->m_requestPlatName);
        platIndexPath.append("/").append(strChannelid);

        if (!CreateNode(platIndexPath, strChannelid, error))
        {
            LOG_ERROR("failed to create platname index, path:%s, error:%s", platIndexPath.c_str(), GET_ERROR_MESSAGE(error).c_str());
            error = ZooKeeperCreateError;
            return false;
        }
    }

    if (ptrMetadata->m_label.m_bkBizId >= 0)
    {
        std::string bizidIndexPath = ZK_CHANNEL_ID_CONFIG_WITH_LABLE_BIZID_PATH(gse::tools::strings::ToString(ptrMetadata->m_label.m_bkBizId));
        bizidIndexPath.append("/").append(strChannelid);
        if (!CreateNode(bizidIndexPath, strChannelid, error))
        {
            LOG_ERROR("failed to create bizid index, path:%s, error:%s", bizidIndexPath.c_str(), GET_ERROR_MESSAGE(error).c_str());
            error = ZooKeeperCreateError;
            return false;
        }
    }

    if (!ptrMetadata->m_label.m_odm.empty())
    {
        std::string odmIndexPath = ZK_CHANNEL_ID_CONFIG_WITH_LABLE_ODM_PATH(ptrMetadata->m_label.m_odm);
        odmIndexPath.append("/").append(strChannelid);

        if (!CreateNode(odmIndexPath, strChannelid, error))
        {
            LOG_ERROR("failed to create odm index, path:%s, error:%s", odmIndexPath.c_str(), GET_ERROR_MESSAGE(error).c_str());
            error = ZooKeeperCreateError;
            return false;
        }
    }

    std::string streamToIdPath = ZK_CHANNEL_ID_CONFIG_WITH_LABLE_PATH();
    streamToIdPath.append("/stream_to_id/");
    streamToIdPath.append(streamToId);
    streamToIdPath.append("/");
    streamToIdPath.append(strChannelid);

    if (!CreateNode(streamToIdPath, strChannelid, error))
    {
        LOG_ERROR("failed to create streamto index, path:%s, error:%s", streamToIdPath.c_str(), GET_ERROR_MESSAGE(error).c_str());
        error = ZooKeeperCreateError;
        return false;
        return false;
    }

    return true;
}

bool ChannelIdZkAPI::CreateChannelIdNotify(const std::string& strChannelId, const std::string& value, std::string& errorMsg, bool create)
{
    std::string channelIdWatchPath = ZK_CHANNEL_ID_WATCH_PATH(strChannelId);

    std::string buff;
    int ret = m_ptrEventZkClient->ZkCreateNormal(channelIdWatchPath, value, buff, m_acl);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("failed to create channelid node, path:%s, error:%s", channelIdWatchPath.c_str(), m_ptrEventZkClient->ZkError(ret));
        return false;
    }

    LOG_DEBUG("successfully create channelid watch node, path:%s, value:%s", channelIdWatchPath.c_str(), value.c_str());
    return true;
}

bool ChannelIdZkAPI::SaveTlogChannelId(const std::shared_ptr<Metadata> ptrMetadata, ApiError& error)
{
    std::string tglogChannelidNotifyPath;
    tglogChannelidNotifyPath.append(std::string(ZK_CHANNEL_ID_TGLOG_NOTIFY_BASE_PATH));
    tglogChannelidNotifyPath.append("/").append(gse::tools::strings::ToString(ptrMetadata->m_label.m_bkBizId));
    tglogChannelidNotifyPath.append("/").append(ptrMetadata->m_label.m_odm);

    Json::Value data;
    data["channelid"] = Json::Value(ptrMetadata->m_channelID);

    Json::FastWriter jfwriter;
    std::string value = jfwriter.write(data);
    bool bret = false;
    int iret = GSE_ERROR;
    bret = m_ptrEventZkClient->ZkExists(tglogChannelidNotifyPath, NULL, NULL, NULL);
    if (bret)
    {
        int ret = m_ptrEventZkClient->ZkSet(tglogChannelidNotifyPath, value, -1, NULL);
        if (ret != GSE_SUCCESS)
        {

            LOG_ERROR("failed set tlog value, path:%s, value:%s", tglogChannelidNotifyPath.c_str());
            return false;
        }
        return true;
    }
    else
    {
        return CreateNode(tglogChannelidNotifyPath, value, error);
    }
}

bool ChannelIdZkAPI::DeleteTglogChannelId(int bizid, std::string& odm)
{
    std::string tglogChannelidNotifyPath;
    tglogChannelidNotifyPath.append(std::string(ZK_CHANNEL_ID_TGLOG_BASE_PATH));
    tglogChannelidNotifyPath.append("/").append(gse::tools::strings::ToString(bizid));
    tglogChannelidNotifyPath.append("/").append(odm);

    Json::Value data;
    data["channelid"] = Json::Value(0);
    Json::FastWriter jfwriter;
    std::string value = jfwriter.write(data);
    bool bret = false;
    int iret = GSE_ERROR;
    iret = m_ptrEventZkClient->ZkSet(tglogChannelidNotifyPath, value, -1, NULL);
    if (iret != GSE_SUCCESS)
    {
        LOG_ERROR("failed to set tlog value, path:%s, value:%s", tglogChannelidNotifyPath.c_str(), value.c_str());
        return false;
    }
    else
    {
        return true;
    }
}

bool ChannelIdZkAPI::GetTglogChannelId(const std::shared_ptr<Metadata> ptrMetadata, uint32_t& realChannelId, ApiError& error)
{
    std::string tglogChannelidNotifyPath;
    tglogChannelidNotifyPath.append(std::string(ZK_CHANNEL_ID_TGLOG_NOTIFY_BASE_PATH));
    tglogChannelidNotifyPath.append("/").append(gse::tools::strings::ToString(ptrMetadata->m_label.m_bkBizId));
    tglogChannelidNotifyPath.append("/").append(ptrMetadata->m_label.m_odm);

    std::string jsonValue;
    realChannelId = -1;
    int iret = m_ptrEventZkClient->ZkGet(tglogChannelidNotifyPath, jsonValue, NULL, NULL, NULL);
    if (iret)
    {
        LOG_DEBUG("TGLOG channelid not config, query path:%s", tglogChannelidNotifyPath.c_str());
        error = ZooKeeperGetError;
        return false;
    }
    else
    {
        Json::Reader reader(Json::Features::strictMode());
        Json::Value propertyCfgJson;
        if (!reader.parse(jsonValue, propertyCfgJson, false))
        {
            LOG_ERROR("failed to parse tglog json config, invalid json [%s]", jsonValue.c_str());
            return false;
        }

        if (propertyCfgJson.isMember("channelid") && propertyCfgJson["channelid"].isInt())
        {
            realChannelId = propertyCfgJson.get("channelid", -1).asInt();
            LOG_DEBUG("TGLOG real channel id(%d) exist", realChannelId);
            if (realChannelId <= 0)
            {
                LOG_DEBUG("TGLOG real channel id(%d) invalid, need create real channel id", realChannelId);
                return false;
            }

            return true;
        }
        LOG_ERROR("tglog config json invalid, missing channelid field, json [%s]", jsonValue.c_str());
        return false;
    }

    return false;
}

bool ChannelIdZkAPI::ReadFilters(const std::string& channelId, std::vector<std::shared_ptr<StreamFilter>>& filters, ApiError& error)
{
    std::vector<std::string> filterNames;
    std::string filtersPath = ZK_CHANNEL_ID_CONFIG_FILTER_PATH(channelId);

    int ret = m_ptrEventZkClient->ZkGetChildren(filtersPath, NULL, NULL, filterNames, NULL);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("it is failed to read the channelid(%s)'s property", filtersPath.c_str());
        error = ZooKeeperGetChildError;
        return false;
    }

    for (auto filterName : filterNames)
    {
        std::string filterValue;
        std::string zkFilterPath = ZK_CHANNEL_ID_CONFIG_WITH_FILTERNAME_PATH(channelId, filterName);
        int ret = m_ptrEventZkClient->ZkGet(zkFilterPath, filterValue, NULL, NULL, NULL);
        LOG_DEBUG("read filters config, zk path:%s, value:%s", zkFilterPath.c_str(), filterValue.c_str());
        if (GSE_SUCCESS != ret)
        {
            LOG_ERROR("it is failed to read the config for the zk node path (%s)", zkFilterPath.c_str());
            error = ZooKeeperGetError;
            return false;
        }

        auto streamFilter = std::make_shared<StreamFilter>();

        Json::Reader reader(Json::Features::strictMode());
        Json::Value propertyCfgJson;
        if (!reader.parse(filterValue, propertyCfgJson, false))
        {
            LOG_ERROR("failed to parse filter config, invalid json [%s]", filterValue.c_str());
            return false;
        }

        if (!streamFilter->Parse(propertyCfgJson, error))
        {
            LOG_ERROR("failed to parse filter json, str:%s, error;%s", filterValue.c_str(), GET_ERROR_MESSAGE(error).c_str());
            return false;
        }

        filters.push_back(streamFilter);
    }

    return true;
}

bool ChannelIdZkAPI::ReadChannels(const std::string& channelId, std::vector<std::shared_ptr<Channel>>& channels, ApiError& error)
{
    std::vector<std::string> channelNames;
    std::string channelsPath = ZK_CHANNEL_ID_CONFIG_CHANNEL_PATH(channelId);

    int ret = m_ptrEventZkClient->ZkGetChildren(channelsPath, NULL, NULL, channelNames, NULL);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("it is failed to read the channels(%s)", channelsPath.c_str());
        error = ZooKeeperGetChildError;
        return false;
    }

    for (auto channelName : channelNames)
    {
        std::string channelValue;
        std::string zkChannelPath = ZK_CHANNEL_ID_CONFIG_WITH_CHANNELNAME_PATH(channelId, channelName);
        int ret = m_ptrEventZkClient->ZkGet(zkChannelPath, channelValue, NULL, NULL, NULL);
        LOG_DEBUG("read channels config, zk path:%s, value:%s", zkChannelPath.c_str(), channelValue.c_str());
        if (GSE_SUCCESS != ret)
        {
            error = ZooKeeperSetError;
            LOG_ERROR("it is failed to read the channel config for the zk node path (%s)", zkChannelPath.c_str());
            return false;
        }

        auto channel = std::make_shared<Channel>();
        Json::Value jsonValue;
        Json::Reader reader(Json::Features::strictMode());
        if (!reader.parse(channelValue, jsonValue, false))
        {
            LOG_ERROR("failed to parse channel json, invalid json [%s]", channelValue.c_str());
            error = ReadFromConfigJsonInvalid;
            return false;
        }

        if (!channel->Parse(jsonValue, error))
        {
            LOG_ERROR("failed to parse channel json, str:%s, error;%s", channelValue.c_str(), GET_ERROR_MESSAGE(error).c_str());
            return false;
        }
        channels.push_back(channel);
    }

    return true;
}

bool ChannelIdZkAPI::ReadMetadata(const std::string& channelId, std::shared_ptr<Metadata> metadata, ApiError& error)
{
    std::string value;
    std::string metadataPath = ZK_CHANNEL_ID_CONFIG_METADATA_PATH(channelId);

    int ret = m_ptrEventZkClient->ZkGet(metadataPath, value, NULL, NULL, NULL);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("it is failed to read the metadata(%s)'s property", metadataPath.c_str());
        error = ZooKeeperGetError;
        return false;
    }

    Json::Value jsonValue;
    Json::Reader reader(Json::Features::strictMode());
    if (!reader.parse(value, jsonValue, false))
    {
        error = ReadFromConfigJsonInvalid;
        LOG_ERROR("failed to parse meta json config, invalid json [%s]", value.c_str());
        return false;
    }

    if (!metadata->ParseMetadata(jsonValue, error))
    {
        LOG_ERROR("failed to parse metadata when read from zk, error:%s", GET_ERROR_MESSAGE(error).c_str());
        return false;
    }

    return true;
}

std::shared_ptr<ChannelIdConfig> ChannelIdZkAPI::ReadChannelIdConfig(const std::string& channelId, const std::string& platName, ApiError& error)
{
    if (NULL == m_ptrEventZkClient)
    {
        error = ZooKeeperUnit;
        LOG_ERROR("it is failed to get the value of the channel id (%s), because of the zookeeper server is not connected", channelId.c_str());
        return nullptr;
    }

    std::string channelIdPath = ZK_CHANNEL_ID_CONFIG_PATH(channelId);
    if (!m_ptrEventZkClient->ZkExists(channelIdPath, NULL, NULL, NULL))
    {
        std::string msg = "channel_id(" + channelId + ") not exist";
        error = ChannelIdNotExist;
        LOG_INFO("channelid not exist, zk path:%s", channelIdPath.c_str());
        return nullptr;
    }

    auto ptrChannelIdConfig = std::make_shared<ChannelIdConfig>();
    LOG_DEBUG("Read channelid config, platname:%s, channelid:%s", platName.c_str(), channelId.c_str());

    if (!ReadMetadata(channelId, ptrChannelIdConfig->m_ptrMetadata, error))
    {
        LOG_ERROR("failed to read metadata config, error:%s", GET_ERROR_MESSAGE(error).c_str());
        return nullptr;
    }

    if (!ReadChannels(channelId, ptrChannelIdConfig->m_channels, error))
    {
        LOG_ERROR("failed to read channel config, error:%s", GET_ERROR_MESSAGE(error).c_str());
        return nullptr;
    }

    if (!ReadFilters(channelId, ptrChannelIdConfig->m_streamFilters, error))
    {
        LOG_ERROR("failed to read streamfilter config, error:%s", GET_ERROR_MESSAGE(error).c_str());
        return nullptr;
    }

    return ptrChannelIdConfig;
}

bool ChannelIdZkAPI::DeleteChannelIdBizIdIndex(const std::string& bizId, const std::string& channelId)
{
    if (NULL == m_ptrEventZkClient)
    {
        LOG_ERROR("failed to delete channelid bizId index, the zookeeper server is not connected");
        return false;
    }

    std::string bizid_name_index_path = ZK_CHANNEL_ID_CONFIG_WITH_LABLE_BIZID_PATH(bizId);
    bizid_name_index_path.append("/").append(channelId);
    int ret = m_ptrEventZkClient->ZkDelete(bizid_name_index_path, -1);
    if (ret != GSE_SUCCESS)
    {
        if (ret == GSE_ZK_NODE_NOTEXIST)
        {
            LOG_INFO("failed to delete bizid index, path:%s, node not exist", bizid_name_index_path.c_str());
            return true;
        }

        LOG_ERROR("failed to delete bizid index, path:%s", bizid_name_index_path.c_str());
        return false;
    }
    return true;
}

bool ChannelIdZkAPI::DeleteChannelIdPlatNameIndex(const std::string& platName, const std::string& channelId)
{
    if (NULL == m_ptrEventZkClient)
    {
        LOG_ERROR("failed to delete channelid plat index, the zookeeper server is not connected");
        return false;
    }

    std::string plat_name_index_path = ZK_CHANNEL_ID_CONFIG_WITH_LABLE_PLATNAME_PATH(platName);
    plat_name_index_path.append("/").append(channelId);
    int ret = m_ptrEventZkClient->ZkDelete(plat_name_index_path, -1);
    if (ret != GSE_SUCCESS)
    {
        if (ret == GSE_ZK_NODE_NOTEXIST)
        {
            LOG_INFO("failed to delete plat_name index, path:%s, node not exist", plat_name_index_path.c_str());
            return true;
        }

        LOG_ERROR("failed to delete plat_name index, path:%s", plat_name_index_path.c_str());
        return false;
    }
    return true;
}

bool ChannelIdZkAPI::DeleteChannelIdOdmNameIndex(const std::string& label, const std::string& channelId)
{
    if (NULL == m_ptrEventZkClient)
    {
        LOG_ERROR("failed to delete channelid odm index, the zookeeper server is not connected");
        return false;
    }
    std::string odm_name_index_path = ZK_CHANNEL_ID_CONFIG_WITH_LABLE_ODM_PATH(label);
    odm_name_index_path.append("/").append(channelId);
    int ret = m_ptrEventZkClient->ZkDelete(odm_name_index_path, -1);
    if (ret != GSE_SUCCESS)
    {
        if (ret == GSE_ZK_NODE_NOTEXIST)
        {
            LOG_INFO("failed to delete odm index, path:%s, node not exist", odm_name_index_path.c_str());
            return true;
        }

        LOG_ERROR("failed to delete odm index, path:%s", odm_name_index_path.c_str());
        return false;
    }
    return true;
}

bool ChannelIdZkAPI::DeleteChannelIdStreamToIndex(const std::string& streamToId, const std::string& channelId)
{
    if (NULL == m_ptrEventZkClient)
    {
        LOG_ERROR("failed to delete channelid stream_to_id index, the zookeeper server is not connected");
        return false;
    }
    std::string channelIdStreamToIndexPath = ZK_CHANNEL_ID_CONFIG_WITH_LABLE_STREAMTO_ID_PATH(streamToId, channelId);
    int ret = m_ptrEventZkClient->ZkDelete(channelIdStreamToIndexPath, -1);
    if (ret != GSE_SUCCESS)
    {
        if (ret == GSE_ZK_NODE_NOTEXIST)
        {
            LOG_INFO("failed to delete stream_to_id index, path:%s, node not exist", channelIdStreamToIndexPath.c_str());
            return true;
        }
        LOG_ERROR("failed to delete stream_to_id index, path:%s", channelIdStreamToIndexPath.c_str());
        return false;
    }

    return true;
}

bool ChannelIdZkAPI::CreateStreamToClusterNode(int64_t streamToId, std::shared_ptr<StreamToCluster> ptrStreamToCluster, ApiError& error)
{
    std::string strStreamToClusterId = gse::tools::strings::ToString(streamToId);
    std::string streamToClusterIdZkPath = ZK_STREAM_ID_CONFIG_PATH(strStreamToClusterId);
    if (NULL == m_ptrEventZkClient)
    {
        error = ZooKeeperUnit;
        LOG_ERROR("it is failed to update stream_to_id (%d), because of the zookeeper server is not connected", streamToId);
        return false;
    }

    // update streamTO
    std::string buff;
    int ret = m_ptrEventZkClient->ZkCreateNormal(streamToClusterIdZkPath, ptrStreamToCluster->ToJSON(), buff, m_acl);
    if (ret != GSE_SUCCESS)
    {
        error = ZooKeeperCreateError;
        LOG_ERROR("failed to create stream cluster config, path:%s, value[%s]", streamToClusterIdZkPath.c_str(), ptrStreamToCluster->ToJSON().c_str());
        return false;
    }

    return true;
}

bool ChannelIdZkAPI::CreateStreamToClusterMetadataNode(int64_t streamToId, std::shared_ptr<StreamToMetadata> ptrMeta, ApiError& error)
{
    std::string strStreamToClusterId = gse::tools::strings::ToString(streamToId);
    std::string streamToClusterIdMetaZkPath = ZK_STREAM_TO_ID_CONFIG_METADATA_PATH(strStreamToClusterId);
    if (NULL == m_ptrEventZkClient)
    {
        error = ZooKeeperUnit;
        LOG_ERROR("it is failed to create stream_to_id (%d), because of the zookeeper server is not connected", streamToId);
        return false;
    }

    // update streamTO
    std::string buff;
    int ret = m_ptrEventZkClient->ZkCreateNormal(streamToClusterIdMetaZkPath, ptrMeta->ToJSON(), buff, m_acl);
    if (ret != GSE_SUCCESS)
    {
        error = ZooKeeperCreateError;
        LOG_ERROR("failed to create stream cluster metadata config, path:%s, value[%s]", streamToClusterIdMetaZkPath.c_str(), ptrMeta->ToJSON().c_str());
        return false;
    }

    return true;
}

bool ChannelIdZkAPI::CreateStreamToIdQueryIndex(std::shared_ptr<StreamToMetadata> ptr_metadata, const std::string& type, ApiError& error)
{
    std::string strStreamToId = gse::tools::strings::ToString(ptr_metadata->m_streamToId);
    std::string strTime = gse::tools::time::GetUTCTimeString();
    if (!ptr_metadata->m_requestPlatName.empty())
    {
        std::string platIndexPath = ZK_STREAM_TO_ID_CONFIG_WITH_LABLE_PLATNAME_PATH(ptr_metadata->m_requestPlatName);
        platIndexPath.append("/").append(strStreamToId);
        if (!CreateNode(platIndexPath, strTime, error))
        {
            error = ZooKeeperCreateError;
            LOG_ERROR("failed to create stream_to_id index, path:%s", platIndexPath.c_str());
            return false;
        }
    }

    if (ptr_metadata->m_label.m_bkBizId >= 0)
    {
        std::string bizidIndexPath = ZK_STREAM_TO_ID_CONFIG_WITH_LABLE_BIZID_PATH(gse::tools::strings::ToString(ptr_metadata->m_label.m_bkBizId));
        bizidIndexPath.append("/").append(strStreamToId);
        if (!CreateNode(bizidIndexPath, strTime, error))
        {
            error = ZooKeeperCreateError;
            LOG_ERROR("failed to create stream_to_id bizid index, path:%s", bizidIndexPath.c_str());
            return false;
        }
    }

    if (!ptr_metadata->m_label.m_odm.empty())
    {
        std::string odmIndexPath = ZK_STREAM_TO_ID_CONFIG_WITH_LABLE_ODM_PATH(ptr_metadata->m_label.m_odm);
        odmIndexPath.append("/").append(strStreamToId);
        if (!CreateNode(odmIndexPath, strTime, error))
        {
            error = ZooKeeperCreateError;
            LOG_ERROR("failed to create stream_to_id odm index, path:%s", odmIndexPath.c_str());
            return false;
        }
    }

    if (!type.empty())
    {
        std::string clusterTypeIndexPath = ZK_STREAM_TO_ID_CONFIG_WITH_TYPE_PATH(type);
        clusterTypeIndexPath.append("/").append(strStreamToId);
        if (!CreateNode(clusterTypeIndexPath, strTime, error))
        {
            error = ZooKeeperCreateError;
            LOG_ERROR("failed to create stream_to_id cluster type index, path:%s", clusterTypeIndexPath.c_str());
            return false;
        }
    }

    return true;
}

bool ChannelIdZkAPI::ReadStreamToMetadata(const std::string& streamToId, std::shared_ptr<StreamToMetadata> metadata, ApiError& error)
{
    std::string streamToClusterIdMetaZkPath = ZK_STREAM_TO_ID_CONFIG_METADATA_PATH(streamToId);

    std::string propertyValue;
    int ret = m_ptrEventZkClient->ZkGet(streamToClusterIdMetaZkPath, propertyValue, NULL, NULL, NULL);
    if (GSE_SUCCESS != ret)
    {
        LOG_ERROR("it is failed to read the meta for the zk node path (%s)", SAFE_CSTR(streamToClusterIdMetaZkPath.c_str()));
        error = ZooKeeperGetError;
        return false;
    }

    Json::Value propertyCfgJson;
    Json::Reader reader(Json::Features::strictMode());
    if (!reader.parse(propertyValue, propertyCfgJson))
    {
        LOG_ERROR("failed to parse meta config json, path:%s, value:%s", streamToClusterIdMetaZkPath.c_str(),
                  propertyValue.c_str());
        error = ReadFromConfigJsonInvalid;
        return false;
    }

    return metadata->ParseMetadata(propertyCfgJson, error);
}

bool ChannelIdZkAPI::ReadStreamToIdZkdata(const std::string& streamToId, std::shared_ptr<StreamToCluster> ptrStreamToConfig, ApiError& error)
{
    std::string propertyValue;
    std::string streamToIdConfigZkPath = ZK_STREAM_ID_CONFIG_PATH(streamToId);

    int ret = m_ptrEventZkClient->ZkGet(streamToIdConfigZkPath, propertyValue, NULL, NULL, NULL);
    if (GSE_SUCCESS != ret)
    {
        error = ZooKeeperGetError;
        LOG_ERROR("failed to read the config for the zk node path (%s)", (streamToIdConfigZkPath.c_str()));
        return false;
    }

    Json::Value propertyCfgJson;
    Json::Reader reader(Json::Features::strictMode());
    if (!reader.parse(propertyValue, propertyCfgJson))
    {
        error = ReadFromConfigJsonInvalid;
        return false;
    }

    if (!ptrStreamToConfig->Parse(propertyCfgJson, error))
    {
        return false;
    }

    return true;
}

std::shared_ptr<StreamToIdConfig> ChannelIdZkAPI::ReadStreamToIdConfig(const std::string& streamToId, ApiError& error)
{
    if (NULL == m_ptrEventZkClient)
    {
        error = ZooKeeperUnit;
        return nullptr;
    }

    auto ptrStreamToIdConfig = std::make_shared<StreamToIdConfig>();

    if (!ReadStreamToMetadata(streamToId, ptrStreamToIdConfig->m_ptrMetadata, error))
    {
        LOG_ERROR("failed to read streamToId metadata, stream_to_id;%s, error:%s", streamToId.c_str(), GET_ERROR_MESSAGE(error).c_str());
        return nullptr;
    }

    if (!ReadStreamToIdZkdata(streamToId, ptrStreamToIdConfig->m_ptrStreamToCluster, error))
    {
        LOG_ERROR("failed to read streamToId config, stream_to_id;%s, error:%s", streamToId.c_str(), GET_ERROR_MESSAGE(error).c_str());
        return nullptr;
    }

    return ptrStreamToIdConfig;
}

bool ChannelIdZkAPI::UpdateStreamToClusterConfig(int64_t streamToId, std::shared_ptr<StreamToCluster> ptrStreamToCluster, ApiError& error)
{
    std::string strStreamToClusterId = gse::tools::strings::ToString(streamToId);
    std::string streamToClusterIdZkPath = ZK_STREAM_ID_CONFIG_PATH(strStreamToClusterId);
    if (NULL == m_ptrEventZkClient)
    {
        error = ZooKeeperUnit;
        LOG_ERROR("it is failed to update stream_to_id (%d), because of the zookeeper server is not connected", streamToId);
        return false;
    }

    if (!ExistStreamToClusterId(strStreamToClusterId))
    {
        error = StreamToIdNotExist;
        return false;
    }

    // update streamTO
    int ret = m_ptrEventZkClient->ZkSet(streamToClusterIdZkPath, ptrStreamToCluster->ToJSON(), -1, nullptr);
    if (ret != GSE_SUCCESS)
    {
        error = ZooKeeperSetError;
        LOG_ERROR("failed to set stream cluster config, path:%s, value[%s]", streamToClusterIdZkPath.c_str(), ptrStreamToCluster->ToJSON().c_str());
    }

    return true;
}

bool ChannelIdZkAPI::UpdateStreamToClusterMeta(int64_t streamToId, std::shared_ptr<StreamToMetadata> ptrMeta, ApiError& error)
{
    std::string strStreamToClusterId = gse::tools::strings::ToString(streamToId);
    std::string streamToClusterIdMetaZkPath = ZK_STREAM_TO_ID_CONFIG_METADATA_PATH(strStreamToClusterId);
    if (NULL == m_ptrEventZkClient)
    {
        error = ZooKeeperUnit;
        LOG_ERROR("it is failed to update stream_to_id (%d), because of the zookeeper server is not connected", streamToId);
        return false;
    }

    // update streamTO
    int ret = m_ptrEventZkClient->ZkSet(streamToClusterIdMetaZkPath, ptrMeta->ToJSON(), -1, nullptr);
    if (ret != GSE_SUCCESS)
    {
        error = ZooKeeperSetError;
        LOG_ERROR("failed to set stream cluster metadata config, path:%s, value[%s]", streamToClusterIdMetaZkPath.c_str(), ptrMeta->ToJSON().c_str());
        return false;
    }

    return true;
}

bool ChannelIdZkAPI::DeleteStreamToID(const std::string& streamToId, ApiError& error)
{
    if (NULL == m_ptrEventZkClient)
    {
        error = ZooKeeperUnit;
        LOG_ERROR("it is failed to delete stream_to_id (%d), because of the zookeeper server is not connected", streamToId.c_str());
        return false;
    }

    std::string streamToIdPath = ZK_STREAM_ID_CONFIG_PATH(streamToId);
    bool bret = m_ptrEventZkClient->ZkExists(streamToIdPath, NULL, NULL, NULL);
    if (!bret)
    {
        error = ZooKeeperExistError;
        LOG_ERROR("failed to check node exist, stream_to_id:%s, not exist, zk path:%s", streamToIdPath.c_str());
        return false;
    }

    if (!DeleteRecursive(streamToIdPath, -1, error))
    {
        LOG_ERROR("failed to delete stream_to_id config, path:%s, error:%s", streamToIdPath.c_str(), GET_ERROR_MESSAGE(error).c_str());
        return false;
    }

    return true;
}

bool ChannelIdZkAPI::DeleteStreamToPlatNameIndex(const std::string& platName, const std::string& streamToId, ApiError& error)
{
    if (NULL == m_ptrEventZkClient)
    {
        error = ZooKeeperUnit;
        return false;
    }

    std::string streamToPlatIndexPath = ZK_STREAM_TO_ID_CONFIG_WITH_LABLE_PLATNAME_PATH(platName);
    streamToPlatIndexPath.append("/").append(streamToId);

    int ret = m_ptrEventZkClient->ZkDelete(streamToPlatIndexPath, -1);
    if (ret != GSE_SUCCESS)
    {
        if (ret == GSE_ZK_NODE_NOTEXIST)
        {
            LOG_INFO("failed to delete streamto plat index, path:%s, node not exist", streamToPlatIndexPath.c_str());
            return true;
        }

        LOG_ERROR("failed to delete streamto plat index, path:%s", streamToPlatIndexPath.c_str());
        return false;
    }

    return true;
}

bool ChannelIdZkAPI::DeleteStreamToBkBizIdIndex(const std::string& bkBizid, const std::string& streamToId, ApiError& error)
{
    if (NULL == m_ptrEventZkClient)
    {
        error = ZooKeeperUnit;
        return false;
    }

    std::string bizidIndexPath = ZK_STREAM_TO_ID_CONFIG_WITH_LABLE_BIZID_PATH(bkBizid);
    bizidIndexPath.append("/").append(streamToId);

    int ret = m_ptrEventZkClient->ZkDelete(bizidIndexPath, -1);

    if (ret != GSE_SUCCESS)
    {
        if (ret == GSE_ZK_NODE_NOTEXIST)
        {
            LOG_INFO("failed to delete streamto bizid index, path:%s, node not exist", bizidIndexPath.c_str());
            return true;
        }

        LOG_ERROR("failed to delete streamto bizid  index, path:%s", bizidIndexPath.c_str());
        return false;
    }

    return true;
}

bool ChannelIdZkAPI::DeleteOriginStreamToID(const std::string& strStreamToId, ApiError& error)
{
    if (NULL == m_ptrEventZkClient)
    {
        error = ZooKeeperUnit;
        return false;
    }

    int streamToId = gse::tools::strings::StringToInt32(strStreamToId);
    int originId = GetOriginStreamToID(streamToId);
    std::string strOriginId = gse::tools::strings::ToString(originId);

    std::stringstream ss;
    ss << std::setw(10) << std::setfill('0') << strOriginId;

    int ret = m_ptrEventZkClient->ZkDelete(ZK_STREAM_TO_ID_ORIGIN_PATH() + "/" + ss.str());
    if (ret != GSE_SUCCESS)
    {
        if (ret == GSE_ZK_NODE_NOTEXIST)
        {
            LOG_INFO("failed to  delete the origin stream_to_id (%s), node not exist", ss.str().c_str());
            return true;
        }

        LOG_ERROR("it is failed to delete the origin stream_to_id (%s)", ss.str().c_str());
        return false;
    }

    return true;
}

bool ChannelIdZkAPI::DeleteClusterTypeIndex(const std::string& clusterType, const std::string& strStreamToId, ApiError& error)
{
    if (NULL == m_ptrEventZkClient)
    {
        error = ZooKeeperUnit;
        return false;
    }

    std::string clusterIndexPath = ZK_STREAM_TO_ID_CONFIG_WITH_TYPE_PATH(clusterType);
    clusterIndexPath.append("/").append(strStreamToId);

    int ret = m_ptrEventZkClient->ZkDelete(clusterIndexPath, -1);

    if (ret != GSE_SUCCESS)
    {
        if (ret == GSE_ZK_NODE_NOTEXIST)
        {
            LOG_INFO("failed to delete streamto cluster type index, path:%s, node not exist", clusterIndexPath.c_str());
            return true;
        }

        LOG_ERROR("failed to delete streamto cluster type index, path:%s", clusterIndexPath.c_str());
        return false;
    }

    return true;
}
bool ChannelIdZkAPI::DeleteRecursive(const std::string& path, int32_t version /*= -1*/, ApiError& error)
{
    if (NULL == m_ptrEventZkClient)
    {
        error = ZooKeeperUnit;
        return false;
    }

    //获取child 结点
    std::vector<string> childNodes;
    childNodes.clear();
    int ret = m_ptrEventZkClient->ZkGetChildren(path, NULL, NULL, childNodes, NULL);
    if (ret != GSE_SUCCESS)
    {
        error = ZooKeeperGetChildError;
        LOG_ERROR("failed to get child node, error:%s", m_ptrEventZkClient->ZkError(ret));
        return false;
    }
    else
    {
        vector<string>::iterator iter = childNodes.begin();
        for (; iter != childNodes.end(); iter++)
        {
            string childPath = path + "/" + (*iter);
            if (!DeleteRecursive(childPath, -1, error))
            {
                LOG_ERROR("failed to recursive delete node");
                return false;
            }
        }

        int ret = m_ptrEventZkClient->ZkDelete(path, -1);
        if (ret != GSE_SUCCESS)
        {
            if (ret == GSE_ZK_NODE_NOTEXIST)
            {
                LOG_INFO("failed to delete zk node:%s, node not exist", path.c_str());
                return true;
            }

            LOG_ERROR("failed to delete zk node, path:%s", path.c_str());
            return false;
        }
        return true;
    }
    return true;
}

bool ChannelIdZkAPI::DeleteOriginChannelID(const std::string& originChannelID, ApiError& error)
{
    if (NULL == m_ptrEventZkClient)
    {
        error = ZooKeeperUnit;
        return false;
    }

    std::stringstream ss;
    ss << std::setw(10) << std::setfill('0') << originChannelID;
    std::string originIdZkPath = ZK_CHANNEL_ID_ORIGIN_PATH() + "/" + ss.str();
    int ret = m_ptrEventZkClient->ZkDelete(originIdZkPath, -1);
    if (ret != GSE_SUCCESS)
    {
        if (ret == GSE_ZK_NODE_NOTEXIST)
        {
            LOG_INFO("failed to delete the origin channel id (%s), node not exist", originChannelID.c_str());
            return true;
        }

        error = ZooKeeperDeleteError;
        LOG_ERROR("failed to delete the origin channel id (%s)", originChannelID.c_str());
        return false;
    }

    return true;
}

bool ChannelIdZkAPI::DeleteChannelId(int64_t channelId, ApiError& error)
{
    if (NULL == m_ptrEventZkClient)
    {
        error = ZooKeeperUnit;
        return false;
    }

    std::string strChannelId = gse::tools::strings::ToString(channelId);

    std::string channelIdPath = ZK_CHANNEL_ID_CONFIG_PATH(strChannelId);
    if (!DeleteRecursive(channelIdPath, -1, error))
    {
        LOG_ERROR("failed to delete channelid, path:%s, error:%s", channelIdPath.c_str(), GET_ERROR_MESSAGE(error).c_str());
        // continue delete originid
    }

    int originChannelId = GetOriginChannelID(channelId);

    return DeleteOriginChannelID(gse::tools::strings::ToString(originChannelId), error);
}

bool ChannelIdZkAPI::DeleteBySpecification(int64_t channelId, std::shared_ptr<DeleteSpecification> deleteSpecification, ApiError& error)
{
    std::string strChannelId = gse::tools::strings::ToString(channelId);

    std::string streamFilterBaesPath = ZK_CHANNEL_ID_CONFIG_FILTER_PATH(strChannelId);
    for (auto filter : deleteSpecification->m_streamFilterNames)
    {
        std::string streamFilterNamePath = streamFilterBaesPath + "/" + filter;
        LOG_DEBUG("delete filter zk path:%s", streamFilterNamePath.c_str());
        int ret = m_ptrEventZkClient->ZkDelete(streamFilterNamePath, -1);
        if (ret != GSE_SUCCESS)
        {
            if (ret == GSE_ZK_NODE_NOTEXIST)
            {
                LOG_INFO("failed to delete filter, path:%s, node not exist", streamFilterNamePath.c_str());
                return true;
            }

            error = ZooKeeperDeleteError;
            LOG_ERROR("failed to delete filter, path:%s", streamFilterNamePath.c_str());
            return false;
        }
    }

    std::string channelBasePath = ZK_CHANNEL_ID_CONFIG_CHANNEL_PATH(strChannelId);
    for (auto channelName : deleteSpecification->m_channelNames)
    {
        std::string channelPath = channelBasePath + "/" + channelName;
        LOG_DEBUG("delete channel zk path:%s", channelPath.c_str());
        int ret = m_ptrEventZkClient->ZkDelete(channelPath, -1);
        if (ret != GSE_SUCCESS)
        {

            error = ZooKeeperDeleteError;
            LOG_ERROR("failed to delete channel, path:%s", channelPath.c_str());
            return false;
        }
    }

    return true;
}

std::shared_ptr<gse::discover::zkapi::ZkApi> ChannelIdZkAPI::GetZkApiClient()
{
    return m_ptrEventZkClient;
}

bool SaveToDataIdZk(const std::string& dataId, const std::string& value, ApiError& error)
{
    std::vector<std::shared_ptr<DataIdZkClient>> dataIdZkClients = DataIdZkApiInst();
    for (auto dataIdZkClient : dataIdZkClients)
    {
        dataIdZkClient->SaveDataId(dataId, value, error);
    }
    return true;
}

bool DeleteDataIdZk(const std::string& dataId, ApiError& error)
{
    std::vector<std::shared_ptr<DataIdZkClient>> dataIdZkClients = DataIdZkApiInst();
    for (auto dataIdZkClient : dataIdZkClients)
    {
        dataIdZkClient->DeleteDataId(dataId, error);
    }
    return true;
}

bool UpdateDataIdZk(const std::string& dataId, const std::string& value, ApiError& error)
{
    std::vector<std::shared_ptr<DataIdZkClient>> dataIdZkClients = DataIdZkApiInst();
    for (auto dataIdZkClient : dataIdZkClients)
    {
        dataIdZkClient->UpdateDataId(dataId, value, error);
    }
    return true;
}

bool CreateClusterIdToZk(const std::string& clusterId, const std::string& value, ApiError& error)
{
    std::vector<std::shared_ptr<DataIdZkClient>> dataIdZkClients = DataIdZkApiInst();
    for (auto dataIdZkClient : dataIdZkClients)
    {
        dataIdZkClient->CreateClusterIdConfig(clusterId, value, error);
    }
    return true;
}

bool DeleteClusterIdFromZk(const std::string& clusterId, ApiError& error)
{
    std::vector<std::shared_ptr<DataIdZkClient>> dataIdZkClients = DataIdZkApiInst();
    for (auto dataIdZkClient : dataIdZkClients)
    {
        dataIdZkClient->DeleteClusterIdConfig(clusterId, error);
    }
    return true;
}

bool UpdateClusterIdToZk(const std::string& dataId, const std::string& value, ApiError& error)
{
    std::vector<std::shared_ptr<DataIdZkClient>> dataIdZkClients = DataIdZkApiInst();
    for (auto dataIdZkClient : dataIdZkClients)
    {
        dataIdZkClient->UpdateClusterIdConfig(dataId, value, error);
    }
    return true;
}

bool DataIdZkClient::DeleteDataId(const std::string& dataId, ApiError& error)
{
    std::string dataIdPath = m_dataIdZkPath + "/" + dataId;
    std::string buffer;
    if (m_ptrEventZkClient == nullptr)
    {
        error = ZooKeeperUnit;
        return false;
    }

    int ret = m_ptrEventZkClient->ZkDelete(dataIdPath, -1);
    if (ret != GSE_SUCCESS)
    {
        if (ret == GSE_ZK_NODE_NOTEXIST)
        {
            LOG_INFO("failed to delete dataid, path:%s, node not exist", dataIdPath.c_str());
            return true;
        }

        error = ZooKeeperDeleteError;
        LOG_ERROR("failed to delete dataid, path:%s", dataIdPath.c_str());
        return false;
    }

    return true;
}

bool DataIdZkClient::CreateNode(const std::string& nodePath, std::string& value, ApiError& error)
{
    if (NULL == m_ptrEventZkClient)
    {
        error = ZooKeeperUnit;
        LOG_ERROR("it is failed to create the node(%s), because of the zookeeper server is not connected", nodePath.c_str());
        return false;
    }

    std::string buff;
    std::vector<std::string> nodes;
    ChannelIdCommonAPI::SplitZkPath(nodePath, nodes);
    std::string path;
    for (std::size_t size = nodes.size(), idx = 0; idx != size; ++idx)
    {
        if (0 != idx)
        {
            path.append("/");
        }
        path.append(nodes.at(idx));
        if (!m_ptrEventZkClient->ZkExists(path, NULL, NULL, NULL))
        {
            int ret = 0;
            if (idx == (size - 1))
            {
                LOG_DEBUG("create leaf node:%s, value:%s", path.c_str(), value.c_str());
                ret = m_ptrEventZkClient->ZkCreateNormal(path, value, buff, m_acl);
            }
            else
            {
                ret = m_ptrEventZkClient->ZkCreateNormal(path, " ", buff, m_acl);
            }

            if (ret != GSE_SUCCESS)
            {
                error = ZooKeeperCreateError;
                LOG_ERROR("failed to create the zk node[%s], ret:%d ", path.c_str(), ret);
                return false;
            }
        }
    }

    return true;
}

bool DataIdZkClient::Init()
{
    ApiError error;
    std::string buff;
    bool ret = CreateNode(m_dataIdZkPath, buff, error);
    ret = CreateNode(m_storageZkPath, buff, error);

    if (m_zkAuth != "")
    {
        m_acl = true;
    }
    return true;
}

bool DataIdZkClient::SaveDataId(const std::string& dataId, const std::string& value, ApiError& error)
{
    std::string dataIdPath = m_dataIdZkPath + "/" + dataId;
    std::string buffer;
    if (m_ptrEventZkClient == nullptr)
    {
        error = ZooKeeperUnit;
        return false;
    }
    LOG_DEBUG("save dataid(%s) path:%s", dataId.c_str(), dataIdPath.c_str());
    int ret = m_ptrEventZkClient->ZkCreateNormal(dataIdPath, value, buffer, m_acl);
    if (ret != GSE_SUCCESS)
    {
        error = ZooKeeperCreateError;
        LOG_ERROR("failed to create dataid, path:%s", dataIdPath.c_str());
        return false;
    }

    return true;
}

bool DataIdZkClient::UpdateDataId(const std::string& dataId, const std::string& value, ApiError& error)
{
    std::string dataIdPath = m_dataIdZkPath + "/" + dataId;
    std::string buffer;
    if (m_ptrEventZkClient == nullptr)
    {
        error = ZooKeeperUnit;
        return false;
    }
    LOG_DEBUG("save dataid(%s) path:%s", dataId.c_str(), dataIdPath.c_str());
    int ret = m_ptrEventZkClient->ZkSet(dataIdPath, value, -1, NULL);
    if (ret != GSE_SUCCESS)
    {
        error = ZooKeeperCreateError;
        LOG_ERROR("failed to create dataid, path:%s", dataIdPath.c_str());
        return false;
    }

    return true;
}

bool DataIdZkClient::CreateClusterIdConfig(const std::string& clusterIndex, const std::string& value, ApiError& error)
{
    std::string clusterIdPath = m_storageZkPath + "/0_" + clusterIndex;
    std::string buffer;
    if (m_ptrEventZkClient == nullptr)
    {
        error = ZooKeeperUnit;
        return false;
    }
    LOG_DEBUG("save clusterid(%s) path:%s", clusterIndex.c_str(), clusterIdPath.c_str());
    int ret = m_ptrEventZkClient->ZkCreateNormal(clusterIdPath, value, buffer, m_acl);
    if (ret != GSE_SUCCESS)
    {
        error = ZooKeeperCreateError;
        LOG_ERROR("failed to create clusterid, path:%s", clusterIdPath.c_str());
        return false;
    }
    return true;
}

bool DataIdZkClient::UpdateClusterIdConfig(const std::string& clusterIndex, const std::string& value, ApiError& error)
{
    std::string clusterIdPath = m_storageZkPath + "/0_" + clusterIndex;
    std::string buffer;
    if (m_ptrEventZkClient == nullptr)
    {
        error = ZooKeeperUnit;
        return false;
    }
    LOG_DEBUG("save clusterid(%s) path:%s", clusterIndex.c_str(), clusterIdPath.c_str());
    int ret = m_ptrEventZkClient->ZkSet(clusterIdPath, value, -1, NULL);
    if (ret != GSE_SUCCESS)
    {
        error = ZooKeeperSetError;
        LOG_ERROR("failed to update clusterid, path:%s", clusterIdPath.c_str());
        return false;
    }
    return true;
}

bool DataIdZkClient::DeleteClusterIdConfig(const std::string& clusterIndex, ApiError& error)
{
    std::string clusterIdPath = m_storageZkPath + "/0_" + clusterIndex;
    std::string buffer;
    if (m_ptrEventZkClient == nullptr)
    {
        error = ZooKeeperUnit;
        return false;
    }
    LOG_DEBUG("delete clusterid(%s) path:%s", clusterIndex.c_str(), clusterIdPath.c_str());
    int ret = m_ptrEventZkClient->ZkDelete(clusterIdPath, -1);
    if (ret != GSE_SUCCESS)
    {
        if (ret == GSE_ZK_NODE_NOTEXIST)
        {
            LOG_INFO("failed to delete stream_to_id, path:%s, node not exist", clusterIdPath.c_str());
            return true;
        }

        error = ZooKeeperDeleteError;
        LOG_ERROR("failed to delete stream_to_id, path:%s", clusterIdPath.c_str());
        return false;
    }
    return true;
}

} // namespace data
} // namespace gse
