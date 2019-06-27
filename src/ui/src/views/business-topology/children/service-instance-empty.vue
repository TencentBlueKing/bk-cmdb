<template>
    <div class="empty">
        <div class="empty-content" v-if="withTemplate && !templates.length">
            <img src="../../../../static/svg/cc-empty.svg" alt="">
            <p class="empty-text">{{$t('BusinessTopology["模板未定义进程"]', { template: (moduleNode || {}).name })}}</p>
            <p class="empty-tips">{{$t('BusinessTopology["模板未定义进程提示"]')}}</p>
            <div class="empty-options">
                <bk-button class="empty-button" type="primary" @click="goToTemplate">跳转模板添加进程</bk-button>
                <bk-button class="empty-button" type="default" @click="handleAddHost">添加主机</bk-button>
            </div>
        </div>
        <div class="empty-content" v-else>
            <i class="bk-icon icon-plus empty-plus"
                @click="handleCreateServiceInstance">
            </i>
            <p class="empty-tips">
                您还没有创建任何服务实例，
                <a class="text-primary" href="javascript:void(0)" @click="handleCreateServiceInstance">立即添加</a>
            </p>
        </div>
        <host-selector
            :visible.sync="visible"
            :module-instance="moduleInstance"
            @host-selected="handleSelectHost">
        </host-selector>
    </div>
</template>

<script>
    import hostSelector from '@/components/ui/selector/host.vue'
    export default {
        components: {
            hostSelector
        },
        data () {
            return {
                visible: false,
                templates: []
            }
        },
        computed: {
            business () {
                return this.$store.getters['objectBiz/bizId']
            },
            moduleNode () {
                const node = this.$store.state.businessTopology.selectedNode
                if (node && node.data.bk_obj_id === 'module') {
                    return node
                }
                return null
            },
            moduleInstance () {
                const instance = this.$store.state.businessTopology.selectedNodeInstance
                if (this.moduleNode && instance) {
                    return instance
                }
                return {}
            },
            withTemplate () {
                return this.moduleNode && this.moduleInstance.service_template_id
            }
        },
        watch: {
            withTemplate (withTemplate) {
                if (withTemplate) {
                    this.getTemplate()
                }
            }
        },
        created () {
            if (this.withTemplate) {
                this.getTemplate()
            }
        },
        methods: {
            async getTemplate () {
                try {
                    const data = await this.$store.dispatch('processTemplate/getBatchProcessTemplate', {
                        params: this.$injectMetadata({
                            service_template_id: this.moduleInstance.service_template_id
                        }),
                        config: {
                            requestId: 'getBatchProcessTemplate',
                            cancelPrevious: true
                        }
                    })
                    this.templates = data.info
                } catch (e) {
                    console.error(e)
                }
            },
            goToTemplate () {
                this.$router.push({
                    name: 'operationalTemplate',
                    params: {
                        templateId: this.moduleInstance.service_template_id
                    },
                    query: {
                        from: {
                            name: this.$route.name,
                            query: {
                                module: this.moduleInstance.bk_module_id
                            }
                        }
                    }
                })
            },
            handleAddHost () {
                this.visible = true
            },
            async handleSelectHost (checked) {
                try {
                    const data = await this.$store.dispatch('serviceInstance/createProcServiceInstanceByTemplate', {
                        params: this.$injectMetadata({
                            name: this.moduleInstance.bk_module_name,
                            bk_module_id: this.moduleInstance.bk_module_id,
                            service_template_id: this.moduleInstance.service_template_id,
                            instances: checked.map(hostId => {
                                return {
                                    bk_host_id: hostId,
                                    processes: this.templates.map(template => {
                                        const processInfo = {}
                                        Object.keys(template.property).forEach(key => {
                                            processInfo[key] = template.property[key].value
                                        })
                                        return {
                                            process_template_id: template.id,
                                            process_info: processInfo
                                        }
                                    })
                                }
                            })
                        })
                    })
                    this.visible = false
                    this.$emit('create-instance-success', data)
                } catch (e) {
                    console.error(e)
                }
            },
            handleCreateServiceInstance () {
                if (this.withTemplate) {
                    this.handleAddHost()
                } else {
                    this.$router.push({
                        name: 'createServiceInstance',
                        params: {
                            moduleId: this.moduleNode.data.bk_inst_id,
                            setId: this.moduleNode.parent.data.bk_inst_id
                        },
                        query: {
                            from: this.$route.fullPath,
                            title: this.moduleNode.name
                        }
                    })
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .empty{
        height: 100%;
        text-align: center;
        &:before {
            content: "";
            width: 0;
            height: 100%;
            @include inlineBlock;
        }
        .empty-content {
            margin: -120px 0 0 0;
            @include inlineBlock;
        }
        .empty-text {
            margin-top: 20px;
            line-height: 29px;
            font-size: 22px;
            color: #63656E;
        }
        .empty-tips {
            margin-top: 10px;
            font-size: 14px;
            line-height: 20px;
        }
        .empty-options {
            margin-top: 25px;
            .empty-button {
                height: 32px;
                margin: 0 4px;
                line-height: 30px;
            }
        }
        .empty-plus {
            @include inlineBlock;
            width: 70px;
            height: 70px;
            margin-bottom: 9px;
            line-height: 70px;
            border: 1px dashed #3A84FF;
            font-size: 19px;
            color: #3A84FF;
            cursor: pointer;
            &:hover {
                font-weight: bold;
                border-style: solid;
                box-shadow: 0 0 2px #3A84FF;
            }
        }
    }
</style>
