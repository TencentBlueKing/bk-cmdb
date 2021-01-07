<template>
    <bk-sideslider
        :width="515"
        :title="title"
        :is-show.sync="isShow"
        :before-close="beforeClose"
        @hidden="handleHidden">
        <bk-form slot="content"
            class="dynamic-group-form"
            form-type="vertical"
            v-bkloading="{ isLoading: $loading() }">
            <bk-form-item :label="$t('业务')" required>
                <cmdb-business-selector class="form-item"
                    disabled
                    :value="bizId">
                </cmdb-business-selector>
            </bk-form-item>
            <bk-form-item :label="$t('分组名称')" required>
                <bk-input class="form-item"
                    v-model.trim="formData.name"
                    v-validate="'required|length:256'"
                    data-vv-name="name"
                    :data-vv-as="$t('查询名称')"
                    :placeholder="$t('请输入xx', { name: $t('查询名称') })">
                </bk-input>
                <p class="form-error" v-if="errors.has('name')">{{errors.first('name')}}</p>
            </bk-form-item>
            <bk-form-item :label="$t('查询对象')" required>
                <bk-select class="form-item"
                    v-model="formData.bk_obj_id"
                    :clearable="false"
                    :disabled="!isCreateMode"
                    @change="handleModelChange">
                    <bk-option v-for="model in searchTargetModels"
                        :id="model.bk_obj_id"
                        :name="model.bk_obj_name"
                        :key="model.bk_obj_id">
                    </bk-option>
                </bk-select>
            </bk-form-item>
            <bk-form-item class="form-condition-tips"
                desc-type="icon"
                desc-icon="icon-cc-tips"
                :label="$t('查询条件')"
                :desc="$t('针对查询内容进行条件过滤')">
                <form-property-list ref="propertyList" @remove="handleRemoveProperty"></form-property-list>
                <bk-button class="form-condition-button" :style="{ marginTop: selectedProperties.length ? '10px' : 0 }"
                    icon="icon-plus-circle"
                    :text="true"
                    @click="showPropertySelector">
                    {{$t('继续添加')}}
                </bk-button>
                <input type="hidden"
                    v-validate="'min_value:1'"
                    data-vv-name="condition"
                    :data-vv-as="$t('查询条件')"
                    v-model="selectedProperties.length">
                <p class="form-error" v-if="errors.has('condition')">{{$t('请添加查询条件')}}</p>
            </bk-form-item>
        </bk-form>
        <div class="dynamic-group-options" slot="footer">
            <cmdb-auth :auth="saveAuth">
                <bk-button class="mr10" slot-scope="{ disabled }"
                    theme="primary"
                    :disabled="disabled"
                    @click="handleConfirm">
                    {{isCreateMode ? $t('提交') : $t('保存')}}
                </bk-button>
            </cmdb-auth>
            <bk-button class="mr10" theme="default" @click="close">{{$t('取消')}}</bk-button>
        </div>
    </bk-sideslider>
</template>

