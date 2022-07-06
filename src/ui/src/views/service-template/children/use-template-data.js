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
import propertyService from '@/service/property/property.js'
import propertyGroupService from '@/service/property/group.js'
import serviceTemplateService from '@/services/service-template'

export const templateDetailRequestId = Symbol()

export default async (bizId, templateId, isFetchTemplate) => {
  const dataReqs = [
    // 模块属性与分组
    propertyService.find({ bk_obj_id: 'module', bk_biz_id: bizId }),
    propertyGroupService.find({ bk_obj_id: 'module', bk_biz_id: bizId }),

    // 进程属性与分组
    propertyService.find({ bk_obj_id: 'process', bk_biz_id: bizId }),
    propertyGroupService.find({ bk_obj_id: 'process', bk_biz_id: bizId }),

    // 服务分类
    store.dispatch('serviceClassification/searchServiceCategory', { params: { bk_biz_id: bizId } }),
  ]

  if (isFetchTemplate) {
    dataReqs.push(serviceTemplateService.getFullOne(
      { bk_biz_id: bizId, id: templateId },
      { requestId: templateDetailRequestId }
    ))
  }

  const templateState = {
    basic: {},
    configProperties: [],
    propertyConfig: {},
    processList: [],
    formDataCopy: {} // 需要纳入表单填写检测的数据拷贝
  }

  const [
    moduleProperties,
    modulePropertyGroup,
    processProperties,
    processPropertyGroup,
    { info: categories },
    templateData
  ] = await Promise.all(dataReqs)

  // 服务分类数据拆分为一二级便于使用
  const categoryList = categories.map(item => ({
    ...item.category,
    displayName: `${item.category.name}（#${item.category.id}）`
  }))

  const primaryCategories = categoryList.filter(category => !category.bk_parent_id)
  const secCategories = categoryList.filter(category => category.bk_parent_id)

  // 编辑态必要的数据初始化
  if (isFetchTemplate) {
    // 进程列表
    templateState.processList = templateData.processes.map(template => ({
      process_id: template.id,
      ...template.property
    })).sort((prev, next) => prev.process_id - next.process_id)

    // 模板表单基础数据
    const { id, service_category_id: categoryId, name } = templateData
    const category = secCategories.find(category => category.id === categoryId) || {}
    templateState.basic.id = id
    templateState.basic.templateName = name
    templateState.basic.primaryCategory = category.bk_parent_id
    templateState.basic.secCategory = categoryId

    // 属性设置
    const propertyConfigList = templateData.attributes || []
    propertyConfigList.forEach((item) => {
      const property = moduleProperties.find(prop => prop.id === item.bk_attribute_id)

      // 已配置属性列表
      templateState.configProperties.push(property)

      // 属性配置值键值对
      templateState.propertyConfig[property.id] = item.bk_property_value
    })

    templateState.formDataCopy = cloneDeep({
      templateName: templateState.basic.templateName,
      primaryCategory: templateState.basic.primaryCategory,
      secCategory: templateState.basic.secCategory,
      propertyConfig: templateState.propertyConfig,
      processList: templateState.processList
    })
  }

  return {
    moduleProperties,
    modulePropertyGroup,
    processProperties,
    processPropertyGroup,
    primaryCategories,
    secCategories,
    ...templateState
  }
}
