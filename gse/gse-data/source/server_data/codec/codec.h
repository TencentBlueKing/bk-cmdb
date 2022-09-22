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

#ifndef _GSE_DATA_CODEC_H_
#define _GSE_DATA_CODEC_H_

#include "datacell.h"
//#include "ops/ops.h"

namespace gse {
namespace data {

class OpsCollection;

class Codec
{
public:
    Codec();
    virtual ~Codec();

public:
    inline void SetOPS(OpsCollection* ptrOPSReport)
    {
        m_ptrOPSReport = ptrOPSReport;
    }

public:
    virtual int DecodeMsg(DataCell* pDataCell) = 0;
    virtual int EncodeMsg(DataCell* pDataCell) = 0;

protected:
    OpsCollection* m_ptrOPSReport;
};

} // namespace data
} // namespace gse
#endif //_GSE_DATA_CODEC_H_
