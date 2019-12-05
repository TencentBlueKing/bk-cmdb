<template>
    <bk-search-select
        :data="searchOptions"
        :filter="true"
        :filter-menu-method="filterMenuMethod"
        :filter-children-method="filterChildrenMethod"
        :condition="ORCondition"
        v-model="searchValue"
        placeholder="关键字/字段值"
        @change="handleChange"
        @menu-select="handleMenuSelect">
        <template slot="nextfix">
            <i class="bk-icon icon-close-circle-shape" v-show="searchValue.length" @click.stop="handleClear"></i>
            <i class="bk-icon icon-search" @click.stop="handleSearch"></i>
        </template>
    </bk-search-select>
</template>
<script>
    import TIMEZONE from '@/components/ui/form/timezone.json'
    import Bus from '@/utils/bus'
    export default {
        components: {

        },
        data () {
            return {
                searchOptions: [],
                searchValue: [],
                properties: [],
                currentMenu: null,
                ORCondition: {
                    name: this.$i18n.locale === 'en' ? 'OR' : '或'
                }
            }
        },
        created () {
            this.initOptions()
        },
        methods: {
            async initOptions () {
                try {
                    const properties = await this.$store.dispatch('hostApply/getProperties')
                    const unsupportType = ['date', 'time', 'objuser']
                    const availableProperties = properties.filter(property => property.host_apply_enabled && !unsupportType.includes(property.bk_property_type))
                    this.searchOptions = availableProperties.map(property => {
                        const type = property.bk_property_type
                        const data = { id: property.id, name: property.bk_property_name, type }
                        if (type === 'enum') {
                            data.children = (property.option || []).map(option => ({ id: option.id, name: option.name, disabled: false }))
                            data.multiable = true
                        } else if (type === 'list') {
                            data.children = (property.option || []).map(option => ({ id: option, name: option, disabled: false }))
                            data.multiable = true
                        } else if (type === 'timezone') {
                            data.children = TIMEZONE.map(timezone => ({ id: timezone, name: timezone, disabled: false }))
                            data.multiable = true
                        } else if (type === 'bool') {
                            data.children = [{ id: true, name: 'true' }, { id: false, name: 'false' }]
                        } else {
                            data.children = []
                        }
                        return data
                    })
                } catch (e) {
                    console.error(e)
                }
            },
            handleChange (values) {
                const keywords = values.filter(value => !value.hasOwnProperty('type') && value.hasOwnProperty('id'))
                if (keywords.length > 1) {
                    keywords.pop()
                    this.searchValue = values.filter(value => !keywords.includes(value))
                }
                this.currentMenu = null
            },
            handleClear () {
                this.searchValue = []
                Bus.$emit('topology-search', { query_filter: { rules: [] } })
            },
            handleSearch () {
                Bus.$emit('topology-search', this.getSearchValue())
            },
            getSearchValue () {
                const params = {
                    query_filter: {
                        condition: 'OR',
                        rules: []
                    }
                }
                const filterGroup = []
                const lastGroup = this.searchValue.reduce((group, value) => {
                    if (!value.hasOwnProperty('id') && group.length) {
                        filterGroup.push(group)
                        return []
                    }
                    return [...group, value]
                }, [])
                filterGroup.push(lastGroup)

                if (filterGroup.length) {
                    filterGroup.forEach(group => {
                        const rule = {
                            condition: 'AND',
                            rules: []
                        }
                        params.query_filter.rules.push(rule)
                        group.forEach(item => {
                            if (item.hasOwnProperty('type')) {
                                if (item.values.length === 1) {
                                    rule.rules.push({
                                        field: String(item.id),
                                        operator: 'contains',
                                        value: item.values[0].id === '*' ? '' : item.values[0].id
                                    })
                                } else {
                                    const subRules = {
                                        condition: 'OR',
                                        rules: []
                                    }
                                    rule.rules.push(subRules)
                                    item.values.forEach(value => {
                                        subRules.rules.push({
                                            field: String(item.id),
                                            operator: 'contains',
                                            value: value.id === '*' ? '' : value.id
                                        })
                                    })
                                }
                            } else {
                                rule.rules.push({
                                    field: 'keyword',
                                    operator: 'contains',
                                    value: item.id
                                })
                            }
                        })
                    })
                }
                return params
            },
            handleMenuSelect (item, index) {
                this.currentMenu = item
            },
            filterMenuMethod (list, filter) {
                return list.filter(item => item.name.toLowerCase().indexOf(filter.toLowerCase()) > -1)
            },
            filterChildrenMethod (list, filter) {
                if (this.currentMenu && this.currentMenu.children && this.currentMenu.children.length) {
                    return this.currentMenu.children.filter(item => item.name.toLowerCase().indexOf(filter.toLowerCase()) > -1)
                }
                return []
            }
        }
    }
</script>
<style lang="scss" scoped>
    .icon-close-circle-shape {
        font-size: 14px;
        margin-right: 6px;
        cursor: pointer;
    }
    .icon-search {
        font-size: 16px;
        margin-right: 10px;
        cursor: pointer;
    }
</style>
