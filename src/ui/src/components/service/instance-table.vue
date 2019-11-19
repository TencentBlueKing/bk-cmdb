<template>
    <div class="service-table-layout">
        <div class="title" @click="localExpanded = !localExpanded">
            <div class="fl">
                <i class="bk-icon icon-down-shape" v-if="localExpanded"></i>
                <i class="bk-icon icon-right-shape" v-else></i>
                {{name}}
            </div>
            <div class="fr">
                <span v-if="topology" class="service-topology">{{topology}}</span>
                <i class="bk-icon icon-close" v-if="deletable" @click.stop="handleDelete"></i>
            </div>
        </div>
        <bk-table
            v-show="localExpanded"
            :data="processFlattenList">
            <bk-table-column v-for="column in header"
                :key="column.id"
                :prop="column.id"
                :label="column.name">
            </bk-table-column>
            <bk-table-column :label="$t('操作')" fixed="right">
                <template slot-scope="{ row, $index }">
                    <a href="javascript:void(0)" class="text-primary mr10" @click="handleEditProcess($index)">
                        {{$t('编辑')}}
                    </a>
                    <a href="javascript:void(0)" class="text-primary"
                        v-if="!sourceProcesses.length"
                        @click="handleDeleteProcess($index)">
                        {{$t('删除')}}
                    </a>
                </template>
            </bk-table-column>
            <template slot="empty" v-if="addible">
                <button class="add-process-button text-primary" @click="handleAddProcess">
                    <i class="bk-icon icon-plus"></i>
                    <span>{{$t('添加进程')}}</span>
                </button>
            </template>
        </bk-table>
        <div class="add-process-options" v-if="addible && !sourceProcesses.length && processList.length">
            <button class="add-process-button text-primary" @click="handleAddProcess">
                <i class="bk-icon icon-plus"></i>
                <span>{{$t('添加进程')}}</span>
            </button>
        </div>
        <bk-sideslider
            v-transfer-dom
            :width="800"
            :title="`${$t('添加进程')}(${name})`"
            :is-show.sync="processForm.show"
            :before-close="handleBeforeClose">
            <cmdb-form slot="content" v-if="processForm.show"
                ref="processForm"
                :type="processForm.type"
                :inst="processForm.instance"
                :properties="processProperties"
                :property-groups="processPropertyGroups"
                :disabled-properties="immutableProperties"
                @on-submit="handleSaveProcess"
                @on-cancel="handleBeforeClose">
                <template slot="bind_ip">
                    <cmdb-input-select
                        :disabled="checkDisabled"
                        :name="'bindIp'"
                        :placeholder="$t('请选择或输入IP')"
                        :options="processBindIp"
                        :validate="validateRules"
                        v-model="bindIp">
                    </cmdb-input-select>
                </template>
            </cmdb-form>
        </bk-sideslider>
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
            },
            addible: {
                type: Boolean,
                default: true
            },
            topology: {
                type: String,
                default: ''
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
                },
                tooltips: {
                    content: this.$t('请为主机添加进程'),
                    placement: 'right'
                },
                processBindIp: [],
                bindIp: ''
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
                const header = []
                display.forEach(id => {
                    const property = this.processProperties.find(property => property.bk_property_id === id)
                    if (property) {
                        header.push({
                            id: property.bk_property_id,
                            name: property.bk_property_name
                        })
                    }
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
            },
            bindIpProperty () {
                return this.processProperties.find(property => property['bk_property_id'] === 'bind_ip') || {}
            },
            validateRules () {
                const rules = {}
                if (this.bindIpProperty.isrequired) {
                    rules['required'] = true
                }
                rules['regex'] = this.bindIpProperty.option
                return rules
            },
            checkDisabled () {
                const property = this.bindIpProperty
                if (this.processForm.type === 'create') {
                    return false
                }
                return !property.editable || property.isreadonly || this.immutableProperties.includes('bind_ip')
            }
        },
        watch: {
            bindIp (value) {
                this.$refs.processForm.values.bind_ip = value
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
                this.getInstanceIpByHost(this.id)
                this.processForm.instance = {}
                this.processForm.type = 'create'
                this.processForm.show = true
                this.$nextTick(() => {
                    this.bindIp = ''
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
            async getInstanceIpByHost (hostId) {
                try {
                    const instanceIpMap = this.$store.state.businessHost.instanceIpMap
                    let res = null
                    if (instanceIpMap.hasOwnProperty(hostId)) {
                        res = instanceIpMap[hostId]
                    } else {
                        res = await this.$store.dispatch('serviceInstance/getInstanceIpByHost', {
                            hostId,
                            config: {
                                requestId: 'getInstanceIpByHost'
                            }
                        })
                        this.$store.commit('businessHost/setInstanceIp', { hostId, res })
                    }
                    this.processBindIp = res.options.map(ip => {
                        return {
                            id: ip,
                            name: ip
                        }
                    })
                } catch (e) {
                    this.processBindIp = []
                    console.error(e)
                }
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
            handleBeforeClose () {
                const changedValues = this.$refs.processForm.changedValues
                if (Object.keys(changedValues).length) {
                    return new Promise((resolve, reject) => {
                        this.$bkInfo({
                            title: this.$t('确认退出'),
                            subTitle: this.$t('退出会导致未保存信息丢失'),
                            extCls: 'bk-dialog-sub-header-center',
                            confirmFn: () => {
                                this.handleCancelCreateProcess()
                            },
                            cancelFn: () => {
                                resolve(false)
                            }
                        })
                    })
                }
                this.handleCancelCreateProcess()
            },
            handleEditProcess (rowIndex) {
                this.getInstanceIpByHost(this.id)
                this.processForm.instance = this.processList[rowIndex]
                this.processForm.rowIndex = rowIndex
                this.processForm.type = 'update'
                this.processForm.show = true

                this.$nextTick(() => {
                    this.bindIp = this.$tools.getInstFormValues(this.processProperties, this.processForm.instance)['bind_ip']
                })
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
        cursor: pointer;
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
        .icon-exclamation {
            font-size: 14px;
            color: #ffffff;
            background: #f0b659;
            border-radius: 50%;
            transform: scale(.6);
        }
        .service-topology {
            padding: 0 5px;
            line-height: 40px;
            font-size: 12px;
            color: $textColor;
            cursor: default;
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
