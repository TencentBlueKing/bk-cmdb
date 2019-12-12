const defaultParams = {
    ip: {
        flag: 'bk_host_innerip|bk_host_outer',
        exact: 0,
        data: []
    }
}
const state = {
    filterList: [],
    filterIP: null,
    filterParams: {
        ...defaultParams
    },
    collection: null,
    collectionList: [],
    propertyList: [],
    isHostSearch: null
}

const getters = {
    isCollection: state => !!state.collection,
    configPropertyList: state => {
        const disableList = ['bk_cpu']
        state.propertyList.forEach(property => {
            property.options = property.option
            property['__extra__'] = {
                disabled: disableList.includes(property.bk_property_id)
            }
        })

        return state.propertyList
    }
}

const mutations = {
    setFilterList (state, list) {
        state.filterList = list
    },
    setFilterIP (state, IP) {
        state.filterIP = IP
    },
    setFilterParams (state, params) {
        state.filterParams = params
    },
    setCollectionList (state, list) {
        state.collectionList = list
    },
    setCollection (state, collection) {
        state.collection = collection
    },
    setIsHostSearch (state, boolean) {
        state.isHostSearch = boolean
    },
    addCollection (state, collection) {
        state.collectionList.push(collection)
    },
    updateCollection (state, updatedData) {
        Object.assign(state.collection, updatedData)
    },
    deleteCollection (state, id) {
        state.collectionList = state.collectionList.filter(collection => collection.id !== id)
    },
    clearFilter (state) {
        state.filterList = []
        state.filterIP = null
        state.filterParams = {
            ...defaultParams
        }
        state.collection = null
        state.isHostSearch = false
    },
    setPropertyList (state, list) {
        state.propertyList = list
    }
}

export default {
    namespaced: true,
    state,
    getters,
    mutations
}
