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

import {
  MENU_RESOURCE_HOST,
  MENU_RESOURCE_BUSINESS,
  MENU_RESOURCE_BUSINESS_SET,
  MENU_RESOURCE_HOST_COLLECTION,
  MENU_RESOURCE_BUSINESS_COLLECTION,
  MENU_RESOURCE_BUSINESS_SET_COLLECTION
} from '@/dictionary/menu-symbol'

// 常用的内置模型ID
export const BUILTIN_MODELS = Object.freeze({
  BUSINESS_SET: 'bk_biz_set_obj',
  BUSINESS: 'biz',
  SET: 'set',
  MODULE: 'module',
  HOST: 'host'
})

// 内置模型ID和名称属性的Key
export const BUILTIN_MODEL_PROPERTY_KEYS = Object.freeze({
  [BUILTIN_MODELS.BUSINESS_SET]: {
    ID: 'bk_biz_set_id',
    NAME: 'bk_biz_set_name'
  },
  [BUILTIN_MODELS.BUSINESS]: {
    ID: 'bk_biz_id',
    NAME: 'bk_biz_name'
  },
  [BUILTIN_MODELS.HOST]: {
    ID: 'bk_host_id',
    NAME: 'bk_host_name'
  },
  [BUILTIN_MODELS.MODULE]: {
    ID: 'bk_module_id',
    NAME: 'bk_module_name'
  },
  [BUILTIN_MODELS.SET]: {
    ID: 'bk_set_id',
    NAME: 'bk_set_name'
  }
})

// 内置模型路由路径参数的Key
export const BUILTIN_MODEL_ROUTEPARAMS_KEYS = Object.freeze({
  [BUILTIN_MODELS.BUSINESS]: 'bizId',
  [BUILTIN_MODELS.BUSINESS_SET]: 'bizSetId',
})

// 内置模型资源目录收藏中用到的Key
export const BUILTIN_MODEL_COLLECTION_KEYS = Object.freeze({
  [BUILTIN_MODELS.HOST]: MENU_RESOURCE_HOST_COLLECTION,
  [BUILTIN_MODELS.BUSINESS]: MENU_RESOURCE_BUSINESS_COLLECTION,
  [BUILTIN_MODELS.BUSINESS_SET]: MENU_RESOURCE_BUSINESS_SET_COLLECTION
})

// 内置模型收藏后的菜单路由名称
export const BUILTIN_MODEL_RESOURCE_MENUS = Object.freeze({
  [BUILTIN_MODELS.HOST]: MENU_RESOURCE_HOST,
  [BUILTIN_MODELS.BUSINESS]: MENU_RESOURCE_BUSINESS,
  [BUILTIN_MODELS.BUSINESS_SET]: MENU_RESOURCE_BUSINESS_SET
})

// 内置模型资源类型用于详情区分实例类型与变更记录查询
export const BUILTIN_MODEL_RESOURCE_TYPES = Object.freeze({
  [BUILTIN_MODELS.HOST]: 'host',
  [BUILTIN_MODELS.BUSINESS]: 'business',
  [BUILTIN_MODELS.BUSINESS_SET]: 'biz_set'
})

/**
 * @constant {String} 未分类模型分组 ID
 */
export const UNCATEGORIZED_GROUP_ID = 'bk_uncategorized'
