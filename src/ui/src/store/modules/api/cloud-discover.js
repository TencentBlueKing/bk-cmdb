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

}

const getters = {

}

const actions = {
    searchCloudTask ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`hosts/cloud/search`, params, config)
    },

    addCloudTask ({ commit, state, dispatch }, { params }) {
        return $http.post(`hosts/cloud/add`, params)
    },

    updateCloudTask ({ commit, state, dispatch }, { params }) {
        return $http.put(`hosts/cloud/update`, params)
    },

    deleteCloudTask ({ commit, state, dispatch }, { taskID }) {
        return $http.delete(`hosts/cloud/delete/${taskID}`)
    },

    startCloudSync ({ commit, state, dispatch }, { params }) {
        return $http.post(`hosts/cloud/startSync`, params)
    },

    searchCloudHistory ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`hosts/cloud/syncHistory`, params, config)
    },

    getResourceConfirm ({ commit, state, dispatch }, { params, config }) {
        return $http.post(`hosts/cloud/searchConfirm`, params, config)
    },

    resourceConfirm ({ commit, state, dispatch }, { params }) {
        return $http.post(`hosts/cloud/resourceConfirm`, params)
    },

    addConfirmHistory ({ commit, state, dispatch }, { params }) {
        return $http.post('hosts/cloud/confirmHistory/add', params)
    },

    searchConfirmHistory ({ commit, state, dispatch }, { params, config }) {
        return $http.post('/hosts/cloud/confirmHistory/search', params, config)
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions
}
