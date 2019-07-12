import $http from '@/api'

const actions = {
    createInstanceLabel (context, { params, config }) {
        return $http.post('createmany/proc/service_instance/labels', params, config)
    },
    deleteInstanceLabel (context, { config }) {
        return $http.delete('deletemany/proc/service_instance/labels', config)
    },
    getHistoryLabel (context, { params, config }) {
        return $http.post('findmany/proc/service_instance/labels/aggregation', params, config)
    }
}

export default {
    namespaced: true,
    actions
}
