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

#ifndef _GSE_DATA_API_METRICS_H_
#define _GSE_DATA_API_METRICS_H_

#include <map>
#include <string>

#include <prometheus/counter.h>
#include <prometheus/gauge.h>
#include <prometheus/registry.h>

#include "conf/dataconf.h"
#include "net/http/http_server.hpp"

namespace gse {
namespace data {

class ChannelIdApiMetrics;
std::shared_ptr<ChannelIdApiMetrics> ChannelIdApiMetricsInst();

class ChannelIdApiMetrics
{
public:
    explicit ChannelIdApiMetrics(std::shared_ptr<DataProcessConfig> configPtr);
    explicit ChannelIdApiMetrics(std::shared_ptr<gse::net::http::HTTPServer> httpServer);
    bool Start(std::shared_ptr<gse::net::http::HTTPServer> htpServer);

    void CounterMetrics(const string &metricName, const std::string &apiName, int count, int errorCode);
    void GuageMetrics(const string &metricName, const std::string &apiName, int count, int errorCode);
    void SumMetrics(const string &metricName, const std::string &apiName, int count, int errorCode);

private:
    std::shared_ptr<prometheus::Registry> m_registry;
    std::map<std::string, std::string> m_labels;
};

} // namespace data
} // namespace gse
#endif // API_METRICS_H
