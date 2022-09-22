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

#ifndef _GSE_DATA_FLOW_H_
#define _GSE_DATA_FLOW_H_

#include "channel/datachannel.h"
#include "conf/dataconf.h"
#include "ops/ops.h"
#include "safe/lock.h"
#include <map>

namespace gse {
namespace data {

class DataFlow
{
public:
    DataFlow();
    virtual ~DataFlow();

    int Start();
    int Stop();
    void Join();
    void UpdateConf(DataFlowConf* pFlowConf);
    void UpdateIDToStorage(IDToStorage* ptrIDToStorage);

private:
    gse::safe::RWLock m_upConfLock;
    DataFlowConf* m_flowConf;
    std::map<std::string, DataChannel*> m_mapDataChannels;

private:
    IDToStorage* m_ptrIDToStorage;

    DataProcessConfig m_gseCfg;

    OpsCollection* m_opsReport;
};

} // namespace data
} // namespace gse
#endif
