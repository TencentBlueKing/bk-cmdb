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
  <section class="delete-layout">
    <cmdb-tips
      :tips-style="{
        background: 'none',
        border: 'none',
        fontSize: '12px',
        lineHeight: '30px',
        padding: 0
      }"
      :icon-style="{
        color: '#63656E',
        fontSize: '14px',
        lineHeight: '30px'
      }">
      {{$t('删除服务实例提示')}}
    </cmdb-tips>
    <bk-table class="table"
      :data="info"
      :cell-style="getCellStyle">
      <bk-table-column :label="$t('操作')">
        <!-- eslint-disable-next-line -->
                <template slot-scope="{ row }">{{$t('删除')}}</template>
      </bk-table-column>
      <bk-table-column :label="$t('服务实例')" prop="name" show-overflow-tooltip></bk-table-column>
      <bk-table-column :label="$t('所属模块')" prop="bk_module_id" show-overflow-tooltip>
        <template slot-scope="{ row }">{{$parent.getModulePath(row.bk_module_id)}}</template>
      </bk-table-column>
    </bk-table>
  </section>
</template>

<script>
  export default {
    name: 'deleted-service-instance',
    props: {
      info: {
        type: Array,
        required: true
      }
    },
    methods: {
      getCellStyle({ columnIndex }) {
        if (columnIndex === 0) {
          return {
            color: '#FF5656'
          }
        }
        return {}
      }
    }
  }
</script>

<style lang="scss" scoped>
    .table {
        margin-top: 8px;
    }
</style>
