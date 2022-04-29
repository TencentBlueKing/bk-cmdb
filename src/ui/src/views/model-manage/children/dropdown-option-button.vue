<template>
  <div
    :class="['dropdown-option-button', { disabled: disabled }, { 'is-active': active }]"
    @click="handleClick"
  >
    <slot></slot>
  </div>
</template>

<script>
  export default {
    name: 'DropdownOptionButton',
    props: {
      // 禁用选项
      disabled: {
        type: Boolean,
        default: false
      },
      // 激活选项
      active: {
        type: Boolean,
        default: false
      },
    },
    methods: {
      handleClick() {
        if (!this.disabled) {
          this.$emit('click')
        }
      },
    },
  }
</script>

<style lang="scss" scoped>
.dropdown-option-button {
  height: 32px;
  line-height: 32px;
  white-space: nowrap;
  color: #63656E;
  padding: 0 12px;
  font-size: 12px;
  user-select: none;

  &:hover,
  &.is-active {
    background-color: #e1ecff;
    color: #3A84FF;
  }

  &.disabled {
    color: $textDisabledColor;
    cursor: not-allowed;
  }
}
</style>
