import $http from '@/api'
const directory = {
    namespaced: true,
    actions: {
        create (context, { params, config }) {
            return $http.post('create/resource/directory', params, config)
        },
        delete (context, { id, config }) {
            return $http.delete(`delete/resource/directory/${id}`, config)
        },
        update (context, { id, params, config }) {
            return $http.put(`update/resource/directory/${id}`, params, config)
        },
        findMany (context, { params, config }) {
            return $http.post('findmany/resource/directory', params, config)
        }
    }
}

export default {
    namespaced: true,
    modules: {
        directory
    }
}
