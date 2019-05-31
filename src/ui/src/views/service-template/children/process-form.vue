<template>
    <div class="form-layout">
        <feature-tips
            class="process-tips"
            :show-tips="true"
            :desc="$t('ProcessManagement[\'添加进程提示\']')">
        </feature-tips>
        <div class="form-groups" ref="formGroups">
            <template v-for="(group, groupIndex) in $sortedGroups">
                <div class="property-group"
                    :key="groupIndex"
                    v-if="checkGroupAvailable(groupedProperties[groupIndex])">
                    <cmdb-collapse
                        :label="group['bk_group_name']"
                        :collapse.sync="groupState[group['bk_group_id']]">
                        <ul class="property-list clearfix">
                            <li class="property-item"
                                :style="{ float: property['bk_property_id'] === 'bk_process_name' ? 'right !important' : 'left !important' }"
                                v-for="(property, propertyIndex) in groupedProperties[groupIndex]"
                                v-if="checkEditable(property)"
                                :key="propertyIndex">
                                <div class="property-name clearfix">
                                    <span class="property-name-text" :class="{ required: property['isrequired'] }">{{property['bk_property_name']}}</span>
                                    <!-- <i class="property-name-tooltips icon-cc-tips"
                                        v-if="property['placeholder'] && property['bk_property_id'] === 'bk_start_param_regex'"
                                        v-tooltip="htmlEncode(property['placeholder'])">
                                    </i> -->
                                    <label class="cmdb-form-checkbox cmdb-checkbox-small" v-if="property['isLocking'] !== undefined">
                                        <input type="checkbox" v-model="values[property['bk_property_id']]['as_default_value']">
                                        <span class="cmdb-checkbox-text">{{$t('ProcessManagement["锁定"]')}}</span>
                                    </label>
                                </div>
                                <div class="property-value">
                                    <component class="form-component" v-if="['bk_func_name', 'bk_process_name'].includes(property['bk_property_id'])"
                                        :is="`cmdb-form-${property['bk_property_type']}`"
                                        :class="{ error: errors.has(property['bk_property_id']) }"
                                        :disabled="type === 'update'"
                                        :options="property.option || []"
                                        :data-vv-name="property['bk_property_id']"
                                        :data-vv-as="property['bk_property_name']"
                                        v-validate="getValidateRules(property)"
                                        v-model.trim="values[property['bk_property_id']]"
                                        @input="handleFuncNameInput(property['bk_property_id'])"
                                        @on-change="handleFuncNameChange(property['bk_property_id'])">
                                    </component>
                                    <component class="form-component" v-else
                                        :is="`cmdb-form-${property['bk_property_type']}`"
                                        :class="{ error: errors.has(property['bk_property_id']) }"
                                        :disabled="checkDisabled(property)"
                                        :options="property.option || []"
                                        :data-vv-name="property['bk_property_id']"
                                        :data-vv-as="property['bk_property_name']"
                                        v-validate="getValidateRules(property)"
                                        v-model.trim="values[property['bk_property_id']]['value']">
                                    </component>
                                    <span class="form-error">{{errors.first(property['bk_property_id'])}}</span>
                                </div>
                            </li>
                        </ul>
                    </cmdb-collapse>
                </div>
            </template>
        </div>
        <div class="form-options"
            v-if="showOptions"
            :class="{ sticky: scrollbar }">
            <slot name="form-options">
                <bk-button class="button-save" type="primary"
                    :disabled="saveDisabled || !hasChange || $loading()"
                    @click="handleSave">
                    {{$t("Common['保存']")}}
                </bk-button>
                <bk-button class="button-cancel" @click="handleCancel">{{$t("Common['取消']")}}</bk-button>
            </slot>
            <slot name="extra-options"></slot>
        </div>
    </div>
</template>

