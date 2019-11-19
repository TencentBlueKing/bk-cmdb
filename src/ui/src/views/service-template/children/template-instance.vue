<template>
    <div class="template-instance-layout">
        <div class="instance-main">
            <div class="options clearfix">
                <div class="fr">
                    <bk-input class="filter-item" right-icon="bk-icon icon-search"
                        clearable
                        :placeholder="$t('请输入模块名称搜索')">
                    </bk-input>
                </div>
            </div>
            <bk-table class="instance-table"
                ref="instanceTable"
                v-bkloading="{ isLoading: $loading(Object.values(request)) }"
                :data="table.data"
                :pagination="table.pagination"
                @page-change="handlePageChange"
                @page-limit-change="handlePageLimitChange">
                <bk-table-column :label="$t('模块名称')" prop="bk_module_name"></bk-table-column>
                <bk-table-column :label="$t('拓扑路径')">
                    <div slot-scope="{ row }" :title="getRowTopoPath(row)">{{getRowTopoPath(row)}}</div>
                </bk-table-column>
                <bk-table-column :label="$t('主机数量')" prop="host_count"></bk-table-column>
                <bk-table-column :label="$t('操作')">
                    <template slot-scope="{ row }">
                        <bk-button text :disabled="isSyncEnable(row)" @click="handleSync(row)">{{$t('去同步')}}</bk-button>
                    </template>
                </bk-table-column>
            </bk-table>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        props: {
            active: Boolean
        },
        data () {
            return {
                table: {
                    filter: '',
                    data: [],
                    syncStatus: [],
                    topoPath: [],
                    pagination: this.$tools.getDefaultPaginationConfig()
                },
                request: {
                    instance: Symbol('instance'),
                    status: Symbol('status')
                }
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
        methods: {
            async refresh () {
                try {
                    const data = await this.getTemplateInstance()
                    const [syncStatus, topoPath] = await Promise.all([
                        this.getSyncStatus(data.info),
                        this.getTopoPath(data.info)
                    ])
                    this.table.data = data.info
                    this.table.pagination.count = data.count
                    this.table.syncStatus = syncStatus
                    this.table.topoPath = topoPath.nodes
                } catch (e) {
                    console.error(e)
                }
            },
            getTemplateInstance () {
                return this.$store.dispatch('serviceTemplate/getServiceTemplateModules', {
                    bizId: this.bizId,
                    serviceTemplateId: this.serviceTemplateId,
                    params: {
                        ...this.$tools.getPageParams(this.table.pagination)
                    },
                    config: {
                        requestId: this.request.instance
                    }
                })
            },
            getSyncStatus (modules) {
                return Promise.resolve([1])
            },
            getTopoPath (modules) {
                return this.$store.dispatch('objectMainLineModule/getTopoPath', {
                    bizId: this.bizId,
                    params: {
                        topo_nodes: modules.map(module => ({ bk_obj_id: 'module', bk_inst_id: module.bk_module_id }))
                    }
                })
            },
            isSyncEnable (row) {
                return this.table.syncStatus.includes(row.id)
            },
            getRowTopoPath (row) {
                const topo = this.table.topoPath.find(topo => topo.topo_node.bk_inst_id === row.bk_module_id)
                if (topo) {
                    return topo.topo_path.map(path => path.bk_inst_name).reverse().join(' / ')
                }
                return '--'
            },
            handleSync (row) {
                console.log(row)
            },
            handlePageChange (current) {
                this.table.pagination.current = current
                this.refresh()
            },
            handlePageLimitChange (limit) {
                this.table.pagination.current = 1
                this.table.pagination.limit = limit
                this.refresh()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .template-instance-layout {
        height: 100%;
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
    }
</style>
