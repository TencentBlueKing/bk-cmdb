import { computed, ref, watch } from '@vue/composition-api'

export default function (pendings = [], initValue = false) {
  let timer = null
  const pending = ref(initValue)
  const realtimePending = computed(() => pendings.some(pending => pending.value))
  watch(realtimePending, (value) => {
    if (value) {
      timer && clearTimeout(timer)
      pending.value = value
    } else {
      timer = setTimeout(() => {
        pending.value = value
      }, 200)
    }
  })
  return pending
}