<script>
    import formMixins from '@/mixins/form'
    import RESIZE_EVENTS from '@/utils/resize-events'
    import featureTips from '@/components/feature-tips/index'
    import { mapGetters, mapMutations } from 'vuex'
    export default {
        components: {
            featureTips
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
            isCreateTemplate: {
                type: Boolean,
                default: true
            },
            showOptions: {
                type: Boolean,
                default: true
            },
            saveDisabled: Boolean
        },
        data () {
            return {
                values: {},
                refrenceValues: {},
                scrollbar: false,
                groupState: {
                    none: true
                },
                autoInput: true
            }
        },
        computed: {
            ...mapGetters('serviceProcess', ['hasProcessName']),
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
                const properties = this.$groupedProperties.map(properties => {
                    const filterProperties = properties.filter(property => !['singleasst', 'multiasst', 'foreignkey'].includes(property['bk_property_type']))
                    filterProperties.map(property => {
                        if (property['bk_property_id'] === 'bk_process_name') {
                            property['isrequired'] = true
                            property['bk_property_name'] = this.$t("ProcessManagement['进程别名']")
                        }
                        if (property['bk_property_id'] === 'bk_func_name') {
                            property['bk_property_name'] = this.$t("ProcessManagement['进程名称']")
                        }
                        if (!['bk_func_name', 'bk_process_name'].includes(property['bk_property_id'])) {
                            property.isLocking = false
                        }
                    })
                    return filterProperties
                })
                return properties
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
            // console.log(this.$sortedGroups)
            // console.log(this.groupedProperties)
        },
        mounted () {
            RESIZE_EVENTS.addResizeListener(this.$refs.formGroups, this.checkScrollbar)
        },
        beforeDestroy () {
            RESIZE_EVENTS.removeResizeListener(this.$refs.formGroups, this.checkScrollbar)
        },
        methods: {
            ...mapMutations('serviceProcess', ['addLocalProcessTemplate']),
            checkScrollbar () {
                const $layout = this.$el
                this.scrollbar = $layout.scrollHeight !== $layout.offsetHeight
            },
            initValues () {
                const inst = {}
                if (this.type === 'update') {
                    Object.keys(this.inst).forEach(key => {
                        inst[key] = this.inst[key] ? this.inst[key]['value'] : this.inst[key]
                    })
                }
                const formValues = this.$tools.getInstFormValues(this.properties, inst)
                Object.keys(formValues).forEach(key => {
                    if (['bk_func_name', 'bk_process_name'].includes(key)) {
                        this.values[key] = this.type === 'update' ? this.inst[key] : formValues[key]
                    } else {
                        this.values[key] = {
                            value: formValues[key],
                            as_default_value: this.type === 'update' ? this.inst[key]['as_default_value'] : false
                        }
                    }
                })
                // console.log(this.values)
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
                if (propertyType === 'float') {
                    rules['float'] = true
                }
                return rules
            },
            handleFuncNameChange (id) {
                this.autoInput = !this.values['bk_process_name']
            },
            handleFuncNameInput (id) {
                if (id === 'bk_func_name' && this.autoInput) {
                    this.values['bk_process_name'] = this.values['bk_func_name']
                }
            },
            handleSave () {
                this.$validator.validateAll().then(result => {
                    if (result) {
                        if (this.type === 'create' && !this.hasProcessName(this.values)) {
                            this.values['id'] = new Date().getTime()
                            this.addLocalProcessTemplate(this.values)
                            this.handleCancel()
                        } else {
                            this.$bkMessage({
                                message: this.$t("ServiceManagement['进程名称已存在']"),
                                theme: 'error'
                            })
                        }
                        // this.$emit('on-submit', this.values, this.changedValues, this.inst, this.type)
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
    .process-tips {
        margin: 10px 20px 0;
    }
    .form-groups{
        padding-left: 20px;
    }
    .property-group{
        padding: 20px 0 10px 0;
        &:first-child {
        padding: 15px 0 10px 0;
        }
    }
    .group-name{
        font-size: 14px;
        font-weight: bold;
        line-height: 14px;
        color: #63656e;
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
                margin: 6px 0 9px;
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
            .property-name-tooltips{
                display: inline-block;
                vertical-align: middle;
                width: 16px;
                height: 16px;
                font-size: 16px;
                color: #c3cdd7;
            }
            .property-value{
                height: 36px;
                line-height: 36px;
                font-size: 12px;
                position: relative;
            }
            .cmdb-form-checkbox {
                padding: 0;
                margin-right: 0;
                float: right;
                line-height: 16px;
                cursor: pointer;
                .cmdb-checkbox-text {
                    font-size: 12px;
                }
            }
        }
    }
    .form-options{
        position: sticky;
        bottom: 0;
        left: 0;
        width: 100%;
        padding: 28px 32px 0;
        font-size: 0;
        &.sticky {
            padding: 10px 32px;
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
    }
</style>
