<template>
    <div class="process-wrapper">
        <div class="process-filter clearfix">
            <bk-button class="process-btn" 
            type="primary"
            :disabled="!table.checked.length"
            @click="handleConfirm">
            <span>批量确认</span>
        </bk-button>
            <div class="filter-text fr">
                <input type="text" class="bk-form-input"
                    v-model.trim="filter.text" @keyup.enter="handlePageChange(1)">
                    <i class="bk-icon icon-search" @click="handlePageChange(1)"></i>
            </div>
            <div class="filter-text fr">
                <bk-selector style="width: 100px;"
                    :list="selector.list"
                    :selected.sync="selector.defaultDemo.selected"
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
            @handlePageChange="handlePageChange"
            @handleCheckAll="handleCheckAll">
                <template slot="operation" slot-scope="{ item }">
                    <span class="text-primary mr20" @click.stop="detail(item)">{{$t('Common["详情"]')}}</span>
                </template>
        </cmdb-table>
        <cmdb-slider :isShow.sync="slider.show" :title="slider.title" :width="560">
            <v-details
                ref="detail"
                slot="content"
                :curPush="curPush"
                @saveSuccess="saveSuccess"
                @cancel="closeSlider">
            </v-details>
        </cmdb-slider>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    import vDetails from './details'
    export default {
        components: {
            vDetails
        },
        data () {
            return {
                selector: {
                    list: [{
                        id: 1,
                        name: '模型'
                    }, {
                        id: 2,
                        name: '来源类型'
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
                        id: 'bk_obj_id',
                        name: '模型'
                    }, {
                        id: 'bk_host_innerip',
                        name: '资源名称'
                    }, {
                        id: 'bk_source_type',
                        name: '来源类型'
                    }, {
                        id: 'bk_source_name',
                        name: '来源名称'
                    }, {
                        id: 'create_time',
                        name: '发现时间'
                    }, {
                        id: 'bk_in_charge',
                        name: '负责人'
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
                    defaultSort: '-bk_resource_id',
                    sort: '-bk_resource_id'
                }
            }
        },
        methods: {
            ...mapActions('cloudDiscover', [
                'searchCloudTask',
                'getResourceConfirm',
                'resourceConfirm'
            ]),
            async getTableData () {
                let pagination = this.table.pagination
                let params = {}
                if (this.selector.checked === 'bk_obj_id') {
                    params['bk_obj_id'] = this.filter.text
                }
                if (this.selector.checked === 'bk_source_type') {
                    params['bk_source_type'] = this.filter.text
                }
                let res = await this.getResourceConfirm({params})
                this.table.list = res.info.map(data => {
                    data['create_time'] = this.$tools.formatTime(data['create_time'], 'YYYY-MM-DD HH:mm:ss')
                    if (data['bk_obj_id'] === 'host') {
                        data['bk_obj_id'] = '主机'
                    } else {
                        data['bk_obj_id'] = '交换机'
                    }
                    return data
                })
                pagination.count = res.count
            },
            handleSizeChange (size) {
                this.table.pagination.size = size
                this.handlePageChange(1)
            },
            handlePageChange (page) {
                this.table.pagination.current = page
                this.getTableData()
            },
            handleConfirm () {
                let params = {}
                params['bk_resource_id'] = this.table.checked
                this.resourceConfirm({params})
                this.getTableData()
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
            async detail (item) {
                this.slider.show = true
                this.slider.title = '查看同步任务详情'
                let params = {}
                params['bk_task_id'] = item['bk_task_id']
                let res = await this.searchCloudTask({params})
                this.curPush = res.info.map(data => {
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
                })[0]
            },
            selected (id) {
                if (id === 1) {
                    this.selector.checked = 'bk_obj_id'
                } else {
                    this.selector.checked = 'bk_source_type'
                }
            },
            handleCheckAll () {
                this.table.checked = this.table.list.map(inst => inst['bk_resource_id'])
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
