<template>
    <bk-table class="cmdb-form-table" :data="list">
        <bk-table-column v-for="col in header"
            :key="col.bk_property_id"
            :prop="col.bk_property_id"
            :label="col.bk_property_name"
            :width="col.bk_property_type === 'bool' ? '90px' : ''">
            <template slot-scope="{ row, $index }">
                <div class="form-item-content">
                    <div class="form-item">
                        <slot name="switch" v-bind="{ row, col, $index }"></slot>
                        <slot name="disabled-tips" v-bind="{ row, col, $index, disabled: checkDisabled(row[col['bk_property_id']], col, $index) }"></slot>
                        <slot
                            :name="`col-${col.bk_property_id}`"
                            v-bind="{ row, sizes, col, disabled: checkDisabled(row[col['bk_property_id']], col, $index), $index }">
                            <component class="form-component"
                                :is="`cmdb-form-${col['bk_property_type']}`"
                                :disabled="checkDisabled(row[col['bk_property_id']], col, $index)"
                                :class="{ error: errors.has(`${col['bk_property_id']}_${$index}`) }"
                                :unit="col.unit"
                                :row="2"
                                :options="col.option || []"
                                :data-vv-name="`${col['bk_property_id']}_${$index}`"
                                :data-vv-as="col['bk_property_name']"
                                :placeholder="getPlaceholder(col)"
                                :auto-select="mode === 'create'"
                                v-bind="sizes"
                                v-validate="getValidateRules(col)"
                                v-model.trim="row[col['bk_property_id']]['value']"
                                v-bk-tooltips="{
                                    disabled: checkTipsDisabled(row[col['bk_property_id']], col, $index),
                                    allowHtml: true,
                                    content: `#disabled-tips-${col['bk_property_id']}_${$index}`,
                                    placement: 'top'
                                }">
                            </component>
                        </slot>
                    </div>
                    <span class="form-item-error">{{errors.first(`${col['bk_property_id']}_${$index}`)}}</span>
                </div>
            </template>
        </bk-table-column>
        <bk-table-column :label="$t('操作')" v-if="operation.show">
            <template slot-scope="{ row, $index }">
                <bk-button
                    class="mr10"
                    theme="primary"
                    :text="true"
                    @click.stop="handleAdd($index)">
                    {{$t('添加')}}
                </bk-button>
                <bk-button
                    theme="primary"
                    :text="true"
                    :disabled="operation.disabled.remove"
                    @click.stop="handleRemove($index)">
                    {{$t('删除')}}
                </bk-button>
            </template>
        </bk-table-column>
        <div slot="empty">
            <button
                class="add-row-button text-primary"
                theme="primary"
                :text="true"
                @click.stop="handleAdd()">
                <i class="bk-icon icon-plus"></i>
                <span>{{$t('立即添加')}}</span>
            </button>
        </div>
    </bk-table>
</template>

<script>
    export default {
        name: 'cmdb-form-table',
        props: {
            value: {
                type: [Array, String],
                default: () => []
            },
            options: {
                type: Array,
                default: () => []
            },
            placeholder: {
                type: String,
                default: ''
            },
            maxlength: {
                type: Number,
                default: 99
            },
            minlength: {
                type: Number,
                default: 0
            },
            size: {
                type: String,
                default: 'small',
                validator (val) {
                    return ['normal', 'large', 'small'].includes(val)
                }
            },
            mode: {
                type: String,
                default: 'create',
                validator (val) {
                    // input自由输入模式
                    // update不允许添加删除行
                    return ['input', 'update', 'create'].includes(val)
                }
            },
            disabledCheck: Function,
            newRowValue: Function
        },
        data () {
            return {
                list: []
            }
        },
        computed: {
            header () {
                return this.options.filter(option => option.editable)
            },
            operation () {
                return {
                    disabled: {
                        add: this.list.length >= this.maxlength,
                        remove: this.list.length <= this.minlength
                    },
                    show: ['input', 'create'].includes(this.mode)
                }
            },
            sizes () {
                const fontSizeMap = { small: 'normal', large: 'large' }
                return {
                    size: this.size,
                    fontSize: fontSizeMap[this.size] || 'medium'
                }
            }
        },
        watch: {
            value: {
                handler (value) {
                    this.list = value || []
                },
                immediate: true
            },
            list: {
                handler (value) {
                    this.$emit('input', value)
                },
                deep: true
            }
        },
        methods: {
            handleAdd (i) {
                this.list.push(this.getNewRow())
            },
            handleRemove (i) {
                this.list.splice(i, 1)
            },
            checkDisabled (field, property, index) {
                const mode = this.mode
                if (mode === 'input') {
                    return false
                }

                if (this.disabledCheck) {
                    return this.disabledCheck(field, property, index)
                }
            },
            getNewRow () {
                const row = {}
                this.header.forEach(prop => {
                    let value = { value: this.getDefaultValue(prop) }
                    if (this.newRowValue) {
                        value = this.newRowValue(prop)
                    }
                    row[prop.bk_property_id] = value
                })
                return row
            },
            getDefaultValue (property) {
                const formValues = this.$tools.getInstFormValues([property])
                return formValues[property.bk_property_id]
            },
            getPlaceholder (property) {
                const placeholderTxt = ['enum', 'list'].includes(property.bk_property_type) ? '请选择xx' : '请输入xx'
                return this.$t(placeholderTxt, { name: property.bk_property_name })
            },
            getValidateRules (property) {
                return this.$tools.getValidateRules(property)
            },
            checkTipsDisabled (field, property, index) {
                const disabled = this.checkDisabled(field, property, index)
                return this.mode === 'create' || !disabled
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cmdb-form-table {
        width: 100%;

        .form-item-content {
            padding: 8px 0;
        }

        /deep/ .form-item {
            display: flex;
            align-items: center;

            .form-checkbox {
                margin-right: 4px;
                overflow: unset;
            }
        }

        .form-item-error {
            font-size: 12px;
            color: #ff5656;
        }

        /deep/ .bk-table-empty-block {
            min-height: 42px;
            .bk-table-empty-text {
                padding: 0;
            }
        }

        .add-row-button {
            line-height: 32px;
            .bk-icon,
            span {
                @include inlineBlock;
            }
            .icon-plus {
                font-size: 20px;
                margin-right: -4px;
            }
        }

        /deep/ .custom-form-error {
            position: relative !important;
        }
    }
</style>
