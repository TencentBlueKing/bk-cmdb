import Vue from 'vue'
import { language } from '@/i18n'
import veeValidate, {Validator} from 'vee-validate'

const customRules = {
    singlechar: {
        validate: value => {
            /* eslint-disable */
            return /^([a-zA-Z0-9]|[\u4e00-\u9fa5]|[\(\)\+\-《》_,，；;“”‘’。\."'\\\/])*$/.test(value)
            /* eslint-enable */
        }
    },
    longchar: {
        validate: value => {
            /* eslint-disable */
            return /^([a-zA-Z0-9]|[\u4e00-\u9fa5]|[\(\)\+\-《》_,，；;“”‘’。\."'\\\/])*$/.test(value)
            /* eslint-enable */
        }
    },
    classifyName: {
        validate: value => {
            return /^([a-zA-Z0-9_ ]|[\u4e00-\u9fa5]|[\uac00-\ud7ff]|[\u0800-\u4e00]){1,20}$/.test(value)
        }
    },
    classifyId: {
        validate: value => {
            return /^[a-z0-9_]{1,20}$/.test(value)
        }
    },
    http: {
        validate: value => {
            return /^http:\/\/[^\s]+/.test(value)
        }
    },
    modelId: {
        validate: value => {
            return /^[a-z\d_]+$/.test(value)
        }
    }
}

const dictionary = {
    'zh_CN': {
        messages: {
            longchar: () => '请输入正确的内容',
            singlechar: () => '请输入正确的内容',
            classifyName: () => '请输入正确的内容',
            classifyId: () => '请输入正确的内容',
            required: (field) => '请输入' + field,
            http: () => '请输入以http://开头的URL',
            modelId: () => '格式不正确，只能包含下划线，数字，英文小写'
        }
    },
    en: {
        messages: {
            longchar: () => 'Please enter the correct content',
            singlechar: () => 'Please enter the correct content',
            classifyName: () => 'Please enter the correct content',
            classifyId: () => 'Please enter the correct content',
            required: (field) => 'Please enter ' + field,
            http: () => 'Please enter a URL beginning with http://',
            modelId: () => '格式不正确，只能包含下划线，数字，英文小写'
        }
    }
}

for (let rule in customRules) {
    Validator.extend(rule, customRules[rule])
}

Validator.localize(dictionary)
Vue.use(veeValidate, {locale: language})
