<template>
    <div class="quick-search-layout">
        <div class="quick-search-toggle" ref="quickSearchToggle" @click="handleToggle">
            {{$t('HostResourcePool["筛选"]')}}
            <i class="bk-icon icon-angle-up"></i>
        </div>
        <cmdb-selector class="filter-selector"
            setting-key="bk_property_id"
            display-key="bk_property_name"
            :list="properties"
            v-model="propertyId"
            @on-selected="handleSelect">
        </cmdb-selector>
        <div v-if="property" class="filter-value-container">
            <cmdb-form-enum class="filter-value fl"
                v-if="property['bk_property_type'] === 'enum'"
                :options="property.option || []"
                :allow-clear="true"
                v-model="value"
                @on-selected="handleSearch">
            </cmdb-form-enum>
            <cmdb-form-int class="filter-value fl"
                v-else-if="property['bk_property_type'] === 'int'"
                v-model="value"
                @keydown.enter="handleSearch">
            </cmdb-form-int>
            <comonent class="filter-value"
                :is="`cmdb-form-${property['bk_property_type']}`"
                v-model.trim="value"
                @keydown.native.enter="handleSearch">
            </comonent>
            <i class="filter-search bk-icon icon-search"
                v-show="!['enum', 'objuser'].includes(property['bk_property_type'])"
                @click="handleSearch"></i>
        </div>
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
        watch: {
            properties () {
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
                if (this.properties.length) {
                    this.selected = this.properties[0]['bk_property_id']
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
                this.$emit('on-search', this.property, this.value)
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
                width: 260px;
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