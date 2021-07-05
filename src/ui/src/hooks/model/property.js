import { reactive, toRefs, watch, isRef } from '@vue/composition-api'
import propertyService from '@/service/property/property'
export default function (options = {}) {
  const state = reactive({
    properties: [],
    pending: false
  })
  const refresh = async (value) => {
    if (!value.bk_obj_id) return
    state.pending = true
    state.properties = await propertyService.find(value)
    state.pending = false
  }
  watch(() => (isRef(options) ? options.value : options), refresh, { immediate: true, deep: true })
  return [toRefs(state), { refresh }]
}
