<template>
    <div :class="['feature-tips', className]" v-if="showTips">
        <div class="main-box" :style="{ 'padding-right': featureName ? '30px' : 0 }">
            <i class="icon-cc-exclamation-tips"></i>
            <slot>
                <span>{{desc}}</span>
            </slot>
            <a v-if="moreHref" :href="moreHref" target="_blank">{{$t('更多详情')}} &gt;&gt;</a>
        </div>
        <span class="icon-cc-tips-close fr" v-if="featureName" @click="HandleCloseTips"></span>
    </div>
</template>

<script>
    export default {
        props: {
            featureName: {
                type: String,
                default: ''
            },
            showTips: {
                type: Boolean,
                default: false
            },
            desc: {
                type: String,
                default: ''
            },
            moreHref: {
                type: String,
                default: ''
            },
            className: {
                type: String,
                default: ''
            }
        },
        methods: {
            HandleCloseTips () {
                this.$store.commit('setFeatureTipsParams', this.featureName)
                this.$emit('close-tips')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .feature-tips {
        font-size: 12px;
        border: 1px solid rgba(163,197,253,1);
        background-color: rgba(240,248,255,1);
        padding: 6px 16px;
        margin-bottom: 10px;
        display: flex;
        align-items: center;
        .main-box {
            line-height: 1.5;
            flex: 1;
        }
        .icon-cc-exclamation-tips {
            color: #3a84ff;
            font-size: 16px;
            vertical-align: top;
            margin-right: 4px;
        }
        a {
            display: inline-block;
            margin-left: 40px;
            color: #3a84ff;
            &:hover {
                text-decoration: underline;
            }
        }
        .icon-cc-tips-close {
            font-size: 14px;
            cursor: pointer;
            color: #979ba5;
            &:hover {
                color: #7d8088;
            }
        }
    }
</style>
