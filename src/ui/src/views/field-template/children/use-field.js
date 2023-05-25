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
import { PROPERTY_TYPES } from '@/dictionary/property-constants'

export default function useField(fieldList, fieldListLocal) {
  const fieldStatus = computed(() => {
    const fieldStatus = {}
    fieldListLocal.value.forEach(({ field }) => {
      fieldStatus[field.id] = {
        new: false,
        changed: false
      }

      // 在接口数据中找不到，表示为新增
      if (!fieldList.find(item => item.id === field.id)) {
        fieldStatus[field.id].new = true
      }
    })

    fieldList.forEach((field) => {
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

  return {
    fieldStatus
  }
}

export const withLockKeys = ['editable', 'isrequired', 'placeholder']

export const excludeFieldType = [PROPERTY_TYPES.INNER_TABLE, PROPERTY_TYPES.ENUMQUOTE]

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

export const normalizeFieldData = (fieldData, isCreate = true) => {
  const fieldList = []
  const defaultData = {
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
  }

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
    }

    fieldList.push(field)
  })

  return fieldList
}

export const normalizeUniqueData = (uniqueData, fieldData, isCreate = true) => {
  const uniqueList = uniqueData.map((item) => {
    const unique = {
      id: item.id,
      keys: item.keys.map(key => fieldData.find(field => field.id === key)?.bk_property_id)
    }

    // TODO: 补充编辑流程中新创建的处理
    if (isCreate) {
      Reflect.deleteProperty(unique, 'id')
    }

    return unique
  })

  return uniqueList
}
