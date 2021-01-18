<template>
    <div class="form-layout">
        <div class="form-groups" ref="formGroups">
            <template v-for="(group, groupIndex) in $sortedGroups">
                <div class="property-group"
                    :key="groupIndex"
                    v-if="checkGroupAvailable(groupedProperties[groupIndex])">
                    <cmdb-collapse
                        :label="group['bk_group_name']"
                        :collapse.sync="groupState[group['bk_group_id']]">
                        <ul class="property-list">
                            <template v-for="(property, propertyIndex) in groupedProperties[groupIndex]">
                                <li class="property-item"
                                    v-if="checkEditable(property)"
                                    :key="propertyIndex">
                                    <div class="property-name">
                                        <span class="property-name-text" :class="{ required: property['isrequired'] }">{{property['bk_property_name']}}</span>
                                        <i class="property-name-tooltips icon-cc-tips"
                                            v-if="property['placeholder']"
                                            v-bk-tooltips="{
                                                trigger: 'click',
                                                content: htmlEncode(property['placeholder'])
                                            }">
                                        </i>
                                    </div>
                                    <div class="property-value clearfix">
                                        <slot :name="property.bk_property_id">
                                            <component class="form-component"
                                                :is="`cmdb-form-${property['bk_property_type']}`"
                                                :class="{ error: errors.has(property['bk_property_id']) }"
                                                :unit="property['unit']"
                                                :row="2"
                                                :disabled="checkDisabled(property)"
                                                :options="property.option || []"
                                                :data-vv-name="property['bk_property_id']"
                                                :data-vv-as="property['bk_property_name']"
                                                v-validate="getValidateRules(property)"
                                                v-model.trim="values[property['bk_property_id']]">
                                            </component>
                                            <span class="form-error"
                                                :title="errors.first(property['bk_property_id'])">
                                                {{errors.first(property['bk_property_id'])}}
                                            </span>
                                        </slot>
                                    </div>
                                </li>
                            </template>
                        </ul>
                    </cmdb-collapse>
                </div>
            </template>
        </div>
    </div>
</template>

<script>
    import formMixins from '@/mixins/form'
    import RESIZE_EVENTS from '@/utils/resize-events'
    export default {
        mixins: [formMixins],
        props: {
            inst: {
                type: Object,
                default () {
                    return {}
                }
            },
            objId: {
                type: String,
                default: ''
            },
            type: {
                default: 'create',
                validator (val) {
                    return ['create', 'update'].includes(val)
                }
            }
        },
        data () {
            return {
                values: {},
                refrenceValues: {},
                scrollbar: false,
                groupState: {
                    none: true
                }
            }
        },
        computed: {
            groupedProperties () {
                return this.$groupedProperties.map(properties => {
                    return properties.filter(property => !['singleasst', 'multiasst', 'foreignkey'].includes(property['bk_property_type']))
                })
            }
        },
        watch: {
            inst (inst) {
                this.initValues()
            },
            properties () {
                this.initValues()
            }
        },
        created () {
            this.initValues()
        },
        mounted () {
            RESIZE_EVENTS.addResizeListener(this.$refs.formGroups, this.checkScrollbar)
        },
        beforeDestroy () {
            RESIZE_EVENTS.removeResizeListener(this.$refs.formGroups, this.checkScrollbar)
        },
        methods: {
            checkScrollbar () {
                const $layout = this.$el
                this.scrollbar = $layout.scrollHeight !== $layout.offsetHeight
            },
            initValues () {
                this.values = this.$tools.getInstFormValues(this.properties, this.inst, this.type === 'create')
                this.refrenceValues = this.$tools.clone(this.values)
            },
            checkGroupAvailable (properties) {
                const availabelProperties = properties.filter(property => {
                    return this.checkEditable(property)
                })
                return !!availabelProperties.length
            },
            checkEditable (property) {
                if (this.type === 'create') {
                    return !property['bk_isapi']
                }
                return property.editable && !property['bk_isapi']
            },
            checkDisabled (property) {
                if (this.type === 'create') {
                    return false
                }
                return !property.editable || property.isreadonly
            },
            htmlEncode (placeholder) {
                let temp = document.createElement('div')
                temp.innerHTML = placeholder
                const output = temp.innerText
                temp = null
                return output
            },
            getPlaceholder (property) {
                const placeholderTxt = ['enum', 'list'].includes(property.bk_property_type) ? '请选择xx' : '请输入xx'
                return this.$t(placeholderTxt, { name: property.bk_property_name })
            },
            getValidateRules (property) {
                return this.$tools.getValidateRules(property)
            },
            uncollapseGroup () {
                this.errors.items.forEach(item => {
                    const property = this.properties.find(property => property['bk_property_id'] === item.field)
                    const group = property['bk_property_group']
                    this.groupState[group] = false
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-layout {
        height: 100%;
        @include scrollbar-y;
    }
    .form-groups {
        padding: 0 20px 20px;
    }
    .property-group {
        padding: 7px 0 10px 0;
        &:first-child {
            padding: 28px 0 10px 0;
        }
    }
    .group-name {
        font-size: 14px;
        line-height: 14px;
        color: #333948;
        overflow: visible;
    }
    .property-list {
        padding: 4px 0;
        display: flex;
        flex-wrap: wrap;
        justify-content: space-between;
        .property-item {
            flex: 0 0 48%;
            margin: 8px 0 0;
            font-size: 12px;
            .property-name {
                display: block;
                margin: 6px 0 9px;
                color: $cmdbTextColor;
                line-height: 16px;
                font-size: 0;
            }
            .property-name-text {
                position: relative;
                display: inline-block;
                max-width: calc(100% - 20px);
                padding: 0 10px 0 0;
                vertical-align: middle;
                font-size: 12px;
                @include ellipsis;
                &.required:after {
                    position: absolute;
                    left: 100%;
                    top: 0;
                    margin: 0 0 0 -10px;
                    content: "*";
                    color: #ff5656;
                }
            }
            .property-name-tooltips {
                display: inline-block;
                vertical-align: middle;
                width: 16px;
                height: 16px;
                font-size: 16px;
                color: #c3cdd7;
            }
            .property-value {
                font-size: 0;
                position: relative;
                .form-component {
                    font-size: 14px;
                }
            }
        }
    }
    .form-error {
        position: absolute;
        top: 100%;
        left: 0;
        line-height: 14px;
        font-size: 12px;
        color: #ff5656;
        max-width: 100%;
        @include ellipsis;
    }
</style>
