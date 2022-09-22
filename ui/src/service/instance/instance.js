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

const find = async ({ bk_obj_id: objId, params, config }) => {
  try {
    const [{ info }, { count }] = await Promise.all([
      http.post(`search/instances/object/${objId}`, params, config),
      http.post(`count/instances/object/${objId}`, params)
    ])
    return { count, info: info || [] }
  } catch (error) {
    console.error(error)
    return { count: 0, info: [] }
  }
}

const findOne = async ({ bk_obj_id: objId, bk_inst_id: instId, config }) => {
  try {
    const { info } = await http.post(`search/instances/object/${objId}`, {
      page: { start: 0, limit: 1 },
      fields: [],
      conditions: {
        condition: 'AND',
        rules: [{
          field: 'bk_inst_id',
          operator: 'equal',
          value: instId
        }]
      }
    }, config)
    const [instance] = info || [null]
    return instance
  } catch (error) {
    console.error(error)
    return null
  }
}

export default {
  find,
  findOne
}
