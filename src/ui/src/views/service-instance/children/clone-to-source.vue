<template>
    <div class="source-layout">
        <!-- <p class="source-tips">{{$t('源实例主机提示')}}</p> -->
        <div class="table-options">
            <bk-button class="options-button"
                :disabled="!checked.length"
                @click="handleBatchEdit">
                {{$t('批量编辑')}}
            </bk-button>
        </div>
        <div class="source-table" ref="sourceTables">
            <bk-table
                :data="cloneProcesses"
                @selection-change="handleSelectChange">
                <bk-table-column type="selection" align="center" width="60" fixed class-name="bk-table-selection"></bk-table-column>
                <bk-table-column v-for="column in header"
                    :key="column.id"
                    :prop="column.id"
                    :label="column.name"
                    show-overflow-tooltip>
                    <template slot-scope="{ row }">
                        <cmdb-property-value v-if="column.id !== 'bind_info'"
                            :value="row[column.id]"
                            :show-unit="false"
                            :property="column.property">
                        </cmdb-property-value>
                        <process-bind-info-value v-else
                            :value="row[column.id]"
                            :property="column.property">
                        </process-bind-info-value>
                    </template>
                </bk-table-column>
                <bk-table-column :label="$t('操作')" fixed="right">
                    <template slot-scope="{ row }">
                        <button class="text-primary mr10" v-if="isRepeat(row)"
                            @click="handleEditProcess(row)">
                            <i class="bk-icon icon-exclamation-circle"></i>
                            {{$t('请编辑')}}
                        </button>
                        <button class="text-primary mr10" v-else
                            @click="handleEditProcess(row)">
                            {{$t('编辑')}}
                        </button>
                    </template>
                </bk-table-column>
            </bk-table>
        </div>
        <div class="page-options" :class="{ 'is-sticky': hasScrollbar }">
            <cmdb-auth :auth="{ type: $OPERATION.C_SERVICE_INSTANCE, relation: [bizId] }">
                <bk-button slot-scope="{ disabled }"
                    class="options-button"
                    theme="primary"
                    :disabled="!!repeatedProcesses.length || disabled"
                    @click="doClone">
                    {{$t('确定')}}
                </bk-button>
            </cmdb-auth>
            <bk-button class="options-button" @click="backToModule">{{$t('取消')}}</bk-button>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import { processTableHeader } from '@/dictionary/table-header'
    import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
    import ProcessBindInfoValue from '@/components/service/process-bind-info-value'
    import ProcessForm from '@/components/service/form/form.js'
    export default {
        name: 'clone-to-source',
        components: {
            ProcessBindInfoValue
        },
        props: {
            sourceProcesses: {
                type: Array,
                default () {
                    return {}
                }
            }
        },
        data () {
            return {
                checked: [],
                cloneProcesses: this.$tools.clone(this.sourceProcesses),
                properties: [],
                propertyUnique: [],
                hasScrollbar: false,
                formValuesReflect: {},
                processFormType: 'single'
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
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
            norepeatProperties () {
                const unique = this.propertyUnique.find(unique => unique.must_check) || {}
                const uniqueKeys = unique.keys || []
                return this.properties.filter(property => uniqueKeys.some(target => target.key_id === property.id))
            },
            repeatedProcesses () {
                return this.cloneProcesses.filter(cloneProcess => {
                    const sourceProcess = this.sourceProcesses.find(sourceProcess => sourceProcess.bk_process_id === cloneProcess.bk_process_id)
                    return this.norepeatProperties.length
                        && this.norepeatProperties.every(property => {
                            return sourceProcess[property.bk_property_id] === cloneProcess[property.bk_property_id]
                        })
                })
            },
            hostId () {
                return parseInt(this.$route.params.hostId)
            },
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
            }
        },
        watch: {
            sourceProcesses (source) {
                this.cloneProcesses = this.$tools.clone(source)
            }
        },
        mounted () {
            addResizeListener(this.$refs.sourceTables, this.resizeHandler)
        },
        beforeDestroy () {
            removeResizeListener(this.$refs.sourceTables, this.resizeHandler)
        },
        async created () {
            try {
                const [properties, propertyUnique] = await Promise.all([
                    this.getProcessProperties(),
                    this.getProcessPropertyUnique()
                ])
                this.properties = properties
                this.propertyUnique = propertyUnique
            } catch (e) {
                console.error(e)
            }
        },
        methods: {
            getProcessProperties () {
                const action = 'objectModelProperty/searchObjectAttribute'
                return this.$store.dispatch(action, {
                    params: {
                        bk_obj_id: 'process',
                        bk_supplier_account: this.$store.getters.supplierAccount
                    },
                    config: {
                        requestId: 'get_service_process_properties',
                        fromCache: true
                    }
                })
            },
            getProcessPropertyUnique () {
                const action = 'objectUnique/searchObjectUniqueConstraints'
                return this.$store.dispatch(action, {
                    objId: 'process',
                    params: {},
                    config: {
                        requestId: 'get_service_process_property_unique',
                        fromCache: true
                    }
                })
            },
            isRepeat (item) {
                return this.repeatedProcesses.some(process => process.bk_process_id === item.bk_process_id)
            },
            handleSelectChange (selection) {
                this.checked = selection.map(row => row.bk_process_id)
            },
            handleBatchEdit () {
                this.processFormType = 'batch'
                ProcessForm.show({
                    type: 'update',
                    title: this.$t('批量编辑'),
                    instance: {},
                    hostId: this.hostId,
                    bizId: this.bizId,
                    submitHandler: this.handleSubmit,
                    invisibleProperties: ['bind_info']
                })
            },
            handleEditProcess (item) {
                this.processFormType = 'single'
                ProcessForm.show({
                    type: 'update',
                    title: `${this.$t('编辑进程')}${item.bk_process_name}`,
                    instance: item,
                    hostId: this.hostId,
                    bizId: this.bizId,
                    submitHandler: this.handleSubmit
                })
            },
            async validateCustomComponent () {
                const { bindInfo } = this.$refs
                const customComponents = [bindInfo]
                const validatePromise = []
                customComponents.forEach(component => {
                    validatePromise.push(component.$validator.validateAll())
                    validatePromise.push(component.$validator.validateScopes())
                })
                const results = await Promise.all(validatePromise)
                return results.every(result => result)
            },
            handleSubmit (values, changedValues, instance) {
                if (this.processFormType === 'single') {
                    Object.assign(instance, changedValues)
                } else {
                    this.cloneProcesses.forEach(instance => {
                        Object.assign(instance, changedValues)
                    })
                }
            },
            async doClone () {
                try {
                    await this.$store.dispatch('serviceInstance/createProcServiceInstanceWithRaw', {
                        params: {
                            name: this.$parent.module.bk_module_name,
                            bk_biz_id: this.bizId,
                            bk_module_id: this.$route.params.moduleId,
                            instances: [
                                {
                                    bk_host_id: this.hostId,
                                    processes: this.getCloneProcessValues()
                                }
                            ]
                        }
                    })
                    this.$success(this.$t('克隆成功'))
                    this.backToModule()
                } catch (e) {
                    console.error(e)
                }
            },
            getCloneProcessValues () {
                return this.cloneProcesses.map(process => {
                    const value = {}
                    this.properties.forEach(property => {
                        value[property.bk_property_id] = process[property.bk_property_id]
                    })
                    return {
                        process_info: value
                    }
                })
            },
            backToModule () {
                this.$routerActions.back()
            },
            resizeHandler () {
                this.$nextTick(() => {
                    const scroller = this.$el.parentElement
                    this.hasScrollbar = scroller.scrollHeight > scroller.offsetHeight
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .source-layout {
        margin: 40px 0 0 0;
    }
    .options-button {
        height: 32px;
        margin: 0 6px 0 0;
        line-height: 30px;
    }
    .table-options {
        margin: 10px 0 0 0;
        padding: 0 20px;
    }
    .page-options {
        margin: 30px 0 0 0;
        padding: 10px 0 10px 20px;
        position: sticky;
        bottom: 0;
        left: 0;
        &.is-sticky {
            background-color: #FFF;
            border-top: 1px solid $borderColor;
            z-index: 100;
        }
    }
    .source-table {
        margin: 10px 0 0 0;
        padding: 0 20px;
    }
    .text-primary {
        .icon-exclamation-circle {
            @include inlineBlock(-1px);
            color: #FFB848;
        }
    }
</style>
