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

#include "ops.h"

#include <atomic>
#include <gperftools/malloc_extension.h>
#include <unordered_map>
#include <vector>

#include "bbx/gse_errno.h"
#include "eventthread/event_thread.h"
#include "log/log.h"

#include "tools/error.h"
#include "tools/macros.h"
#include "tools/pipe.h"

#include "ops/report_client.h"
#include "utils.h"

//#include "prometheusHttpThread.h"
//#include "monitorTag.h"
namespace gse {
namespace data {

const std::string name_ops_queune_msg_count = "ops_queue_msg_count";

OpsStatWorker::OpsStatWorker()
{
    m_ptrEventThread = NULL;
    m_ptrReportClient = NULL;
    m_ptrOpsQueue = new moodycamel::ConcurrentQueue<DataCellOPS *>();
    m_hasOpsServer = false;
    m_threadCount = 8;
    m_opsMsgHandleEventId = -1;
    m_channelID = -1;
    m_agentOpsChannelID = 1;
    m_enableOps = false;
}

OpsStatWorker::~OpsStatWorker()
{
    delete m_ptrOpsQueue;
}

void OpsStatWorker::makeLogicalTag(std::string &channleid, rapidjson::Document &logical_tag)
{
    /*
    "logical_tag": { "tag": "rt1",
    "desc": { "result_table_id": "rt1"
    } },
    */
    rapidjson::Document tag_desc(rapidjson::kObjectType);

    rapidjson::Document logical_tag_obj(rapidjson::kObjectType);
    rapidjson::Value str_channelid;
    str_channelid.SetString(channleid.c_str(), channleid.length(), logical_tag.GetAllocator());
    logical_tag_obj.AddMember("tag", str_channelid, logical_tag.GetAllocator());
    logical_tag_obj.AddMember("desc", tag_desc, logical_tag.GetAllocator());

    logical_tag.AddMember("logical_tag", logical_tag_obj, logical_tag.GetAllocator());
}

void OpsStatWorker::makeCustomTagJson(std::unordered_map<std::string, std::string> &custom_tags, rapidjson::Document &custom_tag)
{
    /*
    "custom_tags":
    {
     "ip": "0.0.0.0",
    "port": "8888",
    "task": "1"
    }
    */
    rapidjson::Document custom_tag_obj(rapidjson::kObjectType);

    for (std::unordered_map<std::string, std::string>::iterator it = custom_tags.begin(); it != custom_tags.end(); it++)
    {
        rapidjson::Value str_custom_key;
        rapidjson::Value str_custom_value;
        LOG_DEBUG("make custom tag json:key:%s, value:%s, tags count:%d", it->first.c_str(), it->second.c_str(), custom_tags.size());
        str_custom_key.SetString(it->first.c_str(), it->first.length(), custom_tag.GetAllocator());
        str_custom_value.SetString(it->second.c_str(), it->second.length(), custom_tag.GetAllocator());
        custom_tag_obj.AddMember(str_custom_key, str_custom_value, custom_tag.GetAllocator());
    }

    custom_tag.AddMember("custom_tags", custom_tag_obj, custom_tag.GetAllocator());

    rapidjson::StringBuffer str_buff;
    rapidjson::Writer<rapidjson::StringBuffer> writer(str_buff);
    custom_tag.Accept(writer);

    LOG_DEBUG("make custom tag :%s", str_buff.GetString());
}

void OpsStatWorker::makeNodeInfoJson(rapidjson::Document &node_info)
{
    //--------add node info tag ----------{
    node_info.AddMember("module", "collector", node_info.GetAllocator());
    node_info.AddMember("component", "ds", node_info.GetAllocator());
    //--------add node info tag ----------}
}

void OpsStatWorker::makePhysicalTagJson(rapidjson::Document &physical_tag_json)
{
    //--------add physical tag ----------{

    rapidjson::Document physical_tag_obj_json(rapidjson::kObjectType);
    std::string str_tag = "ds|" + m_reportIP;
    rapidjson::Value str_tag_value;
    str_tag_value.SetString(str_tag.c_str(), str_tag.length(), physical_tag_json.GetAllocator());
    physical_tag_obj_json.AddMember("tag", str_tag_value, physical_tag_json.GetAllocator());

    rapidjson::Value str_tag_ip;
    str_tag_ip.SetString(m_reportIP.c_str(), m_reportIP.length(), physical_tag_json.GetAllocator());

    rapidjson::Document physical_tag_desc(rapidjson::kObjectType);
    physical_tag_desc.AddMember("module", "ds", physical_tag_json.GetAllocator());
    physical_tag_desc.AddMember("ip", str_tag_ip, physical_tag_json.GetAllocator());
    physical_tag_desc.AddMember("complexid", "0", physical_tag_json.GetAllocator());

    rapidjson::Value str_cluster_name_value;
    str_cluster_name_value.SetString(m_configPtr->GetDataSvrClusterName().c_str(), m_configPtr->GetDataSvrClusterName().length(), physical_tag_json.GetAllocator());
    physical_tag_desc.AddMember("cluster_name", str_cluster_name_value, physical_tag_json.GetAllocator());

    rapidjson::Value str_cluster_instanceid_value;
    str_cluster_instanceid_value.SetString(m_configPtr->GetDataSvrInstanceId().c_str(), m_configPtr->GetDataSvrInstanceId().length(), physical_tag_json.GetAllocator());
    physical_tag_desc.AddMember("instance_id", str_cluster_instanceid_value, physical_tag_json.GetAllocator());

    physical_tag_obj_json.AddMember("desc", physical_tag_desc, physical_tag_json.GetAllocator());

    physical_tag_json.AddMember("physical_tag", physical_tag_obj_json, physical_tag_json.GetAllocator());
    //--------add physical tag ----------}
}

bool OpsStatWorker::PushOpsQueue(DataCellOPS *ptr)
{
    int max_retry_times = 0;
    while (!m_ptrOpsQueue->enqueue(ptr))
    {
        max_retry_times++;
        if (max_retry_times > 5)
        {
            return false;
        }
    }
    return true;
}

void OpsStatWorker::calcDelay(DataCellOPS *ptr_ops, OpsDataStat *dataloss)
{
    int delay = (ptr_ops->m_arrivedTimestamp > ptr_ops->m_createdTimestamp) ? (ptr_ops->m_arrivedTimestamp - ptr_ops->m_createdTimestamp) : 0;

    if (delay >= dataloss->m_delay.m_maxDelay.m_delay)
    {
        dataloss->m_delay.m_maxDelay.m_delay = delay;
        dataloss->m_delay.m_maxDelay.m_local = ptr_ops->m_arrivedTimestamp;
        dataloss->m_delay.m_maxDelay.m_recv = ptr_ops->m_createdTimestamp;
    }

    if (delay <= dataloss->m_delay.m_minDelay.m_delay)
    {
        dataloss->m_delay.m_minDelay.m_local = ptr_ops->m_arrivedTimestamp;
        dataloss->m_delay.m_minDelay.m_recv = ptr_ops->m_createdTimestamp;
        dataloss->m_delay.m_minDelay.m_delay = delay;
    }
    else
    {
        dataloss->m_delay.m_minDelay.m_local = ptr_ops->m_arrivedTimestamp;
        dataloss->m_delay.m_minDelay.m_recv = ptr_ops->m_createdTimestamp;
        dataloss->m_delay.m_minDelay.m_delay = delay;
    }
}

void OpsStatWorker::calcOutputDelay(DataCellOPS *ptr_ops, OpsDataStat *dataloss)
{
    int delay = (ptr_ops->m_outputTimestamp > ptr_ops->m_arrivedTimestamp) ? (ptr_ops->m_outputTimestamp - ptr_ops->m_outputTimestamp) : 0;

    if (delay >= dataloss->m_outputDelay.m_maxDelay.m_delay)
    {
        dataloss->m_outputDelay.m_maxDelay.m_delay = delay;
        dataloss->m_outputDelay.m_maxDelay.m_local = ptr_ops->m_arrivedTimestamp;
        dataloss->m_outputDelay.m_maxDelay.m_recv = ptr_ops->m_outputTimestamp;
    }

    if (delay <= dataloss->m_delay.m_minDelay.m_delay)
    {
        dataloss->m_outputDelay.m_minDelay.m_local = ptr_ops->m_arrivedTimestamp;
        dataloss->m_outputDelay.m_minDelay.m_recv = ptr_ops->m_outputTimestamp;
        dataloss->m_outputDelay.m_minDelay.m_delay = delay;
    }
    else
    {
        dataloss->m_delay.m_minDelay.m_local = ptr_ops->m_arrivedTimestamp;
        dataloss->m_delay.m_minDelay.m_recv = ptr_ops->m_createdTimestamp;
        dataloss->m_delay.m_minDelay.m_delay = delay;
    }
}

void OpsStatWorker::MakeKeyTag(DataCellOPS *ptr_ops, std::string &tag)
{
    // key_tag=dataid|ip|serverip|
    //"%u|%lu|%s", pDataCell->GetChannelID(), agentTimestamp, SAFE_CSTR(pDataCell->GetServerIp().c_str())
    tag.append(ptr_ops->m_szChannelID);
    tag.append("|");

    uint32_t min = toMinutes(ptr_ops->m_createdTimestamp);
    std::string str_arrivedtimestamp = gse::tools::strings::ToString(min);
    tag.append(str_arrivedtimestamp);
    tag.append("|");
    tag.append(ptr_ops->m_szSrcIP);
    return;
}

void OpsStatWorker::CalcRecvCount(DataCellOPS *ptr_ops, MultiDataLoss *count)
{
    // OpsMsgCount *ops_total_count = GetTotalMsgStat(ptr_ops);

    Input *ptr_input = count->m_input;

    if (ptr_input == NULL)
    {
        ptr_input = new Input();
        count->m_input = ptr_input;
    }

    ptr_input->m_incrementCnt += 1;
    ptr_input->m_incrementBytes += ptr_ops->m_bytes;
    ptr_input->m_isValid = true;

    //    if (ops_total_count != NULL)
    //    {
    //        ptr_input->m_totalCnt = ops_total_count->m_input.m_count;
    //    }
    //    else
    //    {
    //        ptr_input->m_totalCnt = 1;
    //    }

    std::string key_tag;
    // MakeKeyTag(ptr_ops, key_tag);
    if (ptr_ops->m_inputTag.empty())
    {
        MakeKeyTag(ptr_ops, key_tag);
    }
    else
    {
        key_tag = ptr_ops->m_inputTag;
    }

    LOG_DEBUG("recv input tag:%s", key_tag.c_str());
    if (ptr_input->m_tags.find(key_tag) == ptr_input->m_tags.end())
    {
        ptr_input->m_tags[key_tag] = 1;
    }
    else
    {
        ptr_input->m_tags[key_tag] += 1;
    }
}

void OpsStatWorker::CalcSendAccount(DataCellOPS *ptr_ops, MultiDataLoss *count)
{
    //    OpsMsgCount *ops_total_count = GetTotalMsgStat(ptr_ops);

    OutputEx *ptr_output = NULL;

    std::map<std::string, OutputEx *>::iterator it;

    it = count->m_output.find(ptr_ops->m_outAddress);

    if (it == count->m_output.end())
    {
        ptr_output = new OutputEx();
        count->m_output[ptr_ops->m_outAddress] = ptr_output;
        LOG_DEBUG("create output stat obj, key:%s", ptr_ops->m_outAddress.c_str());
    }
    else
    {
        ptr_output = it->second;
    }

    if (ptr_output == NULL)
    {
        return;
    }

    ptr_output->m_incrementCnt += 1;
    ptr_output->m_incrementBytes += ptr_ops->m_bytes;
    ptr_output->m_isValid = true;
    //    if (ops_total_count != NULL)
    //    {
    //        ptr_output->m_totalCnt = ops_total_count->m_output.m_count;
    //    }
    //    else
    //    {
    //        ptr_output->m_totalCnt = 1;
    //    }

    ptr_output->m_type = ptr_ops->m_outputType;
    ptr_output->m_address = ptr_ops->m_outAddress;

    std::string output_tag;
    if (!ptr_ops->m_outputTag.empty())
    {
        output_tag = ptr_ops->m_outputTag;
    }
    else
    {
        MakeKeyTag(ptr_ops, output_tag);
    }

    LOG_DEBUG("calc send stat, output tag:%s", ptr_ops->m_outputTag.c_str());

    if (ptr_output->m_tags.find(output_tag) == ptr_output->m_tags.end())
    {
        ptr_output->m_tags[output_tag] = 1;
    }
    else
    {
        ptr_output->m_tags[output_tag] += 1;
    }
}

void OpsStatWorker::CalcDropCount(DataCellOPS *ptr_ops, MultiDataLoss *count)
{
    std::string str_errcode = gse::tools::strings::ToString(ptr_ops->m_errcode);
    DropObj *ptr_drop = NULL;

    std::map<std::string, DropObj *>::iterator it = count->m_drop.find(str_errcode);
    if (it == count->m_drop.end())
    {
        ptr_drop = new DropObj();

        ptr_drop->m_count = 1;
        ptr_drop->m_reason = ptr_ops->m_errmsg;
        ptr_drop->m_errcode = ptr_ops->m_errcode;

        count->m_drop[str_errcode] = ptr_drop;
    }
    else
    {
        ptr_drop = it->second;
        ptr_drop->m_count += 1;
    }
}

/*
PerMinIpStat * OpsCollection::newPerMinIpStat()
{
    PerMinIpStat *per_min_stat = new PerMinIpStat();

    return per_min_stat;
}
*/

void OpsStatWorker::MakeCustomTag(DataCellOPS *ptr_ops, OpsDataStat *dataloss, std::string &key)
{
    if (ptr_ops == NULL || dataloss == NULL)
    {
        LOG_ERROR("ops ptr(%p) or data loss ptr(%p) is null ptr", ptr_ops, dataloss);
        return;
    }

    if (dataloss->m_customTags.find(key) == dataloss->m_customTags.end())
    {
        dataloss->m_customTags.insert(std::make_pair(key, key));
        time_t min_time = toMinutes(ptr_ops->m_createdTimestamp);
        dataloss->m_customTags.insert(std::make_pair("channel_id", ptr_ops->m_szChannelID));
        dataloss->m_customTags.insert(std::make_pair("report_ip", m_reportIP));
        dataloss->m_customTags.insert(std::make_pair("source_ip", ptr_ops->m_szSrcIP));
        dataloss->m_customTags.insert(std::make_pair("listenport", ptr_ops->m_szServerPort));
        dataloss->m_customTags.insert(std::make_pair("instanceid", m_configPtr->GetDataSvrInstanceId()));
        dataloss->m_customTags.insert(std::make_pair("cluster_name", m_configPtr->GetDataSvrClusterName()));
        dataloss->m_customTags.insert(std::make_pair("protocol", ptr_ops->m_channelProtocol));
        dataloss->m_customTags.insert(std::make_pair("create_time", gse::tools::strings::ToString(min_time)));

        LOG_DEBUG("make custom tag, key:%s", key.c_str());
    }
}

void OpsStatWorker::CountMsg(DataCellOPS *ptr_ops, OpsDataStat *dataloss)
{
    if (dataloss == NULL)
    {
        LOG_ERROR("ptr_ops invalid,null ptr");
        return;
    }
    switch (ptr_ops->m_status)
    {
    case EN_OUTPUT_STATE:
        CalcSendAccount(ptr_ops, &dataloss->m_dataLoss);
        break;
    case EN_DEALING_STATE:
        CalcRecvCount(ptr_ops, &dataloss->m_dataLoss);
        break;
    case EN_LOST_STATE:
        CalcDropCount(ptr_ops, &dataloss->m_dataLoss);
        break;
    default:
        LOG_ERROR("recv unkow stat type:%d", ptr_ops->m_status);
        break;
    }
}

OpsMsgCount *OpsStatWorker::GetTotalMsgStat(DataCellOPS *ptr_ops)
{
    std::string key;
    if (ptr_ops->m_status == EN_OUTPUT_STATE)
    {
        key = ptr_ops->m_szChannelID + "_" + ptr_ops->m_szSrcIP + "_" + ptr_ops->m_outputType + ptr_ops->m_outAddress;
    }
    else
    {
        key = ptr_ops->m_szChannelID + "_" + ptr_ops->m_szSrcIP;
    }

    std::unordered_map<std::string, OpsMsgCount *>::iterator it = m_allMsgTotalStat.find(key);
    if (it == m_allMsgTotalStat.end())
    {
        return NULL;
    }
    else
    {
        return it->second;
    }
}

void OpsStatWorker::RecvTotalMsgStat(DataCellOPS *ptr_ops)
{
    std::string key;
    if (ptr_ops->m_status == EN_OUTPUT_STATE)
    {
        key = ptr_ops->m_szChannelID + "_" + ptr_ops->m_szSrcIP + "_" + ptr_ops->m_outputType + ptr_ops->m_outAddress;
    }
    else
    {
        key = ptr_ops->m_szChannelID + "_" + ptr_ops->m_szSrcIP;
    }

    std::unordered_map<std::string, OpsMsgCount *>::iterator it = m_allMsgTotalStat.find(key);

    OpsMsgCount *ptr_ops_msg_count = NULL;
    if (it == m_allMsgTotalStat.end())
    {
        ptr_ops_msg_count = new OpsMsgCount();
        m_allMsgTotalStat.insert(std::make_pair(key, ptr_ops_msg_count));
    }
    else
    {
        ptr_ops_msg_count = it->second;
    }

    //计算total数
    switch (ptr_ops->m_status)
    {
    case EN_DEALING_STATE:
        ptr_ops_msg_count->m_input.m_count += 1;
        ptr_ops_msg_count->m_input.m_bytes += ptr_ops->m_bytes;
        break;
    case EN_OUTPUT_STATE:
        ptr_ops_msg_count->m_output.m_count += 1;
        ptr_ops_msg_count->m_output.m_bytes += ptr_ops->m_bytes;
        break;
    case EN_LOST_STATE:
        ptr_ops_msg_count->m_drop.m_count += 1;
        ptr_ops_msg_count->m_drop.m_bytes += ptr_ops->m_bytes;
        break;
    }
}

void OpsStatWorker::StatRecvPkg(DataCellOPS *ptr_ops)
{
    if (ptr_ops == NULL)
    {
        return;
    }

    // RecvTotalMsgStat(ptr_ops);

    std::string stat_key;
    ptr_ops->OpsKey(stat_key);
    OpsDataStat *dataloss = NULL;
    std::unordered_map<std::string, OpsDataStat *>::iterator it = m_allMsgStat.find(stat_key);
    if (it == m_allMsgStat.end())
    {
        dataloss = new OpsDataStat();
        m_allMsgStat.insert(std::make_pair(stat_key, dataloss));
        LOG_DEBUG("ops dont find ops data stat item, new[%p] by key:%s", dataloss, stat_key.c_str());

        LOG_DEBUG("start make custom tag json by key:%s", stat_key.c_str());
        MakeCustomTag(ptr_ops, dataloss, stat_key);
    }
    else
    {
        dataloss = it->second;
        LOG_DEBUG("ops find ops data stat item[%p] by key:%s", dataloss, stat_key.c_str());
    }

    CountMsg(ptr_ops, dataloss);
    calcDelay(ptr_ops, dataloss);
    calcOutputDelay(ptr_ops, dataloss);
}

void OpsStatWorker::makeDataMonitorJson(OpsDataStat *ptr_data_loss, rapidjson::Document &data_monitor_json)
{
    /*
    "input": {
    "tags": { "front_rt1_1_1577808000": 100, "front_rt1_2_1577808000": 200
    },
    "total_cnt": 4000,
    "total_cnt_increment": 300
    }
    */
    rapidjson::Document data_loss_obj(rapidjson::kObjectType);
    rapidjson::Document input(rapidjson::kObjectType);

    if (ptr_data_loss->m_dataLoss.m_input != NULL && ptr_data_loss->m_dataLoss.m_input->m_isValid)
    {
        rapidjson::Document input_tags(rapidjson::kObjectType);
        rapidjson::Document input_tags_obj(rapidjson::kObjectType);
        for (std::map<std::string, uint64_t>::iterator it = ptr_data_loss->m_dataLoss.m_input->m_tags.begin(); it != ptr_data_loss->m_dataLoss.m_input->m_tags.end(); ++it)
        {
            rapidjson::Value str_key;
            str_key.SetString(it->first.c_str(), it->first.length(), data_monitor_json.GetAllocator());
            input_tags.AddMember(str_key, it->second, data_monitor_json.GetAllocator());
        }
        // input.AddMember("total_cnt", ptr_data_loss->m_dataLoss.m_input->m_totalCnt, data_monitor_json.GetAllocator());
        input.AddMember("total_cnt_increment", ptr_data_loss->m_dataLoss.m_input->m_incrementCnt, data_monitor_json.GetAllocator());
        input.AddMember("increment_bytes", ptr_data_loss->m_dataLoss.m_input->m_incrementBytes, data_monitor_json.GetAllocator());
        input.AddMember("tags", input_tags, data_monitor_json.GetAllocator());

        data_loss_obj.AddMember("input", input, data_monitor_json.GetAllocator());
    }
    else
    {
        data_loss_obj.AddMember("input", input, data_monitor_json.GetAllocator());
    }

    /*
68     //5.4 结构体数组
69     rapidjson::Value ObjectArray(rapidjson::kArrayType);
70     for(int i = 1; i < 4; i++)
71     {
72         rapidjson::Value obj(rapidjson::kObjectType);
73         obj.AddMember("name", "qq849635649", allocator);
74         obj.AddMember("age", i * 10, allocator);
75         ObjectArray.PushBack(obj, allocator);
76     }
77     doc.AddMember("ObjectArray", ObjectArray, allocator);
*/

    rapidjson::Value output_array(rapidjson::kArrayType);
    for (std::map<std::string, OutputEx *>::iterator it_output = ptr_data_loss->m_dataLoss.m_output.begin(); it_output != ptr_data_loss->m_dataLoss.m_output.end(); ++it_output)
    {
        rapidjson::Document output(rapidjson::kObjectType);
        rapidjson::Document output_tags(rapidjson::kObjectType);
        for (std::map<std::string, uint64_t>::iterator it_tags = it_output->second->m_tags.begin(); it_tags != it_output->second->m_tags.end(); ++it_tags)
        {
            rapidjson::Value str_key(rapidjson::kStringType);
            str_key.SetString(it_tags->first.c_str(), it_tags->first.length(), data_monitor_json.GetAllocator());
            output_tags.AddMember(str_key, it_tags->second, data_monitor_json.GetAllocator());
        }

        // output.AddMember("total_cnt", it_output->second->m_totalCnt, data_monitor_json.GetAllocator());
        output.AddMember("total_cnt_increment", it_output->second->m_incrementCnt, data_monitor_json.GetAllocator());
        output.AddMember("increment_bytes", it_output->second->m_incrementBytes, data_monitor_json.GetAllocator());
        output.AddMember("tags", output_tags, data_monitor_json.GetAllocator());
        rapidjson::Value str_type(rapidjson::kStringType);
        str_type.SetString(it_output->second->m_type.c_str(), it_output->second->m_type.length());
        output.AddMember("type", str_type, data_monitor_json.GetAllocator());
        rapidjson::Value str_address(rapidjson::kStringType);
        str_address.SetString(it_output->second->m_address.c_str(), it_output->second->m_address.length());
        output.AddMember("output_address", str_address, data_monitor_json.GetAllocator());
        output_array.PushBack(output, data_monitor_json.GetAllocator());

        LOG_DEBUG("make output json, key:%s", it_output->first.c_str());
    }

    data_loss_obj.AddMember("output", output_array, data_monitor_json.GetAllocator());

    // Add Data drop json
    //--------------------------------------------------

    rapidjson::Value drop_array(rapidjson::kArrayType);

    for (std::map<std::string, DropObj *>::iterator it = ptr_data_loss->m_dataLoss.m_drop.begin(); it != ptr_data_loss->m_dataLoss.m_drop.end(); ++it)
    {
        rapidjson::Document error_code(rapidjson::kObjectType);

        error_code.AddMember("cnt", it->second->m_count, data_monitor_json.GetAllocator());

        rapidjson::Value str_reason_value;
        str_reason_value.SetString(it->second->m_reason.c_str(), it->second->m_reason.length(), data_monitor_json.GetAllocator());

        error_code.AddMember("reason", str_reason_value, data_monitor_json.GetAllocator());
        error_code.AddMember("errcode", it->second->m_errcode, data_monitor_json.GetAllocator());
        // std::string error_key  = it->first;
        // rapidjson::Value str_key;
        // str_key.SetString(error_key.c_str(), error_key.length(), data_monitor_json.GetAllocator());
        // data_drop_obj.AddMember(str_key, error_code, data_monitor_json.GetAllocator());
        drop_array.PushBack(error_code, data_monitor_json.GetAllocator());
    }

    data_loss_obj.AddMember("data_drop", drop_array, data_monitor_json.GetAllocator());

    data_monitor_json.AddMember("data_loss", data_loss_obj, data_monitor_json.GetAllocator());
    //--------------------------------------------------
    rapidjson::Document data_delay_obj(rapidjson::kObjectType);
    data_delay_obj.AddMember("window_time", 60, data_monitor_json.GetAllocator());
    data_delay_obj.AddMember("waiting_time", 30, data_monitor_json.GetAllocator());

    /*
"data_delay":
{ "window_time": 60,
"waiting_time": 30,
"min_delay":
{
"output_time": 1577808000, "data_time": 1577808000, "delay_time" : 0
},
"max_delay": {
"output_time": 1577808000, "data_time": 1577808000, "delay_time" : 0
}
}
*/
    rapidjson::Document data_min_delay(rapidjson::kObjectType);
    data_min_delay.AddMember("output_time", ptr_data_loss->m_delay.m_minDelay.m_recv, data_monitor_json.GetAllocator());
    data_min_delay.AddMember("data_time", ptr_data_loss->m_delay.m_minDelay.m_local, data_monitor_json.GetAllocator());
    data_min_delay.AddMember("delay_time", ptr_data_loss->m_delay.m_minDelay.m_delay, data_monitor_json.GetAllocator());

    rapidjson::Document data_max_delay(rapidjson::kObjectType);
    data_max_delay.AddMember("output_time", ptr_data_loss->m_delay.m_maxDelay.m_recv, data_monitor_json.GetAllocator());
    data_max_delay.AddMember("data_time", ptr_data_loss->m_delay.m_maxDelay.m_local, data_monitor_json.GetAllocator());
    data_max_delay.AddMember("delay_time", ptr_data_loss->m_delay.m_maxDelay.m_delay, data_monitor_json.GetAllocator());

    data_delay_obj.AddMember("min_delay", data_min_delay, data_monitor_json.GetAllocator());
    data_delay_obj.AddMember("max_delay", data_max_delay, data_monitor_json.GetAllocator());

    data_monitor_json.AddMember("data_delay", data_delay_obj, data_monitor_json.GetAllocator());

    rapidjson::Document data_out_delay_obj(rapidjson::kObjectType);
    data_out_delay_obj.AddMember("window_time", 60, data_monitor_json.GetAllocator());
    data_out_delay_obj.AddMember("waiting_time", 30, data_monitor_json.GetAllocator());

    rapidjson::Document data_out_min_delay(rapidjson::kObjectType);
    data_out_min_delay.AddMember("output_time", ptr_data_loss->m_outputDelay.m_minDelay.m_recv, data_monitor_json.GetAllocator());
    data_out_min_delay.AddMember("data_time", ptr_data_loss->m_outputDelay.m_minDelay.m_local, data_monitor_json.GetAllocator());
    data_out_min_delay.AddMember("delay_time", ptr_data_loss->m_outputDelay.m_minDelay.m_delay, data_monitor_json.GetAllocator());

    rapidjson::Document data_out_max_delay(rapidjson::kObjectType);
    data_out_max_delay.AddMember("output_time", ptr_data_loss->m_outputDelay.m_maxDelay.m_recv, data_monitor_json.GetAllocator());
    data_out_max_delay.AddMember("data_time", ptr_data_loss->m_outputDelay.m_maxDelay.m_local, data_monitor_json.GetAllocator());
    data_out_max_delay.AddMember("delay_time", ptr_data_loss->m_outputDelay.m_maxDelay.m_delay, data_monitor_json.GetAllocator());

    data_out_delay_obj.AddMember("min_delay", data_out_min_delay, data_monitor_json.GetAllocator());
    data_out_delay_obj.AddMember("max_delay", data_out_max_delay, data_monitor_json.GetAllocator());

    data_monitor_json.AddMember("data_out_delay", data_out_delay_obj, data_monitor_json.GetAllocator());
}

void OpsStatWorker::ConvertToReportJsonString(std::string &key, OpsDataStat *ops_data_stat, std::string &report_json)
{
    /*
    "time": 1577808000, "version": "3.0", "info": {
    "module": "stream", "component": "flink", "cluster": "default",
    */
    rapidjson::Document root_json(rapidjson::kObjectType);

    std::string channelid = "-1";
    std::size_t found = key.find(":");
    if (found != std::string::npos)
    {
        channelid = key.substr(0, found);
    }

    rapidjson::Document info_json(rapidjson::kObjectType);

    makeNodeInfoJson(info_json);
    makePhysicalTagJson(info_json);
    makeLogicalTag(channelid, info_json);
    makeCustomTagJson(ops_data_stat->m_customTags, info_json);

    root_json.AddMember("info", info_json, root_json.GetAllocator());

    //----------------data monitor-----------------------
    rapidjson::Document data_monitor_json(rapidjson::kObjectType);
    makeDataMonitorJson(ops_data_stat, data_monitor_json);

    rapidjson::Document data_monitor_obj(rapidjson::kObjectType);
    data_monitor_obj.AddMember("data_monitor", data_monitor_json, root_json.GetAllocator());

    root_json.AddMember("metrics", data_monitor_obj, root_json.GetAllocator());

    root_json.AddMember("version", "3.5", root_json.GetAllocator());
    root_json.AddMember("time", ::time(NULL), root_json.GetAllocator());

    rapidjson::StringBuffer str_buff;
    rapidjson::Writer<rapidjson::StringBuffer> writer(str_buff);
    root_json.Accept(writer);

    report_json.assign(str_buff.GetString());
}

void OpsStatWorker::ReportAllChannelIdStat()
{
    std::string report_stat_json;
    for (std::unordered_map<std::string, OpsDataStat *>::iterator it_channelid = m_allMsgStat.begin(); it_channelid != m_allMsgStat.end(); it_channelid++)
    {
        std::string key = it_channelid->first;
        OpsDataStat *ops_data_stat = it_channelid->second;
        ConvertToReportJsonString(key, ops_data_stat, report_stat_json);
        delete it_channelid->second;
        LOG_DEBUG("report ops key:%s, json:[%s]", key.c_str(), report_stat_json.c_str());
        m_opsClient.SendToByUDP(m_channelID, report_stat_json.c_str(), report_stat_json.size());
    }

    m_allMsgStat.clear();
}

uint32_t OpsStatWorker::toMinutes(time_t timestamp)
{
    return (((uint32_t)timestamp / 60) * 60);
}

void OpsStatWorker::splitMonitorKey(uint64_t key, uint32_t &channelid, uint32_t &timestamp)
{
    channelid = key >> 32;
    timestamp = key & 0xFFFFFFFF;
}

void OpsStatWorker::reportToMonitor(evutil_socket_t fd)
{
    if (!m_enableOps)
    {
        return;
    }

    ReportAllChannelIdStat();
}

void OpsStatWorker::handleMonitorReport(evutil_socket_t fd)
{
    DataCellOPS *ptr_datacell_ops = NULL;
    int count = 0;

    while (m_ptrOpsQueue->try_dequeue(ptr_datacell_ops))
    {
        StatRecvPkg(ptr_datacell_ops);
        if (ptr_datacell_ops != NULL)
        {
            delete ptr_datacell_ops;
            ptr_datacell_ops = NULL;
        }
    }
}

void OpsStatWorker::DealData(evutil_socket_t fd, short what, void *args)
{
    if (NULL == args)
    {
        return;
    }

    OpsStatWorker *ptr_this = reinterpret_cast<OpsStatWorker *>(args);
    ptr_this->handleMonitorReport(fd);
}

void OpsStatWorker::perMinReportMonitor(evutil_socket_t fd, short what, void *args)
{
    if (NULL == args)
    {
        return;
    }

    OpsStatWorker *ptr_this = reinterpret_cast<OpsStatWorker *>(args);
    ptr_this->reportToMonitor(fd);
}

void OpsStatWorker::OpsQueueSizeReport()
{
    int thread_id = gse::tools::thread::GetThreadID();
    OPMetric::Instance().ReportQueneCount(const_cast<std::string &>(name_ops_queune_msg_count), thread_id, m_ptrOpsQueue->size_approx());
}

void OpsStatWorker::OpsQueueSizeReportEvent(evutil_socket_t fd, short what, void *args)
{
    if (NULL == args)
    {
        return;
    }

    OpsStatWorker *ptr_this = reinterpret_cast<OpsStatWorker *>(args);
    ptr_this->OpsQueueSizeReport();
}

void OpsStatWorker::SetOpsParam(std::string &report_ip, int channel_id)
{
    m_reportIP = report_ip;
    m_channelID = channel_id;
}
void OpsStatWorker::SetOpsServerAddress(std::vector<OPSServerAddress> ops_servers)
{
    m_opsServerAddress = ops_servers;
}

int OpsStatWorker::Start()
{
    int ret = m_opsClient.OpenUDP(m_opsServerAddress);
    if (ret != GSE_SUCCESS)
    {
        int tmp_errno = gse_socket_errno;
        LOG_ERROR("failed to connect ops udp server, [errorno:%d, errormsg:%s]", tmp_errno, gse::tools::error::ErrnoToStr(tmp_errno).c_str());
        return GSE_ERROR;
    }

    m_ptrEventThread = new EventThread();

    m_opsMsgHandleEventId = m_ptrEventThread->RegisterTimerPersistEvent(OpsStatWorker::DealData, this, 1);
    if (m_opsMsgHandleEventId < 0)
    {
        LOG_ERROR("failed to register event function callback , error code %d", m_opsMsgHandleEventId);
        return m_opsMsgHandleEventId;
    }

    int event_index = m_ptrEventThread->RegisterTimerPersistEvent(OpsStatWorker::perMinReportMonitor, this, 60);
    if (event_index < 0)
    {
        LOG_ERROR("failed to register timer, error code %d", event_index);
        return event_index;
    }

    event_index = m_ptrEventThread->RegisterTimerPersistEvent(OpsStatWorker::OpsQueueSizeReportEvent, this, 5); // report
    if (event_index < 0)
    {
        LOG_ERROR("failed to register timer, error code %d", event_index);
        return event_index;
    }

    int startEventThreadVal = m_ptrEventThread->Start();
    if (GSE_SUCCESS > startEventThreadVal)
    {
        LOG_ERROR("failed to start event thread, error code %d", startEventThreadVal);
        return startEventThreadVal;
    }

    return GSE_SUCCESS;
}

void OpsCollection::SetAgentOpsChannelid(int channelid)
{
    m_agentOpsChannelID = channelid;
}

int OpsCollection::Open(int maxQueueCnt, int thread_count, bool enableops, std::string &reportIP, int channelID, std::vector<OpsAddress> &server_addresses)
{
    //上报到kafka的channelid
    LOG_DEBUG("OPS start, param(max quene:%d, enable:%d, report_ip:%s, channel_id:%d, thread count:%d)", maxQueueCnt, enableops, reportIP.c_str(), channelID, thread_count);
    m_channelID = channelID;
    m_reportIP = reportIP;
    m_enableOps = enableops;

    std::vector<OPSServerAddress> addresses;

    if (!server_addresses.empty())
    {
        for (int i = 0; i < server_addresses.size(); i++)
        {
            OPSServerAddress ops_addr(server_addresses[i].m_ip, server_addresses[i].m_port);
            addresses.push_back(ops_addr);
        }
        m_opsClient.CloseUDP();
        int ret = m_opsClient.OpenUDP(addresses);
        if (ret != GSE_SUCCESS)
        {
            int tmp_errno = gse_socket_errno;
            LOG_ERROR("failed to connect ops udp server, [errorno:%d, errormsg:%s]", tmp_errno, gse::tools::error::ErrnoToStr(tmp_errno).c_str());
            return GSE_ERROR;
        }

        m_hasOpsServer = true;
    }

    for (int i = 0; i < thread_count; i++)
    {
        OpsStatWorker *worker = new OpsStatWorker();
        worker->SetGseConf(m_configPtr);
        worker->SetOpsServerAddress(addresses);
        worker->SetOpsParam(reportIP, channelID);
        int ret = worker->Start();
        if (ret != GSE_SUCCESS)
        {
            LOG_ERROR("failed to start Ops stat worker thread");
            continue;
        }
        m_opsWorkers.push_back(worker);
    }
    return GSE_SUCCESS;
}

int OpsCollection::TransferOpsData(DataCell *ptrData)
{
    return m_opsClient.SendToByUDP(m_agentOpsChannelID, ptrData->GetDataBuf(), ptrData->GetDataBufLen());
}

bool OpsCollection::PushOpsQueue(DataCellOPS *ptr)
{
    if (m_opsWorkers.size() == 0)
    {
        return false;
    }

    std::string ops_key;
    // ptr->OpsKey(ops_key);
    //使用时间戳作为key进行计算hash，确保同一条消息在同一个ops的worker线程上进行数据汇总计算
    int worker_id = stringHash(ptr->m_szCreatedTimestamp.c_str(), ptr->m_szCreatedTimestamp.length()) % m_opsWorkers.size();
    LOG_DEBUG("push msg to ops worker(%d), key:%s", worker_id, ptr->m_szCreatedTimestamp.c_str());
    OpsStatWorker *worker = m_opsWorkers[worker_id];
    bool bret = worker->PushOpsQueue(ptr);
    if (!bret)
    {
        delete ptr;
        return false;
    }

    return true;
}

int OpsCollection::PutOpsData(DataCellOPS *ptrData)
{
    if (NULL == ptrData)
    {
        return GSE_SUCCESS;
    }

    if (!m_enableOps)
    {
        delete ptrData;
        return GSE_SUCCESS;
    }

    if (!m_hasOpsServer)
    {
        delete ptrData;
        return GSE_SUCCESS;
    }

    LOG_DEBUG("report ops data, srcip:%s, protocal:%s, status:%d", ptrData->m_szSrcIP.c_str(), ptrData->m_channelProtocol.c_str(), ptrData->m_status);

    PushOpsQueue(ptrData);

    return GSE_SUCCESS;
}

void OpsCollection::Close()
{
}
} // namespace data
} // namespace gse
