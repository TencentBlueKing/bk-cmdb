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

import store from '@/store'
import cloneDeep from 'lodash/cloneDeep'
import { getValue } from '@/utils/tools'

import propertyService from '@/service/property/property.js'
import propertyGroupService from '@/service/property/group.js'
import setTemplateService from '@/service/set-template/index.js'

export const templateDetailRequestId = Symbol()

export default async function userTemplagteData(bizId, templateId, isFetchTemplate) {
  const dataReqs = [
    // 集群属性与分组
    propertyService.find({ bk_obj_id: 'set', bk_biz_id: bizId }),
    propertyGroupService.find({ bk_obj_id: 'set', bk_biz_id: bizId }),
  ]

  if (isFetchTemplate) {
    dataReqs.push(setTemplateService.getFullOne(
      { bk_biz_id: bizId, id: templateId },
      { requestId: templateDetailRequestId }
    ))
  }

  const templateState = {
    templateName: '',
    configProperties: [],
    propertyConfig: {},
    formDataCopy: {} // 需要纳入表单填写检测的数据拷贝
  }

  const [
    setProperties,
    setPropertyGroup,
    templateData
  ] = await Promise.all(dataReqs)

  // 编辑态必要的数据初始化
  if (isFetchTemplate) {
    templateState.templateName = templateData.name

    // 属性设置
    const propertyConfigList = templateData.attributes || []
    propertyConfigList.forEach((item) => {
      const property = setProperties.find(prop => prop.id === item.bk_attribute_id)

      if (property) {
        // 已配置属性列表
        templateState.configProperties.push(property)

        // 属性配置值map
        templateState.propertyConfig[property.id] = item.bk_property_value
      }
    })

    templateState.formDataCopy = cloneDeep({
      templateName: templateState.templateName,
      propertyConfig: templateState.propertyConfig
    })
  }

  return {
    setProperties,
    setPropertyGroup,
    ...templateState
  }
}

export const getTemplateSyncStatus = async (bizId, templateId) => {
  try {
    const data = await store.dispatch('setTemplate/getSetTemplateStatus', {
      bizId,
      params: {
        set_template_ids: [templateId]
      },
      config: {
        cancelPrevious: true
      }
    })
    const needSync = getValue(data, '0.need_sync')

    return needSync
  } catch (error) {
    console.error(error)
  }
}
