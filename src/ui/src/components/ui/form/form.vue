<template>
    <div class="form-layout">
        <slot name="prepend"></slot>
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
                                <li :class="['property-item', { flex: flexProperties.includes(property['bk_property_id']) }]"
                                    v-if="checkEditable(property)"
                                    :key="propertyIndex">
                                    <div class="property-name" v-if="!invisibleNameProperties.includes(property['bk_property_id'])">
                                        <span class="property-name-text" :class="{ required: isRequired(property) }">{{property['bk_property_name']}}</span>
                                        <i class="property-name-tooltips icon-cc-tips"
                                            v-if="property['placeholder']"
                                            v-bk-tooltips="{
                                                trigger: 'click',
                                                content: htmlEncode(property['placeholder'])
                                            }">
                                        </i>
                                        <form-tips :type="type" :property="property" :render="renderTips"></form-tips>
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
                                                :placeholder="getPlaceholder(property)"
                                                :auto-select="false"
                                                v-bind="$attrs"
                                                v-validate="getValidateRules(property)"
                                                v-model.trim="values[property['bk_property_id']]">
                                            </component>
                                            <form-append :type="type" :property="property" :render="renderAppend"></form-append>
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
        <div class="form-options"
            v-if="showOptions"
            :class="{ sticky: scrollbar }">
            <slot name="form-options">
                <cmdb-auth class="inline-block-middle" :auth="saveAuth">
                    <bk-button slot-scope="{ disabled }"
                        class="button-save"
                        theme="primary"
                        :loading="$loading()"
                        :disabled="disabled || (type !== 'create' && !hasChange)"
                        @click="handleSave">
                        {{type === 'create' ? $t('提交') : $t('保存')}}
                    </bk-button>
                </cmdb-auth>
                <bk-button class="button-cancel" @click="handleCancel">{{$t('取消')}}</bk-button>
            </slot>
            <slot name="extra-options"></slot>
        </div>
    </div>
</template>

<script>
    import formMixins from '@/mixins/form'
    import RESIZE_EVENTS from '@/utils/resize-events'
    import FormTips from './form-tips.js'
    import FormAppend from './form-append.js'
    export default {
        name: 'cmdb-form',
        components: {
            FormTips,
            FormAppend
        },
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
            },
            showOptions: {
                type: Boolean,
                default: true
            },
            saveAuth: {
                type: Object,
                default: null
            },
            renderTips: Function,
            renderAppend: Function,
            flexProperties: {
                type: Array,
                default: () => []
            },
            invisibleNameProperties: {
                type: Array,
                default: () => []
            },
            customValidator: Function
        },
        data () {
            return {
                values: {},
                refrenceValues: {},
                scrollbar: false
            }
        },
        computed: {
            changedValues () {
                const changedValues = {}
                for (const propertyId in this.values) {
                    if (this.values[propertyId] !== this.refrenceValues[propertyId]) {
                        changedValues[propertyId] = this.values[propertyId]
                    }
                }
                return changedValues
            },
            hasChange () {
                return !!Object.keys(this.changedValues).length
            },
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
                return property.editable && !property['bk_isapi'] && !this.uneditableProperties.includes(property.bk_property_id)
            },
            checkDisabled (property) {
                if (this.type === 'create') {
                    return false
                }
                return !property.editable || property.isreadonly || this.disabledProperties.includes(property.bk_property_id)
            },
            isRequired (property) {
                return property.isrequired
                // 后台无对应逻辑，前端屏蔽唯一校验配置中为空必须校验的字段设置为必填的逻辑
                // if (property.isrequired) {
                //     return true
                // }
                // const unique = this.objectUnique.find(unique => unique.must_check)
                // if (unique) {
                //     return unique.keys.some(key => key.key_id === property.id)
                // }
                // return false
            },
            htmlEncode (placeholder) {
                let temp = document.createElement('div')
                temp.innerHTML = placeholder
                const output = temp.innerText
                temp = null
                return output
            },
            getPlaceholder (property) {
                const placeholderTxt = ['enum', 'list', 'organization'].includes(property.bk_property_type) ? '请选择xx' : '请输入xx'
                return this.$t(placeholderTxt, { name: property.bk_property_name })
            },
            getValidateRules (property) {
                const rules = this.$tools.getValidateRules(property)
                if (this.isRequired(property)) {
                    rules.required = true
                }

                if (['bk_set_name', 'bk_module_name', 'bk_inst_name'].includes(property.bk_property_id)) {
                    rules.businessTopoInstNames = true
                }

                return rules
            },
            async handleSave () {
                const validatePromise = [this.$validator.validateAll()]
                if (typeof this.customValidator === 'function') {
                    validatePromise.push(this.customValidator())
                }
                const results = await Promise.all(validatePromise)
                const isValid = results.every(result => result)
                if (isValid) {
                    this.$emit(
                        'on-submit',
                        this.$tools.formatValues(this.values, this.properties),
                        this.$tools.formatValues(this.changedValues, this.properties),
                        this.inst,
                        this.type
                    )
                } else {
                    this.uncollapseGroup()
                }
            },
            uncollapseGroup () {
                this.errors.items.forEach(item => {
                    const compareKey = item.scope || item.field
                    const property = this.properties.find(property => property['bk_property_id'] === compareKey)
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
        @include scrollbar-y;
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
            width: 50%;
            margin: 12px 0 0;
            padding: 0 54px 0 0;
            font-size: 12px;
            flex: 0 0 50%;
            max-width: 50%;
            // flex: 0 1 auto;
            .property-name {
                display: block;
                margin: 2px 0 6px;
                color: $cmdbTextColor;
                line-height: 24px;
                font-size: 0;
            }
            .property-name-text {
                position: relative;
                display: inline-block;
                max-width: calc(100% - 20px);
                padding: 0 10px 0 0;
                vertical-align: middle;
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
                margin-right: 6px;
                color: #c3cdd7;
            }
            .property-value {
                font-size: 0;
                position: relative;
                display: flex;
                /deep/ .control-append-group {
                    .bk-input-text {
                        flex: 1;
                    }
                }
            }
            .form-component {
                flex: 1;
            }

            &.flex {
                flex: 1;
                padding-right: 54px;
                width: 100%;
                max-width: unset;
            }
        }
    }
    .form-options {
        position: sticky;
        bottom: 0;
        left: 0;
        width: 100%;
        padding: 10px 32px 0;
        font-size: 0;
        &.sticky {
            padding: 10px 32px;
            border-top: 1px solid $cmdbBorderColor;
            background-color: #fff;
            z-index: 100;
        }
        .button-save {
            min-width: 76px;
            margin-right: 4px;
        }
        .button-cancel {
            min-width: 76px;
            margin: 0 4px;
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
</style>
