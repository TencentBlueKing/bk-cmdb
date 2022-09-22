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

#include <string>

#include "gsedata_package.h"
#include "log/log.h"
#include "tools/endian.h"
namespace gse { 
namespace data {

GSEDataPackage::GSEDataPackage()
{
    m_ptrValue = (char*)malloc(GSE_DATA_MESSAGE_BLOCK_SIZE);
    memset(m_ptrValue,0, GSE_DATA_MESSAGE_BLOCK_SIZE);
    m_valueLength = GSE_DATA_MESSAGE_BLOCK_SIZE;
    m_valuePosition = GSE_DATA_PUB_HEAD_SIZE;
}

GSEDataPackage::~GSEDataPackage()
{
    if(NULL != m_ptrValue)
    {
        free(m_ptrValue);
    }
}

char* GSEDataPackage::GetDataPointer()
{
    return m_ptrValue;
}

uint32_t GSEDataPackage::GetDataLength()
{
    return m_valuePosition;
}

void GSEDataPackage::SetType(uint16_t type)
{
    GseDataPubHead* ptr_pub_head = (GseDataPubHead*)(m_ptrValue);
    ptr_pub_head->m_extType = gse::tools::endian::HostToNetwork16(type);
}
void GSEDataPackage::SetVersion(uint16_t version)
{
    GseDataPubHead* ptr_pub_head = (GseDataPubHead*)(m_ptrValue);
    ptr_pub_head->m_extVersion = gse::tools::endian::HostToNetwork16(version);
}
void GSEDataPackage::SetMagicNumber(uint32_t magic)
{
    GseDataPubHead* ptr_pub_head = (GseDataPubHead*)(m_ptrValue);
    ptr_pub_head->m_magicNum = gse::tools::endian::HostToNetwork32(magic);
}

void GSEDataPackage::ResetBufferPosition()
{
    m_valuePosition = GSE_DATA_PUB_HEAD_SIZE;
    memset(m_ptrValue, 0, m_valueLength);
}

/*

std::string id = gse::tools::strings::ToString(cnt_id++); // TODO: 为tglog 实时数据添加的数据id头
std::string timestamp = gse::tools::strings::ToString(pDataCell->GetCreationTimestamp());
m_extension.ResetBufferPosition();
if (m_fillChannelid)
{
    m_extension.SetChannelID(pDataCell->GetChannelID());
}
else
{
    m_extension.SetChannelID(0);
}

m_extension.SetTimestamp(pDataCell->GetArrivedTimestamp());
m_extension.AddExtension(id.c_str(), id.size());
m_extension.AddExtension(timestamp.c_str(), timestamp.size());
std::size_t max_count = m_extensions.size();
for(std::size_t idx = 0; idx < max_count; ++idx)
{
    std::string value = m_extensions.at(idx);
    if(value.empty())
    {
        continue;
    }
    m_extension.AddExtension(value.c_str(), value.length());
    //ptr_ops->m_extensionsInfo.push_back(value);
    pDataCell->PushExtension(value);
}

// construct the data package
m_gseDataPackage.ResetBufferPosition();
m_gseDataPackage.SetType(DS_EXT_TYPE);
m_gseDataPackage.SetVersion(DS_EXT_VERSION);
m_gseDataPackage.SetMagicNumber(DS_MAGIC_NUM);
m_gseDataPackage.SetBody(pDataCell->GetDataBuf(), pDataCell->GetDataBufLen(), &m_extension);

*/


void GSEDataPackage::SetId(std::string &id)
{
    m_cntId = id;
}

void GSEDataPackage::SetExterntions(std::vector<std::string>& externstions)
{
    m_extensions = externstions;
}

void GSEDataPackage::SetChannelId(uint32_t channelid)
{
    m_channelId = channelid;
}

void GSEDataPackage::SetCreateTimestamp(uint32_t timestamp)
{
    m_createTimestamp = timestamp;
}

void GSEDataPackage::SetArrivedTimestamp(uint32_t timestamp)
{
    m_arrivedTimeStamp = timestamp;
}

void GSEDataPackage::Pack(const char* ptr_data, uint32_t data_len)
{
    GSEDataPackageExtension  m_extension;
    std::string timestamp = gse::tools::strings::ToString(m_createTimestamp);

    m_extension.ResetBufferPosition();
    m_extension.SetChannelID(m_channelId);
    m_extension.SetTimestamp(m_arrivedTimeStamp);
    m_extension.AddExtension(m_cntId.c_str(), m_cntId.size());
    m_extension.AddExtension(timestamp.c_str(), timestamp.size());
    std::size_t max_count = m_extensions.size();

    for(std::size_t idx = 0; idx < max_count; ++idx)
    {
        std::string value = m_extensions.at(idx);
        if(value.empty())
        {
            continue;
        }
        m_extension.AddExtension(value.c_str(), value.length());
        //pDataCell->PushExtension(value);
    }

    // construct the data package
    ResetBufferPosition();
    SetType(DS_EXT_TYPE);
    SetVersion(DS_EXT_VERSION);
    SetMagicNumber(DS_MAGIC_NUM);
    SetBody(ptr_data, data_len, &m_extension);

    return;
}


void GSEDataPackage::SetBody(const char* ptrValue, uint32_t valueLength, GSEDataPackageExtension* ptrExtension /*= NULL*/)
{
    if(NULL != ptrExtension)
    {
        setExtension(ptrExtension);
    }

    uint32_t need_total_size = m_valuePosition + valueLength;
    uint32_t data_length = need_total_size - (uint32_t)GSE_DATA_PUB_HEAD_SIZE;
    tryReallocBuffer(need_total_size);

    GseDataPubHead* ptr_pub_head = (GseDataPubHead*)(m_ptrValue);
    ptr_pub_head->m_bodyLength = gse::tools::endian::HostToNetwork32(data_length);
    LOG_DEBUG("all data package length is %d, the data length %d , the buffer already has length %d", need_total_size, ntohl(ptr_pub_head->m_bodyLength),m_valueLength);
    memcpy(m_ptrValue + m_valuePosition, ptrValue, valueLength);

    m_valuePosition = need_total_size;
}

void GSEDataPackage::tryReallocBuffer(uint32_t targetSize)
{
    while(targetSize > m_valueLength)
    {
        m_ptrValue = (char*)realloc(m_ptrValue, targetSize);
        if(NULL == m_ptrValue)
        {
            LOG_FATAL("it is failed to realloc the new memory for the gse data package size %u", targetSize);
            // TODO: maybe need a way to recover
            abort();
        }
        m_valueLength += targetSize;
    }
}

void GSEDataPackage::setExtension(GSEDataPackageExtension* ptrExtension)
{
    char* ptr_extension_data = ptrExtension->GetDataPointer();
    uint32_t extension_data_size = ptrExtension->GetDataLength();

    uint32_t need_total_size =  m_valuePosition + extension_data_size;
    LOG_DEBUG("extension data length is %d", extension_data_size);
    tryReallocBuffer(need_total_size);

    memcpy(m_ptrValue + m_valuePosition, ptr_extension_data, extension_data_size);
    m_valuePosition += extension_data_size;
}

}
}
