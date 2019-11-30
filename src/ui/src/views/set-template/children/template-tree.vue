<template>
    <div :class="['template-tree', mode]">
        <div class="node-root clearfix">
            <i class="folder-icon bk-icon icon-down-shape fl" @click="handleCollapse"></i>
            <i class="node-icon fl">{{setName[0]}}</i>
            <span class="root-name" :title="templateName">{{templateName}}</span>
        </div>
        <cmdb-collapse-transition>
            <ul class="node-children" v-show="!collapse">
                <li class="node-child clearfix"
                    v-for="(service, index) in services"
                    :key="service.id"
                    :class="{ selected: selected === service.id }"
                    @click="handleChildClick(service)">
                    <i class="node-icon fl">{{moduleName[0]}}</i>
                    <span class="child-options fr" v-if="mode !== 'view'">
                        <i class="options-view icon icon-cc-show" @click="handleViewService(service)"></i>
                        <bk-popover v-if="serviceExistHost(service.id)">
                            <i class="options-delete icon icon-cc-tips-close disabled"></i>
                            <i18n path="该模块下有主机不可删除" tag="p" class="service-tips" slot="content">
                                <span place="link" @click="handleGoTopoBusiness">{{$t('跳转查看')}}</span>
                            </i18n>
                        </bk-popover>
                        <i v-else class="options-delete icon icon-cc-tips-close" @click="handleDeleteService(index)"></i>
                    </span>
                    <span class="child-name">{{service.name}}</span>
                </li>
                <li class="options-child node-child clearfix"
                    v-if="['create', 'edit'].includes(mode)"
                    @click="handleAddService">
                    <i class="node-icon icon icon-cc-zoom-in fl"></i>
                    <span class="child-name">{{$t('添加服务模板')}}</span>
                </li>
            </ul>
        </cmdb-collapse-transition>
        <bk-dialog
            header-position="left"
            :draggable="false"
            :mask-close="false"
            :width="759"
            :title="dialog.title"
            v-model="dialog.visible"
            @after-leave="handleDialogClose"
            @confirm="handleDialogConfirm">
            <component
                ref="dialogComponent"
                :is="dialog.component"
                :services-host="servicesHost"
                v-bind="dialog.props">
            </component>
            <template slot="footer" v-if="dialog.useCustomFooter">
                <bk-button @click="dialog.visible = false">{{$t('关闭')}}</bk-button>
            </template>
        </bk-dialog>
    </div>
</template>

