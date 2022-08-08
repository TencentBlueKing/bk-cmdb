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
import { normalizationProperty } from '@/service/container/transition.js'

export const find = ({ objId, params }, config) => http.post(`find/container/${objId}/attributes`, params, config)

export const getAll = async ({ objId, params }, config) => {
  try {
    const list = await find({ objId, params }, config)
    return normalizationProperty(list, objId)
  } catch (error) {
    console.error(error)
    return Promise.reject(error)
  }
}

export default {
  find,
  getAll
}
