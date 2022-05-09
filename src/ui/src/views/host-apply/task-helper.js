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

import md5 from 'md5'

const STORAGE_KEY_PREFIX = 'applytask:'

const genKey = salt => `${STORAGE_KEY_PREFIX}${md5(salt)}`

export const setTask = (id, salt) => {
  const key = genKey(salt)
  localStorage.setItem(key, id)
}

export const getTask = (salt) => {
  const key = genKey(salt)
  return localStorage.getItem(key)
}

export const removeTask = (salt) => {
  const key = genKey(salt)
  localStorage.removeItem(key)
}

export const TASK_STATUS = {
  NEW: 'new',
  WAITING: 'waiting',
  EXECUTING: 'executing',
  FINISHED: 'finished',
  FAIL: 'failure'
}
