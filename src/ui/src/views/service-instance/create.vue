<template>
    <div class="create-layout clearfix" v-bkloading="{ isLoading: $loading(Object.values(request)) }">
        <div class="info clearfix mb20">
            <label class="info-label fl">{{$t('已选主机')}}：</label>
            <div class="info-content">
                <bk-button class="select-host-button" theme="default"
                    @click="handleSelectHost">
                    <i class="bk-icon icon-plus"></i>
                    {{$t('添加主机')}}
                </bk-button>
                <i18n class="select-host-count" path="已选择N台主机" v-show="hosts.length">
                    <span place="count" class="count-number">{{hosts.length}}</span>
                </i18n>
            </div>
        </div>
        <div class="info clearfix mb10" ref="changeInfo" v-show="resources.length">
            <label class="info-label fl">{{$t('变更确认')}}：</label>
            <div class="info-content">
                <template v-if="availableTabList.length">
                    <ul class="tab clearfix">
                        <template v-for="(item, index) in availableTabList">
                            <li class="tab-grep fl" v-if="index" :key="index"></li>
                            <li class="tab-item fl"
                                :class="{ active: activeTab === item }"
                                :key="item.id"
                                @click="handleTabClick(item)">
                                <span class="tab-label">{{item.label}}</span>
                                <span :class="['tab-count', { 'unconfirmed': !item.confirmed }]">
                                    {{item.props.count > 999 ? '999+' : item.props.count}}
                                </span>
                            </li>
                        </template>
                    </ul>
                    <div class="tab-component" v-show="activeTab.id === 'createServiceInstance'">
                        <transition-group name="service-table-list" tag="div">
                            <service-instance-table class="service-instance-table"
                                v-for="(instance, index) in instances"
                                ref="serviceInstanceTable"
                                deletable
                                :key="instance.bk_host_id"
                                :index="index"
                                :id="instance.bk_host_id"
                                :name="getName(instance)"
                                :source-processes="getSourceProcesses(instance)"
                                :templates="getServiceTemplates(instance)"
                                :addible="!withTemplate"
                                :editing="getEditState(instance)"
                                :instance="instance"
                                @edit-process="handleEditProcess(instance, ...arguments)"
                                @delete-instance="handleDeleteInstance"
                                @edit-name="handleEditName(instance)"
                                @confirm-edit-name="handleConfirmEditName(instance, ...arguments)"
                                @cancel-edit-name="handleCancelEditName(instance)">
                            </service-instance-table>
                        </transition-group>
                    </div>
                    <div class="tab-component" v-show="activeTab.id === 'hostAttrsAutoApply'">
                        <property-confirm-table class="mt10"
                            ref="confirmTable"
                            max-height="auto"
                            :list="conflictList"
                            :render-icon="true"
                            :show-operation="hasConflict">
                        </property-confirm-table>
                    </div>
                </template>
            </div>
        </div>
        <div class="options" :class="{ 'is-sticky': hasScrollbar }">
            <cmdb-auth class="mr5" :auth="{ type: $OPERATION.C_SERVICE_INSTANCE, relation: [bizId] }">
                <bk-button slot-scope="{ disabled }"
                    theme="primary"
                    :disabled="!hosts.length || disabled"
                    @click="handleConfirm">
                    {{$t('确定')}}
                </bk-button>
            </cmdb-auth>
            <bk-button @click="handleBackToModule">{{$t('取消')}}</bk-button>
        </div>
        <cmdb-dialog v-model="dialog.show" v-bind="dialog.props">
            <component
                :is="dialog.component"
                :confirm-text="$t('确定')"
                v-bind="dialog.componentProps"
                @confirm="handleDialogConfirm"
                @cancel="handleDialogCancel">
            </component>
        </cmdb-dialog>
    </div>
</template>