<script>
    import { MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'
    import serviceTemplateSelector from './service-template-selector.vue'
    import serviceTemplateInfo from './service-template-info.vue'
    export default {
        components: {
            serviceTemplateSelector,
            serviceTemplateInfo
        },
        /* eslint-disable-next-line */
        props: ['mode', 'templateId'],
        data () {
            return {
                templateName: this.$t('模板集群名称'),
                services: [],
                originalServices: [],
                collapse: false,
                selected: null,
                unwatch: null,
                dialog: {
                    visible: false,
                    title: '',
                    useCustomFooter: false,
                    props: {}
                },
                servicesHost: []
            }
        },
        computed: {
            hasChange () {
                if (this.mode !== 'edit') {
                    return false
                }
                if (this.originalServices.length !== this.services.length) {
                    return true
                }
                return this.originalServices.some((service, index) => {
                    const target = this.services[index]
                    return (target && target.id !== service.id) || !target
                })
            },
            setName () {
                const setModel = this.$store.getters['objectModelClassify/getModelById']('set') || {}
                return setModel.bk_obj_name || ''
            },
            moduleName () {
                const moduleModel = this.$store.getters['objectModelClassify/getModelById']('module') || {}
                return moduleModel.bk_obj_name || ''
            },
            sortedServices () {
                return [...this.services].sort((A, B) => {
                    return A.name.localeCompare(B.name, 'zh-Hans-CN', { sensitivity: 'accent' })
                })
            },
            topoNodes () {
                return this.originalServices.map(service => {
                    return {
                        bk_obj_id: 'module',
                        bk_inst_id: service.id
                    }
                })
            }
        },
        watch: {
            hasChange (value) {
                this.$emit('service-change', value)
            },
            services (value) {
                this.$emit('service-selected', value)
            },
            mode () {
                this.selected = null
            }
        },
        created () {
            this.initMonitorTemplateName()
            if (['edit', 'view'].includes(this.mode)) {
                this.getSetTemplateServices()
            }
        },
        beforeDestory () {
            this.unwatch && this.unwatch()
        },
        methods: {
            async getSetTemplateServices () {
                try {
                    const data = await this.$store.dispatch('setTemplate/getSetTemplateServicesStatistics', {
                        bizId: this.$store.getters['objectBiz/bizId'],
                        setTemplateId: this.templateId
                    })
                    this.services = data.map(item => item.service_template)
                    this.servicesHost = data.map(item => {
                        return {
                            service_id: item.service_template.id,
                            host_count: item.host_count
                        }
                    })
                    this.originalServices = [...this.services]
                } catch (e) {
                    console.error(e)
                    this.services = []
                    this.originalServices = []
                    this.servicesHost = []
                }
            },
            serviceExistHost (id) {
                const service = this.servicesHost.find(service => service.service_id === id)
                if (service) {
                    return service.host_count > 0
                }
                return false
            },
            initMonitorTemplateName () {
                this.unwatch = this.$watch(() => {
                    return this.$parent.templateName
                }, value => {
                    if (value) {
                        this.templateName = value
                    } else {
                        this.templateName = this.$t('模板集群名称')
                    }
                }, { immediate: true })
            },
            handleCollapse () {
                this.collapse = !this.collapse
            },
            handleChildClick (service) {
                if (this.mode === 'view') {
                    return false
                }
                this.selected = service.id
            },
            handleAddService () {
                this.selected = null
                this.dialog.props = {
                    selected: this.services.map(service => service.id)
                }
                this.dialog.title = this.$t('添加服务模板')
                this.dialog.useCustomFooter = false
                this.dialog.component = serviceTemplateSelector.name
                this.dialog.visible = true
            },
            handleViewService (service) {
                this.dialog.props = {
                    id: service.id
                }
                this.dialog.title = `【${service.name}】${this.$t('模板服务信息')}`
                this.dialog.useCustomFooter = true
                this.dialog.component = serviceTemplateInfo.name
                this.dialog.visible = true
            },
            handleDialogConfirm () {
                if (this.dialog.component === serviceTemplateSelector.name) {
                    this.services = this.$refs.dialogComponent.getSelectedServices()
                }
            },
            handleDialogClose () {
                this.dialog.component = null
                this.dialog.title = ''
                this.dialog.props = {}
            },
            handleDeleteService (index) {
                this.services.splice(index, 1)
            },
            recoveryService () {
                this.services = [...this.originalServices]
            },
            handleGoTopoBusiness () {
                this.$router.push({
                    name: MENU_BUSINESS_HOST_AND_SERVICE
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    $iconColor: #C4C6CC;
    $fontColor: #63656E;
    $highlightColor: #3A84FF;
    $iconDisabledColor: #D8D8D8;
    .template-tree {
        padding: 10px 0 10px 20px;
        border: 1px solid #C4C6CC;
        background-color: #fff;
        &:not(.view) {
            .node-child:hover {
                background-color: rgba(240,241,245, .6);
                .child-name {
                    color: $highlightColor;
                }
                .child-options {
                    display: block;
                }
            }
        }
    }
    .node-icon {
        position: relative;
        margin: 8px 4px 8px 0px;
        width: 20px;
        height: 20px;
        border-radius: 50%;
        line-height: 20px;
        text-align: center;
        font-size: 12px;
        font-style: normal;
        color: #fff;
        background-color: #97AED6;
        z-index: 2;
    }
    .node-root {
        line-height: 36px;
        cursor: default;
        .folder-icon {
            width: 23px;
            height: 36px;
            line-height: 36px;
            text-align: center;
            font-size: 12px;
            color: $iconColor;
            cursor: pointer;
        }
        .root-name {
            display: block;
            padding: 0 10px 0 0;
            font-size: 14px;
            color: $fontColor;
            @include ellipsis;
        }
    }
    .node-children {
        line-height: 36px;
        margin-left: 32px;
        cursor: default;
        .node-child {
            padding: 0 10px 0 32px;
            position: relative;
            &.selected {
                background-color: #F0F1F5;
            }
            &.selected {
                .node-icon {
                    background-color: $highlightColor;
                }
                .child-name {
                    color: $highlightColor;
                }
                .child-options {
                    display: block;
                }
            }
            &:before {
                position: absolute;
                left: 0px;
                top: -18px;
                content: "";
                width: 25px;
                height: 36px;
                border-left: 1px dashed #DCDEE5;
                border-bottom: 1px dashed #DCDEE5;
                z-index: 1;
            }
            &.options-child {
                cursor: pointer;
                .node-icon {
                    font-size: 18px;
                    background-color: transparent;
                    color: $highlightColor;
                }
                .child-name {
                    color: $highlightColor;
                }
            }
            .child-name {
                display: block;
                padding: 0 10px 0 0;
                font-size: 14px;
                color: $fontColor;
                @include ellipsis;
            }
            .child-options {
                display: none;
                margin-right: 9px;
                font-size: 0;
                color: $iconColor;
                .options-view {
                    font-size: 18px;
                    cursor: pointer;
                    &:hover {
                        color: $highlightColor;
                    }
                }
                .options-delete {
                    width: 24px;
                    height: 24px;
                    margin-left: 14px;
                    font-size: 12px;
                    text-align: center;
                    line-height: 24px;
                    cursor: pointer;
                    &:hover {
                        color: $highlightColor;
                    }
                    &.disabled:hover {
                        color: $iconDisabledColor;
                        cursor: not-allowed;
                    }
                }
            }
        }
    }
    .service-tips {
        span {
            color: $highlightColor;
            cursor: pointer;
        }
    }
</style>
