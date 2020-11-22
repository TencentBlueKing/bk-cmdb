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
    createProcServiceInstancePreview (context, { params, config }) {
        return $http.post('create/proc/service_instance/preview', params, config)
    },
    deleteServiceInstance (context, { config }) {
        return $http.delete('deletemany/proc/service_instance', config)
    },
    removeServiceTemplate (context, { config }) {
        return $http.delete('delete/proc/template_binding_on_module', config)
    },
    getInstanceIpByHost (context, { hostId, config }) {
        return $http.get(`${window.API_HOST}hosts/${hostId}/listen_ip_options`, config)
    },
    previewDeleteServiceInstances (context, { params, config }) {
        return $http.post('deletemany/proc/service_instance/preview', params, config)
    },
    getMoudleProcessList (context, { params, config }) {
        return $http.post(`findmany/proc/process_instance/name_ids`, params, config)
    },
    getProcessListById (context, { params, config }) {
        return $http.post('findmany/proc/process_instance/detail/by_ids', params, config)
    },
    batchUpdateProcess (context, { params, config }) {
        return $http.put('update/proc/process_instance/by_ids', params, config)
    },
    updateServiceInstance (context, { bizId, params, config }) {
        return $http.put(`updatemany/proc/service_instance/biz/${bizId}`, params, config)
    }
}

export default {
    namespaced: true,
    actions
}
