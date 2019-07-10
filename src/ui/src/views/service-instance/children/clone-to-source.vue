<template>
    <div class="source-layout">
        <p class="source-tips">{{$t('BusinessTopology["源实例主机提示"]')}}</p>
        <div class="table-options">
            <bk-button class="options-button"
                :disabled="!checked.length"
                @click="handleBatchEdit">
                {{$t('Common["批量编辑"]')}}
            </bk-button>
        </div>
        <bk-table class="source-table"
            :data="flattenList"
            @selection-change="handleSelectChange">
            <bk-table-column type="selection" align="center" width="60" fixed></bk-table-column>
            <bk-table-column v-for="column in header"
                :key="column.id"
                :prop="column.id"
                :label="column.name">
            </bk-table-column>
            <bk-table-column :label="$t('Common[\'操作\']')" fixed="right">
                <template slot-scope="{ row }">
                    <button class="text-primary mr10" v-if="isRepeat(row)"
                        @click="handleEditProcess(row)">
                        <i class="bk-icon icon-exclamation-circle"></i>
                        {{$t('Common["请编辑"]')}}
                    </button>
                    <button class="text-primary mr10" v-else
                        @click="handleEditProcess(row)">
                        {{$t('Common["编辑"]')}}
                    </button>
                </template>
            </bk-table-column>
        </bk-table>
        <div class="page-options">
            <bk-button class="options-button" theme="primary"
                :disabled="!!repeatedProcesses.length"
                @click="doClone">
                {{$t('Common["确定"]')}}
            </bk-button>
            <bk-button class="options-button" @click="backToModule">{{$t('Common["取消"]')}}</bk-button>
        </div>
        <bk-sideslider
            :is-show.sync="processForm.show"
            :title="processForm.title"
            :width="800"
            :before-close="handleBeforeClose">
            <cmdb-form slot="content"
                ref="processForm"
                v-if="processForm.show"
                :properties="properties"
                :property-groups="propertyGroups"
                :object-unique="processForm.type === 'single' ? [] : propertyUnique"
                :inst="processForm.instance"
                @on-submit="handleSubmit"
                @on-cancel="handleBeforeClose">
            </cmdb-form>
        </bk-sideslider>
    </div>
</template>

<script>
    export default {
        name: 'clone-to-source',
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
                }
            }
        },
        computed: {
            header () {
                const display = [
                    'bk_func_name',
                    'bk_process_name',
                    'bk_start_param_regex',
                    'bind_ip',
                    'port',
                    'work_path'
                ]
                const header = display.map(id => {
                    const property = this.properties.find(property => property.bk_property_id === id) || {}
                    return {
                        id: property.bk_property_id,
                        name: property.bk_property_name
                    }
                })
                return header
            },
            flattenList () {
                return this.$tools.flattenList(this.properties, this.cloneProcesses)
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
            }
        },
        watch: {
            sourceProcesses (source) {
                this.cloneProcesses = this.$tools.clone(source)
            }
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
                this.processForm.title = this.$t('Common["批量编辑"]')
                this.processForm.instance = {}
                this.processForm.show = true
            },
            handleEditProcess (item) {
                this.processForm.type = 'single'
                this.processForm.title = `${this.$t('BusinessTopology["编辑进程"]')}${item.bk_process_name}`
                this.processForm.instance = this.cloneProcesses.find(target => target.bk_process_id === item.bk_process_id)
                this.processForm.show = true
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
                            title: this.$t('Common["退出会导致未保存信息丢失，是否确认？"]'),
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
                        })
                    })
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
                this.$router.replace({
                    name: 'topology',
                    query: {
                        module: this.$route.params.moduleId
                    }
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
    }
    .page-options {
        margin: 30px 0 0 0;
    }
    .source-table {
        margin: 10px 0 0 0;
    }
    .text-primary {
        .icon-exclamation-circle {
            @include inlineBlock(-1px);
            color: #FFB848;
        }
    }
</style>
