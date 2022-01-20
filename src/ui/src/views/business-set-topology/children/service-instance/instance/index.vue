<template>
  <section class="view-instance" v-test-id="'viewInst'">
    <instance-options class="instance-options"></instance-options>
    <instance-list class="instance-list"></instance-list>
  </section>
</template>

<script>
  import InstanceOptions from './options.vue'
  import InstanceList from './instance-list'
  export default {
    name: 'view-instance',
    components: {
      InstanceOptions,
      InstanceList
    },
    data() {
      return {}
    }
  }
</script>
