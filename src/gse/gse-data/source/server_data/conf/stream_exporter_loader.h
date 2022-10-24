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

#ifndef _GSE_DATA_CHANNELID_LOADER_H_
#define _GSE_DATA_CHANNELID_LOADER_H_

#include "discover/zkapi/zk_api.h"

#include "channel_id_config.h"

namespace gse {
namespace data {

class StreamExporterIDLoader final
{
public:
    explicit StreamExporterIDLoader(std::shared_ptr<discover::zkapi::ZkApi> zkCli);
    StreamExporterIDLoader();
    virtual ~StreamExporterIDLoader();

    // init load config
    int LoadStreamExporterConfig();
    void SetStreamExporterManager(std::shared_ptr<ChannelIdStreamExporterManager> manger);

    // ----------

    int GetStreamExporterZkValue(const char *path, std::string &value);
    int SaveStreamExporterConfig(const std::string &id, const std::string &value);
    int DeleteStreamExporterConfig(std::string &id);

private:
    static void GetStreamToIdListCallback(int32_t type, int32_t state, const char *path, void *wctx);
    static void GetStreamToValueCallback(int32_t type, int32_t state, const char *path, void *wctx);

private:
    std::string GetStreamIDFromZkPath(const char *zkPath);

private:
    std::shared_ptr<ChannelIdStreamExporterManager> m_streamExporterManager;
    std::shared_ptr<gse::discover::zkapi::ZkApi> m_zkClient;
};

} // namespace data
} // namespace gse

#endif // CHANNELID_LOADER_H
