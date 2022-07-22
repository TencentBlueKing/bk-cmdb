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
  <div class="import-state">
    <template v-if="state === 'pending'">
      <img class="state-loading" src="../../assets/images/icon/loading.svg" alt="loading">
      <p class="text">{{$t('数据导入中')}}</p>
    </template>
    <template v-else-if="state === 'success'">
      <i class="state-success bk-icon icon-check-circle-shape"></i>
      <p class="text">{{$t('导入成功')}}</p>
      <bk-button class="mt20" theme="default" @click="closeImport">{{$t('关闭')}}</bk-button>
    </template>
    <template v-else-if="state === 'error'">
      <i class="state-error bk-icon icon-close-circle-shape"></i>
      <p class="text">{{$t('导入失败')}}</p>
      <bk-button class="mt20" theme="default" @click="reset">{{$t('重新导入')}}</bk-button>
    </template>
    <ul class="error-list" v-if="instError">
      <li class="error-item"
        v-for="(error, index) in instError"
        :key="index">
        {{error}}
      </li>
    </ul>
    <ul class="error-list" v-if="asstError">
      <li class="error-item"
        v-for="(error, index) in asstError"
        :key="index">
        {{`${$t('第N行', { row: error.row })}：${error.message}`}}
      </li>
    </ul>
  </div>
</template>

<script>
  import useFile from './file'
  import useImport from './index'
  import useStep from './step'
  import { computed } from '@vue/composition-api'
  export default {
    name: 'import-state',
    setup() {
      const [{ state, error }, { clear: clearFile }] = useFile()
      const [, { close: closeImport }] = useImport()
      const [, { reset: resetStep }] = useStep()
      const reset = () => {
        resetStep()
        clearFile()
      }
      const asstError = computed(() => error.value && error.value.asst_error)
      const instError = computed(() => error.value && error.value.error)
      return {
        state,
        reset,
        closeImport,
        asstError,
        instError
      }
    }
  }
</script>

<style lang="scss" scoped>
  .import-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    .state-loading {
      margin: 20px 0 16px 0;
      width: 60px;
    }
    .text {
      font-size: 22px;
      color: #313238;
      line-height: 30px;
    }
    .state-success,
    .state-error {
      font-size: 60px;
      margin: 20px 0 16px 0;
    }
    .state-success {
      color: $successColor;
    }
    .state-error {
      color: $dangerColor;
    }
    .error-list {
      width: 100%;
      background: #ffeded;
      border: 1px solid #ffd2d2;
      border-radius: 3px;
      margin: 18px 0 0 0;
      padding: 18px;
      font-size: 12px;
      line-height: 18px;
    }
  }
</style>
