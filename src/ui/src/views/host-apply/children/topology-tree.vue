<template>
    <div class="topology-tree-wrapper">
        <bk-big-tree class="topology-tree"
            ref="tree"
            v-bind="treeOptions"
            v-bkloading="{
                isLoading: $loading(['getTopologyData', 'getMainLine'])
            }"
            :options="{
                idKey: idGenerator,
                nameKey: 'bk_inst_name',
                childrenKey: 'child'
            }"
            selectable
            :before-select="beforeSelect"
            @select-change="handleSelectChange"
            @check-change="handleCheckChange"
        >
            <div class="node-info clearfix" :class="{ 'is-selected': node.selected }" slot-scope="{ node, data }">
                <i class="node-model-icon fl"
                    :class="{
                        'is-selected': node.selected,
                        'is-template': isTemplate(node),
                        'is-leaf-icon': node.isLeaf
                    }">
                    {{modelIconMap[data.bk_obj_id]}}
                </i>
                <span v-show="applyEnabled(node)" class="config-icon fr"><i class="bk-cc-icon icon-cc-selected"></i></span>
                <div class="info-content">
                    <span class="node-name">{{data.bk_inst_name}}</span>
                </div>
            </div>
        </bk-big-tree>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import ConfirmStore from '@/components/ui/dialog/confirm-store.js'

    const LEAVE_CONFIRM_ID = 'singleModule'

    export default {
        props: {
            treeOptions: {
                type: Object,
                default: () => {}
            }
        },
        data () {
            return {
                treeData: [],
                mainLine: [],
                createInfo: {
                    show: false,
                    visible: false,
                    properties: [],
                    parentNode: null,
                    nextModelId: null
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount', 'isAdminView']),
            business () {
                return this.$store.state.objectBiz.bizId
            },
            propertyMap () {
                return this.$store.state.businessTopology.propertyMap
            },
            mainLineModels () {
                const models = this.$store.getters['objectModelClassify/models']
                return this.mainLine.map(data => models.find(model => model.bk_obj_id === data.bk_obj_id))
            },
            modelIconMap () {
                const map = {}
                this.mainLineModels.forEach(model => {
                    map[model.bk_obj_id] = model.bk_obj_name[0]
                })
                return map
            },
            firstModule () {
                const findModule = function (data) {
                    for (const item of data) {
                        if (item.bk_obj_id === 'module') {
                            return item
                        } else if (item.child) {
                            return findModule(item.child)
                        }
                    }
                }
                return findModule(this.treeData)
            }
        },
        async created () {
            const [data, mainLine] = await Promise.all([
                this.getTopologyData(),
                this.getMainLine()
            ])
            this.treeData = data
            this.mainLine = mainLine
            this.$nextTick(() => {
                this.setDefaultState(data)
            })
        },
        methods: {
            setDefaultState (data) {
                this.$refs.tree.setData(data)
                if (this.firstModule) {
                    const defaultNodeId = this.idGenerator(this.firstModule)
                    this.$refs.tree.setSelected(defaultNodeId, { emitEvent: true })
                    this.$refs.tree.setExpanded(defaultNodeId)
                }
            },
            getTopologyData () {
                return this.$store.dispatch('objectMainLineModule/getInstTopoInstanceNum', {
                    bizId: this.business,
                    config: {
                        requestId: 'getTopologyWithStatData'
                    }
                })
            },
            getMainLine () {
                return this.$store.dispatch('objectMainLineModule/searchMainlineObject', {
                    config: {
                        requestId: 'getMainLine'
                    }
                })
            },
            idGenerator (data) {
                return `${data.bk_obj_id}_${data.bk_inst_id}`
            },
            applyEnabled (node) {
                return this.isModule(node) && node.data.host_apply_enabled
            },
            isTemplate (node) {
                return node.data.service_template_id || node.data.set_template_id
            },
            isModule (node) {
                return node.data.bk_obj_id === 'module'
            },
            async beforeSelect (node) {
                if (this.isModule(node)) {
                    if (this.$parent.editing) {
                        const leaveConfirmResult = await ConfirmStore.popup(LEAVE_CONFIRM_ID)
                        return !leaveConfirmResult
                    }
                    return true
                } else {
                    return false
                }
            },
            handleSelectChange (node) {
                this.$emit('selected', node)
            },
            handleCheckChange (id, checked) {
                this.$emit('checked', id, checked)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .topology-tree-wrapper {
        height: calc(100% - 32px);
        @include scrollbar-y;
    }
    .node-info {
        .node-model-icon {
            width: 22px;
            height: 22px;
            line-height: 21px;
            text-align: center;
            font-style: normal;
            font-size: 12px;
            margin: 8px 8px 0 6px;
            border-radius: 50%;
            background-color: #c4c6cc;
            color: #fff;
            &.is-template {
                background-color: #97aed6;
            }
            &.is-selected {
                background-color: #3a84ff;
            }
            &.is-leaf-icon {
                margin-left: 2px;
            }
        }
        .node-button {
            height: 24px;
            padding: 0 6px;
            margin: 0 20px 0 4px;
            line-height: 22px;
            border-radius: 4px;
            font-size: 12px;
            min-width: auto;
            &.set-template-button {
                @include inlineBlock;
                font-style: normal;
                background-color: #dcdee5;
                color: #ffffff;
                outline: none;
                cursor: not-allowed;
            }
        }
        .config-icon {
            margin: 6px 20px 6px 5px;
            padding: 0 5px;
            height: 18px;
            line-height: 17px;
            color: #979ba5;
            font-size: 26px;
            text-align: center;
            color: #2dcb56;
        }
        .info-content {
            display: flex;
            align-items: center;
            line-height: 36px;
            font-size: 14px;
            .node-name {
                @include ellipsis;
                margin-right: 8px;
            }
        }
    }
</style>
