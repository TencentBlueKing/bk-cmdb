import $http from '@/api'

/**
 * @typedef ParentLayers 父级依赖关系
 * @property {string} resource_type 父级依赖，例如实例依赖模型
 * @property {string} resource_id 父级依赖id
 */

/**
 * @typedef Resource 权限信息
 * @property {string} action 需要鉴权的操作
 * @property {string} resource_type 需要鉴权的资源类型
 * @property {string} [bk_biz_id] 对业务下的操作进行鉴权需要提供业务 ID
 * @property {string} [resource_id] 资源 ID
 * @property {Array.<ParentLayers>} [parent_layers]
 */

/**
 * 鉴定用户是否有某资源的某操作的权限
 * @param {Array.<Resource>} resources 权限信息列表
 * @returns {Promise}
 */
export const verifyAuth = resources => $http.post('auth/verify', {
  resources
})
