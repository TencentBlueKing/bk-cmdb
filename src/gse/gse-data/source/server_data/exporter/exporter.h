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

#ifndef _GSE_DATA_EXPORTER_H_
#define _GSE_DATA_EXPORTER_H_

#include <string>
#include "safe/lock.h"
#include "dataStruct/safe_vector.h"
#include "datacell.h"
#include "conf/confItem.h"
#include "opscollection/ops.h"
namespace gse { 
namespace dataserver {

class ExporterConfigWrapper
{
public:
    typedef enum ExporterConfigTypeEnum_
    {
        UNKNOWN = -1,
        DataFlowConfType = 0,
        DataIDConfType = 1,
        ChannelIDConfType = 2

    } ExporterConfigTypeEnum;

public:
    ExporterConfigWrapper(ExporterConfigTypeEnum type, void *conf)
    {
        m_exporterConfTypeEnum = type;
        switch (m_exporterConfTypeEnum)
        {
//        case ExporterConfigWrapper::ChannelIDConfType:
//            m_conf.m_ptrChannelIDConfig = reinterpret_cast<ChannelIDStorage*>(conf);
//            break;
        case ExporterConfigWrapper::DataIDConfType:
            m_conf.m_ptrDataIDConfig = reinterpret_cast<StorageConfigType*>(conf);
            break;
        case ExporterConfigWrapper::DataFlowConfType:
            m_conf.m_exporterConf = new ExporterConf(*(reinterpret_cast<ExporterConf*>(conf)));
            break;
        case ExporterConfigWrapper::ChannelIDConfType:
            m_conf.m_ptrChannelIdExporterConfig = new ChannelIdExporterConfig(*(reinterpret_cast<ChannelIdExporterConfig*>(conf)));
            break;
        default:
            m_conf.m_exporterConf = NULL;
            m_conf.m_ptrDataIDConfig = NULL;
            m_conf.m_ptrChannelIdExporterConfig = NULL;
            break;
        }
    }

    ~ExporterConfigWrapper()
    {
        switch (m_exporterConfTypeEnum)
        {
        case ExporterConfigWrapper::ChannelIDConfType:
            if (m_conf.m_ptrChannelIdExporterConfig != NULL)
            {
                delete m_conf.m_ptrChannelIdExporterConfig;
                m_conf.m_ptrChannelIdExporterConfig = NULL;
            }
            break;
        case ExporterConfigWrapper::DataIDConfType:
            if (m_conf.m_ptrDataIDConfig != NULL)
            {
                delete m_conf.m_ptrDataIDConfig;
                m_conf.m_ptrDataIDConfig = NULL;
            }
            break;
        case ExporterConfigWrapper::DataFlowConfType:
            if (m_conf.m_exporterConf != NULL)
            {
                delete m_conf.m_exporterConf;
                m_conf.m_exporterConf = NULL;
            }
            break;
        default:
            break;
        }
    }

private:
    //DISALLOW_COPY_AND_ASSIGN(ExporterConfigWrapper);

public:
    ExporterConfigTypeEnum m_exporterConfTypeEnum;

    union {
        // dataflow 本地配置
        ExporterConf *m_exporterConf;
        // dataid 关联的存储配置
        StorageConfigType *m_ptrDataIDConfig;
        // channelid 关联的存储配置
        //ChannelIDStorage *m_ptrChannelIDConfig;
        ChannelIdExporterConfig *m_ptrChannelIdExporterConfig;
    } m_conf;
};

class Exporter
{
public:
    Exporter();
    virtual ~Exporter();

public:
    virtual int Start() = 0;
    virtual int Stop() = 0;
    virtual int Write(DataCell *pDataCell) = 0;


public:
    inline void SetGseConf(DataProcessConfig& cfg)
    {
        m_gseCfg = cfg;
    }

    inline bool IsExistsAndUpdateFilter(const std::string& storageName, ChannelIDFilter* ptrChannelIDFilter)
    {
        LOG_DEBUG("check the exporter (%s) with the storage (%s)", m_name.c_str(), storageName.c_str());
        if(m_name.compare(storageName) == 0)
        {
            AddFilter(ptrChannelIDFilter);
            return true;
        }
        else if(NULL != m_next)
        {
            return m_next->IsExistsAndUpdateFilter(storageName, ptrChannelIDFilter);
        }
        return false;
    }

public:
    int WriteTo(DataCell *pDataCell);
    void UpdateConf(ExporterConfigWrapper::ExporterConfigTypeEnum type, void *conf);
    void SetNext(Exporter *ptrNext);
    void SetNeedStop();
    bool IsNeedStop();
    std::string GetExporterName();
    void SetExporterId(int id);
    int GetExporterId();

public:
    bool Filter(DataCell *pDataCell);
    bool FilterByDataID(DataCell *pDataCell);
    bool FilterByChannelID(DataCell *pDataCell);
    bool FilterByDataFlow(DataCell *pDataCell);
    void AddFilter(ChannelIDFilter* ptrChannelIDFilter);

protected:
    gse::safe::RWLock m_upConfLock;
    ExporterConfigWrapper *m_ptrConfWrapper;
    std::string m_name;
    gse::datastruct::SafeVector<ChannelIDFilter*>  m_filters;

protected:
    OpsCollection*  m_ptrOPSReport;
protected:
    DataProcessConfig  m_gseCfg;
private:
    // 一个Exporter 关联一个存储集群，
    // 当一个ID 的数据需要向多个存储发数据的时候就需要将多个exporter 组织成职责链
    Exporter *m_next;
    int      m_setStopTimestamp;
    int m_id;
};

}
}
#endif
