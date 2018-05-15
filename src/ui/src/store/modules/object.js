import {$axios, $alertMsg} from '@/api/axios'

const state = {
    attribute: {}
}

const getters = {
    attribute: state => state.attribute
}

const actions = {
    getAttribute ({commit, state, rootGetters}, objId, force = false) {
        if (!force && state.attribute.hasOwnProperty(objId)) {
            return Promise.resolve({result: true, data: state.attribute[objId]})
        }
        return $axios.post('object/attr/search', {bk_obj_id: objId, bk_supplier_account: rootGetters.bkSupplierAccount}).then(res => {
            if (res.result) {
                let attribute = {}
                attribute[objId] = res.data
                commit('setAttribute', attribute)
                attribute[objId] = res.data
            } else {
                $alertMsg(res['bk_error_msg'])
            }
        })
    }
}

const mutations = {
    setAttribute (state, attribute) {
        state.attribute = Object.assign({}, state.attribute, attribute)
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
