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
    <div id="userAPI" class="api-wrapper">
        <div class="options clearfix">
            <div class="fl clearfix bizbox">
                <div class="options-business fl">
                    <v-application-selector 
                    :filterable="true"
                    :selected.sync="filter.bkBizId">
                    </v-application-selector>
                </div>
            </div>
            <div class="options-btn-group fl ml10">
                <bk-button type="primary" class="options-btn mr10" style="display: none">{{$t("CustomQuery['新增指引']")}}</bk-button>
                <bk-button type="primary" @click="showUserAPISlider('create')">{{$t("CustomQuery['新增查询']")}}</bk-button>
            </div>
            <div class="options-search fr">
                <input class="bk-form-input" type="text" :placeholder="$t('Inst[\'快速查询\']')"
                v-model.trim="filter.name"
                    @keyup.enter="setCurrentPage(1)"
                >
                <i class="bk-icon icon-search options-search-icon" @click="setCurrentPage(1)"></i>
            </div>
        </div>
        <div class="userAPI-table">
            <v-table ref="userAPITable"
                :header="table.header"
                :list="table.list"
                :defaultSort="table.defaultSort"
                :pagination.sync="table.pagination"
                :loading="table.isLoading"
                :wrapperMinusHeight="150"
                @handlePageChange="setCurrentPage"
                @handleSizeChange="setCurrentSize"
                @handleSortChange="setCurrentSort"
                @handleRowClick="showUserAPIDetails">
            </v-table>
        </div>
        <div class="userAPI-slider">
            <v-sideslider ref="userAPISlider"
                :isShow.sync="slider.isShow"
                :hasQuickClose="true"
                :title="slider.title"
                :hasCloseConfirm="true"
                :isCloseConfirmShow="slider.isCloseConfirmShow"
                @closeSlider="closeSliderConfirm"
                @close="hideUserAPISlider">
                <v-define slot="content" 
                    ref="define"
                    :id="slider.id"
                    :bkBizId="filter.bkBizId"
                    :isShow="slider.isShow"
                    :type="slider.type"
                    @delete="getUserAPIList"
                    @create="handleCreate"
                    @update="getUserAPIList"
                    @cancel="hideUserAPISlider">
                </v-define>
            </v-sideslider>
        </div>
    </div>
</template>

<script>
    import vTable from '@/components/table/table'
    import vSideslider from '@/components/slider/sideslider'
    import vDefine from './children/define'
    import vApplicationSelector from '@/components/common/selector/application'
    import { mapGetters } from 'vuex'
    export default {
        components: {
            vTable,
            vSideslider,
            vDefine,
            vApplicationSelector
        },
        data () {
            return {
                filter: {
                    bkBizId: '',
                    name: ''
                },
                table: {
                    header: [{
                        id: 'id',
                        name: 'ID'
                    }, {
                        id: 'name',
                        name: this.$t("CustomQuery['查询名称']")
                    }, {
                        id: 'create_user',
                        name: this.$t("CustomQuery['创建用户']")
                    }, {
                        id: 'create_time',
                        name: this.$t("CustomQuery['创建时间']")
                    }, {
                        id: 'modify_user',
                        name: this.$t("CustomQuery['修改人']")
                    }, {
                        id: 'last_time',
                        name: this.$t("CustomQuery['修改时间']")
                    }],
                    list: [],
                    sort: '-last_time',
                    defaultSort: '-last_time',
                    pagination: {
                        current: 1,
                        count: 0,
                        size: 10
                    },
                    isLoading: false
                },
                slider: {
                    isShow: false,
                    isCloseConfirmShow: false,
                    type: 'create',
                    id: null,
                    title: {
                        text: this.$t("CustomQuery['新增查询']"),
                        icon: 'icon-cc-edit'
                    }
                }
            }
        },
        computed: {
            ...mapGetters(['bkBizList']),
            /* 构造表格查询参数 */
            searchParams () {
                let params = {
                    start: (this.table.pagination.current - 1) * this.table.pagination.size,
                    limit: this.table.pagination.size,
                    sort: this.table.sort
                }
                this.filter.bkBizId ? params['bk_biz_id'] = this.filter.bkBizId : void (0)
                this.filter.name ? params['condition'] = {'name': this.filter.name} : void (0)
                return params
            }
        },
        watch: {
            'filter.bkBizId' (newID, oldID) {
                this.getUserAPIList()
            }
        },
        methods: {
            closeSliderConfirm () {
                this.slider.isCloseConfirmShow = this.$refs.define.isCloseConfirmShow()
            },
            /* 获取自定义API列表 */
            getUserAPIList () {
                this.table.isLoading = true
                this.$axios.post(`userapi/search/${this.filter.bkBizId}`, this.searchParams).then((res) => {
                    if (res.result) {
                        if (res.data.count) {
                            res.data.info.forEach((listItem) => {
                                listItem['create_time'] = this.$formatTime(listItem['create_time'])
                                listItem['last_time'] = this.$formatTime(listItem['last_time'])
                            })
                            this.table.list = res.data.info
                        } else {
                            this.table.list = []
                        }
                        this.table.pagination.count = res.data.count
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    this.table.isLoading = false
                }).catch((e) => {
                    this.table.isLoading = false
                    this.table.list = []
                    if (e.response && e.response.status === 403) {
                        this.$alertMsg(this.$t("Common['您没有当前业务的权限']"))
                    }
                })
            },
            handleCreate (data) {
                this.slider.id = data['id']
                this.slider.type = 'update'
                this.setCurrentPage(1)
            },
            /* 显示自定义API详情 */
            showUserAPIDetails (userAPI) {
                this.slider.isShow = true
                this.slider.type = 'update'
                this.slider.id = userAPI['id']
                this.slider.title.text = this.$t('CustomQuery["编辑查询"]')
            },
            /* 设置当前页码 */
            setCurrentPage (current) {
                this.table.pagination.current = current
                this.getUserAPIList()
            },
            /* 设置每页显示条数 */
            setCurrentSize (size) {
                this.table.pagination.size = size
                this.setCurrentPage(1)
            },
            /* 设置排序 */
            setCurrentSort (sort) {
                this.table.sort = sort
                this.setCurrentPage(1)
            },
            /* 显示编辑自定义条件侧滑栏 */
            showUserAPISlider (type) {
                this.slider.isShow = true
                this.slider.type = type
                this.slider.title.text = this.$t("CustomQuery['新增查询']")
            },
            /* 隐藏自定义条件侧滑栏 */
            hideUserAPISlider () {
                this.slider.isShow = false
                this.slider.id = null
            }
        }
    }
</script>

<style lang="scss" scoped>
    .api-wrapper{
        height: 100%;
        padding: 20px;
    }
    .options{
        font-size: 14px;
        .bizbox{
            width: 170px;
        }
    }
    .options-search{
        position: relative;
        .bk-form-input{
            width: 320px;
        }
        .options-search-icon{
            position: absolute;
            top: 9px;
            right: 10px;
            font-size: 18px;
            color: #bec6de;
        }
    }
    .options-btn-group{
        font-size: 0;
    }
    .options-btn{
        width: 133px;
    }
    .userAPI-table{
        margin-top: 20px;
    }
</style>
