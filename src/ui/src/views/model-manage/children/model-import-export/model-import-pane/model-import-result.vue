<template>
  <div class="model-import-result">
    <div class="import-icon">
      <img src="@/assets/images/import-package.png" />
    </div>
    <p class="result-main-text" :class="{ [`is-${data.status}`]: data.status }">
      <i18n path="模型导入结果提示">
        <span class="succeed-item-count" place="modelCount">{{data.modelCount}}</span>
        <span class="succeed-item-count" place="relationTypeCount">{{data.relationTypeCount}}</span>
      </i18n>
    </p>
  </div>
</template>
<script>
  import { defineComponent } from '@vue/composition-api'
  import { t } from '@/i18n'

  export default defineComponent({
    nam: 'ModelImportResult',
    props: {
      data: {
        type: Object,
        default: () => ({
          status: 'success',
          modelCount: 0,
          relationTypeCount: 0
        })
      }
    },
    setup(props, { emit }) {
      const retry = () => emit('retry')
      return {
        t,
        retry
      }
    },
  })
</script>

<style lang="scss" scoped>
.model-import-result {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: #fff;
}
.import-icon {
  width: 130px;
  margin-top: 50px;
}
.import-icon img {
  width: 100%;
  height: 100%;
}
.result-main-text{
  font-size: 24px;
  margin-top: 20px;
}
.succeed-item-count{
  color: $successColor;
  margin: 0 5px;
}
</style>
