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

<script lang="ts">
  import { defineComponent, PropType } from '@vue/composition-api'
  import isEqual from 'lodash/isEqual'

  interface IPropertyDiff {
    id: number,
    'inst_value': unknown,
    'template_value': unknown,
    property: Record<string, unknown>
  }

  export default defineComponent({
    props: {
      propertyDiff: {
        type: Array as PropType<IPropertyDiff[]>,
        default: () => ([]),
        required: true
      }
    },
    setup() {
      const getDiffType = (diff: IPropertyDiff) => {
        if (!isEqual(diff.inst_value,  diff.template_value)) {
          return 'changed'
        }
      }
      return {
        getDiffType
      }
    }
  })
</script>

<template>
  <div class="property-config-difference">
    <div class="diff-table">
      <div class="table-head">
        <div class="col before-col">{{$t('属性同步前')}}</div>
        <div class="col after-col">{{$t('属性同步后')}}</div>
      </div>
      <div class="table-body">
        <div class="col before-col">
          <div class="diff-item" v-for="(diff, index) in propertyDiff" :key="index">
            <div class="property-name" v-bk-overflow-tips>{{diff.property.bk_property_name}}</div>
            <cmdb-property-value
              v-bk-overflow-tips
              class="property-value"
              tag="div"
              :value="diff.inst_value"
              :property="diff.property">
            </cmdb-property-value>
          </div>
        </div>
        <div class="col after-col">
          <div class="diff-item" v-for="(diff, index) in propertyDiff" :key="index">
            <div class="property-name">{{diff.property.bk_property_name}}</div>
            <cmdb-property-value
              v-bk-overflow-tips
              :class="['property-value', getDiffType(diff)]"
              tag="div"
              :value="diff.template_value"
              :property="diff.property">
            </cmdb-property-value>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.property-config-difference {
  .diff-item {
    display: flex;
    height: 28px;
    line-height: 28px;
    .property-name {
      width: 110px;
      text-align: right;
      @include ellipsis;
      &::after {
        content: "：";
      }
    }

    .property-value {
      width: calc(100% - 110px);
      @include ellipsis;
      &.changed {
        color: #FF9C01;
      }
    }
  }
}
.diff-table {
  display: grid;
  grid-template-rows: 32px auto;

  .table-head {
    display: grid;
    gap: 4px;
    grid-template-columns: 1fr 1fr;
    font-size: 12px;
    font-weight: 700;
    line-height: 32px;

    .col {
      padding-left: 24px;
      overflow: hidden;
    }
    .before-col {
      background: #F0F1F5;
    }
    .after-col {
      background: #DCDEE5;
    }
  }

  .table-body {
    display: grid;
    gap: 4px;
    grid-template-columns: 1fr 1fr;
    padding: 24px 0;
    font-size: 12px;
    background: #FAFBFD;

    .col {
      padding: 0 24px 0 90px;
      overflow: hidden;
    }
  }
}
</style>
