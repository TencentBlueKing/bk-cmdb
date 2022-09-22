/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

import Vue from 'vue'
import { language } from '@/i18n'
import veeValidate, { Validator } from 'vee-validate'
import cnMessages from 'vee-validate/dist/locale/zh_CN'
import stringLength from 'utf8-byte-length'
import regularRemoteValidate from './regular-remote-validate'
import stringRemoteValidate from './string-remote-validate'
import store from '@/store'
import { PARAMETER_TYPES } from '@/dictionary/parameter-types'

/**
 * 前端内置的验证规则，不包含用户自定义的规则
 */
const buildInVaidationRules = {
  length: {
    validate: (value, [length]) => {
      if (Array.isArray(value)) {
        return value?.length <= length
      }
      return stringLength(value) <= length
    }
  },
  repeat: {
    validate: (value, otherValue) => otherValue.findIndex(item => item === value) === -1
  },
  http: {
    validate: value => /^http(s?):\/\/[^\s]+/.test(value)
  },
  isBigger: {
    validate: (value, [targetValue]) => Number(value) > Number(targetValue)
  },
  oid: {
    validate: value => /^(\d+)?(\.\d+)+$/.test(value)
  },
  hourFormat: {
    validate: value => /^[1-5]?[0-9]$/.test(value)
  },
  dayFormat: {
    validate: value => /^((20|21|22|23|[0-1]\d):[0-5][0-9])?$/.test(value)
  },
  repeatTagKey: {
    validate: (value, otherValue) => otherValue.findIndex(item => item === value) === -1
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
    validate: value => /^(?!bk_).*/.test(value)
  },
  ipSearchRuls: {
    validate: (value) => {
      const list = []
      value.trim().split(/\n|;|；|,|，/)
        .forEach((text) => {
          const ip = text.trim()
          ip.length && list.push(ip)
        })
      let isValid = true
      let currentCloudId = null
      list.forEach((text) => {
        let [, cloudId = ''] = text.split(':').reverse()
        cloudId = cloudId || null
        if (currentCloudId && currentCloudId !== cloudId) {
          isValid = false
        }
        currentCloudId = cloudId
      })
      return isValid
    }
  },
  validRegExp: {
    validate: (value) => {
      try {
        new RegExp(value)
        return true
      } catch {
        return false
      }
    }
  }
}

/**
 * 前端内置的验证规则的提示语国际化字典
 */
const dictionary = {
  zh_CN: {
    messages: {
      regex: field => `请输入合法的${field}`,
      length: (field, [maxLength]) => `请输入${maxLength}个字符以内的内容`,
      required: () => '该字段是必填项',
      http: () => '请输入以http(s)://开头的URL',
      isBigger: () => '必须大于最小值',
      repeat: () => '重复的值',
      oid: () => '请输入正确的内容',
      hourFormat: () => '请输入0-59之间的数字',
      dayFormat: () => '请输入00:00-23:59之间的时间',
      namedCharacter: () => '格式不正确，特殊符号仅支持(:_-)',
      min_value: (field, [val]) => `最小值不可小于${val}`,
      max_value: (field, [val]) => `最大值不可超过${val}`,
      repeatTagKey: () => '标签键不能重复',
      setNameMap: () => '集群名称重复',
      emptySetName: () => '请勿输入空白集群名称',
      setNameLen: () => '请输入256个字符以内的内容',
      businessTopoInstNames: () => '格式不正确，不能包含特殊字符 | / : * , < > " ? #及空格',
      reservedWord: () => '不能以"bk_"开头',
      ipSearchRuls: () => '暂不支持不同云区域的混合搜索',
      validRegExp: () => '请输入合法的正则表达式',
      remoteRegular: () => '请输入合法的正则表达式',
      remoteString: () => '请输入符合自定义校验规则的内容'
    },
    custom: {
      asst: {
        required: '请选择关联模型'
      }
    }
  },
  en: {
    messages: {
      regex: () => 'Please enter a valid $ {field}',
      length: (field, [maxLength]) => `Content length max than ${maxLength}`,
      required: () => 'This field is required',
      http: () => 'Please enter a URL beginning with http(s)://',
      isBigger: () => 'Must be greater than the minimum',
      repeat: () => 'This value should not be repeated',
      oid: () => 'Please enter the correct content',
      hourFormat: () => 'Please enter the number between 0-59',
      dayFormat: () => 'Please enter the time between 00:00-23:59',
      min_value: (field, [val]) => `This value is less than the minimum ${val}`,
      max_value: (field, [val]) => `This value is greater than the maximum ${val}`,
      setNameMap: () => 'Duplicate Set name',
      emptySetName: () => 'Do not enter blank Set name',
      repeatTagKey: () => 'Label key cannot be repeated',
      setNameLen: () => 'Content length max than 256',
      reservedWord: () => 'Can not start with "bk_"',
      ipSearchRuls: () => 'Hybrid search of different cloud regions is not supported at the moment',
      validRegExp: () => 'Please enter valid regular express',
      remoteRegular: () => 'Please input valid regular expression',
      remoteString: () => 'Please input correct content that matchs ths custom rules'
    },
    custom: {
      asst: {
        required: 'Please select the associated model'
      }
    }
  }
}

