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
  <bk-select class="audit-target-selector"
    v-bind="$attrs"
    v-model="localValue">
    <bk-option
      v-for="option in options"
      :key="option.id"
      :id="option.id"
      :name="option.name">
    </bk-option>
  </bk-select>
</template>

<script>
  export default {
    props: {
      value: {
        type: String,
        default: ''
      },
      category: {
        type: String,
        default: 'business',
        validator(category) {
          return ['host', 'business', 'resource', 'other'].includes(category)
        }
      }
    },
    data() {
      return {
        dictionary: [],
        targetMap: Object.freeze({
          host: new Set(['host']),
          business: new Set([
            'dynamic_group',
            'set_template',
            'service_template',
            'service_category',
            'module',
            'set',
            'mainline_instance',
            'service_instance',
            'process',
            'service_instance_label',
            'host_apply',
            'custom_field'
          ]),
          resource: new Set([
            'business',
            'biz_set',
            'model_instance',
            'instance_association',
            'resource_directory',
            'cloud_area',
            'cloud_account',
            'cloud_sync_task',
          ]),
          other: new Set([
            'model_group',
            'model',
            'model_attribute',
            'model_unique',
            'model_association',
            'model_attribute_group',
            'event',
            'association_kind',
            'platform_setting'
          ])
        })
      }
    },
    computed: {
      options() {
        return this.dictionary.filter(target => this.targetMap[this.category].has(target.id))
      },
      localValue: {
        get() {
          return this.value
        },
        set(value) {
          this.$emit('input', value)
          this.$emit('change', value)
        }
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
