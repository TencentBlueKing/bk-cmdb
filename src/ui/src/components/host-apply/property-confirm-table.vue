<template>
    <div class="property-confirm-table">
        <bk-table
            :data="table.list"
            :pagination="table.pagination"
            :max-height="maxHeight || ($APP.height - 220)"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange"
        >
            <bk-table-column :label="$t('内网IP')">
                <template slot-scope="{ row }">
                    <bk-button theme="primary" text @click="handleShowDetails(row)">{{row.expect_host.bk_host_innerip}}</bk-button>
                </template>
            </bk-table-column>
            <bk-table-column :label="$t('主机名称')" prop="expect_host.bk_host_name"></bk-table-column>
            <bk-table-column
                :label="$t('修改值')"
                width="530"
                class-name="table-cell-change-value"
                show-overflow-tooltip
                :render-header="(h, data) => renderTableHeader(h, data, $t('红色为属性冲突值'), { placement: 'right' })">
                <template slot-scope="{ row }">
                    <div class="cell-change-value"><vnodes :vnode="getChangeValue(row)"></vnodes></div>
                </template>
            </bk-table-column>
            <bk-table-column
                v-if="showOperation"
                :label="$t('操作')"
                :render-header="renderIcon ? (h, data) => renderTableHeader(h, data, $t('表格冲突处理提示'), { width: 275 }) : null">
                <template slot-scope="{ row }">
                    <bk-button
                        v-if="row.conflicts.length > 0"
                        theme="primary"
                        text
                        @click="handleShowConflict(row)"
                    >
                        {{$t(row.unresolved_conflict_count > 0 ? '手动修改' : '已修改')}}
                    </bk-button>
                    <span v-else>--</span>
                </template>
            </bk-table-column>
            <cmdb-table-empty slot="empty">
                <div>{{$t('暂无主机新转入的主机将会自动应用模块的主机属性')}}</div>
            </cmdb-table-empty>
        </bk-table>
        <bk-sideslider
            v-transfer-dom
            :width="800"
            :is-show.sync="slider.isShow"
            :title="slider.title"
            @hidden="handleSliderCancel">
            <template slot="content">
                <cmdb-details
                    v-if="slider.content === 'detail'"
                    :show-options="false"
                    :inst="details.inst"
                    :properties="details.properties"
                    :property-groups="details.propertyGroups">
                </cmdb-details>
                <conflict-resolve
                    v-else-if="slider.content === 'conflict'"
                    ref="conflictResolve"
                    :data-row="currentRow"
                    :data-cache="conflictResolveCache"
                    @cancel="handleSliderCancel"
                    @save="handleConflictSave"
                >
                </conflict-resolve>
            </template>
        </bk-sideslider>
    </div>
</template>

