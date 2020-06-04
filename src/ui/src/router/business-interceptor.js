import { MENU_BUSINESS } from '@/dictionary/menu-symbol'

const requestId = Symbol('getAuthorizedBusiness')

let committed = false
export async function getAuthorizedBusiness (app) {
    const { info } = await app.$store.dispatch('objectBiz/getAuthorizedBusiness', {
        requestId: requestId,
        fromCache: true
    })
    if (!committed) {
        app.$store.commit('objectBiz/setAuthorizedBusiness', Object.freeze(info))
        committed = true
    }
    return info
}

export const before = async function (app, to, from, next) {
    const toTopRoute = to.matched[0]
    const fromTopRoute = from.matched[0]
    if (!toTopRoute || toTopRoute.name !== MENU_BUSINESS) {
        if (fromTopRoute && fromTopRoute.name === MENU_BUSINESS) {
            fromTopRoute.meta.view = 'default'
        }
        return true
    }
    if (fromTopRoute && fromTopRoute.name === MENU_BUSINESS && parseInt(to.params.bizId) !== parseInt(from.params.bizId)) {
        window.location.hash = to.fullPath
        window.location.reload()
        return false
    }
    if (toTopRoute.meta.view === 'permission') {
        next()
        return false
    }
    const authorizedList = await getAuthorizedBusiness(app)
    const id = parseInt(to.params.bizId || window.localStorage.getItem('selectedBusiness'))
    const business = authorizedList.find(business => business.bk_biz_id === id)
    const hasURLId = to.params.bizId

    // URL或者缓存中的id对应的业务存在
    if (business) {
        const isSubRoute = to.matched.length > 1
        toTopRoute.meta.view = 'default'
        window.localStorage.setItem('selectedBusiness', id)
        app.$store.commit('objectBiz/setBizId', id)

        if (!isSubRoute) { // 如果是一级路由，则重定向到带业务id的二级路由首页(业务拓扑)
            next({
                path: `/business/${id}/index`,
                replace: true
            })
            return false
        } else if (!hasURLId) { // 如果是二级路由且URL中不包含业务ID，则补充业务ID到URL中
            next({
                name: to.name,
                params: {
                    ...to.params,
                    bizId: id
                },
                query: to.query,
                replace: true
            })
            return false
        }
        return true // 正常的有权限的业务，且URL中带了ID，则直接返回，进行后续的路由逻辑
    }
    // 未找到对应有权限的业务，且URL中有业务ID，则显示一级view的无权限视图
    if (hasURLId) {
        toTopRoute.meta.view = 'permission'
        next()
        return false
    }
    // 缓存无ID，URL无ID，则认为是首次进入业务导航，取有权限业务的第一个写入URL中
    if (authorizedList.length) {
        const [firstBusiness] = authorizedList
        toTopRoute.meta.view = 'default'
        const defaultId = firstBusiness.bk_biz_id
        window.localStorage.setItem('selectedBusiness', defaultId)
        app.$store.commit('objectBiz/setBizId', defaultId)
        next({
            path: `/business/${defaultId}/index`,
            replace: true
        })
        return false
    }
    toTopRoute.meta.view = 'permission'
    next('/business')
    return false
}
