<template>
    <div class="details-layout">
        <div class="info">
            <div class="info-basic">
                <i :class="['info-icon', model.bk_obj_icon]"></i>
                <span class="info-ip">{{host.bk_host_innerip}}</span>
            </div>
            <div class="info-topology clearfix">
                <div class="topology-label fl">{{$t("BusinessTopology['业务拓扑']")}}：</div>
                <div class="topology-details clearfix">
                    <div class="topology-list fl"
                        v-for="(column, index) in ['left', 'right']"
                        ref="topologyList"
                        :key="index"
                        :class="`${column}-list`">
                        <div class="topology-item"
                            v-if="topologyList[column].length"
                            :title="topologyList[column][0]">
                            {{topologyList[column][0]}}
                        </div>
                        <cmdb-collapse-transition>
                            <div class="topology-transiton" v-show="showAll">
                                <div class="topology-item"
                                    v-for="(path, columnIndex) in topologyList[column].slice(1)"
                                    :key="columnIndex"
                                    :title="path">
                                    {{path}}
                                </div>
                            </div>
                        </cmdb-collapse-transition>
                    </div>
                    <a class="view-all fl"
                        href="javascript:void(0)"
                        v-if="topology.length > 2"
                        @click="viewAll">
                        {{$t("Common['更多']")}}
                        <i class="bk-icon icon-angle-down" :class="{ 'is-all-show': showAll }"></i>
                    </a>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                showAll: false
            }
        },
        computed: {
            id () {
                return parseInt(this.$route.params.id)
            },
            business () {
                const business = parseInt(this.$route.params.business)
                if (isNaN(business)) {
                    return -1
                }
                return business
            },
            info () {
                return this.$store.state.hostDetails.info || {}
            },
            host () {
                return this.info.host || {}
            },
            topology () {
                const topology = []
                const sets = this.info.set || []
                const modules = this.info.module || []
                const businesses = this.info.biz || []
                modules.forEach(module => {
                    const set = sets.find(set => set.bk_set_id === module.bk_set_id)
                    if (set) {
                        const business = businesses.find(business => business.bk_biz_id === set.bk_biz_id)
                        if (business) {
                            topology.push(`${business.bk_biz_name} / ${set.bk_set_name} / ${module.bk_module_name}`)
                        }
                    }
                })
                return topology
            },
            topologyList () {
                const list = {
                    left: [],
                    right: []
                }
                this.topology.forEach((item, index) => {
                    list[index % 2 ? 'right' : 'left'].push(item)
                })
                return list
            },
            model () {
                return this.$store.getters['objectModelClassify/getModelById']('host')
            }
        },
        watch: {
            info (info) {
                this.$store.commit('setHeaderTitle', `${this.$t('HostDetails["主机详情"]')}(${info.host.bk_host_innerip})`)
            }
        },
        created () {
            this.$store.commit('setHeaderTitle', this.$t('HostDetails["主机详情"]'))
            this.$store.commit('setHeaderStatus', { back: true })
        },
        async mounted () {
            await this.getHostInfo()
            this.calcTopologyListWidth()
        },
        methods: {
            async getHostInfo () {
                const { info } = await this.$store.dispatch('hostSearch/searchHost', {
                    params: this.getSearchHostParams()
                })
                if (info.length) {
                    this.$store.commit('hostDetails/setHostInfo', info[0])
                    return Promise.resolve()
                }
                this.$router.replace({ name: 404 })
            },
            getSearchHostParams () {
                const hostCondition = {
                    field: 'bk_host_id',
                    operator: '$eq',
                    value: this.id
                }
                return {
                    bk_biz_id: this.business,
                    condition: ['biz', 'set', 'module', 'host'].map(model => {
                        return {
                            bk_obj_id: model,
                            condition: model === 'host' ? [hostCondition] : [],
                            fields: []
                        }
                    }),
                    ip: { flag: 'bk_host_innerip', exact: 1, data: [] }
                }
            },
            viewAll () {
                this.showAll = !this.showAll
            },
            calcTopologyListWidth () {
                // const range = document.createRange();
                //   range.setStart(cellChild, 0);
                //   range.setEnd(cellChild, cellChild.childNodes.length);
                //   const rangeWidth = range.getBoundingClientRect().width;
            }
        }
    }
</script>

<style lang="scss" scoped>
    .details-layout {
        padding: 0;
        .info {
            padding: 11px 24px;
            background:rgba(235, 244, 255, .6);
        }
    }
    .info-basic {
        font-size: 0;
        .info-icon {
            display: inline-block;
            width: 38px;
            height: 38px;
            margin: 0 11px 0 0;
            border: 1px solid #dde4eb;
            border-radius: 50%;
            background-color: #fff;
            vertical-align: middle;
            line-height: 38px;
            text-align: center;
            font-size: 18px;
            color: #3a84ff;
        }
        .info-ip {
            display: inline-block;
            vertical-align: middle;
            line-height: 38px;
            font-size: 16px;
            color: #333948;
        }
    }
    .info-topology {
        line-height: 19px;
        .topology-label {
            padding: 0 0 0 50px;
            font-size: 14px;
            font-weight: bold;
        }
        .topology-details {
            overflow: hidden;
        }
        .view-all {
            margin: 0 0 0 14px;
            font-size: 12px;
            color: #007eff;
            .bk-icon {
                display: inline-block;
                vertical-align: -1px;
                font-size: 12px;
                transition: transform .2s linear;
                &.is-all-show {
                    transform: rotate(-180deg);
                }
            }
        }
    }
    .topology-list {
        max-width: 350px;
        color: #63656e;
        .topology-item {
            margin: 0 0 9px 0;
        }
    }
    .topology-list.right-list {
        margin: 0 0 0 90px;
    }
</style>
