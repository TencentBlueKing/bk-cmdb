/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

import { $Axios, $axios } from '@/api/axios'

const state = {

}

const getters = {

}

const actions = {
    /**
     * 获取角色绑定权限
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {String} bkObjId 对象id
     * @param {String} bkPropertyId 属性id
     * @return {promises} promises 对象
     */
    getRolePrivilege ({ commit, state, dispatch }, { bkSupplierAccount, bkObjId, bkPropertyId }) {
        return $axios.get(`topo/privilege/${bkSupplierAccount}/${bkObjId}/${bkPropertyId}`)
    },

    /**
     * 绑定角色权限
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {String} bkObjId 对象id
     * @param {String} bkPropertyId 属性id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    bindRolePrivilege ({ commit, state, dispatch }, { bkSupplierAccount, bkObjId, bkPropertyId, params }) {
        return $axios.post(`topo/privilege/${bkSupplierAccount}/${bkObjId}/${bkPropertyId}`, params)
    },

    /**
     * 新建用户分组
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createUserGroup ({ commit, state, dispatch }, { bkSupplierAccount, params }) {
        return $axios.post(`topo/privilege/group/${bkSupplierAccount}`, params)
    },

    /**
     * 更新用户分组
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {String} bkGroupId 分组id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updateUserGroup ({ commit, state, dispatch }, { bkSupplierAccount, bkGroupId, params }) {
        return $axios.put(`topo/privilege/group/${bkSupplierAccount}/${bkGroupId}`, params)
    },

    /**
     * 查询用户分组
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchUserGroup ({ commit, state, dispatch }, { bkSupplierAccount, params }) {
        return $axios.post(`topo/privilege/group/${bkSupplierAccount}/search`, params)
    },

    /**
     * 删除用户分组
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {String} bkGroupId 分组id
     * @return {promises} promises 对象
     */
    deleteUserGroup ({ commit, state, dispatch }, { bkSupplierAccount, params }) {
        return $axios.delete(`topo/privilege/group/${bkSupplierAccount}/${bkGroupId}`)
    },

    /**
     * 查询分组权限
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {String} bkGroupId 分组id
     * @return {promises} promises 对象
     */
    searchUserPrivilege ({ commit, state, dispatch }, { bkSupplierAccount, bkGroupId }) {
        return $axios.get(`topo/privilege/group/detail/${bkSupplierAccount}/${bkGroupId}`)
    },

    /**
     * 查询用户权限
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {String} userName 用户名
     * @return {promises} promises 对象
     */
    getUserPrivilege ({ commit, state, dispatch }, { bkSupplierAccount, userName }) {
        return $axios.get(`topo/privilege/group/detail/${bkSupplierAccount}/${userName}`)
    },

    /**
     * 更新分组权限
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {String} bkGroupId 分组id
     * @return {promises} promises 对象
     */
    updateGroupPrivilege ({ commit, state, dispatch }, { bkSupplierAccount, bkGroupId }) {
        return $axios.post(`topo/privilege/group/detail/${bkSupplierAccount}/${bkGroupId}`, params)
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
