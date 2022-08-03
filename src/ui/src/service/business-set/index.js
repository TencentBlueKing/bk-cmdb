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
import { BUILTIN_MODELS, BUILTIN_MODEL_PROPERTY_KEYS } from '@/dictionary/model-constants.js'
import { enableCount, onePageParams } from '../utils.js'

const authorizedRequsetId = Symbol('getAuthorizedBusinessSet')
const MODEL_ID_KEY = BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS_SET].ID

const find = async (params, config) => {
  try {
    const [{ info: list = [] }, { count = 0 }] = await Promise.all([
      http.post('findmany/biz_set', enableCount(params, false), config),
      http.post('findmany/biz_set', enableCount(params, true), config)
    ])
    return { count: count || 0, list: list || [] }
  } catch (error) {
    console.error(error)
    return Promise.reject(error)
  }
}

const findById = async (id, config = {}) => {
  try {
    const { info: [instance = null] } = await http.post('findmany/biz_set', enableCount({
      bk_biz_set_filter: {
        condition: 'AND',
        rules: [{
          field: MODEL_ID_KEY,
          operator: 'equal',
          value: id
        }]
      },
      page: onePageParams()
    }, false), config)

    return instance
  } catch (error) {
    console.error(error)
    return null
  }
}

const findOne = async (params, config = {}) => findById(params[MODEL_ID_KEY], config)

const getAuthorized = async (config) => {
  try {
    const { info: list = [] } = await http.get('findmany/biz_set/with_reduced?sort=bk_biz_set_id', config)
    return list || []
  } catch (error) {
    console.error(error)
    return []
  }
}

const getAuthorizedWithCache = async () => getAuthorized({
  requestId: authorizedRequsetId,
  fromCache: true
})

const previewOfBeforeCreate = async (params, config) => {
  try {
    const [{ info: list = [] }, { count = 0 }] = await Promise.all([
      http.post('find/biz_set/preview', enableCount(params, false), config),
      http.post('find/biz_set/preview', enableCount(params, true), config)
    ])
    return { count: count || 0, list: list || [] }
  } catch (error) {
    console.error(error)
    return { count: 0, list: [] }
  }
}

const previewOfAfterCreate = async (params, config) => {
  try {
    const [{ info: list }, { count = 0 }] = await Promise.all([
      http.post('find/biz_set/biz_list', enableCount(params, false), config),
      http.post('find/biz_set/biz_list', enableCount(params, true), config)
    ])
    return { count: count || 0, list: list || [] }
  } catch (error) {
    console.error(error)
    return { count: 0, list: [] }
  }
}

const create = (data, config) => http.post('create/biz_set', data, config)

const update = (data, config) => http.put('updatemany/biz_set', data, config)

const deleteById = (id, config) => http.post('deletemany/biz_set', {
  bk_biz_set_ids: [id]
}, config)

const getAll = config => http.get('findmany/biz_set/simplify', config)

export default {
  find,
  findById,
  findOne,
  create,
  update,
  deleteById,
  getAll,
  previewOfBeforeCreate,
  previewOfAfterCreate,
  getAuthorized,
  getAuthorizedWithCache
}
