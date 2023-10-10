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

import { computed, unref } from 'vue'
import isEqual from 'lodash/isEqual'
import { PROPERTY_TYPES } from '@/dictionary/property-constants'
import { UNIUQE_TYPES } from '@/dictionary/model-constants'

export default function useUnique(beforeUniqueList, uniqueListLocal) {
  const uniqueStatus = computed(() => {
    const status = {}
    uniqueListLocal.value.forEach((data) => {
      const { id: uniqueId, keys: uniqueKeys } = data
      status[uniqueId] = {
        new: false,
        changed: false
      }

      // 在接口数据中找不到，表示为新增
      const matched = beforeUniqueList.value.find(item => item.id === uniqueId)
      if (!matched) {
        status[uniqueId].new = true
      } else { // 能找到，需要检查是否有变化
        status[uniqueId].changed = !isEqual(matched.keys, uniqueKeys)
      }
    })

    // 判断是否删除，从接口数据中查找
    beforeUniqueList.value.forEach((unique) => {
      if (!status[unique.id]) {
        status[unique.id] = {
          removed: false
        }
      }

      // 在本地数据中找不到，表示为删除
      if (!uniqueListLocal.value.find(item => item.id === unique.id)) {
        status[unique.id].removed = true
      }
    })

    return status
  })

  const removedUniqueList = computed(() => beforeUniqueList.value
    .filter(item => uniqueStatus.value[item.id].removed)
    .map(wrapData))

  const getUniqueByField = (field) => {
    const list = unref(uniqueListLocal)
    const fieldUniqueList = list.filter(item => item.keys.includes(field.id))

    let type = UNIUQE_TYPES.SINGLE
    if (fieldUniqueList.length > 1) {
      type = UNIUQE_TYPES.UNION
    } else if (fieldUniqueList[0]) {
      type = fieldUniqueList[0].keys.length > 1 ? UNIUQE_TYPES.UNION : UNIUQE_TYPES.SINGLE
    }

    return {
      list: fieldUniqueList,
      type
    }
  }

  const clearUniqueByField = (field) => {
    const { list } = getUniqueByField(field)
    list.forEach((unique) => {
      const index = uniqueListLocal.value.findIndex(item => item.id === unique.id)
      uniqueListLocal.value.splice(index, 1)
    })
  }

  const getFieldUniqueWithName = (field, fieldList) => {
    const { list: fieldUniqueList, type: fieldUniqueType } = getUniqueByField(field)
    const fieldUniqueWithNameList = fieldUniqueList.map(item => ({
      ...item,
      names: item.keys.map(key => fieldList.find(field => field.id === key)?.bk_property_name)
    }))
    return {
      list: fieldUniqueWithNameList,
      type: fieldUniqueType
    }
  }

  return {
    uniqueStatus,
    wrapData,
    removedUniqueList,
    getUniqueByField,
    clearUniqueByField,
    getFieldUniqueWithName
  }
}

export const wrapData = data => ({
  type: data.keys.length > 1 ? UNIUQE_TYPES.UNION : UNIUQE_TYPES.SINGLE,
  ...data
})

export const isUniqueExist = (unique, uniqueList) => uniqueList.some(item => isEqual(item.keys, unique.keys))

export const getUniqueName = (unique, fieldList, isTemplate, propIdKey) => {
  const ids = isTemplate ? unique.keys : unique.keys.map(key => key.key_id)
  let dataIdKey = isTemplate ? 'bk_property_id' : 'id'
  dataIdKey = propIdKey ?? dataIdKey
  return ids.map((id) => {
    const property = fieldList.find(property => property[dataIdKey] === id)
    return property ? property.bk_property_name : `unknown(${id})`
  }).join(' + ')
}

export const MAX_UNIQUE_COUNT = 5

export const DIFF_TYPES = {
  NEW: 'new',
  UPDATE: 'update',
  CONFLICT: 'conflict',
  UNBOUND: 'unbound',
  UNCHANGED: 'unchanged'
}

export const singleRuleTypes = [
  PROPERTY_TYPES.SINGLECHAR,
  PROPERTY_TYPES.INT,
  PROPERTY_TYPES.FLOAT
]

export const unionRuleTypes = [
  PROPERTY_TYPES.SINGLECHAR,
  PROPERTY_TYPES.INT,
  PROPERTY_TYPES.FLOAT,
  PROPERTY_TYPES.DATE,
  PROPERTY_TYPES.LIST
]
