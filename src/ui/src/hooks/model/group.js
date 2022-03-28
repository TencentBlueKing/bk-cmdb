import { reactive, toRefs, isRef, watch } from '@vue/composition-api'
import groupService from '@/service/property/group'
export default function (options = {}) {
  const state = reactive({
    groups: [],
    pending: false
  })
  const refresh = async (value) => {
    if (!value.bk_obj_id) return
    state.pending = true
    state.groups = await groupService.find(value)
    state.pending = false
  }
  watch(() => (isRef(options) ? options.value : options), refresh, { immediate: true, deep: true })
  return [toRefs(state), { refresh }]
}
