const state = {
    localProcessTemplate: []
}

const getters = {
    localProcessTemplate: state => state.localProcessTemplate,
    hasProcessName: state => process => {
        return state.localProcessTemplate.find(template => template['bk_func_name'] === process['bk_func_name'])
    }
}

const actions = {}

const mutations = {
    addLocalProcessTemplate (state, process) {
        state.localProcessTemplate.push(process)
    },
    deleteLocalProcessTemplate (state, process) {
        state.localProcessTemplate = state.localProcessTemplate.filter(template => template['id'] !== process['id'])
    },
    clearLocalProcessTemplate (state) {
        state.localProcessTemplate = []
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
