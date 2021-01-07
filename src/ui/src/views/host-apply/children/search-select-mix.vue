<template>
    <bk-search-select
        :data="searchOptions"
        :filter="true"
        :filter-menu-method="filterMenuMethod"
        :filter-children-method="filterChildrenMethod"
        :show-condition="false"
        :show-popover-tag-change="false"
        :strink="false"
        v-model="searchValue"
        :placeholder="$t('关键字/字段值')"
        @change="handleChange"
        @menu-select="handleMenuSelect"
        @key-enter="handleKeyEnter"
        @input-focus="handleFocus"
        @input-click-outside="handleBlur">
        <template slot="nextfix">
            <i class="bk-icon icon-close-circle-shape" v-show="showClear && searchValue.length" @click.stop="handleClear"></i>
        </template>
    </bk-search-select>
</template>
<script>
    import TIMEZONE from '@/components/ui/form/timezone.json'
    import Bus from '@/utils/bus'
    import { mapGetters } from 'vuex'
    export default {
        components: {},
        data () {
            return {
                showClear: false,
                searchOptions: [],
                fullOptions: [],
                searchValue: [],
                properties: [],
                currentMenu: null
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId'])
        },
        watch: {
            searchValue (searchValue) {
                this.searchOptions.forEach(option => {
                    const selected = searchValue.some(value => value.id === option.id && value.name === option.name && value.type === option.type)
                    option.disabled = selected
                })
                this.handleSearch()
            }
        },
        created () {
            this.initOptions()
        },
        methods: {
            async initOptions () {
                try {
                    const properties = await this.$store.dispatch('hostApply/getProperties', { params: { bk_biz_id: this.bizId } })
                    const availableProperties = properties.filter(property => property.host_apply_enabled)
                    this.searchOptions = availableProperties.map(property => {
                        const type = property.bk_property_type
                        const data = { id: property.id, name: property.bk_property_name, type, disabled: false }
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
                    this.fullOptions = this.searchOptions.slice(0)
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
            },
            handleKeyEnter () {
                this.currentMenu = null
            },
            handleFocus () {
                this.showClear = true
            },
            handleBlur () {
                this.showClear = false
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
                        condition: 'AND',
                        rules: []
                    }
                }
                const rules = params.query_filter.rules
                this.searchValue.forEach(item => {
                    if (item.hasOwnProperty('type')) {
                        if (item.values.length === 1) {
                            const value = item.values[0]
                            const isAny = value.id === '*'
                            const rule = { field: String(item.id) }
                            if (isAny) {
                                rule.operator = 'exist'
                            } else {
                                rule.operator = 'contains'
                                rule.value = String(value.id).trim()
                            }
                            rules.push(rule)
                        } else {
                            const subRule = {
                                condition: 'OR',
                                rules: []
                            }
                            item.values.forEach(value => {
                                subRule.rules.push({
                                    field: String(item.id),
                                    operator: 'contains',
                                    value: String(value.id).trim()
                                })
                            })
                            rules.push(subRule)
                        }
                    } else {
                        rules.push({
                            field: 'keyword',
                            operator: 'contains',
                            value: String(item.id).trim()
                        })
                    }
                })
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
