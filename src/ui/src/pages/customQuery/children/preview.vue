/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

<template>
    <div class="userapi-preview-wrapper" v-show="isPreviewShow" v-click-outside="closePreview">
        <div class="userapi-preview">
            <h3 class="preview-title">{{$t("CustomQuery['测试']")}}</h3>
            <v-table class="preview-table"
                :sortable="false"
                :header="table.header"
                :list="table.list"
                :pagination.sync="table.pagination"
                :loading="table.isLoading"
                :maxHeight="325"
                :width="515"
                @handlePageChange="setCurrentPage"
                @handleSizeChange="setCurrentSize">
                <template v-for="({id,name, property}, index) in table.header" :slot="id" slot-scope="{ item }">
                    <template>{{getCellValue(property, item)}}</template>
                </template>
            </v-table>
            <div class="preview-btn-group">
                <bk-button type="primary" class="preview-btn-confirm" @click="closePreview">{{$t("Common['确认']")}}</bk-button>
            </div>
        </div>
    </div>
</template>
<script>
    import vTable from '@/components/table/table'
    export default {
        components: {
            vTable
        },
        props: {
            isPreviewShow: {
                type: Boolean,
                default: false
            },
            apiParams: {
                type: Object,
                required: true
            },
            attribute: {
                type: Object,
                required: true
            }
        },
        data () {
            return {
                table: {
                    header: [{
                        objId: 'host',
                        id: 'bk_host_innerip',
                        name: this.$t("Common['内网IP']"),
                        property: []
                    }, {
                        objId: 'biz',
                        id: 'bk_biz_name',
                        name: this.$t("Common['业务']"),
                        property: []
                    }, {
                        objId: 'set',
                        id: 'bk_set_name',
                        name: this.$t("Hosts['集群']"),
                        property: []
                    }, {
                        objId: 'module',
                        id: 'bk_module_name',
                        name: this.$t("Hosts['模块']"),
                        property: []
                    }, {
                        objId: 'host',
                        id: 'bk_cloud_id',
                        name: this.$t("Hosts['云区域']"),
                        property: []
                    }],
                    list: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        size: 10
                    },
                    isLoading: false
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
            previewFields () {
                let [, , , hostCondition] = this.apiParams['info']['condition']
                let fields = []
                this.table.header.map(header => {
                    fields.push(header.id)
                })
                return [...new Set(fields.concat(hostCondition['fields']))]
            },
            previewParams () {
                let condition = this.$deepClone(this.apiParams['info']['condition'])
                let hostCondition = condition.find(({bk_obj_id: bkObjId}) => {
                    return bkObjId === 'host'
                })
                hostCondition['fields'] = this.previewFields
                let previewParams = {
                    'bk_biz_id': this.apiParams['bk_biz_id'],
                    condition: condition,
                    page: {
                        start: (this.table.pagination.current - 1) * this.table.pagination.size,
                        limit: this.table.pagination.size
                    }
                }
                return previewParams
            }
        },
        watch: {
            isPreviewShow (isShow) {
                if (isShow) {
                    this.getPreviewList()
                } else {
                    this.table.list = []
                }
            }
        },
        methods: {
            setTableHeader (columns) {
                this.table.header.map(header => {
                    const property = this.getColumnProperty(header['id'], header['objId'])
                    this.$set(header, 'property', property)
                })
            },
            getColumnProperty (bkPropertyId, bkObjId) {
                return this.allProperties.find(property => {
                    return property['bk_property_id'] === bkPropertyId && property['bk_obj_id'] === bkObjId
                })
            },
            getCellValue (property, item) {
                if (property) {
                    let bkObjId = property['bk_obj_id']
                    let value = item[bkObjId][property['bk_property_id']]
                    if (property['bk_property_id'] === 'bk_module_name') {
                        let moduleName = []
                        item.module.map(({bk_module_name: bkModuleName}) => {
                            moduleName.push(bkModuleName)
                        })
                        return moduleName.join(',')
                    }
                    if (property['bk_property_id'] === 'bk_set_name') {
                        let setName = []
                        item.set.map(({bk_set_name: bksetName}) => {
                            setName.push(bksetName)
                        })
                        return setName.join(',')
                    }
                    if (property['bk_property_id'] === 'bk_biz_name') {
                        let bizName = []
                        item.biz.map(({bk_biz_name: bkbizName}) => {
                            bizName.push(bkbizName)
                        })
                        return bizName.join(',')
                    }
                    if (property['bk_asst_obj_id'] && Array.isArray(value)) {
                        let tempValue = []
                        value.map(({bk_inst_name: bkInstName}) => {
                            if (bkInstName) {
                                tempValue.push(bkInstName)
                            }
                        })
                        value = tempValue.join(',')
                    } else if (property['bk_property_type'] === 'date') {
                        value = this.$formatTime(value, 'YYYY-MM-DD')
                    } else if (property['bk_property_type'] === 'time') {
                        value = this.$formatTime(value)
                    } else if (property['bk_property_type'] === 'enum') {
                        let option = property.option.find(({id}) => {
                            return id === value
                        })
                        if (option) {
                            value = option.name
                        } else {
                            value = ''
                        }
                    }
                    return value
                }
                return ''
            },
            getPreviewList () {
                this.table.isLoading = true
                this.$axios.post('hosts/search', this.previewParams).then(res => {
                    if (res.result) {
                        this.table.pagination.count = res.data.count
                        this.setTableHeader()
                        this.table.list = res.data.info
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    this.table.isLoading = false
                }).catch(() => {
                    this.table.isLoading = false
                })
            },
            setCurrentPage (current) {
                this.table.pagination.current = current
                this.getPreviewList()
            },
            setCurrentSize (size) {
                this.table.pagination.size = size
                this.setCurrentPage(1)
            },
            closePreview () {
                this.$emit('update:isPreviewShow', false)
            }
        }
    }
</script>
<style lang="scss" scoped>
    .userapi-preview-wrapper{
        position: absolute;
        top: 150px;
        left: 0;
        width: 100%;
        z-index: 99;
    }
    .userapi-preview{
        width: 595px;
        max-height: 466px;
        background-color: #ffffff;
        box-shadow: 0px 2px 8px 0px rgba(0, 0, 0, 0.4);
        text-align: left;
        margin: 0 auto;
        padding: 22px 0 0 0;
        border-radius: 1px;
        .preview-title{
            font-size: 14px;
            line-height: 26px;
            color: #737987;
            margin: 0 40px 15px;
        }
        .preview-table{
            margin: 6px 40px 0;
        }
        .preview-btn-group{
            text-align: right;
            padding: 20px 40px 40px;
            .preview-btn-confirm{
                width: 110px;
            }
        }
    }
</style>
<style lang="scss">
    .preview-table{
        .table-pagination{
            padding: 0 6px !important;
            .bk-page{
                height: 26px;
                margin: 8px 0;
                ul{
                    height: 26px;
                }
                .page-item{
                    min-width: 26px;
                    height: 26px;
                    line-height: 26px;
                }
            }
        }
    }
</style>
