<template>
    <cmdb-sticky-layout>
        <bk-form class="form-layout" form-type="inline" :label-width="300">
            <bk-form-item class="form-item" :label="$t('任务名称')" required>
                <bk-input class="form-meta"
                    :placeholder="$t('请输入xx', { name: $t('任务名称') })"
                    :data-vv-as="$t('任务名称')"
                    data-vv-name="bk_task_name"
                    v-validate="'required|length:256'"
                    v-model.trim="form.bk_task_name">
                </bk-input>
                <p class="form-error" v-if="errors.has('bk_task_name')">{{errors.first('bk_task_name')}}</p>
            </bk-form-item>
            <bk-form-item class="form-item" :label="$t('账户名称')" required>
                <task-account-selector class="form-meta"
                    ref="accountSelector"
                    :disabled="!isCreateMode"
                    :data-vv-as="$t('账户名称')"
                    data-vv-name="bk_account_id"
                    v-model="form.bk_account_id"
                    v-validate="'required'">
                </task-account-selector>
                <bk-link theme="primary" class="form-account-link" @click="handleLinkToCloudAccount">
                    <i class="form-account-link-icon icon-cc-share"></i>
                    <span>{{$t('跳转账户管理')}}</span>
                </bk-link>
                <p class="form-tips" v-if="form.bk_account_id">
                    <i class="icon-cc-exclamation-tips"></i>
                    {{$t('同步频率提示语')}}
                </p>
                <p class="form-error" v-if="errors.has('bk_account_id')">{{errors.first('bk_account_id')}}</p>
            </bk-form-item>
            <bk-form-item class="form-item" :label="$t('资源类型')" required>
                <task-resource-selector class="form-meta"
                    :disabled="true"
                    :data-vv-as="$t('资源类型')"
                    data-vv-name="bk_resource_type"
                    v-model="form.bk_resource_type"
                    v-validate="'required'">
                </task-resource-selector>
                <p class="form-error" v-if="errors.has('bk_resource_type')">{{errors.first('bk_resource_type')}}</p>
            </bk-form-item>
        </bk-form>
        <bk-form class="form-layout" form-type="inline" :label-width="300" v-if="form.bk_account_id">
            <bk-form-item class="form-item" :label="$t('云区域设定')" required>
                <bk-button @click="handleAddVPC">{{$t('添加VPC')}}</bk-button>
                <input type="hidden" v-validate="'min_value:1'" :value="selectedVPC.length" name="vpc-count">
                <p class="form-error" v-if="errors.has('vpc-count')">{{$t('请至少选择一个VPC')}}</p>
            </bk-form-item>
        </bk-form>
        <div class="form-setting-component" v-if="form.bk_account_id">
            <task-form-table
                ref="vpcTable"
                :selected="selectedVPC"
                :account="form.bk_account_id"
                @remove="handleRemove">
            </task-form-table>
        </div>
        <div class="form-options"
            slot="footer"
            slot-scope="{ sticky }"
            :class="{ 'is-sticky': sticky }">
            <cmdb-auth :auth="auth">
                <bk-button theme="primary" slot-scope="{ disabled }"
                    :disabled="disabled"
                    :loading="$loading([request.createTask, request.updateTask, request.createArea])"
                    @click="handleSumbit">
                    {{isCreateMode ? $t('提交') : $t('保存')}}
                </bk-button>
            </cmdb-auth>
            <bk-button class="ml10" @click="handleCancel">{{$t('取消')}}</bk-button>
        </div>
        <cmdb-dialog
            v-model="showVPCSelector"
            :width="850"
            :height="460"
            :show-close-icon="false">
            <task-vpc-selector
                :account="form.bk_account_id"
                :selected="selectedVPC"
                @change="handleVPCChange"
                @cancel="handleVPCCancel">
            </task-vpc-selector>
        </cmdb-dialog>
    </cmdb-sticky-layout>
</template>

