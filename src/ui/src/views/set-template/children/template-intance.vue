<template>
    <div class="template-instance-layout">
        <div class="instance-main">
            <div class="options clearfix">
                <div class="fl">
                    <cmdb-auth :auth="{ type: $OPERATION.U_TOPO, relation: [bizId] }">
                        <bk-button slot-scope="{ disabled }"
                            theme="primary"
                            :disabled="disabled || !checkedList.length"
                            @click="handleBatchSync">
                            {{$t('批量同步')}}
                        </bk-button>
                    </cmdb-auth>
                </div>
                <div class="fr">
                    <bk-select class="filter-item mr10"
                        :clearable="false"
                        v-model="statusFilter"
                        @selected="handleFilter(1, true)">
                        <bk-option v-for="option in filterList"
                            :key="option.id"
                            :id="option.id"
                            :name="option.name">
                        </bk-option>
                    </bk-select>
                    <bk-input class="filter-item" right-icon="bk-icon icon-search"
                        clearable
                        v-model="filterName"
                        :placeholder="$t('请输入集群名称搜索')"
                        @enter="handleSearch"
                        @clear="handleSearch">
                    </bk-input>
                    <icon-button class="ml10"
                        v-bk-tooltips="$t('同步历史')"
                        icon="icon icon-cc-history"
                        @click="routeToHistory">
                    </icon-button>
                </div>
            </div>
            <bk-table class="instance-table"
                ref="instanceTable"
                v-bkloading="{ isLoading: $loading('getSetInstanceData') }"
                :data="displayList"
                :pagination="pagination"
                :max-height="$APP.height - 249"
                @sort-change="handleSortChange"
                @page-change="handlePageChange"
                @page-limit-change="handleSizeChange"
                @selection-change="handleSelectionChange">
                <bk-table-column type="selection" width="50" :selectable="handleSelectable" align="center"></bk-table-column>
                <bk-table-column :label="$t('集群名称')" prop="bk_set_name" show-overflow-tooltip></bk-table-column>
                <bk-table-column :label="$t('拓扑路径')" prop="topo_path" show-overflow-tooltip>
                    <span class="topo-path" slot-scope="{ row }" @click="handlePathClick(row)">{{getTopoPath(row)}}</span>
                </bk-table-column>
                <bk-table-column :label="$t('主机数量')" prop="host_count"></bk-table-column>
                <bk-table-column :label="$t('状态')" prop="status">
                    <template slot-scope="{ row }">
                        <span v-if="row.status === 'syncing'" class="sync-status">
                            <img class="svg-icon" src="../../../assets/images/icon/loading.svg" alt="">
                            {{$t('同步中')}}
                        </span>
                        <span v-else-if="row.status === 'waiting'" class="sync-status">
                            <i class="status-circle waiting"></i>
                            {{$t('待同步')}}
                        </span>
                        <span v-else-if="row.status === 'finished'" class="sync-status">
                            <i class="status-circle success"></i>
                            {{$t('已同步')}}
                        </span>
                        <span v-else-if="row.status === 'failure'"
                            class="sync-status"
                            v-bk-tooltips="{
                                content: row.fail_tips,
                                placement: 'right'
                            }">
                            <i class="status-circle fail"></i>
                            {{$t('同步失败')}}
                        </span>
                        <span v-else>--</span>
                    </template>
                </bk-table-column>
                <bk-table-column :label="$t('上次同步时间')" prop="last_time" sortable="custom" show-overflow-tooltip>
                    <template slot-scope="{ row }">
                        <span>{{row.last_time ? $tools.formatTime(row.last_time, 'YYYY-MM-DD HH:mm:ss') : '--'}}</span>
                    </template>
                </bk-table-column>
                <bk-table-column :label="$t('同步人')" prop="sync_user">
                    <template slot-scope="{ row }">
                        <span>{{row.creator || '--'}}</span>
                    </template>
                </bk-table-column>
                <bk-table-column :label="$t('操作')" width="180">
                    <template slot-scope="{ row }">
                        <cmdb-auth :auth="{ type: $OPERATION.U_TOPO, relation: [bizId] }">
                            <template slot-scope="{ disabled }">
                                <bk-button v-if="row.status === 'failure'"
                                    text
                                    :disabled="disabled"
                                    @click="handleRetry(row)">
                                    {{$t('重试')}}
                                </bk-button>
                                <bk-button v-else
                                    text
                                    :disabled="disabled || ['syncing', 'finished'].includes(row.status)"
                                    @click="handleSync(row)">
                                    {{$t('去同步')}}
                                </bk-button>
                            </template>
                        </cmdb-auth>
                    </template>
                </bk-table-column>
                <cmdb-table-empty slot="empty" :stuff="table.stuff">
                    <div>
                        <i18n path="空集群模板实例提示" tag="div">
                            <bk-button style="font-size: 14px;" text @click="handleLinkServiceTopo" place="link">{{$t('业务拓扑')}}</bk-button>
                        </i18n>
                    </div>
                </cmdb-table-empty>
            </bk-table>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import { MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'
    export default {
        props: {
            templateId: {
                type: [Number, String],
                required: true
            }
        },
        data () {
            return {
                timer: null,
                list: [],
                listWithTopo: [],
                checkedList: [],
                statusFilter: 'all',
                filterName: '',
                pagination: {
                    count: 0,
                    current: 1,
                    ...this.$tools.getDefaultPaginationConfig()
                },
                table: {
                    stuff: {
                        type: 'default',
                        payload: {}
                    }
                },
                listSort: 'last_time',
                instancesInfo: {}
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            filterList () {
                return [{
                    id: 'all',
                    name: this.$t('全部')
                }, {
                    id: 'waiting',
                    name: this.$t('待同步')
                }, {
                    id: 'syncing',
                    name: this.$t('同步中')
                }, {
                    id: 'failure',
                    name: this.$t('同步失败')
                }, {
                    id: 'finished',
                    name: this.$t('已同步')
                }]
            },
            setsId () {
                return this.list.map(item => item.bk_set_id)
            },
            displayList () {
                const list = this.$tools.clone(this.list)
                return list.map(item => {
                    const otherParams = {
                        topo_path: [],
                        host_count: 0
                    }
                    const setInfo = this.listWithTopo.find(set => set.bk_set_id === item.bk_set_id)
                    if (setInfo) {
                        otherParams.topo_path = setInfo.topo_path || []
                        otherParams.host_count = setInfo.host_count || 0
                    }
                    const templateSyncInfo = this.instancesInfo[item.bk_set_id]
                    const syncFail = 500
                    let failTips = []
                    if (templateSyncInfo && templateSyncInfo.status === syncFail) {
                        const failModules = (templateSyncInfo.detail || []).filter(module => module.status === syncFail)
                        failTips = failModules.map(module => {
                            const data = module.data
                            const moduleName = data && data.module_diff && data.module_diff.bk_module_name
                            const errorMsg = module.response && module.response.bk_error_msg
                            const tips = moduleName && errorMsg ? `${moduleName} : ${errorMsg}` : ''
                            return tips
                        }).filter(tips => tips)
                    }
                    return {
                        ...item,
                        ...otherParams,
                        fail_tips: failTips.join('<br />')
                    }
                })
            },
            searchParams () {
                const params = {
                    set_template_id: Number(this.templateId),
                    page: {
                        start: this.pagination.limit * (this.pagination.current - 1),
                        limit: this.pagination.limit,
                        sort: this.listSort
                    }
                }
                if (this.statusFilter !== 'all') {
                    params.status = this.statusFilter
                }
                this.filterName && (params.search = this.filterName)
                return params
            }
        },
        watch: {
            list () {
                if (!this.list.length) {
                    clearInterval(this.timer)
                    this.timer = null
                } else if (!this.timer) {
                    this.polling()
                }
            }
        },
        async created () {
            await this.getData()
            if (this.list.length) {
                this.getSetInstancesWithTopo()
                this.getInstancesInfo()
                this.polling()
            }
        },
        beforeDestroy () {
            clearInterval(this.timer)
        },
        methods: {
            getTopoPath (row) {
                return [...row.topo_path].reverse().map(path => path.bk_inst_name).join(' / ') || '--'
            },
            async getData () {
                const data = await this.getSetInstancesWithStatus('getSetInstanceData')
                this.pagination.count = data.count
                this.list = data.info || []
            },
            async updateStatusData () {
                const data = await this.getSetInstancesWithStatus('updateStatusData', {
                    page: {
                        start: 0,
                        limit: this.pagination.limit,
                        sort: this.listSort
                    },
                    bk_set_ids: this.setsId
                }, {
                    globalError: false
                })
                const backupList = this.$tools.clone(this.checkedList)
                this.list = data.info || []
                const table = this.$refs.instanceTable
                if (table) {
                    this.$nextTick(() => {
                        this.displayList.forEach(row => {
                            table.toggleRowSelection(row, backupList.includes(row.bk_set_id))
                        })
                    })
                }
            },
            getSetInstancesWithStatus (requestId, otherParams, config) {
                return this.$store.dispatch('setTemplate/getSetInstancesWithStatus', {
                    bizId: this.bizId,
                    params: {
                        ...this.searchParams,
                        ...(otherParams || {})
                    },
                    config: {
                        requestId,
                        cancelPrevious: true,
                        ...(config || {})
                    }
                })
            },
            async getSetInstancesWithTopo () {
                try {
                    const data = await this.$store.dispatch('setTemplate/getSetInstancesWithTopo', {
                        bizId: this.bizId,
                        setTemplateId: this.templateId,
                        params: {
                            limit: {
                                start: 0,
                                limit: this.pagination.limit
                            },
                            bk_set_ids: this.setsId
                        },
                        config: {
                            requestId: 'getSetInstancesWithTopo'
                        }
                    })
                    this.listWithTopo = data.info || []
                } catch (e) {
                    console.error(e)
                    this.listWithTopo = []
                }
            },
            async getInstancesInfo () {
                const instancesInfo = await this.$store.dispatch('setSync/getInstancesSyncStatus', {
                    bizId: this.bizId,
                    setTemplateId: this.templateId,
                    params: {
                        bk_set_ids: this.setsId
                    },
                    config: {
                        requestId: 'getInstancesInfo'
                    }
                })
                this.instancesInfo = instancesInfo
            },
            polling () {
                try {
                    if (this.timer) {
                        clearInterval(this.timer)
                        this.timer = null
                    }
                    this.timer = setInterval(() => {
                        this.updateStatusData()
                        this.getInstancesInfo()
                    }, 10000)
                } catch (e) {
                    console.error(e)
                    clearInterval(this.timer)
                    this.timer = null
                }
            },
            handleLinkServiceTopo () {
                this.$routerActions.redirect({ name: MENU_BUSINESS_HOST_AND_SERVICE })
            },
            async handleFilter (current = 1, event) {
                this.pagination.current = current
                await this.getData()

                const searchStatus = this.statusFilter !== 'all' || !!this.filterName
                this.table.stuff.type = event && searchStatus ? 'search' : 'default'

                if (this.list.length) {
                    this.getSetInstancesWithTopo()
                    this.getInstancesInfo()
                }
            },
            handleSelectionChange (selection) {
                this.checkedList = selection.map(item => item.bk_set_id)
            },
            handleSortChange (sort) {
                this.listSort = this.$tools.getSort(sort)
                this.handleFilter()
            },
            handlePageChange (current) {
                this.handleFilter(current)
            },
            handleSizeChange (size) {
                this.pagination.limit = size
                this.handlePageChange(1)
            },
            handleSearch (name) {
                this.filterName = name
                this.handleFilter(1, true)
            },
            handleSelectable (row) {
                return !['syncing', 'finished'].includes(row.status)
            },
            handleBatchSync () {
                this.$store.commit('setFeatures/setSyncIdMap', {
                    id: `${this.bizId}_${this.templateId}`,
                    instancesId: this.checkedList
                })
                this.$routerActions.redirect({
                    name: 'setSync',
                    params: {
                        setTemplateId: this.templateId
                    },
                    history: true
                })
            },
            handleSync (row) {
                this.$store.commit('setFeatures/setSyncIdMap', {
                    id: `${this.bizId}_${this.templateId}`,
                    instancesId: [row.bk_set_id]
                })
                this.$routerActions.redirect({
                    name: 'setSync',
                    params: {
                        setTemplateId: this.templateId
                    },
                    history: true
                })
            },
            async handleRetry (row) {
                try {
                    await this.$store.dispatch('setSync/syncTemplateToInstances', {
                        bizId: this.bizId,
                        setTemplateId: this.templateId,
                        params: {
                            bk_set_ids: [row.bk_set_id]
                        },
                        config: {
                            requestId: 'syncTemplateToInstances'
                        }
                    })
                    row.status = 'syncing'
                    this.$success(this.$t('提交同步成功'))
                } catch (e) {
                    console.error(e)
                }
            },
            routeToHistory () {
                this.$routerActions.redirect({
                    name: 'syncHistory',
                    params: {
                        templateId: this.templateId
                    },
                    history: true
                })
            },
            handlePathClick (row) {
                this.$routerActions.open({
                    name: MENU_BUSINESS_HOST_AND_SERVICE,
                    query: {
                        node: `set-${row.bk_set_id}`
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .template-instance-layout {
        height: 100%;
    }
    .instance-empty {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        height: 100%;
        font-size: 14px;
        text-align: center;
    }
    .instance-main {
        .options {
            padding: 15px 0;
            font-size: 14px;
            color: #63656E;
            .filter-item {
                width: 230px;
            }
            .icon-cc-updating {
                color: #C4C6CC;
            }
        }
        .instance-table {
            .topo-path {
                cursor: pointer;
                &:hover {
                    color: $primaryColor;
                }
            }
            .sync-status {
                color: #63656E;
                .status-circle {
                    display: inline-block;
                    width: 8px;
                    height: 8px;
                    margin-right: 4px;
                    border-radius: 50%;
                    &.waiting {
                        background-color: #3A84FF;
                    }
                    &.success {
                        background-color: #2DCB56;
                    }
                    &.fail {
                        background-color: #EA3536;
                    }
                }
                .svg-icon {
                    @include inlineBlock;
                    margin-top: -4px;
                    width: 16px;
                }
            }
        }
    }
</style>
