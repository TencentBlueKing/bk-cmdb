/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

<template>
    <div class="attribute-wrapper">
        <template v-if="displayType === 'list'">
            <slot name="list"></slot>
            <template v-for="propertyGroup in groupOrder" v-if="bkPropertyGroups.hasOwnProperty(propertyGroup)">
                <div class="attribute-group" v-show="!(propertyGroup === 'none' && isNoneGroupHide)">
                    <h3 class="title">{{propertyGroup === 'none' ? $t("Common['更多属性']") : bkPropertyGroups[propertyGroup]['bkPropertyGroupName']}}</h3>
                    <ul class="clearfix attribute-list">
                        <template v-for="(property, propertyIndex) in bkPropertyGroups[propertyGroup]['properties']">
                            <li class="attribute-item fl" v-if="!property['bk_isapi']" :key="propertyIndex">
                                <template v-if="property['bk_property_type'] !== 'bool'">
                                    <span class="attribute-item-label has-colon" :title="property['bk_property_name']">{{property['bk_property_name']}}</span>
                                    <span class="attribute-item-value" :title="getFieldValue(property)">{{getFieldValue(property) === '' ? '--' : getFieldValue(property)}}</span>
                                </template>
                                <template v-else>
                                    <span class="attribute-item-label" :title="property['bk_property_name']">{{property['bk_property_name']}}</span>
                                    <span class="attribute-item-value bk-form-checkbox">
                                        <input type="checkbox" :checked="getFieldValue(property)" disabled>
                                    </span>
                                </template>
                            </li>
                        </template>
                    </ul>
                </div>
                 <div class="attribute-group-more" v-if="propertyGroup === 'none'">
                    <a href="javascript:void(0)" class="group-more-link" :class="{'open': !isNoneGroupHide}" @click="isNoneGroupHide = !isNoneGroupHide">{{$t("Common['更多属性']")}}</a>
                </div>
            </template>
        </template>
        <template v-else-if="displayType === 'form'">
            <form @submit.prevent v-if="!isMultipleUpdate || (isMultipleUpdate && hasEditableProperties)">
                <template v-for="propertyGroup in groupOrder">
                    <template v-if="checkIsShowGroup(propertyGroup)">
                        <div class="attribute-group" v-show="!(propertyGroup === 'none' && isNoneGroupHide)">
                            <h3 class="title">{{propertyGroup === 'none' ? $t("Common['更多属性']") : bkPropertyGroups[propertyGroup]['bkPropertyGroupName']}}</h3>
                                <ul class="clearfix attribute-list edit">
                                    <template v-for="(property, propertyIndex) in bkPropertyGroups[propertyGroup]['properties']">
                                        <li class="attribute-item fl" :class="property['bk_property_type']" :key="propertyIndex"
                                            v-if="checkIsShowField(property)">
                                            <div>
                                                <label :class="[{'required': property['isrequired']}]" class="bk-form-checkbox bk-checkbox-small">
                                                    <input type="checkbox" v-if="isMultipleUpdate" 
                                                        v-model="multipleEditableFields[property['bk_property_id']]"
                                                        @change="clearFieldValue(property)">
                                                    <span>{{property['bk_property_name']}}</span>
                                                </label>
                                                <i class="icon-tooltips" v-if="property['placeholder']" v-tooltip="htmlEncode(property['placeholder'])"></i>
                                            </div>
                                            <div class="attribute-item-field">
                                                <v-member-selector v-if="property['bk_property_type'] === 'objuser'"
                                                    :disabled="checkIsFieldDisabled(property)"
                                                    :selected.sync="localValues[property['bk_property_id']]" 
                                                    :multiple="true">
                                                </v-member-selector>
                                                <bk-datepicker v-else-if="property['bk_property_type'] === 'date'"
                                                    :timer="false"
                                                    :init-date="localValues[property['bk_property_id']]"
                                                    :disabled="checkIsFieldDisabled(property)"
                                                    @date-selected="setDate(...arguments, property['bk_property_id'])">
                                                </bk-datepicker>
                                                <bk-datepicker v-else-if="property['bk_property_type'] === 'time'"
                                                    :timer="true"
                                                    :init-date="localValues[property['bk_property_id']]"
                                                    :disabled="checkIsFieldDisabled(property)"
                                                    @date-selected="setDate(...arguments, property['bk_property_id'])">
                                                </bk-datepicker>
                                                <template v-else-if="property['bk_property_type'] === 'singleasst' || property['bk_property_type'] === 'multiasst'">
                                                    <v-host v-if="property['bk_asst_obj_id'] === 'host'"
                                                        :multiple="property['bk_property_type'] === 'multiasst'"
                                                        :selected.sync="localValues[property['bk_property_id']]">
                                                    </v-host>
                                                    <v-association v-else
                                                        :multiple="property['bk_property_type'] === 'multiasst'"
                                                        :selected.sync="localValues[property['bk_property_id']]"
                                                        :disabled="checkIsFieldDisabled(property)"
                                                        :asstObjId="property['bk_asst_obj_id']">
                                                    </v-association>
                                                </template>
                                                <v-enumeration v-else-if="property['bk_property_type'] === 'enum'"
                                                    :selected.sync="localValues[property['bk_property_id']]"
                                                    :disabled="checkIsFieldDisabled(property)"
                                                    :options="property['option']">
                                                </v-enumeration>
                                                <v-timezone v-else-if="property['bk_property_type'] === 'timezone'"
                                                    :selected.sync="localValues[property['bk_property_id']]"
                                                    :disabled="checkIsFieldDisabled(property)"
                                                ></v-timezone>
                                                <span class="bk-form-checkbox" v-else-if="property['bk_property_type'] === 'bool'">
                                                    <input
                                                        type="checkbox"
                                                        v-model="localValues[property['bk_property_id']]"
                                                        :disabled="checkIsFieldDisabled(property)">
                                                </span>
                                                <input type="text" class="bk-form-input" v-else-if="property['bk_property_type'] === 'int'"
                                                    :disabled="checkIsFieldDisabled(property)" maxlength="11" 
                                                    v-model.trim.number="localValues[property['bk_property_id']]">
                                                <input v-else
                                                    type="text" class="bk-form-input"
                                                    :disabled="checkIsFieldDisabled(property)" 
                                                    v-model.trim="localValues[property['bk_property_id']]">
                                                <template v-if="getValidateRules(property)">
                                                    <v-validate class="attribute-validate-result"
                                                        v-validate="getValidateRules(property)"
                                                        :name="property['bk_property_name']" 
                                                        :value="property['bk_property_type'] === 'bool' ? localValues[property['bk_property_id']] || false : localValues[property['bk_property_id']]">
                                                    </v-validate>
                                                </template>
                                            </div>
                                        </li>
                                    </template>
                                </ul>
                        </div>
                        <div class="attribute-group-more" v-if="propertyGroup === 'none'">
                            <a href="javascript:void(0)" class="group-more-link" :class="{'open': !isNoneGroupHide}" @click="isNoneGroupHide = !isNoneGroupHide">{{$t("Common['更多属性']")}}</a>
                        </div>
                    </template>
                </template>
            </form>
            <div v-else>
                <p class="attribute-no-multiple">{{$t("Common['无可编辑属性']")}}</p>
            </div>
        </template>
        <template v-if="showBtnGroup">
            <div class="attribute-btn-group" v-if="displayType==='list' && type === 'update'">
                <bk-button type="primary" class="bk-button main-btn" @click.prevent="changeDisplayType('form')" :disabled="unauthorized.update">{{$t("Common['属性编辑']")}}</bk-button>
                <bk-button type="default" :loading="$loading('instDelete')" v-if="type==='update' && showDelete && !isMultipleUpdate" class="del-btn" @click.prevent="deleteObject" :disabled="unauthorized.delete">
                    <template v-if="objId !== 'biz'">
                        {{$t("Common['删除']")}}
                    </template>
                    <template v-else>
                        {{$t("Inst['归档']")}}
                    </template>
                </bk-button>
            </div>
            <div class="attribute-btn-group" v-else-if="!isMultipleUpdate || isMultipleUpdate && hasEditableProperties">
                <bk-button type="primary" :loading="$loading('editAttr')" v-if="type==='create'" class="main-btn" @click.prevent="submit" :disabled="errors.any() || !Object.keys(formData).length || unauthorized.update">{{$t("Common['保存']")}}</bk-button>
                <bk-button type="primary" :loading="$loading('editAttr')" v-if="type==='update'" class="main-btn" @click.prevent="submit" :disabled="errors.any() || !Object.keys(formData).length || unauthorized.update">{{$t("Common['保存']")}}</bk-button>
                <bk-button type="default" v-if="type==='update'" class="vice-btn" @click.prevent="changeDisplayType('list')">{{$t("Common['取消']")}}</bk-button>
            </div>
        </template>
    </div>
