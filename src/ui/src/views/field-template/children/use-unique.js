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

import { unref } from 'vue'
import { UNIUQE_TYPES } from '@/dictionary/model-constants'

export default function useUnique(uniqueList) {
  const getUniqueByField = (field) => {
    const list = unref(uniqueList)
    const fieldUniqueList = list.filter(item => item.keys.includes(field.id))

    let type = UNIUQE_TYPES.SINGLE
    if (fieldUniqueList.length > 1) {
      type = UNIUQE_TYPES.UNION
    } else if (fieldUniqueList[0]) {
      type = fieldUniqueList[0].keys.length > 1 ? UNIUQE_TYPES.UNION : UNIUQE_TYPES.SINGLE
    }

    return {
      list: fieldUniqueList,
      type
    }
  }

  const clearUniqueByField = (field) => {
    const { list } = getUniqueByField(field)
    list.forEach((unique) => {
      const index = uniqueList.value.findIndex(item => item.id === unique.id)
      uniqueList.value.splice(index, 1)
    })
  }

  return {
    getUniqueByField,
    clearUniqueByField
  }
}
