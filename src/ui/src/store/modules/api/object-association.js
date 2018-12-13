import $http from '@/api'

const state = {

}

const getters = {

}

const actions = {
    /**
     * 查询关联类型
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchAssociationType ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`topo/association/type/action/search`, params, config)
    },
    /**
     * 添加关联类型
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createAssociationType ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`topo/association/type/action/create`, params, config)
    },
    /**
     * 编辑关联类型
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} id 自增id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updateAssociationType ({ commit, state, dispatch }, { id, params, config }) {
        return $http.put(`topo/association/type/${id}/action/update`, params, config)
    },
    /**
     * 删除关联类型
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} id 自增id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    deleteAssociationType ({ commit, state, dispatch }, { id, config }) {
        return $http.delete(`topo/association/type/${id}/action/delete`, config)
    },
    /**
     * 查询模型关联
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchObjectAssociation ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`object/association/action/search`, params, config)
    },
    /**
     * 添加模型关联
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createObjectAssociation ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`object/association/action/create`, params, config)
    },
    /**
     * 编辑模型关联
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} id 自增id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updateObjectAssociation ({ commit, state, dispatch }, { id, params, config }) {
        return $http.put(`object/association/${id}/action/update`, params, config)
    },
    /**
     * 删除模型关联
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} id 自增id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    deleteObjectAssociation ({ commit, state, dispatch }, { id, config }) {
        return $http.delete(`object/association/${id}/action/delete`, config)
    },
    /**
     * 根据关联类型查询使用这些关联类型的关联关系列表
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchAssociationListWithAssociationKindList ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`topo/association/type/action/search/batch`, params, config)
    },
    /**
     * 查询实例关联
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchInstAssociation ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`inst/association/action/search`, params, config)
    },
    /**
     * 添加实例关联
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createInstAssociation ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`inst/association/action/create`, params, config)
    },
    /**
     * 删除实例关联
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    deleteInstAssociation ({ commit, state, dispatch }, { id, config }) {
        return $http.delete(`inst/association/${id}/action/delete`, config)
    }
}

const mutations = {

}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
