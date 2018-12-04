<template>
    <div class="process-wrapper">
        <div class="process-filter clearfix">
            <bk-button class="process-btn" type="primary" @click="handleCreate">新建云同步</bk-button>
            <div class="filter-text fr">
                <input type="text" class="bk-form-input" placeholder="任务名称搜索"
                    v-model.trim="filter.text" @keyup.enter="getTableData">
                    <i class="bk-icon icon-search" @click="getTableData"></i>
            </div>
            <div class="filter-text fr">
                <bk-selector style="width: 100px;"
                    :list="list"
                    :selected.sync="defaultDemo.selected"
                    @item-selected="selected">
                </bk-selector>
            </div>
        </div>
        <cmdb-table class="process-table" ref="table"
            :loading="$loading('post_searchProcess_list')"
            :checked.sync="table.checked"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :defaultSort="table.defaultSort"
            :wrapperMinusHeight="300"
            @handleSizeChange="handleSizeChange"
            @handlePageChange="handlePageChange">
                <template slot="status" slot-scope="{ item }">
                    <bk-switcher
                        @change="changeStatus(...arguments, item)"
                        :selected="item.bk_status"
                        :is-Outline="isOutline"
                        on-text="开启"
                        off-text="关闭"
                        :show-text="showText">
                    </bk-switcher>
                </template>
                <template slot="operation" slot-scope="{ item }">
                    <span class="text-primary mr20" @click.stop="detail(item)">{{$t('Common["详情"]')}}</span>
                    <span class="text-danger" @click.stop="deleteConfirm(item)">{{$t('Common["删除"]')}}</span>
                </template>
        </cmdb-table>
        <cmdb-slider :isShow.sync="slider.show" :title="slider.title" :width="720">
            <bk-tab :active-name.sync="tab.active" slot="content">
                <bk-tabpanel name="attribute" :title="$t('Common[\'任务详情\']')" style="width: calc(100% + 40px);margin: 0 -20px;">
                    <v-create v-if="attribute.type === 'create'"
                        ref="detail"
                        :type="attribute.type"
                        :curPush="curPush"
                        @saveSuccess="saveSuccess"
                        @cancel="closeSlider">
                    </v-create>
                    <v-update v-else-if="attribute.type === 'update'"
                        ref="detail"
                        :type="attribute.type"
                        :curPush="curPush"
                        @saveSuccess="saveSuccess"
                        @cancel="closeSlider">   
                    </v-update>
                    <v-task-details v-else-if="attribute.type === 'detail'"
                        ref="detail"
                        :type="attribute.type"
                        :curPush="curPush"
                        @edit="handleEdit"
                        @cancel="closeSlider">
                    </v-task-details>
                </bk-tabpanel>
                <bk-tabpanel name="history" title="同步历史" :show="attribute.type === 'detail'">
                    <v-sync-history
                    :curPush="curPush">
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
                list: [{
                    id: 1,
                    name: '腾讯云'
                }],
                defaultDemo: {
                    selected: 1
                },
                isSelected: false,
                showText: true,
                isOutline: true,
                curPush: {},
                slider: {
                    show: false,
                    title: ''
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
                    header: [ {
                        id: 'bk_task_name',
                        name: '名称'
                    }, {
                        id: 'bk_account_type',
                        name: '账号类型'
                    }, {
                        id: 'bk_account_admin',
                        name: '管理员'
                    }, {
                        id: 'bk_last_sync_time',
                        name: '最近同步时间'
                    }, {
                        id: 'status',
                        name: '是否启用'
                    }, {
                        id: 'operation',
                        name: '操作'
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
        methods: {
            ...mapActions('cloudDiscover', [
                'searchCloudTask', 'deleteCloudTask',
                'startCloudSync'
            ]),
            async getTableData () {
                let pagination = this.table.pagination
                let params = {}
                if (this.filter.text.length !== 0) {
                    params['bk_task_name'] = this.filter.text
                }
                let res = await this.searchCloudTask({params})
                this.table.list = res.info.map(data => {
                    data['bk_last_sync_time'] = this.$tools.formatTime(data['bk_last_sync_time'], 'YYYY-MM-DD HH:mm:ss')
                    if (data['bk_obj_id'] === 'host') {
                        data['bk_obj_id'] = '主机'
                    } else {
                        data['bk_obj_id'] = '交换机'
                    }
                    if (data['bk_period_type'] === 'day') {
                        data['bk_period_type'] = '每天'
                    } else if (data['bk_period_type'] === 'hour') {
                        data['bk_period_type'] = '每小时'
                    } else {
                        data['bk_period_type'] = '每五分钟'
                    }
                    return data
                })
                pagination.count = res.count
            },
            handleSizeChange (size) {
                this.table.pagination.size = size
                this.handlePageChange(1)
            },
            handleCreate () {
                this.slider.show = true
                this.slider.title = '新增云同步任务'
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
                this.curPush = {...item}
                this.slider.show = true
                this.attribute.type = 'detail'
                this.slider.title = '查看同步任务详情'
            },
            deleteConfirm (item) {
                this.$bkInfo({
                    title: this.$tc('EventPush["确认删除该任务"]'),
                    confirmFn: () => {
                        this.deletePush(item['bk_task_id'])
                    }
                })
            },
            async deletePush (taskID) {
                await this.deleteCloudTask({ taskID })
                this.$success(this.$t('EventPush["删除任务成功"]'))
                this.getTableData()
            },
            changeStatus (status, item) {
                let params = {}
                params['bk_task_name'] = item['bk_task_name']
                params['bk_status'] = status
                params['bk_task_id'] = item['bk_task_id']
                this.startCloudSync({params})
            },
            selected (id, data) {
            },
            handleEdit (curPush) {
                this.curPush = curPush
                this.attribute.type = 'update'
                this.slider.title = '修改云同步任务'
            },
            handlePageChange (page) {
                this.table.pagination.current = page
                this.getTableData()
            }
        },
        created () {
            this.getTableData()
        }
    }
</script>

<style lang="scss" scoped>
    .process-wrapper {
        .process-filter {
            .process-btn {
                float: left;
                margin-right: 10px;
            }
        }
        .filter-text{
            position: relative;
            .bk-form-input{
                width: 320px;
            }
            .icon-search{
                position: absolute;
                right: 10px;
                top: 10px;
                cursor: pointer;
            }
        }
        .process-table {
            margin-top: 20px;
        }
    }
</style>
