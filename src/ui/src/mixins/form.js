export default {
    props: {
        properties: {
            type: Array,
            required: true
        },
        propertyGroups: {
            type: Array,
            required: true
        },
        uneditableProperties: {
            type: Array,
            default () {
                return []
            }
        },
        disabledProperties: {
            type: Array,
            default () {
                return []
            }
        },
        objectUnique: {
            type: Array,
            default () {
                return []
            }
        }
    },
    computed: {
        $sortedGroups () {
            const publicGroups = []
            const metadataGroups = []
            this.propertyGroups.forEach(group => {
                if (this.$tools.getMetadataBiz(group)) {
                    metadataGroups.push(group)
                } else {
                    publicGroups.push(group)
                }
            })
            const sortKey = 'bk_group_index'
            publicGroups.sort((groupA, groupB) => groupA[sortKey] - groupB[sortKey])
            metadataGroups.sort((groupA, groupB) => groupA[sortKey] - groupB[sortKey])
            const allGroups = [
                ...publicGroups,
                ...metadataGroups
            ]
            allGroups.forEach((group, index) => {
                group['bk_group_index'] = index
                this.$set(this.groupState, group['bk_group_id'], group['is_collapse'])
            })
            return allGroups
        },
        $sortedProperties () {
            const unique = this.objectUnique.find(unique => unique.must_check) || {}
            const uniqueKeys = unique.keys || []
            const sortKey = 'bk_property_index'
            const properties = this.properties.filter(property => {
                return !property['bk_isapi']
                    && !uniqueKeys.some(key => key.key_id === property.id)
            })
            return properties.sort((propertyA, propertyB) => propertyA[sortKey] - propertyB[sortKey])
        },
        $groupedProperties () {
            return this.$sortedGroups.map(group => {
                return this.$sortedProperties.filter(property => {
                    const inGroup = (['default', 'none'].includes(property['bk_property_group']) && group['bk_group_id'] === 'default') || property['bk_property_group'] === group['bk_group_id']
                    const isAsst = ['singleasst', 'multiasst'].includes(property['bk_property_type'])
                    return inGroup && !isAsst
                })
            })
        }
    },
    data () {
        return {
            groupState: {}
        }
    }
}
