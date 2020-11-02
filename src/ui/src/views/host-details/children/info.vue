<template>
    <div class="info">
        <div class="info-basic">
            <i :class="['info-icon', model.bk_obj_icon]"></i>
            <span class="info-ip">{{hostIp}}</span>
            <span class="info-area">{{cloudArea}}</span>
        </div>
        <div class="info-topology">
            <div class="topology-label">
                <span>{{$t('所属拓扑')}}</span>
                <i class="topology-toggle icon-cc-single-column" v-if="isSingleColumn" @click="toggleDisplayType"></i>
                <i class="topology-toggle icon-cc-double-column" v-else @click="toggleDisplayType"></i>
                <span v-pre style="padding: 0 5px;">:</span>
            </div>
            <ul class="topology-list"
                :class="{ 'is-single-column': isSingleColumn }"
                :style="{
                    height: getListHeight(topologyList) + 'px'
                }">
                <li class="topology-item"
                    v-for="(item, index) in topologyList"
                    :key="index">
                    <span class="topology-path" v-bk-overflow-tips @click="handlePathClick(item)">{{item.path}}</span>
                    <i class="topology-remove-trigger icon-cc-tips-close"
                        v-if="!item.isInternal"
                        v-bk-tooltips="{ content: $t('从该模块移除'), interactive: false }"
                        @click="handleRemove(item.id)">
                    </i>
                </li>
            </ul>
            <a class="view-all"
                href="javascript:void(0)"
                v-if="showMore"
                @click="viewAll">
                {{$t('更多信息')}}
                <i class="bk-icon icon-angle-down" :class="{ 'is-all-show': showAll }"></i>
            </a>
        </div>
    </div>
</template>

<script>
    import { MENU_BUSINESS_TRANSFER_HOST, MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'
    import { mapState } from 'vuex'
    export default {
        name: 'cmdb-host-info',
        data () {
            return {
                displayType: window.localStorage.getItem('host_topology_display_type') || 'double',
                showAll: false,
                topoNodesPath: []
            }
        },
        computed: {
            ...mapState('hostDetails', ['info']),
            isSingleColumn () {
                return this.displayType === 'single'
            },
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
            cloudArea () {
                return (this.host.bk_cloud_id || []).map(cloud => {
                    return `${this.$t('云区域')}：${cloud.bk_inst_name} (ID：${cloud.bk_inst_id})`
                }).join('\n')
            },
            topologyList () {
                const modules = this.info.module || []
                return this.topoNodesPath.map(item => {
                    const instId = item.topo_node.bk_inst_id
                    const module = modules.find(module => module.bk_module_id === instId)
                    return {
                        id: instId,
                        path: item.topo_path.reverse().map(node => node.bk_inst_name).join(' / '),
                        isInternal: module && module.default !== 0
                    }
                }).sort((itemA, itemB) => {
                    return itemA.path.localeCompare(itemB.path, 'zh-Hans-CN', { sensitivity: 'accent' })
                })
            },
            showMore () {
                if (this.isSingleColumn) {
                    return this.topologyList.length > 1
                }
                return this.topologyList.length > 2
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
                    const biz = this.info.biz || []
                    const result = await this.$store.dispatch('objectMainLineModule/getTopoPath', {
                        bizId: biz[0].bk_biz_id,
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
                this.$emit('info-toggle')
            },
            getListHeight (items) {
                const itemHeight = 21
                const itemMargin = 9
                const length = this.isSingleColumn ? items.length : (items.length / 2)
                return (this.showAll ? length : 1) * (itemHeight + itemMargin)
            },
            toggleDisplayType () {
                this.displayType = this.displayType === 'single' ? 'double' : 'single'
                this.$emit('info-toggle')
                window.localStorage.setItem('host_topology_display_type', this.displayType)
            },
            handlePathClick (item) {
                this.$routerActions.open({
                    name: MENU_BUSINESS_HOST_AND_SERVICE,
                    query: {
                        node: `module-${item.id}`
                    }
                })
            },
            handleRemove (moduleId) {
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_TRANSFER_HOST,
                    params: {
                        type: 'remove',
                        module: moduleId
                    },
                    query: {
                        sourceModel: 'module',
                        sourceId: moduleId,
                        resources: this.$route.params.id
                    },
                    history: true
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .info {
        max-height: 450px;
        padding: 11px 0 2px 24px;
        background:rgba(235, 244, 255, .6);
        border-bottom: 1px solid #dcdee5;
        @include scrollbar-y;
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
            font-weight: bold;
            color: #333948;
        }
        .info-area {
            display: inline-block;
            vertical-align: middle;
            height: 18px;
            margin-left: 10px;
            padding: 0 5px;
            line-height: 16px;
            font-size: 12px;
            color: #979BA5;
            border: 1px solid #C4C6CC;
            border-radius: 2px;
            background-color: #fff;
        }
    }
    .info-topology {
        line-height: 19px;
        display: flex;
        .topology-label {
            display: flex;
            align-items: center;
            align-self: baseline;
            padding: 0 0 0 50px;
            font-size: 14px;
            font-weight: bold;
            line-height: 20px;
            .topology-toggle {
                font-size: 16px;
                margin: 0 0 0 5px;
                cursor: pointer;
                &:hover {
                    opacity: .75;
                }
            }
        }
        .topology-list {
            flex: 1;
        }
        .view-all {
            margin: 0 14px;
            font-size: 12px;
            color: #007eff;
            .bk-icon {
                display: inline-block;
                vertical-align: -1px;
                font-size: 20px;
                margin-left: -4px;
                transition: transform .2s linear;
                &.is-all-show {
                    transform: rotate(-180deg);
                }
            }
        }
    }
    .topology-list {
        display: flex;
        flex-wrap: wrap;
        overflow: hidden;
        color: #63656e;
        will-change: height;
        transition: height .2s ease-in;
        max-width: 700px;
        &.is-single-column {
            max-width: 850px;
            display: inline-block;
            flex: none;
            .topology-item {
                width: auto;
            }
        }
        .topology-item {
            width: 50%;
            height: 20px;
            font-size: 0px;
            margin: 0 0 9px 0;
            padding: 0 15px 0 0;
            line-height: 20px;
            &:before {
                content: "";
                height: 100%;
                width: 0;
                display: inline-block;
                vertical-align: middle;
            }
            &:hover {
                .topology-remove-trigger {
                    opacity: 1;
                }
            }
            .topology-path {
                display: inline-block;
                vertical-align: middle;
                font-size: 14px;
                max-width: calc(100% - 30px);
                cursor: pointer;
                @include ellipsis;
                &:hover {
                    color: $primaryColor;
                }
            }
            .topology-remove-trigger {
                opacity: 0;
                font-size: 20px;
                cursor: pointer;
                margin: 0 0 0 10px;
                color: $textColor;
                transform: scale(.5);
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
