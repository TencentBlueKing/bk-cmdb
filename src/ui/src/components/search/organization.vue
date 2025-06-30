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
  <div class="search-value-container" v-if="displayType === 'info'">
    <div class="prepend"><slot name="info-prepend"></slot></div>
    <org-value :property="property" :value="organizationVal" show-on="search"></org-value>
  </div>
  <div v-else>
    <bk-tag-input
      v-model="tagInputVal"
      ref="tagInput"
      :placeholder="$t('点击选择组织')"
      :collapse-tags="true"
      :allow-create="true"
      @focus="handleFocus"
      @removeAll="handleDelete">
    </bk-tag-input>
    <cmdb-form-organization
      :hidden="true"
      ref="organization"
      v-model="organizationVal"
      v-bind="$attrs"
      @confirm="handleOrganizationConfirm">
    </cmdb-form-organization>
  </div>

</template>

<script>
  import activeMixin from './mixins/active'
  import orgValue from '@/components/ui/other/org-value.vue'
  import { parseOrgVal } from '@/utils/tools'
  import store from '@/store'

  export default {
    name: 'cmdb-search-organization',
    components: {
      orgValue
    },
    mixins: [activeMixin],
    props: {
      value: {
        type: [Array, String, Number],
        default: () => []
      },
      property: {
        type: Object,
        default: () => ({})
      },
      multiple: {
        type: Boolean,
        default: true
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
        organizationVal: [], // 组织选择器回显的值 格式[{ type: 'org', id: 1 }]
        tagInputVal: [] // tagInput显示的值 格式['xxx', 'xxx']
      }
    },
    watch: {
      value: {
        async handler(val) {
          if (!val?.[0]) {
            this.tagInputVal = []
            return
          }
          const value = val.join(',')
          const res = await store.dispatch('organization/getDepartment', value)
          this.tagInputVal = res?.map(item => parseOrgVal(item)) ?? []
          this.organizationVal = this.setOrganizationVal(val)
        },
        immediate: true
      }
    },
    methods: {
      setOrganizationVal(value) {
        if (value && !Array.isArray(value)) {
          return [{ id: value, type: 'org' }]
        }
        if (!value) {
          return []
        }
        return value.map(item => ({ id: item, type: 'org' })) || []
      },
      handleFocus() {
        this.$refs.organization?.openEdit()
        this.$refs.tagInput.$refs.input.blur()
      },
      handleDelete() {
        this.emitValue([])
      },
      handleOrganizationConfirm(val) {
        const value = val[0]?.data?.map(item => item.id) ?? []
        this.emitValue(value)
      },
      emitValue(value) {
        this.$emit('input', value)
        this.$emit('change', value)
        this.$emit('confirm')
      }
    }
  }
</script>

<style lang="scss" scoped>
  .search-value-container {
    display: flex;
    .prepend {
      margin-right: 4px;
    }
  }
</style>
