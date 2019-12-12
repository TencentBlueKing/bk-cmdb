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
                        v-for="(item, columnIndex) in topologyList[column]"
                        :key="columnIndex"
                        :title="item.path">
                        <span class="topology-path">{{item.path}}</span>
                        <i class="topology-remove-trigger icon-cc-tips-close"
                            v-if="!item.isInternal"
                            @click="handleRemove(item.id)">
                        </i>
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
    import { MENU_BUSINESS_TRANSFER_HOST } from '@/dictionary/menu-symbol'
    import { mapState } from 'vuex'
    export default {
        name: 'cmdb-host-info',
        data () {
            return {
                topologyListWidth: {
                    left: 'auto',
                    right: 'auto'
                },
                showAll: false,
                topoNodesPath: []
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
                const modules = this.info.module || []
                return this.topoNodesPath.map(item => {
                    const instId = item.topo_node.bk_inst_id
                    const module = modules.find(module => module.bk_module_id === instId)
                    return {
                        id: instId,
                        path: item.topo_path.reverse().map(node => node.bk_inst_name).join(' / '),
                        inInternal: module && module.default !== 0
                    }
                })
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
            async info () {
                await this.getModulePathInfo()
            }
        },
        methods: {
            async getModulePathInfo () {
                try {
                    const modules = this.info.module || []
                    const result = await this.$store.dispatch('objectMainLineModule/getTopoPath', {
                        bizId: this.$store.getters['objectBiz/bizId'],
                        params: {
                            topo_nodes: modules.map(module => ({ bk_obj_id: 'module', bk_inst_id: module.bk_module_id }))
                        }
                    })
                    this.topoNodesPath = result.nodes || []
                } catch (e) {
                    console.error(e)
                    this.topoNodesPath = []
                }
            },
            viewAll () {
                this.showAll = !this.showAll
                this.$emit('info-toggle', this.getListHeight(this.topologyList.left) + 51)
            },
            getListHeight (items) {
                const itemHeight = 21
                const itemMargin = 9
                return (this.showAll ? items.length : 1) * (itemHeight + itemMargin)
            },
            handleRemove (moduleId) {
                this.$router.push({
                    name: MENU_BUSINESS_TRANSFER_HOST,
                    params: {
                        type: 'remove'
                    },
                    query: {
                        sourceModel: 'module',
                        sourceId: moduleId,
                        resources: this.$route.params.id,
                        from: this.$route.fullPath
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .info {
        padding: 11px 24px 2px;
        background:rgba(235, 244, 255, .6);
        border-top: 1px solid #dcdee5;
        border-bottom: 1px solid #dcdee5;
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
        color: #63656e;
        will-change: height;
        transition: height .2s ease-in;
        .topology-item {
            font-size: 0px;
            margin: 0 0 9px 0;
            height: 20px;
            line-height: 20px;
            &:before {
                content: "";
                height: 100%;
                width: 0;
                display: inline-block;
                vertical-align: middle;
            }
            &:hover {
                color: #000000;
                .topology-remove-trigger {
                    opacity: 1;
                }
            }
            .topology-path {
                display: inline-block;
                vertical-align: middle;
                font-size: 14px;
                max-width: 330px;
                @include ellipsis;
            }
            .topology-remove-trigger {
                opacity: 0;
                font-size: 12px;
                cursor: pointer;
                margin: 0 0 0 10px;
                color: $textColor;
                &:hover {
                    color: $primaryColor;
                }
            }
        }
    }
    .topology-list.right-list {
        margin: 0 0 0 90px;
    }
</style>
