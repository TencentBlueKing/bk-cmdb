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
  import { LT_REGEXP, ALL_PROBABLY_IP } from '@/dictionary/regexp'
  import isIP from 'validator/es/lib/isIP'
  import isInt from 'validator/es/lib/isInt'

  const props = defineProps({
    searchContent: {
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
    }
  })

  const emit = defineEmits(['keydown', 'focus', 'blur'])

  onMounted(() => {
    setInputHtml(searchContent.value)
  })

  const searchInput = ref(null)
  const searchContent = ref(props.searchContent)
  const pasteData = reactive({
    cursor: 0,
    length: 0,
    input: false
  })

  const initPasteData = () => {
    pasteData.cursor = 0
    pasteData.length = 0
  }

  const getMatchIp = (ip) => {
    // 先判断是不是常规Ip
    if (isIP(ip)) {
      return ip
    }
    const matchedV4 = ip.split(':')
    const matchedV6 = ip.match(/^(\d+):\[([0-9a-fA-F:.]+)\]$/)
    if (matchedV4.length === 2 && isInt(matchedV4[0]) && isIP(matchedV4[1], 4)) return ip
    if (matchedV6 && isIP(matchedV6[2])) return ip
    return null
  }
  // 获取光标位置
  const getCursorPosition = () => {
    const selection = window.getSelection()
    const element = searchInput.value
    let caretOffset = 0
    // false表示进行了范围选择
    const { isCollapsed } = selection
    // 选中的区域
    if (selection.rangeCount > 0) {
      const range = selection.getRangeAt(0)
      // 克隆一个选中区域
      const preCaretRange = range.cloneRange()
      // 设置选中区域的节点内容为当前节点
      preCaretRange.selectNodeContents(element)
      // 重置选中区域的结束位置
      preCaretRange.setEnd(range.endContainer, range.endOffset)
      const { length } = preCaretRange.toString()
      caretOffset = isCollapsed ? length : length - selection.toString().length
    }
    return caretOffset
  }

  // 设置光标位置
  const setCursorPostion = () => {
    const selection = window.getSelection()
    const parent = document.getElementsByClassName('search-input')[0]
    const child = parent.getElementsByClassName('new-data')[0]
    // 创建一个选中区域
    const range = document.createRange()
    // 选中节点的内容
    range.selectNodeContents(child || parent)
    if (child?.innerHTML?.length > 0) {
      // 粘贴的直接设置光标到粘贴的末尾
      range.setStart(child, child.childNodes.length)
    } else {
      // 非粘贴的数据通过计算文本节点的偏移量来设置光标
      const parentAllNodes = parent.childNodes
      let { cursor } = pasteData
      for (let i = 0;i < parentAllNodes.length;i++) {
        const nowNode = parentAllNodes[i]
        const nodeLength = nowNode?.length ?? nowNode?.innerText?.length
        if (cursor <= nodeLength) {
          range.setStart(nowNode?.firstChild || nowNode, cursor)
          break
        }
        cursor -= nodeLength
      }
    }
    // 设置选中区域为一个点
    range.collapse(true)
    // 移除所有的选中范围
    selection.removeAllRanges()
    // 添加新建的范围
    selection.addRange(range)
  }
  const setSearchContent = (val = '') => {
    searchContent.value = val
    emit('update:search-content', val)
  }
  const setInputHtml = (html = '') => {
    searchInput.value.innerHTML = html
  }
  // 处理数据高光
  const setHighLight = (inputType) => {
    const content = searchContent.value
    if (inputType === 'insertFromPaste') {
      const { cursor, length } = pasteData
      setHighLightPaste(cursor, length, content)
    } else {
      const cursor = getCursorPosition()
      pasteData.cursor = cursor
      setHighLightOther(content)
    }
    setCursorPostion()
  }
  // 粘贴高亮
  const setHighLightPaste = (cursor, length, content) => {
    const start = content.substring(0, cursor).replace(LT_REGEXP, '&lt')
    const end = content.substring(cursor + length).replace(LT_REGEXP, '&lt')
    const paste = `<span class="new-data">${content.substring(cursor, cursor + length).replace(LT_REGEXP, '&lt')}</span>`
    const newHtml = (start + paste + end).replace(ALL_PROBABLY_IP, (val) => {
      // val为所有可能是IP的数据，在这里筛选出符合条件的IP
      const ans = getMatchIp(val)
      return ans ? `<span class="high-light">${val}</span>` : val
    })
    setInputHtml(newHtml)
  }
  // 非粘贴高亮
  const setHighLightOther = (content) => {
    const newHtml = content.replace(LT_REGEXP, '&lt').replace(ALL_PROBABLY_IP, (val) => {
      // val为所有可能是IP的数据，在这里筛选出符合条件的IP
      const ans = getMatchIp(val)
      return ans ? `<span class="high-light">${val}</span>` : val
    })
    setInputHtml(newHtml)
  }

  const parseIp = () => {
    initPasteData()
    const propablyIp = searchContent.value.match(ALL_PROBABLY_IP)
    const ipList = []
    propablyIp?.forEach((ip) => {
      if (getMatchIp(ip)) {
        ipList.push(ip)
      }
    })
    const newHtml = ipList.join('\n')
    setSearchContent(newHtml)
    setInputHtml(newHtml)
  }

  const handlePaste = (event) => {
    const val = event?.clipboardData?.getData('text')?.replace(/\r/g, '')
    pasteData.length = val.length
    pasteData.cursor = getCursorPosition()
  }
  const handleBlur = (event) => {
    parseIp()
    searchInput.value.blur()
    emit('blur', event)
  }
  const handleInput = (event) => {
    const { inputType } = event
    const { innerText } = searchInput.value
    setSearchContent(innerText)
    setHighLight(inputType)
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
    event.preventDefault()
  }
  const handleFocus = (event) => {
    setHighLight('focus')
    emit('focus', event)
  }
  const handleSearch = (event) => {
    if (!props.enterSearch) return
    const { shiftKey, metaKey, ctrlKey } = event
    const agent = window.navigator.userAgent.toLowerCase()
    const isMac = /macintosh|mac os x/i.test(agent)
    const modifierKey = isMac ? metaKey : ctrlKey
    if (!modifierKey && !shiftKey) {
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
    parseIp
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
      cursor: pointer;
      &:hover {
        background: #FFD695;
      }
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
