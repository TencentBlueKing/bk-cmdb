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
    v-if="displayType === 'selector'"
    multiple
    searchable
    display-tag
    selected-style="checkbox"
    v-bind="$attrs"
    v-model="localValue"
    @clear="() => $emit('clear')"
    @toggle="handleToggle">
    <bk-option
      v-for="template in list"
      :key="template.id"
      :id="template.id"
      :name="template.name">
    </bk-option>
  </bk-select>
  <span v-else>
    <slot name="info-prepend"></slot>
    {{info}}
  </span>
</template>

<script>
  import activeMixin from './mixins/active'
  import { mapGetters } from 'vuex'
  import serviceTemplateService from '@/service/service-template/index.js'

  const requestId = Symbol('serviceTemplate')

  export default {
    name: 'cmdb-search-service-template',
    mixins: [activeMixin],
    props: {
      value: {
        type: [Array, String],
        default: () => ([])
      },
      displayType: {
        type: String,
        default: 'selector',
        validator(type) {
          return ['selector', 'info'].includes(type)
        }
      }
    },
    data() {
      return {
        list: [],
        requestId
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      localValue: {
        get() {
          const { value } = this
          if (Array.isArray(value)) {
            return value
          }
          return value.split(',').map(id => parseInt(id, 10))
        },
        set(value) {
          this.$emit('input', value)
          this.$emit('change', value)
        }
      },
      info() {
        const info = []
        this.localValue.forEach((id) => {
          const data = this.list.find(data => data.id === id)
          data && info.push(data.name)
        })
        return info.join(' | ') || '--'
      }
    },
    created() {
      this.getServiceTemplate()
    },
    methods: {
      async getServiceTemplate() {
        try {
          const templates = await serviceTemplateService.findAll({ bk_biz_id: this.bizId }, {
            requestId: this.requestId
          })

          this.list = this.$tools.localSort(templates, 'name')
        } catch (error) {
          console.error(error)
          this.list = []
        }
      }
    }
  }
</script>
