import $http from '@/api'

const actions = {
    getModuleServiceInstances (context, { params, config }) {
        return $http.post('find/proc/service_instance', params, config)
    },
    createProcServiceInstanceWithRaw (context, { params, config }) {
        return $http.post('create/proc/service_instance', params, config)
    },
    createProcServiceInstanceByTemplate (context, { params, config }) {
        return $http.post('create/proc/service_instance', params, config)
    },
    deleteServiceInstance (context, { config }) {
        return $http.delete('deletemany/proc/service_instance', config)
    },
    removeServiceTemplate (context, { config }) {
        return $http.delete('delete/proc/template_binding_on_module', config)
    }
}

export default {
    namespaced: true,
    actions
}
