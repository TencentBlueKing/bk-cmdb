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

#include "dataflow.h"
#include "conf/dataconf.h"
#include "conf/configurator.h"
#include "opscollection/ops.h"

namespace gse { 
namespace dataserver {

class DataServer
{
public:
    virtual ~DataServer();
    DataServer();

public:
    static DataServer& Instance();
public:
    int Start(DataProcessConfig& conf, bool upgrade);
    int Stop();
    void Join();


public:
    OpsCollection *GetOpsReportClient();
    static bool GetUpgradeFlag();
private:

    DataServer(const DataServer& src);
    DataServer& operator=(const DataServer& src);
    void StartOpsReportClient(DataFlowConf *m_dataFlow, DataProcessConfig &gse_cfg);

private:
    DataFlow* m_dataFlow;
    Configurator* m_conftor;
    static bool m_upgradeFlag;
    OpsCollection *m_opsReport;
};

}
}
#endif
