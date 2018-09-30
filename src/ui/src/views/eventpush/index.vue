<template>
    <div class="push-wrapper">
        <div class="btn-wrapper clearfix">
            <bk-button type="primary" @click="createPush">{{$t('EventPush["新增推送"]')}}</bk-button>
        </div>
        <cmdb-table
            rowCursor="default"
            :loading="$loading('searchSubscription')"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :defaultSort="table.defaultSort"
            :wrapperMinusHeight="150"
            @handleSortChange="handleSortChange"
            @handleSizeChange="handleSizeChange"
            @handlePageChange="handlePageChange">
                <template slot="statistics" slot-scope="{ item }">
                    <i class="circle" :class="[{'danger':item.statistics.failure},{'success':!item.statistics.failure}]"></i>
                    {{$t('EventPush[\'失败\']')}} {{item.statistics.failure}} / {{$t('EventPush[\'总量\']')}} {{item.statistics.total}}
                </template>
                <template slot="setting" slot-scope="{ item }">
                    <span class="text-primary mr20" @click.stop="editPush(item)">{{$t('Common["编辑"]')}}</span>
                    <span class="text-danger" @click.stop="deleteConfirm(item)">{{$t('Common["删除"]')}}</span>
                </template>
                <div class="empty-info" slot="data-empty">
                    <p>{{$t("Common['暂时没有数据']")}}</p>
                    <p>{{$t("EventPush['当前并无推送，可点击下方按钮新增']")}}</p>
                    <bk-button class="process-btn" type="primary" @click="createPush">{{$t("EventPush['新增推送']")}}</bk-button>
                </div>
        </cmdb-table>
        <cmdb-slider
            :isShow.sync="slider.isShow"
            :title="slider.title"
            :width="564"
            :beforeClose="handleSliderBeforeClose">
            <v-push-detail
                ref="detail"
                slot="content"
                :type="slider.type"
                :curPush="curPush"
                @saveSuccess="saveSuccess"
                @cancel="closeSlider">
            </v-push-detail>
        </cmdb-slider>
    </div>
</template>

<script>
    import { formatTime } from '@/utils/tools'
    import vPushDetail from './push-detail'
    import { mapActions } from 'vuex'
    export default {
        data () {
            return {
                curPush: {},
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
                        count: 0,
                        size: 10,
                        current: 1
                    },
                    defaultSort: '-last_time',
                    sort: '-last_time'
                },
                slider: {
                    isShow: false,
                    isCloseConfirmShow: false,
                    title: '',
                    type: 'create'
                }
            }
        },
        methods: {
            ...mapActions('eventSub', [
                'searchSubscription',
                'unsubcribeEvent'
            ]),
            handleSliderBeforeClose () {
                if (this.$refs.detail.isCloseConfirmShow()) {
                    return new Promise((resolve, reject) => {
                        this.$bkInfo({
                            title: this.$t('Common["退出会导致未保存信息丢失，是否确认？"]'),
                            confirmFn: () => {
                                resolve(true)
                            },
                            cancelFn: () => {
                                resolve(false)
                            }
                        })
                    })
                }
                return true
            },
            createPush () {
                this.slider.isShow = true
                this.slider.type = 'create'
                this.slider.title = this.$t('EventPush["新增推送"]')
            },
            editPush (item) {
                this.curPush = {...item}
                this.slider.isShow = true
                this.slider.type = 'edit'
                this.slider.title = this.$t('EventPush["编辑推送"]')
            },
            deleteConfirm (item) {
                this.$bkInfo({
                    title: this.$tc('EventPush["删除推送确认"]', item['subscription_name'], {name: item['subscription_name']}),
                    confirmFn: () => {
                        this.deletePush(item['subscription_id'])
                    }
                })
            },
            async deletePush (subscriptionId) {
                await this.unsubcribeEvent({bkBizId: 0, subscriptionId})
                this.$success(this.$t('EventPush["删除推送成功"]'))
                this.getTableData()
            },
            saveSuccess () {
                if (this.slider.type === 'create') {
                    this.handlePageChange(1)
                } else {
                    this.getTableData()
                }
                this.slider.isShow = false
            },
            closeSlider () {
                this.slider.isShow = false
            },
            async getTableData () {
                let pagination = this.table.pagination
                let params = {
                    page: {
                        start: (pagination.current - 1) * pagination.size,
                        limit: pagination.size,
                        sort: this.table.sort
                    }
                }
                let res = await this.searchSubscription({bkBizId: 0, params, config: {requestId: 'searchSubscription'}})
                res.info.map(item => {
                    item['subscription_form'] = item['subscription_form'].split(',')
                    item['last_time'] = formatTime(item['last_time'])
                })
                this.table.list = res.info
                pagination.count = res.count
            },
            handleSortChange (sort) {
                this.table.sort = sort
                this.handlePageChange(1)
            },
            handleSizeChange (size) {
                this.table.pagination.size = size
                this.handlePageChange(1)
            },
            handlePageChange (page) {
                this.table.pagination.current = page
                this.getTableData()
            }
        },
        created () {
            this.getTableData()
        },
        components: {
            vPushDetail
        }
    }
</script>

<style lang="scss" scoped>
    .btn-wrapper {
        margin-bottom: 20px;
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
            background: $cmdbDangerColor;
        }
    }
</style>
