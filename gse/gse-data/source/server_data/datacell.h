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

#ifndef _GSE_DATA_DATACELL_H_
#define _GSE_DATA_DATACELL_H_

#include <string>

#include "log/log.h"
#include "tools/net.h"
#include "tools/strings.h"
#include "tools/time.h"

#include "time_center.h"

namespace gse {
namespace data {

//
// channelid 数据位分布规则：
// 无符号32位整数，高 12位位平台编号，低20位为数据ID
// 对于原有的DataID 平台编号为 0 以兼容存量的ID
//

#ifndef GetPlatNum
#define GetPlatNum(channelID) ((channelID >> 20) & 0x3FF)
#endif

const std::string kLostState = "lost";
const std::string kWaitingState = "waiting";
const std::string kDealingState = "dealing";
const std::string kOutputState = "output";
const std::string kCodecState = "codec";
const std::string kAgentOpsState = "ops";

enum DataCellOpsStatus
{
    EN_DEALING_STATE = 0,
    EN_OUTPUT_STATE = 1,
    EN_LOST_STATE = 2,
    EN_UNKOWN_STATE
};

inline bool is_dataid(uint32_t channelID)
{
    //
    // channelid 数据位分布规则：
    // 无符号32位整数，高 12位位平台编号，低20位为数据ID
    // 对于原有的DataID 平台编号为 0 以兼容存量的ID
    //
    if (channelID == 1)
    {
        return false;
    }
    return false;
    //    return GetPlatNum(channelID) == 0;
}

class DataCellOPS
{
public:
    DataCellOPS()
        : m_szServerPort(""), m_channelProtocol(""), m_szSrcIP(""), m_szServerIP(""), m_state("")
    {
        m_srcPort = 0;
        m_serverPort = 0;
        m_channelID = 0;
        m_szChannelID = "0";
        m_arrivedTimestamp = m_createdTimestamp = TimeCenter::Instance()->GetDateTime();
        m_outputTimestamp = m_arrivedTimestamp;
        m_errcode = 0;
        m_status = EN_UNKOWN_STATE;
        m_bytes = 0;
        m_bizid = 0;
        m_opsType = 0;
        m_bufflen = 0;
    }

    DataCellOPS(const DataCellOPS &ops)
    {
        m_channelID = ops.m_channelID;
        m_srcPort = ops.m_srcPort;
        m_arrivedTimestamp = ops.m_arrivedTimestamp; // data from the protocol head
        m_createdTimestamp = ops.m_createdTimestamp; // creation timestamp
        m_outputTimestamp = ops.m_outputTimestamp;
        m_channelProtocol = ops.m_channelProtocol;
        m_szSrcIP = ops.m_szSrcIP;
        m_state = ops.m_state;
        m_szServerIP = ops.m_szServerIP;
        m_extensionsInfo = ops.m_extensionsInfo; // data from the dynamical
        m_szChannelID = ops.m_szChannelID;
        m_szServerPort = ops.m_szServerPort;
        //延时转成string
        m_szCreatedTimestamp = gse::tools::strings::ToString(m_createdTimestamp);
        m_exportorTag = ops.m_exportorTag;
        m_outputTag = ops.m_outputTag;
        m_errmsg = ops.m_errmsg;
        m_errcode = ops.m_errcode;
        m_status = ops.m_status;
        m_outputType = ops.m_outputType;
        m_outAddress = ops.m_outAddress;
        m_bytes = ops.m_bytes;
        m_bizid = 0;
        m_opsType = 0;
        m_bufflen = 0;
    }

