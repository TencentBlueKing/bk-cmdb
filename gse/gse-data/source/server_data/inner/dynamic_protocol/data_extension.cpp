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

#include <string.h>
#include <math.h>
#include <arpa/inet.h>
#include "log/log.h"
#include "data_extension.h"
namespace gse { 
namespace data {


DataPackageExtension::DataPackageExtension()
{
    m_ptrValue = (char*)malloc(EXTENSION_BLOCK_SIZE);
    m_valueLength = EXTENSION_BLOCK_SIZE;
    m_valuePosition = GSE_DATA_DYNAMIC_EXT_HEAD_SIZE;
    memset(m_ptrValue, 0, m_valueLength);
}

DataPackageExtension::~DataPackageExtension()
{
   if(NULL != m_ptrValue)
   {
       free(m_ptrValue);
       m_ptrValue = NULL;
       m_valueLength = 0;
   }
}

void DataPackageExtension::SetChannelID(uint32_t channelID)
{
    if(NULL == m_ptrValue)
    {
        // TODO: 仅当 内存分配失败的时候会出现此异常
        LOG_ERROR("the extension memory is invalid(null)");
        return;
    }
    GseDataDynamicExt* ptr_dynamic_head = (GseDataDynamicExt*)(m_ptrValue);
    ptr_dynamic_head->m_channelID = htonl(channelID);
}

void DataPackageExtension::SetTimestamp(uint32_t timestamp)
{
    if(NULL == m_ptrValue)
    {
        // TODO: 仅当 内存分配失败的时候会出现此异常
        LOG_ERROR("the extension memory is invalid(null)");
        return;
    }
    GseDataDynamicExt* ptr_dynamic_head = (GseDataDynamicExt*)(m_ptrValue);
    ptr_dynamic_head->m_timestamp = htonl(timestamp);
}


void DataPackageExtension::AddExtension(const char* ptrValue, uint16_t valueLen)
{
    if(NULL == m_ptrValue)
    {
        // TODO: 仅当 内存分配失败的时候会出现此异常
        LOG_ERROR("the extension memory is invalid(null)");
        return;
    }
    GseDataDynamicExt* ptr_dynamic_head = (GseDataDynamicExt*)(m_ptrValue);

    uint32_t extension_total_size = GSE_DATA_DYNAMIC_LVHEAD_SIZE + valueLen;
    uint32_t need_size = m_valuePosition + extension_total_size;
    while( need_size > m_valueLength)
    {
        m_ptrValue = (char*)realloc(m_ptrValue, need_size);
        if(NULL == m_ptrValue)
        {
            LOG_FATAL("it is failed to realloc the new memory for the extension size %u", need_size);
            // TODO: maybe need a way to recover
            m_valueLength = 0;
            m_valuePosition = 0;
            return;
        }
        m_valueLength = need_size;
    }

    GseDataLVHead* ptr_lv_head = (GseDataLVHead*)(m_ptrValue + m_valuePosition);
    ptr_lv_head->m_valueLength = htons(valueLen);

    m_valuePosition += GSE_DATA_DYNAMIC_LVHEAD_SIZE;
    memcpy(m_ptrValue + m_valuePosition, ptrValue, valueLen);
    m_valuePosition += valueLen;
    ptr_dynamic_head->m_dynamicExtBodyLength = htonl( ntohl(ptr_dynamic_head->m_dynamicExtBodyLength) + extension_total_size);
}

void DataPackageExtension::ResetBufferPosition()
{
    if(NULL == m_ptrValue)
    {
        // TODO: 仅当 内存分配失败的时候会出现此异常
        LOG_ERROR("the extension memory is invalid(null)");
        return;
    }
    m_valuePosition = GSE_DATA_DYNAMIC_EXT_HEAD_SIZE;
    memset(m_ptrValue, 0, m_valueLength);
}

char* DataPackageExtension::GetDataPointer()
{
    return m_ptrValue;
}

uint32_t DataPackageExtension::GetDataLength()
{
    return m_valuePosition;
}
}
}
