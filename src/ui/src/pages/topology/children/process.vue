<template>
    <div class="process-wrapper" ref="wrapper">
        <v-table class="process-table" 
            :header="table.header"
            :list="table.list"
            :defaultSort="table.defaultSort"
            :pagination="table.pagination"
            :loading="table.isLoading"
            :height="table.height"
            @handlePageChange="setCurrentPage"
            @handleSizeChange="setCurrentSize"
            @handleSortChange="setCurrentSort">
            <template v-for="({property, id, name}) in table.header" :slot="id" slot-scope="{ item }">
                <template v-if="property['bk_property_type'] === 'enum'">
                    {{getEnumCell(item[id], property)}}
                </template>
                <template v-else-if="!!property['bk_asst_obj_id']">
                    {{getAssociateCell(item[id])}}
                </template>
                <template v-else>{{item[id]}}</template>
            </template>
        </v-table>
        <p class="footer-info">
            <span>{{$t('ProcessManagement["在进程管理中绑定进程到模块，"]')}}</span>
            <router-link to="process">
                {{$t('ProcessManagement["点击此跳转"]')}}
            </router-link>
        </p>
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
                    height: 200,
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
                    this.setTableHeight()
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
                        'bk_property_id': bkPropertyId,
                        'bk_property_name': bkPropertyName
                    } = property
                    let headerItem = {
                        id: bkPropertyId,
                        name: bkPropertyName,
                        property: property
                    }
                    switch (bkPropertyId) {
                        case 'bk_process_name':
                            headerMiddle[0] = headerItem
                            break
                        case 'bk_func_id':
                            headerMiddle[1] = headerItem
                            break
                        case 'bind_ip':
                            headerMiddle[2] = headerItem
                            break
                        case 'port':
                            headerMiddle[3] = headerItem
                            break
                        case 'protocol':
                            headerMiddle[4] = headerItem
                            break
                        case 'bk_func_name':
                            headerMiddle[5] = headerItem
                            break
                    }
                })
                headerMiddle = headerMiddle.filter(header => {
                    return header !== undefined
                })
                this.table.header = headerLead.concat(headerMiddle, headerTail).slice(0, 6)
            }
        },
        methods: {
            setTableHeight () {
                const wrapper = this.$refs.wrapper.getBoundingClientRect()
                const body = document.body.getBoundingClientRect()
                this.table.height = Number(`${body.height - wrapper.top - 41 - 30}`)
            },
            getEnumCell (data, property) {
                let obj = property.option.find(({id}) => {
                    return id === data
                })
                if (obj) {
                    return obj.name
                }
            },
            // 计算关联属性单元格显示的值
            getAssociateCell (data) {
                let label = []
                if (Array.isArray(data)) {
                    data.map(({bk_inst_name: bkInstName}) => {
                        if (bkInstName) {
                            label.push(bkInstName)
                        }
                    })
                }
                return label.join(',')
            },
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
        .footer-info {
            margin-top: 10px;
            font-size: 0;
            span,
            a {
                font-size: 14px;
                line-height: 20px;
                height: 20px;
            }
            a {
                color: #3c96ff;
                &:hover {
                    color: #0082ff;
                }
            }
        }
    }
</style>

