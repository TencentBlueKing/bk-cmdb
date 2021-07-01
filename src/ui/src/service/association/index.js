import http from '@/api'
export const find = async (params = {}) => {
  try {
    const result = await http.post('find/associationtype', params)
    return result
  } catch (error) {
    console.error(error)
    return { count: 0, info: [] }
  }
}

export const findAll = () => find()

export default {
  find,
  findAll
}
