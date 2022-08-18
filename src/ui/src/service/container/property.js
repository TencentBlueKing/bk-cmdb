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
import { normalizationProperty } from '@/service/container/transition.js'
import { CONTAINER_OBJECTS, CONTAINER_OBJECT_INST_KEYS } from '@/dictionary/container.js'


function createIdProperty(objId) {
  const keyMap = {
    [CONTAINER_OBJECTS.CLUSTER]: CONTAINER_OBJECT_INST_KEYS[CONTAINER_OBJECTS.CLUSTER].ID,
    [CONTAINER_OBJECTS.NAMESPACE]: CONTAINER_OBJECT_INST_KEYS[CONTAINER_OBJECTS.NAMESPACE].ID,
    [CONTAINER_OBJECTS.WORKLOAD]: CONTAINER_OBJECT_INST_KEYS[CONTAINER_OBJECTS.WORKLOAD].ID,
    [CONTAINER_OBJECTS.FOLDER]: CONTAINER_OBJECT_INST_KEYS[CONTAINER_OBJECTS.FOLDER].ID,
    [CONTAINER_OBJECTS.POD]: CONTAINER_OBJECT_INST_KEYS[CONTAINER_OBJECTS.POD].ID
  }
  return {
    id: `${objId}_${keyMap[objId]}`,
    bk_obj_id: objId,
    bk_property_id: keyMap[objId] || 'id',
    bk_property_name: 'ID',
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

export const find = ({ objId, params }, config) => http.post(`find/container/${objId}/attributes`, params, config)

export const getMany = async ({ objId, params }, config, injectId = true) => {
  try {
    const list = await find({ objId, params }, config)

    const properties = normalizationProperty(list, objId)

    if (!injectId) {
      return properties
    }

    if (list.some(property => property.is_inject)) {
      return properties
    }

    properties.unshift(createIdProperty(objId))
    return properties
  } catch (error) {
    console.error(error)
    return Promise.reject(error)
  }
}

export default {
  find,
  getMany
}
