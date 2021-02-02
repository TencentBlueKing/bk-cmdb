<template>
    <div class="cmdb-host-topo-path">
        <template v-if="pending">
            <i class="path-pending"></i>
        </template>
        <template v-else>
            <span class="path" v-bk-overflow-tips>
                {{getModulePath(modules[0])}}
            </span>
            <template v-if="!isResourcePool">
                <i class="path-single-link icon-cc-share"
                    v-if="isSingle"
                    v-bk-tooltips="$t('跳转业务拓扑')"
                    @click="handleLinkToTopology(modules[0])">
                </i>
                <span v-else
                    class="path-count"
                    v-bk-tooltips="{
                        content: $refs.tooltipContent,
                        interactive: true,
                        boundary: 'window'
                    }">
                    {{`+${modules.length - 1}`}}
                </span>
            </template>
        </template>
        <div v-if="!isSingle"
            class="path-tooltip-content"
            ref="tooltipContent">
            <div class="path-tooltip-item"
                v-for="moduleId in modules"
                :key="moduleId">
                <span class="path-tooltip-text" :title="getModulePath(moduleId)">{{getModulePath(moduleId)}}</span>
                <i class="path-tooltip-link icon-cc-share"
                    :title="$t('跳转业务拓扑')"
                    @click="handleLinkToTopology(moduleId)">
                </i>
            </div>
        </div>
    </div>
</template>

<script>
    import proxy from './proxy'
    import { MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'
    export default {
        name: 'cmdb-host-topo-path',
        props: {
            host: {
                type: Object,
                required: true
            }
        },
        data () {
            return {
                pending: true,
                nodes: []
            }
        },
        computed: {
            bizId () {
                const [biz] = this.host.biz
                return biz.bk_biz_id
            },
            isResourcePool () {
                const [biz] = this.host.biz
                return biz.default === 1
            },
            modules () {
                return this.host.module.map(module => module.bk_module_id)
            },
            isSingle () {
                return this.modules.length === 1
            }
        },
        watch: {
            host: {
                immediate: true,
                handler () {
                    this.searchPath()
                }
            }
        },
        methods: {
            async searchPath () {
                try {
                    this.pending = true
                    this.nodes = await proxy.search({
                        bk_biz_id: this.bizId,
                        modules: this.modules
                    })
                } catch (error) {
                    console.error(error)
                    this.nodes = []
                } finally {
                    this.pending = false
                    this.$emit('path-ready', this.getFullModulePath())
                }
            },
            getModulePath (moduleId) {
                const node = this.nodes.find(node => node.topo_node.bk_inst_id === moduleId)
                if (!node) {
                    return '--'
                }
                return node.topo_path.map(path => path.bk_inst_name).reverse().join(' / ')
            },
            getFullModulePath () {
                return this.modules.map(moduleId => this.getModulePath(moduleId))
            },
            handleLinkToTopology (moduleId) {
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_HOST_AND_SERVICE,
                    query: {
                        node: `module-${moduleId}`
                    },
                    params: {
                        bizId: this.bizId
                    },
                    history: true
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cmdb-host-topo-path {
        display: flex;
        align-items: center;
        &:hover {
            .path-single-link {
                display: inline-block;
            }
        }
        .path-pending {
            display: inline-block;
            vertical-align: middle;
            width: 16px;
            height: 16px;
            background-color: transparent;
            background-image: url("../../assets/images/icon/loading.svg");
        }
        .path {
            display: block;
            line-height: 24px;
            @include ellipsis;
        }
        .path-single-link {
            display: none;
            flex: 16px 0 0;
            font-size: 12px;
            margin-left: 5px;
            color: $primaryColor;
            cursor: pointer;
            &:hover {
                opacity: .75;
            }
        }
        .path-count {
            padding: 0 6px;
            margin: 0 0 0 5px;
            height: 16px;
            line-height: 16px;
            font-size: 12px;
            border-radius: 8px;
            white-space: nowrap;
            background-color: #dcdee5;
        }
        // 初始化时设置为隐藏，避免影响表格行高计算
        .path-tooltip-content {
            display: none;
        }
    }
    .path-tooltip-content {
        width: 300px;
        .path-tooltip-item {
            display: flex;
            align-items: center;
            &:hover {
                .path-tooltip-link {
                    display: inline-block;
                }
            }
            .path-tooltip-text {
                display: block;
                line-height: 24px;
                @include ellipsis;
            }
            .path-tooltip-link {
                display: none;
                flex: 16px 0 0;
                font-size: 12px;
                margin-left: 5px;
                color: $primaryColor;
                cursor: pointer;
                &:hover {
                    opacity: .75;
                }
            }
        }
    }
</style>
