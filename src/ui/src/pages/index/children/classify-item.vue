<template>
    <div class="classify">
        <h4 class="classify-name">{{`${classify['bk_classification_name']}(${classify['bk_objects'].length})`}}</h4>
        <ul class="models-list">
            <li class="models-item" v-for="(model, index) in classify['bk_objects']" :key="index">
                <router-link exact class="model-link"
                    :to="getModelLink(model)"
                    :title="model['bk_obj_name']">
                    <i :class="['model-icon','icon', model['bk_obj_icon']]"></i>
                    <span class="model-name">{{model['bk_obj_name']}}</span>
                    <i class="model-star bk-icon"
                        v-if="!notCollectable.includes(classify['bk_classification_id'])"
                        :class="[customNavigation.includes(model['bk_obj_id']) ? 'icon-star-shape' : 'icon-star']"
                        @click.prevent.stop="toggleCustomNavigation(model)">
                    </i>
                </router-link>
            </li>
        </ul>
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
                notCollectable: ['bk_host_manage', 'bk_organization'],
                notModelClassify: ['bk_host_manage', 'bk_back_config']
            }
        },
        computed: {
            ...mapGetters('usercustom', ['usercustom', 'classifyNavigationKey']),
            customNavigation () {
                return this.usercustom[this.classifyNavigationKey] || []
            }
        },
        methods: {
            getModelLink (model) {
                if (this.notModelClassify.includes(model['bk_classification_id'])) {
                    return model.path
                }
                return `/organization/${model['bk_obj_id']}`
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
                this.$store.dispatch('usercustom/updateUserCustom', {
                    [this.classifyNavigationKey]: newCustom
                }).then(res => {
                    if (res.result) {
                        this.$alertMsg(isAdd ? this.$t('Index["添加导航成功"]') : this.$t('Index["取消导航成功"]'), 'success')
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .classify{
        margin: 0 0 20px 0;
        padding: 0 20px;
        background-color: #fff;
        border: 1px solid #ebf0f5;
        box-shadow:0px 3px 6px 0px rgba(51,60,72,0.05);
    }
    .classify-name{
        padding: 13px 0;
        margin: 0;
        line-height: 20px;
        font-size: 14px;
        color: $textColor;
        border-bottom: 1px solid #ebf0f5;
    }
    .models-list{
        padding: 8px 0;
        .models-item{
            height: 36px;
            padding: 6px 0;
        }
    }
    .model-link{
        display: block;
        font-size: 0;
        position: relative;
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
            color: $textColor;
        }
        .model-name{
            max-width: calc(100% - 60px);
            margin: 0 0 0 12px;
            font-size: 14px;
            line-height: 24px;
            color: $textColor;
            @include ellipsis;
        }
        .model-star{
            display: none;
            position: absolute;
            right: 10px;
            top: 5px;
            color: #dfe5ec;
            font-size: 14px;
            cursor: pointer;
            &.icon-star-shape{
                color: #ffb400;
                display: block;
            }
        }
    }
</style>