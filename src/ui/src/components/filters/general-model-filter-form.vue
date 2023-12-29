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
  <cmdb-sticky-layout class="filter-layout" slot="content">
    <bk-form class="filter-form" form-type="vertical">
      <bk-form-item class="filter-item"
        v-for="property in selected"
        :key="property.id"
        :class="`filter-item-${property.bk_property_type}`">
        <label class="item-label">
          {{property.bk_property_name}}
        </label>
        <div class="item-content-wrapper">
          <operator-selector class="item-operator"
            v-if="!withoutOperator.includes(property.bk_property_type)"
            :property="property"
            v-model="condition[property.id].operator"
            @change="handleOperatorChange(property, ...arguments)">
          </operator-selector>
          <component class="item-value"
            :is="getComponentType(property)"
            :placeholder="getPlaceholder(property)"
            :ref="`component-${property.id}`"
            v-bind="getBindProps(property)"
            v-model.trim="condition[property.id].value"
            v-bk-tooltips.top="{
              disabled: !property.placeholder,
              theme: 'light',
              trigger: 'click',
              content: property.placeholder
            }"
            @active-change="handleComponentActiveChange(property, ...arguments)">
          </component>
        </div>
      </bk-form-item>
      <bk-form-item>
        <bk-button class="filter-add-button ml10" type="primary" text @click="handleSelectProperty">
          {{$t('添加其他条件')}}
        </bk-button>
      </bk-form-item>
    </bk-form>
    <div class="filter-options"
      slot="footer"
      slot-scope="{ sticky }"
      :class="{ 'is-sticky': sticky }">
      <bk-button class="option-search mr10" theme="primary" :disabled="false" @click="handleSearch">
        {{$t('查询')}}
      </bk-button>
      <bk-button class="option-reset" theme="default" @click="handleReset">{{$t('清空')}}</bk-button>
    </div>
  </cmdb-sticky-layout>
</template>

<script>
  import { mapGetters } from 'vuex'
  import has from 'has'
  import OperatorSelector from './operator-selector.vue'
  import PropertySelector from './general-model-property-selector.js'
  import { setSearchQueryByCondition, resetConditionValue } from './general-model-filter.js'
  import Utils from './utils'

  export default {
    components: {
      OperatorSelector
    },
    props: {
      objId: {
        type: String
      },
      properties: {
        type: Array,
        default: () => ([])
      },
      propertyGroups: {
        type: Array,
        default: () => ([])
      },
      filterSelected: {
        type: Array,
        default: () => ([])
      },
      filterCondition: {
        type: Object,
        default: () => ({})
      }
    },
    data() {
      return {
        withoutOperator: ['date', 'time', 'bool'],
        condition: {},
        selected: []
      }
    },
    computed: {
      ...mapGetters('objectModelClassify', ['getModelById']),
    },
    watch: {
      filterSelected: {
        immediate: true,
        handler() {
          const newCondition = this.$tools.clone(this.filterCondition)
          Object.keys(newCondition).forEach((id) => {
            if (has(this.condition, id)) {
              newCondition[id] = this.condition[id]
            }
          })
          this.condition = newCondition
          this.selected = [...this.filterSelected]
        }
      },
      selected() {
        this.updateCondition()
      }
    },
    methods: {
      updateCondition() {
        const newConditon = {}
        this.selected.forEach((property) => {
          if (has(this.condition, property.id)) {
            newConditon[property.id] = this.condition[property.id]
          } else {
            newConditon[property.id] = Utils.getDefaultData(property)
          }
        })
        this.condition = newConditon
      },
      getComponentType(property) {
        const {
          bk_obj_id: modelId,
          bk_property_id: propertyId,
          bk_property_type: propertyType
        } = property
        const normal = `cmdb-search-${propertyType}`

        // 业务名在包含与非包含操作符时使用输入联想组件
        if (modelId === 'biz' && propertyId === 'bk_biz_name' && this.condition[property.id].operator !== '$regex') {
          return `cmdb-search-${modelId}`
        }

        const isSetName = modelId === 'set' && propertyId === 'bk_set_name'
        const isModuleName = modelId === 'module' && propertyId === 'bk_module_name'
        if (isSetName || isModuleName) {
          return `cmdb-search-${modelId}`
        }
        return normal
      },
      getPlaceholder(property) {
        return Utils.getPlaceholder(property)
      },
      getBindProps(property) {
        return Utils.getBindProps(property)
      },
      handleOperatorChange(property, operator) {
        const { value } = this.condition[property.id]
        const effectValue = Utils.getOperatorSideEffect(property, operator, value)
        this.condition[property.id].value = effectValue
      },
      // 人员选择器参考定位空间不足，备选面板左移了，此处将其通过offset配置移到最右边
      handleComponentActiveChange(property, active) {
        if (!active) {
          return false
        }
        const { id, bk_property_type: type } = property
        if (type !== 'objuser') {
          return false
        }
        const [component] = this.$refs[`component-${id}`]
        try {
          this.$nextTick(() => {
            const reference = component.$el.querySelector('.user-selector-input')
            // eslint-disable-next-line no-underscore-dangle
            reference._tippy.setProps({
              offset: [240, 5]
            })
          })
        } catch (error) {
          console.error(error)
        }
      },
      async handleRemove(property) {
        const index = this.selected.indexOf(property)
        index > -1 && this.selected.splice(index, 1)
      },
      handleSelectProperty() {
        const { objId, properties, propertyGroups, selected: propertySelected } = this
        PropertySelector.show({
          objId,
          properties,
          propertyGroups,
          propertySelected,
          handler: this.updateSelected
        })
      },
      updateSelected(selected) {
        // 将触发updateCondition新的条件项会被生成
        this.selected = selected
      },
      handleSearch() {
        // tag-input组件在blur时写入数据有200ms的延迟，此处等待更长时间，避免无法写入
        this.searchTimer && clearTimeout(this.searchTimer)
        this.searchTimer = setTimeout(() => {
          setSearchQueryByCondition(this.condition, this.selected)
          this.$emit('close')
        }, 300)
      },
      handleReset() {
        this.condition = resetConditionValue(this.condition, this.selected)
      }
    }
  }
</script>

<style lang="scss" scoped>
  .filter-layout {
    height: 100%;
    @include scrollbar-y;
  }

  .filter-form {
    padding: 0 10px;
  }
  .filter-item {
    padding: 2px 10px 10px;
    margin-top: 5px !important;
    &:hover {
      background: #f5f6fa;
      .item-remove {
        opacity: 1;
      }
    }
    .item-label {
      display: block;
      font-size: 14px;
      font-weight: 400;
      line-height: 24px;
      @include ellipsis;
    }
    .item-content-wrapper {
      display: flex;
      align-items: center;
    }
    .item-operator {
      flex: 110px 0 0;
      margin-right: 8px;
      & ~ .item-value {
        max-width: calc(100% - 118px);
      }
    }
    .item-value {
      flex: 1;
    }
    .item-remove {
      position: absolute;
      width: 24px;
      height: 24px;
      display: flex;
      justify-content: center;
      align-items: center;
      right: -10px;
      top: 3px;
      font-size: 20px;
      opacity: 0;
      cursor: pointer;
      color: $textColor;
      &:hover {
        color: $dangerColor;
      }
    }
  }

  .filter-options {
    display: flex;
    align-items: center;
    padding: 10px 20px;
    &.is-sticky {
      border-top: 1px solid $borderColor;
      background-color: #fff;
    }
    .option-collect,
    .option-collect-wrapper {
      & ~ .option-reset {
        margin-left: auto;
      }
    }
  }
</style>
