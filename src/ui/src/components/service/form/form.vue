<template>
    <bk-sideslider
        :width="800"
        :title="internalTitle"
        :is-show.sync="isShow"
        :before-close="beforeClose"
        @hidden="handleHidden">
        <template slot="content">
            <cmdb-details v-if="internalType === 'view'"
                :properties="properties"
                :property-groups="propertyGroups"
                :inst="instance"
                :show-delete="false"
                :edit-auth="{ type: $OPERATION.U_SERVICE_INSTANCE, bk_biz_id: bizId }"
                @on-edit="handleChangeInternalType">
            </cmdb-details>
            <cmdb-form v-else
                ref="form"
                v-bkloading="{ isLoading: pending }"
                :type="internalType"
                :inst="instance"
                :properties="properties"
                :property-groups="propertyGroups"
                :disabled-properties="bindedProperties"
                :render-tips="renderTips"
                @on-submit="handleSaveProcess"
                @on-cancel="handleCancel">
                <template slot="bind_ip">
                    <cmdb-input-select
                        name="bindIP"
                        :disabled="isBindIPDisabled"
                        :placeholder="$t('请选择或输入IP')"
                        :options="bindIPList"
                        :validate="bindIPRules"
                        v-model="bindIP">
                    </cmdb-input-select>
                </template>
            </cmdb-form>
        </template>
    </bk-sideslider>
</template>

<script>
    import { mapGetters } from 'vuex'
    import {
        processPropertyRequestId,
        processPropertyGroupsRequestId
    } from './symbol'
    import RenderTips from './process-form-tips-render'
    export default {
        props: {
            type: String,
            serviceTemplateId: Number,
            processTemplateId: Number,
            instance: Object,
            title: String,
            hostId: Number,
            submitHandler: Function
        },
        data () {
            return {
                isShow: false,
                internalType: this.type,
                internalTitle: this.title,
                properties: [],
                propertyGroups: [],
                bindedProperties: [],
                bindIP: this.instance ? this.instance.bind_ip : '',
                bindIPList: [],
                pending: true
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('objectBiz', ['bizId']),
            bindIPProperty () {
                return this.properties.find(property => property.bk_property_id === 'bind_ip')
            },
            isBindIPDisabled () {
                return this.bindedProperties.includes('bind_ip')
            },
            bindIPRules () {
                if (!this.bindIPProperty) {
                    return {}
                }
                const rules = {}
                if (this.bindIPProperty.isrequired) {
                    rules.required = true
                }
                rules.regex = this.bindIPProperty.option
                return rules
            }
        },
        watch: {
            bindIP (ip) {
                this.$refs.form.values.bind_ip = ip
            },
            internalType (type) {
                this.updateFormWatcher()
            }
        },
        async created () {
            try {
                const request = [
                    this.getProperties(),
                    this.getPropertyGroups()
                ]
                if (this.hostId) {
                    request.push(this.getBindIPList())
                }
                if (this.processTemplateId) {
                    request.push(this.getProcessTemplate())
                }
                await Promise.all(request)
            } catch (error) {
                console.error(error)
            } finally {
                this.pending = false
            }
        },
        mounted () {
            this.updateFormWatcher()
        },
        methods: {
            show () {
                this.isShow = true
            },
            updateFormWatcher () {
                if (this.internalType === 'view') {
                    this.unwatchForm && this.unwatchForm()
                } else {
                    this.$nextTick(() => {
                        const form = this.$refs.form
                        this.unwatchForm = this.$watch(() => {
                            return form.values.bk_func_name
                        }, (newVal, oldValue) => {
                            if (form.values.bk_process_name === oldValue) {
                                form.values.bk_process_name = newVal
                            }
                        })
                    })
                }
            },
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
                } catch (error) {
                    console.error(error)
                    this.properties = []
                }
            },
            async getPropertyGroups () {
                try {
                    this.propertyGroups = await this.$store.dispatch('objectModelFieldGroup/searchGroup', {
                        objId: 'process',
                        params: {},
                        config: {
                            requestId: processPropertyGroupsRequestId,
                            fromCache: true
                        }
                    })
                } catch (error) {
                    console.error(error)
                    this.propertyGroups = []
                }
            },
            async getBindIPList () {
                try {
                    const { options } = await this.$store.dispatch('serviceInstance/getInstanceIpByHost', {
                        hostId: this.hostId,
                        config: {
                            requestId: 'getInstanceIpByHost'
                        }
                    })
                    this.bindIPList = options.map(ip => ({ id: ip, name: ip }))
                } catch (error) {
                    this.bindIPList = []
                    console.error(error)
                }
            },
            async getProcessTemplate () {
                try {
                    const { property } = await this.$store.dispatch('processTemplate/getProcessTemplate', {
                        params: {
                            processTemplateId: this.processTemplateId
                        },
                        config: {
                            cancelPrevious: true
                        }
                    })
                    const bindedProperties = []
                    Object.keys(property).forEach(key => {
                        if (property[key].as_default_value) {
                            bindedProperties.push(key)
                        }
                    })
                    this.bindedProperties = bindedProperties
                } catch (error) {
                    console.error(error)
                }
            },
            handleHidden () {
                this.$emit('close')
            },
            async handleSaveProcess (values, changedValues, instance) {
                try {
                    this.pending = true
                    await this.submitHandler(values, changedValues, instance)
                    this.isShow = false
                } catch (error) {
                    console.error(error)
                } finally {
                    this.pending = false
                }
            },
            async handleCancel () {
                const userConfirm = await this.beforeClose()
                if (!userConfirm) {
                    return false
                }
                if (this.type === 'view') {
                    this.internalType = this.type
                    this.internalTitle = this.title
                } else {
                    this.isShow = false
                }
            },
            beforeClose () {
                if (this.internalType === 'view') return Promise.resolve(true)
                const formChanged = !!Object.values(this.$refs.form.changedValues).length
                if (formChanged) {
                    return new Promise((resolve, reject) => {
                        this.$bkInfo({
                            title: this.$t('确认退出'),
                            subTitle: this.$t('退出会导致未保存信息丢失'),
                            extCls: 'bk-dialog-sub-header-center',
                            confirmFn: () => {
                                resolve(true)
                            },
                            cancelFn: () => resolve(false)
                        })
                    })
                }
                return Promise.resolve(true)
            },
            renderTips (h, { property, type }) {
                if (this.bindedProperties.includes(property.bk_property_id)) {
                    return RenderTips(h, { serviceTemplateId: this.serviceTemplateId })
                }
                return ''
            },
            handleChangeInternalType () {
                this.internalType = 'update'
                this.internalTitle = this.$t('编辑进程')
            }
        }
    }
</script>
