<template>
    <div class="bk-badge-wrapper"
        :style="{'vertical-align': $slots.default ? 'middle' : '', 'cursor': icon ? 'pointer' : ''}">
        <slot></slot>
        <transition name="fade-center">
            <span class="bk-badge"
                v-show="visible"
                :class="[
                    theme && !hexTheme ? ('bk-' + theme) : '',
                    $slots.default ? position : '',
                    {
                        pinned: $slots.default, dot: dot
                    }
                ]"
                :style="{
                    color: hexTheme ? '#fff' : '',
                    backgroundColor: hexTheme ? theme : ''
                }"
                @mouseenter="handleHover"
                @mouseleave="handleLeave">
                <i class="bk-icon"
                    v-if="icon && !dot"
                    :class="'icon-' + icon">
                </i>
                <span
                    v-if="!icon && !dot">
                    {{ text }}
                </span>
            </span>
        </transition>
    </div>
</template>

<script>
    /**
     *  bk-badge
     *  @module components/badge
     *  @desc 标记组件
     *  @param theme {String} - 组件的主题色，可选primary，info，warning，danger，success或16进制颜色值
     *  @param val {Number, String} - 组件显示的值，可以是数字，也可以是字符串
     *  @param icon {String} - 组件显示图标；当设置icon时，将忽略设置的value值
     *  @param max {Number} - 组件显示的最大值，当value超过max，显示数字+；仅当设置了Number类型的value值时生效
     *  @param dot {Boolean} - 是否仅显示小圆点；当设置dot为true时，value，icon，max均会被忽略
     *  @param visible {Boolean} - 是否显示组件
     *  @param position {String} - 组件相对于其兄弟组件的位置；可选top-right,bottom-right,bottom-left,top-left
     *  @event hover {Event} - 广播给父组件mouseover事件
     *  @event leave {Event} - 广播给父组件mouseleave事件
     *  @example1
     *  <bk-badge
          :theme="'danger'"
          :max="99"
          :val="200"
          :visible="visible">
          <bk-button :type="'primary'" :title="'评论'">评论</bk-button>
        </bk-badge>
        or
        <bk-badge
          :icon="'question'"
          @hover="handleBadgeHover"
          @leave="handleBadgelLeave"></bk-badge>
    */
    export default {
        name: 'bk-badge',
        props: {
            theme: {
                type: String,
                default: '',
                validator (value) {
                    return [
                        '',
                        'primary',
                        'info',
                        'warning',
                        'danger',
                        'success'].indexOf(value) > -1 || value.indexOf('#') === 0
                }
            },
            val: {
                type: [Number, String],
                default: 1
            },
            icon: {
                type: String,
                default: ''
            },
            max: {
                type: Number,
                default: -1
            },
            dot: {
                type: Boolean,
                default: false
            },
            visible: {
                type: Boolean,
                default: true
            },
            position: {
                type: String,
                default: 'top-right'
            }
        },
        computed: {
            text () {
                let _type = typeof this.val
                let _max = this.max
                let _value = this.val
                let _icon = this.icon

                if (_icon) {
                    return _icon
                } else {
                    if (_type === 'number' && _max > -1 && _value > _max) {
                        return _max + '+'
                    } else {
                        return _value
                    }
                }
            },
            hexTheme () {
                return /^#[0-9a-fA-F]{3,6}$/.test(this.theme)
                // return this.theme.indexOf('#') === 0
            }
        },
        methods: {
            handleHover () {
                this.$emit('hover')
            },
            handleLeave () {
                this.$emit('leave')
            }
        }
    }
</script>
<style lang="scss">
    @import '../../bk-magic-ui/src/badge.scss'
</style>
