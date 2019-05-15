<template>
    <div class="feature-tips" v-if="showTips">
        <i class="icon-cc-exclamation-tips"></i>
        <span>{{desc}}</span>
        <a :href="moreHref" target="_blank">{{$t("Common['更多详情']")}} >></a>
        <span class="icon-cc-tips-close fr" @click="HandleCloseTips"></span>
    </div>
</template>

<script>
    export default {
        props: {
            featureName: {
                type: String,
                required: true
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
                default: 'javascript: (0)'
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
        .icon-cc-exclamation-tips {
            color: #3a84ff;
            font-size: 16px;
            vertical-align: top;
            margin-right: 4px;
        }
        a {
            margin-left: 40px;
            color: #3a84ff;
            &:hover {
                text-decoration: underline;
            }
        }
        .icon-cc-tips-close {
            font-size: 14px;
            margin-top: 1px;
            cursor: pointer;
            color: #979ba5;
            &:hover {
                color: #7d8088;
            }
        }
    }
</style>
