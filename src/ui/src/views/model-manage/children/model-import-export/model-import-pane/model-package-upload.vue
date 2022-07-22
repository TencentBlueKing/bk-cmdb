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
  <div class="model-package-upload">
    <div class="file-uploader-container model-package-uploader">
      <div class="file-uploader">
        <label
          v-show="!currentFile.file"
          class="upload-area"
          :class="{ 'is-active': isDragging }"
          @dragenter.stop.prevent="handleFileDragenter"
          @dragover.stop.prevent
          @dragleave.prevent="handleFileDragleave"
          @drop.stop.prevent="handleFileDrop"
        >
          <input
            class="upload-input"
            type="file"
            ref="fileInputRef"
            accept="application/zip"
            @change="handleFileChange"
          />
          <i class="upload-icon">
            <bk-icon type="upload-cloud"></bk-icon>
          </i>
          <p class="upload-text">{{ t("将文件拖到此处或点击上传") }}</p>
        </label>

        <div
          class="after-upload"
          :class="[{ [`is-${currentFile.result}`]: currentFile.result !== '' }]"
          v-show="currentFile.file"
        >
          <div class="file-icon">
            <img src="@/assets/images/zip-package.png" />
          </div>
          <div class="file-description">
            <p class="file-name">{{ currentFileName }}</p>
            <p class="upload-result-desc">
              {{ currentFile.resultDesc }}
              <span
                class="retry-decrypt-button"
                v-show="isDecryptPasswordError"
                @click="showDecryptDialog"
              >
                {{ t("重新输入密码") }}
              </span>
            </p>
          </div>
          <div class="file-action">
            <i @click="preProcessFile(currentFile.file)" class="icon-button retry-upload-button icon-cc-refresh"></i>
            <i @click="clearFile" class="icon-button clear-file-button icon-cc-tips-close"></i>
          </div>
        </div>
      </div>
      <p class="upload-tips">{{t('仅支持上传来自蓝鲸配置平台专属导出的模型压缩包')}}</p>
    </div>
    <bk-dialog
      header-position="left"
      :title="t('文件包解密确认')"
      :confirm-fn="confirmDecrypt"
      @cancel="cancelDecrypt"
      v-model="decryptDialogVisible"
    >
      <bk-form
        form-type="vertical"
        ref="decryptPasswordFormRef"
        :model="{ decryptPassword }"
        :rules="{ decryptPassword: decryptPasswordRules }"
      >
        <bk-form-item property="decryptPassword" :label="t('文件包密码')">
          <bk-input
            v-autofocus
            @enter="confirmDecrypt"
            type="password"
            v-model="decryptPassword"
            :placeholder="t('请输入上传的文件包密码，完成后点提交验证')"
          ></bk-input>
        </bk-form-item>
      </bk-form>
    </bk-dialog>
  </div>
</template>

