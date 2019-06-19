<template>
    <div class="process-wrapper">
        <cmdb-table class="process-table" ref="table"
            :loading="loading"
            :checked.sync="table.checked"
            :header="table.header"
            :list="showList"
            :pagination.sync="table.pagination"
            :default-sort="table.defaultSort"
            :wrapper-minus-height="300">
            <template v-for="(header, index) in table.header" :slot="header.id" slot-scope="{ item }">
                <template v-if="header.id === 'operation'">
                    <div :key="index">
                        <button class="text-primary mr10"
                            @click.stop="handleEdite(item['originData'])">
                            {{$t('Common["编辑"]')}}
                        </button>
                        <button class="text-primary"
                            @click.stop="handleDelete(item['originData'])">
                            {{$t('Common["删除"]')}}
                        </button>
                    </div>
                </template>
                <template v-else>
                    <span :key="index">{{item[header.id] ? item[header.id] : '--'}}</span>
                </template>
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
                list = this.$tools.flattenList(this.properties, list)
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
