import http from '@/api'
import { enableCount } from '../utils.js'

const authorizeRequsetId = Symbol('getAuthorizedBusinessSet')

const find = async (params, config) => {
  try {
    const [{ info: list = [] }, { count = 0 }] = await Promise.all([
      http.post('findmany/biz_set', enableCount(params, false), config),
      http.post('findmany/biz_set', enableCount(params, true), config)
    ])
    return { count, list }
  } catch (error) {
    console.error(error)
    return { count: 0, list: [] }
  }
}

const getAuthorized = async (config) => {
  try {
    const { info: list = [] } = await http.get('findmany/biz_set/with_reduced?sort=bk_biz_set_id', config)
    return list
  } catch (error) {
    console.error(error)
    return []
  }
}

const getAuthorizedWithCache = async () => getAuthorized({
  requestId: authorizeRequsetId,
  fromCache: true
})

const previewOfBeforeCreate = async (params, config) => {
  try {
    const [{ info: list = [] }, { count = 0 }] = await Promise.all([
      http.post('find/biz_set/preview', enableCount(params, false), config),
      http.post('find/biz_set/preview', enableCount(params, true), config)
    ])
    return { count, list }
  } catch (error) {
    console.error(error)
    return { count: 0, list: [] }
  }
}

const previewOfAfterCreate = async (params, config) => {
  try {
    const [{ info: list }, { count = 0 }] = await Promise.all([
      http.post('find/biz_set/biz_list', enableCount(params, false), config),
      http.post('find/biz_set/biz_list', enableCount(params, true), config)
    ])
    return { count, list }
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
  previewOfBeforeCreate,
  previewOfAfterCreate,
  getAuthorized,
  getAuthorizedWithCache
}
