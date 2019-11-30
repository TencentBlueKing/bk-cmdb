import Vue from 'vue'
import { language } from '@/i18n'
import veeValidate, { Validator } from 'vee-validate'
import cnMessages from 'vee-validate/dist/locale/zh_CN'
import stringLength from 'utf8-byte-length'

const customRules = {
    singlechar: {
        validate: value => {
            return /^([a-zA-Z0-9]|[\u4e00-\u9fa5]|[\(\)\+\-《》_,，；:;“”‘’。@#\."'\\\/\s？!！～、：＃％%＊*—…＆&·＄$\^（）\[\]『』〔〕｛｝【】￥￡♀‖〖〗「」]){0,256}$/.test(value)
        }
    },
    length: {
        validate: (value, [length]) => {
            return stringLength(value) <= length
        }
    },
    longchar: {
        validate: value => {
            return /^([a-zA-Z0-9]|[\u4e00-\u9fa5]|[\(\)\+\-《》_,，；:;“”‘’。@#\."'\\\/\s？!！～、：＃％%＊*—…＆&·＄$\^（）\[\]『』〔〕｛｝【】￥￡♀‖〖〗「」]){0,2000}$/.test(value)
        }
    },
    associationId: {
        validate: (value) => {
            return /^[a-z_]+$/.test(value)
        }
    },
    classifyName: {
        validate: value => {
            return /^([a-zA-Z0-9_ ]|[\u4e00-\u9fa5]|[\uac00-\ud7ff]|[\u0800-\u4e00]){1,20}$/.test(value)
        }
    },
    classifyId: {
        validate: value => {
            return /^[a-zA-Z0-9_]{1,20}$/.test(value)
        }
    },
    http: {
        validate: value => {
            return /^http(s?):\/\/[^\s]+/.test(value)
        }
    },
    modelId: {
        validate: value => {
            return /^[a-z][a-z\d_]*$/.test(value)
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
            return /^[a-z0-9_]{1,20}$/.test(value)
        }
    },
    float: {
        validate: (value) => {
            return /^[+-]?([0-9]*[.]?[0-9]+|[0-9]+[.]?[0-9]*)([eE][+-]?[0-9]+)?$/.test(value)
        }
    },
    oid: {
        validate: (value) => {
            return /^(\d+)?(\.\d+)+$/.test(value)
        }
    },
    hourFormat: {
        validate: (value) => {
            return /^[1-5]?[0-9]$/.test(value)
        }
    },
    dayFormat: {
        validate: (value) => {
            return /^((20|21|22|23|[0-1]\d):[0-5][0-9])?$/.test(value)
        }
    },
    namedCharacter: {
        validate: (value) => {
            return /^([a-zA-Z0-9]|[\u4e00-\u9fa5]|[-_:]){0,256}$/.test(value)
        }
    },
    instanceTagKey: {
        validate: value => {
            return /^[a-zA-Z]([a-z0-9A-Z\-_.]*[a-z0-9A-Z])?$/.test(value)
        }
    },
    instanceTagValue: {
        validate: value => {
            return /^[a-z0-9A-Z]([a-z0-9A-Z\-_.]*[a-z0-9A-Z])?$/.test(value)
        }
    },
    repeatTagKey: {
        validate: (value, otherValue) => {
            return otherValue.findIndex(item => item === value) === -1
        }
    },
    setNameMap: {
        validate: (value) => {
            const nameList = value.split('\n').filter(name => name)
            const nameSet = new Set(nameList)
            return nameList.length === nameSet.size
        }
    },
    emptySetName: {
        validate: (value) => {
            const values = value.split('\n')
            const list = values.map(text => text.trim()).filter(text => text)
            return values.length === list.length
        }
    }
}

const dictionary = {
    'zh_CN': {
        messages: {
            regex: (field) => {
                return `请输入合法的${field}`
            },
            longchar: () => '请输入正确的长字符内容',
            singlechar: () => '请输入正确的短字符内容',
            length: (field, [maxLength]) => {
                return `请输入${maxLength}个字符以内的内容`
            },
            associationId: () => '格式不正确，只能包含下划线，英文小写',
            classifyName: () => '请输入正确的内容',
            classifyId: () => '请输入正确的内容',
            required: () => '该字段是必填项',
            http: () => '请输入以http(s)://开头的URL',
            modelId: () => '格式不正确，请填写英文开头，下划线，数字，英文小写的组合',
            enumId: () => '请输入正确的内容',
            enumName: () => '请输入正确的内容',
            number: () => '请输入正确的数字',
            float: () => '请输入正确的浮点数',
            isBigger: () => '必须大于最小值',
            repeat: () => '重复的值',
            fieldId: () => '请输入正确的内容',
            oid: () => '请输入正确的内容',
            hourFormat: () => '请输入0-59之间的数字',
            dayFormat: () => '请输入00:00-23:59之间的时间',
            namedCharacter: () => '格式不正确，特殊符号仅支持(:_-)',
            min_value: () => '该值小于最小值',
            max_value: () => '该值大于最大值',
            repeatTagKey: () => '标签键不能重复',
            setNameMap: () => '集群名称重复',
            emptySetName: () => '请勿输入空白集群名称',
            instanceTagValue: () => '请输入英文 / 数字',
            instanceTagKey: () => '请输入英文 / 数字, 以英文开头'
        },
        custom: {
            asst: {
                required: '请选择关联模型'
            }
        }
    },
    en: {
        messages: {
            regex: (field) => {
                return `Please enter a valid $ {field}`
            },
            longchar: () => 'Please enter the correct content',
            singlechar: () => 'Please enter the correct content',
            length: (field, [maxLength]) => {
                return `Content length max than ${maxLength}`
            },
            associationId: () => 'The format is incorrect and can only contain underscores and lowercase letter',
            classifyName: () => 'Please enter the correct content',
            classifyId: () => 'Please enter the correct content',
            required: () => 'This field is required',
            http: () => 'Please enter a URL beginning with http(s)://',
            modelId: () => 'The format is incorrect, can only contain underscores, numbers, and lowercase letter and start with a letter',
            enumId: () => 'Please enter the correct content',
            enumName: () => 'Please enter the correct content',
            number: () => 'Please enter the correct number',
            float: () => 'Please enter the correct float data',
            isBigger: () => 'Must be greater than the minimum',
            repeat: () => 'This value should not be repeated',
            fieldId: () => 'Please enter the correct content',
            oid: () => 'Please enter the correct content',
            hourFormat: () => 'Please enter the number between 0-59',
            dayFormat: () => 'Please enter the time between 00:00-23:59',
            namedCharacter: () => 'Special symbols only support(:_-)',
            min_value: () => 'This value is less than the minimum',
            max_value: () => 'This value is greater than the maximum',
            setNameMap: () => 'Duplicate Set name',
            emptySetName: () => 'Do not enter blank Set name',
            instanceTagValue: () => 'Please enter letter / number',
            instanceTagKey: () => 'Please enter letter / number starts with letter',
            repeatTagKey: () => 'Label key cannot be repeated'
        },
        custom: {
            asst: {
                required: 'Please select the associated model'
            }
        }
    }
}

for (const rule in customRules) {
    Validator.extend(rule, customRules[rule])
}
if (language === 'en') {
    Validator.localize(language)
} else {
    Validator.localize(language, cnMessages)
}
Vue.use(veeValidate, {
    locale: language,
    dictionary
})
