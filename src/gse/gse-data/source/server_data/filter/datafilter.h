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

#ifndef _DATA_FILTER_H_
#define _DATA_FILTER_H_

#include "datacell.h"
#include <string>
#include <vector>
namespace gse { 
namespace data {

class DataFilterCfg
{
public:
    DataFilterCfg(): m_fieldIndex(0){};
    DataFilterCfg& operator = (const DataFilterCfg& cfg)
    {
        if ( this == &cfg ) 
        {
            return *this;
        }

        this->m_delimiter = cfg.m_delimiter;
        this->m_fieldIndex = cfg.m_fieldIndex;
        this->m_matchingWord = cfg.m_matchingWord;

        return *this;
    }
    ~DataFilterCfg(){};

public:
    string m_delimiter;
    int m_fieldIndex;
    string m_matchingWord;
};

class IDataFilter
{
public:
    /**
    brief:          filter data.
    @param	[in]    in          source data.
    @param	[in]    out         filtered data.
    @return         void
    @exception      none
    */
    virtual void Filtering(const vector<DataCell*>& in, vector<DataCell*>& out) = 0;

    /**
    brief:          judge if data is filtered by certain conditions.
    @param	[in]    data        source data.
    @return         bool        true: data is filtered(dropped), false: data is unfiltered(picked)
    @exception      none
    */
    virtual bool IsFiltered(DataCell* data) = 0;

public:
    IDataFilter(){}
    virtual ~IDataFilter(){}
};

}
}
#endif
