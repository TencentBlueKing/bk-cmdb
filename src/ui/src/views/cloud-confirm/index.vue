<template>
    <div>
        <div class="clearfix">
            <bk-button
                theme="primary"
                @click="batchConfirm"
                :disabled="!table.checked.length">
                <span>{{$t('批量确认')}}</span>
            </bk-button>
            <div class="confirm-options-button fr">
                <bk-button class="button-history" v-bk-tooltips.bottom="$t('查看确认记录')" @click="confirmHistory">
                    <i class="icon-cc-history"></i>
                </bk-button>
            </div>
            <div class="confirm-options-filter clearfix fr">
                <bk-select class="confirm-filter-selector fl"
                    v-model="selector.defaultDemo.selected">
                    <bk-option v-for="(option, index) in selector.list"
                        :key="index"
                        :id="option.id"
                        :name="option.name">
                    </bk-option>
                </bk-select>
                <cmdb-form-enum
                    class="confirm-filter-value fl"
                    v-if="filter.type === 'enum'"
                    :options="$tools.getEnumOptions(properties, filter.id)"
                    :allow-clear="true"
                    v-model="filter.value"
                    @on-selected="getTableData">
                </cmdb-form-enum>
                <bk-input
                    class="confirm-filter-value cmdb-form-input fl"
                    type="text"
                    clearable
                    v-else
                    v-model.trim="filter.text"
                    @enter="getTableData">
                </bk-input>
                <i class="confirm-filter-search bk-icon icon-search" @click="getTableData"></i>
            </div>
        </div>
        <bk-table class="cloud-table"
            v-bkloading="{ isLoading: $loading('searchConfirm') }"
            :data="table.list"
            :pagination="table.pagination"
            :max-height="$APP.height - 160"
            @page-limit-change="handleSizeChange"
            @page-change="handlePageChange"
            @selection-change="handleSelectChange">
            <bk-table-column type="selection" fixed width="60" align="center" class-name="bk-table-selection"></bk-table-column>
            <bk-table-column prop="bk_host_innerip" :label="$t('资源名称')"></bk-table-column>
            <bk-table-column prop="bk_resource_type" :label="$t('资源类型')">
                <template slot-scope="{ row }">
                    <span class="attr-changed" v-if="row.bk_resource_type === 'change'">
                        {{$t('变更')}}
                    </span>
                    <span class="new-add" v-else>
                        {{$t('新增')}}
                    </span>
                </template>
            </bk-table-column>
            <bk-table-column prop="bk_obj_id" :label="$t('模型')">
                <template>
                    {{ $t('主机')}}
                </template>
            </bk-table-column>
            <bk-table-column prop="bk_task_name" :label="$t('任务名称')"></bk-table-column>
            <bk-table-column prop="bk_account_type" :label="$t('账号类型')">
                <template slot-scope="{ row }">
                    <span v-if="row.bk_account_type === 'tencent_cloud'">{{$t('腾讯云')}}</span>
                    <span v-else>{{row.bk_account_type}}</span>
                </template>
            </bk-table-column>
            <bk-table-column prop="bk_account_admin" :label="$t('任务维护人')"></bk-table-column>
            <bk-table-column prop="create_time" :label="$t('发现时间')"></bk-table-column>
            <bk-table-column :label="$t('操作')" fixed="right">
                <template slot-scope="{ row }">
                    <span class="text-primary mr20" @click.stop="singleConfirm(row)">{{$t('确认')}}</span>
                </template>
            </bk-table-column>
            <div class="empty-info" slot="empty">
                <p>{{$t('暂时没有数据，请确保先添加云资源发现任务，')}}
                    <span class="text-primary" @click="handleAdd">{{ $t('去添加')}}</span>
                </p>
            </div>
        </bk-table>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    export default {
        data () {
            return {
                selector: {
                    list: [{
                        id: 0,
                        name: this.$t('模型')
                    }, {
                        id: 1,
                        name: this.$t('账号类型')
                    }],
                    defaultDemo: {
                        selected: 1
                    },
                    checked: ''
                },
                isSelected: false,
                showText: true,
                isOutline: true,
                curPush: {},
                slider: {
                    show: false,
                    title: '',
                    type: 'create'
                },
                tab: {
                    active: 'attribute'
                },
                filter: {
                    bizId: '',
                    text: '',
                    businessResolver: null
                },
                table: {
                    list: [],
                    allList: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        size: 10
                    },
                    checked: [],
                    defaultSort: '-bk_resource_id',
                    sort: '-bk_resource_id'
                }
            }
        },
        computed: {
            batchConfirmTips () {
                let tips = {}
                tips = this.$t('您将批量确认') + this.table.checked.length + this.$t('个资源实例，')
                    + this.$t('确认后的资源实例将被录入到主机资源池中')
                return tips
            }
        },
        created () {
            this.getTableData()
        },
        methods: {
            ...mapActions('cloudDiscover', [
                'getResourceConfirm',
                'resourceConfirm',
                'addConfirmHistory'
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
                if (this.filter.text !== '') {
                    attr['$regex'] = this.filter.text
                    attr['$options'] = '$i'
                    if (this.selector.defaultDemo.selected === 0) {
                        params['bk_obj_id'] = 'host'
                    }
                    if (this.selector.defaultDemo.selected === 1) {
                        params['bk_account_type'] = 'tencent_cloud'
                    }
                }
                params['page'] = page
                const res = await this.getResourceConfirm({ params, config: { requestID: 'searchConfirm' } })
                this.table.list = res.info.map(data => {
                    data['create_time'] = this.$tools.formatTime(data['create_time'], 'YYYY-MM-DD HH:mm:ss')
                    return data
                })
                pagination.count = res.count
            },
            confirmHistory () {
                this.$router.push({
                    name: 'confirmHistory'
                })
            },
            handleSelectChange (selection) {
                this.table.checked = selection.map(row => row.bk_resource_id)
            },
            handleSizeChange (size) {
                this.table.pagination.size = size
                this.handlePageChange(1)
            },
            handlePageChange (page) {
                this.table.pagination.current = page
                this.getTableData()
            },
            batchConfirm () {
                const params = {}
                params['bk_resource_id'] = this.table.checked
                this.$bkInfo({
                    title: this.$t('批量资源确认'),
                    subTitle: this.batchConfirmTips,
                    confirmFn: async () => {
                        await this.handleConfirm(params)
                        this.$success(this.$t('资源批量确认确认成功'))
                        this.handlePageChange(1)
                    }
                })
            },
            singleConfirm (item) {
                const params = {}
                const arr = []
                arr.push(item['bk_resource_id'])
                params['bk_resource_id'] = arr
                this.$bkInfo({
                    title: this.$t('资源确认'),
                    subTitle: this.$t('确认后的资源实例将被录入到主机资源池中'),
                    confirmFn: async () => {
                        await this.handleConfirm(params)
                        this.$success(this.$t('确认成功'))
                        this.handlePageChange(1)
                    }
                })
            },
            handleConfirm (params) {
                return Promise.all([
                    this.addConfirmHistory({ params }),
                    this.resourceConfirm({ params })
                ])
            },
            handleCheckAll () {
                this.table.checked = this.table.list.map(inst => inst['bk_resource_id'])
            },
            handleAdd () {
                this.$router.push({
                    name: 'cloudDiscover',
                    params: {
                        type: 'create'
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .confirm-options-button{
        font-size: 0;
        .button-history{
            border-radius: 2px 0 0 2px;
        }
    }
    .confirm-options-filter{
        position: relative;
        margin-right: 10px;
        .confirm-filter-selector{
            width: 115px;
            border-radius: 2px 0 0 2px;
            margin-right: -1px;
        }
        .confirm-filter-value{
            width: 320px;
            border-radius: 0 2px 2px 0;
        }
        .confirm-filter-search{
            position: absolute;
            right: 10px;
            top: 11px;
            cursor: pointer;
        }
    }
    .cloud-table{
        margin-top: 20px;
        span {
            cursor:pointer;
        }
        .attr-changed {
            color: #ffb23a;
        }
        .new-add {
            color: #4f55f3;
        }
    }
</style>
