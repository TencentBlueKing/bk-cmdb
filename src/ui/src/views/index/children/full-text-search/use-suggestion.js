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

import { computed, watch, ref } from '@vue/composition-api'
import useItem from './use-item.js'

export default function useSuggestion(state, root) {
  const { result, showHistory, selectHistory, focusWithin, keyword, forceHide } = state

  const list = computed(() => (result.value.hits || []).slice(0, 8))

  const { normalizationList: suggestion } = useItem(list, root)

  const isShowHistory = computed(() => showHistory.value || selectHistory.value)
  const hasSuggestion = computed(() => suggestion.value.length > 0)
  const hasKeyword = computed(() => keyword.value.length > 0)

  const localForceHide = ref(false)
  watch(keyword, (keyword, old) => {
    localForceHide.value = keyword === old
  })

  const isForceHide = computed(() => forceHide.value || localForceHide.value)

  // eslint-disable-next-line max-len
  const showSuggestion = computed(() => focusWithin.value && !isForceHide.value && !isShowHistory.value && hasKeyword.value && hasSuggestion.value)

  return {
    suggestion,
    showSuggestion
  }
}
