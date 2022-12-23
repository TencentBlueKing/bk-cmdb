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

import i18n from '@/i18n/index.js'
import { CONTAINER_OBJECT_NAMES } from '@/dictionary/container'
import { isWorkload, isFolder, getContainerNodeType, getPropertyType, getPropertyName, getContainerObjectNames } from './common.js'
import Utils from '@/components/filters/utils.js'

export const normalizationTopo = (topoList, refId) => {
  const topo = topoList.map((item) => {
    // 小分类，具体类型
    const { kind } = item

    // 大分类
    const type = getContainerNodeType(kind)

    return {
      bk_inst_id: item.id,
      bk_inst_name: item.name,
      bk_obj_id: kind,
      bk_obj_name: CONTAINER_OBJECT_NAMES[type].FULL,
      default: 0,
      child: [],
      icon_text: getContainerObjectNames(kind).SHORT,
      is_container: true,
      is_workload: isWorkload(kind),
      is_folder: isFolder(type),
      // 上一级的id
      ref_id: refId
    }
  })

  return topo
}

export const normalizationProperty = (propertyList, objId) => {
  const properties = propertyList.map(item => Utils.defineProperty({
    id: `${objId}_${item.field}`,
    bk_obj_id: objId,
    bk_property_id: item.field,
    bk_property_name: getPropertyName(item.field, objId, i18n.locale),
    bk_property_index: Infinity,
    bk_property_type: getPropertyType(item.type),
    required: item.required,
    editable: item.editable,
    option: item.option,
    bk_isapi: false
  }))

  return properties
}

export default {
  normalizationTopo,
  normalizationProperty
}
