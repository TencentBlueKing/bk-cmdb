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

#ifndef _GSE_DATA_PROCESS_H_
#define _GSE_DATA_PROCESS_H_

#include "core/linux/deamon_set.h"
#include "dataserver.h"
#include <memory>
#include <string>
namespace gse {
namespace data {

class DataService : public gse::core::DeamonSet
{
public:
    DataService();
    virtual ~DataService();

    void ServiceStart(int id);
    void ServiceStop();
    void ServiceHandle();
    void Reconfigure();
    void AddPrivateOptions(po::options_description &config);
    void HandlePrivateFlags(po::variables_map &vm);
    const std::string Version();

protected:
private:
    std::shared_ptr<DataServer> m_serverptr;
    std::string m_logPath;
    std::string m_logLevel;
    std::string m_configFile;
};

} // namespace data
} // namespace gse
#endif
