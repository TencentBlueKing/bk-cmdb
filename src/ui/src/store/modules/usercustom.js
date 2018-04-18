import { $axios } from '@/api/axios'
import bkMessage from '@/magicbox/bk-magic/components/message'
const alertMsg = (message, theme = 'error', delay = 3000) => {
    bkMessage({
        message,
        theme,
        delay
    })
}

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
                alertMsg(res.message)
            }
        }).catch(() => {
            state.result = false
            state.promise = null
        })
        return state.promise
    },
    updateUserCustom ({commit, state}, usercustom, visible) {
        return $axios.post('usercustom', JSON.stringify(usercustom)).then(res => {
            if (res.result) {
                commit('setUserCustom', usercustom)
                alertMsg(visible ? '添加导航成功' : '取消导航成功', 'success')
            } else {
                alertMsg(res.message)
            }
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
