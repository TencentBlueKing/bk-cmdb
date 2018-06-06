import {$axios, $alertMsg} from '@/api/axios'

const state = {
    attribute: {},
    topo: []
}

const getters = {
    attribute: state => state.attribute,
    topo: state => state.topo
}

const actions = {
    getAttribute ({commit, state, rootGetters}, {objId, force}) {
        if (!force && state.attribute.hasOwnProperty(objId)) {
            return Promise.resolve({result: true, data: state.attribute[objId]})
        }
        return $axios.post('object/attr/search', {bk_obj_id: objId, bk_supplier_account: rootGetters.bkSupplierAccount}).then(res => {
            if (res.result) {
                let attribute = {}
                attribute[objId] = res.data
                commit('setAttribute', attribute)
            } else {
                $alertMsg(res['bk_error_msg'])
            }
            return res
        })
    },
    getTopo ({commit, state, rootGetters}, force = false) {
        if (!force && state.topo.length) {
            return Promise.resolve({result: true, data: state.topo})
        }
        return $axios.get(`topo/model/${rootGetters.bkSupplierAccount}`).then(res => {
            if (res.result) {
                commit('setTopo', res.data)
            } else {
                $alertMsg(res['bk_error_msg'])
            }
            return res
        })
    }
}

const mutations = {
    setAttribute (state, attribute) {
        state.attribute = Object.assign({}, state.attribute, attribute)
    },
    setTopo (state, topo) {
        state.topo = topo
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
