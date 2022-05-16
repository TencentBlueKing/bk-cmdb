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

export const find = async (modelId, type, config) => {
  try {
    const key = type === 'source' ? 'bk_obj_id' : 'bk_asst_obj_id'
    const result = await http.post('find/objectassociation', { condition: { [key]: modelId } }, config)
    return result
  } catch (error) {
    console.error(error)
    return []
  }
}

export const findAsSource = modelId => find(modelId, 'source')
export const findAsTarget = modelId => find(modelId, 'target')
export const findAll = async (modelId) => {
  const [source, target] = await Promise.all([findAsSource(modelId), findAsTarget(modelId)])
  const all = [...source, ...target]
  const uniqId = [...new Set(all.map(item => item.id))]
  return uniqId.map(id => all.find(item => item.id === id))
}

export default {
  find,
  findAsSource,
  findAsTarget,
  findAll
}
