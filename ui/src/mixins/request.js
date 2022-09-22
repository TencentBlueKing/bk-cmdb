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

import { mapGetters } from 'vuex'
export default {
  computed: {
    ...mapGetters('request', {
      $requestQueue: 'queue',
      $requestCache: 'cache'
    })
  },
  methods: {
    $loading(requestIds) {
      if (typeof requestIds === 'undefined') {
        return !!this.$requestQueue.length
      } if (requestIds instanceof Array) {
        return requestIds.some(requestId => this.$requestQueue.some(request => request.requestId === requestId))
      } if (typeof requestIds === 'string' && requestIds.startsWith('^=')) {
        // eslint-disable-next-line prefer-destructuring
        const requestId = requestIds.split('^=')[1]
        const matchIndex = this.$requestQueue.findIndex(request => (typeof request.requestId === 'string') && request.requestId.startsWith(requestId))
        return matchIndex !== -1
      }
      return this.$requestQueue.some(request => request.requestId === requestIds)
    }
  }
}
