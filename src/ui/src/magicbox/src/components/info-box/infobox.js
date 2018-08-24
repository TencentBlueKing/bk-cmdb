import Vue from 'vue'
import {
    isVNode
} from '../../util.js'

import InfoBox from './infobox.vue'

const InfoBoxConstructor = Vue.extend(InfoBox)
let instance
let instancesArr = []
let count = 0
let zIndex = new Date().getFullYear()

let Info = function (options = {}) {
    let id = 'bkInfoBox' + count++

    if (typeof options === 'string') {
        options = {
            title: options
        }
    }

    instance = new InfoBoxConstructor({
        data: options
    })

    // 解析并挂载内容区的VNode
    if (isVNode(instance.content)) {
        instance.$slots.content = [instance.content]
        instance.content = null
    } else {
        delete instance.$slots.content
    }

    // 解析并挂载不同状态下的title的VNode
    if (isVNode(instance.statusOpts.title)) {
        instance.$slots.statusTitle = [instance.statusOpts.statusTitle]
        instance.statusOpts.statusTitle = null
    } else {
        delete instance.$slots.statusTitle
    }

    // 解析并挂载不同状态下的subtitle的VNode
    if (isVNode(instance.statusOpts.subtitle)) {
        instance.$slots.statusSubtitle = [instance.statusOpts.subtitle]
        instance.statusOpts.subtitle = null
    } else {
        delete instance.$slots.statusSubtitle
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

Info.hide = () => {
    let id = instance.id

    const len = instancesArr.length
    for (let index = 0; index < len; index++) {
        if (id === instancesArr[index].id) {
            instance.viewmodel.hide = true
        }
        instancesArr.splice(index, 1)
        break
    }
    // for (let [index, _instance] of instancesArr.entries()) {
    //     if (id === _instance.id) {
    //         instance.viewmodel.hide = true
    //     }

    //     instancesArr.splice(index, 1)
    //     break
    // }
}
Vue.prototype.$bkInfo = Info
export default Info
