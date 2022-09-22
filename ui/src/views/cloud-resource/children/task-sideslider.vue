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
  <bk-sideslider v-transfer-dom
    :is-show.sync="isShow"
    :title="title"
    :width="900"
    @hidden="handleHidden">
    <component slot="content" ref="component"
      class="slider-content"
      :is="component"
      :container="this"
      v-bind="componentProps">
    </component>
  </bk-sideslider>
</template>

<script>
  import TaskForm from './task-form.vue'
  import TaskDetails from './task-details.vue'
  export default {
    components: {
      [TaskForm.name]: TaskForm,
      [TaskDetails.name]: TaskDetails
    },
    data() {
      return {
        isShow: false,
        title: '',
        component: null,
        componentProps: {}
      }
    },
    methods: {
      show({ mode, props, title }) {
        const componentMap = {
          create: TaskForm.name,
          update: TaskForm.name,
          details: TaskDetails.name
        }
        this.component = componentMap[mode]
        this.componentProps = props
        this.title = title
        this.isShow = true
      },
      hide() {
        this.isShow = false
      },
      handleHidden() {
        this.component = null
      }
    }
  }
</script>

<style lang="scss" scoped>
    .slider-content {
        height: 100%;
        @include scrollbar-y;
    }
</style>
