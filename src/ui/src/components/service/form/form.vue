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
                :edit-auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }"
                :invisible-name-properties="invisibleNameProperties"
                :flex-properties="flexProperties"
                @on-edit="handleChangeInternalType">
            </cmdb-details>
            <cmdb-form v-else
                ref="form"
                v-bkloading="{ isLoading: pending }"
                :type="internalType"
                :inst="instance"
                :properties="visibleProperties"
                :property-groups="propertyGroups"
                :disabled-properties="bindedProperties"
                :invisible-name-properties="invisibleNameProperties"
                :flex-properties="flexProperties"
                :render-append="renderAppend"
                :custom-validator="validateCustomComponent"
                @on-submit="handleSaveProcess"
                @on-cancel="handleCancel">
                <template slot="bind_info">
                    <process-form-property-table
                        ref="bindInfo"
                        v-model="bindInfo"
                        :options="bindInfoProperty.option || []">
                    </process-form-property-table>
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
    import RenderAppend from './process-form-append-render'
    import ProcessFormPropertyTable from './process-form-property-table'
    export default {
        components: {
            ProcessFormPropertyTable
        },
        props: {
            type: String,
            serviceTemplateId: Number,
            processTemplateId: Number,
            instance: {
                type: Object,
                default: () => ({})
            },
            title: String,
            hostId: Number,
            bizId: Number,
            submitHandler: Function,
            invisibleProperties: {
                type: Array,
                default: () => ([])
            }
        },
        provide () {
            return {
                form: this
            }
        },
        data () {
            return {
                isShow: false,
                internalType: this.type,
                internalTitle: this.title,
                properties: [],
                propertyGroups: [],
                bindedProperties: [],
                processTemplate: {},
                pending: true,
                invisibleNameProperties: ['bind_info'],
                flexProperties: ['bind_info'],
                formValuesReflect: {}
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            bindInfoProperty () {
                return this.properties.find(property => property.bk_property_id === 'bind_info') || {}
            },
            bindInfo: {
                get () {
                    return this.formValuesReflect.bind_info || []
                },
                set (values) {
                    this.formValuesReflect.bind_info = values
                }
            },
            visibleProperties () {
                return this.properties.filter(property => !this.invisibleProperties.includes(property.bk_property_id))
            }
        },
        watch: {
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
        beforeDestroy () {
            this.teardownWatcher()
        },
        methods: {
            show () {
                this.isShow = true
            },
            teardownWatcher () {
                this.unwatchName && this.unwatchName()
                this.unwatchFormValues && this.unwatchFormValues()
            },
            updateFormWatcher () {
                if (this.internalType === 'view') {
                    this.teardownWatcher()
                } else {
                    this.$nextTick(() => {
                        const form = this.$refs.form
                        if (!form) {
                            return this.updateFormWatcher() // 递归nextTick等待form创建完成
                        }
                        // watch form组件表单值，用于获取bind_info字段给进程表格字段组件使用
                        this.unwatchFormValues = this.$watch(() => {
                            return form.values
                        }, values => {
                            this.formValuesReflect = values
                        }, { immediate: true })
                        // watch 名称，在用户未修改进程别名时，自动同步进程名称到进程别名
                        this.unwatchName = this.$watch(() => {
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
            async getProcessTemplate () {
                try {
                    this.processTemplate = await this.$store.dispatch('processTemplate/getProcessTemplate', {
                        params: {
                            processTemplateId: this.processTemplateId
                        },
                        config: {
                            cancelPrevious: true
                        }
                    })
                    const { property } = this.processTemplate
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
            async validateCustomComponent () {
                const customComponents = []
                const { bindInfo } = this.$refs
                if (bindInfo) {
                    customComponents.push(bindInfo)
                }
                const validatePromise = []
                customComponents.forEach(component => {
                    validatePromise.push(component.$validator.validateAll())
                    validatePromise.push(component.$validator.validateScopes())
                })
                const results = await Promise.all(validatePromise)
                return results.every(result => result)
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
            renderAppend (h, { property, type }) {
                if (this.bindedProperties.includes(property.bk_property_id)) {
                    return RenderAppend(h, {
                        serviceTemplateId: this.serviceTemplateId,
                        property: property,
                        bizId: this.bizId
                    })
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
