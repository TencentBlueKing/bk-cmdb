<template>
    <div class="table-layout" v-show="show">
        <div class="table-title" @click="localExpanded = !localExpanded">
            <cmdb-form-bool class="title-checkbox"
                :size="16"
                v-model="checked"
                @click.native.stop>
            </cmdb-form-bool>
            <i class="title-icon bk-icon icon-down-shape" v-if="localExpanded"></i>
            <i class="title-icon bk-icon icon-right-shape" v-else></i>
            <span class="title-label">{{instance.name}}</span>
            <cmdb-dot-menu class="instance-menu" @click.native.stop>
                <ul class="menu-list">
                    <li class="menu-item"
                        v-for="(menu, index) in instanceMenu"
                        :key="index">
                        <button class="menu-button"
                            @click="menu.handler">
                            {{menu.name}}
                        </button>
                    </li>
                </ul>
            </cmdb-dot-menu>
            <div class="instance-label clearfix" @click.stop v-if="labelShowList.length">
                <div class="label-title fl">
                    <i class="icon-cc-label"></i>
                    <span>{{$t('BusinessTopology["标签"]')}}</span>
                </div>
                <div class="label-list fl">
                    <span class="label-item" :key="index" v-for="(label, index) in labelShowList">{{`${label.key}：${label.value}`}}</span>
                    <div class="label-item label-tips"
                        ref="tipsLabelContainer"
                        v-if="labelTipsList.length"
                        @mouseenter="handleShowTipsLabel"
                        @mouseleave="handleCloseTipsLabel">
                        <span>...</span>
                        <div class="tips-label-list" ref="tipsLabel" v-click-outside="handleCloseTipsLabel" v-show="tipsLabel.show">
                            <span class="label-item" :key="index" v-for="(label, index) in labelTipsList">{{`${label.key}：${label.value}`}}</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <cmdb-table
            v-show="localExpanded"
            :loading="$loading(Object.values(requestId))"
            :header="header"
            :list="flattenList"
            :empty-height="42"
            :visible="localExpanded"
            :sortable="false"
            :reference-document-height="false">
        </cmdb-table>
    </div>
</template>

