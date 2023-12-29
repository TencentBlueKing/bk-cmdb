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

import { t } from '@/i18n'

export const PROPERTY_TYPE_EXCLAMATION_TIPS = ['innertable', 'bool']

export const PROPERTY_TYPES = Object.freeze({
  SINGLECHAR: 'singlechar',
  INT: 'int',
  FLOAT: 'float',
  ENUM: 'enum',
  DATE: 'date',
  TIME: 'time',
  LONGCHAR: 'longchar',
  OBJUSER: 'objuser',
  TIMEZONE: 'timezone',
  BOOL: 'bool',
  LIST: 'list',
  ORGANIZATION: 'organization',
  ENUMMULTI: 'enummulti',
  ENUMQUOTE: 'enumquote',
  MAP: 'map',
  OBJECT: 'object',
  ARRAY: 'array',
  TABLE: 'table',
  SERVICE_TEMPLATE: 'service-template',
  TOPOLOGY: 'topology',
  FOREIGNKEY: 'foreignkey',
  INNER_TABLE: 'innertable'
})

export const PROPERTY_TYPE_NAMES = Object.freeze({
  [PROPERTY_TYPES.SINGLECHAR]: t('短字符'),
  [PROPERTY_TYPES.INT]: t('数字'),
  [PROPERTY_TYPES.FLOAT]: t('浮点'),
  [PROPERTY_TYPES.ENUM]: t('枚举'),
  [PROPERTY_TYPES.DATE]: t('日期'),
  [PROPERTY_TYPES.TIME]: t('时间'),
  [PROPERTY_TYPES.LONGCHAR]: t('长字符'),
  [PROPERTY_TYPES.OBJUSER]: t('用户'),
  [PROPERTY_TYPES.TIMEZONE]: t('时区'),
  [PROPERTY_TYPES.BOOL]: t('bool'),
  [PROPERTY_TYPES.LIST]: t('列表'),
  [PROPERTY_TYPES.ORGANIZATION]: t('组织'),
  [PROPERTY_TYPES.ENUMMULTI]: t('枚举(多选)'),
  [PROPERTY_TYPES.INNER_TABLE]: t('表格'),
  [PROPERTY_TYPES.ENUMQUOTE]: t('枚举(引用)'),
  [PROPERTY_TYPES.FOREIGNKEY]: t('系统内置类型')
})

export const PROPERTY_TYPE_LIST = [
  {
    id: PROPERTY_TYPES.SINGLECHAR,
    name: PROPERTY_TYPE_NAMES[PROPERTY_TYPES.SINGLECHAR]
  },
  {
    id: PROPERTY_TYPES.INT,
    name: PROPERTY_TYPE_NAMES[PROPERTY_TYPES.INT]
  },
  {
    id: PROPERTY_TYPES.FLOAT,
    name: PROPERTY_TYPE_NAMES[PROPERTY_TYPES.FLOAT]
  },
  {
    id: PROPERTY_TYPES.ENUM,
    name: PROPERTY_TYPE_NAMES[PROPERTY_TYPES.ENUM]
  },
  {
    id: PROPERTY_TYPES.ENUMMULTI,
    name: PROPERTY_TYPE_NAMES[PROPERTY_TYPES.ENUMMULTI]
  },
  {
    id: PROPERTY_TYPES.ENUMQUOTE,
    name: PROPERTY_TYPE_NAMES[PROPERTY_TYPES.ENUMQUOTE]
  },
  {
    id: PROPERTY_TYPES.INNER_TABLE,
    name: PROPERTY_TYPE_NAMES[PROPERTY_TYPES.INNER_TABLE]
  },
  {
    id: PROPERTY_TYPES.DATE,
    name: PROPERTY_TYPE_NAMES[PROPERTY_TYPES.DATE]
  },
  {
    id: PROPERTY_TYPES.TIME,
    name: PROPERTY_TYPE_NAMES[PROPERTY_TYPES.TIME]
  },
  {
    id: PROPERTY_TYPES.LONGCHAR,
    name: PROPERTY_TYPE_NAMES[PROPERTY_TYPES.LONGCHAR]
  },
  {
    id: PROPERTY_TYPES.OBJUSER,
    name: PROPERTY_TYPE_NAMES[PROPERTY_TYPES.OBJUSER]
  },
  {
    id: PROPERTY_TYPES.TIMEZONE,
    name: PROPERTY_TYPE_NAMES[PROPERTY_TYPES.TIMEZONE]
  },
  {
    id: PROPERTY_TYPES.BOOL,
    name: PROPERTY_TYPE_NAMES[PROPERTY_TYPES.BOOL]
  },
  {
    id: PROPERTY_TYPES.LIST,
    name: PROPERTY_TYPE_NAMES[PROPERTY_TYPES.LIST]
  },
  {
    id: PROPERTY_TYPES.ORGANIZATION,
    name: PROPERTY_TYPE_NAMES[PROPERTY_TYPES.ORGANIZATION]
  },
  {
    id: PROPERTY_TYPES.FOREIGNKEY,
    name: PROPERTY_TYPE_NAMES[PROPERTY_TYPES.FOREIGNKEY]
  }
]

export const EDITABLE_TYPES = [
  PROPERTY_TYPES.SINGLECHAR,
  PROPERTY_TYPES.INT,
  PROPERTY_TYPES.FLOAT,
  PROPERTY_TYPES.ENUM,
  PROPERTY_TYPES.DATE,
  PROPERTY_TYPES.TIME,
  PROPERTY_TYPES.LONGCHAR,
  PROPERTY_TYPES.OBJUSER,
  PROPERTY_TYPES.TIMEZONE,
  PROPERTY_TYPES.BOOL,
  PROPERTY_TYPES.LIST,
  PROPERTY_TYPES.ORGANIZATION,
  PROPERTY_TYPES.ENUMMULTI,
  PROPERTY_TYPES.ENUMQUOTE,
  PROPERTY_TYPES.INNER_TABLE
]

export const REQUIRED_TYPES = [
  PROPERTY_TYPES.SINGLECHAR,
  PROPERTY_TYPES.INT,
  PROPERTY_TYPES.FLOAT,
  PROPERTY_TYPES.DATE,
  PROPERTY_TYPES.TIME,
  PROPERTY_TYPES.LONGCHAR,
  PROPERTY_TYPES.OBJUSER,
  PROPERTY_TYPES.TIMEZONE,
  PROPERTY_TYPES.LIST,
  PROPERTY_TYPES.ORGANIZATION,
  PROPERTY_TYPES.INNER_TABLE
]
