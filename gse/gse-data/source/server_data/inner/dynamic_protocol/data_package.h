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

#ifndef SOURCE_COMMON_DYNAMIC_PROTOCOL_PACKAGE_H_
#define SOURCE_COMMON_DYNAMIC_PROTOCOL_PACKAGE_H_

#include <string.h>

#include "inner/dynamic_protocol/data_extension.h"

namespace gse { 
namespace data {

#ifndef DS_MAGIC_NUM
#define DS_MAGIC_NUM   0xdeadbeef
#endif

#ifndef DS_EXT_TYPE
#define DS_EXT_TYPE 0
#endif

#ifndef DS_EXT_VERSION
#define DS_EXT_VERSION 1
#endif

class DataPackage
{
public:
    enum
    {
       GSE_DATA_PUB_HEAD_SIZE = sizeof(GseDataPubHead),
       GSE_DATA_MESSAGE_BLOCK_SIZE = (1024 + 512) // default: each message block size is 1526B
    };
public:
    DataPackage();
    ~DataPackage();

public:
    void SetType(uint16_t type);
    void SetVersion(uint16_t version);
    void SetMagicNumber(uint32_t magic);
    bool SetBody(const char* ptrValue, uint32_t valueLength, DataPackageExtension* ptrExtension = nullptr);
    void ResetBufferPosition();
    char* GetDataPointer();
    uint32_t GetDataLength();

private:
    bool setExtension(DataPackageExtension* ptrExtension);
    bool tryReallocBuffer(uint32_t targetSize);
private:
    char*                     m_ptrValue;
    uint32_t                  m_valueLength;
    uint32_t                  m_valuePosition;
};

}
}
#endif
