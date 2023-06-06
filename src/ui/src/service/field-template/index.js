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
import { enableCount, rollReqUseCount, maxPageParams } from '../utils.js'
import { MAX_FIELD_COUNT } from '@/views/field-template/children/use-field.js'

// 创建模板
const create = (data, config) => http.post('create/field_template', data, config)

// 更新模板
const update = (data, config) => http.put('update/field_template', data, config)

// 更新模板基础信息
const updateBaseInfo = (data, config) => http.put('update/field_template/info', data, config)

// 查询模板列表
const find = async (params, config) => {
  const api = `${window.API_HOST}findmany/field_template`
  try {
    const [{ info: list = [] }, { count = 0 }] = await Promise.all([
      http.post(api, enableCount(params, false), config),
      http.post(api, enableCount(params, true), config)
    ])
    return { count: count || 0, list: list || [] }
  } catch (error) {
    console.error(error)
    return Promise.reject(error)
  }
}

// 查询模板字段数量
const getFieldCount = (data, config) => http.post(`${window.API_HOST}findmany/field_template/attribute/count`, data, config)

// 查询模板绑定的模型数量
const getModelCount = (data, config) => http.post(`${window.API_HOST}findmany/field_template/object/count`, data, config)

// 查询模板简要信息
const findById = (id, config = {}) => http.get(`find/field_template/${id}`, config)

// 查询模板字段
const getFieldList = (data, config) => http.post('findmany/field_template/attribute', {
  page: {
    ...maxPageParams(MAX_FIELD_COUNT),
    sort: 'bk_property_index'
  },
  ...data,
}, config)

// 查询模板唯一校验
const getUniqueList = (data, config) => http.post('findmany/field_template/unique', data, config)

// 查询模板绑定的模型
const getBindModel = (params, config) => rollReqUseCount('findmany/object/by_field_template', params, { limit: 100 }, config)

// 检测字段模板和模型中的字段的区别和冲突，返回值以模型上的字段为维度
const getFieldDifference = (data, config) => http.post('find/field_template/attribute/difference', data, config)

// 检测字段模板和模型中的唯一校验的区别和冲突，返回值以模型上的唯一校验为维度
const getUniqueDifference = (data, config) => http.post('find/field_template/unique/difference', data, config)

// 绑定模型
const bindModel = (data, config) => http.post('update/field_template/bind/object', data, config)

// 提交同步至模型任务
const syncModel = (data, config) => http.post('update/topo/field_template/sync', data, config)

// 查询同步状态
const getSyncStatus = (data, config) => http.post('find/field_template/tasks_status', data, config)

export default {
  create,
  update,
  updateBaseInfo,
  find,
  getFieldCount,
  getModelCount,
  findById,
  getFieldList,
  getUniqueList,
  getBindModel,
  getFieldDifference,
  getUniqueDifference,
  bindModel,
  syncModel,
  getSyncStatus
}
