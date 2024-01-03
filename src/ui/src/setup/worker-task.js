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
import { OPERATION } from '@/dictionary/iam-auth'
import CCWorker from '@/workers/ccworker.js'
import { isViewAuthFreeModel, isViewAuthFreeModelInstance } from '@/service/auth'

const iamWorker = new CCWorker(new Worker(new URL('@/workers/iam.js', import.meta.url)), {
  preverify(payload) {
    store.commit('auth/setAuthedList', payload)
  },
  error(err) {
    console.error(err)
  }
})

const iam = () => {
  const models = store.getters['objectModelClassify/models']

  const authList = [
    { type: OPERATION.R_FULLTEXT_SEARCH },
    { type: OPERATION.R_RESOURCE_HOST },
    { type: OPERATION.R_MODEL_TOPOLOGY },
    { type: OPERATION.R_CLOUD_AREA },
    { type: OPERATION.R_PROJECT },
  ]

  models.filter(item => !item.bk_ishidden)
    .map(item => item.id)
    .forEach((modelId) => {
      // 需要预鉴权模型查看权限的模型
      if (!isViewAuthFreeModel({ id: modelId })) {
        authList.push({ type: OPERATION.R_MODEL, relation: [modelId] })
      }

      // 需要预鉴权模型实例查看权限的模型
      if (!isViewAuthFreeModelInstance({ id: modelId })) {
        authList.push({ type: OPERATION.R_INST, relation: [modelId] })
      }
    })

  iamWorker.dispatch('preverify', {
    authList
  })
}

const run = () => {
  iam()
}

export default {
  run,
  iam
}
