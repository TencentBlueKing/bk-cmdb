<template>
    <div class="process-wrapper">
        <v-table class="process-table" 
            :header="table.header"
            :list="table.list"
            :defaultSort="table.defaultSort"
            :pagination="table.pagination"
            :loading="table.isLoading"
            @handlePageChange="setCurrentPage"
            @handleSizeChange="setCurrentSize"
            @handleSortChange="setCurrentSort">
            <!-- <template v-for="({property, id, name}) in table.header" :slot="id" slot-scope="{ item }">
                <template v-if="property['bk_property_type'] === 'enum'">
                    {{getEnumCell(item[id], property)}}
                </template>
                <template v-else-if="!!property['bk_asst_obj_id']">
                    {{getAssociateCell(item[id])}}
                </template>
                <template v-else>{{item[id]}}</template>
            </template> -->
        </v-table>
    </div>
</template>

<script>
    import vTable from '@/components/table/table'
    import { mapGetters } from 'vuex'
    export default {
        props: {
            bizId: {
                type: Number
            },
            moduleName: {
                type: String
            },
            isShow: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                attribute: [],
                table: {
                    header: [],
                    list: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        size: 10
                    },
                    defaultSort: '-bk_process_id',
                    sort: '-bk_process_id',
                    isLoading: false
                }
            }
        },
        computed: {
            ...mapGetters([
                'bkSupplierAccount'
            ]),
            searchParams () {
                let params = {
                    condition: {
                        bk_module_name: this.moduleName
                    },
                    fields: [],
                    page: {
                        start: (this.table.pagination.current - 1) * this.table.pagination.size,
                        limit: this.table.pagination.size,
                        sort: this.table.sort
                    }
                }
                return params
            }
        },
        watch: {
            isShow (isShow) {
                if (isShow) {
                    this.table.isLoading = true
                    Promise.all([
                        this.getProcessAttribute(),
                        this.getProcessList()
                    ]).finally(() => {
                        this.table.isLoading = false
                    })
                }
            },
            async moduleName () {
                this.table.isLoading = true
                await this.getProcessList()
                this.table.isLoading = false
            },
            attribute (attribute) {
                let headerLead = []
                let headerMiddle = []
                let headerTail = []
                attribute.map(property => {
                    let {
                        isonly,
                        isrequired,
                        'bk_isapi': bkIsapi,
                        'bk_property_id': bkPropertyId,
                        'bk_property_name': bkPropertyName
                    } = property
                    let headerItem = {
                        id: bkPropertyId,
                        name: bkPropertyName,
                        property: property
                    }
                    if (!bkIsapi) {
                        if (isonly && isrequired) {
                            headerLead.push(headerItem)
                        } else if (isonly || isonly) {
                            headerMiddle.push(headerItem)
                        } else {
                            headerTail.push(headerItem)
                        }
                    }
                })
                this.table.header = headerLead.concat(headerMiddle, headerTail).slice(0, 6)
            }
        },
        methods: {
            async getProcessAttribute () {
                let params = {
                    'bk_obj_id': 'process',
                    'bk_supplier_account': this.bkSupplierAccount
                }
                try {
                    const res = await this.$axios.post('object/attr/search', params)
                    if (!res.result) {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    this.attribute = res.data
                } catch (e) {
                    this.$alertMsg(e.message || e.data['bk_error_msg'] || e.statusText)
                }
            },
            async getProcessList () {
                try {
                    const res = await this.$axios.post(`proc/search/${this.bkSupplierAccount}/${this.bizId}`, this.searchParams)
                    if (!res.result) {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    this.table.list = res.data.info
                    this.table.pagination.count = res.data.count
                } catch (e) {
                    this.$alertMsg(e.message || e.data['bk_error_msg'] || e.statusText)
                }
            },
            setCurrentSize (size) {
                this.table.pagination.size = size
                this.setCurrentPage(1)
            },
            setCurrentSort (sort) {
                this.table.sort = sort
                this.setCurrentPage(1)
            },
            setCurrentPage (page) {
                this.table.pagination.current = page
                this.getProcessList()
            }
        },
        components: {
            vTable
        }
    }
</script>

<style lang="scss" scoped>
    .process-wrapper {
        padding: 20px;
    }
</style>

