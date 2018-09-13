<template>
    <div class="process-layout">
        <cmdb-table ref="processTable"
            :header="header"
            :list="list"
            :defaultSort="defaultSort"
            :pagination.sync="pagination"
            :loading="$loading(['post_searchObjectAttribute_process', 'searchProcess'])"
            :wrapperMinusHeight="minusHeight"
            @handlePageChange="handlePageChange"
            @handleSizeChange="handleSizeChange"
            @handleSortChange="handleSortChange">
        </cmdb-table>
        <p class="footer-info">
            <span>{{$t('ProcessManagement["在进程管理中绑定进程到模块，"]')}}</span>
            <router-link class="link" to="process">
                {{$t('ProcessManagement["点击此跳转"]')}}
            </router-link>
        </p>
    </div>
</template>

<script>
    import { mapActions, mapGetters } from 'vuex'
    export default {
        props: {
            business: {
                type: Number,
                required: true
            },
            module: {
                type: Object,
                required: true
            }
        },
        data () {
            return {
                properties: [],
                header: [],
                list: [],
                defaultSort: '-bk_process_id',
                sort: '-bk_process_id',
                pagination: {
                    current: 1,
                    size: 10,
                    count: 0
                },
                minusHeight: 200
            }
        },
        computed: {
            ...mapGetters(['supplierAccount'])
        },
        watch: {
            module () {
                this.search()
            }
        },
        async created () {
            try {
                this.properties = await this.getProperties()
                this.header = await this.getHeader()
                this.search()
            } catch (e) {
                console.log(e)
            }
        },
        mounted () {
            this.calcMinusHeight()
        },
        methods: {
            ...mapActions('objectModelProperty', ['searchObjectAttribute']),
            ...mapActions('procConfig', ['searchProcess']),
            getProperties () {
                return this.searchObjectAttribute({
                    params: {
                        'bk_obj_id': 'process',
                        'bk_supplier_account': this.supplierAccount
                    },
                    config: {
                        requestId: 'post_searchObjectAttribute_process',
                        fromCache: true
                    }
                })
            },
            getHeader () {
                const headerKey = [
                    'bk_process_name',
                    'bk_func_id',
                    'bind_ip',
                    'port',
                    'protocol',
                    'bk_func_name'
                ]
                const header = []
                this.properties.forEach(property => {
                    if (headerKey.includes(property['bk_property_id'])) {
                        header.push({
                            id: property['bk_property_id'],
                            name: property['bk_property_name']
                        })
                    }
                })
                return Promise.resolve(header)
            },
            search () {
                const params = {
                    condition: {
                        'bk_module_name': this.module['bk_inst_name']
                    },
                    fields: [],
                    page: {
                        start: (this.pagination.current - 1) * this.pagination.size,
                        limit: this.pagination.size,
                        sort: this.sort
                    }
                }
                this.searchProcess({
                    bizId: this.business,
                    params,
                    config: {
                        requestId: 'searchProcess'
                    }
                }).then(data => {
                    this.pagination.count = data.count
                    this.list = this.$tools.flatternList(this.properties, data.info)
                })
            },
            handlePageChange (page) {
                this.pagination.page = page
                this.search()
            },
            handleSortChange (sort) {
                this.sort = sort
                this.search()
            },
            handleSizeChange (size) {
                this.pagination.size = size
                this.handlePageChange(1)
            },
            calcMinusHeight () {
                this.minusHeight = this.$refs.processTable.$el.getBoundingClientRect().top + 50
            }
        }
    }
</script>

<style lang="scss" scoped>
    .process-layout{
        margin: 20px 0 0 0;
    }
    .footer-info{
        font-size: 14px;
        margin: 10px 0 0 0;
        .link{
            color: #3c96ff;
            &:hover{
                color: #0082ff;
            }
        }
    }
</style>