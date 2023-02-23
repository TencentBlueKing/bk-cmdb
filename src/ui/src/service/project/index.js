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
import { enableCount } from '../utils.js'

const MODEL_ID_KEY = BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.PROJECT].ID

const findOne = async ({ id }, config) => {
  try {
    const { info } = await http.post('/findmany/project', {
      filter: {
        condition: 'AND',
        rules: [
          {
            field: 'id',
            operator: 'equal',
            value: id
          }
        ]
      },
      page: {
        start: 0,
        limit: 1,
        sort: 'id',
        enable_count: false
      }
    }, config)
    const [instance] = info || [null]
    return instance
  } catch (error) {
    console.error(error)
    return null
  }
}

const find = async ({ params, config }) => {
  try {
    return await http.post('/findmany/project', params, config)
  } catch (error) {
    console.error(error)
    return null
  }
}

const create =  params => http.post('/createmany/project', params)
const getMany = async (params, config) => {
  try {
    const [{ info: list = [] }, { count = 0 }] = await Promise.all([
      http.post('findmany/project', enableCount(params, false), config),
      http.post('findmany/project', enableCount(params, true), config)
    ])
    return { count: count || 0, list: list || [] }
  } catch (error) {
    console.error(error)
    return Promise.reject(error)
  }
}

const findByIds = async (ids, config = {}) => {
  try {
    const { count = 0, info: list = [] } = await http.post('findmany/project', enableCount({
      filter: {
        condition: 'AND',
        rules: [{
          field: MODEL_ID_KEY,
          operator: 'in',
          value: ids
        }]
      },
      page: { start: 0, limit: ids.length }
    }, false), config)

    return { count, list }
  } catch (error) {
    console.error(error)
    return null
  }
}


const update =  params => http.put('/updatemany/project', params)

export default {
  findOne,
  find,
  create,
  update,
  findByIds,
  getMany
}
