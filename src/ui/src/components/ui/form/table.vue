<template>
    <bk-table class="cmdb-form-table" :data="list">
        <bk-table-column v-for="property in header"
            :key="property.bk_property_id"
            :prop="property.bk_property_id"
            :label="property.bk_property_name"
            :width="property.bk_property_type === 'bool' ? '90px' : ''"
            show-overflow-tooltip>
            <template slot-scope="{ row: values, $index }">
                <div class="form-item-content">
                    <div class="form-item">
                        <bk-checkbox
                            v-if="mode === 'create'"
                            class="form-checkbox"
                            v-model="values[property['bk_property_id']]['as_default_value']"
                            @change="handleResetValue(values[property['bk_property_id']]['as_default_value'], property)">
                        </bk-checkbox>
                        <component class="form-component"
                            :is="`cmdb-form-${property['bk_property_type']}`"
                            :disabled="!values[property['bk_property_id']]['as_default_value']"
                            :class="{ error: errors.has(`${property['bk_property_id']}_${$index}`) }"
                            :unit="property.unit"
                            :row="2"
                            :options="property.option || []"
                            :data-vv-name="`${property['bk_property_id']}_${$index}`"
                            :data-vv-as="property['bk_property_name']"
                            :placeholder="getPlaceholder(property)"
                            :auto-select="false"
                            v-bind="sizes"
                            v-validate="getValidateRules(property)"
                            v-model.trim="values[property['bk_property_id']]['value']">
                        </component>
                    </div>
                    <span class="form-item-error">{{errors.first(`${property['bk_property_id']}_${$index}`)}}</span>
                </div>
            </template>
        </bk-table-column>
        <bk-table-column :label="$t('操作')" v-if="!readonly">
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
            <bk-button
                theme="primary"
                :text="true"
                @click.stop="handleAdd()">
                {{$t('立即添加')}}
            </bk-button>
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
                    return ['update', 'create', 'readonly'].includes(val)
                }
            },
            extra: {
                type: Object,
                default: () => ({})
            }
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
                    }
                }
            },
            sizes () {
                const fontSizeMap = { small: 'normal', large: 'large' }
                return {
                    size: this.size,
                    fontSize: fontSizeMap[this.size] || 'medium'
                }
            },
            readonly () {
                return this.mode === 'readonly'
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
            handleResetValue (a, b) {
            },
            getNewRow () {
                const row = {}
                this.header.forEach(prop => {
                    row[prop.bk_property_id] = {
                        value: '',
                        as_default_value: false
                    }
                })
                return row
            },
            getPlaceholder (property) {
                const placeholderTxt = ['enum', 'list'].includes(property.bk_property_type) ? '请选择xx' : '请输入xx'
                return this.$t(placeholderTxt, { name: property.bk_property_name })
            },
            getValidateRules (property) {
                return this.$tools.getValidateRules(property)
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

        /deep/ .bk-table-empty-text {
            padding: 0;
        }
    }
</style>
