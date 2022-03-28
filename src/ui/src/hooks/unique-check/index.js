import { ref, watch } from '@vue/composition-api'
import uniqueCheckService from '@/service/unique-check'
export default function (modelId) {
  const uniqueChecks = ref([])
  const pending = ref(false)
  const refresh = async (value) => {
    if (!value.length) return
    pending.value = true
    uniqueChecks.value = await uniqueCheckService.find(value)
    pending.value = false
  }
  watch(modelId, refresh, { immediate: true })
  return [{ uniqueChecks, pending }, { refresh }]
}
