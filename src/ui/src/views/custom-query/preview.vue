<template>
    <div class="userapi-preview-wrapper">
        <div class="mask" @click="closePreview"></div>
        <div class="userapi-preview">
            <h3 class="preview-title">{{$t('预览查询')}}</h3>
            <i class="bk-icon icon-close" @click="closePreview"></i>
            <div class="preview-table">
                <bk-table
                    ref="table"
                    v-bkloading="{ isLoading: $loading('searchHost') }"
                    :data="table.list"
                    :pagination="table.pagination"
                    :height="379"
                    @page-change="handlePageChange"
                    @page-limit-change="handleSizeChange"
                    @sort-change="handleSortChange">
                    <bk-table-column v-for="column in table.header"
                        :sortable="unSortableProperty.includes(column.id) ? false : 'custom'"
                        :key="column.id"
                        :prop="column.id"
                        :label="column.name"
                        show-overflow-tooltip>
                        <template slot-scope="{ row }">{{getHostCellText(column.property, row)}}</template>
                    </bk-table-column>
                </bk-table>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    export default {
        props: {
            apiParams: {
                type: Object,
                required: true
            },
            attribute: {
                type: Object,
                required: true
            },
            tableHeader: {
                type: Array
            }
        },
        data () {
            return {
                table: {
                    header: [],
                    list: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        limit: 10
                    },
                    sort: ''
                }
            }
        },
        computed: {
            allProperties () {
                let allProperties = []
                for (const key in this.attribute) {
                    allProperties = [...allProperties, ...this.attribute[key].properties]
                }
                return allProperties
            },
            previewParams () {
                const conditions = this.$tools.clone(this.apiParams['info']['condition'])
                const hostCondition = conditions.find(({ bk_obj_id: objId }) => {
                    return objId === 'host'
                })
                hostCondition['fields'] = this.previewFields
                conditions.forEach(model => {
                    const modelCondition = model.condition || []
                    const newConditions = []
                    if (modelCondition.length) {
                        modelCondition.forEach(condition => {
                            const value = condition.value
                            if ((condition.operator === '$multilike' && value !== null && value !== undefined && String(value).length)
                                || condition.operator !== '$multilike') {
                                newConditions.push(condition)
                            }
                        })
                    }
                    model.condition = newConditions
                })
                const previewParams = {
                    'bk_biz_id': this.apiParams['bk_biz_id'],
                    condition: conditions,
                    page: {
                        start: (this.table.pagination.current - 1) * this.table.pagination.limit,
                        limit: this.table.pagination.limit,
                        sort: this.table.sort
                    }
                }
                return previewParams
            },
            unSortableProperty () {
                return ['bk_set_name', 'bk_module_name', 'bk_cloud_id']
            }
        },
        created () {
            this.setTableHeader()
            this.getPreviewList()
        },
        methods: {
            ...mapActions('hostSearch', [
                'searchHost'
            ]),
            getHostCellText (property, item) {
                const objId = property['bk_obj_id']
                const originalValues = item[objId] instanceof Array ? item[objId] : [item[objId]]
                const text = []
                originalValues.forEach(value => {
                    const flattenedText = this.$tools.getPropertyText(property, value)
                    flattenedText ? text.push(flattenedText) : void (0)
                })
                return text.join(',') || '--'
            },
            getColumnProperty (propertyId, objId) {
                return this.allProperties.find(property => {
                    return property['bk_property_id'] === propertyId && property['bk_obj_id'] === objId
                })
            },
            setTableHeader () {
                const headerList = []
                this.tableHeader.map(propertyId => {
                    let header = null
                    if (propertyId === 'bk_set_name') {
                        header = {
                            objId: 'set',
                            id: 'bk_set_name',
                            name: this.$t('集群'),
                            sortable: false
                        }
                    } else if (propertyId === 'bk_module_name') {
                        header = {
                            objId: 'module',
                            id: 'bk_module_name',
                            name: this.$t('模块'),
                            sortable: false
                        }
                    } else if (propertyId === 'bk_biz_name') {
                        header = {
                            objId: 'biz',
                            id: 'bk_biz_name',
                            name: this.$t('业务')
                        }
                    } else {
                        const property = this.attribute.host.properties.find(property => propertyId === property['bk_property_id'])
                        if (property) {
                            header = {
                                objId: 'host',
                                id: property['bk_property_id'],
                                name: property['bk_property_name']
                            }
                        }
                    }
                    if (header) {
                        const property = this.getColumnProperty(propertyId, header.objId)
                        this.$set(header, 'property', property)
                        headerList.push(header)
                    }
                })
                this.table.header = headerList
            },
            async getPreviewList () {
                const res = await this.searchHost({
                    params: this.previewParams,
                    config: {
                        requestId: 'searchHost'
                    }
                })
                this.table.pagination.count = res.count
                this.table.list = res.info
                this.fixPageLimitPosition()
            },
            handlePageChange (current) {
                this.table.pagination.current = current
                this.getPreviewList()
            },
            handleSizeChange (size) {
                this.table.pagination.limit = size
                this.handlePageChange(1)
            },
            handleSortChange (sort) {
                this.table.sort = this.$tools.getSort(sort)
                this.getPreviewList()
            },
            closePreview () {
                this.$emit('close')
            },
            fixPageLimitPosition () {
                this.$nextTick(() => {
                    const limitRefs = this.$refs.table.$el.querySelector('.bk-table-pagination .bk-page-count .bk-tooltip-ref')
                    limitRefs && limitRefs._tippy.set({ boundary: 'window' })
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .userapi-preview-wrapper {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: 2400;
        .mask {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background-color: rgba(0, 0, 0, 0.6);
            z-index: 1;
        }
        .userapi-preview {
            position: absolute;
            width: 880px;
            max-height: 80%;
            min-height: 300px;
            margin: 20px auto;
            top: 50%;
            left: 50%;
            z-index: 2;
            transform: translate(-50%, -50%);
            background: #fff;
            box-shadow: 0 0 8px 4px rgba(0, 0, 0, 0.1);
            border-radius: 2px;
            .preview-title {
                padding: 15px 0 15px 24px;
                font-size: 24px;
                color: #444444;
                font-weight: normal;
            }
            .icon-close {
                position: absolute;
                top: 12px;
                right: 12px;
                cursor: pointer;
                font-size: 14px;
                font-weight: bold;
            }
        }
        .preview-table {
            padding: 0 20px 20px;
        }
    }
</style>
