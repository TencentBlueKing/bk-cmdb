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

import { computed } from 'vue'
import isEqual from 'lodash/isEqual'
import { isEmptyPropertyValue } from '@/utils/tools'
import { PROPERTY_TYPES } from '@/dictionary/property-constants'

export default function useField(beforeFieldList, fieldListLocal) {
  const fieldStatus = computed(() => {
    const fieldStatus = {}
    fieldListLocal.value.forEach((data) => {
      const { field } = data
      fieldStatus[field.id] = {
        new: false,
        changed: false
      }

      // 在接口数据中找不到，表示为新增
      const matched = beforeFieldList.value.find(item => item.id === field.id)
      if (!matched) {
        fieldStatus[field.id].new = true
      } else { // 能找到，需要检查是否有变化
        const before = normalizeFieldData([matched], false)
        // fieldListLocal中会丢失掉bk_template_id，这里补充上，默认使用与认为与原字段一致的值
        const after = normalizeFieldData([{ ...wrapData(data), bk_template_id: matched.bk_template_id }], false)
        fieldStatus[field.id].changed = !isEqual(before, after)
      }
    })

    // 判断是否删除，从接口数据中查找
    beforeFieldList.value.forEach((field) => {
      if (!fieldStatus[field.id]) {
        fieldStatus[field.id] = {
          removed: false
        }
      }

      // 在本地数据中找不到，表示为删除
      if (!fieldListLocal.value.find(item => item.field.id === field.id)) {
        fieldStatus[field.id].removed = true
      }
    })

    return fieldStatus
  })

  const removedFieldList = computed(() => beforeFieldList.value
    .filter(item => fieldStatus.value[item.id].removed)
    .map(unwrapData))

  return {
    fieldStatus,
    removedFieldList
  }
}

export const withLockKeys = ['editable', 'isrequired', 'placeholder']

export const excludeFieldType = [PROPERTY_TYPES.INNER_TABLE, PROPERTY_TYPES.ENUMQUOTE, PROPERTY_TYPES.ENUM]

export const unwrapData = (data) => {
  const finalData = {
    field: {},
    extra: {
      lock: {
        isrequired: true,
        editable: true,
        placeholder: true
      }
    }
  }
  for (const [key, value] of Object.entries(data)) {
    finalData.field[key] = value

    if (withLockKeys.includes(key)) {
      finalData.field[key] = value.value,
      finalData.extra.lock[key] = value.lock
    }
  }

  return finalData
}

export const wrapData = (data) => {
  const { field, extra } = data
  const settingData = {}
  for (const [key, value] of Object.entries(field)) {
    const ignoreKeys = ['bk_property_group']
    if (ignoreKeys.includes(key)) {
      continue
    }

    settingData[key] = value

    if (withLockKeys.includes(key)) {
      settingData[key] = {
        value,
        lock: extra?.lock?.[key] ?? true
      }
    }
  }

  return settingData
}

export const defaultFieldData = () => ({
  id: '',
  bk_template_id: 0,
  bk_property_id: '',
  bk_property_name: '',
  bk_property_type: '',
  unit: '',
  option: '',
  default: '',
  ismultiple: false,
  placeholder: {
    lock: true,
    value: ''
  },
  isrequired: {
    lock: true,
    value: false
  },
  editable: {
    lock: true,
    value: true
  }
})

export const normalizeFieldData = (fieldData, isCreate = true, fieldStatus) => {
  const fieldList = []
  const defaultData = defaultFieldData()

  fieldData.forEach((item) => {
    const field = {
      ...defaultData,
      ...item
    }

    const valideKeys = Object.keys(defaultData)
    for (const [key] of Object.entries(field)) {
      if (!valideKeys.includes(key)) {
        Reflect.deleteProperty(field, key)
      }
    }

    if (isCreate) {
      Reflect.deleteProperty(field, 'id')

      if (isEmptyPropertyValue(field.default)) {
        Reflect.deleteProperty(field, 'default')
      }
    } else {
      // 编辑流程，需要关注状态
      if (fieldStatus?.value?.[field.id]?.new) {
        Reflect.deleteProperty(field, 'id')
      }
      if (isEmptyPropertyValue(field.default)) {
        field.default = null
      }
    }

    fieldList.push(field)
  })

  return fieldList
}

export const normalizeUniqueData = (uniqueData, fieldData, isCreate = true, uniqueStatus) => {
  const uniqueList = uniqueData.map((item) => {
    const unique = {
      id: item.id,
      keys: item.keys.map(key => fieldData.find(field => field.id === key)?.bk_property_id)
    }

    if (isCreate) {
      Reflect.deleteProperty(unique, 'id')
    } else {
      if (uniqueStatus?.value?.[unique.id]?.new) {
        Reflect.deleteProperty(unique, 'id')
      }
    }

    return unique
  })

  return uniqueList
}

export const isFieldSame = (field1, field2, id) => (field1.bk_property_id === field2.bk_property_id
  || field1.bk_property_name === field2.bk_property_name)
  && (!id || field1.id !== id)

export const isFieldExist = (field, fieldList, id) => fieldList.some(item => isFieldSame(item.field, field, id))

export const MAX_FIELD_COUNT = 20

export const DETAILS_FIELDLIST_REQUEST_ID = Symbol()

export const DIFF_TYPES = {
  NEW: 'new',
  UPDATE: 'update',
  CONFLICT: 'conflict',
  UNBOUND: 'unbound',
  UNCHANGED: 'unchanged'
}
