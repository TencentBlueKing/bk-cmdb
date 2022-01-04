import http from '@/api'

const findAllRequsetId = Symbol('findAllRequsetId')

const findOne = async ({ bk_biz_id: bizId, config }) => {
  try {
    const { info } = await http.post(`biz/search/${window.Supplier.account}`, {
      condition: { bk_biz_id: { $eq: bizId } },
      fields: [],
      page: { start: 0, limit: 1 }
    }, config)
    const [instance] = info || [null]
    return instance
  } catch (error) {
    console.error(error)
    return null
  }
}

const findAll = async () => {
  const data = await http.get('biz/simplify?sort=bk_biz_id', {
    requestId: findAllRequsetId,
    fromCache: true
  })

  return Object.freeze(data.info || [])
}

export default {
  findOne,
  findAll
}
