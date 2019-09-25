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
                    <bk-button theme="primary" @click="handleBatchSync">{{$t('批量同步')}}</bk-button>
                    <span class="option-item" v-if="!updating">{{$t('批量同步集群模板实例提示')}}</span>
                    <div class="option-item" v-else>
                        <i class="bk-icon icon-cc-updating"></i>
                        <i18n path="同步任务进行中" tag="span">
                            <bk-button text @click="handleViewSync" place="link">{{$t('立即查看')}}</bk-button>
                        </i18n>
                    </div>
                </div>
                <bk-input class="filter-item" right-icon="bk-icon icon-search"
                    :placeholder="$t('集群名称')"
                    v-model="searchSetName"
                    @enter="handleSearch">
                </bk-input>
                <bk-select class="filter-item mr10" v-model="filterKey">
                    <bk-option v-for="option in filterList"
                        :key="option.id"
                        :id="option.id"
                        :name="option.name">
                    </bk-option>
                </bk-select>
            </div>
            <bk-table class="instance-table"
                v-bkloading="{ isLoading: $loading() }"
                :data="list"
                :pagination="pagination"
                :max-height="$APP.height - 229"
                @page-change="handlePageChange"
                @page-limit-change="handleSizeChange"
                @selection-change="handleSelectionChange">
                <bk-table-column type="selection" width="50" :selectable="handleSelectable"></bk-table-column>
                <bk-table-column :label="$t('拓扑结构')" prop="topo_path">
                    <template slot-scope="{ row }">
                        <span>{{getTopoPath(row)}}</span>
                    </template>
                </bk-table-column>
                <bk-table-column :label="$t('主机数量')" prop=""></bk-table-column>
                <bk-table-column :label="$t('状态')" prop="status">
                    <!-- <template slot-scope="{ row }">
                        <span class="sync-status fail">
                            <i class="bk-icon icon-cc-log-02 "></i>
                            {{$t('同步失败')}}
                        </span>
                        <span class="sync-status success">
                            <i class="bk-icon icon-check-1"></i>
                            {{$t('已完成')}}
                        </span>
                        <span class="sync-status">
                            <img class="svg-icon" src="../../../assets/images/icon/loading.svg" alt="">
                            {{$t('同步中')}}
                        </span>
                    </template> -->
                </bk-table-column>
                <bk-table-column :label="$t('同步人')" prop=""></bk-table-column>
                <bk-table-column :label="$t('操作')" width="180">
                    <template slot-scope="{ row }">
                        <bk-button text @click="handleSync(row)">{{$t('同步')}}</bk-button>
                        <bk-button text @click="handleSync(row)">{{$t('重试')}}</bk-button>
                    </template>
                </bk-table-column>
            </bk-table>
        </div>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                list: [],
                filterList: [],
                checkedList: [],
                filterKey: '',
                searchSetName: '',
                updating: false,
                pagination: {
                    count: 0,
                    current: 1,
                    ...this.$tools.getDefaultPaginationConfig()
                }
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
            }
        },
        async created () {
            await this.getSetTemplateInstances()
        },
        methods: {
            async getSetTemplateInstances () {
                const data = await this.$store.dispatch('setTemplate/getSetTemplateInstances', {
                    bizId: this.business,
                    setTemplateId: 1,
                    params: {
                        limit: this.limit
                    },
                    config: {
                        requestId: 'getSetTemplateInstances'
                    }
                })
                this.pagination.count = data.count
                this.list = data.info || []
            },
            getTopoPath (row) {
                const topoPath = this.$tools.clone(row.topo_path)
                if (topoPath.length) {
                    return topoPath.reverse().map(path => path.InstanceName).join(' / ')
                }
                return '--'
            },
            handleLinkServiceTopo () {},
            handleViewSync () {},
            handleSearch () {},
            handleSelectionChange (selection) {
                this.checkedList = selection.map(item => item.id)
            },
            handlePageChange () {},
            handleSizeChange () {},
            handleSelectable () {},
            handleBatchSync () {
                this.$store.commit('setFeatures/setSyncIdMap', {
                    id: `${this.business}_${1}`,
                    instancesId: this.checkedList
                })
                this.$router.push({
                    name: 'setSync',
                    params: {
                        setTemplateId: 1
                    }
                })
            },
            handleSync () {

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
            .option-item {
                @include inlineBlock;
                margin-left: 10px;
            }
            .filter-item {
                width: 230px;
                float: right;
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
