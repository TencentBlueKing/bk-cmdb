import Meta from '@/router/meta'
import { MENU_MODEL_MANAGEMENT } from '@/dictionary/menu-symbol'

export default [{
    name: MENU_MODEL_MANAGEMENT,
    path: 'index',
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            i18n: '模型'
        }
    })
}, {
    name: 'modelDetails',
    path: 'index/details/:modelId',
    component: () => import('./children/index.vue'),
    meta: new Meta({
        menu: {
            i18n: '模型详情'
        },
        checkAvailable: (to, from, app) => {
            const modelId = to.params.modelId
            const model = app.$store.getters['objectModelClassify/getModelById'](modelId)
            return !!model
        }
    })
}]
