/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

import moment from 'moment-timezone'
import GET_VALUE from 'get-value'
import has from 'has'
import { t } from '@/i18n'
import { CONTAINER_OBJECT_INST_KEYS, CONTAINER_OBJECTS } from '@/dictionary/container'
import { BUILTIN_MODELS, BUILTIN_MODEL_PROPERTY_KEYS } from '@/dictionary/model-constants'
import { TRANSFORM_SPECIAL_HANDLE_OPERATOR } from '@/utils/query-builder-operator'
import { PRESET_TABLE_HEADER_MIN_WIDTH } from '@/dictionary/table-header'
import { PROPERTY_TYPES } from '@/dictionary/property-constants'

/**
 * 获取实例中某个属性的展示值
 * @param {Object} property - 模型具体属性
 * @param {Object} item - 模型实例
 * @return {String} 拍平后的模型属性对应的值
 */

export function getPropertyText(property, item) {
  const propertyId = property.bk_property_id
  const propertyType = property.bk_property_type
  let propertyValue = item[propertyId]
  if (
    propertyValue === null
        || propertyValue === undefined
        || propertyValue === ''
  ) {
    return '--'
  }
  if (propertyType === 'enum') {
    const options = Array.isArray(property.option) ? property.option : []
    const enumOption = options.find(option => option.id === propertyValue)
    propertyValue = enumOption ? enumOption.name : '--'
  } else if (['date', 'time'].includes(propertyType)) {
    propertyValue = formatTime(propertyValue, propertyType === 'date' ? 'YYYY-MM-DD' : 'YYYY-MM-DD HH:mm:ssZZ')
  } else if (propertyType === 'foreignkey') {
    if (Array.isArray(propertyValue)) {
      propertyValue = propertyValue.map(inst => inst.bk_inst_name).join(',')
    } else {
      return String(propertyValue).length ? propertyValue : '--'
    }
  }
  return propertyValue.toString()
}

function getDefaultOptionValue(property) {
  const defaultOption = (property.option || []).find(option => option.is_default)
  if (defaultOption) {
    return defaultOption.id
  }
  return ''
}

function getDefaultOptionMultiValue(property) {
  const defaultOptions = (property.option || []).filter(option => option.is_default)
  return defaultOptions.map(option => option.id)
}

function getDefaultOptionEnumQuoteValue(property) {
  return (property.option || []).map(option => option.bk_inst_id)
}

/**
 * 获取实例的真实值
 * @param {Array} properties - 模型属性
 * @param {Object} inst - 原始实例
 * @param {Boolean} autoSelect - 是否查找默认值作为选中项
 * @return {Object} 实例真实值
 */