<script>
    import { mapGetters } from 'vuex'
    import FormPropertyList from './form-property-list.vue'
    import FormPropertySelector from './form-property-selector.js'
    import RouterQuery from '@/router/query'
    export default {
        components: {
            FormPropertyList
        },
        props: {
            id: [String, Number],
            title: String
        },
        provide () {
            return {
                dynamicGroupForm: this
            }
        },
        data () {
            return {
                isShow: false,
                details: null,
                formData: {
                    name: '',
                    bk_obj_id: 'host'
                },
                selectedProperties: [],
                request: Object.freeze({
                    mainline: Symbol('mainline'),
                    property: Symbol('property')
                }),
                availableModelIds: Object.freeze(['host', 'module', 'set']),
                availableModels: [],
                propertyMap: {}
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('objectBiz', ['bizId']),
            isCreateMode () {
                return !this.id
            },
            searchTargetModels () {
                return this.availableModels.filter(model => ['host', 'set'].includes(model.bk_obj_id))
            },
            saveAuth () {
                if (this.id) {
                    return { type: this.$OPERATION.U_CUSTOM_QUERY, relation: [this.bizId, this.id] }
                }
                return { type: this.$OPERATION.C_CUSTOM_QUERY, relation: [this.bizId] }
            }
        },
        async created () {
            await this.getMainLineModels()
            await this.getModelProperties()
            if (this.id) {
                this.getDetails()
            }
        },
        methods: {
            async getMainLineModels () {
                try {
                    const models = await this.$store.dispatch('objectMainLineModule/searchMainlineObject', {
                        config: {
                            requestId: this.request.mainline,
                            fromCache: true
                        }
                    })
                    // 业务调用方暂时只需要一下三种类型的查询
                    const availableModels = this.availableModelIds.map(modelId => models.find(model => model.bk_obj_id === modelId))
                    this.availableModels = Object.freeze(availableModels)
                } catch (error) {
                    console.error(error)
                }
            },
            async getModelProperties () {
                try {
                    const propertyMap = await this.$store.dispatch('objectModelProperty/batchSearchObjectAttribute', {
                        params: {
                            bk_biz_id: this.bizId,
                            bk_obj_id: { $in: this.availableModels.map(model => model.bk_obj_id) },
                            bk_supplier_account: this.supplierAccount
                        },
                        config: {
                            requestId: this.request.property,
                            fromCache: true
                        }
                    })
                    propertyMap.module.unshift(this.getServiceTemplateProperty())
                    this.propertyMap = Object.freeze(propertyMap)
                } catch (error) {
                    console.error(error)
                    this.propertyMap = {}
                }
            },
            getServiceTemplateProperty () {
                return {
                    id: Date.now(),
                    bk_obj_id: 'module',
                    bk_property_id: 'service_template_id',
                    bk_property_name: this.$t('服务模板'),
                    bk_property_index: -1,
                    bk_property_type: 'service-template',
                    isonly: true,
                    ispre: true,
                    bk_isapi: true,
                    bk_issystem: true,
                    isreadonly: true,
                    editable: false,
                    bk_property_group: null,
                    _is_inject_: true
                }
            },
            async getDetails () {
                try {
                    const details = await this.$store.dispatch('dynamicGroup/details', {
                        bizId: this.bizId,
                        id: this.id
                    })
                    const transformedDetails = this.transformDetails(details)
                    this.formData.name = transformedDetails.name
                    this.formData.bk_obj_id = transformedDetails.bk_obj_id
                    this.details = transformedDetails
                    this.$nextTick(this.setDetailsSelectedProperties)
                    setTimeout(this.$refs.propertyList.setDetailsCondition, 0)
                } catch (error) {
                    console.error(error)
                }
            },
            transformDetails (details) {
                const condition = details.info.condition
                const transformedCondition = []
                condition.forEach(data => {
                    transformedCondition.push({
                        bk_obj_id: data.bk_obj_id,
                        condition: data.condition.reduce((accumulator, current) => {
                            if (['$gte', '$lte'].includes(current.operator)) {
                                // 将相同字段的$gte/$lte两个条件合并为一个range条件，用于表单组件渲染
                                let index = accumulator.findIndex(exist => exist.field === current.field)
                                if (index === -1) {
                                    index = accumulator.push({
                                        field: current.field,
                                        operator: '$range',
                                        value: []
                                    }) - 1
                                }
                                const range = accumulator[index]
                                if (current.operator === '$gte') {
                                    range.value.unshift(current.value)
                                } else {
                                    range.value.push(current.value)
                                }
                            } else if (current.operator === '$eq') {
                                // 将老数据的eq转换为当前支持的数据格式
                                const transformType = ['singlechar', 'longchar', 'enum']
                                const property = this.getConditionProperty(data.bk_obj_id, current.field)
                                if (property && transformType.includes(property.bk_property_type)) {
                                    accumulator.push({
                                        field: current.field,
                                        operator: '$in',
                                        value: Array.isArray(current.value) ? current.value : [current.value]
                                    })
                                } else {
                                    accumulator.push(current)
                                }
                            } else {
                                accumulator.push(current)
                            }
                            return accumulator
                        }, [])
                    })
                })
                return {
                    ...details,
                    info: {
                        condition: transformedCondition
                    }
                }
            },
            getConditionProperty (modelId, field) {
                const properties = this.propertyMap[modelId] || []
                return properties.find(property => property.bk_property_id === field)
            },
            setDetailsSelectedProperties () {
                const conditions = this.details.info.condition
                const properties = []
                conditions.forEach(({ bk_obj_id: modelId, condition }) => {
                    condition.forEach(({ field }) => {
                        const property = this.propertyMap[modelId].find(property => property.bk_property_id === field)
                        property && properties.push(property)
                    })
                })
                this.selectedProperties = properties
            },
            handleModelChange () {
                this.selectedProperties = []
            },
            showPropertySelector () {
                FormPropertySelector.show({
                    selected: this.selectedProperties,
                    handler: this.handlePropertySelected
                }, this)
            },
            handlePropertySelected (selected) {
                this.selectedProperties = selected
            },
            handleRemoveProperty (property) {
                const index = this.selectedProperties.findIndex(target => target.id === property.id)
                if (index > -1) {
                    this.selectedProperties.splice(index, 1)
                }
            },
            async handleConfirm () {
                try {
                    const results = [
                        await this.$validator.validateAll(),
                        await this.$refs.propertyList.$validator.validateAll()
                    ]
                    if (results.some(isValid => !isValid)) {
                        return false
                    }
                    if (this.id) {
                        await this.updateDynamicGroup()
                    } else {
                        await this.createDynamicGroup()
                    }
                    this.close()
                } catch (error) {
                    console.error(error)
                }
            },
            updateDynamicGroup () {
                return this.$store.dispatch('dynamicGroup/update', {
                    bizId: this.bizId,
                    id: this.id,
                    params: {
                        bk_biz_id: this.bizId,
                        bk_obj_id: this.formData.bk_obj_id,
                        name: this.formData.name,
                        info: {
                            condition: this.getSubmitCondition()
                        }
                    }
                })
            },
            createDynamicGroup () {
                return this.$store.dispatch('dynamicGroup/create', {
                    params: {
                        bk_biz_id: this.bizId,
                        bk_obj_id: this.formData.bk_obj_id,
                        name: this.formData.name,
                        info: {
                            condition: this.getSubmitCondition()
                        }
                    }
                })
            },
            getSubmitCondition () {
                const submitConditionMap = {}
                const propertyCondition = this.$refs.propertyList.condition
                Object.values(propertyCondition).forEach(({ property, operator, value }) => {
                    const submitCondition = submitConditionMap[property.bk_obj_id] || []
                    if (operator === '$range') {
                        const [start, end] = value
                        submitCondition.push({
                            field: property.bk_property_id,
                            operator: '$gte',
                            value: start
                        }, {
                            field: property.bk_property_id,
                            operator: '$lte',
                            value: end
                        })
                    } else {
                        submitCondition.push({
                            field: property.bk_property_id,
                            operator,
                            value
                        })
                    }
                    submitConditionMap[property.bk_obj_id] = submitCondition
                })
                return Object.keys(submitConditionMap).map(modelId => {
                    return {
                        bk_obj_id: modelId,
                        condition: submitConditionMap[modelId]
                    }
                })
            },
            close () {
                this.isShow = false
                RouterQuery.set({
                    _t: Date.now()
                })
            },
            show () {
                this.isShow = true
            },
            beforeClose () {
                return true
            },
            handleHidden () {
                this.$emit('close')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .dynamic-group-form {
        padding: 20px;
        height: 100%;
        @include scrollbar-y;
        .form-item {
            width: 100%;
        }
        .form-error {
            position: absolute;
            top: 100%;
            font-size: 12px;
            line-height: 14px;
            color: $dangerColor;
        }
        .form-condition-tips {
            /deep/ .bk-label .bk-label-text{
                display: inline-flex;
                align-items: center;
                .bk-icon {
                    margin-left: 4px;
                }
            }
        }
        .form-condition-button {
            /deep/ > div {
                display: flex;
                align-items: center;
                .bk-icon {
                    top: 0;
                }
            }
        }
    }
    .dynamic-group-options {
        display: flex;
        align-items: center;
        height: 100%;
        width: 100%;
        padding: 0 20px;
        border-top: 1px solid $borderColor;
    }
</style>
