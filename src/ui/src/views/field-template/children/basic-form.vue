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
  import { ref, watch, getCurrentInstance, computed } from 'vue'
  import cloneDeep from 'lodash/cloneDeep'
  import { DUP_CHECK_IDS } from '@/setup/duplicate-remote-validate'

  const props = defineProps({
    data: {
      type: Object,
      default: () => ({})
    }
  })

  const emit = defineEmits(['change'])

  const currentInstance = getCurrentInstance().proxy

  const dupCheckId = DUP_CHECK_IDS.FIELD_TEMPLATE_NAME
  const originName = computed(() => props.data?.name)

  const formData = ref(cloneDeep(props.data))
  watch(() => props.data, (data) => {
    formData.value = cloneDeep(data)
  }, { deep: true })

  const validateAll = async () => {
    const validResults = []
    for (const [key, val] of Object.entries(currentInstance.fields)) {
      if (val.dirty) {
        validResults.push(await currentInstance.$validator.validate(key))
      }
    }

    return validResults.every(valid => valid)
  }

  const handleFormDataChange = () => {
    emit('change')
  }

  defineExpose({
    formData,
    validateAll
  })
</script>

<template>
  <div class="basic-form">
    <bk-form :label-width="140" :model="formData">
      <bk-form-item :label="$t('模板名称')" :required="true" :property="'name'"
        class="cmdb-form-item" :class="{ 'is-error': errors.has('name') }">
        <bk-input
          data-vv-validate-on="change"
          v-validate="`required|length:256|remoteDuplicate:${dupCheckId},${originName},${$t('模板名称已存在')}`"
          name="name"
          :placeholder="$t('请输入模板名称')"
          v-model="formData.name"
          @change="handleFormDataChange">
        </bk-input>
        <p class="form-error" v-if="errors.has('name')">{{errors.first('name')}}</p>
      </bk-form-item>
      <bk-form-item :label="$t('描述')" :property="'description'"
        class="cmdb-form-item" :class="{ 'is-error': errors.has('description') }">
        <bk-input
          type="textarea"
          :rows="4"
          name="description"
          v-validate="'length:2000'"
          v-model="formData.description"
          :placeholder="$t('请输入模板描述')"
          @change="handleFormDataChange">
        </bk-input>
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