export function getInstFormValues(properties, inst = {}, autoSelect = true) {
  const values = {}
  properties.forEach((property) => {
    const propertyId = property.bk_property_id
    const propertyType = property.bk_property_type
    const propertyDefault = property.default
    if (['singleasst', 'multiasst', 'foreignkey'].includes(propertyType)) {
      // const validAsst = (inst[propertyId] || []).filter(asstInst => asstInst.id !== '')
      // values[propertyId] = validAsst.map(asstInst => asstInst['bk_inst_id']).join(',')
    } else if (['date', 'time'].includes(propertyType)) {
      const defaultValue = autoSelect ? propertyDefault : ''
      const formatedTime = formatTime(inst[propertyId], propertyType === 'date' ? 'YYYY-MM-DD' : 'YYYY-MM-DD HH:mm:ssZZ')
      const  value = has(inst, propertyId) ? formatedTime : defaultValue
      values[propertyId] = value || null
    } else if (['int', 'float'].includes(propertyType)) {
      const defaultValue = autoSelect ? propertyDefault : ''
      const  value = has(inst, propertyId) ? inst[propertyId] : defaultValue
      values[propertyId] = value ?? ''
    } else if (['bool'].includes(propertyType)) {
      if ([null, undefined].includes(inst[propertyId]) && autoSelect) {
        values[propertyId] = typeof property.option === 'boolean' ? property.option : false
      } else {
        values[propertyId] = !!inst[propertyId]
      }
    } else if (['enum'].includes(propertyType)) {
      const defaultValue = autoSelect ? getDefaultOptionValue(property) : ''
      values[propertyId] = isNullish(inst[propertyId]) ? defaultValue : inst[propertyId]
    } else if ([PROPERTY_TYPES.ENUMMULTI].includes(propertyType)) {
      const defaultValue = autoSelect ? getDefaultOptionMultiValue(property) : []
      values[propertyId] = isNullish(inst[propertyId]) ? defaultValue : inst[propertyId]
    } else if ([PROPERTY_TYPES.ENUMQUOTE].includes(propertyType)) {
      const defaultValue = autoSelect ? getDefaultOptionEnumQuoteValue(property) : []
      values[propertyId] = isNullish(inst[propertyId]) ? defaultValue : inst[propertyId]
    } else if (['timezone'].includes(propertyType)) {
      const defaultValue = autoSelect ? propertyDefault : ''
      values[propertyId] = isNullish(inst[propertyId]) ? defaultValue : inst[propertyId]
    } else if (['organization'].includes(propertyType)) {
      const defaultValue = autoSelect ? propertyDefault : ''
      const  value = has(inst, propertyId) ? inst[propertyId] : defaultValue
      values[propertyId] = value || null
    } else if (['table'].includes(propertyType)) {
      // table类型的字段编辑和展示目前仅在进程绑定信息被使用，如后期有扩展在其它场景form-table组件与此处都需要调整
      // 接口需要过滤掉不允许编辑及内置的字段
      const tableColumns = property.option?.filter(property => property.editable && !property.bk_isapi)
      // eslint-disable-next-line max-len
      values[propertyId] = (inst[propertyId] || []).map(row => getInstFormValues(tableColumns || [], row, autoSelect))
    } else if (propertyType === PROPERTY_TYPES.INNER_TABLE) {
      const defaultValue = property.option.default || []
      values[propertyId] = isNullish(inst[propertyId]) ? defaultValue : inst[propertyId]
    } else {
      const defaultValue = autoSelect ? propertyDefault : ''
      const  value = has(inst, propertyId) ? inst[propertyId] : defaultValue
      values[propertyId] = value || ''
    }
  })
  return { ...inst, ...values }
}

/**
 * 获取实例的默认值
 * @param {Array} properties - 模型属性
 * @return {Object} 实例默认值
 */
export function getInstFormDefaults(properties) {
  const {
    SINGLECHAR,
    LONGCHAR,
    INT,
    FLOAT,
    OBJUSER
  } = PROPERTY_TYPES
  const defaultValue = {}
  properties.forEach((property) => {
    const {
      bk_property_type: propertyType,
      default: propertyDefault,
      bk_property_id: propertyId
    } = property

    if ([SINGLECHAR, LONGCHAR, INT, FLOAT, OBJUSER].includes(propertyType)) {
      defaultValue[propertyId] = propertyDefault || ''
    }
  })
  return defaultValue
}

export function isEmptyValue(value) {
  return value === '' || value === null || value === void 0
}

export function isNullish(value) {
  return [null, undefined].includes(value)
}


export function formatValue(value, property) {
  if (!(isEmptyValue(value) && property)) {
    return formatPropertyValue(value, property)
  }
  const type = property.bk_property_type
  let formattedValue = value
  switch (type) {
    case 'enum':
    case 'int':
    case 'float':
    case 'list':
    case 'time':
    case PROPERTY_TYPES.ENUMMULTI:
      formattedValue = null
      break
    case 'bool':
      formattedValue = false
      break
    default:
      break
  }
  return formattedValue
}

export function getPropertyDefaultValue(property, value) {
  const propertyValue = formatPropertyValue(value, property)
  const defaultValue = getInstFormValues([property])?.[property.bk_property_id]
  // undefined 认为没有传递属性值，与 null 等假值明确区分开
  return value === undefined ? defaultValue : propertyValue
}

export function formatPropertyValue(value, property) {
  // 枚举引用/多选和组织类型的字段保存时必须转换为数组，在作为form的值使用时如果是单选值不是数组格式在这里统一转换
  const arrayValueTypes = [PROPERTY_TYPES.ENUMQUOTE, PROPERTY_TYPES.ENUMMULTI, PROPERTY_TYPES.ORGANIZATION]
  if (arrayValueTypes.includes(property?.bk_property_type)) {
    return !Array.isArray(value) ? [value] : value
  }
  return value
}

