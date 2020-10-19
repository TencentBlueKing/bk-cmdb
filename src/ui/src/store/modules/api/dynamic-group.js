import $http from '@/api'
export default {
    namespaced: true,
    actions: {
        create (context, { params, config }) {
            return $http.post('dynamicgroup', params, config)
        },
        update (context, { bizId, id, params, config }) {
            return $http.put(`dynamicgroup/${bizId}/${id}`, params, config)
        },
        delete (context, { bizId, id, config }) {
            return $http.delete(`dynamicgroup/${bizId}/${id}`, config)
        },
        details (context, { bizId, id, config }) {
            return $http.get(`dynamicgroup/${bizId}/${id}`, config)
        },
        preview (context, { bizId, id, params, config }) {
            return $http.post(`dynamicgroup/data/${bizId}/${id}`, params, config)
        },
        search (context, { bizId, params, config }) {
            return $http.post(`dynamicgroup/search/${bizId}`, params, config)
        }
    }
}
