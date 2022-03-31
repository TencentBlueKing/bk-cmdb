import http from '@/api'

export const requestIds = {
  getTopology: Symbol('getTopology')
}

const getWithStat = async (bizId, config = {}) => {
  try {
    const res = await http.post(`find/topoinst_with_statistics/biz/${bizId}`, {}, {
      requestId: requestIds.getTopology,
      ...config
    })
    return res
  } catch (error) {
    console.error(error)
  }
}

const geFulltWithStat = async (bizId, config = {}) => {
  try {
    const res = await getWithStat(bizId, {
      ...config,
      params: {
        with_default: 1
      }
    })
    return res
  } catch (error) {
    console.error(error)
  }
}

export default {
  getWithStat,
  geFulltWithStat
}
