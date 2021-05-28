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
    localForceHide.value = (!old || keyword === old)
  })

  const isForceHide = computed(() => forceHide.value || localForceHide.value)

  // eslint-disable-next-line max-len
  const showSuggestion = computed(() => focusWithin.value && !isForceHide.value && !isShowHistory.value && hasKeyword.value && hasSuggestion.value)

  return {
    suggestion,
    showSuggestion
  }
}
