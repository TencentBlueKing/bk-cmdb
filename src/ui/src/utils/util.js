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
