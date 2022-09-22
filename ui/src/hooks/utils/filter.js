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

import { toRef, reactive, watch } from '@vue/composition-api'
import debounce from 'lodash.debounce'

export default function ({ list, keyword, target }) {
  const state = reactive({
    result: []
  })
  const handler = (value) => {
    if (!value) {
      state.result = list.value
      return
    }
    const regexp = new RegExp(value, 'ig')
    state.result = list.value.filter(item => regexp.test(item[target]))
  }
  const filter = debounce(handler, 300, { leading: false, trailing: true })
  watch(keyword, filter)
  watch(list, () => handler(), { immediate: true })
  return [toRef(state, 'result'), { filter }]
}
