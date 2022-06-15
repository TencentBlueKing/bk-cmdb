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
  import { computed, defineComponent, PropType, ref } from '@vue/composition-api'
  import GridLayout from '@/components/ui/other/grid-layout.vue'
  import GridItem from '@/components/ui/other/grid-item.vue'
  import EditableProperty from '@/components/ui/details/editable-property.vue'

  interface IProperty {
    id: number,
    'bk_property_id': string,
    'bk_property_name': string,
    'bk_isapi': boolean,
    'bk_property_group': string
  }

  export default defineComponent({
    components: {
      GridLayout,
      GridItem,
      EditableProperty
    },
    props: {
      properties: {
        type: Array as PropType<IProperty[]>,
        default: [],
        required: true
      },
      propertyIdKey: {
        type: String,
        default: 'bk_attribute_id'
      },
      instance: {
        type: Object,
        default: () => ({})
      },
      auth: Object,
      loadingState: {
        type: Array,
        default: () => ([])
      }
    },
    setup(props, { emit }) {
      const editState = ref({
        property: {},
        value: ''
      })

      const configList = computed(() => {
        const propertyIds = Object.keys(props.instance)
        return propertyIds.map(id => ({
          property: props.properties.find((item: IProperty) => item.bk_property_id === id),
          value: props.instance[id]
        }))
      })

      const isLoading = property => props.loadingState.includes(property)

      const exitEdit = () => {
        editState.value.property = {}
        editState.value.value = ''
      }

      const handleConfirmEdit = (changed: boolean) => {
        if (changed) {
          emit('save', editState.value)
        }
        exitEdit()
      }

      const handleConfirmDel = (property: IProperty) => {
        emit('del', property)
      }

      return {
        configList,
        editState,
        isLoading,
        handleConfirmEdit,
        handleConfirmDel
      }
    }
  })
</script>

<template>
  <div class="property-config-details">
    <grid-layout :min-width="360" :gap="20" class="form-content">
      <grid-item class="form-item" required
        v-for="({ property, value }) in configList"
        :key="property.id"
        :label="property.bk_property_name"
        :label-width="160">
        <editable-property
          :property="property"
          :auth="auth"
          :value="value"
          :loading="isLoading(property)"
          :edit-state.sync="editState"
          @confirm="handleConfirmEdit">
          <template #more-action>
            <bk-popconfirm
              trigger="click"
              :title="$t('确认删除该字段设置？')"
              :content="$t('确认删除模板字段提示')"
              @confirm="handleConfirmDel(property)">
              <bk-icon type="delete" class="property-del-button" />
            </bk-popconfirm>
          </template>
        </editable-property>
      </grid-item>
    </grid-layout>
  </div>
</template>

<style lang="scss" scoped>
  .property-del-button {
    cursor: pointer;
    font-size: 16px;
    margin-left: 4px;

    &:hover {
      color: $primaryColor;
    }
  }
</style>
