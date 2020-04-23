import { MENU_BUSINESS } from '@/dictionary/menu-symbol'
export const requestId = Symbol('getAuthorizedBusiness')
export const before = async function (app, to, from, next) {
    const toTopRoute = to.matched[0]
    const fromTopRoute = from.matched[0]
    if (!toTopRoute || toTopRoute.name !== MENU_BUSINESS) {
        app.$http.cache.delete(requestId)
        if (fromTopRoute && fromTopRoute.name === MENU_BUSINESS) {
            fromTopRoute.meta.view = 'default'
        }
        return
    }
    if (fromTopRoute && fromTopRoute.name === MENU_BUSINESS && parseInt(to.params.bizId) !== parseInt(from.params.bizId)) {
        window.location.hash = to.fullPath
        return window.location.reload()
    }
    if (toTopRoute.meta.view === 'permission') {
        return next()
    }
    const authorizedList = await app.$store.dispatch('objectBiz/getAuthorizedBusiness', {
        requestId: requestId,
        fromCache: true
    })
    app.$store.commit('objectBiz/setAuthorizedBusiness', authorizedList)
    
    const id = parseInt(to.params.bizId || window.localStorage.getItem('selectedBusiness'))
    const business = authorizedList.find(business => business.bk_biz_id === id)
    if (business) {
        const isSubRoute = to.matched.length > 1
        toTopRoute.meta.view = 'default'
        window.localStorage.setItem('selectedBusiness', id)
        app.$store.commit('objectBiz/setBizId', id)
        return !isSubRoute && next(`/business/${id}/index`)
    }
    const hasURLId = to.params.bizId
    if (hasURLId) {
        toTopRoute.meta.view = 'permission'
        return next()
    }
    if (authorizedList.length) {
        toTopRoute.meta.view = 'default'
        const defaultId = authorizedList[0].bk_biz_id
        window.localStorage.setItem('selectedBusiness', defaultId)
        app.$store.commit('objectBiz/setBizId', id)
        return next(`/business/${defaultId}/index`)
    }
    toTopRoute.meta.view = 'permission'
    return next('/business')
}
