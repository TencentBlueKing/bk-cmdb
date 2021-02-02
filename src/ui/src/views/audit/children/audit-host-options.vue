<template>
    <table class="audit-host-options">
        <colgroup>
            <col width="4%">
            <col width="28%">
            <col width="6%">
            <col width="28%">
            <col width="6%">
            <col width="28%">
        </colgroup>
        <tr>
            <td align="right"><label class="option-label">{{$t('业务')}}</label></td>
            <td>
                <audit-business-selector class="option-value"
                    searchable
                    :placeholder="$t('请选择xx', { name: $t('业务') })"
                    v-model="condition.bk_biz_id">
                </audit-business-selector>
            </td>
            <td align="right"><label class="option-label">{{$t('动作')}}</label></td>
            <td>
                <audit-action-selector class="option-value"
                    :target="condition.resource_type"
                    :placeholder="$t('请选择xx', { name: $t('动作') })"
                    :empty-text="$t('请先选择操作对象')"
                    v-model="condition.action">
                </audit-action-selector>
            </td>
            <td align="right"><label class="option-label">{{$t('时间')}}</label></td>
            <td>
                <cmdb-form-date-range class="option-value"
                    font-size="medium"
                    :placeholder="$t('请选择xx', { name: $t('时间') })"
                    :clearable="false"
                    v-model="condition.operation_time">
                </cmdb-form-date-range>
            </td>
        </tr>
        <tr>
            <td align="right"><label class="option-label">{{$t('账号')}}</label></td>
            <td>
                <cmdb-form-objuser class="option-value"
                    v-model="condition.user"
                    :exclude="false"
                    :multiple="false"
                    :placeholder="$t('请输入xx', { name: $t('账号') })">
                </cmdb-form-objuser>
            </td>
            <td align="right"><label class="option-label">{{$t('主机')}}</label></td>
            <td>
                <bk-input class="option-value"
                    v-model="instanceFilter"
                    :placeholder="$t('请输入xx', { name: instanceType === 'resource_id' ? 'ID' : 'IP' })">
                    <bk-select class="option-type" slot="prepend"
                        :clearable="false"
                        v-model="instanceType">
                        <bk-option id="resource_name" :name="$t('IP')"></bk-option>
                        <bk-option id="resource_id" name="ID"></bk-option>
                    </bk-select>
                </bk-input>
            </td>
            <td></td>
            <td>
                <div class="options-button">
                    <bk-button class="mr10" theme="primary" @click="handleSearch(1)">{{$t('查询')}}</bk-button>
                    <bk-button theme="default" @click="handleReset">{{$t('清空')}}</bk-button>
                </div>
            </td>
        </tr>
    </table>
</template>

<script>
    import AuditActionSelector from './audit-action-selector'
    import AuditBusinessSelector from '@/components/audit-history/audit-business-selector'
    import RouterQuery from '@/router/query'
    export default {
        name: 'audit-host-options',
        components: {
            AuditActionSelector,
            AuditBusinessSelector
        },
        data () {
            const today = this.$tools.formatTime(new Date(), 'YYYY-MM-DD')
            const defaultCondition = {
                bk_biz_id: '',
                resource_type: 'host',
                action: [],
                operation_time: [today, today],
                user: '',
                resource_id: '',
                resource_name: '',
                category: 'host'
            }
            return {
                instanceType: 'resource_name',
                defaultCondition,
                condition: {
                    ...defaultCondition
                }
            }
        },
        computed: {
            instanceFilter: {
                get () {
                    return this.condition[this.instanceType]
                },
                set (value) {
                    this.condition[this.instanceType] = value
                }
            }
        },
        watch: {
            instanceType () {
                this.condition.resource_id = ''
                this.condition.resource_name = ''
            }
        },
        created () {
            this.handleSearch()
        },
        methods: {
            handleSearch (isEvent) {
                this.$emit('condition-change', this.condition)
                RouterQuery.set({
                    tab: 'host',
                    page: 1,
                    _t: Date.now(),
                    _e: isEvent
                })
            },
            handleReset () {
                this.condition = { ...this.defaultCondition }
                this.handleSearch()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .audit-host-options {
        width: 100%;
        padding: 5px 0;
        tr {
            td {
                padding: 5px 0;
            }
        }
        .option-label {
            font-size: 14px;
            padding: 0 10px;
            @include ellipsis;
        }
        .option-value {
            width: 100%;
            min-width: 230px;
        }
        .option-type {
            width: 80px;
            margin-top: -1px;
            border-color: #c4c6cc transparent;
            box-shadow: none;
        }
        .options-button {
            display: flex;
            align-items: center;
        }
    }
</style>