</template>
<script>
    import vMemberSelector from '@/components/common/selector/member'
    import vEnumeration from '@/components/common/selector/enumeration'
    import vAssociation from '@/components/common/selector/association'
    import vHost from '@/components/common/hostAssociation/host'
    import vValidate from '@/components/common/validator/validate'
    import { mapGetters } from 'vuex'
    import Authority from '@/mixins/authority'
    const vTimezone = () => import('../timezone/timezone')
    export default {
        mixins: [Authority],
        props: {
            active: {               // 属性展示界面是否激活
                type: Boolean,
                default: false
            },
            type: {
                type: String,       // 属性界面是新增还是修改
                default: 'create'
            },
            formFields: {           // 表单字段, property为该在模型管理中配置的字段各种属性集合
                type: Array,
                default () {
                    return []
                }
            },
            formValues: {           // 当前表单的值，编辑属性时传入的对象，用于初始化编辑表单，key为formFields中的PropertyId
                type: Object,
                default () {
                    return {}
                }
            },
            showDelete: {           // 是否显示删除按钮,  批量修改时可禁用此按钮
                type: Boolean,
                default: true
            },
            showBtnGroup: {         // 是否显示按钮组
                type: Boolean,
                default: true
            },
            isMultipleUpdate: {
                type: Boolean,
                default: false
            },
            isBatchUpdate: {
                type: Boolean,
                default: true
            },
            objId: {
                type: String,
                default: ''
            }
        },
        data () {
            return {
                displayType: 'list', // list | form 当前属性界面为展示或者表单编辑态
                localValues: {},     // 当前表单对应的值
                multipleEditableFields: {},
                groupOrder: [],
                isNoneGroupHide: true,
                selectHostActive: ''
            }
        },
        computed: {
            ...mapGetters(['bkSupplierAccount']),
            localFormFields () {
                let formFields = this.$deepClone(this.formFields)
                formFields.sort((objA, objB) => {
                    return objA['bk_property_index'] - objB['bk_property_index']
                })
                return formFields
            },
            //  属性分组:根据formFields中各property.PropertyGroup 进行属性分组，为'隐藏分组(none)'时，放入更多属性中
            bkPropertyGroups () {
                let bkPropertyGroups = {}
                this.localFormFields.filter(property => !['singleasst', 'multiasst'].includes(property['bk_property_type'])).map(property => {
                    let {
                        bk_property_group: bkPropertyGroup,
                        bk_property_group_name: bkPropertyGroupName
                    } = property
                    if (bkPropertyGroup /* && bkPropertyGroup !== 'none'*/) {    // 暂要显示隐藏分组
                        if (bkPropertyGroups.hasOwnProperty(bkPropertyGroup)) {
                            bkPropertyGroups[bkPropertyGroup]['properties'].push(Object.assign({}, property))
                        } else {
                            bkPropertyGroups[bkPropertyGroup] = {
                                bkPropertyGroup,
                                bkPropertyGroupName,
                                properties: [Object.assign({}, property)]
                            }
                        }
                    }
                })
                return bkPropertyGroups
            },
            groupEditable () {
                let groupEditable = {}
                this.groupOrder.map(group => {
                    if (this.bkPropertyGroups.hasOwnProperty(group)) {
                        groupEditable[group] = this.bkPropertyGroups[group]['properties'].some(property => {
                            if (this.isMultipleUpdate) {
                                return property['editable'] && !property['bk_isapi'] && !property['isonly']
                            } else if (this.type === 'create') {
                                return !property['bk_isapi']
                            } else {
                                return property['editable'] && !property['bk_isapi']
                            }
                        })
                    }
                })
                return groupEditable
            },
            hasEditableProperties () {
                return this.groupOrder.some(group => {
                    return this.groupEditable[group]
                })
            },
            formData () {
                let formData = Object.assign({}, this.localValues)
                let propertyInitValue = {
                    'int': null,
                    'bool': false,
                    'date': null,
                    'time': null,
                    'enum': null,
                    'timezone': null
                }
                this.localFormFields.map(property => {
                    let {
                        bk_property_id: bkPropertyId,
                        bk_property_type: bkPropertyType,
                        bk_asst_obj_id: bkAsstObjId,
                        bk_isapi: bkIsapi,
                        editable
                    } = property
                    // 创建时将未编辑的值填充进表单
                    if (this.type === 'create' && !this.isMultipleUpdate && !formData.hasOwnProperty(bkPropertyId)) {
                        formData[bkPropertyId] = ''
                    }
                    // 删除接口自动生成的字段、更新时删除不可编辑的字段
                    if (bkIsapi || ((this.type === 'update' || this.isMultipleUpdate) && !editable)) {
                        delete formData[bkPropertyId]
                    }
                    if (propertyInitValue.hasOwnProperty(bkPropertyType) && (formData[bkPropertyId] === '')) {
                        formData[bkPropertyId] = propertyInitValue[bkPropertyType]
                    }
                    if (this.isMultipleUpdate && !this.multipleEditableFields[bkPropertyId]) {
                        delete formData[bkPropertyId]
                    }
                    if (Array.isArray(formData[bkPropertyId])) {
                        formData[bkPropertyId] = formData[bkPropertyId].filter(({id}) => !!id).map(({bk_inst_id: bkInstId}) => bkInstId).join(',')
                    }
                })
                // 增量更新，删除未变更的字段
                if (this.isBatchUpdate) {
                    let formDataCopy = Object.assign({}, formData)
                    for (let formDataKey in formDataCopy) {
                        for (let formValuesKey in this.formValues) {
                            if (formDataKey === formValuesKey) {
                                if (Array.isArray(this.formValues[formValuesKey])) {
                                    if (formDataCopy[formDataKey] === this.formValues[formValuesKey].map(({id, bk_inst_id: bkInstId}) => id === '' ? '' : bkInstId).join(',')) {
                                        delete formData[formDataKey]
                                    }
                                } else if (formDataCopy[formDataKey] === this.formValues[formValuesKey]) {
                                    delete formData[formDataKey]
                                } else if (formDataCopy[formDataKey] === '' && this.formValues[formValuesKey] === null) {
                                    delete formData[formDataKey]
                                }
                                break
                            }
                        }
                    }
                }
                
                return formData
            }
        },
        watch: {
            active (active) {
                this.isNoneGroupHide = true
                // 属性界面激活态切换时，未激活时：重置表单数据；激活时，设置属性展示类型
                if (!active) {
                    this.resetData()
                } else {
                    this.displayType = this.type === 'create' ? 'form' : 'list'
                }
            },
            formValues (values) {
                // 属性编辑
                this.setUpdateInitData()
            },
            displayType (displayType) {
                this.isNoneGroupHide = true
                if (displayType === 'list') {
                    this.localValues = {}
                } else if (this.type === 'update') {
                    this.setUpdateInitData()
                }
                if (displayType === 'form') {
                    this.$validator.validateAll().then(() => {
                        this.errors.clear()
                    })
                }
            },
            isMultipleUpdate (isMultipleUpdate) {
                if (isMultipleUpdate) {
                    this.localFormFields.map(property => {
                        let {
                            bk_property_type: bkPropertyType,
                            bk_property_id: bkPropertyId
                        } = property
                        if (bkPropertyType === 'bool') {
                            if (this.localValues.hasOwnProperty(bkPropertyType)) {
                                this.localValues[bkPropertyType] = false
                            } else {
                                this.$set(this.localValues, bkPropertyId, false)
                            }
                        }
                    })
                }
            },
            objId (objId) {
                if (objId) {
                    this.getbkPropertyGroups()
                }
            }
        },
        beforeMount () {
            if (this.objId) {
                this.getbkPropertyGroups()
            }
        },
        methods: {
            isCloseConfirmShow () {
                let isConfirmShow = false
                if (this.displayType === 'list') {
                    return false
                }
                if (this.type === 'create') {
                    for (let key in this.formData) {
                        let property = this.localFormFields.find(({bk_property_type: bkPropertyType, bk_property_id: bkPropertyId}) => {
                            return bkPropertyId === key
                        })
                        if (property['bk_property_type'] === 'enum') {
                            let enumProperty = property.option.find(({id}) => {
                                return id === this.formData[key]
                            })
                            if (enumProperty && !enumProperty['is_default']) {
                                isConfirmShow = true
                                break
                            }
                        } else {
                            if (this.formData[key] !== null && this.formData[key].length) {
                                isConfirmShow = true
                                break
                            }
                        }
                    }
                } else {
                    for (let key in this.formData) {
                        let property = this.localFormFields.find(({bk_property_type: bkPropertyType, bk_property_id: bkPropertyId}) => {
                            return bkPropertyId === key
                        })
                        let value = this.formValues[key]
                        if (property['bk_property_type'] === 'singleasst' || property['bk_property_type'] === 'multiasst') {
                            value = []
                            if (this.formValues.hasOwnProperty(key)) {
                                this.formValues[key].map(formValue => {
                                    value.push(formValue['bk_inst_id'])
                                })
                            }
                            value = value.join(',')
                        }
                        if (value !== this.formData[key] && !(this.formData[key] === '' && !this.formValues.hasOwnProperty(key))) {
                            isConfirmShow = true
                            break
                        }
                    }
                }
                return isConfirmShow
            },
            confirmHost (hostInfo) {
                this.hideSelectHost()
            },
            showSelectHost (propertyId) {
                this.selectHostActive = propertyId
            },
            hideSelectHost () {
                this.selectHostActive = ''
            },
            checkIsShowGroup (group) {
                return this.bkPropertyGroups.hasOwnProperty(group) && this.groupEditable[group]
            },
            checkIsShowField (property) {
                if (this.isMultipleUpdate) {
                    return property['editable'] && !property['bk_isapi'] && !property['isonly']
                } else if (this.type === 'create') {
                    return !property['bk_isapi']
                } else {
                    return property['editable'] && !property['bk_isapi']
                }
            },
            getbkPropertyGroups () {
                this.$axios.post(`objectatt/group/property/owner/${this.bkSupplierAccount}/object/${this.objId}`, {}).then(res => {
                    if (res.result) {
                        let groups = res.data.sort((groupA, groupB) => {
                            return groupA['bk_group_index'] - groupB['bk_group_index']
                        })
                        let groupOrder = groups.map(({bk_group_id: bkGroupId}) => {
                            return bkGroupId
                        })
                        this.groupOrder = [...groupOrder, 'none']
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            setUpdateInitData () {
                let filteredValues = this.filterValues()
                this.localValues = Object.assign({}, filteredValues)
            },
            changeDisplayType (type) {
                this.displayType = type
            },
            setDate (date, bkPropertyId) {
                if (this.localValues.hasOwnProperty(bkPropertyId)) {
                    this.localValues[bkPropertyId] = date
                } else {
                    this.$set(this.localValues, bkPropertyId, date)
                }
            },
            handleHostCancel (bkPropertyId) {
                if (this.localValues.hasOwnProperty(bkPropertyId)) {
                    this.localValues[bkPropertyId] = this.formValues[bkPropertyId]
                }
            },
            // 获取属性列表显示的值
            getFieldValue (property) {
                let {
                    bk_property_id: bkPropertyId,
                    bk_property_type: bkPropertyType
                } = property
                let value = this.formValues[bkPropertyId]
                if (value !== undefined) {
                    if (property['bk_asst_obj_id']) {
                        let associateName = []
                        if (Array.isArray(value)) {
                            value.map(({bk_inst_name: bkInstName}) => {
                                if (bkInstName) {
                                    associateName.push(bkInstName)
                                }
                            })
                        }
                        return associateName.join(',')
                    } else if (bkPropertyType === 'date') {
                        return this.$formatTime(value, 'YYYY-MM-DD')
                    } else if (bkPropertyType === 'time') {
                        return this.$formatTime(value)
                    } else if (bkPropertyType === 'enum') {
                        let obj = property.option.find(({id}) => {
                            return id === value
                        })
                        if (obj) {
                            return obj.name
                        } else {
                            return ''
                        }
                    } else {
                        return value
                    }
                }
            },
            // 判断是否可编辑
            checkIsFieldDisabled (property) {
                let {
                    bk_property_id: bkPropertyId,
                    bk_property_type: bkPropertyType,
                    editable
                } = property
                if (bkPropertyId === 'bk_biz_name' && this.formValues[bkPropertyId] === '蓝鲸') {
                    return true
                } else if (this.isMultipleUpdate) {
                    return !this.multipleEditableFields[bkPropertyId]
                } else if (this.type === 'create') {
                    return false
                } else {
                    return !editable
                }
            },
            // 批量修改，取消勾选时情况对应表单的值
            clearFieldValue (property) {
                let {
                    bk_property_id: bkPropertyId
                } = property
                if (!this.multipleEditableFields[bkPropertyId]) {
                    delete this.localValues[bkPropertyId]
                }
            },
            // 过滤掉编辑态传进来的不在表单项中的属性值
            filterValues () {
                let filteredValues = {}
                Object.keys(this.formValues).map(formPropertyId => {
                    let fieldProperty = this.localFormFields.find(property => {
                        return formPropertyId === property['bk_property_id']
                    })
                    if (fieldProperty) {
                        let {
                            bk_property_id: bkPropertyId,
                            bk_property_type: bkPropertyType
                        } = fieldProperty
                        if (bkPropertyType === 'date') {
                            filteredValues[formPropertyId] = this.$formatTime(this.formValues[formPropertyId], 'YYYY-MM-DD')
                        } else if (bkPropertyType === 'time') {
                            filteredValues[formPropertyId] = this.$formatTime(this.formValues[formPropertyId])
                        } else if (bkPropertyType === 'enum' && this.formValues[formPropertyId] === null) {
                            filteredValues[formPropertyId] = ''
                        } else {
                            filteredValues[formPropertyId] = this.formValues[formPropertyId]
                        }
                    }
                })
                return filteredValues
            },
            resetData () {
                this.displayType = 'list'
                this.localValues = {}
                this.multipleEditableFields = {}
                this.$forceUpdate()
            },
            getValidateRules (property) {
                let rules = {}
                let {
                    bk_property_type: bkPropertyType,
                    option,
                    isrequired
                } = property
                if (isrequired && !this.isMultipleUpdate) {
                    rules['required'] = true
                }
                if (property.hasOwnProperty('option') && option) {
                    if (bkPropertyType === 'int') {
                        if (option.hasOwnProperty('min') && option.min) {
                            rules['min_value'] = option.min
                        }
                        if (option.hasOwnProperty('max') && option.max) {
                            rules['max_value'] = option.max
                        }
                    } else if ((bkPropertyType === 'singlechar' || bkPropertyType === 'longchar') && option !== null) {
                        rules['regex'] = option
                    }
                }
                if (bkPropertyType === 'singlechar') {
                    rules['singlechar'] = true
                }
                if (bkPropertyType === 'longchar') {
                    rules['longchar'] = true
                }
                if (bkPropertyType === 'int') {
                    rules['regex'] = '^(0|[1-9][0-9]*|-[1-9][0-9]*)$'
                }
                return rules
            },
            submit () {
                this.$validator.validateAll().then(isValid => {
                    if (isValid) {
                        if (Object.keys(this.formData).length) {
                            this.$emit('submit', this.formData, Object.assign({}, this.formValues))
                        }
                    }
                })
            },
            deleteObject () {
                this.$emit('delete', Object.assign({}, this.formValues))
            },
            htmlEncode (str) {
                let c = document.createElement('div')
                c.innerHTML = str
                let output = c.innerText
                c = null
                return output
            }
        },
        components: {
            vMemberSelector,
            vTimezone,
            vValidate,
            vHost,
            vEnumeration,
            vAssociation
        }
    }
</script>
<style lang="scss" scoped>
    .attribute-wrapper{
        padding: 0 0 0 32px;
        overflow: visible;
    }
    .attribute-list{
        padding: 4px 0;
        .attribute-item{
            width: 50%;
            font-size: 12px;
            line-height: 16px;
            margin: 12px 0 0 0;
            white-space: nowrap;
            .attribute-item-label{
                width: 116px;
                color: #737987;
                // color: #6b7baa;
                text-align: right;
                display: inline-block;
                overflow: hidden;
                text-overflow: ellipsis;
                margin-right: 10px;
                padding-right: 6px;
                position: relative;
                &.has-colon:after{
                    content: ":";
                    position: absolute;
                    right: 0;
                    top: 0;
                    line-height: 14px;
                }
            }
            .attribute-item-value{
                max-width: calc(100% - 130px);
                display: inline-block;
                overflow: hidden;
                text-overflow: ellipsis;
                color: #333948;
            }
            .attribute-item-value.bk-form-checkbox{
                padding: 0;
                font-size: 0;
                transform: scale(0.889);
                vertical-align: -1px;
                vertical-align: top;
                input[type="checkbox"]{
                    &:checked{
                        background-position: -33px -62px;
                    }
                }
            }
        }
    }
    .attribute-list.edit{
        .attribute-item{
            margin: 0 0 11px 0;
            min-height: 63px;
            label{
                display: inline-block;
                vertical-align: middle;
                margin: 6px 0 9px 0;
                padding: 0;
                color: #737987;
                line-height: 12px;
                font-size: 12px;
                overflow: visible;
                &.required:after{
                    content: '*';
                    color: #ff5656;
                }
                input[type="checkbox"]{
                    transform: scale(0.857);
                }
            }
            .icon-tooltips{
                display: inline-block;
                vertical-align: middle;
                width: 16px;
                height: 16px;
                margin: -2px 0 0 10px;
                background: url('../../common/images/icon/icon-info-tips.png') no-repeat;
            }
            .attribute-item-field{
                width: 310px;
                white-space: normal;
                position: relative;
                color: #333948;
                .bk-date{
                    width: 100%;
                }
                .attribute-validate-result{
                    position: absolute;
                    top: 100%;
                    left : 0;
                }
                .bk-form-input[readonly]{
                    background: #fff;
                }
            }
        }
        .attribute-item.bool{
            label{
                display: inline-block;
            }
            .attribute-item-field{
                display: inline-block;
                height: 36px;
            }
        }
    }
    .attribute-group{
        padding: 17px 0 0 0;
        &:first-child{
            padding: 28px 0 0 0;
        }
        .title{
            margin: 0;
            font-size: 14px;
            line-height: 14px;
            overflow: visible;
            color: #333948;
        }
    }
    .attribute-group-more{
        text-align: center;
        margin-top: 26px;
        .group-more-link{
            color: #6b7baa;
            text-decoration: none;
            font-size: 12px;
            &.open:after{
                transform: rotate(0deg);
            }
            &.open:hover:after{
                transform: rotate(180deg);
            }
            &:hover{
                color: #498fe0;
            }
            &:hover:after{
                background-image: url('../../common/images/icon/icon-result-slide-hover.png');
                transform: rotate(0deg);
            }
            &:after{
                content: '';
                display: inline-block;
                width: 11px;
                height: 10px;
                margin-left: 12px;
                background: url('../../common/images/icon/icon-result-slide.png') no-repeat;
                transform: rotate(180deg);
            }
        }
    }
    .attribute-btn-group{
        position: absolute;
        bottom: 0;
        right: 0;
        height: 62px;
        width: 100%;
        background-color: #f9f9f9;
        padding: 14px 20px;
        button{
            width: 110px;
            height: 34px;
            line-height: 32px;
            margin-right: 6px;
        }
    }
    .attribute-no-multiple{
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate3D(-50%, 0, 0);
    }
</style>