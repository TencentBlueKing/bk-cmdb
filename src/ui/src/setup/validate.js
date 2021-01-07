import Vue from 'vue'
import { language } from '@/i18n'
import veeValidate, { Validator } from 'vee-validate'
import cnMessages from 'vee-validate/dist/locale/zh_CN'
import stringLength from 'utf8-byte-length'

const customRules = {
    length: {
        validate: (value, [length]) => {
            return stringLength(value) <= length
        }
    },
    repeat: {
        validate: (value, otherValue) => {
            return otherValue.findIndex(item => item === value) === -1
        }
    },
    singlechar: {
        validate: value => {
            return /\S*/.test(value)
        }
    },
    longchar: {
        validate: value => {
            return /\S*/.test(value)
        }
    },
    associationId: {
        validate: (value) => {
            return /^[a-zA-Z][\w]*$/.test(value)
        }
    },
    // 未被使用
    classifyName: {
        validate: value => {
            return /^([a-zA-Z0-9_ ]|[\u4e00-\u9fa5]|[\uac00-\ud7ff]|[\u0800-\u4e00]){1,20}$/.test(value)
        }
    },
    // 模型分组id
    classifyId: {
        validate: value => {
            return /^[a-zA-Z][\w]*$/.test(value)
        }
    },
    http: {
        validate: value => {
            return /^http(s?):\/\/[^\s]+/.test(value)
        }
    },
    // 新建模型唯一标识id
    modelId: {
        validate: value => {
            return /^[a-zA-Z][\w]*$/.test(value)
        }
    },
    enumId: {
        validate: value => {
            return /^[a-zA-Z0-9_]*$/.test(value)
        }
    },
    enumName: {
        validate: (value) => {
            return /^([a-zA-Z0-9_]|[\u4e00-\u9fa5]|[()+-《》,，；;“”‘’。\."\' \\/:])*$/.test(value)
        }
    },
    number: {
        validate: (value) => {
            if (!String(value).length) {
                return true
            }
            return /^(\-|\+)?\d+$/.test(value)
        }
    },
    isBigger: {
        validate: (value, [targetValue]) => {
            return Number(value) > Number(targetValue)
        }
    },
    // 新建字段唯一标识
    fieldId: {
        validate: (value) => {
            return /^[a-zA-Z][\w]*$/.test(value)
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
    // 服务分类名称
    namedCharacter: {
        validate: (value) => {
            return /^[a-zA-Z0-9\u4e00-\u9fa5_\-:\(\)]+$/.test(value)
        }
    },
    // 服务实例标签键
    instanceTagKey: {
        validate: value => {
            return /^[a-zA-Z]([a-z0-9A-Z\-_.]*[a-z0-9A-Z])?$/.test(value)
        }
    },
    // 服务实例标签值
    instanceTagValue: {
        validate: value => {
            return /^[a-z0-9A-Z]([a-z0-9A-Z\-_.]*[a-z0-9A-Z])?$/.test(value)
        }
    },
    businessTopoInstNames: {
        validate: value => {
            return /^[^\#\/,\>\<\|]+$/.test(value)
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
    },
    setNameLen: {
        validate: (value) => {
            const nameList = value.split('\n').filter(name => name)
            for (const name of nameList) {
                if (stringLength(name) > 256) return false
            }
            return true
        }
    },
    reservedWord: {
        validate: value => {
            return /^(?!bk_).*/.test(value)
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
            associationId: () => '格式不正确，请填写英文开头，下划线，数字，英文的组合',
            classifyName: () => '请输入正确的内容',
            classifyId: () => '请输入正确的内容',
            required: () => '该字段是必填项',
            http: () => '请输入以http(s)://开头的URL',
            modelId: () => '格式不正确，请填写英文开头，下划线，数字，英文的组合',
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
            instanceTagValue: () => '请输入英文数字的组合',
            instanceTagKey: () => '请输入英文开头数字的组合',
            setNameLen: () => '请输入256个字符以内的内容',
            businessTopoInstNames: () => '格式不正确，不能包含特殊字符 #/,><|',
            reservedWord: () => '不能以"bk_"开头'
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
            associationId: () => 'The format is incorrect, can only contain underscores, numbers, letter and start with a letter',
            classifyName: () => 'Please enter the correct content',
            classifyId: () => 'Please enter the correct content',
            required: () => 'This field is required',
            http: () => 'Please enter a URL beginning with http(s)://',
            modelId: () => 'The format is incorrect, can only contain underscores, numbers, letter and start with a letter',
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
            instanceTagValue: () => 'Please enter letter, number',
            instanceTagKey: () => 'Please enter letter, number starts with letter',
            repeatTagKey: () => 'Label key cannot be repeated',
            setNameLen: () => 'Content length max than 256',
            businessTopoInstNames: () => 'The format is incorrect and cannot contain special characters #/,><|',
            reservedWord: () => 'Can not start with "bk_"'
        },
        custom: {
            asst: {
                required: 'Please select the associated model'
            }
        }
    }
}

// 可配置规则清单
const customConfigRules = [
    {
        number: (value, cb) => {
            if (!String(value).length) {
                return true
            }
            return cb()
        }
    },
    'float',
    'singlechar',
    'longchar',
    'associationId',
    'classifyId',
    'modelId',
    'enumId',
    'enumName',
    'fieldId',
    'namedCharacter',
    'instanceTagKey',
    'instanceTagValue',
    'businessTopoInstNames'
]

const mixinConfig = () => {
    const { validationRules = {} } = window.CMDB_CONFIG || {}
    for (const item of customConfigRules) {
        const useCb = typeof item !== 'string'
        const key = useCb ? Object.keys(item)[0] : item

        const rule = validationRules[key]
        if (!rule) continue

        let validate = value => {
            return new RegExp(rule.value).test(value)
        }
        if (useCb) {
            validate = value => {
                return item[key](value, () => {
                    return new RegExp(rule.value).test(value)
                })
            }
        }

        // 加入到自定义规则列表
        customRules[key] = { validate }
        // 提示语设置
        dictionary['zh_CN']['messages'][key] = (field, args) => {
            // 确保总是获取最新的配置
            const { validationRules } = window.CMDB_CONFIG
            const rule = validationRules[key]
            return rule.i18n.cn.replace(/{field}/g, field)
        }
        dictionary['en']['messages'][key] = (field, args) => {
            const { validationRules } = window.CMDB_CONFIG
            const rule = validationRules[key]
            return rule.i18n.en.replace(/{field}/g, field)
        }
    }
}

export function setupValidator (app) {
    mixinConfig()
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

    app.$store.commit('setValidatorSetuped')
}

export function updateValidator () {
    mixinConfig()
    for (const rule in customRules) {
        Validator.extend(rule, customRules[rule])
    }
}
