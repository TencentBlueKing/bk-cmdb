<template>
    <div class="form-layout">
        <div class="form-groups">
            <template v-for="(group, groupIndex) in $sortedGroups">
                <div class="property-group"
                    :key="groupIndex"
                    v-if="checkGroupAvailable(groupedProperties[groupIndex])">
                    <h3 class="group-name">{{group['bk_group_name']}}</h3>
                    <ul class="property-list">
                        <li class="property-item clearfix"
                            v-for="(property, propertyIndex) in groupedProperties[groupIndex]"
                            v-if="checkEditable(property)"
                            :key="propertyIndex">
                            <div class="property-name fl">
                                <span class="property-name-text" :class="{required: property['isrequired']}">{{property['bk_property_name']}}</span>
                            </div>
                            <div class="property-value fl">
                                <component class="form-component"
                                    v-if="property['bk_property_type'] === 'enum'"
                                    :is="`cmdb-form-${property['bk_property_type']}`"
                                    :class="{error: errors.has(property['bk_property_id'])}"
                                    :disabled="checkDisabled(property)"
                                    :options="property.option || []"
                                    :data-vv-name="property['bk_property_name']"
                                    v-validate="getValidateRules(property)"
                                    v-model.trim="values[property['bk_property_id']]">
                                </component>
                                 <component class="form-component"
                                    v-else
                                    :is="`cmdb-form-${property['bk_property_type']}`"
                                    :class="{error: errors.has(property['bk_property_id'])}"
                                    :disabled="checkDisabled(property)"
                                    :data-vv-name="property['bk_property_name']"
                                    v-validate="getValidateRules(property)"
                                    v-model.trim="values[property['bk_property_id']]">
                                </component>
                                <span class="form-error">{{errors.first(property['bk_property_name'])}}</span>
                            </div>
                            <div class="property-tips fl">
                                <i class="property-name-tooltips bk-icon icon-info-circle-shape" v-if="property['placeholder']" v-tooltip="htmlEncode(property['placeholder'])"></i>
                            </div>
                        </li>
                    </ul>
                </div>
            </template>
        </div>
        <slot name="form-options">
            <div class="form-options" v-if="showOptions">
                <bk-button class="options-btn button-save" type="primary"
                    :disabled="!hasChange || $loading()"
                    @click="handleSave">
                    {{$t("Common['保存']")}}
                </bk-button>
                <bk-button  class="options-btn button-cancel" type="default" @click="handleCancel">{{$t("Common['取消']")}}</bk-button>
                <bk-button  class="options-btn button-delete" type="danger" @click="handleDelete">{{$t("Common['删除']")}}</bk-button>
            </div>
        </slot>
    </div>
</template>

<script>
    import formMixins from '@/mixins/form'
    export default {
        mixins: [formMixins],
        props: {
            inst: {
                type: Object,
                default () {
                    return {}
                }
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
                for (let propertyId in this.values) {
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
                    return properties.filter(property => !['singleasst', 'multiasst'].includes(property['bk_property_type']))
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
                return !property.editable
            },
            htmlEncode (placeholder) {
                let temp = document.createElement('div')
                temp.innerHTML = placeholder
                let output = temp.innerText
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
                        rules['regex'] = option
                    }
                }
                if (['singlechar', 'longchar'].includes(propertyType)) {
                    rules[propertyType] = true
                }
                if (propertyType === 'int') {
                    rules['numeric'] = true
                }
                return rules
            },
            handleSave () {
                this.$validator.validateAll().then(result => {
                    if (result) {
                        this.$emit('on-submit', this.values, this.changedValues, this.inst, this.type)
                    }
                })
            },
            handleCancel () {
                this.$emit('on-cancel')
            },
            handleDelete () {
                this.$emit('on-delete')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-layout{
        height: 100%;
    }
    .form-groups{
        padding: 0 0 0 32px;
    }
    .property-group{
        padding: 17px 0 0 0;
        &:first-child{
            padding: 28px 0 0 0;
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
            margin: 12px 0 0;
            padding: 0 54px 0 0;
            font-size: 14px;
            height: 36px;
            line-height: 36px;
            .property-name{
                position: relative;
                display: block;
                width: 120px;
                padding: 0 16px 0 0;
                color: $cmdbTextColor;
                text-align: right;
                font-size: 0;
                &:after{
                    font-size: 14px;
                    content: ":";
                    position: absolute;
                    right: 10px;
                }
            }
            .property-name-text{
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
            .property-value{
                height: 36px;
                width: calc(100% - 140px);
                max-width: 450px;
                font-size: 12px;
                position: relative;
            }
            .property-tips{
                width: 20px;
                text-align: right;
            }
            .property-name-tooltips{
                display: inline-block;
                vertical-align: middle;
                width: 16px;
                height: 16px;
                font-size: 16px;
                color: #ffb400;
            }
        }
    }
    .form-options{
        padding: 20px 0 0 152px;
        .options-btn{
            margin: 0 10px 0 0;
        }
    }
    .form-error {
        position: absolute;
        top: 100%;
        left: 0;
        line-height: 14px;
        font-size: 12px;
        color: #ff5656;
    }
</style>