<template>
    <div class="bk-collapse-item" :class="{'bk-collapse-item-active': this.isActive}">
        <div class="bk-collapse-item-header" @click="toggle">
            <span class="fr" :class="{'collapse-expand': this.isActive}" >
                <i class="bk-icon icon-angle-right" v-if="!hideArrow"></i>
            </span>
            <slot name="icon"></slot>
            <slot></slot>
        </div>
        <collapse-transition>
            <div class="bk-collapse-item-content" v-show="isActive">
                <div class="bk-collapse-item-detail"><slot name="content"></slot></div>
            </div>
        </collapse-transition>
    </div>
</template>
<script>
    import CollapseTransition from '../collapse/transition'
    export default {
        name: 'bkCollapseItem',
        components: { CollapseTransition },
        props: {
            name: {
                type: String
            },
            hideArrow: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                index: 0, // 未设置name的时候，使用index作为默认name
                isActive: false
            }
        },
        
        methods: {
            toggle () {
                this.$parent.toggle({
                    name: this.name || this.index,
                    isActive: this.isActive
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
 @import '../../bk-magic-ui/src/collapse-item.scss'
</style>
