<template>
    <div class="form-layout">
        <div class="form-groups" v-if="hasAvaliableGroups" ref="formGroups">
            <template v-for="(group, groupIndex) in $sortedGroups">
                <div class="property-group"
                    :key="groupIndex"
                    v-if="groupedProperties[groupIndex].length">
                    <cmdb-collapse
                        :label="group['bk_group_name']"
                        :collapse.sync="groupState[group['bk_group_id']]">
                        <ul class="property-list">
                            <template v-for="(property, propertyIndex) in groupedProperties[groupIndex]">
                                <li class="property-item"
                                    v-if="!uneditableProperties.includes(property.bk_property_id)"
                                    :key="propertyIndex">
                                    <cmdb-auth tag="div" class="property-name" :title="property['bk_property_name']" v-bind="authProps">
                                        <bk-checkbox class="property-name-checkbox" slot-scope="{ disabled }"
                                            :id="`property-name-${property['bk_property_id']}`"
                                            :disabled="disabled"
                                            v-model="editable[property['bk_property_id']]">
                                            <span class="property-name-text"
                                                :for="`property-name-${property['bk_property_id']}`"
                                                :class="{ required: property['isrequired'] && editable[property.bk_property_id] }">
                                                {{property['bk_property_name']}}
                                            </span>
                                        </bk-checkbox>
                                        <i class="property-name-tooltips icon icon-cc-tips"
                                            v-if="property['placeholder']"
                                            v-bk-tooltips="{
                                                trigger: 'click',
                                                content: htmlEncode(property['placeholder'])
                                            }">
                                        </i>
                                    </cmdb-auth>
                                    <div class="property-value">
                                        <component class="form-component"
                                            :is="`cmdb-form-${property['bk_property_type']}`"
                                            :class="{ error: errors.has(property['bk_property_id']) }"
                                            :unit="property['unit']"
                                            :row="2"
                                            :disabled="!editable[property['bk_property_id']]"
                                            :options="property.option || []"
                                            :data-vv-name="property['bk_property_id']"
                                            :auto-select="false"
                                            :placeholder="getPlaceholder(property)"
                                            v-validate="getValidateRules(property)"
                                            v-model.trim="values[property['bk_property_id']]">
                                        </component>
                                        <span class="form-error"
                                            :title="errors.first(property['bk_property_id'])">
                                            {{errors.first(property['bk_property_id'])}}
                                        </span>
                                    </div>
                                </li>
                            </template>
                        </ul>
                    </cmdb-collapse>
                </div>
            </template>
        </div>
        <div class="form-empty" v-else>
            {{$t('暂无可批量更新的属性')}}
        </div>
        <div class="form-options" :class="{ sticky: scrollbar }">
            <slot name="details-options">
                <cmdb-auth class="inline-block-middle" v-bind="authProps">
                    <bk-button slot-scope="{ disabled }"
                        class="button-save"
                        theme="primary"
                        :disabled="disabled || !hasChange || $loading()"
                        @click="handleSave">
                        {{$t('保存')}}
                    </bk-button>
                </cmdb-auth>
                <bk-button class="button-cancel" @click="handleCancel">{{$t('取消')}}</bk-button>
            </slot>
        </div>
    </div>
</template>

