<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <cmdb-sticky-layout class="filter-layout" slot="content" ref="propertyList" v-scroll="{
    targetClass: 'last-item',
    orientation: 'bottom',
    distance: 63
  }">
    <bk-form class="filter-form" form-type="vertical">
      <div class="filter-operate">
        <condition-picker
          ref="conditionPicker"
          class="filter-add"
          :text="$t('添加条件')"
          :selected="selected"
          :property-map="propertyMap"
          :handler="updateSelected"
          :type="2">
        </condition-picker>
        <bk-popconfirm
          :content="$t('确定清空筛选条件')"
          width="280"
          trigger="click"
          :confirm-text="$t('确定')"
          :cancel-text="$t('取消')"
          @confirm="handleClearCondition">
          <bk-button :text="true" class="mr10" theme="primary"
            :disabled="!selected.length">
            {{$t('清空条件')}}
          </bk-button>
        </bk-popconfirm>
      </div>
      <bk-form-item class="filter-item"
        v-for="(property, index) in selected"
        :key="property.id"
        :class="[`filter-item-${property.bk_property_type}`, {
          'last-item': index === selected.length - 1 && scrollToBottom
        }]">
        <label class="item-label">
          {{property.bk_property_name}}
        </label>
        <div class="item-content-wrapper">
          <operator-selector class="item-operator"
            v-if="!withoutOperator.includes(property.bk_property_type)"
            :property="property"
            :custom-type-map="customOperatorTypeMap"
            :symbol-map="operatorSymbolMap"
            :desc-map="operatorDescMap"
            v-model="condition[property.id].operator"
            @change="handleOperatorChange(property, ...arguments)">
          </operator-selector>
          <component class="item-value r0"
            :is="getComponentType(property)"
            :placeholder="getPlaceholder(property)"
            :property="property"
            :is-paste-split="getPasteSplit(property.bk_property_id)"
            :ref="`component-${property.id}`"
            v-bind="getBindProps(property)"
            v-model.trim="condition[property.id].value"
            v-bk-tooltips.top="{
              disabled: !property.placeholder,
              theme: 'light',
              trigger: 'click',
              content: property.placeholder
            }"
            @active-change="handleComponentActiveChange(property, ...arguments)"
            @change="handleChange"
            @inputchange="hanleInputChange"
            @click.native="() => handleClick(`component-${property.id}`)"
            :popover-options="{
              duration: 0,
              onShown: handleShow,
              onHidden: handlePopoverHidden
            }">
          </component>
        </div>
        <i class="item-remove bk-icon icon-close" @click="handleRemove(property)"></i>
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
  import { setSearchQueryByCondition, resetConditionValue } from './general-model-filter.js'
  import Utils from './utils'
  import ConditionPicker from '@/components/condition-picker'
  import { getConditionSelect, updatePropertySelect, isPasteSplit } from '@/utils/util'
  import isEqual from 'lodash/isEqual'
  import useSideslider from '@/hooks/use-sideslider'
  import { QUERY_OPERATOR, QUERY_OPERATOR_OTHER_SYMBOL, QUERY_OPERATOR_OTHER_DESC } from '@/utils/query-builder-operator'
  import { POSITIVE_INTEGER } from '@/dictionary/property-constants'

  export default {
    components: {
      OperatorSelector,
      ConditionPicker
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
      const { IN, NIN, LIKE, CONTAINS_CS, EQ, NE, GTE, LTE, RANGE } = QUERY_OPERATOR
      return {
        scrollToBottom: false,
        withoutOperator: ['date', 'time', 'bool'],
        condition: {},
        originCondition: {},
        selected: [],
        customOperatorTypeMap: {
          float: [EQ, NE, GTE, LTE, RANGE, IN],
          int: [EQ, NE, GTE, LTE, RANGE, IN],
          longchar: [IN, NIN, LIKE, CONTAINS_CS],
          singlechar: [IN, NIN, LIKE, CONTAINS_CS],
          array: [IN, NIN, LIKE, CONTAINS_CS],
          object: [IN, NIN, LIKE, CONTAINS_CS]
        },
        operatorSymbolMap: QUERY_OPERATOR_OTHER_SYMBOL,
        operatorDescMap: QUERY_OPERATOR_OTHER_DESC
      }
    },
    computed: {
      ...mapGetters('objectModelClassify', ['getModelById']),
      propertyMap() {
        const modelPropertyMap = new Map()
        const ignoreProperties = [] // 预留，需要忽略的属性
        // eslint-disable-next-line max-len
        modelPropertyMap.set(this.objId, this.properties.filter(property => !ignoreProperties.includes(property.bk_property_id)))
        return Object.fromEntries(modelPropertyMap)
      },
      hasChange() {
        return !isEqual(this.condition, this.originCondition)
      },
      hasShow() {
        const { isShow } = this.$refs.conditionPicker
        return isShow
      }
    },
    watch: {
      filterSelected: {
        immediate: true,
        handler() {
          this.condition = this.setCondition(this.condition)
          this.selected = [...this.filterSelected]
        }
      },
      selected(val, oldVal) {
        this.updateCondition()
        const { addSelect, deleteSelect } = getConditionSelect(val, oldVal)
        this.scrollToBottom = this.hasAddSelected(val, oldVal, addSelect)
        updatePropertySelect(oldVal, this.handleRemove, addSelect, deleteSelect)
      }
    },
    created() {
      this.originCondition = this.setCondition(this.originCondition)
      const { beforeClose, setChanged } = useSideslider()
      this.beforeClose = beforeClose
      this.setChanged = setChanged
    },
    methods: {
      getPasteSplit(id) {
        return isPasteSplit(id)
      },
      setCondition(nowCondition) {
        const newCondition = this.$tools.clone(this.filterCondition)
        Object.keys(nowCondition).forEach((id) => {
          if (has(nowCondition, id)) {
            newCondition[id] = nowCondition[id]
          }
        })
        return newCondition
      },
      hasAddSelected(val, oldVal, addSelect) {
        return val[0] && oldVal[0] && addSelect.length > 0
      },
      handleClearCondition() {
        this.handleReset()
        this.selected = []
      },
      handleClick(e) {
        const parent = this.$refs[e][0].$el
        this.target = parent.getElementsByClassName('bk-select-tag-container')[0]
          || parent.getElementsByClassName('bk-tag-input')[0]
        if (~this.target?.className.indexOf('is-focus')) {
          // select专属
          return
        }
        this.calcPosition('click')
      },
      handleChange() {
        this.calcPosition()
      },
      hanleInputChange() {
        this.calcPosition()
      },
      handleShow() {
        this.calcPosition()
      },
      handlePopoverHidden() {
        this.$refs.propertyList.$el.classList.remove('over-height')
      },
      calcPosition(type = 'change') {
        if (type === 'click') this.$refs.propertyList.$el.classList.remove('over-height')
        if (!this.target) return

        this.$nextTick(() => {
          const limit = document.querySelector('.sticky-footer').getClientRects()[0].top
          const { bottom } = this.target.getClientRects()[0]
          if (bottom > Math.ceil(limit)) {
            this.$refs.propertyList.$el.classList.add('over-height')
          }
        })
      },
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
          bk_property_type: propertyType,
          id
        } = property
        const {
          operator
        } = this.condition[id]
        const normal = `cmdb-search-${propertyType}`

        // 数字类型int 和 float支持in操作符
        if (Utils.numberUseIn(property, operator)) {
          return 'cmdb-search-singlechar'
        }

        return normal
      },
      getPlaceholder(property) {
        return Utils.getPlaceholder(property)
      },
      getBindProps(property) {
        const props = Utils.getBindProps(property)
        if (POSITIVE_INTEGER.includes(property?.bk_property_id)) {
          if (!props.options) props.options = {}
          props.options.min = 1
        }
        // 数字类型int 和 float支持in操作符
        if (Utils.numberUseIn(property, this.condition[property?.id]?.operator)) {
          props.onlyNumber = true
          props.fuzzy = false
        }

        return props
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
  .over-height {
    .g-expand {
      bottom: 0;
    }
  }
  .filter-layout {
    height: 100%;
    @include scrollbar-y;
  }

  .filter-form {
    padding: 0 14px;
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
      align-items: flex-start;
      min-height: 32px;
    }
    .item-operator {
      flex: 128px 0 0;
      margin-right: 8px;
      & ~ .item-value {
        max-width: calc(100% - 136px);
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
  .filter-operate {
    @include space-between;
    position: sticky;
    top: 0;
    z-index: 9999;
    background: white;
    line-height: 30px;
  }

  .filter-add {
    padding-left: 10px;
  }
  .filter-options {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 24px;
    &.is-sticky {
      border-top: 1px solid $borderColor;
      background-color: #fff;
    }
  }
</style>
