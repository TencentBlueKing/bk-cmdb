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

#include "stream_exporter_loader.h"

#include "api/channelid_def.h"
#include "channel_id_config.h"
#include "tools/rapidjson_macro.h"

namespace gse {
namespace data {

StreamExporterIDLoader::StreamExporterIDLoader(std::shared_ptr<gse::discover::zkapi::ZkApi> zkCli)
    : m_zkClient(zkCli)
{
}
StreamExporterIDLoader::StreamExporterIDLoader() {}

StreamExporterIDLoader::~StreamExporterIDLoader()
{
}

// typedef void (*ZK_WATCH_FUN)(int32_t type, int32_t state, const char *path, void *wctx);
void StreamExporterIDLoader::GetStreamToIdListCallback(int32_t type, int32_t state, const char *path, void *wctx)
{
    LOG_DEBUG("get streamto id list callback, type:%d, state:%d, path:%s", type, state, path);

    StreamExporterIDLoader *self = (StreamExporterIDLoader *)wctx;
    self->LoadStreamExporterConfig();
}

std::string StreamExporterIDLoader::GetStreamIDFromZkPath(const char *zkPath)
{
    std::string path(zkPath);
    std::size_t pos = path.find_last_of("/");
    LOG_DEBUG("split the zkpath(%s), stream id pos is (%u)", zkPath, pos);
    if ((pos + 1) != std::string::npos)
    {
        std::string id = path.substr(pos + 1);
        if (gse::tools::strings::IsNumber(id))
        {
            return id;
        }
        else
        {
            LOG_ERROR("zkpath include stream id, not number, path:%s", zkPath);
            return "";
        }
    }
    LOG_ERROR("zkpath not include stream id, path:%s", zkPath);
    return "";
}

// typedef void (*ZK_WATCH_FUN)(int32_t type, int32_t state, const char *path, void *wctx);
void StreamExporterIDLoader::GetStreamToValueCallback(int32_t type, int32_t state, const char *path, void *wctx)
{
    LOG_DEBUG("get streamto id list callback, type:%d, state:%d, path:%s", type, state, path);
    StreamExporterIDLoader *self = (StreamExporterIDLoader *)wctx;
    std::string streamId = self->GetStreamIDFromZkPath(path);
    if (streamId == "")
    {
        LOG_ERROR("stream id(%s) invalid, path:%s", streamId.c_str(), path);
        return;
    }

    if (type == ZK_DELETED_EVENT_DEF)
    {
        self->DeleteStreamExporterConfig(streamId);
        return;
    }

    std::string value;
    int ret = self->GetStreamExporterZkValue(path, value);
    if (ret != GSE_SUCCESS)
    {
        LOG_ERROR("failed to get stream exporter zk value, ret:%d", ret);
        return;
    }

    self->SaveStreamExporterConfig(streamId, value);

    return;
}

int StreamExporterIDLoader::GetStreamExporterZkValue(const char *path, std::string &value)
{
    int ret = m_zkClient->ZkGet(std::string(path), value, StreamExporterIDLoader::GetStreamToValueCallback, this, nullptr);
    if (ret != GSE_SUCCESS)
    {
        return GSE_ERROR;
    }

    LOG_DEBUG("get streamto(%s)'s config,value(%s)", path, value.c_str());
    return GSE_SUCCESS;
}

int StreamExporterIDLoader::DeleteStreamExporterConfig(std::string &id)
{
    uint32_t uID = gse::tools::strings::StringToUint32(id);
    uint32_t *ptrStreamID = new uint32_t(uID);
    //*ptrStreamID = uID;

    ZkEvent *event = new ZkEvent();
    event->m_eventType = ZK_EVENT_DELETE;
    event->m_msg = (void *)ptrStreamID;

    if (m_streamExporterManager->UpdateExporterConfig(event) != GSE_SUCCESS)
    {
        delete event;
        return GSE_ERROR;
    }

    return GSE_SUCCESS;
}

int StreamExporterIDLoader::SaveStreamExporterConfig(const std::string &id, const std::string &value)
{
    Json::Value configJson;
    Json::Reader reader(Json::Features::strictMode());
    if (!reader.parse(value, configJson))
    {
        LOG_DEBUG("the channel id (%s)'s config is not valid json(%s)", id.c_str(), value.c_str());
        return GSE_ERROR;
    }

    ChannelIdExporterConfig *ptrStreamToIdConfig = new ChannelIdExporterConfig();
    ApiError error;
    if (!ptrStreamToIdConfig->m_streamToCluster.Parse(configJson, error))
    {
        delete ptrStreamToIdConfig;
        LOG_ERROR("failed to parse streamto cluster config, json:%s, error:%s", value.c_str(), GET_ERROR_MESSAGE(error).c_str());
        return GSE_ERROR;
    }

    ptrStreamToIdConfig->m_streamToId = gse::tools::strings::StringToUint32(id);

    ZkEvent *event = new ZkEvent();
    event->m_eventType = ZK_EVENT_CHANGE;
    event->m_msg = (void *)ptrStreamToIdConfig;

    if (m_streamExporterManager->UpdateExporterConfig(event) != GSE_SUCCESS)
    {
        delete event;
        delete ptrStreamToIdConfig;
        return GSE_ERROR;
    }

    return GSE_SUCCESS;
}

int StreamExporterIDLoader::LoadStreamExporterConfig()
{
    std::vector<std::string> idList;
    std::string path = ZK_STREAM_ID_CONFIG_BASE_PATH;

    LOG_DEBUG("streamto id's zk path:%s", path.c_str());

    int ret = m_zkClient->ZkGetChildren(path, StreamExporterIDLoader::GetStreamToIdListCallback, this, idList, nullptr);
    if (ret != GSE_SUCCESS)
    {
        LOG_WARN("failed to get streamto id's list, zk path:%s", path.c_str());
        // TODO: exist
        return ret;
    }

    auto idIter = idList.begin();
    for (; idIter != idList.end(); idIter++)
    {
        std::string idStr = (*idIter);
        if (idStr.compare("index") == 0)
        {
            continue;
        }

        if (!gse::tools::strings::IsNumber(idStr))
        {
            LOG_WARN("stream id(%s) is not number, ignore it", idStr.c_str());
            continue;
        }

        auto intId = gse::tools::strings::StringToUint32(idStr);
        if (m_streamExporterManager->Find(intId))
        {
            LOG_DEBUG("stream id(%s) exist, don't update it", idStr.c_str());
            continue;
        }

        std::string idZkPath = path + "/" + idStr;
        LOG_DEBUG("read the stream id config from zk node :%s", idZkPath.c_str());
        std::string value;
        if (GetStreamExporterZkValue(idZkPath.c_str(), value))
        {
            LOG_ERROR("failed to get stream exporter value, path:%s", idZkPath.c_str());
            continue;
        }

        SaveStreamExporterConfig(idStr, value);
    }

    return GSE_SUCCESS;
}

void StreamExporterIDLoader::SetStreamExporterManager(std::shared_ptr<ChannelIdStreamExporterManager> manger)
{
    m_streamExporterManager = manger;
}

} // namespace data
} // namespace gse
