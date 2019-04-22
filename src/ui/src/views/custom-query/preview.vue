<template>
    <div class="userapi-preview-wrapper">
        <div class="userapi-preview" v-click-outside="closePreview">
            <h3 class="preview-title">{{$t("CustomQuery['预览查询']")}}</h3>
            <i class="bk-icon icon-close" @click="closePreview"></i>
            <cmdb-table
                :loading="$loading('searchHost')"
                :header="table.header"
                :list="table.list"
                :pagination.sync="table.pagination"
                :wrapperMinusHeight="220"
                @handlePageChange="handlePageChange"
                @handleSizeChange="handleSizeChange"
                @handleSortChange="handleSortChange">
                <template v-for="({id,name, property}, index) in table.header" :slot="id" slot-scope="{ item }">
                    <template>{{getHostCellText(property, item)}}</template>
                </template>
            </cmdb-table>
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
                        size: 10
                    },
                    sort: ''
                }
            }
        },
        computed: {
            allProperties () {
                let allProperties = []
                for (let key in this.attribute) {
                    allProperties = [...allProperties, ...this.attribute[key].properties]
                }
                return allProperties
            },
            previewParams () {
                let condition = this.$tools.clone(this.apiParams['info']['condition'])
                let hostCondition = condition.find(({bk_obj_id: objId}) => {
                    return objId === 'host'
                })
                hostCondition['fields'] = this.previewFields
                let previewParams = {
                    'bk_biz_id': this.apiParams['bk_biz_id'],
                    condition: condition,
                    page: {
                        start: (this.table.pagination.current - 1) * this.table.pagination.size,
                        limit: this.table.pagination.size,
                        sort: this.table.sort
                    }
                }
                return previewParams
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
                let text = []
                originalValues.forEach(value => {
                    const flatternedText = this.$tools.getPropertyText(property, value)
                    flatternedText ? text.push(flatternedText) : void (0)
                })
                return text.join(',') || '--'
            },
            getColumnProperty (propertyId, objId) {
                return this.allProperties.find(property => {
                    return property['bk_property_id'] === propertyId && property['bk_obj_id'] === objId
                })
            },
            setTableHeader () {
                let headerList = []
                this.tableHeader.map(propertyId => {
                    let header = null
                    if (propertyId === 'bk_set_name') {
                        header = {
                            objId: 'set',
                            id: 'bk_set_name',
                            name: this.$t("Hosts['集群']"),
                            sortable: false
                        }
                    } else if (propertyId === 'bk_module_name') {
                        header = {
                            objId: 'module',
                            id: 'bk_module_name',
                            name: this.$t("Hosts['模块']"),
                            sortable: false
                        }
                    } else if (propertyId === 'bk_biz_name') {
                        header = {
                            objId: 'biz',
                            id: 'bk_biz_name',
                            name: this.$t("Common['业务']")
                        }
                    } else {
                        let property = this.attribute.host.properties.find(property => propertyId === property['bk_property_id'])
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
            },
            handlePageChange (current) {
                this.table.pagination.current = current
                this.getPreviewList()
            },
            handleSizeChange (size) {
                this.table.pagination.size = size
                this.handlePageChange(1)
            },
            handleSortChange (sort) {
                this.table.sort = sort
                this.getPreviewList()
            },
            closePreview () {
                this.$emit('close')
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
        z-index: 99;
        .userapi-preview {
            position: absolute;
            max-width: 80%;
            min-width: 50%;
            max-height: 80%;
            min-height: 300px;
            margin: 0 auto;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            background: #fff;
            box-shadow: 0 0 8px 4px rgba(0, 0, 0, 0.1);
            border-radius: 2px;
            .preview-title {
                padding-left: 24px;
                line-height: 68px;
                font-size: 20px;
                color: #333948;
                font-weight: normal;
            }
            .icon-close {
                position: absolute;
                top: 12px;
                right: 12px;
                cursor: pointer;
                font-size: 12px;
            }
        }
    }
</style>
