import $http from '@/api'

const actions = {
    getServiceInstanceProcesses (context, { params, config }) {
        return $http.post('findmany/proc/process_instance', params, config)
    },
    updateServiceInstanceProcess (context, { processInstanceId, params, config }) {
        return $http.put(`update/proc/proc_instance/${processInstanceId}`, params, config)
    }
}

export default {
    namespaced: true,
    actions
}
