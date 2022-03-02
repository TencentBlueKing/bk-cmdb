<template>
  <div class="export-status">
    <div class="status">
      <template v-if="pending">
        <img class="status-loading" src="../../assets/images/icon/loading.svg" alt="loading">
        <p class="text">{{$t('数据导入中')}}</p>
      </template>
      <template v-else-if="isFinished">
        <i class="status-success bk-icon icon-check-circle-shape"></i>
        <p class="text">{{$t('数据导出成功')}}</p>
        <div>
          <bk-button class="mt20 mr10" theme="default" @click="resetStep">{{$t('重新导出')}}</bk-button>
          <bk-button class="mt20" theme="default" @click="close">{{$t('关闭')}}</bk-button>
        </div>
      </template>
      <template v-else-if="hasError">
        <i class="status-error bk-icon icon-close-circle-shape"></i>
        <p class="text">{{$t('导入失败')}}</p>
        <div>
          <bk-button class="mt20 mr10" theme="default" @click="restart">{{$t('重试失败')}}</bk-button>
          <bk-button class="mt20" theme="default" @click="close">{{$t('关闭')}}</bk-button>
        </div>
      </template>
    </div>
    <export-task></export-task>
  </div>
</template>

<script>
  import useTask from './task'
  import useState from './state'
  import exportTask from './export-task'
  import { computed } from '@vue/composition-api'
  export default {
    name: 'export-status',
    components: {
      exportTask
    },
    setup() {
      const [{ all, current }, { reset: resetTask, process: restart }] = useTask()
      const hasError = computed(() => all.value.some(task => task.state === 'error'))
      const isFinished = computed(() => all.value.every(task => task.state === 'finished'))
      const pending = computed(() => current.value && current.value.state === 'pending')

      const [{ visible }, { resetPartial }] = useState()
      const resetStep = () => {
        resetPartial()
        resetTask()
      }
      const close = () => {
        visible.value = false
      }

      return { hasError, isFinished, pending, resetStep, close, restart }
    }
  }
</script>

<style lang="scss" scoped>
  .export-status {
    .status {
      display: flex;
      flex-direction: column;
      align-items: center;
      .status-loading {
        margin: 20px 0 16px 0;
        width: 60px;
      }
      .text {
        font-size: 22px;
        color: #313238;
        line-height: 30px;
      }
      .status-success,
      .status-error {
        font-size: 60px;
        margin: 20px 0 16px 0;
      }
      .status-success {
        color: $successColor;
      }
      .status-error {
        color: $dangerColor;
      }
    }
  }
</style>
