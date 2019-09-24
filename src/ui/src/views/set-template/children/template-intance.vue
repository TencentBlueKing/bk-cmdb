<template>
    <div class="template-instance-layout">
        <div class="instance-empty" v-if="list.length">
            <img src="../../../assets/images/empty-content.png" alt="">
            <i18n path="空集群模板实例提示" tag="div">
                <bk-button text @click="handleLinkServiceTopo" place="link">{{$t('服务拓扑')}}</bk-button>
            </i18n>
        </div>
        <div class="instance-main">
            <div class="options clearfix">
                <div class="fl">
                    <bk-button theme="primary">{{$t('批量同步')}}</bk-button>
                    <span class="option-item" v-if="!updating">{{$t('您可以批量同步多个模版实例')}}</span>
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
            <bk-table v-bkloading="{ isLoading: $loading() }"
                :data="list"
                :pagination="pagination"
                :max-height="$APP.height - 229"
                @page-change="handlePageChange"
                @page-limit-change="handleSizeChange"
                @selection-change="handleSelectionChange">
                <bk-table-column type="selection" width="50" :selectable="handleSelectable"></bk-table-column>
                <bk-table-column :label="$t('拓扑结构')" prop=""></bk-table-column>
                <bk-table-column :label="$t('主机数量')" prop=""></bk-table-column>
                <bk-table-column :label="$t('状态')" prop=""></bk-table-column>
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
        methods: {
            handleLinkServiceTopo () {},
            handleViewSync () {},
            handleSearch () {},
            handleSelectionChange () {},
            handlePageChange () {},
            handleSizeChange () {},
            handleSelectable () {},
            handleSync () {}
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
    }
</style>
