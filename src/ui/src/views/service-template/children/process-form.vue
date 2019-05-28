<template>
    <div class="form-layout">
        <feature-tips
            class="process-tips"
            :show-tips="true"
            :desc="$t('ProcessManagement[\'添加进程提示\']')">
        </feature-tips>
        <!-- <div class="form-groups" ref="formGroups">
            <div class="property-group">
                <div class="group-name">{{$t("ProcessManagement['基本属性']")}}</div>
                <ul class="property-list clearfix">
                    <li class="property-item fl">
                        <div class="property-name clearfix">
                            <span class="property-name-text required">{{$t("ProcessManagement['进程名称']")}}</span>
                        </div>
                        <div class="property-value">
                            <cmdb-form-singlechar class="form-component"
                                :class="{ error: errors.has('bk_func_name') }"
                                :data-vv-name="'bk_func_name'"
                                :data-vv-as="$t('ProcessManagement[\'进程名称\']')"
                                v-validate="'required|singlechar'"
                                v-model.trim="values['bk_func_name']">
                            </cmdb-form-singlechar>
                            <span class="form-error">{{errors.first('bk_func_name')}}</span>
                        </div>
                    </li>
                    <li class="property-item fl">
                        <div class="property-name clearfix">
                            <span class="property-name-text required">{{$t("ProcessManagement['进程别名']")}}</span>
                        </div>
                        <div class="property-value">
                            <cmdb-form-singlechar class="form-component"
                                :class="{ error: errors.has('bk_process_name') }"
                                :data-vv-name="'bk_process_name'"
                                :data-vv-as="$t('ProcessManagement[\'进程别名\']')"
                                v-validate="'required|singlechar'"
                                v-model.trim="values['bk_process_name']">
                            </cmdb-form-singlechar>
                            <span class="form-error">{{errors.first('bk_process_name')}}</span>
                        </div>
                    </li>
                    <li class="property-item fl">
                        <div class="property-name clearfix">
                            <span class="property-name-text">{{$t("ProcessManagement['启动参数匹配规则']")}}</span>
                            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                                <input type="checkbox">
                                <span class="cmdb-checkbox-text">{{$t('ProcessManagement["锁定"]')}}</span>
                            </label>
                        </div>
                        <div class="property-value">
                            <cmdb-form-longchar class="form-component"
                                :class="{ error: errors.has('bk_start_param_regex') }"
                                :data-vv-name="'bk_start_param_regex'"
                                :data-vv-as="$t('ProcessManagement[\'启动参数匹配规则\']')"
                                v-validate="'longchar'"
                                v-model.trim="values['bk_start_param_regex']">
                            </cmdb-form-longchar>
                            <span class="form-error">{{errors.first('bk_start_param_regex')}}</span>
                        </div>
                    </li>
                    <li class="property-item fl">
                        <div class="property-name clearfix">
                            <span class="property-name-text">{{$t("ProcessManagement['进程描述']")}}</span>
                            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                                <input type="checkbox">
                                <span class="cmdb-checkbox-text">{{$t('ProcessManagement["锁定"]')}}</span>
                            </label>
                        </div>
                        <div class="property-value">
                            <cmdb-form-singlechar class="form-component"
                                :class="{ error: errors.has('description') }"
                                :data-vv-name="'description'"
                                :data-vv-as="$t('ProcessManagement[\'进程描述\']')"
                                v-validate="'singlechar'"
                                v-model.trim="values['description']">
                            </cmdb-form-singlechar>
                            <span class="form-error">{{errors.first('description')}}</span>
                        </div>
                    </li>
                </ul>
            </div>
            <div class="property-group">
                <div class="group-name">{{$t("ProcessManagement['端口信息']")}}</div>
                <ul class="property-list clearfix">
                    <li class="property-item fl">
                        <div class="property-name clearfix">
                            <span class="property-name-text">{{$t("ProcessManagement['监听IP']")}}</span>
                            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                                <input type="checkbox">
                                <span class="cmdb-checkbox-text">{{$t('ProcessManagement["锁定"]')}}</span>
                            </label>
                        </div>
                        <div class="property-value">
                            <cmdb-form-enum class="form-component" v-model.trim="values['bind_ip']"></cmdb-form-enum>
                        </div>
                    </li>
                    <li class="property-item fl">
                        <div class="property-name clearfix">
                            <span class="property-name-text">{{$t("ProcessManagement['端口']")}}</span>
                            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                                <input type="checkbox">
                                <span class="cmdb-checkbox-text">{{$t('ProcessManagement["锁定"]')}}</span>
                            </label>
                        </div>
                        <div class="property-value">
                            <cmdb-form-singlechar class="form-component"
                                :class="{ error: errors.has('port') }"
                                :data-vv-name="'port'"
                                :data-vv-as="$t('ProcessManagement[\'端口\']')"
                                v-validate="'singlechar'"
                                v-model.trim="values['port']">
                            </cmdb-form-singlechar>
                            <span class="form-error">{{errors.first('port')}}</span>
                        </div>
                    </li>
                    <li class="property-item fl">
                        <div class="property-name clearfix">
                            <span class="property-name-text">{{$t("ProcessManagement['协议']")}}</span>
                            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                                <input type="checkbox">
                                <span class="cmdb-checkbox-text">{{$t('ProcessManagement["锁定"]')}}</span>
                            </label>
                        </div>
                        <div class="property-value">
                            <cmdb-form-enum class="form-component" v-model.trim="values['protocol']"></cmdb-form-enum>
                        </div>
                    </li>
                </ul>
            </div>
            <div class="property-group">
                <div class="group-name">{{$t("ProcessManagement['其他属性']")}}</div>
                <ul class="property-list clearfix">
                    <li class="property-item fl">
                        <div class="property-name clearfix">
                            <span class="property-name-text">{{$t("ProcessManagement['启动用户']")}}</span>
                            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                                <input type="checkbox">
                                <span class="cmdb-checkbox-text">{{$t('ProcessManagement["锁定"]')}}</span>
                            </label>
                        </div>
                        <div class="property-value">
                            <cmdb-form-singlechar class="form-component"
                                :class="{ error: errors.has('user') }"
                                :data-vv-name="'user'"
                                :data-vv-as="$t('ProcessManagement[\'启动用户\']')"
                                v-validate="'singlechar'"
                                v-model.trim="values['user']">
                            </cmdb-form-singlechar>
                            <span class="form-error">{{errors.first('user')}}</span>
                        </div>
                    </li>
                    <li class="property-item fl">
                        <div class="property-name clearfix">
                            <span class="property-name-text">{{$t("ProcessManagement['工作路径']")}}</span>
                            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                                <input type="checkbox">
                                <span class="cmdb-checkbox-text">{{$t('ProcessManagement["锁定"]')}}</span>
                            </label>
                        </div>
                        <div class="property-value">
                            <cmdb-form-longchar class="form-component"
                                :class="{ error: errors.has('work_path') }"
                                :data-vv-name="'work_path'"
                                :data-vv-as="$t('ProcessManagement[\'工作路径\']')"
                                v-validate="'longchar'"
                                v-model.trim="values['work_path']">
                            </cmdb-form-longchar>
                            <span class="form-error">{{errors.first('work_path')}}</span>
                        </div>
                    </li>
                    <li class="property-item fl">
                        <div class="property-name clearfix">
                            <span class="property-name-text">{{$t("ProcessManagement['启动数量']")}}</span>
                            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                                <input type="checkbox">
                                <span class="cmdb-checkbox-text">{{$t('ProcessManagement["锁定"]')}}</span>
                            </label>
                        </div>
                        <div class="property-value">
                            <cmdb-form-int class="form-component"
                                :class="{ error: errors.has('proc_num') }"
                                :data-vv-name="'proc_num'"
                                :data-vv-as="$t('ProcessManagement[\'启动数量\']')"
                                v-validate="'int'"
                                v-model.trim="values['proc_num']">
                            </cmdb-form-int>
                            <span class="form-error">{{errors.first('proc_num')}}</span>
                        </div>
                    </li>
                    <li class="property-item fl">
                        <div class="property-name clearfix">
                            <span class="property-name-text">{{$t("ProcessManagement['启动优先级']")}}</span>
                            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                                <input type="checkbox">
                                <span class="cmdb-checkbox-text">{{$t('ProcessManagement["锁定"]')}}</span>
                            </label>
                        </div>
                        <div class="property-value">
                            <cmdb-form-int class="form-component"
                                :class="{ error: errors.has('priority') }"
                                :data-vv-name="'priority'"
                                :data-vv-as="$t('ProcessManagement[\'启动优先级\']')"
                                v-validate="'int'"
                                v-model.trim="values['priority']">
                            </cmdb-form-int>
                            <span class="form-error">{{errors.first('priority')}}</span>
                        </div>
                    </li>
                    <li class="property-item fl">
                        <div class="property-name clearfix">
                            <span class="property-name-text">{{$t("ProcessManagement['启动超时时长']")}}</span>
                            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                                <input type="checkbox">
                                <span class="cmdb-checkbox-text">{{$t('ProcessManagement["锁定"]')}}</span>
                            </label>
                        </div>
                        <div class="property-value">
                            <cmdb-form-int class="form-component"
                                :class="{ error: errors.has('timeout') }"
                                :data-vv-name="'timeout'"
                                :data-vv-as="$t('ProcessManagement[\'启动超时时长\']')"
                                v-validate="'int'"
                                v-model.trim="values['timeout']">
                            </cmdb-form-int>
                            <span class="form-error">{{errors.first('timeout')}}</span>
                        </div>
                    </li>
                    <li class="property-item fl">
                        <div class="property-name clearfix">
                            <span class="property-name-text">{{$t("ProcessManagement['启动命令']")}}</span>
                            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                                <input type="checkbox">
                                <span class="cmdb-checkbox-text">{{$t('ProcessManagement["锁定"]')}}</span>
                            </label>
                        </div>
                        <div class="property-value">
                            <cmdb-form-longchar class="form-component"
                                :class="{ error: errors.has('start_cmd') }"
                                :data-vv-name="'start_cmd'"
                                :data-vv-as="$t('ProcessManagement[\'启动命令\']')"
                                v-validate="'longchar'"
                                v-model.trim="values['start_cmd']">
                            </cmdb-form-longchar>
                            <span class="form-error">{{errors.first('start_cmd')}}</span>
                        </div>
                    </li>
                    <li class="property-item fl">
                        <div class="property-name clearfix">
                            <span class="property-name-text">{{$t("ProcessManagement['停止命令']")}}</span>
                            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                                <input type="checkbox">
                                <span class="cmdb-checkbox-text">{{$t('ProcessManagement["锁定"]')}}</span>
                            </label>
                        </div>
                        <div class="property-value">
                            <cmdb-form-longchar class="form-component"
                                :class="{ error: errors.has('stop_cmd') }"
                                :data-vv-name="'stop_cmd'"
                                :data-vv-as="$t('ProcessManagement[\'停止命令\']')"
                                v-validate="'longchar'"
                                v-model.trim="values['stop_cmd']">
                            </cmdb-form-longchar>
                            <span class="form-error">{{errors.first('stop_cmd')}}</span>
                        </div>
                    </li>
                    <li class="property-item fl">
                        <div class="property-name clearfix">
                            <span class="property-name-text">{{$t("ProcessManagement['重启命令']")}}</span>
                            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                                <input type="checkbox">
                                <span class="cmdb-checkbox-text">{{$t('ProcessManagement["锁定"]')}}</span>
                            </label>
                        </div>
                        <div class="property-value">
                            <cmdb-form-longchar class="form-component"
                                :class="{ error: errors.has('restart_cmd') }"
                                :data-vv-name="'restart_cmd'"
                                :data-vv-as="$t('ProcessManagement[\'重启命令\']')"
                                v-validate="'longchar'"
                                v-model.trim="values['restart_cmd']">
                            </cmdb-form-longchar>
                            <span class="form-error">{{errors.first('restart_cmd')}}</span>
                        </div>
                    </li>
                    <li class="property-item fl">
                        <div class="property-name clearfix">
                            <span class="property-name-text">{{$t("ProcessManagement['强制停止命令']")}}</span>
                            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                                <input type="checkbox">
                                <span class="cmdb-checkbox-text">{{$t('ProcessManagement["锁定"]')}}</span>
                            </label>
                        </div>
                        <div class="property-value">
                            <cmdb-form-longchar class="form-component"
                                :class="{ error: errors.has('restart_cmd') }"
                                :data-vv-name="'restart_cmd'"
                                :data-vv-as="$t('ProcessManagement[\'重启命令\']')"
                                v-validate="'longchar'"
                                v-model.trim="values['restart_cmd']">
                            </cmdb-form-longchar>
                            <span class="form-error">{{errors.first('restart_cmd')}}</span>
                        </div>
                    </li>
                </ul>
            </div>
        </div> -->
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
                                        :disabled="checkDisabled(property)"
                                        :options="property.option || []"
                                        :data-vv-name="property['bk_property_id']"
                                        :data-vv-as="property['bk_property_name']"
                                        v-validate="getValidateRules(property)"
                                        v-model.trim="values[property['bk_property_id']]">
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
                }
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
                const properties = this.$groupedProperties.map(properties => {
                    const filterProperties = properties.filter(property => !['singleasst', 'multiasst', 'foreignkey'].includes(property['bk_property_type']))
                    filterProperties.map(property => {
                        if (property['bk_property_id'] === 'bk_func_name') {
                            property['bk_property_name'] = '进程名称'
                        }
                        if (property['bk_property_id'] === 'bk_process_name') {
                            property['bk_property_name'] = '进程别名'
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
                        this.values[key] = formValues[key]
                    } else {
                        this.values[key] = {
                            value: formValues[key],
                            as_default_value: false
                        }
                    }
                })
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
            handleSave () {
                this.$validator.validateAll().then(result => {
                    if (result) {
                        console.log(this.values)
                        console.log(this.changedValues)
                        this.$emit('on-submit', this.values, this.changedValues, this.inst, this.type)
                    } else {
                        this.uncollapseGroup()
                    }
                })
            },
            handleSelected (value) {
                
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
