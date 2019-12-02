<template>
    <div class="property-confirm-table">
        <bk-table
            :data="tableList"
            :pagination="pagination"
            :row-style="{ cursor: 'pointer' }"
            :max-height="$APP.height - 240"
            @page-change="handlePageChange"
            @page-limit-change="handleSizeChange"
            @row-click="handleRowClick"
        >
            <bk-table-column :label="$t('内网IP')" prop="expired_host.bk_host_innerip" class-name="is-highlight"></bk-table-column>
            <bk-table-column :label="$t('云区域')" prop="cloud_area.bk_cloud_name"></bk-table-column>
            <bk-table-column :label="$t('固资编号')" prop="expired_host.bk_asset_id"></bk-table-column>
            <bk-table-column :label="$t('主机名称')" prop="expired_host.bk_host_name"></bk-table-column>
            <bk-table-column :label="$t('修改值')" width="430" class-name="table-cell-change-value">
                <template slot-scope="{ row }">
                    <div class="cell-change-value" v-html="getChangeValue(row)"></div>
                </template>
            </bk-table-column>
            <bk-table-column :label="$t('操作')" class-name="is-highlight" :formatter="getOperationColumnText"></bk-table-column>
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
    import { mapGetters, mapState, mapActions } from 'vuex'
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
            }
        },
        data () {
            return {
                pagination: {
                    current: 1,
                    count: 0,
                    ...this.$tools.getDefaultPaginationConfig()
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
                conflictResolveResult: {},
                tableList: []
            }
        },
        computed: {
            ...mapGetters('objectModelClassify', [
                'getModelById'
            ]),
            ...mapGetters('hosts', ['configPropertyList']),
            ...mapState('hosts', ['propertyList']),
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
                this.pagination.count = value
            }
        },
        created () {
            this.getHostPropertyList()
            this.setTableList()
        },
        methods: {
            ...mapActions('objectModelProperty', [
                'searchObjectAttribute'
            ]),
            async getHostPropertyList () {
                try {
                    const data = await this.searchObjectAttribute({
                        params: this.$injectMetadata({
                            bk_obj_id: 'host',
                            bk_supplier_account: this.supplierAccount
                        }),
                        config: {
                            requestId: 'getHostPropertyList',
                            fromCache: true
                        }
                    })

                    this.$store.commit('hosts/setPropertyList', data)
                } catch (e) {
                    console.error(e)
                }
            },
            setTableList () {
                const { start, limit } = this.$tools.getPageParams(this.pagination)
                this.tableList = this.list.slice(start, start + limit)
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
                    text = row.unresolved_conflict_count > 0 ? '处理冲突' : '已处理'
                }
                return text
            },
            handlePageChange (page) {
                this.pagination.current = page
                this.setTableList()
            },
            handleSizeChange (size) {
                this.pagination.limit = size
                this.setTableList()
            },
            handleRowClick (row) {
                this.currentRow = row
                const conflictExist = row.unresolved_conflict_count > 0
                if (conflictExist) {
                    this.handleShowConflict(row)
                } else {
                    this.handleShowDetails(row)
                }
            },
            async handleShowDetails (row) {
                this.slider.title = `属性详情【${row.expired_host.bk_host_innerip}】`
                this.slider.content = 'detail'
                const properties = this.propertyList
                const inst = row
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
                color: red;
            }
        }
    }
</style>
