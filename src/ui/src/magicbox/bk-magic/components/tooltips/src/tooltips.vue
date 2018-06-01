<template>
    <transition name="fade">
        <div class="bk-tooltips" :class="[direction]" :style="[calcPosition, {'max-width': maxWidth + 'px'}]" v-show="visible" @mouseenter.stop="show" @mouseleave.stop="hide">
            <div :class="['bk-tooltips-wrapper', {'bk-tooltips-default': typeof content != 'string'}]">
                <slot>
                    <template v-if="typeof content === 'string'">
                        <span>{{content}}</span>
                    </template>
                    <template v-else>
                        <a v-for="(item,index) in content" @click="clickHandler(item)">{{ item.text }}</a>
                    </template>
                </slot>
            </div>
            <div class="bk-tooltips-arrow"></div>
        </div>
    </transition>
</template>
<script>
/**
 *  bk-tooltips
 *  @module components/tooltips
 *  @desc 内容提示组件
 *  @param isShow {Boolean} - 组件是否显示，仅在指令中使用
 *  @param maxWidth {Number} - 组件最大宽度
 *  @param content {String, VNode} - 提示信息内容
 *  @param direction {String} - 组件出现时相对于触发元素的方向
 *  @param onShow {Function} - 组件显示时的回调函数
 *  @param onClose {Function} - 组件关闭时的回调函数
 *  @example
    <bk-button type="primary" class="mr30"
        v-bktooltips="{
            maxWidth: 200,
            isShow,
            content: '提示信息',
            direction: 'top',
            onShow,
            onClose
        }"
    >内容提示组件</bk-button>
 */

import {
    getActualTop,
    getActualLeft
} from './../../../utils/utils'

export default {
    name: 'bk-tooltips',
    data () {
        return {
            maxWidth: '',
            content: '',
            direction: 'right',
            visible: false, // 组件内部变量
            ele: null, // 组件内部变量
            calcPosition: {}, // 组件内部变量
            timer: 0, // 组件内部变量
            onShow () {},
            onClose () {}
        }
    },
    watch: {
        visible (newVal) {
            if (newVal) {
                this.calcPos()
                this.onShow && this.onShow(this)
            } else {
                this.onClose && this.onClose(this)
            }
        }
    },
    methods: {
        clickHandler (item) {
            item.callback && item.callback()
            this.visible = false
        },
        strToNum (str) {
            return Number(str.replace('px', ''))
        },
        calcPos () {
            let ele = this.ele
            let strToNum = this.strToNum
            let actualTop = getActualTop(ele)
            let actualLeft = getActualLeft(ele)
            let eleComputedStyle = getComputedStyle(ele)
            let eleWidth = strToNum(eleComputedStyle.width)
            let eleHeight = strToNum(eleComputedStyle.height)

            setTimeout(() => {
                let curComputedStyle = getComputedStyle(this.$el)
                let $elWidth = strToNum(curComputedStyle.width)
                let $elHeight = strToNum(curComputedStyle.height)
                let pos = {}
                let top, left
                let gap = 12

                switch (this.direction) {
                    case 'top':
                        top = `${actualTop - $elHeight - gap}`
                        left = `${actualLeft + (eleWidth - $elWidth) / 2}`
                        break
                    case 'right':
                        top = `${actualTop + (eleHeight - $elHeight) / 2}`
                        left = `${actualLeft + eleWidth + gap}`
                        break
                    case 'bottom':
                        top = `${actualTop + eleHeight + gap}`
                        left = `${actualLeft + (eleWidth - $elWidth) / 2}`
                        break
                    case 'left':
                        top = `${actualTop + (eleHeight - $elHeight) / 2}`
                        left = `${actualLeft - $elWidth - gap}`
                        break
                }
                this.$set(this.calcPosition, 'top', `${top}px`)
                this.$set(this.calcPosition, 'left', `${left}px`)
            }, 0)
        },
        show () {
            clearTimeout(this.timer)
            this.timer = -1
            this.visible = true
        },
        hide () {
            this.timer = setTimeout(() => {
                this.visible = false
            }, 200)
        }
    }
}
</script>
