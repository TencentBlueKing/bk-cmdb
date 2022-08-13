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

import { CONTAINER_OBJECTS, WORKLOAD_TYPES } from '@/dictionary/container'

// 根据workload具体类型判断是否为workload
export const isWorkload = type => Object.values(WORKLOAD_TYPES).includes(type)

// 获取容器节点大类型
export const getContainerNodeType = type => (isWorkload(type) ? CONTAINER_OBJECTS.WORKLOAD : type)

// 与传统模型字段类型的映射，原则上在交互形态完全一致的情况下才可以转换
export const typeMapping = {
  string: 'singlechar',
  numeric: 'float',
  mapString: 'map',
  array: 'array',
  object: 'object'
}

export const getNormalType = type => typeMapping[type] || type