<script>
  import { defineComponent, ref, reactive, computed } from '@vue/composition-api'
  import unzip from 'unzip-js'
  import { t } from '@/i18n'
  import { batchImportFileAnalysis } from '@/service/model/import-export.js'
  import { autofocus } from '@/directives/autofocus'

  export default defineComponent({
    name: 'ModelPackageUpload',
    directives: {
      autofocus
    },
    setup(props, { emit }) {
      const isUnzipping = ref(false)
      const isDragging = ref(false)
      const currentFile = reactive({
        file: null,
        result: '',
        resultDesc: ''
      })
      const fileInputRef = ref(null)
      const currentFileName = computed(() => currentFile?.file?.name || '')
      const decryptDialogVisible = ref(false)
      const decryptPasswordFormRef = ref(null)
      const decryptPassword = ref('')
      const isDecryptPasswordError = ref(false)
      const decryptPasswordRules = [
        {
          required: true,
          message: t('请输入文件包密码'),
          trigger: 'blur'
        },
        {
          required: true,
          message: t('密码输入错误'),
          validator: (value) => {
            const lenValid = value.length >= 6 && value.length <= 20
            return lenValid
          },
          trigger: 'blur'
        }
      ]

      const confirmDecrypt = () => {
        decryptPasswordFormRef.value.validate((isPass) => {
          if (isPass) {
            isDecryptPasswordError.value = false
            unzipFile(decryptPassword.value)
              .then(() => {
                closeDecryptDialog()
              })
          }
        })
      }

      const uploadFile = () => {
        fileInputRef.value.click()
      }

      const clearFile = () => {
        fileInputRef.value.value = null
        currentFile.file = ''
        currentFile.result = ''
        currentFile.resultDesc = ''
        emit('unzip', {})
      }

      const closeDecryptDialog = () => {
        decryptPassword.value = ''
        decryptDialogVisible.value = false
      }

      const cancelDecrypt = () => {
        closeDecryptDialog()
        currentFile.result = 'failed'
        currentFile.resultDesc = t('您取消了输入密码')
        isDecryptPasswordError.value = true
      }

      const unzipFile = (password = '') => {
        isUnzipping.value = true

        const data = {
          file: currentFile.file
        }

        if (password) {
          data.password = password
        }

        return batchImportFileAnalysis(data)
          .then((data) => {
            emit('unzip', data)
            currentFile.result = 'succeed'
            currentFile.resultDesc = t('上传成功')
          })
          .catch((err) => {
            isDecryptPasswordError.value = err.bk_error_code === 1111023
            currentFile.result = 'failed'
            currentFile.resultDesc = `${t('上传失败')}：${err.bk_error_msg || err.message}`
          })
          .finally(() => {
            isUnzipping.value = false
          })
      }

      const showDecryptDialog = () => {
        decryptDialogVisible.value = true
      }

      const preProcessFile = (file) => {
        if (!file) return

        isDecryptPasswordError.value = false

        currentFile.file = file

        isFileEncrypted(file, (isEncrypted) => {
          if (isEncrypted) {
            showDecryptDialog()
          } else {
            unzipFile()
          }
        })
      }

      const isFileEncrypted = (file, callback) => {
        let isEncrypted = false

        unzip(file, (err, zipFile) => {
          zipFile.readEntries((err, entries) => {
            if (err) {
              return console.error(err)
            }

            isEncrypted = entries.some(entry => entry.encrypted)

            callback(isEncrypted)
          })
        })
      }

      const handleFileDragenter = () => {
        isDragging.value = true
      }

      const handleFileDragleave = () => {
        isDragging.value = false
      }

      const handleFileDrop = (e) => {
        const dt = e.dataTransfer
        const {
          files: [file]
        } = dt

        isDragging.value = false

        preProcessFile(file)
      }

      const handleFileChange = (e) => {
        const [file] = e.target.files

        preProcessFile(file)
      }

      return {
        t,
        fileInputRef,
        preProcessFile,
        uploadFile,
        clearFile,
        currentFile,
        currentFileName,
        decryptDialogVisible,
        decryptPasswordFormRef,
        decryptPassword,
        isDecryptPasswordError,
        decryptPasswordRules,
        handleFileChange,
        isDragging,
        handleFileDragenter,
        handleFileDragleave,
        handleFileDrop,
        confirmDecrypt,
        cancelDecrypt,
        showDecryptDialog,
        isUnzipping
      }
    }
  })
</script>

<style lang="scss" scoped>
.model-package-uploader {
  width: 680px;
  margin: 0 auto;
}

.file-uploader {
  margin: 50px auto 0;
  height: 80px;
}

.upload-area {
  display: block;
  height: 100%;
  display: flex;
  align-items: center;
  flex-direction: column;
  justify-content: center;
  background-color: #fafbfd;
  border: 1px dashed #c4c6cc;
  border-radius: 2px;
  cursor: pointer;

  &:hover,
  &.is-active {
    border-color: #3a84ff;

    .upload-icon {
      color: #3a84ff;
    }
  }

  .upload-icon {
    font-size: 28px;
    color: #979ba5;
    pointer-events: none;
  }

  .upload-text {
    margin-top: 5px;
    font-size: 12px;
    pointer-events: none;
  }

  .upload-input {
    display: none;
  }
}

.after-upload {
  display: flex;
  height: 100%;
  align-items: center;
  border: 1px solid #c4c6cc;
  border-radius: 2px;

  &.is-failed {
    border-color: $cmdbDangerColor;
    background-color: #fff1f1;

    .upload-result-desc,
    .icon-button {
      color: $cmdbDangerColor;
    }

    .file-action {
      display: flex;
    }
  }

  &.is-succeed {
    background-color: #fff;

    .upload-result-desc {
      color: $cmdbSuccessColor;
    }

    &:hover .file-action {
      display: flex;

      .retry-upload-button {
        display: none;
      }
    }
  }

  .file-icon {
    width: 48px;
    height: 48px;
    margin-left: 16px;

    > img {
      width: 100%;
    }
  }

  .file-description {
    margin-left: 12px;
    font-size: 14px;
  }

  .file-name {
    font-size: 14px;
  }

  .upload-result-desc {
    margin-top: 8px;
    margin-right: 20px;
    font-size: 14px;
  }

  .retry-decrypt-button {
    display: inline-block;
    vertical-align: baseline;
    cursor: pointer;
    color: $primaryColor;
  }

  .file-action {
    display: none;
    align-items: center;
    margin-left: auto;
    font-size: 14px;

    .icon-button {
      display: inline-block;
      padding: 4px;
      cursor: pointer;
      margin-right: 16px;
    }

    .clear-file-button {
      font-size: 12px;
    }
  }
}

.upload-tips {
  font-size: 12px;
  margin-top: 5px;
}
</style>
