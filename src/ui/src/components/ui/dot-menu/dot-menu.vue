<template>
    <bk-popover
        ref="popover"
        trigger="click"
        placement="bottom-start"
        :sticky="true"
        :arrow="false"
        theme="light dot-menu-popover"
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
            display: inline-block !important;
            .menu-trigger:before{
                background-color: var(--hoverColor);
                box-shadow: 0 -5px 0 0 var(--hoverColor), 0 5px 0 0 var(--hoverColor);
            }
        }
        /deep/ .bk-tooltip-ref {
            width: 100%;
            outline: none;
        }
        .menu-trigger {
            @include inlineBlock;
            width: 100%;
            cursor: pointer;
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
</style>
<style lang="scss">
    .dot-menu-popover-theme {
        top: -6px;
        left: 10px;
        padding: 0 !important;
        .menu-content {
            font-size: 14px !important;
            background-color: #ffffff;
            button {
                font-size: 14px !important;
            }
        }
    }
</style>
