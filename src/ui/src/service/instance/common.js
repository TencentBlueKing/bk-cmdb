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

import hostSearchService from '@/service/host/search'
import businessSearchService from '@/service/business/search'
import instanceSearchService from '@/service/instance/search'
import businessSetService from '@/service/business-set/index.js'
import projectSetService from '@/service/project/index.js'
import { BUILTIN_MODELS, BUILTIN_MODEL_PROPERTY_KEYS, BUILTIN_MODEL_ROUTEPARAMS_KEYS } from '@/dictionary/model-constants.js'
import {
  MENU_RESOURCE_INSTANCE_DETAILS,
  MENU_RESOURCE_BUSINESS_DETAILS,
  MENU_RESOURCE_BUSINESS_SET_DETAILS,
  MENU_RESOURCE_HOST_DETAILS,
  MENU_RESOURCE_PROJECT_DETAILS,
  MENU_RESOURCE_BUSINESS_HISTORY
} from '@/dictionary/menu-symbol'
import { foreignkey as foreignkeyFormatter } from '@/filters/formatter.js'

const getService = modelId => ({
  [BUILTIN_MODELS.HOST]: hostSearchService,
  [BUILTIN_MODELS.BUSINESS]: businessSearchService,
  [BUILTIN_MODELS.BUSINESS_SET]: businessSetService,
  [BUILTIN_MODELS.PROJECT]: projectSetService
}[modelId] || instanceSearchService)

export const getIdKey = modelId => ({
  [BUILTIN_MODELS.HOST]: BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.HOST].ID,
  [BUILTIN_MODELS.BUSINESS]: BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS].ID,
  [BUILTIN_MODELS.BUSINESS_SET]: BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS_SET].ID,
  [BUILTIN_MODELS.PROJECT]: BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.PROJECT].ID
}[modelId] || 'bk_inst_id')

export const getNameKey = modelId => ({
  [BUILTIN_MODELS.HOST]: 'bk_host_innerip',
  [BUILTIN_MODELS.BUSINESS]: BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS].NAME,
  [BUILTIN_MODELS.BUSINESS_SET]: BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS_SET].NAME,
  [BUILTIN_MODELS.PROJECT]: BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.PROJECT].NAME
}[modelId] || 'bk_inst_name')

export const getSearchByNameParams = (modelId, value, options = {}) => {
  const nameKey = getNameKey(modelId)
  const idKey = getIdKey(modelId)
  const { fileds = [], page = {} } = options

  const pageParams = {
    start: 0,
    limit: 10,
    sort: idKey,
    ...page
  }

  const fieldParams = [idKey, nameKey, ...fileds]

  const generalParams = {
    fields: fieldParams,
    page: pageParams,
    conditions: value ? {
      condition: 'AND',
      rules: [{
        field: nameKey,
        operator: 'contains',
        value
      }]
    } : undefined
  }

  const params = {
    [BUILTIN_MODELS.HOST]: {
      ip: {
        data: value ? [value] : [],
        exact: 0,
        flag: 'bk_host_innerip|bk_host_outerip'
      },
      condition: [
        {
          bk_obj_id: BUILTIN_MODELS.HOST,
          fields: [...fieldParams, 'bk_cloud_id', 'bk_host_innerip_v6', 'bk_host_outerip_v6']
        }
      ],
      page: pageParams
    },
    [BUILTIN_MODELS.BUSINESS]: {
      fields: fieldParams,
      page: pageParams,
      condition: {
        bk_data_status: { $ne: 'disabled' },
        [nameKey]: value
      }
    },
    [BUILTIN_MODELS.BUSINESS_SET]: {
      fields: fieldParams,
      page: pageParams,
      bk_biz_set_filter: value ? {
        condition: 'AND',
        rules: [{
          field: nameKey,
          operator: 'contains',
          value
        }]
      } : undefined
    },
    [BUILTIN_MODELS.PROJECT]: {
      fields: fieldParams,
      page: pageParams,
      filter: value ? {
        condition: 'AND',
        rules: [{
          field: nameKey,
          operator: 'contains',
          value
        }]
      } : undefined
    }
  }

  return params[modelId] || generalParams
}

export const searchInstanceByName = async (modelId, name, options, config) => {
  const service = getService(modelId)
  const params = getSearchByNameParams(modelId, name, options)

  if (service.find) {
    let request = Promise.resolve([])
    if (modelId === BUILTIN_MODELS.BUSINESS_SET) {
      request = service.find(params, config)
    } else if (modelId === BUILTIN_MODELS.PROJECT) {
      request = service.getMany(params, config)
    } else {
      request = service.find({ bk_obj_id: modelId, params, config })
    }
    try {
      const result = await request
      return result
    } catch (error) {
      console.error(error)
      return Promise.reject(error)
    }
  }
  throw Error('not found model service find method')
}

