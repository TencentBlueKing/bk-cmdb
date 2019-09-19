import $http from '@/api'

const state = {}

const getters = {}

const mutations = {}

const actions = {
    diffTemplateAndInstances ({ commit, state, dispatch, rootGetters }, { bizId, setTemplateId, params, config }) {
        return $http.post(`findmany/topo/set_template/${setTemplateId}/bk_biz_id/${bizId}/diff_with_instances`, params, config)
    },
    syncTemplateToInstances ({ commit }, { bizId, setTemplateId, params, config }) {
        return $http.post(`updatemany/topo/set_template/${setTemplateId}/sync_to_instances/bk_biz_id/${bizId}`, params, config)
    },
    getSetTemplates ({ commit }, { bizId, params, config }) {
        return $http.post(`findmany/topo/set_template/bk_biz_id/${bizId}`, params, config)
    },
    getSingleSetTemplateInfo ({ commit }, { bizId, setTemplateId, config }) {
        return $http.get(`find/topo/set_template/${setTemplateId}/bk_biz_id/${bizId}`, config)
    }
}

export default {
    namespaced: true,
    state,
    getters,
    mutations,
    actions
}
