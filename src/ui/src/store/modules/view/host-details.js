import { getMetadataBiz } from '@/utils/tools'
const state = {
    info: {},
    properties: [],
    propertyGroups: []
}

const getters = {
    groupedProperties: state => {
        const groupedProperties = []
        state.propertyGroups.forEach(group => {
            const properties = state.properties.filter(property => property.bk_property_group === group.bk_group_id)
            if (properties.length) {
                groupedProperties.push({
                    ...group,
                    properties
                })
            }
        })
        return groupedProperties.sort((prev, next) => {
            const prevMetadata = !!getMetadataBiz(prev)
            const nextMetadata = !!getMetadataBiz(next)
            if (
                (prevMetadata && nextMetadata)
                || (!prevMetadata && !nextMetadata)
            ) {
                return prev.bk_group_index - next.bk_group_index
            } else if (prevMetadata) {
                return 1
            } else {
                return -1
            }
        })
    }
}

const mutations = {
    setHostInfo (state, info) {
        state.info = info || {}
    },
    setHostProperties (state, properties) {
        state.properties = properties
    },
    setHostPropertyGroups (state, propertyGroups) {
        state.propertyGroups = propertyGroups
    },
    updateInfo (state, data) {
        Object.assign(state.info.host, data)
    }
}

export default {
    namespaced: true,
    state,
    getters,
    mutations
}
