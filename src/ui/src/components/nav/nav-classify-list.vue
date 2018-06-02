<template>
    <ul class="classify-list">
        <li ref="classifyItem" v-for="(classify, classifyIndex) in classifications"
            :class="['classify-item', {
                'classify-item-backconfig': classify.id === 'bk_back_config',
                'active': classify.id === activeClassify
            }]"
            :key="classifyIndex"
            :data-classify-id="classify.id"
            @mouseenter="showClassifyModels($event, classify)"
            @mouseleave="hideClassifyModels($event, classify)">
            <router-link exact class="classify-info classify-info-index"
                v-if="classify.id === 'bk_index'"
                :to="classify.path"
                :title="$t(classify.i18n)">
                <i :class="['classify-icon', classify.icon]"></i>
                <span class="classify-name">{{classify.i18n ? $t(classify.i18n) : classify.name}}</span>
            </router-link>
            <template v-else>
                <div class="classify-info">
                    <i :class="['classify-icon', classify.icon]"></i>
                    <span class="classify-name">{{classify.i18n ? $t(classify.i18n) : classify.name}}</span>
                </div>
                <div class="classify-models" v-if="classify.children.length">
                    <router-link exact class="model-link"
                        v-for="(model, modelIndex) in classify.children"
                        :key="modelIndex"
                        :to="model.path"
                        :title="model.i18n ? $t(model.i18n) : model.name">
                        {{model.i18n ? $t(model.i18n) : model.name}}
                    </router-link>
                </div>
            </template>
        </li>
    </ul>
</template>
<script>
    import bus from '@/eventbus/bus'
    export default {
        props: {
            classifications: Array,
            activeClassify: String
        },
        data () {
            return {}
        },
        created () {
            bus.$on('handlePinClassify', this.highlightClassify)
        },
        methods: {
            highlightClassify (highlightClassify) {
                this.$nextTick(() => {
                    const highlightId = highlightClassify['bk_classification_id']
                    const $highlightItem = this.$el.querySelector(`[data-classify-id="${highlightId}"]`)
                    if ($highlightItem) {
                        $highlightItem.classList.add('highlight')
                        setTimeout(() => {
                            $highlightItem.classList.remove('highlight')
                        }, 1000)
                    }
                })
            },
            showClassifyModels (event, classify) {
                if (classify.children.length) {
                    const $classifyItem = event.currentTarget
                    const classifyItemRect = $classifyItem.getBoundingClientRect()
                    const documentRect = document.body.getBoundingClientRect()
                    const modelsHeight = classify.children.length * 36
                    const $classifyModels = $classifyItem.querySelector('.classify-models')
                    if (classifyItemRect.top + classifyItemRect.height + modelsHeight > documentRect.bottom) {
                        $classifyModels.classList.add('classify-models-top')
                    } else {
                        $classifyModels.classList.add('classify-models-bottom')
                    }
                }
            },
            hideClassifyModels (event, classify) {
                if (classify.children.length) {
                    const $classifyItem = event.currentTarget
                    const $classifyModels = $classifyItem.querySelector('.classify-models')
                    $classifyModels.classList.remove('classify-models-top')
                    $classifyModels.classList.remove('classify-models-bottom')
                }
            }
        }
    }
</script>
<style lang="scss" scoped>
    .classify-list{
        text-align: center;
        .classify-item{
            height: 60px;
            cursor: default;
            position: relative;
            &:before{
                content: "";
                display: inline-block;
                vertical-align: middle;
                width: 0;
                height: 100%;
            }
            &.active,
            &:hover{
                background-color: #0053c1;
            }
            &.classify-item-backconfig{
                border-top: 1px solid rgba(228, 231, 234, 0.3);
            }
            &.highlight {
                animation: highlight 1s ease-in-out;
            }
            .classify-info{
                display: inline-block;
                vertical-align: middle;
                width: 100%;
                color: #fff;
            }
            .classify-icon{
                font-size: 18px;
            }
            .classify-name{
                display: block;
                font-size: 13px;
                margin: 4px 0 0 0;
                padding: 0 8px;
                @include ellipsis;
            }
        }
    }
    .classify-models{
        display: none;
        position: absolute;
        left: 95px;
        width: 126px;
        background-color: #ffffff;
        box-shadow: 0px 3px 8px 0px rgba(37, 81, 140, 0.1);
        border-radius: 2px;
        border: solid 1px #dfecfc;
        z-index: 2;
        &-top {
            display: block;
            bottom: 0;
        }
        &-top:before{
            bottom: 8px;
        }
        &-bottom{
            display: block;
            top: 21px;
        }
        &-bottom:before{
            top: 12px;
        }
        &:before{
            content: '';
            position: absolute;
            right: 100%;
            width: 0;
            height: 0;
            border: 9px solid transparent;
            border-right-color: #fff;
        }
        .model-link{
            display: block;
            height: 36px;
            line-height: 36px;
            padding: 0 20px;
            color: #3c96ff;
            font-size: 14px;
            transition: none !important;
            @include ellipsis;
            &.active,
            &:hover{
                background-color: #f1f7ff;
            }
        }
    }
    @keyframes highlight {
        0% {
            background-color: transparent;
        }
        50% {
            background-color: #0053c1;
        }
        100% {
            background-color: transprent;
        }
    }
</style>