<template>
    <div class="property-filter clearfix">
        <cmdb-selector class="property-selector fl" style="width: 135px;"
            :list="filteredProperties"
            setting-key="bk_property_id"
            display-key="bk_property_name"
            v-model="localSelected.id"
            @on-selected="handlePropertySelected">
        </cmdb-selector>
        <cmdb-selector class="operator-selector fl" style="width: 135px;"
            :list="operatorOptions"
            setting-key="value"
            display-key="label"
            v-model="localSelected.operator"
            @on-selected="handleOperatorSelected">
        </cmdb-selector>
        <div class="property-value fl" style="width: 245px;"
            v-if="Object.keys(selectedProperty).length">
            <cmdb-form-enum
                v-if="selectedProperty['bk_property_type'] === 'enum'"
                :options="selectedProperty.option || []"
                v-model="localSelected.value">
            </cmdb-form-enum>
            <component
                v-else
                :is="`cmdb-form-${selectedProperty['bk_property_type']}`"
                v-model.trim="localSelected.value">
            </component>
        </div>
    </div>
</template>
<script>
    import { mapGetters, mapActions } from 'vuex'
    export default {
        props: {
            objId: {
                type: String,
                required: true
            },
            excludeType: {
                type: Array,
                default () {
                    return []
                }
            },
            excludeId: {
                type: Array,
                default () {
                    return []
                }
            }
        },
        data () {
            return {
                localSelected: {
                    id: '',
                    operator: '',
                    value: ''
                },
                filteredProperties: [],
                propertyOperator: {
                    'default': ['$eq', '$ne'],
                    'singlechar': ['$regex', '$eq', '$ne'],
                    'longchar': ['$regex', '$eq', '$ne'],
                    'objuser': ['$regex', '$eq', '$ne'],
                    'singleasst': ['$regex', '$eq', '$ne'],
                    'multiasst': ['$regex', '$eq', '$ne']
                },
                operatorLabel: {
                    '$nin': this.$t('不包含'),
                    '$in': this.$t('包含'),
                    '$regex': this.$t('包含'),
                    '$eq': this.$t('等于'),
                    '$ne': this.$t('不等于')
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            selectedProperty () {
                return this.filteredProperties.find(({ bk_property_id: bkPropertyId }) => bkPropertyId === this.localSelected.id) || {}
            },
            operatorOptions () {
                if (this.selectedProperty) {
                    if (['bk_host_innerip', 'bk_host_outerip'].includes(this.selectedProperty['bk_property_id']) || this.objId === 'biz') {
                        return [{ label: this.operatorLabel['$regex'], value: '$regex' }]
                    } else {
                        const propertyType = this.selectedProperty['bk_property_type']
                        const propertyOperator = this.propertyOperator.hasOwnProperty(propertyType) ? this.propertyOperator[propertyType] : this.propertyOperator['default']
                        return propertyOperator.map(operator => {
                            return {
                                label: this.operatorLabel[operator],
                                value: operator
                            }
                        })
                    }
                }
                return []
            }
        },
        watch: {
            filteredProperties (properties) {
                if (properties.length) {
                    this.localSelected.id = properties[0]['bk_property_id']
                    this.$emit('on-property-selected', properties[0]['bk_property_id'], properties[0])
                } else {
                    this.localSelected.id = ''
                    this.$emit('on-property-selected', '', null)
                }
            },
            operatorOptions (operatorOptions) {
                this.localSelected.operator = operatorOptions.length ? operatorOptions[0]['value'] : ''
                this.$emit('handleOperatorSelected', this.localSelected.operator)
            },
            'localSelected.id' (id) {
                this.localSelected.value = ''
            },
            'localSelected.value' (value) {
                this.$emit('on-value-change', value)
            },
            async objId (objId) {
                const properties = await this.searchObjectAttribute({
                    params: this.$injectMetadata({
                        'bk_obj_id': objId,
                        'bk_supplier_account': this.supplierAccount
                    }),
                    config: {
                        requestId: `post_searchObjectAttribute_${objId}`
                    }
                })
                this.filteredProperties = properties.filter(property => {
                    const {
                        bk_isapi: bkIsapi,
                        bk_property_type: bkPropertyType,
                        bk_property_id: bkPropertyId
                    } = property
                    return !bkIsapi && !this.excludeType.includes(bkPropertyType) && !this.excludeId.includes(bkPropertyId)
                })
            }
        },
        methods: {
            ...mapActions('objectModelProperty', ['searchObjectAttribute']),
            handlePropertySelected (value, data) {
                this.$emit('on-property-selected', value, data)
            },
            handleOperatorSelected (value, data) {
                this.$emit('on-operator-selected', value, data)
            }
        }
    }
</script>
<style lang="scss" scoped>
    .property-selector{
        width: 135px;
    }
    .operator-selector{
        width: 135px;
        margin: 0 10px;
    }
    .property-value{
        width: 245px;
    }
</style>
