<template>
    <div class="layout clearfix">
        <div class="wrapper-left fl">
            <h2 class="title">{{title}}</h2>
            <h3 class="subtitle">{{$t('请勾选需要转到的模块')}}</h3>
            <bk-input class="tree-filter" right-icon="icon-search" v-model="filter"></bk-input>
            <bk-big-tree ref="tree" class="topology-tree"
                :default-expand-all="moduleType === 'idle'"
                :options="{
                    idKey: getNodeId,
                    nameKey: 'bk_inst_name',
                    childrenKey: 'child'
                }"
                @node-click="handleNodeClick">
                <template slot-scope="{ node, data }">
                    <i class="internal-node-icon fl"
                        v-if="data.default !== 0"
                        :class="getInternalNodeClass(node, data)">
                    </i>
                    <i v-else class="node-icon fl">{{data.bk_obj_name[0]}}</i>
                    <span :class="['node-checkbox fr', { 'is-checked': checked.includes(node) }]"
                        v-if="data.bk_obj_id === 'module'">
                    </span>
                    <span class="node-name">{{node.name}}</span>
                </template>
            </bk-big-tree>
        </div>
        <div class="wrapper-right fl">
            <div class="selected-info clearfix">
                <i18n class="selected-count fl" path="已选择N个模块">
                    <span class="count" place="count">{{checked.length}}</span>
                </i18n>
                <bk-button class="fr" text theme="primary"
                    v-show="checked.length"
                    @click="handleClearModule">
                    {{$t('清空')}}
                </bk-button>
            </div>
            <ul class="module-list">
                <li class="module-item" v-for="node in checked"
                    :key="node.id">
                    <div class="module-info clearfix">
                        <span class="info-icon fl">{{node.data.bk_obj_name[0]}}</span>
                        <span class="info-name" :title="node.data.bk_inst_name">
                            {{node.data.bk_inst_name}}
                        </span>
                    </div>
                    <div class="module-topology" :title="getNodePath(node)">{{getNodePath(node)}}</div>
                    <i class="bk-icon icon-close" @click="handleDeleteModule(node)"></i>
                </li>
            </ul>
        </div>
        <i class="clearfix"></i>
        <div class="wrapper-footer">
            <bk-button class="mr10" theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
            <bk-button theme="primary" @click="handleNextStep">{{$t('下一步')}}</bk-button>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        name: 'cmdb-module-selector',
        props: {
            defaultChecked: {
                type: Array,
                default () {
                    return []
                }
            },
            title: {
                type: String,
                default: ''
            },
            moduleType: {
                type: String,
                validator (moduleType) {
                    return ['idle', 'business'].includes(moduleType)
                }
            }
        },
        data () {
            return {
                filter: '',
                checked: [],
                request: {
                    internal: Symbol('internal'),
                    business: Symbol('business')
                },
                nodeIconMap: {
                    1: 'icon-cc-host-free-pool',
                    2: 'icon-cc-host-breakdown',
                    default: 'icon-cc-host-free-pool'
                }
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId', 'currentBusiness']),
            ...mapGetters('businessHost', ['topologyModels']),
            internalModelMap () {
                const map = {}
                this.topologyModels.forEach(model => {
                    map[model.bk_obj_id] = model
                })
                return map
            }
        },
        async created () {
            this.getModules()
        },
        methods: {
            async getModules () {
                try {
                    let data
                    if (this.moduleType === 'idle') {
                        data = await this.getInternalModules()
                    } else {
                        data = await this.getBusinessModules()
                    }
                    this.$refs.tree.setData(data)
                    this.$refs.tree.setExpanded(this.getNodeId(data[0]))
                    this.setDefaultChecked()
                } catch (e) {
                    this.$refs.tree.setData([])
                    console.error(e)
                }
            },
            setDefaultChecked () {
                this.$nextTick(() => {
                    this.defaultChecked.forEach(id => {
                        const node = this.$refs.tree.getNodeById(this.getNodeId({
                            bk_obj_id: 'module',
                            bk_inst_id: id
                        }))
                        if (node) {
                            this.checked.push(node)
                            this.$refs.tree.setChecked(node.id)
                        }
                    })
                })
            },
            getInternalModules () {
                return this.$store.dispatch('objectMainLineModule/getInternalTopo', {
                    bizId: this.bizId,
                    config: {
                        requestId: this.request.internal
                    }
                }).then(data => {
                    return [{
                        bk_inst_id: this.currentBusiness.bk_biz_id,
                        bk_inst_name: this.currentBusiness.bk_biz_name,
                        bk_obj_id: 'biz',
                        bk_obj_name: this.internalModelMap.biz.bk_obj_name,
                        default: 0,
                        child: [{
                            bk_inst_id: data.bk_set_id,
                            bk_inst_name: data.bk_set_name,
                            bk_obj_id: 'set',
                            bk_obj_name: this.internalModelMap.set.bk_obj_name,
                            default: 0,
                            child: (data.module || []).map(module => ({
                                bk_inst_id: module.bk_module_id,
                                bk_inst_name: module.bk_module_name,
                                bk_obj_id: 'module',
                                bk_obj_name: this.internalModelMap.module.bk_obj_name,
                                default: module.default
                            }))
                        }]
                    }]
                })
            },
            getBusinessModules () {
                return this.$store.dispatch('objectMainLineModule/getInstTopo', {
                    bizId: this.bizId,
                    config: {
                        requestId: this.request.instance
                    }
                })
            },
            getNodeId (data) {
                return `${data.bk_obj_id}-${data.bk_inst_id}`
            },
            getInternalNodeClass (node, data) {
                const clazz = []
                clazz.push(this.nodeIconMap[data.default] || this.nodeIconMap.default)
                return clazz
            },
            getNodePath (node) {
                const parents = node.parents
                return parents.map(parent => parent.data.bk_inst_name).join(' / ')
            },
            handleNodeClick (node) {
                const data = node.data
                if (data.bk_obj_id !== 'module') {
                    return false
                }
                if (this.checked.includes(node)) {
                    this.checked = this.checked.filter(exist => exist !== node)
                } else {
                    if (this.moduleType === 'idle') {
                        this.checked = [node]
                    } else {
                        this.checked.push(node)
                    }
                }
            },
            handleDeleteModule (node) {
                this.checked = this.checked.filter(exist => exist !== node)
            },
            handleClearModule () {
                this.checked = []
            },
            handleCancel () {
                this.$emit('cancel')
            },
            handleNextStep () {
                if (!this.checked.length) {
                    this.$warn('请选择模块')
                    return false
                }
                this.$emit('confirm', this.checked)
            }
        }
    }
