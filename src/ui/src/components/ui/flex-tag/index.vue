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
  import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
  import { $bkPopover } from '@/magicbox/index.js'

  const props = defineProps({
    list: {
      type: Array,
      default: () => ([])
    },
    isLinkStyle: {
      type: Boolean,
      default: false
    },
    isTextStyle: {
      type: Boolean,
      default: false
    },
    isTagStyle: {
      type: Boolean,
      default: false
    },
    maxWidth: {
      type: String,
      default: '400px'
    },
    fontSize: {
      type: String,
      default: '12px'
    },
    gap: {
      type: String,
      default: '4px'
    },
    height: {
      type: String,
      default: '22px'
    },
    popoverOptions: {
      type: Object,
      default: () => ({})
    },
    popoverMaxHeight: {
      type: String,
      default: '280px'
    },
    forceShowOne: {
      type: Boolean,
      default: false
    }
  })

  const emit = defineEmits(['click'])

  const containerEl = ref(null)
  const plusEl = ref(null)
  let tips = null

  const maxResizeCount = 10
  let execResizeTimes = 0

  const tags = computed(() => props.list.filter(item => item))
  const gapWidth = computed(() => parseInt(props.gap, 10))
  const ellipsisCount = ref(0)

  const tipTagList = ref([])
  const tipTagOffsetIndex = ref(0)
  const tagItemList = ref([])

  const handleClick = (index) => {
    emit('click', index)
  }
  const handleClickText = (tag) => {
    emit('click-text', tag)
  }

  const getTipsInstance = () => {
    if (!tips) {
      tips = $bkPopover(plusEl.value, {
        allowHTML: true,
        placement: 'top',
        boundary: 'window',
        arrow: true,
        theme: `${(props.isLinkStyle || props.isTagStyle) ? 'light' : 'dark'} flex-tag-tooltip`,
        interactive: true,
        ...props.popoverOptions,
        onShow(inst) {
          const contentEl = document.createElement('div')
          const fragment = document.createDocumentFragment()
          contentEl.classList.add('flex-tag-tips-content')
          if (props.isLinkStyle) {
            contentEl.classList.add('is-link')
          }
          if (props.isTagStyle) {
            contentEl.classList.add('is-tag')
          }
          contentEl.style.setProperty('--fontSize', props.fontSize)
          contentEl.style.setProperty('--height', props.height)
          contentEl.style.setProperty('--popoverMaxHeight', props.popoverMaxHeight)

          tipTagList.value.forEach((text, index) => {
            if (props.popoverOptions.appendTo === 'parent') {
              const originalTag = tagItemList.value[index + tipTagOffsetIndex.value]
              // clone出来占住原来的位置，保证重新计算时正确性
              const cloneEl = originalTag.el.cloneNode(true)
              cloneEl.classList.add('clone')
              containerEl.value.appendChild(cloneEl)
              // 移动元素使其成为tooltips的内容，这样能够复用所有的内容样式和事件
              fragment.appendChild(originalTag.el)
            } else {
              // 手工创建
              const itemEl = document.createElement('div')
              itemEl.classList.add('flex-tag-tips-item')
              if (props.isLinkStyle) {
                itemEl.addEventListener('click', () => handleClick(index), false)
              }
              itemEl.textContent = text.name || text
              fragment.appendChild(itemEl)
            }
          })

          contentEl.appendChild(fragment)
          inst.setContent(contentEl)
        },
        onHide(inst) {
          // 将元素替换回去
          if (props.popoverOptions.appendTo === 'parent') {
            tipTagList.value.forEach(() => {
              const cloneTagItems = Array.from(containerEl.value?.querySelectorAll('.tag-item.clone'))
              const popverTagItems = Array.from(inst.popperChildren.content?.querySelectorAll('.tag-item'))
              containerEl.value.replaceChild(popverTagItems.shift(), cloneTagItems.shift())
            })
          }
        }
      })
    }
    return tips
  }

  const resizeHander = () => {
    if (!tags.value.length || execResizeTimes > maxResizeCount) {
      return
    }

    const containerClientWidth = containerEl.value.clientWidth
    const containerScrollWidth = containerEl.value.scrollWidth
    const plusWidth = Math.ceil(plusEl.value.getBoundingClientRect().width)
    const tagItems = Array.from(containerEl.value.querySelectorAll('.tag-item'))

    const tagWidthList = tagItems.map((item, index) => ({
      index,
      el: item,
      width: Math.ceil(item.getBoundingClientRect().width)
    }))

    if (containerScrollWidth > containerClientWidth) {
      // 默认将plus元素宽度计入，这将始终保持plus元素是可见的，但对是否overflow的精度有所降低
      let accWidth = plusWidth

      // 出现断点的那个元素
      let posItem = null

      for (const item of tagWidthList) {
        accWidth = accWidth + item.width + gapWidth.value
        if (props.forceShowOne && item.index === 0) {
          continue
        }
        if (accWidth > containerClientWidth) {
          posItem = item
          ellipsisCount.value = tags.value.length - item.index
          break
        }
      }

      // 将plus元素放到断点元素前面并且显示
      if (posItem) {
        containerEl.value.insertBefore(plusEl.value, posItem.el)
        plusEl.value.classList.add('show')
        if (plusEl.value.previousElementSibling) {
          tagItems.forEach(item => item.classList.remove('is-pos'))
          plusEl.value.previousElementSibling.classList.add('is-pos')
        }
      }

      // 设置tooltips
      const tooltips = getTipsInstance()
      tooltips.setContent('')

      let offsetIndex = 0
      let tipTags = tags.value.slice(offsetIndex)
      if (props.isLinkStyle || props.isTextStyle) {
        offsetIndex = tags.value.length - ellipsisCount.value
        tipTags = tags.value.slice(offsetIndex)
      }

      tipTagList.value = tipTags
      tipTagOffsetIndex.value = offsetIndex
      tagItemList.value = tagWidthList
    } else {
      execResizeTimes = 0
      // 将plus元素放到最后并且隐藏
      containerEl.value.insertBefore(plusEl.value, tagItems[tagItems.length - 1].nextSibling)
      tagItems.forEach(item => item.classList.remove('is-pos'))
      plusEl.value.classList.remove('show')
    }

    execResizeTimes += 1
  }

  const changing = ref(false)

  const execResizeHander = async () => {
    // changing为true时使得所有tagitem可见才能保证宽度的计算准确
    changing.value = true
    await nextTick(resizeHander)
    changing.value = false
  }

  const resizeObserver = new ResizeObserver((entries) => {
    for (const entry of entries) {
      if (entry.target === containerEl.value) {
        execResizeHander()
      }
    }
  })

  watch(tags, () => {
    execResizeTimes = 0
    execResizeHander()
  })

  onMounted(() => {
    resizeObserver.observe(containerEl.value)
    execResizeHander()

    tips?.destroy?.()
  })

  onBeforeUnmount(() => {
    resizeObserver.unobserve(containerEl.value)
  })
