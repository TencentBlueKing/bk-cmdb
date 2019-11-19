import Vue from 'vue'
import i18n from '@/i18n'
import magicbox from 'bk-magic-vue'
import './magicbox.scss'

const magicboxLanguageMap = {
    zh_CN: magicbox.locale.lang.zhCN,
    en: magicbox.locale.lang.enUS
}

i18n.mergeLocaleMessage(i18n.locale, magicboxLanguageMap[i18n.locale])
magicbox.locale.use(magicboxLanguageMap[i18n.locale])
Vue.use(magicbox, {
    'bk-sideslider': {
        quickClose: true,
        width: 800
    },
    'bk-input': {
        fontSize: 'medium'
    },
    i18n: (key, value) => i18n.t(key, value)
})

export const $error = (message, delay = 3000) => {
    magicbox.bkMessage({
        message,
        delay,
        theme: 'error'
    })
}

export const $success = (message, delay = 3000) => {
    magicbox.bkMessage({
        message,
        delay,
        theme: 'success'
    })
}

export const $info = (message, delay = 3000) => {
    magicbox.bkMessage({
        message,
        delay,
        theme: 'primary'
    })
}

export const $warn = (message, delay = 3000) => {
    magicbox.bkMessage({
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
Vue.prototype.$bkInfo = magicbox.bkInfoBox
