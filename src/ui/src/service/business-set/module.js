import $http from '@/api'

/**
 * 获取模块信息
 * @param {number} bizSetId 业务集 ID
 * @param {Object} params 查询参数
 * @param {Object} config 请求配置
 * @returns {Promise}
 */
export const findOne = ({
  bizSetId,
  bizId,
  setId
}, params, config) => $http.post(`findmany/module/biz_set/${bizSetId}/biz/${bizId}/set/${setId}`, params, config)

export const ModuleService = {
  findOne
}
