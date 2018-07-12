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
        <div class="button-wrapper clearfix">
            <bk-button class="fr" type="primary" @click="addPush">{{$t('EventPush["新增推送"]')}}</bk-button>
        </div>
        <v-table
            :header="table.header"
            :list="table.list"
            :defaultSort="table.defaultSort"
            :pagination.sync="table.pagination"
            :loading="table.isLoading"
            :wrapperMinusHeight="150"
            @handlePageChange="setCurrentPage"
            @handleSizeChange="setCurrentSize"
            @handleSortChange="setCurrentSort">
                <template slot="statistics" slot-scope="{ item }">
                    <i class="circle" :class="[{'danger':item.statistics.failure},{'success':!item.statistics.failure}]"></i>
                    {{$t('EventPush[\'失败\']')}} {{item.statistics.failure}} / {{$t('EventPush[\'总量\']')}} {{item.statistics.total}}
                </template>
                <template slot="setting" slot-scope="{ item }">
                    <i class="icon-cc-edit mr20" @click.stop="editEvent(item)"></i>
                    <i class="icon-cc-del" @click.stop="delConfirm(item)"></i>
                </template>
        </v-table>
        <v-sideslider
            :title="sliderTitle"
            :isShow.sync="slider.isShow"
            :hasCloseConfirm="true"
            :isCloseConfirmShow="slider.isCloseConfirmShow"
            @closeSlider="closeSliderConfirm"
        >
            <v-push-detail
                ref="detail"
                :isShow="slider.isShow"
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
    import vTable from '@/components/table/table'
    import vSideslider from '@/components/slider/sideslider'
    import vPushDetail from './children/pushDetail.vue'
    import {mapGetters} from 'vuex'
    export default {
        data () {
            return {
                curEvent: {},
                operationType: '',              // 当前事件推送操作类型  add/edit
                slider: {
                    isShow: false,            // 弹窗状态
                    isCloseConfirmShow: false
                },
                sliderTitle: {
                    text: '',
                    icon: 'icon-cc-edit'
                },
                table: {
                    header: [{
                        id: 'subscription_name',
                        name: this.$t('EventPush["推送名称"]')
                    }, {
                        id: 'system_name',
                        name: this.$t('EventPush["系统名称"]')
                    }, {
                        id: 'operator',
                        name: this.$t('EventPush["操作人"]')
                    }, {
                        id: 'last_time',
                        name: this.$t('EventPush["更新时间"]')
                    }, {
                        id: 'statistics',
                        name: this.$t('EventPush["推送情况（近一周）"]'),
                        sortable: false
                    }, {
                        id: 'setting',
                        name: this.$t('EventPush["配置"]'),
                        sortable: false
                    }],
                    list: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        size: 10
                    },
                    defaultSort: '-last_time',
                    sort: '-last_time'
                }
            }
        },
        computed: {
            ...mapGetters([
                'bkSupplierAccount',
                'language'
            ])
        },
        watch: {
            'language' () {
                this.table.header = [{
                    id: 'subscription_name',
                    name: this.$t('EventPush["推送名称"]')
                }, {
                    id: 'system_name',
                    name: this.$t('EventPush["系统名称"]')
                }, {
                    id: 'operator',
                    name: this.$t('EventPush["操作人"]')
                }, {
                    id: 'last_time',
                    name: this.$t('EventPush["更新时间"]')
                }, {
                    id: 'statistics',
                    name: this.$t('EventPush["推送情况（近一周）"]'),
                    sortable: false
                }, {
                    id: 'setting',
                    name: this.$t('EventPush["配置"]'),
                    sortable: false
                }]
            }
        },
        methods: {
            /*
                获取推送列表
            */
            getTableList () {
                let appid = 0
                let pagination = this.table.pagination
                this.table.isLoading = true
                this.$axios.post(`event/subscribe/search/${this.bkSupplierAccount}/${appid}`, {
                    page: {
                        start: (pagination.current - 1) * pagination.size,
                        limit: pagination.size,
                        sort: this.table.sort
                    }
                }).then(res => {
                    if (res.result) {
                        res.data.info.map(val => {
                            val['subscription_form'] = val['subscription_form'].split(',')
                            val['last_time'] = this.$formatTime(val['last_time'])
                        })
                        
                        this.table.list = res.data.info
                        pagination.count = res.data.count
                    } else {
                        this.$alertMsg(this.$t('EventPush["获取推送列表失败"]'))
                    }
                    this.table.isLoading = false
                }).catch(() => {
                    this.table.isLoading = false
                })
            },
            /*
                编辑某一项事件推送
            */
            editEvent (item) {
                this.curEvent = {...item}
                this.slider.isShow = true
                this.operationType = 'edit'
                this.sliderTitle.text = this.$t('EventPush["编辑推送"]')
            },
            /*
                保存推送成功回调
            */
            saveSuccess (id) {
                if (this.operationType === 'add') {
                    this.operationType = 'edit'
                    this.curEvent['subscription_id'] = id
                    this.setCurrentPage(1)
                } else {
                    this.getTableList()
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
                    title: this.$tc('EventPush["删除推送确认"]', item['subscription_name'], {name: item['subscription_name']}),
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
                        this.$alertMsg(this.$t('EventPush["删除推送成功"]'), 'success')
                        this.getTableList()
                    } else {
                        this.$alertMsg(this.$t('EventPush["删除推送失败"]'))
                    }
                })
            },
            /*
                新增推送
            */
            addPush () {
                this.operationType = 'add'
                this.slider.isShow = true
                this.sliderTitle.text = this.$t('EventPush["新增推送"]')
            },
            closeSlider () {
                this.slider.isShow = false
            },
            closeSliderConfirm () {
                this.slider.isCloseConfirmShow = this.$refs.detail.isCloseConfirmShow()
            },
            setCurrentPage (current) {
                this.table.pagination.current = current
                this.getTableList()
            },
            setCurrentSize (size) {
                this.table.pagination.size = size
                this.setCurrentPage(1)
            },
            setCurrentSort (sort) {
                this.table.sort = sort
                this.setCurrentPage(1)
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
            margin-bottom: 20px;
            .btn-add{
                width: 124px;
                height: 36px;
            }
        }
    }
    .circle{
        display: inline-block;
        vertical-align: baseline;
        width: 8px;
        height: 8px;
        margin-right: 5px;
        border-radius: 50%;
        &.success{
            background: #30d878;
        }
        &.danger{
            background: #ef4c4c;
        }
    }
</style>
