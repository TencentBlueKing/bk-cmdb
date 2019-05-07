<template>
    <div class="push-wrapper">
        <div class="feature-tips">
            <i class="icon-cc-exclamation-tips"></i>
            <span>{{$t("EventPush['事件推送顶部提示']")}}</span>
            <a href="https://docs.bk.tencent.com/cmdb/Introduction.html#EventPush" target="_blank">{{$t("Common['更多详情']")}} >></a>
        </div>
        <div class="btn-wrapper clearfix">
            <bk-button type="primary"
                :disabled="!$isAuthorized(OPERATION.C_EVENT)"
                @click="createPush">
                {{$t('Common["新建"]')}}
            </bk-button>
        </div>
        <cmdb-table
            row-cursor="default"
            :loading="$loading('searchSubscription')"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :default-sort="table.defaultSort"
            :wrapper-minus-height="150"
            @handleSortChange="handleSortChange"
            @handleSizeChange="handleSizeChange"
            @handlePageChange="handlePageChange">
            <template slot="statistics" slot-scope="{ item }">
                <i class="circle" :class="[{ 'danger': item.statistics.failure },{ 'success': !item.statistics.failure }]"></i>
                {{$t('EventPush[\'失败\']')}} {{item.statistics.failure}} / {{$t('EventPush[\'总量\']')}} {{item.statistics.total}}
            </template>
            <template slot="setting" slot-scope="{ item }">
                <span class="text-primary mr20"
                    v-if="$isAuthorized(OPERATION.U_EVENT)"
                    @click.stop="editPush(item)">
                    {{$t('Common["编辑"]')}}
                </span>
                <span class="text-primary disabled mr20" v-else>{{$t('Common["编辑"]')}}</span>
                <span class="text-danger"
                    v-if="$isAuthorized(OPERATION.D_EVENT)"
                    @click.stop="deleteConfirm(item)">
                    {{$t('Common["删除"]')}}
                </span>
                <span class="text-danger disabled" v-else>{{$t('Common["删除"]')}}</span>
            </template>
            <div class="empty-info" slot="data-empty">
                <p>{{$t("Common['暂时没有数据']")}}</p>
                <p>{{$t("EventPush['事件推送功能提示']")}}</p>
            </div>
        </cmdb-table>
        <cmdb-slider
            :is-show.sync="slider.isShow"
            :title="slider.title"
            :width="564"
            :before-close="handleBeforeSliderClose">
            <v-push-detail
                ref="detail"
                slot="content"
                :type="slider.type"
                :cur-push="curPush"
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
    import { OPERATION } from './router.config.js'
    export default {
        components: {
            vPushDetail
        },
        data () {
            return {
                OPERATION,
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
        created () {
            this.$store.commit('setHeaderTitle', this.$t('Nav["事件推送"]'))
            this.getTableData()
        },
        methods: {
            ...mapActions('eventSub', [
                'searchSubscription',
                'unsubcribeEvent'
            ]),
            handleBeforeSliderClose () {
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
                this.curPush = { ...item }
                this.slider.isShow = true
                this.slider.type = 'edit'
                this.slider.title = this.$t('EventPush["编辑推送"]')
            },
            deleteConfirm (item) {
                this.$bkInfo({
                    title: this.$tc('EventPush["删除推送确认"]', item['subscription_name'], { name: item['subscription_name'] }),
                    confirmFn: () => {
                        this.deletePush(item['subscription_id'])
                    }
                })
            },
            async deletePush (subscriptionId) {
                await this.unsubcribeEvent({ bkBizId: 0, subscriptionId })
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
                const pagination = this.table.pagination
                const params = {
                    page: {
                        start: (pagination.current - 1) * pagination.size,
                        limit: pagination.size,
                        sort: this.table.sort
                    }
                }
                const res = await this.searchSubscription({ bkBizId: 0, params, config: { requestId: 'searchSubscription' } })
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
