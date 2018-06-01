import Vue from 'vue'
import ViewModel from './toolbox.vue'
import { isVNode } from './../../../utils/utils'

let Model = Vue.extend(ViewModel)

function toggle (el, isShow) {
    if (isShow) {
        if (!el.$vm) {
            el.$vm = el.viewmodel.$mount()
            Vue.nextTick(() => {
                document.body.appendChild(el.$vm.$el)
            })
        }

        Vue.nextTick(() => {
            el.$vm.visible = true
        })
    } else {
        if (el.$vm) {
            el.$vm.visible = false
        }
    }
}

let install = Vue => {
    Vue.directive('bktooltips', {
        inserted (el, binding) {
            let value = binding.value
            let position = getComputedStyle(el).position
            let options = {}

            if (!position || position !== 'relative' || position !== 'absolute') {
                el.style.position = 'relative'
            }

            for (let key in value) {
                options[key] = value[key]
            }

            options.ele = el
            let tooltips = new Model({
                data: options
            })

            if (isVNode(tooltips.content)) {
                tooltips.$slots.default = [tooltips.content]
                tooltips.content = ''
            } else {
                delete tooltips.$slots.default
            }

            el.viewmodel = tooltips

            // 绑定鼠标事件
            el.addEventListener('mouseenter', (event) => {
                let $vm = el.$vm

                if (!$vm) {
                    toggle(el, true)
                } else {
                    el.$vm.isShow = true
                    toggle(el, el.$vm.isShow)
                }

                event.stopPropagation()
            }, false)

            el.addEventListener('mouseleave', (event) => {
                var tm = el.$vm.$el
                event.stopPropagation()
                el.$vm.isShow = false
                toggle(el, el.$vm.isShow)
                tm.addEventListener('mouseenter', (event) => {
                    el.$vm.isShow = true
                    toggle(el, el.$vm.isShow)
                })
                tm.addEventListener('mouseleave', (event) => {
                    el.$vm.isShow = false
                    toggle(el, el.$vm.isShow)
                })
            }, false)

            toggle(el, binding.value.isShow)
        },
        update (el, binding) {
            if (binding.value.isShow !== binding.oldValue.isShow) {
                toggle(el, binding)
            }
        }
    })
}

export default install
