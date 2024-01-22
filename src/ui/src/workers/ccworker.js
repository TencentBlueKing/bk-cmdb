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

export default class CCWorker {
  constructor(worker, receivers = {}) {
    this.worker = worker

    // 传递主线程中的配置，所有基于此生产的worker都可以在内部先获取配置
    this.dispatch('config', {
      API_PREFIX: window.API_PREFIX,
      IAM_ENABLED: window.Site.authscheme === 'iam'
    })

    // 接收worker中返回的消息，根据type执行对应的receivers处理函数
    this.worker.onmessage = ({ data: { type, payload, error } }) => {
      if (!receivers[type]) {
        console.warn('No receivers handler configured:', type)
        return
      }

      receivers?.[type]?.call(this, type === 'error' ? error : payload)
    }

    return this
  }

  dispatch(type, payload) {
    this.worker.postMessage({ type, payload })
  }

  close() {
    this.worker.terminate()
  }
}
