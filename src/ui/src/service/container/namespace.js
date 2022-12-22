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
import { CONTAINER_OBJECTS, CONTAINER_OBJECT_INST_KEYS } from '@/dictionary/container.js'
import { enableCount, onePageParams } from '../utils.js'

const ID_KEY = CONTAINER_OBJECT_INST_KEYS[CONTAINER_OBJECTS.NAMESPACE].ID

const find = async (params, config) => {
  const api = `findmany/kube/namespace/bk_biz_id/${params.bk_biz_id}`
  try {
    const [{ info: list = [] }, { count = 0 }] = await Promise.all([
      http.post(api, enableCount(params, false), config),
      http.post(api, enableCount(params, true), config)
    ])
    return { count, list }
  } catch (error) {
    console.error(error)
    return Promise.reject(error)
  }
}

const findById = async (id, bizId, config = {}) => {
  try {
    const { info: [instance = null] } = await http.post(`findmany/kube/namespace/bk_biz_id/${bizId}`, enableCount({
      filter: {
        condition: 'AND',
        rules: [{
          field: ID_KEY,
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

const getOne = async (params, config = {}) => findById(params[ID_KEY], params.bizId, config)

export default {
  find,
  findById,
  getOne
}
