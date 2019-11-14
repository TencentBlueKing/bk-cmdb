import moment from 'moment'
import GET_VALUE from 'get-value'

/**
 * 拍平列表
 * @param {Array} properties - 模型属性
 * @param {Array} list - 模型实例列表
 * @return {Array} 拍平后的模型实例列表
 */

export function getFullName (names) {
    if (!names) return ''
    const userList = window.CMDB_USER_LIST || [] // set in setup/preload.js
    const enNames = String(names).split(',')
    const fullNames = enNames.map(enName => {
        const user = userList.find(user => user['english_name'] === enName)
        if (user) {
            return `${user['english_name']}(${user['chinese_name']})`
        }
        return enName
    })
    return fullNames.join(',')
}

export function flattenList (properties, list) {
    if (!list.length) return list
    const flattenedList = clone(list)
    flattenedList.forEach((item, index) => {
        flattenedList[index] = flattenItem(properties, item)
    })
    return flattenedList
}

/**
 * 拍平实例具体的属性
 * @param {Object} properties - 模型具体属性
 * @param {Object} item - 模型实例
 * @return {Object} 拍平后的模型实例
 */
export function flattenItem (properties, item) {
    const flattenedItem = clone(item)
    properties.forEach(property => {
        flattenedItem[property['bk_property_id']] = getPropertyText(property, flattenedItem)
    })
    return flattenedItem
}

/**
 * 获取实例中某个属性的展示值
 * @param {Object} property - 模型具体属性
 * @param {Object} item - 模型实例
 * @return {String} 拍平后的模型属性对应的值
 */

export function getPropertyText (property, item) {
    const propertyId = property['bk_property_id']
    const propertyType = property['bk_property_type']
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
    } else if (['singleasst', 'multiasst'].includes(propertyType)) {
        const values = Array.isArray(propertyValue) ? propertyValue : []
        propertyValue = values.map(inst => inst['bk_inst_name']).join(',')
    } else if (['date', 'time'].includes(propertyType)) {
        propertyValue = formatTime(propertyValue, propertyType === 'date' ? 'YYYY-MM-DD' : 'YYYY-MM-DD HH:mm:ss')
    } else if (propertyType === 'objuser') {
        propertyValue = getFullName(propertyValue)
    } else if (propertyType === 'foreignkey') {
        if (Array.isArray(propertyValue)) {
            propertyValue = propertyValue.map(inst => inst['bk_inst_name']).join(',')
        } else {
            return String(propertyValue).length ? propertyValue : '--'
        }
    }
    return propertyValue.toString()
}

/**
 * 拍平主机列表
 * @param {Object} properties - 模型属性,eg: {host: [], biz: []}
 * @param {Array} list - 模型实例列表
 * @return {Array} 拍平后的模型实例列表
 */
export function flattenHostList (properties, list) {
    if (!list.length) return list
    const flattenedList = clone(list)
    flattenedList.forEach((item, index) => {
        flattenedList[index] = flattenHostItem(properties, item)
    })
    return flattenedList
}

/**
 * 拍平主机实例具体的属性
 * @param {Object} property - 模型具体属性
 * @param {Object} item - 模型实例
 * @return {Object} 拍平后的模型实例
 */
export function flattenHostItem (properties, item) {
    const flattenedItem = clone(item)
    for (const objId in properties) {
        properties[objId].forEach(property => {
            const originalValue = item[objId] instanceof Array ? item[objId] : [item[objId]]
            originalValue.forEach(value => {
                value[property['bk_property_id']] = getPropertyText(property, value)
            })
        })
    }
    return flattenedItem
}

/**
 * 获取实例的真实值
 * @param {Array} properties - 模型属性
 * @param {Object} inst - 原始实例
 * @return {Object} 实例真实值
 */
export function getInstFormValues (properties, inst = {}) {
    const values = {}
    properties.forEach(property => {
        const propertyId = property['bk_property_id']
        const propertyType = property['bk_property_type']
        if (['singleasst', 'multiasst', 'foreignkey'].includes(propertyType)) {
            // const validAsst = (inst[propertyId] || []).filter(asstInst => asstInst.id !== '')
            // values[propertyId] = validAsst.map(asstInst => asstInst['bk_inst_id']).join(',')
        } else if (['date', 'time'].includes(propertyType)) {
            const formatedTime = formatTime(inst[propertyId], propertyType === 'date' ? 'YYYY-MM-DD' : 'YYYY-MM-DD HH:mm:ss')
            values[propertyId] = formatedTime || null
        } else if (['int', 'float'].includes(propertyType)) {
            values[propertyId] = ['', undefined].includes(inst[propertyId]) ? null : inst[propertyId]
        } else if (['bool'].includes(propertyType)) {
            values[propertyId] = !!inst[propertyId]
        } else if (['enum'].includes(propertyType)) {
            values[propertyId] = [null].includes(inst[propertyId]) ? '' : inst[propertyId]
        } else {
            values[propertyId] = inst.hasOwnProperty(propertyId) ? inst[propertyId] : ''
        }
    })
    return values
}

/**
 * 格式化时间
 * @param {String} originalTime - 需要被格式化的时间
 * @param {String} format - 格式化类型
 * @return {String} 格式化后的时间
 */
export function formatTime (originalTime, format = 'YYYY-MM-DD HH:mm:ss') {
    if (!originalTime) {
        return ''
    }
    const formatedTime = moment(originalTime).format(format)
    if (formatedTime === 'Invalid date') {
        return originalTime
    } else {
        return formatedTime
    }
}
/**
 * 从模型属性中获取指定id的属性对象
 * @param {Array} properties - 模型属性
 * @param {String} id - 属性id
 * @return {Object} 模型属性对象
 */
