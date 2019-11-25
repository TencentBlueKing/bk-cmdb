<template>
    <div class="property-confirm-table">
        <bk-table
            :data="data"
            :pagination="pagination"
            :row-style="{ cursor: 'pointer' }"
            :max-height="$APP.height - 240"
            @page-change="handlePageChange"
            @row-click="handleRowClick"
        >
            <bk-table-column :label="$t('内网IP')" prop="expired_host.bk_host_innerip" class-name="is-highlight"></bk-table-column>
            <bk-table-column :label="$t('云区域')" prop="expired_host.bk_cloud_id"></bk-table-column>
            <bk-table-column :label="$t('固资编号')" prop="expired_host.bk_asset_id"></bk-table-column>
            <bk-table-column :label="$t('主机名称')" prop="expired_host.bk_host_name"></bk-table-column>
            <bk-table-column :label="$t('修改值')" width="430" class-name="table-cell-change-value">
                <template slot-scope="{ row }">
                    <div class="cell-change-value" v-html="getChangeValue(row)"></div>
                </template>
            </bk-table-column>
            <bk-table-column :label="$t('操作')">
                <template slot-scope="{}">
                    <bk-button theme="primary" text>处理冲突</bk-button>
                </template>
            </bk-table-column>
        </bk-table>
        <bk-sideslider
            v-transfer-dom
            :width="800"
            :is-show.sync="slider.isShow"
            :title="slider.title"
        >
            <template slot="content">
                <cmdb-details
                    v-if="slider.content === 'detail'"
                    :show-options="false"
                    :inst="details.inst"
                    :properties="details.properties"
                    :property-groups="details.propertyGroups">
                </cmdb-details>
                <conflict-resolve v-else-if="slider.content === 'conflict'"></conflict-resolve>
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
            data: {
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
                    count: this.total,
                    limit: 20
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
                }
            }
        },
        computed: {
            ...mapGetters('objectModelClassify', [
                'getModelById'
            ]),
            ...mapGetters('hosts', ['configPropertyList']),
            ...mapState('hosts', ['propertyList'])
        },
        watch: {
        },
        created () {
            this.getHostPropertyList()
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
                    if (conflicts.find(conflictItem => conflictItem.bk_attribute_id === item.bk_attribute_id)) {
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
            handlePageChange (page) {
                this.pagination.current = page
            },
            handleRowClick (item) {
                const hasConflict = item.diff_value.some(item => item.is_conflict)
                if (hasConflict) {
                    this.handleShowConflict(item)
                } else {
                    this.handleShowDetails(item)
                }
            },
            async handleShowDetails (item) {
                this.slider.title = `属性详情【${item.bk_host_innerip}】`
                this.slider.content = 'detail'
                const properties = this.propertyList
                const inst = item
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
            handleShowConflict (item) {
                this.slider.title = `处理冲突【${item.bk_host_innerip}】`
                this.slider.content = 'conflict'
                this.slider.isShow = true
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
