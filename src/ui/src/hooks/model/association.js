import { ref, watch } from '@vue/composition-api'
import modelAssociationService from '@/service/model/association'
export default function (modelId) {
  const relations = ref([])
  const pending = ref(false)
  const refresh = async (value) => {
    if (!value) return
    pending.value = true
    relations.value = await modelAssociationService.findAll(value)
    pending.value = false
  }
  watch(modelId, refresh, { immediate: true })
  return [{ relations, pending }, refresh]
}
