/**
 * @file 语言逻辑
 * @author ielgnaw <wuji0223@gmail.com>
 */

import Vue from 'vue'
import defaultLang from './lang/zh-CN'
import deepmerge from '../utils/deepmerge'

let curLang = defaultLang

// Vue.locale = () => {}

let merged = false
/**
 * 检测是否使用 vue-i18n，如果使用了，那么会用 vue-i18n 的 $t 来取值
 */
let i18nHandler = function () {
    const i18n = Object.getPrototypeOf(this || Vue).$t
    // vuei18n 6.x 没有 locale 方法
    // if (typeof i18n === 'function' && !!Vue.locale) {
    //     if (!merged) {
    //         merged = true
    //         Vue.locale(
    //             Vue.config.lang,
    //             deepmerge(curLang, Vue.locale(Vue.config.lang) || {}, {clone: true})
    //         )
    //     }
    //     return i18n.apply(this, arguments)
    // }
    if (typeof i18n === 'function') {
        if (!merged && !!Vue.locale) {
            merged = true
            Vue.locale(
                Vue.config.lang,
                deepmerge(curLang, Vue.locale(Vue.config.lang) || {}, {clone: true})
            )
        }
        return i18n.apply(this, arguments)
    }
}

/**
 * 转义特殊字符
 *
 * @param {string} str 待转义字符串
 *
 * @return {string} 结果
 */
export const escape = str => String(str).replace(/([.*+?^=!:${}()|[\]\/\\])/g, '\\$1')

/**
 * 根据语言环境获取对应的值
 *
 * @param {string} path 词语的路径，对应语言包文件里的 key 的路径
 * @param {Object} data 要替换的值
 *
 * @return {string} 对应语言包的值
 */
export const t = function (path, data) {
    let value = i18nHandler.apply(this, arguments)
    if (value !== null && typeof value !== 'undefined') {
        return value
    }

    const arr = path.split('.')
    let current = curLang
    const len = arr.length

    for (let i = 0; i < len; i++) {
        value = current[arr[i]]
        if (i === len - 1) {
            if (data && typeof value === 'string') {
                return value.replace(/\{(?=\w+)/g, '').replace(/(\w+)\}/g, '$1')
                    .replace(new RegExp(Object.keys(data).map(escape).join('|'), 'g'), $0 => data[$0])
            }
            return value
        }
        if (!value) {
            return ''
        }
        current = value
    }
    return ''
}

/**
 * 使用某种语言
 *
 * @param {Object} l 使用的语言包
 */
export const use = l => {
    if (l) {
        curLang = deepmerge(curLang, l)
    }
}

/**
 * 自定义 i18n 的处理函数
 *
 * @param {Function} fn i18n 处理函数
 */
export const i18n = fn => {
    i18nHandler = fn || i18nHandler
}

export default {use, t, i18n}
