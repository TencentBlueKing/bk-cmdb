<template>
    <div class="cmdb-form-process-table">
        <cmdb-form-table
            v-bind="$attrs"
            v-model="localValue"
            :options="options"
            :mode="mode">
            <div class="process-table-content"
                v-for="column in options"
                slot-scope="rowProps"
                :slot="column.bk_property_id"
                :key="`row-${rowProps.index}-${column.bk_property_id}`">
                <bk-popover class="content-value" :disabled="!isLocked(rowProps)">
                    <component
                        size="small"
                        font-size="small"
                        v-validate="$tools.getValidateRules(column)"
                        :disabled="isLocked(rowProps)"
                        :data-vv-name="column.bk_property_id"
                        :data-vv-as="column.bk_property_name"
                        :data-vv-scope="column.bk_property_group || 'bind_info'"
                        :is="getComponentType(column)"
                        :options="column.option || []"
                        :placeholder="getPlaceholder(column)"
                        :value="localValue[rowProps.index][column.bk_property_id]"
                        :auto-select="false"
                        @input="handleColumnValueChange(rowProps, ...arguments)">
                    </component>
                    <i18n path="进程表单锁定提示" slot="content">
                        <bk-link theme="primary" @click="handleRedirect" place="link">{{$t('跳转服务模板')}}</bk-link>
                    </i18n>
                </bk-popover>
            </div>
        </cmdb-form-table>
        <span class="form-error" v-if="validateMsg">{{validateMsg}}</span>
    </div>
</template>

<script>
    import ProcessFormPropertyIp from './process-form-property-ip'
    export default {
        components: {
            ProcessFormPropertyIp
        },
        props: {
            value: {
                type: Array,
                default: () => ([])
            },
            options: {
                type: Array,
                required: true
            }
        },
        inject: ['form'],
        computed: {
            localValue: {
                get () {
                    return (this.value || [])
                },
                set (values) {
                    this.$emit('input', values)
                    this.$emit('change', values)
                }
            },
            lockStates () {
                const property = this.form.processTemplate.property || { bind_info: { value: [] } }
                const values = property.bind_info.value || []
                return values.map(row => {
                    const state = {}
                    Object.keys(row).forEach(key => {
                        // 可能存在as_default_value为null的情况：isapi为true的字段
                        state[key] = !!row[key].as_default_value
                    })
                    return state
                })
            },
            mode () {
                return this.form.serviceTemplateId ? 'info' : 'update'
            },
            validateMsg () {
                const hasError = this.errors.items.some(item => item.scope === 'bind_info')
                return hasError ? this.$t('有未正确定义的监听信息') : null
            }
        },
        methods: {
            isLocked ({ row, column, index }) {
                const rowState = this.lockStates[index]
                return rowState ? rowState[column.property] : false
            },
            getComponentType (property) {
                if (property.bk_property_id === 'ip') {
                    return 'process-form-property-ip'
                }
                return `cmdb-form-${property.bk_property_type}`
            },
            getPlaceholder (property) {
                const placeholderTxt = ['enum', 'list'].includes(property.bk_property_type) ? '请选择xx' : '请输入xx'
                return this.$t(placeholderTxt, { name: property.bk_property_name })
            },
            handleColumnValueChange ({ row, column, index }, value) {
                const rowValue = { ...row }
                rowValue[column.property] = value
                const newValues = [...this.localValue]
                newValues.splice(index, 1, rowValue)
                this.localValue = newValues
            },
            handleRedirect () {
                this.$routerActions.redirect({
                    name: 'operationalTemplate',
                    params: {
                        bizId: this.form.bizId,
                        templateId: this.form.serviceTemplateId
                    },
                    history: true
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cmdb-form-process-table {
        position: relative;
        .process-table-content {
            display: flex;
            align-items: center;
            justify-content: flex-start;
            .content-value {
                width: 100%;
                /deep/ {
                    .bk-tooltip-ref {
                        width: 100%;
                    }
                }
            }
        }
        .form-error {
            position: absolute;
            top: 100%;
            left: 0;
            line-height: 14px;
            font-size: 12px;
            color: $dangerColor;
            max-width: 100%;
            @include ellipsis;
        }
    }
</style>
