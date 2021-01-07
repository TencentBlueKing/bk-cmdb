/**
 * @file tooltips directive
 * @author ielgnaw <wuji0223@gmail.com>
 */

import Tooltips from './tooltips.js'

/**
 * 清除事件
 *
 * @param {Object} el dom 对象
 */
function clearEvent (el) {
    if (el._tooltipsHandler) {
        el.removeEventListener('click', el._tooltipsHandler)
        el.removeEventListener('mouseenter', el._tooltipsHandler)
    }
    if (el._tooltipsMouseleaveHandler) {
        el.removeEventListener('mouseleave', el._tooltipsMouseleaveHandler)
    }
    delete el._tooltipsHandler
    delete el._tooltipsMouseleaveHandler
    delete el._tooltipsOptions
    delete el._tooltipsInstance
}

export default {
    install (Vue, options) {
        options = options || {}
        // 展示方向
        const allPlacements = ['top', 'right', 'bottom', 'left']

        Vue.directive('bktooltips', {
            bind (el, binding) {
                clearEvent(el)

                const {click, light} = binding.modifiers
                const limitPlacementQueue = allPlacements.filter(placement => binding.modifiers[placement])

                el._tooltipsOptions = binding.value

                el._tooltipsHandler = function tooltipsHandler () {
                    if (this._tooltipsOptions == null) {
                        return
                    }
                    const options = this._tooltipsOptions
                    const placements = limitPlacementQueue.length ? limitPlacementQueue : allPlacements
                    const mix = {
                        placements,
                        theme: light ? 'light' : 'dark'
                    }

                    // v-bktooltips directive 使用时可以直接配置需要显示的内容，也可以直接绑定一个配置对象来自定义配置
                    const tipOptions = typeof options === 'object'
                        ? Object.assign(mix, options, {target: this})
                        : Object.assign(mix, {content: String(options), target: this})

                    this._tooltipsInstance = Tooltips(tipOptions)
                }
                el._tooltipsMouseleaveHandler = function tooltipsMouseleaveHandler () {
                    if (this._tooltipsInstance) {
                        this._tooltipsInstance.hiddenTooltips()
                    }
                }
                // 默认触发方式为 hover 触发
                if (click) {
                    el.addEventListener('click', el._tooltipsHandler)
                } else {
                    el.addEventListener('mouseenter', el._tooltipsHandler)
                }
                el.addEventListener('mouseleave', el._tooltipsMouseleaveHandler)
            },

            update (el, {value, oldValue}) {
                if (value === oldValue) {
                    return
                }
                el._tooltipsOptions = value
            },

            unbind (el) {
                const instance = el._tooltipsInstance
                if (instance && instance.destroy) {
                    instance.destroy()
                }
                clearEvent(el)
            }
        })
    }
}
