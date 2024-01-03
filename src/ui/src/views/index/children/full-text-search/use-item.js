/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { computed } from 'vue'
import store from '@/store'
import { t } from '@/i18n'
import routerActions from '@/router/actions'
import { $warn } from '@/magicbox/index.js'
import {
  MENU_RESOURCE_INSTANCE_DETAILS,
  MENU_RESOURCE_BUSINESS_DETAILS,
  MENU_RESOURCE_BUSINESS_SET_DETAILS,
  MENU_RESOURCE_HOST_DETAILS,
  MENU_RESOURCE_BUSINESS_HISTORY,
  MENU_MODEL_DETAILS,
  MENU_BUSINESS_HOST_AND_SERVICE,
  MENU_RESOURCE_PROJECT_DETAILS
} from '@/dictionary/menu-symbol'
import { BUILTIN_MODELS, BUILTIN_MODEL_PROPERTY_KEYS, BUILTIN_MODEL_ROUTEPARAMS_KEYS } from '@/dictionary/model-constants'
import { getPropertyText } from '@/utils/tools'
import { escapeRegexChar } from '@/utils/util'

export default function useItem(list) {
  const getModelById = store.getters['objectModelClassify/getModelById']
  const getModelName = (source) => {
    const model = getModelById(source.bk_obj_id) || {}
    return model.bk_obj_name || ''
  }

  const normalizationList = computed(() => {
    const normalizationList = []
    list.value.forEach((item) => {
      const { key, kind, source } = item
      const newItem = { ...item }
      if (kind === 'instance' && key === BUILTIN_MODELS.HOST) {
        const ip = source.bk_host_innerip || source.bk_host_innerip_v6
        newItem.type = key
        newItem.title = Array.isArray(ip) ? ip.join(',') : ip
        newItem.typeName = t('主机')
        newItem.linkTo = handleGoResourceHost
      } else if (kind === 'instance' && key === BUILTIN_MODELS.BUSINESS) {
        newItem.type = key
        newItem.title = source.bk_biz_name
        newItem.typeName = t('业务')
        newItem.linkTo = handleGoBusiness
      } else if (kind === 'instance' && key === BUILTIN_MODELS.PROJECT) {
        newItem.type = key
        newItem.title = source[BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.PROJECT].NAME]
        newItem.typeName = t('项目')
        newItem.comp = 'project'
        newItem.linkTo = handleGoProject
      } else if (kind === 'instance' && key === BUILTIN_MODELS.BUSINESS_SET) {
        newItem.type = key
        newItem.title = source[BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS_SET].NAME]
        newItem.typeName = t('业务集')
        newItem.comp = 'bizset'
        newItem.linkTo = handleGoBusinessSet
      } else if (kind === 'instance' && key === BUILTIN_MODELS.SET) {
        newItem.type = key
        newItem.title = source.bk_set_name
        newItem.typeName = t('集群')
        newItem.linkTo = source => handleGoTopo('set', source)
      } else if (kind === 'instance' && key === BUILTIN_MODELS.MODULE) {
        newItem.type = key
        newItem.title = source.bk_module_name
        newItem.typeName = t('模块')
        newItem.linkTo = source => handleGoTopo(BUILTIN_MODELS.MODULE, source)
      } else if (kind === 'instance') {
        newItem.type = kind
        newItem.title = source.bk_inst_name
        newItem.typeName = getModelName(source)
        newItem.linkTo = handleGoInstace
      } else if (kind === 'model') {
        newItem.type = kind
        newItem.title = source.bk_obj_name
        newItem.typeName = t('模型')
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
      query: {
        from: 'resource'
      },
      history: true
    }

    if (newTab) {
      routerActions.open(to)
      return
    }

    routerActions.redirect(to)
  }
  const handleGoInstace = (source, newTab = true) => {
    const isPauserd = getModelById(source.bk_obj_id).bk_ispaused

    if (isPauserd) {
      $warn(t('该模型已停用'))
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
      routerActions.open(to)
      return
    }

    routerActions.redirect(to)
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
      routerActions.open(to)
      return
    }

    routerActions.redirect(to)
  }
  const handleGoProject = (source, newTab = true) => {
    const paramKey = BUILTIN_MODEL_ROUTEPARAMS_KEYS[BUILTIN_MODELS.PROJECT]
    const paramVal = source[BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.PROJECT].ID]
    const to = {
      name: MENU_RESOURCE_PROJECT_DETAILS,
      params: { [paramKey]: paramVal },
      history: true
    }

    if (newTab) {
      routerActions.open(to)
      return
    }

    routerActions.redirect(to)
  }
  const handleGoBusinessSet = (source, newTab = true) => {
    const paramKey = BUILTIN_MODEL_ROUTEPARAMS_KEYS[BUILTIN_MODELS.BUSINESS_SET]
    const paramVal = source[BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS_SET].ID]
    const to = {
      name: MENU_RESOURCE_BUSINESS_SET_DETAILS,
      params: { [paramKey]: paramVal },
      history: true
    }

    if (newTab) {
      routerActions.open(to)
      return
    }

    routerActions.redirect(to)
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
      routerActions.open(to)
      return
    }

    routerActions.redirect()
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
      routerActions.open(to)
      return
    }

    routerActions.redirect(to)
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

    const re = new RegExp(escapeRegexChar(words[1]))
    if (re.test(value)) {
      matched = keyword
      break
    }
  }

  return matched
}
