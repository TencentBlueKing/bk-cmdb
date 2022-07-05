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

export const requestIds = {
  findmanyTemplate: Symbol('findmanyTemplate')
}

const find = async (params, config) => {
  try {
    const { count = 0, info: list = [] } = await http.post('findmany/proc/service_template', params, {
      requestId: requestIds.findmanyTemplate,
      ...config
    })
    return { count, list }
  } catch (error) {
    console.error(error)
    return { count: 0, list: [] }
  }
}

const findAll = async (params, config) => {
  try {
    let index = 1
    const size = 1000
    const results = []

    const page = index => ({ ...(params.page || {}), start: (index - 1) * size, limit: size })

    const req = index => http.post('findmany/proc/service_template', {
      ...params,
      page: page(index)
    }, config)

    const { count = 0, info: list = [] } = await req(index)
    results.push(...list)

    const max = Math.ceil(count / size)

    const reqs = []
    while (index < max) {
      index += 1
      reqs.push(req(index))
    }

    const rest = await Promise.all(reqs)
    rest.forEach(({ info: list = [] }) => {
      results.push(...list)
    })

    return results
  } catch (error) {
    console.error(error)
    return []
  }
}

const findAllByIds = async (ids, params, config) => {
  try {
    const size = 1000
    const max = Math.ceil(ids.length / size)

    const req = segment => http.post('findmany/proc/service_template', {
      ...params,
      service_template_ids: segment,
      page: { start: 0, limit: 999999999 } // NoLimit
    }, config)

    const reqs = []
    for (let index = 1; index <= max; index++) {
      const segment = ids.slice((index - 1) * size, size * index)
      reqs.push(req(segment))
    }

    const results = []
    const res = await Promise.all(reqs)
    res.forEach(({ info: list = [] }) => {
      results.push(...list)
    })

    return results
  } catch (error) {
    console.error(error)
    return []
  }
}

// 创建服务模板（全量）
const create = (data, config) => http.post('create/proc/service_template/all_info', data, config)

// 更新服务模板（全量）
const update = (data, config) => http.put('update/proc/service_template/all_info', data, config)

// 获取服务模板完整信息（包括属性设置、进程信息）
const getFullOne = (data, config) => http.post('find/proc/service_template/all_info', data, config)

// 更新属性配置
const updateProperty = (data, config) => http.put('update/proc/service_template/attribute', data, config)

// 删除属性配置
const deleteProperty = (data, config = {}) => http.delete('delete/proc/service_template/attribute', { ...config, data })

// 查询属性配置
const findProperty = (data, config) => http.post('findmany/proc/service_template/attribute', data, config)

export const CONFIG_MODE = {
  MODULE: 'module',
  TEMPLATE: 'template'
}

export default {
  find,
  findAll,
  findAllByIds,
  create,
  update,
  getFullOne,
  updateProperty,
  deleteProperty,
  findProperty
}
