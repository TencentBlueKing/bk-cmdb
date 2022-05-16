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

import { computed, ref, watch } from '@vue/composition-api'
import useRoute from './use-route.js'

export default function useHistory(state, root) {
  const { $store, $routerActions, $route } = root
  const { keyword, focusWithin, forceHide } = state
  const { route } = useRoute(root)

  $store.commit('fullTextSearch/getSearchHistory')
  const historyList = computed(() => $store.state.fullTextSearch.searchHistory)

  const selectHistory = ref(false)
  const selectIndex = ref(-1)
  const localForceHide = ref(false)

  watch(keyword, (value) => {
    if (!value) {
      selectIndex.value = -1
      localForceHide.value = false
    }
  })

  const isSelectHistory = computed(() => !keyword.value.length || selectHistory.value)
  const hasHistory = computed(() => historyList.value.length > 0)
  const isForceHide = computed(() => forceHide.value || localForceHide.value)
  // eslint-disable-next-line max-len
  const showHistory = computed(() => focusWithin.value && !isForceHide.value && isSelectHistory.value && hasHistory.value)

  const onkeydown = (event) => {
    const { keyCode } = event
    const keyCodeMap = { enter: 13, up: 38, down: 40 }
    if (!showHistory.value || !Object.values(keyCodeMap).includes(keyCode)) {
      return
    }
    selectHistory.value = true
    const maxLen = historyList.value.length - 1
    let index = selectIndex.value
    if (keyCode === keyCodeMap.down) {
      index = Math.min(index + 1, maxLen)
    } else if (keyCode === keyCodeMap.up) {
      index = Math.max(index - 1, 0)
    }
    selectIndex.value = index
    keyword.value = historyList.value[selectIndex.value]
  }

  const handlClearHistory = () => {
    $store.commit('fullTextSearch/clearSearchHistory')
  }

  const handleHistorySearch = (history) => {
    localForceHide.value = true
    $routerActions.redirect({
      name: $route.name,
      query: {
        ...route.value.query,
        keyword: history,
        t: Date.now()
      }
    })
  }

  return {
    historyList,
    showHistory,
    selectHistory,
    selectIndex,
    onkeydown,
    handleHistorySearch,
    handlClearHistory
  }
}
