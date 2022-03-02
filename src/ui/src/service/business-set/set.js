import $http from '@/api'

/**
 * 获取集群信息
 * @param {Object} pathParams
 * @param {number} pathParams.bizSetId 业务集 ID
 * @param {number} pathParams.bizId 业务 ID
 * @param {Object} params 查询参数
 * @param {Object} config 请求配置
 * @returns {Promise}
 */
export const findOne = ({
  bizSetId,
  bizId
}, params, config) => $http.post(`findmany/set/biz_set/${bizSetId}/biz/${bizId}`, params, config)

export const SetService = {
  findOne
}
