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

import QS from 'qs'
import RouterQuery from '@/router/query'
import Utils from '@/components/filters/utils'

// 快速搜索的默认query值
const defaultFastQuery = () => ({
  field: '',
  filter: '',
  operator: '',
  fuzzy: ''
})

// 通用的默认query
const defaultBaseQuery = () => ({
  page: '',
  _t: Date.now()
})

// 根据条件Map设置搜索query，空值视为删除
export const setSearchQueryByCondition = (conditionMap = {}, properties = []) => {
  const query = QS.parse(RouterQuery.get('filter_adv'))
  const field = RouterQuery.get('field')
  let clearFastQuery = {}

  Object.keys(conditionMap).forEach((id) => {
    const { operator, value } = conditionMap[id]
    const key = `${id}.${operator.replace('$', '')}`

    if (String(value).length) {
      const property = Utils.findProperty(id, properties)
      query[key] = Array.isArray(value) ? value.join(',') : value

      // 与快速搜索重合，清除快速搜索此优先级更高
      if (field === property.bk_property_id) {
        clearFastQuery = defaultFastQuery()
      }
    } else if (Reflect.has(query, key)) {
      Reflect.deleteProperty(query, key)
    }
  })

  Object.keys(query).forEach((key) => {
    const [id] = key.split('.')
    if (!conditionMap[id]) {
      Reflect.deleteProperty(query, key)
    }
  })

  RouterQuery.set({
    filter_adv: QS.stringify(query, { encode: false }),
    s: 'adv',
    ...clearFastQuery,
    ...defaultBaseQuery()
  })
}

// 移除单个查询条件
export const clearOneSearchQuery = (property, operator) => {
  const query = QS.parse(RouterQuery.get('filter_adv'))
  const field = RouterQuery.get('field')

  // 清除快速搜索项
  if (field === property.bk_property_id) {
    RouterQuery.set({
      filter: '',
      s: 'fast',
      ...defaultBaseQuery()
    })
    return
  }

  // 清除高级搜索项
  const key = `${property.id}.${operator.replace('$', '')}`
  if (Reflect.has(query, key)) {
    Reflect.deleteProperty(query, key)
    RouterQuery.set({
      filter_adv: QS.stringify(query, { encode: false }),
      s: 'adv',
      ...defaultBaseQuery()
    })
  }
}

// 清除所有查询条件
export const clearSearchQuery = () => {
  RouterQuery.set({
    filter_adv: '',
    _t: '',
    s: '',
    page: '',
    ...defaultFastQuery()
  })
}

// 重置所有条件项，用于query被清除后重新生成新的条件项
export const resetConditionValue = (condition, selected) => {
  const newConditon = {}
  Object.keys(condition).forEach((id) => {
    const propertyCondititon = condition[id]
    newConditon[id] = { ...propertyCondititon }

    const property = selected.find(property => property.id.toString() === id.toString())
    const value = Utils.getOperatorSideEffect(property, newConditon[id].operator, [])

    newConditon[id].value = value
  })

  return newConditon
}
