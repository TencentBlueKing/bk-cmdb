<template>
    <bk-search-select
        :data="searchOptions"
        :filter="true"
        :filter-menu-method="filterMethod"
        :filter-children-method="filterMethod"
        v-model="searchValue"
        placeholder="名称关键字 / 或字段：字段值">
        <template slot="nextfix">
            <i class="bk-icon icon-close-circle-shape" v-show="searchValue.length" @click="handleClear"></i>
            <i class="bk-icon icon-search" @click="handleSearch"></i>
        </template>
    </bk-search-select>
</template>
<script>
    import TIMEZONE from '@/components/ui/form/timezone.json'
    export default {
        components: {

        },
        data () {
            return {
                searchOptions: [],
                searchValue: [],
                properties: [],
                typeAny: Symbol('any')
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
                        const any = [{ id: 'any', name: this.$t('任意'), type: this.typeAny }]
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
            handleClear () {
                this.searchValue = []
                this.$emit('clear')
            },
            handleSearch () {
                this.$emit('search', this.getSearchValue())
            },
            getSearchValue () {
                return {}
            },
            filterMethod () {
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
