import { $axios, $alertMsg } from '@/api/axios'

const state = {
    usercustom: {},
    result: false,
    promise: null
}

const getters = {
    usercustom: state => state.usercustom,
    result: state => state.result
}

const actions = {
    getUserCustom ({commit, state}) {
        if (state.promise) {
            return state.promise
        }
        if (state.result) {
            return Promise.resolve({result: true, data: state.usercustom})
        }
        state.promise = $axios.post('usercustom/user/search', {}).then(res => {
            state.result = res.result
            state.promise = null
            if (res.result) {
                commit('setUserCustom', res.data)
            } else {
                $alertMsg(res['bk_error_msg'])
            }
            return res
        }).catch(() => {
            state.result = false
            state.promise = null
        })
        return state.promise
    },
    updateUserCustom ({commit, state}, usercustom) {
        return $axios.post('usercustom', JSON.stringify(usercustom)).then(res => {
            if (res.result) {
                commit('setUserCustom', usercustom)
            } else {
                $alertMsg(res['bk_error_msg'])
            }
            return res
        })
    }
}

const mutations = {
    setUserCustom (state, usercustom) {
        state.usercustom = Object.assign({}, state.usercustom, usercustom)
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
