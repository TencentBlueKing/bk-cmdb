<template>
    <div class="sync-history-layout">
        <div class="options clearfix">
            <bk-date-picker style="width: 300px;" class="fl" type="daterange"></bk-date-picker>
            <bk-input style="width: 240px;" class="fl ml10"
                right-icon="icon-search"
                v-model="pathKeywords"
                :placeholder="$t('拓扑路径关键字')">
            </bk-input>
        </div>
        <bk-table v-bkloading="{ isLoading: $loading() }"
            :data="list"
            :pagination="pagination"
            :max-height="$APP.height - 229"
            @sort-change="handleSortChange"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange">
            <bk-table-column :label="$t('拓扑结构')" prop="topo_path">
                <template slot-scope="{ row }">
                    <span>{{getTopoPath(row)}}</span>
                </template>
            </bk-table-column>
            <bk-table-column :label="$t('状态')" prop="status">
                <template slot-scope="{ row }">
                    <span v-if="!row.details">--</span>
                    <span v-else-if="statusMap.syncing.includes(row.details.status)" class="sync-status">
                        <img class="svg-icon" src="../../assets/images/icon/loading.svg" alt="">
                        {{$t('同步中')}}
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
            <bk-table-column :label="$t('同步时间')" prop="sync_time" sortable="custom">
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
        </bk-table>
    </div>
</template>

<script>
    import { MENU_BUSINESS_SET_TEMPLATE } from '@/dictionary/menu-symbol'
    export default {
        data () {
            return {
                templateName: '',
                pathKeywords: '',
                statusMap: {
                    syncing: [0, 1, 100],
                    success: [200],
                    fail: [500]
                },
                list: [],
                pagination: {
                    count: 0,
                    current: 1,
                    ...this.$tools.getDefaultPaginationConfig()
                },
                listSort: ''
            }
        },
        computed: {
            templateId () {
                return this.$route.params.templateId
            }
        },
        created () {
            this.getSetTemplateInfo()
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
            getHistoryList () {

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
