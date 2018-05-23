import { $axios, $alertMsg } from '@/api/axios'

const getUpdateParams = ({ rootGetters }, payload) => {
    const { updateType, objId, associated, id, value, params, multiple } = payload
    let newAssociation = [...associated]
    if (multiple) {
        if (updateType === 'remove') {
            newAssociation = newAssociation.filter(associatedValue => associatedValue !== value)
        } else {
            newAssociation.push(value)
        }
    } else {
        if (updateType === 'remove') {
            newAssociation = []
        } else {
            newAssociation = [value]
        }
    }
    const updateParams = { ...params }
    updateParams[id] = newAssociation.join(',')
    return updateParams
}

const state = {}

const getters = {}

const actions = {
    /*
    **   payload:
    **        [updateType]       更新类型 remove | new
    **        [objId]      要更新的实例所属模型Id   父节点的模型id
    **        [associated] 当前实例已经关联的实例ID  children => bk_inst_id
    **        [id]         模型关联字段ID           父节点 property => bk_asst_obj_id
    **        [value]      被关联的实例ID           当前inst_id
    **        [params]     自定义参数           父节点 inst_id
    */
    updateAssociation (context, payload) {
        const params = getUpdateParams(context, payload)
        let promise
        switch (payload.objId) {
            case 'host':
                promise = $axios.put('hosts/batch', params)
                break
            case 'biz':
                promise = $axios.put(`biz/${context.rootGetters.bkSupplierAccount}/${payload['bk_biz_id']}`, params)
                break
            default:
                promise = $axios.put(`inst/${context.rootGetters.bkSupplierAccount}/${payload.objId}/${payload['bk_inst_id']}`, params)
        }
        return promise
    }
}

const mutations = {}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