/**
 * 用户可自定义的规则的 key 的集合
 */
const configurableRuleKeys = [
  {
    number: (value, cb) => {
      if (!String(value).length) {
        return true
      }
      return cb()
    }
  },
  ...Object.keys(PARAMETER_TYPES),
  {
    businessTopoInstNames: (value, cb, re) => {
      const values = value.split('\n')
      const list = values.map(text => text.trim()).filter(text => text)
      return list.every(text => re.test(text))
    }
  }
]

/**
 * 混合从远程获取的用户自定义的字段的验证规则
 */
const mixinCustomRules = () => {
  const { globalConfig } = store.state

  // eslint-disable-next-line no-restricted-syntax
  for (const item of configurableRuleKeys) {
    const isFunction = typeof item === 'function'
    const key = isFunction ? Object.keys(item)[0] : item

    if (!globalConfig.config.validationRules[key]) continue

    let validate = (value) => {
      const rule = globalConfig.config.validationRules[key]
      return new RegExp(rule.value).test(value)
    }

    if (isFunction) {
      validate = value => item[key](value, () => {
        const rule = globalConfig.config.validationRules[key]
        return new RegExp(rule.value).test(value)
      }, new RegExp(globalConfig.config.validationRules[key].value))
    }

    // 把用户的自定义规则混入
    buildInVaidationRules[key] = { validate }

    // 提示语设置
    dictionary.zh_CN.messages[key] = (field) => {
      const rule = globalConfig.config.validationRules[key]
      return rule.i18n.cn.replace(/{field}/g, field)
    }

    dictionary.en.messages[key] = (field) => {
      const rule = globalConfig.config.validationRules[key]
      return rule.i18n.en.replace(/{field}/g, field)
    }
  }
}

// 扩展远程验证规则
Validator.extend('remoteRegular', regularRemoteValidate)
Validator.extend('remoteString', stringRemoteValidate, { paramNames: ['regular'] })

export function setupValidator() {
  mixinCustomRules()

  Object.keys(buildInVaidationRules).forEach((ruleKey) => {
    Validator.extend(ruleKey, buildInVaidationRules[ruleKey])
  })

  if (language === 'en') {
    Validator.localize(language)
  } else {
    Validator.localize(language, cnMessages)
  }

  Vue.use(veeValidate, {
    locale: language,
    dictionary
  })
}

export function updateValidator() {
  mixinCustomRules()

  Object.keys(buildInVaidationRules).forEach((ruleKey) => {
    Validator.extend(ruleKey, buildInVaidationRules[ruleKey])
  })
}
