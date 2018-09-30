<template>
    <div class="network-wrapper">
        <div class="filter-wrapper">
            <bk-button type="primary" @click="routeToConfig">
                {{$t('NetworkDiscovery["配置网络发现"]')}}
            </bk-button>
            <div class="filter-content fr">
                <div class="input-box fl">
                    <input type="text" class="cmdb-form-input" 
                    :placeholder="$t('NetworkDiscovery[\'请输入云区域名称\']')"
                    v-model.trim="filter.text"
                    @keyup.enter="getTableData">
                    <i class="filter-search bk-icon icon-search" @click="getTableData"></i>
                </div>
                <bk-button type="default" class="fl" v-tooltip="$t('NetworkDiscovery[\'查看完成历史\']')" @click="routeToHistory">
                    <i class="icon-cc-history"></i>
                </bk-button>
            </div>
        </div>
        <cmdb-table
            class="network-table"
            :loading="$loading('searchNetcollect')"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :defaultSort="table.defaultSort"
            @handleSortChange="handleSortChange"
            @handleSizeChange="handleSizeChange"
            @handlePageChange="handlePageChange">
            <template slot="info" slot-scope="{ item }">
                <div>
                    <span>{{$t("NetworkDiscovery['交换机']")}}({{item.info.switch}})</span>
                    <span>{{$t("Hosts['主机']")}}({{item.info.host}})</span>
                    <span>{{$t("Hosts['关联关系']")}}({{item.info.relation}})</span>
                </div>
            </template>
            <template slot="last_time" slot-scope="{ item }">
                {{$tools.formatTime(item['last_time'], 'YYYY-MM-DD')}}
            </template>
            <template slot="operation" slot-scope="{ item }">
                <span class="text-primary" @click.stop="routeToConfirm">{{$t('NetworkDiscovery["详情确认"]')}}</span>
            </template>
        </cmdb-table>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    export default {
        data () {
            return {
                filter: {
                    text: ''
                },
                table: {
                    header: [{
                        id: 'bk_cloud_name',
                        name: this.$t('Hosts["云区域"]')
                    }, {
                        id: 'info',
                        name: this.$t('NetworkDiscovery["配置信息"]')
                    }, {
                        id: 'last_time',
                        name: this.$t('NetworkDiscovery["发现时间"]')
                    }, {
                        id: 'operation',
                        name: this.$t('Association["操作"]'),
                        sortable: false
                    }],
                    list: [{
                        bk_cloud_name: 'asdf',
                        info: {
                            switch: 1,
                            host: 1,
                            relation: 2
                        },
                        last_time: '2018-04-17T15:00:49.274+08:00'
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
        computed: {
            params () {
                let pagination = this.table.pagination
                let params = {
                    page: {
                        start: (pagination.current - 1) * pagination.size,
                        limit: pagination.size,
                        sort: this.table.sort
                    }
                }
                if (this.filter.text.length) {
                    Object.assign(params, {condition: [{
                        field: 'bk_cloud_name',
                        operator: '$regex',
                        value: this.filter.text
                    }]})
                }
                return params
            }
        },
        methods: {
            ...mapActions('netDiscovery', [
                'searchNetcollect'
            ]),
            routeToConfig () {
                this.$store.commit('setHeaderStatus', {
                    back: true
                })
                this.$router.push('/network-discovery/config')
            },
            routeToConfirm () {
                this.$store.commit('setHeaderStatus', {
                    back: true
                })
                this.$router.push('/network-discovery/confirm')
            },
            routeToHistory () {
                this.$store.commit('setHeaderStatus', {
                    back: true
                })
                this.$router.push('/network-discovery/history')
            },
            async getTableData () {
                const res = await this.searchNetcollect({params: this.params, config: {requestId: 'searchNetcollect'}})
                this.table.pagination.count = res.count
                this.table.list = res.info
            },
            handleSortChange (sort) {
                this.table.sort = sort
                this.handlePageChange(1)
            },
            handleSizeChange (size) {
                this.table.pagination.size = size
                this.handlePageChange(1)
            },
            handlePageChange (page) {
                this.table.pagination.current = page
                this.getTableData()
            }
        }
    }
</script>


<style lang="scss" scoped>
    .network-wrapper {
        background: $cmdbBackgroundColor;
        .filter-wrapper {
            .filter-content {
                .input-box {
                    position: relative;
                    .icon-search {
                        position: absolute;
                        right: 10px;
                        top: 11px;
                        cursor: pointer;
                        font-size: 14px;
                        z-index: 3;
                    }
                }
                .cmdb-form-input {
                    width: 260px;
                }
                .bk-button {
                    margin-left: 10px;
                }
            }
        }
        .network-table {
            margin-top: 20px;
            background: #fff;
        }
    }
</style>