export function getProperty (properties, id) {
    return properties.find(property => property['bk_property_id'] === id)
}

/**
 * 获取指定属性的枚举列表
 * @param {Array} properties - 模型属性
 * @param {String} id - 属性id
 * @return {Array} 枚举列表
 */
export function getEnumOptions (properties, id) {
    const property = getProperty(properties, id)
    return property ? property.option || [] : []
}

/**
 * 获取属性展示优先级
 * @param {Object} property - 模型属性
 * @return {Number} 优先级分数，越小优先级越高
 */
export function getPropertyPriority (property) {
    let priority = 0
    if (property['isonly']) {
        priority--
    }
    if (property['isrequired']) {
        priority--
    }
    return priority
}

/**
 * 获取模型默认表头属性
 * @param {Array} properties - 模型属性
 * @return {Array} 默认展示的前六个模型属性
 */
export function getDefaultHeaderProperties (properties) {
    return [...properties].sort((propertyA, propertyB) => {
        return getPropertyPriority(propertyA) - getPropertyPriority(propertyB)
    }).slice(0, 6)
}

/**
 * 获取模型用户自定义表头
 * @param {Array} properties - 模型属性
 * @param {Array} customColumns - 自定义表头
 * @return {Array} 自定义表头属性
 */
export function getCustomHeaderProperties (properties, customColumns) {
    const columnProperties = []
    customColumns.forEach(propertyId => {
        const columnProperty = properties.find(property => property['bk_property_id'] === propertyId)
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
export function getHeaderProperties (properties, customColumns, fixedPropertyIds = []) {
    let headerProperties
    if (customColumns && customColumns.length) {
        headerProperties = getCustomHeaderProperties(properties, customColumns)
    } else {
        headerProperties = getDefaultHeaderProperties(properties)
    }
    if (fixedPropertyIds.length) {
        headerProperties = headerProperties.filter(property => !fixedPropertyIds.includes(property['bk_property_id']))
        const fixedProperties = []
        fixedPropertyIds.forEach(id => {
            const property = properties.find(property => property['bk_property_id'] === id)
            if (property) {
                fixedProperties.push(property)
            }
        })
        return [...fixedProperties, ...headerProperties]
    }
    return headerProperties
}

/**
 * 深拷贝
 * @param {Object} object - 需拷贝的对象
 * @return {Object} 拷贝后的对象
 */
export function clone (object) {
    return JSON.parse(JSON.stringify(object))
}

/**
 * 获取对象中的metada.label.bk_biz_id属性
 * @param {Object} object - 需拷贝的对象
 * @return {Object} 拷贝后的对象
 */
export function getMetadataBiz (object = {}) {
    const metadata = object.metadata || {}
    const label = metadata.label || {}
    const biz = label['bk_biz_id']
    return biz
}

export function getValidateRules (property) {
    const rules = {}
    const {
        bk_property_type: propertyType,
        option,
        isrequired
    } = property
    if (isrequired) {
        rules.required = true
    }
    if (option) {
        if (propertyType === 'int') {
            if (option.hasOwnProperty('min') && !['', null, undefined].includes(option.min)) {
                rules['min_value'] = option.min
            }
            if (option.hasOwnProperty('max') && !['', null, undefined].includes(option.max)) {
                rules['max_value'] = option.max
            }
        } else if (['singlechar', 'longchar'].includes(propertyType)) {
            rules['regex'] = option
        }
    }
    if (['singlechar', 'longchar'].includes(propertyType)) {
        rules[propertyType] = true
        rules.length = propertyType === 'singlechar' ? 256 : 2000
    } else if (propertyType === 'int') {
        rules.numeric = true
    } else if (propertyType === 'float') {
        rules.float = true
    }
    return rules
}

export function getSort (sort) {
    const order = sort.order
    const prop = sort.prop
    if (!prop) {
        return ''
    }
    if (order === 'descending') {
        return `-${prop}`
    }
    return prop
}

export function getValue () {
    return GET_VALUE(...arguments)
}

export function transformHostSearchParams (params) {
    const transformedParams = clone(params)
    const conditions = transformedParams.condition
    conditions.forEach(item => {
        item.condition.forEach(field => {
            const operator = field.operator
            const value = field.value
            if (['$in', '$multilike'].includes(operator) && !Array.isArray(value)) {
                field.value = value.split('\n').filter(str => str.trim().length).map(str => str.trim())
            }
        })
    })
    return transformedParams
}

const defaultPaginationConfig = window.innerHeight > 750
    ? { limit: 20, 'limit-list': [20, 50, 100, 500] }
    : { limit: 10, 'limit-list': [10, 50, 100, 500] }
export function getDefaultPaginationConfig () {
    return {
        current: 1,
        count: 0,
        ...defaultPaginationConfig
    }
}

export function getPageParams (pagination) {
    return {
        start: (pagination.current - 1) * pagination.limit,
        limit: pagination.limit
    }
}

export default {
    getProperty,
    getPropertyText,
    getPropertyPriority,
    getEnumOptions,
    getDefaultHeaderProperties,
    getCustomHeaderProperties,
    getHeaderProperties,
    flattenList,
    flattenItem,
    flattenHostList,
    flattenHostItem,
    formatTime,
    clone,
    getInstFormValues,
    getMetadataBiz,
    getValidateRules,
    getSort,
    getValue,
    transformHostSearchParams,
    getDefaultPaginationConfig,
    getPageParams
}
