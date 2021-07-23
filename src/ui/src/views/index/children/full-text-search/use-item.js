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
        newItem.linkTo = handleGoTopo
      } else if (kind === 'instance' && key === 'module') {
        newItem.type = key
        newItem.title = source.bk_module_name
        newItem.typeName = root.$t('模块')
        newItem.linkTo = handleGoTopo
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

  const handleGoResourceHost = (host) => {
    root.$routerActions.redirect({
      name: MENU_RESOURCE_HOST_DETAILS,
      params: {
        id: host.bk_host_id
      },
      history: true
    })
  }
  const handleGoInstace = (source) => {
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
    root.$routerActions.redirect({
      name: MENU_RESOURCE_INSTANCE_DETAILS,
      params: {
        objId: source.bk_obj_id,
        instId: source.bk_inst_id
      },
      history: true
    })
  }
  const handleGoBusiness = (source) => {
    const name = source.bk_data_status === 'disabled' ? MENU_RESOURCE_BUSINESS_HISTORY : MENU_RESOURCE_BUSINESS_DETAILS
    root.$routerActions.redirect({
      name,
      params: {
        bizId: source.bk_biz_id,
        bizName: source.bk_biz_name
      },
      history: true
    })
  }
  const handleGoModel = (model) => {
    root.$routerActions.redirect({
      name: MENU_MODEL_DETAILS,
      params: {
        modelId: model.bk_obj_id
      },
      history: true
    })
  }
  const handleGoTopo = (data) => {
    const nodeMap = {
      set: `${data.key}-${data.source.bk_set_id}`,
      module: `${data.key}-${data.source.bk_module_id}`,
    }
    root.$routerActions.redirect({
      name: MENU_BUSINESS_HOST_AND_SERVICE,
      params: {
        bizId: data.source.bk_biz_id
      },
      query: {
        node: nodeMap[data.key]
      },
      history: true
    })
  }

  return {
    normalizationList
  }
}

export const getText = (property, data, thisProperty) => {
  let propertyValue = getPropertyText(property, data.source)

  // if (!Object.keys(data?.highlight).includes(thisProperty)) {
  return propertyValue || '--'
  // }

  // 对highlight属性值做高亮标签处理
  propertyValue = getHighlightValue(propertyValue, data, thisProperty)
  return propertyValue || '--'
}

export const getHighlightValue = (value, data, thisProperty) => {
  const highlightValue = data?.highlight?.[thisProperty]
  if (!highlightValue) {
    return value
  }
  // eslint-disable-next-line prefer-destructuring
  let keyword = Array.isArray(highlightValue) ? highlightValue[0] : highlightValue
  // eslint-disable-next-line prefer-destructuring
  keyword = keyword.match(/<em>(.+?)<\/em>/)[1]
  const reg = new RegExp(`(${keyword})`, 'g')
  return String(value).replace(reg, '<em class="hl">$1</em>')
}
