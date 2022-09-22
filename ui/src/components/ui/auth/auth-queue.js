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

import Vue from 'vue'
import debounce from 'lodash.debounce'
import $http from '@/api'
import { TRANSFORM_TO_INTERNAL } from '@/dictionary/iam-auth'

function equal(source, target) {
  const {
    resource_type: SResourceType,
    resource_id: SResourceId,
    resource_id_ex: SResourceIdEx,
    action: SAction,
    bk_biz_id: SBizId
  } = source
  const {
    resource_type: TResourceType,
    resource_id: TResourceId,
    resource_id_ex: TResourceIdEx,
    action: TAction,
    bk_biz_id: TBizId
  } = target
  const SParentLayers = source.parent_layers || []
  const TParentLayers = target.parent_layers || []
  if (
    SResourceType !== TResourceType
        || SResourceId !== TResourceId
        || SResourceIdEx !== TResourceIdEx
        || SAction !== TAction
        || SBizId !== TBizId
        || SParentLayers.length !== TParentLayers.length
  ) {
    return false
  }
  return SParentLayers.every((_, index) => {
    const SParentLayersMeta = SParentLayers[index]
    const TParentLayersMeta = TParentLayers[index]
    return Object.keys(SParentLayersMeta).every(key => SParentLayersMeta[key] === TParentLayersMeta[key])
  })
}

function unique(data) {
  return data.reduce((queue, meta) => {
    const exist = queue.some(exist => equal(exist, meta))
    if (!exist) {
      queue.push(meta)
    }
    return queue
  }, [])
}

export const AuthRequestId = Symbol('auth_request_id')

const authEnable = window.Site.authscheme === 'iam'
let afterVerifyQueue = []
export function afterVerify(func, once = true) {
  if (authEnable) {
    afterVerifyQueue.push({
      handler: func,
      once
    })
  } else {
    func()
  }
}
function execAfterVerify(authData) {
  afterVerifyQueue.forEach(({ handler }) => handler(authData))
  afterVerifyQueue = afterVerifyQueue.filter(({ once }) => !once)
}

export default new Vue({
  data() {
    return {
      queue: [],
      authComponents: [],
      verify: debounce(this.getAuth, 20)
    }
  },
  watch: {
    queue() {
      this.verify()
    }
  },
  methods: {
    add({ component, data }) {
      this.authComponents.push(component)
      // eslint-disable-next-line new-cap
      const authMetas = TRANSFORM_TO_INTERNAL(data)
      this.queue.push(...authMetas)
    },
    async getAuth() {
      if (!this.queue.length) return
      const queue = unique(this.queue.splice(0))
      const authComponents = this.authComponents.splice(0)
      let authData = []
      try {
        authData = await $http.post('auth/verify', { resources: queue }, { requestId: AuthRequestId })
      } catch (error) {
        console.error(error)
      } finally {
        authComponents.forEach((component) => {
          // eslint-disable-next-line new-cap
          const authMetas = TRANSFORM_TO_INTERNAL(component.auth)
          const authResults = []
          authMetas.forEach((meta) => {
            const result = authData.find((result) => {
              const source = {}
              const target = {}
              Object.keys(meta).forEach((key) => {
                source[key] = meta[key]
                target[key] = result[key]
              })
              return equal(source, target)
            })
            if (result) {
              authResults.push(result)
            }
          })
          component.updateAuth(Object.freeze(authResults), Object.freeze(authMetas))
        })
        execAfterVerify(authData)
      }
    }
  }
})
