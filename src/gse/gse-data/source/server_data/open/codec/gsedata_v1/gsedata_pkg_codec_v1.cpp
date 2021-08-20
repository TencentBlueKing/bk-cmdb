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

#include "gsedata_pkg_codec_v1.h"
#include "protocol.h"
#include "tools/time.h"
#include "protocol.h"

#include "bbx/gse_errno.h"
#include "tools/macros.h"
#include "conf/configurator.h"

namespace gse { 
namespace dataserver {

GseDataPkgCodecV1::GseDataPkgCodecV1()
{
}

GseDataPkgCodecV1::~GseDataPkgCodecV1()
{

}

int GseDataPkgCodecV1::DecodeMsg(DataCell *pDataCell)
{
    int iRet = GSE_SUCCESS;
    if (NULL == pDataCell)
    {
        LOG_WARN("fail to decode gsedata package msg, because the args of datacell is NULL");
        return GSE_ERROR;
    }

    DataMsgHeadV1* ptr_head = (DataMsgHeadV1*)(pDataCell->GetDataBuf());
    uint32_t data_length = ntohl(ptr_head->m_msglen);
    char* ptr_data_body = ptr_head->m_data;
    uint32_t channelid = ntohl(ptr_head->m_channelid);
    LOG_DEBUG("receive channelid id [%u], the data len is [%d]", channelid, data_length);
    pDataCell->SetChannelID(channelid);
    pDataCell->CopyData(ptr_data_body, data_length);

    return GSE_SUCCESS;
}

int GseDataPkgCodecV1::EncodeMsg(DataCell *pDataCell)
{
    return GSE_SUCCESS;
}
}
}