export const searchInstanceByIds = async (modelId, instIds, config) => {
  const service = getService(modelId)

  if (service.find) {
    const request = (modelId === BUILTIN_MODELS.BUSINESS_SET || modelId === BUILTIN_MODELS.PROJECT)
      ? service.findByIds(instIds, config)
      : service.findByIds({ bk_obj_id: modelId, ids: instIds, config })
    const result = await request
    return result
  }
  throw Error('not found model service findByIds method')
}

const getIdValue = (modelId, idKey) => {
  let getId = item => item[idKey]
  if (modelId === BUILTIN_MODELS.HOST) {
    getId = ({ host }) => host[idKey]
  }
  return getId
}
const getNameValue = (modelId, nameKey) => {
  let getName = item => item[nameKey]
  if (modelId === BUILTIN_MODELS.HOST) {
    getName = ({ host }) => {
      const ip = host[nameKey] ?? host.bk_host_innerip_v6
      const cloudArea = foreignkeyFormatter(host.bk_cloud_id)
      return `${ip}(${cloudArea})`
    }
  }
  return getName
}

export const getModelInstanceOptions = async (modelId, instName, instIds, options, config) => {
  try {
    const results = await searchInstanceByName(modelId, instName, options, config)
    const list = results.list ?? results.info

    const nameKey = getNameKey(modelId)
    const idKey = getIdKey(modelId)

    const optionGeter = item => ({
      id: getIdValue(modelId, idKey)(item),
      name: getNameValue(modelId, nameKey)(item)
    })
    const instOptions = list.map(optionGeter)

    const idList = instOptions.map(options => options.id)

    // 没有传递name来搜索，同时传了ids，表示需要根据ids搜索并将结果补充至列表
    if (!instName && instIds?.length) {
      const diffIds = instIds.filter(id => !idList.includes(id))
      if (diffIds.length) {
        const { list: diffList } = await searchInstanceByIds(modelId, diffIds, config) || {}
        instOptions.push(...(diffList || []).map(optionGeter))
      }
    }

    return instOptions
  } catch (error) {
    console.error(error)
    return Promise.reject(error)
  }
}

export const getModelInstanceByIds = async (modelId, instIds, config) => {
  const { list } = await searchInstanceByIds(modelId, instIds, config) || {}
  const idKey = getIdKey(modelId)
  const nameKey = getNameKey(modelId)
  const newList = list.map(item => ({
    modelId,
    id: getIdValue(modelId, idKey)(item),
    name: getNameValue(modelId, nameKey)(item),
    inst: item
  }))

  return newList
}

export const getModelInstanceDetailRoute = (modelId, instId, extra = {}, options = {}) => {
  const { history = false } = options

  const generalRoute = (modelId, instId) => ({
    name: MENU_RESOURCE_INSTANCE_DETAILS,
    params: {
      objId: modelId,
      instId
    },
    history
  })

  const routes = {
    [BUILTIN_MODELS.HOST]: (modelId, instId) => ({
      name: MENU_RESOURCE_HOST_DETAILS,
      params: {
        id: instId
      },
      history
    }),
    [BUILTIN_MODELS.BUSINESS]: (modelId, instId, extra) => (extra.bk_data_status === 'disabled' ? {
      name: MENU_RESOURCE_BUSINESS_HISTORY,
      params: {
        bizName: extra.bk_biz_name
      },
      history
    } : {
      name: MENU_RESOURCE_BUSINESS_DETAILS,
      params: { bizId: instId },
      history
    }),
    [BUILTIN_MODELS.BUSINESS_SET]: (modelId, instId) => {
      const paramKey = BUILTIN_MODEL_ROUTEPARAMS_KEYS[BUILTIN_MODELS.BUSINESS_SET]

      return {
        name: MENU_RESOURCE_BUSINESS_SET_DETAILS,
        params: { [paramKey]: instId },
        history
      }
    },
    [BUILTIN_MODELS.PROJECT]: (modelId, instId) => {
      const paramKey = BUILTIN_MODEL_ROUTEPARAMS_KEYS[BUILTIN_MODELS.PROJECT]

      return {
        name: MENU_RESOURCE_PROJECT_DETAILS,
        params: { [paramKey]: instId },
        history
      }
    }
  }

  const routeGeter = routes[modelId] || generalRoute

  return routeGeter(modelId, instId, extra)
}
