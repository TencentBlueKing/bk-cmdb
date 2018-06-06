<template>
    <div class="property-filter clearfix">
        <bk-select class="property-selector fl" :selected.sync="localSelected.id" @on-selected="handlePropertySelected">
            <bk-select-option v-for="(property, index) in filteredProperties"
                :key="property['bk_property_id']"
                :value="property['bk_property_id']"
                :label="property['bk_property_name']">
            </bk-select-option>
        </bk-select>
        <bk-select class="operator-selector fl" :selected.sync="localSelected.operator" @on-selected="handleOperatorSelected">
            <bk-select-option v-for="(operator, index) in operatorOptions"
                :key="operator.value"
                :value="operator.value"
                :label="operator.label">
            </bk-select-option>
        </bk-select>
        <div class="property-value fl">
            <input type="text" class="bk-form-input" maxlength="11" v-model.number="localSelected.value" v-if="selectedProperty['bk_property_type'] === 'int'">
            <bk-select :selected.sync="localSelected.value" :showClear="true" v-else-if="selectedProperty['bk_property_type'] === 'enum'">
                <bk-select-option v-for="(enumOption, index) in (Array.isArray(selectedProperty.option) ? selectedProperty.option : [])"
                    :key="enumOption.id"
                    :value="enumOption.id"
                    :label="enumOption.name">
                </bk-select-option>
            </bk-select>
            <input type="text" class="bk-form-input" v-model.trim="localSelected.value" v-else>
        </div>
    </div>
</template>
<script>
    import { mapGetters } from 'vuex'
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
                    '$nin': this.$t("Common['不包含']"),
                    '$in': this.$t("Common['包含']"),
                    '$regex': this.$t("Common['包含']"),
                    '$eq': this.$t("Common['等于']"),
                    '$ne': this.$t("Common['不等于']")
                }
            }
        },
        computed: {
            ...mapGetters('object', ['attribute']),
            selectedProperty () {
                return this.filteredProperties.find(({bk_property_id: bkPropertyId}) => bkPropertyId === this.localSelected.id) || {}
            },
            operatorOptions () {
                if (this.selectedProperty) {
                    if (['bk_host_innerip', 'bk_host_outerip'].includes(this.selectedProperty['bk_property_id']) || this.objId === 'biz') {
                        return [{label: this.operatorLabel['$regex'], value: '$regex'}]
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
                    this.$emit('handlePropertySelected', {value: this.localSelected.id, label: properties[0]['bk_property_name']}, 0, 'length')
                } else {
                    this.localSelected.id = ''
                    this.$emit('handlePropertySelected', {value: '', label: ''}, -1, 'nolength')
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
                this.$emit('handleValueChange', value)
            },
            async objId (objId) {
                await this.$store.dispatch('object/getAttribute', {objId})
                this.filteredProperties = this.attribute[objId].filter(property => {
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
            handlePropertySelected (data, index) {
                this.$emit('handlePropertySelected', data, index, 'select')
            },
            handleOperatorSelected (data, index) {
                this.$emit('handleOperatorSelected', this.localSelected.operator)
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