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

#ifndef _GSE_DATA_LOGGER_HPP_
#define _GSE_DATA_LOGGER_HPP_

#include "log/log.h"
#include "log/log_base.h"

/* business log */
#define BLOG_DEBUG(rid, fmt, ...) \
    LOG_FORMAT(gse::log::Log::Instance().GetBusinessLogChannelID(), gse::log::LOG_LEVEL_DEBUG, true, (rid), gse::tools::filesystem::GetBasename(__FILE__), __LINE__, fmt, ##__VA_ARGS__)

#define BLOG_INFO(rid, fmt, ...) \
    LOG_FORMAT(gse::log::Log::Instance().GetBusinessLogChannelID(), gse::log::LOG_LEVEL_INFO, true, (rid), gse::tools::filesystem::GetBasename(__FILE__), __LINE__, fmt, ##__VA_ARGS__)

#define BLOG_WARN(rid, fmt, ...) \
    LOG_FORMAT(gse::log::Log::Instance().GetBusinessLogChannelID(), gse::log::LOG_LEVEL_WARN, true, (rid), gse::tools::filesystem::GetBasename(__FILE__), __LINE__, fmt, ##__VA_ARGS__)

#define BLOG_ERROR(rid, fmt, ...) \
    LOG_FORMAT(gse::log::Log::Instance().GetBusinessLogChannelID(), gse::log::LOG_LEVEL_ERROR, true, (rid), gse::tools::filesystem::GetBasename(__FILE__), __LINE__, fmt, ##__VA_ARGS__)

// key info log, used for business api access log, do error log when it's not zero code.
#define BLOG_INFOK(code, rid, fmt, ...)                                                                                                                                                           \
    do                                                                                                                                                                                            \
    {                                                                                                                                                                                             \
        if (code)                                                                                                                                                                                 \
        {                                                                                                                                                                                         \
            LOG_FORMAT(gse::log::Log::Instance().GetBusinessLogChannelID(), gse::log::LOG_LEVEL_ERROR, true, (rid), gse::tools::filesystem::GetBasename(__FILE__), __LINE__, fmt, ##__VA_ARGS__); \
        }                                                                                                                                                                                         \
        else                                                                                                                                                                                      \
        {                                                                                                                                                                                         \
            LOG_FORMAT(gse::log::Log::Instance().GetBusinessLogChannelID(), gse::log::LOG_LEVEL_INFO, true, (rid), gse::tools::filesystem::GetBasename(__FILE__), __LINE__, fmt, ##__VA_ARGS__);  \
        }                                                                                                                                                                                         \
    } while (0)

// BLOG_FATAL real FATAL level log.
#define BLOG_FATAL(rid, fmt, ...)                                                                                                                                                             \
    do                                                                                                                                                                                        \
    {                                                                                                                                                                                         \
        LOG_FORMAT(gse::log::Log::Instance().GetBusinessLogChannelID(), gse::log::LOG_LEVEL_FATAL, true, (rid), gse::tools::filesystem::GetBasename(__FILE__), __LINE__, fmt, ##__VA_ARGS__); \
        std::cerr << ("fatal:" + gse::tools::filesystem::GetBasename(__FILE__) + ":" + std::to_string(__LINE__)) << std::endl;                                                                \
        exit(1);                                                                                                                                                                              \
    } while (0)

/* system log */
#define SLOG_DEBUG(fmt, ...) \
    LOG_FORMAT(gse::log::Log::Instance().GetSystemLogChannelID(), gse::log::LOG_LEVEL_DEBUG, true, "0", gse::tools::filesystem::GetBasename(__FILE__), __LINE__, fmt, ##__VA_ARGS__)

#define SLOG_INFO(fmt, ...) \
    LOG_FORMAT(gse::log::Log::Instance().GetSystemLogChannelID(), gse::log::LOG_LEVEL_INFO, true, "0", gse::tools::filesystem::GetBasename(__FILE__), __LINE__, fmt, ##__VA_ARGS__)

#define SLOG_WARN(fmt, ...) \
    LOG_FORMAT(gse::log::Log::Instance().GetSystemLogChannelID(), gse::log::LOG_LEVEL_WARN, true, "0", gse::tools::filesystem::GetBasename(__FILE__), __LINE__, fmt, ##__VA_ARGS__)

#define SLOG_ERROR(fmt, ...) \
    LOG_FORMAT(gse::log::Log::Instance().GetSystemLogChannelID(), gse::log::LOG_LEVEL_ERROR, true, "0", gse::tools::filesystem::GetBasename(__FILE__), __LINE__, fmt, ##__VA_ARGS__)

// SLOG_FATAL real FATAL level log.
#define SLOG_FATAL(fmt, ...)                                                                                                                                                              \
    do                                                                                                                                                                                    \
    {                                                                                                                                                                                     \
        LOG_FORMAT(gse::log::Log::Instance().GetSystemLogChannelID(), gse::log::LOG_LEVEL_FATAL, true, "0", gse::tools::filesystem::GetBasename(__FILE__), __LINE__, fmt, ##__VA_ARGS__); \
        std::cerr << ("fatal:" + gse::tools::filesystem::GetBasename(__FILE__) + ":" + std::to_string(__LINE__)) << std::endl;                                                            \
        exit(1);                                                                                                                                                                          \
    } while (0)

#endif // _GSE_CLUSTER_LOGGER_HPP_
