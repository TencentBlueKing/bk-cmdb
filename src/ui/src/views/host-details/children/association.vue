<template>
    <div class="association">
        <div class="options clearfix">
            <div class="fl" v-show="activeView === viewName.list">
                <cmdb-auth class="inline-block-middle mr10"
                    v-if="hasAssociation"
                    :auth="updateAuthResources">
                    <bk-button slot-scope="{ disabled }"
                        theme="primary"
                        class="options-button"
                        :disabled="disabled"
                        @click="showCreate = true">
                        {{$t('新增关联')}}
                    </bk-button>
                </cmdb-auth>
                <span class="inline-block-middle mr10" v-else v-bk-tooltips="$t('当前模型暂未定义可用关联')">
                    <bk-button theme="primary" class="options-button" disabled>
                        {{$t('新增关联')}}
                    </bk-button>
                </span>
                <bk-checkbox v-if="hasAssociation"
                    :size="16" class="options-checkbox"
                    @change="handleExpandAll">
                    <span class="checkbox-label">{{$t('全部展开')}}</span>
                </bk-checkbox>
                <bk-button theme="default" class="options-button" v-show="false">{{$t('批量取消')}}</bk-button>
            </div>
            <div class="fr">
                <bk-button class="options-button options-button-view"
                    :theme="activeView === viewName.list ? 'primary' : 'default'"
                    @click="toggleView(viewName.list)">
                    {{$t('列表')}}
                </bk-button>
                <bk-button class="options-button options-button-view"
                    :theme="activeView === viewName.graphics ? 'primary' : 'default'"
                    @click="toggleView(viewName.graphics)">
                    {{$t('拓扑')}}
                </bk-button>
            </div>
        </div>
        <div class="association-view">
            <component :is="activeView"></component>
        </div>
        <bk-sideslider v-transfer-dom :is-show.sync="showCreate" :width="800" :title="$t('新增关联')">
            <cmdb-host-association-create slot="content" v-if="showCreate"></cmdb-host-association-create>
        </bk-sideslider>
    </div>
</template>

<script>
    import cmdbHostAssociationList from './association-list.vue'
    // import cmdbHostAssociationGraphics from './association-graphics.vue'
    import cmdbHostAssociationGraphics from './association-graphics.new.vue'
    import cmdbHostAssociationCreate from './association-create.vue'
    import { MENU_RESOURCE_HOST_DETAILS } from '@/dictionary/menu-symbol'
    export default {
        name: 'cmdb-host-association',
        components: {
            cmdbHostAssociationList,
            cmdbHostAssociationGraphics,
            cmdbHostAssociationCreate
        },
        data () {
            return {
                viewName: {
                    'list': cmdbHostAssociationList.name,
                    'graphics': cmdbHostAssociationGraphics.name
                },
                activeView: cmdbHostAssociationList.name,
                showCreate: false
            }
        },
        computed: {
            updateAuthResources () {
                const isResourceHost = this.$route.name === MENU_RESOURCE_HOST_DETAILS
                if (isResourceHost) {
                    return this.$authResources({ type: this.$OPERATION.U_RESOURCE_HOST })
                }
                return this.$authResources({ type: this.$OPERATION.U_HOST })
            },
            hasAssociation () {
                const association = this.$store.state.hostDetails.association
                return !!(association.source.length || association.target.length)
            }
        },
        beforeDestroy () {
            this.$store.commit('hostDetails/toggleExpandAll', false)
        },
        methods: {
            toggleView (view) {
                this.activeView = view
            },
            handleExpandAll (expandAll) {
                this.$store.commit('hostDetails/toggleExpandAll', expandAll)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .association {
        height: 100%;
        .association-view {
            height: calc(100% - 62px);
            @include scrollbar-y;
        }
    }
    .options {
        padding: 15px 0;
        font-size: 0;
        .options-button {
            height: 32px;
            line-height: 30px;
            font-size: 14px;
            &.options-button-view {
                margin: 0 0 0 -1px;
                border-radius: 0;
            }
        }
        .options-checkbox {
            margin: 0 0 0 25px;
            line-height: 32px;
            .checkbox-label {
                padding-left: 4px;
            }
        }
    }
</style>
