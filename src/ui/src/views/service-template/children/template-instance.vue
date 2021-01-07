<template>
    <div class="template-instance-layout">
        <div class="instance-main">
            <div class="options clearfix">
                <div class="fl">
                    <bk-button theme="primary"
                        :disabled="!table.selection.length"
                        @click="handleBatchSync">
                        {{$t('批量同步')}}
                    </bk-button>
                </div>
                <div class="fr">
                    <bk-input class="filter-item" right-icon="bk-icon icon-search"
                        clearable
                        :placeholder="$t('请输入拓扑路径')"
                        v-model.trim="table.filter"
                        @change="handleFilter">
                    </bk-input>
                </div>
            </div>
            <bk-table class="instance-table"
                ref="instanceTable"
                v-bkloading="{ isLoading: $loading(Object.values(request)) || table.filtering }"
                :data="table.data"
                :max-height="$APP.height - 250"
                @selection-change="handleSelectionChange">
                <bk-table-column type="selection" :selectable="checkSelectable"></bk-table-column>
                <bk-table-column :label="$t('模块名称')" prop="bk_module_name" show-overflow-tooltip></bk-table-column>
                <bk-table-column :label="$t('拓扑路径')" sortable :sort-method="sortByPath" show-overflow-tooltip>
                    <span slot-scope="{ row }" class="topo-path" @click="handlePathClick(row)">{{row._path_}}</span>
                </bk-table-column>
                <bk-table-column :label="$t('上次同步时间')" sortable :sort-method="sortByTime" show-overflow-tooltip>
                    <template slot-scope="{ row }">{{row.last_time | time}}</template>
                </bk-table-column>
                <bk-table-column :label="$t('操作')">
                    <template slot-scope="{ row }">
                        <span class="latest-sync" v-if="isSyncDisabled(row)" v-bk-tooltips="getSyncDisabledText(row)">{{$t('去同步')}}</span>
                        <bk-button v-else text :disabled="isSyncDisabled(row)" @click="handleSync(row)">{{$t('去同步')}}</bk-button>
                    </template>
                </bk-table-column>
                <cmdb-table-empty slot="empty" :stuff="table.stuff">
                    <div>
                        <i18n path="空服务模板实例提示" tag="div">
                            <bk-button style="font-size: 14px;" text @click="handleToCreatedInstance" place="link">{{$t('创建服务实例')}}</bk-button>
                        </i18n>
                    </div>
                </cmdb-table-empty>
            </bk-table>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import { time } from '@/filters/formatter'
    import debounce from 'lodash.debounce'
    import Bus from '@/utils/bus'
    import { MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'
    export default {
        filters: {
            time
        },
        props: {
            active: Boolean
        },
        data () {
            return {
                table: {
                    filter: '',
                    filtering: false,
                    selection: [],
                    data: [],
                    backup: [],
                    syncStatus: [],
                    stuff: {
                        type: 'default',
                        payload: {}
                    }
                },
                request: {
                    instance: Symbol('instance'),
                    status: Symbol('status'),
                    path: Symbol('path')
                },
                handleFilter: null
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            serviceTemplateId () {
                return this.$route.params.templateId
            }
        },
        watch: {
            active: {
                immediate: true,
                handler (active) {
                    active && this.refresh()
                }
            }
        },
        created () {
            this.handleFilter = debounce(this.filterData, 300)
        },
        methods: {
            async refresh () {
                try {
                    const data = await this.getTemplateInstance()
                    if (data.count) {
                        const [syncStatus, topoPath] = await Promise.all([
                            this.getSyncStatus(data.info),
                            this.getTopoPath(data.info)
                        ])
                        this.table.syncStatus = syncStatus
                        data.info.forEach(module => {
                            const topo = topoPath.nodes.find(topo => topo.topo_node.bk_inst_id === module.bk_module_id)
                            module._path_ = topo.topo_path.map(path => path.bk_inst_name).reverse().join(' / ')
                        })
                    }
                    this.table.data = Object.freeze(data.info)
                    this.table.backup = Object.freeze(data.info)
                    Bus.$emit('module-loaded', data.count)
                } catch (e) {
                    console.error(e)
                }
            },
            getTemplateInstance () {
                return this.$store.dispatch('serviceTemplate/getServiceTemplateModules', {
                    bizId: this.bizId,
                    serviceTemplateId: this.serviceTemplateId,
                    params: {},
                    config: {
                        requestId: this.request.instance
                    }
                })
            },
            getSyncStatus (modules) {
                return this.$store.dispatch('businessSynchronous/searchServiceInstanceDifferences', {
                    params: {
                        bk_module_ids: modules.map(module => module.bk_module_id),
                        service_template_id: this.serviceTemplateId,
                        bk_biz_id: this.bizId
                    },
                    config: {
                        requestId: this.request.status
                    }
                })
            },
            getTopoPath (modules) {
                return this.$store.dispatch('objectMainLineModule/getTopoPath', {
                    bizId: this.bizId,
                    params: {
                        topo_nodes: modules.map(module => ({ bk_obj_id: 'module', bk_inst_id: module.bk_module_id }))
                    },
                    config: {
                        requestId: this.request.path
                    }
                })
            },
            getSyncStatusDifference (row) {
                return this.table.syncStatus.find(difference => difference.bk_module_id === row.bk_module_id)
            },
            isSyncDisabled (row) {
                const difference = this.getSyncStatusDifference(row)
                if (difference) {
                    return !difference.has_difference
                }
                return true
            },
            getSyncDisabledText (row) {
                const difference = this.getSyncStatusDifference(row) || {}
                return difference.unchanged && difference.unchanged.length ? this.$t('已经是最新，无需同步') : this.$t('没有主机无需同步')
            },
            handleSync (row) {
                this.$routerActions.redirect({
                    name: 'syncServiceFromTemplate',
                    params: {
                        modules: row.bk_module_id,
                        template: this.serviceTemplateId
                    },
                    history: true
                })
            },
            checkSelectable (row) {
                return !this.isSyncDisabled(row)
            },
            handleSelectionChange (selection) {
                this.table.selection = selection
            },
            handleBatchSync () {
                this.$routerActions.redirect({
                    name: 'syncServiceFromTemplate',
                    params: {
                        template: this.serviceTemplateId,
                        modules: this.table.selection.map(row => row.bk_module_id).join(',')
                    },
                    history: true
                })
            },
            filterData () {
                this.table.filtering = true
                this.$nextTick(() => {
                    if (this.table.filter) {
                        this.table.data = this.table.backup.filter(row => {
                            const path = row._path_.replace(/\s*(\/)\s*/g, '$1')
                            const filter = this.table.filter.replace(/\s*(\/)\s*/g, '$1')
                            return path.indexOf(filter) > -1
                        })
                    } else {
                        this.table.data = [...this.table.backup]
                    }
                    this.table.stuff.type = this.table.filter ? 'search' : 'default'
                    this.table.filtering = false
                })
            },
            sortByPath (rowA, rowB) {
                return rowA._path_.toLowerCase().localeCompare(rowB._path_.toLowerCase(), 'zh-Hans-CN', { sensitivity: 'accent' })
            },
            sortByTime (rowA, rowB) {
                const timeA = (new Date(rowA.last_time)).getTime()
                const timeB = (new Date(rowB.last_time)).getTime()
                return timeA - timeB
            },
            handleToCreatedInstance () {
                this.$routerActions.redirect({ name: MENU_BUSINESS_HOST_AND_SERVICE })
            },
            handlePathClick (row) {
                this.$routerActions.open({
                    name: MENU_BUSINESS_HOST_AND_SERVICE,
                    query: {
                        node: `module-${row.bk_module_id}`
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .template-instance-layout {
        height: 100%;
        padding: 0 20px;
    }
    .topo-path {
        cursor: pointer;
        &:hover {
            color: $primaryColor;
        }
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
        .latest-sync {
            font-size: 12px;
            cursor: not-allowed;
            color: #DCDEE5;
        }
    }
</style>
