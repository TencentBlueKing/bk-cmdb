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
  <div class="business-scope-settings-form">
    <div class="condition-item">
      <bk-select
        v-model="selectedBusiness"
        :disabled="disabled"
        :placeholder="$t('业务范围选择placeholder')"
        :list="allBusiness"
        class="select-business"
        search-with-pinyin
        display-tag
        multiple
        searchable
        enable-virtual-scroll
        id-key="bk_biz_id"
        display-key="bk_biz_name">
      </bk-select>
    </div>
    <div class="condition-item" v-for="rule in condition" :key="rule.property.id">
      <cmdb-property-selector class="condition-field" v-if="rule.property.id"
        v-model="rule.field"
        :disabled="disabled"
        :properties="getAvailableProperties(rule.property)"
        :searchable="unusedProperties.length > 1"
        :loading="loading.property"
        @change="handleFieldChange">
      </cmdb-property-selector>
      <component :class="['condition-value', rule.property.bk_property_type]" v-if="rule.property.id"
        :is="`cmdb-search-${rule.property.bk_property_type}`"
        :placeholder="getPlaceholder(rule.property)"
        :clearable="true"
        :multiple="true"
        :disabled="disabled"
        v-bind="getBindProps(rule.property)"
        v-model="rule.value">
      </component>
      <i class="bk-icon icon-close" @click="handleRemove(rule)"></i>
    </div>
    <bk-button class="condition-button"
      :disabled="!unusedProperties.length || disabled"
      icon="icon-plus-circle"
      :text="true"
      @click="handleAdd">
      {{$t('添加其他条件')}}
    </bk-button>
  </div>
</template>

<script>
  import { defineComponent, computed, watchEffect, reactive, watch, toRefs, ref } from '@vue/composition-api'
  import Utils from '@/components/filters/utils'
  import cmdbPropertySelector from '@/components/property-selector'
  import propertyService from '@/service/property/property.js'
  import businessService from '@/service/business/search.js'

  export default defineComponent({
    components: {
      cmdbPropertySelector
    },
    props: {
      data: {
        type: Object,
        default: () => ({})
      },
      disabled: {
        type: Boolean,
        default: false
      }
    },
    setup(props, { emit }) {
      const { data: formData } = toRefs(props)

      const loading = reactive({
        property: false
      })

      // 业务属性列表
      const properties = ref([])

      // 所有业务列表
      const allBusiness = ref([])

      // 存放所有条件项
      const condition = ref([])
      // 已选择的业务
      const selectedBusiness = ref([])

      // 初始化表单项的值
      watchEffect(() => {
        const localCondition = formData.value?.condition || []
        const localSelectedBusiness = formData.value?.selectedBusiness || []

        const newCondition = []
        localCondition.forEach((item) => {
          const conditionCopy = { ...item }
          conditionCopy.property = Utils.findProperty(item.field, properties.value) || {}
          newCondition.push(conditionCopy)
        })

        condition.value = newCondition
        selectedBusiness.value = localSelectedBusiness
      })

      // 初始化业务属性和全量业务列表
      watchEffect(async () => {
        loading.property = true
        const [businessProperties, businessList] = await Promise.all([
          propertyService.findBiz(),
          businessService.findAll()
        ])
        loading.property = false

        const allowedPropertyTypes = ['organization', 'enum']
        properties.value = businessProperties.filter(item => allowedPropertyTypes.includes(item.bk_property_type))
        allBusiness.value = businessList
      })

      // 得到一个所有条件用到的propertyMap
      const conditionPropertyMap = computed(() => {
        const propertyMap = new WeakMap()
        condition.value.forEach(item => propertyMap.set(item.property, item.field))
        return propertyMap
      })

      // 未被使用过的属性列表，用于限定属性下拉选项
      const unusedProperties = computed(() => properties.value.filter(item => !conditionPropertyMap.value.has(item)))

      // 根据property动态获取当前可使用的属性列表
      const getAvailableProperties = property => [property, ...unusedProperties.value]

      // 条件项的字段变更时更新对应的property
      const handleFieldChange = (value) => {
        const rule = condition.value.find(rule => rule.field === value)
        rule.property = properties.value.find(item => item.bk_property_id === value)
      }

      watch([condition, selectedBusiness], ([newCondition, newSelectedBusiness]) => {
        emit('change', {
          condition: newCondition,
          selectedBusiness: newSelectedBusiness
        })
      }, { immediate: true, deep: true })

      // 添加条件
      const handleAdd = () => {
        // 从未使用的属性中取第一个作为初始化选项
        const [property] = unusedProperties.value

        const { value } = Utils.getDefaultData(property)
        const { bk_property_id: field } = property

        condition.value.push({ field, value, property })
      }

      // 删除条件
      const handleRemove = (rule) => {
        const index = condition.value.indexOf(rule)
        if (~index) {
          condition.value.splice(index, 1)
        }
      }

      // 条件项的值对应的组件所需的相关方法
      const getPlaceholder = property => Utils.getPlaceholder(property)
      const getBindProps = property => Utils.getBindProps(property)

      return {
        loading,
        condition,
        properties,
        allBusiness,
        selectedBusiness,
        unusedProperties,
        getAvailableProperties,
        handleFieldChange,
        handleAdd,
        handleRemove,
        getPlaceholder,
        getBindProps
      }
    }
  })
</script>

<style lang="scss" scoped>
  .business-scope-settings-form {
    width: 100%;
    .condition-item {
      display: flex;
      align-items: center;
      position: relative;
      padding: 8px 20px 8px 8px;

      &:hover {
        background: #f5f6fa;
        .icon-close {
          opacity: 1;
        }
      }

      .condition-field {
        flex: 150px 0 0;
        margin-right: 8px;
      }
      .condition-value {
        flex: none;
        width: calc(100% - 158px);
        &.organization {
          font-size: 14px;
        }
      }
      .icon-close {
        position: absolute;
        width: 24px;
        height: 24px;
        display: flex;
        justify-content: center;
        align-items: center;
        right: -2px;
        top: 0;
        font-size: 20px;
        opacity: 0;
        cursor: pointer;
        color: $textColor;
        &:hover {
          color: $dangerColor;
        }
      }

      .select-business {
        flex: 1;
      }
    }

    .condition-button {
      margin: 8px 0 8px 8px;
      ::v-deep > div {
          display: flex;
          align-items: center;
          .bk-icon {
              top: 0;
          }
      }
    }
  }
</style>
