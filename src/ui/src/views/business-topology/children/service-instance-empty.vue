<template>
    <div class="empty">
        <template v-bkloading="{ isLoading: $loading('getBatchProcessTemplate') }">
            <div class="empty-content" v-if="withTemplate && !templates.length && !isSearching">
                <img src="../../../../static/svg/cc-empty.svg" alt="">
                <p class="empty-text">{{$t('模板未定义进程', { template: (moduleNode || {}).name })}}</p>
                <p class="empty-tips">{{$t('模板未定义进程提示')}}</p>
                <div class="empty-options">
                    <span style="display: inline-block;"
                        v-cursor="{
                            active: !$isAuthorized($OPERATION.U_SERVICE_TEMPLATE),
                            auth: [$OPERATION.U_SERVICE_TEMPLATE]
                        }">
                        <bk-button class="empty-button" theme="primary"
                            :disabled="!$isAuthorized($OPERATION.U_SERVICE_TEMPLATE)"
                            @click="goToTemplate">
                            {{$t('跳转模板添加进程')}}
                        </bk-button>
                    </span>
                    <span style="display: inline-block;"
                        v-cursor="{
                            active: !$isAuthorized($OPERATION.C_SERVICE_INSTANCE),
                            auth: [$OPERATION.C_SERVICE_INSTANCE]
                        }">
                        <bk-button class="empty-button" theme="default"
                            :disabled="!$isAuthorized($OPERATION.C_SERVICE_INSTANCE)"
                            @click="handleAddHost">
                            {{$t('添加主机')}}
                        </bk-button>
                    </span>
                </div>
            </div>
            <div class="empty-content" v-else-if="!isSearching"
                v-cursor="{
                    active: !$isAuthorized($OPERATION.C_SERVICE_INSTANCE),
                    auth: [$OPERATION.C_SERVICE_INSTANCE]
                }">
                <i class="bk-icon icon-plus empty-plus"
                    @click="handleCreateServiceInstance">
                </i>
                <p class="empty-tips">
                    {{$t('创建实例提示')}}
                    <a class="text-primary" href="javascript:void(0)" @click="handleCreateServiceInstance">{{$t('立即添加')}}</a>
                </p>
            </div>
        </template>
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
                templates: [],
                isSearching: true
            }
        },
        computed: {
            business () {
                return this.$store.getters['objectBiz/bizId']
            },
            currentNode () {
                return this.$store.state.businessTopology.selectedNode
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
            },
            isBlueKing () {
                const node = this.$store.state.businessTopology.selectedNode
                if (node) {
                    return node.tree.nodes[0].data.bk_inst_name === '蓝鲸'
                }
                return false
            }
        },
        watch: {
            withTemplate (withTemplate) {
                if (withTemplate) {
                    this.getTemplate()
                } else {
                    this.isSearching = false
                }
            }
        },
        created () {
            if (this.withTemplate) {
                this.getTemplate()
            } else {
                this.isSearching = false
            }
        },
        methods: {
            async getTemplate () {
                try {
                    this.isSearching = true
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
                    this.isSearching = false
                } catch (e) {
                    console.error(e)
                    this.isSearching = false
                }
            },
            goToTemplate () {
                this.$router.push({
                    name: 'operationalTemplate',
                    params: {
                        templateId: this.moduleInstance.service_template_id,
                        moduleId: this.moduleNode.data.bk_inst_id
                    }
                })
            },
            handleAddHost () {
                this.visible = true
            },
            async handleSelectHost (checked) {
                try {
                    const addNum = checked.length
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
                    if (this.withTemplate) {
                        this.currentNode.data.service_instance_count = this.currentNode.data.service_instance_count + addNum
                        this.currentNode.parents.forEach(node => {
                            node.data.service_instance_count = node.data.service_instance_count + addNum
                        })
                    }
                    this.visible = false
                    this.$success(this.$t('添加成功'))
                    this.$emit('create-instance-success', data)
                } catch (e) {
                    console.error(e)
                }
            },
            handleCreateServiceInstance () {
                if (!this.$isAuthorized(this.$OPERATION.C_SERVICE_INSTANCE)) return
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
            display: inline-block;
        }
        .empty-text {
            margin-top: 20px;
            line-height: 29px;
            font-size: 22px;
            color: #63656E;
            word-break: break-word;
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
        .empty-img {
            display: block;
            margin: 0 auto;
        }
    }
</style>
