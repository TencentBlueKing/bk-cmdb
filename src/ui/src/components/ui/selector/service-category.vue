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
  <bk-select class="service-category-selector"
    v-model="selected"
    searchable
    :multiple="multiple"
    :display-tag="multiple"
    :selected-style="getSelectedStyle"
    :clearable="allowClear"
    :disabled="disabled"
    :placeholder="placeholder"
    :font-size="fontSize"
    :popover-options="{
      boundary: 'window'
    }"
    ref="selector">
    <bk-option-group
      v-for="(group, groupIndex) in firstClassList"
      :name="group.name"
      :key="groupIndex">
      <bk-option v-for="(option, optionIndex) in group.secondCategory"
        :key="optionIndex"
        :id="option.id"
        :name="option.name">
      </bk-option>
    </bk-option-group>
  </bk-select>
</template>

<script>
  import { mapState, mapGetters } from 'vuex'
  import has from 'has'
  export default {
    name: 'cmdb-service-category',
    props: {
      value: {
        type: [Array, String],
        default: () => ([])
      },
      disabled: {
        type: Boolean,
        default: false
      },
      multiple: {
        type: Boolean,
        default: true
      },
      allowClear: {
        type: Boolean,
        default: false
      },
      autoSelect: {
        type: Boolean,
        default: true
      },
      placeholder: {
        type: String,
        default: ''
      },
      fontSize: {
        type: [String, Number],
        default: 'medium'
      }
    },
    data() {
      return {
        selected: this.value || [],
        firstClassList: []
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      ...mapState('businessHost', [
        'categoryMap'
      ]),
      getSelectedStyle() {
        return this.multiple ? 'checkbox' : 'check'
      }
    },
    watch: {
      value(value) {
        this.selected = value || []
      },
      selected(selected) {
        this.$emit('input', selected)
        this.$emit('on-selected', selected)
      }
    },
    created() {
      this.getServiceCategories()
    },
    methods: {
      async getServiceCategories() {
        if (has(this.categoryMap, this.bizId)) {
          this.firstClassList = this.categoryMap[this.bizId]
        } else {
          try {
            const data = await this.$store.dispatch('serviceClassification/searchServiceCategory', {
              params: { bk_biz_id: this.bizId }
            })
            const categories = this.collectServiceCategories(data.info)
            this.firstClassList = categories
            this.$store.commit('businessHost/setCategories', {
              id: this.bizId,
              categories
            })
          } catch (e) {
            console.error(e)
            this.firstClassList = []
          }
        }
      },
      collectServiceCategories(data) {
        const categories = []
        data.forEach((item) => {
          if (!item.category.bk_parent_id) {
            categories.push(item.category)
          }
        })
        categories.forEach((category) => {
          // eslint-disable-next-line max-len
          category.secondCategory = data.filter(item => item.category.bk_parent_id === category.id).map(item => item.category)
        })
        return categories
      }
    }
  }
</script>

<style lang="scss" scoped>
    .service-category-selector {
        width: 100%;
    }
</style>
