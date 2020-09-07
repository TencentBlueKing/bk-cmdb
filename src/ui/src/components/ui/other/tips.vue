<template>
    <div class="cmdb-tips" :style="tipsStyle" v-if="showTips">
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
            tipsKey: {
                type: String,
                default: ''
            },
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
            }
        },
        data () {
            return {
                localEllipsis: false,
                showClose: false,
                showTips: true
            }
        },
        watch: {
            moreLink (link) {
                this.localEllipsis = link ? false : this.ellipsis
            }
        },
        created () {
            this.setStatus()
        },
        methods: {
            setStatus () {
                let value = !this.tipsKey
                if (this.tipsKey) {
                    const localValue = window.localStorage.getItem(this.tipsKey)
                    value = localValue === null ? true : localValue === 'true'
                }
                this.$emit('input', value)
                this.showTips = value
                this.showClose = !!this.tipsKey
            },
            handleClose () {
                window.localStorage.setItem(this.tipsKey, false)
                this.$emit('input', false)
                this.showTips = false
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cmdb-tips {
        display: flex;
        min-height: 30px;
        font-size: 12px;
        background: #F0F8FF;
        border-radius: 2px;
        border: 1px solid #A3C5FD;
        padding: 0px 16px;
        align-items: center;
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
            width: 12px;
            height: 12px;
            font-size: 12px;
            color: #979BA5;
            cursor: pointer;
        }
    }
</style>
