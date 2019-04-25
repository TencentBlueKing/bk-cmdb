import { GET_MODEL_PATH } from '@/views/general-model/router.config'
import MENU, { NAV_COLLECT } from '@/dictionary/menu'
import { clone } from '@/utils/tools'
import { viewRouters } from '@/router'
const state = {
    active: null,
    open: null
}

const getters = {
    active: state => state.active,
    open: state => state.open,
    collectMenus: (state, getters, rootState, rootGetters) => {
        const collectMenus = []
        const usercustom = rootGetters['userCustom/usercustom']
        const collectedModelIds = usercustom.collected_models || []
        collectedModelIds.forEach((id, index) => {
            const model = rootGetters['objectModelClassify/models'].find(model => model.id === id)
            if (model) {
                collectMenus.push({
                    id: model.bk_obj_id,
                    name: model.bk_obj_name,
                    path: GET_MODEL_PATH(model.bk_obj_id),
                    order: index
                })
            }
        })
        return collectMenus
    },
    menus: (state, getters, rootState, rootGetters) => {
        const menus = clone(MENU)
        viewRouters.forEach(route => {
            const meta = route.meta || {}
            const auth = meta.auth || {}
            const menu = meta.menu || false
            const shouldShow = menu
                ? rootGetters.isAdminView
                    ? menu.adminView
                    : menu.businessView
                : true
            if (shouldShow) {
                const authorized = auth.view ? rootGetters['auth/isAuthorized'](auth.view, { type: 'view' }) : true
                if (authorized) {
                    if (menu.parent) {
                        const parent = menus.find(parent => parent.id === menu.parent) || {}
                        const submenu = parent.submenu || []
                        submenu.push(menu)
                    } else {
                        const parent = menus.find(parent => parent.id === menu.id) || {}
                        Object.assign(parent, menu)
                    }
                }
            }
        })
        const collectMenu = menus.find(menu => menu.id === NAV_COLLECT) || {}
        const collectSubmenu = collectMenu.submenu || []
        Array.prototype.push.apply(collectSubmenu, getters.collectMenus)
        const availableMenus = menus.filter(menu => {
            return menu.path
                || (Array.isArray(menu.submenu) && menu.submenu.length)
        })
        availableMenus.forEach(menu => {
            if (Array.isArray(menu.submenu)) {
                menu.submenu.sort((prev, next) => prev.order - next.order)
            }
        })
        return availableMenus
    }
}

const mutations = {
    setActiveMenu (state, menuId) {
        state.active = menuId
    },
    setOpenMenu (state, menuId) {
        state.open = menuId
    }
}

export default {
    namespaced: true,
    state,
    getters,
    mutations
}
