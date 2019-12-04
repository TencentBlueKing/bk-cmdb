<template>
    <div class="conflict-layout">
        <cmdb-tips class="resolve-tips"
            :tips-style="{
                height: '32px',
                'line-height': '32px',
                'font-size': '12px'
            }"
            :icon-style="{ 'font-size': '16px' }">
            {{$t('解决应用字段冲突提示')}}
        </cmdb-tips>
        <div class="conflict-table" ref="conflictTable">
            <bk-table :data="conflictPropertyList">
                <bk-table-column :label="$t('字段名称')" width="160" :resizable="false" prop="bk_property_name"></bk-table-column>
                <bk-table-column :label="$t('所属模块')" :render-header="hanldeColumnRender" :resizable="false">
                    <template slot-scope="{ row }">
                        <div class="conflict-modules">
                            <div
                                v-for="(item, index) in row.__extra__.conflictList"
                                :class="['module-item', { 'selected': item.selected }]"
                                :key="index"
                            >
                                <span v-if="item.is_current">主机当前值</span>
                                <span v-else :title="getModulePath(item.bk_module_id)">{{getModulePath(item.bk_module_id)}}</span>
                                <span :title="item.bk_property_value | formatter(row.bk_property_type, row.option)">
                                    {{item.bk_property_value | formatter(row.bk_property_type, row.option)}}
                                </span>
                                <i class="check-model-value" @click="handlePickValue(row, index, item.bk_property_value)">选定</i>
                            </div>
                        </div>
                    </template>
                </bk-table-column>
                <bk-table-column :label="$t('修改后')"
                    width="230"
                    class-name="conflict-custom-column"
                    :resizable="false">
                    <template slot-scope="{ row }">
                        <component class="property-component"
                            :is="`cmdb-form-${row.bk_property_type}`"
                            :class="[row.bk_property_type, { error: errors.has(row.bk_property_id) }]"
                            :options="row.option || []"
                            :data-vv-name="row.bk_property_id"
                            :data-vv-as="row.bk_property_name"
                            :placeholder="$t('请输入xx', { name: row.bk_property_name })"
                            v-validate="$tools.getValidateRules(row)"
                            v-model.trim="row.__extra__.value"
                        >
                        </component>
                    </template>
                </bk-table-column>
            </bk-table>
        </div>
        <div :class="['footer-btns', { 'sticky': scrollbar }]">
            <bk-button theme="primary" class="mr10" @click="handleConfirm">{{$t('确定')}}</bk-button>
            <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import RESIZE_EVENTS from '@/utils/resize-events'
    export default {
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
                scrollbar: false,
                moduleMap: {}
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
            hanldeColumnRender (h, { column, $index }) {
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
            handlePickValue (row, index, value) {
                row.__extra__.conflictList.forEach((item, i) => (item.selected = index === i))
                row.__extra__.value = value
            },
            restoreConflictPropertyList () {
                this.conflictPropertyList = this.$tools.clone(this.conflictPropertyListSnapshot)
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
            checkScrollbar () {
                const $layout = this.$el
                this.scrollbar = $layout.scrollHeight !== $layout.offsetHeight
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
            .module-item {
                display: flex;
                line-height: 22px;
                span {
                    max-width: 50%;
                    padding-right: 10px;
                    @include ellipsis;
                    &:first-child {
                        width: calc(50% - 30px);
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
    .conflict-custom-column {
        .cell {
            overflow: unset;
        }
    }
</style>
