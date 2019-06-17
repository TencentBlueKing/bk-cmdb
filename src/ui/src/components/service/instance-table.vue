<template>
    <div class="service-table-layout">
        <div class="title">
            <div class="fl" @click="localExpanded = !localExpanded">
                <i class="bk-icon icon-down-shape" v-if="localExpanded"></i>
                <i class="bk-icon icon-right-shape" v-else></i>
                {{name}}
            </div>
            <div class="fr">
                <i class="bk-icon icon-close" v-if="deletable" @click="handleDelete"></i>
            </div>
        </div>
        <cmdb-table
            :header="header"
            :list="processFlattenList"
            :empty-height="58"
            :sortable="false"
            :reference-document-height="false">
            <template slot="data-empty">
                <button class="add-process-button text-primary" @click="handleAddProcess">
                    <i class="bk-icon icon-plus"></i>
                    <span>{{$t('BusinessTopology["添加进程"]')}}</span>
                </button>
            </template>
            <template slot="__operation__" slot-scope="{ rowIndex }">
                <a href="javascript:void(0)" class="text-primary" @click="handleEditProcess(rowIndex)">{{$t('Common["编辑"]')}}</a>
                <a href="javascript:void(0)" class="text-primary" v-if="!sourceProcesses.length"
                    @click="handleDeleteProcess(rowIndex)">{{$t('Common["删除"]')}}
                </a>
            </template>
        </cmdb-table>
        <div class="add-process-options" v-if="!sourceProcesses.length && processList.length">
            <button class="add-process-button text-primary" @click="handleAddProcess">
                <i class="bk-icon icon-plus"></i>
                <span>{{$t('BusinessTopology["添加进程"]')}}</span>
            </button>
        </div>
        <cmdb-slider
            :title="`${$t('BusinessTopology[\'添加进程\']')}(${name})`"
            :is-show.sync="processForm.show">
            <cmdb-form slot="content"
                ref="processForm"
                :type="processForm.type"
                :inst="processForm.instance"
                :properties="processProperties"
                :property-groups="processPropertyGroups"
                :uneditable-properties="immutableProperties"
                @on-submit="handleSaveProcess"
                @on-cancel="handleCancelCreateProcess">
            </cmdb-form>
        </cmdb-slider>
    </div>
</template>

<script>
    export default {
        props: {
            deletable: Boolean,
            expanded: Boolean,
            id: {
                type: Number,
                required: true
            },
            index: {
                type: Number,
                required: true
            },
            name: {
                type: String,
                default: ''
            },
            sourceProcesses: {
                type: Array,
                default () {
                    return []
                }
            },
            templates: {
                type: Array,
                default () {
                    return []
                }
            }
        },
        data () {
            return {
                localExpanded: this.expanded,
                processList: this.$tools.clone(this.sourceProcesses),
                processProperties: [],
                processPropertyGroups: [],
                processForm: {
                    show: false,
                    type: 'create',
                    rowIndex: null,
                    instance: {},
                    unwatch: null
                }
            }
        },
        computed: {
            header () {
                const display = [
                    'bk_process_name',
                    'bind_ip',
                    'port',
                    'work_path',
                    'user'
                ]
                const header = []
                display.map(id => {
                    const property = this.processProperties.find(property => property.bk_property_id === id)
                    if (property) {
                        header.push({
                            id: property.bk_property_id,
                            name: property.bk_property_name
                        })
                    }
                })
                header.push({
                    id: '__operation__',
                    name: this.$t('Common["操作"]')
                })
                return header
            },
            processFlattenList () {
                return this.$tools.flattenList(this.processProperties, this.processList)
            },
            immutableProperties () {
                const properties = []
                if (this.processForm.rowIndex !== null && this.templates.length) {
                    const template = this.templates[this.processForm.rowIndex]
                    Object.keys(template.property).forEach(key => {
                        if (template.property[key].as_default_value) {
                            properties.push(key)
                        }
                    })
                }
                return properties
            }
        },
        created () {
            this.getProcessProperties()
            this.getProcessPropertyGroups()
        },
        methods: {
            async getProcessProperties () {
                try {
                    const action = 'objectModelProperty/searchObjectAttribute'
                    this.processProperties = await this.$store.dispatch(action, {
                        params: {
                            bk_obj_id: 'process',
                            bk_supplier_account: this.$store.getters.supplierAccount
                        },
                        config: {
                            requestId: 'get_service_process_properties',
                            fromCache: true
                        }
                    })
                } catch (e) {
                    console.error(e)
                }
            },
            async getProcessPropertyGroups () {
                try {
                    const action = 'objectModelFieldGroup/searchGroup'
                    this.processPropertyGroups = await this.$store.dispatch(action, {
                        objId: 'process',
                        params: {},
                        config: {
                            requestId: 'get_service_process_property_groups',
                            fromCache: true
                        }
                    })
                } catch (e) {
                    console.error(e)
                }
            },
            handleDelete () {
                this.$emit('delete-instance', this.index)
            },
            handleAddProcess () {
                this.processForm.instance = {}
                this.processForm.type = 'create'
                this.processForm.show = true
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
            handleSaveProcess (values) {
                this.processForm.unwatch && this.processForm.unwatch()
                if (this.processForm.type === 'create') {
                    this.processList.push(values)
                } else {
                    Object.assign(this.processForm.instance, values)
                }
                this.handleCancelCreateProcess()
            },
            handleCancelCreateProcess () {
                this.processForm.show = false
                this.processForm.rowIndex = null
            },
            handleEditProcess (rowIndex) {
                this.processForm.instance = this.processList[rowIndex]
                this.processForm.rowIndex = rowIndex
                this.processForm.type = 'update'
                this.processForm.show = true
            },
            handleDeleteProcess (rowIndex) {
                this.processList.splice(rowIndex, 1)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .title {
        height: 40px;
        padding: 0 10px;
        line-height: 40px;
        border-radius: 2px 2px 0 0;
        background-color: #DCDEE5;
        .bk-icon {
            font-size: 12px;
            font-weight: bold;
            width: 24px;
            height: 24px;
            line-height: 24px;
            text-align: center;
            cursor: pointer;
            @include inlineBlock;
        }
    }
    .add-process-options {
        border: 1px solid $cmdbTableBorderColor;
        border-top: none;
        line-height: 42px;
        font-size: 12px;
        text-align: center;
    }
    .add-process-button {
        line-height: 32px;
        .bk-icon,
        span {
            @include inlineBlock;
        }
    }
</style>
