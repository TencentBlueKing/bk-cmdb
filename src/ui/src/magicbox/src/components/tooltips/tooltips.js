import Vue from 'vue'
import ViewModel from './tooltips.vue'

const TooltipsConstructor = Vue.extend(ViewModel)

const props = ViewModel.props
const defaultOptions = {}
Object.keys(props).forEach(key => {
    const prop = props[key]
    const dv = prop.default
    if (prop && prop.default != null) {
        defaultOptions[key] = typeof dv === 'function' ? dv() : dv
    }
})

let tooltipsInstance = null

export default function tooltips (options) {
    options = options || {}
    // 已存在 tooltips 实例，直接更新属性值
    if (tooltipsInstance && tooltipsInstance.$el.parentNode) {
        Object.assign(tooltipsInstance, defaultOptions, options)
        if (tooltipsInstance.target) {
            tooltipsInstance.updateTooltips()
        } else {
            tooltipsInstance.hiddenTooltips()
        }
        return tooltipsInstance
    }

    // 否则创建一个 tooltips 实例
    tooltipsInstance = new TooltipsConstructor({
        propsData: options
    }).$mount()

    tooltipsInstance.updateTooltips()
    return tooltipsInstance
}
