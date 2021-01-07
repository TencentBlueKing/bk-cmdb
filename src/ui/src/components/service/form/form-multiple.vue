<template>
    <bk-sideslider
        :width="800"
        :title="$t('批量编辑')"
        :is-show.sync="isShow"
        :before-close="beforeClose"
        @hidden="handleHidden">
        <cmdb-form-multiple slot="content"
            ref="form"
            v-bkloading="{ isLoading: pending }"
            :properties="properties"
            :property-groups="propertyGroups"
            :uneditable-properties="bindedProperties"
            @on-submit="handleSaveProcess"
            @on-cancel="beforeClose">
        </cmdb-form-multiple>
    </bk-sideslider>
</template>

<script>
    import { mapGetters } from 'vuex'
    import {
        processPropertyRequestId,
        processPropertyGroupsRequestId
    } from './symbol'
    export default {
        props: {
            serviceTemplateId: Number,
            processTemplateId: Number,
            submitHandler: Function
        },
        data () {
            return {
                isShow: false,
                properties: [],
                propertyGroups: [],
                // 绑定信息暂不支持批量编辑
                bindedProperties: ['bind_info'],
                pending: true
            }
        },
        computed: {
            ...mapGetters(['supplierAccount'])
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
        methods: {
            show () {
                this.isShow = true
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
            beforeClose () {
                if (this.$refs.form.hasChange) {
                    return new Promise((resolve, reject) => {
                        this.$bkInfo({
                            title: this.$t('确认退出'),
                            subTitle: this.$t('退出会导致未保存信息丢失'),
                            extCls: 'bk-dialog-sub-header-center',
                            confirmFn: () => {
                                this.isShow = false
                                resolve(true)
                            },
                            cancelFn: () => resolve(false)
                        })
                    })
                }
                this.isShow = false
                return Promise.resolve(true)
            }
        }
    }
</script>
