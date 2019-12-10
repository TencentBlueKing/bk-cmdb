import moment from 'moment'

const defaultFormatter = value => {
    if (!value) {
        return '--'
    }
    return value
}

const timeFormatter = (value, format = 'YYYY-MM-DD HH:mm:ss') => {
    if (!value) {
        return '--'
    }
    const formated = moment(value).format(format)
    if (formated === 'Invalid date') {
        return value
    }
    return formated
}

const numericFormatter = value => {
    if (isNaN(value) || value === null) {
        return '--'
    }
    return value
}

export function singlechar (value) {
    return defaultFormatter(value)
}

export function longchar (value) {
    return defaultFormatter(value)
}

export function int (value) {
    return numericFormatter(value)
}

export function float (value) {
    return numericFormatter(value)
}

export function date (value) {
    return timeFormatter(value, 'YYYY-MM-DD')
}

export function time (value) {
    return timeFormatter(value, 'YYYY-MM-DD HH:mm:ss')
}

export function objuser (value) {
    if (!value) {
        return '--'
    }
    const userList = window.CMDB_USER_LIST || []
    const user = userList.find(user => user.english_name === value)
    if (user) {
        return `${user.english_name}(${user.chinese_name})`
    }
    return value
}

export function timezone (value) {
    return defaultFormatter(value)
}

export function bool (value) {
    if (['true', 'false'].includes(value.toString())) {
        return value.toString()
    }
    return '--'
}

export function enumeration (value, options) {
    const option = (options || []).find(option => option.id === value)
    if (!option) {
        return '--'
    }
    return option.name
}

export function foreignkey (value) {
    if (Array.isArray(value)) {
        return value.map(inst => inst.bk_inst_name).join(',')
    }
    if (String(value).length) {
        return value
    }
    return '--'
}

const formatterMap = {
    singlechar,
    longchar,
    int,
    float,
    date,
    time,
    objuser,
    timezone,
    bool,
    foreignkey,
    enum: enumeration
}

export default function formatter (value, property, options) {
    const isPropertyObject = typeof property === 'object'
    const type = isPropertyObject ? property.bk_property_type : property
    const propertyOptions = isPropertyObject ? property.option : options
    if (formatterMap.hasOwnProperty(type)) {
        return formatterMap[type](value, propertyOptions)
    }
    return value
}
