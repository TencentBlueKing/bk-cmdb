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

<script lang="ts">
  import { defineComponent, ref, toRef, toRefs, PropType, watch } from 'vue'
  import { getPropertyDefaultValue } from '@/utils/tools.js'
  import GridLayout from '@/components/ui/other/grid-layout.vue'
  import GridItem from '@/components/ui/other/grid-item.vue'
  import PropertyFormElement from '@/components/ui/form/property-form-element.vue'
  import PropertyModal from './property-modal.vue'
  import useProperty from './use-property.js'
  import { isExclmationProperty } from '@/utils/util'

  interface IProperty {
    id: number,
    'bk_property_id': string,
    'bk_isapi': boolean,
    'bk_property_group': string
  }

  interface IPropertyGroup {
    'bk_biz_id': number
  }

  export default defineComponent({
    components: {
      GridLayout,
      GridItem,
      PropertyFormElement,
      PropertyModal
    },
    props: {
      properties: {
        type: Array as PropType<IProperty[]>,
        required: true
      },
      propertyGroups: {
        type: Array as PropType<IPropertyGroup[]>,
        required: true
      },
      selected: {
        type: Array as PropType<IProperty[]>,
        default: () => ([])
      },
      config: {
        type: Object,
        default: () => ({})
      },
      exclude: {
        type: Array as PropType<string[]>,
        default: () => ([])
      },
      formElementSize: String,
      formElementFontSize: String,
      maxColumns: Number
    },
    setup(props, { emit }) {
      const propertyFormEl = ref(null)

      const { sortedGroups, groupedProperties } = useProperty(toRefs(props))

      const propertyModalVisible = ref(false)

      // 当前选中的属性列表
      const selectedList = ref([])
      watch(() => props.selected, (selected) => {
        selectedList.value = selected.slice()
      })

      // 初始化需要展示的属性设置列表
      const {
        displayGroups: configPropertyGroups,
        displayProperties: configGroupedProperties
      } = useProperty({
        properties: selectedList,
        propertyGroups: toRef(props, 'propertyGroups'),
        exclude: toRef(props, 'exclude'),
      })

      // 配置结果
      const propertyConfig = ref({})

      // 使用传递进来的config初始化配置，统一使用属性id作为key
      watch(() => props.config, (config) => {
        for (const [id, value] of Object.entries(config)) {
          propertyConfig.value[id] = value
        }
      })

      // 选中项变化时同步配置列表，新添加的属性会被初始化
      watch(selectedList, (selectedList) => {
        const newConfig = {}
        selectedList.forEach((property) => {
          newConfig[property.id] = getPropertyDefaultValue(property, propertyConfig.value[property.id])
        })
        propertyConfig.value = newConfig
      }, { deep: true })

      const handleSelectField = () => {
        propertyModalVisible.value = true
      }

      const handleRemoveField = (property) => {
        const index = selectedList.value.indexOf(property)
        if (index !== -1) {
          selectedList.value.splice(index, 1)
        }

        emit('change', property)
      }

      const handleChange = (value, property) => {
        emit('change', property, value)
      }

      const isRequired = (property) => {
        const excludeType = ['bool']
        return !excludeType.includes(property.bk_property_type)
      }

      return {
        selectedList,
        sortedGroups,
        groupedProperties,
        propertyModalVisible,
        configPropertyGroups,
        configGroupedProperties,
        propertyConfig,
        propertyFormEl,
        isRequired,
        handleSelectField,
        handleRemoveField,
        handleChange
      }
    },
    methods: {
      isExclmationProperty(type) {
        return isExclmationProperty(type)
      },
      async validateAll() {
        // 获得每一个表单元素的校验方法
        const validates = (this.$refs.propertyFormEl || [])
          .map(formElement => formElement.$validator.validateAll())

        if (validates.length) {
          const results = await Promise.all(validates)
          return results.every(valid => valid)
        }

        return true
      },
      async validate() {
        const propertyFormEls = this.$refs.propertyFormEl || []
        const validates = []
        propertyFormEls.forEach(async (formElement) => {
          const [[key, val]] = Object.entries(formElement.fields)
          // 只检测dirty字段
          if (val.dirty) {
            validates.push(formElement.$validator.validate(key))
          }
        })

        if (validates.length) {
          const results = await Promise.all(validates)
          return results.every(valid => valid)
        }

        return true
      },
      getData() {
        return this.propertyConfig
      }
    }
  })
</script>

<template>
  <div class="property-config">
    <div class="select-trigger">
      <bk-button
        icon="plus"
        @click="handleSelectField">
        {{$t('添加属性字段')}}
      </bk-button>
      <slot name="tips"></slot>
    </div>
    <slot name="selected-list" v-bind="configPropertyGroups">
      <div class="selected-list" v-if="configPropertyGroups.length">
        <cmdb-collapse class="property-group"
          v-for="(group, groupIndex) in configPropertyGroups"
          :label="group.bk_group_name"
          arrow-type="filled"
          :key="groupIndex">
          <grid-layout mode="form"
            class="form-content"
            :min-width="360"
            :max-width="560"
            :gap="24"
            :max-columns="maxColumns">
            <grid-item class="form-item" required
              v-for="property in configGroupedProperties[groupIndex]"
              :key="property.id"
              :label="property.bk_property_name"
              :label-width="120">
              <template #label>
                <div class="label-text" v-bk-overflow-tips>
                  {{property.bk_property_name}}
                </div>
                <i class="property-name-tooltips icon-cc-tips"
                  v-if="property.placeholder && isExclmationProperty(property.bk_property_type)"
                  v-bk-tooltips.top="{
                    theme: 'light',
                    trigger: 'mouseenter',
                    content: property.placeholder
                  }">
                </i>
              </template>
              <property-form-element
                ref="propertyFormEl"
                :must-required="isRequired(property)"
                :property="property"
                :size="formElementSize"
                :font-size="formElementFontSize"
                v-model="propertyConfig[property.id]"
                @change="handleChange">
              </property-form-element>
              <template #append>
                <i class="item-remove bk-icon icon-close" @click="handleRemoveField(property)"></i>
              </template>
            </grid-item>
          </grid-layout>
        </cmdb-collapse>
      </div>
    </slot>
    <property-modal
      :visible.sync="propertyModalVisible"
      :selected-list.sync="selectedList"
      :sorted-groups="sortedGroups"
      :grouped-properties="groupedProperties">
    </property-modal>
  </div>
</template>

<style lang="scss" scoped>
  .property-config {
    .select-trigger {
      display: flex;
      align-items: center;
    }

    .selected-list {
      margin-top: 16px;

      .property-group {
        & + .property-group {
          margin-top: 12px;
        }
      }

      .form-content {
        padding: 24px;
      }

      .form-item {
        position: relative;
        padding: 2px 12px 4px 12px;

        &:hover {
          background: #f5f6fa;

          .item-remove {
            opacity: 1;
          }
        }

        .item-remove {
          position: absolute;
          width: 24px;
          height: 24px;
          display: flex;
          justify-content: center;
          align-items: center;
          right: 0;
          top: 0;
          font-size: 20px;
          opacity: 0;
          cursor: pointer;
          color: $textColor;
          &:hover {
            color: $dangerColor;
          }
        }

        ::v-deep .form-error {
          margin-top: 4px;
        }

        :deep(.item-label) {
          align-items: center;
        }
      }
    }
  }
</style>
