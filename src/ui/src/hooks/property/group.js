import { reactive, toRef } from '@vue/composition-api'
import groupService from '@/service/property/group'
export default function (options) {
  const state = reactive({
    result: []
  })
  ;(async () => {
    state.result = await groupService.find(options)
  })()
  return [toRef(state, 'result')]
}
