<script setup>
  const props = defineProps({
    steps: {
      type: Array,
      default: () => ([
        { title: '基础信息', icon: 1 },
        { title: '字段设置', icon: 2 },
      ])
    },
    current: {
      type: Number
    },
    width: {
      type: String
    }
  })
</script>

<template>
  <cmdb-steps :steps="props.steps" :current="current" :width="props.width"></cmdb-steps>
</template>
