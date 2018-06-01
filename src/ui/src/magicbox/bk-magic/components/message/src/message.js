import Vue from 'vue'
import {
    isVNode
} from './../../../utils/utils'
import Message from './message.vue'

const MessageConstructor = Vue.extend(Message)
let instance // 当前组件实例
let instancesArr = []
let count = 0
let zIndex = new Date().getFullYear()

let Msg = function (options = {}) {
    let id = 'bkMessage' + count++
    let usrClose = options.onClose
    let type = (typeof options).toLowerCase()

    if (type === 'string' || type === 'number') {
        options = {
            message: options
        }
    }

    options.onClose = function () {
        Msg.close(id, usrClose)
    }

    instance = new MessageConstructor({
        data: options
    })

    // 解析并挂载内容区的VNode
    if (isVNode(instance.message)) {
        instance.$slots.default = [instance.message]
        instance.message = null
    }

    instance.id = id
    instance.viewmodel = instance.$mount()
    document.body.appendChild(instance.viewmodel.$el)
    instance.$dom = instance.viewmodel.$el
    instance.$dom.style.zIndex = zIndex++
    instance.viewmodel.isShow = true
    instancesArr.push(instance)
    return instance.viewmodel
}

Msg.close = function (id, usrClose) {
    for (let [index, _instance] of instancesArr.entries()) {
        if (id === _instance.id) {
            usrClose && usrClose(_instance)
        }

        instancesArr.splice(index, 1)
        break
    }
}

export default Msg
