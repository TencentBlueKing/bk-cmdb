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
  <bk-sideslider
    v-transfer-dom
    :is-show.sync="slider.show"
    :title="slider.title"
    :width="800"
    :before-close="handleSliderBeforeClose">
    <div class="form-content" slot="content" style="height: 100%"
      v-bkloading="{ isLoading: $loading(Object.values(request)) }">
      <component :is="slider.component"
        ref="multipleForm"
        v-bind="slider.props"
        @on-submit="handleMultipleSave"
        @on-cancel="handleSliderBeforeClose">
      </component>
    </div>
  </bk-sideslider>
</template>

<script>
  import editMultipleHost from '@/mixins/edit-multiple-host'
  export default {
    mixins: [editMultipleHost]
  }
</script>
