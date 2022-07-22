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
  <div class="import-file">
    <cmdb-tips class="file-tips"
      :icon-style="{ color: '#FF9C01' }"
      :tips-style="{ background: '#fff4e2', border: '1px solid #ffdfac' }">
      {{$t('导入更新提示')}}
    </cmdb-tips>
    <div class="file-trigger" v-if="!file">
      <input type="file"
        accept="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
        @change.prevent="changeFile" />
      <i class="trigger-icon bk-icon icon-upload-cloud"></i>
      <i18n class="trigger-text" path="导入文件拖拽提示">
        <template #clickUpload><span class="trigger-text-link">{{$t('点击上传')}}</span></template>
      </i18n>
    </div>
    <div class="file-info" v-else>
      <i class="file-icon icon-cc-excel"></i>
      <span class="file-name">{{file.name}}</span>
      <span class="file-size">{{formatSize(file.size)}}</span>
      <i class="file-delete bk-icon icon-close" @click="clearFile"></i>
    </div>
    <p class="size-tips">
      {{importState.fileTips || $t('导入文件大小提示')}}
      <a href="javascript:void(0);" class="file-template"
        v-if="importState.template"
        @click="handleDownloadTemplate">
        <img src="../../assets/images/icon/down_model_icon.png">
        {{$t('下载模板')}}
      </a>
    </p>
    <div class="options">
      <bk-button theme="primary" :loading="pending" :disabled="!file" @click="handleNextStep">{{$t('下一步')}}</bk-button>
      <bk-button theme="default" class="ml10" @click="closeImport">{{$t('取消')}}</bk-button>
    </div>
  </div>
</template>

<script>
  import useImport from './index'
  import useStep from './step'
  import useFile from './file'
  import { computed } from '@vue/composition-api'
  export default {
    name: 'import-file',
    setup() {
      const [currentStep, { next: nextStep }] = useStep()
      const [importState, { close: closeImport }] = useImport()
      const [{ file, state }, {
        change: changeFile,
        clear: clearFile,
        setState: setFileState,
        setError: setFileError,
        setResponse: setFileResponse
      }] = useFile()
      const pending = computed(() => state.value === 'resolving')
      return {
        currentStep,
        nextStep,
        importState,
        closeImport,
        file,
        pending,
        changeFile,
        clearFile,
        setFileState,
        setFileError,
        setFileResponse
      }
    },
    methods: {
      async handleNextStep() {
        try {
          this.setFileState('resolving')
          const response = await this.importState.submit({
            file: this.file,
            step: this.currentStep
          })
          this.setFileResponse(response)
          this.nextStep()
          this.setFileState(null)
        } catch (error) {
          console.error(error)
          this.setFileState('error')
          this.setFileError(error)
        }
      },
      handleDownloadTemplate() {
        this.$http.download({ url: this.importState.template })
      },
      formatSize(originalSize) {
        const uints = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
        const index = Math.floor(Math.log(originalSize) / Math.log(1024))
        let size = originalSize / (1024 ** index)
        size = `${size.toFixed(2)}${uints[index]}`
        return size
      }
    }
  }
</script>

<style lang="scss" scoped>
  .import-file {
    .file-tips {
      margin: 20px 0 0 0;
    }
  }
  .file-trigger {
    position: relative;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    height: 100px;
    background: #fafbfd;
    border: 1px dashed #c4c6cc;
    border-radius: 3px;
    margin: 10px 0 0 0;
    &:hover {
      border-color: $primaryColor;
      .trigger-icon {
        color: $primaryColor;
      }
    }
    input[type=file] {
      position: absolute;
      width: 100%;
      height: 100%;
      opacity: 0;
      z-index: 2;
      cursor: pointer;
    }
    .trigger-icon {
      font-size: 24px;
      color: #C4C6CC;
    }
    .trigger-text {
      font-size: 12px;
      color: #63656e;
      line-height: 16px;
      &-link {
        color: $primaryColor;
      }
    }
  }
  .file-info {
    display: flex;
    align-items: center;
    height: 60px;
    border: 1px solid #c4c6cc;
    border-radius: 3px;
    margin: 10px 0 0 0;
    &:hover {
      .file-delete {
        display: inline-block;
      }
    }
    .file-icon {
      font-size: 26px;
      color: #979ba5;
      margin: 0 13px 0 18px;
    }
    .file-name {
      flex: 1;
      height: 20px;
      font-size: 12px;
      font-weight: 700;
      line-height: 20px;
    }
    .file-size {
      font-size: 12px;
      line-height: 16px;
      margin: 0 20px;
    }
    .file-delete {
      display: none;
      font-size: 20px;
      cursor: pointer;
      margin-right: 10px;
      &:hover {
        opacity: .7;
      }
    }
  }
  .size-tips {
    display: flex;
    align-items: center;
    font-size: 12px;
    line-height: 16px;
    margin: 10px 0 0 0;
    .file-template {
      color: $primaryColor;
      display: inline-flex;
      align-items: center;
      margin-left: 10px;
      img {
        width: 14px;
      }
      &:hover {
        opacity: .7;
      }
    }
  }
  .options {
    display: flex;
    margin: 12px 0 0 0;
  }
</style>
