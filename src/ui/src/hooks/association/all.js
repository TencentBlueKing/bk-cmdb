import { ref } from '@vue/composition-api'
import associationService from '@/service/association'
export default function ({ immeidate = false } = {}) {
  const associations = ref([])
  const pending = ref(true)
  const findAll = async () => {
    const { info } = await associationService.findAll()
    associations.value = info
    pending.value = false
  }
  immeidate && findAll()
  return [{ associations, pending }, findAll]
}