<script>
    export default {
        props: {
            instance: {
                type: Object,
                required: true
            },
            expanded: Boolean
        },
        data () {
            return {
                tipsLabel: {
                    show: false,
                    instance: null,
                    id: null
                },
                show: true,
                checked: false,
                localExpanded: this.expanded,
                properties: [],
                header: [],
                list: []
            }
        },
        computed: {
            isModuleNode () {
                const node = this.$store.state.businessTopology.selectedNode
                return node && node.data.bk_obj_id === 'module'
            },
            withTemplate () {
                return this.isModuleNode && !!this.instance.service_template_id
            },
            instanceMenu () {
                const menu = [{
                    name: this.$t('Common["删除"]'),
                    handler: this.handleDeleteInstance
                }]
                return menu
            },
            module () {
                return this.$store.state.businessTopology.selectedNodeInstance
            },
            flattenList () {
                return this.$tools.flattenList(this.properties, this.list.map(data => data.property))
            },
            requestId () {
                return {
                    processList: `get_service_instance_${this.instance.id}_processes`,
                    properties: 'get_service_process_properties',
                    deleteProcess: 'delete_service_process'
                }
            },
            labelList () {
                const list = []
                const labels = this.instance.labels
                labels && Object.keys(labels).forEach(key => {
                    list.push({
                        id: this.instance.id,
                        key: key,
                        value: labels[key]
                    })
                })
                return list
            },
            labelShowList () {
                return this.labelList.slice(0, 3)
            },
            labelTipsList () {
                return this.labelList.slice(3)
            }
        },
        watch: {
            localExpanded (expanded) {
                if (expanded) {
                    this.getServiceProcessList()
                }
            },
            checked (checked) {
                this.$emit('check-change', checked, this.instance)
            }
        },
        async created () {
            await this.getProcessProperties()
            if (this.expanded) {
                this.getServiceProcessList()
            }
        },
        methods: {
            async getProcessProperties () {
                const action = 'objectModelProperty/searchObjectAttribute'
                const properties = await this.$store.dispatch(action, {
                    params: {
                        bk_obj_id: 'process',
                        bk_supplier_account: this.$store.getters.supplierAccount
                    },
                    config: {
                        requestId: this.requestId.properties,
                        fromCache: true
                    }
                })
                this.properties = properties
                this.setHeader()
            },
            async getServiceProcessList () {
                try {
                    this.list = await this.$store.dispatch('processInstance/getServiceInstanceProcesses', {
                        params: this.$injectMetadata({
                            service_instance_id: this.instance.id
                        }),
                        config: {
                            requestId: this.requestId.processList
                        }
                    })
                } catch (e) {
                    this.list = []
                    console.error(e)
                }
            },
            setHeader () {
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
                this.header = header
            },
            handleDeleteInstance () {
                this.$bkInfo({
                    title: this.$t('BusinessTopology["确认删除实例"]'),
                    content: this.$t('BusinessTopology["即将删除实例"]', { name: this.instance.name }),
                    confirmFn: async () => {
                        try {
                            await this.$store.dispatch('serviceInstance/deleteServiceInstance', {
                                config: {
                                    data: this.$injectMetadata({
                                        service_instance_ids: [this.instance.id]
                                    }),
                                    requestId: this.requestId.deleteProcess
                                }
                            })
                            this.$emit('delete-instance', this.instance.id)
                        } catch (e) {
                            console.error(e)
                        }
                    }
                })
            },
            handleCloseTipsLabel () {
                this.tipsLabel.instance && this.tipsLabel.instance.destroy()
                this.tipsLabel.show = false
            },
            handleShowTipsLabel (event) {
                if (this.tipsLabel.show) {
                    this.handleCloseTipsLabel()
                    return
                }
                this.tipsLabel.show = true
                this.tipsLabel.instance = this.$tooltips({
                    duration: -1,
                    theme: 'light',
                    zIndex: 9999,
                    width: 290,
                    placements: ['bottom'],
                    container: this.$refs.tipsLabelContainer,
                    target: event.target
                })
                this.tipsLabel.instance.$el.append(this.$refs.tipsLabel)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .table-layout {
        padding: 0 0 12px 0;
    }
    .table-title {
        height: 40px;
        padding: 0 11px;
        line-height: 40px;
        border-radius: 2px 2px 0 0;
        background-color: #DCDEE5;
        cursor: pointer;
        .title-icon {
            font-size: 12px;
            color: #63656E;
            @include inlineBlock;
        }
        .title-label {
            font-size: 14px;
            color: #313238;
            @include inlineBlock;
        }
    }
    .menu-list {
        min-width: 74px;
        padding: 6px 0;
        .menu-item {
            .menu-button {
                display: block;
                width: 100%;
                height: 32px;
                padding: 0 13px;
                line-height: 32px;
                outline: 0;
                border: none;
                text-align: left;
                color: #63656E;
                font-size: 12px;
                background-color: #fff;
                &:hover {
                    background-color: #E1ECFF;
                    color: #3A84FF;
                }
            }
        }
    }
    .instance-label {
        display: inline-block;
        vertical-align: middle;
        font-size: 12px;
        .icon-cc-label {
            color: #979ba5;
            font-size: 16px;
        }
        .label-list {
            padding-left: 4px;
            line-height: 38px;
            .label-item {
                display: inline-block;
                height: 20px;
                line-height: 20px;
                vertical-align: middle;
                margin-right: 4px;
                padding: 0 6px;
                color: #979ba5;
                background-color: #fafbfd;
                border-radius: 2px;
            }
            .label-tips:hover {
                background-color: #f0f1f5;
            }
            .tips-label-list {
                .label-item {
                    line-height: 18px;
                    color: #979ba5;
                    border: 1px solid #dcdee5;
                    margin: 5px 2px;
                }
            }
        }
    }
</style>
