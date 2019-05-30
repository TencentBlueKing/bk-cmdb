const state = {
    localProcessTemplate: []
}

const getters = {
    localProcessTemplate: state => state.localProcessTemplate
}

const actions = {}

const mutations = {
    hasProcessName (state, process) {
        
    },
    addLocalProcessTemplate (state, process) {
        state.localProcessTemplate.push(process)
    },
    deleteClocalProcessTemplate (state, process) {

    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
