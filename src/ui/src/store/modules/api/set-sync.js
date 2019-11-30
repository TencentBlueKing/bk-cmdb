import $http from '@/api'

const state = {}

const getters = {}

const mutations = {}

const actions = {
    diffTemplateAndInstances ({ commit, state, dispatch, rootGetters }, { bizId, setTemplateId, params, config }) {
        return $http.post(`findmany/topo/set_template/${setTemplateId}/bk_biz_id/${bizId}/diff_with_instances`, params, config)
    },
    syncTemplateToInstances ({ commit }, { bizId, setTemplateId, params, config }) {
        return $http.post(`updatemany/topo/set_template/${setTemplateId}/bk_biz_id/${bizId}/sync_to_instances`, params, config)
    },
    getInstancesSyncStatus ({ commit }, { bizId, setTemplateId, params, config }) {
        return $http.post(`findmany/topo/set_template/${setTemplateId}/bk_biz_id/${bizId}/instances_sync_status`, params, config)
    }
}

export default {
    namespaced: true,
    state,
    getters,
    mutations,
    actions
}
