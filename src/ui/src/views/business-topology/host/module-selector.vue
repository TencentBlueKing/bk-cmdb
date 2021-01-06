<template>
    <div class="module-selector-layout"
        v-bkloading="{ isLoading: $loading(Object.values(request)) }">
        <div class="wrapper">
            <cmdb-resize-layout class="topo-resize-layout"
                direction="right"
                :handler-offset="3"
                :min="281"
                :max="508">
                <div :class="['wrapper-column wrapper-left', { 'has-title': hasTitle }]">
                    <template v-if="hasTitle">
                        <h2 class="title">{{title}}</h2>
                    </template>
                    <bk-input class="tree-filter" clearable right-icon="icon-search" v-model="filter" :placeholder="$t('请输入关键词')"></bk-input>
                    <bk-big-tree ref="tree" class="topology-tree"
                        display-matched-node-descendants
                        :default-expand-all="moduleType === 'idle'"
                        :options="{
                            idKey: getNodeId,
                            nameKey: 'bk_inst_name',
                            childrenKey: 'child'
                        }"
                        :node-height="36"
                        :show-checkbox="isShowCheckbox"
                        @node-click="handleNodeClick"
                        @check-change="handleNodeCheck">
                        <template slot-scope="{ node, data }">
                            <i class="internal-node-icon fl"
                                v-if="data.default !== 0"
                                :class="getInternalNodeClass(node, data)">
                            </i>
                            <i v-else :class="['node-icon fl', { 'is-template': isTemplate(data) }]">{{data.bk_obj_name[0]}}</i>
                            <span :class="['node-checkbox fr', { 'is-checked': checked.includes(node) }]"
                                v-if="moduleType === 'idle' && data.bk_obj_id === 'module'">
                            </span>
                            <span class="node-name" :title="node.name">{{node.name}}</span>
                        </template>
                    </bk-big-tree>
                </div>
            </cmdb-resize-layout>
            <div class="wrapper-column wrapper-right">
                <module-checked-list :checked="checked" @delete="handleDeleteModule" @clear="handleClearModule" />
            </div>
        </div>
        <div class="layout-footer">
            <span class="footer-tips mr10" v-bk-tooltips="confirmTooltips">
                <bk-button theme="primary"
                    :disabled="!checked.length || !hasDifference"
                    :loading="confirmLoading"
                    @click="handleNextStep">
                    {{confirmText || $t('下一步')}}
                </bk-button>
            </span>
            <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import debounce from 'lodash.debounce'
    import ModuleCheckedList from './module-checked-list.vue'
    export default {
        name: 'cmdb-module-selector',
        components: {
            ModuleCheckedList
        },
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
            },
            confirmText: {
                type: String,
                default: ''
            },
            confirmLoading: {
                type: Boolean,
                default: false
            },
            previousModules: {
                type: Array,
                default: () => ([])
            },
            business: {
                type: Object,
                required: true
            }
        },
        data () {
            return {
                filter: '',
                handlerFilter: null,
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
            ...mapGetters('objectModelClassify', ['getModelById']),
            bizId () {
                return this.business.bk_biz_id
            },
            bizName () {
                return this.business.bk_biz_name
            },
            hasDifference () {
                const checkedModules = this.checked.map(node => node.data.bk_inst_id).sort()
                return checkedModules.join(',') !== this.previousModules.join(',')
            },
            hasTitle () {
                return this.title && this.title.length
            },
            confirmTooltips () {
                const tooltips = { disabled: true }
                if (!this.checked.length) {
                    tooltips.content = this.$t('请先选择业务模块')
                    tooltips.disabled = false
                } else if (!this.hasDifference) {
                    tooltips.content = this.$t('模块相同提示')
                    tooltips.disabled = false
                }
                return tooltips
            }
        },
        watch: {
            filter () {
                this.handlerFilter()
            }
        },
        async created () {
            this.handlerFilter = debounce(() => {
                this.$refs.tree.filter(this.filter)
            }, 300)
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
                            this.$refs.tree.setExpanded(node.id)
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
                        bk_inst_id: this.bizId,
                        bk_inst_name: this.bizName,
                        bk_obj_id: 'biz',
                        bk_obj_name: this.getModelById('biz').bk_obj_name,
                        default: 0,
                        child: [{
                            bk_inst_id: data.bk_set_id,
                            bk_inst_name: data.bk_set_name,
                            bk_obj_id: 'set',
                            bk_obj_name: this.getModelById('set').bk_obj_name,
                            default: 0,
                            child: this.$tools.sort((data.module || []), 'default').map(module => ({
                                bk_inst_id: module.bk_module_id,
                                bk_inst_name: module.bk_module_name,
                                bk_obj_id: 'module',
                                bk_obj_name: this.getModelById('module').bk_obj_name,
                                default: module.default
                            }))
                        }]
                    }]
                })
            },
            getBusinessModules () {
                return this.$store.dispatch('objectMainLineModule/getInstTopoInstanceNum', {
                    bizId: this.bizId,
                    config: {
                        requestId: this.request.business
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
            // 选择空闲模块
            handleNodeClick (node) {
                const data = node.data
                if (data.bk_obj_id !== 'module') {
                    return false
                }
                if (this.moduleType === 'idle') {
                    this.checked = [node]
                } else {
                    this.$refs.tree.setChecked(node.id, { checked: !node.checked, emitEvent: true })
                }
            },
            handleNodeCheck (checked, node) {
                if (this.moduleType === 'idle' || node.data.bk_obj_id !== 'module') {
                    return false
                }
                this.checked = checked.map(id => this.$refs.tree.getNodeById(id))
            },
            handleDeleteModule (node) {
                this.checked = this.checked.filter(exist => exist !== node)
                if (this.moduleType === 'business') {
                    this.$refs.tree.setChecked(node.id, { checked: false })
                }
            },
            handleClearModule () {
                if (this.moduleType === 'business') {
                    this.$refs.tree.setChecked(this.checked.map(node => node.id), { checked: false })
                }
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
            },
            isTemplate (data) {
                return data.service_template_id || data.set_template_id
            },
            isShowCheckbox (data) {
                if (this.moduleType === 'idle') {
                    return false
                }
                return data.bk_obj_id === 'module'
            }
        }
    }
</script>

<style lang="scss" scoped>
    .module-selector-layout {
        height: var(--height, 600px);
        min-height: 300px;
        padding: 0 0 50px;
        position: relative;
        .layout-footer {
            position: sticky;
            bottom: 0;
            left: 0;
            width: 100%;
            height: 50px;
            padding: 8px 20px;
            border-top: 1px solid $borderColor;
            background-color: #FAFBFD;
            text-align: right;
            font-size: 0;
            z-index: 100;
            .footer-tips {
                display: inline-block;
                vertical-align: middle;
            }
        }
    }
    .wrapper {
        display: flex;
        height: 100%;
        .topo-resize-layout {
            width: 408px;
            height: 100%;
            border-right: 1px solid #dcdee5;
        }
        .wrapper-column {
            flex: 1;
        }
        .wrapper-left {
            height: calc(100% - 24px);
            margin-top: 24px;
            .title {
                padding: 0 12px 24px 23px;
                font-size: 16px;
                font-weight: normal;
                color: #313238;
                line-height: 22px;
            }
            .tree-filter {
                display: block;
                width: auto;
                margin: 0 12px 0 23px;
            }

            &.has-title {
                height: calc(100% - 18px);
                margin-top: 18px;
                .topology-tree {
                    max-height: calc(100% - 102px);
                }
            }
        }
        .wrapper-right {
            padding: 12px 23px;
            background: #f5f6fa;
            overflow: hidden;
        }
    }
    .topology-tree {
        width: 100%;
        max-height: calc(100% - 32px - 24px);
        padding: 0 0 0 12px;
        margin: 12px 0;
        @include scrollbar;
        .node-icon {
            display: block;
            width: 20px;
            height: 20px;
            margin: 8px 4px 8px 0;
            border-radius: 50%;
            background-color: #C4C6CC;
            line-height: 1.666667;
            text-align: center;
            font-size: 12px;
            font-style: normal;
            color: #FFF;
            &.is-template {
                background-color: #97AED6;
            }
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
</style>
