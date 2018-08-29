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
            const sortKey = 'bk_group_index'
            return [...this.propertyGroups].sort((groupA, groupB) => groupA[sortKey] - groupB[sortKey]).concat(['none'])
        },
        $sortedProperties () {
            const sortKey = 'bk_property_index'
            const properties = this.properties.filter(property => !property['bk_isapi'])
            return properties.sort((propertyA, propertyB) => propertyA[sortKey] - propertyB[sortKey])
        },
        $groupedProperties () {
            return this.$sortedGroups.map(group => {
                return this.$sortedProperties.filter(property => property['bk_property_group'] === group['bk_group_id'])
            })
        }
    }
}
