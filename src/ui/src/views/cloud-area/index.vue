<template>
    <div class="cloud-area-layout">
        <cmdb-tips class="cloud-area-tips" tips-key="cloud-area-tips">
            <i18n path="云区域提示语">
                <bk-button text size="small" place="resource" style="padding: 0" @click="linkResource">{{$t('云资源发现')}}</bk-button>
                <bk-button text size="small" place="agent" style="padding: 0" @click="linkAgent">{{$t('节点管理')}}</bk-button>
            </i18n>
        </cmdb-tips>
        <div class="cloud-area-options">
            <bk-input class="options-filter" clearable
                v-model.trim="filter"
                :placeholder="$t('请输入xx', { name: $t('云区域名称') })">
            </bk-input>
        </div>
        <bk-table class="cloud-area-table"
            v-bkloading="{ isLoading: $loading(request.search) }"
            :data="list"
            :pagination="pagination"
            @sort-change="handleSortChange"
            @page-change="handlePageChange"
            @page-limit-change="handleLimitChange">
            <bk-table-column class-name="is-highlight"
                show-overflow-tooltip
                sortable="custom"
                prop="bk_cloud_name"
                :label="$t('云区域名称')">
            </bk-table-column>
            <bk-table-column :label="$t('状态')" prop="bk_status" sortable="custom">
                <div class="row-status" slot-scope="{ row }">
                    <i :class="['status', { 'is-error': row.bk_status !== '1' }]"></i>
                    {{row.bk_status === '1' ? $t('正常') : $t('异常')}}
                </div>
            </bk-table-column>
            <bk-table-column :label="$t('所属云厂商')" prop="bk_cloud_vendor" sortable="custom">
                <cmdb-vendor slot-scope="{ row }" :type="row.bk_cloud_vendor"></cmdb-vendor>
            </bk-table-column>
            <bk-table-column :label="$t('地域')" prop="bk_region">
                <template slot-scope="{ row }">{{row.bk_region | formatter('singlechar')}}</template>
            </bk-table-column>
            <bk-table-column label="VPC" prop="bk_vpc_name" show-overflow-tooltip>
                <template slot-scope="{ row }">{{getVpcInfo(row) | formatter('singlechar')}}</template>
            </bk-table-column>
            <bk-table-column :label="$t('主机数量')" prop="host_count">
                <template slot-scope="{ row }">{{row.host_count | formatter('int')}}</template>
            </bk-table-column>
            <bk-table-column :label="$t('最近编辑')" prop="last_time" sortable="custom" show-overflow-tooltip>
                <template slot-scope="{ row }">{{row.last_time | formatter('time')}}</template>
            </bk-table-column>
            <bk-table-column :label="$t('编辑人')" prop="bk_last_editor"></bk-table-column>
            <bk-table-column :label="$t('操作')" fixed="right">
                <link-button slot-scope="{ row }"
                    :disabled="!!row.host_count || isLimited(row)"
                    v-bk-tooltips="{
                        disabled: isLimited(row) ? false : !row.host_count,
                        content: isLimited(row) ? $t('系统限定，不能删除') : $t('主机不为空，不能删除')
                    }"
                    @click="handleDelete(row)">
                    {{$t('删除')}}
                </link-button>
            </bk-table-column>
        </bk-table>
    </div>
</template>

<script>
    import CmdbVendor from '@/components/ui/other/vendor'
    import throttle from 'lodash.throttle'
    import { MENU_RESOURCE_CLOUD_RESOURCE } from '@/dictionary/menu-symbol'
    export default {
        components: {
            CmdbVendor
        },
        data () {
            return {
                filter: '',
                list: [],
                pagination: this.$tools.getDefaultPaginationConfig(),
                sort: 'bk_cloud_id',
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
            isLimited (row) {
                return row.bk_cloud_id === 0
            },
            handleSortChange (sort) {
                this.sort = this.$tools.getSort(sort, { prop: 'bk_cloud_id' })
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
            handleDelete (row) {
                const infoInstance = this.$bkInfo({
                    title: this.$t('确认删除xx', { instance: row.bk_cloud_name }),
                    closeIcon: false,
                    confirmFn: async () => {
                        try {
                            await this.$store.dispatch('cloud/area/delete', { id: row.bk_cloud_id })
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
                    const data = await this.$store.dispatch('cloud/area/findMany', {
                        params: {
                            page: {
                                ...this.$tools.getPageParams(this.pagination),
                                sort: this.sort
                            },
                            host_count: true,
                            condition: {
                                bk_cloud_name: this.filter
                            }
                        },
                        config: {
                            requestId: this.request.search,
                            cancelPrevious: true
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
            getVpcInfo (row) {
                const id = row.bk_vpc_id
                const name = row.bk_vpc_name
                if (name && id !== name) {
                    return `${id}(${name})`
                }
                return id
            },
            linkResource () {
                this.$routerActions.redirect({
                    name: MENU_RESOURCE_CLOUD_RESOURCE,
                    history: true
                })
            },
            linkAgent () {
                const agent = window.Site.agent
                if (agent) {
                    const topWindow = window.top
                    const isPaasConsole = topWindow !== window
                    if (isPaasConsole) {
                        topWindow.postMessage(JSON.stringify({
                            action: 'open_other_app',
                            app_code: 'bk_nodeman'
                        }), '*')
                    } else {
                        window.open(agent)
                    }
                } else {
                    this.$warn(this.$t('未配置Agent安装APP地址'))
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cloud-area-layout {
        padding: 0 20px;
    }
    .cloud-area-options {
        margin-top: 10px;
        .options-filter {
            width: 200px;
        }
    }
    .cloud-area-tips {
        margin-top: 10px;
    }
    .cloud-area-table {
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