export function formatValues(values, properties) {
  const formatted = { ...values }
  properties.forEach((property) => {
    const key = property.bk_property_id
    if (has(formatted, key)) {
      formatted[key] = formatValue(formatted[key], property)
    }
  })
  return formatted
}

/**
 * 格式化时间
 * @param {String} originalTime - 需要被格式化的时间
 * @param {String} format - 格式化类型
 * @param {String} timezone - 时区
 * @return {String} 格式化后的时间
 */
export function formatTime(originalTime, format = 'YYYY-MM-DD HH:mm:ssZZ', timezone) {
  if (!originalTime) {
    return ''
  }
  const dateObj = moment(originalTime)

  if (!dateObj.isValid()) return 'Invalid Date'

  const targetTimezone = timezone || window.Site.timezone || moment.tz.guess()

  return dateObj.tz(targetTimezone).format(format)
}
/**
 * 从模型属性中获取指定id的属性对象
 * @param {Array} properties - 模型属性
 * @param {String} id - 属性id
 * @return {Object} 模型属性对象
 */
export function getProperty(properties, id) {
  return properties.find(property => property.bk_property_id === id)
}

/**
 * 获取指定属性的枚举列表
 * @param {Array} properties - 模型属性
 * @param {String} id - 属性id
 * @return {Array} 枚举列表
 */
export function getEnumOptions(properties, id) {
  const property = getProperty(properties, id)
  return property ? property.option || [] : []
}

/**
 * 获取属性展示优先级
 * @param {Object} property - 模型属性
 * @return {Number} 优先级分数，越小优先级越高
 */
export function getPropertyPriority(property) {
  let priority = property.bk_property_index ?? 0
  if (property.isonly) {
    priority = priority - 1
  }
  if (property.isrequired) {
    priority = priority - 1
  }
  return priority
}

/**
 * 获取模型默认表头属性
 * @param {Array} properties - 模型属性
 * @return {Array} 默认展示的前六个模型属性
 */
export function getDefaultHeaderProperties(properties) {
  return [...properties].sort((A, B) => getPropertyPriority(A) - getPropertyPriority(B)).slice(0, 6)
}

/**
 * 获取模型用户自定义表头
 * @param {Array} properties - 模型属性
 * @param {Array} customColumns - 自定义表头
 * @return {Array} 自定义表头属性
 */
export function getCustomHeaderProperties(properties, customColumns) {
  const columnProperties = []
  customColumns.forEach((propertyId) => {
    const columnProperty = properties.find(property => property.bk_property_id === propertyId)
    if (columnProperty) {
      columnProperties.push(columnProperty)
    }
  })
  return columnProperties
}

/**
 * 获取模型表头
 * @param {Array} properties - 模型属性
 * @param {Array} customColumns - 自定义表头
 * @param {Array} fixedPropertyIds - 需固定在表格前面的属性ID
 * @return {Array} 表头属性
 */
export function getHeaderProperties(properties, customColumns, fixedPropertyIds = []) {
  let headerProperties
  if (customColumns && customColumns.length) {
    headerProperties = getCustomHeaderProperties(properties, customColumns)
  } else {
    headerProperties = getDefaultHeaderProperties(properties)
  }
  if (fixedPropertyIds.length) {
    headerProperties = headerProperties.filter(property => !fixedPropertyIds.includes(property.bk_property_id))
    const fixedProperties = []
    fixedPropertyIds.forEach((id) => {
      const property = properties.find(property => property.bk_property_id === id)
      if (property) {
        fixedProperties.push(property)
      }
    })
    return [...fixedProperties, ...headerProperties]
  }
  return headerProperties
}

export function getHeaderPropertyName(property) {
  if (!property?.bk_property_name?.endsWith(`(${property.unit})`) && property.unit) {
    return `${property.bk_property_name}(${property.unit})`
  }
  return property.bk_property_name
}

