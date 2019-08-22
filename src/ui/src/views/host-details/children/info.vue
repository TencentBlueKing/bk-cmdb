<template>
    <div class="info">
        <div class="info-basic">
            <i :class="['info-icon', model.bk_obj_icon]"></i>
            <span class="info-ip">{{hostIp}}</span>
        </div>
        <div class="info-topology clearfix">
            <div class="topology-label fl">{{$t('业务拓扑')}}：</div>
            <div class="topology-details clearfix">
                <ul class="topology-list fl"
                    v-for="(column, index) in ['left', 'right']"
                    ref="topologyList"
                    :key="index"
                    :class="`${column}-list`"
                    :style="{
                        height: getListHeight(topologyList[column]) + 'px'
                    }">
                    <li class="topology-item"
                        v-for="(path, columnIndex) in topologyList[column]"
                        :key="columnIndex"
                        :title="path">
                        {{path}}
                    </li>
                </ul>
                <a class="view-all fl"
                    href="javascript:void(0)"
                    v-if="topology.length > 2"
                    @click="viewAll">
                    {{$t('更多信息')}}
                    <i class="bk-icon icon-angle-down" :class="{ 'is-all-show': showAll }"></i>
                </a>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapState } from 'vuex'
    export default {
        name: 'cmdb-host-info',
        data () {
            return {
                topologyListWidth: {
                    left: 'auto',
                    right: 'auto'
                },
                showAll: false
            }
        },
        computed: {
            ...mapState('hostDetails', ['info']),
            host () {
                return this.info.host || {}
            },
            hostIp () {
                if (Object.keys(this.host).length) {
                    const hostList = this.host.bk_host_innerip.split(',')
                    const host = hostList.length > 1 ? `${hostList[0]}...` : hostList[0]
                    return host
                } else {
                    return ''
                }
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
        methods: {
            viewAll () {
                this.showAll = !this.showAll
                this.$emit('info-toggle', this.getListHeight(this.topologyList.left) + 51)
            },
            getListHeight (items) {
                const itemHeight = 21
                const itemMargin = 9
                return (this.showAll ? items.length : 1) * (itemHeight + itemMargin)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .info {
        padding: 11px 24px 2px;
        background:rgba(235, 244, 255, .6);
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
        overflow: hidden;
        max-width: 350px;
        color: #63656e;
        will-change: height;
        transition: height .2s ease-in;
        .topology-item {
            font-size: 14px;
            margin: 0 0 9px 0;
            @include ellipsis;
        }
    }
    .topology-list.right-list {
        margin: 0 0 0 90px;
    }
</style>
