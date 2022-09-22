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

#ifndef SOURCE_COMMON_DYNAMIC_PROTOCOL_PACKAGE_EXTENSION_H_
#define SOURCE_COMMON_DYNAMIC_PROTOCOL_PACKAGE_EXTENSION_H_

#include "inner/protocol_head.h"
namespace gse { 
namespace data {

class DataPackageExtension
{
public:
    enum
    {
        GSE_DATA_DYNAMIC_EXT_HEAD_SIZE = sizeof(GseDataDynamicExt),
        GSE_DATA_DYNAMIC_LVHEAD_SIZE = sizeof(GseDataLVHead),
        EXTENSION_BLOCK_SIZE = 512 // default 512B
    };
public:
    DataPackageExtension();
    ~DataPackageExtension();

public:
    void SetChannelID(uint32_t channelID);
    void SetTimestamp(uint32_t timestamp);
    void AddExtension(const char* ptrValue, uint16_t valueLen);
    void ResetBufferPosition();
    char* GetDataPointer();
    uint32_t GetDataLength();

private:
   char*              m_ptrValue;
   uint32_t           m_valuePosition;
   uint32_t           m_valueLength;
};

}
}

#endif
