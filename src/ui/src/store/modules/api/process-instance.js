import $http from '@/api'

const actions = {
    getServiceInstanceProcesses (context, { params, config }) {
        return $http.post('findmany/proc/process_instance', params, config)
    }
}

export default {
    namespaced: true,
    actions
}
