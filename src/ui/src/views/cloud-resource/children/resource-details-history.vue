<template>
    <div class="history-layout">
        <div class="history-options">
            <bk-date-picker class="options-date"
                type="datetimerange"
                :placeholder="$t('请选择xx', { name: '时间范围' })"
                v-model="timeRange">
            </bk-date-picker>
        </div>
        <div class="history-table">
            <bk-table ref="table" :data="histories" :max-height="$APP.height - 180">
                <bk-table-column type="expand" width="30" align="center">
                    <cloud-resource-details-history-content
                        slot-scope="{ row }"
                        :list="row.list">
                    </cloud-resource-details-history-content>
                </bk-table-column>
                <bk-table-column :label="$t('操作概要')" prop="summary"></bk-table-column>
                <bk-table-column :label="$t('状态')" prop="xxxx">
                    <div class="row-status" slot-scope="{ row }">
                        <i :class="['status', { 'is-error': row.error }]"></i>
                        异常
                    </div>
                </bk-table-column>
                <bk-table-column :label="$t('时间')" prop="last_time"></bk-table-column>
                <bk-table-column :label="$t('详情')">
                    <template slot-scope="{ row }">
                        <link-button @click="handleView(row)">{{$t('查看详情')}}</link-button>
                    </template>
                </bk-table-column>
            </bk-table>
        </div>
    </div>
</template>

<script>
    import CloudResourceDetailsHistoryContent from './resource-details-history-content.vue'
    export default {
        name: 'cloud-resource-details-history',
        components: {
            [CloudResourceDetailsHistoryContent.name]: CloudResourceDetailsHistoryContent
        },
        data () {
            function repeate (count) {
                return Array(count).fill(Symbol('any')).map((_, index) => ({ isCreate: index % 2 === 0, bk_host_innerip: '192.168.1.1' }))
            }
            return {
                timeRange: [],
                histories: [{ list: repeate(200) }, { list: repeate(100) }]
            }
        },
        methods: {
            handleView (row) {
                this.$refs.table.toggleRowExpansion(row)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .history-layout {
        .history-options {
            padding: 12px 30px;
            .options-date {
                width: 320px;
            }
        }
        .history-table {
            padding: 0 30px;
            /deep/ {
                .bk-table-expanded-cell {
                    padding: 0;
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
            }
        }
    }
</style>
