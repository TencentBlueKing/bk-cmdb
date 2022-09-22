<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<script>
  import FilterTag from './filter-tag.vue'
  import FilterTagItem from './general-model-filter-tag-item.vue'
  import { clearSearchQuery } from './general-model-filter.js'

  export default {
    components: {
      FilterTagItem
    },
    extends: FilterTag,
    provide() {
      return {
        condition: () => this.condition
      }
    },
    props: {
      filterSelected: {
        type: Array,
        default: () => ([])
      },
      filterCondition: {
        type: Object,
        default: () => ({})
      }
    },
    computed: {
      condition() {
        return this.filterCondition
      },
      showIPTag() {
        // 替换继续的值指定为false
        return false
      },
      selected() {
        return this.filterSelected.filter((property) => {
          const { value } = this.condition[property.id]
          return value !== null && value !== undefined && !!value.toString().length
        })
      }
    },
    methods: {
      handleResetAll() {
        clearSearchQuery()
      }
    }
  }
</script>
