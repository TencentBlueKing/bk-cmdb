import store from '@/store'
import i18n from '@/i18n'
import isInt from 'validator/es/lib/isInt'

const getModelById = store.getters['objectModelClassify/getModelById']
export function getLabel(property) {
  const model = getModelById(property.bk_obj_id) || {}
  return `${model.bk_obj_name} - ${property.bk_property_name}`
}

export function getBindProps(property) {
  if (!property) {
    return {}
  }
  const type = property.bk_property_type
  if (['list', 'enum'].includes(type)) {
    return {
      options: property.option || []
    }
  }
  if (type === 'objuser') {
    return {
      fastSelect: true
    }
  }
  return {}
}

export function getPlaceholder(property) {
  if (!property) {
    return ''
  }
  const selectTypes = ['list', 'enum', 'timezone', 'organization', 'date', 'time']
  if (selectTypes.includes(property.bk_property_type)) {
    return i18n.t('请选择xx', { name: property.bk_property_name })
  }
  return i18n.t('请输入xx', { name: property.bk_property_name })
}

export function getDefaultData(property, defaultData = { operator: '$in', value: [] }) {
  const EQ = '$eq'
  const RANGE = '$range'
  const IN = '$in'
  const defaultMap = {
    singlechar: { operator: IN, value: [] },
    int: { operator: EQ, value: '' },
    float: { operator: EQ, value: '' },
    enum: { operator: IN, value: [] },
    date: { operator: RANGE, value: [] },
    time: { operator: RANGE, value: [] },
    longchar: { operator: IN, value: [] },
    objuser: { operator: IN, value: [] },
    timezone: { operator: IN, value: [] },
    bool: { operator: EQ, value: '' },
    list: { operator: IN, value: [] },
    organization: { operator: IN, value: [] },
  }
  return {
    operator: defaultData.operator,
    value: defaultData.value,
    ...defaultMap[property.bk_property_type]
  }
}

export function getOperatorSideEffect(property, operator, value) {
  let effectValue = value
  if (operator === '$range') {
    effectValue = []
  } else {
    const defaultValue = this.getDefaultData(property).value
    const isTypeChanged = (Array.isArray(defaultValue)) !== (Array.isArray(value))
    effectValue = isTypeChanged ? defaultValue : value
  }
  return effectValue
}

export function convertValue(value, operator, property) {
  const { bk_property_type: type } = property
  let convertedValue = Array.isArray(value) ? value : [value]
  convertedValue = convertedValue.map((data) => {
    if (['int', 'foreignkey', 'organization', 'service-template'].includes(type)) {
      return parseInt(data, 10)
    }
    if (type === 'float') {
      return parseFloat(data, 10)
    }
    if (type === 'bool') {
      return data === 'true'
    }
    return data
  })
  if (['$in', '$nin', '$range'].includes(operator)) {
    return convertedValue
  }
  return convertedValue[0]
}

export function findProperty(id, properties) {
  const field = isInt(id) ? 'id' : 'bk_property_id'
  return properties.find(property => property[field].toString() === id.toString())
}

export function findPropertyByPropertyId(propertyId, properties, modelId) {
  if (modelId) {
    return properties.find(property => property.bk_obj_id === modelId && property.bk_property_id === propertyId)
  }
  return properties.find(property => property.bk_property_id === propertyId)
}

