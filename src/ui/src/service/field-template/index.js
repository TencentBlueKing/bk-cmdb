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
import { enableCount, rollReqUseCount } from '../utils.js'

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
const getFieldList = (data, config) => http.post('findmany/field_template/attribute', data, config)

// 查询模板唯一校验
const getUniqueList = (data, config) => http.post('findmany/field_template/unique', data, config)

// 查询模板绑定的模型
const getBindModel = (params, config) => rollReqUseCount('findmany/object/by_field_template', params, { limit: 100 }, config)

// 查询模型绑定的字段模版
const getTemplateList = (data, config) => http.post('findmany/field_template/by_object', data, config)

// 模版解绑模型接口
const updateTemplate = (data, config) => http.post('/update/field_template/unbind/object', data, config)

// 删除字段模板接口
const deleteTemplate = (data, config) => http.delete('/delete/field_template', data, config)


// 克隆字段模版接口
const cloneTemplate = (data, config) => http.post('/create/field_template/clone', data, config)

// 查询对应模版简要信息接口
const getTemplateInfo = (data, config) => http.post('/find/field_template/simplify/by_unique_template_id', data, config)

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
  getTemplateList,
  updateTemplate,
  deleteTemplate,
  cloneTemplate,
  getTemplateInfo
}
