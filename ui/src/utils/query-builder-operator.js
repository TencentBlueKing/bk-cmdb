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

export const QUERY_OPERATOR = Object.freeze({
  EQ: '$eq',
  NE: '$ne',
  IN: '$in',
  NIN: '$nin',
  LT: '$lt',
  GT: '$gt',
  LTE: '$lte',
  GTE: '$gte',
  // 前端构造的操作符，真实数据中会拆分数据为gte, lte向后台传递
  RANGE: '$range',
  NRANGE: '$nrange',
  LIKE: '$regex'
})

export const QUERY_OPERATOR_SYMBOL = {
  [QUERY_OPERATOR.EQ]: '=',
  [QUERY_OPERATOR.NE]: '≠',
  [QUERY_OPERATOR.IN]: 'in',
  [QUERY_OPERATOR.NIN]: 'not in',
  [QUERY_OPERATOR.GT]: '>',
  [QUERY_OPERATOR.LT]: '<',
  [QUERY_OPERATOR.GTE]: '≥',
  [QUERY_OPERATOR.LTE]: '≤',
  [QUERY_OPERATOR.LIKE]: 'like',
  [QUERY_OPERATOR.RANGE]: '≤ ≥'
}

export const QUERY_OPERATOR_DESC = {
  [QUERY_OPERATOR.EQ]: t('等于'),
  [QUERY_OPERATOR.NE]: t('不等于'),
  [QUERY_OPERATOR.LT]: t('小于'),
  [QUERY_OPERATOR.GT]: t('大于'),
  [QUERY_OPERATOR.IN]: t('精确'),
  [QUERY_OPERATOR.NIN]: t('精确'),
  [QUERY_OPERATOR.RANGE]: t('数值范围'),
  [QUERY_OPERATOR.LTE]: t('小于等于'),
  [QUERY_OPERATOR.GTE]: t('大于等于'),
  [QUERY_OPERATOR.LIKE]: t('模糊')
}

const mapping = {
  [QUERY_OPERATOR.EQ]: 'equal',
  [QUERY_OPERATOR.NE]: 'not_equal',
  [QUERY_OPERATOR.IN]: 'in',
  [QUERY_OPERATOR.NIN]: 'not_in',
  [QUERY_OPERATOR.LT]: 'less',
  [QUERY_OPERATOR.LTE]: 'less_or_equal',
  [QUERY_OPERATOR.GT]: 'greater',
  [QUERY_OPERATOR.GTE]: 'greater_or_equal',
  [QUERY_OPERATOR.RANGE]: 'between',
  [QUERY_OPERATOR.NRANGE]: 'not_between',
  [QUERY_OPERATOR.LIKE]: 'contains'
}

export default operator => mapping[operator]
