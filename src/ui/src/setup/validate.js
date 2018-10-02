import Vue from 'vue'
import { language } from '@/i18n'
import veeValidate, {Validator} from 'vee-validate'

const customRules = {
    singlechar: {
        validate: value => {
            /* eslint-disable */
            return /^([a-zA-Z0-9]|[\u4e00-\u9fa5]|[\(\)\+\-《》_,，；:;“”‘’。\."'\\\/\s])*$/.test(value)
            /* eslint-enable */
        }
    },
    longchar: {
        validate: value => {
            /* eslint-disable */
            return /^([a-zA-Z0-9]|[\u4e00-\u9fa5]|[\(\)\+\-《》_,，；:;“”‘’。\."'\\\/\s])*$/.test(value)
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
            return /^http(s?):\/\/[^\s]+/.test(value)
        }
    },
    modelId: {
        validate: value => {
            return /^[a-z\d_]+$/.test(value)
        }
    },
    enumId: {
        validate: value => {
            return /^[a-zA-Z0-9_]{1,20}$/.test(value)
        }
    },
    enumName: {
        validate: (value) => {
            /* eslint-disable */
            return /^([a-zA-Z0-9_]|[\u4e00-\u9fa5]|[()+-《》,，；;“”‘’。\."\' \\/:]){1,15}$/.test(value)
            /* eslint-enable */
        }
    },
    repeat: {
        validate: (value, otherValue) => {
            return otherValue.findIndex(item => item === value) === -1
        }
    },
    number: {
        validate: (value) => {
            return /^(-)?[0-9]*$/.test(value)
        }
    },
    isBigger: {
        validate: (value, [targetValue]) => {
            return Number(value) > Number(targetValue)
        }
    },
    fieldId: {
        validate: (value) => {
            return /^[a-zA-Z0-9_]{1,20}$/.test(value)
        }
    }
}

const dictionary = {
    'zh_CN': {
        messages: {
            regex: () => '请输入符合自定义正则的内容',
            longchar: () => '请输入正确的长字符内容',
            singlechar: () => '请输入正确的短字符内容',
            classifyName: () => '请输入正确的内容',
            classifyId: () => '请输入正确的内容',
            required: () => '该字段是必填项',
            http: () => '请输入以http(s)://开头的URL',
            modelId: () => '格式不正确，只能包含下划线，数字，英文小写',
            enumId: () => '请输入正确的内容',
            enumName: () => '请输入正确的内容',
            number: () => '请输入正确的内容',
            isBigger: () => '必须大于最小值',
            repeat: () => '重复的值',
            fieldId: () => '请输入正确的内容'
        },
        custom: {
            asst: {
                required: '请选择关联模型'
            }
        }
    },
    en: {
        messages: {
            regex: () => 'Please enter the correct content that conform custom regex',
            longchar: () => 'Please enter the correct content',
            singlechar: () => 'Please enter the correct content',
            classifyName: () => 'Please enter the correct content',
            classifyId: () => 'Please enter the correct content',
            required: () => 'This field is required',
            http: () => 'Please enter a URL beginning with http(s)://',
            modelId: () => 'The format is incorrect and can only contain underscores, numbers, and lowercase English',
            enumId: () => 'Please enter the correct content',
            enumName: () => 'Please enter the correct content',
            number: () => 'Please enter the correct content',
            isBigger: () => 'Must be greater than the minimum',
            repeat: () => 'This value should not be repeated',
            fieldId: () => 'Please enter the correct content'
        },
        custom: {
            asst: {
                required: 'Please select the associated model'
            }
        }
    }
}

for (let rule in customRules) {
    Validator.extend(rule, customRules[rule])
}
Validator.localize(language)
Vue.use(veeValidate, {
    locale: language,
    dictionary
})
