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

const storageKey = 'selectedBusinessSet'

// 上次是否使用（访问）了业务集
let recentlyUsed = false

export const setBizSetIdToStorage = (id) => {
  window.localStorage.setItem(storageKey, id)
}

export const getBizSetIdFromStorage = () => window.localStorage.getItem(storageKey)

export const setBizSetRecentlyUsed = used => recentlyUsed = used

export const getBizSetRecentlyUsed = () => recentlyUsed
