<template>
  <step-pane @cancel="cancelExport" :steps="exportSteps">
    <template #default="{ currentStepIndex }">
      <div class="step-container" v-show="currentStepIndex === 1">
        <!-- TODO 添加文案 -->
        <user-notice>
          1.
          模型所属分组以名称方式导出（非ID），即在导入的环境中将其归属到同名的分组中；若不存在同名模型分组，则会创建一个新的分组。<br />
          2.
          导出的关联类型如果在目标环境已存在（唯一标识），请在导入的流程中设置是否进行覆盖更新。<br />
          3.
          当导出的模型在导入环境中已存在（以模型唯一标识为准）时会显示为“模型冲突”，需要用户确认是否进行覆盖更新或不导入。<br />
          4.
          当导出的模型关联关系在导入环境中已存在（以模型和关联类型的唯一标识为准）时，无法导入。
        </user-notice>
      </div>
      <div class="step-container" v-show="currentStepIndex === 2">
        <model-export-editor v-model="confirmedExportData" :data="exportData"></model-export-editor>
      </div>
      <div class="step-container" v-show="currentStepIndex === 3">
        <model-export-setting ref="exportSettingRef" v-model="exportSetting"></model-export-setting>
      </div>
      <div class="step-container" v-show="currentStepIndex === 4">
        <model-export-result @retry-download="handleRetryDownload"></model-export-result>
        <bk-button class="done-export-button" theme="primary" @click="doneExport">{{t('完成')}}</bk-button>
      </div>
    </template>
  </step-pane>
</template>

<script>
  import { defineComponent, ref } from '@vue/composition-api'
  import { t } from '@/i18n'
  import StepPane from '../step-pane.vue'
  import UserNotice from '../user-notice.vue'
  import ModelExportEditor from './model-export-editor.vue'
  import ModelExportSetting from './model-export-setting.vue'
  import ModelExportResult from './model-export-result.vue'
  import { batchExport } from '@/service/model/import-export.js'
  import cloneDeep from 'lodash/cloneDeep'

  export default defineComponent({
    name: 'ModelExportPane',
    components: { StepPane, UserNotice, ModelExportEditor, ModelExportSetting, ModelExportResult },
    props: {
      /**
       * 已选择的需要导出的模型的数据
       */
      data: {
        type: Array,
        required: true,
        default: () => []
      }
    },
    setup({ data }, ctx) {
      const exportData = cloneDeep(data)
      const isDownloading = ref(false)

      const confirmedExportData = ref({})
      const exportSettingRef = ref(null)
      const exportSetting = ref({})

      const downloadExportData = () => {
        const params = {
          object_id: [],
          excluded_asst_id: [],
        }

        params.object_id = confirmedExportData.value.modelGroups
          .reduce((prev, next) => prev.concat(next.bk_objects.map(model => model.id)), [])
        params.excluded_asst_id = Object
          .values(confirmedExportData.value.excludedRelations)
          .reduce((prev, next) => prev.concat(next.map(item => item.id)), [])

        const { expirationTime, password, fileName } = exportSetting.value
        const setValue = (key, value) => {
          if (value !== undefined && value !== null && value !== '') {
            params[key] = value
          }
        }

        setValue('expiration', expirationTime)
        setValue('password', password)
        setValue('file_name', fileName)

        isDownloading.value = true
        return batchExport(params).finally(() => {
          isDownloading.value = false
        })
      }

      // 自动下载导出数据，在点击确认导出后默认自动下载导出数据
      const autoDownloadExportData = (done, stop) => {
        exportSettingRef.value.validate()
          .then(() => {
            downloadExportData()
              .then(() => {
                done()
              })
              .catch((err) => {
                stop()
                console.log(err)
              })
          })
          .catch((err) => {
            stop()
            console.log(err)
          })
      }

      const cancelExport = () => {
        ctx.emit('cancel')
      }

      const doneExport = () => {
        ctx.emit('done')
      }

      const handleRetryDownload = () => {
        if (isDownloading.value) return
        downloadExportData()
      }

      const exportSteps = [
        { title: t('用户须知'), icon: 1, nextButtonText: t('我知道了'), },
        { title: t('导出内容确认'), icon: 2 },
        { title: t('导出设置'), icon: 3, nextHandler: autoDownloadExportData, nextButtonText: t('确认导出') },
        { title: t('开始导出'), icon: 4, nextButtonVisible: false }
      ]

      return {
        exportData,
        exportSteps,
        confirmedExportData,
        exportSetting,
        exportSettingRef,
        cancelExport,
        doneExport,
        downloadExportData,
        handleRetryDownload,
        t
      }
    }
  })
</script>
<style lang="scss" scoped>
.done-export-button {
  display: block;
  width: 120px;
  margin: 20px auto 0;
}
</style>
