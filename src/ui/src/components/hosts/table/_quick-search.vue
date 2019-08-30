<template>
    <div class="quick-search-layout">
        <div class="quick-search-toggle" ref="quickSearchToggle" @click="handleToggle">
            {{$t('筛选')}}
            <i class="bk-icon icon-angle-up"></i>
        </div>
        <cmdb-selector class="filter-selector"
            setting-key="bk_property_id"
            display-key="bk_property_name"
            :list="filteredProperties"
            v-model="propertyId"
            @on-selected="handleSelect">
        </cmdb-selector>
        <div v-if="property" class="filter-value-container">
            <cmdb-form-enum class="filter-value fl"
                v-if="property['bk_property_type'] === 'enum'"
                :options="property.option || []"
                :allow-clear="true"
                v-model.trim="value">
            </cmdb-form-enum>
            <cmdb-form-int class="filter-value fl"
                v-else-if="property['bk_property_type'] === 'int'"
                v-model="value">
            </cmdb-form-int>
            <cmdb-form-bool-input class="filter-value fl"
                v-else-if="property['bk_property_type'] === 'bool'"
                v-model.trim="value">
            </cmdb-form-bool-input>
            <cmdb-form-date-range class="filter-value fl"
                v-else-if="['date', 'time'].includes(property['bk_property_type'])"
                :timer="property['bk_property_type'] === 'time'"
                v-model="value">
            </cmdb-form-date-range>
            <comonent class="filter-value"
                v-else
                :is="`cmdb-form-${property['bk_property_type']}`"
                v-model.trim="value">
            </comonent>
        </div>
        <bk-button theme="primary"
            :loading="$loading()"
            @click="handleSearch">
            {{$t('搜索')}}
        </bk-button>
    </div>
</template>

<script>
    export default {
        props: {
            properties: {
                type: Array,
                required: true
            }
        },
        data () {
            return {
                propertyId: '',
                value: '',
                property: null
            }
        },
        computed: {
            filteredProperties () {
                return this.properties.filter(property => !['singleasst', 'multiasst', 'foreignkey'].includes(property['bk_property_type']))
            },
            type () {
                return this.property ? this.property['bk_property_type'] : ''
            },
            searchValue () {
                if (['objuser'].includes(this.type)) {
                    return this.value.split(',')
                }
                return this.value
            },
            operator () {
                const map = {
                    'singlechar': '$regex',
                    'int': '$eq',
                    'float': '$eq',
                    'enum': '$eq',
                    'date': '$in',
                    'time': '$in',
                    'longchar': '$regex',
                    'objuser': '$in',
                    'timezone': '$eq',
                    'bool': '$eq'
                }
                return map[this.type] || '$eq'
            }
        },
        watch: {
            filteredProperties () {
                this.setDefaultSelected()
            }
        },
        created () {
            this.setDefaultSelected()
        },
        mounted () {
            this.setTogglePosition()
        },
        methods: {
            setDefaultSelected () {
                if (this.filteredProperties.length) {
                    this.selected = this.filteredProperties[0]['bk_property_id']
                } else {
                    this.selected = null
                }
                this.value = ''
            },
            setTogglePosition () {
                const $quickSearchToggle = this.$refs.quickSearchToggle
                const $target = this.$parent.$refs.quickSearchButton.$el
                $quickSearchToggle.style.width = $target.offsetWidth + 'px'
                $quickSearchToggle.style.left = $target.offsetLeft - 1 + 'px'
            },
            handleToggle () {
                this.$emit('on-search', null, '')
                this.$emit('on-toggle')
            },
            handleSelect (propertyId, property) {
                this.property = property
                this.value = ''
            },
            handleSearch () {
                if (String(this.searchValue).length) {
                    this.$emit('on-search', this.property, this.searchValue, this.operator)
                } else {
                    this.$emit('on-search', null, '', '')
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .quick-search-layout {
        position: relative;
        margin: 9px 0 -10px 0;
        padding: 11px;
        background-color: #fafbfd;
        font-size: 0;
        border: 1px solid #dde4eb;
        .quick-search-toggle {
            position: absolute;
            height: 46px;
            line-height: 34px;
            font-size: 14px;
            bottom: 100%;
            background-color: #fafbfd;
            text-align: center;
            border: 1px solid #dde4eb;
            border-bottom: none;
            cursor: pointer;
            z-index: 100;
            .icon-angle-up {
                font-size: 12px;
                top: 0;
            }
        }
        .filter-selector {
            width: 145px;
        }
        .filter-value-container {
            display: inline-block;
            vertical-align: middle;
            position: relative;
            .filter-value {
                min-width: 260px;
                margin: 0 0 0 -1px;
            }
            .filter-search {
                position: absolute;
                right: 9px;
                top: 9px;
                font-size: 18px;
                cursor: pointer;
            }
        }
    }
</style>
