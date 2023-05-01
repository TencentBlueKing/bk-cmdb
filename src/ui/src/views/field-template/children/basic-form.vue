<!--
 * Tencent
  components: { topSteps },is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<script setup>
  import { defineExpose, ref } from 'vue'
  import cloneDeep from 'lodash/cloneDeep'

  const props = defineProps({
    data: {
      type: Object,
      default: () => ({})
    }
  })

  const formData = ref(cloneDeep(props.data))

  defineExpose({
    formData
  })
</script>

<template>
  <bk-form :label-width="200" :model="formData">
    <bk-form-item label="模板名称" :required="true" :property="'name'">
      <bk-input v-model="formData.name"></bk-input>
    </bk-form-item>
    <bk-form-item label="描述">
      <bk-input type="textarea" v-model="formData.desc" placeholder="请输入描述"></bk-input>
    </bk-form-item>
  </bk-form>
</template>
