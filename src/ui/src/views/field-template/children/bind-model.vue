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
  // import { reactive } from 'vue'
  const props = defineProps({
    modelList: {
      type: Array,
      default: () => ([])
    }
  })

  console.log(props)

</script>

<template>
  <div class="bind-model">
    <bk-resize-layout :collapsible="true" style="height: 500px;">
      <div slot="aside">
        model list
      </div>
      <div slot="main">main</div>
    </bk-resize-layout>
  </div>
</template>

<style lang="scss" scoped>
  .bind-model {
    height: 100%;
  }
</style>