<script>
    import { MENU_RESOURCE_CLOUD_ACCOUNT } from '@/dictionary/menu-symbol'
    import TaskFormTable from './task-form-table.vue'
    import TaskVpcSelector from './task-vpc-selector.vue'
    import TaskDetailsInfo from './task-details-info.vue'
    import TaskAccountSelector from './task-account-selector.vue'
    import TaskResourceSelector from './task-resource-selector.vue'
    import Bus from '@/utils/bus.js'
    import symbols from '../common/symbol'
    export default {
        name: 'task-form',
        components: {
            [TaskFormTable.name]: TaskFormTable,
            [TaskVpcSelector.name]: TaskVpcSelector,
            [TaskAccountSelector.name]: TaskAccountSelector,
            [TaskResourceSelector.name]: TaskResourceSelector
        },
        props: {
            task: {
                type: Object,
                default: null
            },
            container: {
                type: Object,
                required: true
            }
        },
        data () {
            const form = {
                bk_task_name: '',
                bk_account_id: '',
                bk_resource_type: 'host'
            }
            if (this.task) {
                Object.assign(form, {
                    bk_task_name: this.task.bk_task_name,
                    bk_account_id: this.task.bk_account_id,
                    bk_resource_type: this.task.bk_resource_type
                })
            }
            return {
                accounts: [],
                form: form,
                selectedVPC: this.task ? [...this.task.bk_sync_vpcs] : [],
                request: {
                    getAccounts: symbols.get('getAccounts'),
                    createTask: symbols.get('createTask'),
                    updateTask: symbols.get('updateTask'),
                    createArea: symbols.get('createArea')
                },
                showVPCSelector: false
            }
        },
        computed: {
            isCreateMode () {
                return this.task === null
            },
            auth () {
                if (this.isCreateMode) {
                    return {
                        type: this.$OPERATION.C_CLOUD_RESOURCE_TASK
                    }
                }
                return {
                    type: this.$OPERATION.U_CLOUD_RESOURCE_TASK,
                    relation: [this.task.bk_task_id]
                }
            }
        },
        watch: {
            selectedVPC () {
                this.errors.remove('vpc-count')
            }
        },
        methods: {
            handleLinkToCloudAccount () {
                this.$router.push({
                    name: MENU_RESOURCE_CLOUD_ACCOUNT
                })
            },
            handleAddVPC () {
                this.showVPCSelector = true
            },
            handleRemove (target) {
                const index = this.selectedVPC.findIndex(vpc => vpc.bk_vpc_id === target.bk_vpc_id)
                index > -1 && this.selectedVPC.splice(index, 1)
            },
            async handleSumbit () {
                const isFormValid = await this.$validator.validateAll()
                const isDirectoryValid = this.$refs.vpcTable && await this.$refs.vpcTable.$validator.validateAll()
                if (!isFormValid || !isDirectoryValid) {
                    return false
                }

                const vendor = this.$refs.accountSelector.accountVendor
                const next = await this.$refs.vpcTable.createCloudArea({
                    bk_cloud_vendor: vendor.id,
                    bk_account_id: this.form.bk_account_id
                })
                if (!next) {
                    return false
                }

                if (this.isCreateMode) {
                    this.doCreate()
                } else {
                    this.doUpdate()
                }
            },
            async doCreate () {
                try {
                    await this.$store.dispatch('cloud/resource/createTask', {
                        params: {
                            ...this.form,
                            bk_sync_vpcs: this.getVPCList()
                        },
                        config: {
                            requestId: this.request.createTask
                        }
                    })
                    this.container.hide()
                    Bus.$emit('request-refresh')
                } catch (e) {
                    console.error(e)
                }
            },
            async doUpdate () {
                try {
                    const params = {
                        bk_task_id: this.task.bk_task_id,
                        bk_task_name: this.form.bk_task_name,
                        bk_sync_vpcs: this.getVPCList()
                    }
                    await this.$store.dispatch('cloud/resource/updateTask', {
                        id: this.task.bk_task_id,
                        params: params,
                        config: {
                            requestId: this.request.updateTask
                        }
                    })
                    this.container.show({
                        detailsComponent: TaskDetailsInfo.name,
                        task: {
                            ...this.task,
                            ...params
                        }
                    })
                    Bus.$emit('request-refresh')
                } catch (e) {
                    console.error(e)
                }
            },
            getVPCList () {
                return this.$refs.vpcTable.list.map(row => {
                    return {
                        bk_vpc_id: row.bk_vpc_id,
                        bk_vpc_name: row.bk_vpc_name,
                        bk_region: row.bk_region,
                        bk_host_count: row.bk_host_count,
                        bk_sync_dir: row.bk_sync_dir,
                        bk_cloud_id: row.bk_cloud_id,
                        destroyed: row.destroyed
                    }
                })
            },
            handleCancel () {
                if (this.isCreateMode) {
                    this.container.hide()
                } else {
                    this.container.show({
                        detailsComponent: TaskDetailsInfo.name,
                        props: {
                            task: { ...this.task }
                        }
                    })
                }
            },
            handleVPCChange (vpcs) {
                this.showVPCSelector = false
                this.selectedVPC = vpcs
            },
            handleVPCCancel () {
                this.showVPCSelector = false
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-layout {
        padding: 10px 25px;
        display: block;
        font-size: 0;
        .form-item:nth-child(2n) {
            padding-left: 32px;
        }
        .form-item:nth-child(n + 3) {
            margin-top: 10px;
        }
        .form-item {
            width: 50%;
            margin-left: 0;
            /deep/ {
                .bk-label,
                .bk-form-content {
                    display: block;
                    float: none;
                    text-align: left;
                }
            }
            .form-meta {
                display: block;
            }
            .form-error {
                position: absolute;
                top: 100%;
                left: 0;
                line-height: 1.5;
                font-size: 12px;
                color: $dangerColor;
            }
        }
        .form-account-link {
            position: absolute;
            bottom: 100%;
            right: 0;
            /deep/ {
                .bk-link-text {
                    display: flex;
                    align-items: center;
                    font-size: 12px;
                }
            }
            .form-account-link-icon {
                height: 20px;
                line-height: 20px;
                margin-right: 4px;
            }
        }
        .form-tips {
            position: absolute;
            top: 100%;
            left: 0;
            font-size: 12px;
            line-height: 1.5;
            color: #979BA5;
            [class^=icon] {
                vertical-align: 1px;
            }
        }
    }
    .form-options {
        font-size: 0;
        margin-top: 10px;
        padding: 0 25px;
        &.is-sticky {
            margin-top: 0;
            padding: 15px 25px;
            border-top: 1px solid $borderColor;
            background-color: #FAFBFD;
        }
    }
    .form-setting-component {
        margin: 10px 25px 20px;
    }
</style>