<script>
    import HostSelector from '@/views/business-topology/host/host-selector-new'
    import serviceInstanceTable from '@/components/service/instance-table.vue'
    import propertyConfirmTable from '@/components/host-apply/property-confirm-table'
    import { mapGetters } from 'vuex'
    import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
    export default {
        name: 'create-service-instance',
        components: {
            [HostSelector.name]: HostSelector,
            serviceInstanceTable,
            propertyConfirmTable
        },
        data () {
            return {
                hasScrollbar: false,
                hosts: [],
                instances: [],
                resources: [],
                confirmParams: {},
                conflictList: [],
                tab: {
                    active: null
                },
                tabList: [{
                    id: 'createServiceInstance',
                    label: this.$t('新增服务实例'),
                    confirmed: false,
                    props: {
                        count: 0
                    }
                }, {
                    id: 'hostAttrsAutoApply',
                    label: this.$t('属性自动应用'),
                    confirmed: false,
                    props: {
                        count: 0
                    }
                }],
                dialog: {
                    show: false,
                    props: {
                        width: 1110,
                        height: 650,
                        showCloseIcon: false
                    },
                    component: null,
                    componentProps: {}
                },
                request: {
                    preview: Symbol('review'),
                    hostInfo: Symbol('hostInfo')
                },
                processChangeState: {}
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('businessHost', ['getDefaultSearchCondition']),
            moduleId () {
                return parseInt(this.$route.params.moduleId)
            },
            setId () {
                return parseInt(this.$route.params.setId)
            },
            withTemplate () {
                const instance = this.instances[0] || {}
                if (instance.service_template && instance.service_template.service_template.id) {
                    return true
                }
                return false
            },
            availableTabList () {
                return this.tabList.filter(tab => tab.props.count > 0)
            },
            activeTab () {
                return this.tabList.find(tab => tab.id === this.tab.active) || this.availableTabList[0]
            },
            hasConflict () {
                const conflictList = this.conflictList.filter(item => item.unresolved_conflict_count > 0)
                return conflictList.length > 0
            }
        },
        watch: {
            activeTab (tab) {
                if (!tab) return
                tab.confirmed = true
            }
        },
        async created () {
            this.resolveData(this.$route)
            if (this.resources.length) {
                this.getSelectedHost()
                this.getPreviewData()
            }
        },
        mounted () {
            addResizeListener(this.$refs.changeInfo, this.resizeHandler)
        },
        beforeDestroy () {
            removeResizeListener(this.$refs.changeInfo, this.resizeHandler)
        },
        async beforeRouteUpdate (to, from, next) {
            this.resolveData(to)
            if (this.resources.length) {
                this.getSelectedHost()
                this.getPreviewData()
            }
            next()
        },
        methods: {
            resolveData (route) {
                const query = route.query || {}
                const resources = query.resources
                if (!resources) {
                    this.resources = []
                } else {
                    this.resources = String(resources).split(',').map(id => Number(id))
                }
                this.confirmParams = {
                    bk_host_ids: this.resources,
                    bk_module_id: this.moduleId
                }
            },
            async getSelectedHost () {
                try {
                    this.$store.commit('setGlobalLoading', this.hasScrollbar)
                    const result = await this.getHostInfo()
                    this.hosts = result.info || []
                } catch (e) {
                    console.error(e)
                } finally {
                    this.$store.commit('setGlobalLoading', false)
                }
            },
            async getPreviewData () {
                try {
                    this.$store.commit('setGlobalLoading', this.hasScrollbar)
                    const data = await this.$store.dispatch('serviceInstance/createProcServiceInstancePreview', {
                        params: {
                            ...this.confirmParams,
                            bk_biz_id: this.bizId
                        },
                        config: {
                            requestId: this.request.preview,
                            globalPermission: false
                        }
                    })
                    this.setCreateServiceInstance(data)
                    this.setHostAttrsAutoApply(data)
                } catch (e) {
                    console.error(e)
                    if (e.code === 9900403) {
                        this.$route.meta.view = 'permission'
                    }
                } finally {
                    this.$store.commit('setGlobalLoading', false)
                }
            },
            getHostInfo () {
                const params = {
                    bk_biz_id: this.bk_biz_id,
                    ip: { data: [], exact: 0, flag: 'bk_host_innerip|bk_host_outerip' },
                    page: {},
                    condition: this.getDefaultSearchCondition()
                }
                const hostCondition = params.condition.find(target => target.bk_obj_id === 'host')
                hostCondition.condition.push({
                    field: 'bk_host_id',
                    operator: '$in',
                    value: this.resources
                })
                return this.$store.dispatch('hostSearch/searchHost', {
                    params,
                    config: {
                        requestId: this.request.hostInfo
                    }
                })
            },
            setCreateServiceInstance (data) {
                const instanceInfo = []
                data.forEach(item => {
                    item.to_add_to_modules.forEach(moduleInfo => {
                        instanceInfo.push({
                            bk_host_id: item.bk_host_id,
                            name: '',
                            editing: { name: false },
                            ...moduleInfo
                        })
                    })
                })
                const tab = this.tabList.find(tab => tab.id === 'createServiceInstance')
                this.instances = instanceInfo
                tab.props.count = this.instances.length
            },
            setHostAttrsAutoApply (data) {
                const conflictInfo = (data || []).map(item => item.host_apply_plan)
                this.conflictList = conflictInfo.filter(item => item.conflicts.length || item.update_fields.length)
                const tab = this.tabList.find(tab => tab.id === 'hostAttrsAutoApply')
                tab.props.count = this.conflictList.length
            },
            getName (instance) {
                if (instance.name) {
                    return instance.name
                }
                const data = this.hosts.find(data => data.host.bk_host_id === instance.bk_host_id)
                if (data) {
                    return data.host.bk_host_innerip
                }
                return '--'
            },
            getServiceTemplates (instance) {
                if (instance.service_template) {
                    return instance.service_template.process_templates
                }
                return []
            },
            getSourceProcesses (instance) {
                const templates = this.getServiceTemplates(instance)
                return templates.map(template => {
                    const value = {}
                    Object.keys(template.property).forEach(key => {
                        if (key === 'bind_info') {
                            value[key] = this.$tools.clone(template.property[key].value) || []
                            value[key].forEach(row => {
                                Object.keys(row).forEach(field => {
                                    if (field === 'ip') {
                                        row[field] = this.getBindIp(instance, row)
                                    } else if (field === 'row_id') {
                                        // 实例数据中使用 template_row_id
                                        row['template_row_id'] = row[field]
                                        delete row[field]
                                    } else if (row[field] !== null && typeof row[field] === 'object') {
                                        row[field] = row[field].value
                                    }
                                })
                            })
                        } else {
                            value[key] = template.property[key].value
                        }
                    })
                    return value
                })
            },
            getBindIp (instance, row) {
                const ipValue = row.ip.value
                const mapping = {
                    1: '127.0.0.1',
                    2: '0.0.0.0'
                }
                if (mapping.hasOwnProperty(ipValue)) {
                    return mapping[ipValue]
                }
                const { host } = this.hosts.find(data => data.host.bk_host_id === instance.bk_host_id)
                // 第一内网IP
                if (ipValue === '3') {
                    const [innerIP] = host.bk_host_innerip.split(',')
                    return innerIP || mapping[1]
                }
                const [outerIP] = host.bk_host_outerip.split(',')
                return outerIP || mapping[1]
            },
            getEditState (instance) {
                return instance.editing
            },
            handleSelectHost () {
                this.dialog.componentProps.exist = this.hosts
                this.dialog.component = HostSelector.name
                this.dialog.show = true
            },
            handleDialogConfirm (selected) {
                this.dialog.show = false
                this.redirect({
                    query: {
                        ...this.$route.query,
                        resources: selected.map(data => data.host.bk_host_id).join(',')
                    }
                })
            },
            handleDialogCancel () {
                this.dialog.show = false
            },
            handleDeleteInstance (index) {
                this.instances.splice(index, 1)
                this.redirect({
                    query: {
                        ...this.$route.query,
                        resources: this.instances.map(data => data.bk_host_id).join(',')
                    }
                })
            },
            /**
             * 解决后端性能问题: 用服务模板生成的实例仅传递有被用户主动触发过编辑的进程信息
             */
            getChangedProcessList (instance, component) {
                if (this.withTemplate) {
                    const processes = []
                    const stateKey = `${instance.bk_module_id}-${instance.bk_host_id}`
                    const changedState = this.processChangeState[stateKey] || new Set()
                    component.processList.forEach((process, listIndex) => {
                        if (!changedState.has(listIndex)) return
                        processes.push({
                            process_template_id: component.templates[listIndex] ? component.templates[listIndex].id : 0,
                            process_info: process
                        })
                    })
                    return processes
                }
                return component.processList.map((process, listIndex) => ({
                    process_template_id: component.templates[listIndex] ? component.templates[listIndex].id : 0,
                    process_info: process
                }))
            },
            /**
             * 解决后端性能问题: 记录用服务模板生成的实例是否触发编辑动作
             */
            handleEditProcess (instance, processIndex) {
                if (!instance.service_template) return
                const key = `${instance.bk_module_id}-${instance.bk_host_id}`
                const state = this.processChangeState[key] || new Set()
                state.add(processIndex)
                this.processChangeState[key] = state
            },
            async handleConfirm () {
                try {
                    const serviceInstanceTables = this.$refs.serviceInstanceTable
                    const confirmTable = this.$refs.confirmTable
                    const params = {
                        bk_module_id: this.moduleId,
                        bk_biz_id: this.bizId
                    }
                    if (serviceInstanceTables) {
                        params.instances = serviceInstanceTables.map(table => {
                            const instance = this.instances.find(instance => instance.bk_host_id === table.id) || {}
                            return {
                                bk_host_id: table.id,
                                service_instance_name: instance.name || '',
                                processes: this.getChangedProcessList(instance, table)
                            }
                        })
                    }
                    if (confirmTable) {
                        params.host_apply_conflict_resolvers = this.getHostApplyConflictResolvers()
                    }

                    await this.$store.dispatch('serviceInstance/createProcServiceInstanceByTemplate', {
                        params: params
                    })

                    this.$success(this.$t('添加成功'))
                    this.handleBackToModule()
                } catch (e) {
                    console.error(e)
                }
            },
            getHostApplyConflictResolvers () {
                const conflictResolveResult = this.$refs.confirmTable.conflictResolveResult
                const conflictResolvers = []
                Object.keys(conflictResolveResult).forEach(key => {
                    const propertyList = conflictResolveResult[key]
                    propertyList.forEach(property => {
                        conflictResolvers.push({
                            bk_host_id: Number(key),
                            bk_attribute_id: property.id,
                            bk_property_value: property.__extra__.value
                        })
                    })
                })
                return conflictResolvers
            },
            handleEditName (instance) {
                this.instances.forEach(instance => (instance.editing.name = false))
                instance.editing.name = true
            },
            handleConfirmEditName (instance, name) {
                instance.name = name
                instance.editing.name = false
            },
            handleCancelEditName (instance) {
                instance.editing.name = false
            },
            handleBackToModule () {
                this.$routerActions.back()
            },
            handleTabClick (tab) {
                this.tab.active = tab.id
            },
            resizeHandler () {
                this.$nextTick(() => {
                    const scroller = this.$el.parentElement
                    this.hasScrollbar = scroller.scrollHeight > scroller.offsetHeight
                })
            },
            redirect (options = {}) {
                const config = {
                    name: 'createServiceInstance',
                    params: this.$route.params,
                    query: this.$route.query
                }
                this.$routerActions.redirect({ ...config, ...options })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .create-layout {
        padding: 15px 0 0 0;
    }

    .info {
        .info-label {
            width: 128px;
            font-size: 14px;
            font-weight: bold;
            color: $textColor;
            text-align: right;
            padding-top: 8px;
        }
        .info-content {
            overflow: hidden;
            padding: 8px 20px 0 8px;
            font-size: 14px;
            .info-count {
                font-weight: bold;
            }
            .module-grep {
                border-top: 1px solid $borderColor;
                margin-top: 10px;
            }
            .edit-trigger {
                @include inlineBlock;
                margin-left: 10px;
                color: $primaryColor;
                cursor: pointer;
                &:hover {
                    color: #1964E1;
                }
            }
        }
    }

    .tab {
        .tab-grep {
            width: 2px;
            height: 19px;
            margin: 0 15px;
            background-color: #C4C6CC;
        }
        .tab-item {
            position: relative;
            color: $textColor;
            font-size: 0;
            cursor: pointer;
            &.active {
                color: $primaryColor;
            }
            &.active:after {
                content: "";
                position: absolute;
                left: 0;
                top: 30px;
                width: 100%;
                height: 2px;
                background-color: $primaryColor;
            }
            .tab-label {
                display: inline-block;
                vertical-align: middle;
                margin-left: 10px;
                margin-right: 4px;
                font-size: 14px;
            }
            .tab-count {
                display: inline-block;
                vertical-align: middle;
                height: 16px;
                padding: 0 5px;
                border-radius: 8px;
                line-height: 14px;
                font-size: 12px;
                color: #FFF;
                background-color: #C4C6CC;
                text-align: center;
                border: 1px solid #fff;

                &.unconfirmed {
                    background-color: #FF5656;
                }
            }
        }
    }
    .tab-component {
        margin-top: 20px;
    }
    .tab-empty {
        height: 60px;
        padding: 0 28px;
        line-height: 60px;
        background-color: #F0F1F5;
        color: $textColor;
        &:before {
            content: "!";
            display: inline-block;
            width: 16px;
            height: 16px;
            line-height: 16px;
            border-radius: 50%;
            text-align: center;
            color: #FFF;
            font-size: 12px;
            background-color: #C4C6CC;
        }
    }

    .select-host-button {
        height: 32px;
        line-height: 30px;
        font-size: 0;
        .bk-icon {
            position: static;
            height: 30px;
            line-height: 30px;
            font-size: 20px;
            font-weight: bold;
            margin: 0 -4px;
            @include inlineBlock(top);
        }
        /deep/ span {
            font-size: 14px;
        }
    }
    .select-host-count {
        color: $textColor;
        .count-number {
            font-weight: bold;
        }
    }

    .service-instance-table + .service-instance-table {
        margin-top: 12px;
    }

    .options {
        position: sticky;
        padding: 10px 0 10px 136px;
        font-size: 0;
        bottom: 0;
        left: 0;
        &.is-sticky {
            background-color: #FFF;
            border-top: 1px solid $borderColor;
            z-index: 100;
        }
    }

    .service-table-list, .service-table-list-leave-active {
        transition: all .7s ease-in;
    }
    .service-table-list-leave-to {
        opacity: 0;
        transform: translateX(30px);
    }
</style>
