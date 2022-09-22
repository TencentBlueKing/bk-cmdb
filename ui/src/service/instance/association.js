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
import i18n from '@/i18n'
import { BUILTIN_MODELS, BUILTIN_MODEL_PROPERTY_KEYS } from '@/dictionary/model-constants.js'

const getIdKey = modelId => ({
  [BUILTIN_MODELS.HOST]: [BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.HOST].ID],
  [BUILTIN_MODELS.BUSINESS]: [BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS].ID],
  [BUILTIN_MODELS.BUSINESS_SET]: [BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS_SET].ID]
}[modelId] || 'bk_inst_id')
const getNameKey = modelId => ({
  [BUILTIN_MODELS.HOST]: 'bk_host_innerip',
  [BUILTIN_MODELS.BUSINESS]: [BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS].NAME],
  [BUILTIN_MODELS.BUSINESS_SET]: [BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS_SET].NAME]
}[modelId] || 'bk_inst_name')
const findInstance = (instances, objId, instId) => {
  const idKey = getIdKey(objId)
  return (instances || []).find(instance => instance[idKey] === instId)
}

const findTopology = async ({
  bk_obj_id: currentModelId,
  bk_inst_id: currentInstId,
  bk_inst_name: currentInstName,
  offset = 0,
  limit = 200
}) => {
  try {
    // eslint-disable-next-line max-len
    const url = `findmany/inst/association/object/${currentModelId}/inst_id/${currentInstId}/offset/${offset}/limit/${limit}/web`
    const result = await http.post(url, {})
    // 忽略实例作为源还是目标，抹平不同模型间的key差异
    const all =  [...(result.data.association.dst || []), ...(result.data.association.src || [])]
    const data = all.map((association) => {
      const {
        bk_obj_id: objId,
        bk_inst_id: instId,
        bk_asst_obj_id: asstObjId,
        bk_asst_inst_id: asstInstId
      } = association
      const isSource = objId === currentModelId && instId === currentInstId
      const instance = isSource
        ? findInstance(result.data.instance[asstObjId], asstObjId, asstInstId)
        : findInstance(result.data.instance[objId], objId, instId)
      const nameKey = isSource ? getNameKey(asstObjId) : getNameKey(objId)
      return {
        id: association.id,
        bk_obj_id: isSource ? asstObjId : objId,
        bk_inst_id: isSource ? asstInstId : instId,
        bk_inst_name: instance ? instance[nameKey] : `${i18n.t('已删除的实例')}(ID: ${isSource ? asstInstId : instId})`,
        bk_asst_id: association.bk_asst_id,
        bk_obj_asst_id: association.bk_obj_asst_id,
        deleted: !instance,
        target: isSource
      }
    })
    const rootIdKey = getIdKey(currentModelId)
    const rootNameKey = getNameKey(currentModelId)
    const rootInstance = (result.data.instance[currentModelId] || []).find(root => root[rootIdKey] === currentInstId)
    return {
      count: result.association_count,
      root: {
        bk_obj_id: currentModelId,
        bk_inst_id: currentInstId,
        bk_inst_name: rootInstance ? rootInstance[rootNameKey] : currentInstName
      },
      data
    }
  } catch (error) {
    console.error(error)
    return { count: 0, data: [], root: { bk_obj_id: currentModelId, bk_inst_id: currentInstId } }
  }
}

const find = async (params, config) => {
  try {
    const [{ info }, [{ count }]] = await Promise.all([
      http.post(`search/instance_associations/object/${params.bk_obj_id}`, params, config),
      http.post(`count/instance_associations/object/${params.bk_obj_id}`, params)
    ])
    return { count, info: info || [] }
  } catch (error) {
    console.error(error)
    return { count: 0, info: [] }
  }
}

const MAX_LIMIT = 500
const findAll = async (params) => {
  try {
    const { count } = await http.post(`count/instance_associations/object/${params.bk_obj_id}`, params)
    if (count === 0) {
      return []
    }
    const requestProxy = Array(Math.ceil(count / MAX_LIMIT)).fill(null)
    const all = await Promise.all(requestProxy.map((_, index) => {
      const page = { start: index * MAX_LIMIT, limit: MAX_LIMIT }
      return http.post(`search/instance_associations/object/${params.bk_obj_id}`, {
        ...params,
        page
      })
    }))
    return all.reduce((acc, { info }) => {
      acc.push(...info)
      return acc
    }, [])
  } catch (error) {
    return []
  }
}

export default {
  find,
  findAll,
  findTopology
}
