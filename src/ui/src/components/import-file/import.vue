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
  <section class="import">
    <bk-steps class="import-steps" :steps="steps" :cur-step="current" v-show="!showState"></bk-steps>
    <import-state v-if="showState"></import-state>
    <keep-alive>
      <import-file v-if="!showState && current === 1"></import-file>
    </keep-alive>
    <import-relation v-if="!showState && current === 2"></import-relation>
  </section>
</template>

<script>
  import importFile from './import-file'
  import importRelation from './import-relation'
  import importState from './import-state'
  import useStep from './step'
  import useFile from './file'
  import { computed } from 'vue'
  export default {
    name: 'host-import',
    components: {
      [importFile.name]: importFile,
      [importRelation.name]: importRelation,
      [importState.name]: importState
    },
    setup() {
      const [current] = useStep()
      const [{ state }] = useFile()
      const showState = computed(() => state.value && state.value !== 'resolving')
      return { current, showState }
    },
    data() {
      return {
        steps: [{ title: this.$t('上传文件'), icon: 1 }, { title: this.$t('选择关联模型'), icon: 2 }]
      }
    }
  }
</script>

<style lang="scss" scoped>
  .import {
    padding: 20px 28px;
    .import-steps {
      width: 350px;
      margin: 0 auto;
    }
  }
</style>
