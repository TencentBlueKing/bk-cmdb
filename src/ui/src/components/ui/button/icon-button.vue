<template>
  <button class="bk-button bk-default icon-button" @click="handleClick($event)">
    <i :class="['icon-wrapper', icon]"></i>
  </button>
</template>

<script>
  export default {
    name: 'icon-button',
    props: {
      disabled: Boolean,
      icon: {
        type: String,
        required: true
      }
    },
    methods: {
      handleClick(event) {
        this.$emit('click', event)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .icon-button {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        width: 32px;
        height: 32px;
        line-height: 30px;
        padding: 0;
        font-size: 0px;
        .icon-wrapper {
            font-size: 0;
            &:before {
                font-size: 14px;
            }
        }
    }
</style>
