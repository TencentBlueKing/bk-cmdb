<template>
    <bk-table class="form-table"
        :data="list">
        <bk-table-column :label="$t('云区域')" prop="bk_cloud_id" width="200" :resizable="false">
            <template slot-scope="{ row }">
                <task-cloud-area-input :disabled="row.bk_cloud_id !== -1"></task-cloud-area-input>
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
                <task-directory-selector class="form-table-selector"
                    v-model="directorySelection[row.bk_vpc_id]"
                    v-validate="'required'"
                    :data-vv-name="`directory-selector-${row.bk_vpc_id}`"
                    :class="{ 'is-error': errors.has(`directory-selector-${row.bk_vpc_id}`) }">
                </task-directory-selector>
            </template>
        </bk-table-column>
    </bk-table>
</template>

<script>
    import TaskFormTableHeader from './task-form-table-header.vue'
    import TaskDirectorySelector from './task-directory-selector.vue'
    import TaskRegionSelector from './task-region-selector.vue'
    import TaskCloudAreaInput from './task-cloud-area-input.vue'
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
                directorySelection: {},
                selection: [],
                directorys: []
            }
        },
        watch: {
            selected: {
                immediate: true,
                handler () {
                    this.updateData()
                }
            }
        },
        methods: {
            updateData () {
                this.list = [...this.selected]
                const newSelection = {}
                this.list.forEach(vpc => {
                    const id = vpc.bk_vpc_id
                    newSelection[id] = this.directorySelection[id] || vpc.bk_sync_dir || 1
                })
                this.directorySelection = newSelection
            },
            vpcFormatter (row, column) {
                const vpcId = row.bk_vpc_id
                const vpcName = row.bk_vpc_name
                if (vpcId !== vpcName) {
                    return `${vpcId}(${vpcName})`
                }
                return vpcId
            },
            handleMultipleSelected (value) {
                Object.keys(this.directorySelection).forEach(id => {
                    this.directorySelection[id] = value
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
            getSyncVPC () {
                return this.list.map(vpc => {
                    return {
                        bk_vpc_id: vpc.bk_vpc_id,
                        bk_vpc_name: vpc.bk_vpc_name,
                        bk_sync_dir: this.directorySelection[vpc.bk_vpc_id],
                        bk_region: vpc.bk_region
                    }
                })
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
    }
</style>
