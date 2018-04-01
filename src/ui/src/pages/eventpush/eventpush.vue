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
    <div class="content-box">
        <div class="button-wrapper">
            <bk-button type="primary" @click="addPush">新增推送</bk-button>
        </div>
        <v-table
            @resetTablePage="getTableListAfterPageSetting"
            :paging="paging"
            :hasCheckbox="false"
            :tableHeader="tableHeader"
            :tableList="tableList"
            @editEvent="editEvent"
            @delEvent="delConfirm"
            @pageTuring="getTableList"
            @handleTableSortClick="handleTableSortClick"
        >
        </v-table>
        <v-sideslider
            :title="sliderTitle"
            :isShow.sync="isSliderShow"
            :hasQuickClose="true"
            @closeSlider="closeSlider"
        >
            <v-push-detail
                :isShow="isSliderShow"
                :type="operationType"
                :curEvent="curEvent"
                @saveSuccess="saveSuccess"
                @cancel="closeSlider"
                slot="content"
            ></v-push-detail>
        </v-sideslider>
    </div>
</template>

<script>
    import vTable from '@/components/table/table2'
    import vSideslider from '@/components/slider/sideslider'
    import vPushDetail from './children/pushDetail.vue'
    import {mapGetters} from 'vuex'
    export default {
        data () {
            return {
                curEvent: {},
                operationType: '',              // 当前事件推送操作类型  add/edit
                isSliderShow: false,            // 弹窗状态
                sliderTitle: {
                    text: '新增推送',
                    icon: 'icon-cc-edit'
                },
                paging: {
                    sort: '-ID'
                },
                tableHeader: [
                    {
                        id: 'subscription_name',
                        name: '名称'
                    },
                    {
                        id: 'system_name',
                        name: '系统'
                    },
                    {
                        id: 'operator',
                        name: '更新人'
                    },
                    {
                        id: 'last_time',
                        name: '更新时间'
                    },
                    {
                        id: '__slot__eventpush__statistics',
                        name: '推送情况（近一周）'
                    },
                    {
                        id: '__slot__eventpush__setting',
                        name: '配置'
                    }
                ],
                tableList: [],
                pageLimit: 10
            }
        },
        computed: {
            ...mapGetters([
                'bkSupplierAccount'
            ])
        },
        methods: {
            /*
                获取推送列表
            */
            getTableList (index) {
                if (!index) {
                    index = 1
                }
                let appid = 0
                this.$axios.post(`event/subscribe/search/${this.bkSupplierAccount}/${appid}`, {
                    page: {
                        start: (index - 1) * this.pageLimit,
                        limit: this.pageLimit,
                        sort: this.paging.sort
                    }
                }).then(res => {
                    if (res.result) {
                        res.data.info.map(val => {
                            val['subscription_form'] = val['subscription_form'].split(',')
                            val['last_time'] = this.$formatTime(val['last_time'])
                        })
                        
                        this.tableList = res.data.info
                        this.paging = {
                            current: index,
                            showItem: this.pageLimit,
                            totalPage: Math.ceil(res.data.count / this.pageLimit),
                            count: res.data.count,
                            sort: this.paging.sort
                        }
                    } else {
                        this.$alertMsg('获取推送列表失败')
                    }
                })
            },
            /*
                编辑某一项事件推送
            */
            editEvent (item) {
                // 记录当前操作项
                this.curEvent = {...item}
                this.isSliderShow = true
                this.operationType = 'edit'
                this.sliderTitle.text = '编辑推送'
            },
            /*
                保存推送成功回调
            */
            saveSuccess (id) {
                if (this.operationType === 'add') {
                    this.operationType = 'edit'
                    this.curEvent['subscription_id'] = id
                    this.getTableList()
                } else {
                    this.getTableList(this.paging.current)
                }
            },
            /*
                删除的确认弹窗
            */
            delConfirm (item) {
                // 记录当前操作项
                this.curEvent = item
                let self = this
                this.$bkInfo({
                    title: `确定删除名称为 ${item['subscription_name']} 的推送？`,
                    confirmFn () {
                        self.delEvent(item)
                    }
                })
            },
            /*
                删除某一项事件推送
            */
            delEvent (item) {
                let appid = 0
                this.$axios.delete(`event/subscribe/${this.bkSupplierAccount}/${appid}/${item['subscription_id']}`).then(res => {
                    if (res.result) {
                        this.$alertMsg('删除推送成功', 'success')
                        this.getTableList(this.paging.current)
                    } else {
                        this.$alertMsg('删除推送失败')
                    }
                })
            },
            /*
                新增推送
            */
            addPush () {
                this.operationType = 'add'
                this.isSliderShow = true
                this.sliderTitle.text = '新增推送'
            },
            /*
                关闭推送弹窗
            */
            closeSlider () {
                this.isSliderShow = false
            },
            handleTableSortClick (sort) {
                if (sort.indexOf('__slot__') === -1) {
                    this.paging.sort = sort
                    this.getTableList()
                }
            },
            /*
                设置分页后重新获取表格数据
            */
            getTableListAfterPageSetting (page) {
                this.pageLimit = page
                this.getTableList()
            }
        },
        mounted () {
            this.getTableList()
        },
        components: {
            vTable,
            vSideslider,
            vPushDetail
        }
    }
</script>

<style lang="scss" scoped>
    .content-box{
        width: 100%;
        padding: 20px;
        .button-wrapper{
            text-align: right;
            margin-bottom: 20px;
            .btn-add{
                width: 124px;
                height: 36px;
            }
        }
    }
</style>
