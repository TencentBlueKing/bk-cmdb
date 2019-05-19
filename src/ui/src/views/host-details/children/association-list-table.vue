<template>
    <div class="table">
        <div class="table-info clearfix">
            <div class="info-title fl">
                {{model.bk_obj_name}}
            </div>
            <div class="info-pagination fr"></div>
        </div>
        <cmdb-table class="association-table"
            :loading="$loading(propertyRequest, instanceRequest)"
            :header="header"
            :list="list"
            :show-footer="false"
            :sortable="false"
            :max-height="462"
            :empty-height="100">
            <template slot="__operation__" slot-scope="{ item }">
                <span class="text-primary"
                    @click="cancelAssociation(item)">
                    {{$t('Association["取消关联"]')}}
                </span>
            </template>
        </cmdb-table>
    </div>
</template>

<script>
    export default {
        name: 'cmdb-host-association-list-table',
        props: {
            type: {
                type: String,
                required: true
            },
            id: {
                type: String,
                required: true
            },
            instances: {
                type: Array,
                required: true
            }
        },
        data () {
            return {
                properties: [],
                list: [],
                pagination: {
                    count: 0,
                    current: 1,
                    size: 10
                }
            }
        },
        computed: {
            hostId () {
                return parseInt(this.$route.params.id)
            },
            model () {
                return this.$store.getters['objectModelClassify/getModelById'](this.id)
            },
            isBusinessModel () {
                return !!this.$tools.getMetadataBiz(this.model)
            },
            propertyRequest () {
                return `get_${this.id}_association_list_table_properties`
            },
            instanceRequest () {
                return `get_${this.id}_association_list_table_instances`
            },
            page () {
                return {
                    limit: this.pagination.size,
                    start: (this.pagination.current - 1) * this.pagination.size
                }
            },
            instanceIds () {
                return this.instances.map(instance => instance.bk_inst_id)
            },
            header () {
                const keyMap = {
                    host: 'bk_host_id',
                    biz: 'bk_biz_id'
                }
                const headerProperties = this.$tools.getDefaultHeaderProperties(this.properties)
                return [{
                    id: keyMap[this.id] || 'bk_inst_id',
                    type: 'checkbox'
                }].concat(headerProperties.map(property => {
                    return {
                        id: property.bk_property_id,
                        name: property.bk_property_name
                    }
                })).concat([{
                    id: '__operation__',
                    name: this.$t('Common["操作"]'),
                    width: 150
                }])
            }
        },
        created () {
            this.getProperties()
            this.getInstances()
        },
        methods: {
            async getProperties () {
                try {
                    this.properties = await this.$store.dispatch('objectModelProperty/searchObjectAttribute', {
                        params: this.$injectMetadata({
                            bk_obj_id: this.id
                        }, {
                            inject: this.isBusinessModel
                        }),
                        config: {
                            fromCache: true,
                            requestId: this.propertyRequest
                        }
                    })
                } catch (e) {
                    console.error(e)
                    this.properties = []
                }
            },
            async getInstances () {
                let promise
                const config = {
                    requestId: this.instanceRequest,
                    cancelPrevious: true,
                    globalError: false
                }
                try {
                    switch (this.id) {
                        case 'host':
                            promise = this.getHostInstances(config)
                            break
                        case 'biz':
                            promise = this.getBusinessInstances(config)
                            break
                        default:
                            promise = this.getModelInstances(config)
                    }
                    const data = await promise
                    this.list = data.info
                    this.pagination.count = data.count
                } catch (e) {
                    console.error(e)
                    this.list = []
                    this.pagination.count = 0
                }
            },
            getHostInstances (config) {
                const models = ['biz', 'set', 'module', 'host']
                const hostCondition = {
                    field: 'bk_host_id',
                    operator: '$in',
                    value: this.instanceIds
                }
                const condition = models.map(model => {
                    return {
                        bk_obj_id: model,
                        fields: [],
                        condition: model === 'host' ? [hostCondition] : []
                    }
                })
                return this.$store.dispatch('hostSearch/searchHost', {
                    params: this.$injectMetadata({
                        bk_biz_id: -1,
                        condition,
                        id: {
                            data: [],
                            exact: 0,
                            flag: 'bk_host_innerip|bk_host_outerip'
                        },
                        page: {
                            ...this.page,
                            sort: 'bk_host_id'
                        }
                    }),
                    config
                }).then(data => {
                    return {
                        count: data.count,
                        info: data.info.map(item => item.host)
                    }
                })
            },
            getBusinessInstances (config) {
                return this.$store.dispatch('objectBiz/searchBusiness', {
                    params: {
                        condition: {
                            bk_biz_id: {
                                $in: this.instanceIds
                            }
                        },
                        fields: [],
                        page: {
                            ...this.page,
                            sort: 'bk_biz_id'
                        }
                    },
                    config
                })
            },
            getModelInstances (config) {
                return this.$store.dispatch('objectCommonInst/searchInst', {
                    objId: this.id,
                    params: {
                        fields: {},
                        condition: {
                            [this.id]: [{
                                field: 'bk_inst_id',
                                operator: '$in',
                                value: this.instanceIds
                            }]
                        },
                        page: {
                            ...this.page,
                            sort: 'bk_inst_id'
                        }
                    },
                    config
                }).then(data => {
                    data = data || {
                        count: 0,
                        info: []
                    }
                    return data
                })
            },
            cancelAssociation (item) {}
        }
    }
</script>

<style lang="scss" scoped>
    .table {
        margin: 0 0 12px 0;
    }
    .table-info {
        height: 42px;
        padding: 0 20px;
        border-radius: 2px;
        line-height: 42px;
        background-color: #DCDEE5;
        font-size: 14px;
    }
</style>
