import $http from '@/api'

/**
 * 获取业务信息
 * @param {Object} pathParams
 * @param {number} pathParams.bizSetId 业务集 ID
 * @param {number} pathParams.bizId 业务 ID
 * @param {Object} params 查询参数
 * @param {Object} config 请求配置
 * @returns {Promise}
 */
export const findOne = ({
  bizSetId,
  bizId,
}, config) => $http.post('find/biz_set/biz_list', {
  bk_biz_set_id: bizSetId,
  filter: {
    condition: 'AND',
    rules: [
      {
        field: 'bk_biz_id',
        operator: 'equal',
        value: bizId
      }
    ]
  },
  page: {
    start: 0,
    limit: 1
  }
}, config)

export const BusinessService = {
  findOne
}