<script>
    import { mapGetters, mapState } from 'vuex'
    import conflictResolve from './conflict-resolve.vue'
    export default {
        components: {
            conflictResolve,
            vnodes: {
                functional: true,
                render: (h, ctx) => ctx.props.vnode
            }
        },
        props: {
            list: {
                type: Array,
                default: () => ([])
            },
            total: {
                type: Number
            },
            maxHeight: {
                type: [Number, String],
                default: 0
            },
            renderIcon: {
                type: Boolean,
                default: false
            },
            showOperation: {
                type: Boolean,
                default: true
            }
        },
        data () {
            return {
                table: {
                    list: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        ...this.$tools.getDefaultPaginationConfig()
                    }
                },
                details: {
                    show: false,
                    title: '',
                    inst: {},
                    properties: [],
                    propertyGroups: []
                },
                slider: {
                    width: 514,
                    isShow: false,
                    content: '',
                    title: ''
                },
                currentRow: {},
                conflictResolveResult: {}
            }
        },
        computed: {
            ...mapGetters('objectModelClassify', [
                'getModelById'
            ]),
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('hostApply', ['configPropertyList']),
            ...mapState('hostApply', ['propertyList']),
            conflictResolveCache () {
                const key = this.currentRow.bk_host_id
                return this.conflictResolveResult[key] || []
            }
        },
        watch: {
            list (data) {
                this.setTableList()
            },
            total (value) {
                this.table.pagination.count = value
            }
        },
        created () {
            this.getHostPropertyList()
            this.setTableList()
        },
        methods: {
            async getHostPropertyList () {
                try {
                    const data = await this.$store.dispatch('hostApply/getProperties', {
                        params: { bk_biz_id: this.bizId },
                        config: {
                            requestId: 'getHostPropertyList'
                        }
                    })

                    this.$store.commit('hostApply/setPropertyList', data)
                } catch (e) {
                    console.error(e)
                }
            },
            setTableList () {
                const { start, limit } = this.$tools.getPageParams(this.table.pagination)
                this.table.list = this.list.slice(start, start + limit)
            },
            getChangeValue (row) {
                const { conflicts, update_fields: updateFields, unresolved_conflict_count: conflictCount } = row

                const resultConflicts = conflicts.map(item => {
                    const property = this.configPropertyList.find(propertyItem => propertyItem.id === item.bk_attribute_id) || {}
                    let content = <span>{property.bk_property_name}：<cmdb-property-value value={item.bk_property_value} property={property} /></span>
                    const conflictExist = conflicts.find(conflictItem => conflictItem.bk_attribute_id === item.bk_attribute_id && conflictItem.unresolved_conflict_exist)
                    if (conflictExist) {
                        content = <span class="conflict-item">{content}</span>
                    }
                    return content
                })
                const resultUpdates = updateFields.map(item => {
                    const property = this.configPropertyList.find(propertyItem => propertyItem.id === item.bk_attribute_id) || {}
                    return <span>{property.bk_property_name}：<cmdb-property-value value={item.bk_property_value} property={property} /></span>
                })
                const conflictSeparator = <span class={`conflict-separator${!conflictCount ? ' resolved' : ''}`}>；</span>

                return (
                    <div>
                        { resultConflicts.reduce((acc, x) => acc === null ? [x] : [acc, conflictSeparator, x], null) }
                        { (resultConflicts.length && resultUpdates.length) ? conflictSeparator : '' }
                        { resultUpdates.reduce((acc, x) => acc === null ? [x] : [acc, '；', x], null) }
                    </div>
                )
            },
            getPropertyGroups () {
                return this.$store.dispatch('objectModelFieldGroup/searchGroup', {
                    objId: 'host',
                    params: { bk_biz_id: this.bizId }
                })
            },
            renderTableHeader (h, data, tips, options = {}) {
                const directive = {
                    content: tips,
                    placement: options.placement || 'top-end'
                }
                if (options.width) {
                    directive.width = options.width
                }
                return <span>{ data.column.label } <i class="bk-cc-icon icon-cc-tips" v-bk-tooltips={ directive }></i></span>
            },
            handlePageChange (page) {
                this.table.pagination.current = page
                this.setTableList()
            },
            handleSizeChange (size) {
                this.table.pagination.limit = size
                this.table.pagination.current = 1
                this.setTableList()
            },
            async handleShowDetails (row) {
                this.currentRow = row
                this.slider.title = `${this.$t('属性详情')}【${row.expect_host.bk_host_innerip}】`
                this.slider.content = 'detail'
                const properties = this.propertyList
                const inst = row.expect_host
                // 云区域数据
                row.cloud_area['bk_inst_name'] = row.cloud_area['bk_cloud_name']
                row.cloud_area['bk_inst_id'] = row.cloud_area['bk_cloud_id']
                inst['bk_cloud_id'] = [row.cloud_area]
                try {
                    const propertyGroups = await this.getPropertyGroups()
                    this.details.inst = inst
                    this.details.properties = properties
                    this.details.propertyGroups = propertyGroups
                    this.slider.isShow = true
                } catch (e) {
                    console.log(e)
                    this.details.inst = {}
                    this.details.properties = []
                    this.details.propertyGroups = []
                    this.slider.isShow = false
                }
            },
            handleShowConflict (row) {
                this.currentRow = row
                this.slider.title = `${this.$t('手动修改')}【${row.expect_host.bk_host_innerip}】`
                this.slider.content = 'conflict'
                this.slider.isShow = true
            },
            handleConflictSave (data) {
                // 将处理结果保存
                this.conflictResolveResult[this.currentRow.bk_host_id] = data

                // 标记冲突状态并更新属性值
                this.currentRow.conflicts.forEach(item => {
                    item.unresolved_conflict_exist = false
                    item.bk_property_value = data.find(property => property.id === item.bk_attribute_id).__extra__.value
                })
                this.currentRow.unresolved_conflict_count = 0
            },
            handleSliderCancel () {
                if (this.slider.content === 'conflict') {
                    if (this.$refs.conflictResolve) {
                        this.$refs.conflictResolve.restoreConflictPropertyList()
                    }
                }
                this.slider.isShow = false
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cell-change-value {
        padding: 8px 0;
        line-height: 1.6;
        word-break: break-word;
    }
</style>
<style lang="scss">
    .table-cell-change-value {
        .cell {
            -webkit-line-clamp: unset !important;
            display: block !important;
            .conflict-item,
            .conflict-separator {
                color: #ea3536;

                &.resolved {
                    color: inherit;
                }
            }

            .icon-cc-tips {
                margin-top: -2px;
            }
        }
    }
</style>
