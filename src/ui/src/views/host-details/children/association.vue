<template>
    <div class="association">
        <div class="options clearfix">
            <div class="fl" v-show="activeView === viewName.list">
                <bk-button type="primary" class="options-button" @click="showCreate = true">{{$t('Association["新增关联"]')}}</bk-button>
                <bk-button type="default" class="options-button" v-show="false">{{$t('Association["批量取消"]')}}</bk-button>
            </div>
            <div class="fr">
                <bk-button class="options-button options-button-view"
                    :type="activeView === viewName.list ? 'primary' : 'default'"
                    @click="toggleView(viewName.list)">
                    {{$t('Association["列表"]')}}
                </bk-button>
                <bk-button class="options-button options-button-view"
                    :type="activeView === viewName.graphics ? 'primary' : 'default'"
                    @click="toggleView(viewName.graphics)">
                    {{$t('Association["拓扑"]')}}
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
        methods: {
            toggleView (view) {
                this.activeView = view
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
    }
</style>
