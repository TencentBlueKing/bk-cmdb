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
            :list="tableList"
            :defaultSort="table.defaultSort"
            @handleSortChange="handleSortChange">
            <template slot="info" slot-scope="{ item }">
                <div>
                    {{getConfigInfo(item)}}
                </div>
            </template>
            <template slot="last_time" slot-scope="{ item }">
                {{$tools.formatTime(item['last_time'], 'YYYY-MM-DD')}}
            </template>
            <template slot="operation" slot-scope="{ item }">
                <span class="text-primary" @click.stop="routeToConfirm(item)">{{$t('NetworkDiscovery["详情确认"]')}}</span>
            </template>
        </cmdb-table>
    </div>
</template>

<script>
    import { mapActions, mapMutations } from 'vuex'
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
                        name: this.$t('NetworkDiscovery["配置信息"]'),
                        sortable: false
                    }, {
                        id: 'last_time',
                        name: this.$t('NetworkDiscovery["发现时间"]')
                    }, {
                        id: 'operation',
                        name: this.$t('Association["操作"]'),
                        sortable: false
                    }],
                    list: [],
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
            tableList () {
                return this.table.list.filter(({bk_cloud_name: cloudName}) => cloudName.includes(this.filter.text))
            }
        },
        created () {
            this.getTableData()
        },
        methods: {
            ...mapMutations('netDiscovery', ['setCloudName']),
            ...mapActions('netDiscovery', [
                'searchNetcollect'
            ]),
            getConfigInfo (item) {
                if (item.statistics) {
                    let str = []
                    Object.keys(item.statistics).map(key => {
                        if (key !== 'associations') {
                            str.push(`${key}(${item.statistics[key]})`)
                        }
                    })
                    return `${str.join(' ')} ${this.$t("Hosts['关联关系']")}(${item.statistics.associations})`
                }
            },
            routeToConfig () {
                this.$store.commit('setHeaderStatus', {
                    back: true
                })
                this.$router.push({name: 'networkDiscoveryConfig'})
            },
            routeToConfirm (item) {
                this.$store.commit('setHeaderStatus', {
                    back: true
                })
                this.setCloudName(item['bk_cloud_name'])
                this.$router.push({
                    name: 'networkDiscoveryConfirm',
                    params: {
                        cloudId: item['bk_cloud_id']
                    }
                })
            },
            routeToHistory () {
                this.$store.commit('setHeaderStatus', {
                    back: true
                })
                this.$router.push({
                    name: 'networkDiscoveryHistory'
                })
            },
            async getTableData () {
                const res = await this.searchNetcollect({params: {}, config: {requestId: 'searchNetcollect'}})
                this.table.list = res
            },
            handleSortChange (sort) {
                let key = sort
                if (sort[0] === '-') {
                    key = sort.substr(1, sort.length - 1)
                }
                this.table.list.sort((itemA, itemB) => {
                    if (itemA[key] === null) {
                        itemA[key] = ''
                    }
                    if (itemB[key] === null) {
                        itemB[key] = ''
                    }
                    return itemA[key].localeCompare(itemB[key])
                })
                if (sort[0] === '-') {
                    this.table.list.reverse()
                }
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
