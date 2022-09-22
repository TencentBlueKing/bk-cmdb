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
  import FilterTagForm from './filter-tag-form.vue'
  import { setSearchQueryByCondition } from './general-model-filter.js'

  export default {
    extends: FilterTagForm,
    props: {
      condition: {
        type: Object,
        default: () => ({})
      }
    },
    computed: {
      operator: {
        get() {
          return this.localOperator || this.condition[this.property.id].operator
        },
        set(operator) {
          this.localOperator = operator
        }
      },
      value: {
        get() {
          if (this.localValue === null) {
            return this.condition[this.property.id].value
          }
          return this.localValue
        },
        set(value) {
          this.localValue = value
        }
      }
    },
    methods: {
      handleConfirm() {
        // 构建单个condition
        const condition = {
          [this.property.id]: {
            operator: this.operator,
            value: this.value
          }
        }
        setSearchQueryByCondition(condition, [this.property])
        this.$emit('confirm')
      }
    }
  }
</script>
