import $http from '@/api'

/**
 * 查询后代
 * @param {number} bizSetId 业务集 ID
 * @param {string} parentModelId 上一级的模型 ID，通常为模型类型比如 set
 * @param {number} parentInstanceId 上一级的实例 ID
 * @returns {Promise}
 */
export const findChildren = ({
  bizSetId,
  parentModelId,
  parentInstanceId,
}, config) => $http.post('find/biz_set/topo_path', {
  bk_biz_set_id: bizSetId,
  bk_parent_obj_id: parentModelId,
  bk_parent_id: parentInstanceId,
}, config)


/**
 * 获取实例数量
 * @param {number} bizSetId 业务集 ID
 * @param {Array} condition 拓扑节点信息，数组最大长度为 20
 * @param {Object} config 请求配置
 * @returns {Promise}
 */
export const getInstanceCount = (bizSetId, condition, config) => $http.post(`count/topoinst/host_service_inst/biz_set/${bizSetId}`, { condition }, config)

/**
 * 获取业务下资源的拓扑路径
 * @param {number} bizSetId 业务集 ID
 * @param {number} bizId 业务 ID
 * @param {Array} condition 实例查询条件
 * @param {Object} config 请求配置
 * @returns {Promise}
 */
export const findTopoPath = ({ bizSetId, bizId }, params, config) => $http.post(`find/topopath/biz_set/${bizSetId}/biz/${bizId}`, params, config)

export const TopologyService = {
  findChildren,
  findTopoPath,
  getInstanceCount
}
