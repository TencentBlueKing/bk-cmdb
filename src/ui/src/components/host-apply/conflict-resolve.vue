<template>
    <div class="conflict-layout">
        <cmdb-tips class="resolve-tips" tips-key="hostApplyConfirmTips">{{$t('策略失效提示语')}}</cmdb-tips>
        <div class="conflict-table-wrapper" ref="conflictTable">
            <bk-table
                ext-cls="conflict-table"
                :data="conflictPropertyList"
                :cell-style="{ background: '#fff' }"
            >
                <bk-table-column :label="$t('字段名称')" width="160" :resizable="false" prop="bk_property_name"></bk-table-column>
                <bk-table-column :label="$t('所属模块')" :render-header="renderColumnHeader" :resizable="false">
                    <template slot-scope="{ row, $index }">
                        <div class="conflict-modules">
                            <div
                                v-for="(item, index) in row.__extra__.conflictList"
                                :class="['module-item', { 'selected': item.selected }]"
                                :key="index"
                                @click="handlePickValue(row, index, item.bk_property_value, $index)"
                            >
                                <span v-if="item.is_current">{{$t('主机当前值')}}</span>
                                <span v-else :title="getModulePath(item.bk_module_id)">{{getModulePath(item.bk_module_id)}}</span>
                                <span :title="item.bk_property_value | formatter(row.bk_property_type, row.option)">
                                    {{item.bk_property_value | formatter(row.bk_property_type, row.option)}}
                                </span>
                            </div>
                        </div>
                    </template>
                </bk-table-column>
                <bk-table-column :label="$t('修改后')"
                    width="230"
                    class-name="table-cell-form-element"
                    :resizable="false"
                >
                    <template slot-scope="{ row, $index }">
                        <property-form-element :property="row" @value-change="value => handlePropertyValueChange(value, row, $index)"></property-form-element>
                    </template>
                </bk-table-column>
            </bk-table>
        </div>
        <div :class="['footer-btns', { 'sticky': scrollbar }]">
            <bk-button theme="primary" class="mr10" :disabled="confirmButtonDisabled" @click="handleConfirm">{{$t('确定')}}</bk-button>
            <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import RESIZE_EVENTS from '@/utils/resize-events'
    import propertyFormElement from './property-form-element'
    export default {
        components: {
            propertyFormElement
        },
        props: {
            dataRow: {
                type: Object,
                default: () => ({})
            },
            dataCache: {
                type: Array,
                default: () => ([])
            }
        },
        data () {
            return {
                conflictPropertyList: [],
                conflictPropertyListSnapshot: [],
                result: {},
                moduleMap: {},
                scrollbar: false,
                confirmButtonDisabled: false
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('hostApply', ['configPropertyList']),
            moduleIds () {
                return this.dataRow.bk_module_ids
            }
        },
        watch: {
            dataRow () {
                this.getData()
            }
        },
        created () {
            this.getData()
        },
        mounted () {
            RESIZE_EVENTS.addResizeListener(this.$refs.conflictTable, this.checkScrollbar)
        },
        beforeDestroy () {
            RESIZE_EVENTS.removeResizeListener(this.$refs.conflictTable, this.checkScrollbar)
        },
        methods: {
            async getData () {
                try {
                    const topopath = await this.getTopopath()
                    const moduleMap = {}
                    topopath.nodes.forEach(node => {
                        moduleMap[node.topo_node.bk_inst_id] = node.topo_path
                    })
                    this.moduleMap = Object.freeze(moduleMap)
                    this.setConflictPropertyList()
                } catch (e) {
                    console.error(e)
                }
            },
            setConflictPropertyList () {
                let conflictPropertyList = []
                if (this.dataCache.length) {
                    conflictPropertyList = this.$tools.clone(this.dataCache)
                } else {
                    const conflicts = this.dataRow.conflicts || []
                    conflicts.forEach(item => {
                        const findProperty = this.configPropertyList.find(property => property.bk_property_id === item.bk_property_id)
                        const property = this.$tools.clone(findProperty)
                        // 主机当前值
                        const hostCurrent = {
                            is_current: true,
                            // 任何一个配置值都未被使用则使用主机当前值
                            selected: item.host_apply_rules.every(item => !item.selected),
                            bk_attribute_id: item.bk_attribute_id,
                            bk_property_value: item.bk_property_value
                        }
                        property.__extra__.conflictList = [hostCurrent, ...item.host_apply_rules]
                        property.__extra__.value = hostCurrent.bk_property_value
                        conflictPropertyList.push(property)
                    })
                }
                this.conflictPropertyList = conflictPropertyList
                this.conflictPropertyListSnapshot = this.$tools.clone(conflictPropertyList)
            },
            getTopopath () {
                return this.$store.dispatch('hostApply/getTopopath', {
                    bizId: this.bizId,
                    params: {
                        topo_nodes: this.moduleIds.map(id => ({ bk_obj_id: 'module', bk_inst_id: id }))
                    }
                })
            },
            getModulePath (id) {
                const info = this.moduleMap[id] || []
                const path = info.map(node => node.bk_inst_name).reverse().join(' / ')
                return path
            },
            getModuleName (id) {
                const topoInfo = this.moduleMap[id] || []
                const target = topoInfo.find(target => target.bk_obj_id === 'module' && target.bk_inst_id === id) || {}
                return target.bk_inst_name
            },
            renderColumnHeader (h, { column, $index }) {
                const style = {
                    width: '50%',
                    display: 'inline-block'
                }
                return h(
                    'div',
                    { class: 'conflict-custom-label' },
                    [
                        h('span', { style }, this.$t('所属模块')),
                        h('span', { style }, this.$t('冲突值'))
                    ]
                )
            },
            checkScrollbar () {
                const $layout = this.$el
                this.scrollbar = $layout.scrollHeight !== $layout.offsetHeight
            },
            restoreConflictPropertyList () {
                this.conflictPropertyList = this.$tools.clone(this.conflictPropertyListSnapshot)
            },
            toggleConfirmButtonDisabled () {
                const everyValidTruthy = this.conflictPropertyList.every(property => property.__extra__.valid !== false)
                this.confirmButtonDisabled = !everyValidTruthy
            },
            handlePickValue (row, index, value, $index) {
                row.__extra__.conflictList.forEach((item, i) => {
                    this.$set(this.conflictPropertyList[$index].__extra__.conflictList[i], 'selected', index === i)
                })
                row.__extra__.value = value
            },
            handleConfirm () {
                this.conflictPropertyListSnapshot = this.$tools.clone(this.conflictPropertyList)
                this.$emit('save', this.conflictPropertyListSnapshot)
                this.$emit('cancel')
            },
            handleCancel () {
                this.restoreConflictPropertyList()
                this.$emit('cancel')
            },
            handlePropertyValueChange (value, row, $index) {
                let hasSelected = false
                row.__extra__.conflictList.forEach((item, i) => {
                    const selected = item.bk_property_value === value
                    this.$set(this.conflictPropertyList[$index].__extra__.conflictList[i], 'selected', selected && !hasSelected)
                    hasSelected = selected
                })
                this.toggleConfirmButtonDisabled()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .conflict-layout {
        height: 100%;
        @include scrollbar-y;
    }
    .resolve-tips {
        margin: 12px 20px 0;
    }
    .conflict-table-wrapper {
        padding: 12px 20px 0;
        .conflict-name {
            font-weight: bold;
        }
        .conflict-modules {
            height: 100%;
            padding: 6px 0;
            .module-item {
                display: flex;
                line-height: 22px;
                padding-left: 6px;
                cursor: pointer;
                span {
                    width: 50%;
                    padding-right: 10px;
                    @include ellipsis;
                    &:nth-child(2) {
                        width: auto;
                        max-width: calc(50% - 30px);
                    }
                }
                &:hover {
                    color: #3a84ff;
                    background-color: #e1ecff;
                }
                &.selected {
                    color: #3a84ff;
                    background-color: #f0f1f5;
                    span:nth-child(2) {
                        max-width: 50%;
                    }
                }
            }
        }
    }
    .footer-btns {
        position: sticky;
        bottom: 0;
        left: 0;
        width: 100%;
        padding: 20px 20px 0;
        background-color: #FFFFFF;
        z-index: 10;
        &.sticky {
            height: 54px;
            line-height: 54px;
            padding: 0 20px;
            margin-top: 20px;
            border-top: 1px solid #dcdee5;
            font-size: 0;
        }
    }
</style>

<style lang="scss">
    .conflict-table {
        overflow: unset;
        .bk-table-body-wrapper {
            overflow: unset !important;
        }

        .table-cell-form-element {
            .cell {
                overflow: unset !important;
                display: block !important;
            }
            .search-input-wrapper {
                position: relative;
            }
            .form-objuser .suggestion-list {
                z-index: 1000;
            }
        }
    }
</style>
