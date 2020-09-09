const state = {
    activeDirectory: null,
    directoryList: []
}

const getters = {
    activeDirectory: state => state.activeDirectory,
    directoryList: state => state.directoryList,
    defaultDirectory: state => state.directoryList.find(directory => directory.default === 1)
}

const mutations = {
    setActiveDirectory (state, active) {
        state.activeDirectory = active
    },
    setDirectoryList (state, list) {
        state.directoryList = list
    },
    addDirectory (state, directory) {
        state.directoryList.splice(1, 0, directory)
    },
    updateDirectory (state, directory) {
        const index = state.directoryList.findIndex(data => data.bk_module_id === directory.bk_module_id)
        if (index > -1) {
            state.directoryList.splice(index, 1, directory)
        }
    },
    deleteDirectory (state, id) {
        const index = state.directoryList.findIndex(target => target.bk_module_id === id)
        if (index > -1) {
            state.directoryList.splice(index, 1)
        }
    },
    refreshDirectoryCount (state, newList = []) {
        state.directoryList.forEach(directory => {
            const newDirectory = newList.find(newDirectory => newDirectory.bk_module_id === directory.bk_module_id)
            Object.assign(directory, newDirectory)
        })
    }
}

export default {
    namespaced: true,
    state,
    getters,
    mutations
}
