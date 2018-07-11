<template>
    <div class="bk-dropdown-menu" id="bk-dropdown-menu1" v-clickoutside="handleClickoutside"
        @click="handleClick"
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
    import clickoutside from './../../../utils/clickoutside'
    export default {
        name: 'bk-dropdown-menu',
        props: {
            align: {
                type: String,
                default: 'left'
            },
            trigger: {
                type: String,
                default: 'mouseover',
                validator (trigger) {
                    return ['click', 'mouseover'].includes(trigger)
                }
            }
        },
        directives: {
            clickoutside
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
                if (this.trigger !== 'mouseover') return
                clearTimeout(this.timer)
                if (this.isShow) return
                this.calcPosition()
                this.isShow = true
                this.$emit('show')
            },
            hide () {
                if (this.trigger !== 'mouseover') return
                this.timer = setTimeout(() => {
                    this.isShow = false
                    this.$emit('hide')
                }, 200)
            },
            calcPosition () {
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
            },
            handleClick () {
                if (this.trigger !== 'click') return
                if (this.isShow) {
                    this.isShow = false
                    this.$emit('hide')
                } else {
                    this.calcPosition()
                    this.isShow = true
                    this.$emit('show')
                }
            },
            handleClickoutside () {
                if (this.isShow) {
                    this.isShow = false
                    this.$emit('hide')
                }
            }
        }
    }
</script>
