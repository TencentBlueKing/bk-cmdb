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
import { definePropertyGroup } from '@/components/filters/utils'

export const getMany = async ({ objId }) => {
  try {
    const propertyDefaultGroup = definePropertyGroup({
      id: 1,
      bk_obj_id: objId,
      bk_group_id: 'default',
      bk_group_index: -1,
      bk_group_name: i18n.locale === 'en' ? 'Default' : '基础信息',
    })

    return [propertyDefaultGroup]
  } catch (error) {
    console.error(error)
    return Promise.reject(error)
  }
}

export default {
  getMany
}
