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

<template>
  <bk-select
    v-bind="$attrs"
    v-model="localValue"
    display-tag
    selected-style="checkbox"
    multiple>
    <bk-option
      v-for="action in actions"
      :key="action.id"
      :id="action.id"
      :name="action.name">
    </bk-option>
  </bk-select>
</template>

<script>
  export default {
    props: {
      value: {
        type: Array,
        default: () => ([])
      },
      target: {
        type: String,
        default: ''
      }
    },
    data() {
      return {
        dictionary: []
      }
    },
    computed: {
      actions() {
        const target = this.dictionary.find(target => target.id === this.target)
        return target ? target.operations : []
      },
      localValue: {
        get() {
          return this.value
        },
        set(values) {
          this.$emit('input', values)
          this.$emit('change', values)
        }
      }
    },
    watch: {
      target() {
        this.localValue = []
      }
    },
    created() {
      this.getAuditDictionary()
    },
    methods: {
      async getAuditDictionary() {
        try {
          this.dictionary = await this.$store.dispatch('audit/getDictionary', {
            fromCache: true
          })
        } catch (error) {
          this.dictionary = []
        }
      }
    }
  }
</script>
