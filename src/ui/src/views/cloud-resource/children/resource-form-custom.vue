<template>
    <bk-table class="form-table"
        :data="list"
        @selection-change="handleSelectionChange">
        <bk-table-column type="selection"></bk-table-column>
        <bk-table-column label="VPC" prop="vpc"></bk-table-column>
        <bk-table-column :label="$t('地域')" prop="location"></bk-table-column>
        <bk-table-column :label="$t('主机数量')" prop="host_count"></bk-table-column>
        <bk-table-column :label="$t('主机录入到')" prop="folder" width="250" :render-header="folderHeaderRender">
            <template slot-scope="{ row }" v-if="isSelected(row)">
                <cloud-resource-folder-selector class="form-table-selector"
                    v-model="folderSelection[row.id]">
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
        data () {
            return {
                list: ['', '', '', ''].map((_, index) => ({ id: index, vpc: 'vpc' + index })),
                selection: [],
                folderSelection: {},
                folders: [{
                    id: 'a',
                    name: '资源池 / LOL新录入'
                }]
            }
        },
        methods: {
            isSelected (row) {
                return this.selection.includes(row)
            },
            handleSelectionChange (selection) {
                const folderSelection = {}
                selection.forEach(row => {
                    if (this.folderSelection.hasOwnProperty(row.id)) {
                        folderSelection[row.id] = this.folderSelection[row.id]
                    } else {
                        folderSelection[row.id] = ''
                    }
                })
                this.folderSelection = folderSelection
                this.selection = selection
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
                            disabled: !this.selection.length
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
