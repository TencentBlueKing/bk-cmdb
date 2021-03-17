<template>
  <i class="loading" v-if="loading"></i>
  <span v-else><slot></slot></span>
</template>

<script>
  export default {
    props: {
      loading: Boolean
    }
  }
</script>

<style lang="scss" scoped>
    .loading {
        display: inline-block;
        width: 16px;
        height: 16px;
        background-color: transparent;
        background-image: url("../../assets/images/icon/loading.svg");
        background-position: center center;
        background-size: 16px;
        background-repeat: no-repeat;
    }
</style>
