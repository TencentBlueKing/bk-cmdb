<template>
    <div class="template-instance-layout">
        <div class="instance-empty" v-if="!list.length">
            <img src="../../../assets/images/empty-content.png" alt="">
            <i18n path="空集群模板实例提示" tag="div">
                <bk-button text @click="handleLinkServiceTopo" place="link">{{$t('服务拓扑')}}</bk-button>
            </i18n>
        </div>
        <div class="instance-main" v-else>
            <div class="options clearfix">
                <div class="fl">
                    <bk-button theme="primary" :disabled="!checkedList.length" @click="handleBatchSync">{{$t('批量同步')}}</bk-button>
                </div>
                <div class="fr">
                    <bk-select class="filter-item mr10"
                        :clearable="false"
                        v-model="filterKey">
                        <bk-option v-for="option in filterList"
                            :key="option.id"
                            :id="option.id"
                            :name="option.name">
                        </bk-option>
                    </bk-select>
                    <bk-input class="filter-item" right-icon="bk-icon icon-search"
                        :placeholder="$t('路径关键字')"
                        v-model="pathKeywords"
                        @enter="handleSearch">
                    </bk-input>
                    <icon-button class="ml10"
                        v-bk-tooltips="$t('同步历史')"
                        icon="icon icon-cc-history"
                        @click="routeToHistory">
                    </icon-button>
                </div>
            </div>
            <bk-table class="instance-table"
                v-bkloading="{ isLoading: $loading('getSetTemplateInstances') }"
                :data="displayList"
                :pagination="pagination"
                :max-height="$APP.height - 229"
                @sort-change="handleSortChange"
                @page-change="handlePageChange"
                @page-limit-change="handleSizeChange"
                @selection-change="handleSelectionChange">
                <bk-table-column type="selection" width="50" :selectable="handleSelectable"></bk-table-column>
                <bk-table-column :label="$t('拓扑结构')" prop="topo_path">
                    <template slot-scope="{ row }">
                        <span>{{getTopoPath(row)}}</span>
                    </template>
                </bk-table-column>
                <bk-table-column :label="$t('主机数量')" prop="host_count"></bk-table-column>
                <bk-table-column :label="$t('状态')" prop="status">
                    <template slot-scope="{ row }">
                        <span v-if="!row.details">--</span>
                        <span v-else-if="statusMap.syncing.includes(row.details.status)" class="sync-status">
                            <img class="svg-icon" src="../../../assets/images/icon/loading.svg" alt="">
                            {{$t('同步中')}}
                        </span>
                        <span v-else-if="statusMap.success.includes(row.details.status) && row.hasDiff">
                            {{$t('待同步')}}
                        </span>
                        <span v-else-if="statusMap.success.includes(row.details.status)" class="sync-status success">
                            <i class="bk-icon icon-check-1"></i>
                            {{$t('已完成')}}
                        </span>
                        <span v-else class="sync-status fail">
                            <i class="bk-icon icon-cc-log-02 "></i>
                            {{$t('同步失败')}}
                        </span>
                    </template>
                </bk-table-column>
                <bk-table-column :label="$t('上次同步时间')" prop="sync_time" sortable="custom">
                    <template slot-scope="{ row }">
                        <span v-if="!row.details">--</span>
                        <span v-else>{{row.details.last_time ? $tools.formatTime(row.details.last_time, 'YYYY-MM-DD HH:mm') : '--'}}</span>
                    </template>
                </bk-table-column>
                <bk-table-column :label="$t('同步人')" prop="sync_user">
                    <template slot-scope="{ row }">
                        <span v-if="!row.details">--</span>
                        <span v-else>{{row.details.user || '--'}}</span>
                    </template>
                </bk-table-column>
                <bk-table-column :label="$t('操作')" width="180">
                    <template slot-scope="{ row }">
                        <bk-button v-if="row.details && statusMap.fail.includes(row.details.status)" text @click="handleRetry(row)">{{$t('重试')}}</bk-button>
                        <bk-button v-else text
                            :disabled="getEditabledStatus(row)"
                            @click="handleSync(row)">
                            {{$t('同步')}}
                        </bk-button>
                    </template>
                </bk-table-column>
            </bk-table>
        </div>
    </div>
</template>

