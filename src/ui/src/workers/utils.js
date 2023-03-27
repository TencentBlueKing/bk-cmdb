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

import customHeaders from '@/api/custom-header.js'

export const postData = async (url, data = {}, config = {}) => {
  const finalConfig = {
    originalResponse: false,
    transformData: true,
    ...config
  }

  const response = await fetch(url, {
    method: 'POST',
    mode: 'cors',
    credentials: 'include',
    headers: customHeaders,
    body: JSON.stringify(data)
  })

  if (!response.ok) {
    throw new Error('Network response was not OK')
  }

  if (finalConfig.originalResponse) {
    return response
  }

  const result = await response.json()

  if (finalConfig.transformData) {
    return result?.data
  }

  return result
}
