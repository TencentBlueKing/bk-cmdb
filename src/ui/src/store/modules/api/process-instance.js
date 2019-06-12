import $http from '@/api'

const actions = {
    getServiceInstanceProcesses (context, { params, config }) {
        return $http.post('findmany/proc/process_instance', params, config)
    },
    updateServiceInstanceProcess ({ rootGetters }, { params, config }) {
        return $http.put('update/proc/process_instance/with_raw', params, config)
    },
    createServiceInstanceProcess (context, { params, config }) {
        return $http.post('create/proc/process_instance/with_raw', params, config)
    },
    deleteServiceInstanceProcess (context, { serviceInstanceId, config }) {
        return $http.delete(`delete/proc/service_instance/${serviceInstanceId}/process`, config)
    }
}

export default {
    namespaced: true,
    actions
}
