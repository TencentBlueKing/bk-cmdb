import merge from 'lodash/merge'

/**
 * 根据是否开启count生成新参数
 * @param {Object} params 基础参数
 * @param {Boolean} flag 是否开启count获取
 * @returns 生成的新参数
 */
export const enableCount = (params = {}, flag = false) => {
  const page = Object.assign(flag ? { start: 0, limit: 0, sort: '' } : {}, { enable_count: flag })
  return merge({}, params, { page })
}

export const onePageParams = () => ({ start: 0, limit: 1 })
