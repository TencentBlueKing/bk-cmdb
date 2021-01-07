<template>
    <div class="cloud-account-layout">
        <cmdb-tips class="cloud-account-tips" tips-key="cloud-account-tips">
            <i18n path="云账户提示语">
                <bk-button text size="small" place="resource" style="padding: 0" @click="linkResource">{{$t('云资源')}}</bk-button>
            </i18n>
        </cmdb-tips>
        <div class="cloud-account-options">
            <div class="options-left">
                <cmdb-auth :auth="{ type: $OPERATION.C_CLOUD_ACCOUNT }">
                    <bk-button theme="primary" slot-scope="{ disabled }"
                        :disabled="disabled"
                        @click="handleCreate">
                        {{$t('新建')}}
                    </bk-button>
                </cmdb-auth>
            </div>
            <div class="options-right">
                <bk-input class="options-filter" clearable
                    v-model.trim="filter"
                    right-icon="icon-search"
                    :placeholder="$t('请输入xx', { name: $t('账户名称') })">
                </bk-input>
            </div>
        </div>
        <bk-table class="cloud-account-table" v-bkloading="{ isLoading: $loading(request.search) }"
            :data="list"
            :pagination="pagination"
            :max-height="$APP.height - 220"
            @sort-change="handleSortChange"
            @page-change="handlePageChange"
            @page-limit-change="handleLimitChange"
            @cell-click="handleCellClick">
            <bk-table-column :label="$t('账户名称')"
                prop="bk_account_name"
                class-name="is-highlight"
                sortable="custom"
                show-overflow-tooltip>
            </bk-table-column>
            <bk-table-column :label="$t('账户类型')" prop="bk_cloud_vendor" sortable="custom">
                <cmdb-vendor slot-scope="{ row }" :type="row.bk_cloud_vendor"></cmdb-vendor>
            </bk-table-column>
            <bk-table-column :label="$t('状态')" prop="status">
                <template slot-scope="{ row }">
                    <span :class="['cloud-account-status', row.status, { pending: row.pending }]"
                        v-bk-tooltips="{
                            content: row.error_message,
                            disabled: !row.error_message
                        }">
                        {{getStatusText(row.status)}}
                    </span>
                </template>
            </bk-table-column>
            <bk-table-column :label="$t('修改人')" prop="bk_last_editor" show-overflow-tooltip>
                <template slot-scope="{ row }">{{row.bk_last_editor | formatter('singlechar')}}</template>
            </bk-table-column>
            <bk-table-column :label="$t('修改时间')" prop="last_time" sortable="custom">
                <template slot-scope="{ row }">{{row.last_time | formatter('time')}}</template>
            </bk-table-column>
            <bk-table-column :label="$t('操作')">
                <template slot-scope="{ row }">
                    <link-button class="mr10" @click="handleView(row)">{{$t('查看')}}</link-button>
                    <cmdb-auth :auth="{ type: $OPERATION.D_CLOUD_ACCOUNT, relation: [row.bk_account_id] }">
                        <link-button slot-scope="{ disabled }"
                            :disabled="!row.bk_can_delete_account || disabled"
                            v-bk-tooltips="{
                                disabled: row.bk_can_delete_account || disabled,
                                content: $t('云账户禁止删除提示')
                            }"
                            @click="handleDelete(row)">
                            {{$t('删除')}}
                        </link-button>
                    </cmdb-auth>
                </template>
            </bk-table-column>
        </bk-table>
        <account-sideslider ref="accountSideslider" @request-refresh="getData"></account-sideslider>
    </div>
</template>

