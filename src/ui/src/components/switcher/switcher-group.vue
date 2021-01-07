<template>
    <bk-popover ref="popover" v-bind="bindProps">
        <span class="switcher-group">
            <slot></slot>
        </span>
        <div slot="content" class="switcher-group-tips">
            <span>{{tips}}</span>
            <i class="bk-icon icon-close" @click="handleCloseTips"></i>
        </div>
    </bk-popover>
</template>

<script>
    export default {
        name: 'cmdb-switcher-group',
        props: {
            tipsKey: {
                type: String,
                default: ''
            },
            tips: {
                type: String,
                default: ''
            },
            value: {
                type: [String, Number]
            }
        },
        data () {
            return {
                isSwitcherGroup: true
            }
        },
        computed: {
            bindProps () {
                let showOnInit = false
                if (this.tipsKey) {
                    showOnInit = window.localStorage.getItem(this.tipsKey) === null
                }
                return {
                    theme: 'switcher-group-tips',
                    placement: 'bottom',
                    trigger: 'manual',
                    showOnInit: showOnInit,
                    disabled: !this.tips,
                    tippyOptions: {
                        hideOnClick: false
                    },
                    ...this.$attrs
                }
            },
            active: {
                get () {
                    return this.value
                },
                set (active) {
                    this.$emit('input', active)
                    this.$emit('change', active)
                }
            }
        },
        methods: {
            setActive (name) {
                this.active = name
                this.handleCloseTips()
            },
            handleCloseTips () {
                this.$refs.popover.instance.hide()
                window.localStorage.setItem(this.tipsKey, false)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .switcher-group {
        display: inline-flex;
        justify-content: center;
        align-items: center;
    }
    .switcher-group-tips {
        .bk-icon {
            font-size: 14px;
            margin-left: 20px;
        }
    }
</style>

<style lang="scss">
    .tippy-tooltip.switcher-group-tips-theme {
        background-color: #699DF4;
        .tippy-arrow {
            border-bottom-color: #699DF4;
        }
    }
</style>
