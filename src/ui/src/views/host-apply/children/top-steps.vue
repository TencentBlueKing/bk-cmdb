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

<template functional>
  <div class="top-steps">
    <bk-steps
      class="steps"
      :steps="[
        { title: '属性设置', icon: 1 },
        { title: '信息确认', icon: 2 },
        { title: '属性应用', icon: 3 }
      ]"
      :cur-step="props.current">
    </bk-steps>
  </div>
</template>

<style lang="scss" scoped>
  .top-steps {
    display: flex;
    height: 50px;
    align-items: center;
    justify-content: center;
    padding: 0 20px;
    background-color: #fff;
    border-bottom: 1px solid $borderColor;
    .steps {
      width: 65%;
    }
  }
</style>
