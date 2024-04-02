<script setup>
  import Notice from '@blueking/notice-component-vue2'
  import '@blueking/notice-component-vue2/dist/style.css'
  import { ref, onMounted, onBeforeUnmount } from 'vue'

  const emit = defineEmits(['show-change', 'size-change'])

  const dataUrl = `${window.API_HOST}notice/get_current_announcements`
  const noticeComp = ref(null)

  const noticeShowChange = (isShow) => {
    emit('show-change', isShow)
  }

  const execResizeHander = () => {
    const el = noticeComp.value.$el
    const height = el.offsetHeight
    const width = el.offsetHeight
    emit('size-change', [width, height])
  }
  const resizeObserver = new ResizeObserver((entries) => {
    for (const entry of entries) {
      if (entry.target === noticeComp.value.$el) {
        execResizeHander()
      }
    }
  })
  onMounted(() => {
    resizeObserver.observe(noticeComp.value.$el)
  })
  onBeforeUnmount(() => {
    resizeObserver.unobserve(noticeComp.value.$el)
  })
</script>

<template>
  <notice ref="noticeComp" :api-url="dataUrl" @show-alert-change="noticeShowChange" />
</template>
