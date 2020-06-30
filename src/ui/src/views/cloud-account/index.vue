<template>
    <div class="cloud-account-layout">
        <cmdb-tips class="cloud-account-tips" tips-key="cloud-account-tips">
            <i18n path="云账户提示语">
                <bk-button text size="small" place="resource" style="padding: 0" @click="linkResource">{{$t('云资源')}}</bk-button>
            </i18n>
        </cmdb-tips>
        <div class="cloud-account-options">
            <div class="options-left">
                <bk-button theme="primary" @click="handleCreate">{{$t('新建')}}</bk-button>
            </div>
            <div class="options-right">
                <bk-input class="options-filter" clearable
                    v-model.trim="filter"
                    :placeholder="$t('请输入xx', { name: $t('账户名称') })">
                </bk-input>
            </div>
        </div>
        <bk-table class="cloud-account-table" v-bkloading="{ isLoading: $loading(request.search) }"
            :data="list"
            :pagination="pagination"
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
            <bk-table-column :label="$t('修改人')" prop="bk_last_editor" show-overflow-tooltip>
                <template slot-scope="{ row }">{{row.bk_last_editor | formatter('singlechar')}}</template>
            </bk-table-column>
            <bk-table-column :label="$t('修改时间')" prop="last_time" sortable="custom">
                <template slot-scope="{ row }">{{row.last_time | formatter('time')}}</template>
            </bk-table-column>
            <bk-table-column :label="$t('操作')">
                <template slot-scope="{ row }">
                    <link-button class="mr10" @click="handleView(row)">{{$t('查看')}}</link-button>
                    <link-button
                        :disabled="!row.bk_can_delete_account"
                        v-bk-tooltips="{
                            disabled: row.bk_can_delete_account,
                            content: $t('云账户禁止删除提示')
                        }"
                        @click="handleDelete(row)">
                        {{$t('删除')}}
                    </link-button>
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
                    search: Symbol('search')
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
            this.getData()
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
                this.sort = this.$tools.getSort(sort, { prop: 'bk_account_id' })
                this.getData()
            },
            handlePageChange (page) {
                this.pagination.current = page
                this.getData()
            },
            handleLimitChange (limit) {
                this.pagination.limit = limit
                this.pagination.current = 1
                this.getData()
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
                    const data = await this.$store.dispatch('cloud/account/findMany', {
                        params: {
                            page: {
                                ...this.$tools.getPageParams(this.pagination),
                                sort: this.sort
                            },
                            condition: {
                                bk_account_name: this.filter
                            }
                        },
                        config: {
                            requestId: this.request.search
                        }
                    })
                    if (data.count && !data.info.length) {
                        this.handlePageChange(this.pagination.current - 1)
                        return
                    }
                    this.list = data.info
                    this.pagination.count = data.count
                } catch (e) {
                    console.error(e)
                    this.list = []
                    this.pagination.count = 0
                }
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
            width: 200px;
        }
    }
    .cloud-account-table {
        margin-top: 10px;
    }
    .row-status {
        display: inline-block;
        .status {
            display: inline-block;
            margin-right: 4px;
            width: 7px;
            height: 7px;
            border-radius: 50%;
            background-color: $successColor;
        }
    }
</style>
