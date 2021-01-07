import Vue from 'vue'
import VueI18n from 'vue-i18n'
import Cookies from 'js-cookie'
import messages from './lang/messages'

Vue.use(VueI18n)

const languageMaps = {
    'zh_cn': 'zh_CN',
    'zh-cn': 'zh_CN',
    'zh': 'zh_CN'
}

let locale = Cookies.get('blueking_language') || 'zh_CN'

locale = languageMaps.hasOwnProperty(locale) ? languageMaps[locale] : locale

const i18n = new VueI18n({
    locale,
    fallbackLocale: 'zh_CN',
    messages,
    missing (locale, path) {
        const parsedPath = i18n._path.parsePath(path)
        return parsedPath[parsedPath.length - 1]
    }
})

export const language = locale

export default i18n
