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

import autocompleteMixin from './autocomplete'
export default {
  mixins: [autocompleteMixin],
  methods: {
    async fuzzySearchMethod(keyword) {
      try {
        const result = await this.$http.get('biz/simplify?sort=bk_biz_id', {
          requestId: `fuzzy_search_${this.type}`,
          cancelPrevious: true
        })
        const list = result.info || []
        const matchRE = new RegExp(keyword, 'i')
        const matched = []
        list.forEach(({ bk_biz_name: name }) => {
          if (matchRE.test(name)) {
            matched.push({ text: name, value: name })
          }
        })
        return Promise.resolve({
          next: false,
          results: matched
        })
      } catch (error) {
        return Promise.reject(error)
      }
    }
  }
}
