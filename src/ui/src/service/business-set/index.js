import http from '@/api'
import { BUILTIN_MODELS, BUILTIN_MODEL_PROPERTY_KEYS } from '@/dictionary/model-constants.js'
import { enableCount, onePageParams } from '../utils.js'

const authorizedRequsetId = Symbol('getAuthorizedBusinessSet')
const MODEL_ID_KEY = BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS_SET].ID

const find = async (params, config) => {
  try {
    const [{ info: list = [] }, { count = 0 }] = await Promise.all([
      http.post('findmany/biz_set', enableCount(params, false), config),
      http.post('findmany/biz_set', enableCount(params, true), config)
    ])
    return { count, list }
  } catch (error) {
    console.error(error)
    return Promise.reject(error)
  }
}

const findById = async (id, config) => {
  try {
    const { info: [instance = null] } = await http.post('findmany/biz_set', enableCount({
      bk_biz_set_filter: {
        condition: 'AND',
        rules: [{
          field: MODEL_ID_KEY,
          operator: 'equal',
          value: id
        }]
      },
      page: onePageParams()
    }, false), config)

    return instance
  } catch (error) {
    console.error(error)
    return null
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
  requestId: authorizedRequsetId,
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

const update = (data, config) => http.put('updatemany/biz_set', data, config)

const deleteById = (id, config) => http.post('deletemany/biz_set', {
  bk_biz_set_ids: [id]
}, config)

export default {
  find,
  findById,
  create,
  update,
  deleteById,
  previewOfBeforeCreate,
  previewOfAfterCreate,
  getAuthorized,
  getAuthorizedWithCache
}
