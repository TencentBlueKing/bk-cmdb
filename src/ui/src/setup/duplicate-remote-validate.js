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

import { onePageParams } from '@/service/utils'
import fieldTemplateService from '@/service/field-template'
import queryBuilderOperator, { QUERY_OPERATOR } from '@/utils/query-builder-operator'

export const DUP_CHECK_IDS = {
  FIELD_TEMPLATE_NAME: 'field_template_name'
}

const requestConfigBase = key => ({
  requestId: `dupcheck_${key}`,
  fromCache: true
})

const dupChecks = {
  [DUP_CHECK_IDS.FIELD_TEMPLATE_NAME]: async (value, oldValue) => {
    const params = {
      template_filter: {
        condition: 'AND',
        rules: [{
          field: 'name',
          operator: queryBuilderOperator(QUERY_OPERATOR.EQ),
          value
        }]
      },
      page: onePageParams()
    }
    const { list: [template = null] } = await fieldTemplateService.find(params, requestConfigBase(`field_template_name_${value}_${oldValue}`))
    if (oldValue) {
      return template && template?.name !== oldValue
    }
    return template?.id > 0
  }
}

export default {
  validate: async (value, [id, oldValue]) => {
    try {
      const isDup = await dupChecks[id](value, oldValue)
      return { valid: !isDup }
    } catch (error) {
      return { valid: false }
    }
  }
}
