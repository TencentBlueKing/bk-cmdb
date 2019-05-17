/**
 * @file deepmerge
 * from https://github.com/KyleAMathews/deepmerge
 */

import defaultIsMergeableObject from './is-mergeable-object'

function emptyTarget (val) {
    return Array.isArray(val) ? [] : {}
}

function cloneUnlessOtherwiseSpecified (value, options) {
    return (options.clone !== false && options.isMergeableObject(value))
        ? deepmerge(emptyTarget(value), value, options)
        : value
}

function defaultArrayMerge (target, source, options) {
    return target.concat(source).map(element => cloneUnlessOtherwiseSpecified(element, options))
}

function mergeObject (target, source, options) {
    const destination = {}
    if (options.isMergeableObject(target)) {
        Object.keys(target).forEach(key => {
            destination[key] = cloneUnlessOtherwiseSpecified(target[key], options)
        })
    }
    Object.keys(source).forEach(key => {
        if (!options.isMergeableObject(source[key]) || !target[key]) {
            destination[key] = cloneUnlessOtherwiseSpecified(source[key], options)
        }
        else {
            destination[key] = deepmerge(target[key], source[key], options)
        }
    })

    return destination
}

function deepmerge (target, source, options) {
    options = options || {}
    options.arrayMerge = options.arrayMerge || defaultArrayMerge
    options.isMergeableObject = options.isMergeableObject || defaultIsMergeableObject

    const sourceIsArray = Array.isArray(source)
    const targetIsArray = Array.isArray(target)
    const sourceAndTargetTypesMatch = sourceIsArray === targetIsArray

    if (!sourceAndTargetTypesMatch) {
        return cloneUnlessOtherwiseSpecified(source, options)
    }
    else if (sourceIsArray) {
        return options.arrayMerge(target, source, options)
    }
    return mergeObject(target, source, options)
}

deepmerge.all = function deepmergeAll (array, options) {
    if (!Array.isArray(array)) {
        throw new Error('first argument should be an array')
    }

    return array.reduce((prev, next) => {
        return deepmerge(prev, next, options)
    }, {})
}

export default deepmerge
