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
  <step-pane @cancel="cancelExport" :steps="exportSteps">
    <template #default="{ currentStepIndex }">
      <div class="step-container" v-show="currentStepIndex === 1">
        <user-notice>
          <pre style="white-space: break-spaces;">{{t('模型导入导出用户须知')}}</pre>
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
  import { defineComponent, ref } from 'vue'
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
