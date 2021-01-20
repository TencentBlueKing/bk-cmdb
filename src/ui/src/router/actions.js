import router from './index'
import { Base64 } from 'js-base64'
import { MENU_BUSINESS } from '@/dictionary/menu-symbol'
export const redirect = function ({ name, params = {}, query = {}, history = false, reload = false, back = false }) {
    const queryBackup = { ...query }
    const currentRoute = router.app.$route

    // 当前页非history模式则先清空历史记录
    if (!currentRoute.query.hasOwnProperty('_f')) {
        window.sessionStorage.setItem('history', JSON.stringify([]))
    }

    // 先取得history列表
    let historyList = []
    try {
        historyList = JSON.parse(window.sessionStorage.getItem('history')) || []
        if (!Array.isArray(historyList)) {
            historyList = [historyList]
        }
    } catch (e) {
        historyList = []
    }

    if (history) {
        const data = {
            name: currentRoute.name,
            params: { ...currentRoute.params },
            query: { ...currentRoute.query }
        }
        const base64 = Base64.encode(JSON.stringify(data))
        queryBackup['_f'] = '1'

        historyList.push(base64)
        window.sessionStorage.setItem('history', JSON.stringify(historyList))
    } else if (back) {
        // 后退操作会注入back，此时从历史记录中删除当前记录
        try {
            const index = historyList.findIndex(item => {
                const history = JSON.parse(Base64.decode(item))
                return history.name === name
            })
            if (index !== -1) {
                historyList.splice(index, 1)
                window.sessionStorage.setItem('history', JSON.stringify(historyList))
            }
        } catch (e) {
            // ignore
        }
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
    let record
    if (router.app.$route.query.hasOwnProperty('_f')) {
        try {
            const historyList = JSON.parse(window.sessionStorage.getItem('history')) || []
            record = historyList.pop()
        } catch (e) {
            // ignore
        }
    }
    if (record) {
        try {
            const route = JSON.parse(Base64.decode(record))
            redirect({ ...route, back: true })
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
