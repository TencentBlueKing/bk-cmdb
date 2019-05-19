<template>
    <div class="association">
        <div class="options clearfix">
            <div class="fl">
                <bk-button type="primary" class="options-button">{{$t('Association["关联管理"]')}}</bk-button>
                <bk-button type="default" class="options-button">{{$t('Association["批量取消"]')}}</bk-button>
            </div>
            <div class="fr">
                <bk-button class="options-button options-button-view"
                    :type="active === 'list' ? 'primary' : 'default'"
                    @click="active = 'list'">
                    {{$t('Association["列表"]')}}
                </bk-button>
                <bk-button class="options-button options-button-view"
                    :type="active === 'topology' ? 'primary' : 'default'"
                    @click="active = 'topology'">
                    {{$t('Association["拓扑"]')}}
                </bk-button>
            </div>
        </div>
        <div class="association-view">
            <component :is="view"></component>
        </div>
    </div>
</template>

<script>
    import cmdbHostAssociationList from './association-list.vue'
    export default {
        name: 'cmdb-host-association',
        components: {
            cmdbHostAssociationList
        },
        data () {
            return {
                active: 'list',
                view: cmdbHostAssociationList.name
            }
        }
    }
</script>

<style lang="scss" scoped>
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
