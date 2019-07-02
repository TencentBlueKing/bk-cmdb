<template>
    <div class="association">
        <div class="options clearfix">
            <div class="fl" v-show="activeView === viewName.list">
                <span class="inline-block-middle" v-if="hasAssociation"
                    v-cursor="{
                        active: !$isAuthorized(updateAuth),
                        auth: [updateAuth]
                    }">
                    <bk-button type="primary" class="options-button"
                        :disabled="!$isAuthorized(updateAuth)"
                        @click="showCreate = true">
                        {{$t('HostDetails["新增关联"]')}}
                    </bk-button>
                </span>
                <cmdb-form-bool v-if="hasAssociation"
                    :size="16" class="options-checkbox"
                    @change="handleExpandAll">
                    <span class="checkbox-label">{{$t('Common["全部展开"]')}}</span>
                </cmdb-form-bool>
                <bk-button type="default" class="options-button" v-show="false">{{$t('HostDetails["批量取消"]')}}</bk-button>
            </div>
            <div class="fr">
                <bk-button class="options-button options-button-view"
                    :type="activeView === viewName.list ? 'primary' : 'default'"
                    @click="toggleView(viewName.list)">
                    {{$t('HostDetails["列表"]')}}
                </bk-button>
                <bk-button class="options-button options-button-view"
                    :type="activeView === viewName.graphics ? 'primary' : 'default'"
                    @click="toggleView(viewName.graphics)">
                    {{$t('HostDetails["拓扑"]')}}
                </bk-button>
            </div>
        </div>
        <div class="association-view">
            <component :is="activeView"></component>
        </div>
        <cmdb-slider :is-show.sync="showCreate">
            <cmdb-host-association-create slot="content" v-if="showCreate"></cmdb-host-association-create>
        </cmdb-slider>
    </div>
</template>

<script>
    import cmdbHostAssociationList from './association-list.vue'
    import cmdbHostAssociationGraphics from './association-graphics.vue'
    import cmdbHostAssociationCreate from './association-create.vue'
    import { OPERATION, RESOURCE_HOST } from '../router.config.js'
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
            updateAuth () {
                const isResourceHost = this.$route.name === RESOURCE_HOST
                if (isResourceHost) {
                    return OPERATION.U_RESOURCE_HOST
                }
                return OPERATION.U_HOST
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
        }
    }
    .options {
        padding: 15px 0;
        font-size: 0;
        .options-button {
            height: 32px;
            margin: 0 11px 0 0;
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
