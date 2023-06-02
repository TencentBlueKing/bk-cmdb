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

<script setup>
  import { defineExpose, ref, watch } from 'vue'
  import cloneDeep from 'lodash/cloneDeep'

  const props = defineProps({
    data: {
      type: Object,
      default: () => ({})
    }
  })

  const formData = ref(cloneDeep(props.data))
  watch(() => props.data, (data) => {
    formData.value = cloneDeep(data)
  }, { deep: true })

  defineExpose({
    formData
  })
</script>

<template>
  <div class="basic-form">
    <bk-form :label-width="140" :model="formData">
      <bk-form-item label="模板名称" :required="true" :property="'name'"
        class="cmdb-form-item" :class="{ 'is-error': errors.has('name') }">
        <bk-input v-validate="'required|length:256'" name="name" v-model="formData.name"></bk-input>
        <p class="form-error" v-if="errors.has('name')">{{errors.first('name')}}</p>
      </bk-form-item>
      <bk-form-item label="描述"
        class="cmdb-form-item" :class="{ 'is-error': errors.has('description') }">
        <bk-input type="textarea" :rows="4" name="description" v-validate="'length:2000'"
          v-model="formData.description" placeholder="请输入描述"></bk-input>
        <p class="form-error" v-if="errors.has('description')">{{errors.first('description')}}</p>
      </bk-form-item>
    </bk-form>
  </div>
</template>

<style lang="scss" scoped>
  .basic-form {
    width: 628px;
    margin: 64px auto 12px;
    position: relative;
    left: -36px;
  }
</style>