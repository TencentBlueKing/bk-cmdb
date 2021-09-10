import { reactive, toRefs, set, del } from '@vue/composition-api'
import useClone from '@/hooks/utils/clone'
const defaultState = {
  visible: false,
  title: '',
  bk_obj_id: null,
  bk_biz_id: null,
  available: () => true,
  submit: () => {},
  count: 0,
  limit: 10000,
  step: 1,
  status: null,
  presetFields: [],
  defaultSelectedFields: [],
  fields: [],
  relations: {},
  exportRelation: false,
  object_unique_id: ''
}

const state = reactive(useClone(defaultState))

const setState = (newState) => {
  Object.assign(state, newState)
}

const resetState = () => setState(useClone(defaultState))
const resetPartial = () => setState({
  step: 1,
  status: null,
  relations: {},
  exportRelation: false,
  object_unique_id: ''
})

const setRelation = (modelId, uniqueId) => set(state.relations, modelId, uniqueId)
const removeRelation = modelId => del(state.relations, modelId)

export default function () {
  return [toRefs(state), { setState, resetState, resetPartial, setRelation, removeRelation }]
}
