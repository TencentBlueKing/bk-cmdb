import http from '@/api'
import { enableCount } from '../utils.js'

const find = async (params, config) => {
  try {
    const [list, count] = await Promise.all([
      http.post('findmany/biz_set', enableCount(params, false), config),
      http.post('findmany/biz_set', enableCount(params, true), config)
    ])
    return { count, list: list || [] }
  } catch (error) {
    console.error(error)
    return { count: 0, list: [] }
  }
}

const previewOfBefore = async (params, config) => {
  try {
    const [list, count] = await Promise.all([
      http.post('find/biz_set/preview', enableCount(params, false), config),
      http.post('find/biz_set/preview', enableCount(params, true), config)
    ])
    return { count, list: list || [] }
  } catch (error) {
    console.error(error)
    return { count: 0, list: [] }
  }
}

const previewOfAfter = async (params, config) => {
  try {
    const [{ info: list }, { count }] = await Promise.all([
      http.post('find/biz_set/biz_list', enableCount(params, false), config),
      http.post('find/biz_set/biz_list', enableCount(params, true), config)
    ])
    return { count, list: list || [] }
  } catch (error) {
    console.error(error)
    return { count: 0, list: [] }
  }
}

const create = (data, config) => http.post('create/biz_set', data, config)

const deleteById = (id, config) => http.post('deletemany/biz_set', {
  bk_biz_set_ids: [id]
}, config)

export default {
  find,
  create,
  deleteById,
  previewOfBefore,
  previewOfAfter
}
