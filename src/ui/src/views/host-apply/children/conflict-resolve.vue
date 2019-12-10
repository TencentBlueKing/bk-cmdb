<template>
    <div class="conflict-layout">
        <feature-tips
            class-name="resolve-tips"
            feature-name="hostApply"
            :show-tips="showFeatureTips"
            :desc="$t('因为主机复用模块的自动应用策略不一致，导致策略失效，需要手动维护不一致的属性。要彻底解决此问题，可以修改复用模块的策略为一致或移除模块的策略配置')"
            @close-tips="showFeatureTips = false"
        >
        </feature-tips>
        <div class="conflict-table" ref="conflictTable">
            <bk-table :data="conflictPropertyList">
                <bk-table-column :label="$t('字段名称')" width="160" :resizable="false" prop="bk_property_name"></bk-table-column>
                <bk-table-column :label="$t('所属模块')" :render-header="renderColumnHeader" :resizable="false">
                    <template slot-scope="{ row, $index }">
                        <div class="conflict-modules">
                            <div
                                v-for="(item, index) in row.__extra__.conflictList"
                                :class="['module-item', { 'selected': item.selected }]"
                                :key="index"
                            >
                                <span v-if="item.is_current">{{$t('主机当前值')}}</span>
                                <span v-else :title="getModulePath(item.bk_module_id)">{{getModulePath(item.bk_module_id)}}</span>
                                <span :title="item.bk_property_value | formatter(row.bk_property_type, row.option)">
                                    {{item.bk_property_value | formatter(row.bk_property_type, row.option)}}
                                </span>
                                <i class="check-model-value" @click="handlePickValue(row, index, item.bk_property_value, $index)">{{$t('选定')}}</i>
                            </div>
                        </div>
                    </template>
                </bk-table-column>
                <bk-table-column :label="$t('修改后')"
                    width="230"
                    class-name="table-cell-form-element"
                    :resizable="false"
                >
                    <template slot-scope="{ row }">
                        <property-form-element :property="row" @value-change="handlePropertyValueChange"></property-form-element>
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
    import featureTips from '@/components/feature-tips/index'
    import RESIZE_EVENTS from '@/utils/resize-events'
    import propertyFormElement from './property-form-element'
    export default {
        components: {
            featureTips,
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
                confirmButtonDisabled: false,
                showFeatureTips: false
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('hostApply', ['configPropertyList']),
            ...mapGetters(['featureTipsParams']),
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
            this.showFeatureTips = this.featureTipsParams['hostApplyConfirm']
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
            handlePropertyValueChange () {
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
        margin: 12px 20px;
    }
    .conflict-table {
        padding: 0 20px;
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
                span {
                    width: 50%;
                    padding-right: 10px;
                    @include ellipsis;
                    &:nth-child(2) {
                        width: auto;
                        max-width: calc(50% - 30px);
                    }
                }
                .check-model-value {
                    display: none;
                    font-style: normal;
                    cursor: pointer;
                    color: #3A84FF;
                }
                &.selected {
                    color: #3A84FF;
                    background-color: #E1ECFF;
                    span:nth-child(2) {
                        max-width: 50%;
                    }
                }
                &:not(.selected):hover {
                    .check-model-value {
                        display: inline-block;
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

        .table-cell-module-path:not(.header-cell) {
            .cell {
                padding-top: 8px;
                padding-bottom: 8px;
                overflow: unset !important;
                display: block !important;
            }
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
