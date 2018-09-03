import $http from '@/api'

const actions = {
    /**
     * 获取实例关联关系
     * @param {String} objId 模型ID
     * @param {String} instId 实例ID
     * @param {Object} config API请求配置
     * @return {Promise} promise 对象
     */
    getInstRelation ({ commit, state, dispatch, rootGetters }, { objId, instId, config }) {
        return $http.post(`inst/association/topo/search/owner/${rootGetters.supplierAccount}/object/${objId}/inst/${instId}`, {}, config)
    }
}

export default {
    namespaced: true,
    actions
}
