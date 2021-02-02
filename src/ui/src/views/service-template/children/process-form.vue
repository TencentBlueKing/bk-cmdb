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
                                <li :class="['property-item', { flex: property.bk_property_type === 'table' }]"
                                    v-if="checkEditable(property)"
                                    :key="propertyIndex">
                                    <div class="property-name clearfix" v-if="!invisibleNameProperties.includes(property['bk_property_id'])">
                                        <span class="property-name-text" :class="{ required: property['isrequired'] }">{{property['bk_property_name']}}</span>
                                        <i class="property-name-tooltips icon-cc-tips"
                                            v-if="property['placeholder']"
                                            v-bk-tooltips="{
                                                trigger: 'click',
                                                content: htmlEncode(property['placeholder'])
                                            }">
                                        </i>
                                    </div>
                                    <div :class="['property-value', { 'is-lock': values[property.bk_property_id].as_default_value }]">
                                        <component class="form-component" ref="formComponent"
                                            :is="getComponentType(property)"
                                            :disabled="getPropertyEditStatus(property)"
                                            :class="{ error: errors.has(property['bk_property_id']) }"
                                            :unit="property.unit"
                                            :row="2"
                                            :options="property.option || []"
                                            :data-vv-name="property['bk_property_id']"
                                            :data-vv-as="property['bk_property_name']"
                                            :placeholder="getPlaceholder(property)"
                                            :auto-select="false"
                                            v-validate="getValidateRules(property)"
                                            v-model.trim="values[property['bk_property_id']]['value']">
                                        </component>
                                        <span class="property-lock-state"
                                            v-if="allowLock(property)"
                                            v-bk-tooltips="{
                                                placement: 'top',
                                                interactive: false,
                                                content: isLocked(property) ? $t('取消锁定') : $t('进程模板锁定提示语'),
                                                delay: [100, 0]
                                            }"
                                            tabindex="-1"
                                            @click="toggleLockState(property)">
                                            <i class="icon-cc-lock-fill" v-if="values[property.bk_property_id].as_default_value"></i>
                                            <i class="icon-cc-lock-line" v-else></i>
                                        </span>
                                        <span class="form-error">{{getFormError(property)}}</span>
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
                <cmdb-auth :auth="auth">
                    <bk-button slot-scope="{ disabled }"
                        class="button-save"
                        theme="primary"
                        :disabled="saveDisabled || $loading() || disabled || btnStatus()"
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
    import { mapMutations } from 'vuex'
    import ProcessFormPropertyTable from './process-form-property-table'
    export default {
        components: {
            ProcessFormPropertyTable
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
            isCreatedService: {
                type: Boolean,
                default: true
            },
            dataIndex: Number,
            showOptions: {
                type: Boolean,
                default: true
            },
            saveDisabled: Boolean,
            hasUsed: {
                type: Boolean,
                default: false
            },
            auth: {
                type: Object,
                default: null
            },
            submitFormat: {
                type: Function,
                default: data => data
            }
        },
        data () {
            return {
                values: {
                    bk_func_name: ''
                },
                refrenceValues: {},
                scrollbar: false,
                invisibleNameProperties: ['bind_info'],
                defaultLocked: ['bk_func_name', 'bk_process_name', 'bind_info']
            }
        },
        computed: {
            groupedProperties () {
                return this.$groupedProperties
            }
        },
        watch: {
            inst (inst) {
                this.initValues()
            },
            properties () {
                this.initValues()
            },
            'values.bk_func_name.value': {
                handler (newVal, oldValue) {
                    if (this.values.bk_process_name.value === oldValue) {
                        this.values.bk_process_name.value = newVal
                    }
                },
                deep: true
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
            ...mapMutations('serviceProcess', ['addLocalProcessTemplate', 'updateLocalProcessTemplate']),
            isLocked (property) {
                return this.values[property.bk_property_id].as_default_value
            },
            allowLock (property) {
                return !this.defaultLocked.includes(property.bk_property_id)
            },
            toggleLockState (property) {
                this.values[property.bk_property_id].as_default_value = !this.isLocked(property)
            },
            getComponentType (property) {
                const type = property.bk_property_type
                if (type === 'table') {
                    return 'process-form-property-table'
                }
                return `cmdb-form-${type}`
            },
            getPropertyEditStatus (property) {
                const uneditable = ['bk_func_name', 'bk_process_name'].includes(property['bk_property_id']) && !this.isCreatedService
                return this.type === 'update' && uneditable
            },
            changedValues () {
                const changedValues = {}
                if (!Object.keys(this.refrenceValues).length) return {}
                Object.keys(this.values).forEach(propertyId => {
                    let isChange = false
                    if (!['sign_id', 'process_id'].includes(propertyId)) {
                        isChange = Object.keys(this.values[propertyId]).some(key => {
                            return JSON.stringify(this.values[propertyId][key]) !== JSON.stringify(this.refrenceValues[propertyId][key])
                        })
                    }
                    if (isChange) {
                        changedValues[propertyId] = this.values[propertyId]
                    }
                })
                return changedValues
            },
            hasChange () {
                return !!Object.keys(this.changedValues()).length
            },
            btnStatus () {
                return this.type === 'create' ? false : !this.hasChange()
            },
            checkScrollbar () {
                const $layout = this.$el
                this.scrollbar = $layout.scrollHeight !== $layout.offsetHeight
            },
            initValues () {
                const restValues = {}
                const formValues = this.$tools.getInstFormValues(this.properties, {}, this.type === 'create')
                Object.keys(formValues).forEach(key => {
                    if (!this.inst.hasOwnProperty(key)) {
                        restValues[key] = {
                            as_default_value: this.defaultLocked.includes(key),
                            value: formValues[key]
                        }
                    }
                })
                this.values = Object.assign({}, this.values, restValues, this.inst)
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
            getFormError (property) {
                if (property.bk_property_type === 'table') {
                    const hasError = this.errors.items.some(item => item.scope === property.bk_property_id)
                    return hasError ? this.$t('有未正确定义的监听信息') : ''
                }
                return this.errors.first(property.bk_property_id)
            },
            callComponentValidator () {
                const componentValidator = []
                const { formComponent = [] } = this.$refs
                formComponent.forEach(component => {
                    componentValidator.push(component.$validator.validateAll())
                    componentValidator.push(component.$validator.validateScopes())
                })
                return componentValidator
            },
            async handleSave () {
                try {
                    const results = await Promise.all([
                        this.$validator.validateAll(),
                        ...this.callComponentValidator()
                    ])
                    const result = results.every(result => result)
                    if (result && !this.hasChange()) {
                        this.$emit('on-cancel')
                    } else if (result && this.isCreatedService) {
                        const cloneValues = this.$tools.clone(this.values)
                        const formatValue = this.submitFormat(cloneValues)
                        if (this.type === 'create') {
                            this.addLocalProcessTemplate(formatValue)
                            this.$emit('on-cancel')
                        } else if (this.type === 'update') {
                            this.updateLocalProcessTemplate({ process: formatValue, index: this.dataIndex })
                            this.$emit('on-cancel')
                        }
                    } else if (result) {
                        this.$emit('on-submit', this.values, this.changedValues(), this.type)
                    } else {
                        this.uncollapseGroup()
                    }
                } catch (error) {
                    console.error(error)
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
                if (this.hasChange()) {
                    return new Promise((resolve, reject) => {
                        this.$bkInfo({
                            title: this.$t('确认退出'),
                            subTitle: this.$t('退出会导致未保存信息丢失'),
                            extCls: 'bk-dialog-sub-header-center',
                            confirmFn: () => {
                                this.$emit('on-cancel')
                            },
                            cancelFn: () => {
                                resolve(false)
                            }
                        })
                    })
                } else {
                    this.$emit('on-cancel')
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-layout {
        height: 100%;
        @include scrollbar-y;
    }
    .process-tips {
        margin: 10px 20px 0;
    }
    .form-groups {
        padding: 0 20px;
    }
    .property-group {
        padding: 20px 0 10px 0;
        &:first-child {
        padding: 15px 0 10px 0;
        }
    }
    .group-name {
        font-size: 14px;
        font-weight: bold;
        line-height: 14px;
        color: #63656e;
        overflow: visible;
    }
    .property-list {
        padding: 4px 0;
        display: flex;
        flex-wrap: wrap;
        .property-item {
            width: 50%;
            margin: 12px 0 0;
            font-size: 12px;
            flex: 0 0 50%;
            &:nth-child(odd) {
                padding-right: 30px;
            }
            &:nth-child(even) {
                padding-left: 30px;
            }
            .property-name {
                display: block;
                margin: 6px 0 10px;
                color: $cmdbTextColor;
                line-height: 16px;
                font-size: 0;
            }
            .property-name-text {
                position: relative;
                display: inline-block;
                vertical-align: middle;
                padding: 0 6px 0 0;
                font-size: 14px;
                @include ellipsis;
                &.required {
                    padding: 0 14px 0 0;
                    &:after {
                        position: absolute;
                        left: 100%;
                        top: 0;
                        margin: 0 0 0 -10px;
                        content: "*";
                        color: #ff5656;
                    }
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
            
            &.flex {
                flex: 1;
                padding-right: 0;
                width: 100%;
            }
        }
    }
    @mixin property-lock-state-visible {
        display: inline-flex;
        border: 1px solid #c4c6cc;
        border-left: none;
    }
    @mixin no-right-radius {
        border-top-right-radius: 0;
        border-bottom-right-radius: 0;
    }
    .property-value {
        font-size: 0;
        position: relative;
        display: flex;
        &:hover,
        &.is-lock {
            .property-lock-state {
               @include property-lock-state-visible;
            }
            .form-component /deep/ {
                .bk-form-input,
                .bk-form-textarea,
                .bk-textarea-wrapper {
                     @include no-right-radius;
                }
            }
            .form-component.bk-select {
                @include no-right-radius;
            }
        }
        .form-component {
            flex: 1;
            &.control-active /deep/ {
                .bk-form-input,
                .bk-form-textarea,
                .bk-textarea-wrapper {
                     @include no-right-radius;
                }
            }
            &.is-focus {
                @include no-right-radius;
            }
            &.control-active ~ .property-lock-state {
                @include property-lock-state-visible;
            }
        }
        .property-lock-state {
            display: none;
            width: 24px;
            align-items: center;
            justify-content: center;
            background-color: #f2f4f8;
            font-size: 14px;
            overflow: hidden;
            cursor: pointer;
        }
    }
    .form-options {
        position: sticky;
        bottom: 0;
        left: 0;
        width: 100%;
        padding: 28px 32px 0;
        font-size: 0;
        z-index: 101;
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