<script>
    import AccountSideslider from './children/account-sideslider.vue'
    import CmdbVendor from '@/components/ui/other/vendor'
    import { MENU_RESOURCE_CLOUD_RESOURCE } from '@/dictionary/menu-symbol'
    import throttle from 'lodash.throttle'
    import RouterQuery from '@/router/query'
    export default {
        components: {
            AccountSideslider,
            CmdbVendor
        },
        data () {
            return {
                filter: '',
                list: [],
                pagination: this.$tools.getDefaultPaginationConfig(),
                sort: 'bk_account_id',
                request: {
                    search: Symbol('search'),
                    status: Symbol('status')
                },
                scheduleSearch: throttle(this.handlePageChange, 800, { leading: false, trailing: true })
            }
        },
        watch: {
            filter () {
                this.scheduleSearch()
            }
        },
        created () {
            this.unwatch = RouterQuery.watch('*', ({
                page,
                limit,
                sort
            }) => {
                this.pagination.current = parseInt(page || this.pagination.current, 10)
                this.pagination.limit = parseInt(limit || this.pagination.limit, 10)
                this.sort = sort || this.sort
                this.getData()
            }, { immediate: true })
        },
        beforeDestroy () {
            this.unwatch && this.unwatch()
        },
        methods: {
            handleCreate () {
                this.$refs.accountSideslider.show({
                    type: 'form',
                    title: this.$t('新建账户'),
                    props: {
                        mode: 'create'
                    }
                })
            },
            handleSortChange (sort) {
                RouterQuery.set({
                    _t: Date.now(),
                    sort: this.$tools.getSort(sort, { prop: 'bk_account_id' })
                })
            },
            handlePageChange (page) {
                RouterQuery.set({
                    _t: Date.now(),
                    page: page
                })
            },
            handleLimitChange (limit) {
                RouterQuery.set({
                    _t: Date.now(),
                    page: 1,
                    limit: limit
                })
            },
            handleCellClick (row, column) {
                if (column.property === 'bk_account_name') {
                    this.handleView(row)
                }
            },
            handleView (row) {
                this.$refs.accountSideslider.show({
                    type: 'details',
                    title: `${this.$t('账户详情')} 【${row.bk_account_name}】`,
                    props: {
                        id: row.bk_account_id
                    }
                })
            },
            async handleDelete (row) {
                const infoInstance = this.$bkInfo({
                    title: this.$t('确认删除xx', { instance: row.bk_account_name }),
                    closeIcon: false,
                    confirmFn: async () => {
                        try {
                            await this.$store.dispatch('cloud/account/delete', { id: row.bk_account_id })
                            infoInstance.buttonLoading = true
                            this.$success('删除成功')
                            this.getData()
                            return true
                        } catch (e) {
                            console.error(e)
                            return false
                        } finally {
                            infoInstance.buttonLoading = false
                        }
                    }
                })
            },
            async getData () {
                try {
                    const params = {
                        page: {
                            ...this.$tools.getPageParams(this.pagination),
                            sort: this.sort
                        },
                        condition: {}
                    }
                    if (this.filter) {
                        params.condition.bk_account_name = this.filter
                        params.is_fuzzy = true
                    }
                    const data = await this.$store.dispatch('cloud/account/findMany', {
                        params: params,
                        config: {
                            requestId: this.request.search
                        }
                    })
                    if (data.count && !data.info.length) {
                        this.handlePageChange(this.pagination.current - 1)
                        return
                    }
                    this.list = data.info.map(account => ({ ...account, pending: true, status: 'normal', error_message: '' }))
                    this.pagination.count = data.count
                    this.list.length && this.getAccountStatus()
                } catch (e) {
                    console.error(e)
                    this.list = []
                    this.pagination.count = 0
                }
            },
            async getAccountStatus () {
                try {
                    const results = await this.$store.dispatch('cloud/account/getStatus', {
                        params: {
                            account_ids: this.list.map(account => account.bk_account_id)
                        },
                        config: {
                            cancelPrevious: true,
                            requestId: this.request.status
                        }
                    })
                    this.list.forEach(account => {
                        const status = results.find(result => result.bk_account_id === account.bk_account_id)
                        if (status && status.err_msg) {
                            account.status = 'error'
                            account.error_message = status.err_msg
                        } else {
                            account.status = 'normal'
                            account.error_message = ''
                        }
                        account.pending = false
                    })
                } catch (error) {
                    this.list.forEach(account => {
                        account.pending = false
                        account.status = 'fail'
                    })
                    console.error(error)
                }
            },
            getStatusText (status) {
                const textMap = {
                    normal: this.$t('正常'),
                    error: this.$t('异常'),
                    fail: '--'
                }
                return textMap[status]
            },
            linkResource () {
                this.$routerActions.redirect({
                    name: MENU_RESOURCE_CLOUD_RESOURCE,
                    history: true
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cloud-account-layout {
        padding: 0 20px;
    }
    .cloud-account-tips {
        margin-top: 10px;
    }
    .cloud-account-options {
        margin-top: 10px;
        display: flex;
        align-items: center;
        justify-content: space-between;
        .options-filter {
            width: 260px;
        }
    }
    .cloud-account-table {
        margin-top: 10px;
    }
    @mixin dot {
        content: "";
        display: inline-block;
        margin-right: 4px;
        width: 7px;
        height: 7px;
        border-radius: 50%;
    }
    .cloud-account-status {
        &.normal:before {
            @include dot;
            background-color: $successColor;
        }
        &.error:before {
            @include dot;
            background-color: $dangerColor;
        }
        &.pending {
            font-size: 0;
            &:before {
                content: "";
                display: inline-block;
                vertical-align: middle;
                width: 16px;
                height: 16px;
                margin: 2px 0;
                background-color: transparent;
                background-image: url("../../assets/images/icon/loading.svg");
            }
        }
    }
</style>
