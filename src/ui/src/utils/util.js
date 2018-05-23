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
