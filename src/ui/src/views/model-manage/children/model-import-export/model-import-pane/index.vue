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
  <div class="model-import-pane">
    <step-pane @cancel="cancelImport" :steps="importSteps">
      <template #default="{ currentStepIndex }">
        <div class="step-container" v-show="currentStepIndex === 1">
          <user-notice>
            <pre style="white-space: break-spaces;">{{t('模型导入导出用户须知')}}</pre>
          </user-notice>
        </div>
        <div class="step-container" v-show="currentStepIndex === 2">
          <model-package-upload @unzip="handlePkgUnzip"></model-package-upload>
        </div>
        <div class="step-container" v-show="currentStepIndex === 3">
          <model-import-editor :data="importData" v-model="confirmedImportData"></model-import-editor>
        </div>
        <div class="step-container" v-show="currentStepIndex === 4">
          <model-import-result :data="importResult"></model-import-result>
          <bk-button class="done-import-button" theme="primary" @click="doneImport">{{t('完成')}}</bk-button>
        </div>
      </template>
    </step-pane>
  </div>
</template>

<script>
  import { computed, defineComponent, ref } from 'vue'
  import { t } from '@/i18n'
  import UserNotice from '../user-notice.vue'
  import StepPane from '../step-pane.vue'
  import ModelPackageUpload from './model-package-upload.vue'
  import ModelImportEditor from './model-import-editor.vue'
  import ModelImportResult from './model-import-result.vue'
  import { batchImport } from '@/service/model/import-export.js'
  import isEmpty from 'lodash/isEmpty'

  export default defineComponent({
    name: 'ModelImportPane',
    components: { UserNotice, StepPane, ModelPackageUpload, ModelImportEditor, ModelImportResult },
    setup(props, ctx) {
      const importData = ref({}) // 从文件导入的数据
      const confirmedImportData = ref({
        confirmedModels: [],
        confirmedRelationTypes: []
      })
      const importResult = computed(() => ({
        modelCount: confirmedImportData.value.confirmedModels?.length || 0,
        relationTypeCount: confirmedImportData.value.confirmedRelationTypes?.length || 0
      }))


      const cancelImport = () => {
        ctx.emit('cancel')
      }

      const doneImport = () => {
        ctx.emit('done')
      }

      // 从文件导入的数据是否存在
      const isImportDataExisted = (done) => {
        if (!isEmpty(importData.value)) {
          done()
        }
      }

      const handlePkgUnzip = (data) => {
        importData.value = data
      }

      const importSteps = [
        {
          title: t('用户须知'),
          icon: 1,
          nextButtonText: t('我知道了'),
        },
        {
          title: t('文件包上传'),
          icon: 2,
          nextHandler: isImportDataExisted,
          nextButtonDisabled: () => isEmpty(importData.value)
        },
        {
          title: t('导入内容确认'),
          icon: 3,
          nextButtonText: t('确认导入'),
          nextButtonDisabledTooltips: t('没有选择任何导入模型或模型全部已在目标环境存在'),
          nextButtonDisabled: () => confirmedImportData.value.confirmedModels?.length === 0,
          nextHandler: (done, stop) => {
            batchImport({
              import_object: confirmedImportData.value.confirmedModels,
              import_asst: confirmedImportData.value.confirmedRelationTypes
            })
              .then(() => {
                done()
              })
              .catch((err) => {
                stop()
                console.log(err)
              })
          },
        },
        { title: t('开始导入'), icon: 4, nextButtonVisible: false }
      ]

      return {
        importResult,
        importSteps,
        importData,
        confirmedImportData,
        cancelImport,
        doneImport,
        handlePkgUnzip,
        t
      }
    }
  })
</script>

<style lang="scss" scoped>
.done-import-button {
  display: block;
  margin: 20px auto 0;
  width: 120px;
}
</style>
