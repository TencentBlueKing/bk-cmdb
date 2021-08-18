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

#include "gsedata_package_v1.h"
#include "log/log.h"
#include "tools/endian.h"

#include "open/protocol_head.h"

namespace gse { 
namespace dataserver {

GSEDataPackageV1::GSEDataPackageV1()
{
    m_ptrValue = (char*)malloc(GSE_DATA_MESSAGE_BLOCK_SIZE);
    memset(m_ptrValue,0, GSE_DATA_MESSAGE_BLOCK_SIZE);
    m_valueLength = GSE_DATA_MESSAGE_BLOCK_SIZE;
    m_valuePosition = GSE_DATA_PUB_HEAD_SIZE;
    m_channelId = 0;
}

GSEDataPackageV1::~GSEDataPackageV1()
{
    if(NULL != m_ptrValue)
    {
        free(m_ptrValue);
    }
}

char* GSEDataPackageV1::GetDataPointer()
{
    return m_ptrValue;
}

uint32_t GSEDataPackageV1::GetDataLength()
{
    return m_valuePosition;
}



void GSEDataPackageV1::SetChannelId(uint32_t channelid)
{
    m_channelId = channelid;
}

void GSEDataPackageV1::ResetBufferPosition()
{
    m_valuePosition = GSE_DATA_PUB_HEAD_SIZE;
    memset(m_ptrValue, 0, m_valueLength);
}

void GSEDataPackageV1::tryReallocBuffer(uint32_t targetSize)
{
    while(targetSize > m_valueLength)
    {
        m_ptrValue = (char*)realloc(m_ptrValue, targetSize);
        if(NULL == m_ptrValue)
        {
            LOG_FATAL("it is failed to realloc the new memory for the gse data package size %u", targetSize);
            return;
        }
        m_valueLength += targetSize;
    }
}

int GSEDataPackageV1::CalcMsgLen(uint32_t data_len)
{
    int total_len = 0;

    if (data_len > 0)
    {
        total_len += data_len + sizeof(TagElement);
    }

    total_len += sizeof(DataMsgHeadV1);

    return total_len;
}
void GSEDataPackageV1::Pack(const char* ptr_data, uint32_t data_len)
{
    int total_len = CalcMsgLen(data_len);
    tryReallocBuffer(total_len);

    DataMsgHeadV1 *ptr_head = (DataMsgHeadV1 *)m_ptrValue;
    ptr_head->m_msgtype = gse::tools::endian::HostToNetwork32(0);
    ptr_head->m_channelid = gse::tools::endian::HostToNetwork32(m_channelId);

    int tag_offset = 0;

    TagElement *ptr_tag = nullptr;
    ptr_tag = (TagElement *)ptr_head->m_data;
    int tag_len = 0;

    ptr_tag = (TagElement *)ptr_head->m_data;
    ptr_tag->m_tag = gse::tools::endian::HostToNetwork32(enMsgContent);
    ptr_tag->m_len = gse::tools::endian::HostToNetwork32(data_len);
    memcpy(ptr_tag->m_value, ptr_data, data_len);
    tag_offset += tag_len;
    ptr_head->m_msglen = gse::tools::endian::HostToNetwork32(tag_offset);

    return;
}

}//
}//
