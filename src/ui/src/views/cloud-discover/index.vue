<template>
    <div class="cloud-wrapper">
        <div class="cloud-filter clearfix">
            <bk-button class="cloud-btn" theme="primary" @click="handleCreate">{{ $t('新建云同步任务')}}</bk-button>
            <div class="cloud-option-filter clearfix fr">
                <bk-select class="cloud-filter-selector fl"
                    v-model="defaultDemo.selected">
                    <bk-option v-for="(option, index) in selectList"
                        :key="index"
                        :id="option.id"
                        :name="option.name">
                    </bk-option>
                </bk-select>
                <bk-input class="cloud-filter-value cmdb-form-input fl"
                    type="text"
                    :placeholder="$t('任务名称搜索')"
                    v-model.trim="filter.text"
                    @enter="getTableData">
                </bk-input>
                <i class="cloud-filter-search bk-icon icon-search" @click="getTableData"></i>
            </div>
        </div>
        <bk-table class="cloud-discover-table"
            v-bkloading="{ isLoading: $loading('searchCloudTask') }"
            :data="table.list"
            :pagination="table.pagination"
            :max-height="$APP.height - 300"
            @page-limit-change="handleSizeChange"
            @page-change="handlePageChange"
            @sort-change="handleSortChange">
            <bk-table-column type="selection" fixed width="60" align="center" class-name="bk-table-selection"></bk-table-column>
            <bk-table-column prop="bk_task_name" :label="$t('任务名称')"></bk-table-column>
            <bk-table-column prop="bk_account_type" :label="$t('账号类型')">
                <template slot-scope="{ row }">
                    <span v-if="row.bk_account_type === 'tencent_cloud'">{{$t('腾讯云')}}</span>
                </template>
            </bk-table-column>
            <bk-table-column prop="bk_last_sync_time" :label="$t('最近同步时间')">
                <template slot-scope="{ row }">
                    <span v-if="row.bk_last_sync_time === ''">--</span>
                    <span v-else>{{ row.bk_last_sync_time }}</span>
                </template>
            </bk-table-column>
            <bk-table-column prop="bk_last_sync_result" :label="$t('最近同步结果')">
                <template slot-scope="{ row }">
                    <span v-if="row.bk_last_sync_time === ''">--</span>
                    <span v-else>
                        {{$t('新增')}} ({{row.new_add}}) / {{$t('变更update')}} ({{row.attr_changed}})
                    </span>
                </template>
            </bk-table-column>
            <bk-table-column prop="bk_account_admin" :label="$t('任务维护人')"></bk-table-column>
            <bk-table-column prop="bk_sync_status"
                :label="$t('状态')"
                width="100">
                <template slot-scope="{ row }">
                    <template v-if="row.bk_status">
                        <div class="bk-spin-loading bk-spin-loading-mini bk-spin-loading-primary">
                            <div class="rotate rotate1"></div>
                            <div class="rotate rotate2"></div>
                            <div class="rotate rotate3"></div>
                            <div class="rotate rotate4"></div>
                            <div class="rotate rotate5"></div>
                            <div class="rotate rotate6"></div>
                            <div class="rotate rotate7"></div>
                            <div class="rotate rotate8"></div>
                        </div>
                        <span>{{$t('同步中')}}</span>
                    </template>
                    <span class="sync-fail" v-else-if="row.bk_sync_status === 'fail'">
                        {{$t('失败')}}
                    </span>
                    <span v-else>--</span>
                </template>
            </bk-table-column>
            <bk-table-column prop="status"
                :label="$t('是否启用')"
                width="90">
                <template slot-scope="{ row }">
                    <bk-switcher
                        :key="row.bk_task_id"
                        @change="changeStatus(...arguments, row)"
                        :selected="row.bk_status"
                        :is-outline="isOutline"
                        size="small"
                        :show-text="showText">
                    </bk-switcher>
                </template>
            </bk-table-column>
            <bk-table-column prop="bk_task_name"
                :label="$t('操作')"
                width="110">
                <template slot-scope="{ row }">
                    <span class="text-primary mr20" @click.stop="detail(row)">{{$t('详情')}}</span>
                    <span class="text-danger" @click.stop="deleteConfirm(row)">{{$t('删除')}}</span>
                </template>
            </bk-table-column>
            <div class="empty-info" slot="empty">
                <p>{{$t('暂时没有数据，请先')}}
                    <span class="text-primary" @click="handleCreate">{{ $t('新建云同步任务')}}</span>
                </p>
            </div>
        </bk-table>
        <bk-sideslider
            v-transfer-dom
            :is-show.sync="slider.show"
            :title="slider.title"
            :before-close="handleSliderBeforeClose"
            :width="680">
            <template v-if="slider.show" slot="content">
                <v-create v-if="attribute.type === 'create'"
                    ref="detail"
                    :type="attribute.type"
                    @saveSuccess="saveSuccess"
                    @cancel="closeSlider">
                </v-create>
                <bk-tab :active.sync="tab.active" type="unborder-card" v-else>
                    <bk-tab-panel name="details" :label="$t('任务详情')" style="width: calc(100% + 40px);margin: 0 -20px;">
                        <v-update v-if="attribute.type === 'update'"
                            ref="detail"
                            :type="attribute.type"
                            :cur-push="curPush"
                            @saveSuccess="saveSuccess"
                            @cancel="closeSlider">
                        </v-update>
                        <v-task-details v-else-if="attribute.type === 'details'"
                            ref="detail"
                            :type="attribute.type"
                            :cur-push="curPush"
                            @edit="handleEdit"
                            @cancel="closeSlider">
                        </v-task-details>
                    </bk-tab-panel>
                    <bk-tab-panel name="history" :label="$t('同步历史')" :visible="['update', 'details'].includes(attribute.type)">
                        <v-sync-history
                            :cur-push="curPush">
                        </v-sync-history>
                    </bk-tab-panel>
                </bk-tab>
            </template>
        </bk-sideslider>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    import vCreate from './create'
    import vUpdate from './update'
    import vSyncHistory from './sync-history'
    import vTaskDetails from './task-details'
    export default {
        components: {
            vCreate,
            vUpdate,
            vSyncHistory,
            vTaskDetails
        },
        data () {
            return {
                selectList: [{
                    id: 'tencent_cloud',
                    name: this.$t('腾讯云')
                }],
                defaultDemo: {
                    selected: 'tencent_cloud'
                },
                isSelected: false,
                showText: true,
                isOutline: true,
                curPush: {},
                slider: {
                    show: false,
                    title: '',
                    isCloseConfirmShow: false
                },
                attribute: {
                    type: null,
                    inst: {
                        details: {},
                        edit: {}
                    }
                },
                tab: {
                    active: 'attribute'
                },
                filter: {
                    taskId: '',
                    text: ''
                },
                table: {
                    
                    list: [],
                    allList: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        limit: 10
                    },
                    checked: [],
                    defaultSort: '-bk_task_id',
                    sort: '-bk_task_id'
                }
            }
        },
        created () {
            const urlType = this.$route.params.type
            if (urlType) {
                this.handleCreate()
            }
            this.getTableData()
        },
        methods: {
            ...mapActions('cloudDiscover', [
                'searchCloudTask',
                'deleteCloudTask',
                'startCloudSync'
            ]),
            async getTableData () {
                const pagination = this.table.pagination
                const params = {}
                const attr = {}
                const page = {
                    start: (pagination.current - 1) * pagination.limit,
                    limit: pagination.limit,
                    sort: this.table.sort
                }
                if (this.filter.text.length !== 0) {
                    attr['$regex'] = this.filter.text
                    attr['$options'] = '$i'
                    params['bk_task_name'] = attr
                }
                params['page'] = page
                const res = await this.searchCloudTask({ params, config: { requestId: 'searchCloudTask' } })
                this.table.list = res.info.map(data => {
                    data['bk_last_sync_time'] = this.$tools.formatTime(data['bk_last_sync_time'], 'YYYY-MM-DD HH:mm:ss')
                    return data
                })
                pagination.count = res.count
            },
            handleCreate () {
                this.tab.active = 'details'
                this.slider.show = true
                this.slider.title = this.$t('新建云同步任务')
                this.attribute.type = 'create'
                this.attribute.inst.edit = {}
            },
            closeSlider () {
                this.slider.show = false
            },
            saveSuccess () {
                if (this.slider.type === 'create') {
                    this.handlePageChange(1)
                } else {
                    this.getTableData()
                }
                this.slider.show = false
            },
            detail (item) {
                this.tab.active = 'details'
                this.curPush = { ...item }
                this.slider.show = true
                this.attribute.type = 'details'
                this.slider.title = item['bk_task_name']
            },
            deleteConfirm (item) {
                if (item.bk_status) {
                    this.$warn(this.$t('请先停止同步'))
                } else {
                    this.$bkInfo({
                        title: this.$tc('确认删除该任务?'),
                        confirmFn: () => {
                            this.deleteTask(item['bk_task_id'])
                        }
                    })
                }
            },
            async deleteTask (taskID) {
                await this.deleteCloudTask({ taskID })
                this.$success(this.$t('删除任务成功'))
                this.getTableData()
            },
            async changeStatus (status, item) {
                const params = {}
                params['bk_task_name'] = item['bk_task_name']
                params['bk_status'] = status
                params['bk_task_id'] = item['bk_task_id']
                if (status) {
                    this.$success(this.$t('启用成功，约有五分钟延迟'))
                }
                await this.startCloudSync({ params })
                this.getTableData()
            },
            handleSliderBeforeClose () {
                if (['create', 'update'].includes(this.attribute.type) && this.$refs.detail.isCloseConfirmShow()) {
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
            handleEdit (curPush) {
                this.curPush = curPush
                this.attribute.type = 'update'
                this.slider.title = this.$t('编辑') + curPush.bk_task_name
            },
            handleSizeChange (size) {
                this.table.pagination.limit = size
                this.handlePageChange(1)
            },
            handleSortChange (sort) {
                this.table.sort = this.$tools.getSort(sort)
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
    .cloud-wrapper {
        .cloud-filter {
            .cloud-btn {
                float: left;
                margin-right: 10px;
            }
            .cloud-option-filter{
                position: relative;
                margin-right: 10px;
                .cloud-filter-selector{
                    width: 115px;
                    border-radius: 2px 0 0 2px;
                    margin-right: -1px;
                }
                .cloud-filter-value{
                    width: 320px;
                    border-radius: 0 2px 2px 0;
                }
                .cloud-filter-search{
                    position: absolute;
                    right: 10px;
                    top: 11px;
                    cursor: pointer;
                }
            }
        }
        .cloud-discover-table {
            margin-top: 20px;
            span {
                cursor:pointer;
            }
            .sync-fail {
                color: #fc2e2e;
            }
        }
    }
</style>
