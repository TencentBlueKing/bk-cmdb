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
            const groups = [...this.propertyGroups].sort((groupA, groupB) => groupA[sortKey] - groupB[sortKey])
            return groups.concat([{
                'bk_group_id': 'none',
                'bk_group_name': this.$t('Common["更多属性"]')
            }])
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