</script>
<script>
  export default {
    name: 'flex-tag'
  }
</script>

<template>
  <ul
    :class="['flex-tag', {
      changing,
      'is-link-style': isLinkStyle,
      'is-text-style': isTextStyle,
      'is-tag-style': isTagStyle
    }]"
    ref="containerEl"
    :style="{
      '--fontSize': fontSize,
      '--gap': gap,
      '--maxWidth': maxWidth,
      '--height': height
    }">
    <li class="tag-item"
      v-for="(tag, index) in tags"
      :key="tag.id || index"
      @click="handleClick(index)">
      <div class="tag-item-text" v-bk-overflow-tips>
        <span @click="handleClickText(tag)">{{tag.name || tag}}</span>
        <slot name="text-append" v-bind="tag"></slot>
      </div>
      <slot name="append" v-bind="tag"></slot>
    </li>
    <li class="more-plus" ref="plusEl" v-show="ellipsisCount">+{{ellipsisCount}}</li>
  </ul>
</template>

<style lang="scss" scoped>
  .flex-tag {
    display: flex;
    gap: var(--gap);
    font-size: var(--fontSize);

    .tag-item {
      color: #63656E;
      background: #F0F1F5;
      border-radius: 2px;
      padding: 0 .6em;
      flex: none;
      height: var(--height);
      line-height: var(--height);
      max-width: var(--maxWidth);
      display: flex;
      align-items: center;

      .tag-item-text {
        width: 100%;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }
    }

    &.is-link-style {
      .tag-item {
        color: #3A84FF;
        background: none;
        border-radius: 0;
        padding: 0;
        cursor: pointer;

        &::after {
          content: ',';
          color: #63656E;
        }

        // 倒数第2个元素为tag-item的最后一个元素，倒数第1个元素为plus
        &:nth-last-of-type(2) {
          &::after {
            display: none;
          }
        }
        &.is-pos {
          &::after {
            content: '...';
          }
        }
      }
    }

    &.is-text-style {
      .tag-item {
        background: none;
        border-radius: 0;
        padding: 0;

        &::after {
          content: '|';
          color: #63656E;
        }

        &:nth-last-of-type(2) {
          &::after {
            display: none;
          }
        }
        &.is-pos {
          &::after {
            content: '...';
          }
        }
      }
    }

    .more-plus {
      display: none;
      align-items: center;
      justify-content: center;
      font-size: var(--fontSize);
      height: var(--height);
      white-space: nowrap;
      background: #F0F1F5;
      border-radius: 2px;
      padding: 0 .6em;
      cursor: pointer;

      & ~ .tag-item {
        display: none;
      }

      &.show {
        display: flex;
      }

      &:hover {
        opacity: .8;
      }
    }

    &.changing {
      .more-plus {
        display: flex;
        & ~ .tag-item {
          display: initial;
          visibility: hidden;
        }
      }
    }
  }
</style>
<style lang="scss">
  .flex-tag-tooltip-theme {
    .flex-tag-tips-content {
      padding: .3em;
      font-size: var(--fontSize);
      display: flex;
      flex-direction: column;
      gap: 6px;
      max-height: var(--popoverMaxHeight);
      @include scrollbar-y;

      &.is-link {
        .flex-tag-tips-item {
          color: #3A84FF;
          cursor: pointer;
        }
      }

      &.is-tag {
        padding: .5em;

        .flex-tag-tips-item {
          margin-right: auto;
          height: var(--height);
          line-height: var(--height);
          color: #63656E;
          background: #F0F1F5;
          border-radius: 2px;
          padding: 0 .6em;
        }
      }

      // clone模式改写部分复用的样式
      .tag-item {
        &::after {
          content: none !important;
        }
      }
    }
  }
</style>
