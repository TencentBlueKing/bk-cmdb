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

#ifndef _PROTOCOL_HEAD_H_
#define _PROTOCOL_HEAD_H_

#include <stdint.h>

namespace gse {
namespace data {
//老版本的消息类型
// plugins and agent msg type
#define REPORT_REDUNDANCY 0x01
#define REPORT_EXT 0x02
#define REPORT_NOT_REDUNDANCY 0x03
#define REPORT_DYNAMICAL_PROTOCOL_TYPE 0x09

#define DEFAULT_MONITOR_TAG_NAME "tag"

#define OPS_ERROR_SKIP 254
#define OPS_ERROR_TRANSFER 255

//新的通用消息类型

#define MSG_CLASS_BIT 30
#define MAX_CLASS_MAX 0xC0000000

#define OPS_MSG_TYPE (0)
#define COMMON_REPORT_TYPE (1 << MSG_CLASS_BIT)
#define CONTROL_MSG_TYPE (2 << MSG_CLASS_BIT)

#define GET_MSG_CLASS(type) ((MAX_CLASS_MAX & type) >> MSG_CLASS_BIT)

#define DATA_GENNERAL_OPS_MSG_TYPE_HEART_ALARM (OPS_MSG_TYPE + 1)  // OPS 心跳失联告警消息
#define DATA_GENNERAL_OPS_MSG_TYPE_DATA_FLOW (OPS_MSG_TYPE + 2)    //数据流水
#define DATA_GENNERAL_OPS_MSG_TYPE_CLUSTER_FLOW (OPS_MSG_TYPE + 3) // CLUSTER 流水数据
#define DATA_GENNERAL_OPS_MSG_TYPE_FILE_FLOW (OPS_MSG_TYPE + 4)    // FILE流水数据
#define DATA_GENNERAL_OPS_MSG_TYPE_TASK_FLOW (OPS_MSG_TYPE + 5)    // TASK流水数据
#define DATA_GENNERAL_OPS_MSG_TYPE_END (OPS_MSG_TYPE + 6)

#define DATA_ACCESS_REQUEST (CONTROL_MSG_TYPE + 1)
#define DATA_ACCESS_RESPONSE (CONTROL_MSG_TYPE + 2)
#define DATA_KEEPALIVE_REQUEST (CONTROL_MSG_TYPE + 3)
#define DATA_KEEPALIVE_RESPONSE (CONTROL_MSG_TYPE + 3)

#define DATA_GENNERAL_COMMON_DATA_REPORT (COMMON_REPORT_TYPE + 1)

#pragma pack(1)

typedef struct Gse_LocalRport_Msg
{
    uint32_t dataid;
    uint32_t utctime;
    uint32_t bodylen;

    //
    //
    // resv[0] is used to save report customer type,GSE_NEEDED_BASIC_* , high 8bit is used to save main type, low 24bit is used to save subtype
    // resv[1] is used to save platid and bizid, high 10bit is used to save platid, low 22bit is used to save bizid
    //
    uint32_t resv[2];

    Gse_LocalRport_Msg()
    {
        dataid = -1;
        bodylen = 0;
        utctime = 0;
        resv[0] = 0;
        resv[1] = 0;
    }
} GseLocalReportMsg;

typedef struct Gse_Report_Msg
{
    uint32_t msgType;
    GseLocalReportMsg lMsg;

    Gse_Report_Msg()
    {
        msgType = -1;
    }
} GseReportMsg;

typedef struct Gse_Report_Ext
{
    uint64_t gseExtId;
    int32_t extType;
    uint32_t bodylen;
    Gse_Report_Ext()
    {
        gseExtId = -1;
        extType = -1;
        bodylen = 0;
    }
} GseReportExt;

typedef struct Gse_Local_Command_Msg
{
    int32_t type;
    uint32_t bodylen;

    Gse_Local_Command_Msg()
    {
        type = 0;
        bodylen = 0;
    }
} GseLocalCommandMsg;

typedef struct Gse_Report_Msg_Ext
{
    GseReportMsg m_msg;
    GseReportExt m_ext;
} GseReportMsgExt;

struct GsePluginReportMsg
{
    uint32_t magic;
    uint32_t id;
    uint32_t flag;
    uint32_t len;
    uint64_t resv;
};

//
// for dynamical protocol
//

typedef struct _DynamicalProtocolHeadExt
{

    uint64_t m_extIndex;
    int32_t m_extType;
    int32_t m_extMetaLen;
    int32_t m_extMetaMaxLen;
    int32_t m_extMetaCount;

    _DynamicalProtocolHeadExt()
        : m_extIndex(0), m_extType(0), m_extMetaLen(0), m_extMetaMaxLen(0), m_extMetaCount(0)
    {
    }

} DynamicalProtocolHeadExt;

typedef struct _DynamicalMetaHead
{

    int32_t m_extMetaKeyLen;
    int32_t m_extMetaValueLen;

    _DynamicalMetaHead()
        : m_extMetaKeyLen(0), m_extMetaValueLen(0)
    {
    }

} DynamicalMetaHead;

typedef struct _DynamicalProtocolHead
{

    GseReportMsg m_msgHead;
    DynamicalProtocolHeadExt m_extHead;

} DynamicalProtocolHead;

//
// for beacon or glog
//
typedef struct _ThirdPartyHead
{
    int m_magic;
    int m_msgtype;
    int m_length;

} ThirdPartyHead;

typedef struct _ThirdPartyHeadExt
{
    uint32_t m_destid;
    uint32_t m_utctime;
    uint32_t m_resv;
} ThirdPartyHeadExt;

typedef struct Gse_Transpond_Msg
{
    uint32_t msgtype;
    uint32_t magic;
    uint32_t utctime;
    uint32_t bodylen;
    uint32_t destid;
    uint32_t resv;
} GseTranspondMsg;

typedef struct Gse_DataIdStatistic
{
    volatile uint64_t total;
    volatile uint64_t sended;
    volatile uint64_t losted;
    volatile uint64_t dataid;
    Gse_DataIdStatistic()
        : total(0), sended(0), losted(0)
    {
    }
} GseDataIdStatistic;
// common msg type under 0x100
#define MSG_TYPE_COMM_START 0x0
#define MSG_TYPE_COMM_END 0xff

#define MSG_TYPE_HEART (uint16_t)(MSG_TYPE_COMM_START + 1)
#define MSG_TYPE_HBRPY (uint16_t)(MSG_TYPE_COMM_START + 2)

#define MSG_TYPE_STATE_START 0xc00
#define MSG_TYPE_STATE_END 0xeff

// state msg type
#define MSG_DATA_REPORT (uint16_t)(MSG_TYPE_STATE_START + 1)
#define MSG_DATA_REPORT_OPS (uint16_t)(MSG_TYPE_STATE_START + 12)

//////////////////////////////////////////////////////////////
//////////          common macro define          /////////////
//////////////////////////////////////////////////////////////

#define MSG_MAX_LEN (64 << 20)

#define MSG_CHECK_NUM ((uint32_t)0xdeadbeef)
#define MSG_CHECK_NUM16 ((uint16_t)0xbeef)

// data msg
//// the msg public head of the gse data
typedef struct Gse_Data_Pub_Head
{
    uint16_t m_extType;
    uint16_t m_extVersion;
    uint32_t m_magicNum;
    uint32_t m_bodyLength;

} GseDataPubHead;
//// the expansion head of BKData
typedef struct Gse_Data_HeadExt
{
    uint32_t m_channelID;
    uint32_t m_timestamp;
} GseDataHeadExt;

//// the dynamic data expansion
typedef struct Gse_Data_DynamicExt_
{
    uint32_t m_channelID;
    uint32_t m_timestamp;
    uint32_t m_dynamicExtBodyLength;
} GseDataDynamicExt;

//// the LV expansion of the dynamic data
typedef struct Gse_Data_LVHead_
{
    uint16_t m_valueLength;
} GseDataLVHead;

//// the msg head of BKData
typedef struct Gse_Data_Head
{
    GseDataPubHead m_pubHead;
    GseDataDynamicExt m_dataHeadExt;
} GseDataHead;

// GseDataTglogMsgHead 用于兼容 原 数据链路（tdbank->tglog->datamore->transit server)推送的数据
typedef struct GseDataTglogMsgHead_
{
    uint32_t m_magic;
    uint32_t m_msgtype;
    uint32_t m_length;
    uint32_t m_destid;
    uint32_t m_utctime;
    uint32_t m_resv;

} GseDataTglogMsgHead;

#ifndef TAG_CHANNEL_ID
#define TAG_CHANNEL_ID 0x01
#endif

#ifndef TAG_UTC_TIMESTAMP
#define TAG_UTC_TIMESTAMP 0x02
#endif

#ifndef TAG_MSG_CONTENT
#define TAG_MSG_CONTENT 0x03
#endif

typedef struct DataGenneralProtocol_
{
    uint16_t m_protver;
    uint32_t m_msgtype;
    uint32_t m_magic;
    uint32_t m_msglen;
    char m_data[0];
} DataGenneralProtocol;

typedef struct TagStruct_
{
    uint32_t m_tag;
    uint32_t m_len;
} TagStruct;

typedef struct TLVMsg_
{
    TagStruct m_tagHead;
    char m_data[0];
} TLVMsg;

typedef struct TLVChannelID_
{
    TagStruct m_tagHead;
    uint32_t m_channel_id;
} TLVChannelID;

typedef struct TLVTimeStamp_
{
    TagStruct m_tagHead;
    uint32_t m_timestamp;
} TLVTimeStamp;

typedef struct TLVString_
{
    uint32_t m_stringLen;
    char m_str[0];
} TLVString;

typedef struct TLVKeyValue_
{
    TagStruct m_tagHead;
    TLVString m_key;
    TLVString m_value;
} TLVKeyValue;

typedef struct TLVTUtcTime_
{
    TagStruct m_tagHead;
    uint32_t m_utcTime;
} TLVTUtcTime;

typedef struct TLVBizId_
{
    TagStruct m_tagHead;
    uint32_t m_bizid;
} TLVBizId;

typedef struct TLVDataContent_
{
    TagStruct m_tagHead;
    char m_data[0];
} TLVDataContent;

typedef struct DataMsg_
{
    TagStruct m_tagHead;
    char m_data[0];
} DataMsg;

#define TLV_HEAD_LEN (sizeof(TagStruct))

#pragma pack()

} // namespace data
} // namespace gse
#endif
