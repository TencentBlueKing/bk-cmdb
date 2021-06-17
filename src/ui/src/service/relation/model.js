import http from '@/api'

export const find = async (modelId, type) => {
  try {
    const key = type === 'source' ? 'bk_obj_id' : 'bk_asst_obj_id'
    const result = await http.post('find/objectassociation', { condition: { [key]: modelId } })
    return result
  } catch (error) {
    console.error(error)
    return []
  }
}

export const findAsSource = modelId => find(modelId, 'source')
export const findAsTarget = modelId => find(modelId, 'target')
export const findAll = async (modelId) => {
  const [source, target] = await Promise.all([findAsSource(modelId), findAsTarget(modelId)])
  const all = [...source, ...target]
  const uniqId = [...new Set(all.map(item => item.id))]
  return uniqId.map(id => all.find(item => item.id === id))
}

export default {
  find,
  findAsSource,
  findAsTarget,
  findAll
}
