import http from '@/api'

export const find = async (params, config) => {
  try {
    const properties = await http.post('find/objectattr', params, config)
    return properties
  } catch (error) {
    console.error(error)
    return []
  }
}

export default {
  find
}
