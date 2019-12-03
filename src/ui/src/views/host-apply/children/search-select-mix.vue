<template>
    <bk-search-select
        :data="searchOptions"
        :filter="true"
        :filter-menu-method="filterMethod"
        :filter-children-method="filterMethod"
        :input-type="getInputType()"
        v-model="searchValue"
        placeholder="关键字/字段值"
        @menu-select="handleMenuSelect"
        @child-check="handleChildCheck"
        @change="handleChange">
        <template slot="nextfix">
            <i class="bk-icon icon-close-circle-shape" v-show="searchValue.length" @click.stop="handleClear"></i>
            <i class="bk-icon icon-search" @click.stop="handleSearch"></i>
        </template>
    </bk-search-select>
</template>
<script>
    import TIMEZONE from '@/components/ui/form/timezone.json'
    import Bus from '@/utils/bus'
    const ANY_ID = 'ANY'
    const ANY_TYPE = Symbol('ANY')
    export default {
        components: {

        },
        data () {
            return {
                searchOptions: [],
                searchValue: [],
                properties: [],
                currentMenu: null
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
                        const data = { id: property.bk_property_id, name: property.bk_property_name, type }
                        const any = [{ id: ANY_ID, name: this.$t('任意'), type: ANY_TYPE }]
                        if (type === 'enum') {
                            data.children = any.concat((property.option || []).map(option => ({ id: option.id, name: option.name })))
                            data.multiable = true
                        } else if (type === 'list') {
                            data.children = any.concat((property.option || []).map(option => ({ id: option, name: option })))
                            data.multiable = true
                        } else if (type === 'timezone') {
                            data.children = any.concat(TIMEZONE.map(timezone => ({ id: timezone, name: timezone })))
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
                const keywords = values.filter(value => !value.hasOwnProperty('type'))
                if (keywords.length > 1) {
                    keywords.pop()
                    this.searchValue = values.filter(value => !keywords.includes(value))
                }
            },
            handleClear () {
                this.searchValue = []
                Bus.$emit('topology-search', {})
            },
            handleSearch () {
                Bus.$emit('topology-search', this.getSearchValue())
            },
            getSearchValue () {
                const params = {}
                const keyword = this.searchValue.filter(value => !value.hasOwnProperty('type'))
                if (keyword.length) {
                    params.keyword = keyword[0].name
                }
                return params
            },
            getInputType () {
                const currentMenu = this.currentMenu
                if (currentMenu) {
                    return ['number', 'float'].includes(currentMenu.type) ? 'number' : 'text'
                }
                return 'text'
            },
            filterMethod () {
                return []
            },
            handleMenuSelect (item, index) {
                this.currentMenu = item
            },
            handleChildCheck (item, index) {
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