export function getHeaderPropertyMinWidth(property, options = {}) {
  const { fontSize = 12, hasSort = false, offset = 30, name, min = 0, preset = {} } = options

  // 预设的固定宽度不需要计算直接使用
  const presetMinWidth = { ...PRESET_TABLE_HEADER_MIN_WIDTH, ...preset }
  if (presetMinWidth[property.bk_property_id]) {
    return presetMinWidth[property.bk_property_id]
  }

  const content = name ?? getHeaderPropertyName(property)

  // 字母数字和空白字符的个数
  const letterCount = (content.match(/[\w\s\\(\\)]/g) ?? []).join('').length

  const totalCount = content?.length ?? 0

  // 分别按字母与非字母计算字符占用的总宽度
  const contentWidth = ((totalCount - letterCount) * fontSize) + (letterCount * fontSize * 0.7)

  const objKeyMap = { ...CONTAINER_OBJECT_INST_KEYS, ...BUILTIN_MODEL_PROPERTY_KEYS }
  const baseWidth = (property.bk_property_id === objKeyMap[property.bk_obj_id]?.ID ?? 'bk_inst_id') ? 50 : contentWidth

  const finalWidth = baseWidth + (hasSort ? 22 : 0) + offset

  return Math.ceil(Math.max(finalWidth, min))
}

/**
 * 深拷贝
 * @param {Object} object - 需拷贝的对象
 * @return {Object} 拷贝后的对象
 */
export function clone(object) {
  return JSON.parse(JSON.stringify(object))
}

export function getValidateEvents(property) {
  const type = property.bk_property_type
  const isChar = ['singlechar', 'longchar'].includes(type)
  const hasRegular = !!property.option
  const isSelectType = [PROPERTY_TYPES.ENUMMULTI, PROPERTY_TYPES.ENUMQUOTE, PROPERTY_TYPES.ORGANIZATION].includes(type)
  if ((isChar && hasRegular) || isSelectType) {
    return {
      'data-vv-validate-on': 'change|blur'
    }
  }
  return {}
}

/**
 * 根据远程返回的属性生成对应的校验规则
 * @param {Object} property 字段属性
 * @param {String} property.bk_property_type 字段类型
 * @param {String} property.option 额外选项
 * @param {String} property.isrequired 是否必须
 * @returns {Array} vee-validate 规则
 */
export function getValidateRules(property) {
  const rules = {}
  const {
    bk_property_type: propertyType,
    option,
    isrequired,
    ismultiple
  } = property

  if (isrequired) {
    rules.required = true
  }

  const isSelectType = [
    PROPERTY_TYPES.ENUMMULTI,
    PROPERTY_TYPES.ENUMQUOTE,
    PROPERTY_TYPES.ORGANIZATION
  ].includes(propertyType)

  if (option) {
    if (['int', 'float'].includes(propertyType)) {
      if (has(option, 'min') && !['', null, undefined].includes(option.min)) {
        rules.min_value = option.min
      }
      if (has(option, 'max') && !['', null, undefined].includes(option.max)) {
        rules.max_value = option.max
      }
    } else if (['singlechar', 'longchar'].includes(propertyType)) {
      rules.remoteString = option
    }
  }
  if (['singlechar', 'longchar'].includes(propertyType)) {
    rules[propertyType] = true
    rules.length = propertyType === 'singlechar' ? 256 : 2000
  } else if (propertyType === 'int') {
    rules.number = true
  } else if (propertyType === 'float') {
    rules.float = true
  } else if (propertyType === 'objuser') {
    rules.length = 2000
  } else if (isSelectType) {
    rules.maxSelectLength = ismultiple ? -1 : 1
  }

  return rules
}

export function getSort(sort, defaultSort = {}) {
  const order = sort.order || defaultSort.order || 'ascending'
  const prop = sort.prop || defaultSort.prop || ''
  if (prop && order === 'descending') {
    return `-${prop}`
  }
  return prop
}

export function getValue() {
  // eslint-disable-next-line prefer-rest-params, new-cap
  return GET_VALUE(...arguments)
}

export function transformHostSearchParams(params) {
  const transformedParams = clone(params)
  const conditions = transformedParams.condition
  conditions.forEach((item) => {
    item.condition.forEach((field) => {
      const { bk_obj_id: objId } = item
      const { operator, value } = field
      field.operator = getOperator(objId, operator)
      if (['$in', '$nin', '$multilike'].includes(operator) && !Array.isArray(value)) {
        field.value = value.split(/\n|;|；|,|，/).filter(str => str.trim().length)
          .map(str => str.trim())
      }
    })
  })
  return transformedParams
}

export function getOperator(objId, operator) {
  if (objId !== BUILTIN_MODELS.HOST) {
    // 主机接口参数condition特殊处理 非host情况下： 1. 操作符为$regex处理成contains_s  2. 操作符为contains处理成$regex
    return transformNoHostOperator(operator)
  }
  return operator
}

