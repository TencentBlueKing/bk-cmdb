<template>
    <div class="network-wrapper">
        <div class="filter-wrapper">
            <bk-button type="primary">
                {{$t('NetworkDiscovery["配置网络发现"]')}}
            </bk-button>
            <div class="filter-content fr">
                <input type="text" class="cmdb-form-input fl" :placeholder="$t('NetworkDiscovery[\'请输入云区域名称\']')">
                <bk-button type="default" class="fl" v-tooltip="$t('NetworkDiscovery[\'查看完成历史\']')">
                    <i class="icon-cc-history"></i>
                </bk-button>
            </div>
        </div>
        <cmdb-table
            class="network-table"
            :loading="$loading('searchSubscription')"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :defaultSort="table.defaultSort">
            <template slot="operation" slot-scope="{ item }">
                <span class="text-primary" @click.stop="detailConfirm">{{$t('NetworkDiscovery["详情确认"]')}}</span>
            </template>
        </cmdb-table>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                table: {
                    header: [{
                        id: 'plat',
                        name: this.$t('Hosts["云区域"]')
                    }, {
                        id: 'config',
                        name: this.$t('NetworkDiscovery["配置信息"]')
                    }, {
                        id: 'time',
                        name: this.$t('NetworkDiscovery["发现时间"]')
                    }, {
                        id: 'operation',
                        name: this.$t('Association["操作"]'),
                        sortable: false
                    }],
                    list: [{
                        plat: 'asdf',
                        config: 'aaa',
                        time: 'ddd'
                    }],
                    pagination: {
                        count: 0,
                        size: 10,
                        current: 1
                    },
                    defaultSort: '-last_time',
                    sort: '-last_time'
                }
            }
        },
        methods: {
            detailConfirm () {
                this.$router.push('/network-discovery/confirm')
            }
        }
    }
</script>


<style lang="scss" scoped>
    .network-wrapper {
        background: #f5f6fa;
        .filter-wrapper {
            .filter-content {
                .cmdb-form-input {
                    width: 260px;
                    margin-right: 10px;
                }
            }
        }
        .network-table {
            margin-top: 20px;
            background: #fff;
        }
    }
</style>
