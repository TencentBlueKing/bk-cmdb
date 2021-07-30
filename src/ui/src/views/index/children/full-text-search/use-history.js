import { computed, ref, watch } from '@vue/composition-api'

export default function useHistory(state, root) {
  const { $store, $routerActions, $route } = root
  const { keyword, focusWithin, forceHide } = state

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
