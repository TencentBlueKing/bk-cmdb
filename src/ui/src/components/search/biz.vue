<template>
  <cmdb-tag-input
    v-model="localValue"
    v-bind="$attrs"
    :fuzzy-search-method="fuzzySearchMethod"
    :exact-search-method="exactSearchMethod"
    @clear="() => $emit('clear')"
    @focus="handleToggle(true, ...arguments)"
    @blur="handleToggle(false, ...arguments)">
  </cmdb-tag-input>
</template>

<script>
  import autocomplete from './mixins/autocomplete-biz'
  export default {
    name: 'cmdb-search-biz',
    mixins: [autocomplete],
    data() {
      return {
        type: 'biz'
      }
    }
  }
</script>
