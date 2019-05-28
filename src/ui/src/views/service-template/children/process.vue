<template>
    <div class="process-wrapper">
        <cmdb-table class="process-table" ref="table"
            :loading="$loading('post_searchProcess_list')"
            :checked.sync="table.checked"
            :header="table.header"
            :list="list"
            :pagination.sync="table.pagination"
            :default-sort="table.defaultSort"
            :wrapper-minus-height="300">
            <template slot="operation" slot-scope="{ item }">
                <button class="text-primary mr10"
                    @click.stop="handleEdite(item)">
                    {{$t('Common["编辑"]')}}
                </button>
                <button class="text-primary"
                    @click.stop="handleDelete(item)">
                    {{$t('Common["删除"]')}}
                </button>
            </template>
        </cmdb-table>
    </div>
</template>

<script>
    export default {
        props: {
            list: {
                type: Array,
                default: () => {
                    return []
                }
            },
            properties: {
                type: Array,
                default: () => {

                }
            }
        },
        data () {
            return {
                table: {
                    header: [
                        {
                            id: 'bk_process_name',
                            name: this.$t("ProcessManagement['进程名称']"),
                            sortable: false
                        }, {
                            id: 'bind_ip',
                            name: this.$t("ProcessManagement['监听IP']"),
                            sortable: false
                        }, {
                            id: 'port',
                            name: this.$t("ProcessManagement['端口']"),
                            sortable: false
                        }, {
                            id: 'work_path',
                            name: this.$t("ProcessManagement['启动路径']"),
                            sortable: false
                        }, {
                            id: 'user',
                            name: this.$t("ProcessManagement['启动用户']"),
                            sortable: false
                        }, {
                            id: 'operation',
                            name: this.$t("Common['操作']"),
                            sortable: false
                        }
                    ],
                    pagination: {
                        current: 1,
                        count: 0,
                        size: 10
                    },
                    defaultSort: '-bk_process_id',
                    sort: '-bk_process_id'
                }
            }
        },
        methods: {
            handleEdite (process) {

            },
            handleDelete (process) {

            }
        }
    }
</script>

<style lang="scss" scoped>

</style>
