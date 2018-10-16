<template>
    <div class="classify">
        <h4 class="classify-name" :title="classify['bk_classification_name']">
            <span class="classify-name-text">{{classify['bk_classification_name']}}</span>
            <span class="classify-name-count">({{classify['bk_objects'].length}})</span>
        </h4>
        <div class="models-layout">
            <div class="models-link" v-for="(model, index) in classify['bk_objects']"
                :key="index"
                :title="model['bk_obj_name']"
                @click="redirect(model)">
                <i :class="['model-icon','icon', model['bk_obj_icon']]"></i>
                <span class="model-name">{{model['bk_obj_name']}}</span>
                <i class="model-star bk-icon"
                    v-if="!notCollectable.includes(classify['bk_classification_id'])"
                    :class="[customNavigation.includes(model['bk_obj_id']) ? 'icon-star-shape' : 'icon-star']"
                    @click.prevent.stop="toggleCustomNavigation(model)">
                </i>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        props: {
            classify: {
                type: Object,
                required: true
            }
        },
        data () {
            return {
                maxCustomNavigationCount: 8,
                notCollectable: ['bk_host_manage', 'bk_organization'],
                notModelClassify: ['bk_host_manage', 'bk_back_config']
            }
        },
        computed: {
            ...mapGetters('userCustom', ['usercustom', 'classifyNavigationKey']),
            ...mapGetters('objectModelClassify', ['authorizedNavigation']),
            customNavigation () {
                return this.usercustom[this.classifyNavigationKey] || []
            },
            usefulNavigation () {
                const usefulNavigation = this.customNavigation.filter(customId => {
                    return this.authorizedNavigation.some(({children}) => children.some(navigation => navigation.id === customId))
                })
                return usefulNavigation
            }
        },
        methods: {
            redirect (model) {
                const path = this.getModelLink(model)
                this.$store.commit('setHeaderStatus', {
                    back: true
                })
                this.$router.push(path)
            },
            getModelLink (model) {
                if (this.notModelClassify.includes(model['bk_classification_id'])) {
                    return model.path
                }
                return model['bk_obj_id'] === 'biz' ? '/business' : `/general-model/${model['bk_obj_id']}`
            },
            toggleCustomNavigation (model) {
                let newCustom
                let oldCustom = this.customNavigation
                let isAdd = false
                if (oldCustom.includes(model['bk_obj_id'])) {
                    newCustom = oldCustom.filter(id => id !== model['bk_obj_id'])
                } else {
                    isAdd = true
                    newCustom = [...oldCustom, model['bk_obj_id']]
                }
                if (isAdd && this.usefulNavigation.length >= this.maxCustomNavigationCount) {
                    this.$warn(this.$t('Index["限制添加导航提示"]', {max: this.maxCustomNavigationCount}))
                    return false
                }
                this.$store.dispatch('userCustom/saveUsercustom', {
                    [this.classifyNavigationKey]: newCustom
                }).then(() => {
                    this.$http.cancel('post_searchUsercustom')
                    this.$success(isAdd ? this.$t('Index["添加导航成功"]') : this.$t('Index["取消导航成功"]'))
                })
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
            height: 36px;
            font-size: 0;
            position: relative;
            padding: 6px 25px;
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
                color: #0082ff;
            }
            &:hover .model-star{
                display: block;
            }
            .model-icon,
            .model-name{
                display: inline-block;
                vertical-align: middle;
            }
            .model-icon{
                font-size: 16px;
                color: $cmdbTextColor;
            }
            .model-name{
                max-width: calc(100% - 60px);
                margin: 0 0 0 12px;
                font-size: 14px;
                line-height: 24px;
                color: $cmdbTextColor;
                @include ellipsis;
            }
            .model-star{
                display: none;
                position: absolute;
                width: 36px;
                height: 36px;
                line-height: 36px;
                text-align: center;
                right: 12px;
                top: 0;
                color: #c3cdd7;
                font-size: 14px;
                cursor: pointer;
                &.icon-star-shape{
                    color: #ffb400;
                    display: block;
                }
            }
        }
    }
</style>