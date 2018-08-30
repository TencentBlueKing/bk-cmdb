import en from './en.json'
import cn from './cn.json'
import magicboxEn from '@/magicbox/src/locale/lang/en-US'
import magicboxCn from '@/magicbox/src/locale/lang/zh-CN'
export default {
    en: Object.assign(magicboxEn, en),
    'zh_CN': Object.assign(magicboxCn, cn)
}
