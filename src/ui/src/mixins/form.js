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
        objectUnique: {
            type: Array,
            default () {
                return []
            }
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
            const unique = this.objectUnique.find(unique => unique.must_check) || {}
            const uniqueKeys = unique.keys || []
            const sortKey = 'bk_property_index'
            const properties = this.properties.filter(property => {
                return !property['bk_isapi'] &&
                    !uniqueKeys.some(key => key.key_id === property.id)
            })
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
