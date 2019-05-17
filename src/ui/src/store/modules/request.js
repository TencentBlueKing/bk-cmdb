import $http from '@/api'

const state = {
    queue: $http.queue.queue,
    cache: $http.cache.cache
}

const getters = {
    queue: state => state.queue,
    cache: state => state.cache
}

export default {
    namespaced: true,
    state,
    getters
}
