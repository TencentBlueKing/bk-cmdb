import Meta from '@/router/meta'
import { MENU_MODEL_MANAGEMENT, MENU_MODEL_DETAILS } from '@/dictionary/menu-symbol'

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
  name: MENU_MODEL_DETAILS,
  path: 'index/details/:modelId',
  component: () => import('./children/index.vue'),
  meta: new Meta({
    menu: {
      i18n: '模型详情'
    },
    layout: {},
    checkAvailable: (to, from, app) => {
      const { modelId } = to.params
      const model = app.$store.getters['objectModelClassify/getModelById'](modelId)
      return !!model
    }
  })
}]
