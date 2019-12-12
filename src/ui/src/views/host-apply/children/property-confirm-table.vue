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
                    <bk-button theme="primary" text @click="handleShowDetails(row)">{{row.expired_host.bk_host_innerip}}</bk-button>
                </template>
            </bk-table-column>
            <bk-table-column :label="$t('云区域')" prop="cloud_area.bk_cloud_name"></bk-table-column>
            <bk-table-column :label="$t('固资编号')" prop="expired_host.bk_asset_id"></bk-table-column>
            <bk-table-column :label="$t('主机名称')" prop="expired_host.bk_host_name"></bk-table-column>
            <bk-table-column :label="$t('修改值')" width="430" class-name="table-cell-change-value">
                <template slot-scope="{ row }">
                    <div class="cell-change-value" v-html="getChangeValue(row)"></div>
                </template>
            </bk-table-column>
            <bk-table-column
                v-if="showOperation"
                :label="$t('操作')"
                :render-header="renderIcon ? (h, data) => renderTableHeader(h, data, $t('表格冲突处理提示')) : null"
            >
                <template slot-scope="{ row }">
                    <bk-button
                        v-if="row.conflicts.length > 0"
                        theme="primary"
                        text
                        @click="handleShowConflict(row)"
                    >
                        {{row.unresolved_conflict_count > 0 ? '手动修改' : '已修改'}}
                    </bk-button>
                    <span v-else>--</span>
                </template>
            </bk-table-column>
            <cmdb-table-empty slot="empty">
                <div>{{$t('暂无主机，新转入的主机将会自动应用模块的主机属性')}}</div>
            </cmdb-table-empty>
        </bk-table>
        <bk-sideslider
            v-transfer-dom
            :width="800"
            :is-show.sync="slider.isShow"
            :title="slider.title"
            @hidden="handleSliderCancel"
        >
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
            conflictResolve
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
                type: Number,
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
                    title: this.$t('拓扑显示设置')
                },
                currentRow: {},
                conflictResolveResult: {}
            }
        },
        computed: {
            ...mapGetters('objectModelClassify', [
                'getModelById'
            ]),
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
                        requestId: 'getHostPropertyList',
                        fromCache: true
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
                const { conflicts, update_fields: updateFields } = row
                const valueMap = {}
                const fieldList = [...conflicts, ...updateFields]
                fieldList.forEach(item => {
                    valueMap[item['bk_property_id']] = item.bk_property_value
                })
                const result = fieldList.map(item => {
                    const property = this.configPropertyList.find(propertyItem => propertyItem.id === item.bk_attribute_id) || {}
                    let content = `${property.bk_property_name}：${this.$tools.getPropertyText(property, valueMap)}`
                    const conflictExist = conflicts.find(conflictItem => conflictItem.bk_attribute_id === item.bk_attribute_id && conflictItem.unresolved_conflict_exist)
                    if (conflictExist) {
                        content = `<span class="conflict">${content}</span>`
                    }
                    return content
                })
                return result.join('；')
            },
            getPropertyGroups () {
                const modelId = 'host'
                const model = this.getModelById(modelId)
                return this.$store.dispatch('objectModelFieldGroup/searchGroup', {
                    objId: modelId,
                    params: this.$injectMetadata({}, {
                        inject: !!this.$tools.getMetadataBiz(model)
                    })
                })
            },
            getOperationColumnText (row) {
                let text = '--'
                if (row.conflicts.length > 0) {
                    text = row.unresolved_conflict_count > 0 ? '手动修改' : '已修改'
                }
                return text
            },
            renderTableHeader (h, data, tips) {
                const directive = {
                    content: tips,
                    placement: 'top-end',
                    width: 275
                }
                return <span>{ data.column.label } <i class="bk-cc-icon icon-cc-tips" v-bk-tooltips={ directive }></i></span>
            },
            handlePageChange (page) {
                this.table.pagination.current = page
                this.setTableList()
            },
            handleSizeChange (size) {
                this.table.pagination.limit = size
                this.setTableList()
            },
            async handleShowDetails (row) {
                this.currentRow = row
                this.slider.title = `属性详情【${row.expired_host.bk_host_innerip}】`
                this.slider.content = 'detail'
                const properties = this.propertyList
                const inst = row.expired_host
                // 云区域数据
                row.cloud_area['bk_inst_name'] = row.cloud_area['bk_cloud_name']
                inst['bk_cloud_id'] = [row.cloud_area]
                try {
                    const propertyGroups = await this.getPropertyGroups()
                    this.details.inst = this.$tools.flattenItem(properties, inst)
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
                this.slider.title = `处理冲突【${row.expired_host.bk_host_innerip}】`
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
</style>
<style lang="scss">
    .table-cell-change-value {
        .cell {
            -webkit-line-clamp: 3;
            .conflict {
                color: #ea3536;
            }
        }
    }
</style>
