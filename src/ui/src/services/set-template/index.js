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

import http from '@/api'

// 创建模板（全量）
const create = (data, config) => http.post('create/topo/set_template/all_info', data, config)

// 更新模板（全量）
const update = (data, config) => http.put('update/topo/set_template/all_info', data, config)

// 获取模板完整信息（包括属性设置、服务拓扑）
const getFullOne = (data, config) => http.post('find/topo/set_template/all_info', data, config)

// 更新属性配置
const updateProperty = (data, config) => http.put('update/topo/set_template/attribute', data, config)

// 删除属性配置
const deleteProperty = (data, config = {}) => http.delete('delete/topo/set_template/attribute', { ...config, data })

// 查询属性配置
const findProperty = (data, config) => http.post('findmany/topo/set_template/attribute', data, config)

// 获取模板与实例对比中被移除的模块是否存在主机
const getRemovedModuleStatus = (bizId, templateId, data, config) => http.post(`findmany/topo/set_template/${templateId}/bk_biz_id/${bizId}/host_with_instances`, data, config)

export default {
  create,
  update,
  getFullOne,
  updateProperty,
  deleteProperty,
  findProperty,
  getRemovedModuleStatus
}
