import $http from '@/api'

const actions = {
    getHostServiceInstances (context, { params, config }) {
        return $http.post('findmany/proc/web/service_instance/with_host', params, config)
    },
    getModuleServiceInstances (context, { params, config }) {
        return $http.post('findmany/proc/web/service_instance', params, config)
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
    },
    getInstanceIpByHost (context, { hostId, config }) {
        return $http.get(`${window.API_HOST}hosts/${hostId}/listen_ip_options`, config)
    }
}

export default {
    namespaced: true,
    actions
}
