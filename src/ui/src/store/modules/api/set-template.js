import $http from '@/api'

const state = {}

const getters = {}

const mutations = {}

const actions = {
    createSetTemplate ({ commit }, { bizId, params, config }) {
        return $http.post(`create/topo/set_template/bk_biz_id/${bizId}`, params, config)
    },
    updateSetTemplate ({ commit }, { bizId, setTemplateId, params, config }) {
        return $http.put(`update/topo/set_template/${setTemplateId}/bk_biz_id/${bizId}`, params, config)
    },
    deleteSetTemplate ({ commit }, { bizId, params, config }) {
        return $http.delete(`deletemany/topo/set_template/bk_biz_id/${bizId}/`, params, config)
    },
    getSetTemplates ({ commit }, { bizId, params, config }) {
        return $http.post(`findmany/topo/set_template/bk_biz_id/${bizId}/web`, params, config)
    },
    getSingleSetTemplateInfo ({ commit }, { bizId, setTemplateId, config }) {
        return $http.get(`find/topo/set_template/${setTemplateId}/bk_biz_id/${bizId}`, config)
    },
    getServiceTemplateInSet ({ commit }, { bizId, setTemplateId }, config) {
        return $http.get(`findmany/topo/set_template/${setTemplateId}/bk_biz_id/${bizId}/service_templates`)
    }
}

export default {
    namespaced: true,
    state,
    getters,
    mutations,
    actions
}
