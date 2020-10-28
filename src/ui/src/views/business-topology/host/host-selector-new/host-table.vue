<template>
    <div class="host-table">
        <div class="search-bar">
            <bk-input
                clearable
                right-icon="icon-search"
                v-model.trim="keyword">
            </bk-input>
        </div>
        <bk-table
            ref="table"
            :data="displayList"
            :max-height="410"
            :outer-border="false"
            :header-border="false"
            @select="handleSelect"
            @select-all="handleSelectAll">
            <bk-table-column type="selection" width="30"></bk-table-column>
            <bk-table-column :label="$t('内网IP')">
                <template slot-scope="{ row }">
                    {{row.host.bk_host_innerip}}
                </template>
            </bk-table-column>
            <bk-table-column :label="$t('云区域')" show-overflow-tooltip>
                <template slot-scope="{ row }">{{row.host.bk_cloud_id | foreignkey}}</template>
            </bk-table-column>
        </bk-table>
    </div>
</template>

<script>
    import debounce from 'lodash.debounce'
    import { foreignkey } from '@/filters/formatter.js'
    export default {
        filters: {
            foreignkey
        },
        props: {
            list: {
                type: Array,
                default: () => ([])
            },
            selected: {
                type: Array,
                default: () => ([])
            }
        },
        data () {
            return {
                keyword: '',
                displayList: []
            }
        },
        watch: {
            list (list) {
                this.displayList = list
            },
            displayList (list) {
                this.setChecked()
            },
            selected () {
                this.setChecked()
            },
            keyword () {
                this.handleFilter()
            }
        },
        created () {
            this.handleFilter = debounce(this.searchList, 300)
        },
        methods: {
            setChecked () {
                this.$nextTick(() => {
                    const ids = [...new Set(this.selected.map(data => data.host.bk_host_id))]
                    const selected = []
                    this.displayList.forEach(row => {
                        if (ids.includes(row.host.bk_host_id)) {
                            selected.push(row)
                        }
                        this.$refs.table.toggleRowSelection(row, ids.includes(row.host.bk_host_id))
                    })
                })
            },
            handleSelect (selection, row) {
                this.handleSelectionChange(selection)
            },
            handleSelectAll (selection) {
                this.handleSelectionChange(selection)
            },
            handleSelectionChange (selection) {
                const ids = [...new Set(selection.map(data => data.host.bk_host_id))]
                const removed = this.displayList.filter(item => !ids.includes(item.host.bk_host_id))
                this.$emit('select-change', { removed, selected: selection })
            },
            searchList () {
                if (this.keyword) {
                    this.displayList = this.list.filter(item => {
                        return new RegExp(this.keyword, 'i').test(item.host.bk_host_innerip)
                    })
                } else {
                    this.displayList = this.list
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .host-table {
        height: 100%;

        .search-bar {
            margin-bottom: 12px;
        }
    }
</style>
