import $http from '@/api'
import Vue from 'vue'
const cloudArea = {
    namespaced: true,
    actions: {
        findMany ({ commit, state, dispatch }, { params, config }) {
            return $http.post(`findmany/cloudarea`, params, config)
        },
        delete (context, { id, config }) {
            return $http.delete(`delete/cloudarea/${id}`, config)
        }
    }
}

const cloudAccount = {
    namespaced: true,
    actions: {
        findMany (context, { params, config }) {
            return $http.post(`findmany/cloud/account`, params, config)
        },
        async findOne (context, { id, config }) {
            return context.dispatch('findMany', {
                params: {
                    condition: {
                        bk_account_id: {
                            $eq: id
                        }
                    }
                },
                config
            }).then(({ info }) => {
                if (!info.length) {
                    return Promise.reject(new Error('Can not find cloud account with id:' + id))
                }
                return info[0]
            })
        },
        verify (context, { params, config }) {
            return $http.post('cloud/account/verify', params, config)
        },
        create (context, { params, config }) {
            return $http.post('create/cloud/account', params, config)
        },
        update (context, { id, params, config }) {
            return $http.put(`update/cloud/account/${id}`, params, config)
        }
    }
}

const cloudResource = {
    namespaced: true,
    state: {
        regions: {}
    },
    actions: {
        createTask (context, { params, config }) {
            return $http.post('create/cloud/sync/task', params, config)
        },
        findTask (context, { params, config }) {
            return $http.post('findmany/cloud/sync/task', params, config)
        },
        findOneTask (context, { id, config }) {
            return context.dispatch('findTask', {
                params: {
                    condition: {
                        bk_task_id: {
                            $eq: id
                        }
                    }
                },
                config
            }).then(({ info }) => {
                if (!info.length) {
                    return Promise.reject(new Error('Can not find cloud task with id:' + id))
                }
                return info[0]
            })
        },
        updateTask (context, { id, params, config }) {
            return $http.put(`update/cloud/sync/task/${id}`, params, config)
        },
        deleteTask (context, { id, config }) {
            return $http.delete(`delete/cloud/sync/task/${id}`, config)
        },
        findRegion ({ state, commit }, { params, config }) {
            if (state.regions.hasOwnProperty(params.bk_account_id)) {
                return Promise.resolve(state.regions[params.bk_account_id])
            }
            return $http.post('findmany/cloud/sync/region', params, config).then(data => {
                commit('setRegions', { account: params.bk_account_id, regions: data.info })
                return data.info
            })
        },
        findVPC (context, { id, params, config }) {
            return $http.post(`findmany/cloud/account/vpc/${id}`, params, config)
        },
        findHistory (context, { params, config }) {
            return $http.post('findmany/cloud/sync/history')
        }
    },
    mutations: {
        setRegions (state, { account, regions }) {
            Vue.set(state.regions, account, regions)
        }
    }
}

export default {
    namespaced: true,
    modules: {
        area: cloudArea,
        account: cloudAccount,
        resource: cloudResource
    }
}
