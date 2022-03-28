<template>
  <div class="model-export-result">
    <div class="download-icon">
      <img src="@/assets/images/package.png" />
    </div>
    <p class="result-main-text">{{t('模型导出成功！请及时保存并妥善保管。')}}</p>
    <p class="result-sub-text">
      {{t('如果页面未出现提示保存窗口，请点击')}}
      <span @click="retryDownload" class="retry-download-button">{{t('重新下载文件')}}</span>
    </p>
  </div>
</template>
<script>
  import { defineComponent } from '@vue/composition-api'
  import { t } from '@/i18n'

  export default defineComponent({
    name: 'ModelExportResult',
    setup(props, { emit }) {
      const retryDownload = () => emit('retry-download')

      return {
        t,
        retryDownload
      }
    },
  })
</script>

<style lang="scss" scoped>
.model-export-result {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: #fff;
}
.download-icon {
  width: 130px;
  margin-top: 50px;
}
.download-icon img {
  width: 100%;
  height: 100%;
}
.result {
  &-main-text {
    font-size: 24px;
    margin-top: 20px;
  }
  &-sub-text {
    margin-top: 10px;
    font-size: 14px;
  }
}
.retry-download-button {
  color: #3a84ff;
  cursor: pointer;
}
</style>
