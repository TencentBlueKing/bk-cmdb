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
            const bizCustomGroups = []
            this.propertyGroups.forEach(group => {
                if (group.hasOwnProperty('bk_biz_id') && group.bk_biz_id > 0) {
                    bizCustomGroups.push(group)
                } else {
                    publicGroups.push(group)
                }
            })
            const sortKey = 'bk_group_index'
            publicGroups.sort((groupA, groupB) => groupA[sortKey] - groupB[sortKey])
            bizCustomGroups.sort((groupA, groupB) => groupA[sortKey] - groupB[sortKey])
            const allGroups = [
                ...publicGroups,
                ...bizCustomGroups
            ]
            allGroups.forEach((group, index) => {
                group.bk_group_index = index
                this.$set(this.groupState, group.bk_group_id, group.is_collapse)
            })
            return allGroups
        },
        $sortedProperties () {
            const unique = this.isMultiple ? this.objectUnique.find(unique => unique.must_check) || {} : {}
            const uniqueKeys = unique.keys || []
            const sortKey = 'bk_property_index'
            const properties = this.properties.filter(property => {
                return !property.bk_isapi
                    && !uniqueKeys.some(key => key.key_id === property.id)
            })
            return properties.sort((propertyA, propertyB) => propertyA[sortKey] - propertyB[sortKey])
        },
        $groupedProperties () {
            return this.$sortedGroups.map(group => {
                return this.$sortedProperties.filter(property => {
                    // 兼容旧数据， 把none 这个分组的属性塞到默认分组去
                    const isNoneGroup = property.bk_property_group === 'none'
                    if (isNoneGroup) {
                        return group.bk_group_id === 'default'
                    }
                    return property.bk_property_group === group.bk_group_id
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
