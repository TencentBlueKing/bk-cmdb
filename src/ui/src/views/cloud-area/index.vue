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
                right-icon="icon-search"
                :placeholder="$t('请输入xx', { name: $t('云区域名称') })">
            </bk-input>
        </div>
        <bk-table class="cloud-area-table"
            v-bkloading="{ isLoading: $loading(request.search) }"
            :data="list"
            :pagination="pagination"
            :max-height="$APP.height - 220"
            @sort-change="handleSortChange"
            @page-change="handlePageChange"
            @page-limit-change="handleLimitChange">
            <bk-table-column
                sortable="custom"
                prop="bk_cloud_name"
                :label="$t('云区域名称')">
                <template slot-scope="{ row }">
                    <cmdb-auth class="cell-name" v-show="row !== rowInEdit"
                        :auth="{ type: $OPERATION.U_CLOUD_AREA, relation: [row.bk_cloud_id] }"
                        :class="{
                            pending: row._pending_,
                            limited: isLimited(row)
                        }"
                        @click="handleEditName($event, row)">
                        <template slot-scope="{ disabled }">
                            <span class="cell-name-text" v-bk-overflow-tips>{{row.bk_cloud_name}}</span>
                            <i :class="['cell-name-icon', 'icon-cc-edit-shape', { disabled }]" v-if="!isLimited(row)"></i>
                        </template>
                    </cmdb-auth>
                    <bk-input class="cell-name-input" size="small" font-size="normal" :value="row.bk_cloud_name" v-if="row === rowInEdit"
                        @enter="handleUpdateName(row, ...arguments)"
                        @blur="handleUpdateName(row, ...arguments)">
                    </bk-input>
                </template>
            </bk-table-column>
            <bk-table-column :label="$t('状态')" prop="bk_status" sortable="custom">
                <div class="row-status" slot-scope="{ row }"
                    v-bk-tooltips="{
                        content: row.bk_status_detail,
                        disabled: row.bk_status === '1' || !row.bk_status_detail
                    }">
                    <i :class="['status', { 'is-error': row.bk_status !== '1' }]"></i>
                    {{row.bk_status === '1' ? $t('正常') : $t('异常')}}
                </div>
            </bk-table-column>
            <bk-table-column :label="$t('所属云厂商')" prop="bk_cloud_vendor" sortable="custom">
                <cmdb-vendor slot-scope="{ row }" :type="row.bk_cloud_vendor"></cmdb-vendor>
            </bk-table-column>
            <bk-table-column :label="$t('地域')" prop="bk_region">
                <template slot-scope="{ row }">
                    <task-region-selector display="info" :account="row.bk_account_id" :value="row.bk_region"></task-region-selector>
                </template>
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
                <template slot-scope="{ row }">
                    <cmdb-auth :auth="{ type: $OPERATION.D_CLOUD_AREA, relation: [row.bk_cloud_id] }">
                        <link-button slot-scope="{ disabled }"
                            :disabled="!isRemovable(row) || disabled"
                            v-bk-tooltips="{
                                disabled: isRemovable(row) || !disabled,
                                content: getRemoveTips(row)
                            }"
                            @click="handleDelete(row)">
                            {{$t('删除')}}
                        </link-button>
                    </cmdb-auth>
                </template>
            </bk-table-column>
        </bk-table>
    </div>
</template>

<script>
    import CmdbVendor from '@/components/ui/other/vendor'
    import throttle from 'lodash.throttle'
    import { MENU_RESOURCE_CLOUD_RESOURCE } from '@/dictionary/menu-symbol'
    import TaskRegionSelector from '@/views/cloud-resource/children/task-region-selector'
    export default {
        components: {
            CmdbVendor,
            TaskRegionSelector
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
                scheduleSearch: throttle(this.handlePageChange, 800, { leading: false, trailing: true }),
                rowInEdit: null
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
            handleEditName (event, row) {
                const cell = event.currentTarget.parentElement
                this.rowInEdit = row
                this.$nextTick(() => {
                    cell.querySelector('input').focus()
                })
            },
            async handleUpdateName (row, value) {
                try {
                    value = value.trim()
                    this.rowInEdit = null
                    if (row.bk_cloud_name === value) {
                        return
                    }
                    this.$set(row, '_pending_', true)
                    await this.$store.dispatch('cloud/area/update', {
                        id: row.bk_cloud_id,
                        params: {
                            bk_cloud_name: value
                        }
                    })
                    row.bk_cloud_name = value
                    this.$delete(row, '_pending_')
                } catch (error) {
                    console.error(error)
                }
            },
            isRemovable (row) {
                return row.host_count === 0 && !this.isLimited(row) && row.sync_task_ids.length === 0
            },
            getRemoveTips (row) {
                if (this.isLimited(row)) {
                    return this.$t('系统限定，不能删除')
                }
                if (row.host_count !== 0) {
                    return this.$t('主机不为空，不能删除')
                }
                if (row.sync_task_ids.length !== 0) {
                    return this.$t('已关联同步任务，不能删除')
                }
                return null
            },
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
                    const params = {
                        page: {
                            ...this.$tools.getPageParams(this.pagination),
                            sort: this.sort
                        },
                        host_count: true,
                        condition: {},
                        sync_task_ids: true
                    }
                    if (this.filter) {
                        params.condition.bk_cloud_name = this.filter
                    }
                    const data = await this.$store.dispatch('cloud/area/findMany', {
                        params: params,
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
            width: 260px;
        }
    }
    .cloud-area-tips {
        margin-top: 10px;
    }
    .cloud-area-table {
        margin-top: 10px;
        .cell-name {
            display: flex;
            align-items: center;
            margin: 0 -5px;
            padding: 0 5px;
            height: 26px;
            line-height: 24px;
            border: 1px solid transparent;
            border-radius: 2px;
            cursor: text;
            &:not(.limited):hover {
                background-color: #DCDEE5;
            }
            &.limited {
                pointer-events: none;
            }
            &.pending {
                pointer-events: none;
                font-size: 0;
                &:before {
                    content: "";
                    display: inline-block;
                    vertical-align: middle;
                    width: 16px;
                    height: 16px;
                    margin: 2px 0;
                    background-image: url("../../assets/images/icon/loading.svg");
                }
            }
            .cell-name-icon {
                display: none;
                flex: 14px 0 0;
                font-size: 14px;
                margin-left: 6px;
                color: $primaryColor;
                cursor: pointer;
                &:hover {
                    opacity: .75;
                }
                &.disabled {
                    color: $textDisabledColor;
                }
            }
            .cell-name-text {
                display: block;
                @include ellipsis;
            }
        }
        .cell-name-input {
            display: block;
            width: auto;
            margin: 0 -5px;
        }
        /deep/ .bk-table-row:hover {
            .cell-name-icon {
                display: inline;
            }
        }
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
            &.is-error {
                background-color: $dangerColor;
            }
        }
    }
</style>
