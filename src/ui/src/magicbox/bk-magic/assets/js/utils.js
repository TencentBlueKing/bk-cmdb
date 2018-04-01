/**
 *  在指定的对象中根据路径找到value
 *  @param obj {Object} - 指定的对象
 *  @param keys {String} - 路径，支持a.b.c的形式
 */
export function findValueInObj (obj, keys = '') {
    let keyArr = keys.split('.')
    let o = obj
    let result = null
    let length = keyArr.length

    for (let [index, key] of keyArr.entries()) {
        if (!o) break

        if (index === length - 1) {
            result = o[key]
            break
        }

        o = o[key]
    }

    return result
}

function convertToNum (el) {
    return Number(el) || el
}

/**
 *  在项为对象的数组中，根据某一组{ key: value }，找到其所在对象中指定key的value
 *  @param arr {Array} - 指定的数组
 *  @param originItem {Object} - 依据的一组{ key: value }
 *  @param targetKey {String} - 指定的key值
 */
export function findValueInArrByRecord (arr, originItem, targetKey) {
    let result
    let key = Object.keys(originItem)[0]
    let item = originItem[key]

    for (let [index, _arr] of arr.entries()) {
        if (convertToNum(_arr[key]) === convertToNum(item)) {
            result = _arr[targetKey]
            break
        }
    }

    return result
}

export function isObject (argv) {
    return Object.prototype.toString.call(argv).toLowerCase() === '[object object]'
}

/**
 *  判断某个值是否在指定的数组中
 *  @param arr {Array} - 指定的数组
 *  @param target {String/Number/Object} - 指定的值
 *  @return result {Object} - 返回的对象，若为找到，该对象中key为result的值是false；若找到，key为result的值为true，
 *              同时key为index的值为当前项在数组中的索引值
 */
export function isInArray (arr, target) {
    let result = {
        result: false
    }
    let stringify = JSON.stringify

    if (!arr.length) return result

    for (let [index, item] of arr.entries()) {
        if (item) {
            let $item = isObject(item) ? stringify(item) : item.toString()
            let $target = isObject(target) ? stringify(target) : target.toString()

            if ($item === $target) {
                result = {
                    result: true,
                    index: index
                }
                break
            }
        }
    }

    return result
}

export function isVNode (node) {
    return typeof node === 'object' && node.hasOwnProperty('componentOptions')
}
