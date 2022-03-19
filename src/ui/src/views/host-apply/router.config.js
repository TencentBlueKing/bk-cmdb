import Meta from '@/router/meta'
import {
  MENU_BUSINESS,
  MENU_BUSINESS_HOST,
  MENU_BUSINESS_HOST_APPLY,
  MENU_BUSINESS_HOST_APPLY_EDIT,
  MENU_BUSINESS_HOST_APPLY_CONFIRM,
  MENU_BUSINESS_HOST_APPLY_RUN,
  MENU_BUSINESS_HOST_APPLY_CONFLICT,
  MENU_BUSINESS_HOST_APPLY_FAILED
} from '@/dictionary/menu-symbol'
import { TASK_STATUS, getTask } from './task-helper.js'

// before拦截钩子，当本地存在未完成的应用任务则直接跳转到应用页面
const beforeInterceptor = async (to, from, app) => {
  const { $store, $router } = app
  const bizId = $store.getters['objectBiz/bizId']
  const localTaskId = getTask(`${$store.getters.userName}${bizId}`)
  if (localTaskId) {
    const result = await $store.dispatch('hostApply/getApplyTaskStatus', {
      params: {
        bk_biz_id: bizId,
        task_ids: [localTaskId]
      },
      config: {
        globalError: false
      }
    })
    const taskStatus = result?.task_info?.find(({ task_id: id }) => id === localTaskId)?.status
    if ([TASK_STATUS.NEW, TASK_STATUS.WAITING, TASK_STATUS.EXECUTING].includes(taskStatus)) {
      $router.push({
        name: MENU_BUSINESS_HOST_APPLY_RUN,
        params: { bizId }
      })
      return Promise.resolve(false)
    }
  }
  return Promise.resolve(true)
}

export default [{
  name: MENU_BUSINESS_HOST_APPLY,
  path: 'host-apply',
  component: () => import('./index.vue'),
  meta: new Meta({
    owner: MENU_BUSINESS,
    menu: {
      i18n: '主机属性自动应用',
      parent: MENU_BUSINESS_HOST
    },
    before: beforeInterceptor
  })
}, {
  name: MENU_BUSINESS_HOST_APPLY_CONFIRM,
  path: 'host-apply/confirm',
  component: () => import('./property-confirm'),
  meta: new Meta({
    owner: MENU_BUSINESS,
    menu: {
      i18n: '主机属性自动应用',
      parent: MENU_BUSINESS_HOST_APPLY,
      relative: MENU_BUSINESS_HOST_APPLY
    },
    before: beforeInterceptor
  })
}, {
  name: MENU_BUSINESS_HOST_APPLY_EDIT,
  path: 'host-apply/edit',
  component: () => import('./edit'),
  meta: new Meta({
    owner: MENU_BUSINESS,
    menu: {
      i18n: '主机属性自动应用',
      parent: MENU_BUSINESS_HOST_APPLY,
      relative: MENU_BUSINESS_HOST_APPLY
    },
    before: beforeInterceptor
  })
}, {
  name: MENU_BUSINESS_HOST_APPLY_RUN,
  path: 'host-apply/run',
  component: () => import('./property-run'),
  meta: new Meta({
    owner: MENU_BUSINESS,
    menu: {
      i18n: '主机属性自动应用',
      parent: MENU_BUSINESS_HOST_APPLY,
      relative: MENU_BUSINESS_HOST_APPLY
    },
    layout: {
      breadcrumbs: {
        back: false
      }
    }
  })
}, {
  name: MENU_BUSINESS_HOST_APPLY_CONFLICT,
  path: 'host-apply/conflict',
  component: () => import('./conflict-list'),
  meta: new Meta({
    owner: MENU_BUSINESS,
    menu: {
      i18n: '主机属性自动应用',
      parent: MENU_BUSINESS_HOST_APPLY,
      relative: MENU_BUSINESS_HOST_APPLY
    }
  })
}, {
  name: MENU_BUSINESS_HOST_APPLY_FAILED,
  path: 'host-apply/failed',
  component: () => import('./failed-list'),
  meta: new Meta({
    owner: MENU_BUSINESS,
    menu: {
      i18n: '主机属性自动应用',
      parent: MENU_BUSINESS_HOST_APPLY,
      relative: MENU_BUSINESS_HOST_APPLY
    }
  })
}]
