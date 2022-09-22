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
import { BUILTIN_MODELS } from '@/dictionary/model-constants.js'

export const find = async (options, config = {}) => {
  try {
    const params = {}
    options.bk_biz_id && (params.bk_biz_id = options.bk_biz_id)
    const groups = await http.post(`find/objectattgroup/object/${options.bk_obj_id}`, params, config)
    const bizGroups = groups.filter(group => !!group.bk_biz_id)
      .sort((previous, next) => previous.bk_group_index - next.bk_group_index)
    const globalGroups = groups.filter(group => !group.bk_biz_id)
      .sort((previous, next) => previous.bk_group_index - next.bk_group_index)
    return [...globalGroups, ...bizGroups]
  } catch (error) {
    console.error(error)
    return []
  }
}

export const findBizSet = config => find({ bk_obj_id: BUILTIN_MODELS.BUSINESS_SET }, config)

export default {
  find,
  findBizSet
}
