<template>
  <section class="view-instance">
    <instance-options class="instance-options"></instance-options>
    <instance-list class="instance-list"></instance-list>
  </section>
</template>

<script>
  import RouterQuery from '@/router/query'
  import Bus from '../common/bus'
  import InstanceOptions from './options'
  import InstanceList from './list'
  export default {
    name: 'view-instance',
    components: {
      InstanceOptions,
      InstanceList
    },
    data() {
      return {}
    },
    created() {
      Bus.$on('delete-complete', this.refreshView)
    },
    beforeDestroy() {
      Bus.$off('delete-complete', this.refreshView)
    },
    methods: {
      refreshView() {
        RouterQuery.set({
          _t: Date.now()
        })
      }
    }
  }
</script>
