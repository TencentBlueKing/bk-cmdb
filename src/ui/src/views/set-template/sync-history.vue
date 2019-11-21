<template>
    <div class="sync-history-layout">
        <div class="options clearfix">
            <bk-date-picker style="width: 300px;" class="fl"
                type="daterange"
                :placeholder="$t('选择日期范围')"
                @change="hanldeFilterByDate">
            </bk-date-picker>
            <bk-input style="width: 240px;" class="fl ml10"
                right-icon="icon-search"
                v-model="searchName"
                clearable
                :placeholder="$t('集群名称')"
                @enter="getData(true)">
            </bk-input>
        </div>
        <bk-table class="history-table"
            v-bkloading="{ isLoading: $loading('getSyncHistory') }"
            :data="displayList"
            :pagination="pagination"
            :max-height="$APP.height - 229"
            @sort-change="handleSortChange"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange">
            <bk-table-column :label="$t('集群名称')" prop="bk_set_name"></bk-table-column>
            <bk-table-column :label="$t('拓扑路径')" prop="topo_path">
                <template slot-scope="{ row }">
                    <span>{{getTopoPath(row)}}</span>
                </template>
            </bk-table-column>
            <bk-table-column :label="$t('状态')" prop="status">
                <template slot-scope="{ row }">
                    <span v-if="row.status === 'syncing'" class="sync-status">
                        <img class="svg-icon" src="../../assets/images/icon/loading.svg" alt="">
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
                    <span v-else-if="row.status === 'failure'" class="sync-status">
                        <i class="status-circle fail"></i>
                        {{$t('同步失败')}}
                    </span>
                    <span v-else>--</span>
                </template>
            </bk-table-column>
            <bk-table-column :label="$t('同步时间')" prop="last_time" sortable="custom">
                <template slot-scope="{ row }">
                    <span>{{row.last_time ? $tools.formatTime(row.last_time, 'YYYY-MM-DD HH:mm:ss') : '--'}}</span>
                </template>
            </bk-table-column>
            <bk-table-column :label="$t('同步人')" prop="sync_user">
                <template slot-scope="{ row }">
                    <span>{{row.creator || '--'}}</span>
                </template>
            </bk-table-column>
            <cmdb-table-empty slot="empty" :stuff="table.stuff"></cmdb-table-empty>
        </bk-table>
    </div>
</template>

<script>
    import { MENU_BUSINESS_SET_TEMPLATE } from '@/dictionary/menu-symbol'
    export default {
        data () {
            return {
                templateName: '',
                searchName: '',
                searchDate: [],
                list: [],
                listWithTopo: [],
                pagination: {
                    count: 0,
                    current: 1,
                    ...this.$tools.getDefaultPaginationConfig()
                },
                table: {
                    stuff: {
                        type: 'default',
                        payload: {
                            emptyText: this.$t('bk.table.emptyText')
                        }
                    }
                },
                listSort: 'last_time'
            }
        },
        computed: {
            business () {
                return this.$store.state.objectBiz.bizId
            },
            templateId () {
                return this.$route.params.templateId
            },
            setsId () {
                const ids = this.list.map(item => item.bk_set_id)
                return [...new Set(ids)]
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
                    return {
                        ...item,
                        ...otherParams
                    }
                })
            },
            searchParams () {
                const params = {
                    set_template_id: Number(this.templateId),
                    search: this.searchName,
                    page: {
                        start: this.pagination.limit * (this.pagination.current - 1),
                        limit: this.pagination.limit,
                        sort: this.listSort
                    }
                }
                if (this.searchDate.length) {
                    params.start_time = this.searchDate[0] || ''
                    params.end_time = this.searchDate[1] || ''
                }
                return params
            }
        },
        created () {
            this.getSetTemplateInfo()
            this.getData()
        },
        methods: {
            setBreadcrumbs () {
                this.$store.commit('setBreadcrumbs', [{
                    label: this.$t('集群模板'),
                    route: {
                        name: MENU_BUSINESS_SET_TEMPLATE
                    }
                }, {
                    label: this.templateName,
                    route: {
                        name: 'setTemplateConfig',
                        params: {
                            mode: 'view',
                            templateId: this.templateId
                        },
                        query: {
                            tab: 'instance'
                        }
                    }
                }, {
                    label: this.$t('同步历史')
                }])
            },
            getTopoPath (row) {
                const topoPath = this.$tools.clone(row.topo_path)
                if (topoPath.length) {
                    const setIndex = topoPath.findIndex(path => path.ObjectID === 'set')
                    if (setIndex > -1) {
                        topoPath.splice(setIndex, 1)
                    }
                    const sortPath = topoPath.sort((prev, next) => prev.bk_inst_id - next.bk_inst_id)
                    return sortPath.map(path => path.bk_inst_name).join(' / ')
                }
                return '--'
            },
            async getData (event) {
                await this.getHistoryList()

                if (event) {
                    this.table.stuff.type = 'search'
                }

                this.setsId.length && this.getSetInstancesWithTopo()
            },
            async getSetTemplateInfo () {
                try {
                    const info = await this.$store.dispatch('setTemplate/getSingleSetTemplateInfo', {
                        bizId: this.$store.getters['objectBiz/bizId'],
                        setTemplateId: this.templateId
                    })
                    this.templateName = info.name
                    this.setBreadcrumbs()
                } catch (e) {
                    console.error(e)
                }
            },
            async getHistoryList () {
                try {
                    const data = await this.$store.dispatch('setTemplate/getSyncHistory', {
                        bizId: this.business,
                        params: this.searchParams,
                        config: {
                            requestId: 'getSyncHistory'
                        }
                    })
                    this.pagination.count = data.count
                    this.list = data.info || []
                } catch (e) {
                    console.error(e)
                    this.list = []
                }
            },
            async getSetInstancesWithTopo () {
                try {
                    const data = await this.$store.dispatch('setTemplate/getSetInstancesWithTopo', {
                        bizId: this.business,
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
            handleSortChange (sort) {
                this.listSort = this.$tools.getSort(sort)
                this.handlePageChange(1)
            },
            handlePageChange (current) {
                this.pagination.current = current
                this.getData()
            },
            handleSizeChange (size) {
                this.pagination.limit = size
                this.handlePageChange(1)
            },
            hanldeFilterByDate (daterange) {
                daterange = daterange.filter(date => date)
                this.searchDate = daterange.map((date, index) => {
                    return index === 0 ? (date + ' 00:00:00') : (date + ' 23:59:59')
                })
                this.getData(true)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .sync-history-layout {
        padding: 0 20px;
    }
    .options {
        padding-bottom: 15px;
    }
    .history-table {
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
</style>
