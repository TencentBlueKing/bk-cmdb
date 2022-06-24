<script lang="ts">
  import { computed, defineComponent, inject } from '@vue/composition-api'

  export default defineComponent({
    props: {
      label: String,
      labelWidth: {
        type: [Number, String]
      },
      direction: {
        type: String
      },
      required: {
        type: Boolean,
        default: false
      },
      gap: {
        type: Number
      }
    },
    setup(props) {
      const mode = inject('mode', 'detail')

      const labelContainerWidth = computed(() => {
        let { labelWidth } = props
        if (!labelWidth) {
          labelWidth = mode === 'detail' ? '160' : '100%'
        }

        return isNaN(Number(labelWidth)) ? labelWidth : `${labelWidth}px`
      })

      const itemGap = computed(() => (isNaN(Number(props.gap)) ? props.gap : `${props.gap}px`))

      return {
        mode,
        labelContainerWidth,
        itemGap
      }
    }
  })
</script>

<template>
  <div :class="['cmdb-grid-item', { required }, direction, mode]"
    :style="{ '--label-width': `${labelContainerWidth}`, '--flex-direction': direction, '--item-gap': itemGap }">
    <div class="item-label">
      <slot name="label">
        <div class="label-text" v-bk-overflow-tips>{{label}}</div>
      </slot>
    </div>
    <div class="item-content">
      <slot></slot>
    </div>
    <slot name="append"></slot>
  </div>
</template>

<style lang="scss" scoped>
  .cmdb-grid-item {
    display: flex;

    .item-label {
      display: flex;
      width: var(--label-width);
      font-size: 12px;
      color: #63656E;

      .label-text {
        @include ellipsis;
      }
    }

    .item-content {
      flex: 1;
      position: relative;
    }

    &.detail {
      align-items: center;
      flex-direction: var(--flex-direction, row);
      .item-label {
        margin-right: calc(var(--item-gap, 8px) / 2);
        .label-text {
          flex: 1;
          text-align: right;
        }

        &::after {
          content: "：";
        }
      }

      .item-content {
        margin-left: calc(var(--item-gap, 8px) / 2);
      }
    }

    &.form {
      flex-direction: var(--flex-direction, column);

      &:not(.row) .item-label {
        margin-bottom: calc(var(--item-gap, 8px) / 2);
      }

      &:not(.row) .item-content {
        margin-top: calc(var(--item-gap, 8px) / 2);
      }

      &.row {
        align-items: center;

        .item-label {
          margin-right: calc(var(--item-gap, 8px) / 2);
          .label-text {
            flex: 1;
            text-align: right;
          }
        }
        .item-content {
          margin-left: calc(var(--item-gap, 8px) / 2);
        }
      }
    }

    &.required {
      .label-text {
        &::after {
          content: "*";
          color: #ff5656;
          padding: 0 2px;
        }
      }
    }
  }
</style>
