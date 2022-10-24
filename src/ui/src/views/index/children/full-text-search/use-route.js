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

import { reactive, toRefs, watch } from 'vue'

export const pickQuery = (query = {}, include = [], exclude = []) => {
  let queryKeys = Object.keys(query)
  if (include.length) {
    queryKeys = queryKeys.filter(item => include.includes(item))
  }
  if (exclude.length) {
    queryKeys = queryKeys.filter(item => !exclude.includes(item))
  }
  const newQuery = {}
  queryKeys.forEach(key => newQuery[key] = query[key])
  return newQuery
}

export default function useRoute(root) {
  const state = reactive({ route: root.$route })

  watch(
    () => root.$route,
    (route) => {
      state.route = route
    }
  )

  return {
    ...toRefs(state)
  }
}
