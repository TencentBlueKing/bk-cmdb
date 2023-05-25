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

// 创建模板
const create = (data, config) => http.post('create/field_template', data, config)

// 更新模板
const update = (data, config) => http.put('update/field_template', data, config)

// 查询模板列表
const find = (data, config) => http.post(`${window.API_HOST}findmany/field_template`, data, config)

// 删除属性配置
const deleteProperty = (data, config = {}) => http.delete('delete/topo/set_template/attribute', { ...config, data })

// 查询属性配置
const findProperty = (data, config) => http.post('findmany/topo/set_template/attribute', data, config)

export default {
  create,
  update,
  find,
  deleteProperty,
  findProperty
}
