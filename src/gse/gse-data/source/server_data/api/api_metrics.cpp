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

#include "api_metrics.h"

#include "bbx/gse_errno.h"
#include "bbx/prometheus/prometheus_exporter.h"
#include "net/http/http_server.hpp"
#include "tools/error.h"
#include "tools/strings.h"

#include "dataserver.h"

namespace gse {
namespace data {

std::shared_ptr<ChannelIdApiMetrics> ChannelIdApiMetricsInst()
{
    return DataServer::GetChannelIdApiMetricsInst();
}

ChannelIdApiMetrics::ChannelIdApiMetrics(std::shared_ptr<DataProcessConfig> configPtr)
{
    m_registry = std::make_shared<prometheus::Registry>();
    m_labels = {
        {"module", "data"},
        {"node", configPtr->GetAdvertiseIp()},
        {"cluster_name", configPtr->GetDataSvrClusterName()}};
}

bool ChannelIdApiMetrics::Start(std::shared_ptr<gse::net::http::HTTPServer> htpServer)
{
    auto exporter = OPMetric::OPMetricInst()->GetPromethusExporter();
    exporter->RegisterMetricCollectable(m_registry);
    return true;
}

void ChannelIdApiMetrics::CounterMetrics(const std::string &metricName, const std::string &apiName, int count, int errorCode)
{
    int err = 0;
    // runtime error stat as error
    if (GSE_DATA_ERROR_RUNTIME_ERROR == errorCode)
    {
        err = errorCode;
    }

    auto &counterFamily = prometheus::BuildCounter().Name(metricName).Labels(m_labels).Register(*m_registry);
    std::map<std::string, std::string> codeLabel = {{"code", gse::tools::strings::ToString(err)},
                                                    {"apiname", apiName}};
    auto &counter = counterFamily.Add(codeLabel);
    counter.Increment(count);
}

void ChannelIdApiMetrics::GuageMetrics(const string &metricName, const std::string &apiName, int count, int errorCode)
{
    auto &gaugeFamily = prometheus::BuildGauge().Name(metricName).Labels(m_labels).Register(*m_registry);
    std::map<std::string, std::string> codeLabel = {{"code", gse::tools::strings::ToString(errorCode)},
                                                    {"apiname", apiName}};
    auto &gauge = gaugeFamily.Add(codeLabel);
    gauge.Set(count);
}

void ChannelIdApiMetrics::SumMetrics(const string &metricName, const string &apiName, int count, int errorCode)
{
    auto &sumFamily = prometheus::BuildSummary()
                          .Name(metricName)
                          .Labels(m_labels)
                          .Register(*m_registry);

    std::map<std::string, std::string> codeLabel = {{"code", gse::tools::strings::ToString(errorCode)},
                                                    {"api", apiName}};
    auto &sum = sumFamily.Add(codeLabel, prometheus::Summary::Quantiles{{0.5, 0.05}, {0.9, 0.01}, {0.99, 0.001}});
    sum.Observe(count);
}

} // namespace data
} // namespace gse
