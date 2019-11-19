<template>
    <div class="push-wrapper">
        <feature-tips
            :feature-name="'eventpush'"
            :show-tips="showFeatureTips"
            :desc="$t('事件推送顶部提示')"
            :more-href="'https://docs.bk.tencent.com/cmdb/Introduction.html#EventPush'"
            @close-tips="showFeatureTips = false">
        </feature-tips>
        <div class="btn-wrapper clearfix">
            <cmdb-auth class="inline-block-middle" :auth="$authResources({ type: $OPERATION.C_EVENT })">
                <bk-button slot-scope="{ disabled }"
                    theme="primary"
                    :disabled="disabled"
                    @click="createPush">
                    {{$t('新建')}}
                </bk-button>
            </cmdb-auth>
        </div>
        <bk-table
            v-bkloading="{ isLoading: $loading('searchSubscription') }"
            :data="table.list"
            :pagination="table.pagination"
            :max-height="$APP.height - 229"
            @sort-change="handleSortChange"
            @page-limit-change="handleSizeChange"
            @page-change="handlePageChange">
            <bk-table-column prop="subscription_name" :label="$t('订阅名称')" sortable="custom">
            </bk-table-column>
            <bk-table-column prop="system_name" :label="$t('系统名称')" sortable="custom">
            </bk-table-column>
            <bk-table-column prop="operator" :label="$t('操作人')" sortable="custom">
            </bk-table-column>
            <bk-table-column prop="last_time" :label="$t('更新时间')" sortable="custom">
            </bk-table-column>
            <bk-table-column prop="operation" :label="$t('推送情况（近一周）')">
                <template slot-scope="{ row }">
                    <i class="circle"
                        :class="{
                            'danger': row.operation.failure,
                            'success': !row.operation.failure
                        }">
                    </i>
                    {{$t('失败')}} {{row.operation.failure}} / {{$t('总量')}} {{row.operation.total}}
                </template>
            </bk-table-column>
            <bk-table-column prop="setting" :label="$t('配置')">
                <template slot-scope="{ row }">
                    <cmdb-auth class="mr10" :auth="$authResources({ type: $OPERATION.U_EVENT })">
                        <bk-button slot-scope="{ disabled }"
                            theme="primary"
                            text
                            :disabled="disabled"
                            @click.stop="editPush(row)">
                            {{$t('编辑')}}
                        </bk-button>
                    </cmdb-auth>
                    <cmdb-auth class="mr10" :auth="$authResources({ type: $OPERATION.D_EVENT })">
                        <bk-button slot-scope="{ disabled }"
                            theme="primary"
                            text
                            :disabled="disabled"
                            @click.stop="deleteConfirm(row)">
                            {{$t('删除')}}
                        </bk-button>
                    </cmdb-auth>
                </template>
            </bk-table-column>
            <cmdb-table-empty slot="empty" :stuff="table.stuff">
                <template>
                    <p>{{$t('暂时没有数据')}}</p>
                    <p>{{$t('事件推送功能提示')}}</p>
                </template>
            </cmdb-table-empty>
        </bk-table>
        <bk-sideslider
            v-transfer-dom
            :is-show.sync="slider.isShow"
            :title="slider.title"
            :width="564"
            :before-close="handleBeforeSliderClose">
            <v-push-detail
                ref="detail"
                slot="content"
                v-if="slider.isShow"
                :type="slider.type"
                :cur-push="curPush"
                @saveSuccess="saveSuccess"
                @cancel="closeSlider">
            </v-push-detail>
        </bk-sideslider>
    </div>
</template>

<script>
    import { formatTime } from '@/utils/tools'
    import featureTips from '@/components/feature-tips/index'
    import vPushDetail from './push-detail'
    import { mapActions, mapGetters } from 'vuex'
    export default {
        components: {
            vPushDetail,
            featureTips
        },
        data () {
            return {
                showFeatureTips: false,
                curPush: {},
                table: {
                    list: [],
                    pagination: {
                        count: 0,
                        current: 1,
                        ...this.$tools.getDefaultPaginationConfig()
                    },
                    defaultSort: '-last_time',
                    sort: '-last_time',
                    stuff: {
                        type: 'default',
                        payload: {}
                    }
                },
                slider: {
                    isShow: false,
                    isCloseConfirmShow: false,
                    title: '',
                    type: 'create'
                }
            }
        },
        computed: {
            ...mapGetters(['featureTipsParams'])
        },
        created () {
            this.getTableData()
            this.showFeatureTips = this.featureTipsParams['eventpush']
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
                            title: this.$t('确认退出'),
                            subTitle: this.$t('退出会导致未保存信息丢失'),
                            extCls: 'bk-dialog-sub-header-center',
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
                this.slider.title = this.$t('新增订阅')
            },
            editPush (item) {
                this.curPush = { ...item }
                this.slider.isShow = true
                this.slider.type = 'edit'
                this.slider.title = this.$t('编辑订阅')
            },
            deleteConfirm (item) {
                this.$bkInfo({
                    title: this.$tc('删除推送确认', item['subscription_name'], { name: item['subscription_name'] }),
                    confirmFn: () => {
                        this.deletePush(item['subscription_id'])
                    }
                })
            },
            async deletePush (subscriptionId) {
                await this.unsubcribeEvent({ bkBizId: 0, subscriptionId })
                this.$success(this.$t('删除推送成功'))
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
                        start: (pagination.current - 1) * pagination.limit,
                        limit: pagination.limit,
                        sort: this.table.sort
                    }
                }
                try {
                    const res = await this.searchSubscription({
                        bkBizId: 0,
                        params,
                        config: {
                            requestId: 'searchSubscription',
                            globalPermission: false
                        }
                    })
                    if (res.count && !res.info.length) {
                        this.table.pagination.current -= 1
                        this.getTableData()
                    }
                    res.info.map(item => {
                        item['subscription_form'] = item['subscription_form'].split(',')
                        item['last_time'] = formatTime(item['last_time'])
                    })
                    this.table.list = res.info
                    pagination.count = res.count
                } catch ({ permission }) {
                    if (permission) {
                        this.table.stuff = {
                            type: 'permission',
                            payload: { permission }
                        }
                    }
                }
            },
            handleSortChange (sort) {
                this.table.sort = this.$tools.getSort(sort)
                this.handlePageChange(1)
            },
            handleSizeChange (size) {
                this.table.pagination.limit = size
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
    .push-wrapper {
        padding: 0 20px;
    }
    .btn-wrapper {
        margin-bottom: 14px;
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
