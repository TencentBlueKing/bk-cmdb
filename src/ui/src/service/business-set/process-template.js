import $http from '@/api'

/**
 * 查询业务集下的进程模板
 * @param {number} bizSetId 业务集 ID
 * @param {Array} params 查询参数
 * @param {Object} config 请求配置
 * @returns {Promise}
 */
export const findAll = (bizSetId, params, config) => $http.post(`findmany/proc/biz_set/${bizSetId}/proc_template`, params, config)

export const findOne = ({ bizSetId, processTemplateId }, config) => $http.post(`find/proc/biz_set/${bizSetId}/proc_template/id/${processTemplateId}`, undefined, config)

export const ProcessTemplateService = {
  findAll,
  findOne
}
