import moment from 'moment'
/**
 * 拍平列表
 * @param {Array} properties - 模型属性
 * @param {Array} list - 模型实例列表
 * @return {Array} 拍平后的模型实例列表
 */
export function flatternList (properties, list) {
    if (!list.length) return list
    const flatternedList = clone(list)
    properties.forEach(property => {
        flatternedList.forEach((item, index) => {
            flatternedList[index] = flatternItem(property, item)
        })
    })
    return flatternedList
}

/**
 * 拍平实例具体的属性
 * @param {Object} property - 模型具体属性
 * @param {Object} item - 模型实例
 * @return {Object} 拍平后的模型实例
 */
export function flatternItem (property, item) {
    const flatternedItem = clone(item)
    const propertyId = property['bk_property_id']
    const propertyType = property['bk_property_type']
    if (propertyType === 'enum') {
        const enumOption = (property.option || []).find(option => option.id === flatternedItem[propertyId])
        flatternedItem[propertyId] = enumOption ? enumOption.name : null
    } else if (['singleasst', 'multiasst'].includes(propertyType)) {
        flatternedItem[propertyId] = (flatternedItem[propertyId] || []).map(inst => inst['bk_inst_name']).join(',')
    } else if (['date', 'time'].includes(propertyType)) {
        flatternedItem[propertyId] = formatTime(flatternedItem[propertyId], propertyType === 'date' ? 'YYYY-MM-DD' : 'YYYY-MM-DD HH:mm:ss')
    }
    return flatternedItem
}

/**
 * 拍平主机列表
 * @param {Object} properties - 模型属性,eg: {host: [], biz: []}
 * @param {Array} list - 模型实例列表
 * @return {Array} 拍平后的模型实例列表
 */
export function flatternHostList (properties, list) {
    if (!list.length) return list
    const flatternedList = clone(list)
    const objIds = Object.keys(properties)
    objIds.forEach(objId => {
        properties[objId].forEach(property => {
            flatternedList.forEach((item, index) => {
                flatternedList[index] = flatternHostItem(property, item)
            })
        })
    })
    return flatternedList
}

/**
 * 拍平主机实例具体的属性
 * @param {Object} property - 模型具体属性
 * @param {Object} item - 模型实例
 * @return {Object} 拍平后的模型实例
 */
export function flatternHostItem (property, item) {
    const flatternedItem = clone(item)
    const objId = property['bk_obj_id']
    const propertyId = property['bk_property_id']
    const propertyType = property['bk_property_type']
    if (propertyType === 'enum') {
        if (flatternedItem[objId] instanceof Array) {
            flatternedItem[objId].forEach(subItem => {
                const enumOption = (property.option || []).find(option => option.id === subItem[propertyId])
                subItem[propertyId] = enumOption ? enumOption.name : null
            })
        } else {
            const enumOption = (property.option || []).find(option => option.id === flatternedItem[objId][propertyId])
            flatternedItem[objId][propertyId] = enumOption ? enumOption.name : null
        }
    } else if (['singleasst', 'multiasst'].includes(propertyType)) {
        if (flatternedItem[objId] instanceof Array) {
            flatternedItem[objId].forEach(subItem => {
                subItem[propertyId] = (subItem[propertyId] || []).map(inst => inst['bk_inst_name']).join(',')
            })
        } else {
            flatternedItem[objId][propertyId] = (flatternedItem[objId][propertyId] || []).map(inst => inst['bk_inst_name']).join(',')
        }
    } else if (['date', 'time'].includes(propertyType)) {
        if (flatternedItem[objId] instanceof Array) {
            flatternedItem[objId].forEach(subItem => {
                subItem[propertyId] = formatTime(subItem[propertyId], propertyType === 'date' ? 'YYYY-MM-DD' : 'YYYY-MM-DD HH:mm:ss')
            })
        } else {
            flatternedItem[objId][propertyId] = formatTime(flatternedItem[objId][propertyId], propertyType === 'date' ? 'YYYY-MM-DD' : 'YYYY-MM-DD HH:mm:ss')
        }
    }
    return flatternedItem
}

/**
 * 获取主机表格展示的文本
 * @param {Object} inst - 主机实例
 * @param {String} objId - 主机属性模型ID
 * @param {String} propertyId - 属性ID
 * @return {Object} 实例真实值
 */
export function getHostCellText (inst, objId, propertyId) {
    const valueObj = inst[objId]
    const values = []
    if (valueObj instanceof Array) {
        valueObj.forEach(value => {
            if (!['', null].includes(value[propertyId])) {
                values.push(value[propertyId])
            }
        })
    } else {
        values.push(valueObj[propertyId])
    }
    return values.join(',') || '--'
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
        if (['singleasst', 'multiasst'].includes(propertyType)) {
            const validAsst = (inst[propertyId] || []).filter(asstInst => asstInst.id !== '')
            values[propertyId] = validAsst.map(asstInst => asstInst['bk_inst_id']).join(',')
        } else if (['date', 'time'].includes(propertyType)) {
            values[propertyId] = formatTime(inst[propertyId], propertyType === 'date' ? 'YYYY-MM-DD' : 'YYYY-MM-DD HH:mm:ss')
        } else if (['int'].includes(propertyType)) {
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
    let formatedTime = moment(originalTime).format(format)
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

export default {
    getProperty,
    getPropertyPriority,
    getEnumOptions,
    getDefaultHeaderProperties,
    getCustomHeaderProperties,
    getHeaderProperties,
    flatternList,
    flatternItem,
    flatternHostList,
    flatternHostItem,
    getHostCellText,
    formatTime,
    clone,
    getInstFormValues
}
