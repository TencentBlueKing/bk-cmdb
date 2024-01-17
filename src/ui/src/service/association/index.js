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

import http from '@/api'
export const find = async (params = {}) => {
  try {
    const result = await http.post('find/associationtype', params)
    return result
  } catch (error) {
    console.error(error)
    return { count: 0, info: [] }
  }
}

export const getAssociationCount = async ({ params = {}, config }) => {
  try {
    const result = await http.post('count/topoassociationtype', params, config)
    return result
  } catch (error) {
    console.error(error)
  }
}

export const getInstAssociation = async ({ params = {}, config }) => {
  try {
    const result = await http.post('find/instassociation', params, config)
    return result
  } catch (error) {
    console.error(error)
  }
}

export const getInstAssociationWithBiz = async ({ bizId, params = {}, config }) => {
  try {
    const result = await http.post(`find/instassociation/biz/${bizId}`, params, config)
    return result
  } catch (error) {
    console.error(error)
  }
}

export const findAll = () => find()

export default {
  find,
  findAll,
  getAssociationCount,
  getInstAssociation,
  getInstAssociationWithBiz
}
