<template>
    <div class="empty">
        <div class="empty-content" v-if="withTemplate && !templates.length">
            <img src="../../../../static/svg/cc-empty.svg" alt="">
            <p class="empty-text">{{$t('模板未定义进程', { template: (moduleNode || {}).name })}}</p>
            <p class="empty-tips">{{$t('模板未定义进程提示')}}</p>
            <div class="empty-options">
                <cmdb-auth :auth="$authResources({ type: $OPERATION.U_SERVICE_TEMPLATE })">
                    <bk-button slot-scope="{ disabled }"
                        class="empty-button"
                        theme="primary"
                        :disabled="disabled"
                        @click="goToTemplate">
                        {{$t('跳转模板添加进程')}}
                    </bk-button>
                </cmdb-auth>
            </div>
        </div>
        <cmdb-auth v-else :auth="$authResources({ type: $OPERATION.C_SERVICE_INSTANCE })">
            <template slot-scope="{ disabled }">
                <i class="bk-icon icon-plus empty-plus"
                    @click="handleCreateServiceInstance(disabled)">
                </i>
                <p class="empty-tips">
                    {{$t('创建实例提示')}}
                    <a class="text-primary" href="javascript:void(0)" @click="handleCreateServiceInstance(disabled)">{{$t('立即添加')}}</a>
                </p>
            </template>
        </cmdb-auth>
        <cmdb-dialog v-model="dialog.show" :width="850" :height="460">
            <component
                :is="dialog.component"
                v-bind="dialog.componentProps"
                @confirm="handleDialogConfirm"
                @cancel="handleDialogCancel">
            </component>
        </cmdb-dialog>
    </div>
</template>

<script>
    import HostSelector from '../host/host-selector.vue'
    export default {
        components: {
            [HostSelector.name]: HostSelector
        },
        props: {
            active: Boolean
        },
        data () {
            return {
                visible: false,
                isSearching: true,
                dialog: {
                    show: false,
                    component: null,
                    componentProps: {}
                }
            }
        },
        computed: {
            business () {
                return this.$store.getters['objectBiz/bizId']
            },
            currentNode () {
                return this.$store.state.businessHost.selectedNode
            },
            moduleNode () {
                const node = this.$store.state.businessHost.selectedNode
                if (node && node.data.bk_obj_id === 'module') {
                    return node
                }
                return null
            },
            withTemplate () {
                return this.moduleNode && this.moduleNode.data.service_template_id
            },
            isBlueKing () {
                const node = this.$store.state.businessHost.selectedNode
                if (node) {
                    return node.tree.nodes[0].data.bk_inst_name === '蓝鲸'
                }
                return false
            },
            templates () {
                return this.$parent.templates
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
                            requestId: 'getBatchProcessTemplate_empty'
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
                        templateId: this.moduleNode.data.service_template_id,
                        moduleId: this.moduleNode.data.bk_inst_id
                    }
                })
            },
            handleAddHost () {
                this.dialog.component = HostSelector.name
                this.dialog.show = true
            },
            handleCreateServiceInstance (disabled) {
                if (disabled) return
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
            },
            handleDialogConfirm (selected) {
                this.dialog.show = false
                this.$router.push({
                    name: 'createServiceInstance',
                    params: {
                        setId: this.currentNode.parent.data.bk_inst_id,
                        moduleId: this.currentNode.data.bk_inst_id
                    },
                    query: {
                        resources: selected.map(item => item.host.bk_host_id).join(',')
                    }
                })
            },
            handleDialogCancel () {
                this.dialog.show = false
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
