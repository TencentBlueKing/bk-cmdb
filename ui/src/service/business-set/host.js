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
import has from 'has'
import { transformHostSearchParams, localSort } from '@/utils/tools'

/**
 * 查询业务集下的主机列表
 * @param {number} bizSetId 业务集 ID
 * @param {Object} params 查询参数
 * @param {Obejct} config 请求配置
 * @returns {Promise}
 */
export const findAll = (bizSetId, params, config) => $http.post(`findmany/hosts/biz_set/${bizSetId}`, transformHostSearchParams(params), config).then((data) => {
  if (has(data, 'info')) {
    data.info.forEach((host) => {
      localSort(host.module, 'bk_module_name')
      localSort(host.set, 'bk_set_name')
    })
  }
  return data
})


/**
 * 查询单个主机
 * @param {number} bizSetId 业务集 ID
 * @param {Object} params 查询参数
 * @param {Obejct} config 请求配置
 * @returns {Promise}
 */
export const findOne = (bizSetId, params, config) => $http.post(
  `findmany/hosts/biz_set/${bizSetId}`,
  params,
  config
)

export const HostService = {
  findAll,
  findOne
}
