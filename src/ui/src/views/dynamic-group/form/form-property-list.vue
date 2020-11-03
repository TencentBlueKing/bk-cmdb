<template>
    <div class="form-property-list">
        <bk-form form-type="vertical" :label-width="400">
            <bk-form-item
                v-for="property in properties"
                :key="property.id"
                :label="getPropertyLabel(property)">
                <div class="form-property-item">
                    <form-operator-selector class="item-operator"
                        v-if="!withoutOperator.includes(property.bk_property_type)"
                        :type="property.bk_property_type"
                        v-model="condition[property.id].operator"
                        @change="handleOperatorChange(property, ...arguments)">
                    </form-operator-selector>
                    <component class="item-value"
                        :is="`cmdb-search-${property.bk_property_type}`"
                        :placeholder="getPlaceholder(property)"
                        :data-vv-name="property.bk_property_id"
                        :data-vv-as="property.bk_property_name"
                        v-bind="getBindProps(property)"
                        v-model="condition[property.id].value"
                        v-validate="'required'">
                    </component>
                    <i class="item-remove bk-icon icon-close" @click="handleRemove(property)"></i>
                </div>
                <p class="form-error" v-if="errors.has(property.bk_property_id)">{{errors.first(property.bk_property_id)}}</p>
            </bk-form-item>
        </bk-form>
    </div>
</template>

<script>
    import FormOperatorSelector from './form-operator-selector.vue'
    export default {
        components: {
            FormOperatorSelector
        },
        inject: ['dynamicGroupForm'],
        data () {
            return {
                condition: {},
                withoutOperator: ['date', 'time', 'bool', 'service-template']
            }
        },
        computed: {
            properties () {
                return this.dynamicGroupForm.selectedProperties
            },
            availableModels () {
                return this.dynamicGroupForm.availableModels
            },
            modelMap () {
                const modelMap = {}
                this.availableModels.forEach(model => {
                    modelMap[model.bk_obj_id] = model
                })
                return modelMap
            },
            details () {
                return this.dynamicGroupForm.details
            }
        },
        watch: {
            properties: {
                immediate: true,
                handler () {
                    this.UpdateCondition()
                }
            }
        },
        methods: {
            getDefaultData (property) {
                const defaultMap = {
                    bool: {
                        operator: '$eq',
                        value: ''
                    },
                    date: {
                        operator: '$range',
                        value: []
                    },
                    float: {
                        operator: '$eq',
                        value: ''
                    },
                    int: {
                        operator: '$eq',
                        value: ''
                    },
                    time: {
                        operator: '$range',
                        value: []
                    },
                    'service-template': {
                        operator: '$in',
                        value: []
                    }
                }
                return {
                    operator: '$in',
                    value: [],
                    ...defaultMap[property.bk_property_type]
                }
            },
            setDetailsCondition () {
                Object.values(this.condition).forEach(condition => {
                    const modelId = condition.property.bk_obj_id
                    const propertyId = condition.property.bk_property_id
                    const detailsCondition = this.details.info.condition.find(detailsCondition => detailsCondition.bk_obj_id === modelId)
                    const detailsFieldData = detailsCondition.condition.find(data => data.field === propertyId)
                    condition.operator = detailsFieldData.operator
                    condition.value = detailsFieldData.value
                })
            },
            UpdateCondition () {
                const newConditon = {}
                this.properties.forEach(property => {
                    if (this.condition.hasOwnProperty(property.id)) {
                        newConditon[property.id] = this.condition[property.id]
                    } else {
                        newConditon[property.id] = {
                            property: property,
                            ...this.getDefaultData(property)
                        }
                    }
                })
                this.condition = newConditon
            },
            handleOperatorChange (property, operator) {
                if (operator === '$range') {
                    this.condition[property.id].value = []
                } else {
                    const defaultValue = this.getDefaultData(property).value
                    const currentValue = this.condition[property.id].value
                    const isTypeChanged = (Array.isArray(defaultValue)) !== (Array.isArray(currentValue))
                    this.condition[property.id].value = isTypeChanged ? defaultValue : currentValue
                }
            },
            handleRemove (property) {
                this.$emit('remove', property)
            },
            getBindProps (property) {
                if (['list', 'enum'].includes(property.bk_property_type)) {
                    return {
                        options: property.option || []
                    }
                }
                return {}
            },
            getPropertyLabel (property) {
                const modelId = property.bk_obj_id
                const propertyName = property.bk_property_name
                const modelName = this.modelMap[modelId].bk_obj_name
                return `${modelName} - ${propertyName}`
            },
            getPlaceholder (property) {
                const selectTypes = ['list', 'enum', 'timezone', 'organization', 'date', 'time']
                if (selectTypes.includes(property.bk_property_type)) {
                    return this.$t('请选择xx', { name: property.bk_property_name })
                }
                return this.$t('请输入xx', { name: property.bk_property_name })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-property-list {
        .form-property-item {
            display: flex;
            align-items: center;
            &:hover {
                .item-remove {
                    visibility: visible;
                }
            }
            .item-operator {
                flex: 110px 0 0;
                margin-right: 10px;
            }
            .item-value {
                flex: 1;
                margin: 0 10px 0 0;
            }
            .item-remove {
                font-size: 20px;
                visibility: hidden;
                cursor: pointer;
            }
        }
        .form-error {
            position: absolute;
            top: 100%;
            font-size: 12px;
            line-height: 14px;
            color: $dangerColor;
        }
    }
</style>
