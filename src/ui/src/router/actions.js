import router from './index'
import { Base64 } from 'js-base64'
import { MENU_BUSINESS } from '@/dictionary/menu-symbol'
export const redirect = function ({ name, params = {}, query = {}, history = false, reload = false }) {
    const queryBackup = { ...query }
    if (history) {
        const currentRoute = router.app.$route
        const data = {
            name: currentRoute.name,
            params: { ...currentRoute.params },
            query: { ...currentRoute.query }
        }
        const base64 = Base64.encode(JSON.stringify(data))
        queryBackup['_f'] = base64
    }
    const to = {
        name,
        params,
        query: queryBackup
    }
    if (reload) {
        const href = router.resolve(to).href
        window.location.href = href
        window.location.reload()
    } else {
        const resolved = router.resolve(to).resolved
        // 注入bizId，未改造的页面跳转，可能会遗漏了bizId的设置
        if (resolved.matched.length && resolved.matched[0].name === MENU_BUSINESS && !params.bizId) {
            to.params.bizId = router.app.$route.params.bizId
            console.warn('路由跳转未提供参数bizId, 已自动注入当前URL中的bizId')
        }
        router.replace(to)
    }
}

export const back = function () {
    const queryStr = router.app.$route.query._f
    if (queryStr) {
        try {
            const route = JSON.parse(Base64.decode(queryStr))
            redirect(route)
        } catch (error) {
            router.go(-1)
        }
    } else {
        router.go(-1)
    }
}

export const open = function (to) {
    const href = router.resolve(to).href
    window.open(href)
}

export default {
    redirect,
    back,
    open
}
