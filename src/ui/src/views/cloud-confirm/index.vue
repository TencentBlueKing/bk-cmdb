<template>
    <div>
        <div class="clearfix">
            <bk-button
                type="primary"
                @click="batchConfirm"
                :disabled="!table.checked.length">
                <span>{{$t("Cloud['批量确认']")}}</span>
            </bk-button>
            <div class="confirm-options-button fr">
                <bk-button class="button-history" v-tooltip.bottom="$t('Cloud[\'查看确认记录\']')" @click="confirmHistory">
                    <i class="icon-cc-history"></i>
                </bk-button>
            </div>
            <div class="confirm-options-filter clearfix fr">
                <bk-selector
                    class="confirm-filter-selector fl"
                    :list="selector.list"
                    :selected.sync="selector.defaultDemo.selected">
                </bk-selector>
                <cmdb-form-enum
                    class="confirm-filter-value fl"
                    v-if="filter.type === 'enum'"
                    :options="$tools.getEnumOptions(properties, filter.id)"
                    :allow-clear="true"
                    v-model="filter.value"
                    @on-selected="getTableData">
                </cmdb-form-enum>
                <input
                    class="confirm-filter-value cmdb-form-input fl"
                    type="text"
                    v-else
                    v-model.trim="filter.text"
                    @keydown.enter="getTableData">
                <i class="confirm-filter-search bk-icon icon-search" @click="getTableData"></i>
            </div>
        </div>
        <cmdb-table class="cloud-table" ref="table"
            :loading="$loading('searchConfirm')"
            :header="table.header"
            :checked.sync="table.checked"
            :list="table.list"
            :pagination.sync="table.pagination"
            :defaultSort="table.defaultSort"
            :wrapperMinusHeight="157"
            @handleSizeChange="handleSizeChange"
            @handlePageChange="handlePageChange">
                <template slot="bk_resource_type" slot-scope="{ item }">
                    <span class="attr-changed" v-if="item.bk_resource_type === 'change'">
                        {{$t('Cloud["变更"]')}}
                    </span>
                    <span class="new-add" v-else>
                        {{$t('Cloud["新增"]')}}
                    </span>
                </template>
                <template slot="bk_account_type" slot-scope="{ item }">
                    <span>{{$t('Cloud["腾讯云"]')}}</span>
                </template>
                <template slot="bk_obj_id" slot-scope="{ item }">
                    {{ $t('Hosts["主机"]')}}
                </template>
                <template slot="operation" slot-scope="{ item }">
                    <span class="text-primary mr20" @click.stop="singleConfirm(item)">{{$t('Hosts["确认"]')}}</span>
                </template>
                <div class="empty-info" slot="data-empty">
                    <p>{{$t("Cloud['暂时没有数据，请确保先添加云资源发现任务，']")}}
                        <span class="text-primary" @click="handleAdd">{{ $t('Cloud["去添加"]')}}</span>
                    </p>
                </div>
        </cmdb-table>
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
                        name: this.$t('Cloud["模型"]')
                    }, {
                        id: 1,
                        name: this.$t('Cloud["账号类型"]')
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
                    header: [ {
                        id: 'bk_resource_id',
                        type: 'checkbox'
                    }, {
                        id: 'bk_host_innerip',
                        name: this.$t('Cloud["资源名称"]')
                    }, {
                        id: 'bk_resource_type',
                        name: this.$t('Cloud["资源类型"]')
                    }, {
                        id: 'bk_obj_id',
                        name: this.$t('Cloud["模型"]')
                    }, {
                        id: 'bk_task_name',
                        name: this.$t('Cloud["任务名称"]')
                    }, {
                        id: 'bk_account_type',
                        name: this.$t('Cloud["账号类型"]')
                    }, {
                        id: 'bk_account_admin',
                        sortable: false,
                        name: this.$t('Cloud["任务维护人"]')
                    }, {
                        id: 'create_time',
                        width: 160,
                        name: this.$t('Cloud["发现时间"]')
                    }, {
                        id: 'operation',
                        sortable: false,
                        width: 100,
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
                    defaultSort: '-bk_resource_id',
                    sort: '-bk_resource_id'
                }
            }
        },
        computed: {
            batchConfirmTips () {
                let tips = {}
                tips = this.$t("Cloud['您将批量确认']") + this.table.checked.length + this.$t("Cloud['个资源实例，']") +
                    this.$t("Cloud['确认后的资源实例将被录入到主机资源池中']")
                return tips
            }
        },
        created () {
            this.$store.commit('setHeaderTitle', this.$t('Cloud["资源确认"]'))
            this.getTableData()
        },
        methods: {
            ...mapActions('cloudDiscover', [
                'getResourceConfirm',
                'resourceConfirm',
                'addConfirmHistory'
            ]),
            async getTableData () {
                let pagination = this.table.pagination
                let params = {}
                let attr = {}
                let page = {
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
                let res = await this.getResourceConfirm({params, config: {requestID: 'searchConfirm'}})
                this.table.list = res.info.map(data => {
                    data['create_time'] = this.$tools.formatTime(data['create_time'], 'YYYY-MM-DD HH:mm:ss')
                    return data
                })
                pagination.count = res.count
            },
            confirmHistory () {
                this.$router.push({name: 'resourceConfirmHistory'})
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
                let params = {}
                params['bk_resource_id'] = this.table.checked
                this.$bkInfo({
                    title: this.$t("Cloud['批量资源确认']"),
                    content: this.batchConfirmTips,
                    confirmFn: async () => {
                        await this.handleConfirm(params)
                        this.$success(this.$t('Cloud["资源批量确认确认成功"]'))
                        this.handlePageChange(1)
                    }
                })
            },
            singleConfirm (item) {
                let params = {}
                let arr = []
                arr.push(item['bk_resource_id'])
                params['bk_resource_id'] = arr
                this.$bkInfo({
                    title: this.$t("Cloud['资源确认']"),
                    content: this.$t("Cloud['确认后的资源实例将被录入到主机资源池中']"),
                    confirmFn: async () => {
                        await this.handleConfirm(params)
                        this.$success(this.$t('Cloud["确认成功"]'))
                        this.handlePageChange(1)
                    }
                })
            },
            handleConfirm (params) {
                return Promise.all([
                    this.addConfirmHistory({params}),
                    this.resourceConfirm({params})
                ])
            },
            handleCheckAll () {
                this.table.checked = this.table.list.map(inst => inst['bk_resource_id'])
            },
            handleAdd () {
                this.$router.push({name: 'cloud', params: { type: 'create' }})
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
