<template>
    <div class="bk-dropdown-menu" id="bk-dropdown-menu1"
        @mouseover="show"
        @mouseout="hide">
        <div class="bk-dropdown-trigger" >
            <slot name="dropdown-trigger"></slot>
        </div>

        <div
            :class="[
                'bk-dropdown-content',
                {
                    'is-show': isShow,
                    'right-align': align == 'right',
                    'center-align': align == 'center'
                }
            ]"
	    :style="menuStyle">
            <slot name="dropdown-content"></slot>
        </div>
    </div>
</template>

<script>
    export default {
        name: 'bk-dropdown-menu',
        props: {
            align: {
                type: String,
                default: 'left'
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
