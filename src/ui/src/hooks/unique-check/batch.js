import { ref, watch } from '@vue/composition-api'
import uniqueCheckService from '@/service/unique-check'
export default function (models) {
  const uniqueChecks = ref([])
  const pending = ref(false)
  const refresh = async (value) => {
    if (!value.length) return
    pending.value = true
    uniqueChecks.value = await uniqueCheckService.findMany(value)
    pending.value = false
  }
  watch(models, refresh, { immediate: true })
  return [{ uniqueChecks, pending }, { refresh }]
}