const transformNoHostOperator = operator => TRANSFORM_SPECIAL_HANDLE_OPERATOR?.[operator] ?? operator

const defaultPaginationConfig = window.innerHeight > 750
  ? { limit: 20, 'limit-list': [20, 50, 100, 500] }
  : { limit: 10, 'limit-list': [10, 50, 100, 500] }
export function getDefaultPaginationConfig(customConfig = {}, useQuery = true) {
  const RouterQuery = require('@/router/query').default
  const config = {
    count: 0,
    current: useQuery ? parseInt(RouterQuery.get('page', 1), 10) : 1,
    limit: useQuery ? parseInt(RouterQuery.get('limit', defaultPaginationConfig.limit), 10) : defaultPaginationConfig.limit,
    'limit-list': customConfig['limit-list'] || defaultPaginationConfig['limit-list']
  }
  return config
}

export function getPageParams(pagination) {
  return {
    start: (pagination.current - 1) * pagination.limit,
    limit: pagination.limit
  }
}

export function localSort(data, compareKey) {
  return data.sort((A, B) => {
    if (has(A, compareKey) && has(B, compareKey)) {
      return A[compareKey].localeCompare(B[compareKey], 'zh-Hans-CN', { sensitivity: 'accent', caseFirst: 'lower' })
    }
    return 0
  })
}

export function sort(data, compareKey) {
  return [...data].sort((A, B) => A[compareKey] - B[compareKey])
}

export function versionSort(data, compareKey) {
  return data.sort((a, b) => {
    let i = 0
    const arr1 = a[compareKey].split('.')
    const arr2 = b[compareKey].split('.')
    while (true) {
      const s1 = arr1[i]
      const s2 = arr2[i]
      i = i + 1
      if (s1 === undefined || s2 === undefined) {
        return arr2.length - arr1.length
      }
      if (s1 === s2) continue
      return s2 - s1
    }
  })
}

/**
 * 递归对拓扑树进行自然排序
 * @param {array} topoTree 拓扑
 * @param {string} compareKey 需要对比的属性的 Key
 * @param {string} [childrenKey] 后代属性的 Key
 */
export function sortTopoTree(topoTree, compareKey, childrenKey) {
  if (!Array.isArray(topoTree)) {
    throw Error('topoTree must be type of array')
  }
  if (!compareKey) {
    throw Error('compareKey is required')
  }

  topoTree.sort((a, b) => {
    if (has(a, compareKey) && has(b, compareKey)) {
      const valueA = a[compareKey]
      const valueB = b[compareKey]

      if (/[a-zA-Z0-9]/.test(valueA) || /[a-zA-Z0-9]/.test(valueB)) {
        if (valueA > valueB) return 1
        if (valueA < valueB) return -1
        return 0
      }

      return valueA.localeCompare(valueB)
    }
    return 0
  })

  if (childrenKey) {
    topoTree?.forEach((node) => {
      if (node[childrenKey]) {
        sortTopoTree(node[childrenKey], compareKey, childrenKey)
      }
    })
  }
}

export function isEmptyPropertyValue(originalValue) {
  return originalValue === ''
    || originalValue === null
    || originalValue === undefined
    || (Array.isArray(originalValue) && originalValue.length === 0)
}

export function getPropertyCopyValue(originalValue, propertyType, options = {}) {
  if (isEmptyPropertyValue(originalValue)) {
    return '--'
  }
  const type = typeof propertyType === 'string' ? propertyType : propertyType.bk_property_type
  let value
  switch (type) {
    case 'date':
      value = formatTime(originalValue, 'YYYY-MM-DD')
      break
    case 'time':
      value = formatTime(originalValue, 'YYYY-MM-DD HH:mm:ssZZ')
      break
    case 'foreignkey': {
      if (options.isFullCloud) {
        value = (originalValue || []).map(cloud => `${cloud.bk_inst_name}[${cloud.bk_inst_id}]`).join(',')
      } else {
        value = (originalValue || []).map(cloud => cloud.bk_inst_id).join(',')
      }
      break
    }
    case 'list':
    case 'objuser':
    case 'organization':
      value = originalValue.toString()
      break
    case 'map': {
      const pair = []
      for (const [key, val] of Object.entries(originalValue)) {
        pair.push(`${key}: ${val}`)
      }
      value = pair.join('\n')
      break
    }
    case 'enum':
      value =  propertyType.option.find(item => item.id === originalValue).name
      break
    case 'enummulti':
      value = originalValue.map(value => propertyType.option.find(item => item.id === value).name).join('\n')
      break
    case 'array':
    case 'object':
      value = JSON.stringify(originalValue, null, 2)
      break
    default:
      value = originalValue
  }
  return value
}

