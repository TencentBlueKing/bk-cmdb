import Vue from 'vue'
import ViewModel from './loading.vue'

let Model = Vue.extend(ViewModel)

function toggle (el, binding) {
    if (!el.$vm) {
        el.$vm = el.viewmodel.$mount()
        el.appendChild(el.$vm.$el)
    }

    if (binding.value.isLoading) {
        Vue.nextTick(() => {
            el.$vm.isShow = true
        })
    } else {
        el.$vm.isShow = false
    }

    let title = binding.value.title

    if (title) {
        el.$vm.title = title
    }
}

let install = Vue => {
    Vue.directive('bkloading', {
        inserted (el, binding) {
            let value = binding.value
            let position = getComputedStyle(el).position
            let options = {}

            if (!position || position !== 'relative' || position !== 'absolute') {
                el.style.position = 'relative'
            }

            for (let key in value) {
                if (key !== 'isLoading') {
                    options[key] = value[key]
                }
            }

            options.type = 'directive'

            let loading = new Model({
                data: options
            })

            el.viewmodel = loading
            toggle(el, binding)
        },
        update (el, binding) {
            toggle(el, binding)
        }
    })
}

export default install