function createTimeCondition() {
  return { oper: 'and', rules: [] }
}
export function transformCondition(condition, properties, header) {
  const conditionMap = {
    host: [],
    module: [],
    set: [],
    biz: [],
    object: []
  }
  const timeCondition = {
    host: createTimeCondition(),
    module: createTimeCondition(),
    set: createTimeCondition(),
    biz: createTimeCondition(),
    object: createTimeCondition()
  }
  Object.keys(condition).forEach((id) => {
    const property = findProperty(id, properties)
    const { operator, value } = condition[id]
    if (value === null || !value.toString().length) return
    // 时间类型的字段需要上升一层单独处理
    if (property.bk_property_type === 'time') {
      const [start, end] = value
      timeCondition[property.bk_obj_id].rules.push({
        field: property.bk_property_id,
        start,
        end
      })
      return
    }
    const submitCondition = conditionMap[property.bk_obj_id]
    if (operator === '$range') {
      const [start, end] = value
      submitCondition.push({
        field: property.bk_property_id,
        operator: '$gte',
        value: start
      }, {
        field: property.bk_property_id,
        operator: '$lte',
        value: end
      })
    } else {
      submitCondition.push({
        field: property.bk_property_id,
        operator,
        value
      })
    }
  })
  return Object.keys(conditionMap).map((modelId) => {
    const condition = {
      bk_obj_id: modelId,
      fields: header.filter(property => property.bk_obj_id === modelId).map(property => property.bk_property_id),
      condition: conditionMap[modelId]
    }
    if (timeCondition[modelId].rules.length) {
      condition.time_condition = timeCondition[modelId]
    }
    return condition
  })
}

export function splitIP(raw) {
  const list = []
  raw.trim().split(/\n|;|；|,|，/)
    .forEach((text) => {
      const ip = text.trim()
      ip.length && list.push(ip)
    })
  return list
}

export function transformIP(raw) {
  const transformedIP = {
    data: [],
    condition: null
  }
  const list = splitIP(raw)
  list.forEach((text) => {
    const [IP, cloudId] = text.split(':').reverse()
    transformedIP.data.push(IP)
    // 当前的查询接口对于形如 0:ip0  1:ip1 的输入
    // 拆分后实际的查询结果是云区域id与ip的排列组合形式:0+ip0, 0+ip1, 1+ip0, 1+ip1
    // 因此实际传入的云区域id不能重复，只用设置一次conditon即可
    if (cloudId && !transformedIP.condition) {
      transformedIP.condition = {
        field: 'bk_cloud_id',
        operator: '$eq',
        value: parseInt(cloudId, 10)
      }
    }
  })
  return transformedIP
}

const operatorSymbolMap = {
  $eq: '=',
  $ne: '≠',
  $in: '*=',
  $nin: '*≠',
  $gt: '>',
  $lt: '<',
  $gte: '≥',
  $lte: '≤',
  $regex: 'Like',
  $range: '≤ ≥'
}
export function getOperatorSymbol(operator) {
  return operatorSymbolMap[operator]
}

export function getDefaultIP() {
  return {
    text: '',
    inner: true,
    outer: true,
    exact: true
  }
}

export function defineProperty(definition) {
  return Object.assign({}, {
    id: null,
    bk_obj_id: null,
    bk_property_id: null,
    bk_property_name: null,
    bk_property_index: -1,
    bk_property_type: 'singlechar',
    isonly: true,
    ispre: true,
    bk_isapi: true,
    bk_issystem: true,
    isreadonly: true,
    editable: false,
    bk_property_group: 'default',
    is_custom: true
  }, definition)
}

export function getUniqueProperties(preset, dynamic) {
  const unique = dynamic.filter(property => !preset.includes(property))
  const full = [...preset, ...unique]
  const ids = [...new Set(full.map(property => property.id))]
  return ids.map(id => full.find(property => property.id === id))
}

function getPropertyPriority(property) {
  let priority = 0
  if (property.isonly) {
    priority = priority - 1
  }
  if (property.isrequired) {
    priority = priority - 1
  }
  return priority
}
export function getInitialProperties(properties) {
  // eslint-disable-next-line max-len
  return [...properties].sort((propertyA, propertyB) => getPropertyPriority(propertyA) - getPropertyPriority(propertyB)).slice(0, 6)
}

export default {
  getLabel,
  getPlaceholder,
  getBindProps,
  getDefaultData,
  getOperatorSideEffect,
  convertValue,
  findProperty,
  findPropertyByPropertyId,
  transformCondition,
  transformIP,
  getOperatorSymbol,
  splitIP,
  getDefaultIP,
  defineProperty,
  getUniqueProperties,
  getInitialProperties
}
