import $http from '@/api'

const getUpdateParams = (payload) => {
    const { updateType, objId, relation, id, value, params, multiple } = payload
    let newRelation = [...relation]
    if (multiple) {
        if (updateType === 'remove') {
            newRelation = newRelation.filter(associatedValue => associatedValue !== value)
        } else {
            newRelation.push(value)
        }
    } else {
        if (updateType === 'remove') {
            newRelation = []
        } else {
            newRelation = [value]
        }
    }
    const updateParams = { ...params }
    updateParams[id] = newRelation.join(',')
    return updateParams
}

const actions = {
    /**
     * 获取实例关联关系
     * @param {String} objId 模型ID
     * @param {String} instId 实例ID
     * @param {Object} config API请求配置
     * @return {Promise} promise 对象
     */
    getInstRelation ({ commit, state, dispatch, rootGetters }, { objId, instId, config }) {
        return $http.post(`inst/association/topo/search/owner/${rootGetters.supplierAccount}/object/${objId}/inst/${instId}`, {}, config)
    },

    updateInstRelation ({commit, state, dispatch, rootGetters}, {params, config}) {
        const updateParams = getUpdateParams(params)
        let promise
        switch (params.objId) {
            case 'host':
                promise = $http.put('hosts/batch', updateParams)
                break
            case 'biz':
                promise = $http.put(`biz/${rootGetters.supplierAccount}/${params['bk_biz_id']}`, updateParams)
                break
            default:
                promise = $http.put(`inst/${rootGetters.supplierAccount}/${params.objId}/${params['bk_inst_id']}`, updateParams)
        }
        return promise
    }
}

export default {
    namespaced: true,
    actions
}
