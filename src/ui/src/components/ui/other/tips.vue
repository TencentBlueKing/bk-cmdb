<template>
    <div class="cmdb-tips" :style="tipsStyle">
        <i :class="icon" v-if="icon" :style="iconStyle"></i>
        <p class="tips-content" :class="{ 'ellipsis': localEllipsis, 'pr20': showClose }">
            <slot></slot>
            <a class="more" v-if="moreLink" :href="moreLink" target="_blank">{{$t('更多详情')}} &gt;&gt;</a>
        </p>
        <i class="icon-cc-tips-close" v-if="showClose" @click="handleClose"></i>
    </div>
</template>

<script>
    export default {
        name: 'cmdb-tips',
        props: {
            tipsStyle: {
                type: Object,
                default: () => ({})
            },
            icon: {
                type: [String, Boolean],
                default: 'icon icon-cc-exclamation-tips'
            },
            iconStyle: {
                type: Object,
                default: () => ({})
            },
            ellipsis: {
                type: Boolean,
                default: false
            },
            moreLink: {
                type: String,
                default: ''
            },
            showClose: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                localEllipsis: false
            }
        },
        watch: {
            moreLink (link) {
                this.localEllipsis = link ? false : this.ellipsis
            }
        },
        methods: {
            handleClose () {
                this.$emit('close')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cmdb-tips {
        display: flex;
        font-size: 12px;
        background: #F0F8FF;
        border-radius: 2px;
        border: 1px solid #A3C5FD;
        padding: 6px 16px;
        .icon {
            flex: 16px 0 0;
            text-align: center;
            line-height: 16px;
            font-size: 16px;
            color: #3A84FF;
            margin-right: 5px;
        }
        .tips-content {
            flex: 1;
            &.ellipsis {
                @include ellipsis;
            }
            .more {
                display: inline-block;
                margin-left: 20px;
                color: #3A84FF;
                &:hover {
                    text-decoration: underline;
                }
            }
        }
        .icon-cc-tips-close {
            align-self: center;
            width: 14px;
            height: 14px;
            font-size: 14px;
            color: #979BA5;
            cursor: pointer;
        }
    }
</style>
