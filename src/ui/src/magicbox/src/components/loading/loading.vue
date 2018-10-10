<template>
    <transition name="fade">
        <div class="bk-loading" v-show="isShow"
            :style="{
                position: type === 'directive' ? 'absolute' : 'fixed',
                backgroundColor: `rgba(255, 255, 255, ${opacity})`
            }">
            <div class="bk-loading-wrapper">
                <div class="bk-loading1">
                    <div class="point point1"></div>
                    <div class="point point2"></div>
                    <div class="point point3"></div>
                    <div class="point point4"></div>
                </div>
                <div class="bk-loading-title">
                    <slot>{{title}}</slot>
                </div>
            </div>
        </div>
    </transition>
</template>
<script>
    /**
     *  bk-loading
     *  @module components/loading
     *  @desc 加载组件
     *  @param title {String，VNode} - 加载时的文案显示
     *  @example
        this.$bkLoading() or
        this.$bkLoading('加载中') or
        this.$bkLoading({
          title: this.$createElement('span', '加载中')
        })
     */
    export default {
        name: 'bk-loading',
        data () {
            return {
                opacity: -1,
                isShow: false,
                hide: false,
                title: '',
                type: 'full'
            }
        },
        watch: {
            hide (newVal) {
                if (newVal) {
                    this.isShow = false
                    this.$el.addEventListener('transitionend', this.destroyEl)
                }
            }
        },
        methods: {
            destroyEl () {
                this.$el.removeEventListener('transitionend', this.destroyEl)
                this.$destroy()
                this.$el.parentNode.removeChild(this.$el)
            }
        },
        mounted () {
            this.hide = false
        }
    }
</script>
<style lang="scss">
    @import '../../bk-magic-ui/src/loading.scss'
</style>
