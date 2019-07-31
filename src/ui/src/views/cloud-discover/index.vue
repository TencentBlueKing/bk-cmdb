<template>
    <div class="cloud-wrapper">
        <div class="cloud-filter clearfix">
            <bk-button class="cloud-btn" type="primary" @click="handleCreate">{{ $t('Cloud["新建云同步任务"]')}}</bk-button>
            <div class="cloud-option-filter clearfix fr">
                <bk-selector class="cloud-filter-selector fl"
                    :list="selectList"
                    :selected.sync="defaultDemo.selected">
                </bk-selector>
                <input class="cloud-filter-value cmdb-form-input fl"
                    type="text"
                    :placeholder="$t('Cloud[\'任务名称搜索\']')"
                    v-model.trim="filter.text"
                    @keyup.enter="getTableData">
                <i class="cloud-filter-search bk-icon icon-search" @click="getTableData"></i>
            </div>
        </div>
        <cmdb-table class="cloud-discover-table" ref="table"
            :loading="$loading('searchCloudTask')"
            :checked.sync="table.checked"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :default-sort="table.defaultSort"
            :wrapper-minus-height="300"
            @handleSizeChange="handleSizeChange"
            @handlePageChange="handlePageChange"
            @handleSortChange="handleSortChange">
            <template slot="bk_sync_status" slot-scope="{ item }">
                <template v-if="item.bk_status">
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
                    <span>{{$t('Cloud["同步中"]')}}</span>
                </template>
                <span class="sync-fail" v-else-if="item.bk_sync_status === 'fail'">
                    {{$t('EventPush["失败"]')}}
                </span>
                <span v-else>--</span>
            </template>
            <template slot="status" slot-scope="{ item }">
                <bk-switcher
                    :key="item.bk_task_id"
                    @change="changeStatus(...arguments, item)"
                    :selected="item.bk_status"
                    :is-outline="isOutline"
                    size="small"
                    :show-text="showText">
                </bk-switcher>
            </template>
            <template slot="bk_account_type" slot-scope="{ item }">
                <span v-if="item.bk_account_type === 'tencent_cloud'">{{$t('Cloud[\'腾讯云\']')}}</span>
            </template>
            <template slot="bk_last_sync_time" slot-scope="{ item }">
                <span v-if="item.bk_last_sync_time === ''">--</span>
                <span v-else>{{ item.bk_last_sync_time }}</span>
            </template>
            <template slot="bk_last_sync_result" slot-scope="{ item }">
                <span v-if="item.bk_last_sync_time === ''">--</span>
                <span v-else>
                    {{$t('Cloud[\'新增\']')}} ({{item.new_add}}) / {{$t('Cloud[\'变更\']')}} ({{item.attr_changed}})
                </span>
            </template>
            <template slot="operation" slot-scope="{ item }">
                <span class="text-primary mr20" @click.stop="detail(item)">{{$t('Cloud["详情"]')}}</span>
                <span class="text-danger" @click.stop="deleteConfirm(item)">{{$t('Common["删除"]')}}</span>
            </template>
            <div class="empty-info" slot="data-empty">
                <p>{{$t("Cloud['暂时没有数据，请先']")}}
                    <span class="text-primary" @click="handleCreate">{{ $t('Cloud["新建云同步任务"]')}}</span>
                </p>
            </div>
        </cmdb-table>
        <cmdb-slider
            :is-show.sync="slider.show"
            :title="slider.title"
            :before-close="handleSliderBeforeClose"
            :width="680">
            <v-create v-if="attribute.type === 'create'"
                slot="content"
                ref="detail"
                :type="attribute.type"
                @saveSuccess="saveSuccess"
                @cancel="closeSlider">
            </v-create>
            <bk-tab :active-name.sync="tab.active" slot="content" v-else>
                <bk-tabpanel name="details" :title="$t('Cloud[\'任务详情\']')" style="width: calc(100% + 40px);margin: 0 -20px;">
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
                </bk-tabpanel>
                <bk-tabpanel name="history" :title="$t('Cloud[\'同步历史\']')" :show="['update', 'details'].includes(attribute.type)">
                    <v-sync-history
                        :cur-push="curPush">
                    </v-sync-history>
                </bk-tabpanel>
            </bk-tab>
        </cmdb-slider>
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
                    name: this.$t('Cloud["腾讯云"]')
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
                    header: [{
                        id: 'bk_task_name',
                        name: this.$t('Cloud["任务名称"]')
                    }, {
                        id: 'bk_account_type',
                        name: this.$t('Cloud["账号类型"]')
                    }, {
                        id: 'bk_last_sync_time',
                        name: this.$t('Cloud["最近同步时间"]')
                    }, {
                        id: 'bk_last_sync_result',
                        sortable: false,
                        name: this.$t('Cloud["最近同步结果"]')

                    }, {
                        id: 'bk_account_admin',
                        sortable: false,
                        name: this.$t('Cloud["任务维护人"]')
                    }, {
                        id: 'bk_sync_status',
                        sortable: false,
                        width: 100,
                        name: this.$t('ProcessManagement["状态"]')
                    }, {
                        id: 'status',
                        sortable: false,
                        width: 90,
                        name: this.$t('Cloud["是否启用"]')
                    }, {
                        id: 'operation',
                        sortable: false,
                        width: 110,
                        name: this.$t('Common["操作"]')
                    }],
                    list: [],
                    allList: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        size: 10
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
                    start: (pagination.current - 1) * pagination.size,
                    limit: pagination.size,
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
                this.slider.title = this.$t('Cloud["新建云同步任务"]')
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
                    this.$warn(this.$t('Cloud["请先停止同步"]'))
                } else {
                    this.$bkInfo({
                        title: this.$tc('Cloud["确认删除该任务?"]'),
                        confirmFn: () => {
                            this.deleteTask(item['bk_task_id'])
                        }
                    })
                }
            },
            async deleteTask (taskID) {
                await this.deleteCloudTask({ taskID })
                this.$success(this.$t('Cloud["删除任务成功"]'))
                this.getTableData()
            },
            async changeStatus (status, item) {
                const params = {}
                params['bk_task_name'] = item['bk_task_name']
                params['bk_status'] = status
                params['bk_task_id'] = item['bk_task_id']
                if (status) {
                    this.$success(this.$t('Cloud["启用成功，约有五分钟延迟"]'))
                }
                await this.startCloudSync({ params })
                this.getTableData()
            },
            handleSliderBeforeClose () {
                if (['create', 'update'].includes(this.attribute.type) && this.$refs.detail.isCloseConfirmShow()) {
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
            handleEdit (curPush) {
                this.curPush = curPush
                this.attribute.type = 'update'
                this.slider.title = this.$t('Cloud["编辑"]') + curPush.bk_task_name
            },
            handleSizeChange (size) {
                this.table.pagination.size = size
                this.handlePageChange(1)
            },
            handleSortChange (sort) {
                this.table.sort = sort
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
