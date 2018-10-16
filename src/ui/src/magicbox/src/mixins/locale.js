/**
 * @file 语言处理的 mixin，给组件加上一个 t 方法，组件在需要根据语言切换的地方，只要加入这个 mixin 并在输出的地方使用 t(key) 即可，
 *       例如 t(datePicker.today)
 * @author hieiwang <wuji0223@gmail.com>
 */

import {t} from '../locale'

export default {
    methods: {
        t (...args) {
            return t.apply(this, args)
        }
    }
}
