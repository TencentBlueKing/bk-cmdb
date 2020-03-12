import $http from '@/api'

const actions = {
    getDirectoryList (context, { params, config }) {
        return $http.post('findmany/resource/directory', params, config)
    },
    createDirectory (context, { params, config }) {
        return $http.post('create/resource/directory', params, config)
    },
    updateDirectory (context, { moduleId, params, config }) {
        return $http.put(`update/resource/directory/${moduleId}`, params, config)
    },
    deleteDirectory (context, { moduleId, config }) {
        return $http.delete(`delete/resource/directory/${moduleId}`, config)
    },
    changeHostsDirectory (context, { params, config }) {
        return $http.post(`host/transfer/resource/directory`, params, config)
    },
    assignHostsToBusiness (context, { params, config }) {
        return $http.post('hosts/modules/resource/idle', params, config)
    }
}

export default {
    namespaced: true,
    actions
}
