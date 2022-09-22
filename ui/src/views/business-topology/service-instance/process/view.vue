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
  <section class="view-process" v-test-id="'viewProcess'">
    <process-options></process-options>
    <process-list class="process-list"></process-list>
  </section>
</template>

<script>
  import ProcessOptions from './options'
  import ProcessList from './list'
  export default {
    name: 'view-process',
    components: {
      ProcessOptions,
      ProcessList
    },
    data() {
      return {}
    }
  }
</script>
