<template>
    <div class="table-field">
        <v-table v-if="type === 'list'"
            :width="width"
            :emptyHeight="40"
            :header="header"
            :list="localValue"
            :valueKey="'list_header_name'"
            :labelKey="'list_header_describe'"
            :sortable="false">
        </v-table>
        <v-table v-else-if="type === 'form'"
            :width="width"
            :emptyHeight="40"
            :colBorder="true"
            :rowBorder="true"
            :header="header"
            :list="localValue"
            :valueKey="'list_header_name'"
            :labelKey="'list_header_describe'"
            :sortable="false">
            <template v-for="(head, index) in header" :slot="head['list_header_name']" slot-scope="{ item, rowIndex }">
                <div v-if="head['list_header_name'] !== operationId">
                    <input class="bk-form-input" type="text"
                        :value="item[head['list_header_name']]"
                        @change="updateValue($event, head, rowIndex)">
                </div>
                <div v-else class="field-operation">
                    <i class="bk-icon icon-plus-circle-shape" @click="addRow(rowIndex)"></i>
                    <i class="icon-cc-del" v-if="rowIndex > 0" @click="deleteRow(rowIndex)"></i>
                </div>
            </template>
        </v-table>
    </div>
</template>

<script>
    import vTable from '@/components/table/table'
    export default {
        components: {
            vTable
        },
        props: {
            property: {
                type: Object,
                required: true,
                validator (property) {
                    return Array.isArray(property.option)
                }
            },
            value: {
                required: true,
                validator (value) {
                    return Array.isArray(value) || !value
                }
            },
            type: {
                type: String,
                default: 'list', // list | form
                validator (type) {
                    return ['list', 'form'].includes(type)
                }
            },
            width: {
                type: Number,
                default: 675
            }
        },
        data () {
            return {
                localValue: [],
                operationId: Symbol('operation id')
            }
        },
        computed: {
            header () {
                let header = [...this.property.option]
                if (this.type === 'form') {
                    header.push({
                        width: 100,
                        'list_header_name': this.operationId,
                        'list_header_describe': ''
                    })
                }
                return header
            }
        },
        watch: {
            value (value) {
                this.setLocalValue()
            },
            type (type) {
                this.setLocalValue()
            }
        },
        created () {
            this.setLocalValue()
        },
        methods: {
            setLocalValue () {
                let value = this.value || []
                let localValue = this.$deepClone(value)
                if (this.type === 'form' && !localValue.length) {
                    localValue.push({})
                }
                this.localValue = localValue
            },
            updateValue (event, head, rowIndex) {
                this.localValue[rowIndex][head['list_header_name']] = event.target.value
                this.$emit('update:value', this.$deepClone(this.localValue))
            },
            addRow (rowIndex) {
                this.localValue.splice(rowIndex, 0, {})
                this.$emit('update:value', this.$deepClone(this.localValue))
            },
            deleteRow (rowIndex) {
                this.localValue.splice(rowIndex, 1)
                this.$emit('update:value', this.$deepClone(this.localValue))
            }
        }
    }
</script>

<style lang="scss" scoped>
    .field-operation{
        text-align: center;
        font-size: 14px;
        .bk-icon{
            margin-right: 4px;
        }
    }
</style>