    DataCellOPS &operator=(const DataCellOPS &ops)
    {
        m_channelID = ops.m_channelID;
        m_srcPort = ops.m_srcPort;
        m_serverPort = ops.m_serverPort;
        m_arrivedTimestamp = ops.m_arrivedTimestamp; // data from the protocol head
        m_createdTimestamp = ops.m_createdTimestamp; // creation timestamp
        m_outputTimestamp = ops.m_outputTimestamp;
        m_channelProtocol = ops.m_channelProtocol;
        m_szSrcIP = ops.m_szSrcIP;
        m_state = ops.m_state;
        m_szServerIP = ops.m_szServerIP;
        m_extensionsInfo = ops.m_extensionsInfo; // data from the dynamical
        m_szChannelID = ops.m_szChannelID;
        m_szServerPort = ops.m_szServerPort;
        //延时转成string
        m_szCreatedTimestamp = gse::tools::strings::ToString(m_createdTimestamp);
        m_bytes = ops.m_bytes;
        m_errmsg = ops.m_errmsg;
        m_errcode = ops.m_errcode;
        m_exportorTag = ops.m_exportorTag;
        m_outputTag = ops.m_outputTag;
        m_status = ops.m_status;
        m_outputType = ops.m_outputType;
        m_outAddress = ops.m_outAddress;
        return *this;
    }
    void OpsKey(std::string &ops_key)
    {
        uint64_t min_stamp = m_createdTimestamp / 60 * 60; // 收到的时间戳
        ops_key.append(m_szChannelID);
        ops_key.append(":");
        ops_key.append(m_szSrcIP);
        ops_key.append(gse::tools::strings::ToString(min_stamp));
        if (m_extensionsInfo.size() > 0)
        {
            for (size_t i = 0; i < m_extensionsInfo.size(); i++)
            {
                ops_key.append(m_extensionsInfo[i]);
            }
        }
    }

public:
    int64_t m_bytes;
    uint32_t m_channelID;
    std::string m_szChannelID;
    uint16_t m_srcPort;
    uint16_t m_serverPort;
    std::string m_szServerPort;       // 为了OPS 模块不做转换可以拥有更高的执行效率
    uint32_t m_arrivedTimestamp;      // data from the protocol head
    uint32_t m_createdTimestamp;      // creation timestamp
    uint32_t m_outputTimestamp;       // 处理完成时间
    std::string m_szCreatedTimestamp; // 为 OPS 能够更换的做数据转换，可以拥有更高的执行效率
    std::string m_outputTag;          // 对账使用key
    uint32_t m_bizid;                 // 兼容V1.0 的bizid 云区域等符合ID
    std::string m_channelProtocol;
    std::string m_szSrcIP;
    std::string m_szServerIP;
    std::vector<std::string> m_extensionsInfo; // data from the dynamical
    std::string m_state;                       // dealing , lost etc..
    int m_status;
    uint32_t m_opsType; // ops type, default promethus, 1, new ,send to TGDP kafka
    uint32_t m_bufflen;
    std::string m_datatype;
    std::string m_errmsg;
    int m_errcode;
    std::string m_exportorTag;
    std::string m_inputTag;
    std::string m_outputType;
    std::string m_outAddress;
};

class DataCell
{
public:
    DataCell()
        : m_dataBufLen(0),
          m_dataBuf(nullptr),
          m_partition(-1),
          m_encodeDataBufLen(0),
          m_encodeDataBuf(nullptr),
          m_isOps(false),
          m_opsServiceId(0)

    {
    }

    ~DataCell()
    {
        if (NULL != m_dataBuf)
        {
            delete[] m_dataBuf;
            m_dataBuf = nullptr;
            m_dataBufLen = 0;
        }

        if (m_encodeDataBuf != nullptr)
        {
            delete[] m_encodeDataBuf;
            m_encodeDataBuf = nullptr;
            m_encodeDataBufLen = 0;
        }
    }

public:
    // inline functions
    inline void GetExtensionString(std::string &extension)
    {
        extension.assign(m_extensions);
    }

public:
    // 返回值需要在始用后通过 delete 释放
    DataCellOPS *ToOPS(int status);

public:
    bool IsDataID();

    void SetChannelProtocol(const std::string &channelProtocol);
    std::string GetChannelProtocol();

    uint32_t GetSourceIp() const;
    void SetSourceIp(uint32_t ip);

    std::string GetStrServerIp();
    uint32_t GetServerIp();

    uint16_t GetSourcePort() const;
    void SetSourcePort(uint16_t port);

    void SetServerPort(uint16_t port);
    uint16_t GetServerPort() const;

    char *GetDataBuf() const;
    uint32_t GetDataBufLen() const;

    void SetSourceIp(const std::string &ip);
    void SetServerIP(const std::string &ip);
    std::string GetSourceIp();

    void SetChannelID(uint32_t channelID);
    void SetErrorMsg(const std::string &errmsg, int errorcode);
    int GetErrorCode();
    uint32_t GetChannelID();

    void SetExportorName(const std::string &exportorname);

    void SetCreationTimestamp(uint32_t timestamp);
    uint32_t GetCreationTimestamp();

    void SetArrivedTimestamp(uint32_t timestamp);
    uint32_t GetArrivedTimestamp() const;

    void SetOutputTimestamp(uint32_t timestamp);
    uint32_t GetOutputTimeStamp() const;
    int GetDelaySeconds() const;

    void PushExtension(std::string &extenInfo);
    std::size_t GetExtensionSize() const;
    std::string GetExtensionByIndex(std::size_t index);
    void GetExtensions(std::vector<std::string> &extensions);

    int CopyData(const char *buf, uint32_t len);
    void DealLineBreak();
    void AppendLineBreak();

public:
    void SetDataKey(const std::string &dataKey);
    void SetOutputTag(const std::string &tag);
    void SetInputTag(const std::string &tag);
    void SetOutputAddress(const std::string &address);
    void SetOutputType(const std::string &type);
    void GetDataKey(std::string &dataKey);
    void ClearTableNames();
    void AddTableName(const std::string &tableName);
    void GetTableName(std::vector<std::string> &tableName);
    void SetPartition(int partition);
    int GetPartition();
    void SetBizID(uint32_t bizid);
    uint32_t GetBizID();
    uint32_t GetBufferLen();
    bool IsOpsMsg();
    void SetOpsMsg(bool isOpsMsg);

    int GetOpsServiceId();
    void SetOpsServiceId(int opsServiceId);

private:
    uint32_t m_dataBufLen;
    char *m_dataBuf;
    DataCellOPS m_dataCellOPS;
    std::vector<std::string> m_filterValues;

    uint32_t m_encodeDataBufLen;
    char *m_encodeDataBuf;
    bool m_isOps;
    int m_opsServiceId;

private:
    std::string m_dataKeyForValue;
    std::vector<std::string> m_tableNames;
    int m_partition;

    std::string m_extensions;
};
} // namespace data
} // namespace gse

#endif
