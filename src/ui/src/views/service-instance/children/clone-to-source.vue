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
                    <template slot-scope="{ row }">{{row[column.id] | formatter(column.property)}}</template>
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
            <cmdb-auth :auth="$authResources({ type: $OPERATION.C_SERVICE_INSTANCE })">
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
        <bk-sideslider
            v-transfer-dom
            :is-show.sync="processForm.show"
            :title="processForm.title"
            :width="800"
            :before-close="handleBeforeClose">
            <cmdb-form slot="content"
                ref="processForm"
                type="update"
                v-if="processForm.show"
                :properties="properties"
                :property-groups="propertyGroups"
                :object-unique="processForm.type === 'single' ? [] : propertyUnique"
                :inst="processForm.instance"
                :invisible-name-properties="invisibleNameProperties"
                :flex-properties="flexProperties"
                :uneditable-properties="uneditableProperties"
                @on-submit="handleSubmit"
                @on-cancel="handleBeforeClose">
                <template slot="bind_info">
                    <process-bind-table v-bind="processTableProps" @change="handleProcessTableChange"></process-bind-table>
                </template>
            </cmdb-form>
        </bk-sideslider>
    </div>
</template>

<script>
    import { processTableHeader } from '@/dictionary/table-header'
    import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
    import ProcessBindTable from '@/components/service/form/process-bind-table'
    export default {
        name: 'clone-to-source',
        components: {
            ProcessBindTable
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
                propertyGroups: [],
                propertyUnique: [],
                processForm: {
                    show: false,
                    type: 'single',
                    title: '',
                    instance: {}
                },
                hasScrollbar: false,
                invisibleNameProperties: ['bind_info'],
                flexProperties: ['bind_info'],
                uneditableProperties: []
            }
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
            processTableProps () {
                return {
                    hostId: this.hostId,
                    properties: this.properties,
                    list: this.processForm.instance ? this.processForm.instance['bind_info'] : []
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
                const [properties, propertyGroups, propertyUnique] = await Promise.all([
                    this.getProcessProperties(),
                    this.getProcessPropertyGroups(),
                    this.getProcessPropertyUnique()
                ])
                this.properties = properties
                this.propertyGroups = propertyGroups
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
            getProcessPropertyGroups () {
                const action = 'objectModelFieldGroup/searchGroup'
                return this.$store.dispatch(action, {
                    objId: 'process',
                    params: {},
                    config: {
                        requestId: 'get_service_process_property_groups',
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
                this.processForm.type = 'batch'
                this.processForm.title = this.$t('批量编辑')
                this.processForm.instance = {}
                this.processForm.show = true
                // 绑定信息暂不支持批量编辑
                this.uneditableProperties.push('bind_info')
                this.$nextTick(() => {
                    const { processForm } = this.$refs
                    this.processForm.unwatch = processForm.$watch(() => {
                        return processForm.values.bk_func_name
                    }, (newVal, oldValue) => {
                        if (processForm.values.bk_process_name === oldValue) {
                            processForm.values.bk_process_name = newVal
                        }
                    })
                })
            },
            handleEditProcess (item) {
                this.processForm.type = 'single'
                this.processForm.title = `${this.$t('编辑进程')}${item.bk_process_name}`
                this.processForm.instance = item
                this.processForm.show = true
                const index = this.uneditableProperties.findIndex(property => 'bind_info')
                if (index !== -1) {
                    this.uneditableProperties.splice(index, 1)
                }
                this.$nextTick(() => {
                    const { processForm } = this.$refs
                    this.processForm.unwatch = processForm.$watch(() => {
                        return processForm.values.bk_func_name
                    }, (newVal, oldValue) => {
                        if (processForm.values.bk_process_name === oldValue) {
                            processForm.values.bk_process_name = newVal
                        }
                    })
                })
            },
            handleSubmit (values, changedValues) {
                if (this.processForm.type === 'single') {
                    Object.assign(this.processForm.instance, changedValues)
                } else {
                    this.cloneProcesses.forEach(instance => {
                        Object.assign(instance, changedValues)
                    })
                }
                this.processForm.show = false
            },
            handleCloseProcessForm () {
                this.processForm.show = false
                this.processForm.instance = {}
            },
            handleBeforeClose () {
                const changedValues = this.$refs.processForm.changedValues
                if (Object.keys(changedValues).length) {
                    return new Promise((resolve, reject) => {
                        this.$bkInfo({
                            title: this.$t('确认退出'),
                            subTitle: this.$t('退出会导致未保存信息丢失'),
                            extCls: 'bk-dialog-sub-header-center',
                            confirmFn: () => {
                                this.handleCloseProcessForm()
                            },
                            cancelFn: () => {
                                resolve(false)
                            }
                        })
                    })
                }
                this.handleCloseProcessForm()
            },
            async doClone () {
                try {
                    await this.$store.dispatch('serviceInstance/createProcServiceInstanceWithRaw', {
                        params: this.$injectMetadata({
                            name: this.$parent.module.bk_module_name,
                            bk_module_id: this.$route.params.moduleId,
                            instances: [
                                {
                                    bk_host_id: this.$route.params.hostId,
                                    processes: this.getCloneProcessValues()
                                }
                            ]
                        }, { injectBizId: true })
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
            },
            handleProcessTableChange (list) {
                this.$refs.processForm.values.bind_info = list.map(item => {
                    const row = {}
                    Object.keys(item).forEach(key => {
                        row[key] = item[key].value
                    })
                    return row
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
