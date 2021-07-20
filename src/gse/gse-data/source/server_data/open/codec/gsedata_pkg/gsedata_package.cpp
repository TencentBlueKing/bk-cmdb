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

#include "gsedata_package.h"
#include "log/log.h"
#include "tools/endian.h"

#include "open/protocol_head.h"

namespace gse { 
namespace dataserver {

GSEDataPackage::GSEDataPackage()
{
    m_ptrValue = (char*)malloc(GSE_DATA_MESSAGE_BLOCK_SIZE);
    memset(m_ptrValue,0, GSE_DATA_MESSAGE_BLOCK_SIZE);
    m_valueLength = GSE_DATA_MESSAGE_BLOCK_SIZE;
    m_valuePosition = GSE_DATA_PUB_HEAD_SIZE;
    m_createTimestamp = 0;
    m_arrivedTimeStamp = 0;
    m_msgLen = 0;
    m_channelId = 0;
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
    return m_msgLen;
}

void GSEDataPackage::SetId(std::string &id)
{
    m_id = id;
}

void GSEDataPackage::SetExterntions(std::vector<std::string> &externstions)
{
    m_externtions = externstions;
}

void GSEDataPackage::SetChannelId(uint32_t channelid)
{
    m_channelId = channelid;
}

void GSEDataPackage::ResetBufferPosition()
{
    m_valuePosition = GSE_DATA_PUB_HEAD_SIZE;
    memset(m_ptrValue, 0, m_valueLength);
}


void GSEDataPackage::SetCreateTimestamp(uint32_t timestamp)
{
    m_createTimestamp = timestamp;
}

void GSEDataPackage::SetArrivedTimestamp(uint32_t timestamp)
{
    m_arrivedTimeStamp = timestamp;
}

void GSEDataPackage::tryReallocBuffer(uint32_t targetSize)
{
    while(targetSize > m_valueLength)
    {
        m_ptrValue = (char*)realloc(m_ptrValue, targetSize);
        if(NULL == m_ptrValue)
        {
            LOG_FATAL("it is failed to realloc the new memory for the gse data package size %u", targetSize);
            return;
        }

        m_msgLen = targetSize;
    }
}

int GSEDataPackage::CalcMsgLen(uint32_t data_len)
{
    int total_len = 0;
    if (m_id.size() > 0)
    {
        total_len += m_id.size() + sizeof(TagElement);
    }

    if (m_externtions.size() > 0)
    {
        for (int i = 0;i < m_externtions.size(); i++)
        {
            total_len += m_externtions[i].size() + sizeof(TagLV);
        }
        total_len += sizeof(TagElement);

    }

    if (data_len > 0)
    {
        total_len += data_len + sizeof(TagElement);
    }

    total_len += sizeof(DataMsgHead);

    return total_len;
}
void GSEDataPackage::Pack(const char* ptr_data, uint32_t data_len)
{
    int total_len = CalcMsgLen(data_len);
    m_msgLen = total_len;

    tryReallocBuffer(total_len);

    DataMsgHead *ptr_head = (DataMsgHead *)m_ptrValue;
    ptr_head->m_msgtype = gse::tools::endian::HostToNetwork32(0);
    ptr_head->m_channelid = gse::tools::endian::HostToNetwork32(m_channelId);
    ptr_head->m_magic = gse::tools::endian::HostToNetwork32(MGAIC_NUM);

    int tag_offset = 0;


    TagElement *ptr_tag = nullptr;
    ptr_tag = (TagElement *)ptr_head->m_data;
    int tag_len = 0;
    ptr_tag->m_tag = gse::tools::endian::HostToNetwork32(enExtentionId);
    ptr_tag->m_len = gse::tools::endian::HostToNetwork32(m_id.size());
    memcpy(ptr_tag->m_value, m_id.c_str(), m_id.size());
    tag_len = sizeof(TagElement) + m_id.size();
    tag_offset += tag_len;

    ptr_tag = (TagElement *)(ptr_head->m_data + ptr_head->m_msglen);

    ptr_tag->m_tag = gse::tools::endian::HostToNetwork32(enExterntionInfo);
    int offset = 0;
    int extenstion_len = 0;
    for (int i = 0; i< m_externtions.size(); i++)
    {
        TagLV *lv = (TagLV *)(ptr_tag->m_value + offset);
        lv->m_len = gse::tools::endian::HostToNetwork32(m_externtions[i].size());
        memcpy(lv->m_value, m_externtions[i].c_str(), m_externtions[i].size());
        offset += m_externtions[i].size() + sizeof(TagLV);
    }

    ptr_tag->m_len = gse::tools::endian::HostToNetwork32(offset);
    tag_len = sizeof(TagElement) + offset;
    tag_offset += tag_len;

    ptr_tag = (TagElement *)(ptr_head->m_data + ptr_head->m_msglen);
    ptr_tag->m_tag = gse::tools::endian::HostToNetwork32(enMsgContent);
    ptr_tag->m_len = gse::tools::endian::HostToNetwork32(data_len);
    memcpy(ptr_tag->m_value, ptr_data, data_len);
    tag_offset += tag_len;
    ptr_head->m_msglen = gse::tools::endian::HostToNetwork32(tag_offset);

    return;
}

}//
}//
