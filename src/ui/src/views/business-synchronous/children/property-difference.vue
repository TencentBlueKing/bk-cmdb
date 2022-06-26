<script lang="ts">
  import { defineComponent, PropType } from '@vue/composition-api'

  interface IPropertyDiff {
    id: number,
    'inst_value': unknown,
    'template_value': unknown,
    property: Record<string, unknown>
  }

  export default defineComponent({
    props: {
      moduleId: {
        type: Number,
        required: true
      },
      templateId: {
        type: Number,
        required: true
      },
      propertyDiff: {
        type: Array as PropType<IPropertyDiff[]>,
        default: () => ([]),
        required: true
      }
    },
    setup() {
      const getDiffType = (diff: IPropertyDiff) => {
        if (diff.inst_value !== diff.template_value) {
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
  <div class="property-difference">
    <div class="comparison-table">
      <div class="table-head">
        <div class="col before-col">属性同步前</div>
        <div class="col after-col">属性值同步后</div>
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
.property-difference {
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
      &.changed {
        color: #FF9C01;
      }
    }
  }
}
.comparison-table {
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
      padding-left: 90px;
    }
  }
}
</style>
