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
 * 查询后代
 * @param {number} bizSetId 业务集 ID
 * @param {string} parentModelId 上一级的模型 ID，通常为模型类型比如 set
 * @param {number} parentInstanceId 上一级的实例 ID
 * @returns {Promise}
 */
export const findChildren = ({
  bizSetId,
  parentModelId,
  parentInstanceId,
}, config) => $http.post('find/biz_set/topo_path', {
  bk_biz_set_id: bizSetId,
  bk_parent_obj_id: parentModelId,
  bk_parent_id: parentInstanceId,
}, config)


/**
 * 获取实例数量
 * @param {number} bizSetId 业务集 ID
 * @param {Array} condition 拓扑节点信息，数组最大长度为 20
 * @param {Object} config 请求配置
 * @returns {Promise}
 */
export const getInstanceCount = (bizSetId, condition, config) => $http.post(`count/topoinst/host_service_inst/biz_set/${bizSetId}`, { condition }, config)

/**
 * 获取业务下资源的拓扑路径
 * @param {number} bizSetId 业务集 ID
 * @param {number} bizId 业务 ID
 * @param {Array} condition 实例查询条件
 * @param {Object} config 请求配置
 * @returns {Promise}
 */
export const findTopoPath = ({ bizSetId, bizId }, params, config) => $http.post(`find/topopath/biz_set/${bizSetId}/biz/${bizId}`, params, config)

export const TopologyService = {
  findChildren,
  findTopoPath,
  getInstanceCount
}
