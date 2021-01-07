import Vue from 'vue'
import i18n from '@/i18n'
import magicbox from './src'
import magicboxLocal from './src/locale'
import magicboxEn from './src/locale/lang/en-US'
import magicboxCn from './src/locale/lang/zh-CN'

i18n.setLocaleMessage('zh_CN', Object.assign({}, magicboxCn, i18n.messages['zh_CN']))
i18n.setLocaleMessage('en', Object.assign({}, magicboxEn, i18n.messages.en))

magicboxLocal.i18n((key, value) => i18n.t(key, value))

Vue.use(magicbox)

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
