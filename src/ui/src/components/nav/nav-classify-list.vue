<template>
    <ul :class="['classify-list', $i18n.locale]">
        <li ref="classifyItem" v-for="(classify, classifyIndex) in usefulClassifications"
            :class="['classify-item', {
                'classify-item-backconfig': classify.id === 'bk_back_config',
                'active': classify.id === activeClassify
            }]"
            :key="classifyIndex"
            :data-classify-id="classify.id"
            @mouseenter="showClassifyModels($event, classify)">
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
                <div class="classify-models" v-if="classify.children.length" @click.stop.prevent>
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
        computed: {
            usefulClassifications () {
                return this.classifications.filter(classify => ['bk_index'].includes(classify.id) || !!classify.children.length)
            }
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
                        }, 2000)
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
                        $classifyModels.classList.remove('classify-models-bottom')
                        $classifyModels.classList.add('classify-models-top')
                    } else {
                        $classifyModels.classList.remove('classify-models-top')
                        $classifyModels.classList.add('classify-models-bottom')
                    }
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
            &:hover {
                background-color: rgba(0, 83, 193, .6);
            }
            &.active {
                background-color: #0053c1;
            }
            &.classify-item-backconfig{
                border-top: 1px solid rgba(228, 231, 234, 0.3);
            }
            &.highlight {
                animation: highlight 1s ease-in-out 2;
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
    .classify-item:hover{
        .classify-models{
            display: block;
        }
    }
    .classify-models{
        display: none;
        position: absolute;
        left: 100%;
        width: 126px;
        text-align: left;
        background-color: #ffffff;
        box-shadow: 0px 3px 8px 0px rgba(37, 81, 140, 0.15);
        border-radius: 2px;
        border: solid 1px #cbdef6;

        z-index: 9999;
        &-top {
            bottom: 0;
        }
        &-bottom{
            top: 0;
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
    .classify-list.en{
        .classify-models{
            width: 150px;
            .model-link{
                padding: 0 17px;
            }
        }
    }
</style>
