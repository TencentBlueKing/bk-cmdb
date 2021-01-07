<template>
    <div class="process-wrapper">
        <bk-table class="process-table"
            v-bkloading="{ isLoading: loading }"
            :data="showList">
            <bk-table-column v-for="column in header"
                :key="column.id"
                :prop="column.id"
                :label="column.name"
                show-overflow-tooltip>
                <template slot-scope="{ row }">
                    <cmdb-property-value
                        v-if="column.id !== 'bind_info'"
                        :show-on="'cell'"
                        :value="row[column.id]"
                        :property="column.property">
                    </cmdb-property-value>
                    <process-bind-info-value v-else
                        :value="row[column.id]"
                        :property="column.property">
                    </process-bind-info-value>
                </template>
            </bk-table-column>
            <bk-table-column :label="$t('操作')" prop="operation" v-if="showOperation">
                <template slot-scope="{ row, $index }">
                    <cmdb-auth :auth="auth">
                        <bk-button slot-scope="{ disabled }"
                            class="mr10"
                            theme="primary"
                            :disabled="disabled"
                            :text="true"
                            @click.stop="handleEdit(row._original_, $index)">
                            {{$t('编辑')}}
                        </bk-button>
                    </cmdb-auth>
                    <cmdb-auth :auth="auth">
                        <bk-button slot-scope="{ disabled }"
                            theme="primary"
                            :disabled="disabled"
                            :text="true"
                            @click.stop="handleDelete(row._original_, $index)">
                            {{$t('删除')}}
                        </bk-button>
                    </cmdb-auth>
                </template>
            </bk-table-column>
        </bk-table>
    </div>
</template>

<script>
    import { processTableHeader } from '@/dictionary/table-header'
    import ProcessBindInfoValue from '@/components/service/process-bind-info-value'
    export default {
        components: {
            ProcessBindInfoValue
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
            },
            showOperation: Boolean
        },
        data () {
            return {}
        },
        computed: {
            header () {
                const header = processTableHeader.map(id => {
                    const property = this.properties.find(property => property.bk_property_id === id) || {}
                    return {
                        id: property.bk_property_id,
                        name: this.$tools.getHeaderPropertyName(property),
                        property
                    }
                })
                return header
            },
            showList () {
                const list = this.list.map(template => {
                    const result = {}
                    Object.keys(template).map(key => {
                        const type = typeof template[key]
                        if (type === 'object') {
                            result[key] = template[key]['value']
                        } else {
                            result[key] = template[key]
                        }
                    })
                    result._original_ = template
                    return result
                })
                list.sort((prev, next) => prev.process_id - next.process_id)
                return list
            }
        },
        methods: {
            handleEdit (process, index) {
                this.$emit('on-edit', process, index)
            },
            handleDelete (process, index) {
                this.$emit('on-delete', process, index)
            }
        }
    }
</script>

<style lang="scss" scoped>

</style>
