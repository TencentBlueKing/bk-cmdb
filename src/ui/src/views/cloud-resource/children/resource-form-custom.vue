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
                <bk-select class="form-table-selector"
                    searchable
                    size="small"
                    :placeholder="$t('请选择xx', { name: $t('资源目录') })"
                    @selected="handleSelected(...arguments, row)">
                    <bk-option v-for="folder in folders"
                        :key="folder.id"
                        :id="folder.id"
                        :name="folder.name">
                    </bk-option>
                    <a href="javascript:void(0)" class="extension-link" slot="extension">
                        <i class="bk-icon icon-plus-circle"></i>
                        {{$t('申请其他目录权限')}}
                    </a>
                </bk-select>
            </template>
        </bk-table-column>
    </bk-table>
</template>

<script>
    import CloudResourceFormCustomHeader from './resource-form-custom-header.vue'
    export default {
        name: 'cloud-resource-form-custom',
        data () {
            return {
                list: [{}, {}, {}, {}],
                selection: [],
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
                this.selection = selection
            },
            handleSelected () {},
            handleMultipleSelected (value) {
                console.log(value)
            },
            folderHeaderRender (h, data) {
                return h('div', [
                    h(CloudResourceFormCustomHeader, {
                        props: {
                            data: data,
                            folders: this.folders,
                            batchSelectHandler: this.handleMultipleSelected
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