<script>
    import formMixins from '@/mixins/form'
    import RESIZE_EVENTS from '@/utils/resize-events'
    export default {
        name: 'cmdb-form-multiple',
        mixins: [formMixins],
        props: {
            saveAuth: {
                type: [Object, Array],
                default: null
            }
        },
        data () {
            return {
                isMultiple: true,
                values: {},
                refrenceValues: {},
                editable: {},
                scrollbar: false,
                groupState: {
                    'none': true
                }
            }
        },
        computed: {
            changedValues () {
                const changedValues = {}
                for (const propertyId in this.values) {
                    const property = this.getProperty(propertyId)
                    if (
                        ['bool'].includes(property['bk_property_type'])
                        || this.values[propertyId] !== this.refrenceValues[propertyId]
                    ) {
                        changedValues[propertyId] = this.values[propertyId]
                    }
                }
                return changedValues
            },
            hasChange () {
                let hasChange = false
                for (const propertyId in this.editable) {
                    if (this.editable[propertyId]) {
                        hasChange = true
                        break
                    }
                }
                return hasChange
            },
            groupedProperties () {
                return this.$groupedProperties.map(properties => {
                    return properties.filter(property => {
                        const editable = property.editable
                        const isapi = property['bk_isapi']
                        const isonly = property.isonly
                        const isAsst = ['singleasst', 'multiasst'].includes(property['bk_property_type'])
                        return editable && !isapi && !isonly && !isAsst && !this.uneditableProperties.includes(property.bk_property_id)
                    })
                })
            },
            hasAvaliableGroups () {
                return this.groupedProperties.some(properties => !!properties.length)
            },
            authProps () {
                if (this.saveAuth) {
                    return {
                        auth: this.saveAuth
                    }
                }
                return {
                    auth: [],
                    ignore: true
                }
            }
        },
        watch: {
            properties () {
                this.initValues()
                this.initEditableStatus()
            }
        },
        created () {
            this.initValues()
            this.initEditableStatus()
        },
        mounted () {
            if (this.$refs.formGroups) {
                RESIZE_EVENTS.addResizeListener(this.$refs.formGroups, this.checkScrollbar)
            }
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
                this.values = this.$tools.getInstFormValues(this.properties, {}, false)
                this.refrenceValues = this.$tools.clone(this.values)
            },
            initEditableStatus () {
                const editable = {}
                this.groupedProperties.forEach(properties => {
                    properties.forEach(property => {
                        editable[property['bk_property_id']] = false
                    })
                })
                this.editable = editable
            },
            htmlEncode (placeholder) {
                let temp = document.createElement('div')
                temp.innerHTML = placeholder
                const output = temp.innerText
                temp = null
                return output
            },
            getProperty (id) {
                return this.properties.find(property => property['bk_property_id'] === id)
            },
            getPlaceholder (property) {
                const placeholderTxt = ['enum', 'list'].includes(property.bk_property_type) ? '请选择xx' : '请输入xx'
                return this.$t(placeholderTxt, { name: property.bk_property_name })
            },
            getValidateRules (property) {
                if (!this.editable[property.bk_property_id]) {
                    return {}
                }
                return this.$tools.getValidateRules(property)
            },
            getMultipleValues () {
                const multipleValues = {}
                for (const propertyId in this.editable) {
                    if (this.editable[propertyId]) {
                        multipleValues[propertyId] = this.values[propertyId]
                    }
                }
                return this.$tools.formatValues(multipleValues, this.properties)
            },
            handleSave () {
                this.$validator.validateAll().then(result => {
                    if (result) {
                        this.$emit('on-submit', this.getMultipleValues())
                    } else {
                        this.uncollapseGroup()
                    }
                })
            },
            uncollapseGroup () {
                this.errors.items.forEach(item => {
                    const property = this.properties.find(property => property['bk_property_id'] === item.field)
                    const group = property['bk_property_group']
                    this.groupState[group] = false
                })
            },
            handleCancel () {
                this.$emit('on-cancel')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-layout {
        height: 100%;
        @include scrollbar;
    }
    .form-groups {
        padding: 0 0 0 32px;
    }
    .property-group {
        padding: 7px 0 10px 0;
        &:first-child{
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
        .property-item {
            margin: 12px 0 0;
            padding: 0 54px 0 0;
            font-size: 12px;
            flex: 0 0 50%;
            max-width: 50%;
            .property-name {
                display: flex;
                margin: 6px 0 10px;
                color: $cmdbTextColor;
                font-size: 0;
                line-height: 18px;
            }
            .property-name-checkbox {
                margin: 0 6px 0 0;
                max-width: calc(100% - 30px);
                display: flex;

                /deep/ .bk-checkbox-text {
                    width: calc(100% - 30px);
                    flex: 1;
                }
            }
            .property-name-text {
                position: relative;
                display: inline-block;
                max-width: 100%;
                padding: 0 10px 0 0;
                vertical-align: top;
                font-size: 14px;
                @include ellipsis;
                &.required:after{
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
                /deep/ .control-append-group {
                    .bk-input-text {
                        flex: 1;
                    }
                }
            }
        }
    }
    .form-options {
        position: sticky;
        bottom: 0;
        left: 0;
        width: 100%;
        padding: 28px 32px 0;
        &.sticky {
            padding: 10px 32px;
            border-top: 1px solid $cmdbBorderColor;
            background-color: #fff;
        }
        .button-save {
            min-width: 76px;
            margin-right: 4px;
        }
        .button-cancel {
            min-width: 76px;
            background-color: #fff;
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
    .form-empty {
        height: 100%;
        text-align: center;
        &:before{
            content: "";
            display: inline-block;
            vertical-align: middle;
            height: 100%;
            width: 0;
        }
    }
</style>
