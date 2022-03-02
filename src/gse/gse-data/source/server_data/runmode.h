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

#ifndef __GSE_DATASERVER_RUN_MODE_H__
#define __GSE_DATASERVER_RUN_MODE_H__

namespace gse {
namespace dataserver {
/**
 * @brief DataServer 运行模式定义
 */
enum DataServerRunMode
{

    /**
     * @brief DataServer 模式
     */
    DATASERVER_RUN = 0x01,

    /**
     * @brief OPS 模式
     */
    DATASERVER_OPS = 0x02,

    /**
     * @brief REST API 模式
     */
    DATASERVER_RESTAPI = 0x03,

    // porxy mode
    DATASERVER_PROXY = 0x04,

    // thirdparty mode
    DATASERVER_TGLOG = 0x05,
    DATASERVER_BEACON = 0x06

};

}
}
#endif

