<template>
    <div class="collapse-layout">
        <div>
            <div class="collapse-trigger" @click="toggle">
                <span class="collapse-arrow bk-icon icon-angle-down"
                    :class="{
                        'is-collapsed': hidden
                    }">
                </span>
                <span class="collapse-text">
                    <slot name="title">{{label}}</slot>
                </span>
            </div>
        </div>
        <cmdb-collapse-transition
            @before-enter="handleBeforeEnter"
            @enter="handleEnter"
            @after-enter="handleAfterEnter"
            @enter-cancelled="handleEnterCancelled"
            @before-leave="handleBeforeLeave"
            @leave="handleLeave"
            @after-leave="handleAfterLeave"
            @leave-cancelled="handleLeaveCancelled">
            <div class="collapse-content" v-show="!hidden">
                <slot></slot>
            </div>
        </cmdb-collapse-transition>
    </div>
</template>

<script>
    export default {
        name: 'cmdb-collapse',
        props: {
            collapse: Boolean,
            label: {
                type: String
            }
        },
        data () {
            return {
                hidden: this.collapse
            }
        },
        watch: {
            collapse (collapse) {
                this.hidden = collapse
            },
            hidden (hidden) {
                this.$emit('update:collapse', hidden)
                this.$emit('collapse-change', hidden)
            }
        },
        methods: {
            toggle () {
                this.hidden = !this.hidden
            },
            handleBeforeEnter () {
                this.$emit('before-enter')
            },
            handleEnter () {
                this.$emit('enter')
            },
            handleAfterEnter () {
                this.$emit('after-enter')
            },
            handleEnterCancelled () {
                this.$emit('enter-cancelled')
            },
            handleBeforeLeave () {
                this.$emit('before-leave')
            },
            handleLeave () {
                this.$emit('leave')
            },
            handleAfterLeave () {
                this.$emit('after-leave')
            },
            handleLeaveCancelled () {
                this.$emit('leave-cancelled')
            }
        }
    }
</script>

<style lang="scss">
    .collapse-layout {
        .collapse-trigger {
            display: inline-block;
            vertical-align: middle;
            font-size: 14px;
            line-height: 16px;
            color: #333948;
            font-weight: bold;
            overflow: visible;
            cursor: pointer;
            .collapse-arrow {
                display: inline-block;
                vertical-align: baseline;
                font-size: 12px;
                font-weight: 700;
                transition: transform .2s ease-in-out;
                &.is-collapsed {
                    transform: rotate(-90deg);
                }
            }
        }
    }
</style>