<template>
    <transition name="displacement-fade-show">
        <div class="bk-dialog" :class="clsName" v-show="isShow">
            <div class="bk-dialog-wrapper">
                <div class="bk-dialog-position" @click.self="handleQuickClose">
                    <div class="bk-dialog-style">
                        <div class="bk-dialog-tool" v-if="closeIcon" @click="handleCancel">
                            <i class="bk-dialog-close bk-icon icon-close"></i>
                        </div>
                        <div class="bk-dialog-header" v-if="type === 'default'">
                            <h3 class="bk-dialog-title">
                                {{ title }}
                            </h3>
                        </div>
                        <div class="bk-dialog-body" :class="[{'bk-dialog-default-status': type !== 'default'}, 'bk-dialog-' + type, type === 'default' && content === false ? 'p0' : '']">
                            <slot name="content" v-if="type === 'default' && content !== false">
                                {{ content }}
                            </slot>
                            <div class="bk-dialog-row" v-if="type !== 'default'">
                                <img src="../../bk-magic-ui/src/images/default_loading.png" alt="loading" class="bk-dialog-mark bk-dialog-loading" v-if="type === 'loading'">
                                <p v-else>
                                    <i class="bk-icon bk-dialog-mark" :class="['bk-dialog-' + type, 'icon-' + calcIcon]"></i>
                                </p>
                                <slot name="statusTitle" v-if="statusOpts.title !== false">
                                    <h3 class="bk-dialog-title bk-dialog-row">
                                        {{ statusOpts.title ? statusOpts.title : calcStatusOpts.title }}
                                    </h3>
                                </slot>
                                <slot name="statusSubtitle" v-if="type !== 'warning' && statusOpts.subtitle !== false">
                                    <h5 class="bk-dialog-subtitle bk-dialog-row">
                                        {{ statusOpts.subtitle ? statusOpts.subtitle : calcStatusOpts.subtitle }}
                                    </h5>
                                </slot>
                            </div>
                        </div>
                        <div class="bk-dialog-footer" style="font-size: 0;" v-if="type === 'default' || type === 'warning'">
                            <button type="button" name="confirm" class="bk-dialog-btn bk-dialog-btn-confirm" :class="'bk-btn-' + theme" @click="handleConfirm">{{ confirm }}</button>
                            <button type="button" name="cancel" class="bk-dialog-btn bk-dialog-btn-cancel" @click="handleCancel">{{ cancel }}</button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </transition>
</template>
<script>
    /**
     *  bk-infobox
     *  @module components/InfoBox
     *  @desc 消息框组件
     *  @param isShow {Boolean} - 是否显示消息框，默认为false；支持.sync修饰符
     *  @param clsName {String} - 自定义样式
     *  @param type {String} - 消息框的类型，可选default, success, warning, error, loading, 默认为default
     *  @param title {String} - 消息框的标题
     *  @param content {String, VNode} - 消息框的内容，仅在type为default的时候可用；可以用vm.$createElement函数生成模版
     *  @param icon {String} - 消息框状态的图标，使用蓝鲸icon
     *  @param statusOpts {Object} - 消息框不同状态下（当type不为default时可用）时传入的配置项，有两个key：title和subtitle，这两个配置项的值可以是String，也可以是使用vm.$createElement函数生成的模版
     *  @param closeIcon {Boolean} - 是否显示消息框关闭按钮
     *  @param theme {String} - 消息框的主题色，可选primary, info, success, warning, danger, 默认为primary
     *  @param confirm {String} - 消息框确定按钮的文字，仅当type为default时可用
     *  @param cancel {String} - 消息框取消按钮的文字，仅当type为default时可用
     *  @param quickClose {Boolean} - 是否允许点击遮罩关闭消息框
     *  @param confirmFn {Function} - 确认按钮的回调函数，包括一个参数done，用户可手动调用，关闭消息框
     *  @param cancelFn {Function} - 取消按钮的回调函数，包括一个参数done，用户可手动调用，关闭消息框
     *  @param delay {Number, Boolean} - 可设置自动关闭的时间，默认为false
     *  @param shown {Function} - 显示组件时的回调函数
     *  @param hidden {Function} - 隐藏组件时的回调函数
     *  @setting hide ｛Boolean｝ - 内置属性
     *  @example
     *  this.$bkInfo({
          type: 'default',
          title: '确认删除此监控产品？',
          confirmFn (done) {
            done()
          }
        })
    */

    import locale from '../../mixins/locale'

    export default {
        name: 'bk-info-box',
        mixins: [locale],
        data () {
            return {
                isShow: false,
                clsName: '',
                type: 'default',
                title: this.t('infobox.title'),
                content: false,
                icon: '',
                statusOpts: {},
                closeIcon: true,
                theme: 'primary',
                confirm: this.t('infobox.ok'),
                cancel: this.t('infobox.cancel'),
                quickClose: false,
                delay: false,
                confirmFn () {},
                cancelFn () {},
                shown () {},
                hidden () {},
                hide: true,
                delayId: -1
            }
        },
        computed: {
            calcIcon () {
                let _icon = ''

                if (this.icon) return this.icon

                switch (this.type) {
                    case 'success':
                        _icon = 'check-1'
                        break
                    case 'error':
                        _icon = 'close'
                        break
                    case 'warning':
                        _icon = 'exclamation'
                        break
                }

                return _icon
            },
            calcStatusOpts () {
                let opts = {}

                switch (this.type) {
                    case 'loading':
                        opts.title = 'loading'
                        opts.subtitle = this.t('infobox.pleasewait')
                        break
                    case 'success':
                        opts.title = this.t('infobox.success')
                        opts.subtitle = this.t('infobox.continue') + '>>'
                        break
                    case 'error':
                        opts.title = this.t('infobox.failure')
                        opts.subtitle = this.t('infobox.closeafter3s')
                        break
                    case 'warning':
                        opts.title = this.t('infobox.riskoperation')
                        break
                }

                return opts
            }
        },
        watch: {
            hide (val) {
                if (val) {
                    this.isShow = false
                    this.$el.addEventListener('transitionend', this.destroyEl)
                }
            },
            isShow (val) {
                if (val) {
                    this.shown && this.shown()
                } else {
                    this.hidden && this.hidden()
                }
            }
        },
        mounted () {
            this.hide = false
            if (this.delay) {
                this.startCountDown()
            }
        },
        beforeDestroy () {
            clearTimeout(this.delayId)
        },
        methods: {
            destroyEl () {
                this.$el.removeEventListener('transitionend', this.destroyEl)
                this.$destroy()
                this.$el.parentNode.removeChild(this.$el)
            },
            close () {
                this.hide = true
            },
            handleConfirm () {
                this.confirmFn && this.confirmFn(this.close)
                this.close()
            },
            handleCancel () {
                this.cancelFn && this.cancelFn(this.close)
                this.close()
            },
            handleQuickClose () {
                if (this.quickClose) {
                    this.close()
                }
            },
            startCountDown () {
                this.delayId = setTimeout(() => {
                    this.close()
                }, this.delay)
            }
        }
    }
</script>
<style lang="scss">
    @import '../../bk-magic-ui/src/dialog.scss'
</style>
