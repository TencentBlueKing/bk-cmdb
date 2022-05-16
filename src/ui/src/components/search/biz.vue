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
  <cmdb-tag-input
    v-model="localValue"
    v-bind="$attrs"
    :fuzzy-search-method="fuzzySearchMethod"
    :exact-search-method="exactSearchMethod"
    @clear="() => $emit('clear')"
    @focus="handleToggle(true, ...arguments)"
    @blur="handleToggle(false, ...arguments)">
  </cmdb-tag-input>
</template>

<script>
  import autocomplete from './mixins/autocomplete-biz'
  export default {
    name: 'cmdb-search-biz',
    mixins: [autocomplete],
    data() {
      return {
        type: 'biz'
      }
    }
  }
</script>
