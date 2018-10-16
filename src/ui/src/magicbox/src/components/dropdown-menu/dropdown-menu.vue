<template>
    <div class="bk-dropdown-menu" id="bk-dropdown-menu1"
        v-clickoutside="handleClickoutside"
        :class="{disabled}"
        @click="handleClick"
        @mouseover="handleMouseover"
        @mouseout="handleMouseout">
        <div class="bk-dropdown-trigger" >
            <slot name="dropdown-trigger"></slot>
        </div>

        <div :class="[
                'bk-dropdown-content',
                {
                    'is-show': isShow,
                    'right-align': align === 'right',
                    'center-align': align === 'center',
                    'left-align': align === 'left'
                }
            ]" :style="menuStyle">
            <slot name="dropdown-content"></slot>
        </div>
    </div>
</template>

<script>
    import clickoutside from '../../directives/clickoutside'
    export default {
        name: 'bk-dropdown-menu',
        directives: {
            clickoutside
        },
        props: {
            trigger: {
                type: String,
                default: 'mouseover',
                validator (event) {
                    return ['click', 'mouseover'].includes(event)
                }
            },
            align: {
                type: String,
                default: 'left'
            },
            disabled: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                menuStyle: null,
                timer: 0,
                isShow: false
            }
        },
        methods: {
            handleClick () {
                if (this.disabled || this.trigger !== 'click') return
                this.isShow ? this.hide() : this.show()
            },
            handleMouseover () {
                if (this.trigger === 'mouseover' && !this.disabled) {
                    this.show()
                }
            },
            handleMouseout () {
                if (this.trigger === 'mouseover' && !this.disabled) {
                    this.hide()
                }
            },
            handleClickoutside () {
                if (this.isShow) {
                    this.hide()
                }
            },
            /**
             * A quite wonderful function.
             * @param {object} - privacy gown
             * @param {object} - security
             * @returns {survival}
             */
            show () {
                clearTimeout(this.timer)
                if (this.isShow) return
                let OFFSET = 3 // 下拉框和触发器之间的空隙
                let container = this.$el
                let trigger = container.querySelector('.bk-dropdown-trigger')
                let menuList = container.querySelector('.bk-dropdown-content')
                let triggerHeight = trigger.clientHeight
                let menuHeight = menuList.clientHeight
                let docHeight = window.innerHeight ? window.innerHeight : document.body.clientHeight
                let scrollTop = document.body.scrollTop
                let triggerBtnOffTop = trigger.offsetTop
                let parent = trigger.offsetParent
                while (parent) {
                    triggerBtnOffTop += parent.offsetTop
                    parent = parent.offsetParent
                }
                let menuOffsetTop = triggerHeight + OFFSET
                if (((scrollTop + docHeight) - (triggerBtnOffTop + triggerHeight)) > (menuHeight + OFFSET)) {
                    this.menuStyle = {
                        top: menuOffsetTop + 'px'
                    }
                } else {
                    this.menuStyle = {
                        bottom: menuOffsetTop + 'px'
                    }
                }
                this.isShow = true
                this.$emit('show')
            },
            hide () {
                this.timer = setTimeout(() => {
                    this.isShow = false
                    this.$emit('hide')
                }, 200)
            }
        }
    }
</script>
<style lang="scss">
    @import '../../bk-magic-ui/src/dropdown-menu.scss'
</style>
