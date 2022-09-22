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
  <div class="apply-layout">
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
      {{$t('转移属性变化确认提示')}}
    </cmdb-tips>
    <property-confirm-table
      class="table"
      ref="confirmTable"
      max-height="auto"
      :list="list"
      :render-icon="true"
      :show-operation="!!conflictList.length">
    </property-confirm-table>
  </div>
</template>

<script>
  import propertyConfirmTable from '@/components/host-apply/property-confirm-table'
  export default {
    name: 'host-attrs-auto-apply',
    components: {
      propertyConfirmTable
    },
    props: {
      info: {
        type: Array,
        required: true
      }
    },
    computed: {
      conflictList() {
        return this.info.filter(item => item.unresolved_conflict_count > 0)
      },
      list() {
        return this.conflictList.length ? this.conflictList : this.info
      }
    }
  }
</script>

<style lang="scss" scoped>
    .table {
        margin-top: 8px;
    }
</style>
