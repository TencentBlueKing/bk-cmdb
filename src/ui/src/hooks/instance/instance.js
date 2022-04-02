import { reactive, isRef, watch, toRefs } from '@vue/composition-api'
import hostSearchService from '@/service/host/search'
import businessSearchService from '@/service/business/search'
import instanceSearchService from '@/service/instance/search'
import businessSetService from '@/service/business-set/index.js'
import { BUILTIN_MODELS, BUILTIN_MODEL_PROPERTY_KEYS } from '@/dictionary/model-constants.js'

const getService = ({ bk_obj_id: objId }) => {
  const modelServiceMapping = {
    [BUILTIN_MODELS.HOST]: hostSearchService,
    [BUILTIN_MODELS.BUSINESS]: businessSearchService,
    [BUILTIN_MODELS.BUSINESS_SET]: businessSetService
  }
  return modelServiceMapping[objId] || instanceSearchService
}

const getServiceOptions = (options) => {
  const idMapping = {
    [BUILTIN_MODELS.HOST]: [BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.HOST].ID],
    [BUILTIN_MODELS.BUSINESS]: [BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS].ID],
    [BUILTIN_MODELS.BUSINESS_SET]: [BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS_SET].ID]
  }
  return { ...options, [idMapping[options.bk_obj_id] || 'bk_inst_id']: options.bk_inst_id }
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
