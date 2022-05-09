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
  props: {
    bizId: {
      type: Number,
      required: true
    }
  },
  methods: {
    async fuzzySearchMethod(keyword) {
      try {
        const list = await this.$http.post('findmany/object/instances/names', {
          bk_obj_id: this.type,
          bk_biz_id: this.bizId,
          name: keyword
        }, {
          requestId: `fuzzy_search_${this.type}`,
          cancelPrevious: true
        })
        return Promise.resolve({
          next: false,
          results: list.map(name => ({ text: name, value: name }))
        })
      } catch (error) {
        return Promise.reject(error)
      }
    }
  }
}
