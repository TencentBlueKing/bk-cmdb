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

import has from 'has'
import { computed, reactive, set } from '@vue/composition-api'

export default ({ properties, propertyGroups, exclude }) => {
  const groupCollapseState = reactive({})

  // 已排序的字段分组
  const sortedGroups = computed(() => {
    const publicGroups = []
    const bizCustomGroups = []
    propertyGroups.value?.forEach((group) => {
      if (has(group, 'bk_biz_id') && group.bk_biz_id > 0) {
        bizCustomGroups.push(group)
      } else {
        publicGroups.push(group)
      }
    })
    const sortKey = 'bk_group_index'
    publicGroups.sort((groupA, groupB) => groupA[sortKey] - groupB[sortKey])
    bizCustomGroups.sort((groupA, groupB) => groupA[sortKey] - groupB[sortKey])
    const allGroups = [
      ...publicGroups,
      ...bizCustomGroups
    ]
    allGroups.forEach((group, index) => {
      group.bk_group_index = index
      set(groupCollapseState, group.bk_group_id, group.is_collapse)
    })
    return allGroups
  })

  // 已排序的字段
  const sortedProperties = computed(() => {
    const sortKey = 'bk_property_index'
    const propertyList = properties.value
      ?.filter(property => !exclude.value.includes(property.bk_property_id) && !property.bk_isapi)
    return propertyList.sort((propertyA, propertyB) => propertyA[sortKey] - propertyB[sortKey])
  })

  // 按分组聚合的属性列表
  const groupedProperties = computed(() => sortedGroups.value
    .map(group => sortedProperties.value?.filter((property) => {
      // 兼容旧数据，把none这个分组的属性塞到默认分组去
      const isNoneGroup = property.bk_property_group === 'none'
      if (isNoneGroup) {
        return group.bk_group_id === 'default'
      }
      return property.bk_property_group === group.bk_group_id
    })))

  // 去除空项后的分组
  const displayGroups = computed(() => sortedGroups.value
    .filter((item, index) => groupedProperties.value[index]?.length > 0))

  // 去除空项后的分组属性列表
  const displayProperties = computed(() => groupedProperties.value.filter(properties => properties.length > 0))

  return {
    groupCollapseState,
    sortedGroups,
    sortedProperties,
    groupedProperties,
    displayGroups,
    displayProperties
  }
}
