<template>
    <div class="layout" v-bkloading="{
        isLoading: $loading(Object.values(request)) || loading
    }">
        <div class="info clearfix mb20">
            <label class="info-label fl">{{$t('已选主机')}}：</label>
            <i18n tag="div" path="N台主机" class="info-content">
                <b class="info-count" place="count">{{resources.length}}</b>
            </i18n>
        </div>
        <div class="info clearfix mb20" v-if="type !== 'remove'">
            <label class="info-label fl">{{$t('转移到')}}：</label>
            <div class="info-content">
                <ul class="module-list">
                    <li class="module-item" v-for="(id, index) in targetModules"
                        :key="index"
                        :class="{
                            'is-business-module': type === 'business'
                        }"
                        v-bk-tooltips="getModulePath(id)">
                        <span class="module-icon" v-if="type === 'business'">{{$i18n.locale === 'en' ? 'M' : '模'}}</span>
                        {{getModuleName(id)}}
                        <span class="module-mask"
                            v-if="type === 'idle'"
                            @click="handleChangeModule">
                            {{$t('点击修改')}}
                        </span>
                    </li>
                    <li class="module-item is-trigger"
                        v-if="type === 'business'"
                        @click="handleChangeModule">
                        <i class="icon-cc-plus"></i>
                    </li>
                </ul>
                <div class="module-grep"></div>
            </div>
        </div>
        <div class="info clearfix mb10" ref="changeInfo">
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
                                <span :class="['tab-count', { 'has-badge': !item.confirmed }]">
                                    {{item.props.info.length}}
                                </span>
                            </li>
                        </template>
                    </ul>
                    <component class="tab-component"
                        v-for="item in availableTabList"
                        v-bind="item.props"
                        v-show="activeTab === item"
                        :ref="item.id"
                        :key="item.id"
                        :is="item.component">
                    </component>
                </template>
                <div class="tab-empty" v-else-if="isSameModule">
                    {{$t('相同模块转移提示')}}
                </div>
                <div class="tab-empty" v-else-if="isEmptyChange">
                    {{$t('无转移确认信息提示')}}
                </div>
            </div>
        </div>
        <div class="options" :class="{ 'is-sticky': hasScrollbar }" v-show="!loading">
            <bk-button theme="primary" :disabled="isSameModule" @click="handleConfrim">{{confirmText}}</bk-button>
            <bk-button class="ml10" theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
        </div>
        <cmdb-dialog v-model="dialog.show" :width="dialog.width" :height="460">
            <component
                :is="dialog.component"
                :confirm-text="$t('确定')"
                v-bind="dialog.props"
                @cancel="handleDialogCancel"
                @confirm="handleDialogConfirm">
            </component>
        </cmdb-dialog>
    </div>
</template>

