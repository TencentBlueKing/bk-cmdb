<template>
    <div class="cmdb-form-table">
        <bk-table :data="localValue">
            <bk-table-column
                v-for="property in columns"
                :prop="property.bk_property_id"
                :label="property.bk_property_name"
                :key="property.bk_property_id"
                :min-width="property.bk_property_type === 'bool' ? 80 : 150"
                :render-header="(h, data) => renderHeader(h, data, property)">
                <template slot-scope="rowProps">
                    <custom-content-render
                        v-if="$scopedSlots[property.bk_property_id]"
                        v-bind="rowProps"
                        :render="$scopedSlots[property.bk_property_id]">
                    </custom-content-render>
                    <component class="form-component"
                        v-else
                        v-validate="$tools.getValidateRules(property)"
                        size="small"
                        :data-vv-name="property.bk_property_id"
                        :data-vv-as="property.bk_property_name"
                        :is="`cmdb-form-${property.bk_property_type}`"
                        :unit="property.unit"
                        :row="2"
                        :disabled="disabled"
                        :options="property.option || []"
                        :auto-select="false"
                        :value="localValue[rowProps.$index][property.bk_property_id]"
                        @input="handleColumnValueChange(rowProps, ...arguments)">
                    </component>
                </template>
            </bk-table-column>
            <bk-table-column
                v-if="mode !== 'info'"
                min-width="100"
                :label="$t('操作')">
                <template slot-scope="{ row, $index }">
                    <bk-button class="mr10" theme="primary" text
                        @click.stop="handleAddRow($event, $index)">
                        {{$t('添加')}}
                    </bk-button>
                    <bk-button theme="primary" text
                        @click.stop="handleDeleteRow($event, $index)">
                        {{$t('删除')}}
                    </bk-button>
                </template>
            </bk-table-column>
            <div slot="empty">
                <span v-if="mode === 'info'">{{$t('暂无数据')}}</span>
                <bk-button class="empty-button"
                    v-else
                    theme="primary"
                    text
                    icon="icon-plus"
                    @click.stop="handleAddRow($event, -1)">
                    {{$t('立即添加')}}
                </bk-button>
            </div>
        </bk-table>
    </div>
</template>

<script>
    import CustomContentRender from './table-custom-content-render'
    export default {
        name: 'cmdb-form-table',
        components: {
            CustomContentRender
        },
        props: {
            value: {
                type: Array,
                default: () => ([])
            },
            options: {
                type: Array,
                required: true
            },
            mode: {
                type: String,
                default: 'create',
                validator (mode) {
                    return ['create', 'update', 'info'].includes(mode)
                }
            },
            disabled: Boolean
        },
        data () {
            return {}
        },
        computed: {
            columns () {
                if (this.mode === 'update') {
                    return this.options.filter(property => property.editable && !property.bk_isapi)
                }
                if (this.mode === 'create') {
                    return this.options.filter(property => !property.bk_isapi)
                }
                return this.options.filter(property => property.editable)
            },
            localValue: {
                get () {
                    return this.value || []
                },
                set (value) {
                    this.$emit('input', value)
                    this.$emit('change', value)
                }
            }
        },
        methods: {
            renderHeader (h, { column, $index }, property) {
                if (!property.placeholder) {
                    return property.bk_property_name
                }
                const directive = {
                    name: 'bkTooltips',
                    content: property.placeholder,
                    placement: 'top',
                    trigger: 'click'
                }
                const style = {
                    'text-decoration': 'underline',
                    'text-decoration-style': 'dashed'
                }
                return (<span v-bk-tooltips={ directive } style={ style }>{ property.bk_property_name }</span>)
            },
            handleColumnValueChange ({ row, column, $index }, value) {
                const newRowValue = { ...row }
                newRowValue[column.property] = value
                const newValues = [...this.localValue]
                newValues.splice($index, 1, newRowValue)
                this.localValue = newValues
            },
            handleAddRow (event, index) {
                const newRowIndex = index + 1
                const newRowValue = this.$tools.getInstFormValues(this.columns, {}, false)
                const newValues = [...this.localValue]
                newValues.splice(newRowIndex, 0, newRowValue)
                this.localValue = newValues
                this.$emit('add-row', newRowValue, newRowIndex)
            },
            handleDeleteRow (event, index) {
                const newValues = [...this.localValue]
                const [deleteRow] = newValues.splice(index, 1)
                this.localValue = newValues
                this.$emit('delete-row', deleteRow, index)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cmdb-form-table {
        /deep/ .bk-table-empty-block {
            min-height: 42px;
            .bk-table-empty-text {
                padding: 0;
            }
        }
    }
    .empty-button {
        /deep/ {
            > div {
                display: flex;
                align-items: center;
                justify-content: center;
            }
            .left-icon {
                font-size: 20px;
            }
        }
    }
</style>
