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
import i18n from '@/i18n/index.js'
import { normalizationProperty } from '@/service/container/transition.js'
import { CONTAINER_OBJECTS, CONTAINER_OBJECT_INST_KEYS } from '@/dictionary/container.js'
import { rollReqUseTotalCount } from '@/service/utils'
import { getPropertyName } from './common.js'
import { defineProperty as defineModelProperty } from '@/components/filters/utils.js'

function createIdProperty(objId, isPrependName) {
  const keyMap = {
    [CONTAINER_OBJECTS.CLUSTER]: CONTAINER_OBJECT_INST_KEYS[CONTAINER_OBJECTS.CLUSTER].ID,
    [CONTAINER_OBJECTS.NAMESPACE]: CONTAINER_OBJECT_INST_KEYS[CONTAINER_OBJECTS.NAMESPACE].ID,
    [CONTAINER_OBJECTS.WORKLOAD]: CONTAINER_OBJECT_INST_KEYS[CONTAINER_OBJECTS.WORKLOAD].ID,
    [CONTAINER_OBJECTS.FOLDER]: CONTAINER_OBJECT_INST_KEYS[CONTAINER_OBJECTS.FOLDER].ID,
    [CONTAINER_OBJECTS.POD]: CONTAINER_OBJECT_INST_KEYS[CONTAINER_OBJECTS.POD].ID,
    [CONTAINER_OBJECTS.NODE]: CONTAINER_OBJECT_INST_KEYS[CONTAINER_OBJECTS.NODE].ID,
    [CONTAINER_OBJECTS.CONTAINER]: CONTAINER_OBJECT_INST_KEYS[CONTAINER_OBJECTS.CONTAINER].ID
  }
  const propertyIdKey = keyMap[objId] || 'id'
  return {
    id: `${objId}_${keyMap[objId]}`,
    bk_obj_id: objId,
    bk_property_id: propertyIdKey,
    bk_property_name: isPrependName ? (getPropertyName(propertyIdKey, objId, i18n.locale) || 'ID') : 'ID',
    bk_property_index: -1,
    bk_property_type: 'int',
    isonly: true,
    ispre: true,
    bk_isapi: false,
    bk_issystem: true,
    isreadonly: true,
    editable: false,
    bk_property_group: 'default',
    is_inject: true
  }
}

export const find = ({ objId, params }, config) => http.get(`find/kube/${objId}/attributes`, params, config)

export const getMany = async ({ objId, params }, config, injectId = true, isPrependName = false) => {
  try {
    const list = await find({ objId, params }, config)

    const properties = normalizationProperty(list, objId)

    if (!injectId) {
      return properties
    }

    if (list.some(property => property.is_inject)) {
      return properties
    }

    properties.unshift(createIdProperty(objId, isPrependName))
    return properties
  } catch (error) {
    console.error(error)
    return Promise.reject(error)
  }
}

export const getMapValue = async (params, total, config) => {
  const result = await rollReqUseTotalCount(
    `find/kube/field/map_str_val/bk_biz_id/${params.bk_biz_id}`,
    params,
    { limit: 200, total },
    config,
    'post',
    data => data.info
  )

  const mergeResult = {}
  const { fields = [] } = params
  result.forEach((row) => {
    fields.forEach((field) => {
      const data = row[field] || []
      if (mergeResult[field]) {
        // 去重
        const newItems = data.filter(item => !mergeResult[field].some(x => x.key === item.key && x.val === item.val))
        mergeResult[field].push(...newItems)
      } else {
        mergeResult[field] = data.slice()
      }
    })
  })

  return mergeResult
}

export const getPodTopoNodeProps = () => {
  const propIds = ['cluster_uid', 'namespace', 'ref']
  const objId = CONTAINER_OBJECTS.POD
  return propIds.map(id => defineModelProperty({
    id: `${objId}_${id}`,
    bk_obj_id: objId,
    bk_property_id: id,
    bk_property_name: getPropertyName(id, objId, i18n.locale),
    bk_property_index: 0,
    bk_property_type: 'singlechar'
  }))
}

export default {
  find,
  getMany,
  getMapValue,
  getPodTopoNodeProps
}