</script>

<style lang="scss" scoped>
    $leftPadding: 0 12px 0 23px;
    .layout {
        height: 542px;
    }
    .wrapper-left {
        width: 350px;
        height: 490px;
        border-right: 1px solid $borderColor;
        .title {
            margin-top: 15px;
            padding: $leftPadding;
            font-size: 24px;
            font-weight: normal;
            color: #444444;
            line-height:32px;
        }
        .subtitle {
            margin-top: 10px;
            padding: $leftPadding;
            font-size: 14px;
            font-weight: normal;
            color: $textColor;
            line-height: 20px;
        }
        .tree-filter {
            display: block;
            width: auto;
            margin: 10px 12px 0 23px;
        }
    }
    .topology-tree {
        width: 100%;
        max-height: 370px;
        padding: 0 0 0 23px;
        @include scrollbar;
        .node-icon {
            display: block;
            width: 20px;
            height: 20px;
            margin: 8px 4px 8px 0;
            vertical-align: middle;
            border-radius: 50%;
            background-color: #C4C6CC;
            line-height: 1.666667;
            text-align: center;
            font-size: 12px;
            font-style: normal;
            color: #FFF;
            &.is-selected {
                background-color: #3A84FF;
            }
        }
        .node-name {
            height: 36px;
            line-height: 36px;
            overflow: hidden;
            @include ellipsis;
        }
        .node-checkbox {
            width: 16px;
            height: 16px;
            margin: 10px 17px 0 10px;
            background: #FFF;
            border-radius: 50%;
            border: 1px solid #979BA5;
            &.is-checked {
                padding: 3px;
                border-color: $primaryColor;
                background-color: $primaryColor;
                background-clip: content-box;
            }
        }
        .internal-node-icon{
            width: 20px;
            height: 20px;
            line-height: 20px;
            text-align: center;
            margin: 8px 4px 8px 0;
        }
    }

    .wrapper-right {
        width: 370px;
        padding: 57px 23px 0;
        .selected-info {
            font-size: 14px;
            line-height: 20px;
            color: $textColor;
        }
    }
    .module-list {
        max-height: 400px;
        margin-top: 12px;
        @include scrollbar-y;
        .module-item {
            position: relative;
            margin-top: 12px;
            .icon-close {
                position: absolute;
                top: 6px;
                right: 0px;
                width: 28px;
                height: 28px;
                line-height: 28px;
                text-align: center;
                color: #D8D8D8;
                cursor: pointer;
                &:hover {
                    color: #979BA5;
                }
            }
        }
    }
    .module-info {
        .info-icon {
            width: 20px;
            height: 20px;
            border-radius: 50%;
            background-color: #C4C6CC;
            line-height: 1.666667;
            text-align: center;
            font-size: 12px;
            font-style: normal;
            color: #FFF;
        }
        .info-name {
            display: block;
            width: 250px;
            padding-left: 10px;
            font-size:14px;
            color: $textColor;
            line-height:20px;
            @include ellipsis;
        }
    }
    .module-topology {
        width: 250px;
        padding-left: 30px;
        margin-top: 3px;
        font-size: 12px;
        color: #C4C6CC;
        @include ellipsis;
    }
    .wrapper-footer {
        height: 51px;
        padding: 9px 20px;
        border-top: 1px solid $borderColor;
        background-color: #FAFBFD;
        text-align: right;
        font-size: 0;
    }
</style>
