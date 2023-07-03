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
const getFieldCount = (data, config) => http.post('findmany/field_template/attribute/count', data, config)

// 查询模板绑定的模型数量
const getModelCount = (data, config) => http.post(`${window.API_HOST}findmany/field_template/object/count`, data, config)

// 查询模板简要信息
const findById = (id, config = {}) => http.get(`find/field_template/${id}`, config)

// 查询模板字段，查看权限
const getFieldList = (data, config) => http.post('findmany/field_template/attribute', {
  page: {
    ...maxPageParams(MAX_FIELD_COUNT),
    sort: 'bk_property_index'
  },
  ...data,
}, config)

// 查询模板唯一校验，查看权限
const getUniqueList = (data, config) => http.post('findmany/field_template/unique', data, config)

// 查询模板绑定的模型
const getBindModel = (params, config) => rollReqUseCount('findmany/object/by_field_template', params, { limit: 100 }, config)

// 滚动查询模板绑定的模型
const rollBindModel = (params, limit, config) => rollReqUseCount(
  'findmany/object/by_field_template',
  params,
  { limit, generator: true },
  config
)

// 查询模型绑定的字段模版
const getModelBindTemplate = (data, config) => http.post('findmany/field_template/by_object', data, config)

// 模版解绑模型接口
const unbind = (data, config) => http.post('/update/field_template/unbind/object', data, config)

// 删除字段模板接口
const deleteTemplate = (data, config) => http.delete('/delete/field_template', data, config)

// 克隆字段模版接口
const cloneTemplate = (data, config) => http.post('/create/field_template/clone', data, config)

// 查询模型属性对应模版简要信息接口
const getFieldBindTemplate = (data, config) => http.post('/find/field_template/simplify/by_attr_template_id', data, config)

// 检测字段模板和模型中的字段的区别和冲突，返回值以模型上的字段为维度
const getFieldDifference = (data, config) => http.post('find/field_template/attribute/difference', data, config)

// 检测字段模板和模型中的唯一校验的区别和冲突，返回值以模型上的唯一校验为维度
const getUniqueDifference = (data, config) => http.post('find/field_template/unique/difference', data, config)

// 绑定模型
const bindModel = (data, config) => http.post('update/field_template/bind/object', data, config)

// 提交同步至模型任务
const syncModel = (data, config) => http.post('update/topo/field_template/sync', data, config)

// 查询同步任务状态
const getTaskSyncStatus = (data, config) => http.post('find/field_template/tasks_status', data, config)

// 查询同步任务状态
const getModelSyncStatus = (data, config) => http.post('find/field_template/model/status', data, config)

// 查询模板与模型差异状态
const getModelDiffStatus = (data, config) => http.post('find/field_template/sync/status', data, config)

// 获取字段绑定的模板的字段列表
const getTemplateFieldListByField = async (field, config) => {
  try {
    const bindTemplate = await getFieldBindTemplate({
      bk_template_id: field.bk_template_id,
      bk_attribute_id: field.id
    }, config)
    const templateFieldList = await getFieldList({ bk_template_id: bindTemplate.id })
    return templateFieldList?.info ?? []
  } catch (error) {
    return Promise.reject(error)
  }
}

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
  rollBindModel,
  getModelBindTemplate,
  unbind,
  deleteTemplate,
  cloneTemplate,
  getFieldBindTemplate,
  getFieldDifference,
  getUniqueDifference,
  bindModel,
  syncModel,
  getTaskSyncStatus,
  getModelSyncStatus,
  getModelDiffStatus,
  getTemplateFieldListByField
}
