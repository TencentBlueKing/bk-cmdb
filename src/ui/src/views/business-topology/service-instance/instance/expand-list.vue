<template>
    <bk-table v-if="!serviceInstance.pending"
        :data="list"
        :outer-border="false"
        :header-cell-style="{ backgroundColor: '#fff' }"
        v-bind="dynamicProps"
        v-bkloading="{ isLoading: $loading(request.list) }">
        <bk-table-column v-for="property in header"
            :key="property.bk_property_id"
            :label="property.bk_property_name"
            :prop="property.bk_property_id"
            :show-overflow-tooltip="property.bk_property_id !== 'bind_info'">
            <template slot-scope="{ row }">
                <cmdb-property-value v-if="property.bk_property_id !== 'bind_info'"
                    :theme="property.bk_property_id === 'bk_func_name' ? 'primary' : 'default'"
                    :value="row.property[property.bk_property_id]"
                    :show-unit="false"
                    :property="property"
                    @click.native="handleView(row)">
                </cmdb-property-value>
                <process-bind-info-value v-else
                    :value="row.property[property.bk_property_id]"
                    :property="property">
                </process-bind-info-value>
            </template>
        </bk-table-column>
        <bk-table-column width="150" :resizable="false">
            <div class="options-wrapper" slot-scope="{ row }">
                <cmdb-auth class="mr10" :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }">
                    <bk-button slot-scope="{ disabled }"
                        theme="primary" text
                        :disabled="disabled"
                        @click="handleEdit(row)">
                        {{$t('编辑')}}
                    </bk-button>
                </cmdb-auth>
                <cmdb-auth :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }" v-if="!row.relation.process_template_id">
                    <bk-button slot-scope="{ disabled }"
                        theme="primary" text
                        :disabled="disabled"
                        @click="handleDelete(row)">
                        {{$t('删除')}}
                    </bk-button>
                </cmdb-auth>
            </div>
        </bk-table-column>
    </bk-table>
</template>

<script>
    import { processPropertyRequestId } from '@/components/service/form/symbol'
    import { processTableHeader } from '@/dictionary/table-header'
    import { mapGetters } from 'vuex'
    import Form from '@/components/service/form/form.js'
    import ProcessBindInfoValue from '@/components/service/process-bind-info-value'
    import Bus from '../common/bus'
    export default {
        components: {
            ProcessBindInfoValue
        },
        props: {
            serviceInstance: Object
        },
        data () {
            return {
                properties: [],
                header: [],
                list: [],
                request: {
                    list: Symbol('getList'),
                    delete: Symbol('delete')
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('objectBiz', ['bizId']),
            dynamicProps () {
                const dynamicProps = {}
                const paddingHeight = 43
                const rowHeight = 42
                if (this.list.length && this.list.length < 3) {
                    dynamicProps.height = paddingHeight + rowHeight * (this.list.length + 1)
                }
                return dynamicProps
            }
        },
        created () {
            this.getProperties()
            this.getList()
            Bus.$on('refresh-process-list', this.handleRefresh)
        },
        beforeDestroy () {
            Bus.$off('refresh-process-list', this.handleRefresh)
        },
        methods: {
            async getProperties () {
                try {
                    this.properties = await this.$store.dispatch('objectModelProperty/searchObjectAttribute', {
                        params: {
                            bk_obj_id: 'process',
                            bk_supplier_account: this.supplierAccount
                        },
                        config: {
                            requestId: processPropertyRequestId,
                            fromCache: true
                        }
                    })
                    this.setHeader()
                } catch (error) {
                    console.error(error)
                }
            },
            setHeader () {
                const header = []
                processTableHeader.forEach(id => {
                    const property = this.properties.find(property => property.bk_property_id === id)
                    if (property) {
                        header.push(property)
                    }
                })
                this.header = header
            },
            handleRefresh (target) {
                if (target !== this.serviceInstance) {
                    return
                }
                this.getList()
            },
            async getList () {
                try {
                    this.list = await this.$store.dispatch('processInstance/getServiceInstanceProcesses', {
                        params: {
                            bk_biz_id: this.bizId,
                            service_instance_id: this.serviceInstance.id
                        },
                        config: {
                            requestId: this.request.list,
                            cancelPrevious: true,
                            cancelWhenRouteChange: true
                        }
                    })
                } catch (error) {
                    console.error(error)
                } finally {
                    this.$emit('update-list', this.list)
                }
            },
            handleView (row) {
                Form.show({
                    type: 'view',
                    title: this.$t('查看进程'),
                    instance: row.property,
                    hostId: row.relation.bk_host_id,
                    bizId: this.bizId,
                    serviceTemplateId: this.serviceInstance.service_template_id,
                    processTemplateId: row.relation.process_template_id,
                    submitHandler: this.editSubmitHandler
                })
            },
            handleEdit (row) {
                Form.show({
                    type: 'update',
                    title: this.$t('编辑进程'),
                    instance: row.property,
                    hostId: row.relation.bk_host_id,
                    bizId: this.bizId,
                    serviceTemplateId: this.serviceInstance.service_template_id,
                    processTemplateId: row.relation.process_template_id,
                    submitHandler: this.editSubmitHandler
                })
            },
            async editSubmitHandler (values, changedValues, instance) {
                try {
                    await this.$store.dispatch('processInstance/updateServiceInstanceProcess', {
                        params: {
                            bk_biz_id: this.bizId,
                            processes: [{ ...instance, ...values }]
                        }
                    })
                    this.getList()
                } catch (error) {
                    console.error(error)
                }
            },
            handleDelete (row) {
                this.$bkInfo({
                    title: this.$t('确定删除该进程'),
                    confirmFn: async () => {
                        try {
                            await this.$store.dispatch('processInstance/deleteServiceInstanceProcess', {
                                config: {
                                    data: {
                                        bk_biz_id: this.bizId,
                                        process_instance_ids: [row.property.bk_process_id]
                                    },
                                    requestId: this.request.delete
                                }
                            })
                            if (this.list.length === 1) {
                                this.$emit('update-list', [])
                            } else {
                                this.getList()
                            }
                        } catch (error) {
                            console.error(error)
                        }
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .options-wrapper {
        display: none;
    }
    /deep/ {
        .bk-table-row:hover {
            .options-wrapper {
                display: block;
            }
        }
    }
</style>
