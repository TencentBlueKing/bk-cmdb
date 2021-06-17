import { ref, watch } from '@vue/composition-api'
import modelRelationService from '@/service/relation/model'
export default function (modelId) {
  const relations = ref([])
  const pending = ref(false)
  const refresh = async (value) => {
    if (!value) return
    pending.value = true
    relations.value = await modelRelationService.findAll(value)
    pending.value = false
  }
  watch(modelId, refresh, { immediate: true })
  return [{ relations, pending }, refresh]
}
