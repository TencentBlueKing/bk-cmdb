import { computed } from '@vue/composition-api'
import {
  MENU_RESOURCE_INSTANCE_DETAILS,
  MENU_RESOURCE_BUSINESS_DETAILS,
  MENU_RESOURCE_HOST_DETAILS,
  MENU_RESOURCE_BUSINESS_HISTORY,
  MENU_MODEL_DETAILS,
  MENU_BUSINESS_HOST_AND_SERVICE
} from '@/dictionary/menu-symbol'
import { getPropertyText } from '@/utils/tools'

export default function useItem(list, root) {
  const getModelById = root.$store.getters['objectModelClassify/getModelById']
  const getModelName = (source) => {
    const model = getModelById(source.bk_obj_id) || {}
    return model.bk_obj_name || ''
  }

  const normalizationList = computed(() => {
    const normalizationList = []
    list.value.forEach((item) => {
      const { key, kind, source } = item
      const newItem = { ...item }
      if (kind === 'instance' && key === 'host') {
        newItem.type = key
        newItem.title = Array.isArray(source.bk_host_innerip) ? source.bk_host_innerip.join(',') : source.bk_host_innerip
        newItem.typeName = root.$t('主机')
        newItem.linkTo = handleGoResourceHost
      } else if (kind === 'instance' && key === 'biz') {
        newItem.type = key
        newItem.title = source.bk_biz_name
        newItem.typeName = root.$t('业务')
        newItem.linkTo = handleGoBusiness
      } else if (kind === 'instance' && key === 'set') {
        newItem.type = key
        newItem.title = source.bk_set_name
        newItem.typeName = root.$t('集群')
        newItem.linkTo = source => handleGoTopo('set', source)
      } else if (kind === 'instance' && key === 'module') {
        newItem.type = key
        newItem.title = source.bk_module_name
        newItem.typeName = root.$t('模块')
        newItem.linkTo = source => handleGoTopo('module', source)
      } else if (kind === 'instance') {
        newItem.type = kind
        newItem.title = source.bk_inst_name
        newItem.typeName = getModelName(source)
        newItem.linkTo = handleGoInstace
      } else if (kind === 'model') {
        newItem.type = kind
        newItem.title = source.bk_obj_name
        newItem.typeName = root.$t('模型')
        newItem.linkTo = handleGoModel
      }
      normalizationList.push(newItem)
    })

    return normalizationList
  })

  const handleGoResourceHost = (host, newTab = true) => {
    const to = {
      name: MENU_RESOURCE_HOST_DETAILS,
      params: {
        id: host.bk_host_id
      },
      history: true
    }

    if (newTab) {
      root.$routerActions.open(to)
      return
    }

    root.$routerActions.redirect(to)
  }
  const handleGoInstace = (source, newTab = true) => {
    const model = getModelById(source.bk_obj_id)
    const isPauserd = getModelById(source.bk_obj_id).bk_ispaused
    if (model.bk_classification_id === 'bk_biz_topo') {
      root.$bkMessage({
        message: root.$t('主线模型无法查看'),
        theme: 'warning'
      })
      return
    } if (isPauserd) {
      root.$bkMessage({
        message: root.$t('该模型已停用'),
        theme: 'warning'
      })
      return
    }

    const to = {
      name: MENU_RESOURCE_INSTANCE_DETAILS,
      params: {
        objId: source.bk_obj_id,
        instId: source.bk_inst_id
      },
      history: true
    }

    if (newTab) {
      root.$routerActions.open(to)
      return
    }

    root.$routerActions.redirect(to)
  }
  const handleGoBusiness = (source, newTab = true) => {
    let to = {
      name: MENU_RESOURCE_BUSINESS_DETAILS,
      params: { bizId: source.bk_biz_id },
      history: true
    }

    if (source.bk_data_status === 'disabled') {
      to = {
        name: MENU_RESOURCE_BUSINESS_HISTORY,
        params: { bizName: source.bk_biz_name },
        history: true
      }
    }

    if (newTab) {
      root.$routerActions.open(to)
      return
    }

    root.$routerActions.redirect(to)
  }
  const handleGoModel = (model, newTab = true) => {
    const to = {
      name: MENU_MODEL_DETAILS,
      params: {
        modelId: model.bk_obj_id
      },
      history: true
    }

    if (newTab) {
      root.$routerActions.open(to)
      return
    }

    root.$routerActions.redirect()
  }
  const handleGoTopo = (key, source, newTab = true) => {
    const nodeMap = {
      set: `${key}-${source.bk_set_id}`,
      module: `${key}-${source.bk_module_id}`,
    }

    const to = {
      name: MENU_BUSINESS_HOST_AND_SERVICE,
      params: {
        bizId: source.bk_biz_id
      },
      query: {
        node: nodeMap[key]
      },
      history: true
    }

    if (newTab) {
      root.$routerActions.open(to)
      return
    }

    root.$routerActions.redirect(to)
  }

  return {
    normalizationList
  }
}

export const getText = (property, data) => {
  let propertyValue = getPropertyText(property, data.source)

  // 对highlight属性值做高亮标签处理
  propertyValue = getHighlightValue(propertyValue, data)
  return propertyValue || '--'
}

export const getHighlightValue = (value, data) => {
  const keywords = data?.highlight?.keywords
  if (!keywords || !keywords.length) {
    return value
  }

  // 用匹配到的高亮词（不一定等于搜索词）去匹配给定的值，如果命中则返回完整高亮词替代原本的值
  let matched = value
  // eslint-disable-next-line no-restricted-syntax
  for (const keyword of keywords) {
    const words = keyword.match(/<em>(.+?)<\/em>/)
    if (!words) {
      continue
    }

    const re = new RegExp(words[1])
    if (re.test(value)) {
      matched = keyword
      break
    }
  }

  return matched
}
