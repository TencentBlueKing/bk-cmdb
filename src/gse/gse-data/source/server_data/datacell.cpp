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

#include "datacell.h"

#include <string.h>

#include "bbx/gse_errno.h"
#include "tools/error.h"
#include "tools/macros.h"
#include "tools/net.h"

namespace gse {
namespace data {
DataCellOPS* DataCell::ToOPS(int status)
{
    m_dataCellOPS.m_status = status;
    return new DataCellOPS(m_dataCellOPS);
}

void DataCell::SetChannelProtocol(const std::string& channelProtocol)
{
    m_dataCellOPS.m_channelProtocol = channelProtocol;
}

std::string DataCell::GetChannelProtocol()
{
    return m_dataCellOPS.m_channelProtocol;
}

std::string DataCell::GetStrServerIp()
{
    return m_dataCellOPS.m_szServerIP;
}

uint16_t DataCell::GetSourcePort() const
{
    return m_dataCellOPS.m_srcPort;
}

void DataCell::SetSourcePort(uint16_t port)
{
    m_dataCellOPS.m_srcPort = port;
}

void DataCell::SetServerPort(uint16_t port)
{
    m_dataCellOPS.m_serverPort = port;
    m_dataCellOPS.m_szServerPort = gse::tools::strings::ToString(port);
}

uint16_t DataCell::GetServerPort() const
{
    return m_dataCellOPS.m_serverPort;
}

char* DataCell::GetDataBuf() const
{
    return m_dataBuf;
}

uint32_t DataCell::GetDataBufLen() const
{
    return m_dataBufLen;
}

void DataCell::SetSourceIp(const std::string& ip)
{
    m_dataCellOPS.m_szSrcIP = ip;
}

void DataCell::SetServerIP(const std::string& ip)
{
    m_dataCellOPS.m_szServerIP = ip;
}

void DataCell::SetExportorName(const std::string& exportorname)
{
    m_dataCellOPS.m_exportorTag = exportorname;
}

std::string DataCell::GetSourceIp()
{
    return m_dataCellOPS.m_szSrcIP;
}

void DataCell::SetChannelID(uint32_t channelID)
{
    m_dataCellOPS.m_channelID = channelID;
    m_dataCellOPS.m_szChannelID = gse::tools::strings::ToString(channelID);
}

void DataCell::SetErrorMsg(const std::string& errmsg, int errorcode)
{
    m_dataCellOPS.m_errmsg = errmsg;
    m_dataCellOPS.m_errcode = errorcode;
}

int DataCell::GetErrorCode()
{
    return m_dataCellOPS.m_errcode;
}

uint32_t DataCell::GetChannelID()
{
    return m_dataCellOPS.m_channelID;
}

void DataCell::SetArrivedTimestamp(uint32_t timestamp)
{
    m_dataCellOPS.m_arrivedTimestamp = timestamp;
}

uint32_t DataCell::GetArrivedTimestamp() const
{
    return m_dataCellOPS.m_arrivedTimestamp;
}

void DataCell::SetOutputTimestamp(uint32_t timestamp)
{
    m_dataCellOPS.m_outputTimestamp = timestamp;
}
uint32_t DataCell::GetOutputTimeStamp() const
{
    return m_dataCellOPS.m_outputTimestamp;
}

void DataCell::SetCreationTimestamp(uint32_t timestamp)
{
    m_dataCellOPS.m_createdTimestamp = timestamp;
    m_dataCellOPS.m_szCreatedTimestamp = gse::tools::strings::ToString(timestamp);
}

uint32_t DataCell::GetCreationTimestamp()
{
    return m_dataCellOPS.m_createdTimestamp;
}

int DataCell::GetDelaySeconds() const
{
    return m_dataCellOPS.m_arrivedTimestamp - m_dataCellOPS.m_createdTimestamp;
}

void DataCell::PushExtension(std::string& extenInfo)
{
    m_dataCellOPS.m_extensionsInfo.push_back(extenInfo);
    // m_extensions.append(extenInfo);
    std::string tagexterninfo = "[" + extenInfo + "]";
    m_extensions.append(tagexterninfo);
}

std::size_t DataCell::GetExtensionSize() const
{
    return m_dataCellOPS.m_extensionsInfo.size();
}

std::string DataCell::GetExtensionByIndex(std::size_t index)
{
    if (index < m_dataCellOPS.m_extensionsInfo.size())
    {
        return m_dataCellOPS.m_extensionsInfo.at(index);
    }
    return std::string("");
}

void DataCell::GetExtensions(std::vector<std::string>& extensions)
{
    std::size_t maxIdx = m_dataCellOPS.m_extensionsInfo.size();
    for (std::size_t idx = 0; idx < maxIdx; ++idx)
    {
        extensions.push_back(m_dataCellOPS.m_extensionsInfo.at(idx));
    }
}

int DataCell::CopyData(const char* buf, uint32_t len)
{
    m_dataCellOPS.m_bytes = len;
    // storage \n\0
    char* pBuf = new char[len + 2];
    if (NULL == pBuf)
    {
        int err = gse_errno;
        LOG_WARN("malloc data cell buffer failed, errno:%d, strerr:%s", err, gse::tools::error::ErrnoToStr(err).c_str());
        return GSE_ERROR;
    }
    pBuf[len] = 0;
    LOG_DEBUG("data cell will copy the data's length is %u", len);
    ::memcpy(pBuf, buf, len);

    if (NULL != m_dataBuf)
    {
        delete[] m_dataBuf;
        m_dataBuf = NULL;
    }

    m_dataBuf = pBuf;
    m_dataBufLen = len;

    return GSE_SUCCESS;
}

void DataCell::DealLineBreak()
{
    if (NULL == m_dataBuf || m_dataBufLen == 0)
    {
        return;
    }

    if (0x0a == m_dataBuf[m_dataBufLen - 1])
    {
        m_dataBuf[m_dataBufLen - 1] = 0;
        m_dataBufLen = m_dataBufLen - 1;
    }
}

void DataCell::AppendLineBreak()
{
    if (NULL == m_dataBuf || m_dataBufLen == 0)
    {
        return;
    }

    if (0x0a != m_dataBuf[m_dataBufLen - 1])
    {
        m_dataBuf[m_dataBufLen] = 0x0a;
        m_dataBufLen += 1;
        m_dataBuf[m_dataBufLen] = 0;
    }
}

uint32_t DataCell::GetBufferLen()
{
    return m_dataBufLen;
}

bool DataCell::IsDataID()
{
    return is_dataid(m_dataCellOPS.m_channelID);
}

void DataCell::SetDataKey(const std::string& dataKey)
{
    m_dataKeyForValue = dataKey;
}

void DataCell::SetOutputTag(const std::string& tag)
{
    m_dataCellOPS.m_outputTag = tag;
}

void DataCell::SetInputTag(const std::string& tag)
{
    m_dataCellOPS.m_inputTag = tag;
}

void DataCell::SetOutputAddress(const std::string& address)
{
    m_dataCellOPS.m_outAddress = address;
}

void DataCell::SetOutputType(const std::string& type)
{
    m_dataCellOPS.m_outputType = type;
}

void DataCell::GetDataKey(std::string& dataKey)
{
    dataKey = m_dataKeyForValue;
}

void DataCell::ClearTableNames()
{
    m_tableNames.clear();
}

void DataCell::AddTableName(const std::string& tableName)
{
    m_tableNames.clear();
    m_tableNames.push_back(tableName);
}
void DataCell::GetTableName(std::vector<std::string>& tableName)
{
    tableName = m_tableNames;
}

void DataCell::SetPartition(int partition)
{
    m_partition = partition;
}

int DataCell::GetPartition()
{
    return m_partition;
}

void DataCell::SetBizID(uint32_t bizid)
{
    m_dataCellOPS.m_bizid = bizid;
}

uint32_t DataCell::GetBizID()
{
    return m_dataCellOPS.m_bizid;
}

bool DataCell::IsOpsMsg()
{
    return m_isOps;
}

void DataCell::SetOpsMsg(bool isOpsMsg)
{
    m_isOps = isOpsMsg;
}

int DataCell::GetOpsServiceId()
{
    return m_opsServiceId;
}

void DataCell::SetOpsServiceId(int opsServiceId)
{
    m_opsServiceId = opsServiceId;
}

} // namespace data
} // namespace gse
