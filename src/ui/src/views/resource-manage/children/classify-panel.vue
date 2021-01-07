<template>
    <div class="classify">
        <h4 class="classify-name" :title="classify['bk_classification_name']">
            <span class="classify-name-text">{{classify['bk_classification_name']}}</span>
        </h4>
        <div class="models-layout">
            <div class="models-link" v-for="(model, index) in classify['bk_objects']"
                :key="index"
                :title="model['bk_obj_name']"
                @click="redirect(model)">
                <i :class="['model-icon','icon', model['bk_obj_icon'], { 'nonpre-mode': !model['ispre'] }]"></i>
                <span class="model-name">{{model['bk_obj_name']}}</span>
                <i class="model-star bk-icon"
                    :class="[isCollected(model) ? 'icon-star-shape' : 'icon-star']"
                    @click.prevent.stop="toggleCustomNavigation(model)">
                </i>
                <span class="model-instance-count">{{getInstanceCount(model)}}</span>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import {
        MENU_RESOURCE_HOST,
        MENU_RESOURCE_BUSINESS,
        MENU_RESOURCE_INSTANCE,
        MENU_RESOURCE_COLLECTION,
        MENU_RESOURCE_HOST_COLLECTION,
        MENU_RESOURCE_BUSINESS_COLLECTION
    } from '@/dictionary/menu-symbol'
    export default {
        props: {
            classify: {
                type: Object,
                required: true
            },
            instanceCount: {
                type: Array,
                required: true
            },
            collection: {
                type: Array,
                required: true
            }
        },
        data () {
            return {
                maxCustomNavigationCount: 8
            }
        },
        computed: {
            ...mapGetters('userCustom', ['usercustom']),
            collectedCount () {
                return this.collection.length
            }
        },
        methods: {
            getInstanceCount (model) {
                const data = this.instanceCount.find(data => data.bk_obj_id === model.bk_obj_id)
                if (data) {
                    return data.instance_count
                }
                return 0
            },
            redirect (model) {
                const map = {
                    host: MENU_RESOURCE_HOST,
                    biz: MENU_RESOURCE_BUSINESS
                }
                if (map.hasOwnProperty(model.bk_obj_id)) {
                    this.$routerActions.redirect({
                        name: map[model.bk_obj_id]
                    })
                } else {
                    this.$routerActions.redirect({
                        name: MENU_RESOURCE_INSTANCE,
                        params: {
                            objId: model.bk_obj_id
                        }
                    })
                }
            },
            isCollected (model) {
                return this.collection.includes(model.bk_obj_id)
            },
            toggleCustomNavigation (model) {
                if (['host', 'biz'].includes(model.bk_obj_id)) {
                    this.toggleDefaultCollection(model)
                } else {
                    let isAdd = false
                    let newCollection
                    const oldCollection = this.usercustom[MENU_RESOURCE_COLLECTION] || []
                    if (oldCollection.includes(model['bk_obj_id'])) {
                        newCollection = oldCollection.filter(id => id !== model.bk_obj_id)
                    } else {
                        isAdd = true
                        newCollection = [...oldCollection, model.bk_obj_id]
                    }
                    if (isAdd && this.collectedCount >= this.maxCustomNavigationCount) {
                        this.$warn(this.$t('限制添加导航提示', { max: this.maxCustomNavigationCount }))
                        return false
                    }
                    const promise = this.$store.dispatch('userCustom/saveUsercustom', {
                        [MENU_RESOURCE_COLLECTION]: newCollection
                    })
                    promise.then(() => {
                        this.$success(isAdd ? this.$t('添加导航成功') : this.$t('取消导航成功'))
                    })
                }
            },
            async toggleDefaultCollection (model) {
                const isCollected = this.collection.includes(model.bk_obj_id)
                if (!isCollected && this.collection.length >= this.maxCustomNavigationCount) {
                    this.$warn(this.$t('限制添加导航提示', { max: this.maxCustomNavigationCount }))
                } else {
                    try {
                        const key = model.bk_obj_id === 'host' ? MENU_RESOURCE_HOST_COLLECTION : MENU_RESOURCE_BUSINESS_COLLECTION
                        await this.$store.dispatch('userCustom/saveUsercustom', {
                            [key]: !isCollected
                        })
                        this.$success(isCollected ? this.$t('取消导航成功') : this.$t('添加导航成功'))
                    } catch (e) {
                        console.error(e)
                    }
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .classify{
        margin: 0 0 20px 0;
        background-color: #fff;
        border: 1px solid #ebf0f5;
        box-shadow:0px 3px 6px 0px rgba(51,60,72,0.05);
    }
    .classify-name{
        padding: 13px 5px;
        margin: 0 20px;
        line-height: 20px;
        font-size: 0;
        color: $cmdbTextColor;
        border-bottom: 1px solid #ebf0f5;
        &-text {
            display: inline-block;
            padding: 0 2px 0 0;
            vertical-align: middle;
            max-width: calc(100% - 40px);
            font-size: 14px;
            @include ellipsis;
        }
        &-count {
            display: inline-block;
            width: 40px;
            vertical-align: middle;
            font-size: 14px;
        }
    }
    .models-layout{
        padding: 8px 0;
        .models-link{
            display: block;
            height: 38px;
            font-size: 0;
            position: relative;
            padding: 7px 25px;
            cursor: pointer;
            &:hover{
                background-color: #ecf3ff;
            }
            &:before{
                content: "";
                display: inline-block;
                height: 100%;
                vertical-align: middle;
            }
            &:hover .model-icon,
            &:hover .model-name{
                color: #3A84FF;
            }
            &:hover .model-star{
                display: inline-block;
            }
            .model-icon,
            .model-name{
                display: inline-block;
                vertical-align: middle;
            }
            .model-icon{
                font-size: 16px;
                color: #798AAD;
            }
            .nonpre-mode {
                color: #3A84FF !important;
            }
            .model-name{
                max-width: calc(100% - 100px);
                margin: 0 0 0 12px;
                font-size: 14px;
                line-height: 24px;
                color: $cmdbTextColor;
                @include ellipsis;
            }
            .model-instance-count {
                float: right;
                width: 35px;
                font-size: 14px;
                line-height: 24px;
                color: #C4C6CC;
                text-align: right;
                @include inlineBlock;
            }
            .model-star{
                display: none;
                width: 24px;
                height: 24px;
                margin-left: 5px;
                line-height: 24px;
                text-align: center;
                font-size: 14px;
                cursor: pointer;
                vertical-align: middle;
                &.icon-star-shape{
                    color: #FFB400;
                    display: inline-block;
                }
            }
        }
    }
</style>
