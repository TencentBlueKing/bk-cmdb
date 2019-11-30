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
    deleteSetTemplate ({ commit }, { bizId, config }) {
        return $http.delete(`deletemany/topo/set_template/bk_biz_id/${bizId}`, config)
    },
    getSetTemplates ({ commit }, { bizId, params, config }) {
        return $http.post(`findmany/topo/set_template/bk_biz_id/${bizId}/web`, params, config)
    },
    getSingleSetTemplateInfo ({ commit }, { bizId, setTemplateId, config }) {
        return $http.get(`find/topo/set_template/${setTemplateId}/bk_biz_id/${bizId}`, config)
    },
    getSetTemplateServices ({ commit }, { bizId, setTemplateId }, config) {
        return $http.get(`findmany/topo/set_template/${setTemplateId}/bk_biz_id/${bizId}/service_templates`)
    },
    getSetInstancesWithStatus ({ commit }, { bizId, params, config }) {
        return $http.post(`findmany/topo/set_template_sync_status/bk_biz_id/${bizId}`, params, config)
    },
    getSetInstancesWithTopo ({ commit }, { bizId, setTemplateId, params, config }) {
        return $http.post(`findmany/topo/set_template/${setTemplateId}/bk_biz_id/${bizId}/sets/web`, params, config)
    },
    getSyncHistory ({ commit }, { bizId, params, config }) {
        return $http.post(`findmany/topo/set_template_sync_history/bk_biz_id/${bizId}`, params, config)
    },
    getSetTemplateServicesStatistics ({ commit }, { bizId, setTemplateId }, config) {
        return $http.get(`findmany/topo/set_template/${setTemplateId}/bk_biz_id/${bizId}/service_templates/with_statistics`, config)
    }
}

export default {
    namespaced: true,
    state,
    getters,
    mutations,
    actions
}
