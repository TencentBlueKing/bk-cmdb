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

#ifndef _GSE_GSEDATA_PKG_CODEC_V1_H_
#define _GSE_GSEDATA_PKG_CODEC_V1_H_

#include <map>
#include <curl/curl.h>
#include <curl/easy.h>

#include "datacell.h"
#include "codec/codec.h"
namespace gse { 
namespace dataserver {


// 兼容V1.0 版本的DS的数据通信协议
class GseDataPkgCodecV1 : public Codec
{
public:
    typedef int (GseDataPkgCodecV1::*DecodeFuncPtr)(DataCell* pDataCell);
    typedef std::map<int, DecodeFuncPtr> DecodeTypeFuncMap;

public:
    GseDataPkgCodecV1();
    virtual ~GseDataPkgCodecV1();

public:
    int DecodeMsg(DataCell *pDataCell);
    int EncodeMsg(DataCell *pDataCell);


};

}
}
#endif // _GSE_GSEDATA_PKG_CODEC_H_
