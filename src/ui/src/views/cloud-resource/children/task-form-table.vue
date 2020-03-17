<template>
    <bk-table class="form-table"
        :data="list">
        <bk-table-column label="VPC" prop="vpc" :formatter="vpcFormatter"></bk-table-column>
        <bk-table-column :label="$t('地域')" prop="bk_region_name" show-overflow-tooltip>
            <template slot-scope="{ row }">{{getRegionName(row)}}</template>
        </bk-table-column>
        <bk-table-column :label="$t('主机数量')" prop="bk_host_count"></bk-table-column>
        <bk-table-column :label="$t('主机录入到')" prop="directory" width="250" :render-header="directoryHeaderRender">
            <template slot-scope="{ row }">
                <task-directory-selector class="form-table-selector"
                    v-model="directorySelection[row.bk_vpc_id]"
                    :batch-select-handler="handleDirectoryBatchSelect">
                </task-directory-selector>
            </template>
        </bk-table-column>
    </bk-table>
</template>

<script>
    import TaskFormTableHeader from './task-form-table-header.vue'
    import TaskDirectorySelector from './task-directory-selector.vue'
    export default {
        name: 'task-form-table',
        components: {
            TaskDirectorySelector
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
                directorys: [],
                request: {
                    region: Symbol('region')
                },
                regions: []
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
        created () {
            this.getRegions()
        },
        methods: {
            async getRegions () {
                try {
                    this.regions = await this.$store.dispatch('cloud/resource/findRegion', {
                        params: {
                            bk_account_id: this.account,
                            with_host_account: false
                        },
                        config: {
                            requestId: this.request.region
                        }
                    })
                } catch (e) {
                    console.error(e)
                }
            },
            updateData () {
                this.list = [...this.selected]
                const vpcIds = this.list.map(vpc => vpc.bk_vpc_id)
                const newSelection = {}
                vpcIds.forEach(id => {
                    newSelection[id] = this.directorySelection[id] || ''
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
            handleDirectoryBatchSelect (id) {
                for (const key in this.directorySelection) {
                    this.directorySelection[key] = id
                }
            },
            getSyncVPC () {
                return this.list.map(vpc => {
                    return {
                        bk_vpc_id: vpc.bk_vpc_id,
                        bk_sync_dir: this.directorySelection[vpc.bk_vpc_id]
                    }
                })
            },
            getRegionName (vpc) {
                const region = this.regions.find(region => region.bk_region === vpc.bk_region)
                return region ? region.bk_region_name : '--'
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-table {
        .form-table-selector {
            width: 100%;
        }
    }
</style>
