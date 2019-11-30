<template>
    <div class="form-layout">
        <feature-tips
            class="process-tips"
            :show-tips="true"
            :desc="$t('添加进程提示')">
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
                            <li class="property-item fl"
                                v-for="(property, propertyIndex) in groupedProperties[groupIndex]"
                                v-if="checkEditable(property)"
                                :key="propertyIndex">
                                <div class="property-name clearfix">
                                    <bk-checkbox class="form-checkbox"
                                        v-if="property['isLocking'] !== undefined"
                                        v-model="values[property['bk_property_id']]['as_default_value']"
                                        @change="handleResetValue(values[property['bk_property_id']]['as_default_value'], property)">
                                        <span class="property-name-text" :class="{ required: property['isrequired'] }">{{property['bk_property_name']}}</span>
                                    </bk-checkbox>
                                    <span v-else class="property-name-text" :class="{ required: property['isrequired'] }">{{property['bk_property_name']}}</span>
                                    <i class="property-name-tooltips icon-cc-tips"
                                        v-if="property['placeholder']"
                                        v-bk-tooltips="htmlEncode(property['placeholder'])">
                                    </i>
                                </div>
                                <div class="property-value">
                                    <component class="form-component"
                                        :is="`cmdb-form-${property['bk_property_type']}`"
                                        :disabled="getPropertyEditStatus(property)"
                                        :class="{ error: errors.has(property['bk_property_id']) }"
                                        :options="property.option || []"
                                        :data-vv-name="property['bk_property_id']"
                                        :data-vv-as="property['bk_property_name']"
                                        :placeholder="getPlaceholder(property)"
                                        :auto-select="false"
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
                <cmdb-auth :auth="$authResources(auth)">
                    <bk-button slot-scope="{ disabled }"
                        class="button-save"
                        theme="primary"
                        :disabled="saveDisabled || $loading() || disabled"
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
            isCreatedService: {
                type: Boolean,
                default: true
            },
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
                default: () => ({})
            }
        },
        data () {
            return {
                ipOption: [
                    {
                        'name': '127.0.0.1',
                        'type': 'text',
                        'is_default': true,
                        'id': '1'
                    },
                    {
                        'id': '2',
                        'name': '0.0.0.0',
                        'type': 'text',
                        'is_default': false
                    }
                    // {
                    //     'name': '第一内网IP',
                    //     'type': 'text',
                    //     'is_default': false,
                    //     'id': '3'
                    // },
                    // {
                    //     'name': '第一外网IP',
                    //     'type': 'text',
                    //     'is_default': false,
                    //     'id': '4'
                    // }
                ],
                values: {
                    bk_func_name: ''
                },
                refrenceValues: {},
                scrollbar: false
            }
        },
        computed: {
            ...mapGetters('serviceProcess', ['hasProcessName']),
            groupedProperties () {
                const properties = this.$groupedProperties.map(properties => {
                    const filterProperties = properties.filter(property => !['singleasst', 'multiasst', 'foreignkey'].includes(property['bk_property_type']))
                    filterProperties.map(property => {
                        if (!['bk_func_name', 'bk_process_name'].includes(property['bk_property_id'])) {
                            property.isLocking = false
                        }
                        if (['bind_ip'].includes(property['bk_property_id'])) {
                            property.bk_property_type = 'enum'
                            property.option = this.ipOption
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
            getPropertyEditStatus (property) {
                return (this.type === 'update' && ['bk_func_name', 'bk_process_name'].includes(property['bk_property_id']))
                    || !this.values[property['bk_property_id']]['as_default_value']
            },
            changedValues () {
                const changedValues = {}
                if (!this.values['bind_ip']['value']) {
                    this.$set(this.values.bind_ip, 'value', '')
                    this.$set(this.refrenceValues.bind_ip, 'value', '')
                }
                Object.keys(this.values).forEach(propertyId => {
                    const isChange = Object.keys(this.values[propertyId]).some(key => {
                        return JSON.stringify(this.values[propertyId][key]) !== JSON.stringify(this.refrenceValues[propertyId][key])
                    })
                    if (isChange) {
                        changedValues[propertyId] = this.values[propertyId]
                    }
                })
                return changedValues
            },
            hasChange () {
                return !!Object.keys(this.changedValues()).length
            },
            checkScrollbar () {
                const $layout = this.$el
                this.scrollbar = $layout.scrollHeight !== $layout.offsetHeight
            },
            initValues () {
                const inst = {}
                if (this.type === 'update') {
                    Object.keys(this.inst).forEach(key => {
                        const type = typeof this.inst[key]
                        if (type === 'object') {
                            inst[key] = this.inst[key] ? this.inst[key]['value'] : this.inst[key]
                        } else {
                            inst[key] = this.inst[key]
                        }
                    })
                }
                const formValues = this.$tools.getInstFormValues(this.properties, inst)
                Object.keys(formValues).forEach(key => {
                    this.$set(this.values, key, {
                        value: formValues[key],
                        as_default_value: this.type === 'update'
                            ? this.inst[key] ? this.inst[key]['as_default_value'] : false
                            : ['bk_func_name', 'bk_process_name'].includes(key)
                    })
                })
                if (this.isCreatedService && this.type === 'update') {
                    this.values['sign_id'] = inst['sign_id']
                } else if (this.type === 'update') {
                    this.values['process_id'] = inst['process_id']
                }
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
            handleSave () {
                this.$validator.validateAll().then(result => {
                    if (!this.hasChange()) {
                        this.$emit('on-cancel')
                        return
                    }
                    if (result && this.isCreatedService) {
                        if (this.type === 'create' && !this.hasProcessName(this.values)) {
                            this.values['sign_id'] = new Date().getTime()
                            this.addLocalProcessTemplate(this.values)
                            this.$emit('on-cancel')
                        } else if (this.type === 'update') {
                            this.updateLocalProcessTemplate(this.values)
                            this.$emit('on-cancel')
                        } else {
                            this.$bkMessage({
                                message: this.$t('进程名称已存在'),
                                theme: 'error'
                            })
                        }
                    } else if (result) {
                        this.$emit('on-submit', this.values, this.changedValues(), this.type)
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
            },
            handleResetValue (status, property) {
                if (!status) {
                    const type = property['bk_property_type']
                    if (['bool'].includes(type)) {
                        this.values[property['bk_property_id']]['value'] = false
                    } else if (['int'].includes(type)) {
                        this.values[property['bk_property_id']]['value'] = null
                    } else {
                        this.values[property['bk_property_id']]['value'] = ''
                    }
                }
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
        padding: 0 20px;
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
            font-size: 12px;
            &:nth-child(odd) {
                padding-right: 30px;
            }
            &:nth-child(even) {
                padding-left: 30px;
            }
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
                vertical-align: middle;
                padding: 0 14px 0 0;
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
                height: 32px;
                line-height: 32px;
                font-size: 12px;
                position: relative;
                /deep/ .control-append-group {
                    .bk-input-text {
                        flex: 1;
                    }
                }
            }
            .form-checkbox {
                outline: 0;
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
