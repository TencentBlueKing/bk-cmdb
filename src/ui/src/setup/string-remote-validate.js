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
import Vue from 'vue'
import debounce from 'lodash.debounce'
const observer = new Vue({
  data() {
    return {
      queue: [],
      batchValidate: debounce(this.validate, 100)
    }
  },
  watch: {
    queue() {
      this.batchValidate()
    }
  },
  methods: {
    /**
     * { data: { regular: 'xxx', content: 'xxx' }, resolve: function }
     */
    add(item) {
      this.queue.push(item)
    },
    async validate() {
      if (!this.queue.length) return
      const queue = this.queue.splice(0)
      try {
        const results = await http.post(`${window.API_HOST}regular/verify_regular_content_batch`, {
          items: queue.map(({ data }) => data)
        }, {
          globalError: false
        })
        queue.forEach(({ resolve }, index) => resolve({ valid: results[index] }))
      } catch (error) {
        queue.forEach(({ resolve }) => resolve({ valid: false }))
      }
    }
  }
})
export default {
  validate: async (content, { regular }) => new Promise((resolve) => {
    observer.add({
      data: { content, regular },
      resolve
    })
  })
}
