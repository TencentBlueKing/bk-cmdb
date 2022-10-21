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
namespace gse { 
namespace data {

GSEDataPackageV1::GSEDataPackageV1()
{
    m_ptrValue = (char*)malloc(GSE_DATA_MESSAGE_BLOCK_SIZE);
    memset(m_ptrValue,0, GSE_DATA_MESSAGE_BLOCK_SIZE);
    m_valueLength = GSE_DATA_MESSAGE_BLOCK_SIZE;
    m_valuePosition = GSE_DATA_PUB_HEAD_SIZE;
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

void GSEDataPackageV1::Pack(const char* ptr_data, uint32_t data_len)
{

    DynamicalProtocolHead head;
    // msghead
    head.m_msgHead.msgType = REPORT_DYNAMICAL_PROTOCOL_TYPE;

    head.m_msgHead.lMsg.dataid = m_channelId;

    head.m_msgHead.lMsg.utctime = time(NULL);
    head.m_msgHead.lMsg.resv[0] = 0;
    head.m_msgHead.lMsg.resv[1] = 0;

    DynamicalProtocol dynamicalProtocol;
    int ret = dynamicalProtocol.setDynamicalProtocolHead(head, ptr_data, data_len);
    if (ret != 0)
    {
        LOG_ERROR("set dynamical protocol head failed");
        return;
    }

    int32_t len = 0;
    char *ptrData = NULL;
    dynamicalProtocol.binaryData(ptrData, len);

    return;
}


void GSEDataPackageV1::tryReallocBuffer(uint32_t targetSize)
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

}
}
