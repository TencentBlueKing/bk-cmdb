<template>
    <div class="process-wrapper">
        <bk-table class="process-table"
            v-bkloading="{ isLoading: loading }"
            :list="showList"
            :max-height="$APP.height - 300">
            <bk-table-column v-for="column in table.header"
                :key="column.id"
                :prop="column.id"
                :label="column.name">
                <template slot-scope="{ row }">
                    {{row[column.id] || '--'}}
                </template>
            </bk-table-column>
            <bk-table-column :label="$t('Common[\'操作\']')">
                <template slot-scope="{ row }">
                    <button class="text-primary mr10"
                        @click.stop="handleEdite(row['originData'])">
                        {{$t('Common["编辑"]')}}
                    </button>
                    <button class="text-primary"
                        @click.stop="handleDelete(row['originData'])">
                        {{$t('Common["删除"]')}}
                    </button>
                </template>
            </bk-table-column>
        </bk-table>
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
                    return []
                }
            },
            loading: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                table: {
                    header: [
                        {
                            id: 'bk_func_name',
                            name: this.$t("ProcessManagement['进程名称']"),
                            sortable: false
                        }, {
                            id: 'bk_process_name',
                            name: this.$t("ProcessManagement['进程别名']"),
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
                        }
                    ]
                }
            }
        },
        computed: {
            showList () {
                let list = this.list.map(template => {
                    const result = {}
                    // const property = template['property']
                    Object.keys(template).map(key => {
                        const type = typeof template[key]
                        if (type === 'object') {
                            result[key] = template[key]['value']
                        } else {
                            result[key] = template[key]
                        }
                    })
                    result['originData'] = template
                    return result
                })
                list = this.$tools.flattenList(this.properties, list).sort((prev, next) => prev.process_id - next.process_id)
                return list
            }
        },
        methods: {
            handleEdite (process) {
                this.$emit('on-edit', process)
            },
            handleDelete (process) {
                this.$emit('on-delete', process)
            }
        }
    }
</script>

<style lang="scss" scoped>

</style>
