import $http from '@/api'

const actions = {
    getModuleServiceInstances (context, { params, config }) {
        return $http.post('find/proc/service_instance', params, config)
    },
    createProcServiceInstanceWithRaw (context, { params, config }) {
        return $http.post('create/proc/service_instance/with_raw', params, config)
    },
    createProcServiceInstanceByTemplate (context, { params, config }) {
        return $http.post('create/proc/service_instance/with_template', params, config)
    },
    deleteServiceInstance (context, { serviceInstanceId, config }) {
        return $http.delete('delete/proc/service_instance', config)
    }
}

export default {
    namespaced: true,
    actions
}
