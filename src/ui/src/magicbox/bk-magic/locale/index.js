import ZH from './lang/zh'
import EN from './lang/en'

let langCode = 'zh'
const langs = {
    zh: ZH,
    en: EN
}

function t (key) {
    let langKeys = Object.keys(langs)
    if (!langKeys.indexOf(langCode)) {
        langCode = 'zh'
    }
    let lang = langs[langCode]
    let result = lang

    let paths = key.split('.')
    for (let path of paths) {
        result = result[path]
    }

    if (result) {
        return result
    } else {
        return ''
    }
}

function use (l) {
    langCode = l
}
export default {
    t,
    use
}
