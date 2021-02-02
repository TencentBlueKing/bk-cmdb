<template>
    <table class="audit-resource-options">
        <colgroup>
            <col width="4%">
            <col width="28%">
            <col width="6%">
            <col width="28%">
            <col width="6%">
            <col width="28%">
            <col width="10%" v-if="isModelInstance">
        </colgroup>
        <tr>
            <td align="right"><label class="option-label">{{$t('操作对象')}}</label></td>
            <td>
                <audit-target-selector class="option-value"
                    searchable
                    category="resource"
                    :placeholder="$t('请选择xx', { name: $t('操作对象') })"
                    v-model="condition.resource_type">
                </audit-target-selector>
            </td>
            <template v-if="isModelInstance">
                <td align="right"><label class="option-label">{{$t('模型')}}</label></td>
                <td>
                    <audit-model-selector class="option-value"
                        searchable
                        :placeholder="$t('请选择xx', { name: $t('模型') })"
                        v-model="condition.bk_obj_id">
                    </audit-model-selector>
                </td>
            </template>
            <td align="right"><label class="option-label">{{$t('动作')}}</label></td>
            <td>
                <audit-action-selector class="option-value"
                    :target="condition.resource_type"
                    :placeholder="$t('请选择xx', { name: $t('动作') })"
                    :empty-text="$t('请先选择操作对象')"
                    v-model="condition.action">
                </audit-action-selector>
            </td>
            <template v-if="!isModelInstance">
                <td align="right"><label class="option-label">{{$t('时间')}}</label></td>
                <td>
                    <cmdb-form-date-range class="option-value"
                        font-size="medium"
                        :placeholder="$t('请选择xx', { name: $t('时间') })"
                        v-model="condition.operation_time">
                    </cmdb-form-date-range>
                </td>
            </template>
        </tr>
        <tr>
            <template v-if="isModelInstance">
                <td align="right"><label class="option-label">{{$t('时间')}}</label></td>
                <td>
                    <cmdb-form-date-range class="option-value"
                        font-size="medium"
                        :placeholder="$t('请选择xx', { name: $t('时间') })"
                        :clearable="false"
                        v-model="condition.operation_time">
                    </cmdb-form-date-range>
                </td>
            </template>
            <td align="right"><label class="option-label">{{$t('账号')}}</label></td>
            <td>
                <cmdb-form-objuser class="option-value"
                    v-model="condition.user"
                    :exclude="false"
                    :multiple="false"
                    :placeholder="$t('请输入xx', { name: $t('账号') })">
                </cmdb-form-objuser>
            </td>
            <td align="right"><label class="option-label">{{$t('实例')}}</label></td>
            <td>
                <bk-input class="option-value"
                    v-model="instanceFilter"
                    :placeholder="$t('请输入xx', { name: instanceType === 'resource_id' ? 'ID' : $t('名称') })">
                    <bk-select class="option-type" slot="prepend"
                        :clearable="false"
                        v-model="instanceType">
                        <bk-option id="resource_name" :name="$t('名称')"></bk-option>
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
    import AuditTargetSelector from './audit-target-selector'
    import AuditActionSelector from './audit-action-selector'
    import AuditModelSelector from './audit-model-selector'
    import RouterQuery from '@/router/query'
    export default {
        name: 'audit-resource-options',
        components: {
            AuditTargetSelector,
            AuditActionSelector,
            AuditModelSelector
        },
        data () {
            const today = this.$tools.formatTime(new Date(), 'YYYY-MM-DD')
            const defaultCondition = {
                bk_biz_id: '',
                resource_type: '',
                action: [],
                operation_time: [today, today],
                user: '',
                resource_id: '',
                resource_name: '',
                category: 'resource',
                bk_obj_id: ''
            }
            return {
                instanceType: 'resource_name',
                defaultCondition,
                condition: { ...defaultCondition }
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
            },
            isModelInstance () {
                return this.condition.resource_type === 'model_instance'
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
                    tab: 'resource',
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
    .audit-resource-options {
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
            max-width: 400px;
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
