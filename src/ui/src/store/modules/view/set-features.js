import Vue from 'vue'

const state = {
    syncIdMap: {}
}

const getters = {
    // syncIds: (state) => (templateId) => {
    //     const sessionSyncIdMap = sessionStorage.getItem('setSyncIdMap')
    //     if (!Object.keys(state.syncIdMap).length && sessionSyncIdMap) {
    //         const syncIdMap = JSON.parse(sessionSyncIdMap)
    //         // commit('resetSyncIdMap', syncIdMap)
    //         return syncIdMap[templateId]
    //     }
    //     return state.syncIdMap[templateId]
    // }
}

const mutations = {
    deleteInstancesId (state, data) {
        const curInstancesId = state.syncIdMap[data.id]
        const newInstancesId = curInstancesId.filter(id => id !== data.deleteId)
        Vue.set(state.syncIdMap, data.id, newInstancesId)
        sessionStorage.setItem('setSyncIdMap', JSON.stringify(state.syncIdMap))
    },
    setSyncIdMap (state, data) {
        Vue.set(state.syncIdMap, data.id, data.instancesId)
        sessionStorage.setItem('setSyncIdMap', JSON.stringify(state.syncIdMap))
    },
    resetSyncIdMap (state, data) {
        state.syncIdMap = data
    }
}

export default {
    namespaced: true,
    state,
    getters,
    mutations
}
