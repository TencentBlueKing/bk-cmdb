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

const host = {
    namespaced: true,
    modules: {
        transfer: {
            namespaced: true,
            actions: {
                directory (context, { params, config }) {
                    return $http.post('host/transfer/resource/directory', params, config)
                },
                idle (context, { params, config }) {
                    return $http.post('hosts/modules/resource/idle', params, config)
                }
            }
        }
    }
}

export default {
    namespaced: true,
    modules: {
        directory,
        host
    }
}
