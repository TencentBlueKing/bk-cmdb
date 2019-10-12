<template>
    <div class="sync-history-layout">
        <div class="options clearfix">
            <bk-date-picker style="width: 300px;" class="fl"
                type="daterange"
                :placeholder="$t('选择日期范围')"
                v-model="searchDate">
            </bk-date-picker>
            <bk-input style="width: 240px;" class="fl ml10"
                right-icon="icon-search"
                v-model="searchName"
                :placeholder="$t('集群名称')">
            </bk-input>
        </div>
        <bk-table v-bkloading="{ isLoading: $loading('getSyncHistory') }"
            :data="list"
            :pagination="pagination"
            :max-height="$APP.height - 229"
            @sort-change="handleSortChange"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange">
            <bk-table-column :label="$t('集群名称')" prop="bk_set_name"></bk-table-column>
            <bk-table-column :label="$t('拓扑结构')" prop="topo_path">
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
                    <span v-else-if="row.status === 'waiting'">
                        {{$t('待同步')}}
                    </span>
                    <span v-else-if="row.status === 'finished'" class="sync-status success">
                        <i class="bk-icon icon-check-1"></i>
                        {{$t('已同步')}}
                    </span>
                    <span v-else-if="row.status === 'failure'" class="sync-status fail">
                        <i class="bk-icon icon-cc-log-02 "></i>
                        {{$t('同步失败')}}
                    </span>
                    <span v-else>--</span>
                </template>
            </bk-table-column>
            <bk-table-column :label="$t('同步时间')" prop="sync_time" sortable="custom">
                <template slot-scope="{ row }">
                    <span>{{row.last_time ? $tools.formatTime(row.last_time, 'YYYY-MM-DD HH:mm:ss') : '--'}}</span>
                </template>
            </bk-table-column>
            <bk-table-column :label="$t('同步人')" prop="sync_user">
                <template slot-scope="{ row }">
                    <span>{{row.creator || '--'}}</span>
                </template>
            </bk-table-column>
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
                pagination: {
                    count: 0,
                    current: 1,
                    ...this.$tools.getDefaultPaginationConfig()
                },
                listSort: 'sync_time'
            }
        },
        computed: {
            templateId () {
                return this.$route.params.templateId
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
                return params
            }
        },
        created () {
            this.getSetTemplateInfo()
            this.getHistoryList()
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
                    const sortPath = topoPath.sort((prev, next) => prev.InstanceID - next.InstanceID)
                    return sortPath.map(path => path.InstanceName).join(' / ')
                }
                return '--'
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
                        bizId: this.$store.state.objectBiz.bizId,
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
            handleSortChange (sort) {
                this.listSort = this.$tools.getSort(sort)
                this.handlePageChange(1)
            },
            handlePageChange (current) {
                this.pagination.current = current
                this.getHistoryList()
            },
            handleSizeChange (size) {
                this.pagination.limit = size
                this.getHistoryList(1)
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
</style>
