import $http from '@/api'

/**
 * 查询业务集下的服务进程列表
 * @param {*} bizSetId
 * @param {*} params
 * @param {*} config
 * @returns
 */
export const findAll = (bizSetId, params, config) => $http.post(`findmany/proc/biz_set/${bizSetId}/process_instance/name_ids`, params, config)

/**
 * 根据服务实例查询进程
 * @param {number} bizSetId 业务集 ID
 * @param {Object} params 查询参数
 * @param {Object} config 请求配置
 * @returns {Promise}
 */
export const findProcessByServiceInstance = (bizSetId, params, config) => $http.post(`findmany/proc/biz_set/${bizSetId}/process_instance`, params, config)

/**
 * 根据进程查询对应的服务实例
 * @param {number} bizSetId 业务集 ID
 * @param {Object} params 查询参数
 * @param {Object} config 请求配置
 * @returns {Promise}
 */
export const findServiceInstanceByProcess = (bizSetId, params, config) => $http.post(`/findmany/proc/biz_set/${bizSetId}/process_instance/detail/by_ids`, params, config)

export const ProcessInstanceService = {
  findAll,
  findProcessByServiceInstance,
  findServiceInstanceByProcess
}