// 使用独立组件展示value的类型
export function isUseComplexValueType(property) {
  const types = [
    PROPERTY_TYPES.OBJUSER,
    PROPERTY_TYPES.TABLE,
    PROPERTY_TYPES.SERVICE_TEMPLATE,
    PROPERTY_TYPES.ORGANIZATION,
    PROPERTY_TYPES.MAP,
    PROPERTY_TYPES.ENUMQUOTE
  ]
  return types.includes(property.bk_property_type)
}

export function isShowOverflowTips(property) {
  const otherTypes = [PROPERTY_TYPES.TOPOLOGY]
  return !isUseComplexValueType(property) && !otherTypes.includes(property.bk_property_type)
}

export function getPropertyPlaceholder(property) {
  if (!property) {
    return ''
  }
  const placeholderTxt = [
    PROPERTY_TYPES.ENUM,
    PROPERTY_TYPES.ENUMMULTI,
    PROPERTY_TYPES.ENUMQUOTE,
    PROPERTY_TYPES.ORGANIZATION,
    PROPERTY_TYPES.LIST,
    PROPERTY_TYPES.DATE,
    PROPERTY_TYPES.TIME,
    PROPERTY_TYPES.TIMEZONE
  ].includes(property.bk_property_type) ? '请选择xx' : '请输入xx'
  return t(placeholderTxt, { name: property.bk_property_name })
}

export function getPropertyDefaultEmptyValue(_property) {
  return ''
}

export function isPropertySortable(property) {
  if (property.bk_obj_id === BUILTIN_MODELS.HOST) {
    return ![PROPERTY_TYPES.FOREIGNKEY, PROPERTY_TYPES.TOPOLOGY, PROPERTY_TYPES.INNER_TABLE]
      .includes(property.bk_property_type)
  }

  return ![PROPERTY_TYPES.INNER_TABLE].includes(property.bk_property_type)
}

export function isIconTipProperty(type) {
  return ['innertable', 'bool'].includes(type)
}

/**
 * 判断是否为容器字段
 */
export function isContainerObjects(objId) {
  return Object.values(CONTAINER_OBJECTS).includes(objId)
}

/**
 * 下拉框是否要展示全部选项
 * 需要展示的： 由用户自定义数据的枚举、枚举(多选)、列表三个类型的字段
 */
export function getSelectAll(property = {}) {
  if (!property) return false

  const fields = [
    PROPERTY_TYPES.ENUM,
    PROPERTY_TYPES.ENUMMULTI,
    PROPERTY_TYPES.LIST
  ]
  const { bk_property_type: type, ispre } = property
  return !ispre && fields.includes(type)
}

// 组织选择器回显
export function parseOrgVal(data) {
  const { orgPath, name, organization_path: path } = data
  if (orgPath) return orgPath
  if (path) return `${path}/${name}`
  return name
}

export default {
  getProperty,
  getPropertyText,
  getPropertyPriority,
  getEnumOptions,
  getDefaultHeaderProperties,
  getCustomHeaderProperties,
  getHeaderProperties,
  getHeaderPropertyName,
  formatTime,
  clone,
  getInstFormValues,
  getInstFormDefaults,
  formatValue,
  formatValues,
  getValidateRules,
  getValidateEvents,
  getSort,
  getValue,
  transformHostSearchParams,
  getDefaultPaginationConfig,
  getPageParams,
  localSort,
  sort,
  getPropertyCopyValue,
  isEmptyPropertyValue,
  getHeaderPropertyMinWidth,
  isShowOverflowTips,
  isUseComplexValueType,
  getPropertyPlaceholder,
  getPropertyDefaultValue,
  versionSort,
  isPropertySortable,
  isIconTipProperty,
  isContainerObjects,
  getSelectAll,
  parseOrgVal
}
