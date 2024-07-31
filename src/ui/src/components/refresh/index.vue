<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <div class="cmdb-refresh">
    <bk-button
      :size="size"
      class="button-refresh"
      icon="refresh"
      @click="handleFresh">
      {{$t('刷新')}}
    </bk-button>
    <div class="refresh-time" v-if="showTime">
      {{ data.text }}
    </div>
  </div>
</template>

<script setup>
  import { reactive, onBeforeUnmount, onMounted, ref } from 'vue'
  import { t } from '@/i18n'
  import interval from '@/utils/interval'

  const props = defineProps({
    showTime: {
      type: Boolean,
      default: true
    },
    size: {
      type: String,
      default: 'normal'
    }
  })
  const emit = defineEmits(['refresh'])

  const init = () => {
    data.text = t('刚刚刷新')
    data.time = getTime()
    timeInterval.clear()
    timeInterval.set()
  }
  const getTime = () => new Date().valueOf()
  const updateRefreshTime = () => {
    const now = getTime()
    const { time } = data
    const left = parseInt((now - time) / 60000, 10)
    let text = t('刚刚刷新')  // 小于1分钟
    if (left < 60 && left >= 1) {
      // 大于1分钟 小于60分钟
      text = t('x分钟前刷新', {
        time: left
      })
    }
    if (left >= 60) {
      text = t('大于1小时前刷新')
      timeInterval.clear()
    }
    data.text = text
  }
  const setCanRefresh = (val = false) => {
    canRefresh.value = val
  }

  const handleFresh = () => {
    if (!canRefresh.value) return
    emit('refresh')
    if (!props.showTime) return
    setCanRefresh()
    init()
  }

  const data = reactive({
    time: null,
    text: null
  })
  const timeInterval = reactive({
    set: () => ({}),
    clear: () => ({}),
  })
  const canRefresh = ref(true)

  onMounted(() => {
    if (!props.showTime) return
    const { clearTimeInterval, setTimeInterval } = interval(updateRefreshTime, 60000)
    timeInterval.set = setTimeInterval
    timeInterval.clear = clearTimeInterval
    init()
  })
  onBeforeUnmount(() => {
    timeInterval.clear()
  })

  defineExpose({
    setCanRefresh,
    init
  })

</script>

<style lang="scss" scoped>
  :deep(.button-refresh) {
    >div {
      display: flex;
      align-items: center;
      justify-content: space-between;
    }

    .icon-refresh {
      font-size: inherit;
      margin-right: 5px;
    }
  }
  .refresh-time {
    display: inline-block;
    font-size: 14px;
    margin-left: 3px;
  }
</style>
