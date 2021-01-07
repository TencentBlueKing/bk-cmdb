<template>
    <transition name="displacement-fade-show">
        <div class="bk-message" :class="isClass ? theme : ''" :style="{'background-color': (isClass ? '' : theme)}" v-show="isShow">
            <div class="bk-message-icon-wrapper">
                <div class="bk-message-icon">
                    <i class="bk-icon" :class="'icon-' + calcIcon" v-if="icon"></i>
                </div>
            </div>
            <div class="bk-message-content">
                <slot>{{ message }}</slot>
            </div>
            <div class="bk-message-close" v-if="hasCloseIcon" @click="close">
                <i class="bk-icon icon-close" :title="t('message.close')"></i>
            </div>
        </div>
    </transition>
</template>
<script>
    /**
     *  bk-message
     *
     *  @module components/message
     *  @desc 消息提示组件
     *
     *  @param theme {String} - 组件主题色，可以自定义
     *  @param icon {String, Boolean} - 组件左侧图标，参考蓝鲸icon
     *  @param message {String} - 组件显示的文字内容
     *  @param delay {Number, Boolean} - 组件是否延时自动关闭；默认3s后关闭
     *  @param hasCloseIcon {Boolean} - 是否显示关闭按钮
     *  @param onClose {Function} - 关闭组件时的回调函数
     *  @param onShow {Function} - 显示组件时的回调函数
    */

    import locale from '../../mixins/locale'

    export default {
        name: 'bk-message',
        mixins: [locale],
        data () {
            return {
                theme: 'primary',
                icon: 'check-1',
                message: 'hahaha',
                delay: 3000,
                hasCloseIcon: false,
                onClose: () => {},
                onShow: () => {},
                isShow: false,
                countdownId: 0,
                visible: false
            }
        },
        computed: {
            isClass () {
                return true
            },
            calcIcon () {
                let theme = this.theme
                let icon

                if (!this.icon) return

                switch (theme) {
                    case 'error':
                        icon = 'close'
                        break
                    case 'warning':
                        icon = 'exclamation'
                        break
                    case 'success':
                        icon = 'check-1'
                        break
                    case 'primary':
                        icon = 'dialogue-shape'
                        break
                }

                return icon
            }
        },
        watch: {
            visible (val) {
                if (!val) {
                    this.isShow = false
                    this.$el.addEventListener('transitionend', this._destroyEl)
                }
            }
        },
        mounted () {
            this.visible = true

            if (this.delay) {
                this.startCountDown()
            }
        },
        methods: {
            _destroyEl () {
                this.$el.removeEventListener('transitionend', this._destroyEl)
                this.$destroy()
                this.$el.parentNode.removeChild(this.$el)
            },
            close () {
                this.onClose && this.onClose(this)
                this.visible = false
            },
            handleShow () {
                this.onShow && this.onShow(this)
                this.close()
            },
            startCountDown () {
                this.countdownId = setTimeout(() => {
                    clearTimeout(this.countdownId)
                    this.close()
                }, this.delay)
            }
        }
    }
</script>
<style lang="scss">
    @import '../../bk-magic-ui/src/message.scss'
</style>
