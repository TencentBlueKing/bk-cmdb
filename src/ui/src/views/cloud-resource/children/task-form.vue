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
                <bk-select class="form-meta"
                    :searchable="true"
                    :readonly="!isCreateMode"
                    :placeholder="$t('请选择xx', { name: $t('账户名称') })"
                    :data-vv-as="$t('账户名称')"
                    :loading="$loading(request.getAccounts)"
                    data-vv-name="bk_account_id"
                    v-model="form.bk_account_id"
                    v-validate="'required'">
                    <bk-option v-for="account in accounts"
                        :key="account.bk_account_id"
                        :name="account.bk_account_name"
                        :id="account.bk_account_id">
                    </bk-option>
                </bk-select>
                <link-button class="form-account-link" @click="handleLinkToCloudAccount">
                    <i class="icon-cc-share"></i>
                    {{$t('跳转账户管理')}}
                </link-button>
                <p class="form-tips" v-if="form.bk_account_id">
                    <i class="icon-cc-exclamation-tips"></i>
                    {{$t('同步频率提示语')}}
                </p>
                <p class="form-error" v-if="errors.has('bk_account_id')">{{errors.first('bk_account_id')}}</p>
            </bk-form-item>
            <bk-form-item class="form-item" :label="$t('资源类型')" required>
                <bk-select class="form-meta"
                    :readonly="true"
                    :placeholder="$t('请选择xx', { name: $t('资源类型') })"
                    :data-vv-as="$t('资源类型')"
                    data-vv-name="bk_resource_type"
                    v-model="form.bk_resource_type"
                    v-validate="'required'">
                    <bk-option id="host" :name="$t('主机')"></bk-option>
                </bk-select>
                <p class="form-error" v-if="errors.has('bk_resource_type')">{{errors.first('bk_resource_type')}}</p>
            </bk-form-item>
        </bk-form>
        <bk-form class="form-layout" form-type="inline" :label-width="300" v-if="form.bk_account_id">
            <bk-form-item class="form-item" :label="$t('云区域设定')" required>
                <bk-button @click="handleAddVPC">{{$t('添加VPC')}}</bk-button>
            </bk-form-item>
        </bk-form>
        <div class="form-setting-component" v-if="form.bk_account_id">
            <task-form-table
                ref="vpcTable"
                :selected="selectedVPC"
                :account="form.bk_account_id">
            </task-form-table>
        </div>
        <div class="form-options"
            slot="footer"
            slot-scope="{ sticky }"
            :class="{ 'is-sticky': sticky }">
            <bk-button theme="primary"
                :loading="$loading([request.createTask])"
                @click="handleSumbit">
                {{$t('提交')}}
            </bk-button>
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
    import ResourceFormTable from './task-form-table.vue'
    import ResourceVpcSelector from './task-vpc-selector.vue'
    import TaskDetailsInfo from './task-details-info.vue'
    export default {
        name: 'task-form',
        components: {
            [ResourceFormTable.name]: ResourceFormTable,
            [ResourceVpcSelector.name]: ResourceVpcSelector
        },
        props: {
            task: {
                type: Object,
                default: null
            },
            mode: {
                type: String,
                default: 'create'
            },
            container: {
                type: Object,
                required: true
            }
        },
        data () {
            return {
                accounts: [],
                form: {
                    bk_task_name: '',
                    bk_account_id: '',
                    bk_resource_type: 'host'
                },
                selectedVPC: [],
                request: {
                    getAccounts: Symbol('getAccounts'),
                    createTask: Symbol('createTask')
                },
                showVPCSelector: false
            }
        },
        computed: {
            isCreateMode () {
                return this.mode === 'create'
            }
        },
        created () {
            this.getAccounts()
        },
        methods: {
            async getAccounts () {
                try {
                    const { info: accounts } = await this.$store.dispatch('cloud/account/findMany', {
                        params: {},
                        config: {
                            requestId: this.request.getAccounts
                        }
                    })
                    this.accounts = accounts
                } catch (e) {
                    console.error(e)
                    this.accounts = []
                }
            },
            handleLinkToCloudAccount () {
                this.$router.push({
                    name: MENU_RESOURCE_CLOUD_ACCOUNT
                })
            },
            handleAddVPC () {
                this.showVPCSelector = true
            },
            handleSumbit () {
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
                            bk_sync_vpcs: this.$refs.vpcTable.getSyncVPC()
                        },
                        config: {
                            requestId: this.request.createTask
                        }
                    })
                    this.container.hide('request-refresh')
                } catch (e) {
                    console.error(e)
                }
            },
            async doUpdate () {
                try {
                    await Promise.resolve()
                    this.container.show({
                        detailsComponent: TaskDetailsInfo.name,
                        props: {
                            mission: this.mission
                        }
                    })
                } catch (e) {
                    console.error(e)
                }
            },
            handleCancel () {
                if (this.isCreateMode) {
                    this.container.hide()
                } else {
                    this.container.show({
                        detailsComponent: TaskDetailsInfo.name,
                        props: {
                            mission: this.mission
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
        .form-item:nth-child(2n) {
            margin-left: 32px;
        }
        .form-item:nth-child(n + 3) {
            margin-top: 10px;
        }
        .form-item {
            width: 300px;
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
            font-size: 12px;
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
