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
  <div class="model-import-result">
    <div class="import-icon">
      <img src="@/assets/images/import-package.png" />
    </div>
    <p class="result-main-text" :class="{ [`is-${data.status}`]: data.status }">
      <i18n path="模型导入结果提示">
        <template #modelCount><span class="succeed-item-count">{{data.modelCount}}</span></template>
      </i18n>
    </p>
  </div>
</template>
<script>
  import { defineComponent } from '@vue/composition-api'
  import { t } from '@/i18n'

  export default defineComponent({
    nam: 'ModelImportResult',
    props: {
      data: {
        type: Object,
        default: () => ({
          status: 'success',
          modelCount: 0,
          // 暂时不显示，之后这里可能还要加上关联关系数
          relationTypeCount: 0
        })
      }
    },
    setup(props, { emit }) {
      const retry = () => emit('retry')
      return {
        t,
        retry
      }
    },
  })
</script>

<style lang="scss" scoped>
.model-import-result {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: #fff;
}
.import-icon {
  width: 130px;
  margin-top: 50px;
}
.import-icon img {
  width: 100%;
  height: 100%;
}
.result-main-text{
  font-size: 24px;
  margin-top: 20px;
}
.succeed-item-count{
  color: $successColor;
  margin: 0 5px;
}
</style>
