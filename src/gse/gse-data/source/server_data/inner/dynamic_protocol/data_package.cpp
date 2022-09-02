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

#include <arpa/inet.h>

#include "data_package.h"
#include "log/log.h"
namespace gse { 
namespace data {


DataPackage::DataPackage()
{
    m_ptrValue = (char*)malloc(GSE_DATA_MESSAGE_BLOCK_SIZE);
    memset(m_ptrValue,0, GSE_DATA_MESSAGE_BLOCK_SIZE);
    m_valueLength = GSE_DATA_MESSAGE_BLOCK_SIZE;
    m_valuePosition = GSE_DATA_PUB_HEAD_SIZE;
}

DataPackage::~DataPackage()
{
    if(NULL != m_ptrValue)
    {
        free(m_ptrValue);
    }
}

char* DataPackage::GetDataPointer()
{
    return m_ptrValue;
}

uint32_t DataPackage::GetDataLength()
{
    return m_valuePosition;
}

void DataPackage::SetType(uint16_t type)
{
    GseDataPubHead* ptr_pub_head = (GseDataPubHead*)(m_ptrValue);
    ptr_pub_head->m_extType = htons(type);
}
void DataPackage::SetVersion(uint16_t version)
{
    GseDataPubHead* ptr_pub_head = (GseDataPubHead*)(m_ptrValue);
    ptr_pub_head->m_extVersion = htons(version);
}
void DataPackage::SetMagicNumber(uint32_t magic)
{
    GseDataPubHead* ptr_pub_head = (GseDataPubHead*)(m_ptrValue);
    ptr_pub_head->m_magicNum = htonl(magic);
}

void DataPackage::ResetBufferPosition()
{
    m_valuePosition = GSE_DATA_PUB_HEAD_SIZE;
    memset(m_ptrValue, 0, m_valueLength);
}

bool DataPackage::SetBody(const char* ptrValue, uint32_t valueLength, DataPackageExtension* ptrExtension /*= NULL*/)
{
    if(NULL != ptrExtension)
    {
        setExtension(ptrExtension);
    }

    uint32_t need_total_size = m_valuePosition + valueLength;
    uint32_t data_length = need_total_size - (uint32_t)GSE_DATA_PUB_HEAD_SIZE;
    if(!tryReallocBuffer(need_total_size))
    {
        return false;
    }

    GseDataPubHead* ptr_pub_head = (GseDataPubHead*)(m_ptrValue);
    ptr_pub_head->m_bodyLength = htonl(data_length);
    LOG_DEBUG("all data package length is %d, the data length %d , the buffer already has length %d", need_total_size, ntohl(ptr_pub_head->m_bodyLength),m_valueLength);
    memcpy(m_ptrValue + m_valuePosition, ptrValue, valueLength);

    m_valuePosition = need_total_size;
    return true;
}

bool DataPackage::tryReallocBuffer(uint32_t targetSize)
{
    while(targetSize > m_valueLength)
    {
        m_ptrValue = (char*)realloc(m_ptrValue, targetSize);
        if(NULL == m_ptrValue)
        {
            LOG_FATAL("it is failed to realloc the new memory for the gse data package size %u", targetSize);
            // TODO: maybe need a way to recover
            return false;
        }
        m_valueLength += targetSize;
    }
    return true;
}

bool DataPackage::setExtension(DataPackageExtension* ptrExtension)
{
    char* ptr_extension_data = ptrExtension->GetDataPointer();
    uint32_t extension_data_size = ptrExtension->GetDataLength();
    if(NULL == ptr_extension_data)
    {
        // Note:仅当内存分配失败会出现此异常
        LOG_ERROR("the extension memory is invalid(null)");
        return false;
    }

    uint32_t need_total_size =  m_valuePosition + extension_data_size;
    LOG_DEBUG("extension data length is %d", extension_data_size);
    if(!tryReallocBuffer(need_total_size))
    {
        return false;
    }

    memcpy(m_ptrValue + m_valuePosition, ptr_extension_data, extension_data_size);
    m_valuePosition += extension_data_size;

    return true;
}

}
}
