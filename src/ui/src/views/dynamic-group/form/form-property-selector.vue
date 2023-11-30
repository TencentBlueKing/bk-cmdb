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
  <addCondition
    ref="addConditionComp"
    :selected="selected"
    :disabled-property-map="disabledPropertyMap"
    :models="models"
    :property-map="propertyMap"></addCondition>
</template>

<script setup>
  import { computed, ref, inject, reactive } from 'vue'
  import addCondition from '@/components/add-condition'

  const props = defineProps({
    selected: {
      type: Array,
      default: () => ([])
    },
    handler: Function
  })
  const addConditionComp = ref('')
  const dynamicGroupForm = inject('dynamicGroupForm')

  const disabledPropertyMap = reactive(dynamicGroupForm.disabledPropertyMap)
  const propertyMap = computed(() => dynamicGroupForm.propertyMap)

  const target = computed(() => dynamicGroupForm.formData.bk_obj_id)
  const models = computed(() => {
    if (target.value === 'host') {
      return dynamicGroupForm.availableModels
    }
    return dynamicGroupForm.availableModels.filter(model => model.bk_obj_id === target.value)
  })

  defineExpose({
    confirm: () => {
      props.handler && props.handler([...addConditionComp?.value?.localSelected])
    }
  })


</script>
