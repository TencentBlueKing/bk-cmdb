/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

import $http from '@/api'

/**
 * 查询业务集下的服务进程列表
 * @param {*} bizSetId
 * @param {*} params
 * @param {*} config
 * @returns
 */
export const findAll = (bizSetId, params, config) => $http.post(`findmany/proc/biz_set/${bizSetId}/process_instance/name_ids`, params, config)

/**
 * 根据服务实例查询进程
 * @param {number} bizSetId 业务集 ID
 * @param {Object} params 查询参数
 * @param {Object} config 请求配置
 * @returns {Promise}
 */
export const findProcessByServiceInstance = (bizSetId, params, config) => $http.post(`findmany/proc/biz_set/${bizSetId}/process_instance`, params, config)

/**
 * 根据进程查询对应的服务实例
 * @param {number} bizSetId 业务集 ID
 * @param {Object} params 查询参数
 * @param {Object} config 请求配置
 * @returns {Promise}
 */
export const findServiceInstanceByProcess = (bizSetId, params, config) => $http.post(`/findmany/proc/biz_set/${bizSetId}/process_instance/detail/by_ids`, params, config)

export const ProcessInstanceService = {
  findAll,
  findProcessByServiceInstance,
  findServiceInstanceByProcess
}
