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
  <bk-table-column :label="$t('状态')" prop="status">
    <template slot-scope="{ row }">
      <span v-if="isSyncing(row.status)" class="sync-status">
        <img class="svg-icon" src="../../../assets/images/icon/loading.svg" alt="">
        {{$t('同步中')}}
      </span>
      <span v-else-if="row.status === 'need_sync'" class="sync-status">
        <i class="status-circle waiting"></i>
        {{$t('待同步')}}
      </span>
      <span v-else-if="row.status === 'finished'" class="sync-status">
        <i class="status-circle success"></i>
        {{$t('已同步')}}
      </span>
      <span v-else-if="row.status === 'failure'"
        class="sync-status"
        v-bk-tooltips="{
          disabled: !row.fail_tips,
          content: row.fail_tips,
          placement: 'right'
        }">
        <i class="status-circle fail"></i>
        {{$t('同步失败')}}
      </span>
      <span v-else>--</span>
    </template>
  </bk-table-column>
</template>

<script>
  /**
   * 组件分别在集群模板实例、服务模板实例中使用，用于展示不同状态的变化，修改时要注意一起修改。
   */
  export default {
    name: 'InstanceStatusColumn',
    methods: {
      /**
       * 判断实例是否正在同步中
       */
      isSyncing(status) {
        return ['new', 'waiting', 'executing'].includes(status)
      },
    },
  }
</script>

<style lang="scss" scoped>
.sync-status {
    color: #63656E;
    .status-circle {
        display: inline-block;
        width: 8px;
        height: 8px;
        margin-right: 4px;
        border-radius: 50%;
        &.waiting {
            background-color: #3A84FF;
        }
        &.success {
            background-color: #2DCB56;
        }
        &.fail {
            background-color: #EA3536;
        }
    }
    .svg-icon {
        @include inlineBlock;
        margin-top: -4px;
        width: 16px;
    }
}
</style>
