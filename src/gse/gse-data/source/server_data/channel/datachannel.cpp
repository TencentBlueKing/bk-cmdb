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

#include "datachannel.h"

#include "receiver/httpreceiver.h"
//##include "receiver/kcpreceiver.h"
#include "bbx/gse_errno.h"
#include "dataserver.h"
#include "exporter/dsproxy/proxy_exporter.h"
#include "exporter/kafka/kafkaexporter.h"
#include "exporter/log/logexporter.h"
#include "exporter/pulsar/pulsarexporter.h"
#include "exporter/redis/redisexporter.h"
#include "protocol.h"
#include "receiver/tcpreceiver.h"
#include "receiver/udpreceiver.h"
#include "time_center.h"
#include "tools/macros.h"

namespace gse {
namespace data {

#define MAX_WORKER_NUM 64

DataChannel::DataChannel()
{
    m_channelConf = NULL;
    m_receiver = NULL;
    m_codec = NULL;
    m_nextWorkerId = 0;
    m_numWorker = 6;
    m_isKeepTiming = false;
    m_ptrEventThread = NULL;
    m_opsReport = NULL;
    m_conftor = NULL;
    m_ptrEventThread = NULL;
    m_ptrIDToStorage = NULL;
}

DataChannel::~DataChannel()
{
    Stop();

    if (m_channelConf != NULL)
    {
        delete m_channelConf;
        m_channelConf = NULL;
    }

    if (m_codec != NULL)
    {
        delete m_codec;
        m_codec = NULL;
    }
}

void DataChannel::cleanExporter(evutil_socket_t fd, short what, void *args)
{

    DataChannel *ptr_this = reinterpret_cast<DataChannel *>(args);

    int max_count = ptr_this->m_toDeletedExporter.Size();
    LOG_DEBUG("datachannel[%s] clean the invalid exporters by timer, invalid exporter size:%d", SAFE_CSTR(ptr_this->m_channelName.c_str()), max_count);

    while (--max_count >= 0)
    {
        Exporter *ptr_exporter = NULL;
        ptr_this->m_toDeletedExporter.Pop(ptr_exporter);
        if (NULL == ptr_exporter)
        {
            break;
        }

        if (!ptr_exporter->IsNeedStop())
        {
            LOG_DEBUG("push expoter(%p:%s) to delete list", ptr_exporter, ptr_exporter->GetExporterName().c_str());
            ptr_this->m_toDeletedExporter.Push(ptr_exporter);
            continue;
        }

        std::string name = ptr_exporter->GetExporterName();
        LOG_DEBUG("will stop the invalid exporter (%s) by timer", SAFE_CSTR(name.c_str()));
        ptr_exporter->Stop();
        LOG_DEBUG("real delete exporer %p, name:%s", ptr_exporter, ptr_exporter->GetExporterName().c_str());
        delete ptr_exporter;

        LOG_DEBUG("finish stop the exporter (%s) by timer", SAFE_CSTR(name.c_str()));
    }
}

int DataChannel::Start()
{
    if (NULL == m_channelConf)
    {
        LOG_WARN("the pointer of channel configure is NULL, can't start channel, please check channel configure");
        return GSE_SYSTEMERROR;
    }

    Stop();

    // init
    if (GSE_SUCCESS != init())
    {
        LOG_WARN("fail to init datachannel[%s]", SAFE_CSTR(m_channelName.c_str()));
    }

    // channel worker
    if (GSE_SUCCESS != StartChannelWorker())
    {
        LOG_WARN("fail to start channelworker for datachannel[%s]", SAFE_CSTR(m_channelName.c_str()));
        Stop();
        return GSE_ERROR;
    }

    // receiver
    if (GSE_SUCCESS != startReceiver())
    {
        LOG_WARN("fail to start receiver[%s] for datachannel[%s]", SAFE_CSTR(m_receiverName.c_str()), SAFE_CSTR(m_channelName.c_str()));
        Stop();
        return GSE_ERROR;
    }
    LOG_INFO("success to start receiver[%s] in channel[%s]", SAFE_CSTR(m_receiverName.c_str()), SAFE_CSTR(m_channelName.c_str()));

    if (GSE_SUCCESS != CreateDataflowExporter())
    {
        LOG_WARN("fail to start exporter for datachannel[%s]", SAFE_CSTR(m_channelName.c_str()));
        Stop();
        return GSE_ERROR;
    }

    LOG_DEBUG("datachannel[%s] has start. the size of channelworker is [%d]", SAFE_CSTR(m_channelName.c_str()), m_vChannelWorkers.size());

    m_ptrEventThread = new EventThread();
    int registerResponseVal = m_ptrEventThread->RegisterTimerPersistEvent(DataChannel::cleanExporter, this, 60); // 1 minutes
    if (GSE_SUCCESS > registerResponseVal)
    {
        LOG_ERROR("failed to register timer, error code %d", registerResponseVal);
        return registerResponseVal;
    }

    int startEventThreadVal = m_ptrEventThread->Start();
    if (GSE_SUCCESS > startEventThreadVal)
    {
        LOG_ERROR("failed to start event thread, error code %d", startEventThreadVal);
        return startEventThreadVal;
    }

    return GSE_SUCCESS;
}

int DataChannel::Stop()
{
    if (NULL != m_ptrEventThread)
    {
        m_ptrEventThread->Stop();
        m_ptrEventThread->Join();
        delete m_ptrEventThread;
        m_ptrEventThread = NULL;
    }

    // stop receiver
    StopReceiver();

    // stop channel worker
    StopChannelWorker();

    // stop exporter
    ClearDataflowExporter();

    return GSE_SUCCESS;
}

void DataChannel::Join()
{
    if (m_receiver != NULL)
    {
        m_receiver->Join();
    }

    ChannelWorker *pWorker = NULL;
    for (size_t i = 0; i < m_vChannelWorkers.size(); i++)
    {

        pWorker = m_vChannelWorkers[i];
        if (NULL == pWorker)
        {
            continue;
        }

        pWorker->Join();
        LOG_DEBUG("channel worker[%d] joined joined,", i);
    }
}

void DataChannel::SetOps(OpsCollection *ptr_ops_report)
{
    m_opsReport = ptr_ops_report;
}

void DataChannel::UpdateConf(const ChannelConf &channelConf)
{
    m_upConfLock.WLock();
    if (m_channelConf != NULL)
    {
        delete m_channelConf;
    }

    m_channelConf = new ChannelConf(channelConf);
    m_upConfLock.UnLock();
}

int DataChannel::DecodeMsg(DataCell *pDataCell)
{
    if (m_codec != NULL)
    {
        int iRet = m_codec->DecodeMsg(pDataCell);
        if (iRet != GSE_SUCCESS)
        {
            LOG_ERROR("Message decoding failed, channeid:%d, client ip:%s, port:%d",
                      pDataCell->GetChannelID(), pDataCell->GetSourceIp().c_str(), pDataCell->GetSourcePort());
            pDataCell->SetErrorMsg("decode failed", OPS_ERROR_DECODE);
            return iRet;
        }
    }

    return GSE_SUCCESS;
}

void DataChannel::dataflowExporterWrite(DataCell *pDataCell)
{
    for (size_t i = 0; i < m_vDataflowExporter.size(); i++)
    {
        Exporter *pExporter = m_vDataflowExporter[i];
        if (pExporter != NULL)
        {
            if (GSE_SUCCESS != pExporter->Write(pDataCell))
            {
                LOG_WARN("fail to write data into exporter[%s]", SAFE_CSTR(pExporter->GetExporterName().c_str()));
            }
        }
        // LOG_DEBUG("send data for the channel id (%d) by the exporter(%s)", pDataCell->GetChannelID(), SAFE_CSTR(pExporter->GetExporterName().c_str()));
    }
}

void DataChannel::SetDataCellChannelidByServiceId(DataCell *pDataCell)
{
    if (pDataCell->IsOpsMsg())
    {
        pDataCell->SetChannelID(m_conftor->GetChannelIdByOpsServiceId(pDataCell->GetOpsServiceId()));
    }
}

int DataChannel::ExportData(DataCell *pDataCell)
{
    pDataCell->SetOutputTimestamp(TimeCenter::Instance()->GetDateTime());

    dataflowExporterWrite(pDataCell);

    if (pDataCell->GetChannelID() == 0)
    {
        return GSE_SUCCESS;
    }

    if (DataServer::GetChannelIdManagerInst() != NULL)
    {
        if (!DataServer::GetChannelIdManagerInst()->WriteByFilter(pDataCell))
        {
            LOG_DEBUG("Write data msg to plat exporter");
            m_conftor->GetPlatIdManager()->WriteByFilter(pDataCell);
        }
        return GSE_SUCCESS;
    }

    return GSE_SUCCESS;
}

int DataChannel::init()
{
    if (NULL == m_channelConf)
    {
        LOG_WARN("the pointer of channel configure is NULL, can't start channel, please check channel configure");
        return GSE_SYSTEMERROR;
    }

    DecodeTypeEnum decodeType = D_TYPE_UNKNOWN;
    m_upConfLock.RLock();
    decodeType = m_channelConf->m_decodeType;
    m_channelName = m_channelConf->m_name;
    m_receiverName = m_channelConf->m_receiverName;
    m_upConfLock.UnLock();

    m_codec = CodecFactory::CreateCodec(decodeType);

    if (NULL != m_codec)
    {
    }

    return GSE_SUCCESS;
}

int DataChannel::StartChannelWorker()
{
    if (NULL == m_channelConf)
    {
        LOG_WARN("the pointer of channel configure is NULL, can't start channel, please check channel configure");
        return GSE_SYSTEMERROR;
    }

    m_upConfLock.RLock();
    m_numWorker = m_channelConf->m_numChannelWorker;
    m_upConfLock.UnLock();
    // start channel worker
    if (m_numWorker <= 0)
    {
        m_numWorker = 6;
    }

    if (m_numWorker > MAX_WORKER_NUM)
    {
        LOG_WARN("the config of channel worker number[%d] is more then MAX_WORKER_NUM[%d]. will set channel woker number to MAX_WORKER_NUM", m_numWorker, MAX_WORKER_NUM);
        m_numWorker = MAX_WORKER_NUM;
    }

    uint32_t workerRealNum = m_numWorker;

    for (uint32_t i = 0; i < m_numWorker; i++)
    {
        ChannelWorker *pWorker = new ChannelWorker(i, this);
        if (GSE_SUCCESS != pWorker->Init())
        {
            LOG_WARN("fail to initialize the channel worker[%d]", i);
            delete pWorker;
            workerRealNum -= 1;
            continue;
        }

        if (GSE_SUCCESS != pWorker->Start())
        {
            LOG_WARN("fail to start the channel worker[%d]", i);
            delete pWorker;
            workerRealNum -= 1;
            continue;
        }

        m_vChannelWorkers.push_back(pWorker);
    }

    if (workerRealNum <= 0)
    {
        LOG_WARN("fail to start channel worker, there is zero worker be created");
        return GSE_ERROR;
    }

    m_numWorker = workerRealNum;

    return GSE_SUCCESS;
}

int DataChannel::StopChannelWorker()
{
    ChannelWorker *pWorker = NULL;
    for (size_t i = 0; i < m_vChannelWorkers.size(); i++)
    {
        pWorker = m_vChannelWorkers[i];
        if (NULL == pWorker)
        {
            continue;
        }

        pWorker->stop();
        pWorker->Join();
        delete pWorker;
        LOG_DEBUG("channel work[%i] stop", i);
    }
    m_vChannelWorkers.clear();

    return GSE_SUCCESS;
}

int DataChannel::startReceiver()
{
    ReceiverProtocolEnum recvProto = R_PROTO_UNKNOWN;
    m_upConfLock.RLock();
    if (m_channelConf->m_receiverConf != NULL)
    {
        recvProto = m_channelConf->m_receiverConf->m_protocol;
    }
    switch (recvProto)
    {
    case R_PROTO_TCP:
    {
        TcpReceiver *tcpReceiver = new TcpReceiver();
        tcpReceiver->SetSheduler(DataServer::GetScheduler());
        m_receiver = tcpReceiver;
    }
    break;
    case R_PROTO_UDP:
    {
        m_receiver = new UdpReceiver();
    }
    break;
    case R_PROTO_KCP:
    {
        // m_receiver = new KcpReceiver();
    }
    break;
    case R_PROTO_HTTP:
    {
        m_receiver = new HttpReceiver();
    }
    break;
    default:
        LOG_WARN("unknow the protocol[%d] type of receiver[%s] in channel[%s]", recvProto, SAFE_CSTR(m_receiverName.c_str()), SAFE_CSTR(m_channelName.c_str()));
        break;
    }

    if (m_receiver != NULL)
    {
        m_receiver->SetConfig(m_configPtr);
        m_receiver->UpdateConf(*m_channelConf->m_receiverConf);
    }

    m_upConfLock.UnLock();
    LOG_DEBUG("start receiver:%s", m_channelConf->m_receiverConf->m_name.c_str());
    if (NULL == m_receiver)
    {
        LOG_ERROR("failed to start receiver:%s", m_channelConf->m_receiverConf->m_name.c_str());
        return GSE_ERROR;
    }

    m_receiver->SetRecvCallBack(DataChannel::handlRecvDataCell, this);

    return m_receiver->Start();
}

int DataChannel::StopReceiver()
{
    if (m_receiver != nullptr)
    {
        m_receiver->Stop();

        delete m_receiver;
        m_receiver = nullptr;
    }

    return GSE_SUCCESS;
}

int DataChannel::CreateDataflowExporter()
{
    m_upConfLock.WLock();
    for (size_t i = 0; i < m_channelConf->m_vExporterConf.size(); i++)
    {
        ExporterConf *pExporterConf = m_channelConf->m_vExporterConf[i];
        Exporter *pExporter = NULL;
        switch (pExporterConf->m_type)
        {
        case E_TYPE_LOG:
        {
            pExporter = new LogExporter();
            LOG_DEBUG("create log exporter(%s) by dataflow config", pExporterConf->m_name.c_str());
        }
        break;
        case E_TYPE_BKDATA:
        {
            //
        }
        break;
        case E_TYPE_KAFKA:
        {
            pExporter = new KafkaExporter();
            LOG_DEBUG("create kafka exporter(%s) by dataflow config", pExporterConf->m_name.c_str());
        }
        break;
        case E_TYPE_REDIS:
        {
            pExporter = new RedisExporter();
            LOG_DEBUG("create redis exporter(%s) by dataflow config", pExporterConf->m_name.c_str());
        }
        break;
        case E_TYPE_DS_PROXY:
        {
            pExporter = new ProxyExporter();
            LOG_DEBUG("create proxy exporter(%s) by dataflow config", pExporterConf->m_name.c_str());
        }
        break;

        case E_TYPE_PULSAR:
        {
            pExporter = new PulsarExporter();
            LOG_DEBUG("create pulsar exporter(%s) by dataflow config", pExporterConf->m_name.c_str());
        }
        break;
        default:
            LOG_WARN("unknow the type[%d] of exporter[%s] in channel[%s]", pExporterConf->m_type, SAFE_CSTR(pExporterConf->m_name.c_str()), SAFE_CSTR(m_channelName.c_str()));
            break;
        }

        if (pExporter != NULL)
        {
            pExporter->SetGseConf(m_configPtr);
            pExporter->UpdateConf(ExporterConfigWrapper::DataFlowConfType, pExporterConf);
            if (GSE_SUCCESS != pExporter->Start())
            {
                LOG_WARN("fail to start exporter[%s] in channel[%s]", SAFE_CSTR(pExporterConf->m_name.c_str()), SAFE_CSTR(m_channelName.c_str()));
                delete pExporter;
                continue;
            }
            else
            {
                m_vDataflowExporter.push_back(pExporter);
                LOG_INFO("success to start exporter[%s] in channel[%s]", SAFE_CSTR(pExporterConf->m_name.c_str()), SAFE_CSTR(m_channelName.c_str()));
            }
        }
        else
        {
            LOG_ERROR("failed to create exporter");
        }
    }
    m_upConfLock.UnLock();

    return GSE_SUCCESS;
}

void DataChannel::ClearDataflowExporter()
{
    Exporter *pExporter = NULL;
    for (size_t i = 0; i < m_vDataflowExporter.size(); i++)
    {
        pExporter = m_vDataflowExporter[i];
        if (NULL == pExporter)
        {
            continue;
        }

        pExporter->Stop();
        delete pExporter;
    }

    m_vDataflowExporter.clear();
}

void DataChannel::handlRecvDataCell(DataCell *pDataCell, void *pCaller)
{
    DataChannel *pDataChannel = (DataChannel *)pCaller;
    pDataChannel->localHanlerRecvDataCell(pDataCell);
}

void DataChannel::localHanlerRecvDataCell(DataCell *pDataCell)
{
    if (m_isKeepTiming)
    {
        directExport(pDataCell);
    }
    else
    {
        dispatchChannelWorker(pDataCell);
    }
}

int DataChannel::directExport(DataCell *pDataCell)
{
    int iRet = GSE_SUCCESS;
    iRet = DecodeMsg(pDataCell);
    if (iRet != GSE_SUCCESS)
    {
        delete pDataCell;
        return iRet;
    }

    iRet = ExportData(pDataCell);
    if (iRet != GSE_SUCCESS)
    {
        delete pDataCell;
        return iRet;
    }

    // success
    delete pDataCell;
    return GSE_SUCCESS;
}

int DataChannel::dispatchChannelWorker(DataCell *pDataCell)
{
    ChannelWorker *pWorker = NULL;
    int workerIndex = (m_nextWorkerId++) % m_numWorker;
    pWorker = m_vChannelWorkers[workerIndex];
    if (NULL == pWorker)
    {
        LOG_WARN("the worker[%d] is invalid, failed to dispatch datacell[0x%x] to worker", workerIndex, pDataCell);
        delete pDataCell;
        return GSE_ERROR;
    }

    pWorker->Notify(pDataCell);
    return GSE_SUCCESS;
}
} // namespace data
} // namespace gse
