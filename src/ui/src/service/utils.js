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

import merge from 'lodash/merge'
import http from '@/api'

/**
 * 根据是否开启count生成新参数
 * @param {Object} params 基础参数
 * @param {Boolean} flag 是否开启count获取
 * @returns 生成的新参数
 */
export const enableCount = (params = {}, flag = false) => {
  const page = Object.assign(flag ? { start: 0, limit: 0, sort: '' } : {}, { enable_count: flag })
  return merge({}, params, { page })
}

export const onePageParams = () => ({ start: 0, limit: 1 })

export const maxPageParams = (max = 500) => ({ start: 0, limit: max })

// 使用enableCount的方式滚动拉取接口数据（先获取总数）
export const rollReqUseCount = async (
  url,
  params = {},
  options = {},
  config = {},
  method = 'post'
) => {
  const { start = 1, limit = 1000, countKey = 'count', listKey = 'info' } = options

  let index = start
  const size = limit
  const results = []

  // 先获取到总数
  const { [countKey]: total = 0 } = await http[method](url, enableCount(params, true), config)

  // 分页方法
  const page = index => ({ ...(params.page || {}), start: (index - 1) * size, limit: size })

  // 请求列表的req
  const req = index => http[method](url, enableCount({
    ...params,
    page: page(index)
  }, false), config)

  // 列表一共要拉取多少次
  const max = Math.ceil(total / size)

  // 循环组装得到所有的req
  const reqs = []
  while (index <= max) {
    reqs.push(req(index))
    index += 1
  }

  const all = await Promise.all(reqs)
  all.forEach(({ [listKey]: list = [] }) => {
    results.push(...list)
  })

  return results
}

// 常规方法滚动拉取接口数据（先拉取第1页）
export const rollReq = async (
  url,
  params = {},
  options = {},
  config = {},
  method = 'post'
) => {
  const { start = 1, limit = 1000, countKey = 'count', listKey = 'info' } = options

  let index = start
  const size = limit
  const results = []

  const page = index => ({ ...(params.page || {}), start: (index - 1) * size, limit: size })

  const req = index => http[method](url, {
    ...params,
    page: page(index)
  }, config)

  // 先拉起始页，通常是第1页，同时得到总数
  const { [countKey]: total = 0, [listKey]: list = [] } = await req(index)
  results.push(...list)

  // 一共要拉取多少次
  const max = Math.ceil(total / size)

  const reqs = []
  while (index < max) {
    index += 1
    reqs.push(req(index))
  }

  const rest = await Promise.all(reqs)
  rest.forEach(({ [listKey]: list = [] }) => {
    results.push(...list)
  })

  return results
}

// 给定指定的总数自动分页获取所有数据
export const rollReqUseTotalCount = async (
  url,
  params = {},
  options = {},
  config = {},
  method = 'post',
  getter
) => {
  const { start = 1, limit = 1000, total, listKey = 'info' } = options

  let index = start
  const size = limit
  const results = []

  // 分页方法
  const page = index => ({ ...(params.page || {}), start: (index - 1) * size, limit: size })

  // 请求列表的req
  const req = index => http[method](url, enableCount({
    ...params,
    page: page(index)
  }, false), config)

  // 列表一共要拉取多少次
  const max = Math.ceil(total / size)

  // 循环组装得到所有的req
  const reqs = []
  while (index <= max) {
    reqs.push(req(index))
    index += 1
  }

  const all = await Promise.all(reqs)

  if (getter) {
    all.forEach((data) => {
      results.push(getter(data))
    })
  } else {
    all.forEach(({ [listKey]: list = [] }) => {
      results.push(...list)
    })
  }

  return results
}

// 通过分割指定的参数数据拆分请求
export const rollReqByDataKey = async (
  url,
  params = {},
  options = {},
  config = {},
  method = 'post',
  getter
) => {
  const { limit = 1000, dataKey = 'ids' } = options

  const results = []

  let index = 0

  // 循环组装得到所有的req
  const reqs = []

  while (index < params?.[dataKey]?.length) {
    reqs.push(http[method](url, {
      ...params,
      [dataKey]: params?.[dataKey]?.slice(index, index + limit)
    }, config))
    index = index + limit
  }

  const all = await Promise.all(reqs)

  if (getter) {
    all.forEach((data) => {
      results.push(getter(data))
    })
  } else {
    all.forEach((data) => {
      results.push(...data)
    })
  }

  return results
}
