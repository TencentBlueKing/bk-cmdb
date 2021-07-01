import hostSearchService from '@/service/host/search'
import businessSearchService from '@/service/business/search'
import instanceSearchService from '@/service/instance/search'
import { reactive, isRef, watch, toRefs } from '@vue/composition-api'

const getService = ({ bk_obj_id: objId }) => {
  const modelServiceMapping = {
    host: hostSearchService,
    biz: businessSearchService,
  }
  return modelServiceMapping[objId] || instanceSearchService
}

const getServiceOptions = (options) => {
  if (options.bk_obj_id === 'host') {
    return { ...options, bk_host_id: options.bk_inst_id }
  }
  if (options.bk_obj_id === 'biz') {
    return { ...options, bk_biz_id: options.bk_inst_id }
  }
  return options
}
/**
 * options.bk_obj_id 模型id
 * options.bk_inst_id 实例id
 * options.bk_biz_id 业务id
 */
export default function (options) {
  const state = reactive({
    instance: {},
    pending: false
  })
  const refresh = async (value) => {
    if (!value.bk_inst_id) return
    state.pending = true
    const service = getService(value)
    const serviceOptions = getServiceOptions(value)
    const instance = await service.findOne(serviceOptions)
    state.instance = instance || {}
    state.pending = false
  }
  watch(() => (isRef(options) ? options.value : options), refresh, { immediate: true, deep: true })
  return [toRefs(state), { refresh }]
}
