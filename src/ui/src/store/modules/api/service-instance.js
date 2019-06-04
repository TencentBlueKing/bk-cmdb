import $http from '@/api'

const actions = {
    getModuleServiceInstances (context, { params, config }) {
        return $http.post('find/proc/service_instance', params, config)
    }
}

export default {
    namespaced: true,
    actions
}
