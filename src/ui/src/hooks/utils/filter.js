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