<script>
    import CreateServiceInstance from './children/create-service-instance.vue'
    import DeletedServiceInstance from './children/deleted-service-instance.vue'
    import MoveToIdleHost from './children/move-to-idle-host.vue'
    import ModuleSelector from '@/views/business-topology/host/module-selector.vue'
    import HostAttrsAutoApply from './children/host-attrs-auto-apply.vue'
    import {
        MENU_BUSINESS_TRANSFER_HOST,
        MENU_BUSINESS_HOST_AND_SERVICE
    } from '@/dictionary/menu-symbol'
    import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
    import { mapGetters } from 'vuex'
    export default {
        components: {
            [CreateServiceInstance.name]: CreateServiceInstance,
            [DeletedServiceInstance.name]: DeletedServiceInstance,
            [MoveToIdleHost.name]: MoveToIdleHost,
            [ModuleSelector.name]: ModuleSelector,
            [HostAttrsAutoApply.name]: HostAttrsAutoApply
        },
        data () {
            return {
                hasScrollbar: false,
                hostCount: 108,
                hostInfo: [],
                tab: {
                    active: null
                },
                dialog: {
                    width: 720,
                    show: false,
                    component: null,
                    props: {}
                },
                tabList: [{
                    id: 'createServiceInstance',
                    label: this.$t('新增服务实例'),
                    confirmed: false,
                    component: CreateServiceInstance.name,
                    props: {
                        info: []
                    }
                }, {
                    id: 'deletedServiceInstance',
                    label: this.$t('删除服务实例'),
                    confirmed: false,
                    component: DeletedServiceInstance.name,
                    props: {
                        info: []
                    }
                }, {
                    id: 'moveToIdleHost',
                    label: this.$t('移动到空闲机的主机'),
                    confirmed: false,
                    component: MoveToIdleHost.name,
                    props: {
                        info: []
                    }
                }, {
                    id: 'hostAttrsAutoApply',
                    label: this.$t('属性自动应用'),
                    confirmed: false,
                    component: HostAttrsAutoApply.name,
                    props: {
                        info: []
                    }
                }],
                request: {
                    preview: Symbol('review'),
                    module: Symbol('module'),
                    confirm: Symbol('confirm'),
                    mainline: Symbol('mainline'),
                    host: Symbol('host'),
                    internal: Symbol('internal')
                },
                targetModules: [],
                resources: [],
                type: this.$route.params.type,
                confirmParams: {},
                moduleMap: {},
                isSameModule: false,
                isEmptyChange: false,
                loading: true
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('businessHost', [
                'getDefaultSearchCondition'
            ]),
            confirmText () {
                const textMap = {
                    remove: this.$t('确认移除'),
                    idle: this.$t('确认转移'),
                    business: this.$t('确认转移')
                }
                return textMap[this.type]
            },
            availableTabList () {
                const map = {
                    remove: ['deletedServiceInstance', 'moveToIdleHost'],
                    idle: ['deletedServiceInstance'],
                    business: ['createServiceInstance', 'deletedServiceInstance', 'hostAttrsAutoApply']
                }
                const available = map[this.type]
                return this.tabList.filter(tab => available.includes(tab.id) && tab.props.info.length > 0)
            },
            activeTab () {
                return this.tabList.find(tab => tab.id === this.tab.active) || this.availableTabList[0]
            }
        },
        watch: {
            availableTabList (tabList) {
                tabList.forEach(tab => {
                    if (tab !== this.activeTab) {
                        tab.confirmed = false
                    }
                })
            },
            activeTab (tab) {
                if (!tab) return
                tab.confirmed = true
            }
        },
        async created () {
            this.resolveData(this.$route)
            this.setBreadcrumbs()
            await this.getTopologyModels()
            this.getPreviewData()
        },
        mounted () {
            addResizeListener(this.$refs.changeInfo, this.resizeHandler)
        },
        beforeDestroy () {
            removeResizeListener(this.$refs.changeInfo, this.resizeHandler)
        },
        beforeRouteUpdate (to, from, next) {
            this.resolveData(to)
            this.$nextTick(this.setBreadcrumbs)
            this.getPreviewData()
            next()
        },
        methods: {
            resolveData (route) {
                this.type = route.params.type
                const query = route.query || {}
                const targetModules = query.targetModules
                if (!targetModules) {
                    this.targetModules = []
                } else {
                    this.targetModules = String(targetModules).split(',').map(id => Number(id))
                }

                const resources = query.resources
                if (!resources) {
                    this.resources = []
                } else {
                    this.resources = String(resources).split(',').map(id => Number(id))
                }

                const isTransfer = ['idle', 'business'].includes(this.type)

                const params = {
                    bk_host_ids: this.resources,
                    remove_from_node: {
                        bk_inst_id: isTransfer ? this.bizId : Number(query.sourceId),
                        bk_obj_id: isTransfer ? 'biz' : query.sourceModel
                    }
                }
                if (this.type === 'idle') {
                    params.default_internal_module = this.targetModules[0]
                } else if (this.targetModules.length) {
                    params.add_to_modules = this.targetModules
                }
                this.confirmParams = params
            },
            setBreadcrumbs () {
                const titleMap = {
                    idle: this.$t('转移到空闲模块'),
                    business: this.$t('转移到业务模块'),
                    remove: this.$t('移除主机')
                }
                this.$store.commit('setTitle', titleMap[this.type])
                this.$store.commit('setBreadcrumbs', [{
                    label: this.$t('业务拓扑'),
                    route: {
                        name: MENU_BUSINESS_HOST_AND_SERVICE
                    }
                }, {
                    label: titleMap[this.type]
                }])
            },
            async getTopologyModels () {
                try {
                    const topologyModels = await this.$store.dispatch('objectMainLineModule/searchMainlineObject', {
                        config: {
                            requestId: this.request.mainline
                        }
                    })
                    this.$store.commit('businessHost/setTopologyModels', topologyModels)
                } catch (e) {
                    console.error(e)
                }
            },
            async getPreviewData () {
                try {
                    this.loading = true
                    const data = await this.$http.post(
                        `host/transfer_with_auto_clear_service_instance/bk_biz_id/${this.bizId}/preview`,
                        this.confirmParams,
                        {
                            requestId: this.request.preview
                        }
                    )
                    this.setConfirmState(data)
                    this.setModulePathInfo(data)
                    this.setHostAttrsAutoApply(data)
                    this.setCreateServiceInstance(data)
                    this.setDeletedServiceInstance(data)
                    if (this.type === 'remove') {
                        this.setMoveToIdleHost(data)
                    }
                    this.loading = false
                } catch (e) {
                    console.error(e)
                    this.loading = false
                }
            },
            setConfirmState (data) {
                // 是否是相同的模块转换
                this.isSameModule = data.every(datum => !(datum.to_add_to_modules.length || datum.to_remove_from_modules.length))
                // 是否溢出的是空服务实例（前端流转不会创建空服务实例，但是ESB会）
                this.isEmptyChange = data.every(datum => {
                    const hasAdd = !datum.to_add_to_modules.length
                    const hasRemoveInstance = datum.to_remove_from_modules.some(module => !module.service_instances.length)
                    return !(hasAdd || hasRemoveInstance)
                })
            },
            async setModulePathInfo (data) {
                try {
                    const moduleIds = [...this.targetModules]
                    data.forEach(datum => {
                        moduleIds.push(...datum.to_add_to_modules.map(datum => datum.bk_module_id))
                        moduleIds.push(...datum.to_remove_from_modules.map(datum => datum.bk_module_id))
                    })
                    const uniqueModules = [...new Set(moduleIds)]
                    const result = await this.$store.dispatch('objectMainLineModule/getTopoPath', {
                        bizId: this.bizId,
                        params: {
                            topo_nodes: uniqueModules.map(id => ({ bk_obj_id: 'module', bk_inst_id: id }))
                        }
                    })
                    const moduleMap = {}
                    result.nodes.forEach(node => {
                        moduleMap[node.topo_node.bk_inst_id] = node.topo_path
                    })
                    this.moduleMap = Object.freeze(moduleMap)
                } catch (e) {
                    console.error(e)
                }
            },
            getModulePath (id) {
                const info = this.moduleMap[id] || []
                const path = info.map(node => node.bk_inst_name).reverse().join(' / ')
                return path
            },
            setHostAttrsAutoApply (data) {
                const conflictInfo = (data || []).map(item => item.host_apply_plan)
                const tab = this.tabList.find(tab => tab.id === 'hostAttrsAutoApply')
                tab.props.info = Object.freeze(conflictInfo)
            },
            setCreateServiceInstance (data) {
                const instanceInfo = []
                data.forEach(item => {
                    item.to_add_to_modules.forEach(moduleInfo => {
                        instanceInfo.push({
                            bk_host_id: item.bk_host_id,
                            ...moduleInfo
                        })
                    })
                })
                const tab = this.tabList.find(tab => tab.id === 'createServiceInstance')
                tab.props.info = Object.freeze(instanceInfo)
                this.setHostInfo(data)
            },
            async setHostInfo (data) {
                try {
                    const result = await this.$store.dispatch('hostSearch/searchHost', {
                        params: this.getSearchHostParams(data),
                        config: {
                            requestId: this.request.host
                        }
                    })
                    this.hostInfo = result.info
                } catch (e) {
                    console.error(e)
                }
            },
            getSearchHostParams (data) {
                const params = {
                    bk_biz_id: this.bizId,
                    ip: { data: [], exact: 0, flag: 'bk_host_innerip|bk_host_outerip' },
                    page: {},
                    condition: this.getDefaultSearchCondition()
                }
                const hostCondition = params.condition.find(target => target.bk_obj_id === 'host')
                hostCondition.condition.push({
                    field: 'bk_host_id',
                    operator: '$in',
                    value: data.map(item => item.bk_host_id)
                })
                return params
            },
            setDeletedServiceInstance (data) {
                const deletedServiceInstance = []
                data.forEach(item => {
                    item.to_remove_from_modules.forEach(module => {
                        deletedServiceInstance.push(...module.service_instances)
                    })
                })
                const tab = this.tabList.find(tab => tab.id === 'deletedServiceInstance')
                tab.props.info = Object.freeze(deletedServiceInstance)
            },
            getModuleName (id) {
                const topoInfo = this.moduleMap[id] || []
                const target = topoInfo.find(target => target.bk_obj_id === 'module' && target.bk_inst_id === id) || {}
                return target.bk_inst_name
            },
            async setMoveToIdleHost (data) {
                try {
                    const internalTopology = await this.getInternalTopology()
                    const internalModules = internalTopology.module
                    const idleModule = internalModules.find(module => module.default === 1)
                    const idleHost = []
                    data.forEach(item => {
                        const finalModules = item.final_modules
                        const isIdleModule = finalModules.length && finalModules[0] === idleModule.bk_module_id
                        if (isIdleModule) {
                            idleHost.push(item.bk_host_id)
                        }
                    })
                    const tab = this.tabList.find(tab => tab.id === 'moveToIdleHost')
                    tab.props.info = Object.freeze(idleHost)
                } catch (e) {
                    console.error(e)
                }
            },
            getInternalTopology () {
                return this.$store.dispatch('objectMainLineModule/getInternalTopo', {
                    bizId: this.bizId,
                    config: {
                        requestId: this.request.internal
                    }
                })
            },
            handleRemoveModule (id) {
                const targetModules = this.targetModules.filter(exist => exist !== id)
                this.$router.replace({
                    name: MENU_BUSINESS_TRANSFER_HOST,
                    params: {
                        type: 'business'
                    },
                    query: {
                        ...this.$route.query,
                        targetModules: targetModules.join(',')
                    }
                })
            },
            resizeHandler () {
                this.$nextTick(() => {
                    const scroller = this.$el.parentElement
                    this.hasScrollbar = scroller.scrollHeight > scroller.offsetHeight
                })
            },
            handleTabClick (tab) {
                this.tab.active = tab.id
            },
            handleChangeModule () {
                this.dialog.props = {
                    moduleType: this.type,
                    title: this.type === 'idle' ? this.$t('转移主机到空闲模块') : this.$t('转移主机到业务模块'),
                    defaultChecked: this.targetModules
                }
                this.dialog.width = 720
                this.dialog.component = ModuleSelector.name
                this.dialog.show = true
            },
            handleDialogCancel () {
                this.dialog.show = false
            },
            handleDialogConfirm () {
                if (this.dialog.component === ModuleSelector.name) {
                    this.gotoTransferPage(...arguments)
                    this.dialog.show = false
                }
            },
            gotoTransferPage (modules) {
                this.$router.replace({
                    name: MENU_BUSINESS_TRANSFER_HOST,
                    params: {
                        type: this.dialog.props.moduleType
                    },
                    query: {
                        ...this.$route.query,
                        targetModules: modules.map(node => node.data.bk_inst_id).join(',')
                    }
                })
            },
            async handleConfrim () {
                try {
                    const params = { ...this.confirmParams }
                    const createComponent = this.$refs.createServiceInstance && this.$refs.createServiceInstance[0]
                    const hostAttrsComponent = this.$refs.hostAttrsAutoApply && this.$refs.hostAttrsAutoApply[0]
                    if (createComponent || hostAttrsComponent) {
                        params.options = {}
                        if (createComponent) {
                            params.options.service_instance_options = createComponent.$refs.serviceInstance.map((component, index) => {
                                const instance = createComponent.instances[index]
                                return {
                                    bk_module_id: instance.bk_module_id,
                                    bk_host_id: instance.bk_host_id,
                                    processes: component.processList.map((process, listIndex) => ({
                                        process_template_id: component.templates[listIndex] ? component.templates[listIndex].id : 0,
                                        process_info: process
                                    }))
                                }
                            })
                        }
                        if (hostAttrsComponent) {
                            const conflictResolveResult = hostAttrsComponent.$refs.confirmTable.conflictResolveResult
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
                            params.options.host_apply_conflict_resolvers = conflictResolvers
                        }
                    }
                    await this.$http.post(
                        `host/transfer_with_auto_clear_service_instance/bk_biz_id/${this.bizId}`, params, {
                            requestId: this.request.confirm
                        }
                    )
                    const success = this.type === 'remove' ? '移除成功' : '转移成功'
                    this.$success(this.$t(success))
                    this.redirect()
                } catch (e) {
                    console.error(e)
                }
            },
            handleCancel () {
                this.redirect()
            },
            redirect () {
                if (this.$route.query.from) {
                    this.$router.replace(this.$route.query.from)
                } else {
                    this.$router.replace({
                        name: MENU_BUSINESS_HOST_AND_SERVICE,
                        query: {
                            node: `${this.$route.query.sourceModel}-${this.$route.query.sourceId}`
                        }
                    })
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .info {
        .info-label {
            width: 128px;
            font-size: 14px;
            font-weight: bold;
            color: $textColor;
            text-align: right;
        }
        .info-content {
            overflow: hidden;
            padding: 0 20px 0 8px;
            font-size: 14px;
            .info-count {
                font-weight: bold;
            }
            .module-grep {
                border-top: 1px solid $borderColor;
                margin-top: 10px;
            }
        }
    }
    .module-list {
        font-size: 0;
        .module-item {
            position: relative;
            display: inline-block;
            vertical-align: middle;
            height: 26px;
            max-width: 150px;
            line-height: 24px;
            padding: 0 15px;
            margin: 0 10px 8px 0;
            border: 1px solid #C4C6CC;
            border-radius: 13px;
            color: $textColor;
            font-size: 12px;
            outline: none;
            cursor: default;
            @include ellipsis;
            &.is-business-module {
                padding: 0 12px 0 25px;
            }
            &.is-trigger {
                width: 40px;
                padding: 0;
                text-align: center;
                font-size: 0;
                cursor: pointer;
                .icon-cc-plus {
                    font-size: 14px;
                }
            }
            &:hover {
                border-color: $primaryColor;
                color: $primaryColor;
                .module-mask {
                    display: block;
                }
                .module-icon {
                    background-color: $primaryColor;
                }
            }
            .module-mask {
                display: none;
                position: absolute;
                left: 0;
                top: 0;
                width: 100%;
                height: 100%;
                color: #fff;
                background-color: rgba(0, 0, 0, 0.53);
                text-align: center;
                cursor: pointer;
            }
            .module-icon {
                position: absolute;
                left: 2px;
                top: 2px;
                width: 20px;
                height: 20px;
                border-radius: 50%;
                line-height: 20px;
                text-align: center;
                color: #FFF;
                font-size: 12px;
                background-color: #C4C6CC;
            }
            .module-remove {
                position: absolute;
                right: 4px;
                top: 4px;
                width: 16px;
                height: 16px;
                border-radius: 50%;
                text-align: center;
                line-height: 16px;
                color: #FFF;
                font-size: 0px;
                background-color: #C4C6CC;
                cursor: pointer;
                &:before {
                    display: inline-block;
                    vertical-align: middle;
                    font-size: 20px;
                    transform: translateX(-2px) scale(.5);
                }
            }
        }
    }
    .tab {
        .tab-grep {
            width: 2px;
            height: 19px;
            margin: 0 8px;
            background-color: #C4C6CC;
        }
        .tab-item {
            position: relative;
            color: $textColor;
            font-size: 0;
            cursor: pointer;
            &.active {
                color: $primaryColor;
                .tab-count {
                    color: #FFF;
                    background-color: $primaryColor;
                }
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
                margin-right: 7px;
                font-size: 14px;
            }
            .tab-count {
                position: relative;
                display: inline-block;
                vertical-align: middle;
                height: 16px;
                padding: 0 5px;
                border-radius: 4px;
                line-height: 16px;
                font-size: 12px;
                color: #FFF;
                background-color: #979BA5;
                &.has-badge:after {
                    position: absolute;
                    top: -3px;
                    right: -3px;
                    width: 6px;
                    height: 6px;
                    border-radius: 50%;
                    border: 1px solid #FFF;
                    background-color: $dangerColor;
                    content: "";
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
</style>
