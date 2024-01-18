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
  <div class="editable-block">
    <div class="search-input"
      ref="searchInput"
      :placeholder="placeholder"
      :focusTip="focusTip"
      :blurTip="blurTip"
      contenteditable="plaintext-only"
      spellcheck="false"
      @blur="handleBlur"
      @keydown.enter="handleSearch"
      @focus="handleFocus"
      @input="handleInput"
      @paste="handlePaste">
    </div>
    <i class="search-close bk-icon icon-close-circle-shape" @mousedown="handleClear" v-if="searchContent"></i>
  </div>
</template>

<script setup>
  import { reactive, ref, onMounted } from 'vue'
  import { LT_REGEXP, ALL_PROBABLY_IP, AREA_IPV6_IP, AREA_IPV4_IP, IPV6_IP, IPV4_IP } from '@/dictionary/regexp'
  import { getCursorPosition, setCursorPosition } from '@/utils/util'

  const props = defineProps({
    value: {
      type: String,
      default: ''
    },
    placeholder: {
      type: String,
      default: ''
    },
    focusTip: {
      type: String,
      default: ''
    },
    blurTip: {
      type: String,
      default: ''
    },
    enterSearch: {
      type: Boolean,
      default: true
    },
    noBlurClass: {
      type: String,
      default: 'search-btn'
    },
    blurParse: {
      type: Boolean,
      default: true
    }
  })

  const emit = defineEmits(['keydown', 'focus', 'blur', 'updateValue'])

  onMounted(() => {
    setInputHtml(searchContent.value)
  })

  const searchInput = ref(null)
  const searchContent = ref(props.value)
  const pasteData = reactive({
    cursor: 0,
    length: 0,
    input: false
  })

  const initPasteData = () => {
    pasteData.cursor = 0
    pasteData.length = 0
  }
  const getMatchIP = (ip) => {
    const reg = []
    const match = []

    const matchedAreaV6 = ip.match(AREA_IPV6_IP)
    if (matchedAreaV6) {
      match.push(...matchedAreaV6)
      reg.push(AREA_IPV6_IP)
    }

    const matchedAreaV4 = ip.match(AREA_IPV4_IP)
    if (matchedAreaV4) {
      match.push(...matchedAreaV4)
      reg.push(AREA_IPV4_IP)
    }

    const matchedV6 = ip.match(IPV6_IP)
    if (matchedV6 && !reg.includes(AREA_IPV6_IP)) {
      match.push(...matchedV6)
      reg.push(IPV6_IP)
    }

    const matchedV4 = ip.match(IPV4_IP)
    if (matchedV4 && !reg.includes(AREA_IPV4_IP)) {
      match.push(...matchedV4)
      reg.push(IPV4_IP)
    }

    return [reg, ip, match]
  }

  const replaceRealIP = (allReg, ip) => {
    allReg.forEach(reg => ip = ip.replace(reg, '<span class="high-light">$<ip></span>'))
    return ip
  }
  const getNewHtml = content => content.replace(ALL_PROBABLY_IP, (val) => {
    // val为所有可能是IP的数据，在这里筛选出符合条件的IP
    const [reg, ip, matched] = getMatchIP(val)
    return matched[0] ? replaceRealIP(reg, ip) : val
  })

  const setSearchContent = (val = '') => {
    searchContent.value = val
    emit('updateValue', val)
  }
  const setInputHtml = (html = '') => {
    searchInput.value.innerHTML = html
  }
  // 处理数据高光
  const setHighlight = () => {
    const content = searchContent.value
    const cursor = getCursorPosition(searchInput.value)
    pasteData.cursor = cursor
    setInputHtml(getNewHtml(content.replace(LT_REGEXP, '&lt')))
    setCursorPosition(searchInput.value, cursor)
  }

  const parseIP = () => {
    initPasteData()
    const propablyIP = searchContent.value.match(ALL_PROBABLY_IP)
    const ipList = new Set()
    propablyIP?.forEach((ip) => {
      const [, , matched] = getMatchIP(ip)
      if (matched[0]) {
        matched.forEach(ip => ipList.add(ip))
      }
    })

    // 如果一个IP都没有并且blurParse为false，则内容不解析，保持原状
    if (!(ipList.size || props.blurParse)) {
      return
    }
    const newHtml = Array.from(ipList).join('\n')
    setSearchContent(newHtml)
    setInputHtml(newHtml)
  }

  const handlePaste = (event) => {
    const val = event?.clipboardData?.getData('text')?.replace(/\r/g, '')
    pasteData.length = val.length
  }
  const handleBlur = (event) => {
    parseIP()
    emit('blur', event)
  }
  const handleInput = () => {
    const { innerText } = searchInput.value
    setSearchContent(innerText)
    setHighlight()
    setTimeout(() => {
      const { scrollHeight, clientHeight, scrollTop } = searchInput.value
      const bottom = scrollHeight - clientHeight - scrollTop
      // 防止光标被遮挡
      if (bottom > 0 && bottom < 30) {
        searchInput.value.scrollTop = scrollHeight - 10
      }
    }, 0)
  }
  const handleClear = (event) => {
    initPasteData()
    setSearchContent()
    setInputHtml()
    event?.preventDefault()
  }
  const handleFocus = (event) => {
    setHighlight()
    emit('focus', event)
  }
  const handleSearch = (event) => {
    if (!props.enterSearch) return
    const { shiftKey, metaKey, ctrlKey } = event
    const agent = window.navigator.userAgent.toLowerCase()
    const isMac = /macintosh|mac os x/i.test(agent)
    const modifierKey = isMac ? metaKey : ctrlKey
    if (!modifierKey && !shiftKey) {
      parseIP()
      emit('search')
      event.preventDefault()
    }
  }

  const focus = () => {
    searchInput.value.focus()
  }
  defineExpose({
    focus,
    searchContent,
    clear: handleClear
  })
