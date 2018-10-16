import Vue from 'vue'
import {isVNode} from '../../util'
import LoadingView from './loading.vue'
const LoadingConstructor = Vue.extend(LoadingView)
let instance

let Loading = function (options = {}) {
    if (typeof options === 'string') {
        options = {
            title: options
        }
    }

    options.opacity = options.opacity || 0.9

    instance = new LoadingConstructor({
        data: options
    })

    if (isVNode(instance.title)) {
        instance.$slots.default = [instance.title]
        instance.title = null
    } else {
        delete instance.$slots.default
    }

    instance.viewmodel = instance.$mount()
    document.body.appendChild(instance.viewmodel.$el)
    instance.$dom = instance.viewmodel.$el
    instance.viewmodel.isShow = true

    return instance.viewmodel
}

Loading.hide = function () {
    instance.viewmodel.hide = true
}

Vue.prototype.$bkLoading = Loading

export default Loading
