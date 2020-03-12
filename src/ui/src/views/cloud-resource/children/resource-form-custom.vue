<template>
    <bk-table class="form-table"
        :data="list">
        <bk-table-column label="VPC" prop="vpc" :formatter="vpcFormatter"></bk-table-column>
        <bk-table-column :label="$t('地域')" prop="bk_region_name" show-overflow-tooltip></bk-table-column>
        <bk-table-column :label="$t('主机数量')" prop="bk_host_count"></bk-table-column>
        <bk-table-column :label="$t('主机录入到')" prop="folder" width="250" :render-header="folderHeaderRender">
            <template slot-scope="{ row }">
                <cloud-resource-folder-selector class="form-table-selector"
                    v-model="folderSelection[row.bk_vpc_id]">
                </cloud-resource-folder-selector>
            </template>
        </bk-table-column>
    </bk-table>
</template>

<script>
    import CloudResourceFormCustomHeader from './resource-form-custom-header.vue'
    import CloudResourceFolderSelector from './resource-folder-selector.vue'
    export default {
        name: 'cloud-resource-form-custom',
        components: {
            CloudResourceFolderSelector
        },
        props: {
            selected: Array
        },
        data () {
            return {
                list: [],
                selection: [],
                folderSelection: {},
                folders: [{
                    id: 'a',
                    name: '资源池 / LOL新录入'
                }]
            }
        },
        watch: {
            selected (selected) {
                this.list = [...selected]
            }
        },
        methods: {
            vpcFormatter (row, column) {
                const vpcId = row.bk_vpc_id
                const vpcName = row.bk_vpc_name
                if (vpcId !== vpcName) {
                    return `${vpcId}(${vpcName})`
                }
                return vpcId
            },
            handleMultipleSelected (value) {
                Object.keys(this.folderSelection).forEach(id => {
                    this.folderSelection[id] = value
                })
            },
            folderHeaderRender (h, data) {
                return h('div', [
                    h(CloudResourceFormCustomHeader, {
                        props: {
                            data: data,
                            folders: this.folders,
                            batchSelectHandler: this.handleMultipleSelected,
                            disabled: !this.list.length
                        }
                    })
                ])
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