</script>

<style lang="scss" scoped>
@mixin tip {
  color: #C4C6CC;
  position: sticky;
  left: 0px;
  bottom: 0px;
  font-size: 12px;
  line-height: 17px;
  display: block;
  width: 100%;
  background: white;
  padding-bottom: 4px;
}
.editable-block {
  flex: 1;
  max-width: 646px;
  position: relative;
}
.search-input[contenteditable]:empty::before {
  content: attr(placeholder);
  color: #C4C6CC;
  cursor: text;
  font-size: 12px;
  position: absolute;
  left: 16px;
}
.search-input[contenteditable]:focus {
  border-color: #3A84FF;
}
.search-input[contenteditable]:not(:empty):focus::after {
  content: attr(focusTip);
  @include tip;
}
.search-input[contenteditable]:not(:empty):not(:focus)::after {
  content: attr(blurTip);
  @include tip;
}
.search-input {
  max-width: 646px;
  background: white;
  min-height: 100%;
  height: max-content;
  padding: 5px 32px 0px 16px;
  border: 1px solid #c4c6cc;
  border-radius: 0 0 0 2px;
  font-size: 14px;
  line-height: 30px;
  max-height: 400px;
  position: relative;
  @include scrollbar-y;
  /deep/ {
    .bk-textarea-wrapper {
      border: 0;
      border-radius: 0 0 0 2px;
    }
    .bk-form-textarea {
      min-height: 42px;
      line-height: 30px;
      font-size: 14px;
      border: 1px solid #C4C6CC;
      padding: 5px 32px 5px 16px;
      border-radius: 0 0 0 2px;
    }
    .right-icon {
      right: 20px !important;
    }
    .high-light {
      display: inline-block;
      background: #FFE8C3;
      border-radius: 2px;
    }
  }
}
.search-close {
  position: absolute;
  right: 10px;
  top: 13px;
  cursor: pointer;
  color: #c4c6cc;
}
</style>
