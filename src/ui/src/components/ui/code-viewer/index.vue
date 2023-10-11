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

<script setup>
  import { onMounted, getCurrentInstance, ref } from 'vue'
  import Prism from 'prismjs'
  import { t } from '@/i18n'
  import { downloadFile } from '@/utils/util'
  import { $error } from '@/magicbox/index.js'
  import '@/assets/prism-vsc-dark-plus.css'

  Prism.manual = true

  const props = defineProps({
    code: {
      type: String,
      default: ''
    },
    lang: {
      type: String,
      default: 'json'
    },
    filename: {
      type: String,
      default: ''
    }
  })

  const { proxy } = getCurrentInstance()

  const languages = {
    json: 'JSON',
    javascript: 'JavaScript'
  }

  const showCopyTips = ref(false)
  const handleCopy = () => {
    proxy.$copyText(props.code).then(() => {
      showCopyTips.value = true
      const timer = setTimeout(() => {
        showCopyTips.value = false
        clearTimeout(timer)
      }, 500)
    }, () => {
      $error(t('复制失败'))
    })
  }

  const handleDownload = () => downloadFile(props.code, props.filename)

  onMounted(() => {
    Prism.highlightAll()
  })
</script>

<template>
  <div class="code-viewer">
    <div class="toolbar">
      <div class="title">{{ languages[lang] || lang }}</div>
      <div class="extra">
        <div class="toolbar-item">
          <i
            class="icon-cc-copy"
            @click="handleCopy">
          </i>
          <transition name="fade">
            <span class="copy-tips"
              :style="{ width: $i18n.locale === 'en' ? '100px' : '70px' }"
              v-show="showCopyTips">
              {{$t('复制成功')}}
            </span>
          </transition>
        </div>
        <div class="toolbar-item">
          <i
            class="icon-cc-download"
            @click="handleDownload">
          </i>
        </div>
      </div>
    </div>
    <pre class="code-content line-numbers"><code :class="`language-${lang}`">{{ code }}</code></pre>
  </div>
</template>

<style lang="scss" scoped>
.code-viewer {
  --toolbar-height: 40px;
  height: 100%;

  .toolbar {
    display: flex;
    align-items: center;
    height: var(--toolbar-height);
    background: #242424;
    border-bottom: 1px solid #0a0a0a;
    padding: 0 1em;

    .title {
      font-size: 14px;
      color: #C4C6CC;
    }
    .extra {
      margin-left: auto;
      display: flex;
      gap: 12px;
    }
    .toolbar-item {
      color: #979BA5;
      cursor: pointer;
      position: relative;

      .copy-tips {
        position: absolute;
        top: -24px;
        left: -24px;
        min-width: 70px;
        height: 26px;
        line-height: 26px;
        font-size: 12px;
        color: #ffffff;
        text-align: center;
        background-color: #000;
        border-radius: 2px;
      }
      .fade-enter-active, .fade-leave-active {
        transition: all 0.5s;
      }
      .fade-enter {
        top: -14px;
        opacity: 0;
      }
      .fade-leave-to {
        top: -28px;
        opacity: 0;
      }
    }
  }

  .code-content {
    height: calc(100% - var(--toolbar-height));
    margin: 0;
  }
}
</style>
