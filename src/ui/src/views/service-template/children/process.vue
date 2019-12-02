<template>
    <div class="process-wrapper">
        <bk-table class="process-table"
            v-bkloading="{ isLoading: loading }"
            :data="showList"
            :max-height="$APP.height - 300">
            <bk-table-column v-for="column in table.header"
                :key="column.id"
                :prop="column.id"
                :label="column.name">
                <template slot-scope="{ row }">
                    <span v-if="column.id === 'bind_ip'">{{row[column.id] | ipText}}</span>
                    <span v-else>{{row[column.id] || '--'}}</span>
                </template>
            </bk-table-column>
            <bk-table-column :label="$t('操作')" prop="operation" v-if="$parent.isFormMode">
                <template slot-scope="{ row }">
                    <cmdb-auth :auth="$authResources(auth)">
                        <bk-button slot-scope="{ disabled }"
                            class="mr10"
                            theme="primary"
                            :disabled="disabled"
                            :text="true"
                            @click.stop="handleEdite(row['originData'])">
                            {{$t('编辑')}}
                        </bk-button>
                    </cmdb-auth>
                    <cmdb-auth :auth="$authResources(auth)">
                        <bk-button slot-scope="{ disabled }"
                            theme="primary"
                            :disabled="disabled"
                            :text="true"
                            @click.stop="handleDelete(row['originData'])">
                            {{$t('删除')}}
                        </bk-button>
                    </cmdb-auth>
                </template>
            </bk-table-column>
        </bk-table>
    </div>
</template>

<script>
    export default {
        filters: {
            ipText (value) {
                if (['1', '2'].includes(value)) {
                    const ip = ['127.0.0.1', '0.0.0.0']
                    return ip[value - 1]
                }
                return value || '--'
            }
        },
        props: {
            auth: {
                type: Object,
                default: () => ({})
            },
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
                            name: this.$t('进程名称'),
                            sortable: false
                        }, {
                            id: 'bk_process_name',
                            name: this.$t('进程别名'),
                            sortable: false
                        }, {
                            id: 'bind_ip',
                            name: this.$t('监听IP'),
                            sortable: false
                        }, {
                            id: 'port',
                            name: this.$t('端口'),
                            sortable: false
                        }, {
                            id: 'work_path',
                            name: this.$t('启动路径'),
                            sortable: false
                        }, {
                            id: 'user',
                            name: this.$t('启动用户'),
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
