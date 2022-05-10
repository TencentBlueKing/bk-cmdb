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

import merge from 'lodash/merge'

/**
 * 根据是否开启count生成新参数
 * @param {Object} params 基础参数
 * @param {Boolean} flag 是否开启count获取
 * @returns 生成的新参数
 */
export const enableCount = (params = {}, flag = false) => {
  const page = Object.assign(flag ? { start: 0, limit: 0, sort: '' } : {}, { enable_count: flag })
  return merge({}, params, { page })
}

export const onePageParams = () => ({ start: 0, limit: 1 })
