import $http from '@/api'

const actions = {
    getServiceInstanceProcesses (context, { params, config }) {
        return $http.post('findmany/proc/process_instance', params, config)
    },
    updateServiceInstanceProcess ({ rootGetters }, { business, processInstanceId, params, config }) {
        return $http.put(`proc/${rootGetters.supplierAccount}/${business}/${processInstanceId}`, params, config)
    },
    createServiceInstanceProcess (context, { params, config }) {
        return $http.post('create/proc/process_instance/with_raw', params, config)
    }
}

export default {
    namespaced: true,
    actions
}
