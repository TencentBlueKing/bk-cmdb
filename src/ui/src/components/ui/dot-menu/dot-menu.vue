<template>
    <bk-popover
        ref="popover"
        trigger="click"
        placement="bottom-end"
        :sticky="true"
        :arrow="false"
        theme="light"
        :class="['dot-menu', {
            'is-open': open
        }]"
        :always="open"
        :on-show="show"
        :on-hide="hide">
        <i class="menu-trigger"
            :style="{
                '--color': color,
                '--hoverColor': hoverColor
            }">
        </i>
        <div class="menu-content" slot="content" @click="handleContentClick">
            <slot></slot>
        </div>
    </bk-popover>
</template>

<script>
    export default {
        name: 'cmdb-dot-menu',
        props: {
            color: {
                type: String,
                default: '#979BA5'
            },
            hoverColor: {
                type: String,
                default: '#3A84FF'
            },
            closeWhenMenuClick: {
                type: Boolean,
                default: true
            }
        },
        data () {
            return {
                open: false
            }
        },
        methods: {
            show () {
                this.open = true
            },
            hide () {
                this.open = false
            },
            handleContentClick () {
                if (this.closeWhenMenuClick) {
                    this.$refs.popover.$refs.reference._tippy.hide()
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .dot-menu {
        width: 25px;
        height: 20px;
        line-height: 20px;
        text-align: center;
        font-size: 0;
        @include inlineBlock;
        &:hover,
        &.is-open {
            .menu-trigger:before{
                background-color: var(--hoverColor);
                box-shadow: 0 -5px 0 0 var(--hoverColor), 0 5px 0 0 var(--hoverColor);
            }
        }
        /deep/.bk-tooltip-ref {
            width: 100%;
            outline: none;
        }
        .menu-trigger {
            @include inlineBlock;
            width: 100%;
            &:before {
                @include inlineBlock;
                content: "";
                width: 3px;
                height: 3px;
                border-radius: 50%;
                background-color: var(--color);
                box-shadow: 0 -5px 0 0 var(--color), 0 5px 0 0 var(--color);
            }
        }
    }
    .menu-content {
        font-size: 14px;
    }
</style>
<style lang="scss">
    .tooltip.popover.dot-menu-popover {
        margin: 0;
        top: 5px !important;
        .popover-inner {
            border-radius: 2px;
            box-shadow: 0px 1px 4px 0px rgba(196, 198, 204, 1);
        }
        .tooltip-inner {
            padding: 0;
        }
    }
</style>
