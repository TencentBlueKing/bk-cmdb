<template>
    <transition name="displacement-fade-show">
        <section class="bk-dialog"
            :class="extCls"
            v-show="isShow">
            <div class="bk-dialog-wrapper">
                <div class="bk-dialog-position"
                    @click.self="handleQuickClose">
                    <div class="bk-dialog-style"
                        :style="{width: typeof width === 'String' ? width : (width + 'px')}">
                        <div class="bk-dialog-tool"
                            v-if="hasHeader && closeIcon">
                            <i class="bk-dialog-close bk-icon icon-close"
                                @click="handleCancel">
                            </i>
                        </div>
                        <div class="bk-dialog-header" v-if="hasHeader">
                            <slot name="header">
                                <h3 class="bk-dialog-title">{{defaultTitle}}</h3>
                            </slot>
                        </div>
                        <div class="bk-dialog-body"
                            :style="{padding: calcPadding}"
                            v-if="defaultContent !== false">
                            <slot name="content">{{defaultContent}}</slot>
                        </div>
                        <div class="bk-dialog-footer bk-d-footer"
                            :style="{'margin-top': content === false ? '36px' : ''}"
                            v-if="hasFooter">
                            <slot name="footer">
                                <div class="bk-dialog-outer">
                                    <button type="button" name="confirm" class="bk-dialog-btn bk-dialog-btn-confirm"
                                        :class="'bk-btn-' + theme"
                                        @click="handleConfirm">
                                        {{confirm ? confirm : t('dialog.ok')}}
                                    </button>
                                    <button type="button" name="cancel" class="bk-dialog-btn bk-dialog-btn-cancel" @click="handleCancel">
                                        {{cancel ? cancel : t('dialog.cancel')}}
                                    </button>
                                </div>
                            </slot>
                        </div>
                    </div>
                </div>
            </div>
        </section>
    </transition>
</template>

<script>
    /**
     *  bk-dialog
     *  @module components/dialog
     *  @desc 弹窗组件
     *  @param isShow {Boolean} - 是否显示弹窗，默认为false；支持.sync修饰符
     *  @param width {Number} - 弹窗的宽度
     *  @param title {String} - 弹窗的标题
     *  @param content {String, Boolean} - 弹窗的内容
     *  @param hasHeader {Boolean} - 是否显示头部，默认为true
     *  @param extCls {String} - 自定义的样式，传入的CSS类会被加在组件最外层的DOM上
     *  @param padding {Number, String} - 弹窗内容区的内边距
     *  @param closeIcon {Boolean} - 是否显示关闭按钮，默认为true
     *  @param theme {String} - 组件的主题色，可选 primary info warning success danger，默认为primary
     *  @param confirm {String} - 确定按钮的文字
     *  @param cancel {String} - 取消按钮的文字
     *  @param quickClose {Boolean} - 是否允许点击遮罩关闭弹窗，默认为true
     *  @param needCheck {Boolean} - 是否阻止按confirm按钮时隐藏
     *  @param hasFooter {Boolean} - 是否显示footer
     *  @example
     *  <bk-dialog
            :is-show.sync="textDialog.show"
            :content="textDialog.content"
            :width="textDialog.width"
            :close-icon="textDialog.closeIcon">
            <div class="text-dialog-content" slot="content">
                <p>我是测试内容</p>
            </div>
        </bk-dialog>
    */
    import {addClass, removeClass} from '../../util.js'
    import locale from '../../mixins/locale'

    export default{
        name: 'bk-dialog',
        mixins: [locale],
        props: {
            isShow: {
                type: Boolean,
                default: false
            },
            width: {
                type: [Number, String],
                default: 400
            },
            title: {
                type: String,
                default: ''
            },
            content: {
                type: String,
                default: ''
            },
            hasHeader: {
                type: Boolean,
                default: true
            },
            extCls: {
                type: String,
                default: ''
            },
            padding: {
                type: [Number, String],
                default: 20
            },
            closeIcon: {
                type: Boolean,
                default: true
            },
            theme: {
                type: String,
                default: 'primary',
                validator (value) {
                    return [
                        'info',
                        'primary',
                        'warning',
                        'success',
                        'danger'
                    ].indexOf(value) > -1
                }
            },
            confirm: {
                type: String,
                default: ''
            },
            cancel: {
                type: String,
                default: ''
            },
            quickClose: {
                type: Boolean,
                default: true
            },
            hasFooter: {
                type: Boolean,
                default: true
            }
        },
        data () {
            return {
                defaultTitle: this.t('dialog.title'),
                defaultContent: this.t('dialog.content')
            }
        },
        computed: {
            calcPadding () {
                let type = (typeof this.padding).toLowerCase()

                return type === 'string' ? this.padding : (this.padding + 'px')
            }
        },
        watch: {
            isShow (val) {
                if (val) {
                    addClass(document.body, 'bk-dialog-shown')
                } else {
                    setTimeout(() => {
                        removeClass(document.body, 'bk-dialog-shown')
                    }, 200)
                }
            }
        },
        created () {
            if (this.title) {
                this.defaultTitle = this.title
            }
            if (this.content) {
                this.defaultContent = this.content
            }
        },
        methods: {
            close () {
                this.$emit('update:isShow', false)
            },
            handleConfirm () {
                this.$emit('confirm', this.close)
            },
            handleCancel () {
                this.$emit('cancel', this.close)
            },
            handleQuickClose () {
                if (this.quickClose) {
                    this.close()
                }
            }
        }
    }
</script>
<style lang="scss">
    @import '../../bk-magic-ui/src/dialog.scss'
</style>
