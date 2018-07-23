import sortUnicode from '@/common/js/sortUnicode'
/**
 * 获取主机关联关系
 * @param data - 主机列表信息
 * @return list - 列表 ['XXX业务 > xxx集群 > XXX模块']
 */
export function getHostRelation (data) {
    let list = []
    data.module.map(module => {
        let set = data.set.find(set => {
            return set['bk_set_id'] === module['bk_set_id']
        })
        if (set) {
            let biz = data.biz.find(biz => {
                return biz['bk_biz_id'] === set['bk_biz_id']
            })
            if (biz) {
                list.push(`${biz['bk_biz_name']} > ${set['bk_set_name']} > ${module['bk_module_name']}`)
            }
        }
    })
    return list
}
const HEX_TO_RGB = (hex) => {
    let rgb = []
    hex = hex.substr(1)
    if (hex.length === 3) {
        hex = hex.replace(/(.)/g, '$1$1')
    }
    hex.replace(/../g, function (color) {
        rgb.push(parseInt(color, 0x10))
    })
    return rgb
}
const GET_FILE_EXTENSION = (fileName) => {
    return fileName.substr((~-fileName.lastIndexOf('.') >>> 0) + 2)
}
const GET_BASE_64_IMAGE = (image, color) => {
    let canvas = document.createElement('canvas')
    let ctx = canvas.getContext('2d')
    ctx.clearRect(0, 0, canvas.width, canvas.height)
    canvas.width = image.width
    canvas.height = image.height
    ctx.drawImage(image, 0, 0, image.width, image.height)
    const imageData = ctx.getImageData(0, 0, image.width, image.height)
    const rgbColor = HEX_TO_RGB(color)
    for (let i = 0; i < imageData.data.length; i += 4) {
        imageData.data[i] = rgbColor[0]
        imageData.data[i + 1] = rgbColor[1]
        imageData.data[i + 2] = rgbColor[2]
    }
    ctx.putImageData(imageData, 0, 0)
    return canvas.toDataURL(`image/${GET_FILE_EXTENSION(image.src)}`)
}

export function generateObjIcon (image, options) {
    if (image instanceof Image) {
        const base64Image = GET_BASE_64_IMAGE(image, options.iconColor)
        return `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="100" height="100">
                    <circle cx="50" cy="50" r="49" fill="${options.backgroundColor}"/>
                    <svg xmlns="http://www.w3.org/2000/svg" stroke="rgba(0, 0, 0, 0)" viewBox="0 0 18 18" x="35" y="-12" fill="${options.iconColor}" width="35" >
                        <image width="15" xlink:href="${base64Image}"></image>
                    </svg>
                    <foreignObject x="0" y="58" width="100%" height="100%">
                        <div xmlns="http://www.w3.org/1999/xhtml" style="font-size:14px">
                            <div style="color:${options.fontColor};text-align: center;width: 60px;overflow:hidden;white-space:nowrap;text-overflow:ellipsis;margin:0 auto">${options.name}</div>
                        </div>
                    </foreignObject>
                </svg>`
    } else {
        options = image
        return `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="100" height="100">
                    <circle cx="50" cy="50" r="49" fill="${options.backgroundColor}"/>
                    <foreignObject x="0" y="43" width="100%" height="100%">
                        <div xmlns="http://www.w3.org/1999/xhtml" style="font-size:14px">
                            <div style="color:${options.fontColor};text-align: center;width: 60px;overflow:hidden;white-space:nowrap;text-overflow:ellipsis;margin:0 auto">${options.name}</div>
                        </div>
                    </foreignObject>
                </svg>`
    }
}

/**
 * 加载图片
 * @param {String} src - 图片路径
 * @param {Function} successFunc - 成功回调
 * @param {Function} failFunc - 失败回调
 */
export function loadImage (src, successFunc, failFunc) {
    const image = new Image()
    image.onload = () => {
        if ('naturalHeight' in image) {
            if (image.naturalHeight + image.naturalWidth === 0) {
                image.onerror()
                return
            }
        } else if (image.width + image.height === 0) {
            image.onerror()
            return
        }
        successFunc(image)
    }
    image.onerror = () => {
        failFunc(image)
    }
    image.src = src
}

/**
 * 封装加载图片方法
 * @param {String} url - 图片路径
 */
export function getImgUrl (url) {
    return new Promise((resolve, reject) => {
        loadImage(url, img => resolve(img), e => reject(e))
    })
}

/**
 * 数组按中英文首字母排序
 * @param {Array} array - 需要设置排序的数组
 * @param {String} field - 根据该字段内容进行排序 默认为null 只有数组中为对象时需要传该值
 * @param {Boolean} order - 正序为true 倒序为false 默认为true
 * @return {Array} 排序后的数组
 */
export function sortArray (array, field = null, order = true) {
    let arrayCopy = JSON.parse(JSON.stringify(array))
    arrayCopy.map(item => {
        let str = field === null ? item : item[field]
        let sortKey = ''
        for (let i = 0; i < str.length; i++) {
            let code = str.charCodeAt(i)
            if (code < 40869 && code >= 19968) { // 为中文
                sortKey += sortUnicode.charAt(code - 19968)
            } else {
                sortKey += str[i]
            }
        }
        item['_sortKey'] = sortKey
    })
    arrayCopy.sort((itemA, itemB) => order ? itemA['_sortKey'].localeCompare(itemB['_sortKey']) : itemB['_sortKey'].localeCompare(itemA['_sortKey']))
    arrayCopy.map(item => delete item['_sortKey'])
    return arrayCopy
}
