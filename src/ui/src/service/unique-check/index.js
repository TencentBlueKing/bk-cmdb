import http from '@/api'
export const find = async (modelId) => {
  try {
    const result = await http.post(`find/objectunique/object/${modelId}`)
    return result
  } catch (error) {
    console.error(error)
    return []
  }
}

export const findMany = async models => Promise.all(models.map(modelId => find(modelId)))

export default {
  find,
  findMany
}
