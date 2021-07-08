<template>
    <cmdb-sticky-layout class="form-layout">
        <div class="form-groups" ref="formGroups">
            <template v-for="(group, groupIndex) in $sortedGroups">
                <div class="property-group"
                    :key="groupIndex"
                    v-if="checkGroupAvailable(groupedProperties[groupIndex])">
                    <cmdb-collapse
                        :label="group['bk_group_name']"
                        :collapse.sync="groupState[group['bk_group_id']]">
                        <ul class="property-list clearfix">
                            <li class="property-item fl"
                                v-for="(property, propertyIndex) in groupedProperties[groupIndex]"
                                v-if="checkEditable(property)"
                                :key="propertyIndex">
                                <div class="property-name">
                                    <span class="property-name-text" :class="{ required: property['isrequired'] }">{{property['bk_property_name']}}</span>
                                    <i class="property-name-tooltips icon-cc-tips"
                                        v-if="property['placeholder']"
                                        v-bk-tooltips="htmlEncode(property['placeholder'])">
                                    </i>
                                </div>
                                <div class="property-value clearfix">
                                    <slot :name="property.bk_property_id">
                                        <component class="form-component"
                                            :is="`cmdb-form-${property['bk_property_type']}`"
                                            :class="{ error: errors.has(property['bk_property_id']) }"
                                            :unit="property['unit']"
                                            :disabled="checkDisabled(property)"
                                            :options="property.option || []"
                                            :data-vv-name="property['bk_property_id']"
                                            :data-vv-as="property['bk_property_name']"
                                            :placeholder="$t('请输入xx', { name: property.bk_property_name })"
                                            v-bind="$tools.getValidateEvents(property)"
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
                        </ul>
                    </cmdb-collapse>
                </div>
            </template>
        </div>
        <div class="form-options"
            v-if="showOptions"
            slot="footer"
            slot-scope="{ sticky }"
            :class="{ sticky: sticky }">
            <slot name="form-options">
                <span class="inline-block-middle"
                    v-cursor="{
                        active: !$isAuthorized(saveAuth),
                        auth: [saveAuth]
                    }">
                    <bk-button class="button-save" theme="primary"
                        :disabled="!$isAuthorized(saveAuth) || !hasChange || $loading()"
                        @click="handleSave">
                        {{$t('保存')}}
                    </bk-button>
                </span>
                <bk-button class="button-cancel" @click="handleCancel">{{$t('取消')}}</bk-button>
            </slot>
            <slot name="extra-options"></slot>
        </div>
    </cmdb-sticky-layout>
</template>

<script>
    import formMixins from '@/mixins/form'
    export default {
        name: 'cmdb-form',
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
                type: [String, Array],
                default: ''
            }
        },
        data () {
            return {
                values: {},
                refrenceValues: {}
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
        methods: {
            initValues () {
                this.values = this.$tools.getInstFormValues(this.properties, this.inst)
                const timer = setTimeout(() => {
                    this.refrenceValues = this.$tools.clone(this.values)
                    clearTimeout(timer)
                })
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
            htmlEncode (placeholder) {
                let temp = document.createElement('div')
                temp.innerHTML = placeholder
                const output = temp.innerText
                temp = null
                return output
            },
            getValidateRules (property) {
                const rules = {}
                const {
                    bk_property_type: propertyType,
                    option,
                    isrequired
                } = property
                if (isrequired) {
                    rules.required = true
                }
                if (option) {
                    if (propertyType === 'int') {
                        if (option.hasOwnProperty('min') && !['', null, undefined].includes(option.min)) {
                            rules['min_value'] = option.min
                        }
                        if (option.hasOwnProperty('max') && !['', null, undefined].includes(option.max)) {
                            rules['max_value'] = option.max
                        }
                    } else if (['singlechar', 'longchar'].includes(propertyType)) {
                        rules['remoteString'] = option
                    }
                }
                if (['singlechar', 'longchar'].includes(propertyType)) {
                    rules[propertyType] = true
                    rules.length = propertyType === 'singlechar' ? 256 : 2000
                }
                if (propertyType === 'float') {
                    rules['float'] = true
                }
                return rules
            },
            handleSave () {
                this.$validator.validateAll().then(result => {
                    if (result) {
                        this.$emit('on-submit', { ...this.values }, { ...this.changedValues }, this.inst, this.type)
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
    .form-layout{
        height: 100%;
        @include scrollbar-y;
    }
    .form-groups{
        padding: 0 0 0 32px;
    }
    .property-group{
        padding: 7px 0 10px 0;
        &:first-child{
            padding: 28px 0 10px 0;
        }
    }
    .group-name{
        font-size: 14px;
        line-height: 14px;
        color: #333948;
        overflow: visible;
    }
    .property-list{
        padding: 4px 0;
        .property-item{
            width: 50%;
            margin: 12px 0 0;
            padding: 0 54px 0 0;
            font-size: 12px;
            .property-name{
                display: block;
                margin: 6px 0 10px;
                color: $cmdbTextColor;
                line-height: 16px;
                font-size: 0;
            }
            .property-name-text{
                position: relative;
                display: inline-block;
                max-width: calc(100% - 20px);
                padding: 0 10px 0 0;
                vertical-align: middle;
                font-size: 12px;
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
            .property-name-tooltips{
                display: inline-block;
                vertical-align: middle;
                width: 16px;
                height: 16px;
                font-size: 16px;
                color: #c3cdd7;
            }
            .property-value{
                height: 32px;
                line-height: 32px;
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
    .form-options{
        position: sticky;
        bottom: 0;
        left: 0;
        width: 100%;
        font-size: 0;
        padding: 10px 32px;
        &.sticky {
            border-top: 1px solid $cmdbBorderColor;
            background-color: #fff;
        }
        .button-save{
            min-width: 76px;
            margin-right: 4px;
        }
        .button-cancel{
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
