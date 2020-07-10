<template>
    <bk-table class="form-table"
        v-bkloading="{
            isLoading: $loading(request.createArea),
            title: this.$t('正在创建云区域')
        }"
        :row-class-name="getRowClass"
        :data="list">
        <bk-table-column :label="$t('云区域')" prop="bk_cloud_id" width="200" :resizable="false">
            <template slot-scope="{ row }">
                <task-cloud-area-input v-if="row.bk_cloud_id === -1"
                    v-model="row.bk_cloud_name"
                    :key="row.bk_vpc_id"
                    :id="row.bk_cloud_id"
                    :error="row.bk_cloud_error"
                    @input="handleAreaChange(row, ...arguments)">
                </task-cloud-area-input>
                <div class="info-cloud-area" v-else>
                    <task-cloud-area-input display="info" :id="row.bk_cloud_id"></task-cloud-area-input>
                    <span class="info-destroyed"
                        v-if="row.destroyed"
                        v-bk-tooltips="$t('VPC已销毁')">
                        {{$t('已失效')}}
                    </span>
                </div>
            </template>
        </bk-table-column>
        <bk-table-column label="VPC" prop="vpc" :formatter="vpcFormatter" show-overflow-tooltip></bk-table-column>
        <bk-table-column :label="$t('地域')" prop="bk_region_name" show-overflow-tooltip>
            <task-region-selector slot-scope="{ row }"
                display="info"
                :account="account"
                :value="row.bk_region">
            </task-region-selector>
        </bk-table-column>
        <bk-table-column :label="$t('主机数量')" prop="bk_host_count"></bk-table-column>
        <bk-table-column :label="$t('主机录入到')" prop="directory" width="200" :render-header="directoryHeaderRender" :resizable="false">
            <template slot-scope="{ row }">
                <task-directory-selector display="info" :value="row.bk_sync_dir" v-if="row.destroyed"></task-directory-selector>
                <task-directory-selector class="form-table-selector"
                    v-else
                    v-model="row.bk_sync_dir"
                    v-validate="'required'"
                    :key="row.bk_vpc_id"
                    :data-vv-name="`directory-selector-${row.bk_vpc_id}`"
                    :class="{ 'is-error': errors.has(`directory-selector-${row.bk_vpc_id}`) }">
                </task-directory-selector>
            </template>
        </bk-table-column>
        <bk-table-column :label="$t('操作')" width="80" :resizable="false">
            <template slot-scope="{ row }">
                <bk-button text @click="handleRemove(row)">{{$t('删除')}}</bk-button>
            </template>
        </bk-table-column>
        <template slot="empty">{{$t('请添加VPC')}}</template>
    </bk-table>
</template>

<script>
    import TaskFormTableHeader from './task-form-table-header.vue'
    import TaskDirectorySelector from './task-directory-selector.vue'
    import TaskRegionSelector from './task-region-selector.vue'
    import TaskCloudAreaInput from './task-cloud-area-input.vue'
    import Bus from '@/utils/bus'
    import symbols from '../common/symbol'
    export default {
        name: 'task-form-table',
        components: {
            TaskDirectorySelector,
            TaskRegionSelector,
            TaskCloudAreaInput
        },
        props: {
            selected: Array,
            account: Number
        },
        data () {
            return {
                list: [],
                selection: [],
                directorys: [],
                request: {
                    createArea: symbols.get('createArea')
                }
            }
        },
        watch: {
            selected: {
                immediate: true,
                handler () {
                    this.updateList()
                }
            }
        },
        methods: {
            updateList () {
                const oldList = this.list
                this.list = this.selected.map(vpc => {
                    const newRow = { ...vpc }
                    const oldRow = oldList.find(row => row.bk_vpc_id === vpc.bk_vpc_id)
                    newRow.bk_sync_dir = oldRow ? oldRow.bk_sync_dir : vpc.bk_sync_dir
                    newRow.bk_cloud_id = oldRow ? oldRow.bk_cloud_id : vpc.bk_cloud_id
                    // bk_cloud_name、bk_cloud_error用于创建新的云区域
                    newRow.bk_cloud_name = oldRow ? oldRow.bk_cloud_name : ''
                    newRow.bk_cloud_error = oldRow ? oldRow.bk_cloud_error : false
                    return newRow
                })
            },
            vpcFormatter (row, column) {
                const vpcId = row.bk_vpc_id
                const vpcName = row.bk_vpc_name
                if (vpcId !== vpcName) {
                    return `${vpcId}(${vpcName})`
                }
                return vpcId
            },
            handleRemove (row) {
                this.$emit('remove', row)
            },
            handleMultipleSelected (value) {
                this.list.forEach(row => {
                    row.bk_sync_dir = value
                })
            },
            directoryHeaderRender (h, data) {
                return h('div', [
                    h(TaskFormTableHeader, {
                        props: {
                            data: data,
                            batchSelectHandler: this.handleMultipleSelected,
                            disabled: !this.list.length
                        }
                    })
                ])
            },
            handleAreaChange (row, cloudName) {
                row.bk_cloud_name = cloudName
                row.bk_cloud_error = false
            },
            // task-form保存时调用，创建新的云区域
            async createCloudArea (accountInfo) {
                try {
                    const newAreaList = this.list.filter(row => row.bk_cloud_id === -1)
                    if (!newAreaList.length) {
                        return Promise.resolve(true)
                    }
                    const valid = this.validate(newAreaList)
                    if (!valid) {
                        return Promise.resolve(false)
                    }
                    
                    const results = await this.$store.dispatch('cloud/area/batchCreate', {
                        params: {
                            data: newAreaList.map(row => {
                                return {
                                    bk_cloud_name: row.bk_cloud_name,
                                    bk_vpc_id: row.bk_vpc_id,
                                    bk_vpc_name: row.bk_vpc_name,
                                    bk_region: row.bk_region,
                                    bk_cloud_vendor: accountInfo.bk_cloud_vendor,
                                    bk_account_id: accountInfo.bk_account_id
                                }
                            })
                        },
                        config: {
                            requestId: this.request.createArea
                        }
                    })
                    
                    let hasError = false
                    results.forEach((result, index) => {
                        if (result.bk_cloud_id === -1) {
                            newAreaList[index].bk_cloud_error = result.err_msg
                            hasError = true
                        } else {
                            newAreaList[index].bk_cloud_id = result.bk_cloud_id
                            newAreaList[index].bk_cloud_error = false
                        }
                    })
                    Bus.$emit('refresh-cloud-area')
                    return Promise.resolve(!hasError)
                } catch (error) {
                    console.error(error)
                    return Promise.resolve(false)
                }
            },
            validate (list) {
                let valid = true
                list.forEach(row => {
                    if (!row.bk_cloud_name) {
                        row.bk_cloud_error = this.$t('请填写云区域')
                        valid = false
                    }
                })
                return valid
            },
            getRowClass ({ row }) {
                if (row.destroyed) {
                    return 'is-destroyed'
                }
                return ''
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-table {
        .form-table-selector {
            width: 100%;
            &.is-error {
                border-color: $dangerColor;
            }
        }
        /deep/ {
            .bk-table-row.is-destroyed {
                color: #C4C6CC;
            }
        }
        .info-cloud-area {
            display: flex;
            align-items: center;
            justify-content: flex-start;
            white-space: nowrap;
            .info-destroyed {
                margin-left: 4px;
                font-size: 12px;
                line-height: 18px;
                color: #EA3536;
                padding: 0 4px;
                border-radius: 2px;
                background-color: #FEDDDC;
            }
        }
    }
</style>
