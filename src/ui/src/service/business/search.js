import http from '@/api'
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

export default {
  findOne
}
