export default {
    props: {
        properties: {
            type: Array,
            required: true
        },
        propertyGroups: {
            type: Array,
            required: true
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
                ...metadataGroups,
                {
                    'bk_group_id': 'none',
                    'bk_group_name': this.$t('Common["更多属性"]')
                }
            ]
            allGroups.forEach((group, index) => {
                group['bk_group_index'] = index
            })
            return allGroups
        },
        $sortedProperties () {
            const sortKey = 'bk_property_index'
            const properties = this.properties.filter(property => !property['bk_isapi'])
            return properties.sort((propertyA, propertyB) => propertyA[sortKey] - propertyB[sortKey])
        },
        $groupedProperties () {
            return this.$sortedGroups.map(group => {
                return this.$sortedProperties.filter(property => {
                    const inGroup = property['bk_property_group'] === group['bk_group_id']
                    const isAsst = ['singleasst', 'multiasst'].includes(property['bk_property_type'])
                    return inGroup && !isAsst
                })
            })
        }
    }
}
