import Vue from 'vue'
import magicBox from './src/index'

const languageMaps = {
    'zh_CN': 'zhCN',
    'en': 'enUS'
}

Vue.use(magicBox)

const Message = Vue.prototype.$bkMessage

let messageInstance = null

export const $error = (message, delay = 3000) => {
    messageInstance && messageInstance.close()
    messageInstance = Message({
        message,
        delay,
        theme: 'error'
    })
}

export const $success = (message, delay = 3000) => {
    messageInstance && messageInstance.close()
    messageInstance = Message({
        message,
        delay,
        theme: 'success'
    })
}

export const $info = (message, delay = 3000) => {
    messageInstance && messageInstance.close()
    messageInstance = Message({
        message,
        delay,
        theme: 'primary'
    })
}

export const $warn = (message, delay = 3000) => {
    messageInstance && messageInstance.close()
    messageInstance = Message({
        message,
        delay,
        theme: 'warning',
        hasCloseIcon: true
    })
}

Vue.prototype.$error = $error
Vue.prototype.$success = $success
Vue.prototype.$info = $info
Vue.prototype.$warn = $warn
