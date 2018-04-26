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
                :tableHeader="table.header"
                :tableList="table.list"
                :pagination="table.pagination"
                :isLoading="table.isLoading"
                :maxHeight="table.maxHeight"
                @handlePageTurning="setCurrentPage"
                @handlePageSizeChange="setCurrentSize">
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
            }
        },
        data () {
            return {
                table: {
                    header: [{
                        id: 'bk_host_innerip',
                        name: this.$t("Common['内网IP']")
                    }, {
                        id: 'bk_biz_name',
                        name: this.$t("Common['业务']")
                    }, {
                        id: 'bk_set_name',
                        name: this.$t("Hosts['集群']")
                    }, {
                        id: 'bk_module_name',
                        name: this.$t("Hosts['模块']")
                    }, {
                        id: 'bk_cloud_id',
                        name: this.$t("Hosts['云区域ID']")
                    }],
                    list: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        size: 10
                    },
                    isLoading: false,
                    maxHeight: 0
                }
            }
        },
        computed: {
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
                hostCondition.fields = this.previewFields
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
            getPreviewList () {
                this.table.isLoading = true
                this.$axios.post('hosts/search', this.previewParams).then(res => {
                    if (res.result) {
                        this.table.pagination.count = res.data.count
                        this.initTableList(res.data.info)
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    this.table.isLoading = false
                }).catch(() => {
                    this.table.isLoading = false
                })
            },
            initTableList (list) {
                let tableList = []
                list.forEach((item, index) => {
                    let cellItem = {}
                    this.table.header.map(({id}) => {
                        Object.keys(item).map(bkObjId => {
                            if (item[bkObjId].hasOwnProperty(id)) {
                                let cellValue = item[bkObjId][id]
                                if (Array.isArray(cellValue)) {
                                    cellItem[id] = cellValue.map(({bk_inst_name: bkInstName}) => {
                                        return bkInstName
                                    }).join(',')
                                } else {
                                    cellItem[id] = cellValue
                                }
                            }
                        })
                    })
                    tableList.push(cellItem)
                })
                this.table.list = tableList
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
        z-index: 3;
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
            padding: 0 40px;
            margin-top: 6px;
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
    .userapi-preview .preview-table .table-scrollbar{
        max-height: 250px !important;
    }
    .userapi-preview .preview-table .min-height-control{
        min-height: 250px;
    }
</style>