<script>
    import { MENU_BUSINESS_SERVICE_TOPOLOGY } from '@/dictionary/menu-symbol'
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
                statusMap: {
                    syncing: [0, 1, 100],
                    success: [200],
                    fail: [500]
                },
                list: [],
                diffList: [],
                syncMap: {},
                checkedList: [],
                filterKey: -1,
                pathKeywords: '',
                pagination: {
                    count: 0,
                    current: 1,
                    ...this.$tools.getDefaultPaginationConfig()
                },
                listSort: ''
            }
        },
        computed: {
            business () {
                return this.$store.state.objectBiz.bizId
            },
            limit () {
                return {
                    start: this.pagination.limit * (this.pagination.current - 1),
                    limit: this.pagination.limit
                }
            },
            filterList () {
                return [{
                    id: -1,
                    name: this.$t('全部')
                }, {
                    id: 0,
                    name: this.$t('待同步')
                }, {
                    id: 1,
                    name: this.$t('同步中')
                }, {
                    id: 2,
                    name: this.$t('同步失败')
                }, {
                    id: 3,
                    name: this.$t('已同步')
                }]
            },
            setsId () {
                return this.list.map(item => item.bk_set_id)
            },
            displayList () {
                const list = this.$tools.clone(this.list)
                return list.map(item => {
                    let hasDiff = false
                    const diff = this.diffList.find(diff => diff.bk_set_id === item.bk_set_id)
                    if (diff) {
                        const diffModules = (diff.module_diffs || []).filter(module => module.diff_type !== 'unchanged')
                        hasDiff = !!diffModules.length
                    }
                    return {
                        hasDiff,
                        ...item,
                        details: this.syncMap[item.bk_set_id] ? this.syncMap[item.bk_set_id] : null
                    }
                })
            }
        },
        async created () {
            await this.getSetTemplateInstances()
            if (this.list.length) {
                this.getSyncStatus()
                this.getTemplateDiff()
                this.polling()
            }
        },
        beforeDestroy () {
            clearInterval(this.timer)
        },
        methods: {
            getEditabledStatus (row) {
                return (row.details && this.statusMap.syncing.includes(row.details.status)) || !row.hasDiff
            },
            getTopoPath (row) {
                const topoPath = this.$tools.clone(row.topo_path)
                if (topoPath.length) {
                    return topoPath.reverse().map(path => path.InstanceName).join(' / ')
                }
                return '--'
            },
            async getSetTemplateInstances () {
                const data = await this.$store.dispatch('setTemplate/getSetTemplateInstances', {
                    bizId: this.business,
                    setTemplateId: this.templateId,
                    params: {
                        limit: this.limit,
                        sort: this.listSort
                    },
                    config: {
                        requestId: 'getSetTemplateInstances'
                    }
                })
                this.pagination.count = data.count
                this.list = data.info || []
            },
            async getSyncStatus () {
                try {
                    const data = await this.$store.dispatch('setSync/getInstancesSyncStatus', {
                        bizId: this.business,
                        setTemplateId: this.templateId,
                        params: {
                            bk_set_ids: this.setsId
                        },
                        config: {
                            requestId: 'getInstancesSyncStatus',
                            cancelPrevious: true
                        }
                    })
                    this.syncMap = data || {}
                } catch (e) {
                    console.error(e)
                    clearInterval(this.timer)
                    this.timer = null
                    this.syncMap = {}
                }
            },
            async getTemplateDiff () {
                try {
                    const data = await this.$store.dispatch('setSync/diffTemplateAndInstances', {
                        bizId: this.business,
                        setTemplateId: this.templateId,
                        params: {
                            bk_set_ids: this.setsId
                        },
                        config: {
                            requestId: 'diffTemplateAndInstances'
                        }
                    })
                    this.diffList = data || []
                } catch (e) {
                    console.error(e)
                    clearInterval(this.timer)
                    this.timer = null
                    this.diffList = []
                }
            },
            polling () {
                try {
                    if (this.timer) {
                        clearInterval(this.timer)
                        this.timer = null
                    }
                    this.timer = setInterval(() => {
                        this.getSyncStatus()
                        this.getTemplateDiff()
                    }, 10000)
                } catch (e) {
                    console.error(e)
                    clearInterval(this.timer)
                    this.timer = null
                }
            },
            handleLinkServiceTopo () {
                this.$router.push({ name: MENU_BUSINESS_SERVICE_TOPOLOGY })
            },
            handleSearch () {},
            handleSelectionChange (selection) {
                this.checkedList = selection.map(item => item.bk_set_id)
            },
            handleSortChange (sort) {
                this.listSort = this.$tools.getSort(sort)
                this.handlePageChange(1)
            },
            handlePageChange (current) {
                if (this.timer === null) {
                    this.polling()
                }
                this.pagination.current = current
                this.getSetTemplateInstances()
            },
            handleSizeChange (size) {
                this.pagination.limit = size
                this.handlePageChange(1)
            },
            handleSelectable (row) {
                if (row.details && this.statusMap.syncing.includes(row.details.status)) {
                    return false
                }
                return row.hasDiff
            },
            handleBatchSync () {
                this.$store.commit('setFeatures/setSyncIdMap', {
                    id: `${this.business}_${this.templateId}`,
                    instancesId: this.checkedList
                })
                this.$router.push({
                    name: 'setSync',
                    params: {
                        setTemplateId: this.templateId
                    }
                })
            },
            handleSync (row) {
                this.$store.commit('setFeatures/setSyncIdMap', {
                    id: `${this.business}_${this.templateId}`,
                    instancesId: [row.bk_set_id]
                })
                this.$router.push({
                    name: 'setSync',
                    params: {
                        setTemplateId: this.templateId
                    }
                })
            },
            async handleRetry (row) {
                try {
                    await this.$store.dispatch('setSync/syncTemplateToInstances', {
                        bizId: this.business,
                        setTemplateId: this.templateId,
                        params: {
                            bk_set_ids: [row.bk_set_id]
                        },
                        config: {
                            requestId: 'syncTemplateToInstances'
                        }
                    })
                    this.$success(this.$t('重试同步中'))
                } catch (e) {
                    console.error(e)
                }
            },
            routeToHistory () {
                this.$router.push({
                    name: 'syncHistory',
                    params: {
                        templateId: this.templateId
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
            .sync-status {
                color: #63656E;
                .bk-icon {
                    margin-top: -2px;
                }
                .svg-icon {
                    @include inlineBlock;
                    margin-top: -4px;
                    width: 16px;
                }
                &.fail {
                    color: #EA3536;
                    .bk-icon {
                        color: #63656E;
                    }
                }
                &.success {
                    color: #2DCB56;
                }
            }
        }
    }
</style>
