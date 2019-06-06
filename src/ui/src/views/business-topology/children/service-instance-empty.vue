<template>
    <div class="empty">
        <div class="empty-content" v-if="withTemplate">
            <img src="../../../../static/svg/cc-empty.svg" alt="">
            <p class="empty-text">Agent模版尚未定义进程，无法创建服务</p>
            <p class="empty-tips">您可以先跳转模版添加进程或直接添加主机，后续再添加模版进程</p>
            <div class="empty-options">
                <bk-button class="empty-button" type="primary">跳转模板添加进程</bk-button>
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
                visible: false
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
                const instance = this.$store.state.businessTopology.selectedNodeInstance
                if (this.moduleNode && instance) {
                    return !!instance.template_id
                }
                return false
            }
        },
        methods: {
            handleAddHost () {
                this.visible = true
            },
            async handleSelectHost (checked) {
                await this.$store.dispatch('hostRelation/transferHostModule', {
                    params: {
                        bk_biz_id: this.business,
                        bk_host_id: checked,
                        bk_module_id: [this.moduleNode.data.bk_inst_id],
                        is_increment: true
                    }
                })
                this.visible = false
            },
            handleCreateServiceInstance () {
                this.$router.push({
                    name: 'createServiceInstance',
                    params: {
                        moduleId: this.moduleNode.data.bk_inst_id,
                        setId: this.moduleNode.parent.data.bk_inst_id
                    }
                })
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
        }
    }
</style>
