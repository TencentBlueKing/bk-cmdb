/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

import $http from '@/api'

const state = {
    privilege: {
        'model_config': {},
        'sys_config': {
            'global_busi': null,
            'back_config': null
        }
    },
    roles: []
}

const getters = {
    privilege: state => state.privilege,
    roles: state => state.roles
}

const actions = {
    /**
     * 获取角色绑定权限
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkObjId 对象id
     * @param {String} bkPropertyId 属性id
     * @return {promises} promises 对象
     */
    getRolePrivilege ({ commit, state, dispatch, rootGetters }, { bkObjId, bkPropertyId }) {
        return $http.get(`topo/privilege/${rootGetters.supplierAccount}/${bkObjId}/${bkPropertyId}`)
    },

    /**
     * 绑定角色权限
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkObjId 对象id
     * @param {String} bkPropertyId 属性id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    bindRolePrivilege ({ commit, state, dispatch, rootGetters }, { bkObjId, bkPropertyId, params, config }) {
        return $http.post(`topo/privilege/${rootGetters.supplierAccount}/${bkObjId}/${bkPropertyId}`, params, config)
    },

    /**
     * 新建用户分组
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createUserGroup ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.post(`topo/privilege/group/${rootGetters.supplierAccount}`, params, config)
    },

    /**
     * 更新用户分组
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkGroupId 分组id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updateUserGroup ({ commit, state, dispatch, rootGetters }, { bkGroupId, params, config }) {
        return $http.put(`topo/privilege/group/${rootGetters.supplierAccount}/${bkGroupId}`, params, config)
    },

    /**
     * 查询用户分组
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchUserGroup ({ commit, state, dispatch, rootGetters }, { params, config }) {
        return $http.post(`topo/privilege/group/${rootGetters.supplierAccount}/search`, params, config)
    },

    /**
     * 删除用户分组
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkGroupId 分组id
     * @return {promises} promises 对象
     */
    deleteUserGroup ({ commit, state, dispatch, rootGetters }, { bkGroupId }) {
        return $http.delete(`topo/privilege/group/${rootGetters.supplierAccount}/${bkGroupId}`)
    },

    /**
     * 查询分组权限
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkGroupId 分组id
     * @return {promises} promises 对象
     */
    searchUserPrivilege ({ commit, state, dispatch, rootGetters }, { bkGroupId, config }) {
        return $http.get(`topo/privilege/group/detail/${rootGetters.supplierAccount}/${bkGroupId}`, config)
    },

    /**
     * 查询用户权限
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} userName 用户名
     * @return {promises} promises 对象
     */
    getUserPrivilege ({ commit, state, dispatch, rootGetters }, httpConfig) {
        return $http.get(`topo/privilege/user/detail/${rootGetters.supplierAccount}/${rootGetters.userName}`, httpConfig).then(privilege => {
            commit('setUserPrivilege', privilege)
        })
    },

    /**
     * 更新分组权限
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkGroupId 分组id
     * @return {promises} promises 对象
     */
    updateGroupPrivilege ({ commit, state, dispatch, rootGetters }, { bkGroupId, params, config }) {
        return $http.post(`topo/privilege/group/detail/${rootGetters.supplierAccount}/${bkGroupId}`, params, config)
    }
}

const mutations = {
    setUserPrivilege (state, privilege) {
        state.privilege = privilege
    },
    setRoles (state, roles) {
        state.roles = roles
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
