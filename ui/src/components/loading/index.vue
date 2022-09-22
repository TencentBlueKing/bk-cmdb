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
  <i class="loading" v-if="props.loading"></i>
  <span v-else><slot></slot></span>
</template>

<style lang="scss" scoped>
  .loading {
    display: inline-block;
    width: 16px;
    height: 16px;
    background-color: transparent;
    background-image: url("../../assets/images/icon/loading.svg");
    background-position: center center;
    background-size: 16px;
    background-repeat: no-repeat;
  }
</style>
