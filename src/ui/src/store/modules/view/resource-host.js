import Vue from 'vue'

const state = {
    activeDirectory: null,
    dirList: []
}

const getters = {
    activeDirectory: state => state.activeDirectory,
    dirList: state => state.dirList
}

const mutations = {
    setActiveDirectory (state, active) {
        state.activeDirectory = active
    },
    setDirList (state, list) {
        state.dirList = list
    },
    setHostCount (state, params = {}) {
        const index = state.dirList.findIndex(dir => dir.bk_inst_id === params.id)
        if (index > -1) {
            const curDir = state.dirList[index]
            const curCount = curDir.host_count
            curDir.host_count = params.type === 'add' ? curCount + params.count : curCount - params.count
            Vue.$set(state.dirList, index, curDir)
        }
    }
}

export default {
    namespaced: true,
    state,
    getters,
    mutations
}
